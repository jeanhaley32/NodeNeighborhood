package worker

// The work encapsulates a job's task to perform.
// It also contains the payload, which is the result
// of the task.
// The work is designed to be executed in a goroutine.
// The work is considered done when the done flag is set to true.
type work struct {
	task    task
	done    bool
	payload []byte
	e       error
	vars    VMap
}

type Work interface {
	SetPayload(p []byte)
	GetPayload() []byte
	Error() error
	IsDone() bool
	GetVmap() *VMap
	execute()
}

type (
	// The task Signature expected by the worker.
	TaskSignature func(VMap) (error, []byte)
	task          interface { // The task interface.
		Func() TaskSignature
	}
)

// Set work payload.
func (w *work) SetPayload(p []byte) {
	w.payload = p
}

// Get work payload.
func (w *work) GetPayload() []byte {
	return w.payload
}

// Return Error value of the job.
func (w *work) Error() error {
	return w.e
}

func (w *work) IsDone() bool {
	return w.done
}

// Returns the "variable map" used by the job.
// A pointer is returned, modify the map via
// the SetKeys, Get, and Delete methods.
func (w *work) GetVmap() *VMap {
	return &w.vars
}

func (w *work) execute() {
	defer func() {
		w.done = true
	}()
	err, payload := w.task.Func()(w.vars)
	w.SetPayload(payload)
	if err != nil {
		w.e = err
		return
	}
	w.e = nil
}
