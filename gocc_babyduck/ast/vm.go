package ast

import (
	"fmt"
	"strconv"
)

type VirtualMachine struct {
	Memory     map[int]interface{}
	Cuadruplos *Queue
	IP         int
	FunDir     *HashMap
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

	for _, key := range cte_table.Keys() { // itera sobre las claves 2, 1,5, "hola"
		cte_data, _ := cte_table.Get(key)

		cte_value := cte_data.(VariableInfo).Name //valor de la constante
		cte_type := cte_data.(VariableInfo).Type  //tipo de la constante
		dir := cte_data.(VariableInfo).Address    //direccion de memoria

		val := Cast_Value(cte_value, cte_type) //convierte el valor a su tipo correspondiente
		memory[dir] = val                      // a memory["5"] se le asigna la dirección de memoria
	}
	return memory
}

// crea la VM
func NewVirtualMachine(quads *Queue, cte_table *HashMap, fun_dir *HashMap) *VirtualMachine {

	memory := CopyCtes(cte_table) //carga ctes

	vm := &VirtualMachine{
		Memory:     memory,
		Cuadruplos: quads,
		IP:         0,
		FunDir:     fun_dir, //??
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
		panic("No se puede convertir a float64")
	}
}

func (vm *VirtualMachine) Run() {

	for vm.IP < vm.Cuadruplos.Size() {
		quad := vm.Cuadruplos.GetItem(vm.IP)

		izq := vm.Memory[quad.Izq] //obtengo lo que hay en direccion izq 2
		der := vm.Memory[quad.Der] //obtengo lo que hay en direccion der 1.5

		switch quad.Operador {
		case 1: //+
			if isFloat(izq) || isFloat(der) {
				vm.Memory[quad.Res] = toFloat(izq) + toFloat(der)
			} else {
				vm.Memory[quad.Res] = izq.(int) + der.(int)
			}
		case 2: //-
			if isFloat(izq) || isFloat(der) {
				vm.Memory[quad.Res] = toFloat(izq) - toFloat(der)
			} else {
				vm.Memory[quad.Res] = izq.(int) - der.(int)
			}
		case 3: //*
			if isFloat(izq) || isFloat(der) {
				vm.Memory[quad.Res] = toFloat(izq) * toFloat(der)
			} else {
				vm.Memory[quad.Res] = izq.(int) * der.(int)
			}
		case 4: // /
			if izq.(int) == 0 {
				panic("division por cero")
			} else if isFloat(izq) || isFloat(der) {
				vm.Memory[quad.Res] = toFloat(izq) / toFloat(der)
			} else {
				vm.Memory[quad.Res] = izq.(int) / der.(int)
			}
		case 5: // >
			if isFloat(izq) || isFloat(der) {
				vm.Memory[quad.Res] = toFloat(izq) > toFloat(der)
			} else {
				vm.Memory[quad.Res] = izq.(int) > der.(int)
			}
		case 6: // <
			if isFloat(izq) || isFloat(der) {
				vm.Memory[quad.Res] = toFloat(izq) < toFloat(der)
			} else {
				vm.Memory[quad.Res] = izq.(int) < der.(int)
			}
		case 7: // =
			vm.Memory[quad.Res] = izq
		case 8: // !=
			if isFloat(izq) || isFloat(der) {
				vm.Memory[quad.Res] = toFloat(izq) != toFloat(der)
			} else {
				vm.Memory[quad.Res] = izq.(int) != der.(int)
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
		default:
			panic("Operador no reconocido")
		}
		vm.IP++
	}
}
