package main

import (
	"testing"

	"gocc_babyduck/ast"
	"gocc_babyduck/lexer"
	"gocc_babyduck/parser"
)

type TI struct {
	src    string // El código o fragmento de entrada
	expect int64  // El valor esperado del resultado
}

/*
var testData = []*TI{
	//VARIABLE DOBLEMENTE DECLARADA A NIVEL GLOBAL
	{"program A1; var x: int; x: float; main { x = 5; } end", -1},

	//VARIABLE DOBLEMENTE DECLARADA GLOBAL Y LOCAL
	{"program A1; var a,b: int; void fun1 (s : int) [ var a: int; { x = 2; }]; main { x = 5; } end", 0}, //Ok: funcion con asignación sencilla en el cuerpo, //MAL: 2 funciones con el mismo nombre

	//FUNCION DOBLEMENTE DECLARADA
	{"program A1; var a,b: int; void fun1 (z : int) [{}]; void fun1 (z : int) [{}]; main { x = 5; } end", -1},

	{"program A1; var x: int; main { x = 5; } end", 0}, // Caso ok
	// Cambios en ID
	{"program s1A; var x: int; main { x = 5; } end", 0}, //Ok: ID puede empezar con letra minuscula
	{"program saa; var x: int; main { x = 5; } end", 0}, //Ok: ID puede tener solo letras minusculas
	{"program SAA; var x: int; main { x = 5; } end", 0}, //Ok: ID puede tener solo letras mayusculas
	{"program S22; var x: int; main { x = 5; } end", 0}, //Ok: ID puede tener numeros seguidos despues de que empieza con letra
	{"program S; var x: int; main { x = 5; } end", 0},   //Ok: ID puede ser una sola letra
	// Cambios en var
	{"program A1; main { } end", 0},                                //Ok: si var es vacia
	{"program A1; var a,b: int; main { x = 5; } end", 0},           //Ok: var puede tener varios ID separados por coma
	{"program A1; var a,b: int; c: float; main { x = 5; } end", 0}, //Ok: var puede tener varios ID separados por coma y diferentes tipos
	{"program A1; var a,b: int; c: char; main { x = 5; } end", -1}, //Bad: no hay tipo char
	// Cambios en funciones
	{"program A1; var a,b: int; void fun1 (z : int) [{}]; main { x = 5; } end", 0}, //Ok: funcion sin variables locales ni cuerpo
	// funciones con assign
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2; }]; main { x = 5; } end", 0},         //Ok: funcion con asignación sencilla en el cuerpo
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 > 1; }]; main { x = 5; } end", 0},     //Ok: funcion de asignacion con comparacion > se supone que no se vale operador con asignacion
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 < 3; }]; main { x = 5; } end", 0},     //Ok: funcion de asignacion con comparacion <
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 != 2.1; }]; main { x = 5; } end", 0},  //Ok: funcion de asignacion con comparacion !=
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 == 2.1; }]; main { x = 5; } end", -1}, //Bad: no existe el operador ==
	// funciones con assign y diferentes operaciones
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = (2 > 1); }]; main { x = 5; } end", 0},   //Ok: Uso de parentesis para operaciones
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = (2+2) > 1; }]; main { x = 5; } end", 0}, //Ok: Uso de parentesis para operaciones compuestas combinadas con comparacion
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 - 1; }]; main { x = 5; } end", 0},     //Ok: Operacion de resta
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 * 1; }]; main { x = 5; } end", 0},     //Ok: Operacion de multiplicacion
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 / 1; }]; main { x = 5; } end", 0},     //Ok: Operacion de division
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = 2 % 1; }]; main { x = 5; } end", -1},    //Bad: no existe el operador %
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = +2; }]; main { x = 5; } end", 0},        //Ok: Operacion de suma unaria
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = -2; }]; main { x = 5; } end", 0},        //Ok: Operacion de resta unaria
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { a = +b; }]; main { x = 5; } end", 0},        //Ok: Operacion de suma unaria con variable
	// funciones con condicionales
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { if (2){}; }]; main { x = 5; } end", 0},                       //Ok: if de cuerpo vacio sin else
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { if (2){}else{}; }]; main { x = 5; } end", 0},                 //Ok if de cuerpo vacio con else con cuerpo vacio
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { if (2>1){a = +2;}else{a = -2;}; }]; main { x = 5; } end", 0}, //Ok: if de comparacion con operacion en cuerpo y else con operacion en cuerpo
	// funciones con while
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { while (b) do {}; }]; main { x = 5; } end", 0},               //Ok: while de cuerpo vacio
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { while ((b+2)!=a) do {a = +2;}; }]; main { x = 5; } end", 0}, //Ok: while de comparacion con operacion en cuerpo
	// funciones con f calls
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { idFcall ((2 * 1) != 0); }]; main { x = 5; } end", 0},        //Ok: llamada a funcion con unico argumento
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { idFcall ((2 * 1) != 0, a > b); }]; main { x = 5; } end", 0}, //Ok: llamada a funcion con varios argumentos
	// funciones con print
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { print (\"hola mundo\"); }]; main { x = 5; } end", 0}, //Ok: print con constante string
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { print (a+b); }]; main { x = 5; } end", 0},            //Ok: print con una expresion
	{"program A1; var a,b: int; void fun1 (z : int) [ var varFun1: int; { print (a+b, (b+2)!=a); }]; main { x = 5; } end", 0},  //Ok: print con varias expresiones
}*/

/*var testData = []*TI{
	{
		`program xyz;
			var a,b,z: int;
				c,d: float;
			void funcion1 (param1 : int)
			[ var varLocal: int; { b = varLocal + 2; }];
			main { print("Hello", 42, c + d); c = 4 - 2 /  (d * 1.5); } end`,
		0,
	},
}*/

/*var testData = []*TI{
	{
		`program xyz;
			var a,b: int;
			void funcion1 (param1 : int)
			[ var varLocal: int; { b = varLocal + 2; }];
			main {
				if (a < b)
				{
        			print("a es menor que b");
    			}
    			else
    			{
					print("a no es menor que b");
    			};

    			print("Fin del ifelse");
			} end`,
		0,
	},
}*/

/*var testData = []*TI{
	{
		`program xyz;
			var a,b: int;
			void funcion1 (param1 : int)
			[ var varLocal: int; { b = varLocal + 2; }];
			main {
				if (a < b)
				{
        			print("a es menor que b");
    			};
    			print("Fin del ifelse");
			} end`,
		0,
	},
}*/

var testData = []*TI{
	{
		`program xyz;
    var a, b: int;

    main {
        a = 0;
        b = 3;

        while (a < b) do {
            print("a es", a);
            a = a + 1;
        };

        print("Fin del ciclo");
    }
end
`,
		0,
	},
}

func TestParser(t *testing.T) {
	p := parser.NewParser()

	for _, ts := range testData {
		// Resetear el estado semantico antes de cada prueba
		ast.ResetSemanticState()

		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)

		if ts.expect == 0 {
			if err != nil {
				t.Errorf("Expected no error but got: %v for input: %s", err, ts.src)
			} else {
				ast.ImprimirCuadruplos()
			}
		} else if ts.expect == -1 {
			if err == nil {
				t.Errorf("Expected error but got none for input: %s", ts.src)
			} else {
				t.Logf("Correctly got expected error: %v", err)
			}
		}
	}
}
