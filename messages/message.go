package messages

type Message struct {
	ID        string      `json:"unique_id"`
	Body      interface{} `json:"body"`
	QueueName string
}
