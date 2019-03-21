package matcher

import (
	"bytes"
	"errors"
)

func preProcess(pattern *string, negate *bool) error {
	i := 0
	if len(*pattern) <= i || (*pattern)[i] != '{' {
		return nil
	}
	removeSpaces := false
LOOP:
	for {
		i++
		if len(*pattern) <= i {
			return errors.New("unterminated prefix in '" + *pattern + "'")
		}
		switch (*pattern)[i] {
		case '}':
			break LOOP
		case '!':
			*negate = true
		case '#':
			removeSpaces = true
		default:
			return errors.New("unknown prefix in '" + *pattern + "'")
		}
	}
	i++
	*pattern = (*pattern)[i:]
	if removeSpaces {
		var err error
		*pattern, err = removeExtras(*pattern)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeExtras removes unescaped spaces from a string.
func removeExtras(s string) (string, error) {
	var buf bytes.Buffer

	for i := 0; i < len(s); i++ {
		r := s[i]
		if r == ' ' {
			continue
		}
		if r == '\\' {
			i++
			if i >= len(s) {
				return "", errors.New("pattern termminated with escape")
			}
			if s[i] != ' ' {
				buf.WriteByte(r)
			}
			buf.WriteByte(s[i])
			continue
		}
		buf.WriteByte(s[i])
	}
	return buf.String(), nil
}
