package worker

// Worker is a package designed to execute a task in a goroutine,
// with a variable map and task specific metrics. It is
// designed to be used with the delegator package. The task
// is defined as an enum, which allows for easy addition of new tasks.
import (
	"sync"
	"time"
	"workpath/delegator"

	"github.com/google/uuid"

	"log"
)

// Define enums used to indicate state of the job.
type state int64

const (
	running state = iota
	completed
	failed
)

func (a state) String() string {
	switch a {
	case running:
		return "still running"
	case completed:
		return "completed successfully"
	case failed:
		return "failed"
	}
	return ""
}

type job struct {
	id        uint32    // Unique id of the job.
	created   time.Time // job creation time.
	started   time.Time // task start time.
	completed time.Time // task completion time.
	work      work      // The work to be done.
}

// Worker interface represents a job that can be executed.
type Worker interface {
	GetState() state
	ID() uint32
	Run(done chan delegator.Directive)
	logError()
	Announce()
	ObtainWork() *work
	RunTime() time.Duration
}

// return the state of the job.
func (j *job) GetState() state {
	w := &j.work
	if w.IsDone() {
		if w.Error() != nil {
			return failed
		}
		return completed
	}
	return running
}

// Returns the unique id of the job.
func (j *job) ID() uint32 {
	return j.id
}

// Returns the time the job was started.
func (j *job) StartTime() time.Time {
	return j.started
}

// returns the time the job was completed.
func (j *job) CompletedTime() time.Time {
	return j.completed
}

// returns the time the job was created.
func (j *job) CreatedTime() time.Time {
	return j.created
}

// returns the elapsed time of the job.
func (j *job) RunTime() time.Duration {
	return j.completed.Sub(j.started)
}

// Runs the job in a goroutine.
func (j *job) Run(done chan delegator.Directive) {
	w := &j.work
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		j.started = w.execute()
		wg.Done()
	}()
	wg.Wait()
	j.completed = time.Now()
	j.Announce()
	d := delegator.NewDoneDirective(j.id)
	done <- d
}

// handle any errors returned by the task.
// If the error is not nil, it is logged.
func (j *job) logError() {
	if j.work.Error() != nil {
		log.Println(j.work.Error())
	}
}

// Creates a new, nil variable map.
func NewVmap() *VMap {
	return &VMap{}
}

// constructs a new job.
func NewJob(t task, v map[string]any) *job {
	j := &job{}
	w := &j.work
	w.SetTask(t)
	w.SetVmap(v)
	j.id = uuid.New().ID()
	j.created = time.Now()
	return j
}

// Log the completion of the job.
func (j *job) Announce() {
	log.Printf("worker %d finished task \"%v\" "+
		"%v with a Runtime of %vms, roundtrip time: %vms\n",
		j.id,
		j.work.task,
		j.GetState().String(),
		j.RunTime().Abs().Microseconds(),
		j.roundtrip().Abs().Microseconds())
	j.logError()
}

func (j job) roundtrip() time.Duration {
	return time.Since(j.created)
}
