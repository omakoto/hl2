package highlighter

import (
	"github.com/omakoto/hl2/src/hl"
	"github.com/omakoto/hl2/src/hl/colors"
	"github.com/omakoto/hl2/src/hl/matcher"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/util"
)

const InitialState = "INIT"

type decorativeLine struct {
	Marker []byte
	Colors *term.RenderedColors
}

func newDecorativeLine(context hl.Context, marker string, colors *colors.Colors) *decorativeLine {
	t := context.Term()
	c := term.NewRenderedColors(t, colors)
	return &decorativeLine{
		Marker: []byte(marker),
		Colors: c,
	}
}

type Rule struct {
	context hl.Context

	matcher    matcher.Matcher
	preMatcher matcher.Matcher

	after  int
	before int

	show bool
	hide bool
	stop bool

	matchColors *term.RenderedColors
	lineColors  *term.RenderedColors

	preLine  *decorativeLine
	postLine *decorativeLine

	states    []string
	nextState string
}

func newRule(context hl.Context) *Rule {
	return &Rule{context: context}
}

func (r *Rule) isForState(state string) bool {
	if len(r.states) == 0 {
		return true
	}
	for i := 0; i < len(r.states); i++ {
		if r.states[i] == state {
			return true
		}
	}
	return false
}

func (r *Rule) SetMatcherString(pattern string) error {
	m, err := matcher.CompileWithContext(r.context, pattern)
	if err != nil {
		return err
	}
	r.matcher = m
	return nil
}

func (r *Rule) SetPreMatcherString(pattern string) error {
	m, err := matcher.CompileWithContext(r.context, pattern)
	if err != nil {
		return err
	}
	r.preMatcher = m
	return nil
}

func (r *Rule) SetBefore(n int) {
	r.before = n
}

func (r *Rule) SetAfter(n int) {
	r.after = n
}

func (r *Rule) SetShow(v bool) {
	r.show = v
}

func (r *Rule) SetHide(v bool) {
	r.hide = v
}

func (r *Rule) SetStop(v bool) {
	r.stop = v
}

func (r *Rule) SetStates(states []string) {
	r.states = states
}

func (r *Rule) SetNextState(s string) {
	r.nextState = s
}

func (r *Rule) SetMatchColorsString(colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.matchColors = term.NewRenderedColors(r.context.Term(), c)
	return nil
}

func (r *Rule) SetLineColorsString(colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.lineColors = term.NewRenderedColors(r.context.Term(), c)
	return nil
}

func (r *Rule) SetPreLineString(marker, colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.preLine = newDecorativeLine(r.context, marker, c)
	return nil
}

func (r *Rule) SetPostLineString(marker, colorsStr string) error {
	c, err := colors.FromString(colorsStr)
	if err != nil {
		return err
	}
	r.postLine = newDecorativeLine(r.context, marker, c)
	return nil
}

func (r *Rule) MustSetMatcherString(pattern string) {
	util.Must(func() error { return r.SetMatcherString(pattern) })
}

func (r *Rule) MustSetPreMatcherString(pattern string) {
	util.Must(func() error { return r.SetPreMatcherString(pattern) })
}

func (r *Rule) MustSetMatchColorsString(colorsStr string) {
	util.Must(func() error { return r.SetMatchColorsString(colorsStr) })
}

func (r *Rule) MustSetLineColorsString(colorsStr string) {
	util.Must(func() error { return r.SetLineColorsString(colorsStr) })
}

func (r *Rule) MustSetPreLineString(marker, colorsStr string) {
	util.Must(func() error { return r.SetPreLineString(marker, colorsStr) })
}

func (r *Rule) MustSetPostLineString(marker, colorsStr string) {
	util.Must(func() error { return r.SetPostLineString(marker, colorsStr) })
}
