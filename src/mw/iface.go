package mw

const LoopbackIpAddr = "127.0.0.1"
const LoopbackBroadcast = "127.255.255.255"
const LoopbackNetmask = "255.0.0.0"
const LoopbackNetwork = "127.0.0.0"

// An Iface is a single iface.
type Iface struct {
	Family    AddrFamily
	Unicast   IP
	Netmask   IP
	Broadcast IP
	Dev       IDevice
}
