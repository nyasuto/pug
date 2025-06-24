package phase1

import (
	"testing"
)

// BenchmarkEvalInteger benchmarks integer arithmetic evaluation
func BenchmarkEvalInteger(t *testing.B) {
	input := "((2 + 3) * 4 - 1) / 2"

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalFloat benchmarks floating point arithmetic evaluation
func BenchmarkEvalFloat(t *testing.B) {
	input := "((2.5 + 3.7) * 4.1 - 1.2) / 2.0"

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalFibonacci benchmarks recursive function calls
func BenchmarkEvalFibonacci(t *testing.B) {
	input := `
	let fib = fn(x) {
		if (x < 2) {
			return x;
		} else {
			return fib(x - 1) + fib(x - 2);
		}
	};
	fib(10);
	`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalFactorial benchmarks iterative function with variable assignment
func BenchmarkEvalFactorial(t *testing.B) {
	input := `
	let factorial = fn(n) {
		let result = 1;
		let i = 1;
		while (i <= n) {
			result = result * i;
			i = i + 1;
		}
		return result;
	};
	factorial(10);
	`

	// Note: This uses a conceptual while loop syntax
	// The actual implementation might need adjustment based on language features
	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalStringOperations benchmarks string concatenation
func BenchmarkEvalStringOperations(t *testing.B) {
	input := `"Hello" + " " + "World" + "!"`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalArrayOperations benchmarks built-in array functions
func BenchmarkEvalArrayOperations(t *testing.B) {
	input := `
	let arr = [1, 2, 3, 4, 5];
	let newArr = push(arr, 6);
	let firstItem = first(newArr);
	let lastItem = last(newArr);
	let restItems = rest(newArr);
	len(restItems);
	`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalClosures benchmarks closure creation and execution
func BenchmarkEvalClosures(t *testing.B) {
	input := `
	let makeAdder = fn(x) {
		return fn(y) { x + y };
	};
	let add5 = makeAdder(5);
	add5(10);
	`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalComplexExpression benchmarks a complex nested expression
func BenchmarkEvalComplexExpression(t *testing.B) {
	input := `
	let compute = fn(a, b, c) {
		let temp1 = a * b + c;
		let temp2 = temp1 / (a + b);
		if (temp2 > c) {
			return temp2 * 2;
		} else {
			return temp2 + c;
		}
	};
	compute(5, 10, 3);
	`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalVariableAccess benchmarks variable access in nested scopes
func BenchmarkEvalVariableAccess(t *testing.B) {
	input := `
	let x = 10;
	let y = 20;
	let z = 30;
	let compute = fn() {
		let a = x + y;
		let b = a * z;
		return b / x;
	};
	compute();
	`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalBuiltinFunctions benchmarks built-in function calls
func BenchmarkEvalBuiltinFunctions(t *testing.B) {
	input := `
	len("Hello World");
	type(42);
	type("string");
	type(true);
	`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkEvalEnvironmentOperations benchmarks environment operations
func BenchmarkEvalEnvironmentOperations(t *testing.B) {
	for i := 0; i < t.N; i++ {
		env := NewEnvironment()

		// Simulate heavy environment usage
		for j := 0; j < 100; j++ {
			env.Set("var"+string(rune(j)), &Integer{Value: int64(j)})
		}

		// Access variables
		for j := 0; j < 100; j++ {
			env.Get("var" + string(rune(j)))
		}
	}
}

// BenchmarkEvalParseAndEval benchmarks combined parsing and evaluation
func BenchmarkEvalParseAndEval(t *testing.B) {
	input := `
	let quicksort = fn(arr) {
		if (len(arr) <= 1) {
			return arr;
		}
		
		let pivot = first(arr);
		let less = [];
		let greater = [];
		
		// This is a conceptual implementation
		// Real implementation would need array manipulation
		return arr;
	};
	
	quicksort([3, 1, 4, 1, 5, 9, 2, 6]);
	`

	for i := 0; i < t.N; i++ {
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()
		env := NewEnvironment()
		Eval(program, env)
	}
}
