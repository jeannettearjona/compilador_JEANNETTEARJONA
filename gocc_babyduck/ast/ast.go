package ast

import "fmt"

type VariableInfo struct {
	Name string
	Type string
	//Scope string
}

type FunctionInfo struct {
	Name       string
	Parameters []VariableInfo
	VarTable   *HashMap
}

var GlobalVarTable = NewHashMap()
var FunctionDirectory = NewHashMap()
var CurrentFunction *FunctionInfo

//var funcionActual string = ""

//tabla puede ser local o global, dependiendo de donde se declare la variable
func DeclaracionVar(ids []string, tipo string) (interface{}, error) {
	var targetTable *HashMap

	if CurrentFunction == nil {
		targetTable = GlobalVarTable
	} else {
		targetTable = CurrentFunction.VarTable
	}

	for _, id := range ids {
		if targetTable.Contains(id) {
			return nil, fmt.Errorf("variable '%s' ya declarada", id)
		}
		targetTable.Add(id, VariableInfo{Name: id, Type: tipo})
	}
	return targetTable, nil
}

func ValidateParams(params []VariableInfo) error {
	paramSet := NewHashMap()
	for _, param := range params {
		if paramSet.Contains(param.Name) {
			return fmt.Errorf("error: parámetro '%s' duplicado en la función", param.Name)
		}
		paramSet.Add(param.Name, param)
	}
	return nil
}

func ProcessFuncDecl(name string, params []VariableInfo, localVars *HashMap) (interface{}, error) {
	if FunctionDirectory.Contains(name) {
		return nil, fmt.Errorf("error: función '%s' ya declarada", name)
	}

	if err := ValidateParams(params); err != nil {
		return nil, err
	}

	// Agregar los parámetros a la tabla de variables locales
	for _, param := range params {
		if localVars.Contains(param.Name) {
			return nil, fmt.Errorf("variable '%s' ya declarada en función '%s'", param.Name, name)
		}
		localVars.Add(param.Name, param)
	}

	CurrentFunction = &FunctionInfo{
		Name:       name,
		Parameters: params,
		VarTable:   localVars,
	}

	FunctionDirectory.Add(name, *CurrentFunction)

	// Ya no la necesitas después
	CurrentFunction = nil

	return nil, nil
}

func DeclaracionVarFromTable(ids []string, tipo string, table *HashMap) (interface{}, error) {
	for _, id := range ids {
		if table.Contains(id) {
			return nil, fmt.Errorf("variable '%s' ya declarada", id)
		}
		table.Add(id, VariableInfo{Name: id, Type: tipo})
	}
	return table, nil
}

/*func IniciarFuncion(nombre string, parametros []VariableInfo, localVars *HashMap) error {
	// Verificar si la función ya está declarada
	if FunctionDirectory.Contains(nombre) {
		return fmt.Errorf("función '%s' ya declarada", nombre)
	}

	if err := ValidateParams(parametros); err != nil {
		return err
	}

	CurrentFunction = &FunctionInfo{
		Name:       nombre,
		Parameters: parametros,
		VarTable:   NewHashMap(), // Nueva tabla para variables locales
	}
	FunctionDirectory.Add(nombre, CurrentFunction)
	return nil
}*/

func FinalizarFuncion() {
	CurrentFunction = nil
}

func ResetSemanticState() {
	GlobalVarTable = NewHashMap()
	FunctionDirectory = NewHashMap()
}
