package messages

import "time"

type Message struct {
	ID   string      `json:"id"`
	TTL  time.Time   `json:"ttl"`
	Body interface{} `json:"body"`
}
