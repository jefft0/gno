package main

import "fmt"

func foo(a uint64, b float32) {
	fmt.Printf("%T\n", a)
	println(a)
}

func main() {
	x := 11
	foo(1<<x, 1>>x)
}

// Error:
// main/files/types/shift_d12.gno:12:2: operator >> not defined on: Float32Kind

// TypeCheckError:
// main/files/types/shift_d12.gno:12:12: invalid operation: shifted operand 1 (type float32) must be integer
