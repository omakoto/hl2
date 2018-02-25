package matcher

import (
	"reflect"
	"testing"
)

const (
	NoError = true
	Error   = false
)

func TestRegex_Matches(t *testing.T) {
	tests := []struct {
		pattern  string
		target   string
		result   [][]int
		flags    Flags
		compiles bool
	}{
		{"", "", [][]int{{0, 0}}, None, NoError},
		{"x", "y", nil, None, NoError},
		{"(", "y", nil, None, Error},
		{"^", "xyz", [][]int{{0, 0}}, None, NoError},
		{"^x", "xyzx", [][]int{{0, 1}}, None, NoError},
		{"x", "xyzx", [][]int{{0, 1}, {3, 4}}, None, NoError},
		{"xy", "xyzxy", [][]int{{0, 2}, {3, 5}}, None, NoError},
		{"{!}x", "abcde", [][]int{{0, 5}}, None, NoError},
		{"{!}a", "abcde", nil, None, NoError},
		{"{!#} a  b", "abcde", nil, None, NoError},
		{"{!#}a c", "a cde", [][]int{{0, 5}}, None, NoError},
		{"{!#}a\\ c", "a cde", nil, None, NoError},
		{"x(y)", "xyzxy", [][]int{{1, 2}, {4, 5}}, None, NoError},
		{"x(y)x(z)", "xyxzYYxyxz", [][]int{{1, 2}, {3, 4}, {7, 8}, {9, 10}}, None, NoError},
		{"y", "xyzXYZ", [][]int{{1, 2}, {4, 5}}, IgnoreCase, NoError},
		{"{#}x y z", "xyz", [][]int{{0, 3}}, IgnoreCase, NoError},
	}
	for _, v := range tests {
		re, err := CompileGo(v.pattern, v.flags)
		if !v.compiles {
			if re != nil {
				t.Errorf("p='%s' t='%s' -> re must be null", v.pattern, v.target)
			}
			if err == nil {
				t.Errorf("p='%s' t='%s' -> pattern not expected to compile", v.pattern, v.target)
			}
			continue
		}
		if err != nil {
			t.Errorf("p='%s' t='%s' -> pattern expected to complie, but it didn't: %s", v.pattern, v.target, err)
			continue
		}
		res := re.Matches([]byte(v.target))
		if res == nil && v.result == nil {
			continue
		}
		if res == nil || v.result == nil || !reflect.DeepEqual(res, v.result) {
			t.Errorf("p='%s' t='%s' -> result must be %+v, but was %+v", v.pattern, v.target, v.result, res)
		}
	}
}
