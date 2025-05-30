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

/*var testData = []*TI{
	{
		`program test;
			var a, b: int;
			    c: float;
			void funcion1 (param1 : int, param2 : int, param3: float)
			[ var varLocal: int; { varLocal = 4; b = varLocal - 2; print(b); a = param1 - param2; print(a);}];
			main {
    			a = 20 + 6;
				b = 18;
				c = 55.55;
				print(a);
				print(c);
				if(c > 1.5 * a){
					print("ENTRO AL IF", 2+6);
				};
				print("donde deberia caer gotof");
				funcion1(10, 5, 3.2);
				print("END");
			}
		end`,
		0,
	},
}*/

/*var testData = []*TI{
	{
		`program testFactorial;
		var n, resultado: int;

		void factorial(num: int)
		[
			var i, result: int;
			{
				result = 1;
				i = 1;
				while (i < num+1) do {
					result = result * i;
					i = i + 1;
				};
				resultado = result;
			}
		];

		void printFactorial()
		[
		{
			print("El factorial es", resultado);
		}
		];

		main {
			n = 5;
			resultado = 0;
			factorial(n);
			printFactorial();
		}
	end`,
		0,
	},
}*/

var testData = []*TI{
	{
		`program test8;
			var y: int;
			void funcion() 
			[ var x: int;
				{print("func", y);}
			];
						
			void second() 
			[ var x: int;
				{funcion(); 
				x = 9;}
			];
						
			main {
				y = 10;
							
				second();
			}
			end`,
		0,
	},
}

/*var testData = []*TI{
	{
		`program testFibonacci;
		var n, resultado: int;

		void fibonacciIter(num: int)
		[
			var a, b, i, temp: int;
			{
				a = 0;
				b = 1;
				i = 0;
				while (i < num) do {
					temp = b;
					b = a + b;
					a = temp;
					i = i + 1;
				};
				resultado = a;
			}
		];

		main {
			n = 10;
			resultado = 0;
			fibonacciIter(n);
			print("Fibonacci de", n, "es", resultado);
		}
	end`,
		0,
	},
}*/

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
