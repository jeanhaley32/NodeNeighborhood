package worker

// A work object encapsulates the task a job is to perform.
// It wraps up the task, the payload, and all of the necessary
// variables needed to perform the task. This can then be passed
// to the delegator to be executed as a worker, or from the worker
// to an ingestion point.
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
	TaskSignature *func(*VMap) ([]byte, error)
	task          interface { // The task interface.
		Func() TaskSignature
	}
)

func (w *work) SetTask(t task) {
	w.task = t
}

// Set work payload.
func (w *work) SetPayload(p []byte) {
	w.payload = p
}

// Initializes the variable map.
func (w *work) SetVmap(v map[string]any) {
	vMap := NewVmap()
	for k, v := range v {
		vMap.Set(k, v)
	}
	w.vars = *vMap
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
	t := w.task.Func()
	payload, err := (*t)(w.GetVmap())
	w.SetPayload(payload)
	if err != nil {
		w.e = err
		return
	}
	w.e = nil
}
