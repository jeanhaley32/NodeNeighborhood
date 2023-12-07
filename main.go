package main

import (
	"fmt"
	"workpath/derectives"
	"workpath/worker"
)

func main() {
	dchan := make(chan derectives.Directive)
	w := worker.NewWorker(HelloWorld, nil)
	w.Run(dchan)
	if w.Error() != nil {
		fmt.Println(w.Error())
	}
	<-dchan
}
