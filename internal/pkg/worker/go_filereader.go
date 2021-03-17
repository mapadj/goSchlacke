package worker

import (
	"bufio"
	"os"
)

// StartFileReader  File Reader Thread
func StartFileReader(f *os.File) <-chan []byte {
	out := make(chan []byte)
	scanner := bufio.NewScanner(f)
	go func() {
		for scanner.Scan() {
			out <- scanner.Bytes()
		}
		close(out)
	}()
	return out
}
