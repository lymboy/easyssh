package util

import "strings"

// LevenshteinDistance calculates the minimum number of edits
// (insertions, deletions, substitutions) needed to transform s1 into s2.
func LevenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create a matrix to store distances
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}

	// Initialize first row
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill in the rest of the matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// SimilarNames returns names with edit distance <= maxDistance from the given name.
func SimilarNames(name string, candidates []string, maxDistance int) []string {
	var result []string
	for _, candidate := range candidates {
		if LevenshteinDistance(strings.ToLower(name), strings.ToLower(candidate)) <= maxDistance {
			result = append(result, candidate)
		}
	}
	return result
}
