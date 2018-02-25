package highlighter

import (
	"bufio"
	"bytes"
	"github.com/omakoto/hl2/src/hl/rules"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/util"
	"io"
)

var (
	EmptyBytes   = []byte("")
	HiddenMarker = []byte("---\n")
)

type MatchResult struct {
	rule      *rules.Rule
	positions [][]int
}

type ColorsCache struct {
	cache []*term.RenderedColors
}

func newColorsCache() ColorsCache {
	return ColorsCache{cache: make([]*term.RenderedColors, 4096)}
}

func (c *ColorsCache) prepare(lineByteCount int) {
	size := cap(c.cache)
	if size < lineByteCount {
		for size < lineByteCount {
			size *= 2
		}
		c.cache = make([]*term.RenderedColors, size)
	}

	for i := 0; i < lineByteCount; i++ {
		c.cache[i] = nil
	}
}

func (c *ColorsCache) applyColors(start, end int, colors *term.RenderedColors) {
	for i := start; i < end; i++ {
		prev := c.cache[i]
		c.cache[i] = colors
		c.cache[i].SetNext(prev)
	}
}

func (c *ColorsCache) applyColorsMulti(startEnds [][]int, colors *term.RenderedColors) {
	for i := 0; i < len(startEnds); i++ {
		c.applyColors(startEnds[i][0], startEnds[i][1], colors)
	}
}

func (c *ColorsCache) getFg(index int) []byte {
	if c.cache[index] != nil {
		return c.cache[index].FgCode()
	}
	return nil
}

func (c *ColorsCache) getBg(index int) []byte {
	if c.cache[index] != nil {
		return c.cache[index].BgCode()
	}
	return nil
}

type Runtime struct {
	h *Highlighter

	wr io.Writer

	colorsCache  ColorsCache
	matchesCache []MatchResult
	writeCache   bytes.Buffer

	maxBefore      int
	remainingAfter int
	numHiddenLines int

	hiddenMarkWritten bool

	beforeBuffer *util.BytesRingBuffer

	state string
}

func (h *Highlighter) NewRuntime(wr io.Writer) *Runtime {
	r := Runtime{h: h}

	r.wr = wr
	r.matchesCache = make([]MatchResult, len(r.h.rules))

	for _, rule := range r.h.rules {
		if r.maxBefore < rule.Before {
			r.maxBefore = rule.Before
		}
	}
	r.beforeBuffer = util.NewStringRingBuffer(r.maxBefore)
	r.colorsCache = newColorsCache()

	return &r
}

