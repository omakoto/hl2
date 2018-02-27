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
	context hl.Context

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

func NewRule(context hl.Context) *Rule {
	return &Rule{context: context}
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

func (r *Rule) SetMatcher(pattern string) error {
	m, err := matcher.CompileWithContext(r.context, pattern)
	if err != nil {
		return err
	}
	r.Matcher = m
	return nil
}

func (r *Rule) SetPreMatcher(pattern string) error {
	m, err := matcher.CompileWithContext(r.context, pattern)
	if err != nil {
		return err
	}
	r.PreMatcher = m
	return nil
}

func (r *Rule) SetBefore(n int) {
	r.Before = n
}

func (r *Rule) SetAfter(n int) {
	r.After = n
}

func (r *Rule) SetShow(v bool) {
	r.Show = v
}

func (r *Rule) SetHide(v bool) {
	r.Hide = v
}

func (r *Rule) SetStop(v bool) {
	r.Stop = v
}

func (r *Rule) SetStates(states []string) {
	r.States = states
}

func (r *Rule) SetNextState(s string) {
	r.NextState = s
}

func (r *Rule) SetMatchColors(colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.MatchColors = term.NewRenderedColors(r.context.Term(), c)
	return nil
}

func (r *Rule) SetLineColors(colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.LineColors = term.NewRenderedColors(r.context.Term(), c)
	return nil
}

func (r *Rule) SetPreLine(marker, colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.PreLine = NewDecorativeLine(r.context, marker, c)
	return nil
}

func (r *Rule) SetPostLine(marker, colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.PostLine = NewDecorativeLine(r.context, marker, c)
	return nil
}
