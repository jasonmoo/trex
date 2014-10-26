package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

const (
	Flag1 uint64 = 1 << iota
	Flag2
	Flag3
)

var (
	TestWords = map[string]uint64{
		"in the evening, where the rain grows and the heart knows.  yes in the evening.  ": Flag1 | Flag2 | Flag3,

		strings.Repeat("X", 5):  Flag2,
		strings.Repeat("X", 10): Flag2,
		strings.Repeat("X", 20): Flag2,
		strings.Repeat("X", 40): Flag2,
		"three flags":           Flag1 | Flag2 | Flag3,
		"flag":                  Flag1,
		"Flag":                  Flag2,

		"\n": 0,
		" ":  0,
	}

	root *Node
)

func TestAddNodeAndWalk(t *testing.T) {

	root = NewNode()

	var i int
	for term, flags := range TestWords {
		t := NewToken(term, i)
		t.SetFlags(flags)
		i++
		root.Add(t)
	}

	root.Walk(func(n *Node) error {

		// fmt.Printf("%#v\n", n)

		for _, token := range n.Tokens {
			flags, exists := TestWords[token.Text]
			if !exists {
				t.Errorf("Expected exists, got not exists for %s", token.Text)
			} else if flags != token.Flags {
				t.Errorf("Expected %s, got %s on flags for %s", strconv.FormatUint(flags, 2), strconv.FormatUint(token.Flags, 2), token.Text)
			}
		}

		return nil
	})

}

// func TestWaht(t *testing.T) {

// 	var line = "three fLAgs have flagged \nin the fucking Flag "

// 	lex := NewLexer(strings.NewReader(line), root)
// 	lex.CaseInsensitive = true
// 	lex.EmitUnmatchedTokens = true

// 	for lex.Lex() {
// 		fmt.Printf("%#v\n", lex.Token())
// 	}

// 	if lex.Error() != nil {
// 		log.Fatal(lex.Error())
// 	}

// 	os.Exit(0)

// }

func TestSearch(t *testing.T) {

	var (
		line     = "three flags have flagged \nin the fucking Flag "
		expected = []*Token{
			&Token{Matched: true, Position: 1, Text: "three flags", Flags: 0x7},
			&Token{Matched: true, Position: 4, Text: " ", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "h", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "a", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "v", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "e", Flags: 0x0},
			&Token{Matched: true, Position: 4, Text: " ", Flags: 0x0},
			&Token{Matched: true, Position: 2, Text: "flag", Flags: 0x1},
			&Token{Matched: false, Position: 0, Text: "g", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "e", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "d", Flags: 0x0},
			&Token{Matched: true, Position: 4, Text: " ", Flags: 0x0},
			&Token{Matched: true, Position: 3, Text: "\n", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "i", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "n", Flags: 0x0},
			&Token{Matched: true, Position: 4, Text: " ", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "t", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "h", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "e", Flags: 0x0},
			&Token{Matched: true, Position: 4, Text: " ", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "f", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "u", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "c", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "k", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "i", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "n", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "g", Flags: 0x0},
			&Token{Matched: true, Position: 4, Text: " ", Flags: 0x0},
			&Token{Matched: true, Position: 6, Text: "Flag", Flags: 0x2},
		}
	)

	lex := NewLexer(strings.NewReader(line), root)
	// lex.CaseInsensitive = true
	lex.EmitUnmatchedTokens = true

	for _, expectedToken := range expected {
		if !lex.Lex() {
			t.Error("Unable to lex enough tokens from string for test")
			break
		}
		token := lex.Token()

		if expectedToken.Text != token.Text {
			t.Errorf("expected text: %v, got: %v\n", expectedToken.Text, token.Text)
		} else if expectedToken.Matched != token.Matched {
			t.Errorf("expected matched: %v, got: %v\n", expectedToken.Matched, token.Matched)
		} else if expectedToken.Flags != token.Flags {
			t.Errorf("expected flags: %v, got: %v\n", expectedToken.Flags, token.Flags)
		}
	}

	if lex.Error() != nil {
		t.Error(lex.Error())
	}

}

