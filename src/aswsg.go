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
	"ASWSG-VERSION": "0.1",
	"ASWSG-AUTHOR":  "Alexander Kulbartsch",
	"ASWSG-LICENSE": "GPL V3",
	// inline formating, pairs end on -1 respective -2
	"ASWSG-ESCAPE":   "\\",
	"ASWSG-VAR-1":    "{{",
	"ASWSG-VAR-2":    "}}",
	"ASWSG-LINK-1":   "[[",
	"ASWSG-LINK-2":   "]]",
	"ASWSG-BOLD-1":   "*",
	"ASWSG-BOLD-2":   "*",
	"ASWSG-EMP-1":    "//",
	"ASWSG-EMP-2":    "//",
	"ASWSG-CODE-1":   "``",
	"ASWSG-CODE-2":   "``",
	"ASWSG-STRIKE-1": "~~",
	"ASWSG-STRIKE-2": "~~",
	// line level formating at begin of line, using one of the characters
	"ASWSG-DEFINE":     "@",
	"ASWSG-INCLUDE":    "+",
	"ASWSG-LIST":       "*-",
	"ASWSG-CITE":       ">",
	"ASWSG-HEADER":     "=!",
	"ASWSG-NUMERATION": "#0123456789",
	"ASWSG-TABLE":      "|",
	// single multi char in one line alone, at least 3
	"ASWSG-LINE":    "-", // horizontal line
	"ASWSG-ML-CODE": "%", // start/end code block
	"ASWSG-ML-CITE": ">",
}

// Var handling

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
	// TODO: check if value exists
	return v[strings.ToUpper(key)]
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

func parseAndSetVar(line string) {
	//var validID = regexp.MustCompile(`^@(.+):(.+)`)
	//if validID.MatchString(line) {
	//if line[0] == siteVars["ASWSG-VAR"][0]) {
	if strings.ContainsAny(line[0:1], siteVars.GetVal("ASWSG-VAR")) {
		siteVars.ParseAndSetVar(line[1:])
	}
}

// line

func parseLine(line string) string {
	switch {
	case len(line) == 0:
		return "/0"
	case strings.ContainsAny(line[0:1], siteVars["ASWSG-VAR"]):
		parseAndSetVar(line)
		return "/var"
	case strings.ContainsAny(line[0:1], siteVars["ASWSG-HEADER"]):
		return "<h1>" + line[1:] + "</h1>"
	}
	return line
}

// inline

func StringBracketsSplit(text string, b1 string, b2 string, escape string) (a string, b string, c string) {
	m := strings.Index(text, b1)
	if m == -1 { // ToDo: Check for Escape rune
		return text, "", "1"
	}
	n := strings.Index(text[m+1:], b2) + m + 1
	if n == -1 || n <= m { // ToDo: Check for Escape rune
		return text, "", "2"
	}
	return text[0:m], text[m+len(b1) : n], text[n+len(b2):]
}

//HACK: stream parsing

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

func TestSBS(text string) {
	a, b, c := StringBracketsSplit(text, "{{", "}}", "\\")
	fmt.Println(a, "*", b, "*", c)
}

func TestSBS2(text string) {
	a, b, c := StringBracketsSplit(text, "_", "_", "\\")
	fmt.Println(a, "*", b, "*", c)
}

func parseFile(filename string, lineState string) ([]string, string) {
	// ToDo
	// 1. read file
	// 2. replace vars
	// 3. parse line
	// 4. parse inline

	// TODO dummy, remove
	result := []string{"test", "empty"}
	return result, ""
}

func main() {

	setDefaultSiteVars()

	parseAndSetCommandLineVars()

	parseFile("", " ")

	// Tests

	fmt.Println("Hi!")

	fmt.Println(parseLine("@test:OK"))
	fmt.Println(parseLine("@ASWSG-VAR:$@"))
	fmt.Println(parseLine("@FOO:foo"))
	fmt.Println(parseLine("$BAA:baa"))
	fmt.Println(parseLine("= Welcome"))
	fmt.Println(parseLine("== To the Future ="))
	fmt.Println(parseLine(""))
	fmt.Println(parseLine("Bla bla"))

	// just some test
	fmt.Println("---- my vars ----")
	for key, value := range siteVars {
		fmt.Println(key, ":", value)
	}

	fmt.Println("inline.")

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

	fmt.Println("By!")

}

//EOF
