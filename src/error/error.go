package error

const (
	AlreadyBound Err = iota
	ArpIncomplete
	BacklogFull
	CantAllocatePcb
	CantCloseIOResource
	CantCreateEndpoint
	CantCreateEpollInstance
	CantModifyIOResourceParameter
	CantOpenIOResource
	ChecksumMismatch
	DeviceNotOpened
	Error
	Exist
	InterfaceNotFound
	Interrupted
	InvalidPacket
	InvalidPacketLength
	InvalidPcbState
	InvalidProtocolVersion
	NeedRetry
	NetworkAddressNotMatch
	NoDataToRead
	NotFound
	OK
	PacketTooLong
	PcbNotFound
	ReadFromBufError
	RouteNotFound
	SyscallError
	TtlExpired
	UnsupportedProtocol
	WriteToBufError
)

var errors = map[Err]string{
	AlreadyBound:                  "ALREADY_BOUND",
	ArpIncomplete:                 "ARP_INCOMPLETE",
	BacklogFull:                   "BACKLOG_FULL",
	CantAllocatePcb:               "CANT_ALLOCATE_PCB",
	CantCloseIOResource:           "CANT_CLOSE_IO_RESOURCE",
	CantCreateEndpoint:            "CANT_CREATE_ENDPOINT",
	CantCreateEpollInstance:       "CANT_CREATE_EPOLL_INSTANCE",
	CantModifyIOResourceParameter: "CANT_MODIFY_IO_RESOURCE_PARAMETER",
	CantOpenIOResource:            "CANT_OPEN_IO_RESOURCE",
	ChecksumMismatch:              "CHECKSUM_MISMATCH",
	DeviceNotOpened:               "DEVICE_NOT_OPENED",
	Error:                         "ERROR",
	Exist:                         "EXIST",
	InterfaceNotFound:             "INTERFACE_NOT_FOUND",
	Interrupted:                   "INTERRUPTED",
	InvalidPacket:                 "INVALID_PACKET",
	InvalidPacketLength:           "INVALID_PACKET_LENGTH",
	InvalidPcbState:               "INVALID_PCB_STATE",
	InvalidProtocolVersion:        "INVALID_PROTOCOL_VERSION",
	NeedRetry:                     "NEED_RETRY",
	NetworkAddressNotMatch:        "NETWORK_ADDRESS_NOT_MATCH",
	NoDataToRead:                  "NO_DATA_TO_READ",
	NotFound:                      "NOT_FOUND",
	OK:                            "OK",
	PacketTooLong:                 "PACKET_TOO_LONG",
	PcbNotFound:                   "PCB_NOT_FOUND",
	ReadFromBufError:              "READ_FROM_BUFFER_ERROR",
	RouteNotFound:                 "ROUTE_NOT_FOUND",
	SyscallError:                  "SYSTEM_CALL_ERROR",
	TtlExpired:                    "TTL_EXPIRED",
	UnsupportedProtocol:           "UNSUPPORTED_PROTOCOL",
	WriteToBufError:               "WRITE_TO_BUFFER_ERROR",
}

type Err int

func (v Err) Error() string {
	return v.String()
}

func (v Err) String() string {
	return errors[v]
}
