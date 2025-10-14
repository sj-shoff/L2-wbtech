package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func groupAnagrams(data []string) map[string][]string {
	lowerData := make([]string, len(data))
	for i, word := range data {
		lowerData[i] = strings.ToLower(word)
	}

	anagramMap := make(map[string][]string)

	for _, word := range lowerData {
		runes := []rune(word)
		slices.Sort(runes)
		sorted := string(runes)

		anagramMap[sorted] = append(anagramMap[sorted], word)
	}

	result := make(map[string][]string)

	for _, words := range anagramMap {
		if len(words) <= 1 {
			continue
		}

		uniqueWords := removeDuplicates(words)
		sort.Strings(uniqueWords)

		firstWord := uniqueWords[0]
		result[firstWord] = uniqueWords
	}

	return result
}

func removeDuplicates(words []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(words))

	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
		}
	}

	return result
}

func main() {
	res := groupAnagrams([]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"})
	for key, value := range res {
		fmt.Printf("%s: %v\n", key, value)
	}
}
