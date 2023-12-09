## Synopsis
The Work path is responsible for performing various tasks in the Node Neighborhood application.
It handles background processing, such as data processing, file handling, and other computational tasks.
The module provides an interface for managing and monitoring worker processes, allowing for efficient and scalable execution of tasks.
 

 ## Job
 ### Job
 A job is an uninstantiated worker, it contains the `work` to be done, and handles the actuation of that `task`. 
 ```go
    type job struct {
        id   uint32 // Unique id of the job
        work work   // The work to be done.
    }

    type Worker interface {
        GetState() state
        ID() uint32
        Run(done chan delegator.Directive)
        logError()
    }
 ```
 ### Work
 `Work` is a bundled struct containing the `task`, and all relevant information for that `task` to execute.. 
 ```go
    type work struct {
        task    task
        done    bool
        payload []byte
        e       error
        vars    VMap
    }

    type Work interface {
        SetPayload(p []byte)
        GetPayload() []byte
        Error() error
        IsDone() bool
        GetVmap() *VMap
        execute()
    }
 ```
 ### task
  A `task` maps directly to a function. This relationship is created through a `task` interface that sets the expectation for a `TaskSignature`. The `TaskSignature` can be referenced externally in order for an external library to define a function.
 ``` go
    type (
        // The task Signature expected by the worker.
        TaskSignature func(VMap) (error, []byte)
        task          interface { // The task interface.
            Func() TaskSignature
        }
    )
 ```

 ## Work Delegator

