package log

import (
	"regexp"
	"strings"
	"testing"
)

func format(s string) string {
	ret := strings.Replace(s, "\t", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	return ret
}

func TestI(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;34m\\[INFO]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} info$")
	got := CaptureLogOutput(func() {
		I("info")
	})
	got = format(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;34m\\[INFO]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} info log$")
	got = CaptureLogOutput(func() {
		I("info %v", "log")
	})
	got = format(got)
	if !want.MatchString(got) {
		t.Errorf("I() = %v; want %v", got, want.String())
	}
}

func TestW(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;33m\\[WARN]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} warning$")
	got := CaptureLogOutput(func() {
		W("warning")
	})
	got = format(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;33m\\[WARN]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} warning log$")
	got = CaptureLogOutput(func() {
		W("warning %v", "log")
	})
	got = format(got)
	if !want.MatchString(got) {
		t.Errorf("W() = %v; want %v", got, want.String())
	}
}

func TestE(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;31m\\[ERROR]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} error$")
	got := CaptureLogOutput(func() {
		E("error")
	})
	got = format(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}

	want, _ = regexp.Compile("^\u001B\\[1;31m\\[ERROR]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} error log$")
	got = CaptureLogOutput(func() {
		E("error %v", "log")
	})
	got = format(got)
	if !want.MatchString(got) {
		t.Errorf("E() = %v; want %v", got, want.String())
	}
}

func TestF(t *testing.T) {
	// No tests exists because log.Fatal*() calls os.Exit(1).
}
