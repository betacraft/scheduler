package sqs

import (
	"encoding/json"
	"log"
	"runtime"
	"strconv"
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
		log.Print("error getting region:", err)
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
			log.Print("error getting message:", q.Name, err)
			continue
		}
		if len(msgs.Messages) < 1 {
			log.Print("no messages received:", q.Name, err)
			continue
		}
		message := msgs.Messages[0]
		_, err = q.DeleteMessage(&message)
		if err != nil {
			log.Print("error deleting message:", q.Name, err)
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
			log.Print("error unmarshalling job", err)
			continue
		}
		now := time.Now().UTC()
		log.Print("now for job: ", j.Type, j.ID, now)
		log.Print("exectime for job: ", j.Type, j.ID, j.ExecTime)

		// run the job if exec time is less than current time
		// or equal to current time
		if now.After(j.ExecTime) || now.Equal(j.ExecTime) {
			// panic point if GetExecutor returns nil
			err := j.GetExecutor().Execute(j)
			if err != nil { // don't enqueue if err is found
				log.Print("error executing job", err)
			}
			j.ExecTime = time.Now().Add(time.Duration(j.Interval) * time.Millisecond).UTC()
		} else {
			log.Print("job not executed as exectime is more: ", j.Type, j.ID)
		}

		if j.IsRecurring {
			log.Print("next exectime for job: ", j.Type, j.ID, j.ExecTime)
			err = jobs.Enqueue(j)
			if err != nil {
				log.Print("error enqueuing job: ", err)
				continue
			}
		}
	}
}

func (d *sqsdoer) Enqueue(j *jobs.Job) error {
	s, err := SQS(j.QueueRegion)
	if err != nil {
		log.Print("error getting region:", err)
		return err
	}
	q, err := s.GetQueue(j.Queue)
	if err != nil {
		log.Print("error getting queue", err)
		return err
	}
	res, err := json.Marshal(j)
	if err != nil {
		log.Print("Error marshaling job", err)
		return err
	}
	attrs := map[string]string{
		"DelaySeconds": getDelaySeconds(j.Interval),
	}

	log.Print("message being sent:", string(res), j.ID)
	_, err = q.SendMessageWithAttributes(string(res), attrs)

	if err != nil {
		log.Print("error sending message", err)
		return err
	}
	return nil
}

//expects interval to be in milisecs, returns "900" if >= 900000
// else returns "x" where x is < 90000/1000
func getDelaySeconds(interval int64) string {
	var inSeconds int64
	inSeconds = interval / 1000
	if inSeconds >= 900 {
		return "900"
	}
	return strconv.FormatInt(inSeconds, 10)
}
