package main

// We define the task type as an enum.
// This allows us to easily add new tasks.

import (
	"context"
	"fmt"
	"time"

	"github.com/jeanhaley32/nodeneighborhood/worker"
)

type task int64

// First, define your task in the enum below
const (
	HelloWorld task = iota // The task to be performed.
	Timeout
)

// Defining the functions to be returned by the task enum.
// We are returning a pointer to the function, so that we can define
// more complex functions, without worrying about passing the function
// around.
// The general template for a task function is:
//
//	func (t task) Func() worker.TaskSignature {
//		t := func(c context.Context) ([]byte, error) {
//			// Do something
//			return nil, nil
//		}
//		return &t
//	}
//
// We encase this ina  switch statement to allow for multiple tasks.
// The function must return a pointer to the function, and the function
func (t task) Func() worker.TaskSignature {
	switch t {
	case HelloWorld:
		t := func(c context.Context) ([]byte, error) {
			_, err := fmt.Println("Hello World")
			return nil, err
		}
		return &t
	case Timeout:
		t := func(c context.Context) ([]byte, error) {
			for {
				select {
				case <-c.Done():
					return nil, c.Err()
				default:
					time.Sleep(1 * time.Second)
				}
			}
		}
		return &t
	}
	return nil
}

func (t task) String() string {
	switch t {
	case HelloWorld:
		return "Hello World"
	case Timeout:
		return "Mock Timeout test"
	}
	return "No task defined"
}
