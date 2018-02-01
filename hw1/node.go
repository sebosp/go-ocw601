package main

import (
	"errors"
	"fmt"
)

// Node specifies a tree structure that will have the math tree stored
type Node struct {
	Value string
	Left  *Node
	Right *Node
	Top   *Node
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

func (n *Node) printLeaves(depth int) {
	if n.Left != nil {
		n.Left.printLeaves(depth + 1)
	}
	fmt.Printf(" '%s'[%d] ", n.Value, depth)
	if n.Right != nil {
		n.Right.printLeaves(depth + 1)
	}
}

// Print the state of the tree
func (n *Node) Print() {
	fmt.Println("Printing tree state")
	n.printLeaves(0)
	fmt.Println("")
}

// BuildTokenTree based on the expTokens input, which is
// the output from tokenize
func (n *Node) BuildTokenTree(expTokens []string) error {
	temp := n
	pos := 0
	err := errors.New("Placeholder error")
	fmt.Printf("Working on %+v\n", expTokens)
	for _, token := range expTokens {
		switch token {
		case `(`:
			temp, err = temp.Insert(`.`)
			if err != nil {
				return fmt.Errorf("Unable to descend on '(' at pos: %d", pos)
			}
		case `)`:
			if temp.Top != nil {
				temp = temp.Top
			} else {
				return fmt.Errorf("Unable to ascend to parent on ')' at pos: %d", pos)
			}
		case `+`, `-`, `*`, `/`:
			if temp.Left == nil {
				return fmt.Errorf("Unexpected operator on end-leave at pos: %d", pos)
			}
			temp.Value = token
			temp, err = temp.Insert(`.`)
			if err != nil {
				return fmt.Errorf("Unable to create right operand at pos: %d", pos)
			}
		default:
			if temp.Left != nil || temp.Right != nil {
				return fmt.Errorf("Unexpected value on non end-leave at pos: %d", pos)
			}
			if temp.Top == nil {
				return fmt.Errorf("Unexpected top position for a value at pos: %d", pos)
			}
			temp.Value = token
			temp = temp.Top
		}
		pos++
	}
	return nil
}
