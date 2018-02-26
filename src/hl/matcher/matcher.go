package matcher

import (
	"github.com/omakoto/hl2/src/hl"
)

type Matcher interface {
	Matches(target []byte) [][]int
	String() string
}

var NoPcre = false

func CompileWithContext(context hl.Context, pattern string) (Matcher, error) {
	flags := None
	if context.IgnoreCase() {
		flags |= IgnoreCase
	}
	return Compile(pattern, flags)
}

func Compile(pattern string, flags Flags) (Matcher, error) {
	if NoPcre {
		return CompileGo(pattern, flags)
	}
	return CompilePcre(pattern, flags)
}
