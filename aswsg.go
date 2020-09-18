// Another|Alexanders Static WebSite Generator
// (c) 2016-2020 Alexander Kulbartsch
//
// Lines beginning with "@" (and no further white space)
// are interpreted as variables in the form "@var: value".
// White spaces after the dubble colon is optional and gets removed.
//
// Lines beginning with a "+" (and no further white space)
// are interpreted in the form "+filename".
// The named file will be included here.
//
// \ escapes the special line characters (and will be removed)
//
// A {{variable}} in the text will be replaced by the named variable
//

// TODO: remember line number for meta Message (will be done with context type)
// TODO: use OUT-FILE (will be done with context type)
// TODO: inherit of html lines(/blocks?). identified by starting with an "<". (Ending with a empty line?)

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type siteContextType struct {
	vars           SimpleVars
	inStream       io.Reader // ???
	outStream      io.Writer // ???
	paragraphState string
	lineNumber     uint32
	blockMode      string
}

var siteContext siteContextType

var paragraphTags = map[string]string{
	" ": "",
	"B": "li",
	"C": "cite",
	"D": "td",
	"L": "ul", // "ul style=\"list-style-type:circle\"",
	"M": "cite",
	"N": "ol",
	"O": "pre",
	"P": "p",
	"R": "tl",
	"T": "table",
	"b": "b",
}

type indexedLinksType struct {
	index string
	link  string
}

var indexedLinks []indexedLinksType
var linkIndex int = 1


// message handling

// Message logs to stderr
func Message(filename string, line int, severity string, messagetext string) {
	if !strings.ContainsAny(severity, siteContext.vars.GetVal("ASWSG-MESSAGE-FILTER")) {
		log.Println(filename, ":", line, ":", severity, ":", messagetext)
	}
}


// main

// site context helper functions

func (c siteContextType) addStringToOutput(s string) (err error) {
	// ToDo implement
	return
}

func setDefaultSiteVars() {
	siteContext.vars = SimpleVars{ // was: var siteVars
		"ASWSG-VERSION": "0.5",
		"ASWSG-AUTHOR":  "Alexander Kulbartsch",
		"ASWSG-LICENSE": "GPL V3",

		// control vars
		"ASWSG-MESSAGE-FILTER":       "Dd",  // D = Debug
		"ASWSG-AUTO-GENERATE-ANCHOR": "T",   // T = true, everything else is false

		// inline formating, pairs end on -1 respective -2
		"ASWSG-VAR-1":    "{{", // special: variable to be replaced
		"ASWSG-VAR-2":    "}}",
		"ASWSG-LINK-1-1": "[[", // special: link internal
		"ASWSG-LINK-1-2": "]]",
		"ASWSG-LINK-1-3": "|",
		"ASWSG-LINK-2-1": "[", // special: link
		"ASWSG-LINK-2-2": ")",
		"ASWSG-LINK-2-3": "](",
		"ASWSG-LINK-3-1": "[[", // special: indexed-link
		"ASWSG-LINK-3-2": "]]",
		"ASWSG-BOLD-1":   "*", // inline: bold
		"ASWSG-BOLD-2":   "*",
		"ASWSG-EMP-1":    "//", // inline: emphasised
		"ASWSG-EMP-2":    "//",
		"ASWSG-CODE-1":   "``", // inline: code
		"ASWSG-CODE-2":   "``",
		"ASWSG-STRIKE-1": "~~", // inline: strike through
		"ASWSG-STRIKE-2": "~~",
		"ASWSG-UNDERL-1": "__", // TODO inline: underline
		"ASWSG-UNDERL-2": "__",

		// line level formating (for paragraphs) at begin of line, using one of the characters
		"ASWSG-DEFINE":   "@",  // special: define var
		"ASWSG-INCLUDE":  "+",  // special: include parsed file
		"ASWSG-CONTINUE": "\\", // special: if at end of line, continue (join) with next line
		// "ASWSG-RAWFILE": "<",  // special: include raw file - won't implemented this way, but as command. This special character will be used to identify raw HTML code. See ASWSG-RAWHMTL.
		"ASWSG-RAWHMTL": "<",  // TODO special: ram html line (this may have leading white spaces)
		"ASWSG-RAWLINE": "$",  // special: raw (html) line
		"ASWSG-ESCAPE":  "\\", // special: escape char for paragraph
		// ... paragraph: initial state: __ (empty)
		// ... paragraph: _P_aragraph
		"ASWSG-LIST":       "*-",          // paragraph: _L_ist and _B_ullets
		"ASWSG-CITE":       ">",           // paragraph: _C_ite
		"ASWSG-NUMERATION": "#0123456789", // paragraph: _N_umbered list and _B_ullets
		"ASWSG-COMMAND":    "(",           // single line command, optionally closed by an ")", should not be changed // TODO implement commands
		"ASWSG-TABLE":      "|",           // paragraph: _T_able and _R_ows and D_ata // TODO implement table
		"ASWSG-HEADER":     "=!",          // one liner: header

		// single multi char in one line alone, at least 3
		"ASWSG-LINE":    "-", // special: horizontal line
		"ASWSG-ML-CODE": "%", // TODO start/end block: code c_O_de
		"ASWSG-ML-CITE": ">", // TODO start/end block: cite _M_ention
		"ASWSG-ML-RAW":  "$", // TODO start/end block: raw line (i.e. for HTML code)
	}
	_ = siteContext.vars.SetVar("TimeStampFormat", "2006-01-02 15:04:05 UTC+ 07:00")
	_ = siteContext.vars.SetVar("DateFormat", "2006-01-02")
	_ = siteContext.vars.SetVar("TimeFormat", "15:04:05")
	_ = siteContext.vars.SetVar("now", time.Now().Format(siteContext.vars.GetVal("TimeStampFormat")))
	_ = siteContext.vars.SetVar("today", time.Now().Format(siteContext.vars.GetVal("DateFormat")))
	_ = siteContext.vars.SetVar("time", time.Now().Format(siteContext.vars.GetVal("TimeFormat")))
}

