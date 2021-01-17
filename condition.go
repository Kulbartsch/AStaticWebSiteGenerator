// conditions

package main

import (
	"strings"
)

type CondType string

const (
	None  CondType = "NONE"
	Set            = "SET"
	Equal          = "EQUAL"
)

// parse conditions
func parseCondition(command string) (result bool) {
	var function, parameter1, parameter2 string
	cond := strings.Trim(command, " \t\n)")
	i := strings.IndexAny(cond, " \t")
	if i == -1 {
		function = strings.ToUpper(cond)
	} else {
		function = strings.ToUpper(WhiteSpaceTrim(cond[:i]))
		parameters := WhiteSpaceTrim(cond[i:])
		j := strings.IndexAny(parameters, " \t")
		if j == -1 {
			parameter1 = parameters
			parameter2 = ""
		} else {
			parameter1 = WhiteSpaceTrim(parameters[:j])
			parameter2 = WhiteSpaceTrim(parameters[j:])
		}
	}
	switch function {
	case "COND-IF-SET":
		if len(parameter1) == 0 {
			Message("", 0, "E", "condition IF-SET parameter missing")
			return true
		}
		if len(parameter2) != 0 {
			Message("", 0, "W", "condition IF-SET extra parameter(s) ignored")
		}
		conditionSet(Set, parameter1, "", false)
	case "COND-IF-NOT-SET":
		if len(parameter1) == 0 {
			Message("", 0, "E", "condition IF-NOT-SET parameter missing")
			return true
		}
		if len(parameter2) != 0 {
			Message("", 0, "W", "condition IF-NOT-SET extra parameter(s) ignored")
		}
		conditionSet(Set, parameter1, "", true)
	case "COND-IF-EQUAL":
		if len(parameter1) == 0 {
			Message("", 0, "E", "condition IF-EQUAL parameter missing")
			return true
		}
		conditionSet(Equal, parameter1, parameter2, false)
	case "COND-IF-NOT-EQUAL":
		if len(parameter1) == 0 {
			Message("", 0, "E", "condition IF-NOT-EQUAL parameter missing")
			return true
		}
		conditionSet(Equal, parameter1, parameter2, true)
	case "COND-ELSE":
		if len(parameter1) != 0 {
			Message("", 0, "W", "condition ELSE extra parameter(s) ignored")
		}
		conditionNegate()
	case "COND-END":
		if len(parameter1) != 0 {
			Message("", 0, "W", "condition END extra parameter(s) ignored")
		}
		conditionSet(None, "", "", false)
	default:
		// no condition command
		return false
	}
	conditionValidate()
	return true
}

func conditionSet(t CondType, p1 string, p2 string, not bool) {
	siteContext.vars.SetVar("aswsg-cond-type", string(t))
	siteContext.vars.SetVar("aswsg-cond-var", p1)
	siteContext.vars.SetVar("aswsg-cond-val", p2)
	if not {
		siteContext.vars.SetVar("aswsg-cond-not", "T")
	} else {
		siteContext.vars.SetVar("aswsg-cond-not", "F")
	}
}

func conditionNegate() {
	if IsVarTrue("aswsg-cond-not") {
		siteContext.vars.SetVar("aswsg-cond-not", "F")
	} else {
		siteContext.vars.SetVar("aswsg-cond-not", "T")
	}
}

func conditionValidate() {
	var condi bool
	switch siteContext.vars.GetVal("aswsg-cond-type") {
	case Set:
		if siteContext.vars.ExistsVal(siteContext.vars.GetVal("aswsg-cond-var")) {
			condi = true
		} else {
			condi = false
		}
	case Equal:
		if  siteContext.vars.ExistsVal(siteContext.vars.GetVal("aswsg-cond-var")) &&
			siteContext.vars.ExistsVal("aswsg-cond-val") &&
			siteContext.vars.GetVal(siteContext.vars.GetVal("aswsg-cond-var")) ==
			siteContext.vars.GetVal("aswsg-cond-val") {
			condi = true
		} else {
			condi = false
		}
	default: // no (valid) condition is always fulfilled
		siteContext.conditionFulfilled = true
		return
	}
	if IsVarTrue("aswsg-cond-not") {
		if condi {
			condi = false
		} else {
			condi = true
		}
	}
	siteContext.conditionFulfilled = condi
}

// EOF
