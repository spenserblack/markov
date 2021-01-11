package sentence

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
)

const chainEnder string = "\x00"

// Markov is a Markov chain container for creating a sentence.
type Markov struct {
	mutex     sync.Mutex
	chain     map[string][]string
	prefixLen int
}

// Generate returns a random sentence using the Markov chain.
func (generator *Markov) Generate() string {
	var builder strings.Builder
	chainStarters := generator.chain[""]
	starter := chainStarters[rand.Intn(len(chainStarters))]
	lastWords := strings.Split(starter, " ")
	lastWordsLen := len(lastWords)
	builder.WriteString(starter)

	for {
		key := strings.Join(lastWords, " ")
		nextValues, nextValuesExist := generator.chain[key]

		if !nextValuesExist {
			return builder.String()
		}

		nextValue := nextValues[rand.Intn(len(nextValues))]

		if nextValue == chainEnder {
			return builder.String()
		}

		for i := 0; i < lastWordsLen-1; i++ {
			lastWords[i] = lastWords[i+1]
		}
		lastWords[lastWordsLen-1] = nextValue

		builder.WriteRune(' ')
		builder.WriteString(nextValue)
	}
}

// LimitedGenerate return a random sentence using the Markov chain, with a maximum
// number of tokens to generate before returning.
//
// Useful if the chain has a chance of entering infinite generation, or to simply
// prevent an overly long sentence.
func (generator *Markov) LimitedGenerate(maxTokens int) (output string, err error) {
	if maxTokens < generator.prefixLen {
		err = errors.New("maxTokens cannot be less than the number of tokens used in the prefix")
		return
	}

	var builder strings.Builder
	chainStarters := generator.chain[""]
	starter := chainStarters[rand.Intn(len(chainStarters))]
	lastWords := strings.Split(starter, " ")
	lastWordsLen := len(lastWords)
	builder.WriteString(starter)

	for i := generator.prefixLen; i < maxTokens; i++ {
		key := strings.Join(lastWords, " ")
		nextValues, nextValuesExist := generator.chain[key]

		if !nextValuesExist {
			break
		}

		nextValue := nextValues[rand.Intn(len(nextValues))]

		if nextValue == chainEnder {
			break
		}

		for i := 0; i < lastWordsLen-1; i++ {
			lastWords[i] = lastWords[i+1]
		}
		lastWords[lastWordsLen-1] = nextValue

		builder.WriteRune(' ')
		builder.WriteString(nextValue)
	}

	output = builder.String()
	return
}

// New feeds data to a markov chain and returns the sentence generator.
//
// Each sentence in `sentences` should be a string of space-separated words to
// be used when building the chain -- order of the words determines how each next
// word in a generated sentence is decided.
// `prefixLen` is the number of words to be used as a "key" to deciding the next
// word. For example, if `prefixLen` is 2 and the generated text is "I made a
// chain" then "I made" was a key to "a" and "made a" was a key to "chain" in
// the sentence.
func New(sentences []string, prefixLen int) (generator *Markov) {
	generator = new(Markov)
	generator.chain = make(map[string][]string)
	generator.prefixLen = prefixLen
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

				generator.mutex.Lock()
				if i == 0 {
					generator.chain[""] = append(generator.chain[""], prefix)
				}

				if suffixes, ok := generator.chain[prefix]; ok {
					generator.chain[prefix] = append(suffixes, suffix)
				} else {
					generator.chain[prefix] = []string{suffix}
				}
				generator.mutex.Unlock()
			}
			lastPrefix := strings.Join(splitWords[len(splitWords)-adjustedPrefixLen:], " ")
			generator.mutex.Lock()
			generator.chain[lastPrefix] = append(generator.chain[lastPrefix], chainEnder)
			generator.mutex.Unlock()
		}(words)
	}

	waiter.Wait()
	return
}
