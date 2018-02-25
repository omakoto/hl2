package term

import (
	"bytes"
	"github.com/omakoto/hl2/src/hl/colors"
	"strconv"
)

var (
	EmptyBytes = []byte("")
	CsiStart   = []byte("\x1b[")
	CsiEnd     = []byte("m")
	CsiReset   = []byte("\x1b[0m")
)

type Term interface {
	Width() int

	CsiReset() []byte

	addColor(b *bytes.Buffer, c colors.Color, base int)

	renderFg(c colors.Color, attrs colors.Attribute) []byte
	renderBg(c colors.Color) []byte
}

var _ = Term((*DumbTerm)(nil))
var _ = Term((*ConsoleTerm)(nil))
var _ = Term((*Rgb8Term)(nil))
var _ = Term((*Rgb24Term)(nil))

type DumbTerm struct {
}

func (*DumbTerm) Width() int {
	return DefaultTermWidth
}

func (*DumbTerm) CsiReset() []byte {
	return EmptyBytes
}

func (*DumbTerm) renderFg(c colors.Color, attrs colors.Attribute) []byte {
	return EmptyBytes
}

func (*DumbTerm) renderBg(c colors.Color) []byte {
	return EmptyBytes
}

func (*DumbTerm) addColor(b *bytes.Buffer, c colors.Color, base int) {
}

func Color256ToColor8(r, g, b uint8) uint8 {
	r5 := uint8(int(r) * 5 / 255)
	g5 := uint8(int(g) * 5 / 255)
	b5 := uint8(int(b) * 5 / 255)
	return 16 + 36*r5 + 6*g5 + b5
}

func max8(x, y uint8) uint8 {
	if x > y {
		return x
	}
	return y
}

func Color256ToIndex(r, g, b uint8) uint8 {
	max := max8(r, max8(g, b))
	if max == 0 {
		return 0
	}
	half := max / 2
	var ret uint8 = 0

	if r >= half {
		ret += 1
	}
	if g >= half {
		ret += 2
	}
	if b >= half {
		ret += 4
	}
	return ret
}

func addCsiAttributeCode(buffer *bytes.Buffer, attrs colors.Attribute) {
	if attrs == colors.NoAttributes {
		return
	}
	first := true
	add := func(flag colors.Attribute, ch rune) {
		if (attrs & flag) != 0 {
			if !first {
				buffer.WriteRune(';')
			}
			first = false
			buffer.WriteRune(ch)
		}
	}
	add(colors.Intense, '1')
	add(colors.Italic, '3')
	add(colors.Underline, '4')
	add(colors.Strike, '9')
	add(colors.Faint, '2')
}

type ConsoleTerm struct {
	width int
}

func (t *ConsoleTerm) Width() int {
	return t.width
}

func (t *ConsoleTerm) CsiReset() []byte {
	return CsiReset
}

func (t *ConsoleTerm) addColor(b *bytes.Buffer, c colors.Color, base int) {
	if !c.IsNone() {
		var index uint8
		if c.IsIndex() {
			index = c.Index()
		} else if c.IsRgb() {
			index = Color256ToIndex(c.R(), c.B(), c.G())
		}
		b.WriteString(strconv.Itoa(base + int(index)))
	}
}

func renderFgInner(t Term, c colors.Color, attrs colors.Attribute) []byte {
	if c == colors.NoColor && attrs == colors.NoAttributes {
		return EmptyBytes
	}
	b := bytes.Buffer{}
	b.Write(CsiStart)
	addCsiAttributeCode(&b, attrs)
	if c != colors.NoColor && attrs != colors.NoAttributes {
		b.WriteByte(';')
	}
	t.addColor(&b, c, 30)
	b.Write(CsiEnd)
	return b.Bytes()
}

func renderBgInner(t Term, c colors.Color) []byte {
	if c == colors.NoColor {
		return EmptyBytes
	}
	b := bytes.Buffer{}
	b.Write(CsiStart)
	t.addColor(&b, c, 40)
	b.Write(CsiEnd)
	return b.Bytes()
}

func (t *ConsoleTerm) renderFg(c colors.Color, attrs colors.Attribute) []byte {
	return renderFgInner(t, c, attrs)
}

func (t *ConsoleTerm) renderBg(c colors.Color) []byte {
	return renderBgInner(t, c)
}

type Rgb8Term struct {
	width int
}

func (t *Rgb8Term) Width() int {
	return t.width
}

func (*Rgb8Term) CsiReset() []byte {
	return CsiReset
}

func (t *Rgb8Term) addColor(b *bytes.Buffer, c colors.Color, base int) {
	if c.IsIndex() {
		index := c.Index()
		b.WriteString(strconv.Itoa(base + int(index)))
		return
	}
	if c.IsRgb() {
		b.WriteString(strconv.Itoa(base + 8))
		b.WriteString(";5;")
		b.WriteString(strconv.Itoa(int(Color256ToColor8(c.R(), c.G(), c.B()))))
	}
}

func (t *Rgb8Term) renderFg(c colors.Color, attrs colors.Attribute) []byte {
	return renderFgInner(t, c, attrs)
}

func (t *Rgb8Term) renderBg(c colors.Color) []byte {
	return renderBgInner(t, c)
}

type Rgb24Term struct {
	width int
}

func (t *Rgb24Term) Width() int {
	return t.width
}

func (*Rgb24Term) CsiReset() []byte {
	return CsiReset
}

func (t *Rgb24Term) addColor(b *bytes.Buffer, c colors.Color, base int) {
	if c.IsIndex() {
		index := c.Index()
		b.WriteString(strconv.Itoa(base + int(index)))
		return
	}
	if c.IsRgb() {
		b.WriteString(strconv.Itoa(base + 8))
		b.WriteString(";2;")
		b.WriteString(strconv.Itoa(int(c.R())))
		b.WriteString(";")
		b.WriteString(strconv.Itoa(int(c.G())))
		b.WriteString(";")
		b.WriteString(strconv.Itoa(int(c.B())))
	}
}

func (t *Rgb24Term) renderFg(c colors.Color, attrs colors.Attribute) []byte {
	return renderFgInner(t, c, attrs)
}

func (t *Rgb24Term) renderBg(c colors.Color) []byte {
	return renderBgInner(t, c)
}

func NewDumbTerm() *DumbTerm {
	return &DumbTerm{}
}

func NewConsoleTerm(width int) *ConsoleTerm {
	return &ConsoleTerm{width: width}
}

func NewRgb8Term(width int) *Rgb8Term {
	return &Rgb8Term{width: width}
}

func NewRgb24Term(width int) *Rgb24Term {
	return &Rgb24Term{width: width}
}
