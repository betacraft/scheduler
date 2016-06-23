package sqs

import (
	"encoding/json"
	"log"
	"runtime"
	"time"

	"github.com/betacraft/goamz/sqs"
	"github.com/betacraft/scheduler/jobs"
)

var maxUsableProcs int

func Init() {
	maxUsableProcs = runtime.NumCPU()
	log.Print("registering job scheduler with sqs")
}

func init() {
	d := new(sqsdoer)
	jobs.RegisterDoer(d)
}

type sqsdoer struct {
}

func (d *sqsdoer) Monitor(c jobs.Config) {
	s, err := SQS(c.RegionName)
	if err != nil {
		log.Fatal("error getting region:", err)
		return
	}
	q, err := s.GetQueue(c.QueueName)
	if err != nil {
		log.Print("error getting queue", err)
		return
	}

	messages := make(chan sqs.Message, 2*maxUsableProcs)
	go processJob(q, messages)
	done := make(chan bool)
	for {
		msgs, err := q.ReceiveMessage(1)
		if err != nil {
			log.Fatal("error getting message:", q.Name, err)
			continue
		}
		if len(msgs.Messages) < 1 {
			log.Print("no messages received:", q.Name, err)
			continue
		}
		message := msgs.Messages[0]
		_, err = q.DeleteMessage(&message)
		if err != nil {
			log.Fatal("error deleting message:", q.Name, err)
			continue
		}
		log.Print("deleted message with receipt:", message.MessageId, c.QueueName)
		messages <- message
	}
	<-done
}

func processJob(q *sqs.Queue, messages chan sqs.Message) {
	for {
		msg := <-messages
		j := &jobs.Job{}
		err := json.Unmarshal([]byte(msg.Body), j)
		if err != nil {
			log.Fatal("error unmarshalling job", err)
			continue
		}
		now := time.Now().UTC()
		log.Print("now for job", j.Type, j.ID, now)
		log.Print("exectime for job", j.Type, j.ID, j.ExecTime)

		// run the job if exec time is less than current time
		// or equal to current time
		// or if current time is 1 minute behind exectime
		if now.After(j.ExecTime) || now.Equal(j.ExecTime) || now.After(j.ExecTime.Add(-1*time.Minute)) {
			// panic point if GetExecutor returns nil
			err := j.GetExecutor().Execute(j)
			if err != nil { // don't enqueue if err is found
				log.Fatal("error executing job", err)
			}
			j.ExecTime = time.Now().Add(time.Duration(j.Interval) * time.Millisecond).UTC()
		} else {
			log.Print("job not executed as exectime is more", j.Type, j.ID)
		}

		if j.IsRecurring {
			log.Print("next exectime for job", j.Type, j.ID, j.ExecTime)
			err = jobs.Enqueue(j)
			if err != nil {
				log.Fatal("error enqueuing job")
				continue
			}
		}
	}
}

func (d *sqsdoer) Enqueue(j *jobs.Job) error {
	s, err := SQS(j.QueueRegion)
	if err != nil {
		log.Fatal("error getting region:", err)
		return err
	}
	q, err := s.GetQueue(j.Queue)
	if err != nil {
		log.Fatal("error getting queue", err)
		return err
	}
	res, err := json.Marshal(j)
	if err != nil {
		log.Fatal("Error marshaling job", err)
		return err
	}
	log.Print("message being sent:", string(res), j.ID)
	_, err = q.SendMessage(string(res))
	if err != nil {
		log.Fatal("error sending message", err)
		return err
	}
	return nil
}
