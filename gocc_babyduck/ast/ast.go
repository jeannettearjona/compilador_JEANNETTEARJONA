package ast

import "fmt"

type VariableInfo struct {
	Name    string
	Type    string
	Address int
}

type FunctionInfo struct {
	Name              string
	Parameters        []VariableInfo
	VarTable          *HashMap
	Address           int
	Counter_LocalVars int
	Counter_Params    int
	FunStart_Quad     int //donde inicia la funcion
	Counter_Temps     int //??
}

var GlobalVarTable = NewHashMap()
var ConstantsVarTable = NewHashMap()
var FunctionDirectory = NewHashMap()
var CurrentFunction *FunctionInfo = nil

var CurrentCalledFunction *FunctionInfo

var Prog_MemoryManager = NewMemoryManager()

var ParamCounter int
var MainADDRESS int

func Create_VarList(id_list []string, type_vars string) []VariableInfo {
	//lista de objetos VariableInfo de tamaño de la cantidad de ids
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

	//recorre lista de variables
	for _, variable := range variables {
		//checa si la variable ya existe
		if GlobalVarTable.Contains(variable.Name) {
			return nil, fmt.Errorf("variable '%s' ya declarada en GlobalVarTable", variable.Name)
		}

		//obtener la dirección de memoria de var global
		direccion := Prog_MemoryManager.GetGlobalVarMem(variable.Type)
		variable.Address = direccion

		GlobalVarTable.Add(variable.Name, variable)
	}
	return GlobalVarTable, nil
}

func Declare_Function(fun_name string, params []VariableInfo, localVars *HashMap) (*FunctionInfo, error) {
	//checar si la funcion ya existe
	if FunctionDirectory.Contains(fun_name) {
		return nil, fmt.Errorf("funcion '%s' ya esta declarada en FunctionDirectory", fun_name)
	}

	//obtener direccion de memoria para funcion
	direccion := Prog_MemoryManager.GetGlobalVarMem("void")

	//count de parametros
	num_params := len(params)

	//donde empieza la funcion (una linea despues de la declaracion)
	//start_quad := Cuadruplos.Size()

	funcion := &FunctionInfo{
		Name:              fun_name,
		Parameters:        params,
		VarTable:          localVars,
		Address:           direccion,
		Counter_LocalVars: 0,
		Counter_Params:    num_params,
		FunStart_Quad:     0,
		Counter_Temps:     0,
	}

	FunctionDirectory.Add(fun_name, funcion)
	CurrentFunction = funcion
	Prog_MemoryManager.ResetTemps() //?????

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

		if variable.Address == 0 {
			//obtener la dirección de memoria
			direccion := Prog_MemoryManager.GetLocalVarMem(variable.Type)
			variable.Address = direccion
		}

		CurrentFunction.VarTable.Add(variable.Name, variable)
	}
	//CurrentFunction.Counter_LocalVars = CurrentFunction.VarTable.Size() - CurrentFunction.Counter_Params
	return nil
}

func Declare_Constant(value, tipo string) int {
	if !ConstantsVarTable.Contains(value) {
		//obtener la direccion de memoria de cte
		direccion := Prog_MemoryManager.GetConstVarMem(tipo)
		variable := VariableInfo{
			Name:    value,
			Type:    tipo,
			Address: direccion,
		}
		//agregar la variable a la tabla de constantes
		ConstantsVarTable.Add(variable.Name, variable)
		return direccion
	}
	var_info, _ := ConstantsVarTable.Get(value)
	return var_info.(VariableInfo).Address
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

		//buscar en la tabla de variables locales
		if CurrentFunction.VarTable.Contains(var_name) {
			varInfo, _ := CurrentFunction.VarTable.Get(var_name)
			return varInfo.(VariableInfo), nil
		}

		//buscar en la tabla de parametros
		for _, param := range CurrentFunction.Parameters {
			if param.Name == var_name {
				return param, nil
			}
		}

	}

	//buscar en globales
	if GlobalVarTable.Contains(var_name) {
		varInfo, _ := GlobalVarTable.Get(var_name)
		return varInfo.(VariableInfo), nil
	}

	return VariableInfo{}, fmt.Errorf("variable '%s' no encontrada", var_name)
}

func BuscarFuncion(fun_name string) error {
	if !FunctionDirectory.Contains((fun_name)) {
		return fmt.Errorf("función '%s' no existe", fun_name)
	}
	return nil
}

func IsStackEmpty() error {
	if Operadores.IsEmpty() || Operandos.IsEmpty() {
		return fmt.Errorf("stacks de operadores-operandos están vacíos")
	}

	return nil
}

