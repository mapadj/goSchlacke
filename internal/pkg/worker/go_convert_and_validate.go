package worker

import (
	"errors"
	"log"

	"github.com/mapadj/goSchlacke/internal/pkg/tables"
)

// StartConvertAndValidateThread Convert and Validate
func StartConvertAndValidateThread(in <-chan tables.Importable, importChoserParams tables.ImportChoserParams, errorChannel chan error) <-chan tables.Importable {

	out := make(chan tables.Importable)
	go func() {

		// Consume input channel
		for x := range in {

			// convert and validate
			err := x.ConvertAndValidate()
			if err != nil {
				importChoserParams.ImportTxResult.NumberOfFailes++
				//TODO: Deactivate log in production.
				// log the error, if someone wants to know, what's wrong
				log.Printf("ErrorCount: %d Conversion Error: %v+ -> %v+ \n", importChoserParams.ImportTxResult.NumberOfFailes, x, err)

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

	}()
	return out
}
