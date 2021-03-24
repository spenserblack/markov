package chain

import "testing"

// TestEmptyLine tries to generate a word from an input with an empty line to
// confirm that it will not cause an Out-of-Bounds error.
func TestEmptyLine(t *testing.T) {
	var feed []string = []string{""}

	chain, _ := NewWordChain(feed, 3)
	generator := chain.Generator()

	for _, err := generator(); ; _, err = generator() {
		if err == nil {
			continue
		}
		if err != ErrStopIter {
			t.Fatalf("Received an error besides ErrStopIter: %v", err)
		}
	}
}
