![CI](https://github.com/spenserblack/markov/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/spenserblack/markov)](https://goreportcard.com/report/github.com/spenserblack/markov)

My attempt at setting up a Markov chain without fully understanding how it works.

## Install

```bash
go get -u github.com/spenserblack/markov/cmd/...
```

## Build Binary

```bash
make markov
```

## Binary Usage

The binary takes a filename which should point to a file containing tokens
that should be fed to the chain separated by newlines (see the [example input files]).

```bash
# print help
markov -h

# feed a list of sentences and generate a random sentence
markov <filename>

# feed a list of words and generate a random word
markov -w <filename>

# specify the number of previous tokens to use when mapping the next token
markov -n <token number> <remaining args>
```

## Use in your Go project

*This project uses `math/rand`. You may want to use `math/rand.Seed` to get more random results.*

*If you don't know until runtime if a random sentence or random word should be generated, the
`github.com/spenserblack/markov/pkg/generator.Generator` interface can be used.*

### Generate random sentence

In a sentence, each word is a token. For this project, a word is defined as a substring of a string
that has been split on spaces.

```go
package main

import (
	"fmt"
	"github.com/spenserblack/markov/pkg/generator/sentence"
)

func main() {
	sentences := []string{"foo bar baz", "foo bar bar", "foo foo baz"}
	fmt.Println(sentence.New(sentences, 1).Generate())
}
```

What this does is tell the generator to look 1 token back (so 1 word back) to decide the next word.

- the beginning of a sentence has a
  - 3/3 chance of being "foo"
- "foo" has a
  - 2/4 chance of being followed by "bar"
  - 1/4 chance of being followed by "foo"
  - 1/4 chance of being followed by "baz"
- "bar" has a
  - 1/3 chance of being followed by "bar"
  - 1/3 chance of ending a sentence
  - 1/3 chance of ending a sentence
- "baz" has a
  - 2/2 chance of ending a sentence

### Generate random word

In a word, each letter (or `rune`) is a token.

```go
package main

import (
	"fmt"
	"github.com/spenserblack/markov/pkg/generator/word"
)

func main() {
	words := []string{"Go", "Golang", "Good morning", "Good day"}
	fmt.Println(word.New(words, 1).Generate())
}
```

[example input files]: ./examples/resources
