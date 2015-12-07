// Another|Alexanders Static WebSite Generator
// (c) 2015 Alexander Kulbartsch
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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	//"regexp"
)

type SimpleVars map[string]string

var siteVars = SimpleVars{
	"ASWSG-VERSION": "0.2",
	"ASWSG-AUTHOR":  "Alexander Kulbartsch",
	"ASWSG-LICENSE": "GPL V3",
	// inline formating, pairs end on -1 respective -2
	"ASWSG-VAR-1":    "{{", // special: variable to be replaced
	"ASWSG-VAR-2":    "}}",
	"ASWSG-LINK-1":   "[[", // special: link
	"ASWSG-LINK-2":   "]]",
	"ASWSG-BOLD-1":   "*",  // inline: bold
	"ASWSG-BOLD-2":   "*",
	"ASWSG-EMP-1":    "//", // inline: emphasised
	"ASWSG-EMP-2":    "//",
	"ASWSG-CODE-1":   "``", // inline: code
	"ASWSG-CODE-2":   "``",
	"ASWSG-STRIKE-1": "~~", // inline: strike through
	"ASWSG-STRIKE-2": "~~",
	// line level formating (for paragraphs) at begin of line, using one of the characters
	"ASWSG-DEFINE":     "@",  // special: define var
	"ASWSG-INCLUDE":    "+",  // special: include
	"ASWSG-ESCAPE":     "\\", // special: escape char for paragraph
	                          // paragraph: initial state: __ (empty)
							  // paragraph: _P_aragraph
	"ASWSG-LIST":       "*-", // paragraph: _L_ist and _B_ullets
	"ASWSG-CITE":       ">",  // paragraph: _C_ite
	"ASWSG-NUMERATION": "#0123456789", // paragraph: _N_umbered list and _B_ullets
	// "ASWSG-TABLE":      "|",  // paragraph: _T_able and _R_ows and _D_ata // TODO implement
	"ASWSG-HEADER":     "=!", // one liner: header
	// single multi char in one line alone, at least 3
	"ASWSG-LINE":    "-", // special: horizontal line
	"ASWSG-ML-CODE": "%", // start/end block: code c_O_de
	"ASWSG-ML-CITE": ">", // start/end block: cite _M_ention
}

var paragraphTags = map[string]string {
	" ": "",
	"B": "li",
	"C": "cite",
	"D": "td",
	"L": "ul",
	"M": "cite",
	"N": "ol",
	"O": "pre",
	"P": "p",
	"R": "tl",
	"T": "table",
	"b": "b",
}

var paragraphState string


// simpleVar handling

// WhiteSpaceTrim
func WhiteSpaceTrim(in string) string {
	return strings.Trim(in, " \t\n")
}

func (v SimpleVars) SetVar(key, val string) (ok bool) {
	tkey := WhiteSpaceTrim(key)
	if len(tkey) == 0 {
		return false
	}
	v[strings.ToUpper(tkey)] = WhiteSpaceTrim(val)
	return true
}

func (v SimpleVars) GetVal(key string) (result string) {
	tkey := WhiteSpaceTrim(key)
	if len(tkey) == 0 {
		return ""
	}
	return v[strings.ToUpper(key)]
}

func (v SimpleVars) ExistsVal(key string) (result bool) {
	tkey := WhiteSpaceTrim(key)
	if len(tkey) == 0 {
		return false
	}
	_, result = v[strings.ToUpper(key)]
	return
}

func (v SimpleVars) ParseAndSetVar(toparse string) (ok bool) {
	dp := strings.Index(toparse, ":")
	if dp < 1 || dp == len(toparse) {
		return false
	}
	v.SetVar(toparse[0:(dp)], toparse[(dp+1):])
	return true
}

// message handling

func Message(filename string, line int, severity string, messagetext string) {
	fmt.Println(filename, ":", line, ":", severity, ":", messagetext)
}

// main

func setDefaultSiteVars() {
	_ = siteVars.SetVar("TimeStampFormat", "2006-01-02 15:04:05 UTC+ 07:00")
	_ = siteVars.SetVar("DateFormat", "2006-01-02")
	_ = siteVars.SetVar("TimeFormat", "15:04:05")
	_ = siteVars.SetVar("now", time.Now().Format(siteVars.GetVal("TimeStampFormat")))
	_ = siteVars.SetVar("today", time.Now().Format(siteVars.GetVal("DateFormat")))
}

func parseAndSetCommandLineVars() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		Message("$CMDLINEARG$", i, "D", arg)
		if strings.Index(arg, ":") >= 0 {
			if siteVars.ParseAndSetVar(arg) != true {
				Message("", 0, "w", "Can't parse variable: "+arg)
			}
		} else {
			// TODO process none variable setting parameter
		}
	}
}

