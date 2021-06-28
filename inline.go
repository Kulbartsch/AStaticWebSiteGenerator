package main

import (
	"strconv"
	"strings"
)

// inline

// HTMLAttrib Tags as map
type HTMLAttrib map[string]string

func generateHTMLTagWithAttributes(tag string, openTag bool, attrib HTMLAttrib) (resultHTMLTag string) {
	// TODO check using HTML library
	var attributes string
	if len(tag) == 0 {
		return
	}
	var closeTag string
	if openTag == false {
		closeTag = "/"
	} else {
		closeTag = ""
		for k, v := range attrib {
			attributes += " "
			attributes += k + "=\"" + v + "\" "
		}
	}
	resultHTMLTag = "<" + closeTag + tag + attributes + ">"
	return
}

func generateHTMLTag(tag string, openTag bool) (resultHTMLTag string) {
	return generateHTMLTagWithAttributes(tag, openTag, HTMLAttrib{})
}

func surroundWithHTMLTagWithAttributes(tag string, s string, attrib HTMLAttrib) string {
	return generateHTMLTagWithAttributes(tag, true, attrib) + s + generateHTMLTag(tag, false)
}

func surroundWithHTMLTag(tag string, s string) string {
	return surroundWithHTMLTagWithAttributes(tag, s, HTMLAttrib{})
}

func generateTag(tagKind string, openTag bool) (resultTag string) {
	if len(tagKind) == 0 {
		return
	}
	resultTag = generateHTMLTag(paragraphTags[tagKind], openTag)
	return
}

// Proceeses the inner part of an [[text|URL]] (ASWSG-LINK-1-x) an generates a complete <a> tag.
// If text contains a "|"" (pipe) the left part is the displayed content and
// the right part is the href.
func parseLink1(text string) string {
	if len(text) == 0 {
		return ""
	}
	var link, display string
	var attrib HTMLAttrib
	i := strings.Index(text, siteContext.vars.GetVal("ASWSG-LINK-1-3"))
	if i == -1 {
		link = text
		display = link
	} else {
		display = text[:i]
		link = text[i+1:]
	}
	if display == "" {
		display = link
	}

	attrib = HTMLAttrib{"href": link}
	return surroundWithHTMLTagWithAttributes("a", display, attrib) // tag string, s string, attrib HTMLAttrib)
}

// Proceeses the inner part of an [text](link) format (ASWSG-LINK-2-x) an generates a complete <a> tag.
// If text contains no "](" (ASWSG-LINK-2-3) the link processing will be canceled
// and the complete inner text returned.
func parseLink2(text string) string {
	if len(text) == 0 || len(siteContext.vars.GetVal("ASWSG-LINK-2-1")) == 0 || len(siteContext.vars.GetVal("ASWSG-LINK-2-2")) == 0 || len(siteContext.vars.GetVal("ASWSG-LINK-2-3")) == 0 {
		return ""
	}
	var link, display string
	var attrib HTMLAttrib
	i := strings.Index(text, siteContext.vars.GetVal("ASWSG-LINK-2-3"))
	if i == -1 {
		return siteContext.vars.GetVal("ASWSG-LINK-2-1") + text + siteContext.vars.GetVal("ASWSG-LINK-2-2")
	}
	display = text[:i]
	link = text[i+2:]
	if display == "" {
		display = link
	}
	attrib = HTMLAttrib{"href": link, "rel": "external"}
	return surroundWithHTMLTagWithAttributes("a", display, attrib) // tag string, s string, attrib HTMLAttrib)
}

// Processes the inner part of an [[link]] format (ASWSG-LINK-3-x) and
// generates a numbered <a> tag referring to an index.
func parseLink3(text string) string {
	if len(text) == 0 {
		return ""
	}
	var attrib HTMLAttrib
	display := "[" + strconv.Itoa(linkIndex) + "]"
	linkIndex++
	indxLink := "#" + display
	destLink := indexedLinksType{display, text}
	indexedLinks = append(indexedLinks, destLink)
	attrib = HTMLAttrib{"href": indxLink}
	return surroundWithHTMLTagWithAttributes("a", display, attrib) // tag string, s string, attrib HTMLAttrib)
}

// StringBracketsSplit splits a string into a part before the brackets (b1) in the brackets and after the brackets (b2)
func StringBracketsSplit(text string, b1 string, b2 string, escape string) (a string, b string, c string) {
	m := strings.Index(text, b1)
	if m == -1 { // TODO: maybe Check for code rune
		return text, "", ""
	}
	n := strings.Index(text[m+1:], b2) + m + 1
	if n == -1 || n <= m { // ToDo: maybe Check for Escape rune
		return text, "", ""
	}
	return text[0:m], text[m+len(b1) : n], text[n+len(b2):]
}

func parseInLine(rawLine string) (parsedLine string) {

	didParse := false
	parsedLine = rawLine
	var t1, t2, t3 string

	// check bold
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-BOLD-1"), siteContext.vars.GetVal("ASWSG-BOLD-2"), siteContext.vars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("b", t2) + t3
	}

	// check emphasised
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-EMP-1"), siteContext.vars.GetVal("ASWSG-EMP-2"), siteContext.vars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 && Right(t1, 1) != ":" { // TOFIX: Workaround to not mess with HTML links containing "://"
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("em", t2) + t3
	}

	// check strike
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-STRIKE-1"), siteContext.vars.GetVal("ASWSG-STRIKE-2"), siteContext.vars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("del", t2) + t3
	}

	// check underline
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-UNDERL-1"), siteContext.vars.GetVal("ASWSG-UNDERL-2"), siteContext.vars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("u", t2) + t3
	}

	// check code
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-CODE-1"), siteContext.vars.GetVal("ASWSG-CODE-2"), siteContext.vars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("code", t2) + t3
	}

	// check link (1)
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-LINK-1-1"), siteContext.vars.GetVal("ASWSG-LINK-1-2"), siteContext.vars.GetVal("ASWSG-ESCAPE")) // TOFIX implement link func
	l2 := strings.Index(t2, siteContext.vars.GetVal("ASWSG-LINK-1-3"))
	if len(t2) > 0 && l2 >= 0 {
		didParse = true
		parsedLine = t1 + parseLink1(t2) + t3
	}

	// check link (2)
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-LINK-2-1"), siteContext.vars.GetVal("ASWSG-LINK-2-2"), siteContext.vars.GetVal("ASWSG-ESCAPE")) // TOFIX implement link func
	l2 = strings.Index(t2, siteContext.vars.GetVal("ASWSG-LINK-2-3"))
	if len(t2) > 0 && l2 >= 0 {
		didParse = true
		parsedLine = t1 + parseLink2(t2) + t3
	}

	// check link (3)
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteContext.vars.GetVal("ASWSG-LINK-3-1"), siteContext.vars.GetVal("ASWSG-LINK-3-2"), siteContext.vars.GetVal("ASWSG-ESCAPE")) // TOFIX implement link func
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + parseLink3(t2) + t3
	}

	if didParse == true {
		parsedLine = parseInLine(parsedLine)
	}

	return
}

// EOF