package log

import (
	"regexp"
	"strings"
	"testing"
)

func TestD(t *testing.T) {
	want, _ := regexp.Compile(`^\[I] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Debug$`)
	got := CaptureLogOutput(func() {
		I("Debug")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("D() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^\[I] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} DebugHelloWorld$`)
	got = CaptureLogOutput(func() {
		I("Debug", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("D() = %v; want %v", got, want.String())
	}
}

func TestI(t *testing.T) {
	want, _ := regexp.Compile(`^\[I] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Info$`)
	got := CaptureLogOutput(func() {
		I("Info")
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
	want, _ := regexp.Compile(`^\[1;33m\[W] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Warning\[0m$`)
	got := CaptureLogOutput(func() {
		W("Warning")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^\[1;33m\[W] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Warning\[0m\[1;33mHello\[0m\[1;33mWorld\[0m$`)
	got = CaptureLogOutput(func() {
		W("Warning", "Hello", "World")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}
}

func TestE(t *testing.T) {
	want, _ := regexp.Compile(`^\[1;31m\[E] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Error\[0m$`)
	got := CaptureLogOutput(func() {
		E("Error")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile(`^\[1;31m\[E] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} Error\[0m\[1;31mHello\[0m\[1;31mWorld\[0m$`)
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
