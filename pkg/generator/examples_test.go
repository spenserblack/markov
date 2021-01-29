package generator_test

import (
	"fmt"
	"github.com/spenserblack/markov/pkg/generator"
)

// ExampleNewByteGenerator feeds "Hello, World!" and "Hello, Go!" into a bytes
// generator.
func ExampleNewByteGenerator() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	generator.NewByteGenerator(feed, 1)
}

// Example feeds "Hello, World!" and "Hello, Go!" into a bytes generator, and
// outputs the first word.
func Example() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	g, err := generator.NewByteGenerator(feed, 1)

	if err != nil {
		panic(err)
	}

	next := g.Generate()

	fmt.Println(string(next()))
	// Output: Hello,
}
