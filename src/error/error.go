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
	InvalidHeader
	InvalidPacket
	NoDataToRead
	NotFound
)

type Error struct {
	Code int
	Msg  string
}

func (p *Error) Error() string {
	return fmt.Sprintf("%v (code: %v)", p.Msg, p.Code)
}
