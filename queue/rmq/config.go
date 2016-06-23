package rmq

// RMQConfig is used as the struct to keep a queue configuration.
type RMQConfig struct {
	// Name of the queue to be created
	QueueName string

	// Routing Key for queue
	RoutingKey string
}

func NewRMQConfig(exchangeName, queueName, routingKey string) RMQConfig {
	return RMQConfig{QueueName: queueName, RoutingKey: routingKey}
}
