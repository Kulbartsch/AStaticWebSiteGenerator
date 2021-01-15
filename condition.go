// conditions

package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)


// parse conditions
func parseCondition(command string) (result bool) {
	var function, parameter1, parameter2 string
	cond := strings.Trim(command, " \t\n)")
	i := strings.IndexAny(cond, " \t")
	if i == -1 {
		function = strings.ToUpper(cmd)
	} else {
		function = strings.ToUpper(WhiteSpaceTrim(cmd[:i]))
		parameters := WhiteSpaceTrim(cmd[i:])
		j := strings.IndexAny(parameters, " \t")
		if j == -1 {
			parameter1 = parameters
		} else {
			parameter1 = WhiteSpaceTrim(parameters[:j])
			parameter2 = WhiteSpaceTrim(parameters[j:])
		}
	}
	switch function {
	case "COND-IF-SET":
		result = conditionIfSet(parameter1, parameter2)
	case "COND-IF-NOT-SET":
		result = conditionIfNotSet(parameter1, parameter2)
	case "COND-IF-EQUAL":
		result = conditionIfEqual(parameter1, parameter2)
	case "COND-IF-NOT-EQUAL":
		result = conditionIfNotEqual(parameter1, parameter2)
	case "COND-ELSE":
		result = conditionElse(parameter1, parameter2)
	case "COND-END":
		result = conditionEnd(parameter1, parameter2)
	default:
		Message("", 0, "W", "unknown condition command ignored (function/parameter): "+function+"/"+parameter+" = "+cmd)
		break
	}
	conditionValidate();
	return true
}
