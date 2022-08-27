package main

import (
	"testing"
)

func Test_parseGeminiLink(t *testing.T) {
	type args struct {
		geminiLink string
	}

	type siteContextType struct {
		vars SimpleVars
	}
	var siteContext siteContextType
	// _ = siteContext.vars.SetVar("ASWSG-GEMINI-LINK", "=>")
	siteContext.vars = SimpleVars{
		"ASWSG-GEMINI-LINK": "=>",
	}

	tests := []struct {
		name            string
		args            args
		wantLink        string
		wantDescription string
	}{
		{"empty", args{""}, "", ""},
		{"standard", args{"=> link description text"}, "link", "description text"},
		{"compact", args{"=>http://test.com/"}, "http://test.com/", "http://test.com/"},
		{"withTabs", args{"=> ftp:bla.fasel 	some\ttext and so  "}, "ftp:bla.fasel", "some\ttext and so"},
		{"noArrow", args{"/test   go here  "}, "/test", "go here"},
		{"toShort", args{"=>"}, "", ""},
		{"noArrow", args{"ftp://somwhere.go	 some where  "}, "ftp://somwhere.go", "some where"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLink, gotDescription := parseGeminiLink(tt.args.geminiLink)
			if gotLink != tt.wantLink {
				t.Errorf("parseGeminiLink() gotLink = %v, want %v", gotLink, tt.wantLink)
			}
			if gotDescription != tt.wantDescription {
				t.Errorf("parseGeminiLink() gotDescription = %v, want %v", gotDescription, tt.wantDescription)
			}
		})
	}
}
