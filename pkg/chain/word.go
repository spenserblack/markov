package chain

import (
	"errors"
	"sync"
	"unicode/utf8"
)

// WordChain represents a Markov chain that can be used to generate a sequence
// of runes.
type WordChain struct {
	chain *BytesChain
}

// Generator returns a generator of random runes using the Markov chain.
//
// Returns a StopIteration error if/when generation has completed.
func (chain *WordChain) Generator() func() (next rune, stop error) {
	g := chain.chain.Generator()

	return func() (next rune, stop error) {
		bytes, err := g()
		if err != nil {
			return next, err
		}
		next, _ = utf8.DecodeRune(bytes)

		if next == utf8.RuneError {
			stop = errors.New("Could not decode bytes to rune. Was valid UTF-8 used?")
		}
		return
	}
}

// NewWordChain feeds data to a markov chain and return the word generator.
//
// Each word in `words` should be a string of letters to be used when building the
// chain -- order of the letters determines how each next letter in a generated
// word is decided.
// `prefixLen` is the number of letters to be used as a "key" to deciding the next
// letter. For example, if `prefixLen` is 2 and the generated text is "abcd" then
// "ab" was a key to "c" and "bc" was a key to "d" in the word.
func NewWordChain(words []string, prefixLen int) (wordChain *WordChain, err error) {
	wordChain = new(WordChain)

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

	wordChain.chain, err = NewBytesChain(bytes, prefixLen)

	return
}
