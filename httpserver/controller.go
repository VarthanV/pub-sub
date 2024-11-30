package httpserver

import (
	"github.com/VarthanV/pub-sub/broker"
	"gorm.io/gorm"
)

type HTTPController struct {
	db     *gorm.DB
	broker broker.Broker
}

func NewController(db *gorm.DB, b broker.Broker) *HTTPController {
	return &HTTPController{
		db:     db,
		broker: b,
	}
}
