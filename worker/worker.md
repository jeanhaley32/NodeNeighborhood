```go
type job struct {
	id      uint32                                   // Unique id of the job
	task    func(context.Context) (success, payload) // The task to be performed
	success success                                  // Whether the task was successful or not.
	vars    map[string]any                           // Variables used by the task.
	p       payload                                  // The returned value of the task.
	complete  bool          // Whether the task is complete or not.
	metrics struct {                                 // Metrics of the task.
		created   time.Time     // The time the task was created.
		start     time.Time     // start time of task. Used to calculate the time taken to complete the task.
		completed time.Time     // The time the task was completed.
	}
	chBundle struct { // Bundle of channels used for communication.
		parent         workch // The parent channel.
		localComms     workch // The local channel.
		targetIngester workch // The target channel.
	}
}

```

Job Struct contains the necessary task and components to create a worker, it also contains relative datapoints to be used to derive metrics. 

Work items are comprised of:

- id - unique ID representing this item of work. 128 bit UUID [rfc41222](https://datatracker.ietf.org/doc/html/rfc4122)
    - generating this ID using the [google/uuid](https://www.github.com/google/uuid) package. This will generate a unique id with a mathematical improbabilty of repeating.
    
    > "A UUID is 128 bits long, and can guarantee
    uniqueness across space and time." - rfc41222
    > 
    - Since this is mathematically trusted to provide a unique ID, i'm going to trust it to not cause duplicate work ids. I am open to this being a potential needle in a haystack problem in the future. In which case we'll need to implement a system of indexing all creating work ids, and ensuring we have no duplicates.
    - an alternative method, and one potentially "simpler" and less computationally intensive, is to create a routine that return an ever increasing integer whenever it's pinged. This routine increments upwards everytime it's called, and in turn generates a number in sequence that would ensure uniqueness at least for the instance in which the service is running. This would not be a good system to use if we planned on archiving work, as every restart of the server would return that int to zero.
    - If the promise of UUIDv4 is accurate, we should never have to worry about generating the same number twice. As a high school dropout, the mathematical reality of this boggles my mind, and I can only imagine that the probability of two duplicate ID's existing at the same time is so astronomically low, it's akin to the chances of an errant gamma ray burst ripping apart the earth's atmosphere, theoretically possible, but so unlikely that it's not worth considering.
- task - function that represents the task to be completed.
    - Action to be executed.
    - takes in context.Context, This will set a runtime deadline, and is a vessel for any other runtime specific context that may be deemed necessary.
    - It adopts variables from the local vars component.
    - after processing it's work, it returns a boolean indicating success, and a payload.
    - the payload is an array of bytes, that is then sent to the ingester for ingestion.
- success - bool representing the success for failure of the jobs primary task. This is set based on the return value of the onboard task.
- complete - bool representing the tasks state, complete to not completed.
- vars - a map of `map[string]any` , vars contains the variables used by the task.
- payload - This is the payload returned by the work objects task.
- Metrics - Useful Data-points used to determine performance metrics.
    - start
        - Start Time, so we can derive statics and see what time this piece of work was started.
        - There may be value in differentiating “work generation time” and “work execution start time.”
        - finish
        - Time that the worker finished its operation.
        - There may also be value in adding a time to indicate when work was completed and successfully ingested.
    - Created
        - Time job object was created
    - completion
        - Time task was completed
- ChBundle - Bundle of channels used for communication between routines.
    - localComms - Channel the worker will listen on.
    - parent - Channel of the parent process that spawned the worker
    - targetIngester - Target channel to send payload.



# Methods

## execute()
 - non-public
 - Executes the job's task.
 - populates job payload. 
 - Sets state information. 
   - Task start time
   - Task Completion Time. 
   - Task runtime
   - Task Completion Bool
   - Task success Bool
## Run()
