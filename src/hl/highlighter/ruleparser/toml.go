package ruleparser

import (
	"github.com/BurntSushi/toml"
	"github.com/omakoto/hl2/src/hl"
	"github.com/omakoto/hl2/src/hl/rules"
	"github.com/omakoto/hl2/src/hl/util"
)

type SingleRule struct {
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
	Rules []SingleRule `toml:"rule"`
}

func ParseFile(context hl.Context, filename string) ([]*rules.Rule, error) {

	var r RuleFile
	util.Debugf("Reading rules from '%s'...\n", filename)

	_, err := toml.DecodeFile(filename, &r)
	if err != nil {
		return nil, err
	}

	util.Dump("Rules=", r)

	var ret = make([]*rules.Rule, 0)

	for _, ir := range r.Rules {
		or, err := ir.ToRule(context)
		if err != nil {
			return nil, err
		}
		ret = append(ret, or)
	}
	return ret, err
}
