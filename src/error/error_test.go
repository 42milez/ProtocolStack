package error

import "testing"

func TestE_Error(t *testing.T) {
	for i := 0; i <= int(Error); i++ {
		got := E(i).Error()
		if got == "UNKNOWN" {
			t.Errorf("E(%d).Error() didn't return appropriate string", i)
		}
	}
}
