package matcher

import (
	"testing"

	pcre "github.com/Jemmic/go-pcre2"
)

func TestPcre(t *testing.T) {
	re, err := pcre.Compile(`ab(\d+)`, 0)
	if err != nil {
		t.Errorf("Compile returned %v", err)
	}
	m := re.NewMatcher()

	res := m.Exec([]byte("xab123XYZab456"), 0)
	if res < 0 {
		t.Errorf("Exec returned %v", res)
	}
	groups := m.Groups()
	if groups != 1 {
		t.Errorf("Groups returned %v", groups)
	}
	first := string(m.Group(0))
	if first != "ab123" {
		t.Errorf("Group(0): %v", first)
	}
	//t.Errorf("GroupIndices(0): %v", m.GroupIndices(0)) // GroupIndices(0): [1 6]
	//t.Errorf("GroupIndices(1): %v", m.GroupIndices(1)) // GroupIndices(1): [3 6]
}

func TestPcreUtf8(t *testing.T) {
	re, err := pcre.Compile(`[③-⑥]+`, pcre.UTF|pcre.NO_UTF_CHECK)
	if err != nil {
		t.Errorf("Compile returned %v", err)
	}
	m := re.NewMatcher()

	res := m.Exec([]byte("①②③④⑤⑥⑦⑧⑨⑩"), 0)
	if res < 0 {
		t.Errorf("Exec returned %v", res)
	}
	groups := m.Groups()
	if groups != 0 {
		t.Errorf("Groups returned %v", groups)
	}
	//first := string(m.Group(0))
	//if first != "ab123" {
	//	t.Errorf("Group(0): %v", first)
	//}
	//t.Errorf("GroupIndices(0): %v", m.GroupIndices(0)) // GroupIndices(0): [6 18]
}
