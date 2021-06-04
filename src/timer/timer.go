package timer

const (
	Stop Signal = iota
)
const (
	Any State = iota
	Stopped
	Running
)

type Condition struct {
	CurrentState State
	DesiredState State
}

type Signal int
type State int
