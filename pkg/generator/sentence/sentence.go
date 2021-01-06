package sentence

import (
	"math/rand"
	"strings"
	"sync"
)

// Markov chain container for creating a sentence.
type Markov struct {
	mutex     sync.Mutex
	chain     map[string][]string
	prefixLen int
}

// Generate a random sentence from the Markov chain.
func (markov *Markov) Generate() string {
	var builder strings.Builder
	chainStarters := markov.chain[""]
	starter := chainStarters[rand.Intn(len(chainStarters))]
	lastWords := strings.Split(starter, " ")
	lastWordsLen := len(lastWords)
	builder.WriteString(starter)

	for {
		key := strings.Join(lastWords, " ")
		nextValues := markov.chain[key]

		if len(nextValues) == 0 {
			return builder.String()
		}

		nextValue := nextValues[rand.Intn(len(nextValues))]

		for i := 0; i < lastWordsLen-1; i++ {
			lastWords[i] = lastWords[i+1]
		}
		lastWords[lastWordsLen-1] = nextValue

		builder.WriteRune(' ')
		builder.WriteString(nextValue)
	}
}

// Feed data to create a Markov chain.
// Each sentence in `sentences` should be a string of space-separated words to
// be used when building the chain -- order of the words determines how each next
// word in a generated sentence is decided.
// `prefixLen` is the number of words to be used as a "key" to deciding the next
// word. For example, if `prefixLen` is 2 and the generated text is "I made a
// chain" then "I made" was a key to "a" and "made a" was a key to "chain" in
// the sentence.
func New(sentences []string, prefixLen int) *Markov {
	markov := new(Markov)
	markov.chain = make(map[string][]string)
	markov.prefixLen = prefixLen
	var waiter sync.WaitGroup

	for _, words := range sentences {
		// Let waiter know that goroutine will start
		waiter.Add(1)

		go func(sentence string) {
			// Let waiter know that goroutine has finished
			defer waiter.Done()

			splitWords := strings.Split(sentence, " ")

			var adjustedPrefixLen int
			if splitWordsLen := len(splitWords); prefixLen >= splitWordsLen {
				adjustedPrefixLen = splitWordsLen - 1
			} else {
				adjustedPrefixLen = prefixLen
			}

			for i, suffix := range splitWords[adjustedPrefixLen:] {
				prefix := strings.Join(splitWords[i:i+adjustedPrefixLen], " ")

				markov.mutex.Lock()
				if i == 0 {
					markov.chain[""] = append(markov.chain[""], prefix)
				}

				if suffixes, ok := markov.chain[prefix]; ok {
					markov.chain[prefix] = append(suffixes, suffix)
				} else {
					markov.chain[prefix] = []string{suffix}
				}
				markov.mutex.Unlock()
			}
		}(words)
	}

	waiter.Wait()

	return markov
}
