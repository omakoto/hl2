package matcher

const (
	NoFlags    Flags = 0
	IgnoreCase Flags = 1 << iota
)

type Matcher interface {
	Matches(target []byte) [][]int
	String() string
}

var NoPcre = false

func Compile(pattern string, flags Flags) (Matcher, error) {
	if NoPcre {
		return CompileGo(pattern, flags)
	}
	return CompilePcre(pattern, flags)
}
