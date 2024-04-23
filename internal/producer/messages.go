package producer

// Message queue names
type MessageQueueNames struct {
	ReportType string
	QueueName  string
}

// Message to produce
type Message struct {
	Userid uint
	Format string
	Status int
	Data   []byte
}
