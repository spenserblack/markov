package generator

import (
	"crypto/sha1"
	"errors"
	"math/rand"
	"sync"
)

type markovChain = map[string][][]byte

// ByteGenerator uses a Markov to create a randomized sequence of tokens.
type ByteGenerator struct {
	mutex         sync.Mutex
	chain         markovChain
	chainStarters [][]byte
	prefixLen     int
}

func (generator *ByteGenerator) Generate() (output [][]byte) {
	starter := generator.chainStarters[rand.Intn(len(generator.chainStarters))]
	var lastBytes [][]byte = [][]byte{starter}
	lastBytesLen := len(lastBytes)
	output = append(output, starter)
	h := sha1.New()

	for ; ; h.Reset() {
		for _, bytes := range lastBytes {
			h.Write(bytes)
		}
		key := string(h.Sum(nil))

		nextValues, nextValuesExist := generator.chain[key]

		if !nextValuesExist {
			return
		}

		var nextValue []byte = nextValues[rand.Intn(len(nextValues))]

		if nextValue == nil {
			return
		}

		for i := 0; i < lastBytesLen-1; i++ {
			lastBytes[i] = lastBytes[i+1]
		}
		lastBytes[lastBytesLen-1] = nextValue

		output = append(output, nextValue)
	}
}

func (generator *ByteGenerator) LimitedGenerate(maxTokens int) (output [][]byte, err error) {
	if maxTokens < generator.prefixLen {
		err = errors.New("maxTokens cannot be less than the number of tokens used in the prefix")
		return
	}

	starter := generator.chainStarters[rand.Intn(len(generator.chainStarters))]
	var lastBytes [][]byte = [][]byte{starter}
	lastBytesLen := len(lastBytes)
	output = append(output, starter)
	h := sha1.New()

	for i := generator.prefixLen; i < maxTokens; i++ {
		h.Reset()

		for _, bytes := range lastBytes {
			h.Write(bytes)
		}
		key := string(h.Sum(nil))

		nextValues, nextValuesExist := generator.chain[key]

		if !nextValuesExist {
			return
		}

		var nextValue []byte = nextValues[rand.Intn(len(nextValues))]

		if nextValue == nil {
			return
		}

		for i := 0; i < lastBytesLen-1; i++ {
			lastBytes[i] = lastBytes[i+1]
		}
		lastBytes[lastBytesLen-1] = nextValue

		output = append(output, nextValue)
	}
	return
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
func New(feed [][][]byte, prefixLen int) (generator *ByteGenerator, err error) {
	if prefixLen < 1 {
		err = errors.New("prefixLen must be 1 or greater")
		return
	}

	generator = new(ByteGenerator)
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
