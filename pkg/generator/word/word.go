// Package word is a utility to create a random sentence generator
package word

import (
	gen "github.com/spenserblack/markov/pkg/generator"
	"strings"
	"sync"
	"unicode/utf8"
)

type wordGenerator struct {
	generator *gen.ByteGenerator
}

// Generate returns a random word using the Markov chain.
//
// If maxTokens is < 0, then generation will continue until its "natural"
// end from the chain deciding that a token should end the chain.
// Enforcing a maximum number of tokens can be helpful if the chain has a
// chance of generating infinitely, or to simply prevent the generated
// word from being overly long.
func (generator *wordGenerator) Generate(maxTokens int) string {
	var builder strings.Builder

	g := generator.generator.Generate()

	for i, next := 0, g(); i != maxTokens && next != nil; i++ {
		for _, b := range next {
			builder.WriteByte(b)
		}
		next = g()
	}

	return builder.String()
}

// New feeds data to a markov chain and return the word generator.
//
// Each word in `words` should be a string of letters to be used when building the
// chain -- order of the letters determines how each next letter in a generated
// word is decided.
// `prefixLen` is the number of letters to be used as a "key" to deciding the next
// letter. For example, if `prefixLen` is 2 and the generated text is "abcd" then
// "ab" was a key to "c" and "bc" was a key to "d" in the word.
func New(words []string, prefixLen int) (generator gen.StringGenerator, err error) {
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

	g.generator, err = gen.NewByteGenerator(bytes, prefixLen)
	generator = g

	return
}
