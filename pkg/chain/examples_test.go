package chain_test

import (
	"fmt"
	"github.com/spenserblack/markov/pkg/chain"
)

// ExampleNewByteGenerator feeds "Hello, World!" and "Hello, Go!" into a bytes
// chain.
func ExampleNewByteGenerator() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	chain.NewByteGenerator(feed, 1)
}

// Example feeds "Hello, World!" and "Hello, Go!" into a bytes chain, and
// outputs the first word.
func Example() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	g, err := chain.NewByteGenerator(feed, 1)

	if err != nil {
		panic(err)
	}

	next := g.Generate()

	fmt.Println(string(next()))
	// Output: Hello,
}
