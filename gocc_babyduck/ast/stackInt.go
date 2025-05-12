package ast

type StackInt struct {
	items []int
}

// Crear un nuevo stack
/*func NewStackInt() *StackInt {
	return &StackInt{
		items: []int{},
	}
}*/

// CHECK IF EMPTY
func (s *StackInt) IsEmpty() bool {
	return len(s.items) == 0
}

// Agregar un elemento a la pila
func (s *StackInt) Push(value int) {
	s.items = append(s.items, value)
}

// Eliminar y devolver el elemento superior de la pila
func (s *StackInt) Pop() int {
	if s.IsEmpty() {
		panic("error: la pila está vacía")
	}
	top := s.Top()
	s.items = s.items[:len(s.items)-1]
	return top
}

// Devolver el elemento superior sin eliminarlo
func (s *StackInt) Top() int {
	if len(s.items) == 0 {
		panic("error: la pila está vacía")
	}
	return s.items[len(s.items)-1]
}

// Obtener el tamaño de la pila
func (s *StackInt) Size() int {
	return len(s.items)
}