func TestSearchInsensitive(t *testing.T) {

	var (
		line     = "three fLAgs have flagged \nin the fucking Flag "
		expected = []*Token{
			&Token{Matched: true, Position: 0, Text: "three flags", Flags: 0x7},
			&Token{Matched: true, Position: 7, Text: " ", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "h", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "a", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "v", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "e", Flags: 0x0},
			&Token{Matched: true, Position: 7, Text: " ", Flags: 0x0},
			&Token{Matched: true, Position: 5, Text: "flag", Flags: 0x1},
			&Token{Matched: false, Position: 0, Text: "g", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "e", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "d", Flags: 0x0},
			&Token{Matched: true, Position: 7, Text: " ", Flags: 0x0},
			&Token{Matched: true, Position: 1, Text: "\n", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "i", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "n", Flags: 0x0},
			&Token{Matched: true, Position: 7, Text: " ", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "t", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "h", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "e", Flags: 0x0},
			&Token{Matched: true, Position: 7, Text: " ", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "f", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "u", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "c", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "k", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "i", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "n", Flags: 0x0},
			&Token{Matched: false, Position: 0, Text: "g", Flags: 0x0},
			&Token{Matched: true, Position: 7, Text: " ", Flags: 0x0},
			&Token{Matched: true, Position: 6, Text: "Flag", Flags: 0x2},
		}
	)

	lex := NewLexer(strings.NewReader(line), root)
	lex.CaseInsensitive = true
	lex.EmitUnmatchedTokens = true

	for _, expectedToken := range expected {
		if !lex.Lex() {
			t.Error("Unable to lex enough tokens from string for test")
			break
		}
		token := lex.Token()

		if expectedToken.Text != token.Text {
			t.Errorf("expected text: %v, got: %v\n", expectedToken.Text, token.Text)
		} else if expectedToken.Matched != token.Matched {
			t.Errorf("expected matched: %v, got: %v\n", expectedToken.Matched, token.Matched)
		} else if expectedToken.Flags != token.Flags {
			t.Errorf("expected flags: %v, got: %v\n", expectedToken.Flags, token.Flags)
		}
	}

	if lex.Error() != nil {
		t.Error(lex.Error())
	}

}

func BenchmarkLexerAllHit(b *testing.B) {

	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	words := NewNode()
	words.Add(NewToken("\n", 0))
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for i := 0; scanner.Scan(); i++ {
		words.Add(NewToken(scanner.Text(), i))
	}

	file.Seek(0, os.SEEK_SET)

	lex := NewLexer(bytes.NewReader(data), words)

	b.ResetTimer()

	var i int
	for ; lex.Lex(); i++ {
	}

	b.N = i

}

func BenchmarkLexerAllMiss(b *testing.B) {

	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	lex := NewLexer(bytes.NewReader(data), NewNode())

	b.ResetTimer()

	var i int
	for ; lex.Lex(); i++ {
	}

	b.N = i

}

func BenchmarkSearch5(b *testing.B) {

	term := strings.Repeat("X", 5)

	for i := 0; i < b.N; i++ {

		root.Search(term)

	}

}
func BenchmarkSearch10(b *testing.B) {

	term := strings.Repeat("X", 10)

	for i := 0; i < b.N; i++ {

		root.Search(term)

	}

}

func BenchmarkSearch20(b *testing.B) {

	term := strings.Repeat("X", 20)

	for i := 0; i < b.N; i++ {

		root.Search(term)

	}

}
func BenchmarkSearch40(b *testing.B) {

	term := strings.Repeat("X", 40)

	for i := 0; i < b.N; i++ {

		root.Search(term)

	}

}

func BenchmarkSearch80(b *testing.B) {

	for i := 0; i < b.N; i++ {

		root.Search("in the evening, where the rain grows and the heart knows.  yes in the evening.  ")

	}

}

func BenchmarkSearchInsensitive5(b *testing.B) {

	term := strings.Repeat("x", 5)

	for i := 0; i < b.N; i++ {

		root.SearchInsensitive(term)

	}

}
func BenchmarkSearchInsensitive10(b *testing.B) {

	term := strings.Repeat("x", 10)

	for i := 0; i < b.N; i++ {

		root.SearchInsensitive(term)

	}

}

func BenchmarkSearchInsensitive20(b *testing.B) {

	term := strings.Repeat("x", 20)

	for i := 0; i < b.N; i++ {

		root.SearchInsensitive(term)

	}

}
func BenchmarkSearchInsensitive40(b *testing.B) {

	term := strings.Repeat("x", 40)

	for i := 0; i < b.N; i++ {

		root.SearchInsensitive(term)

	}

}

func BenchmarkSearchInsensitive80(b *testing.B) {

	for i := 0; i < b.N; i++ {

		root.SearchInsensitive("in the eveNing, where the rain grows and the heart knows.  yes in the evening.  ")

	}

}
