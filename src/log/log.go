package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const dtFormat = "2006/02/01 15:04:05"

var mtx sync.Mutex
var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

func I(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	dt := time.Now().Format(dtFormat)
	if s != "" {
		_, _ = fmt.Fprintf(stdout, "[I] %s %s\n", dt, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(stdout, "                        %s\n", v)
	}
}

func W(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	dt := time.Now().Format(dtFormat)
	if s != "" {
		_, _ = fmt.Fprintf(stdout, "\u001B[1;33m[W] %s %s\u001B[0m\n", dt, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(stdout, "\u001B[1;33m                        %s\u001B[0m\n", v)
	}
}

func E(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	dt := time.Now().Format(dtFormat)
	if s != "" {
		_, _ = fmt.Fprintf(stderr, "\u001B[1;31m[E] %s %s\u001B[0m\n", dt, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(stderr, "\u001B[1;31m                        %s\u001B[0m\n", v)
	}
}

func F(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	dt := time.Now().Format(dtFormat)
	if s != "" {
		_, _ = fmt.Fprintf(stderr, "\u001B[1;31m[F] %s %s\u001B[0m\n", dt, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(stderr, "\u001B[1;31m                        %s\u001B[0\n", v)
	}
	os.Exit(1)
}

func CaptureLogOutput(f func()) string {
	defer resetOutput()

	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	setOutput(writer)

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()

	f()
	_ = writer.Close()
	ret := <-out
	_ = reader.Close()

	return ret
}

func DisableOutput() {
	setOutput(ioutil.Discard)
}

func EnableOutput() {
	resetOutput()
}

func resetOutput() {
	stdout = os.Stdout
	stderr = os.Stderr
}

func setOutput(writer io.Writer) {
	stdout = writer
	stderr = writer
}
