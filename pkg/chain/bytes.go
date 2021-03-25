package chain

import (
	"crypto/sha1"
	"errors"
	"math/rand"
	"sync"
)

type markovChain = map[string][][]byte

// BytesChain uses a Markov chain to create a randomized sequence of tokens.
type BytesChain struct {
	chain         markovChain
	chainStarters [][][]byte
}

// Generator returns a generator function that returns a random sequence of
// bytes each time it is called. The random bytes are chosen from the previous
// bytes returned. If not enough bytes have been returned yet, then random bytes
// are chosen that are marked as being able to start a chain.
//
// Each []byte is a token in the chain.
func (bytesChain *BytesChain) Generator() func() (next []byte, stop error) {
	lastBytes := bytesChain.chainStarters[rand.Intn(len(bytesChain.chainStarters))]

	h := sha1.New()

	return func() (next []byte, stop error) {
		defer h.Reset()
		if len(lastBytes) != 0 {
			next = lastBytes[0]
		}

		if next == nil {
			stop = ErrStopIter
			return
		}

		for _, bytes := range lastBytes {
			h.Write(bytes)
		}
		key := string(h.Sum(nil))

		nextValue := []byte(nil)
		if nextValues, ok := bytesChain.chain[key]; ok {
			nextValue = nextValues[rand.Intn(len(nextValues))]
		}

		for i, v := range lastBytes[1:] {
			lastBytes[i] = v
		}

		lastBytes[len(lastBytes)-1] = nextValue

		return
	}
}

// NewBytesChain feeds data to a markov chain and returns the generator.
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
func NewBytesChain(feed [][][]byte, prefixLen int) (generator *BytesChain, err error) {
	if prefixLen < 1 {
		err = errors.New("prefixLen must be 1 or greater")
		return
	}

	generator = new(BytesChain)
	var chain struct {
		sync.Mutex
		val markovChain
	}
	var chainStarters struct {
		sync.Mutex
		val [][][]byte
	}
	chain.val = make(markovChain)
	chainStarters.val = make([][][]byte, 0, len(feed))

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
			chainStarters.Lock()
			chainStarters.val = append(chainStarters.val, prefix)
			chainStarters.Unlock()

			for i, suffix := range sequence[adjustedPrefixLen:] {
				var prefix [][]byte = sequence[i : i+adjustedPrefixLen]

				for _, byteSlice := range prefix {
					h.Write(byteSlice)
				}

				key := string(h.Sum(nil))

				chain.Lock()
				chain.val[key] = append(chain.val[key], suffix)
				chain.Unlock()
				h.Reset()
			}

			var lastPrefix [][]byte = sequence[len(sequence)-adjustedPrefixLen:]
			for _, byteSlice := range lastPrefix {
				h.Write(byteSlice)
			}
			key := string(h.Sum(nil))

			chain.Lock()
			chain.val[key] = append(chain.val[key], nil)
			chain.Unlock()

		}(sequence)
	}

	waiter.Wait()
	generator.chain = chain.val
	generator.chainStarters = chainStarters.val
	return
}
