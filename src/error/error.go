package error

const (
	OK E = iota
	CantCreateEndpoint
	CantCreateEpollInstance
	CantModifyIOResourceParameter
	CantOpenIOResource
	CantRead
	InterfaceNotFound
	Interrupted
	InvalidHeader
	InvalidPacket
	NoDataToRead
	NotFound
	Terminated
	Error
)

type E int

func (v E) Error() string {
	switch v {
	case OK:
		return "OK"
	case CantCreateEndpoint:
		return "CANT_CREATE_ENDPOINT"
	case CantCreateEpollInstance:
		return "CANT_CREATE_EPOLL_INSTANCE"
	case CantModifyIOResourceParameter:
		return "CANT_MODIFY_IO_RESOURCE_PARAMETER"
	case CantOpenIOResource:
		return "CANT_OPEN_IO_RESOURCE"
	case CantRead:
		return "CANT_READ"
	case InterfaceNotFound:
		return "INTERFACE_NOT_FOUND"
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
	case Terminated:
		return "TERMINATED"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
