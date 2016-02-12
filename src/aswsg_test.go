package main

import (
	"testing"
)


func TestStringBracketsSplit(t *testing.T) {
	var tests = []struct {
		t   string // input text
		s1  string // separator one
		s2  string
		o1  string
		o2  string
		o3  string
	}{
		{"Hallo_{{name}}!", "{{", "}}", "Hallo_", "name", "!"},
		{"Hallo_{{na{{me}}","{{", "}}", "Hallo_", "na{{me", ""},
		{"Hallo_{{name}} und {{nummer}}!","{{", "}}", "Hallo_", "name", " und {{nummer}}!"},
		{"4 Hallo_{{name))!","{{", "}}", "4 Hallo_{{name))!", "", ""},
		{"Hallo_{{}}!","{{", "}}", "Hallo_", "", "!"},
		{"{{ 6 name}}","{{", "}}", "", " 6 name", ""},
		{"Hallo_{{}}!","{{", "}}", "Hallo_", "", "!"},
		{"brave world","{{", "}}", "brave world", "", ""},
		// {"Hallo_\\{{name}}!","{{", "}}", "Hallo_\\{{name}}!", "", ""}, // Check Escape - currently not implemented // TODO implement

		{"Hallo-{{_name__))!","_", "_", "Hallo-{{", "name", "_))!"},
		{"Hallo _name_!","_", "_", "Hallo ", "name", "!"},
		{"Hallo __name__!","_", "_", "Hallo ", "", "name__!"},
	}
	for _, test := range tests {
		a1, a2, a3 := StringBracketsSplit(test.t, test.s1, test.s2, "\\")
		if a1 != test.o1 || a2 != test.o2 || a3 != test.o3  {
			t.Errorf("StringBracketsSplit(%q,%q,%q,%q) = %v, %v, %v; want %v, %v, %v", test.t, test.s1, test.s2, "\\", a1, a2, a3, test.o1, test.o2, test.o3)
		}
	}
}


func TestRight(t *testing.T) {
	var tests = []struct {
		t   string // input text
		l   int    // length
		o   string // output
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
		old       string    // old paragraph state
		new       string    // new paragraph state
		refresh   bool      // refresh inner
		result    []string
	}{
		{"", "L", false,       []string {"<ul>", } },
		{"L", "L", false,      []string { } },
		{"L", "L", true,       []string {"</ul>", "<ul>", } },
		{"P", "PNL", true,     []string {"<ol>", "<ul>", } },
		{"NL", "", true,       []string {"</ul>", "</ol>", } },
		{"PLLL", "PLN", true,  []string {"</ul>", "</ul>", "<ol>", } },
		{"", "", true,         []string { } },
	}
	for _, test := range tests {
		a := changeParagraphs(test.old, test.new, test.refresh)
		if len(a) != len(test.result) {
			t.Errorf("changeParagraphs(%q,%q,%q) = %v (len=%q); want %v (len=%q)", test.old, test.new, test.refresh, a, len(a), test.result, len(test.result))
		} else {
			for i, v := range a {
				if v != test.result[i] {
					t.Errorf("changeParagraphs(%q,%q,%q) = %v; want %v", test.old, test.new, test.refresh, a, test.result)
				}
			}
		}
	}
}

//  Testing parseLink(text) (result string)
func TestParseLink(t *testing.T) {
	var tests = []struct {
		text      string    // old paragraph state
		html      string    // new paragraph state
	}{
		{"",                               "" },
		{"example.org",                    "<a href=\"example.org\" >example.org</a>" },
		{"click|example.org",              "<a href=\"example.org\" >click</a>" },
		{"c|",                             "<a href=\"\" >c</a>" },
		{"|e",                             "<a href=\"e\" ></a>" },
		{"|",                              "<a href=\"\" ></a>" },
	}
	for _, test := range tests {
		a := parseLink(test.text)
		if len(a) != len(test.html) {
			t.Errorf("parseLink(%q) = %q; want %q", test.text, a, test.html)
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
