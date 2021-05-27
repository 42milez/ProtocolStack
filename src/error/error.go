package error

const (
	OK int = iota
	AlreadyOpened
	CantConvert
	CantOpen
	CantProcess
	CantRead
	CantRegister
	CantSend
	CantWriteToBuffer
	CantWriteToFile
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
	switch p.Code {
	case OK:
		return "OK"
	case AlreadyOpened:
		return "ALREADY_OPENED"
	case CantConvert:
		return "CANT_CONVERT"
	case CantProcess:
		return "CANT_PROCESS"
	case CantRead:
		return "CANT_READ"
	case CantSend:
		return "CANT_SEND"
	case CantWriteToBuffer:
		return "CANT_WRITE_TO_BUFFER"
	case CantWriteToFile:
		return "CANT_WRITE_TO_FILE"
	case Failed:
		return "FAILED"
	case Interrupted:
		return "INTERRUPTED"
	case InvalidHeader:
		return "INVALID_HEADER"
	case InvalidPacket:
		return "INVALID_PACKET"
	case NoDataToRead:
		return "NO_DATA_TO_READ"
	case NotFound:
		return "NOT_FOUND"
	default:
		return "UNKNOWN"
	}
}
