package rabbitmqrepo

import "github.com/streadway/amqp"



type Monitor struct{ 
	ch *amqp.Channel
}

func NewMonitor(channel *amqp.Channel) (*Monitor, error) {
	return &Monitor{
		ch: channel,

	}, nil
}







