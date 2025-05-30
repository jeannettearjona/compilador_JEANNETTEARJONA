package ast

type TipoDato struct {
	Int    int
	Float  int
	Bool   int
	String int
	Void   int
	Param  int
}

type MemoryManager struct {
	Global TipoDato
	Local  TipoDato
	Temp   TipoDato
	Const  TipoDato
}

const (
	Gi = 1000
	Gf = 2000

	Li = 3000
	Lf = 4000

	Ti = 5000
	Tf = 6000
	Tb = 7000

	Ci = 8000
	Cf = 9000
	Cs = 10000

	VOID = 100
	par  = 50
)

const (
	SUMA      = 1
	RESTA     = 2
	MULT      = 3
	DIV       = 4
	GT        = 5
	LT        = 6
	EQ        = 7
	NEQ       = 8
	PRINT     = 9
	GOTOF     = 10
	GOTO      = 11
	ERA       = 12
	ENDFUNC   = 13
	END       = 14
	parametro = 15
	GOSUB     = 16
)

var CodigoNum_Operador = map[string]int{
	"+":         SUMA,
	"-":         RESTA,
	"*":         MULT,
	"/":         DIV,
	">":         GT,
	"<":         LT,
	"=":         EQ,
	"!=":        NEQ,
	"print":     PRINT,
	"GOTOF":     GOTOF,
	"GOTO":      GOTO,
	"ERA":       ERA,
	"ENDFUNC":   ENDFUNC,
	"END":       END,
	"parametro": parametro,
	"GOSUB":     GOSUB,
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		Global: TipoDato{
			Int:   Gi,
			Float: Gf,
			Void:  VOID,
		},
		Local: TipoDato{
			Int:   Li,
			Float: Lf,
			Param: par,
		},
		Temp: TipoDato{
			Int:   Ti,
			Float: Tf,
			Bool:  Tb,
		},
		Const: TipoDato{
			Int:    Ci,
			Float:  Cf,
			String: Cs,
		},
	}
}

func (mm *MemoryManager) GetGlobalVarMem(dataType string) int {
	switch dataType {
	case "int":
		dir_vir := mm.Global.Int
		mm.Global.Int++
		return dir_vir
	case "float":
		dir_vir := mm.Global.Float
		mm.Global.Float++
		return dir_vir
	case "void":
		dir_vir := mm.Global.Void
		mm.Global.Void++
		return dir_vir
	default:
		panic("tipo de dato no soportado en MEMORIA GLOBAL")
	}
}

func (mm *MemoryManager) GetLocalVarMem(dataType string) int {
	switch dataType {
	case "int":
		dir_vir := mm.Local.Int
		mm.Local.Int++
		return dir_vir
	case "float":
		dir_vir := mm.Local.Float
		mm.Local.Float++
		return dir_vir
	case "param":
		dir_vir := mm.Local.Param
		mm.Local.Param++
		return dir_vir
	default:
		panic("tipo de dato no soportado en MEMORIA LOCAL")
	}
}

func (mm *MemoryManager) GetTempVarMem(dataType string) int {
	switch dataType {
	case "int":
		dir_vir := mm.Temp.Int
		mm.Temp.Int++
		return dir_vir
	case "float":
		dir_vir := mm.Temp.Float
		mm.Temp.Float++
		return dir_vir
	case "bool":
		dir_vir := mm.Temp.Bool
		mm.Temp.Bool++
		return dir_vir
	default:
		panic("tipo de dato no soportado en MEMORIA TEMPORAL")
	}
}

func (mm *MemoryManager) GetConstVarMem(dataType string) int {
	switch dataType {
	case "int":
		dir_vir := mm.Const.Int
		mm.Const.Int++
		return dir_vir
	case "float":
		dir_vir := mm.Const.Float
		mm.Const.Float++
		return dir_vir
	case "string":
		dir_vir := mm.Const.String
		mm.Const.String++
		return dir_vir
	default:
		panic("tipo de dato no soportado en MEMORIA CTES")
	}
}

func (mm *MemoryManager) Get_TotalTempCount() int {
	return (mm.Temp.Int - Ti) + (mm.Temp.Float - Tf) + (mm.Temp.Bool - Tb)
}

func (mm *MemoryManager) ResetTemps() {
	mm.Temp.Int = Ti
	mm.Temp.Float = Tf
	mm.Temp.Bool = Tb
}
