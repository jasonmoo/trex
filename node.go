package trex

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	Node struct {
		Runes  map[rune]*Node
		Tokens map[string]*Token
	}
)

var NonMatchError = errors.New("No match found")

func NewNode() *Node {
	return &Node{}
}

func (n *Node) Add(token *Token) {

	// ensure text is atomic
	token.Text = string(append([]byte(nil), token.Text...))

	text := strings.ToLower(token.Text)
	node := n

	for _, r := range text {
		if node.Runes == nil {
			node.Runes = map[rune]*Node{r: NewNode()}
		} else if _, exists := node.Runes[r]; !exists {
			node.Runes[r] = NewNode()
		}
		node = node.Runes[r]
	}

	if node.Tokens == nil {
		node.Tokens = map[string]*Token{token.Text: token}
	} else {
		node.Tokens[token.Text] = token
	}

}

// depth first
func (n *Node) Walk(f func(n *Node) error) error {

	for _, node := range n.Runes {
		if err := node.Walk(f); err != nil {
			return err
		}
	}
	return f(n)

}

func (n *Node) SearchInsensitive(line string) (*Token, error) {
	return n.search(line, false)
}
func (n *Node) Search(line string) (*Token, error) {
	return n.search(line, true)
}

func (n *Node) search(line string, case_sensitive bool) (*Token, error) {

	var (
		node = n

		tokens    map[string]*Token
		token_len int
	)

	for i, r := range line {

		if next, exists := node.Runes[unicode.ToLower(r)]; exists {

			node = next

			if len(node.Tokens) > 0 {
				tokens = node.Tokens
				token_len = i + utf8.RuneLen(r)
			}

		} else {
			break
		}

	}

	if tokens == nil {
		return nil, NonMatchError
	}

	if t, exists := tokens[line[:token_len]]; exists {
		return t, nil
	}

	if !case_sensitive {
		// grab first in map
		for _, t := range node.Tokens {
			return t, nil
		}
	}

	return nil, NonMatchError

}

func LoadGrams(r io.Reader, root *Node, n, nn int) error {

	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	grams := (nn - n) + 1
	words := make([]string, grams)
	pos := 0

	for i, _ := range words {
		s.Scan()
		words[i] = s.Text()
	}

	if err := s.Err(); err != nil {
		return err
	}

	for {

		for i := 1; i <= len(words); i++ {
			token := NewToken(strings.Join(words[:i], " "), pos)
			token.Gram = i
			root.Add(token)
		}

		pos++

		words = append([]string(nil), words[1:]...)

		if s.Scan() {
			words = append(words, s.Text())
		} else if len(words) == 0 {
			break
		}

	}

	if err := s.Err(); err != nil {
		return err
	}

	return nil

}
