package error

import "fmt"

const (
	OK int = iota	// 0
	AlreadyOpened
	CantCreate
	CantOpen
	CantRead
	CantRegister	// 5
	Failed
	Interrupted
	InvalidData
	InvalidHeader
	IoctlFailed		// 10
	NoDataToRead
)

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v (code: %v)", e.Msg, e.Code)
}
