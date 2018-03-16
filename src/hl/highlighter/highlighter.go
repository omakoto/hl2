package highlighter

import (
	"fmt"
	"github.com/omakoto/hl2/src/hl/matcher"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/util"
)

// Highlighter defines a highlighter specification.
// Use NewHighlighter() or NewHighlighterWithTerm() to create a new instance.
type Highlighter struct {
	term term.Term

	ignoreCase   bool
	defaultHide  bool
	noSkipMarker bool

	defaultBefore int
	defaultAfter  int

	rules []*Rule
}

// NewHighlighter creates a new Highlighter instance with the auto-detected Term.
func NewHighlighter() *Highlighter {
	h := &Highlighter{}
	h.term = term.NewDefaultTerm()
	return h
}

// NewHighlighter creates a new Highlighter instance with a given Term.
func NewHighlighterWithTerm(t term.Term) *Highlighter {
	h := &Highlighter{}
	h.term = t
	return h
}

func (h *Highlighter) Term() term.Term {
	return h.term
}

func (h *Highlighter) IgnoreCase() bool {
	return h.ignoreCase
}

func (h *Highlighter) SetIgnoreCase(ignoreCase bool) {
	h.ignoreCase = ignoreCase
}

func (h *Highlighter) MatcherFlags() matcher.Flags {
	if h.ignoreCase {
		return matcher.IgnoreCase
	}
	return matcher.NoFlags
}

func (h *Highlighter) DefaultHide() bool {
	return h.defaultHide
}

func (h *Highlighter) SetDefaultHide(defaultHide bool) {
	h.defaultHide = defaultHide
}

func (h *Highlighter) NoSkipMarker() bool {
	return h.noSkipMarker
}

func (h *Highlighter) SetNoSkipMarker(noSkipMarker bool) {
	h.noSkipMarker = noSkipMarker
}

func (h *Highlighter) DefaultAfter() int {
	return h.defaultAfter
}

func (h *Highlighter) SetDefaultAfter(defaultAfter int) {
	h.defaultAfter = defaultAfter
}

func (h *Highlighter) DefaultBefore() int {
	return h.defaultBefore
}

func (h *Highlighter) SetDefaultBefore(defaultBefore int) {
	h.defaultBefore = defaultBefore
}

func (h *Highlighter) getRules() []*Rule {
	if h.rules == nil {
		h.rules = make([]*Rule, 0)
	}
	return h.rules
}

func (h *Highlighter) LoadToml(ruleFile string) error {
	return h.parseTomlFile(ruleFile)
}

func (h *Highlighter) addRule(r *Rule) {
	h.rules = append(h.getRules(), r)
}

func (h *Highlighter) NewRule() *Rule {
	r := newRule(h)
	h.rules = append(h.getRules(), r)
	return r
}

func (h *Highlighter) AddSimpleRule(pattern, colorsStr string) error {
	rule, err := simpleToRule(h, pattern, colorsStr)
	if err != nil {
		return err
	}
	util.Dump("Adding simple rule: ", rule)
	rule.after = h.defaultAfter
	rule.before = h.defaultBefore

	h.addRule(rule)
	return nil
}

var rangeRuleNext = 0

func (h *Highlighter) AddSimpleRangeRules(patternStart, colorsStart, patternEnd, colorsEnd string) error {
	start, err := simpleToRule(h, patternStart, colorsStart)
	if err != nil {
		return err
	}
	end, err := simpleToRule(h, patternEnd, colorsEnd)
	if err != nil {
		return err
	}

	util.Dump("Adding simple rule range start: ", start)
	util.Dump("Adding simple rule range end: ", end)

	implicitState := fmt.Sprintf("*range-rule-%d", rangeRuleNext)
	rangeRuleNext++

	// End rule.
	end.after = h.defaultAfter
	end.SetNextState(InitialState)
	end.SetStates([]string{implicitState})
	end.SetShow(true)
	h.addRule(end)

	// Add a rule to show all lines between start and end.
	intermediate := newRule(h)
	m, _ := matcher.Compile("^", h.MatcherFlags())
	intermediate.matcher = m
	intermediate.SetStates([]string{implicitState})
	intermediate.SetShow(true)
	h.addRule(intermediate)

	// Start rule.
	start.before = h.defaultBefore
	start.SetNextState(implicitState)
	start.SetShow(true)
	h.addRule(start)

	return nil
}

func (h *Highlighter) MustAddSimpleRule(pattern, colorsStr string) {
	util.Must(func() error { return h.AddSimpleRule(pattern, colorsStr) })
}

func (h *Highlighter) MustAddSimpleRangeRules(patternStart, colorsStart, patternEnd, colorsEnd string) {
	util.Must(func() error { return h.AddSimpleRangeRules(patternStart, colorsStart, patternEnd, colorsEnd) })
}
