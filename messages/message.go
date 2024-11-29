package messages

import "time"

type Message struct {
	ID   string
	TTL  time.Time
	Body []byte
}
