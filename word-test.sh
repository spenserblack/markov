#!/bin/bash
LEN=${1:-"1"}
go run cmd/markov/main.go -n $LEN -w "Go Golang Goodbye GoodMorning"
