package worker

const (
	Any State = iota
	Running
	Stopped
)

type State int

type Message struct {
	Current State
	Desired State
}
