
## Synopsis 
The work delegator service is designed to efficiently manage and delegate tasks in a concurrent environment. It operates as a goroutine and utilizes two channels: `work` and `admin`.

The `work` channel is responsible for receiving `Job` objects, which represent tasks to be executed. Upon receiving a job, the `delegator` service converts it into a `goroutine` and adds it to a list of running `work`. This allows multiple tasks to be executed concurrently.

The `admin` channel is used for communication between the running jobs and the delegator service. When a job finishes its work, it sends a signal through the admin channel to notify the delegator. This signal prompts the delegator to remove the completed job from its list of running jobs.

By utilizing these two channels, the `work delegator` service effectively manages and tracks the execution of `tasks`, ensuring efficient task delegation and proper synchronization between the delegator and the running jobs.

## Module Defined Signatures
```Go
type Delegator struct{
    work map[uin32]*Worker
}

type worker interface{

}

type directive struct{
    action string // action to be performed
    target uint32 // target of action, if any. uses uint32 unique identifier. 
}

```

## Glossary
### Directive:
`directive` is a data structure used by the `work delegator` service to represent an action. It comprises two fields: `action`, which defines the action type, and `target`, an optional uint32 identifier for the action’s target. The `directive` is communicated via the `admin` channel, allowing the `delegator` to instruct running jobs and control task execution.

Defined Directives:
 - done: Signals the completion of a target worker. 