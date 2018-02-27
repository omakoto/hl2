package highlighter

import (
	"errors"
	"github.com/omakoto/hl2/src/hl"
	"github.com/omakoto/hl2/src/hl/colors"
	"github.com/omakoto/hl2/src/hl/rules"
	"github.com/omakoto/hl2/src/hl/term"
	"strings"
)

func simpleToRule(context hl.Context, pattern, colorsStr string) (*rules.Rule, error) {
	rule := rules.NewRule(context)

	rule.Show = true

	// Pattern
	err := rule.SetMatcher(pattern)
	if err != nil {
		return nil, err
	}

	// Colors
	vals := strings.Split(colorsStr, "@")
	if len(vals) > 3 || len(vals[0]) > 0 {
		return nil, errors.New("Invalid pattern; too many @'s in '" + colorsStr + "', or it doesn't start with @.")
	}

	toColors := func(spec string) (*term.RenderedColors, error) {
		if len(spec) == 0 {
			return nil, nil
		}
		c, err := colors.FromString(spec)
		if err != nil {
			return nil, err
		}
		return term.NewRenderedColors(context.Term(), c), nil
	}

	if len(vals) > 1 {
		rule.MatchColors, err = toColors(vals[1])
		if err != nil {
			return nil, err
		}
	}
	if len(vals) > 2 {
		rule.LineColors, err = toColors(vals[2])
		if err != nil {
			return nil, err
		}
	}

	return rule, nil
}
