package rabbitmqrepo

import (
	"bytes"
	"context"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pkg/rabbitmq"
	"encoding/gob"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type MonitorMQ struct {
	ch *amqp.Channel
}

func NewMonitorMQ(rmq *rabbitmq.RabbitMQ) *MonitorMQ {
	return &MonitorMQ{
		ch: rmq.Channel,
	}
}

func (m *MonitorMQ) SendAllUserMetrics(ctx context.Context, usersAllMetrics []domain.UserMetricsMessage) error {
	return m.publish(ctx, "monitor.billing.all_users", usersAllMetrics)
}

func (m *MonitorMQ) publish(ctx context.Context, routingKey string, event interface{}) error {

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		zap.L().Error("gob.NewEncoder(&b).Encode(event)", zap.Error(err))
		return err
	}

	err := m.ch.Publish(
		"monitor-billing", // exchange
		routingKey,        // routing key
		false,
		false,
		amqp.Publishing{
			AppId:       "monitor-rest-server",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		zap.L().Error("m.ch.Publish: ", zap.Error(err))
		return err
	}

	return nil
}
