package main

import (
	"strconv"
	"strings"
)

// "ASWSG-TABLE-HEADERLINES":    "1",      // number of headers, when parsing a table.
// "ASWSG-TABLE-ALIGNMENT":      "LL",     // L -> <th style="text-align:left">, C = center, R = right, other/default = left

func parseTableLine(line string) string {
	siteContext.tableLine += 1
	if Right(line, 1) == siteContext.vars.GetVal("ASWSG-TABLE") {
		line = line[:len(line)-1]
	}
	cells := parseTableCells(line)
	hl, _ := strconv.Atoi(siteContext.vars.GetVal("ASWSG-TABLE-HEADERLINES"))
	if siteContext.tableLine <= hl {
		return bulidTableRow(cells, "th")
	} else {
		return bulidTableRow(cells, "td")
	}
}

func parseTableCells(line string) []string {
	fields := strings.Split(line, siteContext.vars.GetVal("ASWSG-TABLE"))
	var cells []string
	for _, f := range fields {
		cells = append(cells, parseInLine(strings.Trim(f, " \t")))
	}
	return cells
}

func getColumnAligmnet(column int) string {
	ta := siteContext.vars.GetVal("ASWSG-TABLE-ALIGNMENT")
	if column >= len(ta) {
		Message("", 0, "I", "Table alignment to short: "+ta+". Using default L.")
		return "left"
	}
	a := ta[column]
	switch a {
	case 'C', 'c':
		return "center"
	case 'R', 'r':
		return "right"
	case 'L', 'l':
		return "left"
	}
	Message("", 0, "I", "Invalid table alignment value: "+string(a)+". Allowed are LCR.")
	return "left"
}

func bulidTableRow(cells []string, tag string) string {
	row := "<tr>"
	for i, c := range cells {
		a := HTMLAttrib{"style": "text-align:" + getColumnAligmnet(i)}
		row = row + surroundWithHTMLTagWithAttributes(tag, c, a)
	}
	return row + "</tr>"
}

// EOF
