// Package chain contains random generators that use markov chains
package chain

import "errors"

// ErrStopIter signifies that the generator should stop
var ErrStopIter error = errors.New("Generation has completed")
