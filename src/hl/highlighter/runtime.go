package highlighter

import (
	"bytes"
	"github.com/omakoto/hl2/src/hl/term"
	"github.com/omakoto/hl2/src/hl/textio"
	"github.com/omakoto/hl2/src/hl/util"
	"github.com/pborman/getopt/v2"
	"io"
)

var (
	noCrSupport = getopt.BoolLong("no-cr-aware", 0, "Don't treat CRs as line terminator too. (faster)")

	emptyBytes   = []byte("")
	hiddenMarker = []byte("---\n")
)

type matchResult struct {
	rule      *Rule
	positions [][]int
}

type colorsCache struct {
	cache []*term.RenderedColors
}

func newColorsCache() colorsCache {
	return colorsCache{cache: make([]*term.RenderedColors, 4096)}
}

func (c *colorsCache) prepare(lineByteCount int) {
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

func (c *colorsCache) applyColors(start, end int, colors *term.RenderedColors) {
	for i := start; i < end; i++ {
		prev := c.cache[i]
		if prev == nil {
			c.cache[i] = colors
		} else {
			clone := *colors
			c.cache[i] = &clone
			c.cache[i].SetNext(prev)
		}
	}
}

func (c *colorsCache) applyColorsMulti(startEnds [][]int, colors *term.RenderedColors) {
	for i := 0; i < len(startEnds); i++ {
		c.applyColors(startEnds[i][0], startEnds[i][1], colors)
	}
}

func (c *colorsCache) getFg(index int) []byte {
	if c.cache[index] != nil {
		return c.cache[index].FgCode()
	}
	return nil
}

func (c *colorsCache) getBg(index int) []byte {
	if c.cache[index] != nil {
		return c.cache[index].BgCode()
	}
	return nil
}

// Runtime defines a Highlighter execution context.
// Multiple Runtime's can be created for the same Highlighter instance.
type Runtime struct {
	h *Highlighter

	wr io.Writer

	colorsCache  colorsCache
	matchesCache []matchResult
	writeCache   bytes.Buffer

	maxBefore      int
	remainingAfter int
	numHiddenLines int

	hiddenMarkWritten bool

	beforeBuffer *util.BytesRingBuffer

	state string
}

// NewRuntime creates a new Runtime. Output will be written to wr.
func (h *Highlighter) NewRuntime(wr io.Writer) *Runtime {
	r := Runtime{h: h}

	r.wr = wr
	r.matchesCache = make([]matchResult, len(r.h.rules))

	for _, rule := range r.h.rules {
		if r.maxBefore < rule.before {
			r.maxBefore = rule.before
		}
	}
	r.beforeBuffer = util.NewStringRingBuffer(r.maxBefore)
	r.colorsCache = newColorsCache()

	return &r
}

// Finish finalizes the output.
func (r *Runtime) Finish() error {
	if r.numHiddenLines > 0 {
		err := r.maybeWriteHiddenMarker()
		if err != nil {
			return err
		}
	}
	return nil
}

// ColorReader reads from rd and applies filter on the output.
func (r *Runtime) ColorReader(rd io.Reader, callFinish bool) error {
	br := textio.NewLineReader(rd, !*noCrSupport)

	for {
		bytes, err := br.ReadLine()
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
	if callFinish {
		r.Finish()
	}
	return nil
}

func (r *Runtime) clearMatchesCache() {
	for i := 0; i < len(r.matchesCache); i++ {
		r.matchesCache[i] = matchResult{}
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
	_, err := r.wr.Write(hiddenMarker)
	return err
}

// ColorBytes applies filter on a given byte array and write to wr.
func (r *Runtime) ColorBytes(b []byte) error {
	var lineTerminator []byte
	lastIndex := len(b) - 1
	last := b[lastIndex]
	if last == '\r' || last == '\n' {
		lineTerminator = b[lastIndex : lastIndex+1]
		b = b[0:lastIndex]
	}
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
			if rule.preLine != nil {
				r.writeDecorativeLine(rule.preLine)
			}
		}
	}

	// Print body.
	w := r.writeCache
	w.Truncate(0)

	// First, apply the line colors.
	for i := numMatches - 1; i >= 0; i-- {
		rule := matches[i].rule
		if rule.lineColors != nil {
			r.colorsCache.applyColors(0, numBytes, rule.lineColors)
		}
	}
	// Then, apply the match colors.
	for i := numMatches - 1; i >= 0; i-- {
		rule := matches[i].rule
		if rule.matchColors != nil {
			r.colorsCache.applyColorsMulti(matches[i].positions, rule.matchColors)
		}
	}

	// Finally print the built line.
	lastFg := emptyBytes
	lastBg := emptyBytes
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
	w.Write(lineTerminator)

	if show || r.remainingAfter > 0 {
		r.printBody(w.Bytes())
	} else {
		r.hideLine(w.Bytes())
	}

	// Post-line
	if show {
		for i := numMatches - 1; i >= 0; i-- {
			rule := r.matchesCache[i].rule
			if rule.postLine != nil {
				r.writeDecorativeLine(rule.postLine)
			}
		}
	}

	return nil
}

func (r *Runtime) findMatches(b []byte, defaultShow bool) (matches []matchResult, show bool, after int, before int) {
	show = defaultShow

	numMatches := 0
	for i := 0; i < len(r.h.rules); i++ {
		rule := r.h.rules[i]

		if !rule.isForState(r.state) {
			continue
		}
		if rule.preMatcher != nil && rule.preMatcher.Matches(b) == nil {
			continue
		}
		m := rule.matcher.Matches(b)
		if m == nil {
			continue
		}
		util.Debugf("Matched=%s\n", rule.matcher)

		//util.Debugf("Matched=%v [%s @ %s]\n", m, rule.MatchColors, rule.LineColors)
		if rule.nextState != "" {
			r.state = rule.nextState
			util.Debugf("Next state=%s\n", r.state)
		}
		r.matchesCache[numMatches] = matchResult{rule: rule, positions: m}
		numMatches++
		if rule.hide {
			show = false
		}
		if rule.show {
			show = true
		}
		if rule.after > 0 && show {
			thisAfter := rule.after + 1 // +1 because the current line consumes 1.
			if after < thisAfter {
				after = thisAfter
			}
		}
		if before < rule.before {
			before = rule.before
		}

		if rule.stop {
			break
		}
	}
	matches = r.matchesCache[0:numMatches]
	return
}

func (r *Runtime) writeDecorativeLine(d *decorativeLine) {
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
