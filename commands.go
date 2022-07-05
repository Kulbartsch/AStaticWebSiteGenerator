// Commands

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
		// do nothing
	case "DUMP-VARS":
		result = commandDumpVars(parameter)
	case "MESSAGE":
		result = commandMessage(parameter)
	case "ANCHOR":
		result = append(result, "<a id=\""+parameter+"\"></a>")
	case "LINK-INDEX":
		result = commandLinkIndex(parameter)
	case "INCLUDE-FILE-CRUDE":
		result = commandIncludeFileCrude(parameter)
	case "INCLUDE-FILE-RAW":
		result = commandIncludeFileRaw(parameter)
	case "INCLUDE-CSV":
		result = commandIncludeCSV(parameter)
	case "INCLUDE-SCRIPT":
		result = commandIncludeScript(parameter)
	case "GT-AS-BLOCKQUOTE":
		result = commandGtAsBlockQuote(parameter)
	default:
		Message("", 0, "W", "unknown command ignored (function/parameter): "+function+"/"+parameter+" = "+cmd)
		break
	}
	return
}

// command dump-vars  (to log)
func commandDumpVars(p string) (r []string) {
	Message("", 0, "I", "---- my vars ----")
	for key, value := range siteContext.vars {
		Message("", 0, "I", key+":"+value)
	}
	return
}

// command message  (to log)
func commandMessage(p string) (r []string) {
	Message("", 0, "I", p)
	return
}

// command link-index
func commandLinkIndex(p string) (r []string) {
	for _, il := range indexedLinks {
		attrib := HTMLAttrib{"href": il.link, "rel": "external"}
		r = append(r, il.index+" "+surroundWithHTMLTagWithAttributes("a", il.link, attrib)+"<br />")
	}
	indexedLinks = nil
	return
}

// command include-file-crude <filename>        (include raw files, but with variable replacing)
func commandIncludeFileCrude(p string) (r []string) {
	var err error
	r, err = readTextFile(p, true)
	if err != nil {
		Message("", 0, "E", "Problem reading file: "+p)
		return nil
	}
	return
}

// command include-file-raw <filename>        (include raw files, with NO variable replacing)
func commandIncludeFileRaw(p string) (r []string) {
	var err error
	r, err = readTextFile(p, false)
	if err != nil {
		Message("", 0, "E", "Problem reading file: "+p)
		return nil
	}
	return
}

// command include-csv <filename>        (include a csv file, with NO variable replacing)
func commandIncludeCSV(p string) (r []string) {
	f, err := os.Open(p)
	if err != nil {
		Message("", 0, "E", "Problem reading CSV-file: "+p)
		return nil
	}
	crdr := csv.NewReader(f)
	crdr.Comma = []rune(siteContext.vars.GetVal("ASWSG-CSV-COMMA"))[0]
	crdr.Comment = []rune(siteContext.vars.GetVal("ASWSG-CSV-COMMENT"))[0]
	records, err2 := crdr.ReadAll()
	if err2 != nil {
		Message("", 0, "E", "Problem parsing CSV-file: "+p)
		return nil
	}
	r = append(r, generateHTMLTag("table", true))
	for i, cells := range records {
		hl, _ := strconv.Atoi(siteContext.vars.GetVal("ASWSG-TABLE-HEADERLINES"))
		siteContext.tableLine *= 1
		if hl >= i+1 {
			r = append(r, bulidTableRow(cells, "th"))
		} else {
			r = append(r, bulidTableRow(cells, "td"))
		}
	}
	r = append(r, generateHTMLTag("table", false))
	return
}

// command include-script <script> <parameters...>
func commandIncludeScript(p string) (r []string) {
	var command string
	var parameters []string
	cwp := strings.Trim(p, " \t\n)") // Command With Parameters
	i := strings.IndexAny(cwp, " \t")
	if i == -1 {
		command = WhiteSpaceTrim(cwp)
	} else {
		command = WhiteSpaceTrim(cwp[:i])
		parameters = strings.Split(WhiteSpaceTrim(cwp[i:]), " ")
	}
	out, err := exec.Command(command, parameters...).Output()
	if err != nil {
		Message("", 0, "E", "Problem executing: "+p)
		log.Println(err)
		Message("", 0, "E", "... Command: "+command)
		for _, p := range parameters {
			Message("", 0, "E", "... Parameters: "+p)
		}
	} else {
		r = append(r, string(out[:]))
	}
	return
}

// command commandGtAsBlockQuote <true/false>
func commandGtAsBlockQuote(p string) (r []string) {
	param := strings.ToUpper(p)
	if len(param) == 0 || param == "T" || param == "TRUE" {
		paragraphTags["C"] = "blockquote"
	} else {
		paragraphTags["C"] = "cite"
	}
	return
}

// command template (name?)
/* func commandTemplate(p string) (r []string) {
        // maybe some implementation later
	return
} */

// TODO command dump-context  (to log)

// include File as raw, or crude with replacing variables
func readTextFile(path string, crude bool) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Opening file error", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Closing file error", err)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if crude {
			lines = append(lines, replaceInlineVars(scanner.Text()))
		} else { // raw
			lines = append(lines, scanner.Text())
		}
	}
	return lines, scanner.Err()
}

// EOF
