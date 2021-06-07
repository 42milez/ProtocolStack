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

const red = "31"
const yellow = "33"
const dtFormat = "2006/02/01 15:04:05"

var mtx sync.Mutex
var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

func doPrint(prefix string, s string, dt string, args ...string) {
	if s != "" {
		_, _ = fmt.Fprintf(stdout, "%s %s %s\n", prefix, dt, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(stdout, "                        %s\n", v)
	}
}

func doColorPrint(color string, prefix string, s string, dt string, args ...string) {
	if s != "" {
		_, _ = fmt.Fprintf(stdout, "\u001B[1;%sm%s %s %s\u001B[0m\n", color, prefix, dt, s)
	}
	for _, v := range args {
		_, _ = fmt.Fprintf(stdout, "\u001B[1;%sm                        %s\u001B[0m\n", color, v)
	}
}

func I(s string, args []string) {
	defer mtx.Unlock()
	mtx.Lock()
	doPrint("[I]", s, time.Now().Format(dtFormat), args...)
}

func W(s string, args []string) {
	defer mtx.Unlock()
	mtx.Lock()
	doColorPrint(yellow, "[W]", s, time.Now().Format(dtFormat), args...)
}

func E(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	doColorPrint(red, "[E]", s, time.Now().Format(dtFormat), args...)
}

func F(s string, args ...string) {
	defer mtx.Unlock()
	mtx.Lock()
	doColorPrint(red, "[F]", s, time.Now().Format(dtFormat), args...)
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
