package error

const (
	OK                            E = "OK"
	CantCreateEndpoint            E = "CANT_CREATE_ENDPOINT"
	CantCreateEpollInstance       E = "CANT_CREATE_EPOLL_INSTANCE"
	CantModifyIOResourceParameter E = "CANT_MODIFY_IO_RESOURCE_PARAMETER"
	CantOpenIOResource            E = "CANT_OPEN_IO_RESOURCE"
	ChecksumMismatch              E = "CHECKSUM_MISMATCH"
	Exist                         E = "EXIST"
	InterfaceNotFound             E = "INTERFACE_NOT_FOUND"
	Interrupted                   E = "INTERRUPTED"
	InvalidPacket                 E = "INVALID_PACKET"
	InvalidProtocolVersion        E = "INVALID_PROTOCOL_VERSION"
	IpPacketTooLong               E = "IP_PACKET_TOO_LONG"
	NetworkAddressNotMatch        E = "NETWORK_ADDRESS_NOT_MATCH"
	NoDataToRead                  E = "NO_DATA_TO_READ"
	NotFound                      E = "NOT_FOUND"
	RouteNotFound                 E = "ROUTE_NOT_FOUND"
	Terminated                    E = "TERMINATED"
	TtlExpired                    E = "TTL_EXPIRED"
	UnsupportedProtocol           E = "UNSUPPORTED_PROTOCOL"
	Error                         E = "ERROR"
)

type E string

func (v E) Error() string {
	return string(v)
}
