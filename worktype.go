package main

// We define the task type as an enum.
// This allows us to easily add new tasks.

import (
	"fmt"
	"workpath/worker"
)

type task int64

const (
	HelloWorld task = iota // The task to be performed.
)

// Defining the functions to be returned by the task enum.
// We are returning a pointer to the function, so that we can define
// more complex functions, without worrying about passing the function
// around.
func (t task) Func() worker.TaskSignature {
	switch t {
	case HelloWorld:
		t := func(v *worker.VMap) ([]byte, error) {
			_, err := fmt.Println("Hello World")
			return nil, err
		}
		f := &t
		return f
	}
	return nil
}

func (t task) String() string {
	switch t {
	case HelloWorld:
		return "Hello World"
	}
	return "No task defined"
}
