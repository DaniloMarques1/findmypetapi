package lib

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

const (
	NOTIFICATION_EXCHANGE = "NOTIFICATION_EXCHANGE"
	NOTIFICATION_QUEUE    = "NOTIFICATION_QUEUE"
)

type Producer interface {
	Setup() error
	Publish(postId, commentId string) error // TODO how to handle errors
}

type AmqpProducer struct {
	connection *amqp.Connection
	db         *sql.DB
}

type Message struct {
	PostId             string `json:"post_id"`
	PostAuthorEmail    string `json:"post_author_email"`
	PostAuthorName     string `json:"post_author_name"`
	CommentAuthorEmail string `json:"comment_author_email"`
	CommentAuthorName  string `json:"comment_author_name"`
}

func NewAmqpProducer(connection *amqp.Connection, db *sql.DB) (*AmqpProducer, error) {
	producer := AmqpProducer{connection: connection, db: db}

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

	err = p.declareExchange(channel)
	if err != nil {
		log.Printf("Exchange declarion went wrong %v\n", err)
		return err
	}

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
	_, err := channel.QueueDeclare(NOTIFICATION_QUEUE, true, false, false, false, nil)
	if err != nil {
		log.Printf("Error declaring queue %v\n", err)
		return err
	}

	return nil
}

func (p *AmqpProducer) Publish(postId, commentId string) error {
	channel, err := p.connection.Channel()
	if err != nil {
		log.Printf("Error getting channel %v\n", err)
		return err
	}
	defer channel.Close()

	msg, err := p.getMessage(postId, commentId)
	if err != nil {
		log.Printf("Error generating message %v\n", err)
		return err
	}

	// no reason to generate a message if who commented
	// was the post author
	// TODO test it out
	if len(msg) == 0 {
		return nil
	}
	err = channel.Publish("", NOTIFICATION_QUEUE, false,
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

func (p *AmqpProducer) getMessage(postId, commentId string) ([]byte, error) {
	stmt, err := p.db.Prepare(`
		select pu.email, pu.name, cou.email, cou.name
		from comment as c
		join post as p on c.post_id = p.id
		join userpet as cou on c.author_id = cou.id
		join userpet as pu on p.author_id = pu.id
		where c.id = $1
	`)
	if err != nil {
		log.Printf("Error statement %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	var msg Message
	err = stmt.QueryRow(commentId).Scan(&msg.PostAuthorEmail,
		&msg.PostAuthorName, &msg.CommentAuthorEmail,
		&msg.CommentAuthorName)
	msg.PostId = postId
	if err != nil {
		log.Printf("Error querying %v\n", err)
		return nil, err
	}
	if msg.PostAuthorEmail == msg.CommentAuthorEmail {
		return []byte(""), nil
	}

	mBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return mBytes, nil
}
