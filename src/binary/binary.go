package binary

import (
	goBinary "encoding/binary"
	"unsafe"
)

var Endian goBinary.ByteOrder

func Swap16(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}

func byteOrder() goBinary.ByteOrder {
	x := 0x0100
	p := unsafe.Pointer(&x)
	if 0x01 == *(*byte)(p) {
		return goBinary.BigEndian
	} else {
		return goBinary.LittleEndian
	}
}

func init() {
	Endian = byteOrder()
}
