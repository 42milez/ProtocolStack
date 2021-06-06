package log

import (
	"regexp"
	"strings"
	"testing"
)

func TestI(t *testing.T) {
	want, _ := regexp.Compile(`^\[I] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} info$`)
	got := CaptureLogOutput(func() {
		I("info")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^\[I] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} InfoHelloWorld$`)
	got = CaptureLogOutput(func() {
		I("Info", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}
}

func TestW(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;33m\\[W] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Warning\u001B\\[0m$")
	got := CaptureLogOutput(func() {
		W("Warning")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;33m\\[W] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Warning\u001B\\[0m\u001B\\[1;33mHello\u001B\\[0m\u001B\\[1;33mWorld\u001B\\[0m$")
	got = CaptureLogOutput(func() {
		W("Warning", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}
}

func TestE(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;31m\\[E] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Error\u001B\\[0m$")
	got := CaptureLogOutput(func() {
		E("Error")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;31m\\[E] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Error\u001B\\[0m\u001B\\[1;31mHello\u001B\\[0m\u001B\\[1;31mWorld\u001B\\[0m$")
	got = CaptureLogOutput(func() {
		E("Error", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}
}

func TestF(t *testing.T) {
	// No tests exists because log.Fatal*() calls os.Exit(1).
}

func Trim(s string) (ret string) {
	ret = strings.Replace(s, "", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	ret = strings.Replace(ret, "                        ", "", -1)
	return
}
