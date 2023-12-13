package main

import (
	"fmt"
	"time"
	"workpath/delegator"
	"workpath/worker"
)

var (
	taskList = []task{HelloWorld, Timeout}
)

func main() {
	for _, t := range taskList {
		dchan := make(chan delegator.Directive)
		w := worker.NewJob(t, nil)
		go w.Run(dchan)
		<-dchan
		time.Sleep(3 * time.Second)
	}
	fmt.Println("done")
}
