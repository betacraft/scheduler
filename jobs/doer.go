package jobs

type Doer interface {
	Enqueue(j *Job) error
	Monitor(c Config)
}

var doer Doer

func RegisterDoer(d Doer) {
	doer = d
}

func Enqueue(j *Job) error { return doer.Enqueue(j) }
func Monitor(c Config)     { doer.Monitor(c); return }
