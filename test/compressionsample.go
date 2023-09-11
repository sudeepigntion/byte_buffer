package main

import (
	"fmt"
	"strings"
)

// Run-Length Encode a string
func compress(input string) string {
	var result strings.Builder
	count := 1

	for i := 1; i < len(input); i++ {
		if input[i] == input[i-1] {
			count++
		} else {
			result.WriteByte(input[i-1])
			result.WriteString(fmt.Sprintf("%d", count))
			count = 1
		}
	}

	// Append the last character and its count
	result.WriteByte(input[len(input)-1])
	result.WriteString(fmt.Sprintf("%d", count))

	return result.String()
}

// Run-Length Decode a compressed string
func decompress(input string) string {
	var result strings.Builder
	i := 0

	for i < len(input) {
		char := input[i]
		i++

		// Find the count of the character
		count := 0
		for i < len(input) && input[i] >= '0' && input[i] <= '9' {
			count = count*10 + int(input[i]-'0')
			i++
		}

		// Append the character count times to the result
		result.WriteString(strings.Repeat(string(char), count))
	}

	return result.String()
}

func main() {
	original := "AAABBBCCCCDDDD"
	fmt.Println("Original:  ", original)

	compressed := compress(original)
	fmt.Println("Compressed:", compressed)

	decompressed := decompress(compressed)
	fmt.Println("Decompressed:", decompressed)
}
