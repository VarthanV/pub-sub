package models

type Queue struct {
	Base
	Name    string `gorm:"uniqueIndex"`
	Durable bool
}
