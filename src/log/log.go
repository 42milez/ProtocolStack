package log

import (
	"bytes"
	"io"
	"io/ioutil"
	goLog "log"
	"os"
)

func I(s string, v ...interface{}) {
	if len(v) == 0 {
		i.Println(s)
	} else {
		i.Printf(s+"\n", v...)
	}
}

func W(s string, v ...interface{}) {
	if len(v) == 0 {
		w.Println(s)
	} else {
		w.Printf(s+"\n", v...)
	}
}

func E(s string, v ...interface{}) {
	if len(v) == 0 {
		e.Println(s)
	} else {
		e.Printf(s+"\n", v...)
	}
}

func F(s string, v ...interface{}) {
	if len(v) == 0 {
		f.Fatalln(s)
	} else {
		f.Fatalf(s+"\n", v...)
	}
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
