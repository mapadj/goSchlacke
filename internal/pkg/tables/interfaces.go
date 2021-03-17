package tables

import (
	"context"
	"os"

	db "github.com/mapadj/goSchlacke/internal/pkg/db/sqlc"
)

type Importable interface {
	ConvertAndValidate() (err error)
}

type ImportHandler interface {
	NewContainer() Importable
}

// ImportChoserParams params
type ImportChoserParams struct {
	Ctx            context.Context
	ImportTxParams ImportTxParams
	ImportTxResult *ImportTxResult
	File           *os.File
}

// ImportTxParams params
type ImportTxParams struct {
	Queries              *db.Queries
	FilePath             string
	Table                string
	DatVersion           string
	MaxFailRateInPerCent int
	MaxErrorCount        int
}

// ImportTxResult params
type ImportTxResult struct {
	NumberOfLines  int
	NumberOfFailes int
	Success        bool
	Inserts        int
	Updates        int
}
