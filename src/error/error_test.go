package error

import "testing"

func TestE_Error(t *testing.T) {
	want := "OK"
	got := OK.Error()
	if got != want {
		t.Errorf("OK.Error() = %s; want %s", got, want)
	}
}
