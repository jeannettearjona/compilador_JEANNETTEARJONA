package ast

type HashMap struct {
	data map[string]interface{} // Cambié uint8 a interface{} para almacenar cualquier tipo de dato
}

func NewHashMap() *HashMap {
	return &HashMap{
		data: make(map[string]interface{}),
	}
}

// IS EMPTY
func (m *HashMap) IsEmpty() bool {
	return len(m.data) == 0
}

// ADD
func (m *HashMap) Add(key string, value interface{}) {
	m.data[key] = value
}

// GET
func (m *HashMap) Get(key string) (interface{}, bool) {
	val, exists := m.data[key]
	return val, exists
}

// REMOVE
func (m *HashMap) Remove(key string) {
	delete(m.data, key)
}

// MODIFY
func (m *HashMap) Modify(key string, value interface{}) {
	m.data[key] = value
}

// PRINT
func (m *HashMap) Print() {
	if !m.IsEmpty() {
		for key, val := range m.data {
			println(key, ":", val)
		}
	} else {
		println("El HashMap esta vacio")
	}
}

// CONTAINS
func (m *HashMap) Contains(key string) bool {
	_, exists := m.data[key]
	return exists
}

func (h *HashMap) Keys() []string {
	keys := make([]string, 0, len(h.data))
	for key := range h.data {
		keys = append(keys, key)
	}
	return keys
}
