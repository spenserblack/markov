// Package sentence is a utility to create a random sentence generator
package sentence

import (
	"errors"
	"github.com/spenserblack/markov/pkg/chain"
	"strings"
	"sync"
)

type sentenceGenerator struct {
	generator *chain.ByteGenerator
}

// Generate returns a generator of random words that make up a sentence, using
// the Markov chain.
func (generator *sentenceGenerator) Generate() func() (next string, stop error) {
	g := generator.generator.Generate()

	return func() (next string, stop error) {
		if bytes := g(); bytes != nil {
			next = string(bytes)
		} else {
			stop = errors.New("Generation has completed")
		}
		return
	}
}

// New feeds data to a markov chain and returns the sentence generator.
//
// Each sentence in `sentences` should be a string of space-separated words to
// be used when building the chain -- order of the words determines how each next
// word in a generated sentence is decided.
// `prefixLen` is the number of words to be used as a "key" to deciding the next
// word. For example, if `prefixLen` is 2 and the generated text is "I made a
// chain" then "I made" was a key to "a" and "made a" was a key to "chain" in
// the sentence.
func New(sentences []string, prefixLen int) (generator *sentenceGenerator, err error) {
	g := new(sentenceGenerator)

	bytes := make([][][]byte, len(sentences), len(sentences))
	var waiter sync.WaitGroup

	for i, sentence := range sentences {
		waiter.Add(1)
		go func(index int, sentence string) {
			defer waiter.Done()

			words := strings.Split(sentence, " ")
			bytes[index] = make([][]byte, 0, len(words))

			for _, word := range words {
				bytes[index] = append(bytes[index], []byte(word))
			}

		}(i, sentence)
	}

	waiter.Wait()

	g.generator, err = chain.NewByteGenerator(bytes, prefixLen)
	generator = g

	return
}
