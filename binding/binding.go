package binding

import "github.com/VarthanV/pub-sub/queue"

type Binding struct {
	Key    string
	Queues []*queue.Queue
}
