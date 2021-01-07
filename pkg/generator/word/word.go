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
func (generator *Markov) Generate() string {
	var builder strings.Builder
	starter := generator.chainStarters[rand.Intn(len(generator.chainStarters))]
	lastRunes := []rune(starter)
	lastRunesLen := len(lastRunes)
	builder.WriteString(starter)

	for {
		key := string(lastRunes)
		nextValues, nextValuesExist := generator.chain[key]

		if !nextValuesExist {
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
func New(words []string, prefixLen int) (generator *Markov) {
	generator = new(Markov)
	generator.chain = make(map[string][]rune)
	generator.prefixLen = prefixLen
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

				generator.mutex.Lock()
				if i == 0 {
					generator.chainStarters = append(generator.chainStarters, prefix)
				}

				if suffixes, ok := generator.chain[prefix]; ok {
					generator.chain[prefix] = append(suffixes, suffix)
				} else {
					generator.chain[prefix] = []rune{suffix}
				}
				generator.mutex.Unlock()
			}
		}(word)
	}

	waiter.Wait()
	return
}
