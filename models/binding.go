package models

type Binding struct {
	Base
	Key    string
	Queues []Queue `gorm:"many2many:binding_queue"`
}
