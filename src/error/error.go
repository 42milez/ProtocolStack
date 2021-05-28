package error

const (
	OK                            E = "OK"
	CantCreateEndpoint            E = "CANT_CREATE_ENDPOINT"
	CantCreateEpollInstance       E = "CANT_CREATE_EPOLL_INSTANCE"
	CantModifyIOResourceParameter E = "CANT_MODIFY_IO_RESOURCE_PARAMETER"
	CantOpenIOResource            E = "CANT_OPEN_IO_RESOURCE"
	InterfaceNotFound             E = "INTERFACE_NOT_FOUND"
	Interrupted                   E = "INTERRUPTED"
	InvalidPacket                 E = "INVALID_PACKET"
	NoDataToRead                  E = "NO_DATA_TO_READ"
	NotFound                      E = "NOT_FOUND"
	Terminated                    E = "TERMINATED"
	Error                         E = "ERROR"
)

type E string

func (v E) Error() string {
	return string(v)
}
