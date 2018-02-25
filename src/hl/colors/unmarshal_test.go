package colors

import "testing"

const (
	NoError = true
	Error   = false
)

func TestColors_MarshalText(t *testing.T) {
	tests := []struct {
		source   string
		expected string
		noError  bool
	}{
		{`x`, ``, Error},
		{``, `Colors{Color{none}/Color{none}}`, NoError},
		{`bred`, `Colors{Attribute{b}, Color{index:1}/Color{none}}`, NoError},
		{`500`, `Colors{Color{r:255, g:0, b:0}/Color{none}}`, NoError},
		{`010`, `Colors{Color{r:0, g:51, b:0}/Color{none}}`, NoError},
		{`002`, `Colors{Color{r:0, g:0, b:102}/Color{none}}`, NoError},
		{`/500`, `Colors{Color{none}/Color{r:255, g:0, b:0}}`, NoError},
		{`/123`, `Colors{Color{none}/Color{r:51, g:102, b:153}}`, NoError},
		{`ff0000`, `Colors{Color{r:255, g:0, b:0}/Color{none}}`, NoError},
	}
	for _, v := range tests {
		var c Colors
		err := c.UnmarshalText([]byte(v.source))

		if err != nil {
			if !v.noError {
				continue
			}
			t.Errorf("Unexpected error '%s', source='%s'", err, v.source)
			continue
		}
		if !v.noError {
			t.Errorf("Error expected, but didn't happen, source='%s'", v.source)
			continue
		}
		if c.String() != v.expected {
			t.Errorf("Source='%s', expected='%s', actual='%s'", v.source, v.expected, c.String())
		}
	}
}
