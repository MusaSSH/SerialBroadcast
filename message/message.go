package message

import (
	"context"
	"time"

	"github.com/MusaSSH/SerialBroadcast/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

type Message struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func (m Message) Publish(msg []byte) error {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*10)
	defer cf()
	err := m.channel.PublishWithContext(ctx, "", m.queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	return err
}

func Build() fx.Option {
	return fx.Provide(func(c config.Config) (Message, error) {
		conn, err := amqp.Dial(c.AmqpAddr)
		if err != nil {
			return Message{}, err
		}

		cn, err := conn.Channel()
		if err != nil {
			return Message{}, err
		}

		q, err := cn.QueueDeclare(c.AmqpQueueName, false, false, false, false, amqp.Table{
			"x-message-ttl": c.AmqpTtl,
		})
		if err != nil {
			return Message{}, err
		}

		return Message{
			connection: conn,
			channel:    cn,
			queue:      q,
		}, nil
	})
}