func parseAndSetCommandLineVars() {
	destinationVar := "IN-FILE"
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		Message("$CMDLINEARG$", i, "D", arg)
		if strings.Index(arg, ":") >= 0 {
			if siteContext.vars.ParseAndSetVar(arg) != true {
				Message("", i, "w", "Can't parse variable: "+arg)
			}
		} else {
			if siteContext.vars.SetVar(destinationVar, arg) != true {
				Message("", i, "w", "Can't parse variable: "+arg)
			} else {
				if destinationVar == "IN-FILE" {
					destinationVar = "OUT-FILE"
				} else {
					Message("$CMDLINEARG$", i, "W", "To much non variable parameters (ignored): "+arg)
				}
			}
		}
	}
}

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

// Proceeses the inner part of an [[link]] format (ASWSG-LINK-3-x) and
// generates a numbered <a> tag refering to an index.
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

func parseAndSetVar(line string) (varParsed bool) {
	if strings.ContainsAny(line[0:1], siteContext.vars.GetVal("ASWSG-DEFINE")) {
		siteContext.vars.ParseAndSetVar(line[1:])
		return true
	}
	return false
}

func replaceInlineVars(line string) string {
	// TODO change to use interface Simplevars
	t1, t2, t3 := StringBracketsSplit(line, siteContext.vars.GetVal("ASWSG-VAR-1"), siteContext.vars.GetVal("ASWSG-VAR-2"), siteContext.vars.GetVal("ASWSG-ESCAPE"))
	if !siteContext.vars.ExistsVal(t2) {
		return line
	}
	return replaceInlineVars(t1 + siteContext.vars.GetVal(t2) + t3)
}


// line

