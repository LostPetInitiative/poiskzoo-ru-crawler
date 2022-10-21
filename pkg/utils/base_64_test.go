package utils

import (
	"testing"
)

func TestStringEncode(t *testing.T) {
	var s string = "Привет!"
	var encoded string = Base64Encode([]byte(s))

	if encoded != "0J/RgNC40LLQtdGCIQ==" {
		t.Fail()
	}
}
