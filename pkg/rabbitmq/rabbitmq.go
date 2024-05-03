package rabbitmq

import (
	"dogker/lintang/monitor-service/config"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(cfg *config.Config) *RabbitMQ {
	zap.L().Info("rmq address: " + cfg.RabbitMQ.RMQAddress)
	conn, err := amqp.Dial(cfg.RabbitMQ.RMQAddress)
	if err != nil {
		zap.L().Fatal("error: cannot connect to rabbitmq: " + err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("error can't get rabbitmq cahnnel: " + err.Error())
	}

	err = channel.ExchangeDeclare(
		"monitor-billing",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("err: channel.ExchangeDeclare : " + err.Error())
	}

	// bikin queue sekalian
	queue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		zap.L().Fatal("channel.QueueDeclare")
	}
	err = channel.QueueBind(
		queue.Name,
		"monitor.billing.all_users",
		"monitor-billing",
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("channel.QueueBind")
	}

	err = channel.Qos(
		1, 0,
		false,
	)
	if err != nil {
		zap.L().Fatal("err: channel.Qos" + err.Error())
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}

}

func (r *RabbitMQ) Close() error {
	return r.Connection.Close()
}
