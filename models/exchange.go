package models

import "github.com/VarthanV/pub-sub/exchange"

type Exchange struct {
	Base
	Name         string                `gorm:"uniqueIndex" json:"name,omitempty"`
	ExchangeType exchange.ExchangeType `json:"exchange_type,omitempty"`
	Bindings     []Binding             `gorm:"many2many:exchange_bindings"`
}
