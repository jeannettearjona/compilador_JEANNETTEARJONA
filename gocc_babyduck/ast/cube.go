package ast

import "fmt"

type SemanticCube map[string]map[string]map[string]string

func (cube SemanticCube) GetResultType(leftType, rightType, operator string) (string, error) {
	if op, ok := cube[operator]; ok {
		if left, ok := op[leftType]; ok {
			if result, ok := left[rightType]; ok {
				return result, nil
			}
		}
	}
	return "", fmt.Errorf("type mismatch: cannot apply operator '%s' between '%s' and '%s'", operator, leftType, rightType)
}

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
			"int": "bool",
		},
		"float": {
			"float": "bool",
		},
	},
	"<": {
		"int": {
			"int": "bool",
		},
		"float": {
			"float": "bool",
		},
	},
	"!=": {
		"int": {
			"int": "bool",
		},
		"float": {
			"float": "bool",
		},
	},
}
