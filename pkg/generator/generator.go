package generator

// Generator is a type that uses a Markov chain to generate a randomized string.
type Generator interface {
	// Generate returns a random output using a Markov chain.
	Generate() string
	// LimitedGenerate returns a random output using a markov chain,
	// with a set maximum number of tokens.
	LimitedGenerate(maxTokens int) (output string, err error)
}
