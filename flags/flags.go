package flags

import (
	"errors"
	"strings"
)

// t - sort  by modification date
// r - reverse
// a - show all files including .
// R - recursive
// l - long listing format

var flagsList = []string{
	"t", "r", "a", "R", "l",
}

var errorFlag = errors.New("incorrect flag")

type Flag struct {
	flags []string
	many  bool
}

func NewFlags(args []string, len int) (*Flag, error) {
	var flags []string

	for _, l := range args {
		maybeFlags := strings.Split(l, "")
		flags = append(flags, maybeFlags[1:]...)
	}

	if validate(flags) {
		return &Flag{flags: flags, many: len > 1}, nil
	}

	return nil, errorFlag
}

func (f *Flag) Contains(flag string) bool {
	return contains(flag, f.flags)
}

func (f *Flag) IsMany() bool {
	return f.many
}

func contains(element string, array []string) bool {
	for _, l := range array {
		if l == element {
			return true
		}
	}
	return false
}

func validate(flags []string) bool {
	for _, l := range flags {
		if !contains(l, flagsList) {
			return false
		}
	}
	return true
}
