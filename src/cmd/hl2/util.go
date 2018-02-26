package main

import (
	"fmt"
	"os"
)

func Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(Name+": "+format, args...)
	fmt.Fprint(os.Stderr, msg)

	if msg[len(msg)-1] != '\n' {
		fmt.Fprint(os.Stderr, "\n")
	}

	os.Exit(1)
}
