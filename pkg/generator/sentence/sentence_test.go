package sentence

import (
	"strings"
	"testing"
)

// TestMaximumHit creates a sentence generator that *would* generate too many tokens,
// checks that generation stops when the set maximum number of tokens is reached instead.
func TestMaximumHit(t *testing.T) {
	generator, err := New([]string{"a b c d e f g"}, 1)

	if err != nil {
		t.Fatalf(`Failed to create generator: %v`, err)
	}

	output, err := generator.LimitedGenerate(3)
	wantedLength := len(strings.Split(output, " "))
	if wantedLength != 3 || err != nil {
		t.Fatalf(`LimitedGenerate(3) = %q, %v, want <string of length 3>, nil`, output, err)
	}
}

// TestTooSmallMax calls sentence.LimitedGenerate with a maximum that is less than
// the generator's prefix length, checking for an error.
func TestTooSmallMax(t *testing.T) {
	generator, err := New([]string{"a b c d e f g"}, 3)

	if err != nil {
		t.Fatalf(`Failed to create generator: %v`, err)
	}

	output, err := generator.LimitedGenerate(2)

	if len(output) != 0 || err == nil {
		t.Fatalf(`LimitedGenerate(2) = %q, %v, want "", error`, output, err)
	}
}
