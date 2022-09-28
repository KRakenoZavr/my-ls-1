package flags

import (
	"errors"
	"testing"
)

func TestNewFlags(t *testing.T) {
	// 1 normal flag
	flag, err := NewFlags([]string{"-l"}, 1)
	if err != nil {
		t.Error("should not be error")
	}

	if !contains("l", flag.flags) {
		t.Error("should contain", "l")
	}

	// incorrect flag
	_, err = NewFlags([]string{"-d"}, 1)
	if !errors.Is(err, errorFlag) {
		t.Error("should be error")
	}

	// couple normal flag
	flag1, err := NewFlags([]string{"-l", "-aR"}, 1)
	if err != nil {
		t.Error("should not be error")
	}
	if !contains("l", flag1.flags) || !contains("R", flag1.flags) || !contains("a", flag1.flags) {
		t.Error("should contain", "l")
	}
}
