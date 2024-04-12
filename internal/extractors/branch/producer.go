package branch

import (
	"github.com/CesarDelgadoM/extractor-reports/config"
	"github.com/CesarDelgadoM/extractor-reports/internal/producer"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
)

// Initializator to create a branch producer
func NewBranchProducer(config *config.Branch, rabbit *stream.RabbitMQ) producer.IProducer {

	opts := &producer.ProducerOpts{
		ExchangeName: config.ExchangeName,
		ExchangeType: config.ExchangeType,
		ContentType:  config.ContentType,
	}

	p := producer.NewProducer(opts, rabbit)

	p.Exchange(&stream.ExchangeOpts{
		Name:    opts.ExchangeName,
		Kind:    opts.ExchangeType,
		Durable: true,
	})

	return p
}
