package producer

import "encoding/json"

// Message to produce
type Message struct {
	Userid uint
	Type   string
	Format string
	Status int
	Data   any
}

func (m *Message) ToBytes() []byte {
	bytes, _ := json.Marshal(&m)
	return bytes
}
