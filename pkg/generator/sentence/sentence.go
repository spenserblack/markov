// Package sentence is a utility to create a random sentence generator
package sentence

import (
	gen "github.com/spenserblack/markov/pkg/generator"
	"strings"
	"sync"
)

type sentenceGenerator struct {
	generator *gen.ByteGenerator
}

// Generate returns a random sentence using the Markov chain.
//
// If maxTokens is <= 0, then generation will continue until its "natural"
// end from the chain deciding that a token should end the chain.
// Enforcing a maximum number of tokens can be helpful if the chain has a
// chance of generating infinitely, or to simply prevent the generated
// sentence from being overly long.
func (generator *sentenceGenerator) Generate(maxTokens int) string {
	var builder strings.Builder
	c := make(chan []byte)
	tokenCounter := 1

	go generator.generator.Generate(c)

	for bytes := range c {
		for _, b := range bytes {
			builder.WriteByte(b)
		}
		if tokenCounter == maxTokens {
			break
		}
		builder.WriteRune(' ')
		tokenCounter++
	}

	return builder.String()
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
func New(sentences []string, prefixLen int) (generator gen.StringGenerator, err error) {
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

	g.generator, err = gen.NewByteGenerator(bytes, prefixLen)
	generator = g

	return
}
