package models

import (
	"github.com/VarthanV/pub-sub/pkg/exchange"
	"gorm.io/gorm"
)

type ExchangeModel struct {
	gorm.Model
	exchange.Exchange
	Name string `gorm:"uniqueIndex"`
}
