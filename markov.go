package main

import (
	"fmt"
	"strings"
)

type Markov struct {
	Chain         map[string][]string
	chainStarters []string
}

func main() {
	markov := NewSentence("hello goodbye hello go hello world", 1)
	fmt.Println(markov)
}

func NewSentence(words string, prefixLen int) Markov {
	chain := make(map[string][]string)
	chainStarters := make([]string, 1)

	splitWords := strings.Split(words, " ")

	for i, suffix := range splitWords[prefixLen:] {
		prefix := strings.Join(splitWords[i:i+prefixLen], " ")

		if i == 0 {
			chainStarters = append(chainStarters, prefix)
		}

		if suffixes, ok := chain[prefix]; ok {
			chain[prefix] = append(suffixes, suffix)
		} else {
			chain[prefix] = make([]string, 1)
			chain[prefix][0] = suffix
		}
	}

	markov := Markov{chain, chainStarters}

	return markov
}
