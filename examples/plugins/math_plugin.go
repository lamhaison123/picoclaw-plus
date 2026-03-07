//go:build tinygo

// PicoClaw Example WASM Plugin
// This is a simple math plugin demonstrating WASM plugin capabilities

package main

//export add
func add(a, b int32) int32 {
	return a + b
}

//export subtract
func subtract(a, b int32) int32 {
	return a - b
}

//export multiply
func multiply(a, b int32) int32 {
	return a * b
}

//export divide
func divide(a, b int32) int32 {
	if b == 0 {
		return 0 // Handle division by zero
	}
	return a / b
}

//export factorial
func factorial(n int32) int32 {
	if n <= 1 {
		return 1
	}
	result := int32(1)
	for i := int32(2); i <= n; i++ {
		result *= i
	}
	return result
}

//export fibonacci
func fibonacci(n int32) int32 {
	if n <= 1 {
		return n
	}
	a, b := int32(0), int32(1)
	for i := int32(2); i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

//export power
func power(base, exp int32) int32 {
	if exp == 0 {
		return 1
	}
	result := base
	for i := int32(1); i < exp; i++ {
		result *= base
	}
	return result
}

// Main function required by TinyGo
func main() {}
