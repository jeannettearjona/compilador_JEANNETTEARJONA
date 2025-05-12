package ast

type Stack struct {
	items []string
}

// CHECK IF EMPTY
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// ADD
func (s *Stack) Push(value string) {
	s.items = append(s.items, value)
}

// REMOVE
func (s *Stack) Pop() string {
	if s.IsEmpty() {
		return ""
	}
	top := s.Top()
	s.items = s.items[:len(s.items)-1]
	return top
}

// TOP ELEMENT
func (s *Stack) Top() string {
	if s.IsEmpty() {
		return ""
	}
	return s.items[len(s.items)-1]
}

// PRINT
func (s *Stack) Print() {
	if !s.IsEmpty() {
		for i := len(s.items) - 1; i >= 0; i-- {
			println(s.items[i])
		}
	} else {
		println("El stack esta vacio")
	}
}
