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

func expect_nil(t *testing.T, m Match, expr string) {
	if m != nil {
		t.Errorf("For `%s` a nil return value is expected", expr)
	}
}

func check_bounds(t *testing.T, p Pattern) {
	m := p.Match("a", -1)
	expect_nil(t, m, "p.Match('a',-1)")

	m = p.Match("a", 5)
	expect_nil(t, m, "p.Match('a',5)")
}

// ==============================================================================
func TestPn_pass(t *testing.T) {
	// Test P(2) match that gets two characters
	p := P(2)
	check_bounds(t, p)
	m := p.Match("test", 0)
	expect_match(t, m, "te", 0, 2, "P(2).Match('test', 0)")
}

func TestPn_fail(t *testing.T) {
	// Test P(n) match where too few characters are left
	p := P(2)
	m := p.Match("str", 2)
	expect_nil(t, m, "P(2).Match('str',2)")
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
	expect_nil(t, m, "P(-1).match('str',1)")
}

// ==============================================================================
func TestPtrue(t *testing.T) {
	p := P(true)
	check_bounds(t, p)
	m := p.Match("asdf", 0)
	expect_match(t, m, "", 0, 0, "P(true).Match('asdf',0)")
}

func TestPfalse(t *testing.T) {
	p := P(false)
	m := p.Match("asdf", 0)
	expect_nil(t, m, "P(false).Match('asdf',0)")
}

// ==============================================================================
func TestPstr_pass(t *testing.T) {
	p := P("test")
	check_bounds(t, p)
	m := p.Match("a test", 2)
	expect_match(t, m, "test", 2, 6, "P('test').Match('a test', 2)")
}

func TestPstr_fail(t *testing.T) {
	p := P("test")
	m := p.Match("a test", 0)
	expect_nil(t, m, "P('test').Match('a test', 0)")
}

// ==============================================================================
func TestPfn_pass(t *testing.T) {
	p := P(func(s string, i int) int { return 3 })
	check_bounds(t, p)
	m := p.Match("test", 0)
	expect_match(t, m, "tes", 0, 3, "P(fn...).Match('test',0)")
}

func TestPfn_fail(t *testing.T) {
	p := P(func(s string, i int) int { return 7 })
	m := p.Match("test", 0)
	expect_nil(t, m, "P(fn...).Match('test', 0)")
}

func TestPfn_fail2(t *testing.T) {
	p := P(func(s string, i int) int { return 0 })
	m := p.Match("test", 1)
	expect_nil(t, m, "P(fn...).Match('test', 1)")
}

// ==============================================================================
func TestS(t *testing.T) {
	p := S("abc")
	check_bounds(t, p)
	m := p.Match("cba", 2)
	expect_match(t, m, "a", 2, 3, "S('abc').Match('cba',2)")
	m = p.Match("cba", 1)
	expect_match(t, m, "b", 1, 2, "S('abc').Match('cba',1)")
	m = p.Match("cba", 0)
	expect_match(t, m, "c", 0, 1, "S('abc').Match('cba',0)")
	m = p.Match("q", 0)
	expect_nil(t, m, "S('q').Match('q', 0)")
}

// ==============================================================================
func TestR(t *testing.T) {
	p := R("az")
	check_bounds(t, p)
	m := p.Match("aAbz", 0)
	expect_match(t, m, "a", 0, 1, "R('az').Match('aAbz',0)")
	m = p.Match("aAbz", 2)
	expect_match(t, m, "b", 2, 3, "R('az').Match('aAbz',2)")
	m = p.Match("aAbz", 3)
	expect_match(t, m, "z", 3, 4, "R('az').Match('aAbz',3)")
	m = p.Match("aAbz", 1)
	expect_nil(t, m, "R('az').Match('aAbz',1)")
}

// ==============================================================================
func TestSOL(t *testing.T) {
	p := SOL()
	check_bounds(t, p)
	m := p.Match("test", 0)
	expect_match(t, m, "", 0, 0, "SOL().Match('test',0)")
	m = p.Match("test", 1)
	expect_nil(t, m, "SOL().Match('test',1)")
	m = p.Match("test\n123", 5)
	expect_match(t, m, "", 5, 5, "SOL().Match('test\\n123',5)")
	m = p.Match("test\r123", 5)
	expect_match(t, m, "", 5, 5, "SOL().Match('test\\r123',5)")
	m = p.Match("test\r\n123", 5)
	expect_nil(t, m, "SOL().Match('test\\r\\n123',5)")
}

// ==============================================================================
func TestEOL(t *testing.T) {
	p := EOL()
	check_bounds(t, p)
	m := p.Match("test", 4)
	expect_match(t, m, "", 4, 4, "EOL().Match('test',4)")
	m = p.Match("test", 0)
	expect_nil(t, m, "EOL().Match('test',0)")
	m = p.Match("test\n123", 4)
	expect_match(t, m, "", 4, 4, "EOL().Match('test\\n123',4)")
	m = p.Match("test\r\n123", 4)
	expect_match(t, m, "", 4, 4, "EOL().Match('test\\r\\n123',4)")
	m = p.Match("test\r\n123", 5)
	expect_nil(t, m, "EOL().Match('test\\r\\n123',5)")
}
