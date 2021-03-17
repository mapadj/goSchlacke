package worker

import (
	"errors"
	"log"

	"github.com/mapadj/goSchlacke/internal/pkg/tables"
)

// StartConvertAndValidateThread Convert and Validate
func StartConvertAndValidateThread(in <-chan tables.ImportContainer, importChoserParams tables.ImportChoserParams, errorChannel chan error) <-chan tables.ImportContainer {

	out := make(chan tables.ImportContainer, 1000)
	go func() {
		i := 0
		// Consume input channel
		for x := range in {
			i++
			// convert and validate
			err := x.ConvertAndValidate()
			if err != nil {
				importChoserParams.ImportTxResult.NumberOfFailes++
				//TODO: Deactivate log in production.
				// log the error, if someone wants to know, what's wrong
				log.Printf("Line: %d \t ErrorCount: %d \t Conversion Error: %v+ -> %v+ \n", i, importChoserParams.ImportTxResult.NumberOfFailes, x, err)

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
			out <- x
		}
		close(out)
		log.Println("FINISHED CONVERTER")

	}()
	return out
}
