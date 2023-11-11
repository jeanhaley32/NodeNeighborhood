package main

import "time"

// Worker is a simple goroutine that performs a task, and sends the result to the writer.
// It then signals to the delegator that the work is done.

type worker struct {
	target   string           // UNique Identifier for the target.
	halt     chan interface{} // channel to receive halt signal.
	writerCh chan interface{} // Channel to send complete work to the writer.
	done     chan workItem    // Done chanel, used to signal to the delegator that work is done.
}

// Work waits for a halt signal, and if none is received runs a 'task' that returns
// an empty interface, and sends the result to the writer.

func (w *worker) run(task func(target string) interface{}) {
	startTime := time.Now()
	for {
		select {
		case <-w.halt:
			return
		default:
			// perform task, and send response to the writer.
			w.writerCh <- task(w.target)
			w.done <- workItem{
				target:   w.target,              // Send the unique identifier back to the delegator.
				action:   done,                  // Signal that work is done.
				execTime: time.Since(startTime), // Calculate the time it took to perform the task.
			} // Signal that work is done.
		}
	}
}
