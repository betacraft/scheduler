package rmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/betacraft/scheduler/jobs"
	"github.com/go-ini/ini"
	"github.com/streadway/amqp"
)

func init() {
	d := new(rmqdoer)
	jobs.RegisterDoer(d)
}

type rmqdoer struct {
}

func (d *rmqdoer) Enqueue(j *jobs.Job) error {
	res, err := json.Marshal(j)
	if err != nil {
		log.Print("Error marshaling job", err)
		return err
	}
	headers := amqp.Table{}
	headers["x-delay"] = j.Interval
	pub := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/json",
		Body:         res,
		Headers:      headers,
	}
	err = pubCh.Publish("droidcloud", j.RoutingKey, false, false, pub)
	if err != nil {
		pubCh.Close()
		pubCh, err = conn.Channel()
		if err != nil { // recreate channel
			conn.Close()
			restartConn()
			pubCh, err = conn.Channel()
		}
		err = pubCh.Publish("droidcloud", j.RoutingKey, false, false, pub)
	}
	log.Print(fmt.Sprintf("Enqueued JobID: %s, JobType: %s, delay: %d", j.ID, j.Type, j.Interval))
	return err
}

func (d *rmqdoer) Monitor(c jobs.Config) {
	for {
		consume(c.QueueName)
		restartConn()
		restartPubChannel()
		restartConsumerChannel()
	}
}

func consume(qname string) {
	log.Print("starting consumer for ", qname)
	consName := fmt.Sprintf("%s-consumer", qname)
	msgs, err := consumerCh.Consume(
		qname,    // queue
		consName, // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)

	if err != nil {
		log.Print("queue consumer could not be initiated:", err)
		return
	}

	done := make(chan bool)
	go func() {
		for {
			// Check if channel is up
			d := <-msgs
			log.Print("received message to execute")
			defer rescue()               // recover in case of panicks, and wait for other messages
			go func(del amqp.Delivery) { // start a goroutine to handle a message delivery
				j := &jobs.Job{}
				err := json.Unmarshal(del.Body, &j)

				// Ack message always
				defer func(d amqp.Delivery) {
					log.Print(fmt.Sprintf("Ack delivery, JobID: %s, JobType: %s, consumer: %s", j.ID, j.Type, d.ConsumerTag))
					d.Ack(false)
				}(del)

				// Do nothing if unmarshalling fails
				if err != nil {
					log.Print("Error converting message body to job: ", err)
					return
				}

				// panic point if GetExecutor returns nil
				err = j.GetExecutor().Execute(j)
				if err != nil { // do not enqueue if execute returns error
					log.Print("error executing job: ", err)
					return
				}

				log.Print("succesfully executed job")
				if j.IsRecurring {
					log.Print(fmt.Sprintf("re-enqueing, JobID: %s, JobType: %s", j.ID, j.Type))
					err = jobs.Enqueue(j)
					if err == nil { // succesfully enqueued
						return
					}

					// failure enqueueing
					switch err.(type) {
					case amqp.Error:
						err, _ = err.(amqp.Error)
						done <- true
					default:
						log.Print("marshal error do nothing")
					}
				}
			}(d)
		}
	}()
	<-done
}

func rescue() {
	if r := recover(); r != nil {
		log.Print("recover from panic: ", r)
	}
}

func restartConn() {
	cfg, _ := ini.Load("config.ini")
	config, _ := cfg.GetSection(os.Getenv("ENV"))
	var err error
	conn.Close()
	conn, err = DialConn(config.Key("rabbitmq_dial_url").String())
	if err != nil {
		log.Print("error re-initiating connection: ", err)
		panic(err)
	}
	log.Print("re-initiated connection")
}

func restartPubChannel() {
	pubCh.Close()
	var err error
	pubCh, err = conn.Channel()
	if err != nil {
		log.Print("error re-initiating publisher channel: ", err)
		panic(err)
	}
	log.Print("re-initiated publisher channel")
}

func restartConsumerChannel() {
	consumerCh.Close()
	var err error
	consumerCh, err = conn.Channel()
	if err != nil {
		log.Print("error re-initiating consumer channel: ", err)
		panic(err)
	}
	log.Print("re-initiated consumer channel")
}
