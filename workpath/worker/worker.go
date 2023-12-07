package worker

// Worker is a package designed to execute a task in a goroutine,
// with a variable map and task specific metrics. It is designed to be
// used with the delegator package.
// The task is defined as an enum, which allows for easy addition of new tasks.
import (
	"github.com/google/uuid"
)

type (
	VMap          map[string]any             // Variables used by the task.
	TaskSignature func(VMap) (error, []byte) // The task Signature expected by the worker.

)

type job struct {
	id      uint32 // Unique id of the job
	task    task   // The task to be performed
	vars    VMap   // Variables used by the task.
	payload []byte // The returned value of the task.
	e       error  // The error returned by the task.

}

type task interface { // The task interface.
	Func() TaskSignature // Returns the task signature.
}

// Adds an entry to the varible map.
func (v *VMap) Set(key string, value any) {
	(*v)[key] = value
}

// Returns the value associated with the key.
func (v *VMap) Get(key string) any {
	return (*v)[key]
}

// Removes the entry associated with the key.
func (v *VMap) Delete(key string) {
	delete(*v, key)
}

// Appends multiple entries to the variable map.
func (v *VMap) SetKeys(m map[string]any) {
	for key, val := range m {
		(*v)[key] = val
	}
}

// Creates a new, nil variable map.
func NewVmap() *VMap {
	return &VMap{}
}

// constructs a new job.
func NewWorker(t task, v map[string]any) *job {
	return &job{
		id:   uuid.New().ID(),
		task: t,
		vars: v,
	}
}

// Returns the unique id of the job.
func (j *job) ID() uint32 {
	return j.id
}

// Return Error value of the job.
func (j *job) Error() error {
	return j.e
}

// Returns the "variable map" used by the job.
// A pointer is returned, which allows the user to modify the map via
// the SetKeys, Get, and Delete methods.
func (j *job) GetVmap() *VMap {
	return &j.vars
}

func (j *job) execute(ch chan any) {
	err, payload := j.task.Func()(j.vars)
	if err != nil {
		j.e = err
		return
	}
	j.e = nil
	j.payload = payload
	ch <- nil // Signal that the job is done.
}

// Run executes the job's task as a goroutine.
// The job's error value is set to the error returned by the task.
// The job's payload value is set to the payload returned by the task.
// Run takes in a waitgroup, which is used to wait for the job to finish.
func (j *job) Run() chan any {
	ch := make(chan any)
	go j.execute(ch)
	return ch
}
