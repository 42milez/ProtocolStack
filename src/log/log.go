package log

import (
	"bytes"
	"io"
	"io/ioutil"
	goLog "log"
	"os"
	"regexp"
)

func format(s string) (ret string) {
	if m, _ := regexp.MatchString("^\\w", s); m {
		ret = "▶ " + s
	} else {
		ret = s
	}
	return
}

func I(s string) {
	i.Println(format(s))
}

func W(s string) {
	w.Println(format(s))
}

func E(s string) {
	e.Println(format(s))
}

func F(s string) {
	f.Fatalln(format(s))
}

func CaptureLogOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	defer func() {
		resetOutput()
	}()

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
	i.SetOutput(os.Stdout)
	w.SetOutput(os.Stdout)
	e.SetOutput(os.Stderr)
	f.SetOutput(os.Stderr)
}

func setOutput(writer io.Writer) {
	i.SetOutput(writer)
	w.SetOutput(writer)
	e.SetOutput(writer)
	f.SetOutput(writer)
}

var i *goLog.Logger
var w *goLog.Logger
var e *goLog.Logger
var f *goLog.Logger

func init() {
	i = goLog.New(os.Stdout, "\u001B[1;34m[INFO]\u001B[0m ", goLog.LstdFlags)
	w = goLog.New(os.Stdout, "\u001B[1;33m[WARN]\u001B[0m ", goLog.LstdFlags)
	e = goLog.New(os.Stderr, "\u001B[1;31m[ERROR]\u001B[0m ", goLog.LstdFlags)
	f = goLog.New(os.Stderr, "\u001B[1;31m[FATAL]\u001B[0m ", goLog.LstdFlags)
}
