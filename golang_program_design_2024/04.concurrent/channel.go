package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	initChannel()
	bufferedChannel()
	channelBufferedAndCapacity()
	selectChannel()
	rangeLoopChannel()
	concurrentChannel()
	time.Sleep(2 * time.Second)
	handleChannelError()
}

func initChannel() {
	ch := make(chan int)             // unbuffered channel
	chBuffered := make(chan int, 10) // buffered channel with capacity 10

	// sending data without a receiver: When you send data to an unbuffered channel, the program deadlocks.
	// this is because the send operation blocks until there is a receiver to receive the data.
	// ch <- 3  // This will cause a deadlock

	// start a goroutine to receive data, preventing deadlock
	go func() {
		value := <-ch
		fmt.Println(value)
	}()

	// send data to the channel
	ch <- 3

	// fmt.Println(<-ch) // this will block until data is sent from ch
	go func() {
		// use goroutine to receive data
		fmt.Println(<-ch)
	}()
	ch <- 4

	// close the channel
	close(ch)

	// optionally use chBuffered
	// here we simply send and receive a value
	chBuffered <- 5
	valueBuffered := <-chBuffered
	fmt.Println(valueBuffered)

	// close the buffered channel
	close(chBuffered)
}

func bufferedChannel() {
	// Create a buffered channel with capacity 2
	ch := make(chan int, 2)

	// Send two values to the channel
	ch <- 1
	ch <- 2

	// the buffer is now full; the next send would block until a value is received
	// ch <- 3 // Uncommenting this line would cause the program to block

	// receive and print the values
	fmt.Println(<-ch) // Output: 1
	fmt.Println(<-ch) // Output: 2

	// now we can send again as there is space in the buffer
	ch <- 3
	fmt.Println(<-ch) // output: 3

	// close the channel
	close(ch)
}

func channelBufferedAndCapacity() {
	ch1 := make(chan int)
	go func() {
		ch1 <- 1 // it will block here if there is no goroutine receiving
		close(ch1)
	}()

	ch2 := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			ch2 <- i // this won't block unless the channel is already full.
		}
		close(ch2)
	}()
}

// The select statement is very useful when choosing between multiple channels,
// similar to a switch statement but with each case statement being a channel operation.
// It can listen for data flow on a channel, and when multiple channels are ready simultaneously,
// select will randomly choose one to execute.
func selectChannel() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			ch1 <- i * 2
		}
		close(ch1) // don't forget to close the channel
	}()

	go func() {
		for i := 0; i < 3; i++ {
			ch2 <- i * 3
		}
		close(ch2)
	}()

	for i := 0; i < 5; i++ {
		select {
		case v1 := <-ch1:
			fmt.Println("received from ch1: ", v1)
		case v2 := <-ch2:
			fmt.Println("received from ch2: ", v2)
		}
	}
}

// when handling an unknown amount of data, using the range keyword allows for continuously receiving data
// from a Channel until it is closed.
func rangeLoopChannel() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	for n := range ch {
		fmt.Println("received: ", n)
	}
}

// operation1 simulates an operation that takes 1 second to complete
func operation1(ctx context.Context) {
	time.Sleep(1 * time.Second)
	select {
	case <-ctx.Done(): // Check if context is done
		fmt.Println("operation1 canceled")
		return
	default:
		fmt.Println("operation1 completed")
	}
}

// operation2 simulates an operation that takes 2 seconds to complete
func operation2(ctx context.Context) {
	time.Sleep(2 * time.Second)
	select {
	case <-ctx.Done(): // Check if context is done
		fmt.Println("operation2 canceled")
		return
	default:
		fmt.Println("operation2 completed")
	}
}

// the code uses `context.WithTimeout` to create a context that automatically cancels
// itself by sending a cancellation signal after the set duration elapses.
func concurrentChannel() {
	// Create a context with a timeout of 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the cancel function is called to release resources

	// Start operation1 and operation2 as goroutines
	go operation1(ctx)
	go operation2(ctx)

	// Wait for the context to be done
	<-ctx.Done()
	fmt.Println("main: context done")
}

// performTask simulates a task that either succeeds or fails based on its ID
func performTask(id int, errCh chan<- error) {
	// Simulate task execution
	if id%2 == 0 { // Tasks with even IDs are considered to fail
		time.Sleep(2 * time.Second)        // Simulate a delay to represent task processing time
		errCh <- errors.New("task failed") // Send an error to the error channel
	} else { // Tasks with odd IDs are considered to succeed
		fmt.Printf("task %d completed successfully\n", id)
		errCh <- nil // Send a nil value to indicate successful completion
	}
}

func handleChannelError() {
	tasks := 5
	// Create a buffered channel to hold errors from all tasks
	// The buffer size equals the number of tasks to prevent blocking
	errCh := make(chan error, tasks)

	// Launch each task as a separate goroutine
	for i := 0; i < tasks; i++ {
		go performTask(i, errCh)
	}

	// Collect and handle errors from the channel
	for i := 0; i < tasks; i++ {
		err := <-errCh // Receive error from the channel
		if err != nil {
			fmt.Printf("received error: %s\n", err) // Print the received error
		}
	}
	fmt.Println("finished processing all tasks") // Indicate that all tasks have been processed
}
