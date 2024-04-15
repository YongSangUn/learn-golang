package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	structInit()
	structJson()
	structCopy()
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	// "omitempty" means that if a field is empty or missing, it will be omitted from the JSON.
	Emails []string `json:"emails,omitempty"`
}

func structInit() {
	// A pointer to a newly allocated Person type variable,
	// whose member variables are initialized with zero values.
	p1 := new(Person)
	fmt.Printf("%+v\n", p1)

	// Initialize with fields.
	p2 := Person{
		Name:   "Alice",
		Age:    30,
		Emails: []string{"alice@example.com", "alice123@example.com"},
	}
	fmt.Printf("%+v\n", p2)

	// Initialize without field_name. However, it is necessary to ensure that the
	// order of the initial values of each member variable is the same as the
	// order when defining the structure, and that no field can be omitted.
	p3 := Person{"Bob", 25, []string{"bob@example.com"}}
	fmt.Printf("%+v\n", p3)

	// Anonymous struct
	p4 := struct {
		Name string
		Age  int
	}{
		Name: "Eve",
		Age:  40,
	}
	fmt.Println("Name:", p4.Name)
}

func structJson() {
	// serialization and
	p1 := Person{
		Name:   "John Doe",
		Age:    30,
		Emails: []string{"john@example.com", "j.doe@example.com"},
	}
	jsonData, err := json.Marshal(p1)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("JSON format: %s\n", jsonData)

	// deserialization.
	var p2 Person
	if err := json.Unmarshal(jsonData, &p2); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	fmt.Printf("Recovered Struct: %#v\n", p2)
	// multline string(use raw string literal)
	jsonString := `{
		"name": "John Doe",
		"age": 30,
		"city": "San Francisco"
	}`
	var p3 Person
	// key:<city> will not unmarshal.
	if err := json.Unmarshal([]byte(jsonString), &p3); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	fmt.Printf("%#v\n", p3)
}

type User struct {
	Name string
	Age  int
}

type Data struct {
	Numbers []int
}

func structCopy() {
	// Struct copy by assignment.
	// Deep copy: If a structure contains only primitive types (such as int,
	//   string, etc.), copying is a deep copy.
	p1 := struct {
		Name string
		Age  int
	}{
		Name: "Eve",
		Age:  40,
	}
	p2 := p1
	p1.Name = "Tom"
	fmt.Printf("%#v\n%#v\n", p1, p2)

	u1 := User{"Eve", 40}
	u2 := u1
	u1.Name = "Tom"
	fmt.Printf("%#v\n%#v\n", u1, u2)

	// Shallow copy: If the structure contains a reference type (such as a slice or map),
	//   then the copy will be a shallow copy, and the original instance and the
	//   newly copied instance will share the memory of the reference type.
	original := Data{Numbers: []int{1, 2, 3}}
	copied := original
	copied.Numbers[0] = 100
	fmt.Println("Original:", original.Numbers) // output: Original: [100 2 3]
	fmt.Println("Copied:", copied.Numbers)     // output: Copied: [100 2 3]

	// This problem can be avoided by explicitly copying the content of a slice
	//   into a new slice, thus achieving a true deep copy.
	newNumbers := make([]int, len(original.Numbers))
	copy(newNumbers, original.Numbers)
	copied2 := Data{Numbers: newNumbers}
	fmt.Printf("%+v\n", copied2)
}
