package word

import "testing"

// TestMaximumHit creates a word generator that *would* generate too many tokens,
// checks that generation stops when the set maximum number of tokens is reached instead.
func TestMaximumHit(t *testing.T) {
	generator, err := New([]string{"abcdefg"}, 1)

	if err != nil {
		t.Fatalf(`Failed to create generator: %v`, err)
	}

	output, err := generator.LimitedGenerate(3)
	want := "abc"
	if output != want || err != nil {
		t.Fatalf(`LimitedGenerate(3) = %q, %v, want %q, nil`, output, err, want)
	}
}

// TestTooSmallMax calls word.LimitedGenerate with a maximum that is less than
// the generator's prefix length, checking for an error.
func TestTooSmallMax(t *testing.T) {
	generator, err := New([]string{"abcdefg"}, 3)

	if err != nil {
		t.Fatalf(`Failed to create generator: %v`, err)
	}

	output, err := generator.LimitedGenerate(2)

	if len(output) != 0 || err == nil {
		t.Fatalf(`LimitedGenerate(2) = %q, %v, want "", error`, output, err)
	}
}
