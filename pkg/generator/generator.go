package generator

import (
	"crypto/sha1"
	"errors"
	"sync"
)

// Generator is a type that uses a Markov chain to generate a randomized string.
type Generator interface {
	// Generate returns a random output using a Markov chain.
	Generate() string
	// LimitedGenerate returns a random output using a markov chain,
	// with a set maximum number of tokens.
	LimitedGenerate(maxTokens int) (output string, err error)
}

type markovChain = map[string][][]byte

// MarkovGenerator uses a Markov to create a randomized sequence of tokens.
type MarkovGenerator struct {
	mutex         sync.Mutex
	chain         markovChain
	chainStarters [][]byte
	prefixLen     int
}

// New feeds data to a markov chain and returns the generator.
//
// The 3-Dimensional slice of bytes can be a bit confusing, but here's the
// logic behind it:
//
// - The 1st dimension is the list of inputs. For example, a list of sentences.
//
// - The 2nd dimension is the split inputs. For example, the words making up a
// sentence.
//
// - The 3rd dimension are the particles that each token is composed of. For
// example, the letters in a word.
//
// So, if you want to feed the sentences "Hello, World!" and "Hello, Go!" into
// the chain, you would use [][][]byte{
//	{[]byte("Hello,"), []byte("World!")},
//	{[]byte("Hello,"), []byte("Go!")},
// }
func New(feed [][][]byte, prefixLen int) (generator *MarkovGenerator, err error) {
	if prefixLen < 1 {
		err = errors.New("prefixLen must be 1 or greater")
		return
	}

	generator = new(MarkovGenerator)
	generator.chain = make(markovChain)
	generator.prefixLen = prefixLen

	var waiter sync.WaitGroup

	for _, sequence := range feed {
		// Let waiter know that goroutine will start
		waiter.Add(1)

		go func(sequence [][]byte) {
			// Let waiter know that goroutine has finished
			defer waiter.Done()

			var adjustedPrefixLen int
			if prefixLen >= len(sequence) {
				adjustedPrefixLen = len(sequence) - 1
			} else {
				adjustedPrefixLen = prefixLen
			}

			for i, suffix := range sequence[adjustedPrefixLen:] {
				var prefix [][]byte = sequence[i : i+adjustedPrefixLen]
				h := sha1.New()

				var flattenedByteSlice []byte

				for _, byteSlice := range prefix {
					for _, b := range byteSlice {
						flattenedByteSlice = append(flattenedByteSlice, b)
					}
				}

				h.Write(flattenedByteSlice)

				key := string(h.Sum(nil))

				generator.mutex.Lock()
				if i == 0 {
					generator.chainStarters = append(generator.chainStarters, flattenedByteSlice)
				}

				generator.chain[key] = append(generator.chain[key], suffix)
				generator.mutex.Unlock()
			}

			var lastPrefix [][]byte = sequence[len(sequence)-adjustedPrefixLen:]
			h := sha1.New()
			for _, byteSlice := range lastPrefix {
				h.Write(byteSlice)
			}
			key := string(h.Sum(nil))

			generator.mutex.Lock()
			generator.chain[key] = append(generator.chain[key], nil)
			generator.mutex.Unlock()

		}(sequence)
	}

	waiter.Wait()
	return
}
