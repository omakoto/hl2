package highlighter

import (
	"errors"
	"github.com/BurntSushi/toml"
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
	or := h.NewRule()

	or.SetShow(fr.Show)
	or.SetHide(fr.Hide)
	or.SetStop(fr.Stop)

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
	or.SetNextState(fr.NextState)
	or.SetStates(fr.States)

	// Colors
	err = or.SetMatchColors(fr.Colors)
	if err != nil {
		return err
	}

	// Line colors
	err = or.SetLineColors(fr.LineColors)
	if err != nil {
		return err
	}

	// Pre/post lines
	if fr.PreLine != "" {
		err = or.SetPreLine(fr.PreLine, fr.PreLineColors)
		if err != nil {
			return err
		}
	}

	if fr.PostLine != "" {
		err = or.SetPostLine(fr.PostLine, fr.PostLineColors)
		if err != nil {
			return err
		}
	}

	// After / before
	if fr.Hide {
		if fr.After > 0 || fr.Before > 0 {
			return errors.New("hidden rules can't have after/before")
		}
	}
	or.SetAfter(h.DefaultAfter())
	or.SetBefore(h.DefaultBefore())
	if fr.After > 0 {
		or.after = fr.After
	}
	if fr.Before > 0 {
		or.before = fr.Before
	}
	return nil
}
