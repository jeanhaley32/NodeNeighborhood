## Synopsis
`directive` is a data structure used by the `work delegator` service to represent an action. It comprises two fields: `action`, which defines the action type, and `target`, an optional uint32 identifier for the action’s target. The `directive` is communicated via the `admin` channel, allowing the `delegator` to instruct running jobs and control task execution.

```Go
type directive struct{
    action string // action to be performed
    target uint32 // target of action, if any. uses uint32 unique identifier. 
}

```
### Defined Directives:
- done: Signals the completion of a target worker
