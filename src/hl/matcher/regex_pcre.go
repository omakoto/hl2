package matcher

import (
	"fmt"
	"github.com/d4l3k/go-pcre"
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

func CompilePcre(pattern string, flags Flags) (*matcherPcre, error) {
	negate := false

	realPattern := pattern
	err := preProcess(&realPattern, &negate)
	if err != nil {
		return nil, err
	}
	pcreFlags := pcre.UTF8 | pcre.NO_UTF8_CHECK
	if (flags & IgnoreCase) != 0 {
		pcreFlags |= pcre.CASELESS
	}
	pat, err := pcre.Compile(realPattern, pcreFlags)
	if err != nil {
		return nil, err
	}

	return &matcherPcre{srcPattern: pattern, realPattern: realPattern, negate: negate, pattern: &pat}, nil
}

func (r *matcherPcre) Matches(target []byte) [][]int {
	if r.negate {
		if !r.pattern.Matcher(target, 0).Matches() {
			return [][]int{{0, len(target)}}
		}
		return nil
	}
	m := r.pattern.NewMatcher()

	res := make([][]int, 0)

	start := 0
	for {
		if start > len(target) {
			break
		}
		//util.Debugf("Start=%d\n", start)
		//util.Debugf("  Target=\"%s\"\n", string(target[start:]))
		result := m.Exec(target[start:], 0)
		if result == pcre.ERROR_NOMATCH || result == pcre.ERROR_BADUTF8 {
			break
		}
		if result < 0 {
			panic(fmt.Sprintf("pcre_exec returned %d for pattern \"%s\" (len %db), target \"%s\", index %d", result, r.realPattern, len([]byte(r.realPattern)), string(target), start))
		}
		// matchStart := m.GroupIndices(0)[0]
		matchEnd := m.GroupIndices(0)[1] + start
		//util.Debugf("  Match=%d, %d\n", matchStart, matchEnd)

		groups := m.Groups()
		if groups == 0 {
			indexes := make([]int, 2)
			indexes[0] = m.GroupIndices(0)[0] + start
			indexes[1] = m.GroupIndices(0)[1] + start
			res = append(res, indexes)
		} else {
			for j := 1; j <= groups; j++ {
				gi := m.GroupIndices(j)
				if gi != nil {
					indexes := make([]int, 2)
					indexes[0] = gi[0] + start
					indexes[1] = gi[1] + start
					res = append(res, indexes)
				}
			}
		}

		if matchEnd == start {
			matchEnd++
		}
		start = matchEnd
	}
	if len(res) == 0 {
		return nil
	}
	return res
}
