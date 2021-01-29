package chain_test

import (
	"fmt"
	"github.com/spenserblack/markov/pkg/chain"
)

// ExampleNewBytesChain feeds "Hello, World!" and "Hello, Go!" into a bytes
// chain.
func ExampleNewBytesChain() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	chain.NewBytesChain(feed, 1)
}

// ExampleBytesChain feeds "Hello, World!" and "Hello, Go!" into a bytes chain,
// and outputs the first word.
func ExampleBytesChain() {
	feed := [][][]byte{
		{[]byte("Hello,"), []byte("World!")},
		{[]byte("Hello,"), []byte("Go!")},
	}

	byteChain, err := chain.NewBytesChain(feed, 1)

	if err != nil {
		panic(err)
	}

	next := byteChain.Generator()

	bytes, err := next()

	if err == chain.StopIteration {
		panic("We should be able to generate at least 1 byte slice :(")
	}
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
	// Output: Hello,
}

// ExampleWordChain feeds "Hello!" and "Test" into a word chain,and outputs the
// second rune (which should always be "e").
func ExampleWordChain() {
	feed := []string{"Hello!", "Test"}

	// Look 3 tokens back to generate the next token.
	wordChain, err := chain.NewWordChain(feed, 3)

	next := wordChain.Generator()

	next()

	r, err := next()

	if err == chain.StopIteration {
		panic("We should be able to generate at least 2 runes :(")
	}
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r))
	// Output: e
}

// ExampleSentenceChain feeds "It's a test." and "A passing test." into a word
// chain,and outputs the third word (which should always be "test.").
func ExampleSentenceChain() {
	feed := []string{
		"It's a test.",
		"A passing test.",
	}

	// Look 2 tokens back to generate the next token.
	sentenceChain, err := chain.NewSentenceChain(feed, 2)

	next := sentenceChain.Generator()

	next()
	next()

	word, err := next()

	if err == chain.StopIteration {
		panic("We should be able to generate at least 3 words :(")
	}

	fmt.Println(string(word))
	// Output: test.
}
