package main

import (
	"errors"
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
		{"1 1", true},
		{"( 1 + 2(", true},
		{"( + 2)", true},
		{"( 1 +)", true},
		{"( 1 + + )", true},
		{"+ +", true},
		{")", true},
		{"())", true},
		{"( 1 + 2) (", true},
		{"((", true},
		{"))", true},
		{"+", true},
		{"(e =)", true},
		{"(  +", true},
		{"( + + )", true},
		{"( 1 + + )", true},
		{"1 2", true},
		{"( a = 2 )", false},
		{"( 1 + 3 )", false},
		{"( 1 + ( ( 2 + 3 ) + ( a / b ) ) )", false},
		{"1", false},
	}
	for _, c := range cases {
		Root := &Node{Value: ""}
		expTokens, err := tokenize(c.input + "\n")
		if err != nil {
			t.Errorf("Unable to tokenize '%s': %s", c.input, err)
			continue
		}
		err = Root.BuildTokenTree(expTokens)
		if c.isError && err == nil {
			t.Errorf("Expecting error on %#v, got success", c.input)
		}
		//fmt.Printf("%s = ", c.input)
		//Root.Print()
	}
}

func TestRunTree(t *testing.T) {
	env := map[string]float64{}
	cases := []struct {
		input   string
		output  float64
		isError bool
	}{
		{"( n = 2 )", 2.0, false},
		{"( b = (1 + 1 ) )", 2.0, false},
		{"( a = (4 * n ) )", 8.0, false},
		{"( 1 + n )", 3.0, false},
		{"( 1 + ( ( 2 + 3 ) + ( a / b ) ) )", 10.0, false},
		{"( 1 + ( z = 2) )", 0.0, true},
		{"+", 0.0, true},
		{"(=)", 0.0, true},
	}
	for _, c := range cases {
		err := errors.New("Placeholder error")
		out := 0.0
		Root := &Node{Value: ""}
		expTokens, err := tokenize(c.input + "\n")
		if err != nil {
			t.Errorf("Tokenize error %s\n", err)
			continue
		}
		err = Root.BuildTokenTree(expTokens)
		if err != nil && !c.isError {
			t.Errorf("Tokenize error %s\n", err)
			continue
		}
		out, err = Root.RunTree(env)
		if (c.isError && err == nil) || (!c.isError && err != nil) {
			expecting := "error"
			got := "success"
			if !c.isError {
				expecting = "success"
			}
			if err != nil {
				expecting = "error"
			}
			Root.Print()
			t.Errorf("Expecting %s on %s, got %s", expecting, c.input, got)
		} else {
			// Check if assignment took place
			if Root.Value == "=" {
				if env[Root.Left.Value] != c.output {
					t.Errorf("On assigment '%s' expected %f got: %f", c.input, c.output, out)
				}
			}
			if c.output != out {
				t.Errorf("On input '%s' expected %f got: %f", c.input, c.output, out)
			}
		}
	}
	envCases := []struct {
		key   string
		value float64
	}{
		{"n", 2.0},
		{"b", 2.0},
		{"a", 8.0},
	}
	for _, c := range envCases {
		storedVal, exists := env[c.key]
		if !exists {
			t.Errorf("Key '%s' not found", c.key)
		}
		if storedVal != c.value {
			t.Errorf("Expected env[%s]=%f, got: %f", c.key, c.value, c.value)
		}
	}
}
