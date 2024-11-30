package broker

import "github.com/sirupsen/logrus"

func (b *broker) processStream(ch <-chan interface{}, bufferSize int) {
	writer := NewWriter(b.db, bufferSize)
	for val := range ch {
		_, err := writer.Write(val)
		if err != nil {
			logrus.Error("error in persisting ", err)
		}
	}
}
