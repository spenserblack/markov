.PHONY: clean

markov:
	go build -o markov cmd/markov/main.go

clean:
	rm markov