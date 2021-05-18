package error

import "fmt"

const (
	OK int = iota
	AlreadyOpened
	CantConvert
	CantCreate
	CantOpen
	CantRead
	CantRegister
	Failed
	Interrupted
	InvalidData
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
