package broker

import (
	"context"
	"sync"

	"github.com/VarthanV/pub-sub/binding"
	"github.com/VarthanV/pub-sub/errors"
	"gorm.io/gorm"

	"github.com/VarthanV/pub-sub/exchange"
	"github.com/VarthanV/pub-sub/queue"
)

// Broker orchestrates the  whole pub-sub process
type Broker struct {
	db *gorm.DB
	mu sync.Mutex
	// mapping of exchanges available in the current run-time
	exchanges map[string]*exchange.Exchange
	// mapping of queues available in the current run-time
	queues map[string]*queue.Queue
	// we will be persisting the exchanges info ,queues info to a persistent storage to rebuild on crash
	// fixed amount of workers will be bound to this channel , buffered writers will be used and the data
	// will be flushed when the buffer gets full or when the ticker ticks whichever happens first
	persistQueue chan interface{}
}

func New() *Broker {
	b := &Broker{
		mu:           sync.Mutex{},
		exchanges:    make(map[string]*exchange.Exchange),
		queues:       make(map[string]*queue.Queue),
		persistQueue: make(chan interface{}, 100),
	}
	return b
}

func (b *Broker) Start(ctx context.Context) {

	go func() {
		var (
			persistentWg sync.WaitGroup
		)
		// only one start instance will be running so safely say we will only be closing the
		// chan
		defer close(b.persistQueue)

		workersSize := 10 // aribtary chosen might tweak later
		for i := 0; i < workersSize; i++ {
			persistentWg.Add(1)
			go func() {
				defer persistentWg.Done()
				b.processStream(b.persistQueue, workersSize)
			}()
		}

		persistentWg.Wait()
	}()

}

// CreateExchange creates an exchange with the given type
func (b *Broker) CreateExchange(name string, exchangeType exchange.ExchangeType) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.exchanges[name]
	if ok {
		errors.Handle(errors.ErrExchangeAlreadyExists)
	}

	b.exchanges[name] = exchange.New(name, exchangeType)
	return nil
}

// CreateQueue creates a queue, if needs to survive the rebuidling during restart
// can be configured using the `durable` flag
func (b *Broker) CreateQueue(name string, durable bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.queues[name]
	if ok {
		errors.Handle(errors.ErrExchangeAlreadyExists)
	}
	b.queues[name] = queue.New(name, durable)
	return nil
}

// BindQueue binds a queue to an exchange  by the given binding key
func (b *Broker) BindQueue(queueName, exchangeName, bindingKey string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	exchange, exists := b.exchanges[exchangeName]
	if !exists {
		errors.Handle(errors.ErrExchangeDoesnotExist)
	}

	q, exists := b.queues[queueName]
	if !exists {
		errors.Handle(errors.ErrQueueDoesnotExist)
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
	return nil
}
