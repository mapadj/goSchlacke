package worker

import (
	"context"
	"log"
	"reflect"

	"github.com/mapadj/goSchlacke/internal/pkg/tables"
)

func Upserter(errorChannel chan error, input <-chan tables.Importable, upsertMethod reflect.Value, ctx context.Context) error {

enough:
	// This is the Query Loop. It waits for receives Data from channel X until it closes.
	for {
		select {
		case err := <-errorChannel:
			return err
		// Wait for data
		case val, ok := <-input:
			// check data health
			if !ok {
				log.Println(val, ok, "loop broke")
				break enough
			}

			inputParams := []reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(val),
			}

			// Upsert
			returnValue := upsertMethod.Call(inputParams)

			// Get 2nd return Parameter (Error) and drop the first.
			err := returnValue[1].Interface().(error)

			if err != nil {
				log.Println("Upsert failed: ", err)
				return err
			}
		}
	}
	return nil
}
