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

type Nchars struct {
	nchrs int
}

func (P Nchars) Match(str string, index int) Match {
	isLessThanN := P.nchrs < 0
	n := abs(P.nchrs)
	if isLessThanN {
		n = n - 1
	}

	// Index is out of bounds of the string
	if index < 0 || len(str) < index {
		return nil
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

func P(n int) Pattern {
	return Nchars{n}
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
