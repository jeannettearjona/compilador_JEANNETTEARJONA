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

func Fill_VM_Memory(cte_table *HashMap) map[int]interface{} {
	memory := make(map[int]interface{})

	for _, valor_cte := range cte_table.Keys() { //itera sobre las claves de la tabla de constantes 5, 3.14, hola
		memory_dir, _ := cte_table.Get(valor_cte) //devuelve la dirección de memoria, dirAny = 8000
		dir := memory_dir.(int)

		var cte interface{} //cte: int, float, string

		//como las keys (valor_cte) son strings, se convierte al tipo correspondiente usando su direccion virtual
		if dir >= 8000 && dir < 9000 { //ctes int
			cte_int, _ := strconv.Atoi(valor_cte)
			cte = cte_int
		} else if dir >= 9000 && dir < 10000 { //ctes float
			cte_float, _ := strconv.ParseFloat(valor_cte, 64)
			cte = cte_float
		} else {
			cte = valor_cte //ctes string
		}

		memory[dir] = cte //a memoria[1000] se le asigna el valor de la cte
	}
	return memory
}

// Inicializa la VM con los cuádruplos y la memoria cargada
func NewVirtualMachine(quads *Queue, cte_table *HashMap, fun_dir *HashMap) *VirtualMachine {

	memory := Fill_VM_Memory(cte_table) // Carga la memoria con las constantes

	vm := &VirtualMachine{
		Memory:     memory,
		Cuadruplos: quads,
		IP:         0,
		FunDir:     fun_dir,
	}

	return vm
}

func (vm *VirtualMachine) Run() {

	for vm.IP < vm.Cuadruplos.Size() {
		quad := vm.Cuadruplos.GetItem(vm.IP)
		//quad = operador izq der res
		izq := vm.Memory[quad.Izq]
		der := vm.Memory[quad.Der]

		switch quad.Operador {
		case 1: //+
			vm.Memory[quad.Res] = izq.(int) + der.(int)
		case 2: //-
			vm.Memory[quad.Res] = izq.(int) - der.(int)
		case 3: //*
			vm.Memory[quad.Res] = izq.(int) * der.(int)
		case 4: // /
			if izq.(int) == 0 {
				panic("division por cero")
			}
			vm.Memory[quad.Res] = izq.(int) / der.(int)
		case 5: // >
			if izq.(int) > der.(int) {
				vm.Memory[quad.Res] = 1
			} else {
				vm.Memory[quad.Res] = 0
			}
		case 6: // >
			if izq.(int) < der.(int) {
				vm.Memory[quad.Res] = 1
			} else {
				vm.Memory[quad.Res] = 0
			}
		case 7: // =
			vm.Memory[quad.Res] = izq
		case 8: // !=
			if izq != der {
				vm.Memory[quad.Res] = 1
			} else {
				vm.Memory[quad.Res] = 0
			}
		case 9: //print
			fmt.Println(izq)
		default:
			panic("Operador no reconocido")
		}
		vm.IP++
	}
}
