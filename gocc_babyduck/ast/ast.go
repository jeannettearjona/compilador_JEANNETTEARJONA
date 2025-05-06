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

func DeclaracionVar(variables []VariableInfo) (*HashMap, error) {
	//mapaVar := NewHashMap()

	for _, variable := range variables {
		if GlobalVarTable.Contains(variable.Name) {
			return nil, fmt.Errorf("variable '%s' ya declarada", variable.Name)
		}
		GlobalVarTable.Add(variable.Name, variable)
	}
	return GlobalVarTable, nil
}

func DeclaracionVarLocal(variables []VariableInfo) (*HashMap, error) {
	mapaVarLocal := NewHashMap()

	for _, variable := range variables {
		if mapaVarLocal.Contains(variable.Name) {
			return nil, fmt.Errorf("variable '%s' ya declarada en la funcion '%s'", variable.Name, CurrentFunction.Name)
		}
		mapaVarLocal.Add(variable.Name, variable)
	}
	return mapaVarLocal, nil
}

/*func DeclararFuncion(name string, params []VariableInfo, localVars *HashMap) (*HashMap, error) { //NO PASAR LOCAL VARS COMO HASHMAP

	//Registrar parametros en la tabla local (porque parametros son variables locales)
	for _, param := range params {
		if localVars.Contains(param.Name) {
			return nil, fmt.Errorf("error: parámetro '%s' duplicado en la función", param.Name)
		}
		localVars.Add(param.Name, param)
	}

	//Registrar funcion en la tabla de funciones
	if FunctionDirectory.Contains(name) {
		return nil, fmt.Errorf("error: función '%s' ya declarada", name)
	}

	CurrentFunction = &FunctionInfo{
		Name:       name,
		Parameters: params,
		VarTable:   localVars,
	}
	FunctionDirectory.Add(name, *CurrentFunction)

	// Ya no la necesitas después
	CurrentFunction = nil

	return localVars, nil
}*/
func DeclararFuncion(name string, params []VariableInfo, localVars *HashMap) (*HashMap, error) {
	// Verificar si la función ya está declarada
	if FunctionDirectory.Contains(name) {
		return nil, fmt.Errorf("la función '%s' ya está declarada", name)
	}

	// Inicializar la función actual
	CurrentFunction = &FunctionInfo{
		Name:       name,
		Parameters: params,
		VarTable:   localVars,
	}

	// Registrar la función en el directorio de funciones
	FunctionDirectory.Add(name, *CurrentFunction)

	// Devolver el mapa de variables locales
	return nil, nil
}

/*func DeclaracionVarOLD(ids []string, tipo string) (*HashMap, error) {
	mapaVar := NewHashMap()

	for _, id := range ids {
		if mapaVar.Contains(id) {
			return nil, fmt.Errorf("variable '%s' ya declarada", id)
		}
		mapaVar.Add(id, VariableInfo{Name: id, Type: tipo})
	}
	return mapaVar, nil
}*/

func FinalizarFuncion() {
	CurrentFunction = nil
}

func ResetSemanticState() {
	GlobalVarTable = NewHashMap()
	FunctionDirectory = NewHashMap()
}
