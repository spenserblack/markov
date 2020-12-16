package main

import (
	"fmt"
	"strings"
)

type Markov struct {
	Chain map[string][]string
}

func main() {
	markov := NewSentence("hello goodbye hello go hello world", 1)
	fmt.Println(markov)
}

func NewSentence(words string, prefixLen int) Markov {
	chain := make(map[string][]string)
	markov := Markov{chain}

	splitWords := strings.Split(words, " ")

	for i, suffix := range splitWords[prefixLen:] {
		prefix := strings.Join(splitWords[i:i+prefixLen], " ")

		if suffixes, ok := chain[prefix]; ok {
			chain[prefix] = append(suffixes, suffix)
		} else {
			chain[prefix] = make([]string, 1)
			chain[prefix][0] = suffix
		}
	}

	return markov
}
