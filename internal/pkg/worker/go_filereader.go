package worker

import (
	"bufio"
	"log"
	"os"
)

// StartFileReader  File Reader Thread
func StartFileReader(f *os.File) <-chan []byte {
	out := make(chan []byte, 1000)
	scanner := bufio.NewScanner(f)
	go func() {
		for scanner.Scan() {
			out <- scanner.Bytes()
		}
		close(out)
		log.Println("FINISHED FILEREADER")
	}()
	return out
}
