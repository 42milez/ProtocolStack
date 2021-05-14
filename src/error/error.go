package error

import "fmt"

const (
	OK int = iota
	AlreadyOpened
	CantCreate
	CantOpen
	CantRead
	CantRegister
	Failed
	Interrupted
	InvalidHeader
	IoctlFailed
	NoDataToRead
)

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v (code: %v)", e.Msg, e.Code)
}
