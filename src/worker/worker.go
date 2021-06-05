package worker

const (
	Any State = iota
	Error
	Running
	Stopped
)

type ID int
type State int

type Message struct {
	ID      ID
	Current State
	Desired State
}
