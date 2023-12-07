package delegator
// Delegator Processes jobs from a `work` channel, starting workers, and keeping track of them.
// Delegator will also handle andy communications with running wokers.
// Actions are objects received on the `admin` channel.
// Jobs are received on the `work` channel. 
import ()

type worker interface {
	// define the worker interface.	
}

type Delegator struct {
	workers map[uint32]*Worker // A map of all the workers.
}