package derectives

type action int64

const (
	// define the actions.
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
	// define the directive struct.
	a      action
	target uint32
}

func (d *Directive) Action() string {
	return d.a.String()
}

func (d *Directive) Target() uint32 {
	return d.target
}

func NewDoneDirective(target uint32) *Directive {
	return &Directive{
		a:      done,
		target: target,
	}
}
