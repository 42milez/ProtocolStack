## ProtocolStack
[![CI](https://github.com/42milez/ProtocolStack/actions/workflows/ci.yml/badge.svg)](https://github.com/42milez/ProtocolStack/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/42milez/ProtocolStack/branch/main/graph/badge.svg?token=ALHDIWP6KH)](https://codecov.io/gh/42milez/ProtocolStack)  
This repository is for learning network programming in Go.

## Instructions
### 1. Prepare virtual machine
ProtocolStack needs TAP device for its capability. For the reason, this project uses a virtual machine (Linux) to debug the application. Software required is as follows: 

- [Mutagen](https://github.com/mutagen-io/mutagen)
- [Vagrant](https://www.vagrantup.com/)
- [VirtualBox](https://www.virtualbox.org/)

Note: Mutagen synchronizes files between local system and virtual machine.

### 2. Control virtual machine
`vm.sh` controls the virtual machine. The available commands are as follows:

- `start` start the virtual machine
- `stop` stop the virtual machine
- `restart` restart the virtual machine

You can perform these commands like bellow:

```shell
> ./vm.sh start   # Start a VM and create a Mutagen session.
> ./vm.sh stop    # Stop the VM and terminate the Mutagen session.
> ./vm.sh restart # Restart the VM and recreate the Mutagen session.
```

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
  - [DNF, the next-generation replacement for YUM](https://dnf.readthedocs.io/en/latest/index.html)
  - Docker
    - [Compose file](https://docs.docker.com/compose/compose-file/)
    - [Dockerfile](https://docs.docker.com/engine/reference/builder/)
    - [Overview of Docker Compose](https://docs.docker.com/compose/)
  - GitHub Actions
    - [Reference](https://docs.github.com/en/actions/reference)
  - Mutagen
    - [File synchronization](https://mutagen.io/documentation/synchronization)
  - Vagrant
    - [Vagrantfile](https://www.vagrantup.com/docs/vagrantfile)
- Go
  - [Command compile](https://golang.org/cmd/compile/)
  - [Command vet](https://golang.org/cmd/vet)
  - [Compiler And Runtime Optimizations](https://github.com/golang/go/wiki/CompilerOptimizations)
  - [Data Race Detector](https://golang.org/doc/articles/race_detector)
  - [Effective Go](https://golang.org/doc/effective_go)
  - [Frequently Asked Questions (FAQ)](https://golang.org/doc/faq)
    - [Should I define methods on values or pointers?](https://golang.org/doc/faq#methods_on_values_or_pointers)
      - [Value vs Pointer Receivers](https://h12.io/article/value-vs-pointer-receivers)
  - [Getting Started with Code Coverage for Golang](https://about.codecov.io/blog/getting-started-with-code-coverage-for-golang/)
  - [Package testing](https://golang.org/pkg/testing)
  - Packages
    - [github.com/google/go-cmp/cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp)
  - [The Go Memory Model](https://golang.org/ref/mem)
  - [The Go Programming Language Specification](https://golang.org/ref/spec)
- Makefile
  - [A Good Makefile for Go](https://kodfabrik.com/journal/a-good-makefile-for-go)
  - [GNU Make Manual](https://www.gnu.org/software/make/manual/)
- Papers
  - [Fifty Shades of Congestion Control: A Performance and Interactions Evaluation](https://arxiv.org/abs/1903.03852)
- RFC
  - [791: Internet Protocol](https://tools.ietf.org/html/rfc791)
  - [793: Transmission Control Protocol](https://tools.ietf.org/html/rfc793)
  - [5952: 5. Text Representation of Special Addresses](https://tools.ietf.org/html/rfc5952#section-5)
- Source Code
  - [pandax381 / microps](https://github.com/pandax381/microps)
  - [torvalds / linux](https://github.com/torvalds/linux)
    - [net / ethernet](https://github.com/torvalds/linux/tree/master/net/ethernet)
    - [net / ipv4](https://github.com/torvalds/linux/tree/master/net/ipv4)
- Wikipedia
  - Computer Network
      - [Address Resolution Protocol](https://en.wikipedia.org/wiki/Address_Resolution_Protocol)
      - [Data link layer](https://en.wikipedia.org/wiki/Data_link_layer)
      - [Ethernet frame](https://en.wikipedia.org/wiki/Ethernet_frame)
      - [Internet Protocol](https://en.wikipedia.org/wiki/Internet_Protocol)
      - [Network Layer](https://en.wikipedia.org/wiki/Network_layer)
      - [TUN/TAP](https://en.wikipedia.org/wiki/TUN/TAP)
      - [TCP congestion control](https://en.wikipedia.org/wiki/TCP_congestion_control)
      - [Transmission Control Protocol](https://en.wikipedia.org/wiki/Transmission_Control_Protocol)
      - [User Datagram Protocol](https://en.wikipedia.org/wiki/User_Datagram_Protocol)
  - Programming
    - [Continuation](https://en.wikipedia.org/wiki/Continuation)
    - [Cooperative multitasking](https://en.wikipedia.org/wiki/Cooperative_multitasking)
    - [Critical section](https://en.wikipedia.org/wiki/Critical_section)
    - [Fair-share scheduling](https://en.wikipedia.org/wiki/Fair-share_scheduling)
    - [Fork-join model](https://en.wikipedia.org/wiki/Fork%E2%80%93join_model)
    - [Monitor (synchronization)](https://en.wikipedia.org/wiki/Monitor_(synchronization))
    - [Pipeline stall](https://en.wikipedia.org/wiki/Pipeline_stall)
    - [Preemption (computing)](https://en.wikipedia.org/wiki/Preemption_(computing))
    - [Task (computing)](https://en.wikipedia.org/wiki/Task_(computing))
    - [Thread pool](https://en.wikipedia.org/wiki/Thread_pool)
    - [Threading models](https://en.wikipedia.org/wiki/Thread_(computing)#Threading_models)
    - [Work stealing](https://en.wikipedia.org/wiki/Work_stealing)
  - Others
    - [Process substitution](https://en.wikipedia.org/wiki/Process_substitution)

## Notes
- GO
  - Asynchronous Preemption
    - [Go: Goroutine and Preemption](https://medium.com/a-journey-with-go/go-goroutine-and-preemption-d6bc2aa2f4b7)
    - [Goroutine preemptive scheduling with new features of go 1.14](https://developpaper.com/goroutine-preemptive-scheduling-with-new-features-of-go-1-14)
  - Context
    - [Go: Context and Cancellation by Propagation](https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c)
- Programming
  - [History of programming languages](https://en.wikipedia.org/wiki/History_of_programming_languages)

## Memos
- Use the context types properly acording to the rubric.
  - [context.Background()](https://github.com/golang/go/blob/a72622d028077643169dc48c90271a82021f0534/src/context/context.go#L208)
  - [context.TODO()](https://github.com/golang/go/blob/a72622d028077643169dc48c90271a82021f0534/src/context/context.go#L216)
