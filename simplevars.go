// variable handling

package main

// IDEA: maybe put this in own library

import "strings"

// SimpleVars is a Structure holding vars
type SimpleVars map[string]string

// simpleVar handling

// SetVar key to val
func (v SimpleVars) SetVar(key, val string) (ok bool) {
	tkey := WhiteSpaceTrim(key)
	if len(tkey) == 0 {
		return false
	}
	v[strings.ToUpper(tkey)] = WhiteSpaceTrim(val)
	return true
}

// GetVal get value for key or empty string
func (v SimpleVars) GetVal(key string) (result string) {
	tkey := WhiteSpaceTrim(key)
	if len(tkey) == 0 {
		return ""
	}
	s, r := v[strings.ToUpper(key)]
	if r == false {
		// don't use message to avoid endless loops on error
		// Message("", 0, "W", "Key '"+key+"' does not exist.") /
		println("Warning: Key '" + key + "' does not exist.")
	}
	return s
}

// ExistsVal checks key for existence
func (v SimpleVars) ExistsVal(key string) (result bool) {
	tkey := WhiteSpaceTrim(key)
	if len(tkey) == 0 {
		return false
	}
	_, result = v[strings.ToUpper(key)]
	return
}

// ParseAndSetVar parse value setting line and does so.
func (v SimpleVars) ParseAndSetVar(toparse string) (ok bool) {
	dp := strings.Index(toparse, ":")
	if dp < 1 || dp == len(toparse) {
		Message("", 0, "W", "Could not set variable. Colon missing?")
		return false
	}
	v.SetVar(toparse[0:(dp)], toparse[(dp+1):])
	return true
}

func IsVarTrue(variable string) bool {
	val := strings.ToUpper(siteContext.vars.GetVal(variable))
	if val == "T" || val == "TRUE" {
		return true
	} else {
		return false
	}
}

// EOF
