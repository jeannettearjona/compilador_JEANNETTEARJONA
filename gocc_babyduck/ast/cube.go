package ast

//import "fmt"

type SemanticCube map[string]map[string]map[string]string

//esto no se si es necesario, pero lo dejo por si acaso
/*func (cube SemanticCube) GetResultType(operator, leftType, rightType string) (string, error) {
	if op, ok := cube[operator]; ok {
		if left, ok := op[leftType]; ok {
			if result, ok := left[rightType]; ok {
				return result, nil
			}
		}
	}
	return "", fmt.Errorf("invalid operation: %s %s %s", leftType, operator, rightType)
}*/

var DefaultSemanticCube = SemanticCube{
	"=": {
		"int": {
			"int": "int",
		},
		"float": {
			"float": "float",
		},
	},
	"+": {
		"int": {
			"int":   "int",
			"float": "float",
		},
		"float": {
			"int":   "float",
			"float": "float",
		},
	},
	"-": {
		"int": {
			"int":   "int",
			"float": "float",
		},
		"float": {
			"int":   "float",
			"float": "float",
		},
	},
	"*": {
		"int": {
			"int":   "int",
			"float": "float",
		},
		"float": {
			"int":   "float",
			"float": "float",
		},
	},
	"/": {
		"int": {
			"int":   "int",
			"float": "float",
		},
		"float": {
			"int":   "float",
			"float": "float",
		},
	},
	">": {
		"int": {
			"int": "int",
		},
		"float": {
			"float": "int",
		},
	},
	"<": {
		"int": {
			"int": "int",
		},
		"float": {
			"float": "int",
		},
	},
	">=": {
		"int": {
			"int": "int",
		},
		"float": {
			"float": "int",
		},
	},
	"<=": {
		"int": {
			"int": "int",
		},
		"float": {
			"float": "int",
		},
	},
	"!=": {
		"int": {
			"int": "int",
		},
		"float": {
			"float": "int",
		},
	},
}