// changeParagraphs returns the necessary HTML Tags to close the previous state and initiate the new one.
// if both states are the same refreshInner forces the inner tag to be closed and opened
func changeParagraphs(oldParagraphState string, newParagraphState string, refreshInner bool) (resultLines []string) {
	if oldParagraphState == newParagraphState && len(oldParagraphState) > 0 {
		if refreshInner {
			resultLines = append(resultLines, generateTag(Right(oldParagraphState, 1), false))
			resultLines = append(resultLines, generateTag(Right(newParagraphState, 1), true))
		}
		return
	}
	intermediateState := oldParagraphState
	Message("", 0, "D", "intermediateState: '"+intermediateState+"', newParagraphState: '"+newParagraphState+"'")
	// close previous paragraph state(s)
	for len(intermediateState) > 0 {
		iSLen := len(intermediateState)
		if iSLen <= len(newParagraphState) && intermediateState == newParagraphState[:iSLen] {
			break
		}
		x := len(intermediateState) - 1
		resultLines = append(resultLines, generateTag(Right(intermediateState, 1), false))
		intermediateState = intermediateState[:x]
	}
	// open new paragraph state(s)
	for len(intermediateState) < len(newParagraphState) /* && intermediateState != newParagraphState */ {
		x := len(intermediateState)
		addState := newParagraphState[x : x+1]
		resultLines = append(resultLines, generateTag(addState, true))
		intermediateState = intermediateState + addState
	}
	return
}


// parse paragraph line + parse inline
func parseCommonParagraphControls(line string, currentParagraphState string) (resultLines []string, resultingParagraphState string) {
	firstChar := line[0:1]
	resultingParagraphState = ""
	refreshInner := false
	surroundWith := ""
	controlChars := siteContext.vars.GetVal("ASWSG-LIST") + siteContext.vars.GetVal("ASWSG-CITE") + siteContext.vars.GetVal("ASWSG-NUMERATION")
	// Message("", 0, "D", "inline:" + line)
	// Message("", 0, "D", "  PS:" + currentParagraphState)

	// parse LIST, CITE and NUMERATION
	for _, r := range line {
		a := string(r)
		Message("", 0, "D", "  controlChar:"+a)
		if ContainsOnly(a, controlChars) {
			switch {
			case ContainsOnly(a, siteContext.vars.GetVal("ASWSG-LIST")):
				resultingParagraphState = resultingParagraphState + "L"
				surroundWith = "li"
			case ContainsOnly(a, siteContext.vars.GetVal("ASWSG-CITE")):
				resultingParagraphState = resultingParagraphState + "C"
			case ContainsOnly(a, siteContext.vars.GetVal("ASWSG-NUMERATION")):
				resultingParagraphState = resultingParagraphState + "N"
				surroundWith = "li"
			default:
				Message("", 0, "A", "should not happen (controlChar not found)")
				break
			}
			line = WhiteSpaceTrim(line[1:])
		} else {
			break
		}
		// Message("", 0, "D", "  resultState:"+ resultingParagraphState )
	}

	// parse Escape
	if firstChar == siteContext.vars.GetVal("ASWSG-ESCAPE") {
		line = line[1:]
	}

	// parse paragraph
	if len(resultingParagraphState) == 0 && len(line) != 0 {
		resultingParagraphState = "P"
	}

	var newLine string
	if surroundWith != "" {
		newLine = surroundWithHTMLTag(surroundWith, parseInLine(line))
	} else {
		newLine = parseInLine(line)
	}

	resultLines = append(changeParagraphs(currentParagraphState, resultingParagraphState, refreshInner), newLine)

	return
}


