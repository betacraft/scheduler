package rmq

type RMQConfig struct {
	// Name of the queue to be created
	QueueName string

	// Routing Key for queue
	RoutingKey string
}

func NewRMQConfig(queueName, routingKey string) RMQConfig {
	return RMQConfig{QueueName: queueName, RoutingKey: routingKey}
}
