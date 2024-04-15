package main

import "fmt"

/*
The underlying principle of defer is to use a stack (Last In First Out principle) to store each deferred function.
When a defer statement is encountered, the Go language does not immediately execute the function after the statement,
but pushes it into a dedicated stack. Only when the outer function is about to return,
these deferred functions will be executed in the order of the stack, that is,
the function in the last declared defer statement will be executed first.
*/

func example() {
	defer fmt.Println("world") // deferred
	fmt.Println("hello")
}

// Execute in the order of last in first out(LIFO).
func multipleDefers() {
	defer fmt.Println("First defer")
	defer fmt.Println("Second defer")
	defer fmt.Println("Third defer")

	fmt.Println("Function body")
}

func main() {
	example()
	multipleDefers()

}
