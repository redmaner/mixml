package miuires

import (
	"bufio"
	"bytes"
	"io"
)

// NewElementScanner returns a pointer to a new bufio.Scanner with the ScanElements
// split function enabled
func NewElementScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Split(ScanElements)
	return s
}

// ScanElements is a split function for a bufio.Scanner that returns each slice of data that ends with
// a greater than sign '>'. This function is a custom function that enables to read XML encoded payloads.
func ScanElements(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '>'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
