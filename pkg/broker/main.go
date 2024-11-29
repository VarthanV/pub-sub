package broker

import (
	"sync"

	"github.com/VarthanV/pub-sub/pkg/binding"
	"github.com/VarthanV/pub-sub/pkg/errors"

	"github.com/VarthanV/pub-sub/pkg/exchange"
	"github.com/VarthanV/pub-sub/pkg/queue"
)

// Broker orchestrates the  whole pub-sub process
type Broker struct {
	mu        sync.Mutex
	exchanges map[string]*exchange.Exchange
	queues    map[string]*queue.Queue
}

func New() *Broker {
	return &Broker{
		mu:        sync.Mutex{},
		exchanges: make(map[string]*exchange.Exchange),
		queues:    make(map[string]*queue.Queue),
	}
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
