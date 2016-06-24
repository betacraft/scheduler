## Scheduler

Can schedule, recurring jobs and then execute them. Uses RabbitMQ or AWS SQS. More documentation on the way !

Examples will be updated soon

## Note
Broken as of now. Currently Not working

## Dependencies
* github.com/mitchellh/mapstructure
* github.com/betacraft/goamz/sqs
* github.com/streadway/amqp
* github.com/go-ini/ini

## Godocs
* [Jobs package](https://godoc.org/github.com/betacraft/scheduler/jobs)
* [RabbitMQ implementation](https://godoc.org/github.com/betacraft/scheduler/queue/rmq)
* [AWS SQS implementation](https://godoc.org/github.com/betacraft/scheduler/queue/sqs)

## TODOs:
* Write examples
* Write documentation for sqs and rmq
* Write documentation for setting up rmq with delayed_message_plugin
* Make making delayed queue in sqs, idempotent, currently if a queue is created, and again create is called, the call fails
* Improve logging
