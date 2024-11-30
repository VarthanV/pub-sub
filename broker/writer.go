package broker

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type writer struct {
	// buffer size corresponds to item length in the array
	bufferSize         int
	buffer             []interface{}
	db                 *gorm.DB
	checkPointDuration time.Duration
	ctx                context.Context
}

func NewWriter(ctx context.Context, db *gorm.DB, size int,
	duration time.Duration) *writer {
	return &writer{
		bufferSize:         size,
		db:                 db,
		checkPointDuration: duration,
		ctx:                ctx,
	}
}

func (w *writer) syncToDB() error {
	var (
		err error
	)
	logrus.Info("######## Checkpoint to db #########")
	for _, v := range w.buffer {
		err = w.db.Model(v).Create(v).
			Clauses(clause.OnConflict{
				DoNothing: true,
			}).Error
		if err != nil {
			logrus.Error("error in creating val ", err)
		}
	}

	// If errror is not nil we dont want to clear the buffer we would like
	// to persist it at next cycle
	if err == nil {
		w.buffer = nil
	}

	return nil
}

func (w *writer) Write(val interface{}) (n int, err error) {
	go func() {
		tick := time.NewTicker(w.checkPointDuration)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				w.syncToDB()
			case <-w.ctx.Done():
				logrus.Info("exiting due to done ctx")
				w.syncToDB()
				return
			}
		}
	}()

	w.buffer = append(w.buffer, val)
	if len(w.buffer) >= w.bufferSize {
		// flush
		err = w.syncToDB()

	}
	return len(w.buffer), nil
}
