package log

import (
	"regexp"
	"testing"
)

func TestI(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;34m\\[INFO]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} hello$")
	got := CaptureLogOutput(func () {
		I("hello")
	})
	if ! want.MatchString(got) {
		t.Errorf("EtherDump() = %v; want %v", got, want.String())
	}
}

func TestW(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;33m\\[WARN]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} hello$")
	got := CaptureLogOutput(func () {
		W("hello")
	})
	if ! want.MatchString(got) {
		t.Errorf("EtherDump() = %v; want %v", got, want.String())
	}
}

func TestE(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;31m\\[ERROR]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} hello$")
	got := CaptureLogOutput(func () {
		E("hello")
	})
	if ! want.MatchString(got) {
		t.Errorf("EtherDump() = %v; want %v", got, want.String())
	}
}

func TestF(t *testing.T) {
	want, _ := regexp.Compile("^\u001B\\[1;31m\\[FATAL]\u001B\\[0m [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} hello$")
	got := CaptureLogOutput(func () {
		F("hello")
	})
	if ! want.MatchString(got) {
		t.Errorf("EtherDump() = %v; want %v", got, want.String())
	}
}
