// general utility functions

package main

// MAYBE: put this in own library

import "strings"

// WhiteSpaceTrim trims space, tabs and new lines
func WhiteSpaceTrim(in string) string {
	if len(in) == 0 {
		return ""
	}
	return strings.Trim(in, " \t\n")
}

// counts the occurence of the first character of a string at the beginning,
// and returns the first character, it's count and the string trimmed from
// leading and ending whitespaces plus first character
func firstCharCountAndTrim(line string) (firstChar string, count int, content string) {
	if len(line) == 0 {
		return "", 0, ""
	}
	firstChar = line[0:1]
	for count = 1; line[count] == firstChar[0]; count++ {
	}
	content = strings.Trim(line[count:], " \t"+firstChar)
	return firstChar, count, content
}

// ContainsOnly returns true if s contains only runes in the string only.
// An empty string s is always true. If otherwise only is empty the result is false.
func ContainsOnly(s, only string) bool {
	var isin bool
	if len(s) == 0 {
		return true
	}
	if len(only) == 0 { // no none
		return false
	}
	for _, c := range s {
		isin = false
		for _, m := range only {
			if c == m {
				isin = true
			}
		}
		if isin == false {
			return false
		}
	}
	return true
}

// Right returns the right most l char(s) of s in r
func Right(s string, l int) (r string) {
	le := len(s)
	if le == 0 || l == 0 {
		return ""
	}
	le -= l
	if le < 0 {
		le = 0
	}
	r = s[le:]
	return
}

func checkBlockModeToggle(line string) string {
	l := strings.TrimRight(line, " \t")
	if len(l) < 3 {
		return ""
	}
	blockModes := siteContext.vars.GetVal("ASWSG-ML-CODE") +
		siteContext.vars.GetVal("ASWSG-ML-CITE") +
		siteContext.vars.GetVal("ASWSG-ML-CRUDE") +
		siteContext.vars.GetVal("ASWSG-ML-COMMENT")
	if ContainsOnly(l, blockModes) && ContainsOnly(l, l[0:1]) {
		return l[0:1]
	}
	return ""
}

//EOF
