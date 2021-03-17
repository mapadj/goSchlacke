package worker

import (
	"context"
	"errors"
	"log"
	"reflect"

	"github.com/mapadj/goSchlacke/internal/pkg/tables"
)

func Upserter(errorChannel chan error, input <-chan tables.ImportContainer, upsertMethod reflect.Value, ctx context.Context) error {

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
				break enough
			}

			value := reflect.ValueOf(val).Elem()
			field := value.FieldByName("ConvertedAndValidated").Elem().Interface()
			if !ok {
				return errors.New("field not found")
			}

			inputParams := []reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(field),
			}

			// Upsert
			returnValue := upsertMethod.Call(inputParams)

			// Get 2nd return Parameter (Error) and drop the first.
			rerr := returnValue[1].Interface()
			if rerr != nil {
				err := rerr.(error)
				log.Println("Upsert failed: ", err)
				return err
			}

		}
	}
	log.Println("FINISHED UPSERTER")
	return nil
}
