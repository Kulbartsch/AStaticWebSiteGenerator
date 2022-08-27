package main

import (
	"reflect"
	"testing"
)

func TestReverseStringArray(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty", args{[]string{}}, []string{}},
		{"one element", args{[]string{"a"}}, []string{"a"}},
		{"three elements", args{[]string{"a", "b", "c"}}, []string{"c", "b", "a"}},
		{"four elements", args{[]string{"a", "b", "c", "d"}}, []string{"d", "c", "b", "a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseStringArray(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseStringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToValidHtmlAnchor(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add more test cases.
		{"empty", args{""}, ""},
		{"only valid", args{"abc_DEF-123"}, "abc_DEF-123"},
		{"unvalid chars", args{"Fußnote 123.8$"}, "Fu_note_123_8_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToValidHtmlAnchor(tt.args.s); got != tt.want {
				t.Errorf("ToValidHtmlAnchor() = %v, want %v", got, tt.want)
			}
		})
	}
}
