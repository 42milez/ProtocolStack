package error

const (
	ArpIncomplete                 E = "ARP_INCOMPLETE"
	CantCloseIOResource           E = "CANT_CLOSE_IO_RESOURCE"
	CantCreateEndpoint            E = "CANT_CREATE_ENDPOINT"
	CantCreateEpollInstance       E = "CANT_CREATE_EPOLL_INSTANCE"
	CantModifyIOResourceParameter E = "CANT_MODIFY_IO_RESOURCE_PARAMETER"
	CantOpenIOResource            E = "CANT_OPEN_IO_RESOURCE"
	ChecksumMismatch              E = "CHECKSUM_MISMATCH"
	DeviceNotOpened               E = "DEVICE_NOT_OPENED"
	Error                         E = "ERROR"
	Exist                         E = "EXIST"
	InterfaceNotFound             E = "INTERFACE_NOT_FOUND"
	Interrupted                   E = "INTERRUPTED"
	InvalidPacket                 E = "INVALID_PACKET"
	InvalidProtocolVersion        E = "INVALID_PROTOCOL_VERSION"
	NetworkAddressNotMatch        E = "NETWORK_ADDRESS_NOT_MATCH"
	NoDataToRead                  E = "NO_DATA_TO_READ"
	NotFound                      E = "NOT_FOUND"
	OK                            E = "OK"
	PacketTooLong                 E = "PACKET_TOO_LONG"
	ReadFromBufError              E = "READ_FROM_BUFFER_ERROR"
	RouteNotFound                 E = "ROUTE_NOT_FOUND"
	SyscallError                  E = "SYSTEM_CALL_ERROR"
	TtlExpired                    E = "TTL_EXPIRED"
	UnsupportedProtocol           E = "UNSUPPORTED_PROTOCOL"
	WriteToBufError               E = "WRITE_TO_BUFFER_ERROR"
)

type E string

func (v E) Error() string {
	return string(v)
}
