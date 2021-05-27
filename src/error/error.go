package error

const (
	OK E = iota
	AlreadyOpened
	CantConvert
	CantInitialize
	CantOpen
	CantProcess
	CantRead
	CantRegister
	CantSend
	CantWriteToBuffer
	CantWriteToFile
	Error
	Failed
	InterfaceNotFound
	Interrupted
	InvalidHeader
	InvalidPacket
	NoDataToRead
	NotFound
	Terminated
)

type E int

//
//func (v E) Error() string {
//	switch v {
//	case OK:
//		return "OK"
//	case AlreadyOpened:
//		return "ALREADY_OPENED"
//	case CantConvert:
//		return "CANT_CONVERT"
//	case CantInitialize:
//		return "CANT_INITIALIZE"
//	case CantProcess:
//		return "CANT_PROCESS"
//	case CantRead:
//		return "CANT_READ"
//	case CantSend:
//		return "CANT_SEND"
//	case CantWriteToBuffer:
//		return "CANT_WRITE_TO_BUFFER"
//	case CantWriteToFile:
//		return "CANT_WRITE_TO_FILE"
//	case Error:
//		return "ERROR"
//	case Failed:
//		return "FAILED"
//	case Interrupted:
//		return "INTERRUPTED"
//	case InvalidHeader:
//		return "INVALID_HEADER"
//	case InvalidPacket:
//		return "INVALID_PACKET"
//	case NoDataToRead:
//		return "NO_DATA_TO_READ"
//	case NotFound:
//		return "NOT_FOUND"
//	case Terminated:
//		return "TERMINATED"
//	default:
//		return "UNKNOWN"
//	}
//}
