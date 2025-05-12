package ast

import "fmt"

type VariableInfo struct {
	Name string
	Type string
}

type FunctionInfo struct {
	Name       string
	Parameters []VariableInfo
	VarTable   *HashMap
}

var GlobalVarTable = NewHashMap()    //tabla de variables globales ya esta disponible desde el inicio
var FunctionDirectory = NewHashMap() //tabla de funciones ya esta disponible desde el inicio
//var CurrentFunction *FunctionInfo
var CurrentFunction *FunctionInfo = nil

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

func InsertarVariableLocal(variables []VariableInfo) error {
	if CurrentFunction == nil {
		return fmt.Errorf("no hay función actual activa para insertar variable local")
	}

	for _, variable := range variables {
		if CurrentFunction.VarTable.Contains(variable.Name) {
			return fmt.Errorf("variable '%s' ya declarada en la funcion '%s'", variable.Name, CurrentFunction.Name)
		}
		CurrentFunction.VarTable.Add(variable.Name, variable)
	}
	return nil
}

func DeclararFuncion(name string, params []VariableInfo, localVars *HashMap) (*FunctionInfo, error) {
	// Verificar si la función ya está declarada
	if FunctionDirectory.Contains(name) {
		return nil, fmt.Errorf("la función '%s' ya está declarada", name)
	}

	// Inicializar la función actual
	/*CurrentFunction = &FunctionInfo{
		Name:       name,
		Parameters: params,
		VarTable:   localVars,
	}*/

	funcion := &FunctionInfo{
		Name:       name,
		Parameters: params,
		VarTable:   localVars,
	}

	// Registrar la función en el directorio de funciones
	FunctionDirectory.Add(name, *funcion)

	CurrentFunction = funcion

	// Devolver el mapa de variables locales
	return funcion, nil
}

func FinalizarFuncion() {
	CurrentFunction = nil
}

func ResetSemanticState() {
	GlobalVarTable = NewHashMap()
	FunctionDirectory = NewHashMap()
}

//Avance 3
var Operadores Stack
var Operandos Stack
var Tipos Stack

var Cuadruplos Queue

var PJumps StackInt

func BuscarVariable(name string) (VariableInfo, error) {

	// Buscar en la tabla de variables locales de la función actual

	/*if CurrentFunction == nil {
		return VariableInfo{}, fmt.Errorf("no hay funcion actual")
	}*/

	if CurrentFunction != nil {
		if CurrentFunction.VarTable.Contains(name) {
			varInfo, _ := CurrentFunction.VarTable.Get(name)
			return varInfo.(VariableInfo), nil
		}
	}

	// Buscar en la tabla de variables globales
	if GlobalVarTable.Contains(name) {
		varInfo, _ := GlobalVarTable.Get(name)
		return varInfo.(VariableInfo), nil
	}

	return VariableInfo{}, fmt.Errorf("variable '%s' no encontrada", name)
}

var TempCounter int = 0

func GetTemp() string {
	// Generar un nuevo nombre temporal
	temp := fmt.Sprintf("t%d", TempCounter)
	TempCounter++
	return temp
}

func GenerateQuadrupleForExp() error {

	if Operadores.IsEmpty() || Operandos.IsEmpty() {
		return fmt.Errorf("error: no hay operadores u operandos para generar el cuadruplo")
	}

	right_operand := Operandos.Pop()
	right_type := Tipos.Pop()
	left_operand := Operandos.Pop()
	left_type := Tipos.Pop()

	op := Operadores.Pop()

	resultType, err := DefaultSemanticCube.GetResultType(left_type, right_type, op)

	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	temp := GetTemp()

	// Generar el cuadruplo
	quad := NewQuadruple(op, left_operand, right_operand, temp)
	Cuadruplos.Enqueue(quad)

	Operandos.Push(temp)
	Tipos.Push(resultType)

	return nil
}

func GenerateQuadrupleForAssign(variable VariableInfo) error {
	if Operandos.IsEmpty() || Tipos.IsEmpty() {
		return fmt.Errorf("error: no hay operandos o tipos para generar el cuadruplo")
	}

	//lo que se asigna
	right := Operandos.Pop()
	right_type := Tipos.Pop()

	//la variable a la que se asigna
	left_type := variable.Type

	_, err := DefaultSemanticCube.GetResultType(left_type, right_type, "=")

	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	// Generar el cuadruplo
	quad := NewQuadruple("=", right, "", variable.Name)
	Cuadruplos.Enqueue(quad)

	return nil
}

func VerificarCondicion() error {
	tipoVerificacion := Tipos.Pop()
	if tipoVerificacion != "bool" {
		return fmt.Errorf("error: la condicion no es de tipo booleano")
	}
	return nil
}

func GoToF_IF() error {
	condicion := Operandos.Pop()
	tipo := Tipos.Pop()

	//si NO es una comparativa de int-int o float-float
	if tipo != "bool" {
		return fmt.Errorf("error: la condicion no es de tipo booleano")
	}

	//si ES una comparativa de int-int o float-float
	// Generar el cuadruplo
	quad := NewQuadruple("GoToF", condicion, "", "")
	Cuadruplos.Enqueue(quad)

	//pending jump
	PJumps.Push(Cuadruplos.Size() - 1)

	return nil
}

func GoToEndIf() {
	// Generar el cuadruplo
	quad := NewQuadruple("GoTo", "", "", "")
	Cuadruplos.Enqueue(quad)

	//pending jump
	PJumps.Push(Cuadruplos.Size() - 1)
}

/*func BackPatch() error{
	// Verificar si hay saltos pendientes

	if PJumps.IsEmpty() {
		return fmt.Errorf("error: no hay saltos pendientes para hacer backpatch")
	}

	// Obtener el índice del salto pendiente
	jumpIndex := PJumps.Pop()

	// Obtener el cuadruplo en la posición del salto pendiente
	quad := Cuadruplos.GetItem(jumpIndex)

	// Modificar el cuadruplo para que apunte a la posición correcta
	quad.Res = fmt.Sprintf("%d", Cuadruplos.Size())

	// Reinsertar el cuadruplo modificado
	Cuadruplos.Update(jumpIndex, quad)

	return nil
}*/

func ImprimirCuadruplos() {
	// Imprimir los cuadruplos generados
	fmt.Println("Cuadruplos generados:")
	for _, quad := range Cuadruplos.Print() {
		fmt.Printf("%s %s %s %s\n", quad.Operador, quad.Izq, quad.Der, quad.Res)
	}
}
