package colors

import (
	"errors"
	"github.com/omakoto/hl2/src/hl/util"
	"regexp"
	"strings"
)

var (
	colorPat = `(?:(black|red|green|yellow|blue|magenta|cyan|white)|(\d{3})|([0-9a-f]{2}),?([0-9a-f]{2}),?([0-9a-f]{2}))`
	colorsRe = regexp.MustCompile(`^(?i)\s*(?:([bifus]*)\s*` + colorPat + `)?\s*(?:\/\s*` + colorPat + `)?\s*$`)
)

func FromString(s string) (*Colors, error) {
	c := Colors{}
	err := c.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Colors) UnmarshalText(text []byte) error {
	captures := colorsRe.FindSubmatch(text)

	if captures == nil {
		return errors.New("Invalid color spec '" + string(text) + "'")
	}

	var attrs Attribute
	attrsSpec := string(captures[1])
	setAttr := func(r rune, attr Attribute) {
		if strings.ContainsRune(attrsSpec, r) {
			attrs |= attr
		}
	}

	if attrsSpec != "" {
		setAttr('b', Intense)
		setAttr('i', Italic)
		setAttr('f', Faint)
		setAttr('u', Underline)
		setAttr('s', Strike)
	}
	c.attrs = attrs

	c.fg = parseColor(captures[2], captures[3], captures[4], captures[5], captures[6])
	c.bg = parseColor(captures[7], captures[8], captures[9], captures[10], captures[11])

	util.Debugf("Input='%s'\n", text)
	util.Dump("Matches=", captures)
	util.Dump("Colors=", *c)

	return nil
}

func parseColor(name, rgb216, r8, g8, b8 []byte) Color {
	if name != nil {
		if name[0] == 'r' { // red
			return NewIndexColor(1)
		}
		if name[0] == 'g' { // green
			return NewIndexColor(2)
		}
		if name[0] == 'y' { // yellow
			return NewIndexColor(3)
		}
		if name[0] == 'm' { // magenta
			return NewIndexColor(5)
		}
		if name[0] == 'c' { // cyan
			return NewIndexColor(6)
		}
		if name[0] == 'w' { // white
			return NewIndexColor(7)
		}
		if name[2] == 'a' { // black
			return NewIndexColor(0)
		}
		return NewIndexColor(4) // blue
	}
	if rgb216 != nil {
		return NewRgb216Color(rgb216[0]-'0', rgb216[1]-'0', rgb216[2]-'0')
	}
	if r8 == nil {
		return NoColor
	}
	hexToDec := func(b byte) uint8 {
		if '0' <= b && b <= '9' {
			return b - '0'
		}
		if 'a' <= b && b <= 'f' {
			return b - 'a' + 10
		}
		if 'A' <= b && b <= 'F' {
			return b - 'A' + 10
		}
		panic("Invalid hex")
	}
	parseHex := func(v []byte) uint8 {
		return uint8(hexToDec(v[0])*16 + hexToDec(v[1]))
	}

	return NewRgb888Color(parseHex(r8), parseHex(g8), parseHex(b8))
}
