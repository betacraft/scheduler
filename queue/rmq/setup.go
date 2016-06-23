package rmq

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var consumerCh, pubCh *amqp.Channel

func GetAMQPConn() *amqp.Connection {
	return conn
}

func GetPubChannel() *amqp.Channel {
	return consumerCh
}

func GetConsumerChannel() *amqp.Channel {
	return pubCh
}

func DialConn(conUrl string) (*amqp.Connection, error) {
	var err error
	conn, err = amqp.Dial(conUrl)
	return conn, err
}

func InitPubChannel() (*amqp.Channel, error) {
	var err error
	pubCh, err = conn.Channel()
	return pubCh, err
}

func InitConsumerChannel() (*amqp.Channel, error) {
	var err error
	consumerCh, err = conn.Channel()
	return consumerCh, err
}

func Setup(configs []RMQConfig) error {
	args := amqp.Table{}
	args["x-delayed-type"] = "topic" // topic based routing
	err := pubCh.ExchangeDeclare(
		"droidcloud",        // name
		"x-delayed-message", // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		args,                // arguments
	)

	if err != nil {
		log.Fatal("Error creating exchange ", err)
		return err
	}
	for _, v := range configs {
		err = declareAndBind(v.QueueName, v.RoutingKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func declareAndBind(queueName, routingKey string) error {
	// declare queue for apk scrape
	q, err := pubCh.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal("Error creating queue ", err)
		return err
	}
	err = pubCh.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		"droidcloud", // exchange
		false,
		nil)
	if err != nil {
		log.Fatal("Error binding queue ", err)
		return err
	}
	return nil
}
