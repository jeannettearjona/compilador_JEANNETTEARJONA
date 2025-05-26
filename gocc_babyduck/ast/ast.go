package ast

import "fmt"

type VariableInfo struct {
	Name    string
	Type    string
	Address int
}

type FunctionInfo struct {
	Name       string
	Parameters []VariableInfo
	VarTable   *HashMap
}

var GlobalVarTable = NewHashMap()
var FunctionDirectory = NewHashMap()
var ConstantsVarTable = NewHashMap()

//var CurrentFunction *FunctionInfo
var CurrentFunction *FunctionInfo = nil
var Prog_MemoryManager = NewMemoryManager()

func Create_VarList(id_list []string, type_vars string) []VariableInfo {
	variables := make([]VariableInfo, len(id_list))
	for i, id := range id_list {
		variables[i] = VariableInfo{
			Name: id,
			Type: type_vars,
		}
	}
	return variables
}

func Declare_GlobalVars(variables []VariableInfo) (*HashMap, error) {

	for _, variable := range variables {
		if GlobalVarTable.Contains(variable.Name) {
			return nil, fmt.Errorf("variable '%s' ya declarada en GlobalVarTable", variable.Name)
		}

		//obtener la dirección de memoria
		direccion := Prog_MemoryManager.GetGlobalVarMem(variable.Type)
		variable.Address = direccion

		GlobalVarTable.Add(variable.Name, variable)
	}
	return GlobalVarTable, nil
}

func Declare_Function(fun_name string, params []VariableInfo, localVars *HashMap) (*FunctionInfo, error) {
	// Verificar si la función ya está declarada
	if FunctionDirectory.Contains(fun_name) {
		return nil, fmt.Errorf("funcion '%s' ya está declarada en FunctionDirectory", fun_name)
	}

	funcion := &FunctionInfo{
		Name:       fun_name,
		Parameters: params,
		VarTable:   localVars,
	}

	FunctionDirectory.Add(fun_name, *funcion)

	CurrentFunction = funcion

	return funcion, nil
}

func Declare_LocalVars(variables []VariableInfo) error {

	if CurrentFunction == nil {
		return fmt.Errorf("no hay función actual activa para insertar variable local")
	}

	for _, variable := range variables {
		if CurrentFunction.VarTable.Contains(variable.Name) {
			return fmt.Errorf("variable '%s' ya declarada en la funcion '%s'", variable.Name, CurrentFunction.Name)
		}

		//obtener la dirección de memoria
		direccion := Prog_MemoryManager.GetLocalVarMem(variable.Type)
		variable.Address = direccion

		CurrentFunction.VarTable.Add(variable.Name, variable)
	}
	return nil
}

func Declare_Constant(value, tipo string) int {
	if !ConstantsVarTable.Contains(value) {
		direccion := Prog_MemoryManager.GetConstVarMem(tipo)
		ConstantsVarTable.Add(value, direccion)
		return direccion
	}
	direccion, _ := ConstantsVarTable.Get(value)
	return direccion.(int)
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
var Operandos StackInt
var Tipos Stack

var Cuadruplos Queue

var PJumps StackInt

func BuscarVariable(var_name string) (VariableInfo, error) {
	//buscar en locales
	if CurrentFunction != nil {
		if CurrentFunction.VarTable.Contains(var_name) {
			varInfo, _ := CurrentFunction.VarTable.Get(var_name)
			return varInfo.(VariableInfo), nil
		}
	}

	//buscar en globales
	if GlobalVarTable.Contains(var_name) {
		varInfo, _ := GlobalVarTable.Get(var_name)
		return varInfo.(VariableInfo), nil
	}

	return VariableInfo{}, fmt.Errorf("variable '%s' no encontrada", var_name)
}

func IsStackEmpty() error {
	if Operadores.IsEmpty() || Operandos.IsEmpty() {
		return fmt.Errorf("stacks de operadores-operandos están vacíos")
	}

	return nil
}

func GenerateQuadrupleForExp() error {

	//pop de lado izq y der de la operacion
	right_operand := Operandos.Pop()
	right_type := Tipos.Pop()
	left_operand := Operandos.Pop()
	left_type := Tipos.Pop()

	op := Operadores.Pop()

	//verificar si se vale la operacion entre los dos tipos
	resultType, err := DefaultSemanticCube.GetResultType(left_type, right_type, op)

	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	temp := Prog_MemoryManager.GetTempVarMem(resultType)
	op_code := CodigoNum_Operador[op]

	// Generar el cuadruplo
	quad := NewQuadruple(op_code, left_operand, right_operand, temp)
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

	op_code := CodigoNum_Operador["="]

	// Generar el cuadruplo
	quad := NewQuadruple(op_code, right, 0, variable.Address)
	Cuadruplos.Enqueue(quad)

	return nil
}

func GenerateQuad_GOTOF() error {
	condicion := Operandos.Pop()
	tipo_condicion := Tipos.Pop()

	//la cond solo puede ser bool
	if tipo_condicion != "bool" {
		return fmt.Errorf("la condicion no es de tipo booleano")
	}

	op_code := CodigoNum_Operador["GOTOF"]
	quad := NewQuadruple(op_code, condicion, 0, 0)
	Cuadruplos.Enqueue(quad)

	//indice del GOTOF
	PJumps.Push(Cuadruplos.Size() - 1)

	return nil
}

func GenerateQuad_GOTO() error {
	op_code := CodigoNum_Operador["GOTO"]
	quad := NewQuadruple(op_code, 0, 0, 0)
	Cuadruplos.Enqueue(quad)

	//indice del GOTO
	PJumps.Push(Cuadruplos.Size() - 1)
	return nil
}

func Fill_QuadJumps(in_blank, jump_to int) error {
	quadToFill := Cuadruplos.GetItem(in_blank)
	quadToFill.Res = jump_to
	Cuadruplos.Update(in_blank, quadToFill)
	return nil
}

func GenerateQuadrupleForPrint(valor int) error {

	op_code := CodigoNum_Operador["print"]
	quad := NewQuadruple(op_code, valor, 0, 0)
	Cuadruplos.Enqueue(quad)
	return nil
}

func ImprimirCuadruplos() {

	fmt.Println("Cuadruplos generados:")
	for i, quad := range Cuadruplos.Print() {
		fmt.Printf("%d: %d %d %d %d\n", i, quad.Operador, quad.Izq, quad.Der, quad.Res)
	}

}

func VerificarCondicion(tipo string) error {

	if tipo != "bool" {
		return fmt.Errorf("error: la condicion no es de tipo booleano")
	}
	return nil
}
