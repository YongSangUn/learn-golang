package main

import "fmt"

func main() {
	fmt.Println("-> mult params")
	fmt.Println(sum(1, 2, 3, 4))

	fmt.Println("-> pass var")
	passvar()
}

// mult params
func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

// pass by value and pass by reference(ptr)
func double(val int) {
	val *= 2
}
func doublePtr(val *int) {
	*val *= 2
}
func passvar() {
	value := 123
	double(value)
	fmt.Println(value)
	doublePtr(&value)
	fmt.Println(value)
}
