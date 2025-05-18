package ast

import (
	"fmt"
)

const (
	// Globales
	GlobalIntStart   = 1000
	GlobalFloatStart = 2000

	// Locales
	LocalIntStart   = 3000
	LocalFloatStart = 4000

	// Temporales
	TempIntStart   = 5000
	TempFloatStart = 6000
	TempBoolStart  = 7000

	// Constantes
	ConstIntStart    = 8000
	ConstFloatStart  = 9000
	ConstStringStart = 10000
)

type MemoryManager struct {
	globalInt   int
	globalFloat int
	localInt    int
	localFloat  int
	tempInt     int
	tempFloat   int
	tempBool    int
	constInt    int
	constFloat  int
	constString int
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		globalInt:   GlobalIntStart,
		globalFloat: GlobalFloatStart,
		localInt:    LocalIntStart,
		localFloat:  LocalFloatStart,
		tempInt:     TempIntStart,
		tempFloat:   TempFloatStart,
		tempBool:    TempBoolStart,
		constInt:    ConstIntStart,
		constFloat:  ConstFloatStart,
		constString: ConstStringStart,
	}
}

func (mm *MemoryManager) GetGlobalVarMem(dataType string) int {
	switch dataType {
	case "int":
		addr := mm.globalInt
		mm.globalInt++
		return addr
	case "float":
		addr := mm.globalFloat
		mm.globalFloat++
		return addr
	default:
		panic("error: tipo de dato no soportado en la memoria global")
	}
}

func (mm *MemoryManager) GetLocalVarMem(dataType string) int {
	switch dataType {
	case "int":
		addr := mm.localInt
		mm.localInt++
		return addr
	case "float":
		addr := mm.localFloat
		mm.localFloat++
		return addr
	default:
		panic("error: tipo de dato no soportado en la memoria local")
	}
}

func (mm *MemoryManager) GetTempVarMem(dataType string) int {
	switch dataType {
	case "int":
		addr := mm.tempInt
		mm.tempInt++
		return addr
	case "float":
		addr := mm.tempFloat
		mm.tempFloat++
		return addr
	case "bool":
		addr := mm.tempBool
		mm.tempBool++
		return addr
	default:
		//panic("error: tipo de dato no soportado en la memoria temporal")
		panic(fmt.Sprintf("error: tipo de dato no soportado en la memoria temporal porque es %s", dataType))
	}
}

func (mm *MemoryManager) GetConstVarMem(dataType string) int {
	switch dataType {
	case "int":
		addr := mm.constInt
		mm.constInt++
		return addr
	case "float":
		addr := mm.constFloat
		mm.constFloat++
		return addr
	default:
		panic("error: tipo de dato no soportado en la memoria de constantes")
	}
}

func (mm *MemoryManager) GetStringConstMem() int {
	addr := mm.constString
	mm.constString++
	return addr
}
