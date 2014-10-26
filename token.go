package trex

type (
	Token struct {
		Matched  bool
		Position int
		Text     string
		Flags    uint64
	}
)

func NewToken(term string, pos int) *Token {
	return &Token{
		Matched:  true,
		Position: pos,
		Text:     term,
	}
}

func (t *Token) HasFlags(flags uint64) bool {
	return t.Flags&flags == flags
}
func (t *Token) SetFlags(flags uint64) {
	t.Flags |= flags
}
func (t *Token) UnsetFlags(flags uint64) {
	t.Flags &= ^flags
}
