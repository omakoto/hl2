package highlighter

import (
	"errors"
	"github.com/omakoto/hl2/src/hl"
	"strings"
)

func simpleToRule(context hl.Context, pattern, colorsStr string) (*Rule, error) {
	rule := newRule(context)

	rule.SetShow(true)

	// Pattern
	err := rule.SetMatcherString(pattern)
	if err != nil {
		return nil, err
	}

	// Colors
	vals := strings.Split(colorsStr, "@")
	if len(vals) > 3 || len(vals[0]) > 0 {
		return nil, errors.New("Invalid pattern; too many @'s in '" + colorsStr + "', or it doesn't start with @.")
	}

	if len(vals) > 1 {
		rule.SetMatchColorsString(vals[1])
		if err != nil {
			return nil, err
		}
	}
	if len(vals) > 2 {
		rule.SetLineColorsString(vals[2])
		if err != nil {
			return nil, err
		}
	}

	return rule, nil
}
