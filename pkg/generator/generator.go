package generator

// To be implemented by types, specifically Markov chains, that generate a
// random string output.
type Generator interface {
	// Generate a random output using a Markov chain.
	Generate() string
}
