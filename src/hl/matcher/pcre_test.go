package matcher

import (
	"testing"

	"github.com/dlclark/regexp2"
)

func TestPcre(t *testing.T) {
	re, err := regexp2.Compile(`ab(\d+)`, 0)
	if err != nil {
		t.Errorf("Compile returned %v", err)
		return
	}

	m, err := re.FindStringMatch("xab123XYZab456")
	if err != nil {
		t.Errorf("FindStringMatch returned error: %v", err)
		return
	}
	if m == nil {
		t.Errorf("FindStringMatch returned nil")
		return
	}

	groups := m.Groups()
	if len(groups) != 2 { // groups[0] = full match, groups[1] = capture group
		t.Errorf("expected 2 groups, got %v", len(groups))
	}

	if m.String() != "ab123" {
		t.Errorf("expected match 'ab123', got %v", m.String())
	}
}

func TestPcreUtf8(t *testing.T) {
	re, err := regexp2.Compile(`[③-⑥]+`, 0)
	if err != nil {
		t.Errorf("Compile returned %v", err)
		return
	}

	m, err := re.FindStringMatch("①②③④⑤⑥⑦⑧⑨⑩")
	if err != nil {
		t.Errorf("FindStringMatch returned error: %v", err)
		return
	}
	if m == nil {
		t.Errorf("FindStringMatch returned nil")
		return
	}

	groups := m.Groups()
	if len(groups) != 1 { // only groups[0] (full match), no capture groups
		t.Errorf("expected 1 group, got %v", len(groups))
	}
}
