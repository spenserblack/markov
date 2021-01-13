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
func (generator *wordGenerator) Generate() string {
	var builder strings.Builder

	for _, bytes := range generator.generator.Generate() {
		for _, b := range bytes {
			builder.WriteByte(b)
		}
	}

	return builder.String()
}

// LimitedGenerate return a random word using the Markov chain, with a maximum
// number of tokens to generate before returning.
//
// Useful if the chain has a chance of entering infinite generation, or to simply
// prevent an overly long word.
func (generator *wordGenerator) LimitedGenerate(maxTokens int) (output string, err error) {
	var builder strings.Builder

	bytes2d, err := generator.generator.LimitedGenerate(maxTokens)

	if err != nil {
		return
	}

	for _, bytes := range bytes2d {
		for _, b := range bytes {
			builder.WriteByte(b)
		}
	}

	output = builder.String()

	return
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
			runesAsBytes := make([][]byte, len(runes), len(runes))

			for i, r := range runes {
				runeLen := utf8.RuneLen(r)
				buf := make([]byte, runeLen, runeLen)
				utf8.EncodeRune(buf, r)
				runesAsBytes[i] = buf
			}

			bytes[index] = runesAsBytes
		}(i, word)
	}

	waiter.Wait()

	g.generator, err = gen.New(bytes, prefixLen)
	generator = g

	return
}
