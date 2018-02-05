package main

import (
	"errors"
	"fmt"
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
		expTokens, err := tokenize(c.input + "\n")
		if err != nil {
			t.Errorf("Unable to tokenize '%s': %s", c.input, err)
			continue
		}
		Root := &Node{Value: "", tokens: expTokens}
		err = Root.BuildTokenTree()
		if c.isError && err == nil {
			t.Errorf("Expecting error on %#v, got success", c.input)
		}
		//fmt.Printf("%s = ", c.input)
		//Root.Print()
	}
}

func TestRunTree(t *testing.T) {
	env := map[string]*Node{}
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
		expTokens, err := tokenize(c.input + "\n")
		if err != nil {
			t.Errorf("Tokenize error %s\n", err)
			continue
		}
		Root := &Node{Value: "", tokens: expTokens}
		err = Root.BuildTokenTree()
		if err != nil && !c.isError {
			t.Errorf("Tokenize error %s\n", err)
			continue
		}
		fmt.Printf("Running test %s\n", c.input)
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
			fmt.Println(Root.String())
			t.Errorf("Expecting %s on %s, got %s", expecting, c.input, got)
		} else {
			// Check if assignment took place
			if Root.Value == "=" && !c.isError {
				envRoot := env[Root.Left.Value]
				if envRoot.Value == "" {
					envRoot.BuildTokenTree()
				}
				envOut, envErr := envRoot.RunTree(env)
				if envErr != nil {
					t.Errorf("env[%s] failed: %s", Root.Left.Value, envErr)
				}
				if envOut != c.output {
					t.Errorf("On assigment '%s' expected %f got: %f", c.input, c.output, envOut)
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
		envRoot, exists := env[c.key]
		if !exists {
			t.Errorf("Key '%s' not found", c.key)
		}
		if envRoot.Value == "" {
			envRoot.BuildTokenTree()
		}
		envOut, envErr := envRoot.RunTree(env)
		if envErr != nil {
			t.Errorf("env[%s] failed: %s", envRoot.Left.Value, envErr)
		}
		if envOut != c.value {
			t.Errorf("Expected env[%s]=%f, got: %f", c.key, c.value, envOut)
		}
	}
}
func TestString(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{"( n = 2 )", "Assign(Var(n),Num(" + fmt.Sprintf("%f", 2.0) + "))"},
		{"( b = (1 + 1 ) )", "Assign(Var(b),Sum(Num(" + fmt.Sprintf("%f", 1.0) + "),Num(" + fmt.Sprintf("%f", 1.0) + ")))"},
		{"( a = (4 * n ) )", "Assign(Var(a),Prod(Num(" + fmt.Sprintf("%f", 4.0) + "),Var(n)))"},
		{"( 1 + n )", "Sum(Num(" + fmt.Sprintf("%f", 1.0) + "),Var(n))"},
		{"( 1 + ( ( 2 + 3 ) + ( a / b ) ) )", "Sum(Num(" + fmt.Sprintf("%f", 1.0) + "),Sum(Sum(Num(" + fmt.Sprintf("%f", 2.0) + "),Num(" + fmt.Sprintf("%f", 3.0) + ")),Quot(Var(a),Var(b))))"},
	}
	for _, c := range cases {
		expTokens, _ := tokenize(c.input + "\n")
		Root := &Node{Value: "", tokens: expTokens}
		_ = Root.BuildTokenTree()
		out := Root.String()
		if out != c.output {
			t.Errorf("Expecting: \n\t'%s'\n got:\n\t'%s'", c.output, out)
		}
	}
}
