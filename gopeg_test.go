package gopeg_test

import (
	. "gopeg"
	"testing"
)

func expect_match(t *testing.T, m Match, expected_str string, expected_start int, expected_end int, expr string) {
	if m == nil {
		t.Errorf("Match expected for `%s.Match()`. Got nil.", expr)
		return
	}
	str := m.GetValue()
	if str != expected_str {
		t.Errorf("For `%s.GetValue()` '%s' was expected. Got '%s'", expr, expected_str, str)
	}
	start := m.Start()
	if start != expected_start {
		t.Errorf("For `%s.Start()`, %d was expected. Got %d", expr, expected_start, start)
	}
	end := m.End()
	if end != expected_end {
		t.Errorf("For `%s.End()`, %d was expected. Got %d", expr, expected_end, end)
	}
}

// ==============================================================================
func TestPn_pass(t *testing.T) {
	// Test P(2) match that gets two characters
	p := P(2)
	m := p.Match("test", 0)
	expect_match(t, m, "te", 0, 2, "P(2).Match('test', 0)")
}

func TestPn_fail(t *testing.T) {
	// Test P(n) match where too few characters are left
	p := P(2)
	m := p.Match("str", 2)
	if m != nil {
		t.Errorf("For `P(2).Match('str',2)` a nil return value is expected")
	}
}

func TestPn_minus_pass(t *testing.T) {
	// Test P(-1) that matches the end of the string
	p := P(-1)
	m := p.Match("str", 3)
	expect_match(t, m, "", 3, 3, "P(-1).Match('str',3)")
}

func TestPn_minus_pass2(t *testing.T) {
	// Test P(-1) that matches the end of the string
	p := P(-4)
	m := p.Match("str", 0)
	expect_match(t, m, "str", 0, 3, "P(-4).Match('str', 0)")
}

func TestPn_minus_fail(t *testing.T) {
	// Test P(-1) that matches the end of the string
	p := P(-1)
	m := p.Match("str", 1)
	if m != nil {
		t.Errorf("`P(-1).match('str',1)` should return nil. Got %s", m.GetValue())
	}
}

// ==============================================================================
func TestPtrue(t *testing.T) {
	p := P(true)
	m := p.Match("asdf", 0)
	expect_match(t, m, "", 0, 0, "P(true).Match('asdf',0)")
}

func TestPfalse(t *testing.T) {
	p := P(false)
	m := p.Match("asdf", 0)
	if m != nil {
		t.Error("`P(false).Match('asdf',0)` should return nil")
	}
}

// ==============================================================================
func TestPstr_pass(t *testing.T) {
	p := P("test")
	m := p.Match("a test", 2)
	expect_match(t, m, "test", 2, 6, "P('test').Match('a test', 2)")
}
