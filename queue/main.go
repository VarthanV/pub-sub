package queue

import (
	"sync"

	"github.com/VarthanV/pub-sub/messages"
	"gorm.io/gorm"
)

type Queue struct {
	mu       sync.Mutex
	Name     string
	messages []messages.Message
	// Queue which flushes messages added to this queue periodically
	// and flushes the buffer
	MessagesToPersist []*messages.Message
	Durable           bool
	db                *gorm.DB
}

func New(name string, durable bool) *Queue {

	return &Queue{
		messages:          make([]messages.Message, 0),
		Name:              name,
		Durable:           durable,
		MessagesToPersist: nil,
	}
}

func (qe *Queue) SetDB(db *gorm.DB) {
	qe.mu.Lock()
	defer qe.mu.Unlock()
	qe.db = db
}

func (qe *Queue) Enqueue(msg messages.Message) {
	qe.mu.Lock()
	defer qe.mu.Unlock()
	qe.messages = append(qe.messages, msg)
	if qe.Durable {
		qe.MessagesToPersist = append(qe.MessagesToPersist, &msg)
	}
}

func (qe *Queue) IsEmpty() bool {
	qe.mu.Lock()
	defer qe.mu.Unlock()
	return len(qe.messages) == 0
}

func (qe *Queue) Dequeue() *messages.Message {
	qe.mu.Lock()
	defer qe.mu.Unlock()

	if len(qe.messages) == 0 {
		return nil
	}

	elem := qe.messages[0]
	qe.messages = qe.messages[1:]
	return &elem
}
