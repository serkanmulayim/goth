package testutils

import (
	"testing"
)

func AssertEqualsString(t *testing.T, out string, expected string) {
	if out != expected {
		t.Errorf("out=%s, wanted=%s", out, expected)
	}
}

func AssertEqualsInt(t *testing.T, out int, expected int) {
	if out != expected {
		t.Errorf("out=%d, wanted=%d", out, expected)
	}
}

func AssertEqualsInt64(t *testing.T, out int64, expected int64) {
	if out != expected {
		t.Errorf("out=%d, wanted=%d", out, expected)
	}
}
