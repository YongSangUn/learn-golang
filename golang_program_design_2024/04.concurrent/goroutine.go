package main

import (
	"fmt"
	"runtime"
)

/*
Concurrency:

	Refers to the handling of multiple tasks within the same time period, but only one task is executed at any given moment. Tasks switch rapidly between each other, giving the user the illusion that they are being executed simultaneously. Concurrency is suitable for single-core processors.

Parallelism:

	Refers to the actual simultaneous execution of multiple tasks at the same moment, which requires the support of multi-core processors.

The Go language is designed with concurrent design as one of its main goals. It achieves an efficient concurrent programming model through Goroutines and Channels.

The Go runtime manages Goroutines, scheduling them across multiple system threads to achieve parallel processing.
*/

/*
The scheduling of goroutines is handled by the scheduler within the go runtime.
GO's scheduler uses m:n scheduling technology (multiple goroutines mapped to multiple os threads).

Three important entities: M (Machine), P (Processor), and G (Goroutine):

M(corresponds to the kernel thread):

	Represents the machine or thread, it is an abstraction of the OS kernel thread.

P(represents the context during scheduling):

	Is a collection of resources needed to execute a Goroutine. Each P has a local Goroutine queue.

G(is the specific Goroutine):

	Represents a Goroutine, which includes information such as the Goroutine's execution stack and instruction set.
*/

func main() {
	goroutineHello()
}
func init() {
	// The default value is the number of CPU cores on the machine.
	runtime.GOMAXPROCS(2)
}

func sayHello() {
	fmt.Println("hello")
}
func goroutineHello() {
	go sayHello()

	// the function will continue to execute without waiting for sayHello() to complete.
	// therefore, we need `time.Sleep` to pause the main goroutine so that the print statement in `sayHello` has a chance to be executed.
	// time.Sleep(1 * time.Second)
	fmt.Println("Main process")
}
