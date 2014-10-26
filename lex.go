package trex

import (
	"io"
	"unicode"
	"unicode/utf8"
)

type (
	Lexer struct {
		EmitUnmatchedTokens bool
		CaseInsensitive     bool

		r    io.RuneReader
		root *Node

		buf []rune

		t   *Token
		err error
	}
)

func NewLexer(r io.RuneReader, n *Node) *Lexer {
	return &Lexer{r: r, root: n, buf: []rune{}}
}

func (l *Lexer) Lex() bool {

	var (
		node = l.root

		tokens map[string]*Token

		term string
		r    rune
		i    int
	)

	// first try to read the buffer
	for i, r = range l.buf {
		if next, exists := node.Runes[unicode.ToLower(r)]; exists {

			// node is consistent between loops
			node = next

			if len(node.Tokens) > 0 {
				tokens = node.Tokens
				term = string(l.buf[:i+1])
			}

		} else {
			goto RETURN
		}
	}

	for {
		r, _, l.err = l.r.ReadRune()
		if l.err != nil {
			if l.err == io.EOF {
				l.err = nil
			}
			return false
		}

		l.buf = append(l.buf, r)

		if next, exists := node.Runes[unicode.ToLower(r)]; exists {

			node = next

			if len(node.Tokens) > 0 {
				tokens = node.Tokens
				term = string(l.buf)
			}

		} else {
			goto RETURN
		}
	}

RETURN:

	if len(term) == 0 {
		term = string(l.buf[0])
	}

	if tokens == nil {
		// nothing matched so emit the rune or nil
		if l.EmitUnmatchedTokens {
			l.t = &Token{Text: term}
		} else {
			l.t = nil
		}
	} else if t, exists := tokens[term]; exists {
		l.t = t
	} else if l.CaseInsensitive {
		// grab first in map
		for _, l.t = range node.Tokens {
			break
		}
	} else {
		// token matched but wrong case
		if l.EmitUnmatchedTokens {
			l.t = &Token{Text: term}
		} else {
			l.t = nil
		}
	}

	// whatever we decided the term was, ensure that
	// it's trimmed from teh buffer

	if ct := utf8.RuneCountInString(term); len(l.buf) > ct {
		l.buf = append([]rune(nil), l.buf[ct:]...) // new copy for protect mem leaks
	} else if len(l.buf) == ct {
		l.buf = l.buf[:0]
	}

	return true

}

func (l *Lexer) Token() *Token {
	return l.t
}

func (l *Lexer) Error() error {
	return l.err
}
