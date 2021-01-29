// Package chain contains random generators that use markov chains
package chain

import "errors"

// StopIteration signifies that the generator should stop
var StopIteration error = errors.New("Generation has completed")
