package ast

import (
	"fmt"
	"strconv"
)

func getDefault(typ string) interface{} {
	switch typ {
	case "int":
		return 0
	case "float":
		return 0.0
	case "bool":
		return false
	default:
		return nil
	}
}

type VirtualMachine struct {
	Memory        map[int]interface{}
	Cuadruplos    *Queue
	IP            int
	FunDir        *HashMap
	ActiveMemory  map[int]interface{}   //funcion activa
	LocalStack    []map[int]interface{} //local stack
	Callstack     []int                 //para guardar las direcciones de memoria de las funciones ???
	PendingMemory map[int]interface{}   //construccion de la memoria local pendiente
}

func Cast_Value(val_name, val_type string) interface{} {
	switch val_type {
	case "int":
		val, _ := strconv.Atoi(val_name)
		return val
	case "float":
		val, _ := strconv.ParseFloat(val_name, 64)
		return val
	case "string":
		return val_name
	default:
		panic(fmt.Sprintf("Tipo de dato '%s' no soportado", val_type))
	}
}

func CopyCtes(cte_table *HashMap) map[int]interface{} {

	memory := make(map[int]interface{})

	for _, key := range cte_table.Keys() {
		cte_data, _ := cte_table.Get(key)

		cte_value := cte_data.(VariableInfo).Name //valor de la constante
		cte_type := cte_data.(VariableInfo).Type  //tipo de la constante
		dir := cte_data.(VariableInfo).Address    //direccion de memoria

		val := Cast_Value(cte_value, cte_type) //convierte el valor a su tipo correspondiente
		memory[dir] = val                      // a memory["5"] se le asigna la dirección de memoria
	}
	return memory
}

func CopyFuncs(fun_dir *HashMap, memory map[int]interface{}) {
	for _, key := range fun_dir.Keys() {
		funcInfo, _ := fun_dir.Get(key)
		info := funcInfo.(*FunctionInfo) // O la estructura que uses

		memory[info.Address] = info.Name
	}
}

// crea la VM
func NewVirtualMachine(quads *Queue, cte_table *HashMap, fun_dir *HashMap) *VirtualMachine {

	//inicializa memoria global
	memory := CopyCtes(cte_table)
	CopyFuncs(fun_dir, memory)

	//crea instancia de vm
	vm := &VirtualMachine{
		Memory:     memory,
		Cuadruplos: quads,
		IP:         0,
		FunDir:     fun_dir,
	}

	return vm
}

func isFloat(val interface{}) bool {
	_, ok := val.(float64)
	return ok
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	default:
		panic("no se puede convertir a float")
	}
}

func (vm *VirtualMachine) GetValue(dir int) interface{} {
	if dir == 0 { //checa si es 0 para que no se rompa con GOTO o asi
		return 0
	} else if val, ok := vm.ActiveMemory[dir]; ok { //checa si el valor esta en local
		return val
	} else if val, ok := vm.Memory[dir]; ok { //checa si el valor esta en global
		return val
	}
	panic(fmt.Sprintf("no se encontro el valor en ninguna memoria para la direccion %d", dir))
}

func (vm *VirtualMachine) SetValue(dir int, val interface{}) {
	if _, exists := vm.Memory[dir]; exists { //si la direccion res YA existe en global, asignar el valor ahi
		vm.Memory[dir] = val
	} else if vm.ActiveMemory != nil { //si hay local, asignar a la direccion de res en local
		vm.ActiveMemory[dir] = val
	} else {
		vm.Memory[dir] = val //si no existe en ninguna, es nueva y se asigna a la marcha en global
		return
	}
}

