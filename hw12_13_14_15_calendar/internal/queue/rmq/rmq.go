package rmqqueue

import (
	"context"
	"fmt"

	"github.com/LightAir/otus_home_work/hw12_13_14_15_calendar/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rmq struct {
	config config.Config
	conn   *amqp.Connection
	ctx    context.Context
}

func NewRmq(config *config.Config) *Rmq {
	return &Rmq{
		config: *config,
	}
}

func (rmq *Rmq) Connect(ctx context.Context) error {
	rmq.ctx = ctx

	c := rmq.config.Rmq
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Pswd, c.Host, c.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	rmq.conn = conn

	return nil
}

func getQueue(name string, ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
}

func (rmq *Rmq) Sent(body []byte, name string) error {
	ch, err := rmq.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed open channel: %w", err)
	}

	q, err := getQueue(name, ch)
	if err != nil {
		return fmt.Errorf("failed queue declare: %w", err)
	}

	err = ch.PublishWithContext(
		rmq.ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed publish: %w", err)
	}
	defer ch.Close()

	return nil
}

func (rmq *Rmq) Receive(name string, callback func(body []byte)) error {
	ch, err := rmq.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed open channel: %w", err)
	}

	q, err := getQueue(name, ch)
	if err != nil {
		return fmt.Errorf("failed queue declare: %w", err)
	}

	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed open channel: %w", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msg {
			callback(d.Body)
		}
	}()

	<-forever

	return nil
}

func (rmq *Rmq) Close() error {
	if !rmq.conn.IsClosed() {
		return rmq.conn.Close()
	}

	return nil
}
