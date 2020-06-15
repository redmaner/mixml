#!/bin/bash

# Set environemnt variables
export GOOS="linux"
export GOHOSTARCH="amd64"

# Setup directories
mkdir -p $PWD/out

# Go build
go build -o ./out/mixml_linux64 ./src/cmd/mixml 
