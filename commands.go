// Commands

package main

import (
	"strings"
	"bufio"
	"log"
	"os"
	"os/exec"
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
		// nothing
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
	case "INCLUDE-SCRIPT":
		result = includeScript(parameter)
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
		Message("", 0, "E", "Problem reading file: " + p)
		return nil
	 }
	return 
} 

// command include-file-raw <filename>        (include raw files, with NO variable replacing)
func commandIncludeFileRaw(p string) (r []string) {
	var err error
	r, err = readTextFile(p, false)
	if err != nil { 
		Message("", 0, "E", "Problem reading file: " + p)
		return nil
	 }
	return 
} 

// command include-script <script> <parameters...>
func includeScript(p string) (r []string) {
	var command, parameter string
	cwp := strings.Trim(p, " \t\n)")
	i := strings.IndexAny(cwp, " \t")
	if i == -1 {
		command = WhiteSpaceTrim(cwp)
	} else {
		command = WhiteSpaceTrim(cwp[:i])
		parameter = WhiteSpaceTrim(cwp[i:])
	}
	out, err := exec.Command(command, parameter).Output()
	if err != nil { 
		Message("", 0, "E", "Problem executing: " + p)
		log.Println(err)
	} else {
		r = append(r, string(out[:]))
	}
	return
} 


// command template (name?)
/* func commandTemplate(p string) (r []string) {
        // maybe some implementation later
	return
} */

// TODO command dump-context  (to log)

// TODO command include-files <start-of-filname>     (includes all files beginning with given name, normal parsing)

// TODO command include-script                       (run an OS script, including its stdout)


// include File as raw, or crude with with replacing variables 
func readTextFile(path string, crude bool) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if crude {
			lines = append(lines, replaceInlineVars(scanner.Text()))
		} else {  // raw
			lines = append(lines, scanner.Text())
		}
	}
	return lines, scanner.Err()
}


// EOF