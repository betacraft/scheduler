package jobs

import (
	"encoding/json"
	"log"
	"reflect"
	"time"
)

// Config is used by the Monitor() method
type Config struct {
	// Used in Monitor method
	QueueName string

	// This is required only for sqs based
	// schedler not for rmq, but is mandatory
	// for sqs based scheduler
	RegionName string
}

// This interface must be implemented by the user,
// New() will just return a new instnace of the struct
// and Execute() method will have the logic of running the job
type Executor interface {
	// Should return a new executor instance
	New() Executor

	// Execute will have all the logic associated with the job,
	// Job struct is passed as argument, so that any changes in the
	// Job related data like IsRecurring can be changed by the user
	Execute(j *Job) error
}

var executors map[string]*Executor

func init() {
	executors = map[string]*Executor{}
}

// Register method registers the job with a type of executor,
// This must be called for each new type of job, before jobs of
// that type are submitted
func RegisterExecutor(jobType string, executor Executor) {
	executors[jobType] = &executor
}

// Job defines the attributes required to run the job,
type Job struct {
	// ID, identifies a job uniquely, this must be set by the user, uuid.New() will suffice
	ID string `json:"id"`

	// Time when the the job is submitted, must be set by the user
	EnqueueTime time.Time `json:"enqueue_time"`

	// Type defines the type of interface that implements the Execute method,
	// this must be same as the key that is used to register the Executor interface
	Type string `json:"type"`

	// Interval is the time gap after which the is to be executed from the EnqueueTime,
	// This must be provided in milliseconds
	// For eg if the job is to be exectuded in 10 seconds, then 10000 should be the value
	Interval int64 `json:"interval"`

	// Routing Key is also used for rmq implementation it is used for routing of the job,
	// to the specific queue
	RoutingKey string `json:"routing_key"`

	// Queue name defines the name of the queue to which the message is sent, used for both.
	// Alhthugh, it is not required for rmq based implementation, it is recommended to be provided
	// for debuging purposes
	Queue string `json:"queue"`

	// This is required only for sqs based implementation,
	// This value must correspond to the region where the Queue is present,
	// and the name must be the friendly name, check queue package for the names
	// of all the regions
	QueueRegion string `json:"queue_region"`

	// If this is true then the job is executed after every specific interval,
	// this interval value is taken from the Interval attribute
	IsRecurring bool `json:"is_recurring"`

	// This is the time at which a job is to be executed,
	// this is used in sqs implementation and must be provided,
	// while submitting the job, it should be equal to EnqueueTime + Interval
	ExecTime time.Time `json:"exec_time"`

	// It is an interface, which could hold job specific data,
	// For eg: if a push notification is to be sent for a user, it could contain
	// UserId, and related data
	JobData map[string]interface{} `json:"job_data"`
}

func (j *Job) GetExecutor() Executor {
	// panic point, if exector not registered
	// log.Println("Executors registered: ", executors)
	// log.Println("Job Type: ", j.Type)
	data, err := json.Marshal(j.JobData)
	if err != nil {
		return nil
	}
	e := *executors[j.Type]
	executor := e.New()
	log.Println("type of interface: ", reflect.TypeOf(executor))
	err = json.Unmarshal(data, &executor)
	if err != nil {
		log.Println(err)
		return nil
	}
	return executor
}
