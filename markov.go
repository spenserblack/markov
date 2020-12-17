package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Markov struct {
	Chain         map[string][]string
	chainStarters []string
	prefixLen     int
}

func main() {
	markov := NewSentence("hello goodbye hello go hello world", 1)
	fmt.Println(markov.generate())
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (markov Markov) generate() string {
	starter := markov.chainStarters[rand.Intn(len(markov.chainStarters))]
	output := starter

	for {
		splitWords := strings.Split(output, " ")
		lastWords := splitWords[len(splitWords)-markov.prefixLen:]
		key := strings.Join(lastWords, " ")
		nextValues := markov.Chain[key]

		if len(nextValues) == 0 {
			return output
		}

		output = fmt.Sprintf("%v %v", output, nextValues[rand.Intn(len(nextValues))])
	}
}

func NewSentence(words string, prefixLen int) Markov {
	chain := make(map[string][]string)
	chainStarters := make([]string, 1)

	splitWords := strings.Split(words, " ")

	var lastPrefix string

	for i, suffix := range splitWords[prefixLen:] {
		prefix := strings.Join(splitWords[i:i+prefixLen], " ")
		lastPrefix = suffix

		if i == 0 {
			chainStarters = append(chainStarters, prefix)
		}

		if suffixes, ok := chain[prefix]; ok {
			chain[prefix] = append(suffixes, suffix)
		} else {
			chain[prefix] = []string{suffix}
		}
	}

	chain[lastPrefix] = []string{""}
	chain[""] = make([]string, 0, 0)

	markov := Markov{chain, chainStarters, prefixLen}

	return markov
}
