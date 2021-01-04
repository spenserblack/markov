// Use a Markov chain to generate randomized sentences.
package main

import (
	"flag"
	"fmt"
	"github.com/spenserblack/markov/pkg/generator"
	"github.com/spenserblack/markov/pkg/sentence"
	"github.com/spenserblack/markov/pkg/word"
	"math/rand"
	"strings"
	"time"
)

func main() {
	var markov generator.Generator
	prefixLen := flag.Int("n", 1, "length of words to use as a key")
	genWord := flag.Bool("w", false, "generate a word instead of a sentence")
	printHelp := flag.Bool("h", false, "print this help message")
	flag.Parse()

	feed := flag.Arg(0)

	if *printHelp || feed == "" {
		println("usage: markov [OPTIONS] [TEXT]")
		println("TEXT: the text to feed to the markov chain")
		flag.PrintDefaults()
		return
	}

	if *genWord {
		markov = word.New(strings.Split(feed, " "), *prefixLen)
	} else {
		markov = sentence.New(strings.Split(feed, "\n"), *prefixLen)
	}

	fmt.Println(markov.Generate())
}

func init() {
	rand.Seed(time.Now().Unix())
}
