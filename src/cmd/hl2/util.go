package main

import (
	"fmt"
	"os"
)

func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, Name+": "+format, args...)
	os.Exit(1)
}
