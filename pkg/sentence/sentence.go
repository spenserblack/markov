package sentence

import (
	"math/rand"
	"strings"
	"sync"
)

// Markov chain container for creating a sentence.
type Markov struct {
	sync.Mutex
	chain     map[string][]string
	prefixLen int
}

// Generate a random sentence from the Markov chain.
func (markov *Markov) Generate() string {
	var builder strings.Builder
	chainStarters := markov.chain[""]
	starter := chainStarters[rand.Intn(len(chainStarters))]
	builder.WriteString(starter)

	for {
		splitWords := strings.Split(builder.String(), " ")
		lastWords := splitWords[len(splitWords)-markov.prefixLen:]
		key := strings.Join(lastWords, " ")
		nextValues := markov.chain[key]

		if len(nextValues) == 0 {
			return builder.String()
		}

		builder.WriteRune(' ')
		builder.WriteString(nextValues[rand.Intn(len(nextValues))])
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
	markov := Markov{}
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

			for i, suffix := range splitWords[prefixLen:] {
				prefix := strings.Join(splitWords[i:i+prefixLen], " ")

				markov.Lock()
				if i == 0 {
					markov.chain[""] = append(markov.chain[""], prefix)
				}

				if suffixes, ok := markov.chain[prefix]; ok {
					markov.chain[prefix] = append(suffixes, suffix)
				} else {
					markov.chain[prefix] = []string{suffix}
				}
				markov.Unlock()
			}
		}(words)
	}

	waiter.Wait()

	return &markov
}
