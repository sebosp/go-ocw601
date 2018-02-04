package main

import (
	"reflect"
	"testing"
)

func TestSplitToken(t *testing.T) {
	cases := []struct {
		inputString string
		start       int
		end         int
		result      string
	}{
		{"aoeu", 0, 1, "a"},
		{"aoeu", 0, -2, "ao"},
		{"aoeu", 3, 3, ""},
		{"aoeu", 2, 4, "eu"},
	}
	for _, c := range cases {
		res, err := substr(
			c.inputString,
			c.start,
			c.end,
		)
		if err != nil {
			t.Errorf(
				"Failed to split input '%s':  '%s'",
				c.inputString,
				err,
			)
		}
		if res != c.result {
			t.Errorf("Expecting %s, got %s", c.result, res)
		}
	}
	_, err := substr("aoeu", 3, 2)
	if err == nil {
		t.Errorf("Expecting error on start > end")
	}
	_, err = substr("aoeu", 10, 2)
	if err == nil {
		t.Errorf("Expecting error on end > len")
	}
}
func TestTokenize(t *testing.T) {
	cases := []struct {
		input  string
		result []string
	}{
		{"()", []string{"(", ")"}},
		{"(  )", []string{"(", ")"}},
		{"( ( ) )", []string{"(", "(", ")", ")"}},
		{"(1+2)", []string{"(", "1", "+", "2", ")"}},
		{"(a=2)", []string{"(", "a", "=", "2", ")"}},
	}
	for _, c := range cases {
		res, err := tokenize(c.input)
		if err != nil {
			t.Errorf(
				"Failed to tokenize input '%s':  '%s'",
				c.input,
				err,
			)
		}
		if !reflect.DeepEqual(res, c.result) {
			t.Errorf("Expecting '%#v', got '%#v'", c.result, res)
		}
	}
}
