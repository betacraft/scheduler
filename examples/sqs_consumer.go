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
