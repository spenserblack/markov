package sentence

import (
	bytegenerator "github.com/spenserblack/markov/pkg/generator"
	"strings"
	"sync"
)

// Markov is a Markov chain container for creating a sentence.
type Markov struct {
	generator *bytegenerator.ByteGenerator
}

// Generate returns a random sentence using the Markov chain.
func (generator *Markov) Generate() string {
	var builder strings.Builder

	for _, bytes := range generator.generator.Generate() {
		for _, b := range bytes {
			builder.WriteByte(b)
		}
		builder.WriteRune(' ')
	}

	return builder.String()
}

// LimitedGenerate return a random sentence using the Markov chain, with a maximum
// number of tokens to generate before returning.
//
// Useful if the chain has a chance of entering infinite generation, or to simply
// prevent an overly long sentence.
func (generator *Markov) LimitedGenerate(maxTokens int) (output string, err error) {
	var builder strings.Builder

	bytes2d, err := generator.generator.LimitedGenerate(maxTokens)

	if err != nil {
		return
	}

	for _, bytes := range bytes2d {
		for _, b := range bytes {
			builder.WriteByte(b)
		}
		builder.WriteRune(' ')
	}

	output = builder.String()

	return
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
func New(sentences []string, prefixLen int) (generator *Markov, err error) {
	generator = new(Markov)

	bytes := make([][][]byte, len(sentences), len(sentences))
	var waiter sync.WaitGroup

	for i, sentence := range sentences {
		waiter.Add(1)
		go func(index int, sentence string) {
			defer waiter.Done()

			words := strings.Split(sentence, " ")

			for _, word := range words {
				bytes[index] = append(bytes[index], []byte(word))
			}

		}(i, sentence)
	}

	waiter.Wait()

	generator.generator, err = bytegenerator.New(bytes, prefixLen)

	return
}
