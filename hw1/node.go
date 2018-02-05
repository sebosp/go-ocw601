package main

import (
	"errors"
	"fmt"
	"strconv"
)

// Node specifies a tree structure that will have the math tree stored
type Node struct {
	tokens []string
	Value  string
	Left   *Node
	Right  *Node
	Top    *Node
}

// Insert sets a value, first filling left, then right
func (n *Node) Insert(value string) (*Node, error) {
	if n.Left == nil {
		n.Left = &Node{Value: value}
		n.Left.Top = n
		return n.Left, nil
	}
	if n.Right == nil {
		n.Right = &Node{Value: value}
		n.Right.Top = n
		return n.Right, nil
	}
	return nil, fmt.Errorf("Node is full")
}

// String returns the state of the tree
func (n *Node) String() string {
	ret := ""
	switch n.Value {
	case `+`, `-`, `*`, `/`, `=`:
		switch n.Value {
		case `+`:
			ret = "Sum("
		case `-`:
			ret = "Diff("
		case `*`:
			ret = "Prod("
		case `/`:
			ret = "Quot("
		case `=`:
			ret = "Assign("
		}
		if n.Left != nil {
			ret += n.Left.String()
		}
		ret += ","
		if n.Right != nil {
			ret += n.Right.String()
		}
	default:
		val, err := strconv.ParseFloat(n.Value, 64)
		if err != nil {
			ret = "Var(" + n.Value
		} else {
			ret = "Num(" + fmt.Sprintf("%f", val)
		}
	}
	ret += ")"
	return ret
}

func (n *Node) hasUnsetOperands() bool {
	if n.Value == "" {
		return true
	}
	if n.Left != nil {
		if n.Left.hasUnsetOperands() {
			return true
		}
	}
	if n.Right != nil {
		if n.Right.hasUnsetOperands() {
			return true
		}
	}
	return false
}

// BuildTokenTree based on the 'tokens' root definition, which is
// the output from tokenize
func (n *Node) BuildTokenTree() error {
	temp := n
	pos := 0
	err := errors.New("Placeholder error")
	openParens := 0
	for _, token := range n.tokens {
		switch token {
		case `(`:
			temp, err = temp.Insert("")
			if err != nil {
				return fmt.Errorf("Unable to descend on '(' at pos %d: %s", pos, err)
			}
			openParens++
		case `)`:
			if temp.Top != nil {
				temp = temp.Top
			} else {
				if temp.Value != `=` {
					return fmt.Errorf("Unable to ascend to parent on ')' at pos %d", pos)
				}
			}
			openParens--
		case `+`, `-`, `*`, `/`, `=`:
			if temp.Top == nil {
				return fmt.Errorf("Unexpected top position for an operand at pos: %d", pos)
			}
			temp = temp.Top
			if temp.Left == nil {
				return fmt.Errorf("Empty left operand for operator %s at pos: %d", token, pos)
			}
			if temp.Value != "" {
				return fmt.Errorf("Operator already set at pos: %d", pos)
			}
			temp.Value = token
			temp, err = temp.Insert("")
			if err != nil {
				return fmt.Errorf("Unable to create right operand at pos %d: %s", pos, err)
			}
		default:
			if temp.Left != nil || temp.Right != nil {
				return fmt.Errorf("Unexpected value on non end-leave at pos: %d", pos)
			}
			if temp.Value != "" {
				return fmt.Errorf("Operand already set at pos: %d", pos)
			}
			temp.Value = token
		}
		pos++
	}
	if openParens > 0 {
		return fmt.Errorf("Non matching parens")
	}
	if n.hasUnsetOperands() {
		return fmt.Errorf("Unset Operands found")
	}
	return nil
}

// RunTree recurses through the tree to resolve it
func (n *Node) RunTree(env map[string]*Node) (float64, error) {
	err := errors.New("Placeholder error")
	switch n.Value {
	case `=`:
		if n.Top != nil {
			return 0.0, fmt.Errorf("Assignment on non-root")
		}
		if n.Left == nil || n.Right == nil {
			return 0.0, fmt.Errorf("Not enough operands")
		}
		// We need to drop the outer containing parens and the left val assign: '(','x','=', so the length of the array must be more than 4
		if len(n.tokens) < 5 {
			return 0.0, fmt.Errorf("Not enough operands")
		}
		envRoot, exists := env[n.Left.Value]
		if exists || envRoot == nil {
			// Overwrite the current ref.
			env[n.Left.Value] = &Node{Value: "", tokens: n.tokens[3 : len(n.tokens)-1]}
			envRoot = env[n.Left.Value]
		}
		if envRoot.Value == "" {
			envRoot.BuildTokenTree()
		}
		envOut, envErr := envRoot.RunTree(env)
		if envErr != nil {
			return 0.0, fmt.Errorf("env[%s] failed: %s", n.Left.Value, envErr)
		}
		return envOut, nil
	case `+`, `-`, `*`, `/`:
		left := 0.0
		right := 0.0
		if n.Left == nil || n.Right == nil {
			return 0.0, fmt.Errorf("Not enough operands")
		}
		left, err = n.Left.RunTree(env)
		if err != nil {
			return 0.0, fmt.Errorf("Error on left operand: %s", err)
		}
		right, err = n.Right.RunTree(env)
		if err != nil {
			return 0.0, fmt.Errorf("Error on right operand: %s", err)
		}
		switch n.Value {
		case `+`:
			return left + right, nil
		case `-`:
			return left - right, nil
		case `*`:
			return left * right, nil
		case `/`:
			return left / right, nil
		}
	default:
		if n.Left != nil || n.Right != nil {
			return 0.0, fmt.Errorf("%s is not an operand", n.Value)
		}
		val := 0.0
		val, err = strconv.ParseFloat(n.Value, 64)
		if err != nil {
			envRoot, exists := env[n.Value]
			if !exists || envRoot == nil {
				return 0.0, fmt.Errorf("env[%s] unset", n.Value)
			}
			if envRoot.Value == "" {
				envRoot.BuildTokenTree()
			}
			envOut, envErr := envRoot.RunTree(env)
			if envErr != nil {
				return 0.0, fmt.Errorf("Unable to resolve env[%s]: %s", n.Value, envErr)
			}
			return envOut, nil
		}
		return val, nil
	}
	return 0.0, fmt.Errorf("Unexpected end of function reach on n.Value = '%s'", n.Value)
}
