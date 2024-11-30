package broker

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type writer struct {
	// buffer size corresponds to item length in the array
	bufferSize int
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
		// flush
		for _, v := range w.buffer {
			err = w.db.Create(v).Error
			if err != nil {
				logrus.Error("error in creating val ", err)
			}
		}

		// If errror is not nil we dont want to clear the buffer we would like
		// to persist it at next cycle
		if err != nil {
			w.buffer = nil
		}
	}
	return len(w.buffer), nil
}
