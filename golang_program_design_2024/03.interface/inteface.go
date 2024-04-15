package main

import (
	"fmt"
	"math"
)

/*
An interface is an abstract type, and when all the methods in the interface are
	implemented, this interface is implicitly implemented.

definition of the interfaces:

	type interfaceName interface {
		methodName(parameterList) returnTypeList
	}

Detailed Explanation:
	https://draveness.me/golang/docs/part2-foundation/ch04-basic/golang-interface/
*/

func main() {
	interface_test()
	intefaceAdvanced()
}

/*
Interface and polymorphism:
Guidelines:

	https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces

Example by:

	https://coolshell.cn/articles/8460.html#接口和多态
*/
type shape interface {
	area() float64 // You almost never need a pointer to an interface.
	perimeter() float64
}

type rect struct {
	width, height float64
}

func (r *rect) area() float64 {
	return r.width * r.height
}
func (r *rect) perimeter() float64 {
	return 2 * (r.width + r.height)
}

type circle struct {
	radius float64
}

func (c *circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c *circle) perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func interface_test() {
	r := rect{width: 2.9, height: 4.8}
	c := circle{radius: 4.3}
	s := []shape{&r, &c} // achieved through interfaces
	for _, sh := range s {
		fmt.Println(sh)
		fmt.Println(sh.area())
		fmt.Println(sh.perimeter())
	}
}

func printAny(v interface{}) {
	fmt.Println(v)
}

// ------------------------ sep ------------------------

func intefaceAdvanced() {
	// empty interface, user for dynamic type processing.
	var a1 interface{} = 123
	var a2 interface{} = "hello world"
	var a3 interface{} = struct{ name string }{name: "golang"}

	printAny(a1)
	printAny(a2)
	printAny(a3)

	// type assertions, Operation that checks and converts a given value of a specified type from an interface.
	aInt, ok := a1.(int)
	fmt.Println(aInt, ok)
	aStr, ok := a2.(string)
	fmt.Println(aStr, ok)
	aStruct, ok := a3.(string) // converts false, type is `struct{ name string }`
	fmt.Println(aStruct, ok)
}
