package main

import "fmt"

func main() {
	fmt.Println("-> callback")
	traverse([]int{1, 2, 3}, func(n int) {
		fmt.Println(n * n)
	})

	// closure
	seqFunc := sequenceGenerator()
	fmt.Println(seqFunc()) // output 1
	fmt.Println(seqFunc()) // output 2
}

// as callback func
func traverse(numbers []int, callback func(int)) {
	for _, num := range numbers {
		callback(num)
	}
}

// A closure is a function value that references variables from outside its body.
// The function may access and assign to the referenced variables;
// in this sense the function is "bound" to the variables.
func sequenceGenerator() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
