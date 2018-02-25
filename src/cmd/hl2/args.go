package main

import (
	"github.com/omakoto/hl2/src/hl/highlighter"
	"github.com/omakoto/hl2/src/hl/highlighter/simpleparser"
	"strings"
)

var (
	defaultColors = []string{
		"@b055@/012",
		"@b550@/110",
		"@b505@/101",
		"@b511@/100",
		"@b151@/010",
	}
)

func parseArgs(h *highlighter.Highlighter, args []string, execute bool, commandTerminator string) error {
	pos := 0
	if execute {
		extractCommandLine(h, args, &pos, commandTerminator)
	}
	peek := func(nth int) string {
		next := pos + nth
		if next < len(args) {
			return args[next]
		}
		return ""
	}
	// [ Pattern [@fg-color@bg-color] [,] Pattern [@fg-color@bg-color] ] ...

	defaultColorIndex := 0

	nextPatternAndColors := func() (string, string) {
		pattern := peek(0)
		if pattern == "" {
			pattern = "^$"
		}
		pos++
		colors := ""
		if strings.HasPrefix(peek(0), "@") {
			colors = peek(0)
			pos++
		} else if pattern != "^$" {
			colors = defaultColors[defaultColorIndex%len(defaultColors)]
			defaultColorIndex++
		}
		return pattern, colors
	}

	for pos < len(args) {
		pattern, colors := nextPatternAndColors()

		if peek(0) != "," {
			h.AddSimpleRule(simpleparser.NewSimple(pattern, colors))
		} else {
			pos++
			pattern2, colors2 := nextPatternAndColors()
			h.AddSimpleRangeRules(simpleparser.NewSimple(pattern, colors), simpleparser.NewSimple(pattern2, colors2))
		}
	}
	return nil
}

func extractCommandLine(h *highlighter.Highlighter, args []string, pos *int, commandTerminator string) {
	commandLine := make([]string, 0)
	for *pos < len(args) {
		a := args[*pos]
		*pos++
		if a == commandTerminator {
			break
		}
		commandLine = append(commandLine, a)
	}
	h.SetCommandLine(commandLine)
}
