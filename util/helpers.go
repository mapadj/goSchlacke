package util

import (
	"bytes"
	"io"
)

// Check checks fatal error
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// LineCounter counts lines
func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
