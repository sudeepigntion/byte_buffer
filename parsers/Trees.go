package parsers

import "strings"

type TreeNode struct {
	ArrayCount int
	Name       string
	Value      string
	Children   []*TreeNode
}

type RootClass struct {
	Name       string
	ArrayCount int
}

// StructField represents a struct field with name and type.
type StructField struct {
	Name string
	Type string
}

func CountSquareBrackets(input string) (bool, int) {
	// Split the input into characters
	characters := strings.Split(input, "")

	// Initialize counters for open and close square brackets
	openCount := 0
	closeCount := 0

	// Iterate through the characters and count square brackets
	for _, char := range characters {
		if char == "[" {
			openCount++
		} else if char == "]" {
			closeCount++
		}
	}

	// Return the minimum count (as open and close brackets must match)
	return (openCount == closeCount), openCount
}
