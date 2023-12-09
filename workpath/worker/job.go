package worker

import "sync"

type work struct {
	task    task   // The task to be performed.
	done    bool   // Indicates if the task has been run.
	payload []byte // The returned value of the task.
	e       error  // the error returned by the task.
	vars    VMap   // Variables used by the task.
}

type (
	// The task Signature expected by the worker.
	TaskSignature func(VMap) (error, []byte)
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

func (w *work) execute(wg *sync.WaitGroup) {
	defer func() {
		w.done = true
		wg.Done()
	}()
	err, payload := w.task.Func()(w.vars)
	w.SetPayload(payload)
	if err != nil {
		w.e = err
		return
	}
	w.e = nil
}
