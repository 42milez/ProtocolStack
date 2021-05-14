package logger

import (
	"log"
	"os"
)

var i *log.Logger
var w *log.Logger
var e *log.Logger
var f *log.Logger

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
		f.Println(s)
	} else {
		f.Printf(s+"\n", v...)
	}
}

func init() {
	i = log.New(os.Stdout, "\u001B[1;34m[INFO]\u001B[0m ", log.LstdFlags)
	w = log.New(os.Stdout, "\u001B[1;33m[WARN]\u001B[0m ", log.LstdFlags)
	e = log.New(os.Stderr, "\u001B[1;31m[ERROR]\u001B[0m ", log.LstdFlags)
	f = log.New(os.Stderr, "\u001B[1;31m[FATAL]\u001B[0m ", log.LstdFlags)
}