func (vm *VirtualMachine) Run() {

	for vm.IP < vm.Cuadruplos.Size() {
		quad := vm.Cuadruplos.GetItem(vm.IP)

		izq := vm.GetValue(quad.Izq) //obtengo lo que hay en direccion izq
		der := vm.GetValue(quad.Der) //obtengo lo que hay en direccion der

		switch quad.Operador {
		case 1: //+
			if isFloat(izq) || isFloat(der) { //si uno es float hace cast del que no es float
				vm.SetValue(quad.Res, toFloat(izq)+toFloat(der))
			} else {
				vm.SetValue(quad.Res, izq.(int)+der.(int))
			}
		case 2: //-
			if isFloat(izq) || isFloat(der) {
				vm.SetValue(quad.Res, toFloat(izq)-toFloat(der))
			} else {
				vm.SetValue(quad.Res, izq.(int)-der.(int))
			}
		case 3: //*
			if isFloat(izq) || isFloat(der) {
				vm.SetValue(quad.Res, toFloat(izq)*toFloat(der))
			} else {
				vm.SetValue(quad.Res, izq.(int)*der.(int))
			}
		case 4: // /
			if toFloat(der) == 0.0 {
				panic("division por cero")
			} else if isFloat(izq) || isFloat(der) {
				vm.SetValue(quad.Res, toFloat(izq)/toFloat(der))
			} else {
				vm.SetValue(quad.Res, izq.(int)/der.(int))
			}
		case 5: // >
			if isFloat(izq) || isFloat(der) {
				vm.SetValue(quad.Res, toFloat(izq) > toFloat(der))
			} else {
				vm.SetValue(quad.Res, izq.(int) > der.(int))
			}
		case 6: // <
			if isFloat(izq) || isFloat(der) {
				vm.SetValue(quad.Res, toFloat(izq) < toFloat(der))
			} else {
				vm.SetValue(quad.Res, izq.(int) < der.(int))
			}
		case 7: // =
			vm.SetValue(quad.Res, izq)
		case 8: // !=
			if isFloat(izq) || isFloat(der) {
				vm.SetValue(quad.Res, toFloat(izq) != toFloat(der))
			} else {
				vm.SetValue(quad.Res, izq.(int) != der.(int))
			}
		case 9: //print
			fmt.Println(izq)
		case 10: //GOTOF
			condicion := izq.(bool)
			if !condicion {
				vm.IP = quad.Res
				continue
			}
		case 11: //GOTO
			vm.IP = quad.Res
			continue

		case 12: //ERA
			func_name := izq.(string)

			//obtener el objeto de functionInfo
			func_data, _ := vm.FunDir.Get(func_name)
			func_info := func_data.(*FunctionInfo)

			vm.PendingMemory = make(map[int]interface{})

			//inicializa params
			for _, param := range func_info.Parameters {
				vm.PendingMemory[param.Address] = getDefault(param.Type)
			}

			//inicializa variables locales
			for _, varData := range func_info.VarTable.data {
				varInfo := varData.(VariableInfo)
				vm.PendingMemory[varInfo.Address] = getDefault(varInfo.Type)
			}

			//for i := 0; i < func_info.Counter_Temps; i++ {
			// Puedes usar 5000+i si sabes que todos los temps empiezan ahí
			//	vm.PendingMemory[5000+i] = 0 // o nil si no te importa el tipo
			//}

		case 13: //ENDFUNC
			if len(vm.Callstack) == 0 {
				panic("No hay funciones en la pila de llamadas")
			}
			// DEBUG: Imprimir memoria local antes de restaurar
			fmt.Println("Memoria local antes de restaurar global:", vm.ActiveMemory)

			n := len(vm.LocalStack) - 1
			vm.ActiveMemory = vm.LocalStack[n] //restaura la memoria activa
			vm.LocalStack = vm.LocalStack[:n]  //elimina la memoria local de la funcion
			m := len(vm.Callstack) - 1
			vm.IP = vm.Callstack[m]         //restaura la direccion de retorno
			vm.Callstack = vm.Callstack[:m] //elimina la direccion de retorno
			fmt.Print("Regresa de la funcion a la direccion de memoria: ", vm.IP, "\n")
			continue
		case 14: //END
			return

		case 15: //parametro
			vm.PendingMemory[quad.Res] = vm.Memory[quad.Izq] //asigna el parametro a la memoria pendiente

		case 16: //GOSUB
			func_name := izq.(string)
			fmt.Println("llamada a la funcion:", func_name)

			//obtener el objeto de functionInfo
			func_data, _ := vm.FunDir.Get(func_name)
			func_info := func_data.(*FunctionInfo)

			//guarda dir de retorno (sig. quad despues de GOSUB)
			vm.Callstack = append(vm.Callstack, vm.IP+1)

			// NUEVO??? Guarda la memoria activa actual
			vm.LocalStack = append(vm.LocalStack, vm.ActiveMemory)
			fmt.Println("Memoria local guardada LO QUE HABIA EN ACTIVE MEMORY:", vm.ActiveMemory)

			//activa la memoria local preparada
			vm.ActiveMemory = vm.PendingMemory
			fmt.Println("Memoria activa:", vm.ActiveMemory)
			vm.PendingMemory = nil
			vm.IP = func_info.FunStart_Quad //cambia la direccion de memoria a la funcion
			fmt.Println("Llamada a la funcion, nueva direccion de memoria:", vm.IP)
			//fmt.Print("La direccion de la funcion que va a llamar es: ", func_info.Address, "\n")
			//return
			continue

		default:
			panic("Operador no reconocido")
		}
		vm.IP++
	}
}
