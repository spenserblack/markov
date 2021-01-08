package generator

// Generator is a type that uses a Markov chain to generate a randomized string.
type Generator interface {
	// Generate returns a random output using a Markov chain.
	Generate() string
}
