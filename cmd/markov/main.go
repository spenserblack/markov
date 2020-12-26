// Use a Markov chain to generate randomized sentences.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Markov chain container for creating a sentence.
type SentenceMarkov struct {
	sync.Mutex
	chain         map[string][]string
	chainStarters []string
	prefixLen     int
}

// To be implemented by types, specifically Markov chains, that generate a
// random string output.
type Generator interface {
	// Generate a random output using a Markov chain.
	Generate() string
}

func main() {
	var markov Generator
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

	markov = NewSentence(strings.Split(feed, "\n"), *prefixLen)
	fmt.Println(markov.Generate())
}

func init() {
	rand.Seed(time.Now().Unix())
}

// Generate a random sentence from the Markov chain.
func (markov SentenceMarkov) Generate() string {
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
func NewSentence(sentences []string, prefixLen int) SentenceMarkov {
	markov := SentenceMarkov{}
	markov.chain = make(map[string][]string)
	markov.prefixLen = prefixLen
	var waiter sync.WaitGroup

	for _, words := range sentences {
		// Let waiter know that goroutine will start
		waiter.Add(1)

		go func(sentence string) {
			// Let waiter know that goroutine has finished
			defer waiter.Done()

			splitWords := strings.Split(sentence, " ")

			for i, suffix := range splitWords[prefixLen:] {
				prefix := strings.Join(splitWords[i:i+prefixLen], " ")

				markov.Lock()
				if i == 0 {
					markov.chainStarters = append(markov.chainStarters, prefix)
				}

				if suffixes, ok := markov.chain[prefix]; ok {
					markov.chain[prefix] = append(suffixes, suffix)
				} else {
					markov.chain[prefix] = []string{suffix}
				}
				markov.Unlock()
			}
		}(words)
	}

	waiter.Wait()

	return markov
}
