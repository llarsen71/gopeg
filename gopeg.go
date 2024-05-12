// nolint:all
package gopeg

//==============================================================================

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

//==============================================================================

type Match interface {
	GetValue() string
	Start() int
	End() int
}

type IMatch struct {
	full   string
	start_ int
	end_   int
}

func (match *IMatch) GetValue() string {
	return match.full[match.start_:match.end_]
}

func (match *IMatch) Start() int {
	return match.start_
}

func (match *IMatch) End() int {
	return match.end_
}

//==============================================================================

type Pattern interface {
	Match(str string, index int) Match
	Or(...Union) Pattern
	And(...Union) Pattern
	Not(Union) Pattern
	Rep(int) Pattern
}

type BasePattern struct {
	self Pattern
}

func (P BasePattern) Match(str string, index int) Match {
	return nil
}

func Unionize(u0 Union, u []Union) []Union {
	U := make([]Union, len(u)+1)
	U[0] = u0
	for i := 1; i < len(U); i++ {
		U[i] = u[i-1]
	}
	return U
}

func (P BasePattern) Or(p ...Union) Pattern {
	u := Unionize(P.self, p)
	return Or(u...)
}

func (P BasePattern) And(p ...Union) Pattern {
	u := Unionize(P.self, p)
	return And(u...)
}

func (P BasePattern) Not(u Union) Pattern {
	return And(P.self, Not(u))
}

func (P BasePattern) Rep(n int) Pattern {
	return Rep(P.self, n)
}

// ==============================================================================
type Union interface{} // Use for duck typing

func P(val Union) Pattern {
	switch v := val.(type) {
	case int:
		return newIntPattern(v)
	case bool:
		return newBoolPattern(v)
	case string:
		return newStringPattern(v)
	case func(string, int) int:
		return newFnPattern(v)
	case Pattern:
		return v
	default:
		return nil
	}
	return nil
}

//==============================================================================

type StringPattern struct {
	BasePattern
	str string
}

func newStringPattern(str string) Pattern {
	P := new(StringPattern)
	P.self = P
	P.str = str
	return P
}

func (P StringPattern) Match(str string, index int) Match {
	// Index is out of bounds of the string
	if index < 0 || len(str) < index {
		return nil
	}

	sz := len(P.str)
	left := len(str) - index
	if left < sz || str[index:index+sz] != P.str {
		return nil
	}

	return &IMatch{str, index, index + sz}
}

//==============================================================================

type IntPattern struct {
	BasePattern
	nchrs int
}

func newIntPattern(n int) Pattern {
	P := new(IntPattern)
	P.self = P
	P.nchrs = n
	return P
}

func (P IntPattern) Match(str string, index int) Match {
	// Index is out of bounds of the string
	if index < 0 || len(str) < index {
		return nil
	}

	isLessThanN := P.nchrs < 0
	n := abs(P.nchrs)
	if isLessThanN {
		n = n - 1
	}

	nchr := len(str) - index
	if (isLessThanN && nchr > n) || (!isLessThanN && nchr < n) {
		return nil
	}

	if isLessThanN {
		return &IMatch{str, index, index + min(nchr, n)}
	}

	return &IMatch{str, index, index + P.nchrs}
}

// ==============================================================================
type BoolPattern struct {
	BasePattern
	isTrue bool
}

func newBoolPattern(isTrue bool) Pattern {
	P := new(BoolPattern)
	P.self = P
	P.isTrue = isTrue
	return P
}

func (P BoolPattern) Match(str string, index int) Match {
	// Index is out of bounds of the string
	if index < 0 || len(str) < index {
		return nil
	}
	if !P.isTrue {
		return nil
	}

	return &IMatch{str, index, index}
}

// ==============================================================================
type FnPattern struct {
	BasePattern
	fn func(string, int) int
}

func newFnPattern(fn func(string, int) int) Pattern {
	P := new(FnPattern)
	P.self = P
	P.fn = fn
	return P
}

func (P FnPattern) Match(str string, index int) Match {
	if index < 0 || len(str) < index {
		return nil
	}
	i := P.fn(str, index)
	if i < index || len(str) < i {
		return nil
	}
	return &IMatch{str, index, i}
}

//==============================================================================

type SPattern struct {
	BasePattern
	set string
}

func (P SPattern) Match(str string, index int) Match {
	if index < 0 || len(str)-1 < index {
		return nil
	}
	s := str[index]
	for i := 0; i < len(P.set); i++ {
		if P.set[i] == s {
			return &IMatch{str, index, index + 1}
		}
	}
	return nil
}

func S(set string) Pattern {
	P := new(SPattern)
	P.self = P
	P.set = set
	return P
}

//==============================================================================

type Range struct {
	from string
	to   string
}

func (R Range) inRange(str string, index int) bool {
	s := string(str[index])
	return R.from <= s && s <= R.to
}

func newRange(str string) Range {
	R := Range{string(str[0]), string(str[1])}
	return R
}

type RPattern struct {
	BasePattern
	rng []Range
}

//------------------------------------------------------------------------------

