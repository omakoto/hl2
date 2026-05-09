package matcher

import (
	"github.com/dlclark/regexp2"
)

type matcherPcre struct {
	srcPattern  string
	realPattern string
	negate      bool
	pattern     *regexp2.Regexp
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
	var re2Flags regexp2.RegexOptions
	if (flags & IgnoreCase) != 0 {
		re2Flags |= regexp2.IgnoreCase
	}
	pat, err := regexp2.Compile(realPattern, re2Flags)
	if err != nil {
		return nil, err
	}

	return &matcherPcre{srcPattern: pattern, realPattern: realPattern, negate: negate, pattern: pat}, nil
}

func (r *matcherPcre) Matches(target []byte) [][]int {
	s := string(target)

	// regexp2 returns rune positions; build mapping to byte offsets
	runeToBytePos := make([]int, 0, len(s))
	for bytePos := range s {
		runeToBytePos = append(runeToBytePos, bytePos)
	}
	runeToBytePos = append(runeToBytePos, len(target))

	byteOffset := func(runeIdx int) int {
		if runeIdx >= 0 && runeIdx < len(runeToBytePos) {
			return runeToBytePos[runeIdx]
		}
		return len(target)
	}

	if r.negate {
		m, _ := r.pattern.FindStringMatch(s)
		if m == nil {
			return [][]int{{0, len(target)}}
		}
		return nil
	}

	res := make([][]int, 0)
	for m, _ := r.pattern.FindStringMatch(s); m != nil; m, _ = r.pattern.FindNextMatch(m) {
		groups := m.Groups()
		if len(groups) <= 1 {
			res = append(res, []int{
				byteOffset(groups[0].Index),
				byteOffset(groups[0].Index + groups[0].Length),
			})
		} else {
			for _, g := range groups[1:] {
				if len(g.Captures) > 0 {
					res = append(res, []int{
						byteOffset(g.Index),
						byteOffset(g.Index + g.Length),
					})
				}
			}
		}
	}

	if len(res) == 0 {
		return nil
	}
	return res
}
