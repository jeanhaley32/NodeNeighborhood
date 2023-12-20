    WARNING: changes made to this document, should reflect in the top level workpath readme.md. ../readme.md

## Workers

A running `worker` is defined as a goroutine executed and tracked by the `delegator`, running a `job` embedded with a `task`. 

The `job` and `task` object bundle contains the`operation`for the worker to  execute. We reference an `operation`, and embed the necessary `context` to accomplish that operation. 
### Job
 A Job is the Top level encapsulation of a task, it contains a unique `id`, time stamps for the job's `creation`, `start`, and `completion` times. and it contains the `task` to be executed. 
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

### `ID`
A `unique ID` is used to reference a `job` throughout it's journey within system. 
- `Unique ID` is generated with the `uuid` package. This is a well used package that should solve the need for each `ID` to be independently unique, without the need for inter-commmunication and varification of that uniqueness. 
### `Time stamps`
Several `time stamps` are collected
   - `start time` marks the start of the job's task when first executed. 
   - `Created Time` marks the time the `job` object was generated, this is the very beginning of the `job's` journey throughout the system. 
   - `Completed Time` is the time the job's `task`completed it's operation. 

   We can derive the time it took to run the operation `Run Time`, and the how long the job has been traversing our system up until it gets to a specific point `Round Trip Time`. 
### Task
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
### `operation`
References an `operation` interface that consites of `.Func()`, and a `.String()` methods.  
- `.Func()` should return a function that conforms to the `TaskSignature` defined as:
```
*func(context.Context) ([]byte, error)
```
An operation's `TaskSignature` should take in context, and return a payload of `[]byte`, and an `error`. 

- `.String()` Will return a stringified name for the operation run.

### `done`
 Set upon the completion of a `task's` `operation`. 
### `payload`
 Contains the payload returned by a `task's` `operation`
### `error`
 Contains the `error` returned by the `task's` `operation`
### `context`
 `context` uses Golangs context library to pass arbitrary variables, and set a deadline for operations referenced within the task. 

 - variables are stored within a `context` as a `value` in a key/value map. You can set a value by creating a new context with the `.WithValue` operation, referencing the previous context targeted as a `parent context`, and you can reference that value with the `.Value('key')` operation.
