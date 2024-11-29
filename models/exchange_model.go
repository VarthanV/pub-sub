package models

import (
	"github.com/VarthanV/pub-sub/exchange"
)

type ExchangeModel struct {
	Base
	exchange.Exchange
	Name string `gorm:"uniqueIndex"`
}
