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

//==============================================================================

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

//==============================================================================
type Union interface {
}

func P(val Union) Pattern {
	switch v := val.(type) {
	case int:
		return IntPattern{v}
	case bool:
		return BoolPattern{v}
	case string:
		return StringPattern{v}
	default:
		return nil
	}
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
