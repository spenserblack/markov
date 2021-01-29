package chain_test

import (
	"fmt"
	"github.com/spenserblack/markov/pkg/chain"
)

// ExampleNewByteChain feeds "Hello, World!" and "Hello, Go!" into a bytes
// chain.
func ExampleNewByteChain() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	chain.NewByteChain(feed, 1)
}

// ExampleByteChain feeds "Hello, World!" and "Hello, Go!" into a bytes chain,
// and outputs the first word.
func ExampleByteChain() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	g, err := chain.NewByteChain(feed, 1)

	if err != nil {
		panic(err)
	}

	next := g.Generate()

	fmt.Println(string(next()))
	// Output: Hello,
}
