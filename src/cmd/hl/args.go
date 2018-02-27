package main

import (
	"github.com/omakoto/hl2/src/hl/highlighter"
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

func parseArgs(h *highlighter.Highlighter, args []string, takeInput bool, argumentSeparator string) (inputArgs []string, err error) {
	pos := 0
	if takeInput {
		inputArgs = extractInputArgs(args, &pos, argumentSeparator)
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

		if peek(0) != argumentSeparator {
			h.AddSimpleRule(pattern, colors)
		} else {
			pos++
			pattern2, colors2 := nextPatternAndColors()
			h.AddSimpleRangeRules(pattern, colors, pattern2, colors2)
			h.SetDefaultHide(true) // Range patterns imply -n.
		}
	}
	return
}

func extractInputArgs(args []string, pos *int, argumentSeparator string) []string {
	commandLine := make([]string, 0)
	for *pos < len(args) {
		a := args[*pos]
		*pos++
		if a == argumentSeparator {
			break
		}
		commandLine = append(commandLine, a)
	}
	return commandLine
}
