package matcher

import (
	"go.elara.ws/pcre"
)

type matcherPcre struct {
	srcPattern  string
	realPattern string
	negate      bool
	pattern     *pcre.Regexp
}

var _ = Matcher((*matcherPcre)(nil))

func (r *matcherPcre) String() string {
	return r.srcPattern
}

func CompilePcre(pattern string, flags Flags) (Matcher, error) {
	negate := false

	realPattern := pattern
	err := preProcess(&realPattern, &negate)
	if err != nil {
		return nil, err
	}
	pcreFlags := pcre.UTF | pcre.NoUTFCheck
	if (flags & IgnoreCase) != 0 {
		pcreFlags |= pcre.Caseless
	}
	pat, err := pcre.CompileOpts(realPattern, pcreFlags)
	if err != nil {
		return nil, err
	}

	return &matcherPcre{srcPattern: pattern, realPattern: realPattern, negate: negate, pattern: pat}, nil
}

func (r *matcherPcre) Matches(target []byte) [][]int {
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
