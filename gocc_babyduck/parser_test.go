package main

import (
	"testing"

	"fmt"
	"gocc_babyduck/ast"
	"gocc_babyduck/lexer"
	"gocc_babyduck/parser"
)

type TI struct {
	src    string // El código o fragmento de entrada
	expect int64  // El valor esperado del resultado
}

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
				c,d: float;
			void funcion1 (param1 : int)
			[ var varLocal: int; { b = varLocal + 2; }];
			main {
				if (a < b)
				{
        			print("a es menor que b");
					print("Hello", 42, c + d);
					c = 4 - 2 /  (d * 1.5);
    			}
    			else
    			{
					print("a no es menor que b");
					c = d - 1.5;
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
				c,d: float;
			void funcion1 (param1 : int, param2 : int)
			[ var varLocal: int; { b = varLocal + 2; a = param1 + param2; }];
			main {
				if (a < b)
				{
        			print("a es menor que b");
    			};
    			print("Fin del ifelse");
				c = 4 - 2 /  (d * 1.5);
			} end`,
		0,
	},
}*/

/*var testData = []*TI{
	{
		`program xyz;
    var a, b: int;

    main {
        a = 0;
        b = 3;

        while (a < b) do {
            print("a es", a);
        };

        print("Fin del ciclo");
    }
end
`,
		0,
	},
}*/

/*var testData = []*TI{
	{
		`program test;
			var a, b, c: int;

			main {
    			a = 2;
    			b = 5;

    			if (a < b) {
        			print("a menor que b");
        			a = a + 1;
    			} else {
        			print("a no es menor que b");
        			a = a - 1;
    			};

    			while (a < b) do {
        			print("en el while");
        			a = a + 1;
    			};
			}
		end`,
		0,
	},
}*/

var testData = []*TI{
	{
		`program test;
			var a, b: float;
			main {
    			a = 20 + 6.5;
				b = 20 - 6.5;
				print(a);
				print(b);
			}
		end`,
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
				vm := ast.NewVirtualMachine(&ast.Cuadruplos, ast.ConstantsVarTable, ast.FunctionDirectory)
				vm.Run()
				fmt.Println(vm.Memory)
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
