package highlighter

import (
	"fmt"
	"github.com/omakoto/hl2/src/hl"
	"github.com/omakoto/hl2/src/hl/matcher"
	"github.com/omakoto/hl2/src/hl/rules"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/util"
)

type Highlighter struct {
	term term.Term

	ignoreCase  bool
	defaultHide bool

	defaultBefore int
	defaultAfter  int

	rules []*rules.Rule
}

var _ = hl.Context((*Highlighter)(nil))

func NewHighlighter() *Highlighter {
	h := &Highlighter{}
	h.SetTerm(term.NewDefaultTerm())
	return h
}

func (h *Highlighter) Term() term.Term {
	return h.term
}

func (h *Highlighter) SetTerm(t term.Term) {
	if t == nil {
		t = term.NewDefaultTerm()
	}
	h.term = t
}

func (h *Highlighter) IgnoreCase() bool {
	return h.ignoreCase
}

func (h *Highlighter) SetIgnoreCase(ignoreCase bool) {
	h.ignoreCase = ignoreCase
}

func (h *Highlighter) DefaultHide() bool {
	return h.defaultHide
}

func (h *Highlighter) SetDefaultHide(defaultHide bool) {
	h.defaultHide = defaultHide
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

func (h *Highlighter) getRules() []*rules.Rule {
	if h.rules == nil {
		h.rules = make([]*rules.Rule, 0)
	}
	return h.rules
}

func (h *Highlighter) LoadToml(ruleFile string) error {
	return h.parseTomlFile(ruleFile)
}

func (h *Highlighter) AddRule(r *rules.Rule) {
	h.rules = append(h.getRules(), r)
}

func (h *Highlighter) NewRule() *rules.Rule {
	r := rules.NewRule(h)
	h.rules = append(h.getRules(), r)
	return r
}

func (h *Highlighter) AddSimpleRule(pattern, colorsStr string) error {
	rule, err := simpleToRule(h, pattern, colorsStr)
	if err != nil {
		return err
	}
	util.Dump("Adding simple rule: ", rule)
	rule.After = h.defaultAfter
	rule.Before = h.defaultBefore

	h.AddRule(rule)
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
	end.After = h.defaultAfter
	end.NextState = rules.InitialState
	end.States = []string{implicitState}
	end.Show = true
	h.AddRule(end)

	// Add a rule to show all lines between start and end.
	intermediate := rules.NewRule(h)
	m, _ := matcher.CompileWithContext(h, "^")
	intermediate.Matcher = m
	intermediate.States = []string{implicitState}
	intermediate.Show = true
	h.AddRule(intermediate)

	// Start rule.
	start.Before = h.defaultBefore
	start.NextState = implicitState
	start.Show = true
	h.AddRule(start)

	return nil
}
