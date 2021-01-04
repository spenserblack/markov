My attempt at setting up a Markov chain without fully understanding how it works.

## Build Binary

```bash
go build -o markov cmd/markov/main.go
```

## Binary Usage

```bash
# print help
./markov -h

# feed a list of sentences and generate a random sentence
./markov <newline-separated sentences>

# feed a list of words and generate a random word
./markov -w <space-separated list of words>

# specify the number of previous tokens to use when mapping the next token
./markov -n <token number> <remaining args>
```
