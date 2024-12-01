package broker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/VarthanV/pub-sub/binding"
	"github.com/VarthanV/pub-sub/errors"
	"github.com/VarthanV/pub-sub/messages"
	"github.com/VarthanV/pub-sub/models"
	"github.com/VarthanV/pub-sub/pkg/config"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/VarthanV/pub-sub/exchange"
	"github.com/VarthanV/pub-sub/queue"
)

type Broker interface {
	Start(ctx context.Context) error
	CreateExchange(ctx context.Context, name string, exchangeType exchange.ExchangeType) error
	CreateQueue(ctx context.Context, name string, durable bool) error
	BindQueue(ctx context.Context, queueName, exchangeName, bindingKey string) error
	Subscribe(ctx context.Context, queueName string, conn *websocket.Conn) error
	Unsubscribe(ctx context.Context, queueName string, conn *websocket.Conn) error
	PublishMessage(ctx context.Context, exchangeName string, routingKey string, message interface{}) error
}

const (
	BrokerStateClosed int32 = iota
	BrokerStateRunning
	BrokerStateSettingUp
)

// broker orchestrates the  whole pub-sub process
type broker struct {
	mu sync.RWMutex

	db *gorm.DB

	cfg *config.Config
	// mapping of exchanges available in the current run-time
	exchanges map[string]*exchange.Exchange
	// mapping of queues available in the current run-time
	queues map[string]*queue.Queue

	subscriptions      map[string][]*websocket.Conn
	realtimeUpdatesSem chan struct{}
	currentState       atomic.Int32
}

func New(db *gorm.DB, cfg *config.Config) Broker {
	b := &broker{
		mu:            sync.RWMutex{},
		exchanges:     make(map[string]*exchange.Exchange),
		queues:        make(map[string]*queue.Queue),
		db:            db,
		subscriptions: make(map[string][]*websocket.Conn),
		realtimeUpdatesSem: make(chan struct{},
			cfg.WorkerConfig.MaxWorkersAllowedConcurrentlyForRealtimeUpdates),
		cfg: cfg,
	}
	return b
}

func (b *broker) Start(ctx context.Context) error {
	var (
		errorChan     = make(chan error, 1)
		isStaredChan  = make(chan bool, 1)
		startDeadline <-chan time.Time
	)

	if b.currentState.Load() == int32(BrokerStateRunning) {
		// currently only allowing only one instance of broker
		// to be running
		logrus.Fatal("broker already running")
	}

	go func() {
		defer close(isStaredChan)
		defer close(errorChan)

		// Rebuild broker

		// Setup tokens in sem for configuring workers
		for i := 0; i <
			b.cfg.WorkerConfig.MaxWorkersAllowedConcurrentlyForRealtimeUpdates; i++ {
			b.realtimeUpdatesSem <- struct{}{}
		}

		isStaredChan <- true
	}()

	// See if deadline already there in the context else we derive a custom deadline for starting
	// and make it done

	_, ok := ctx.Deadline()
	if !ok {
		// FIXME: take the deadline from config
		logrus.Infof("deadline not enforced in base ctx ,setting deadline for %d  minutes", 1)
		startDeadline = time.After(time.Minute * 1)
	}

	for {
		select {
		case e, ok := <-errorChan:
			if ok {
				logrus.Error("error in startup ", e)
				return e
			}

		case v, ok := <-isStaredChan:
			if ok && v {
				b.currentState.Store(BrokerStateRunning)
				logrus.Info("Started broker.....")
				return nil
			}

		case <-ctx.Done():
			logrus.Error("context done")
			return fmt.Errorf("context done")

		case <-startDeadline:
			logrus.Error("start deadline exceeded")
			return fmt.Errorf("start deadline exceeded")

		}
	}
}

// CreateExchange creates an exchange with the given type
func (b *broker) CreateExchange(ctx context.Context, name string, exchangeType exchange.ExchangeType) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.exchanges[name]
	if ok {
		return errors.Handle(errors.ErrExchangeAlreadyExists)
	}

	b.exchanges[name] = exchange.New(name, exchangeType)
	logrus.Info("Created exchange ", name)
	err := b.db.WithContext(ctx).Create(&models.Exchange{
		Name:         name,
		ExchangeType: exchangeType,
	}).Error
	if err != nil {
		logrus.Error("error in creating exchange ", err)
		return err
	}
	return nil
}

// CreateQueue creates a queue, if needs to survive the rebuidling during restart
// can be configured using the `durable` flag
func (b *broker) CreateQueue(ctx context.Context, name string, durable bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.queues[name]
	if ok {
		return errors.Handle(errors.ErrExchangeAlreadyExists)
	}
	b.queues[name] = queue.New(name, durable)
	logrus.Info("Created queue ", name)
	err := b.db.WithContext(ctx).Create(&models.Queue{
		Name:    name,
		Durable: durable,
	}).Error
	if err != nil {
		logrus.Error("error in creating queue ", err)
		return err
	}
	return nil
}

