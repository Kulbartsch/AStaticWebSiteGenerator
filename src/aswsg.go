// Another|Alexanders Static WebSite Generator
// (c) 2016 Alexander Kulbartsch
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
// TODO: extract tools
// TODO: extract var handling

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	//"regexp"
	"io"
	"strconv"
)

type siteContextType struct {
	vars           SimpleVars
	inStream       io.Reader // ???
	outStream      io.Writer // ???
	paragraphState string
	lineNumber     uint32
	blockMode      string
}

type SimpleVars map[string]string

var siteVars = SimpleVars{
	"ASWSG-VERSION": "0.2",
	"ASWSG-AUTHOR":  "Alexander Kulbartsch",
	"ASWSG-LICENSE": "GPL V3",
	// control vars
	"ASWSG-MESSAGE-FILTER": "Dd",
	// inline formating, pairs end on -1 respective -2
	"ASWSG-VAR-1":    "{{", // special: variable to be replaced
	"ASWSG-VAR-2":    "}}",
	"ASWSG-LINK-1-1": "[[", // special: link
	"ASWSG-LINK-1-2": "]]",
	"ASWSG-LINK-1-3": "|",
	"ASWSG-LINK-2-1": "[", // special: link
	"ASWSG-LINK-2-2": ")",
	"ASWSG-LINK-2-3": "](",
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
	"ASWSG-DEFINE":  "@", // special: define var
	"ASWSG-INCLUDE": "+", // special: include parsed file
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

// general tool functions

// WhiteSpaceTrim
func WhiteSpaceTrim(in string) string {
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

// right returns the right most l char(s) of s in r
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

// Commands

// parse commands
func parseCommands(command string) (result []string) {
	var function, parameter string
	cmd := strings.Trim(command, " \t\n)")
	i := strings.IndexAny(cmd, " \t")
	if i == -1 {
		function = strings.ToUpper(cmd)
	} else {
		function = strings.ToUpper(WhiteSpaceTrim(cmd[:i]))
		parameter = WhiteSpaceTrim(cmd[i:])
	}
	switch function {
        case "COMMENT":
                // nothing
	case "DUMP-VARS":
		result = commandDumpVars(parameter)
	case "MESSAGE":
		result = commandMessage(parameter)
	default:
		Message("", 0, "W", "unknown command ignored (function/parameter): "+function+"/"+parameter+" = "+cmd)
		break
	}
	return
}

// command dump-vars  (to log)
func commandDumpVars(p string) (r []string) {
	Message("", 0, "I", "---- my vars ----")
	for key, value := range siteVars {
		Message("", 0, "I", key+":"+value)
	}
	return
}

// command comment (name?) 
/* func commandComment(p string) (r []string) {
        // maybe some implementation later
	return
} */


// TODO command dump-context  (to log)

// command message  (to log)
func commandMessage(p string) (r []string) {
	Message("", 0, "I", p)
	return
}

// TODO command interactive  (enter interactive mode = read from io.stdin)

// TODO command execute-shell-command  <command with parameters>

// TODO command include-raw-file <filename>  (include a raw file)

// TODO command include-crude-file <filename>  (include raw files, but with with variable replacing)

// TODO command execute-script <filename>  (run a script ... maybe in the future)

// simpleVar handling

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
	s, r := v[strings.ToUpper(key)]
	if r == false {
		Message("", 0, "W", "Key '"+key+"' does not exist.")
	}
	return s
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
	if !strings.ContainsAny(severity, siteVars.GetVal("ASWSG-MESSAGE-FILTER")) {
		log.Println(filename, ":", line, ":", severity, ":", messagetext)
	}
}

// main

func setDefaultSiteVars() {
	_ = siteVars.SetVar("TimeStampFormat", "2006-01-02 15:04:05 UTC+ 07:00")
	_ = siteVars.SetVar("DateFormat", "2006-01-02")
	_ = siteVars.SetVar("TimeFormat", "15:04:05")
	_ = siteVars.SetVar("now", time.Now().Format(siteVars.GetVal("TimeStampFormat")))
	_ = siteVars.SetVar("today", time.Now().Format(siteVars.GetVal("DateFormat")))
	_ = siteVars.SetVar("time", time.Now().Format(siteVars.GetVal("TimeFormat")))
}

func parseAndSetCommandLineVars() {
	destinationVar := "IN-FILE"
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		Message("$CMDLINEARG$", i, "D", arg)
		if strings.Index(arg, ":") >= 0 {
			if siteVars.ParseAndSetVar(arg) != true {
				Message("", i, "w", "Can't parse variable: "+arg)
			}
		} else { // doesn't seem to be a variable
			if destinationVar == "" {
				Message("$CMDLINEARG$", i, "W", "To much non variable parameters (ignored): "+arg)
			} else {
				if siteVars.SetVar(destinationVar, arg) != true {
					Message("", i, "w", "Can't set '"+destinationVar+"' to '"+arg+"'")
				}
				if destinationVar == "IN-FILE" {
					destinationVar = "OUT-FILE"
				} else {
					destinationVar = ""
				}
			}
		}
	}
}

