# Worker Notes:

```
type job struct {
	id      int             // Unique id of the job
	task    func() payload  // The task to be performed
	ctx     context.Context // The context of the task. Used to set the initial values of the task.
	p       payload         // The returned value of the task.
	metrics struct {        // Metrics of the task.
		created   time.Time     // The time the task was created.
		start     time.Time     // start time of task. Used to calculate the time taken to complete the task.
		taskTime  time.Duration // The Time taken to complete the task.
		completed time.Time     // The time the task was completed.
		complete  bool          // Whether the task is complete or not.
	}
	chBundle struct { // Bundle of channels used for communication.
		parent         workch // The parent channel.
		localComms     workch // The local channel.
		targetIngester workch // The target channel.
	}
```

## Elemenents of the job struct
- id - unique ID representing this item of work. 128 bit UUID [rfc41222](https://datatracker.ietf.org/doc/html/rfc4122)
  - generating this ID using the [google/uuid](https://www.github.com/google/uuid) package. This will generate a unique id with a mathematical improbabilty of repeating.
  > "A UUID is 128 bits long, and can guarantee
   uniqueness across space and time." - [rfc41222](https://datatracker.ietf.org/doc/html/rfc4122)
  - Since this is mathematically trusted to provide a unique ID, i'm going to trust it to not cause duplicate work ids. I am open to this being a potential needle in a haystack problem in the future. In which case we'll need to implement a system of indexing all creating work ids, and ensuring we have no duplicates. 
  - an alternative method, and one potentially "simpler" and less computationally intensive, is to create a routine that return an ever increasing integer whenever it's pinged. This routine increments upwards everytime it's called, and in turn generates a number in sequence that would ensure uniqueness at least for the instance in which the service is running. This would not be a good system to use if we planned on archiving work, as every restart of the server would return that int to zero. 
  - If the promisE of UUIDv4 is accurate, we should never have to worry about generating the same number twice. As a highschool dropout, the mathematical reality of this boggles my mind, and I can only imagine the probability is low enough, the chances of two numbers existing at the same time akin to the chances of an errant gamma ray burst ripping apart the earth's atmosphere, theoretically possible, but so unlikely that it's not worth considering. 
- task - function that represents the task to be completed. 
  - This is the task is the action that the worker generated from the work object performs. 
  - It pulls variables from the job structs local "context"
  - after processing it's work, it returns a boolean indicating success, and a payload. 
  - the payload is an array of bytes, that is then sent to the ingester for ingestion. 

- success - bool representing the success for failure of the jobs primary task. This is set based on the return value of the onboard task.
- ctx - ctx is a context element used to store variables that the job requires to complete the task provided. 
  - Using the golang native '[context](https://pkg.go.dev/context)' library to achieve this.
> "Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes." - [context](https://pkg.go.dev/context)