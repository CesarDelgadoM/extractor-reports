package databus

import (
	"github.com/CesarDelgadoM/extractor-reports/internal/producer"
	"github.com/CesarDelgadoM/extractor-reports/internal/utils"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
)

const (
	queuenames = "queues-names-queue"
	bindnames  = "queues-names-bind"
)

type IDataBus interface {
	PublishQueueName(msg producer.MessageQueueNames)
}

type dataBus struct {
	producer producer.IProducer
}

func NewDataBus(rabbitmq *stream.RabbitMQ) IDataBus {

	opts := &producer.ProducerOpts{
		ExchangeType: "direct",
		ContentType:  "application/json",
	}

	return &dataBus{
		producer: producer.NewProducer(opts, rabbitmq),
	}
}

// Publish the name of the queues, for a listener consumer
func (db *dataBus) PublishQueueName(msg producer.MessageQueueNames) {
	queue := db.producer.Queue(&stream.QueueOpts{
		Name:    queuenames,
		Durable: true,
	})

	db.producer.BindQueue(&stream.BindOpts{
		Name: queue.Name,
		Key:  bindnames,
	})

	db.producer.Publish(&stream.PublishOpts{
		RoutingKey: bindnames,
		Body:       utils.ToBytes(msg),
	})
}
