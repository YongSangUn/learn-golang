package main

import (
	"fmt"
	"strconv"
	"sync"
)

/*
make(map[keyType]valueType)

Reference type: **

	map is a reference type. After it is created, a reference to the underlying data structure is actually obtained.

Dynamic resize:

	Similar to slice, it will dynamically expand as the data increases

Key uniqueness:

	Each key in a map is unique. If a value is stored using the same key, the new value will overwrite the original value.

Unordered collection:

	The elements in a map are unordered. Each time you traverse a map, the order of the key-value pairs may be different.
*/
func main() {
	mapInit()
	mapOpt()
	mapAdcanced()
}

// Notes on map initialization that the zero value of an uninitialized map is nil.
// at this point key-value pairs cannot be stored(must be init using make()), otherwise a runtime panic will be triggered.
func mapInit() {
	fmt.Println("- step01")
	var m map[string]int
	fmt.Println(m == nil)

	// running after panic.
	defer func() {
		fmt.Println("- step04")
		m = make(map[string]int)
		m["one"] = 3
		value, ok := m["one"]
		fmt.Println(value, ok)
	}()

	// running when panic
	defer func() {
		fmt.Println("- step03")
		if r := recover(); r != nil {
			// Recovered from panic
			fmt.Println("panic:", r)
		}
	}()

	fmt.Println("- step02")
	fmt.Println("m['one'] = 3")
	m["one"] = 3
}

func mapOpt() {
	// If the key does not exist, the zero value of the value type is returned.
	scores := map[string]int{
		"Tom":   90,
		"Jerry": 85,
	}
	tomS := scores["Tom"]
	bobS := scores["Bob"]
	fmt.Println(tomS, bobS)

	// Check if key exists.
	_, exists := scores["Bob"]
	fmt.Printf("Is Bob in map scores: %v\n", exists)

	// add and delete key-value
	scores["Alice"] = 92
	scores["Bob"] = 75
	delete(scores, "Bob")

	// range
	for k, v := range scores {
		// No guarantee that map will iterate in the same order each time.
		fmt.Println(k, v)
	}
}
func mapAdcanced() {
	// Specify a reasonable initial capacity for the map in advance to reduce
	// the overhead caused by the dynamic expansion of the map at runtime.
	myMap := make(map[string]int, 100)
	for i := 0; i < 102; i++ {
		myMap[fmt.Sprintln("no.%s", strconv.Itoa(i))] = i
	}
	fmt.Println(len(myMap))

	// map is a reference type, when map is assigned to another variable,
	// they both reference the same underlying data structure.
	nums := map[string]int{"one": 1, "two": 2}
	newNums := nums
	newNums["three"] = 3
	fmt.Printf("%+v", nums) // Write out the newly increased key-value pair

	// sync.Map
	var mySyncMap sync.Map
	mySyncMap.Store("Alice", 23)
	mySyncMap.Store("Bob", 25)
	value, ok := mySyncMap.Load("Alice")
	fmt.Println(value, ok)
	// range method for iterating over sync.Map
	mySyncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true // continue to iterate
	})
}
