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

	// Setup queue, with 30 seconds delay
	sqsConf := sqs.NewSQSConfig("APSoutheast", "test-queue", sqs.MIN_QUEUE_DELAY)
	err := sqs.Setup(sqsConf)
	if err != nil {
		log.Println("Error creating queue", err)
		return
	}
	log.Println("Queue created successfully")

	// Create a Job
	// Check that job data is of the form of CustomExector{} in consumer
	// CustomExector implements the Executor interface, also stored some job
	// related data
	jobData := map[string]interface{}{
		"name": "John Snow",
	}

	j := &jobs.Job{
		ID:          uuid.New(),
		EnqueueTime: time.Now().UTC(),
		Type:        "CustomerExecutor", // this will be same while registering
		JobData:     jobData,
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
