package producer

import (
	"context"
	"time"

	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
	"github.com/CesarDelgadoM/extractor-reports/pkg/stream"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IChannel interface {
	Exchange(opts *stream.ExchangeOpts)
	BindQueue(opts *stream.BindOpts)
	Publish(opts *stream.PublishOpts)
	Queue(opts *stream.QueueOpts) amqp.Queue
	Close()
}

type ProducerOpts struct {
	ExchangeName string
	ExchangeType string
	ContentType  string
}

type producer struct {
	opts *ProducerOpts
	ch   *amqp.Channel
}

func NewProducer(opts *ProducerOpts, rabbit *stream.RabbitMQ) IChannel {
	p := producer{
		opts: opts,
		ch:   rabbit.Channel(),
	}

	p.Exchange(&stream.ExchangeOpts{
		Name:    p.opts.ExchangeName,
		Kind:    p.opts.ExchangeType,
		Durable: true,
	})

	return &p
}

func (p *producer) Exchange(opts *stream.ExchangeOpts) {
	err := p.ch.ExchangeDeclare(
		opts.Name,
		opts.Kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		zap.Log.Error("Failed to create exchange: ", err)
	}
}

func (p *producer) BindQueue(opts *stream.BindOpts) {
	err := p.ch.QueueBind(
		opts.Name,
		opts.Key,
		p.opts.ExchangeName,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		zap.Log.Info("Failed to create queue bind: ", err)
	}
}

func (p *producer) Publish(opts *stream.PublishOpts) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	publishing := amqp.Publishing{
		ContentType: p.opts.ContentType,
		Body:        opts.Body,
	}

	err := p.ch.PublishWithContext(
		ctx,
		p.opts.ExchangeName,
		opts.RoutingKey,
		opts.Mandatory,
		opts.Immediate,
		publishing,
	)
	if err != nil {
		zap.Log.Error("Failed to pusblih message: ", err)
	}
}

func (p *producer) Queue(opts *stream.QueueOpts) amqp.Queue {
	queue, err := p.ch.QueueDeclare(
		opts.Name,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		zap.Log.Error("Failed to create queue: ", queue)
	}

	return queue
}

func (p *producer) Close() {
	if err := p.ch.Close(); err != nil {
		zap.Log.Error("Filed to close channel: ", err)
	}
}
