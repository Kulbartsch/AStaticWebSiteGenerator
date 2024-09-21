package main

import "testing"

func TestStringBracketsSplit(t *testing.T) {
	type args struct {
		text   string
		b1     string
		b2     string
		escape string
	}
	tests := []struct {
		name  string
		args  args
		wantA string
		wantB string
		wantC string
	}{
		// TODO: Add test cases.
		{"empty-01", args{"", "", "", ""}, "", "", ""},
		{"empty-02", args{"", "(", ")", ""}, "", "", ""},
		{"empty-03", args{"", "([", "])", "\\"}, "", "", ""},
		{"unchanged-01", args{"abc", "(", ")", ""}, "abc", "", ""},
		{"unchanged-02", args{"abc", "([", "])", "\\"}, "abc", "", ""},
		{"unchanged-03", args{"abc(xyz", "(", ")", "\\"}, "abc(xyz", "", ""},
		{"unchanged-04", args{"abc([xyz", "([", "])", "\\"}, "abc([xyz", "", ""},
		{"straight-01", args{"abc(mn)xyz", "(", ")", ""}, "abc", "mn", "xyz"},
		{"straight-02", args{"abc([mn])xyz", "([", "])", ""}, "abc", "mn", "xyz"},
		{"straight-03", args{"äb𑜸c([mµ])xy😇z", "([", "])", ""}, "äb𑜸c", "mµ", "xy😇z"},
		{"edge-01", args{"(mn)xyz", "(", ")", ""}, "", "mn", "xyz"},
		{"edge-02", args{"([mn])xyz", "([", "])", ""}, "", "mn", "xyz"},
		{"edge-03", args{"abc(mn)", "(", ")", ""}, "abc", "mn", ""},
		{"edge-04", args{"abc([mn])", "([", "])", ""}, "abc", "mn", ""},
		{"edge-05", args{"abc()xyz", "(", ")", ""}, "abc", "", "xyz"},
		{"edge-06", args{"abc([])xyz", "([", "])", ""}, "abc", "", "xyz"},
		{"edge-07", args{"(mn)", "(", ")", ""}, "", "mn", ""},
		{"edge-08", args{"([mn])", "([", "])", ""}, "", "mn", ""},
		{"edge-09", args{"abc()", "(", ")", ""}, "abc", "", ""},
		{"edge-10", args{"abc([])", "([", "])", ""}, "abc", "", ""},
		{"edge-11", args{"()xyz", "(", ")", ""}, "", "", "xyz"},
		{"edge-12", args{"([])xyz", "([", "])", ""}, "", "", "xyz"},
		{"edge-13", args{"()", "(", ")", ""}, "", "", ""},
		{"edge-14", args{"([])", "([", "])", ""}, "", "", ""},
		{"symmetric-01", args{"abc_mn_xyz", "_", "_", ""}, "abc", "mn", "xyz"},
		{"symmetric-02", args{"abc__mn__xyz", "__", "__", ""}, "abc", "mn", "xyz"},
		{"symmetric-03", args{"____", "__", "__", ""}, "", "", ""},
		{"symmetric-03", args{"___", "__", "__", ""}, "___", "", ""},
		{"symmetric-04", args{"abc__mn_", "__", "__", ""}, "abc__mn_", "", ""},
		{"symmetric-05", args{"   ___  __", "__", "__", ""}, "   ", "_  ", ""},
		{"symmetric-03", args{"____", "_", "_", ""}, "", "", "__"},
		{"classic-01", args{"Hallo_{{name}}!", "{{", "}}", ""}, "Hallo_", "name", "!"},
		{"classic-02", args{"Hallo_{{na{{me}}", "{{", "}}", ""}, "Hallo_", "na{{me", ""},
		{"classic-03", args{"Hallo_{{name}} und {{nummer}}!", "{{", "}}", ""}, "Hallo_", "name", " und {{nummer}}!"},
		{"classic-04", args{"4 Hallo_{{name))!", "{{", "}}", ""}, "4 Hallo_{{name))!", "", ""},
		{"classic-05", args{"Hallo_{{}}!", "{{", "}}", ""}, "Hallo_", "", "!"},
		{"classic-06", args{"{{ 6 name}}", "{{", "}}", ""}, "", " 6 name", ""},
		{"classic-07", args{"Hallo_{{}}!", "{{", "}}", ""}, "Hallo_", "", "!"},
		{"classic-08", args{"brave world", "{{", "}}", ""}, "brave world", "", ""},
		{"classic-09", args{"Hallo-{{_name__))!", "_", "_", ""}, "Hallo-{{", "name", "_))!"},
		{"classic-10", args{"Hallo _name_!", "_", "_", ""}, "Hallo ", "name", "!"},
		{"classic-11", args{"Hallo __name__!", "_", "_", ""}, "Hallo ", "", "name__!"},
		{"why-01", args{" (that's it in **bold** type)", "**", "**", ""}, " (that's it in ", "bold", " type)"},
		// {"Hallo_\\{{name}}!","{{", "}}", ""}, "Hallo_\\{{name}}!", "", ""}, // Check Escape - currently not implemented // TODO implement
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB, gotC := StringBracketsSplit(tt.args.text, tt.args.b1, tt.args.b2, tt.args.escape)
			if gotA != tt.wantA {
				t.Errorf("StringBracketsSplit() %s gotA = %v, want %v", tt.name, gotA, tt.wantA)
			}
			if gotB != tt.wantB {
				t.Errorf("StringBracketsSplit() %s gotB = %v, want %v", tt.name, gotB, tt.wantB)
			}
			if gotC != tt.wantC {
				t.Errorf("StringBracketsSplit() %s gotC = %v, want %v", tt.name, gotC, tt.wantC)
			}
		})
	}
}
