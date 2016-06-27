package sqs

const (
	MAX_QUEUE_DELAY = "900" // 15 minutes
	MIN_QUEUE_DELAY = "0"   // No delay
)

type SQSConfig struct {
	// Region where the queue is to be made
	// This should be the friendly name given in the
	// RegionNames map
	RegionName string

	// Name of the queue to be created
	QueueName string

	// Delay should be provided in seconds
	// Maximum seconds can be 900 i.e. 15 minutes
	Delay string
}

func NewSQSConfig(regionName, queueName, delay string) SQSConfig {
	return SQSConfig{RegionName: regionName, QueueName: queueName, Delay: delay}
}

// Takes a list of Config structs
// The method must be called after InitSQSRegions() is called
func Setup(configs ...SQSConfig) error {
	for _, v := range configs {
		_, err := CreateQueueWithDelay(v.QueueName, v.RegionName, v.Delay)
		if err != nil {
			// if err return, even in one of the queues creation
			// TODO: change this to idempotent
			return err
		}
	}
	return nil
}
