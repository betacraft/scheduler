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
	ce := CustomExecutor{}
	jobs.RegisterExecutor("CustomerExecutor", ce)

	conf := jobs.Config{
		QueueName:  "test-queue",  // Note that this is same as the sqs_publisher example
		RegionName: "APSoutheast", // Note the region is same as the region in the Setup() call in sqs_publisher example
	}

	// Start Monitoring Job
	done := make(chan bool)
	go jobs.Monitor(conf)
	<-done
}

// Implementation of Executor implementation
type CustomExecutor struct {
	Name string `json:"name"`
}

func (ce CustomExecutor) New() jobs.Executor {
	return CustomExecutor{}
}

func (ce CustomExecutor) Execute(j *jobs.Job) error {
	// print job info
	fmt.Println(j.ID)
	fmt.Println(ce.Name)

	// job info could be changed here
	// not changing in state of the job so, it would be recurring
	return nil
}
