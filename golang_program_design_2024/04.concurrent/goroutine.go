package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
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
	safeGoroutine()
	anonymousFuncGoroutine()
	useChanStopGoroutine()
	useContentStopGoroutine()
}
func init() {
	// The default value is the number of CPU cores on the machine.
	runtime.GOMAXPROCS(2)
}

func sayHello() {
	fmt.Println("hello")
}
func goroutineHello() {
	// Use the 'go' keyword to create a goroutine.
	// func sayHello() will be executed asynchronously in a new goroutine.
	go sayHello()

	// the function will continue to execute without waiting for sayHello() to complete.
	// therefore, we need `time.Sleep` to pause the main goroutine so that the print statement in `sayHello` has a chance to be executed.
	fmt.Println("Main process")
	time.Sleep(1 * time.Second) // do not use sleep to wait for a goroutine to finish.
}

func worker(done chan bool) {
	fmt.Println("worker starting...")
	time.Sleep(time.Second)
	fmt.Println("done.")
	done <- true
}

// Goroutines should have clear start and end points, and avoid creating goroutines without termination conditions.
func safeGoroutine() {
	// a channel can be understood as a simple message queue, use "<-" to read and write queue data.
	done := make(chan bool, 1) // as done signal
	go worker(done)
	// wait for the goroutine to finish
	<-done
}

func anonymousFuncGoroutine() {
	done := make(chan bool, 1)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("task done.")

		done <- true
	}()

	// The main goroutine waits for the done signal.
	<-done

	fmt.Println("The main goroutine receives the done signal and continues to execute.")
}

// In most cases, the termination of the main program implicitly ends all goroutines.
// However, in long-running services, we may need to proactively stop a goroutine.
func useChanStopGoroutine() {
	stop := make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case <-stop:
				fmt.Println("Got the stop signal, stop...")
			default:
				fmt.Printf("Start Loop%d\n", i)
				// time.Sleep(time.Microsecond)
				i += 1
			}
		}
	}()

	stop <- struct{}{} // Send stop signal
}

func useContentStopGoroutine() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Got the stop signal. Shutting down...")
				return
			default:
				// execute normal operation
			}
		}
	}(ctx)

	// when you want to stop the goroutine
	cancel()
}
