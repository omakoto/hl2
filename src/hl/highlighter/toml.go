package highlighter

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/omakoto/hl2/src/hl/colors"
	"github.com/omakoto/hl2/src/hl/rules"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/util"
)

type FileRule struct {
	Pattern string `toml:"pattern"`
	When    string `toml:"when"`

	Colors     string `toml:"color"`
	LineColors string `toml:"line_color"`

	PreLine        string `toml:"pre_line"`
	PreLineColors  string `toml:"pre_line_color"`
	PostLine       string `toml:"post_line"`
	PostLineColors string `toml:"post_line_color"`

	Show bool `toml:"show"`
	Hide bool `toml:"hide"`
	Stop bool `toml:"stop"`

	NextState string   `toml:"next_state"`
	States    []string `toml:"states"`

	After  int `toml:"after"`
	Before int `toml:"before"`
}

type RuleFile struct {
	Rules []FileRule `toml:"rule"`
}

func (h *Highlighter) parseTomlFile(filename string) error {

	var r RuleFile
	util.Debugf("Reading rules from '%s'...\n", filename)

	_, err := toml.DecodeFile(filename, &r)
	if err != nil {
		return err
	}

	util.Dump("Rules=", r)

	for _, fr := range r.Rules {
		err := h.addSingleRule(&fr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Highlighter) addSingleRule(fr *FileRule) error {
	or := rules.NewRule(h)

	or.Show = fr.Show
	or.Hide = fr.Hide
	or.Stop = fr.Stop

	// Matcher
	err := or.SetMatcher(fr.Pattern)
	if err != nil {
		return err
	}

	// Prematcher
	if fr.When != "" {
		err := or.SetPreMatcher(fr.When)
		if err != nil {
			return err
		}
	}

	// States
	or.NextState = fr.NextState
	or.States = fr.States

	// Colors
	c, err := colors.FromString(fr.Colors)
	if err != nil {
		return err
	}
	or.MatchColors = term.NewRenderedColors(h.Term(), c)

	// Line colors
	c, err = colors.FromString(fr.LineColors)
	if err != nil {
		return err
	}
	or.LineColors = term.NewRenderedColors(h.Term(), c)

	// Pre/post lines
	if fr.PreLine != "" {
		c, err = colors.FromString(fr.PreLineColors)
		if err != nil {
			return err
		}
		or.PreLine = rules.NewDecorativeLine(h, fr.PreLine, c)
	}

	if fr.PostLine != "" {
		c, err = colors.FromString(fr.PostLineColors)
		if err != nil {
			return err
		}
		or.PostLine = rules.NewDecorativeLine(h, fr.PostLine, c)
	}

	// After / before
	if fr.Hide {
		if fr.After > 0 || fr.Before > 0 {
			return errors.New("hidden rules can't have after/before")
		}
	}
	or.After = h.DefaultAfter()
	or.Before = h.DefaultBefore()
	if fr.After > 0 {
		or.After = fr.After
	}
	if fr.Before > 0 {
		or.Before = fr.Before
	}

	h.AddRule(or)
	return nil
}
