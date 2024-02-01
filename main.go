package main

import (
	"log"
	"time"

	"github.com/jeanhaley32/logger"
	"github.com/jeanhaley32/nodeneighborhood/worker"
)

var (
	taskList = []task{HelloWorld, Timeout}
)

func main() {
	l := logger.StartLogger(log.Default())
	for _, t := range taskList {
		dchan := make(chan delegator.Directive)
		w := worker.NewJob(t, nil)
		go w.Run(dchan)
		<-dchan
		time.Sleep(3 * time.Second)
	}
	l.Info("Done")
}
