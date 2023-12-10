package worker

// Worker is a package designed to execute a task in a goroutine,
// with a variable map and task specific metrics. It is
// designed to be used with the delegator package. The task
// is defined as an enum, which allows for easy addition of new tasks.
import (
	"fmt"
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
	created   time.Time // The time the job was created.
	started   time.Time // The time the job was started.
	completed time.Time // The time the job was completed.
	work      work      // The work to be done.
}

// Worker interface represents a job that can be executed.
type Worker interface {
	GetState() state
	ID() uint32
	Run(done chan delegator.Directive)
	logError()
}

// return the state of the job.
func (j *job) GetState() state {
	w := &j.work
	if w.IsDone() {
		if w.Error() != nil {
			return failed
		}
		return running
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
		j.Announce()
		fmt.Println("waiting to send completion message to delegator")
		done <- delegator.NewDoneDirective(j.id)
		fmt.Println("passed to delegator")
	}()
	wg.Add(1)
	func() { go w.execute(); wg.Done() }()
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
	j := &job{}
	w := &j.work
	w.SetTask(t)
	w.SetVmap(v)
	j.created = time.Now()
	j.id = uuid.New().ID()
	return j
}

func (j *job) Announce() {
	// Log the completion of the job.
	log.Printf("Job %d %v", j.id, j.GetState().String())
	j.logError()
}

// returns pointer to the job's work.
func (j *job) ObtainWork() *work {
	return &j.work
}