// site context helper functions

func (c siteContextType) addStringToOutput(s string) (err error) {
	// ToDo implement
	return
}

// inline

type HTMLAttrib map[string]string

func generateHTMLTagWithAttributes(tag string, openTag bool, attrib HTMLAttrib) (resultHTMLTag string) {
	// TODO check using HTML library
	var attributes string
	if len(tag) == 0 {
		return
	}
	var close string
	if openTag == false {
		close = "/"
	} else {
		close = ""
		for k, v := range attrib {
			attributes = " "
			attributes += k + "=\"" + v + "\" "
		}
	}
	resultHTMLTag = "<" + close + tag + attributes + ">"
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

// Proceeses the inner part of an [[text]] (ASWSG-LINK-1-x) an generates a complete <a> tag.
// If text contains a "|"" (pipe) the left part is the displayed content and
// the right part is the href.
func parseLink1(text string) string {
	if len(text) == 0 {
		return ""
	}
	var link, display string
	var attrib HTMLAttrib
	i := strings.Index(text, siteVars.GetVal("ASWSG-LINK-1-3"))
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
	if len(text) == 0 {
		return ""
	}
	var link, display string
	var attrib HTMLAttrib
	i := strings.Index(text, siteVars.GetVal("ASWSG-LINK-2-3"))
	if i == -1 {
                return siteVars.GetVal("ASWSG-LINK-2-1") + text + siteVars.GetVal("ASWSG-LINK-2-2")
	} else {
		display = text[:i]
		link = text[i+2:]
	}
	if display == "" {
		display = link
	}
	attrib = HTMLAttrib{"href": link}
	return surroundWithHTMLTagWithAttributes("a", display, attrib) // tag string, s string, attrib HTMLAttrib)
}

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
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteVars.GetVal("ASWSG-BOLD-1"), siteVars.GetVal("ASWSG-BOLD-2"), siteVars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("b", t2) + t3
	}

	// check emphasised
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteVars.GetVal("ASWSG-EMP-1"), siteVars.GetVal("ASWSG-EMP-2"), siteVars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 && Right(t1, 1) != ":" { // TOFIX: Workaround to not mess with HTML links containing "://"
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("em", t2) + t3
	}

	// check strike
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteVars.GetVal("ASWSG-STRIKE-1"), siteVars.GetVal("ASWSG-STRIKE-2"), siteVars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("del", t2) + t3
	}

	// check code
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteVars.GetVal("ASWSG-CODE-1"), siteVars.GetVal("ASWSG-CODE-2"), siteVars.GetVal("ASWSG-ESCAPE"))
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + surroundWithHTMLTag("code", t2) + t3
	}

	// check link (1)
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteVars.GetVal("ASWSG-LINK-1-1"), siteVars.GetVal("ASWSG-LINK-1-2"), siteVars.GetVal("ASWSG-ESCAPE")) // TOFIX implement link func
	if len(t2) > 0 {
		didParse = true
		parsedLine = t1 + parseLink1(t2) + t3
	}

	// check link (2)
	t1, t2, t3 = StringBracketsSplit(parsedLine, siteVars.GetVal("ASWSG-LINK-2-1"), siteVars.GetVal("ASWSG-LINK-2-2"), siteVars.GetVal("ASWSG-ESCAPE")) // TOFIX implement link func
	l2 := strings.Index(t2, siteVars.GetVal("ASWSG-LINK-2-3"),)
	if len(t2) > 0 && l2 >= 0 {
		didParse = true
		parsedLine = t1 + parseLink2(t2) + t3
	}

	if didParse == true {
		parsedLine = parseInLine(parsedLine)
	}

	return
}

