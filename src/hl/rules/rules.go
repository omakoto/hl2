package rules

import (
	"github.com/omakoto/hl2/src/hl"
	"github.com/omakoto/hl2/src/hl/colors"
	"github.com/omakoto/hl2/src/hl/matcher"
	"github.com/omakoto/hl2/src/hl/term"
)

const InitialState = "INIT"

type DecorativeLine struct {
	Marker []byte
	Colors *term.RenderedColors
}

func NewDecorativeLine(context hl.Context, marker string, colors *colors.Colors) *DecorativeLine {
	t := context.Term()
	c := term.NewRenderedColors(t, colors)
	return &DecorativeLine{
		Marker: []byte(marker),
		Colors: c,
	}
}

type Rule struct {
	Matcher    matcher.Matcher
	PreMatcher matcher.Matcher

	After  int
	Before int

	Show bool
	Hide bool
	Stop bool

	MatchColors *term.RenderedColors
	LineColors  *term.RenderedColors

	PreLine  *DecorativeLine
	PostLine *DecorativeLine

	States    []string
	NextState string
}

func (r *Rule) IsForState(state string) bool {
	if len(r.States) == 0 {
		return true
	}
	for i := 0; i < len(r.States); i++ {
		if r.States[i] == state {
			return true
		}
	}
	return false
}
