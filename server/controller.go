package server

import (
	"github.com/VarthanV/pub-sub/broker"
	"gorm.io/gorm"
)

type Controller struct {
	db     *gorm.DB
	broker broker.Broker
}

func NewController(db *gorm.DB, b broker.Broker) *Controller {
	return &Controller{
		db:     db,
		broker: b,
	}
}
