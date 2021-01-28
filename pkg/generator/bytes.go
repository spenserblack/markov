package generator

import (
	"crypto/sha1"
	"errors"
	"math/rand"
	"sync"
)

type markovChain = map[string][][]byte

// ByteGenerator uses a Markov chain to create a randomized sequence of tokens.
type ByteGenerator struct {
	chain         markovChain
	chainStarters [][][]byte
}

// Generate creates a randomized sequence of bytes using the Markov chain.
//
// Each []byte is a token in the chain.
//
// For example, if Generate was used to create a random sentence, then each
// []byte would a word in the sentence.
func (generator *ByteGenerator) Generate() func() []byte {
	lastBytes := generator.chainStarters[rand.Intn(len(generator.chainStarters))]

	h := sha1.New()

	return func() []byte {
		defer h.Reset()
		next := lastBytes[0]

		for _, bytes := range lastBytes {
			h.Write(bytes)
		}
		key := string(h.Sum(nil))

		nextValue := []byte(nil)
		if nextValues, ok := generator.chain[key]; ok {
			nextValue = nextValues[rand.Intn(len(nextValues))]
		}

		for i, v := range lastBytes[1:] {
			lastBytes[i] = v
		}

		lastBytes[len(lastBytes)-1] = nextValue

		return next
	}
}

// NewByteGenerator feeds data to a markov chain and returns the generator.
//
// The 3-Dimensional slice of bytes can be a bit confusing, but here's the
// logic behind it:
//
// - The 3rd dimension is the list of inputs. For example, a list of sentences.
//
// - The 2nd dimension is the split inputs. For example, the words making up a
// sentence.
//
// - The 1st dimension are the particles that each token is composed of. For
// example, the letters in a word.
func NewByteGenerator(feed [][][]byte, prefixLen int) (generator *ByteGenerator, err error) {
	if prefixLen < 1 {
		err = errors.New("prefixLen must be 1 or greater")
		return
	}

	generator = new(ByteGenerator)
	generator.chain = make(markovChain)
	generator.chainStarters = make([][][]byte, 0, len(feed))
	var chainMutex, chainStarterMutex sync.Mutex

	var waiter sync.WaitGroup

	for _, sequence := range feed {
		// Let waiter know that goroutine will start
		waiter.Add(1)

		go func(sequence [][]byte) {
			// Let waiter know that goroutine has finished
			defer waiter.Done()
			h := sha1.New()

			var adjustedPrefixLen int
			if prefixLen > len(sequence) {
				adjustedPrefixLen = len(sequence)
			} else {
				adjustedPrefixLen = prefixLen
			}

			var prefix [][]byte = sequence[:adjustedPrefixLen]
			chainStarterMutex.Lock()
			generator.chainStarters = append(generator.chainStarters, prefix)
			chainStarterMutex.Unlock()

			for i, suffix := range sequence[adjustedPrefixLen:] {
				var prefix [][]byte = sequence[i : i+adjustedPrefixLen]

				for _, byteSlice := range prefix {
					h.Write(byteSlice)
				}

				key := string(h.Sum(nil))

				chainMutex.Lock()
				generator.chain[key] = append(generator.chain[key], suffix)
				chainMutex.Unlock()
				h.Reset()
			}

			var lastPrefix [][]byte = sequence[len(sequence)-adjustedPrefixLen:]
			for _, byteSlice := range lastPrefix {
				h.Write(byteSlice)
			}
			key := string(h.Sum(nil))

			chainMutex.Lock()
			generator.chain[key] = append(generator.chain[key], nil)
			chainMutex.Unlock()

		}(sequence)
	}

	waiter.Wait()
	return
}
