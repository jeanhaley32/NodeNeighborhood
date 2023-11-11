package main

import (
	"fmt"
	"time"
)

// The Delegator ingests workType objects, and performs the appropriate action.
type WorkType uint32

// Defines the types of work that can be performed.
const (
	start WorkType = iota
	halt
	shutdown
	done
)

// Stringer returns a string representation of the workType object.
func (w WorkType) String() string {
	switch w {
	case start:
		return "start"
	case done:
		return "done"
	case halt:
		return "stop"
	case shutdown:
		return "shutdown"
	default:
		return "unknown"
	}
}

// workItem defines an object of work, that includes A Unique Identifier, and an action to be performed on that identifier.
type workItem struct {
	target   string
	action   WorkType
	execTime time.Duration
}

// Stringer returns a string representation of the work object.
// TODO(JeanHaley) - This should be used for future logging.
func (w *workItem) String() string {
	return fmt.Sprintf("UID: %s, action: %s", w.target, w.action.String())
}

// Defines the struct for the Delegator routine.
type delegator struct {
	work    chan workItem // Channel to receive work on.
	workers map[string]*worker
}

// Delegate Starts a new delegator routine, and returns a channel for communication.
func Delegate() chan<- workItem {
	d := &delegator{
		work:    make(chan workItem),
		workers: make(map[string]*worker),
	}

	go d.run()
	// TODO(jeanhaley) - Flagging this for potential point of failure.
	// I don't know if this go-routine will continue to run after we return the work channel.
	// My hope is that since we didn't defer the close, it will continue to run. And we can
	// shutdown the process through the work channel.
	return d.work
}

// run is the main loop for the delegator routine.
func (d *delegator) run() {
	shutdownRun := false
	for { // Loop forever.
		if shutdownRun && len(d.workers) == 0 { // If we have shutdown, and there are no more workers.
			close(d.work) // Close the work channel.
			return        // Exit the loop.
		}
		select { // Wait for work.
		case w := <-d.work: // Receive work.
			switch w.action { // Perform the appropriate action.
			case start: // Start a new worker.
				d.workers[w.target] = &worker{ // Create a new worker.
					target:   w.target,               // Set the target.
					halt:     make(chan interface{}), // Create a halt channel.
					writerCh: make(chan interface{}), // Create a channel to send completed work to the writer.
					done:     d.work,                 // Set the done channel to the delegator work channel
				}
				go d.workers[w.target].run(
					// This is the task that the worker will perform
					func(target string) interface{} {
						time.Sleep(1 * time.Second) // Simulate work.
						return nil                  // Return empty interface.
					},
				) // Start the worker.
			case halt:
				d.workers[w.target].halt <- struct{}{} // Send halt signal to worker
				delete(d.workers, w.target)            // Remove worker from map.
			case shutdown: // if we receive the shutdown signal, send halt to all workers
				for _, worker := range d.workers {
					worker.halt <- struct{}{}
				}
				shutdownRun = true // Set shutdown flag to true.
			default:
				fmt.Printf("Unknown action: %s\n", w.action)
			}
		}
	}
}
