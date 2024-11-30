package models

import "gorm.io/datatypes"

type Message struct {
	Base
	UniqueID  string         `json:"id"`
	Body      datatypes.JSON `json:"body"`
	QueueName string         `json:"name,omitempty"`

	Queue Queue `gorm:"foreignKey:QueueName;references:Name"`
}
