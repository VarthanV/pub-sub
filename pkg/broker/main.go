package broker

import (
	"sync"

	"github.com/VarthanV/pub-sub/pkg/binding"
	"github.com/VarthanV/pub-sub/pkg/errors"

	"github.com/VarthanV/pub-sub/pkg/exchange"
	"github.com/VarthanV/pub-sub/pkg/queue"
)

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

	exchange.Bindings = append(exchange.Bindings, binding.Binding{
		Key:    bindingKey,
		Queues: []*queue.Queue{q},
	})
	return nil

}
