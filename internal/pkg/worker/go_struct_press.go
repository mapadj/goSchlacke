package worker

import (
	"log"

	fl "github.com/ianlopshire/go-fixedwidth"
	"github.com/mapadj/goSchlacke/internal/pkg/tables"
)

// StructPress splits every line into structs and sends result Data Converter
func StructPress(in <-chan []byte, factory tables.ImportHandler) <-chan tables.Importable {

	out := make(chan tables.Importable)
	go func() {
		for n := range in {
			d := factory.NewContainer()
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
