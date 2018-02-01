package main

import (
	"testing"
)

func TestInsert(t *testing.T) {
	fullNode := &Node{Value: `+`}
	leaf, err := fullNode.Insert(`1`)
	if err != nil {
		t.Errorf("Expecting success on left leaf, got: %s", err)
	}
	if fullNode.Left != leaf {
		t.Errorf("Expecting new leaf to be left")
	}
	leaf, err = fullNode.Insert(`2`)
	if err != nil {
		t.Errorf("Expecting success on right leaf, got: %s", err)
	}
	if fullNode.Right != leaf {
		t.Errorf("Expecting new leaf to be right")
	}
	res, err := fullNode.Insert(`error`)
	if err == nil {
		t.Errorf("Expecting error on fullNode, inserted %s", res.Value)
	}
}

func TestBuildTokenTree(t *testing.T) {
	cases := []struct {
		input   string
		isError bool
	}{
		{")", true},
		{"())", true},
		{"( 1 + 2)", false},
		{"( 1 + ( ( 2 + 3 ) + ( a / b ) ) )", false},
	}
	for _, c := range cases {
		Root := &Node{Value: `.`}
		expTokens, err := tokenize(c.input)
		if err != nil {
			t.Errorf("Unable to tokenize '%s'", c.input)
		}
		err = Root.BuildTokenTree(expTokens)
		if c.isError && err == nil {
			t.Errorf("Expecting error on %#v, got success", c.input)
		}
		Root.Print()
	}
}
