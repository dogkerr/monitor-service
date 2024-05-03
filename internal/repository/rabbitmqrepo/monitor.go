package rabbitmqrepo

import (
	"bytes"
	"context"
	"dogker/lintang/monitor-service/domain"
	"encoding/gob"
	"time"

	"github.com/streadway/amqp"
)

type MonitorMQ struct {
	ch *amqp.Channel
}

func NewMonitorMQ(channel *amqp.Channel) *MonitorMQ {
	return &MonitorMQ{
		ch: channel,
	}
}

func (m *MonitorMQ) SendAllUserMetrics(ctx context.Context, usersAllMetrics domain.AllUsersMetricsMessage) error {
	return m.publish(ctx, "monitor.billing", usersAllMetrics)
}

func (m *MonitorMQ) publish(ctx context.Context, routingKey string, event interface{}) error {

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
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
		return err
	}

	return nil
}