func (P RPattern) Match(str string, index int) Match {
	if index < 0 || len(str)-1 < index {
		return nil
	}
	for _, R := range P.rng {
		if R.inRange(str, index) {
			return &IMatch{str, index, index + 1}
		}
	}
	return nil
}

func R(strs ...string) Pattern {
	R := make([]Range, len(strs))
	for i := 0; i < len(R); i++ {
		R[i] = newRange(strs[i])
	}

	P := new(RPattern)
	P.self = P
	P.rng = R
	return P
}

//==============================================================================

func SOL() Pattern {
	fn := func(str string, i int) int {
		if i < 0 || len(str) <= i {
			return -1
		}
		if i == 0 {
			return 0
		}
		s := string(str[i-1])
		if s == "\n" {
			return i
		}
		if s == "\r" && string(str[i]) != "\n" {
			return i
		}
		return -1
	}

	p := FnPattern{BasePattern{}, fn}
	return p
}

//==============================================================================

func EOL() Pattern {
	fn := func(str string, i int) int {
		if i < 0 || len(str) < i {
			return -1
		}
		if i == len(str) {
			return i
		}
		if i == 0 {
			return -1
		}
		s := string(str[i])
		if s == "\r" {
			return i
		}
		if s == "\n" && string(str[i-1]) != "\r" {
			return i
		}
		return -1
	}

	p := FnPattern{BasePattern{}, fn}
	return p
}

//==============================================================================

type OrPattern struct {
	BasePattern
	p []Pattern
}

func (P OrPattern) Match(str string, index int) Match {
	// Loop over the orred patterns and return the first Match
	for i := 0; i < len(P.p); i++ {
		m := P.p[i].Match(str, index)
		if m != nil {
			return m
		}
	}

	// No matches
	return nil
}

func Or(u ...Union) Pattern {
	p := make([]Pattern, len(u))
	for i := 0; i < len(u); i++ {
		// TODO: Handle when a Pattern fails
		p[i] = P(u[i])
	}
	P := new(OrPattern)
	P.self = P
	P.p = p
	return P
}

//==============================================================================

type AndPattern struct {
	BasePattern
	p []Pattern
}

func (P AndPattern) Match(str string, index int) Match {
	var i int = index
	for k := 0; k < len(P.p); k++ {
		m := P.p[k].Match(str, i)
		if m == nil {
			return nil
		}
		i = m.End()
	}
	return &IMatch{str, index, i}
}

func And(u ...Union) Pattern {
	p := make([]Pattern, len(u))
	for i := 0; i < len(u); i++ {
		// TODO: Handle when a Pattern fails
		p[i] = P(u[i])
	}

	P := new(AndPattern)
	P.self = P
	P.p = p
	return P
}

//==============================================================================

type NotPattern struct {
	BasePattern
	p Pattern
}

func (P NotPattern) Match(str string, index int) Match {
	m := P.p.Match(str, index)
	if m != nil {
		return nil
	}
	return &IMatch{str, index, index}
}

func Not(p Union) Pattern {
	P_ := new(NotPattern)
	P_.self = P_
	P_.p = P(p)
	return P_
}

//==============================================================================

type RepPattern struct {
	BasePattern
	p Pattern
	n int
}

func (P RepPattern) AtLeast(str string, index int) Match {
	i := index
	for j := 0; ; j++ {
		m := P.p.Match(str, i)
		if m == nil {
			// We must match at least n values or this fails
			if j < P.n {
				return nil
			}
			break
		}
		i = m.End()
	}
	return &IMatch{str, index, i}
}

func (P RepPattern) AtMost(str string, index int) Match {
	i := index
	n := -P.n
	for j := 0; j < n; j++ {
		mi := P.p.Match(str, i)
		if mi == nil {
			break
		}
		i = mi.End()
	}
	// Always succeeds
	return &IMatch{str, index, i}
}

func (P RepPattern) Match(str string, index int) Match {
	if P.n < 0 {
		return P.AtMost(str, index)
	}
	return P.AtLeast(str, index)
}

func Rep(p Union, n int) Pattern {
	P_ := new(RepPattern)
	P_.self = P_
	P_.p = P(p)
	P_.n = n
	return P_
}

//==============================================================================

type PatternRef interface {
	SetPattern(p Pattern)
}

type VPattern struct {
	BasePattern
	ref string
	p   Pattern
}

func (P *VPattern) Match(str string, index int) Match {
	if P.p == nil {
		return nil
	}
	return P.p.Match(str, index)
}

func (P *VPattern) SetPattern(p Pattern) {
	P.p = p
}

func V(ref string) Pattern {
	P := new(VPattern)
	P.ref = ref
	P.self = P
	return P
}

//==============================================================================

var Whitespace Pattern = S(" \t")
var Whitespace0 Pattern = Rep(Whitespace, 0)
var Whitespace1 Pattern = Rep(Whitespace, 1)
var Alpha Pattern = R("az", "AZ")
var Digit Pattern = R("09")
var Newline Pattern = Or("\r\n", "\r", "\n")
var Quote Pattern = S(`"'`)
