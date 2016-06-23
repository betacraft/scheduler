package rmq

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var consumerCh, pubCh *amqp.Channel

// Gives the amqp connection i.e. rabbitmq connection
// should be used as read only, should not be edited at the
// user's side
func GetAMQPConn() *amqp.Connection {
	return conn
}

// Gives the rabbitmq publisher channel
// should be used as read only, should not be edited at the
// user's side.
func GetPubChannel() *amqp.Channel {
	return consumerCh
}

// Gives the rabbitmq consumer channel.
// Should be used as read only, should not be edited at the
// user's side.
func GetConsumerChannel() *amqp.Channel {
	return pubCh
}

// Takes connection url of the rabbitmq, and creates a connection
func DialConn(conUrl string) (*amqp.Connection, error) {
	var err error
	conn, err = amqp.Dial(conUrl)
	return conn, err
}

// Initiates the publisher channel, and returns it
// once Initiated, pubCh should not be tampered with
func InitPubChannel() (*amqp.Channel, error) {
	var err error
	pubCh, err = conn.Channel()
	return pubCh, err
}

// Inititates Consumer channel, and returns it.
// Consumer channel returned here should not be tampered with.
// this is used by the Monitor() method which must be ran as a goroutine
func InitConsumerChannel() (*amqp.Channel, error) {
	var err error
	consumerCh, err = conn.Channel()
	return consumerCh, err
}

// Takes list of RMQConfig, this is used to the queues
// with a delayed exchange, ExchangeName is used for the
// exchange creation. Setup call is idempotent, so it can be
// called as many time as wished, although one time should suffice.
// Note that all the exchange created is a delayed exchange, and
// is of type topic, so routing is based on topic,
// other types of exchange are not supported yet.
func Setup(exchangeName string, configs []RMQConfig) error {
	args := amqp.Table{}
	args["x-delayed-type"] = "topic" // topic based routing
	err := pubCh.ExchangeDeclare(
		exchangeName,        // name
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
		err = declareAndBind(exchangeName, v.QueueName, v.RoutingKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func declareAndBind(exchangeName, queueName, routingKey string) error {
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
		exchangeName, // exchange
		false,
		nil)
	if err != nil {
		log.Fatal("Error binding queue ", err)
		return err
	}
	return nil
}
