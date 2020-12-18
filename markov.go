package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Markov struct {
	chain         map[string][]string
	chainStarters []string
	prefixLen     int
}

func main() {
	prefixLen := flag.Int("n", 1, "length of words to use as a key")
	printHelp := flag.Bool("h", false, "print this help message")
	flag.Parse()

	feed := flag.Arg(0)

	if *printHelp || feed == "" {
		println("usage: markov [OPTIONS] [TEXT]")
		println("TEXT: the text to feed to the markov chain")
		flag.PrintDefaults()
		return
	}

	markov := NewSentence(feed, *prefixLen)
	fmt.Println(markov.Generate())
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (markov Markov) Generate() string {
	starter := markov.chainStarters[rand.Intn(len(markov.chainStarters))]
	output := starter

	for {
		splitWords := strings.Split(output, " ")
		lastWords := splitWords[len(splitWords)-markov.prefixLen:]
		key := strings.Join(lastWords, " ")
		nextValues := markov.chain[key]

		if len(nextValues) == 0 {
			return output
		}

		output = fmt.Sprintf("%v %v", output, nextValues[rand.Intn(len(nextValues))])
	}
}

func NewSentence(words string, prefixLen int) Markov {
	chain := make(map[string][]string)
	chainStarters := make([]string, 0)

	splitWords := strings.Split(words, " ")

	var lastPrefix string

	for i, suffix := range splitWords[prefixLen:] {
		prefix := strings.Join(splitWords[i:i+prefixLen], " ")
		lastPrefix = prefix

		if i == 0 {
			chainStarters = append(chainStarters, prefix)
		}

		if suffixes, ok := chain[prefix]; ok {
			chain[prefix] = append(suffixes, suffix)
		} else {
			chain[prefix] = []string{suffix}
		}
	}

	chain[lastPrefix] = append(chain[lastPrefix], "")
	chain[""] = make([]string, 0, 0)

	markov := Markov{chain, chainStarters, prefixLen}

	return markov
}
