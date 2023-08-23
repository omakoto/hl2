package term

import (
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	DefaultTermWidth = 80
)

var (
	TermWidth = DefaultTermWidth
)

func GetTermWidth() int {
	width, _, err := term.GetSize(1)
	if err != nil {
		return DefaultTermWidth
	}
	return width
}

func NewDefaultTerm() Term {
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
