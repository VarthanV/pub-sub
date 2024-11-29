package messages

import "time"

type Message struct {
	TTL  time.Duration
	Body []byte
}
