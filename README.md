## Scheduler

Can schedule, recurring jobs and then execute them. Uses RabbitMQ or AWS SQS. More documentation on the way !

Examples will be updated soon

## Build Status
[![CircleCI](https://circleci.com/gh/betacraft/scheduler.svg?style=svg)](https://circleci.com/gh/betacraft/scheduler)


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

## Examples(SQS)

### Before running the example set the ENV variables
```
export AWS_ACCESS=your_access_key
export AWS_SECRET=your_secret_key
```

### Publisher (Submitting a job )
```Go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/betacraft/scheduler/jobs"
	"github.com/betacraft/scheduler/queue/sqs"
	"github.com/pborman/uuid"
)

func main() {
	awsAccess := os.Getenv("AWS_ACCESS")
	awsSecret := os.Getenv("AWS_SECRET")

	// register sqs implementation
	sqs.Init()

	// Init SQS SDK with creds
	sqs.InitSQSRegions(awsAccess, awsSecret)

	// Setup queue, with minimum delay i.e 0 second
	sqsConf := sqs.NewSQSConfig("APSoutheast", "test-queue", sqs.MIN_QUEUE_DELAY)
	err := sqs.Setup(sqsConf)
	if err != nil {
		log.Println("Error creating queue", err)
		return
	}
	log.Println("Queue created successfully")

	// Create a JobData
	// Check that job data is of the form of CustomJob{} in the consumer example
	// CustomJob implements the Executor interface, also stores some job related data
	type JobData struct {
		Name string `json:"name"`
	}
	jobData := JobData{Name: "John Snow"}

	// Create a Job, and fill up the required values,
	// then submit using enqueue
	j := &jobs.Job{
		ID:          uuid.New(),
		EnqueueTime: time.Now().UTC(),
		Type:        "CustomerExecutor", // this will be same while registering
		JobData:     &jobData,
		Interval:    45000,         // this is used for delaying the message in queue as well as the execution time is set in accordance
		RoutingKey:  "",            // not required for sqs
		Queue:       "test-queue",  // same as the queue created above
		QueueRegion: "APSoutheast", // QueueRegion is same as region in Setup(), not required if rmq
		IsRecurring: true,
		ExecTime:    time.Now().UTC().Add(3000 * time.Millisecond), // Setting the exectution time for forst submission, will be set by interval from next time onwards
	}

	// Submit the job using enqueue method
	err = jobs.Enqueue(j)
	if err != nil {
		log.Println("Error submitting job", err)
		return
	}

	// Check this ID, it should be printed when job is executed
	fmt.Println("Job successfully submitted with id", j.ID)
}
```


### Consumer(This will have the executor implementation)
```Go
package main

import (
	"fmt"
	"os"

	"github.com/betacraft/scheduler/jobs"
	"github.com/betacraft/scheduler/queue/sqs"
)

func main() {
	awsAccess := os.Getenv("AWS_ACCESS")
	awsSecret := os.Getenv("AWS_SECRET")

	// register sqs implementation
	sqs.Init()

	// Init SQS SDK with creds
	sqs.InitSQSRegions(awsAccess, awsSecret)

	// Register the job type with its executor
	ce := &CustomJob{}
	jobs.RegisterExecutor("CustomerExecutor", ce)

	// Start Monitoring Job
	conf := jobs.Config{
		QueueName:  "test-queue",  // Note that this is same as the sqs_publisher example
		RegionName: "APSoutheast", // Note the region is same as the region in the Setup() call in sqs_publisher example
	}
	done := make(chan bool)
	go jobs.Monitor(conf)
	<-done
}

// Custom implementation of Executor interface
type CustomJob struct {
	Name string `json:"name"`
}

func (ce *CustomJob) New() jobs.Executor {
	return &CustomJob{}
}

func (ce *CustomJob) Execute(j *jobs.Job) error {
	// print job info
	fmt.Println(j.ID)
	fmt.Println(ce.Name)

	// job info could be changed here
	// not changing in state of the job so, it would be recurring
	return nil
}
```


## For issues
* Raise them on Github
* Email at (abhishek@betacraft.co, abhishek.bhattacharjee11@gmail.com)
