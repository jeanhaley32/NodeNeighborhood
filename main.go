package main

import (
	"time"
	"workpath/delegator"
	"workpath/worker"
)

func main() {
	dchan := make(chan delegator.Directive)
	w := worker.NewJob(HelloWorld, nil)
	for {
		w.Run(dchan)
		<-dchan
		time.Sleep(3 * time.Second)
	}
}
