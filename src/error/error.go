package error

type Error int

const(
	OK Error = iota
	AlreadyOpened
	CantOpen
	CantRead
	CantRegister
	Fatal
	IoctlFailed
	NoDataToRead
)

func (e Error) Error() string {
	switch e {
	case OK:
		return "OK"
	case AlreadyOpened:
		return "ALREADY_OPENED"
	case CantOpen:
		return "CANT_OPEN"
	case CantRead:
		return "CANT_READ"
	case CantRegister:
		return "CANT_REGISTER"
	case Fatal:
		return "FATAL"
	case IoctlFailed:
		return "IOCTL_FAILED"
	case NoDataToRead:
		return "NO_DATA_TO_READ"
	default:
		return "UNKNOWN_ERROR"
	}
}
