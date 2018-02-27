package hl

import (
	"github.com/omakoto/hl2/src/hl/term"
)

type Context interface {
	Term() term.Term

	IgnoreCase() bool
	DefaultHide() bool

	DefaultAfter() int
	DefaultBefore() int

	//AddRule(r *rules.Rule) error
	//AddSimpleRule(pattern, colors string) error
	//AddSimpleRangeRules(patternStart, colorsStart, patternEnd, colorsEnd string) error
}
