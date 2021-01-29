// Use a Markov chain to generate randomized sentences.
package main

import (
	"flag"
	"fmt"
	"github.com/spenserblack/markov/pkg/chain"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

var prefixLen = flag.Int("prefix", 1, "number of tokens to use to map following token")
var maxTokens = flag.Int("max", -1, "maximum number of tokens to generate. A negative number signifies no maximum")
var genWord = flag.Bool("word", false, "generate a word instead of a sentence")
var printHelp = flag.Bool("help", false, "print this help message")

func main() {
	var err error
	flag.Parse()

	feedFile := flag.Arg(0)

	if *printHelp || feedFile == "" {
		println("usage: markov [OPTIONS] [FILENAME]")
		println("FILENAME: a text file containing newline-separated groups of tokens to feed into the markov chain")
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
		markov, err := chain.NewWordChain(strings.Split(feed, "\n"), *prefixLen)
		if err != nil {
			panic(err)
		}
		generator := markov.Generator()

		for next, err := generator(); ; next, err = generator() {
			if err == nil {
				fmt.Print(string(next))
				continue
			}
			if err != chain.ErrStopIter {
				panic(err)
			}
			break
		}
	} else {
		markov, err := chain.NewSentenceChain(strings.Split(feed, "\n"), *prefixLen)
		if err != nil {
			panic(err)
		}
		generator := markov.Generator()

		for next, err := generator(); err == nil; next, err = generator() {
			fmt.Print(next)
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

func init() {
	rand.Seed(time.Now().Unix())
}
