package binding

import "github.com/VarthanV/pub-sub/pkg/queue"

type Binding struct {
	Key   string
	Queue []queue.Queue
}
