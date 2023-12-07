package delegator

import ()

type worker interface {
	Run() chan any
}

type directive interface {
	Action() string // returns the action to be performed.
	Target() uint32 // returns the target of the action.
}


type Delegator struct {
	work    chan *worker // channel to receive jobs on.
	admin 	chan *directive // channel to receive directives on.
	workers map[uint32]*worker // A map of all the workers.
}