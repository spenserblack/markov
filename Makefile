.PHONY: clean format test vet

markov:
	go build -o markov cmd/markov/main.go

clean:
	rm markov

format:
	gofmt -s -w -l .

test:
	go test ./...

vet:
	go vet ./...
