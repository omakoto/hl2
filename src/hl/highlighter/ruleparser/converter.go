package ruleparser

import (
	"errors"
	"github.com/omakoto/hl2/src/hl"
	"github.com/omakoto/hl2/src/hl/colors"
	"github.com/omakoto/hl2/src/hl/rules"
	"github.com/omakoto/hl2/src/hl/term"
)

func (ir *SingleRule) ToRule(context hl.Context) (*rules.Rule, error) {
	or := rules.NewRule(context)

	or.Show = ir.Show
	or.Hide = ir.Hide
	or.Stop = ir.Stop

	// Matcher
	err := or.SetMatcher(ir.Pattern)
	if err != nil {
		return nil, err
	}

	// Prematcher
	if ir.When != "" {
		err := or.SetPreMatcher(ir.When)
		if err != nil {
			return nil, err
		}
	}

	// States
	or.NextState = ir.NextState
	or.States = ir.States

	// Colors
	c, err := colors.FromString(ir.Colors)
	if err != nil {
		return nil, err
	}
	or.MatchColors = term.NewRenderedColors(context.Term(), c)

	// Line colors
	c, err = colors.FromString(ir.LineColors)
	if err != nil {
		return nil, err
	}
	or.LineColors = term.NewRenderedColors(context.Term(), c)

	// Pre/post lines
	if ir.PreLine != "" {
		c, err = colors.FromString(ir.PreLineColors)
		if err != nil {
			return nil, err
		}
		or.PreLine = rules.NewDecorativeLine(context, ir.PreLine, c)
	}

	if ir.PostLine != "" {
		c, err = colors.FromString(ir.PostLineColors)
		if err != nil {
			return nil, err
		}
		or.PostLine = rules.NewDecorativeLine(context, ir.PostLine, c)
	}

	// After / before
	if ir.Hide {
		if ir.After > 0 || ir.Before > 0 {
			return nil, errors.New("hidden rules can't have after/before")
		}
	}
	or.After = context.DefaultAfter()
	or.Before = context.DefaultBefore()
	if ir.After > 0 {
		or.After = ir.After
	}
	if ir.Before > 0 {
		or.Before = ir.Before
	}

	return or, nil
}
