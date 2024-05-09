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
	if index < 0 || len(str)-1 < index+P.nchrs {
		return nil
	}
	return IMatch{str, index, index + P.nchrs}
}

func P(n int) Pattern {
	return Nchars{n}
}
