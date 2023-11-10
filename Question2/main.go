package main

import (
	"fmt"
	"sort"
)

func reorganizeString(s string) string {
	n := len(s)
	charCount := make(map[rune]int)

	// Count the occurrences of each character
	for _, char := range s {
		charCount[char]++
	}

	// Create a slice of character-frequency pairs and sort it by frequency
	pairs := make([]Pair, 0, len(charCount))
	for char, count := range charCount {
		pairs = append(pairs, Pair{char, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})

	// Check if it's possible to rearrange characters
	if pairs[0].Count > (n+1)/2 {
		return ""
	}

	// Rearrange characters
	result := make([]rune, n)
	index := 0
	for _, pair := range pairs {
		for i := 0; i < pair.Count; i++ {
			if index >= n {
				index = 1
			}
			result[index] = pair.Char
			index += 2
		}
	}

	return string(result)
}

// Pair represents a character and its count
type Pair struct {
	Char  rune
	Count int
}

func main() {
	// Example 1
	s1 := "aab"
	fmt.Println("Input:", s1)
	fmt.Println("Output:", reorganizeString(s1))

	// Example 2
	s2 := "aaab"
	fmt.Println("Input:", s2)
	fmt.Println("Output:", reorganizeString(s2))
}
