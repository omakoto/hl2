package term

import (
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

const (
	DefaultTermWidth = 80
)

var (
	TermWidth = DefaultTermWidth
)

func GetTermWidth() int {
	width, _, err := terminal.GetSize(1)
	if err != nil {
		return DefaultTermWidth
	}
	return width
}

func NewTerm() Term {
	var t Term
	termEnv := os.Getenv("TERM")
	if strings.HasPrefix(termEnv, "xterm") {
		if os.Getenv("COLORTERM") == "truecolor" {
			t = NewRgb24Term(TermWidth)
		} else {
			t = NewRgb8Term(TermWidth)
		}
	} else if termEnv != "" {
		t = NewConsoleTerm(TermWidth)
	} else {
		t = NewDumbTerm()
	}
	return t
}
