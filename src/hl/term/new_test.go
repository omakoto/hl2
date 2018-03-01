package term

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewDefaultTerm(t *testing.T) {
	os.Setenv("TERM", "")
	os.Setenv("COLORTERM", "")
	assert.IsType(t, &DumbTerm{}, NewDefaultTerm())

	os.Setenv("TERM", "vt100")
	assert.IsType(t, &ConsoleTerm{}, NewDefaultTerm())

	os.Setenv("TERM", "xterm")
	assert.IsType(t, &Rgb8Term{}, NewDefaultTerm())

	os.Setenv("TERM", "xterm256")
	assert.IsType(t, &Rgb8Term{}, NewDefaultTerm())

	os.Setenv("TERM", "xterm256")
	os.Setenv("COLORTERM", "truecolor")
	assert.IsType(t, &Rgb24Term{}, NewDefaultTerm())
}
