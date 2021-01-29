// Package word is a utility to create a random sentence generator
package word

import (
	"errors"
	"github.com/spenserblack/markov/pkg/chain"
	"sync"
	"unicode/utf8"
)

type wordGenerator struct {
	chain *chain.ByteChain
}

// StopIteration signifies that the generator should stop
var StopIteration error = errors.New("Generation has completed")

// Generate returns a generator of random runes using the Markov chain.
//
// Returns a StopIteration error if/when generation has completed.
func (generator *wordGenerator) Generate() func() (next rune, stop error) {
	g := generator.chain.Generate()

	return func() (next rune, stop error) {
		bytes := g()
		if bytes == nil {
			return next, StopIteration
		}
		next, _ = utf8.DecodeRune(bytes)

		if next == utf8.RuneError {
			stop = errors.New("Could not decode bytes to rune. Was valid UTF-8 used?")
		}
		return
	}
}

// New feeds data to a markov chain and return the word generator.
//
// Each word in `words` should be a string of letters to be used when building the
// chain -- order of the letters determines how each next letter in a generated
// word is decided.
// `prefixLen` is the number of letters to be used as a "key" to deciding the next
// letter. For example, if `prefixLen` is 2 and the generated text is "abcd" then
// "ab" was a key to "c" and "bc" was a key to "d" in the word.
func New(words []string, prefixLen int) (generator *wordGenerator, err error) {
	g := new(wordGenerator)

	bytes := make([][][]byte, len(words), len(words))

	var waiter sync.WaitGroup

	for i, word := range words {
		waiter.Add(1)
		go func(index int, word string) {
			defer waiter.Done()

			runes := []rune(word)
			runesAsBytes := make([][]byte, 0, len(runes))

			for _, r := range runes {
				runeLen := utf8.RuneLen(r)
				buf := make([]byte, runeLen, runeLen)
				utf8.EncodeRune(buf, r)
				runesAsBytes = append(runesAsBytes, buf)
			}

			bytes[index] = runesAsBytes
		}(i, word)
	}

	waiter.Wait()

	g.chain, err = chain.NewByteChain(bytes, prefixLen)
	generator = g

	return
}
