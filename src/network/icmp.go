package network

const IcmpHeaderSize = 8 // byte

type IcmpPacket struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	Content  uint32
}
