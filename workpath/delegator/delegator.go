package delegator

type worker interface {
	Run() chan any
}

type directive interface {
	Action() string
	Target() uint32
}

type Delegator struct {
	// channel to receive jobs on.
	work chan *worker
	// channel to receive directives on.
	admin chan *directive
	// A map of all the workers.
	workers map[uint32]*worker
}
