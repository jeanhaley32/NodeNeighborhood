package main

// We define the task type as an enum.
// This allows us to easily add new tasks.

import (
	"fmt"
	"workpath/worker"
)

type TaskType int64

const (
	HelloWorld TaskType = iota // The task to be performed.
)

func (t TaskType) Func() worker.TaskSignature {
	switch t {
	case HelloWorld:
		return func(_ worker.VMap) (error, []byte) {
			_, err := fmt.Println("Hello World")
			return err, nil
		}
	}
	return nil
}
