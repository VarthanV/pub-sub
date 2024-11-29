package queue_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/VarthanV/pub-sub/messages"
	"github.com/VarthanV/pub-sub/queue"
)

func TestQueue(t *testing.T) {
	// Create a new queue
	q := queue.New("foo", false)

	// Test if the queue is empty initially
	if !q.IsEmpty() {
		t.Errorf("Expected queue to be empty, but it is not")
	}

	// Create messages
	msg1 := messages.Message{
		ID:   "1",
		TTL:  time.Now().Add(5 * time.Second),
		Body: []byte("Test message 1"),
	}

	msg2 := messages.Message{
		ID:   "2",
		TTL:  time.Now().Add(10 * time.Second),
		Body: []byte("Test message 2"),
	}

	// Enqueue the first message
	q.Enqueue(msg1)

	// Test if the queue is no longer empty
	if q.IsEmpty() {
		t.Errorf("Expected queue to not be empty, but it is")
	}

	// Enqueue the second message
	q.Enqueue(msg2)

	// Dequeue the first message
	dequeuedMsg := q.Dequeue()
	if dequeuedMsg == nil {
		t.Fatalf("Dequeued message is nil, expected %+v", msg1)
	}
	if dequeuedMsg.ID != msg1.ID || dequeuedMsg.TTL != msg1.TTL || !bytes.Equal(dequeuedMsg.Body, msg1.Body) {
		t.Errorf("Dequeued message mismatch. Got %+v, expected %+v", dequeuedMsg, msg1)
	}

	// Dequeue the second message
	dequeuedMsg = q.Dequeue()
	if dequeuedMsg == nil {
		t.Fatalf("Dequeued message is nil, expected %+v", msg2)
	}
	if dequeuedMsg.ID != msg2.ID || dequeuedMsg.TTL != msg2.TTL || !bytes.Equal(dequeuedMsg.Body, msg2.Body) {
		t.Errorf("Dequeued message mismatch. Got %+v, expected %+v", dequeuedMsg, msg2)
	}

	// Test if the queue is empty again
	if !q.IsEmpty() {
		t.Errorf("Expected queue to be empty, but it is not")
	}

	// Try to dequeue from an empty queue
	dequeuedMsg = q.Dequeue()
	if dequeuedMsg != nil {
		t.Errorf("Expected nil when dequeuing from an empty queue, but got %+v", dequeuedMsg)
	}
}
