// Use a Markov chain to generate randomized sentences.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Markov chain container.
type Markov struct {
	chain         map[string][]string
	chainStarters []string
	prefixLen     int
}

func main() {
	prefixLen := flag.Int("n", 1, "length of words to use as a key")
	printHelp := flag.Bool("h", false, "print this help message")
	flag.Parse()

	feed := flag.Arg(0)

	if *printHelp || feed == "" {
		println("usage: markov [OPTIONS] [TEXT]")
		println("TEXT: the text to feed to the markov chain")
		flag.PrintDefaults()
		return
	}

	markov := NewSentence(strings.Split(feed, "\n"), *prefixLen)
	fmt.Println(markov.Generate())
}

func init() {
	rand.Seed(time.Now().Unix())
}

// Generate a random sentence from the Markov chain.
func (markov Markov) Generate() string {
	starter := markov.chainStarters[rand.Intn(len(markov.chainStarters))]
	output := starter

	for {
		splitWords := strings.Split(output, " ")
		lastWords := splitWords[len(splitWords)-markov.prefixLen:]
		key := strings.Join(lastWords, " ")
		nextValues := markov.chain[key]

		if len(nextValues) == 0 {
			return output
		}

		output = fmt.Sprintf("%v %v", output, nextValues[rand.Intn(len(nextValues))])
	}
}

// Feed data to create a Markov chain.
// Each sentence in `sentences` should be a string of space-separated words to
// be used when building the chain -- order of the words determines how each next
// word in a generated sentence is decided.
// `prefixLen` is the number of words to be used as a "key" to deciding the next
// word. For example, if `prefixLen` is 2 and the generated text is "I made a
// chain" then "I made" was a key to "a" and "made a" was a key to "chain" in
// the sentence.
func NewSentence(sentences []string, prefixLen int) Markov {
	chain := make(map[string][]string)
	chainStarters := make([]string, 0)

	for _, words := range sentences {
		splitWords := strings.Split(words, " ")

		var lastPrefix string

		for i, suffix := range splitWords[prefixLen:] {
			prefix := strings.Join(splitWords[i:i+prefixLen], " ")
			lastPrefix = prefix

			if i == 0 {
				chainStarters = append(chainStarters, prefix)
			}

			if suffixes, ok := chain[prefix]; ok {
				chain[prefix] = append(suffixes, suffix)
			} else {
				chain[prefix] = []string{suffix}
			}
		}

		chain[lastPrefix] = append(chain[lastPrefix], "")
	}

	chain[""] = make([]string, 0, 0)

	markov := Markov{chain, chainStarters, prefixLen}

	return markov
}
