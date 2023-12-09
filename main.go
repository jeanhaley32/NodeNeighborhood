package main

import (
	"fmt"
	"time"
	"workpath/delegator"
	"workpath/worker"
)

func main() {
	dchan := make(chan delegator.Directive)
	w := worker.NewJob(HelloWorld, nil)
	for {
		w.Run(dchan)
		fmt.Println("Job sent to worker")
		<-dchan
		fmt.Println("Job completed")
		time.Sleep(3 * time.Second)
	}
}
