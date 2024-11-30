package broker

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (b *broker) processStream(ctx context.Context, ch <-chan interface{}, bufferSize int, duration time.Duration) {
	writer := NewWriter(ctx, b.db, bufferSize, duration)
	for val := range ch {
		_, err := writer.Write(val)
		if err != nil {
			logrus.Error("error in persisting ", err)
		}
	}
}
