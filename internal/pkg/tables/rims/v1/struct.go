package rims

import (
	db "github.com/mapadj/goSchlacke/internal/pkg/db/sqlc"
	"github.com/mapadj/goSchlacke/internal/pkg/tables"
)

/*
Feldname 	Typ 	Länge 	Beschreibung

code 		Num 	5 		Eindeutiger Code
width 		Num		5		Felgenmaulweite n Zoll
height		Char	1		Code für Felgenhöhe
one_piece	Bool	1		X = True, x = False
diameter	Num		2		Felgendurchmesser in Zoll
material	Char	1		S = Stahl, L = Leichtmetall
*/

// RimsV1LineLength Length of one Line
const RimsV1LineLength = 15

// ImportTable define the format
type importTable struct {
	tables.ImportStruct
	Code     string `fixed:"1,5"`
	Width    string `fixed:"6,10"`
	Height   string `fixed:"11,11"`
	OnePiece string `fixed:"12,12"`
	Diameter string `fixed:"13,14"`
	Material string `fixed:"15,15"`
}

// Container

type RimsV1Container struct {
	importTable           *importTable
	ConvertedAndValidated *db.UpsertRimsV1Params
}

func (c RimsV1Container) GetImportStruct() tables.ImportStruct {
	return c.importTable
}

func (c RimsV1Container) GetUpsertStruct() interface{} {
	return c.ConvertedAndValidated
}

// Handler
type RimsV1Handler struct {
}

func NewHandler() tables.ImportHandler {
	return &RimsV1Handler{}
}

func (handler RimsV1Handler) NewContainer() tables.ImportContainer {
	return &RimsV1Container{
		importTable:           &importTable{},
		ConvertedAndValidated: &db.UpsertRimsV1Params{},
	}
}
