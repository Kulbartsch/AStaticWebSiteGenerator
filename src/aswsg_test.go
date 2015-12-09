package main_test

import (
    "aswsg"
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
		{"Hallo_\\{{name}}!","{{", "}}", "Hallo_\\{{name}}!", "", ""},

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

// TODO implement "right" test
/*
fmt.Println("---- other ----")
fmt.Println("right of 'aBc': " + right("aBc", 1))
fmt.Println("right of nothing: " + right("", 1))
fmt.Println("right 2 of '4321': " + right("4321", 2))
fmt.Println("---- By! ----")
*/