// BindQueue binds a queue to an exchange  by the given binding key
func (b *broker) BindQueue(ctx context.Context, queueName, exchangeName, bindingKey string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	exchange, exists := b.exchanges[exchangeName]
	if !exists {
		errors.Handle(errors.ErrExchangeDoesnotExist)
	}

	q, exists := b.queues[queueName]
	if !exists {
		return errors.Handle(errors.ErrQueueDoesnotExist)
	}
	if exchange.Bindings == nil {
		exchange.Bindings = make(map[string]*binding.Binding)
	}

	bi, exists := exchange.Bindings[bindingKey]
	if !exists {
		bi = &binding.Binding{
			Key:    bindingKey,
			Queues: []*queue.Queue{},
		}
		exchange.Bindings[bindingKey] = bi
	}

	for _, bq := range bi.Queues {
		if q == bq {
			return nil
		}
	}
	bi.Queues = append(bi.Queues, q)

	// Get queue with the name
	modelQ := &models.Queue{}

	err := b.db.Where(&models.Queue{Name: queueName}).
		Last(&modelQ).Error
	if err != nil {
		logrus.Error("error in getting q ", err)
		return err
	}

	err = b.db.Create(&models.Binding{
		Key:          bindingKey,
		ExchangeName: exchangeName,
		Queues:       []models.Queue{{Base: models.Base{ID: modelQ.ID}}},
	}).Error
	return err
}

// Subscribe subscribes a websocket connection to a topic for exchange of realtime updates.
func (b *broker) Subscribe(ctx context.Context, queueName string, conn *websocket.Conn) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	logrus.Info("subscribed to queue ", queueName)
	b.subscriptions[queueName] = append(b.subscriptions[queueName], conn)
	return nil
}

// PublishMessage publishes message.
func (b *broker) PublishMessage(ctx context.Context, exchangeName string, routingKey string, message interface{}) error {
	b.mu.RLock()
	defer b.mu.RLock()
	e, exists := b.exchanges[exchangeName]
	if !exists {
		return errors.Handle(errors.ErrExchangeDoesnotExist)
	}

	logrus.Info("Found exchange...")

	switch e.ExchangeType {
	case exchange.ExchangeTypeFanOut:
		logrus.Info("Fan out data to connections")
		// Send to all queues and corresponding subscribers
		for _, binding := range e.Bindings {
			for _, q := range binding.Queues {
				msg := messages.Message{
					ID:   uuid.NewString(),
					Body: message,
				}
				// Enqueue the message to the queue
				q.Enqueue(msg)

				logrus.Info("found connections ", len(b.subscriptions[q.Name]))
				for _, conn := range b.subscriptions[q.Name] {
					go func() {
						<-b.realtimeUpdatesSem
						conn.SetWriteDeadline(time.Now().Add(time.Minute))
						err := conn.WriteJSON(msg)
						if err != nil {
							logrus.Error("error in writing to conn ", err)
						}
						b.realtimeUpdatesSem <- struct{}{}
					}()
				}

			}

		}

	case exchange.ExchangeTypeDirect:
		logrus.Info("Directing to excatly bound routing key")
		binding, exists := e.Bindings[routingKey]
		if !exists {
			// Drop messages
			logrus.Warnf("No binding found for routing key: %s in exchange: %s", routingKey, exchangeName)
			return nil
		}

		for _, q := range binding.Queues {
			// Enqueue the message to the queue
			msg := messages.Message{
				ID:   uuid.NewString(),
				Body: message,
			}
			q.Enqueue(msg)

			logrus.Info("found connections ", len(b.subscriptions[q.Name]))

			for _, conn := range b.subscriptions[q.Name] {
				go func() {
					<-b.realtimeUpdatesSem
					conn.SetWriteDeadline(time.Now().Add(time.Minute))
					err := conn.WriteJSON(msg)
					if err != nil {
						logrus.Error("error in writing to conn ", err)
					}
					b.realtimeUpdatesSem <- struct{}{}
				}()
			}

		}
	}

	return nil

}

// Unsubscribe unsubscribes a subscriber from the relatime updates of the queue.
func (b *broker) Unsubscribe(ctx context.Context, queueName string, conn *websocket.Conn) error {
	connections, ok := b.subscriptions[queueName]
	if !ok {
		if err := errors.Handle(errors.ErrSubscriptionDoesnotExist); err != nil {
			return err
		}
	}

	indexToDelete := -1
	for i, c := range connections {
		if c.NetConn() == conn.NetConn() {
			indexToDelete = i
			err := c.Close()
			if err != nil {
				logrus.Error("error in closing connection ", err)
			}
		}
	}

	if indexToDelete >= 0 {
		b.subscriptions[queueName] = deleteElement(connections, indexToDelete)
	}

	return nil
}
