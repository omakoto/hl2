package matcher

import (
	"regexp"
)

type Flags int

type matcherGo struct {
	srcPattern  string
	realPattern string
	negate      bool
	pattern     *regexp.Regexp
}

var _ = Matcher((*matcherGo)(nil))

func (r *matcherGo) String() string {
	return r.srcPattern
}

func CompileGo(pattern string, flags Flags) (Matcher, error) {
	negate := false

	realPattern := pattern
	err := preProcess(&realPattern, &negate)
	if err != nil {
		return nil, err
	}
	if (flags & IgnoreCase) != 0 {
		realPattern = "(?i)" + realPattern
	}
	pat, err := regexp.Compile(realPattern)
	if err != nil {
		return nil, err
	}

	return &matcherGo{srcPattern: pattern, realPattern: realPattern, negate: negate, pattern: pat}, nil
}

func (r *matcherGo) Matches(target []byte) [][]int {
	if r.negate {
		if !r.pattern.Match(target) {
			return [][]int{{0, len(target)}}
		}
		return nil
	}
	matches := r.pattern.FindAllSubmatchIndex(target, -1)
	if len(matches) == 0 {
		return nil
	}
	captures := (len(matches[0]) / 2) - 1
	if captures == 0 {
		return matches
	}

	ret := make([][]int, 0, captures*len(matches))
	for i := 0; i < len(matches); i++ {
		for j := 0; j < captures; j++ {
			ret = append(ret, []int{matches[i][2+j*2], matches[i][2+j*2+1]})
		}
	}
	return ret
}
