package main

import (
	"fmt"
	"unsafe"
)

func main() {
	arrayInit()
	sliceInit()
	sliceAppend()
	issues()
}

func arrayInit() {
	// var myArray [n]T
	var nums1 = [5]int{1, 2, 3, 4, 5}
	fmt.Println(nums1)

	// Here ... means that the length of the array is calculated by the compiler.
	var nums2 = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(nums2)
}

func sliceInit() {
	slice1 := []int{1, 2, 3}
	// Create an array of length 5 with a capacity of 10
	slice2 := make([]int, 5, 10)
	fmt.Println(slice1, slice2)
	//Capacity is the maximum number of elements that the array underlying the slice can contain, starting from the start of the slice.
	fmt.Println(cap(slice2))

	// array and slice
	array := [5]int{10, 20, 30, 40, 50}
	slice := array[1:4]
	//Note that slicing does not actually copy the values of the array, it just points to a contiguous section of the original array.
	// Therefore, modification of the slice will also affect the underlying array and vice versa.
	slice[0] = 1
	fmt.Println(array, slice, cap(slice))
}

func sliceAppend() {
	slice := []int{1, 2, 3}
	fmt.Println(slice)
	slice = append(slice, 4)
	slice = append(slice, 5, 6)
	fmt.Println(slice)

	// If the underlying array's capacity is insufficient, the append operation will result in the slice pointing to a new, larger array.
	slice1 := make([]int, 4, 6) // When capacity is not enough, it will expand 6
	fmt.Printf("%d, %p, %p\n", cap(slice1), &slice1, unsafe.Pointer(&slice1[0]))

	slice1 = append(slice1, 3, 4)
	fmt.Printf("%d, %p, %p\n", cap(slice1), &slice1, unsafe.Pointer(&slice1[0])) // &slice1[0] not change.
	slice1 = append(slice1, 5, 6)
	fmt.Printf("%d, %p, %p\n", cap(slice1), &slice1, unsafe.Pointer(&slice1[0])) // cap=6+6 > 5, &slice1[0] was changed.
	//
}

func issues() {
	// index out of array.
	var arr [5]int
	index := 10 // 假设有一个超出范围的索引
	if index < len(arr) {
		fmt.Println(arr[index])
	} else {
		fmt.Println("index out of range.")
	}

	// sliced memory leak
	original := make([]int, 1000000)
	smallSlice := make([]int, 10)
	copy(smallSlice, original[:10])
	fmt.Println(cap(original), cap(smallSlice))
}
