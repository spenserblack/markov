package chain

import (
	"strings"
	"sync"
)

type SentenceChain struct {
	chain *ByteChain
}

// Generate returns a generator of random words that make up a sentence, using
// the Markov chain.
func (chain *SentenceChain) Generate() func() (next string, stop error) {
	g := chain.chain.Generate()

	return func() (next string, stop error) {
		if bytes := g(); bytes != nil {
			next = string(bytes)
		} else {
			stop = StopIteration
		}
		return
	}
}

// NewSentenceChain feeds data to a markov chain and returns the sentence generator.
//
// Each sentence in `sentences` should be a string of space-separated words to
// be used when building the chain -- order of the words determines how each next
// word in a generated sentence is decided.
// `prefixLen` is the number of words to be used as a "key" to deciding the next
// word. For example, if `prefixLen` is 2 and the generated text is "I made a
// chain" then "I made" was a key to "a" and "made a" was a key to "chain" in
// the sentence.
func NewSentenceChain(sentences []string, prefixLen int) (sentenceChain *SentenceChain, err error) {
	sentenceChain = new(SentenceChain)

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

	sentenceChain.chain, err = NewByteChain(bytes, prefixLen)

	return
}
