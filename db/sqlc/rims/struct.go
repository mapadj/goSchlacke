package rims

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
type ImportTable struct {
	Code     string `fixed:"1,5"`
	Width    string `fixed:"6,10"`
	Height   string `fixed:"11,11"`
	OnePiece string `fixed:"12,12"`
	Diameter string `fixed:"13,14"`
	Material string `fixed:"15,15"`
}

func NewImportTable() (ret interface{}) {
	return ImportTable{}
}
