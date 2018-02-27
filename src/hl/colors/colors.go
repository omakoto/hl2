package colors

import (
	"bytes"
	"fmt"
)

const (
	noColor     = 0
	rgbColor    = 1
	indexOffset = 128
)

var (
	// NoColor represents a transparent color.
	NoColor = Color{}

	// EmptyColors represents a transparent Colors.
	EmptyColors = Colors{}
)

// Color represents an index color, or a RGB888 color.
type Color struct {
	index uint8
	r     uint8
	g     uint8
	b     uint8
}

// String converts a color to a debug string.
func (c *Color) String() string {
	if c.index == noColor {
		return "Color{none}"
	}
	if c.index >= indexOffset {
		return fmt.Sprintf("Color{index:%d}", c.index-indexOffset)
	}
	return fmt.Sprintf("Color{r:%d, g:%d, b:%d}", c.r, c.g, c.b)
}

// NewIndexColor creates a new indexed color [0-7].
func NewIndexColor(index uint8) Color {
	if index > 7 {
		panic(fmt.Sprintf("Color index out of range: %d", index))
	}
	return Color{index: index + indexOffset}
}

func color6to256(v uint8) uint8 {
	return uint8(int(v) * 255 / 5)
}

func newRgb216Color(r, g, b uint8) Color {
	return Color{index: rgbColor, r: color6to256(r), g: color6to256(g), b: color6to256(b)}
}

func newRgb888Color(r, g, b uint8) Color {
	return Color{index: rgbColor, r: r, g: g, b: b}
}

// IsNone returns whether its NoColor or not.
func (c *Color) IsNone() bool {
	return c.index == noColor
}

// IsIndex returns whether it's an index color or not.
func (c *Color) IsIndex() bool {
	return c.index >= indexOffset
}

// Index returns the color index. Panics if it's not an index color.
func (c *Color) Index() uint8 {
	if !c.IsIndex() {
		panic("Not index color.")
	}
	return uint8(c.index - indexOffset)
}

// IsIndex returns whether it's a RGB color or not.
func (c *Color) IsRgb() bool {
	return c.index == rgbColor
}

// Index returns the R value [0-255]. Panics if it's not a RGB color.
func (c *Color) R() uint8 {
	if !c.IsRgb() {
		panic("Not RGB color.")
	}
	return c.r
}

// Index returns the G value [0-255]. Panics if it's not a RGB color.
func (c *Color) G() uint8 {
	if !c.IsRgb() {
		panic("Not RGB color.")
	}
	return c.g
}

// Index returns the B value [0-255]. Panics if it's not a RGB color.
func (c *Color) B() uint8 {
	if !c.IsRgb() {
		panic("Not RGB color.")
	}
	return c.b
}

// Attribute attributes text attributes, such as "bold" and "italic".
type Attribute int

const (
	NoAttributes Attribute = 0
	Intense      Attribute = 1 << iota
	Italic
	Underline
	Strike
	Faint
)

func (v Attribute) String() string {
	if v == NoAttributes {
		return "Attribute{none}"
	}

	var buffer bytes.Buffer
	addFlagRune := func(flag Attribute, r rune) {
		if (v & flag) != 0 {
			buffer.WriteRune(r)
		}
	}

	buffer.WriteString("Attribute{")
	addFlagRune(Intense, 'b')
	addFlagRune(Italic, 'i')
	addFlagRune(Underline, 'u')
	addFlagRune(Strike, 's')
	addFlagRune(Faint, 'f')
	buffer.WriteString("}")

	return buffer.String()
}

// Colors is Attribute + a foreground color + a background color.
type Colors struct {
	attrs Attribute

	fg Color
	bg Color
}

// NewColors creates a new Colors.
func NewColors(fg, bg Color, attrs Attribute) Colors {
	return Colors{attrs: attrs, fg: fg, bg: bg}
}

// Attributes return the attributes of a Colors.
func (c *Colors) Attributes() Attribute {
	return c.attrs
}

// Fg returns the foreground Color of a Colors.
func (c *Colors) Fg() Color {
	return c.fg
}

// Fg returns the background Color of a Colors.
func (c *Colors) Bg() Color {
	return c.bg
}

func (c *Colors) String() string {
	var buf bytes.Buffer

	buf.WriteString("Colors{")
	if c.attrs != NoAttributes {
		buf.WriteString(c.attrs.String())
		buf.WriteString(", ")
	}
	buf.WriteString(c.fg.String())
	buf.WriteString("/")
	buf.WriteString(c.bg.String())

	buf.WriteString("}")

	return buf.String()
}
