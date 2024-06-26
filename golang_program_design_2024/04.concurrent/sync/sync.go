package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
In concurrent programming, when multiple goroutines share resources, it is necessary to ensure that the resource is accessed by only one goroutine at any given moment to guarantee data consistency and state synchronization.

Golang supports multiple synchronization mechanisms:

- Mutual exclusion locks (sync.Mutex) and read-write mutual exclusion locks (sync.RWMutex)
- Channels
- WaitGroups
- Atomic functions (atomic package)
- Condition variables (sync.Cond)
*/

func main() {
	syncMutex()
	syncRWMutex()
	syncCond()
	syncAtomic()
	syncOnce()

	if err := errGroup(); err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	}
}

func syncMutex() {
	// without sync
	var counter1 int
	for i := 0; i < 1000; i++ {
		go func() {
			counter1++
		}()
	}
	time.Sleep(2 * time.Second) // wait goroutines done.
	fmt.Printf("without sync counter1: %d\n", counter1)

	// use sync.Mutex
	var mu sync.Mutex
	var counter2 int
	// need to use sync.WaitGroup or other ways to wait for all coroutines to complete.
	// sync.WaitGroup is a synchronization mechanism to wait for a group
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1) // call the Add() of WaitGroup to increase the counter
		go func() {
			defer wg.Done() // Done() (Add(-1)) to decrease the counter when each goroutine is finished.
			mu.Lock()
			counter2++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("use sync counter2: %d\n", counter2)
}

func syncRWMutex() {
	// Without sync.RWMutex
	fmt.Println("=== Without sync.RWMutex ===")
	unsafeCounter := 0
	for i := 0; i < 10; i++ {
		go func() {
			unsafeCounter++
			fmt.Printf("Write: %d\n", unsafeCounter)
		}()
		go func() {
			fmt.Printf("Read: %d\n", unsafeCounter)
		}()
	}
	time.Sleep(time.Millisecond * 100)

	// With sync.RWMutex
	fmt.Println("\n=== With sync.RWMutex ===")
	var mu sync.RWMutex
	safeCounter := 0
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			mu.Lock()
			safeCounter++
			fmt.Printf("Write: %d\n", safeCounter)
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			mu.RLock()
			fmt.Printf("Read: %d\n", safeCounter)
			mu.RUnlock()
		}()
	}
	wg.Wait()
}

/*
sync.Cond is less frequently used in Go concurrent programming for the following reasons:

- higher complexity compared to Mutex and RWMutex
- specific use cases (mainly for "wait-notify" patterns)
- often replaceable by channels
- requires deeper knowledge of concurrent programming
*/

// TODO: sync.Cond
func syncCond() {
}

// syncAtomic simulates concurrent visitors to a website using atomic operations
// and displays the changing visitor count over time.
func syncAtomic() {
	var visitorCount int32
	var wg sync.WaitGroup
	done := make(chan bool)

	// Goroutine to periodically print the visitor count
	go func() {
		for { // for{ ... } keeps the goroutine running indefinitely.
			select { // used for non-blocking channel operations.
			case <-done:
				return
			default:
				fmt.Printf("Current visitors: %d\n", atomic.LoadInt32(&visitorCount))
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Simulate 100 concurrent visitors
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// safely reads the current value of visitorCount, controls the frequency of updates.
			atomic.AddInt32(&visitorCount, 1)
			time.Sleep(time.Duration(500+i*10) * time.Millisecond) // Varying visit durations
			atomic.AddInt32(&visitorCount, -1)
		}()
	}

	wg.Wait()
	done <- true // a true value is sent on the done channel, the goroutine returns and terminates.
	fmt.Printf("Final visitor count: %d\n", visitorCount)
}

func syncOnce() {
	// sync.Once ensures that the function it wraps is executed only once,
	// even if called multiple times concurrently.
	var once sync.Once

	// simulates a resource-intensive operation
	expensiveOperation := func() {
		fmt.Println("Performing expensive operation...")
		time.Sleep(2 * time.Second)
		fmt.Println("Expensive operation completed.")
	}
	var wg sync.WaitGroup

	// launch 5 goroutines that all try to execute the expensive operation
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("Goroutine %d is trying to execute the operation\n", id)

			// once.Do ensures the wrapped function is called only once
			// across all goroutines
			once.Do(expensiveOperation)

			fmt.Printf("Goroutine %d has finished\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All goroutines have finished execution")
}

func errGroup() error {
	fmt.Println("==> errGroup")
	// Create a new context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure all resources are freed at the end

	// Create a new errgroup.Group
	// This group will help us manage multiple goroutines and collect their errors
	group, ctx := errgroup.WithContext(ctx)

	// Slice of URLs we want to fetch
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.invalid-url-for-error-demo.com",
	}

	// For each URL, start a goroutine to fetch it
	for _, url := range urls {
		url := url // Create a new variable to avoid closure problems

		// Add a new goroutine to the group
		group.Go(func() error {
			// Check if the context has been cancelled before making the request
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Make the HTTP GET request
				resp, err := http.Get(url)
				if err != nil {
					return fmt.Errorf("failed to fetch %s: %v", url, err)
				}
				defer resp.Body.Close()

				fmt.Printf("Successfully fetched %s, status: %s\n", url, resp.Status)
				return nil
			}
		})
	}

	// Wait for all goroutines to complete and collect any error
	if err := group.Wait(); err != nil {
		return fmt.Errorf("one of the goroutines failed: %v", err)
	}

	fmt.Println("All URLs were fetched successfully")
	return nil
}
