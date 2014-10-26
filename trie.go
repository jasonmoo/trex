package main

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	Node struct {
		Runes map[rune]*Node   // []*Node
		Terms map[string]*Term // []*Term
	}
	Term struct {
		Text  string
		Flags uint64
	}
)

var NonMatchError = errors.New("No match found")

func NewTerm(term string, flags uint64) *Term {
	return &Term{Text: term, Flags: flags}
}

func (t *Term) HasFlags(flags uint64) bool {
	return t.Flags&flags == flags
}
func (t *Term) SetFlags(flags uint64) {
	t.Flags |= flags
}
func (t *Term) UnsetFlags(flags uint64) {
	t.Flags &= ^flags
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) Add(term *Term) {

	text := strings.ToLower(term.Text)
	node := n

	for _, r := range text {
		if node.Runes == nil {
			node.Runes = map[rune]*Node{r: NewNode()}
		} else if _, exists := node.Runes[r]; !exists {
			node.Runes[r] = NewNode()
		}
		node = node.Runes[r]
	}

	if node.Terms == nil {
		node.Terms = map[string]*Term{term.Text: term}
	} else {
		node.Terms[term.Text] = term
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

func (n *Node) SearchInsensitive(line string) (*Term, error) {
	return n.search(line, false)
}
func (n *Node) Search(line string) (*Term, error) {
	return n.search(line, true)
}

func (n *Node) search(line string, case_sensitive bool) (*Term, error) {

	var (
		node = n

		terms    map[string]*Term
		term_len int
	)

	for i, r := range line {

		if next, exists := node.Runes[unicode.ToLower(r)]; exists {

			node = next

			if len(node.Terms) > 0 {
				terms = node.Terms
				term_len = i + utf8.RuneLen(r)
			}

		} else {
			break
		}

	}

	if terms == nil {
		return nil, NonMatchError
	}

	if t, exists := terms[line[:term_len]]; exists {
		return t, nil
	}

	if !case_sensitive {
		// grab first in map
		for _, t := range node.Terms {
			return t, nil
		}
	}

	return nil, NonMatchError

}

// func LoadFile(path string, delimiter rune, bit_flags map[rune]uint64) (*Node, error) {

// 	if len(flags) > 64 {
// 		return nil, errors.New("Flag map must contain <= 64 entries")
// 	}

// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	root := NewNode()

// 	scanner := bufio.NewScanner(file)

// 	for scanner.Scan() {

// 		var (
// 			line  = scanner.Text()
// 			i     = strings.IndexRune(line, delimiter)
// 			term  = line[:i]
// 			flags = line[i+1:]
// 		)

// 		// fmt.Println(string(term), string(tags))
// 		t := NewTerm(term, 0)
// 		for _, r := range flags {
// 			t.SetFlags(bit_flags[r])
// 		}
// 		root.Add(t)

// 	}

// 	return root, nil

// }
