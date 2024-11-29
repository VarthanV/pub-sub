package exchange

import "github.com/VarthanV/pub-sub/pkg/binding"

type ExchangeType string

const (
	ExchangeTypeDirect ExchangeType = "direct"
	ExchangeTypeFanOut ExchangeType = "fanout"
)

type Exchange struct {
	Name         string
	ExchangeType ExchangeType
	Bindings     []binding.Binding
}

func New(name string, exchangeType ExchangeType) *Exchange {
	return &Exchange{
		Name:         name,
		ExchangeType: exchangeType,
		Bindings:     make([]binding.Binding, 0),
	}
}
