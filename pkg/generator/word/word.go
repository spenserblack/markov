package word

import (
	"math/rand"
	"strings"
	"sync"
)

// Markov chain container for creating a word.
type Markov struct {
	mutex         sync.Mutex
	chain         map[string][]rune
	chainStarters []string
	prefixLen     int
}

// Generate a random word from the Markov chain.
func (markov *Markov) Generate() string {
	var builder strings.Builder
	starter := markov.chainStarters[rand.Intn(len(markov.chainStarters))]
	lastRunes := []rune(starter)
	lastRunesLen := len(lastRunes)
	builder.WriteString(starter)

	for {
		key := string(lastRunes)
		nextValues := markov.chain[key]

		if len(nextValues) == 0 {
			return builder.String()
		}

		nextValue := nextValues[rand.Intn(len(nextValues))]

		for i := 0; i < lastRunesLen-1; i++ {
			lastRunes[i] = lastRunes[i+1]
		}
		lastRunes[lastRunesLen-1] = nextValue

		builder.WriteRune(nextValue)
	}
}

// Feed data to create a Markov chain.
// Each word in `words` should be a string of letters to be used when building the
// chain -- order of the letters determines how each next letter in a generated
// word is decided.
// `prefixLen` is the number of letters to be used as a "key" to deciding the next
// letter. For example, if `prefixLen` is 2 and the generated text is "abcd" then
// "ab" was a key to "c" and "bc" was a key to "d" in the word.
func New(words []string, prefixLen int) *Markov {
	markov := new(Markov)
	markov.chain = make(map[string][]rune)
	markov.prefixLen = prefixLen
	var waiter sync.WaitGroup

	for _, word := range words {
		// Let waiter know that goroutine will start
		waiter.Add(1)

		go func(word string) {
			// Let waiter know that goroutine has finished
			defer waiter.Done()

			var adjustedPrefixLen int
			if wordLen := len(word); prefixLen >= wordLen {
				adjustedPrefixLen = wordLen - 1
			} else {
				adjustedPrefixLen = prefixLen
			}

			for i, suffix := range word[adjustedPrefixLen:] {
				prefix := word[i : i+adjustedPrefixLen]

				markov.mutex.Lock()
				if i == 0 {
					markov.chainStarters = append(markov.chainStarters, prefix)
				}

				if suffixes, ok := markov.chain[prefix]; ok {
					markov.chain[prefix] = append(suffixes, suffix)
				} else {
					markov.chain[prefix] = []rune{suffix}
				}
				markov.mutex.Unlock()
			}
		}(word)
	}

	waiter.Wait()

	return markov
}
