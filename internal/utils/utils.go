package utils

import (
	"github.com/CesarDelgadoM/extractor-reports/internal/producer"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
)

const (
	queuenames = "queues-names-queue"
	bindnames  = "queues-names-bind"
)

// Publish the name of the queues, for a listener consumer
func PublisQueueName(producer producer.IChannel, queuename string) {
	queue := producer.Queue(&stream.QueueOpts{
		Name: queuenames,
	})

	producer.BindQueue(&stream.BindOpts{
		Name: queue.Name,
		Key:  bindnames,
	})

	producer.Publish(&stream.PublishOpts{
		RoutingKey: bindnames,
		Body:       []byte(queuename),
	})
}
