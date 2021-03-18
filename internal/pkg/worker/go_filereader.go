package worker

import (
	"bufio"
	"log"
	"os"
)

// StartFileReader  File Reader Thread
func StartFileReader(f *os.File) <-chan []byte {
	out := make(chan []byte, 10000)
	scanner := bufio.NewScanner(f)
	go func() {
		for scanner.Scan() {
			x := scanner.Bytes()
			log.Println(x)
			out <- x
		}
		close(out)
		log.Println("FINISHED FILEREADER")
	}()
	return out
}
