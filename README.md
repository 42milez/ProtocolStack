## ProtocolStack
[![CI](https://github.com/42milez/ProtocolStack/actions/workflows/ci.yml/badge.svg)](https://github.com/42milez/ProtocolStack/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/42milez/ProtocolStack/branch/main/graph/badge.svg?token=ALHDIWP6KH)](https://codecov.io/gh/42milez/ProtocolStack)  
This repository is for learning network programming in Go.

## Requirements
- OS: Linux
- Go: 1.14.x or higher

## Instructions
ProtocolStack needs a TAP device for its capability. For the reason, this project uses a virtual machine (Linux) to run the application. Also, you can perform debugging remotely with the VM if you use GoLand.

### 1. Prepare virtual machine
Software required is as follows: 

- [Mutagen](https://github.com/mutagen-io/mutagen)
- [Vagrant](https://www.vagrantup.com/)
- [VirtualBox](https://www.virtualbox.org/)

Note: Mutagen synchronizes files between local system and virtual machine.

### 2. Control virtual machine
`vm.sh` controls the virtual machine. The available commands are as follows:

- `start` start the virtual machine
- `stop` stop the virtual machine
- `restart` restart the virtual machine

You can perform the commands below:

```shell
> ./vm.sh start   # Start a VM and create a Mutagen session.
> ./vm.sh stop    # Stop the VM and terminate the Mutagen session.
> ./vm.sh restart # Restart the VM and recreate the Mutagen session.
```

### 3. Remote debugging with GoLand
See the [instruction](https://github.com/42milez/ProtocolStack/wiki/Remote-Debugging-with-GoLand) for more detail.

Related information:
- [What Are Run Targets & How To Run Code Anywhere](https://blog.jetbrains.com/go/2021/04/29/what-are-run-targets-and-how-to-run-code-anywhere/)
- [How to use Docker to compile and run Go code from GoLand](https://blog.jetbrains.com/go/2021/04/30/how-to-use-docker-to-compile-go-from-goland/)
- [Compile and run Go code using WSL 2 and GoLand](https://blog.jetbrains.com/go/2021/05/05/compile-and-run-go-code-using-wsl-2-and-goland/)

### 4. Compile
It is able to compile the application with `make` as below:
```shell
> make compile
```

Note: `make` supports the commands below:
- `build` build project
- `clean` clean up caches
- `compile` clean up caches, resolve dependencies, and build the program
- `fmt` run formatter
- `gen` generate source code
- `lint` run linters (golangci-lint)
- `resolve` resolve dependencies
- `test` run all tests

## Supported Protocols

- [x] Ethernet
- [x] ARP
- [x] IP
  - [x] v4
  - [ ] v6
- [x] ICMP
- [ ] TCP
- [ ] UDP

## Examples

### ICMP Echo Reply
1. Reply to ARP request
2. Reply to ICMP echo request

<details>
<summary>Log</summary>

```text
[vagrant@ps ~]$ ./app/bin/tcp_server
[I] 2021/06/01 13:49:33 --------------------------------------------------
[I] 2021/06/01 13:49:33  INITIALIZE DEVICES
[I] 2021/06/01 13:49:33 --------------------------------------------------
[I] 2021/06/01 13:49:33 ▶ Device was registered
[I] 2021/06/01 13:49:33 	type:      LOOPBACK
[I] 2021/06/01 13:49:33 	name:      net0 ()
[I] 2021/06/01 13:49:33 	addr:      00:00:00:00:00:00
[I] 2021/06/01 13:49:33 	broadcast: 00:00:00:00:00:00
[I] 2021/06/01 13:49:33 	peer:      00:00:00:00:00:00
[I] 2021/06/01 13:49:33 ▶ Interface was attached
[I] 2021/06/01 13:49:33 	ip:     127.0.0.1
[I] 2021/06/01 13:49:33 	device: net0 ()
[I] 2021/06/01 13:49:33 ▶ Route was registered
[I] 2021/06/01 13:49:33 	network:  127.0.0.0
[I] 2021/06/01 13:49:33 	netmask:  255.0.0.0
[I] 2021/06/01 13:49:33 	unicast:  127.0.0.1
[I] 2021/06/01 13:49:33 	next hop: 0.0.0.0
[I] 2021/06/01 13:49:33 	device:   net0 ()
[I] 2021/06/01 13:49:33 ▶ Device was registered
[I] 2021/06/01 13:49:33 	type:      ETHERNET
[I] 2021/06/01 13:49:33 	name:      net1 (tap0)
[I] 2021/06/01 13:49:33 	addr:      0b:16:21:2c:37:42
[I] 2021/06/01 13:49:33 	broadcast: ff:ff:ff:ff:ff:ff
[I] 2021/06/01 13:49:33 	peer:      00:00:00:00:00:00
[I] 2021/06/01 13:49:33 ▶ Interface was attached
[I] 2021/06/01 13:49:33 	ip:     192.0.2.2
[I] 2021/06/01 13:49:33 	device: net1 (tap0)
[I] 2021/06/01 13:49:33 ▶ Route was registered
[I] 2021/06/01 13:49:33 	network:  192.0.0.0
[I] 2021/06/01 13:49:33 	netmask:  255.255.255.0
[I] 2021/06/01 13:49:33 	unicast:  192.0.2.2
[I] 2021/06/01 13:49:33 	next hop: 0.0.0.0
[I] 2021/06/01 13:49:33 	device:   net1 (tap0)
[I] 2021/06/01 13:49:33 ▶ Default gateway was registered
[I] 2021/06/01 13:49:33 	network:  0.0.0.0
[I] 2021/06/01 13:49:33 	netmask:  0.0.0.0
[I] 2021/06/01 13:49:33 	unicast:  192.0.2.2
[I] 2021/06/01 13:49:33 	next hop: 192.0.2.1
[I] 2021/06/01 13:49:33 	device:   net1 (tap0)
[I] 2021/06/01 13:49:33 --------------------------------------------------
[I] 2021/06/01 13:49:33  START SERVER
[I] 2021/06/01 13:49:33 --------------------------------------------------
[I] 2021/06/01 13:49:33 ▶ Device was opened
[I] 2021/06/01 13:49:33 	type: LOOPBACK
[I] 2021/06/01 13:49:33 	name: net0 ()
[I] 2021/06/01 13:49:33 ▶ Device was opened
[I] 2021/06/01 13:49:33 	type: ETHERNET
[I] 2021/06/01 13:49:33 	name: net1 (tap0)
[I] 2021/06/01 13:49:33 ▶ Event occurred
[I] 2021/06/01 13:49:33 	events: 1
[I] 2021/06/01 13:49:33 	device: net1 (tap0)
[I] 2021/06/01 13:49:33 ▶ Ethernet frame was received: 54 bytes
[I] 2021/06/01 13:49:33 ▶ Event occurred
[I] 2021/06/01 13:49:33 	events: 1
[I] 2021/06/01 13:49:33 	device: net1 (tap0)
[I] 2021/06/01 13:49:33 ▶ Ethernet frame was received: 54 bytes
[I] 2021/06/01 13:49:34 ▶ Event occurred
[I] 2021/06/01 13:49:34 	events: 1
[I] 2021/06/01 13:49:34 	device: net1 (tap0)
[I] 2021/06/01 13:49:34 ▶ Ethernet frame was received: 42 bytes
[I] 2021/06/01 13:49:34 ▶ Incoming ethernet frame
[I] 2021/06/01 13:49:34 	dst:  ff:ff:ff:ff:ff:ff
[I] 2021/06/01 13:49:34 	src:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:34 	type: 0x0806 (ARP)
[I] 2021/06/01 13:49:34 ▶ Incoming ARP packet
[I] 2021/06/01 13:49:34 	hardware type:           Ethernet (10Mb)
[I] 2021/06/01 13:49:34 	protocol Type:           IPv4
[I] 2021/06/01 13:49:34 	hardware address length: 6
[I] 2021/06/01 13:49:34 	protocol address length: 4
[I] 2021/06/01 13:49:34 	opcode:                  REQUEST (1)
[I] 2021/06/01 13:49:34 	sender hardware address: e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:34 	sender protocol address: 192.0.2.1
[I] 2021/06/01 13:49:34 	target hardware address: 00:00:00:00:00:00
[I] 2021/06/01 13:49:34 	target protocol address: 192.0.2.2
[I] 2021/06/01 13:49:34 ▶ Outgoing ARP packet (REPLY):
[I] 2021/06/01 13:49:34 	hardware type:           Ethernet (10Mb)
[I] 2021/06/01 13:49:34 	protocol Type:           IPv4
[I] 2021/06/01 13:49:34 	hardware address length: 6
[I] 2021/06/01 13:49:34 	protocol address length: 4
[I] 2021/06/01 13:49:34 	opcode:                  REPLY (2)
[I] 2021/06/01 13:49:34 	sender hardware address: 01:00:e2:20:76:2d
[I] 2021/06/01 13:49:34 	sender protocol address: 192.0.2.2
[I] 2021/06/01 13:49:34 	target hardware address: e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:34 	target protocol address: 192.0.2.1
[I] 2021/06/01 13:49:34 ▶ Outgoing Ethernet frame
[I] 2021/06/01 13:49:34 	dst:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:34 	src:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:34 	type: 0x0806 (ARP)
[I] 2021/06/01 13:49:34 	payload: 00 01 08 00 06 04 00 02 01 00 e2 20 76 2d c0 00 02 02 e2 20
[I] 2021/06/01 13:49:34 ▶ Ethernet frame was sent: 60 bytes (payload: 28 bytes)
[I] 2021/06/01 13:49:34 ▶ Event occurred
[I] 2021/06/01 13:49:34 	events: 1
[I] 2021/06/01 13:49:34 	device: net1 (tap0)
[I] 2021/06/01 13:49:34 ▶ Ethernet frame was received: 98 bytes
[I] 2021/06/01 13:49:34 ▶ Incoming ethernet frame
[I] 2021/06/01 13:49:34 	dst:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:34 	src:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:34 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:34 ▶ Incoming IP packet
[I] 2021/06/01 13:49:34 	version:             IPv4
[I] 2021/06/01 13:49:34 	ihl:                 5
[I] 2021/06/01 13:49:34 	type of service:     0b00000000
[I] 2021/06/01 13:49:34 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:34 	id:                  8328
[I] 2021/06/01 13:49:34 	flags:               0b010
[I] 2021/06/01 13:49:34 	fragment offset:     0
[I] 2021/06/01 13:49:34 	ttl:                 64
[I] 2021/06/01 13:49:34 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:34 	checksum:            0x0000
[I] 2021/06/01 13:49:34 	source address:      192.0.2.1
[I] 2021/06/01 13:49:34 	destination address: 192.0.2.2
[I] 2021/06/01 13:49:34 ▶ Incoming ICMP packet
[I] 2021/06/01 13:49:34 	type:     8 (Echo)
[I] 2021/06/01 13:49:34 	code:     0
[I] 2021/06/01 13:49:34 	checksum: 0x85f5
[I] 2021/06/01 13:49:34 	id:       48
[I] 2021/06/01 13:49:34 	seq:      1
[I] 2021/06/01 13:49:34 ▶ Outgoing ICMP packet
[I] 2021/06/01 13:49:34 	type:     0 (Echo Reply)
[I] 2021/06/01 13:49:34 	code:     0
[I] 2021/06/01 13:49:34 	checksum: 0x8df5
[I] 2021/06/01 13:49:34 ▶ Outgoing IP packet
[I] 2021/06/01 13:49:34 	version:             IPv4
[I] 2021/06/01 13:49:34 	ihl:                 5
[I] 2021/06/01 13:49:34 	type of service:     0b00000000
[I] 2021/06/01 13:49:34 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:34 	id:                  0
[I] 2021/06/01 13:49:34 	flags:               0b000
[I] 2021/06/01 13:49:34 	fragment offset:     0
[I] 2021/06/01 13:49:34 	ttl:                 255
[I] 2021/06/01 13:49:34 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:34 	checksum:            0x37a5
[I] 2021/06/01 13:49:34 	source address:      192.0.2.2
[I] 2021/06/01 13:49:34 	destination address: 192.0.2.1
[I] 2021/06/01 13:49:34 ▶ Outgoing Ethernet frame
[I] 2021/06/01 13:49:34 	dst:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:34 	src:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:34 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:34 	payload: 45 00 00 54 00 00 00 00 ff 01 37 a5 c0 00 02 02 c0 00 02 01
[I] 2021/06/01 13:49:34 			 00 00 8d f5 00 30 00 01 ee 3a b6 60 00 00 00 00 01 6b 0d 00
[I] 2021/06/01 13:49:34 			 00 00 00 00 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f
[I] 2021/06/01 13:49:34 			 20 21 22 23 24 25 26 27 28 29 2a 2b 2c 2d 2e 2f 30 31 32 33
[I] 2021/06/01 13:49:34 ▶ Ethernet frame was sent: 98 bytes (payload: 84 bytes)
[I] 2021/06/01 13:49:35 ▶ Event occurred
[I] 2021/06/01 13:49:35 	events: 1
[I] 2021/06/01 13:49:35 	device: net1 (tap0)
[I] 2021/06/01 13:49:35 ▶ Ethernet frame was received: 98 bytes
[I] 2021/06/01 13:49:35 ▶ Incoming ethernet frame
[I] 2021/06/01 13:49:35 	dst:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:35 	src:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:35 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:35 ▶ Incoming IP packet
[I] 2021/06/01 13:49:35 	version:             IPv4
[I] 2021/06/01 13:49:35 	ihl:                 5
[I] 2021/06/01 13:49:35 	type of service:     0b00000000
[I] 2021/06/01 13:49:35 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:35 	id:                  8545
[I] 2021/06/01 13:49:35 	flags:               0b010
[I] 2021/06/01 13:49:35 	fragment offset:     0
[I] 2021/06/01 13:49:35 	ttl:                 64
[I] 2021/06/01 13:49:35 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:35 	checksum:            0x0000
[I] 2021/06/01 13:49:35 	source address:      192.0.2.1
[I] 2021/06/01 13:49:35 	destination address: 192.0.2.2
[I] 2021/06/01 13:49:35 ▶ Incoming ICMP packet
[I] 2021/06/01 13:49:35 	type:     8 (Echo)
[I] 2021/06/01 13:49:35 	code:     0
[I] 2021/06/01 13:49:35 	checksum: 0xc9cc
[I] 2021/06/01 13:49:35 	id:       48
[I] 2021/06/01 13:49:35 	seq:      2
[I] 2021/06/01 13:49:35 ▶ Outgoing ICMP packet
[I] 2021/06/01 13:49:35 	type:     0 (Echo Reply)
[I] 2021/06/01 13:49:35 	code:     0
[I] 2021/06/01 13:49:35 	checksum: 0xd1cc
[I] 2021/06/01 13:49:35 ▶ Outgoing IP packet
[I] 2021/06/01 13:49:35 	version:             IPv4
[I] 2021/06/01 13:49:35 	ihl:                 5
[I] 2021/06/01 13:49:35 	type of service:     0b00000000
[I] 2021/06/01 13:49:35 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:35 	id:                  1
[I] 2021/06/01 13:49:35 	flags:               0b000
[I] 2021/06/01 13:49:35 	fragment offset:     0
[I] 2021/06/01 13:49:35 	ttl:                 255
[I] 2021/06/01 13:49:35 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:35 	checksum:            0x37a4
[I] 2021/06/01 13:49:35 	source address:      192.0.2.2
[I] 2021/06/01 13:49:35 	destination address: 192.0.2.1
[I] 2021/06/01 13:49:35 ▶ Outgoing Ethernet frame
[I] 2021/06/01 13:49:35 	dst:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:35 	src:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:35 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:35 	payload: 45 00 00 54 00 01 00 00 ff 01 37 a4 c0 00 02 02 c0 00 02 01
[I] 2021/06/01 13:49:35 			 00 00 d1 cc 00 30 00 02 ef 3a b6 60 00 00 00 00 bc 92 0d 00
[I] 2021/06/01 13:49:35 			 00 00 00 00 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f
[I] 2021/06/01 13:49:35 			 20 21 22 23 24 25 26 27 28 29 2a 2b 2c 2d 2e 2f 30 31 32 33
[I] 2021/06/01 13:49:35 ▶ Ethernet frame was sent: 98 bytes (payload: 84 bytes)
[I] 2021/06/01 13:49:36 ▶ Event occurred
[I] 2021/06/01 13:49:36 	events: 1
[I] 2021/06/01 13:49:36 	device: net1 (tap0)
[I] 2021/06/01 13:49:36 ▶ Ethernet frame was received: 98 bytes
[I] 2021/06/01 13:49:36 ▶ Incoming ethernet frame
[I] 2021/06/01 13:49:36 	dst:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:36 	src:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:36 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:36 ▶ Incoming IP packet
[I] 2021/06/01 13:49:36 	version:             IPv4
[I] 2021/06/01 13:49:36 	ihl:                 5
[I] 2021/06/01 13:49:36 	type of service:     0b00000000
[I] 2021/06/01 13:49:36 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:36 	id:                  8587
[I] 2021/06/01 13:49:36 	flags:               0b010
[I] 2021/06/01 13:49:36 	fragment offset:     0
[I] 2021/06/01 13:49:36 	ttl:                 64
[I] 2021/06/01 13:49:36 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:36 	checksum:            0x0000
[I] 2021/06/01 13:49:36 	source address:      192.0.2.1
[I] 2021/06/01 13:49:36 	destination address: 192.0.2.2
[I] 2021/06/01 13:49:36 ▶ Incoming ICMP packet
[I] 2021/06/01 13:49:36 	type:     8 (Echo)
[I] 2021/06/01 13:49:36 	code:     0
[I] 2021/06/01 13:49:36 	checksum: 0x1679
[I] 2021/06/01 13:49:36 	id:       48
[I] 2021/06/01 13:49:36 	seq:      3
[I] 2021/06/01 13:49:36 ▶ Outgoing ICMP packet
[I] 2021/06/01 13:49:36 	type:     0 (Echo Reply)
[I] 2021/06/01 13:49:36 	code:     0
[I] 2021/06/01 13:49:36 	checksum: 0x1e79
[I] 2021/06/01 13:49:36 ▶ Outgoing IP packet
[I] 2021/06/01 13:49:36 	version:             IPv4
[I] 2021/06/01 13:49:36 	ihl:                 5
[I] 2021/06/01 13:49:36 	type of service:     0b00000000
[I] 2021/06/01 13:49:36 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:36 	id:                  2
[I] 2021/06/01 13:49:36 	flags:               0b000
[I] 2021/06/01 13:49:36 	fragment offset:     0
[I] 2021/06/01 13:49:36 	ttl:                 255
[I] 2021/06/01 13:49:36 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:36 	checksum:            0x37a3
[I] 2021/06/01 13:49:36 	source address:      192.0.2.2
[I] 2021/06/01 13:49:36 	destination address: 192.0.2.1
[I] 2021/06/01 13:49:36 ▶ Outgoing Ethernet frame
[I] 2021/06/01 13:49:36 	dst:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:36 	src:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:36 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:36 	payload: 45 00 00 54 00 02 00 00 ff 01 37 a3 c0 00 02 02 c0 00 02 01
[I] 2021/06/01 13:49:36 			 00 00 1e 79 00 30 00 03 f0 3a b6 60 00 00 00 00 6e e5 0d 00
[I] 2021/06/01 13:49:36 			 00 00 00 00 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f
[I] 2021/06/01 13:49:36 			 20 21 22 23 24 25 26 27 28 29 2a 2b 2c 2d 2e 2f 30 31 32 33
[I] 2021/06/01 13:49:36 ▶ Ethernet frame was sent: 98 bytes (payload: 84 bytes)
[I] 2021/06/01 13:49:37 ▶ Event occurred
[I] 2021/06/01 13:49:37 	events: 1
[I] 2021/06/01 13:49:37 	device: net1 (tap0)
[I] 2021/06/01 13:49:37 ▶ Ethernet frame was received: 98 bytes
[I] 2021/06/01 13:49:37 ▶ Incoming ethernet frame
[I] 2021/06/01 13:49:37 	dst:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:37 	src:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:37 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:37 ▶ Incoming IP packet
[I] 2021/06/01 13:49:37 	version:             IPv4
[I] 2021/06/01 13:49:37 	ihl:                 5
[I] 2021/06/01 13:49:37 	type of service:     0b00000000
[I] 2021/06/01 13:49:37 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:37 	id:                  8684
[I] 2021/06/01 13:49:37 	flags:               0b010
[I] 2021/06/01 13:49:37 	fragment offset:     0
[I] 2021/06/01 13:49:37 	ttl:                 64
[I] 2021/06/01 13:49:37 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:37 	checksum:            0x0000
[I] 2021/06/01 13:49:37 	source address:      192.0.2.1
[I] 2021/06/01 13:49:37 	destination address: 192.0.2.2
[I] 2021/06/01 13:49:37 ▶ Incoming ICMP packet
[I] 2021/06/01 13:49:37 	type:     8 (Echo)
[I] 2021/06/01 13:49:37 	code:     0
[I] 2021/06/01 13:49:37 	checksum: 0x2771
[I] 2021/06/01 13:49:37 	id:       48
[I] 2021/06/01 13:49:37 	seq:      4
[I] 2021/06/01 13:49:37 ▶ Outgoing ICMP packet
[I] 2021/06/01 13:49:37 	type:     0 (Echo Reply)
[I] 2021/06/01 13:49:37 	code:     0
[I] 2021/06/01 13:49:37 	checksum: 0x2f71
[I] 2021/06/01 13:49:37 ▶ Outgoing IP packet
[I] 2021/06/01 13:49:37 	version:             IPv4
[I] 2021/06/01 13:49:37 	ihl:                 5
[I] 2021/06/01 13:49:37 	type of service:     0b00000000
[I] 2021/06/01 13:49:37 	total length:        84 bytes (payload: 64 bytes)
[I] 2021/06/01 13:49:37 	id:                  3
[I] 2021/06/01 13:49:37 	flags:               0b000
[I] 2021/06/01 13:49:37 	fragment offset:     0
[I] 2021/06/01 13:49:37 	ttl:                 255
[I] 2021/06/01 13:49:37 	protocol:            1 (ICMP)
[I] 2021/06/01 13:49:37 	checksum:            0x37a2
[I] 2021/06/01 13:49:37 	source address:      192.0.2.2
[I] 2021/06/01 13:49:37 	destination address: 192.0.2.1
[I] 2021/06/01 13:49:37 ▶ Outgoing Ethernet frame
[I] 2021/06/01 13:49:37 	dst:  e2:20:76:2d:5d:2d
[I] 2021/06/01 13:49:37 	src:  01:00:e2:20:76:2d
[I] 2021/06/01 13:49:37 	type: 0x0800 (IPv4)
[I] 2021/06/01 13:49:37 	payload: 45 00 00 54 00 03 00 00 ff 01 37 a2 c0 00 02 02 c0 00 02 01
[I] 2021/06/01 13:49:37 			 00 00 2f 71 00 30 00 04 f1 3a b6 60 00 00 00 00 5c ec 0d 00
[I] 2021/06/01 13:49:37 			 00 00 00 00 10 11 12 13 14 15 16 17 18 19 1a 1b 1c 1d 1e 1f
[I] 2021/06/01 13:49:37 			 20 21 22 23 24 25 26 27 28 29 2a 2b 2c 2d 2e 2f 30 31 32 33
[I] 2021/06/01 13:49:37 ▶ Ethernet frame was sent: 98 bytes (payload: 84 bytes)
```
</details>

## References
- Articles
  - [Demystifying memory management in modern programming languages](https://deepu.tech/memory-management-in-programming)
- Blogs
  - [golangspec](https://medium.com/golangspec)
  - [Vincent Blanchon](https://medium.com/@blanchon.vincent)
- Books
  - [Hands-On High Performance with Go](https://www.packtpub.com/product/hands-on-high-performance-with-go/9781789805789)
  - [Mastering Go - Second Edition](https://www.packtpub.com/product/mastering-go-second-edition/9781838559335)
  - [Multiplayer Game Programming: Architecting Networked Games](https://www.oreilly.com/library/view/multiplayer-game-programming/9780134034355)
  - [TCP/IP Illustrated](https://en.wikipedia.org/wiki/TCP/IP_Illustrated)
- Docs
  - Codecov
    - [About the Codecov yaml](https://docs.codecov.io/docs/codecov-yaml)
    - [Codecov Delta](https://docs.codecov.io/docs/codecov-delta)
  - Docker
    - [Compose file](https://docs.docker.com/compose/compose-file/)
    - [Dockerfile](https://docs.docker.com/engine/reference/builder/)
    - [Overview of Docker Compose](https://docs.docker.com/compose/)
  - fedora
      - [DNF, the next-generation replacement for YUM](https://dnf.readthedocs.io/en/latest/index.html)
      - [Understanding and administering systemd](https://docs.fedoraproject.org/en-US/quick-docs/understanding-and-administering-systemd/)
  - GitHub Actions
    - [Reference](https://docs.github.com/en/actions/reference)
  - Mutagen
    - [File synchronization](https://mutagen.io/documentation/synchronization)
  - Slack
    - Builder
      - [Block Kit](https://app.slack.com/block-kit-builder/T82HLGL84#%7B%22blocks%22:%5B%5D%7D)
      - [Message](https://api.slack.com/docs/messages/builder?msg=%7B%22text%22%3A%22I%20am%20a%20test%20message%22%2C%22attachments%22%3A%5B%7B%22text%22%3A%22And%20here%E2%80%99s%20an%20attachment!%22%7D%5D%7D)
    - [Reference guides for app features](https://api.slack.com/reference)
  - UI
    - [Error Message Guidelines](https://docs.microsoft.com/en-us/windows/win32/debug/error-message-guidelines?redirectedfrom=MSDN)
    - [Messages: UI Text Guidelines](https://docs.microsoft.com/en-us/previous-versions/windows/desktop/bb226791(v=vs.85))
  - Vagrant
    - [Vagrantfile](https://www.vagrantup.com/docs/vagrantfile)
- Go
  - Articles
    - [A Closer Look at Golang From an Architect’s Perspective](https://thenewstack.io/a-closer-look-at-golang-from-an-architects-perspective)
    - [Go at Google: Language Design in the Service of Software Engineering](https://talks.golang.org/2012/splash.article)
  - Docs
      - [Command compile](https://golang.org/cmd/compile/)
      - [Command vet](https://golang.org/cmd/vet)
      - [Compiler And Runtime Optimizations](https://github.com/golang/go/wiki/CompilerOptimizations)
      - [Data Race Detector](https://golang.org/doc/articles/race_detector)
      - [Effective Go](https://golang.org/doc/effective_go)
      - [Frequently Asked Questions (FAQ)](https://golang.org/doc/faq)
        - [Should I define methods on values or pointers?](https://golang.org/doc/faq#methods_on_values_or_pointers)
          - [Value vs Pointer Receivers](https://h12.io/article/value-vs-pointer-receivers)
      - [Getting Started with Code Coverage for Golang](https://about.codecov.io/blog/getting-started-with-code-coverage-for-golang/)
      - [Go by Example](https://gobyexample.com/)
      - [Package testing](https://golang.org/pkg/testing)
      - [The Go Memory Model](https://golang.org/ref/mem)
      - [The Go Programming Language Specification](https://golang.org/ref/spec)
  - Packages
    - [github.com/google/go-cmp/cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp)
- Makefile
  - [A Good Makefile for Go](https://kodfabrik.com/journal/a-good-makefile-for-go)
  - [GNU Make Manual](https://www.gnu.org/software/make/manual/)
- Papers
  - [Fifty Shades of Congestion Control: A Performance and Interactions Evaluation](https://arxiv.org/abs/1903.03852)
- RFC
  - [791: Internet Protocol](https://tools.ietf.org/html/rfc791)
  - [793: Transmission Control Protocol](https://tools.ietf.org/html/rfc793)
  - [5952: 5. Text Representation of Special Addresses](https://tools.ietf.org/html/rfc5952#section-5)
- Open Source Software
  - [avelino / awesome-go](https://github.com/avelino/awesome-go)
  - [pandax381 / microps](https://github.com/pandax381/microps)
  - [torvalds / linux](https://github.com/torvalds/linux)
    - [net / ethernet](https://github.com/torvalds/linux/tree/master/net/ethernet)
    - [net / ipv4](https://github.com/torvalds/linux/tree/master/net/ipv4)
- Wikipedia
  - Computer Network
      - [0.0.0.0](https://en.wikipedia.org/wiki/0.0.0.0)
      - [Address Resolution Protocol](https://en.wikipedia.org/wiki/Address_Resolution_Protocol)
      - [Data link layer](https://en.wikipedia.org/wiki/Data_link_layer)
      - [Ethernet frame](https://en.wikipedia.org/wiki/Ethernet_frame)
      - [Internet Protocol](https://en.wikipedia.org/wiki/Internet_Protocol)
      - [Network Layer](https://en.wikipedia.org/wiki/Network_layer)
      - [TUN/TAP](https://en.wikipedia.org/wiki/TUN/TAP)
      - [TCP congestion control](https://en.wikipedia.org/wiki/TCP_congestion_control)
      - [Transmission Control Protocol](https://en.wikipedia.org/wiki/Transmission_Control_Protocol)
      - [User Datagram Protocol](https://en.wikipedia.org/wiki/User_Datagram_Protocol)
  - Memory
    - [Memory barrier](https://en.wikipedia.org/wiki/Memory_barrier)
    - [Memory management](https://en.wikipedia.org/wiki/Memory_management)
    - [Memory model (programming)](https://en.wikipedia.org/wiki/Memory_model_(programming))
    - [Memory safety](https://en.wikipedia.org/wiki/Memory_safety)
    - [Region-based memory management](https://en.wikipedia.org/wiki/Region-based_memory_management)
  - Others
    - [Communicating sequential processes](https://en.wikipedia.org/wiki/Communicating_sequential_processes)
    - [Continuation](https://en.wikipedia.org/wiki/Continuation)
    - [Cooperative multitasking](https://en.wikipedia.org/wiki/Cooperative_multitasking)
    - [Critical section](https://en.wikipedia.org/wiki/Critical_section)
    - [Ellipsis (computer programming)](https://en.wikipedia.org/wiki/Ellipsis_(computer_programming))
    - [Fair-share scheduling](https://en.wikipedia.org/wiki/Fair-share_scheduling)
    - [Fork-join model](https://en.wikipedia.org/wiki/Fork%E2%80%93join_model)
    - [Generic programming](https://en.wikipedia.org/wiki/Generic_programming)
    - [Monitor (synchronization)](https://en.wikipedia.org/wiki/Monitor_(synchronization))
    - [Pipeline stall](https://en.wikipedia.org/wiki/Pipeline_stall)
    - [Preemption (computing)](https://en.wikipedia.org/wiki/Preemption_(computing))
    - [Process substitution](https://en.wikipedia.org/wiki/Process_substitution)
    - [Reflective programming](https://en.wikipedia.org/wiki/Reflective_programming)
    - [Task (computing)](https://en.wikipedia.org/wiki/Task_(computing))
    - [Template metaprogramming](https://en.wikipedia.org/wiki/Template_metaprogramming)
    - [Thread pool](https://en.wikipedia.org/wiki/Thread_pool)
    - [Threading models](https://en.wikipedia.org/wiki/Thread_(computing)#Threading_models)
    - [Work stealing](https://en.wikipedia.org/wiki/Work_stealing)

## Notes
- GO
  - Asynchronous Preemption
    - [Go: Goroutine and Preemption](https://medium.com/a-journey-with-go/go-goroutine-and-preemption-d6bc2aa2f4b7)
    - [Goroutine preemptive scheduling with new features of go 1.14](https://developpaper.com/goroutine-preemptive-scheduling-with-new-features-of-go-1-14)
  - Context
    - [Go: Context and Cancellation by Propagation](https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c)
- Programming
  - [History of programming languages](https://en.wikipedia.org/wiki/History_of_programming_languages)
- Question and answer site
  - [C++ : how to close a tcp socket (server) when receiving SIGKILL](https://stackoverflow.com/questions/21329861/c-how-to-close-a-tcp-socket-server-when-receiving-sigkill)
  - [Understanding INADDR_ANY for socket programming](https://stackoverflow.com/questions/16508685/understanding-inaddr-any-for-socket-programming)
  - [What is the difference between AF_INET and PF_INET in socket programming?](https://stackoverflow.com/questions/6729366/what-is-the-difference-between-af-inet-and-pf-inet-in-socket-programming)

## Memos
- Use the context types properly acording to the rubric.
  - [context.Background()](https://github.com/golang/go/blob/a72622d028077643169dc48c90271a82021f0534/src/context/context.go#L208)
  - [context.TODO()](https://github.com/golang/go/blob/a72622d028077643169dc48c90271a82021f0534/src/context/context.go#L216)
- Definition of errno in Go
  - [go/src/syscall/zerrors_linux_amd64.go](https://github.com/golang/go/blob/master/src/syscall/zerrors_linux_amd64.go#L1205)