//
func parseLine(line string, paragraphState string) (resultLines []string, newParagraphState string) {

	newParagraphState = paragraphState
	lineLength := len(line)

	// empty line
	if lineLength == 0 {
		newParagraphState = ""
		resultLines = changeParagraphs(paragraphState, newParagraphState, false)
		return resultLines, newParagraphState
	}

	// replace inline vars
	line = replaceInlineVars(line)

	// parse vars
	if parseAndSetVar(line) == true {
		return resultLines, newParagraphState
	}

	// TODO block mode raw/crude

	// TODO block mode code

	// parse commands
	if line[0:1] == siteContext.vars.GetVal("ASWSG-COMMAND") {
		resultLines = append(resultLines, parseCommands(line[1:])...)
		return
	}

	// Entering Markup Mode

	// parse raw/crude (variables where allready replaced)
	if line[0:1] == siteContext.vars.GetVal("ASWSG-RAWLINE") {
		resultLines = append(resultLines, line[1:])
		return
	}

	// process includes
	if strings.ContainsAny(line[0:1], siteContext.vars.GetVal("ASWSG-INCLUDE")) {
		parsedLines, parsedParagraph, err := parseFile(line[1:], newParagraphState)
		// TODO restore linenumber and filename
		if err != nil {
			Message(line[1:], -1, "E", err.Error())
		}
		resultLines = parsedLines
		newParagraphState = parsedParagraph
		return resultLines, newParagraphState
	}

	// TODO block mode cite

	// parse Markup one liner

	// parse one liner: header
	if strings.ContainsAny(line[0:1], siteContext.vars.GetVal("ASWSG-HEADER")) {
		newParagraphState = ""
		// header parser
		fc, count, content := firstCharCountAndTrim(line)
		if !ContainsOnly(fc, siteContext.vars.GetVal("ASWSG-HEADER")) {
			Message("", 0, "A", "should not happen - expected ASWSG-HEADER character")
		}
		level := strconv.Itoa(count)
		// ToC ~~~
		anchor := ""
		if siteContext.vars.GetVal("ASWSG-AUTO-GENERATE-ANCHOR") == "T" {
			anchor = " id=\"" + strings.ReplaceAll(content, "\"", "'") + "\"" // remove " in content
			// toc_line := strings.Repeat(siteContext.vars.GetVal("ASWSG-LIST"), level) + " (" + content + ")[" + anchorText + "]"
			// TODO add toc_line anchor to list
		}
		//
		resultLines = append(changeParagraphs(paragraphState, newParagraphState, false), "<h"+level+anchor+">"+content+"</h"+level+">")
		return resultLines, newParagraphState
	}

	// parse one liner: horizontal line
	if ContainsOnly(strings.TrimRight(line, " \t"), siteContext.vars.GetVal("ASWSG-LINE")) && len(strings.TrimRight(line, " \t")) >= 3 {
		newParagraphState = " "
		resultLines = append(changeParagraphs(paragraphState, newParagraphState, false), "<hr \\>")
		return resultLines, newParagraphState
	}

	// TODO parse Table

	resultLines, newParagraphState = parseCommonParagraphControls(line, paragraphState)

	return resultLines, newParagraphState
}


// core logic /////////////////////////

func parseFile(filename string, startParagraphState string) ([]string, string, error) {
	var lines []string

	paragraphState := startParagraphState

	file, err := os.Open(filename)
	if err != nil {
		return nil, paragraphState, err
	}
	defer file.Close()

	// ToDo set vars
	//    - filename        ok
	//    - fqfn
	//    - basefn
	//    - file change date + time
	//    - file creation date + time
	fileStat, statError := file.Stat()
	if statError != nil {
		return nil, paragraphState, err
	}
	siteContext.vars.SetVar("filename", fileStat.Name())

	scanner := bufio.NewScanner(file)
	// TODO check for errors

	var result []string

	var continued string

	// stream scanner
	for scanner.Scan() {
		// TODO check for errors
		// TODO handle line numbers

		_ = "breakpoint"

		var oneInputLine string

		oneInputLine = continued + scanner.Text() // TODO check for Error

		if Right(oneInputLine, 1) == siteContext.vars.GetVal("ASWSG-CONTINUE") {
			continued = oneInputLine[:len(oneInputLine)-1]
			continue
		}
		continued = ""

		lines, paragraphState = parseLine(oneInputLine, paragraphState)
		result = append(result, lines...)

	}

	return result, paragraphState, scanner.Err()

}


func main() {

	var parsedText []string
	var err error

	_ = "breakpoint"

	setDefaultSiteVars()

	Message("", 0, "D", "---- ASWSG start ----")

	paragraphState := ""

	parseAndSetCommandLineVars()

	// TODO set original file name

	_ = "breakpoint"

	parsedText, paragraphState, err = parseFile(siteContext.vars.GetVal("IN-FILE"), paragraphState)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	Message("", 9999, "D", "---- Resulting paragraph style :"+paragraphState)

	// cleanup unclosed paragraphs
	parsedText = append(parsedText, changeParagraphs(paragraphState, "", false)...)

	// Output
	// TODO use out file / stream
	for _, l := range parsedText {
		fmt.Println(l)
	}

	Message("", 0, "D", "---- bye ----")

}

//EOF