func (r *Runtime) Finish() error {
	if r.numHiddenLines > 0 {
		err := r.maybeWriteHiddenMarker()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) ColorReader(rd io.Reader) error {
	br := bufio.NewReader(rd)

	for {
		bytes, err := br.ReadBytes('\n')
		if len(bytes) > 0 {
			e2 := r.ColorBytes(bytes)
			if e2 != nil {
				return e2
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	r.Finish()
	return nil
}

func (r *Runtime) clearMatchesCache() {
	for i := 0; i < len(r.matchesCache); i++ {
		r.matchesCache[i] = MatchResult{}
	}
}

func (r *Runtime) printBefore(numBefore int) error {
	util.Debugf("[printBefore: %d]\n", numBefore)
	if r.numHiddenLines > numBefore {
		err := r.maybeWriteHiddenMarker()
		if err != nil {
			return err
		}
	}

	var e error
	r.beforeBuffer.For(numBefore, func(bytes []byte) {
		_, e = r.wr.Write(bytes)
	})
	if e != nil {
		return e
	}
	r.beforeBuffer.Clear()
	return nil
}

func (r *Runtime) printBody(bytes []byte) error {
	_, err := r.wr.Write(bytes)
	if err != nil {
		return err
	}
	if r.remainingAfter > 0 {
		r.remainingAfter--
	}

	r.beforeBuffer.Clear()

	r.numHiddenLines = 0
	r.hiddenMarkWritten = false
	return nil
}

func (r *Runtime) hideLine(bytes []byte) error {
	// hidden.
	r.beforeBuffer.Add(bytes)
	r.numHiddenLines++

	util.Debugf("[Hidden: %d]\n", r.numHiddenLines)

	if r.numHiddenLines > r.maxBefore {
		r.maybeWriteHiddenMarker()
	}
	return nil
}

func (r *Runtime) maybeWriteHiddenMarker() error {
	if r.hiddenMarkWritten {
		return nil
	}
	r.hiddenMarkWritten = true
	_, err := r.wr.Write(HiddenMarker)
	return err
}

func (r *Runtime) ColorBytes(b []byte) error {
	b = bytes.TrimRight(b, "\r\n \t")
	numBytes := len(b)

	r.colorsCache.prepare(numBytes)
	r.clearMatchesCache()

	// Find the matches.
	matches, show, after, before := r.findMatches(b, !r.h.defaultHide)
	if show {
		r.remainingAfter = after
	}
	numMatches := len(matches)

	// Before
	if show {
		err := r.printBefore(before)
		if err != nil {
			return err
		}
	}

	// Pre-line

	if show {
		for i := 0; i < numMatches; i++ {
			rule := matches[i].rule
			if rule.PreLine != nil {
				r.writeDecorativeLine(rule.PreLine)
			}
		}
	}

	// Print body.
	w := r.writeCache
	w.Truncate(0)

	// First, apply the line colors.
	for i := numMatches - 1; i >= 0; i-- {
		rule := matches[i].rule
		if rule.LineColors != nil {
			r.colorsCache.applyColors(0, numBytes, rule.LineColors)
		}
	}
	// Then, apply the match colors.
	for i := numMatches - 1; i >= 0; i-- {
		rule := matches[i].rule
		if rule.MatchColors != nil {
			r.colorsCache.applyColorsMulti(matches[i].positions, rule.MatchColors)
		}
	}

	// Finally print the built line.
	lastFg := EmptyBytes
	lastBg := EmptyBytes
	for i := 0; i < numBytes; i++ {
		fg := r.colorsCache.getFg(i)
		bg := r.colorsCache.getBg(i)
		if !bytes.Equal(lastFg, fg) || !bytes.Equal(lastBg, bg) {
			w.Write(r.h.Term().CsiReset())

			w.Write(fg)
			w.Write(bg)
			lastFg = fg
			lastBg = bg
		}
		w.WriteByte(b[i])
	}
	if len(lastFg) > 0 || len(lastBg) > 0 {
		w.Write(r.h.Term().CsiReset())
	}
	w.WriteByte('\n')

	if show || r.remainingAfter > 0 {
		r.printBody(w.Bytes())
	} else {
		r.hideLine(w.Bytes())
	}

	// Post-line
	if show {
		for i := numMatches - 1; i >= 0; i-- {
			rule := r.matchesCache[i].rule
			if rule.PostLine != nil {
				r.writeDecorativeLine(rule.PostLine)
			}
		}
	}

	return nil
}

func (r *Runtime) findMatches(b []byte, defaultShow bool) (matches []MatchResult, show bool, after int, before int) {
	show = defaultShow

	numMatches := 0
	for i := 0; i < len(r.h.rules); i++ {
		rule := r.h.rules[i]

		if !rule.IsForState(r.state) {
			continue
		}
		if rule.PreMatcher != nil && rule.PreMatcher.Matches(b) == nil {
			continue
		}
		m := rule.Matcher.Matches(b)
		if m == nil {
			continue
		}
		util.Debugf("Matched=%s\n", rule.Matcher)

		//util.Debugf("Matched=%v [%s @ %s]\n", m, rule.MatchColors, rule.LineColors)
		if rule.NextState != "" {
			r.state = rule.NextState
			util.Debugf("Next state=%s\n", r.state)
		}
		r.matchesCache[numMatches] = MatchResult{rule: rule, positions: m}
		numMatches++
		if rule.Hide {
			show = false
		}
		if rule.Show {
			show = true
		}
		if rule.After > 0 && show {
			thisAfter := rule.After + 1 // +1 because the current line consumes 1.
			if after < thisAfter {
				after = thisAfter
			}
		}
		if before < rule.Before {
			before = rule.Before
		}

		if rule.Stop {
			break
		}
	}
	matches = r.matchesCache[0:numMatches]
	return
}

func (r *Runtime) writeDecorativeLine(d *rules.DecorativeLine) {
	w := r.writeCache

	fg := d.Colors.FgCode()
	bg := d.Colors.BgCode()
	w.Truncate(0)
	w.Write(fg)
	w.Write(bg)
	for i := r.h.Term().Width() / len(d.Marker); i > 0; i-- {
		w.Write(d.Marker)
	}
	if len(fg) > 0 || len(bg) > 0 {
		w.Write(r.h.Term().CsiReset())
	}
	w.WriteByte('\n')
	r.wr.Write(w.Bytes())
}
