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
	chainStarters [][][]byte
	prefixLen     int
}

// Generate creates a randomized sequence of bytes using the Markov chain.
//
// The bytes are in a [][]byte. Each []byte is a token in the chain.
//
// For example, if Generate was used to create a random sentence, then each
// []byte would a word in the sentence.
func (generator *ByteGenerator) Generate() (output [][]byte) {
	starter := generator.chainStarters[rand.Intn(len(generator.chainStarters))]

	output = starter

	h := sha1.New()

	for ; ; h.Reset() {
		var adjustedPrefixLen int

		if generator.prefixLen >= len(output) {
			adjustedPrefixLen = len(output)
		} else {
			adjustedPrefixLen = generator.prefixLen
		}

		for _, bytes := range output[len(output)-adjustedPrefixLen:] {
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

		output = append(output, nextValue)
	}
}

// LimitedGenerate returns a random sequence of bytes using the Markov chain,
// with a maximum number of []byte tokens to generate before returning.
//
// The bytes are in a [][]byte. Each []byte is a token in the chain.
//
// For example, if Generate was used to create a random sentence, then each
// []byte would a word in the sentence.
//
// Useful if the chain has a chance of entering infinite generation, or to simply
// prevent an overly long sequence of tokens.
func (generator *ByteGenerator) LimitedGenerate(maxTokens int) (output [][]byte, err error) {
	if maxTokens < generator.prefixLen {
		err = errors.New("maxTokens cannot be less than the number of tokens used in the prefix")
		return
	}

	starter := generator.chainStarters[rand.Intn(len(generator.chainStarters))]

	output = make([][]byte, maxTokens, maxTokens)
	outputIndex := copy(output, starter)

	// Shrink output to the final length
	defer func() {
		output = output[:outputIndex]
	}()

	h := sha1.New()

	for i := generator.prefixLen; i < maxTokens; i++ {
		h.Reset()

		var adjustedPrefixLen int

		if generator.prefixLen >= outputIndex {
			adjustedPrefixLen = outputIndex
		} else {
			adjustedPrefixLen = generator.prefixLen
		}

		for _, bytes := range output[outputIndex-adjustedPrefixLen:] {
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

		output[outputIndex] = nextValue
		outputIndex++
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
	generator.chainStarters = make([][][]byte, len(feed), len(feed))
	chainStarterCounter := 0
	generator.prefixLen = prefixLen

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
			generator.mutex.Lock()
			generator.chainStarters[chainStarterCounter] = prefix
			chainStarterCounter++
			generator.mutex.Unlock()

			for i, suffix := range sequence[adjustedPrefixLen:] {
				var prefix [][]byte = sequence[i : i+adjustedPrefixLen]

				for _, byteSlice := range prefix {
					h.Write(byteSlice)
				}

				key := string(h.Sum(nil))

				generator.mutex.Lock()
				generator.chain[key] = append(generator.chain[key], suffix)
				generator.mutex.Unlock()
				h.Reset()
			}

			var lastPrefix [][]byte = sequence[len(sequence)-adjustedPrefixLen:]
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
