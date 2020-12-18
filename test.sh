#!/bin/bash
LEN=${1:-"1"}
go run markov.go -n $LEN "This is a sentence.
This is a test.
That is a string.
That is random."
