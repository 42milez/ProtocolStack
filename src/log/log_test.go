package log

import (
	"regexp"
	"strings"
	"testing"
)

func trim(s string) string {
	ret := strings.Replace(s, "\t", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	return ret
}

func TestI(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;34m\\[I]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} ▶ info$")
	got := CaptureLogOutput(func() {
		I("info")
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;34m\\[I]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} info$")
	got = CaptureLogOutput(func() {
		I("\tinfo")
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}
}

func TestW(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;33m\\[W]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} ▶ warning$")
	got := CaptureLogOutput(func() {
		W("warning")
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;33m\\[W]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} warning$")
	got = CaptureLogOutput(func() {
		W("\twarning")
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}
}

func TestE(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;31m\\[E]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} ▶ error$")
	got := CaptureLogOutput(func() {
		E("error")
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;31m\\[E]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} error$")
	got = CaptureLogOutput(func() {
		E("\terror")
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}
}

func TestF(t *testing.T) {
	// No tests exists because log.Fatal*() calls os.Exit(1).
}
