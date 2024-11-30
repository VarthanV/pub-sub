package models

type Binding struct {
	Base
	Key          string  `json:"key,omitempty"`
	Queues       []Queue `gorm:"many2many:binding_queue"`
	ExchangeName string  `json:"exchange_name,omitempty"`

	Exchange Exchange `gorm:"foreignKey:ExchangeName;references:Name"`
}
