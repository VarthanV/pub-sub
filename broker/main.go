package broker

import (
	"context"
	"sync"
	"time"

	"github.com/VarthanV/pub-sub/binding"
	"github.com/VarthanV/pub-sub/errors"
	"github.com/VarthanV/pub-sub/models"
	"github.com/VarthanV/pub-sub/pkg/config"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/VarthanV/pub-sub/exchange"
	"github.com/VarthanV/pub-sub/queue"
)

// broker orchestrates the  whole pub-sub process

type Broker interface {
	Start(ctx context.Context)
	CreateExchange(ctx context.Context, name string, exchangeType exchange.ExchangeType) error
	CreateQueue(ctx context.Context, name string, durable bool) error
	BindQueue(ctx context.Context, queueName, exchangeName, bindingKey string) error
	Subscribe(ctx context.Context, queueName string, conn *websocket.Conn) error
}

type broker struct {
	mu sync.Mutex

	db *gorm.DB

	cfg *config.Config
	// mapping of exchanges available in the current run-time
	exchanges map[string]*exchange.Exchange
	// mapping of queues available in the current run-time
	queues map[string]*queue.Queue
	// we will be persisting the exchanges info ,queues info to a persistent storage to rebuild on crash
	// fixed amount of workers will be bound to this channel , buffered writers will be used and the data
	// will be flushed when the buffer gets full or when the ticker ticks whichever happens first
	persistQueue chan interface{}

	subscriptions map[string][]*websocket.Conn
}

func New(db *gorm.DB, cfg *config.Config) Broker {
	b := &broker{
		mu:            sync.Mutex{},
		exchanges:     make(map[string]*exchange.Exchange),
		queues:        make(map[string]*queue.Queue),
		persistQueue:  make(chan interface{}, 100),
		db:            db,
		subscriptions: make(map[string][]*websocket.Conn),
		cfg:           cfg,
	}
	return b
}

func (b *broker) Start(ctx context.Context) {

	logrus.Info("Started broker.....")
	go func() {
		var (
			persistentWg sync.WaitGroup
		)
		// only one start instance will be running so safely say we will only be closing the
		// chan

		workersSize := 10 // aribtary chosen might tweak later
		for i := 0; i < workersSize; i++ {
			persistentWg.Add(1)
			go func() {
				defer persistentWg.Done()
				b.processStream(ctx, b.persistQueue, workersSize,
					time.Second*time.Duration(b.cfg.SyncConfiguration.CheckpointInSeconds))
			}()
		}

		persistentWg.Wait()
	}()

	<-ctx.Done()

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
	b.persistQueue <- &models.Exchange{
		Name:         name,
		ExchangeType: exchangeType,
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
	b.persistQueue <- &models.Queue{
		Name:    name,
		Durable: durable,
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

	b.persistQueue <- &models.Binding{
		ExchangeName: exchangeName,
		Key:          bindingKey,
		Queues:       []models.Queue{{Name: exchange.Name}},
	}
	return nil
}

// Subscribe subscribes a websocket connection to a topic for exchange of realtime updates.
func (b *broker) Subscribe(ctx context.Context, queueName string, conn *websocket.Conn) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscriptions[queueName] = append(b.subscriptions[queueName], conn)
	return nil
}
