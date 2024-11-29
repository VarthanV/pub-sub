package broker

import (
	"gorm.io/gorm"
)

type writer struct {
	bufferSize int // buffer size corresponds to item length in the array
	buffer     []interface{}
	db         *gorm.DB
}

func NewWriter(db *gorm.DB, size int) *writer {
	return &writer{
		bufferSize: size,
		db:         db,
	}
}

func (w *writer) Write(val interface{}) (n int, err error) {
	w.buffer = append(w.buffer, val)
	if len(w.buffer) >= w.bufferSize {
		tx := w.db.CreateInBatches(w.buffer, len(w.buffer))
		if tx.Error != nil {
			return 0, tx.Error
		}

		w.buffer = nil
	}
	return len(w.buffer), nil
}
