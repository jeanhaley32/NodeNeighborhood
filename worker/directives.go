package github.com/jeanhaley32/nodeneighborhood/worker

type action int64

const (
	done action = iota
)

func (a action) String() string {
	switch a {
	case done:
		return "done"
	}
	return ""
}

type Directive struct {
	a      action
	target uint32
}

func (d *Directive) Action() string {
	return d.a.String()
}

func (d *Directive) Target() uint32 {
	return d.target
}

// Creates a New Done directive.
// The target is the completed job's id.
// If sent to the delegator, the delegator will remove the job from the
func NewDoneDirective(target uint32) Directive {
	return Directive{
		a:      done,
		target: target,
	}
}
