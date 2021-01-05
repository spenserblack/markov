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

[example input files]: ./examples/resources
