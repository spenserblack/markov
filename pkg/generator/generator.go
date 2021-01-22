// Package generator contains random generators that use markov chains
package generator

// StringGenerator is a type that uses a Markov chain to generate a randomized string.
type StringGenerator interface {
	// Generate returns a random output using a Markov chain, with an
	// optional maximum number of tokens.
	Generate(maxTokens int) string
}
