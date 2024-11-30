package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Queue struct {
	Base
	Name    string `gorm:"uniqueIndex"`
	Durable bool
}

func (q *Queue) BeforeCreate(tx *gorm.DB) error {
	tx.Clauses(clause.OnConflict{
		DoNothing: true,
		DoUpdates: clause.Assignments(map[string]interface{}{"updated_at": q.UpdatedAt}),
	}, clause.Returning{Columns: []clause.Column{{Name: "id"}}})
	return nil
}
