package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	workch  chan job                                 // channel for the communication of jobs.
	payload []byte                                   // The returned value of the task.
	success bool                                     // Whether the task was successful or not.
	task    func(context.Context) (success, payload) // The task to be performed
)
type job struct {
	id      int                                      // Unique id of the job
	task    func(context.Context) (success, payload) // The task to be performed
	success success                                  // Whether the task was successful or not.
	vars    map[string]any                           // Variables used by the task.
	p       payload                                  // The returned value of the task.
	metrics struct {                                 // Metrics of the task.
		created   time.Time     // The time the task was created.
		start     time.Time     // start time of task. Used to calculate the time taken to complete the task.
		taskTime  time.Duration // The Time taken to complete the task.
		completed time.Time     // The time the task was completed.
		complete  bool          // Whether the task is complete or not.
	}
	chBundle struct { // Bundle of channels used for communication.
		parent         workch // The parent channel.
		localComms     workch // The local channel.
		targetIngester workch // The target channel.
	}
}

// newJob creates a new job.
func NewJob(t task, v map[string]any) *job {
	j := &job{
		id:      uuid.New().ID(),
		task:    t,
		success: false,
		vars:    v,
	}
	j.metrics.created = time.Now()
	return j
}
