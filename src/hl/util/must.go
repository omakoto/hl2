package util

import "fmt"

func Must(f func() error) {
	err := f()
	if err != nil {
		panic(fmt.Sprintf("Must function failed: %s", err))
	}
}
