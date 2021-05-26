package binary

import "unsafe"

const BigEndian = 4321
const LittleEndian = 1234

func ByteOrder() int {
	x := 0x0100
	p := unsafe.Pointer(&x)
	if 0x01 == *(*byte)(p) {
		return BigEndian
	} else {
		return LittleEndian
	}
}

func Swap16(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}
