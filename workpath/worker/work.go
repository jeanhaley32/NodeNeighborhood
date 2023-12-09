package worker

// Worker is a package designed to execute a task in a goroutine,
// with a variable map and task specific metrics. It is
// designed to be used with the delegator package. The task
// is defined as an enum, which allows for easy addition of new tasks.
import (
	"sync"
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
		return "running"
	case completed:
		return "completed"
	case failed:
		return "failed"
	}
	return ""
}

type job struct {
	id   uint32 // Unique id of the job
	work work   // The work to be done.
}

// The task interface.
type task interface {
	// defines the task signature.
	Func() TaskSignature
}

// return the state of the job.
func (j *job) GetState() state {
	w := &j.work
	if !w.IsDone() {
		return running
	}
	if w.Error() != nil {
		return failed
	}
	return completed
}

// Returns the unique id of the job.
func (j *job) ID() uint32 {
	return j.id
}

// Runs the job in a goroutine.
func (j *job) Run(done chan delegator.Directive) {
	w := &j.work
	var wg sync.WaitGroup
	defer func() {
		j.logError()
		// Log the completion of the job.
		log.Printf("Job %d finished %d", j.id, j.GetState())
		done <- delegator.NewDoneDirective(j.id)
	}()
	wg.Add(1)
	go w.execute(&wg)
	wg.Wait()
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
	w := work{
		task: t,
		vars: v,
		done: false,
	}
	return &job{
		id:   uuid.New().ID(),
		work: w,
	}
}
