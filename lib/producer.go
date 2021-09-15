package lib

import (
	"log"

	"github.com/streadway/amqp"
)

const (
	NOTIFICATION_EXCHANGE = "NOTIFICATION_EXCHANGE"
	COMMENT_QUEUE         = "COMMENT_QUEUE"
	//STATUS_CHANGE_QUEUE = "COMMENT_QUEUE" // TODO
)

type Producer interface {
	Setup() error
	Publish(msg []byte) error // TODO how to handle errors
}

type AmqpProducer struct {
	connection *amqp.Connection
}

func NewAmqpProducer(connection *amqp.Connection) (*AmqpProducer, error) {
	producer := AmqpProducer{connection: connection}

	err := producer.Setup()
	if err != nil {
		log.Printf("Error setting up %v\n", err)
		return nil, err
	}

	return &producer, nil
}

func (p *AmqpProducer) Setup() error {
	channel, err := p.connection.Channel()
	if err != nil {
		log.Printf("Error creating channel %v\n", err)
		return err
	}
	defer channel.Close()

	/*
		// TODO not working, look later
			err = p.declareExchange(channel)
			if err != nil {
				log.Printf("Exchange declarion went wrong %v\n", err)
				return err
			}
	*/

	err = p.declareQueue(channel)
	if err != nil {
		return err
	}

	return nil
}

func (p *AmqpProducer) declareExchange(channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(NOTIFICATION_EXCHANGE, amqp.ExchangeDirect,
		true, false, false, false, nil)
	if err != nil {
		log.Printf("Error declaring exchange %v\n", err)
		return err
	}

	return nil
}

func (p *AmqpProducer) declareQueue(channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(COMMENT_QUEUE, true, false, false, false, nil)
	if err != nil {
		log.Printf("Error declaring queue %v\n", err)
		return err
	}

	return nil
}

func (p *AmqpProducer) Publish(msg []byte) error {
	channel, err := p.connection.Channel()
	if err != nil {
		log.Printf("Error getting channel %v\n", err)
		return err
	}
	defer channel.Close()

	err = channel.Publish("", COMMENT_QUEUE, false,
		false, amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	if err != nil {
		log.Printf("Error publishing message %v\n", err)
		return err
	}

	return nil
}
