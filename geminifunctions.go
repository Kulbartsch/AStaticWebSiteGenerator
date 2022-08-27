package main

import "strings"

// parseGeminiLink parses a gemini style link "=> link description text".
// The "=>" is optional (assumed to be already removed if missing).
func parseGeminiLink(geminiLink string) (link string, description string) {
	if len(geminiLink) < 3 {
		return "", ""
	}
	var gemLnk string
	if strings.HasPrefix(geminiLink, "=>") {
		gemLnk = WhiteSpaceTrim(geminiLink[2:])
	} else {
		gemLnk = WhiteSpaceTrim(geminiLink)
	}
	if len(geminiLink) == 0 {
		return "", ""
	}
	i := strings.IndexAny(gemLnk, " \t")
	if i == -1 {
		link = gemLnk
		description = ""
	} else {
		link = WhiteSpaceTrim(gemLnk[:i])
		description = WhiteSpaceTrim(gemLnk[i:])
	}
	if len(description) == 0 {
		description = link
	}
	return
}
