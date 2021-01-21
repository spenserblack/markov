package word

import "testing"

// TestMaximumHit creates a word generator that *would* generate too many tokens,
// checks that generation stops when the set maximum number of tokens is reached instead.
func TestMaximumHit(t *testing.T) {
	generator, err := New([]string{"abcdefg"}, 1)

	if err != nil {
		t.Fatalf(`Failed to create generator: %v`, err)
	}

	output := generator.Generate(3)
	want := "abc"
	if output != want || err != nil {
		t.Fatalf(`Generate(3) = %q, %v, want %q, nil`, output, err, want)
	}
}
