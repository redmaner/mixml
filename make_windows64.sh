#!/bin/bash

# Set environemnt variables
export GOOS="windows"
export GOHOSTARCH="amd64"

# Setup directories
mkdir -p $PWD/out

# Go build
go build -o ./out/mixml_windows64.exe ./src/cmd/mixml 
