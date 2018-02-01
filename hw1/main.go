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
		if expChar == " " || expChar == `\n` {
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
	Root := &Node{Value: `.`}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	res, err := tokenize(text)
	if err != nil {
		fmt.Printf("Got error %s", err)
	} else {
		fmt.Printf("Result = %#v", res)
	}
	err = Root.BuildTokenTree(res)
}
