package util

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

var (
	Debug = false
)

func Debugging() bool {
	return Debug
}

func Debugf(format string, args ...interface{}) {
	if Debug {
		fmt.Printf(format, args...)
	}
}

func Dump(prefix string, arg interface{}) {
	if Debug {
		fmt.Print(prefix)
		spew.Dump(arg)
	}
}
