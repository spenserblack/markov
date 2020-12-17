#!/bin/bash
LEN=${1:-"1"}
go run markov.go -n $LEN "Hello World Hello Go Goodbye World Goodbye Go Goodbye Markov Hello Markov Go Markov Go Go Go"
