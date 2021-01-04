#!/bin/bash
LEN=${1:-"1"}
go run cmd/markov/main.go -n $LEN "This is a sentence.
This is a test.
That is a string.
That is random."
