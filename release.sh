#!/usr/bin/env bash

set -e
dir=$( pwd )

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/pancakeswap-v1.0.0 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/pancakeswap-v1.0.0.exe main.go