package binary

import (
	goBinary "encoding/binary"
	"unsafe"
)

var Endian goBinary.ByteOrder

func Swap16(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}

func Swap32(v uint32) uint32 {
	return (v&0xff000000)>>24 | (v&0x00ff0000)>>8 | (v&0x0000ff00)<<8 | (v&0x000000ff)<<24
}

func byteOrder() goBinary.ByteOrder {
	x := 0x0100
	p := unsafe.Pointer(&x)
	if 0x01 == *(*byte)(p) {
		return goBinary.LittleEndian
	} else {
		return goBinary.BigEndian
	}
}

func init() {
	Endian = byteOrder()
}
