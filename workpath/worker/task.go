package worker

import (
	"context"
	"time"
)

// A task object encapsulates the function a job is to perform.
// It wraps up the opertaion, along with the payload, and passes
// contextual variables through the context.Context object.
type task struct {
	op      operation
	done    bool
	payload []byte
	e       error
	ctx     context.Context
}

type Task interface {
	SetPayload(p []byte)
	GetPayload() []byte
	Error() error
	IsDone() bool
	GetContext() *context.Context
	execute()
}

type (
	// The task Signature expected by the worker.
	TaskSignature *func(context.Context) ([]byte, error)
	operation     interface { // The task interface.
		Func() TaskSignature
		String() string
	}
)

func (t *task) SetTask(op operation) {
	t.op = op
}

// Set work payload.
func (t *task) SetPayload(p []byte) {
	t.payload = p
}

// Initializes the variable map.
func (t *task) SetContext(ctx context.Context) {
	t.ctx = ctx
}

// Get work payload.
func (t *task) GetPayload() []byte {
	return t.payload
}

// Return Error value of the job.
func (t *task) Error() error {
	return t.e
}

func (t *task) IsDone() bool {
	return t.done
}

// Returns Context used to execute the task.
// if the task's context is nil, a new context is created.
func (t *task) GetContext() *context.Context {
	if t.ctx == nil {
		t.ctx = context.Background()
	}
	return &t.ctx
}

func (t *task) execute() time.Time {
	defer func() {
		t.done = true
	}()
	startTime := time.Now()
	select {
	case <-t.ctx.Done():
		t.e = t.ctx.Err()
		t.payload = nil
	default:
		t.payload, t.e = (*t.op.Func())(t.ctx)
	}
	return startTime
}
