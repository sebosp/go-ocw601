package main

import (
	"bufio"
	"fmt"
	"os"
)

// substr splits a string from start to end indexes
func substr(s string, start int, end int) (string, error) {
	if end < 0 {
		end = len(s) + end
	}
	if end > len(s) || start > end {
		return "", fmt.Errorf(
			"Bad params: len=%d, start=%d, end=%d",
			len(s), start, end,
		)
	}
	res := ""
	if start != end {
		res = s[start:end]
	}
	return res, nil
}

//tokenize transforms an input into individual math operators
func tokenize(exp string) ([]string, error) {
	var tokens []string
	literalStart := 0
	literalEnd := 0
	validTokens := [7]string{`(`, `)`, `+`, `-`, `*`, `/`, `=`}
	for _, expCharOrd := range exp {
		literalEnd++
		expChar := string(expCharOrd)
		if expChar == " " || expChar == "\n" {
			res, err := substr(exp, literalStart, literalEnd-1)
			literalStart = literalEnd
			if err != nil {
				return []string{}, err
			}
			if len(res) > 0 {
				tokens = append(tokens, res)
			}
		}
		for _, tokenCharOrd := range validTokens {
			tokenChar := string(tokenCharOrd)
			if expChar == tokenChar {
				res, err := substr(exp, literalStart, literalEnd-1)
				if err != nil {
					return []string{}, err
				}
				if len(res) > 0 {
					tokens = append(tokens, res)
				}
				tokens = append(tokens, expChar)
				literalStart = literalEnd
				break
			}
		}
	}
	return tokens, nil
}

func main() {
	env := make(map[string]*Node)
	reader := bufio.NewReader(os.Stdin)
	text := ""
	output := 0.0
	for {
		fmt.Printf("%% ")
		text, _ = reader.ReadString('\n')
		if text == "quit\n" {
			break
		}
		res, err := tokenize(text)
		if err != nil {
			fmt.Printf("Tokenize error %s\n", err)
			continue
		}
		Root := &Node{Value: "", tokens: res}
		err = Root.BuildTokenTree()
		if err != nil {
			fmt.Printf("BuildTokenTree error %s\n", err)
			continue
		}
		output, err = Root.RunTree(env)
		if err != nil {
			fmt.Printf("RunTree error %s\n", err)
			continue
		}
		fmt.Printf("result: %f\n", output)
		fmt.Printf("env = %+v\n", env)
	}
}
