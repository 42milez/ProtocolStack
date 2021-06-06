package worker

const (
	Any State = iota
	Error
	Running
	Stopped
)

type State uint32

type Message struct {
	ID      uint32
	Current State
	Desired State
}
