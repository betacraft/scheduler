package sqs

import (
	"errors"
	"log"

	"github.com/betacraft/goamz/aws"
	"github.com/betacraft/goamz/sqs"
)

var SQSRegions map[string]*sqs.SQS

func InitSQSRegions(aws_access, aws_secret string) {
	auth := aws.Auth{AccessKey: aws_access, SecretKey: aws_secret}
	SQSRegions = make(map[string]*sqs.SQS)
	for k, v := range RegionNames {
		SQSRegions[k] = sqs.New(auth, aws.Regions[v])
	}
}

func SQS(regionName string) (*sqs.SQS, error) {
	s, ok := SQSRegions[regionName]
	if ok == false {
		return nil, errors.New("Region Name Not found")
	}
	return s, nil
}

func CreateQueue(regionName, queueName string) (*sqs.Queue, error) {
	s, err := SQS(regionName)
	if err != nil {
		return nil, err
	}
	attrs := map[string]string{
		"VisibilityTimeout":             "30",
		"ReceiveMessageWaitTimeSeconds": "10",
		"MessageRetentionPeriod":        "900", // 15 minutes
	}
	q, err := s.CreateQueueWithAttributes(queueName, attrs)
	return q, err
}

// Give the queueName, regionName and delay in seconds as string
// max delay possible is 900 seconds i.e 15 minutes
func CreateQueueWithDelay(queueName, regionName, delay string) (*sqs.Queue, error) {
	s, err := SQS(regionName)
	if err != nil {
		return nil, err
	}
	attrs := map[string]string{
		"VisibilityTimeout":             "30",
		"ReceiveMessageWaitTimeSeconds": "20",
		"DelaySeconds":                  delay,
	}
	q, err := s.CreateQueueWithAttributes(queueName, attrs)
	return q, err

}

func DeleteQueue(regionName, queueName string) error {
	s, err := SQS(regionName)
	if err != nil {
		return err
	}
	q, err := s.GetQueue(queueName)
	if err != nil {
		return err
	}
	_, err = q.Delete()
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func SendMessage(regionName, queueName, msg string) error {
	s, err := SQS(regionName)
	if err != nil {
		log.Fatal("error getting sqs object", err)
		return err
	}
	q, err := s.GetQueue(queueName)
	if err != nil {
		log.Fatal("error getting queue with name: ", queueName, err)
		return err
	}
	_, err = q.SendMessage(msg)
	if err != nil {
		log.Fatal("error sending message ", err)
		return err
	}
	return nil
}
