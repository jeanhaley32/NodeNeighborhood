package main

import (
	"fmt"
	"worker"
)

func main() {
	worker := worker.NewWorker(HelloWorld, nil)
	wchan := worker.Run()
	if worker.Error() != nil {
		fmt.Println(worker.Error())
	}
	<-wchan
}