func parseAndSetVar(line string) (varParsed bool) {
	if strings.ContainsAny(line[0:1], siteVars.GetVal("ASWSG-DEFINE")) {
		siteVars.ParseAndSetVar(line[1:])
		return true
	}
	return false
}

func replaceInlineVars(line string) string {
	// TODO change to use interface Simplevars
	t1, t2, t3 := StringBracketsSplit(line, siteVars.GetVal("ASWSG-VAR-1"), siteVars.GetVal("ASWSG-VAR-2"), siteVars.GetVal("ASWSG-ESCAPE"))
	if !siteVars.ExistsVal(t2) {
		return line
	}
	return replaceInlineVars(t1 + siteVars.GetVal(t2) + t3)
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
	controlChars := siteVars.GetVal("ASWSG-LIST") + siteVars.GetVal("ASWSG-CITE") + siteVars.GetVal("ASWSG-NUMERATION")
	// Message("", 0, "D", "inline:" + line)
	// Message("", 0, "D", "  PS:" + currentParagraphState)

	// parse LIST, CITE and NUMERATION
	for _, r := range line {
		a := string(r)
		Message("", 0, "D", "  controlChar:"+a)
		if ContainsOnly(a, controlChars) {
			switch {
			case ContainsOnly(a, siteVars.GetVal("ASWSG-LIST")):
				resultingParagraphState = resultingParagraphState + "L"
				surroundWith = "li"
			case ContainsOnly(a, siteVars.GetVal("ASWSG-CITE")):
				resultingParagraphState = resultingParagraphState + "C"
			case ContainsOnly(a, siteVars.GetVal("ASWSG-NUMERATION")):
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
	if firstChar == siteVars.GetVal("ASWSG-ESCAPE") {
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

	// TODO block mode raw

	// TODO block mode code

	// parse commands
	if line[0:1] == siteVars.GetVal("ASWSG-COMMAND") {
		resultLines = append(resultLines, parseCommands(line[1:])...)
		return
	}

	// Entering Markup Mode

	// parse raw
	if line[0:1] == siteVars.GetVal("ASWSG-RAWLINE") {
		resultLines = append(resultLines, line[1:])
		return
	}

	// process includes
	if strings.ContainsAny(line[0:1], siteVars.GetVal("ASWSG-INCLUDE")) {
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
	if strings.ContainsAny(line[0:1], siteVars.GetVal("ASWSG-HEADER")) {
		newParagraphState = ""
		// header parser
		fc, count, content := firstCharCountAndTrim(line)
		if !ContainsOnly(fc, siteVars.GetVal("ASWSG-HEADER")) {
			Message("", 0, "A", "should not happen - expected ASWSG-HEADER character")
		}
		level := strconv.Itoa(count)
		resultLines = append(changeParagraphs(paragraphState, newParagraphState, false), "<h"+level+">"+content+"</h"+level+">")
		return resultLines, newParagraphState
	}

	// parse one liner: horizontal line
	if ContainsOnly(strings.TrimRight(line, " \t"), siteVars.GetVal("ASWSG-LINE")) && len(strings.TrimRight(line, " \t")) >= 3 {
		newParagraphState = " "
		resultLines = append(changeParagraphs(paragraphState, newParagraphState, false), "<hr \\>")
		return resultLines, newParagraphState
	}

	// TODO parse Table

	resultLines, newParagraphState = parseCommonParagraphControls(line, paragraphState)

	return resultLines, newParagraphState
}

//TODO remove ?
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

// core logic

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
	file_stat, stat_error := file.Stat()
	if stat_error != nil {
		return nil, paragraphState, err
	}
	siteVars.SetVar("filename", file_stat.Name())

	scanner := bufio.NewScanner(file)
	// TODO check for errors

	var result []string

	for scanner.Scan() {
		// TODO check for errors
		// TODO handle line numbers
		lines, paragraphState = parseLine(scanner.Text(), paragraphState)
		result = append(result, lines...)
	}

	return result, paragraphState, scanner.Err()

}

func main() {

	var parsedText []string
	var err error

	Message("", 0, "D", "---- ASWSG start ----")

	setDefaultSiteVars()

	paragraphState := ""

	parseAndSetCommandLineVars()

	// TODO set original file name

	parsedText, paragraphState, err = parseFile(siteVars.GetVal("IN-FILE"), paragraphState)
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
