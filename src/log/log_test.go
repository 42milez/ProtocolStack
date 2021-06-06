package log

import (
	"regexp"
	"strings"
	"testing"
)

func TestI(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;34m\\[I] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} info\u001B\\[0m$")
	got := CaptureLogOutput(func() {
		I("info")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;34m\\[I] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} info\u001B\\[0m                        a                        b$")
	got = CaptureLogOutput(func() {
		I("info", "a", "b")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}
}

func TestW(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;33m\\[W] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} warning\u001B\\[0m$")
	got := CaptureLogOutput(func() {
		W("warning")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;33m\\[W] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} warning\u001B\\[0m                        a                        b$")
	got = CaptureLogOutput(func() {
		W("warning", "a", "b")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}
}

func TestE(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;31m\\[E] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} error\u001B\\[0m$")
	got := CaptureLogOutput(func() {
		E("error")
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;31m\\[E] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} error\u001B\\[0m                        a                        b$")
	got = CaptureLogOutput(func() {
		E("\terror", "a", "b")
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
	ret = strings.Replace(s, "\t", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	return
}
