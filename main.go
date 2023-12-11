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
		go w.Run(dchan)
		d := <-dchan
		fmt.Printf("received Directive: %v %v\n\n", d.Target(), d.Action())
		time.Sleep(3 * time.Second)
	}
}
