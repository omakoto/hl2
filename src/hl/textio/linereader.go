package textio

import (
	"bufio"
	"github.com/pborman/getopt/v2"
	"io"
)

type LineReader interface {
	ReadLine() (line []byte, err error)
}

type bufioLineReader struct {
	br *bufio.Reader
}

func (b *bufioLineReader) ReadLine() (line []byte, err error) {
	return b.br.ReadBytes('\n')
}

var (
	defaultBufferSize = getopt.IntLong("read-buffer-size", 0, 4096, "Specify read buffer size.")
)

type lineReader struct {
	reader io.Reader

	// Next position to buf from buf.
	next int

	// # of bytes buf by the last buf.
	avail int

	// error returned by the last buf.
	err error

	// Buffer to store buf data. The capacity is always defaultBufferSize.
	buf []byte

	// When the current line doesn't fit in buf we use this buffer.
	line []byte
}

func NewLineReader(r io.Reader, crAware bool) LineReader {
	if crAware {
		return &lineReader{
			reader: r,
			buf:    make([]byte, *defaultBufferSize),
			line:   make([]byte, 0, *defaultBufferSize),
		}
	}
	return &bufioLineReader{bufio.NewReader(r)}
}

func findCrOrLf(data []byte, start, end int) int {
	for start < end {
		switch data[start] {
		case '\r', '\n':
			return start
		}
		start++
	}
	return -1
}

func (l *lineReader) ReadLine() ([]byte, error) {
	for {
		//util.Debugf("ReadLine: %d / %d\n", l.next, l.avail)
		if l.next >= l.avail {
			if l.err != nil {
				return nil, l.err
			}

			l.next = 0

			l.avail, l.err = l.reader.Read(l.buf)
			//util.Debugf("Read: %d byte(s) read, err=%v, data=\"%s\"\n", l.avail, l.err, string(l.buf))
			if l.avail == 0 {
				if len(l.line) > 0 {
					return l.line, l.err
				}
				return nil, l.err
			}
		}
		start := l.next
		lfPos := findCrOrLf(l.buf, start, l.avail)
		//util.Debugf("Next segment: %d - %d\n", start, lfPos)
		if lfPos >= 0 {
			l.next = lfPos + 1
			if len(l.line) > 0 {
				l.line = append(l.line, l.buf[start:lfPos]...)
				ret := l.line
				//util.Debugf("Return: \"%s\"\n", string(ret))
				l.line = l.line[0:0]
				return ret, nil
			}
			ret := l.buf[start : lfPos+1]
			//util.Debugf("Return: \"%s\"\n", string(ret))
			return ret, nil
		}
		l.next = l.avail
		l.line = append(l.line, l.buf[start:l.avail]...)
		//util.Debugf("Line buffer: \"%s\"\n", string(l.line))
	}
}
