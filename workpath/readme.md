## TODO
TODO(jeanhaley) Create unit tests for worker module.
## Synopsis
The Work path is responsible for executing individual tasks in the Node Neighborhood package. It receives `work` in the form of `jobs`, each `job` contains a `task` that points to an `operation` to perform, and `context` needed to perform that operation with a `deadline` 
 

 ## Job
 A Job is the Top level encapsulation of a task, it contains a unique `id`, time stampes for the job's `creation`, `start`, and `completion` times. and it contains the `task` to be executed. 
 ```go
    type job struct {
        id        uint32    // Unique id of the job.
        created   time.Time // job creation time.
        started   time.Time // task start time.
        completed time.Time // task completion time.
        task      task      // The work to be done.
    }
    type Worker interface {
        GetState() state
        ID() uint32
        StartTime() time.Time
        CompletedTime() time.Time
        CreatedTime() time.Time
        RunTime() time.Duration
        Run(done chan delegator.Directive)
        Announce()
    }

 ```

### id
A unique ID is used to reference a job throughout the system. 
- Unique ID is generated with the uuid packages, This should resolve the need for these ID's to be unique, without any repitition. 
### Time stamps
Several timestamps are collected
   - start time marks the start of the job's task when first executed. 
   - Created Time marks the time the job object is generated, this marks the very beginning of the job's journey throughout the system. 
   - Completed Time is the time the job's task completed it's operation. 

   We can derive RunTime and RoundTrip time from this information, the time it took to run the operation, and the how long the job has been traversing our system up until it gets to a specific point. 
## task
`task` is encapsulated within a `job`, and contains a reference to an `operation`, and a `context`. It executed that operation using the `context` to pass contextual variables so that the operation has all of the information it needs to run. 

``` go
type task struct {
	op      operation
	done    bool
	payload []byte
	e       error
	ctx     context.Context
}

type (
	// The task Signature expected by the worker.
	TaskSignature *func(context.Context) ([]byte, error)
	operation     interface { // The task interface.
		Func() TaskSignature
		String() string
	}
)
```
### operation
References an `operation` interface that consites of `.Func()`, and a `.String()` methods.  
- `.Func()` should return a function that conforms to the `TaskSignature` defined as:
```
*func(context.Context) ([]byte, error)
```
An operation's `TaskSignature` should take in context, and return a payload of `[]byte`, and an `error`. 

- `.String()` Will return a stringified name for the operation run.

### done
 Set upon the completion of a `task's` `operation`. 
### payload
 Contains the payload returned by a `task's` `operation`
### error
 Contains the `error` returned by the `task's` `operation`
### context
 `context` uses Golangs context library to pass arbitrary variables, and set a deadline for operations referenced within the task. 

 - variables are stored within a `context` as a `value` in a key/value map. You can set a value by creating a new context with the `.WithValue` operation, referencing the previous context targeted as a `parent context`, and you can reference that value with the `.Value('key')` operation.