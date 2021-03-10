package db

import (
	"bufio"
	"errors"
	"log"
	"os"

	fl "github.com/ianlopshire/go-fixedwidth"
)

/*

This function takes a file pointer and returns a string channel, which delivers the the data lines as string

*/

// FileReader  File Reader Thread
func FileReader(f *os.File) <-chan []byte {
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

// StructPress splits every line into structs and sends result Data Converter
func StructPress(in <-chan []byte, importChoserParams ImportChoserParams) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		for n := range in {
			d := importChoserParams.Functions.ImportTableFactory()
			err := fl.Unmarshal(n, &d)
			if err != nil {
				log.Fatal(err)
			}
			out <- d
		}
		close(out)
	}()
	return out
}

// StartConvertAndValidateThread Convert and Validate
func StartConvertAndValidateThread(in <-chan interface{}, importChoserParams ImportChoserParams, errorChannel chan error) <-chan interface{} {

	out := make(chan interface{})
	go func() {

		// Consume input channel
		for n := range in {

			// convert and validate
			c, err := importChoserParams.Functions.Convert()
			if err != nil {
				importChoserParams.ImportTxResult.NumberOfFailes++
				//TODO: Deactivate log in production.
				// log the error, if someone wants to know, what's wrong
				log.Printf("ErrorCount: %d Conversion Error: %v+ -> %v+ \n", importChoserParams.ImportTxResult.NumberOfFailes, n, err)

				// Check fail rate
				if importChoserParams.ImportTxResult.NumberOfFailes > importChoserParams.ImportTxParams.MaxErrorCount {
					// close channel and return
					errorChannel <- errors.New("MAX ERROR REACHED")
					close(errorChannel)
					break
				}
				// Skip to next round.
				continue
			}
			// Send converted and validated object to next instance
			out <- c
		}
		close(out)

	}()
	return out
}