func firstCharCountAndTrim(line string) (firstChar string, count int, content string) {
	if len(line) == 0 {
		return "", 0, ""
	}
	firstChar = line[0:1]
	for count = 1; line[count] == firstChar[0]; count++ {
	}
	content = strings.Trim(line[count:], " \t")
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

// right returns the right most char of s in r
func right(s string, l int) (r string) {
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



// inline

func generateHTMLTag(tag string, openTag bool) (resultHTMLTag string) {
	// TODO check using HTML library
	if len(tag) == 0 {
		return
	}
	var close string
	if openTag == false {
		close = "/"
	} else {
		close = ""
	}
	resultHTMLTag = "<" + close + tag + ">"
	return
}

func surroundWithHTMLTag(s string, tag string) string {
	return generateHTMLTag(tag, true) + s + generateHTMLTag(tag, false)
}

func generateTag(tagKind string, openTag bool) (resultTag string) {
	if len(tagKind) == 0 {
		return
	}
	resultTag = generateHTMLTag(paragraphTags[tagKind], openTag)
	return
}

func StringBracketsSplit(text string, b1 string, b2 string, escape string) (a string, b string, c string) {
	m := strings.Index(text, b1)
	if m == -1 { // TODO: maybe Check for code rune
		return text, "", "1"
	}
	n := strings.Index(text[m+1:], b2) + m + 1
	if n == -1 || n <= m { // ToDo: maybe Check for Escape rune
		return text, "", "2"
	}
	return text[0:m], text[m+len(b1) : n], text[n+len(b2):]
}

func parseInLine(rawLine string) (parsedLine string) {

	didParse := false
	parsedLine = rawLine

	// check bold
	t1, t2, t3 := StringBracketsSplit(parsedLine, siteVars.GetVal("ASWSG-BOLD-1"), siteVars.GetVal("ASWSG-BOLD-2"), siteVars.GetVal("ASWSG-ESC"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("b", t2) + t3
	}

	if didParse == true {
		parsedLine = parseInLine(parsedLine)
	}

	return
}


func parseAndSetVar(line string) (varParsed bool) {
	if strings.ContainsAny(line[0:1], siteVars.GetVal("ASWSG-VAR")) {
		siteVars.ParseAndSetVar(line[1:])
		return true
	}
   return false
}


func replaceInlineVars(line string) (string) {
	t1, t2, t3 := StringBracketsSplit(line, siteVars.GetVal("ASWSG-VAR-1"), siteVars.GetVal("ASWSG-VAR-2"), siteVars.GetVal("ASWSG-ESC"))
	if !siteVars.ExistsVal(t2) {
		return line
	}
	return t1 + siteVars.GetVal(t2) + t3
}


// line

// changeParagraphs returns the necessary HTML Tags to close the previous state and initiate the new one.
// if both states are the same refreshInner forces the inner tag to be closed and opened
func changeParagraphs(oldParagraphState string, newParagraphState string, refreshInner bool) (resultLines []string) {
	if oldParagraphState == newParagraphState {
		if refreshInner {
			resultLines = append(resultLines, generateTag(right(oldParagraphState, 1), true))
			resultLines = append(resultLines, generateTag(right(newParagraphState, 1), false))
		}
		return
	}
	intermediateState := oldParagraphState
	for len(intermediateState) > 0 && intermediateState != newParagraphState[:len(intermediateState)] {
		x := len(intermediateState) - 1
		resultLines = append(resultLines, generateTag(right(intermediateState, 1), true))
		intermediateState = intermediateState[:x]
	}
	for len(intermediateState) < len(newParagraphState) && intermediateState != newParagraphState {
		x := len(intermediateState)
		addState := newParagraphState[x:x+1]
		resultLines = append(resultLines, generateTag(addState, false))
		intermediateState = intermediateState + addState
	}
	return
}


// parse paragraph line + parse inline
func parseCommonParagraphControls(line string, oldParagraphState string, newParagraphState string) (resultLines []string, resultingParagraphState string) {
	firstChar := line[0:1]
	resultingParagraphState = newParagraphState
	if strings.ContainsAny(firstChar, siteVars.GetVal("ASWSG-LIST")) {
		resultingParagraphState = resultingParagraphState + "L"
		return parseCommonParagraphControls(line[1:], oldParagraphState, resultingParagraphState)
	}

	if len(resultingParagraphState) == 0 {
		resultingParagraphState = "P"
	}

	resultLines = append( changeParagraphs(paragraphState, resultingParagraphState, false), parseInLine(line) )

	return
}



func parseLine(line string, paragraphState string) (resultLines []string, newParagraphState string) {

	newParagraphState = paragraphState

	// empty line
	if len(line) == 0 {
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

	// process includes
	if strings.ContainsAny(line[0:1], siteVars.GetVal("ASWSG-INCLUDE")) {
		parsedLines, parsedParagraph, err := parseFile(line[1:], newParagraphState)
		if err != nil {
			Message(line[1:], -1, "E", err.Error())
		}
		resultLines = parsedLines
		newParagraphState = parsedParagraph
		return resultLines, newParagraphState
	}

	// parse one liner: header
	if strings.ContainsAny(line[0:1], siteVars.GetVal("ASWSG-HEADER")) {
		newParagraphState = ""
		// TODO implement real header parser ...
		resultLines = append( changeParagraphs(paragraphState, newParagraphState, false), "<h1>" + line[1:] + "</h1>" )
		return resultLines, newParagraphState
	}

	// parse one liner: horizontal line
	if ContainsOnly(strings.TrimRight(line, " \t"), siteVars.GetVal("ASWSG-LINE")) && len(strings.TrimRight(line, " \t")) >= 3 {
		newParagraphState = " "
		resultLines = append( changeParagraphs(paragraphState, newParagraphState, false), "<hr \\>")
		return resultLines, newParagraphState
	}

	// parse paragraph line (list, ...) + parse inline
	resultLines, newParagraphState = parseCommonParagraphControls(line, paragraphState, newParagraphState)

	return resultLines, newParagraphState
}


//TODO remove
func ReadTextFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// test
// TODO move tests to test file

func TestSBS(text string) {
	a, b, c := StringBracketsSplit(text, "{{", "}}", "\\")
	fmt.Println(a, "*", b, "*", c)
}

func TestSBS2(text string) {
	a, b, c := StringBracketsSplit(text, "_", "_", "\\")
	fmt.Println(a, "*", b, "*", c)
}


// core logic

func parseFile(filename string, startParagraphState string) ([]string, string, error) {
	var lines []string

	paragraphState = startParagraphState

	file, err := os.Open(filename)
	if err != nil {
		return nil, paragraphState, err
	}
	defer file.Close()

	// ToDo set vars
	//    - filename        ok
	//    - fqfn
	//    - file last change date + time
	//    - file creation date + time
	file_stat, stat_error := file.Stat()
	if stat_error != nil {
		return nil, paragraphState, err
	}
	siteVars.SetVar("filename", file_stat.Name())

	scanner := bufio.NewScanner(file)

	var result []string

	for scanner.Scan() {
		lines, paragraphState = parseLine(scanner.Text(), paragraphState)
		result = append(result, lines...)
	}

	return result, paragraphState, scanner.Err()

}

func main() {

	var parsedText []string
	var err error

	setDefaultSiteVars()

	paragraphState = ""

	parseAndSetCommandLineVars()

	parsedText, paragraphState, err = parseFile(siteVars.GetVal("IN-FILE"), paragraphState)
	if err != nil {
		fmt.Println("Error:", err.Error());
	}

	// cleanup unclosed paragraphs
	parsedText = append(parsedText, changeParagraphs(paragraphState, "", false)...)

	// Output
	// TODO use out file
	for _, l := range parsedText {
		fmt.Println( l )
	}

	// DEBUG remove
	fmt.Println("---- Resulting paragraph style :", paragraphState)



	// Tests

	fmt.Println("---- Test Line Parsing ----")

	fmt.Println(parseLine("@test:OK", ""))
	fmt.Println(parseLine("@ASWSG-VAR:$@", ""))
	fmt.Println(parseLine("@FOO:foo", ""))
	fmt.Println(parseLine("$BAA:baa", ""))
	fmt.Println(parseLine("= Welcome", ""))
	fmt.Println(parseLine("== To the Future =", ""))
	fmt.Println(parseLine("", ""))
	fmt.Println(parseLine("Bla bla", ""))

	// just some test
	fmt.Println("---- my vars ----")
	for key, value := range siteVars {
		fmt.Println(key, ":", value)
	}

	fmt.Println("---- inline test ----")

	TestSBS("1 Hallo_{{name}}!")
	TestSBS("2 Hallo_{{na{{me}}!")
	TestSBS("3 Hallo_{{name}} und {{nummer}}!")
	TestSBS("4 Hallo_{{name))!")
	TestSBS("5 Hallo_{{}}!")
	TestSBS("{{ 6 name}}")
	TestSBS("7 brave world")
	TestSBS2("11 Hallo-{{_name__))!")
	TestSBS2("12 Hallo _name_!")
	TestSBS2("13 Hallo __name__!")
	TestSBS("14 Hallo_\\{{name}}!")

	fmt.Println("---- other ----")

	fmt.Println("right of 'aBc': " + right("aBc", 1))
	fmt.Println("right of nothing: " + right("", 1))
	fmt.Println("right 2 of '4321': " + right("4321", 2))
	fmt.Println("---- By! ----")


}

//EOF
