package main

import (
	"strings"
)

const badWordReplacement = "****"

func getBadWords() []string {
	return []string{"kerfuffle", "sharbert", "fornax"}
}

func replaceBadWords(text, replacement string, badWords []string) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				words[i] = replacement			
			}
		}
	}
	return strings.Join(words, " ")
}
