package gopeg

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

func (match IMatch) GetValue() string {
	return match.full[match.start_:match.end_]
}

func (match IMatch) Start() int {
	return match.start_
}

func (match IMatch) End() int {
	return match.end_
}

//==============================================================================

type Pattern interface {
	Match(str string, index int) Match
}

//==============================================================================

type StringPattern struct {
	str string
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

	return IMatch{str, index, index + sz}
}

//==============================================================================

type IntPattern struct {
	nchrs int
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
		return IMatch{str, index, index + min(nchr, n)}
	}

	return IMatch{str, index, index + P.nchrs}
}

// ==============================================================================
type BoolPattern struct {
	isTrue bool
}

func (P BoolPattern) Match(str string, index int) Match {
	// Index is out of bounds of the string
	if index < 0 || len(str) < index {
		return nil
	}
	if !P.isTrue {
		return nil
	}

	return IMatch{str, index, index}
}

// ==============================================================================
type FnPattern struct {
	fn func(string, int) int
}

func (P FnPattern) Match(str string, index int) Match {
	if index < 0 || len(str) < index {
		return nil
	}
	i := P.fn(str, index)
	if i < index || len(str) < i {
		return nil
	}
	return IMatch{str, index, i}
}

// ==============================================================================
type Union interface{} // Use for duck typing

func P(val Union) Pattern {
	switch v := val.(type) {
	case int:
		return IntPattern{v}
	case bool:
		return BoolPattern{v}
	case string:
		return StringPattern{v}
	case func(string, int) int:
		return FnPattern{v}
	case Pattern:
		return v
	default:
		return nil
	}
	return nil
}

//==============================================================================

type SPattern struct {
	set string
}

func (P SPattern) Match(str string, index int) Match {
	if index < 0 || len(str) < index {
		return nil
	}
	s := str[index]
	for i := 0; i < len(P.set); i++ {
		if P.set[i] == s {
			return IMatch{str, index, index + 1}
		}
	}
	return nil
}

func S(str string) Pattern {
	return SPattern{str}
}

//==============================================================================

type RPattern struct {
	from string
	to   string
}

func (P RPattern) Match(str string, index int) Match {
	if index < 0 || len(str) < index {
		return nil
	}
	s := string(str[index])
	if P.from <= s && s <= P.to {
		return IMatch{str, index, index + 1}
	}
	return nil
}

func R(str string) Pattern {
	return RPattern{string(str[0]), string(str[1])}
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

	p := FnPattern{fn}
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

	p := FnPattern{fn}
	return p
}

//==============================================================================

type OrPattern struct {
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
	return OrPattern{p}
}

//==============================================================================

type AndPattern struct {
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
	return IMatch{str, index, i}
}

func And(u ...Union) Pattern {
	p := make([]Pattern, len(u))
	for i := 0; i < len(u); i++ {
		// TODO: Handle when a Pattern fails
		p[i] = P(u[i])
	}
	return AndPattern{p}
}

//==============================================================================

type NotPattern struct {
	p Pattern
}

func (P NotPattern) Match(str string, index int) Match {
	m := P.p.Match(str, index)
	if m != nil {
		return nil
	}
	return IMatch{str, index, index}
}

func Not(p Union) Pattern {
	return NotPattern{P(p)}
}

//==============================================================================

type RepPattern struct {
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
	return IMatch{str, index, i}
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
	return IMatch{str, index, i}
}

func (P RepPattern) Match(str string, index int) Match {
	if P.n < 0 {
		return P.AtMost(str, index)
	}
	return P.AtLeast(str, index)
}

func Rep(p Union, n int) Pattern {
	return RepPattern{P(p), n}
}

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
