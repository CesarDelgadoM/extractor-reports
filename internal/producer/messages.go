package producer

// Message queue names
type MessageQueueNames struct {
	TypeReport string
	QueueName  string
}

// Message to produce
type Message struct {
	Userid uint
	Format string
	Status int
	Data   any
}
