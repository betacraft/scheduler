## Scheduler

Can schedule, recurring jobs and then execute them. Uses rabbitmq or sqs. More documentation on the way !

Examples will be updated soon

## Dependencies
* github.com/mitchellh/mapstructure
* github.com/betacraft/goamz/sqs
* github.com/streadway/amqp
* github.com/go-ini/ini

## TODOs:
* Write examples
* Write documentation for sqs and rmq
* Write documentation for setting up rmq with delayed_message_plugin
* Make making delayed queue in sqs, idempotent, currently if a queue is created, and again create is called, the call fails