func GenerateQuadrupleForExp() (int, error) {

	//pop de lado izq y der de la operacion
	right_operand := Operandos.Pop()
	right_type := Tipos.Pop()
	left_operand := Operandos.Pop()
	left_type := Tipos.Pop()

	op := Operadores.Pop()

	//verificar si se vale la operacion entre los dos tipos
	resultType, err := DefaultSemanticCube.GetResultType(left_type, right_type, op)

	if err != nil {
		return 0, fmt.Errorf("error: %s", err)
	}

	temp := Prog_MemoryManager.GetTempVarMem(resultType)
	op_code := CodigoNum_Operador[op]

	// Generar el cuadruplo
	quad := NewQuadruple(op_code, left_operand, right_operand, temp)
	Cuadruplos.Enqueue(quad)

	Operandos.Push(temp)
	Tipos.Push(resultType)

	return temp, nil
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
	PJumps.Push(Cuadruplos.Size() - 1) //pushea el indice del quad GOTO
	return nil
}

func GenerateQuad_TOMAIN() error {
	op_code := CodigoNum_Operador["GOTO"]
	quad := NewQuadruple(op_code, 0, 0, 0)
	Cuadruplos.Enqueue(quad)
	MainADDRESS = Cuadruplos.Size() - 1 //guarda el indice del GOTO a main
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

func GenerateQuad_ERA(fun_name string) error {

	fun_info, _ := FunctionDirectory.Get(fun_name)
	fun_address := fun_info.(*FunctionInfo).Address

	op_code := CodigoNum_Operador["ERA"]
	quad := NewQuadruple(op_code, fun_address, 0, 0)
	Cuadruplos.Enqueue(quad)

	//indice del ERA
	PJumps.Push(Cuadruplos.Size() - 1)

	//
	funcion := fun_info.(*FunctionInfo)
	CurrentCalledFunction = funcion
	//CurrentFunction = funcion

	ParamCounter = 0

	return nil
}

func VerificarCondicion(tipo string) error {

	if tipo != "bool" {
		return fmt.Errorf("error: la condicion no es de tipo booleano")
	}
	return nil
}

func ValidarYGenerarParametros(expectedParams []VariableInfo) error {
	fmt.Printf("llego a ValidarYGenerarParametros con %d parametros esperados\n", len(expectedParams))
	//recibe los parametros esperados y los compara con los que se reciben
	if Operandos.Size() != len(expectedParams) || Tipos.Size() != len(expectedParams) {
		return fmt.Errorf("faltan argumentos para la función")
	}

	//stacks en orden
	var Params_EnOrden StackInt
	var Tipos_EnOrden Stack
	var ParamDirs_EnOrden StackInt

	//validacion de parametros
	for i := len(expectedParams) - 1; i >= 0; i-- {

		argument := Operandos.Pop()
		argType := Tipos.Pop()
		expectedType := expectedParams[i].Type

		//nuevo nuevo
		param_dir := expectedParams[i].Address
		ParamDirs_EnOrden.Push(param_dir)
		//------

		//checa que sean del mismo tipo
		if argType != expectedType {
			return fmt.Errorf("tipo incorrecto en parámetro %d: se esperaba %s, se recibió %s", i+1, expectedType, argType)
		}

		//si todo bien
		Params_EnOrden.Push(argument)
		Tipos_EnOrden.Push(argType)
	}

	//generacion de cuadruplos con los parametros en orden
	for !Params_EnOrden.IsEmpty() {
		//si si son
		param := Params_EnOrden.Pop()
		param_dir := ParamDirs_EnOrden.Pop()
		op := CodigoNum_Operador["parametro"]

		//par := Prog_MemoryManager.GetLocalVarMem("param")
		quad := NewQuadruple(op, param, 0, param_dir)
		Cuadruplos.Enqueue(quad)

		ParamCounter++
	}

	if !Operandos.IsEmpty() || !Tipos.IsEmpty() {
		return fmt.Errorf("demasiados argumentos en la llamada a la función")
	}

	return nil
}

func GenerateQuad_GOSUB() error {
	if CurrentCalledFunction == nil {
		return fmt.Errorf("no hay una función actual activa para generar GOSUB")
	}
	/*fun_address := CurrentFunction.Address
	start_fun := CurrentFunction.FunStart_Quad
	op_code := CodigoNum_Operador["GOSUB"]
	quad := NewQuadruple(op_code, fun_address, 0, start_fun)
	Cuadruplos.Enqueue(quad)

	CurrentFunction = nil
	ParamCounter = 0*/

	fun_address := CurrentCalledFunction.Address
	start_fun := CurrentCalledFunction.FunStart_Quad
	op_code := CodigoNum_Operador["GOSUB"]
	quad := NewQuadruple(op_code, fun_address, 0, start_fun)
	Cuadruplos.Enqueue(quad)

	// Ahora sí cambiamos el contexto
	CurrentFunction = CurrentCalledFunction
	CurrentCalledFunction = nil

	ParamCounter = 0

	return nil
}
