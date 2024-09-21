package main

import (
	"testing"
)

func TestRight(t *testing.T) {
	var tests = []struct {
		t string // input text
		l int    // length
		o string // output
	}{
		{"aBc", 1, "c"},
		{"", 1, ""},
		{"4321", 2, "21"},
		{"Hugo", 9, "Hugo"},
		{"This is the only universe we know", 7, "we know"},
	}
	for _, test := range tests {
		a := Right(test.t, test.l)
		if a != test.o {
			t.Errorf("Right(%q,%q) = %v; want %v", test.t, test.l, a, test.o)
		}
	}
}

//func changeParagraphs(oldParagraphState string, newParagraphState string, refreshInner bool) (resultLines []string) {

func TestChangeParagraphs(t *testing.T) {
	var tests = []struct {
		old     string // old paragraph state
		new     string // new paragraph state
		refresh bool   // refresh inner
		result  []string
	}{
		{"", "L", false, []string{"<ul>"}},
		{"L", "L", false, []string{}},
		{"L", "L", true, []string{"</ul>", "<ul>"}},
		{"P", "PNL", true, []string{"<ol>", "<ul>"}},
		{"NL", "", true, []string{"</ul>", "</ol>"}},
		{"PLLL", "PLN", true, []string{"</ul>", "</ul>", "<ol>"}},
		{"", "", true, []string{}},
	}
	for _, test := range tests {
		a := changeParagraphs(test.old, test.new, test.refresh)
		if len(a) != len(test.result) {
			t.Errorf("changeParagraphs(%q,%q,%t) = %v (len=%q); want %v (len=%q)", test.old, test.new, test.refresh, a, len(a), test.result, len(test.result))
		} else {
			for i, v := range a {
				if v != test.result[i] {
					t.Errorf("changeParagraphs(%q,%q,%t) = %v; want %v", test.old, test.new, test.refresh, a, test.result)
				}
			}
		}
	}
}

// Testing parseLink1(text) (result string)
func TestParseLink1(t *testing.T) {
	var tests = []struct {
		text string // old paragraph state
		html string // new paragraph state
	}{
		{"", ""},
		{"example.org", "<a href=\"example.org\" >example.org</a>"},
		{"click|example.org", "<a href=\"example.org\" >click</a>"},
		{"c|", "<a href=\"\" >c</a>"},
		{"|e", "<a href=\"e\" >e</a>"},
		{"|", "<a href=\"\" ></a>"},
	}
	for _, test := range tests {
		a := parseLink1(test.text)
		if len(a) != len(test.html) {
			t.Errorf("parseLink1(%q) = %q; want %q", test.text, a, test.html)
		}
	}
}

// Testing parseLink2(text) (result string)
func TestParseLink2(t *testing.T) {
	var tests = []struct {
		text string // old paragraph state
		html string // new paragraph state
	}{
		{"", ""},
		{"example.org", "[example.org)"},
		{"click](example.org", "<a href=\"example.org\" >click</a>"},
		{"c](", "<a href=\"\" >c</a>"},
		{"](e", "<a href=\"e\" >e</a>"},
		{"](", "<a href=\"\" ></a>"},
	}
	for _, test := range tests {
		a := parseLink2(test.text)
		if len(a) != len(test.html) {
			t.Errorf("parseLink2(%q) = %q; want %q", test.text, a, test.html)
		}
	}
}

/* Tests

Message("", 0, "D", "---- Test Line Parsing ----")

fmt.Println(parseLine("@test:OK", ""))
fmt.Println(parseLine("@ASWSG-VAR:$@", ""))
fmt.Println(parseLine("@FOO:foo", ""))
fmt.Println(parseLine("$BAA:baa", ""))
fmt.Println(parseLine("= Welcome", ""))
fmt.Println(parseLine("== To the Future =", ""))
fmt.Println(parseLine("", ""))
fmt.Println(parseLine("Bla bla", ""))

*/

// EOF
