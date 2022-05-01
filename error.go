package main

import (
	"fmt"
	"strconv"
	"strings"
)

// FindUserCode finds the line number range of code strings
func FindUserCode(code []string) (start, end int) {
	if len(code) == 0 {
		return
	}

	for start = 0; start < len(code); start++ {
		if strings.Contains(code[start], "---- START USER CODE ----") {
			break
		}
	}

	for end = len(code) - 1; end >= 0; end-- {
		if strings.Contains(code[end], "---- END USER CODE ----") {
			break
		}
	}

	return
}

// PrintErrorLine prints out the code surrounding an error from the compiler
func PrintErrorLine(sourceCode string, errText string) bool {
	parts := strings.SplitN(errText, ":", 4)
	if len(parts) != 4 || !strings.Contains(parts[0], "main.go") {
		return false // Wrong number of parts / filename is incorrect
	}
	code := strings.Split(sourceCode, "\n")
	start, end := FindUserCode(code)
	line, err := strconv.Atoi(parts[1])
	if err == nil {
		line-- // Lines are 1 indexed
	}
	if err != nil || line <= start || line >= end {
		return false // Couldn't parse line, or line out of user code range
	}
	fmt.Println(code[line])
	col, err := strconv.Atoi(parts[2])
	if err != nil || col <= 0 || col >= len(code[line]) {
		return false // Couldn't parse col;, or col out of user code line length
	}
	indent := ""
	for _, chr := range code[line][:col] {
		c := string(chr)
		if c == " " || c == "\t" {
			indent += c
		} else {
			indent += " "
		}
	}
	fmt.Printf("%s|\n", indent)
	fmt.Printf("%s\\- %s (line %d)\n", indent, parts[3], line-start)
	return true
}
