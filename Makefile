.PHONY: clean test

markov:
	go build -o markov cmd/markov/main.go

clean:
	rm markov

test:
	go test ./...
