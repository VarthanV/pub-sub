package messages

import "time"

type Message struct {
	ID   string      `json:"unique_id"`
	Body interface{} `json:"body"`
	TTL  time.Duration
}
