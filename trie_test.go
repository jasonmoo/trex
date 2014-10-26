package main

import (
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

	for term, flags := range TestWords {
		root.Add(NewTerm(term, flags))
	}

	root.Walk(func(n *Node) error {

		for _, term := range n.Terms {
			flags, exists := TestWords[term.Text]
			if !exists {
				t.Errorf("Expected exists, got not exists for %s", term.Text)
			} else if flags != term.Flags {
				t.Errorf("Expected %s, got %s on flags for %s", strconv.FormatUint(flags, 2), strconv.FormatUint(term.Flags, 2), term.Text)
			}
		}

		return nil
	})

}

func TestSearch(t *testing.T) {

	var (
		line     = "three flags have flagged \nin the fucking Flag "
		expected = []*Term{
			&Term{Text: "three flags", Flags: 0x7},
			&Term{Text: " ", Flags: 0x0},
			nil,
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			&Term{Text: "flag", Flags: 0x1},
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			&Term{Text: "\n", Flags: 0x0},
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			&Term{Text: "Flag", Flags: 0x2},
			&Term{Text: " ", Flags: 0x0},
		}
	)

	for _, expectedTerm := range expected {
		term, err := root.Search(line)
		if err != nil && err != NonMatchError {
			t.Error(err)
		}

		if expectedTerm == nil {
			if term != nil {
				t.Errorf("expected: %#v, got: %#v\n", expectedTerm, term)
			}
		} else if expectedTerm.Text != term.Text {
			t.Errorf("expected: %#v, got: %#v\n", expectedTerm, term)
		} else if expectedTerm.Flags != term.Flags {
			t.Errorf("expected: %#v, got: %#v\n", expectedTerm, term)
		}

		if term == nil {
			line = line[1:]
		} else {
			line = line[len(term.Text):]
		}
	}

}

func TestSearchInsensitive(t *testing.T) {

	var (
		line     = "three Three FLAGS Flagged \nin the fucking flag "
		expected = []*Term{
			nil,
			nil,
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			&Term{Text: "three flags", Flags: 0x7},
			&Term{Text: " ", Flags: 0x0},
			&Term{Text: "Flag", Flags: 0x2},
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			&Term{Text: "\n", Flags: 0x0},
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			&Term{Text: " ", Flags: 0x0},
			&Term{Text: "flag", Flags: 0x1},
		}
	)

	for _, expectedTerm := range expected {
		term, err := root.SearchInsensitive(line)
		if err != nil && err != NonMatchError {
			t.Error(err)
		}

		// fmt.Printf("%#v\n", term)

		if expectedTerm == nil {
			if term != nil {
				t.Errorf("expected: %#v, got: %#v\n", expectedTerm, term)
			}
		} else if strings.ToLower(expectedTerm.Text) != strings.ToLower(term.Text) {
			t.Errorf("expected: %#v, got: %#v\n", expectedTerm, term)
		} else if expectedTerm.Flags != term.Flags {
			t.Errorf("expected: %#v, got: %#v\n", expectedTerm, term)
		}

		if term == nil {
			line = line[1:]
		} else {
			line = line[len(term.Text):]
		}
	}

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
