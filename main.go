package main

import (
	"fmt"
	"workpath/worker"
)

func main() {
	w := worker.NewWorker(HelloWorld, nil)
	wchan := w.Run()
	if w.Error() != nil {
		fmt.Println(w.Error())
	}
	<-wchan
}
