package ast

type Quadruple struct {
	Operador string
	Izq      string
	Der      string
	Res      string
}

func NewQuadruple(op, izq, der, res string) Quadruple {
	return Quadruple{
		Operador: op,
		Izq:      izq,
		Der:      der,
		Res:      res,
	}
}

type Queue struct {
	items []Quadruple
}

// CHECK IF EMPTY
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// ADD
func (q *Queue) Enqueue(value Quadruple) {
	q.items = append(q.items, value)
}

// REMOVE
func (q *Queue) Dequeue() {
	if q.IsEmpty() {
		return
	}
	//value := q.items[0]
	q.items = q.items[1:]
	//return value
}

// FRONT ELEMENT
func (q *Queue) Front() Quadruple {
	if q.IsEmpty() {
		return Quadruple{}
	}
	return q.items[0]
}

// SIZE
func (q *Queue) Size() int {
	return len(q.items)
}

// PRINT QUADRUPLES
func (q *Queue) Print() []Quadruple {
	return q.items
}

/*func (q *Queue) GetItem(index int) Quadruple {
	if index < 0 || index >= len(q.items) {
		return Quadruple{}
	}
	return q.items[index]
}

func (q *Queue) Update(index int, value Quadruple) {
	if index < 0 || index >= len(q.items) {
		return
	}
	q.items[index] = value
}
*/
