// Use a Markov chain to generate randomized sentences.
package main

import (
	"flag"
	"fmt"
	"github.com/spenserblack/markov/pkg/generator"
	"github.com/spenserblack/markov/pkg/generator/sentence"
	"github.com/spenserblack/markov/pkg/generator/word"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func main() {
	var markov generator.Generator
	prefixLen := flag.Int("n", 1, "number of tokens to use to map following token")
	maxTokens := flag.Int("max", -1, "maximum number of tokens to generate. a negative number signifies no maximum")
	genWord := flag.Bool("w", false, "generate a word instead of a sentence")
	printHelp := flag.Bool("h", false, "print this help message")
	flag.Parse()

	feedFile := flag.Arg(0)

	if *printHelp || feedFile == "" {
		println("usage: markov [OPTIONS] [TEXT]")
		println("TEXT: the text to feed to the markov chain")
		flag.PrintDefaults()
		return
	}

	feedBytes, err := ioutil.ReadFile(feedFile)

	if err != nil {
		panic(err)
	}

	if feedBytes[len(feedBytes)-1] == '\n' {
		feedBytes = feedBytes[:len(feedBytes)-1]
	}

	feed := string(feedBytes)

	if *genWord {
		markov = word.New(strings.Split(feed, "\n"), *prefixLen)
	} else {
		markov = sentence.New(strings.Split(feed, "\n"), *prefixLen)
	}

	if *maxTokens < 0 {
		fmt.Println(markov.Generate())
	} else {
		panic("unimplemented")
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
