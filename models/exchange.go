package models

import (
	"github.com/VarthanV/pub-sub/exchange"
)

type Exchange struct {
	Base
	exchange.Exchange
	Name string `gorm:"uniqueIndex"`
}
