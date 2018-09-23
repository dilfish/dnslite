package dnslite

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	ret, err := GetIP("1.1.1.1:53")
	if err != nil || ret != "1.1.1.1" {
		t.Error("expect 1.1.1.1, got", ret, err)
	}
}
