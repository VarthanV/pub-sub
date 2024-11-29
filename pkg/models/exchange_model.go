package models

import (
	"github.com/VarthanV/pub-sub/pkg/exchange"
)

type ExchangeModel struct {
	Base
	exchange.Exchange
	Name string `gorm:"uniqueIndex"`
}
