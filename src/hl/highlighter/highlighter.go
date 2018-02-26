package highlighter

import (
	"fmt"
	"github.com/omakoto/hl2/src/hl"
	"github.com/omakoto/hl2/src/hl/highlighter/ruleparser"
	"github.com/omakoto/hl2/src/hl/highlighter/simpleparser"
	"github.com/omakoto/hl2/src/hl/matcher"
	"github.com/omakoto/hl2/src/hl/rules"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/util"
)

type Highlighter struct {
	term term.Term

	commandLine []string

	ignoreCase  bool
	defaultHide bool

	defaultBefore int
	defaultAfter  int

	rules []*rules.Rule
}

var _ = hl.Context((*Highlighter)(nil))

func NewHighlighter(t term.Term, ignoreCase, defaultHide bool, defaultBefore, defaultAfter int) *Highlighter {
	return &Highlighter{
		term:          t,
		ignoreCase:    ignoreCase,
		defaultHide:   defaultHide,
		defaultBefore: defaultBefore,
		defaultAfter:  defaultAfter,
		rules:         make([]*rules.Rule, 0),
	}
}

func (h *Highlighter) Term() term.Term {
	return h.term
}

func (h *Highlighter) IgnoreCase() bool {
	return h.ignoreCase
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

func (h *Highlighter) DefaultBefore() int {
	return h.defaultBefore
}

func (h *Highlighter) LoadToml(ruleFile string) error {
	rules, err := ruleparser.ParseFile(h, ruleFile)
	if err != nil {
		return err
	}
	h.rules = append(h.rules, rules...)
	return nil
}

func (h *Highlighter) AddRule(r *rules.Rule) {
	h.rules = append(h.rules, r)
}

func (h *Highlighter) SetCommandLine(commandLine []string) {
	h.commandLine = commandLine
}

func (h *Highlighter) CommandLine() []string {
	return h.commandLine
}

func (h *Highlighter) AddSimpleRule(simple *simpleparser.Simple) error {
	util.Dump("Adding simple rule: ", simple)
	rule, err := simple.ToRule(h)
	if err != nil {
		return err
	}
	rule.After = h.defaultAfter
	rule.Before = h.defaultBefore

	h.AddRule(rule)
	return nil
}

var rangeRuleNext = 0

func (h *Highlighter) AddSimpleRangeRules(start, end *simpleparser.Simple) error {
	util.Dump("Adding simple rule range start: ", start)
	util.Dump("Adding simple rule range end: ", end)

	implicitState := fmt.Sprintf("*range-rule-%d", rangeRuleNext)
	rangeRuleNext++

	// End rule.
	er, err := end.ToRule(h)
	if err != nil {
		return err
	}
	er.After = h.defaultAfter
	er.NextState = rules.InitialState
	er.States = []string{implicitState}
	er.Show = true
	h.AddRule(er)

	// Add a rule to show all lines between start and end.
	intermediate := rules.Rule{}
	m, _ := matcher.CompileWithContext(h, "^")
	intermediate.Matcher = m
	intermediate.States = []string{implicitState}
	intermediate.Show = true
	h.AddRule(&intermediate)

	// Start rule.
	sr, err := start.ToRule(h)
	if err != nil {
		return err
	}
	sr.Before = h.defaultBefore
	sr.NextState = implicitState
	sr.Show = true
	h.AddRule(sr)

	return nil
}
