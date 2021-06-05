package worker

const (
	Any State = iota
	Error
	Running
	Stopped
)

type State int

type Message struct {
	Current State
	Desired State
}
