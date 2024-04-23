package databus

import (
	"github.com/CesarDelgadoM/extractor-reports/config"
	"github.com/CesarDelgadoM/extractor-reports/internal/producer"
	"github.com/CesarDelgadoM/extractor-reports/internal/utils"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
)

const (
	queuenames     = "queues-names-queue"
	bindqueuenames = "queues-names-bind"
)

type IDataBus interface {
	PublishQueueName(msg producer.MessageQueueNames)
}

type dataBus struct {
	producer producer.IProducer
}

func NewDataBus(config *config.DataBus, rabbitmq *stream.RabbitMQ) IDataBus {

	opts := &producer.ProducerOpts{
		ExchangeType: config.ExchangeType,
		ExchangeName: config.ExchangeName,
		ContentType:  config.ContentType,
	}

	p := producer.NewProducer(opts, rabbitmq)

	p.Exchange(&stream.ExchangeOpts{
		Name:    opts.ExchangeName,
		Kind:    opts.ExchangeType,
		Durable: true,
	})

	return &dataBus{
		producer: p,
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
		Key:  bindqueuenames,
	})

	db.producer.Publish(&stream.PublishOpts{
		RoutingKey: bindqueuenames,
		Body:       utils.ToBytes(msg),
	})
}
