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
	lineNumber     int
	blockMode      string
	tableLine      int
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
		var fn string
		var ln int
		if filename == "" {
			fn = siteContext.vars.GetVal("filename")
		} else {
			fn = filename
		}
		if line == 0 {
			ln = siteContext.lineNumber
		} else {
			ln = line
		}
		log.Println(fn, ":", ln, ":", severity, ":", messagetext)
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
		"ASWSG-VERSION": "0.7",
		"ASWSG-AUTHOR":  "Alexander Kulbartsch",
		"ASWSG-LICENSE": "AGPL V3",

		// control vars
		"ASWSG-MESSAGE-FILTER":       "Dd",     // D = Debug
		"ASWSG-AUTO-GENERATE-ANCHOR": "T",      // T = true, everything else is false
		"ASWSG-TABLE-HEADERLINES":    "1",      // number of headers, when parsing a table.
		"ASWSG-TABLE-ALIGNMENT":      "LL",     // L -> <th style="text-align:left">, C = center, R = right, other/default = left

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


// vars

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


// parse one line 
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
		tmpLine := siteContext.lineNumber
    tmpFilename := siteContext.vars.GetVal("filename")
		parsedLines, parsedParagraph, err := parseFile(line[1:], newParagraphState)
		siteContext.lineNumber = tmpLine
		siteContext.vars.SetVar("filename", tmpFilename)
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

	// parse Table
	if strings.ContainsAny(line[0:1], siteContext.vars.GetVal("ASWSG-TABLE")) {
		if paragraphState != "T" {
			siteContext.tableLine = 0
		}
		newParagraphState = "T"
		tline := parseTableLine(WhiteSpaceTrim(line[1:]))
		resultLines = append(changeParagraphs(paragraphState, newParagraphState, false), tline)
		return resultLines, newParagraphState
	}

	
	// parse one liner: horizontal line
	if ContainsOnly(strings.TrimRight(line, " \t"), siteContext.vars.GetVal("ASWSG-LINE")) && len(strings.TrimRight(line, " \t")) >= 3 {
		newParagraphState = " "
		resultLines = append(changeParagraphs(paragraphState, newParagraphState, false), "<hr \\>")
		return resultLines, newParagraphState
	}

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

	// TODO set vars
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
	siteContext.lineNumber = 0

	// stream scanner
	for scanner.Scan() {
		// TODO check for errors ?

		var oneInputLine string

		oneInputLine = continued + scanner.Text() // TODO check for Error
		siteContext.lineNumber += 1
		siteContext.vars.SetVar("linenumber", strconv.Itoa(siteContext.lineNumber))

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

	setDefaultSiteVars()

	Message("", 0, "D", "---- ASWSG start ----")

	paragraphState := ""

	parseAndSetCommandLineVars()

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
