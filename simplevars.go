package main

// TODO: put this in own library

import "strings"

// "SimpleVars is a Structure holding vars"
type SimpleVars map[string]string

// simpleVar handling

// "SimpleVars.SetVar key to val"
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
