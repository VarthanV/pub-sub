package models

import "gorm.io/gorm"

type Base struct {
	gorm.Model
	ID uint `json:"int" gorm:"primaryKey"`
}
