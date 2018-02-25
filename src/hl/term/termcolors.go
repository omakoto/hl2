package term

import "github.com/omakoto/hl2/src/hl/colors"

type RenderedColors struct {
	colors *colors.Colors
	fgCode []byte
	bgCode []byte

	next *RenderedColors
}

func (r *RenderedColors) String() string {
	return r.colors.String()
}

func NewRenderedColors(t Term, c *colors.Colors) *RenderedColors {
	fg := t.renderFg(c.Fg(), c.Attributes())
	bg := t.renderBg(c.Bg())
	return &RenderedColors{
		colors: c,
		fgCode: fg,
		bgCode: bg,
	}
}

func (r *RenderedColors) Colors() *colors.Colors {
	return r.colors
}

func (r *RenderedColors) FgCode() []byte {
	if len(r.fgCode) > 0 {
		return r.fgCode
	}
	if r.next != nil {
		return r.next.FgCode()
	}
	return []byte("")
}

func (r *RenderedColors) BgCode() []byte {
	if len(r.bgCode) > 0 {
		return r.bgCode
	}
	if r.next != nil {
		return r.next.BgCode()
	}
	return []byte("")
}

func (r *RenderedColors) SetNext(next *RenderedColors) {
	r.next = next
}
