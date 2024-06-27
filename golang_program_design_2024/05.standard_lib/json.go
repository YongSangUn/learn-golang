package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

/*
Serialization is the process of converting a data structure or object into a format that can be stored or transmitted, such as a byte stream, JSON, or XML.

Deserialization is the reverse process, converting serialized data back into the original data structure or object.

Purpose:

- Data Storage: Save data to files or databases.
- Data Transmission: Transfer data between systems.
- Caching: Simplify storage format for complex objects.

Serialization and deserialization ensure data integrity and consistency during storage, transmission, and recovery across different systems.
*/
func main() {
	marshaling()
	structTagTest()
	unmarshaling()
}

func jsonM(data interface{}) {
	if jsonData, err := json.Marshal(data); err != nil {
		log.Fatalf("marshaling failed: %s", err)
	} else {
		fmt.Println(string(jsonData))
	}
}

func marshaling() {

	// marshaling struct
	jsonM(struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  "Tom",
		Email: "tom@test.com",
	})

	// marshaling map
	jsonM(map[string]interface{}{
		"name": "Jack",
		"age":  20,
	})

	// marshaling slice
	jsonM([]string{"name", "age"})

	time.Sleep(100 * time.Millisecond)
}

// Struct tags provide metadata for struct fields to control JSON serialization behavior.
type User struct {
	Name         string   `json:"name"`
	Biographical string   `json:"bio,omitempty"` // the "omitempty" option excludes empty fields from JSON serialization.
	Password     string   `json:"-"`             // the "-" tag tells json.Marshal to ignore this field.
	Email        []string `json:"email"`
}

func structTagTest() {
	jsonM(User{
		Name:     "Jackson",
		Password: "P@ssw0rd",
	})
	jsonM(User{
		Name:         "Jackson",
		Biographical: "This is Jackson.",
		Password:     "P@ssw0rd",
	})
}

func jsonUnm(str string, v any) {
	if err := json.Unmarshal([]byte(str), v); err != nil {
		log.Fatalf("unmarshaling failed: %s", err)
	} else {
		fmt.Printf("%#v\n", v)
	}
}

func unmarshaling() {
	jsonData1 := `{
	    "name": "Dick",
	    "bio":  "This is Dick <:",
		"email": ["dick@test.com", "dick@mail.com"]
	}`

	// unmarshaling struct
	var user User
	jsonUnm(jsonData1, &user)

	// dynamic unmarshaling
	jsonData2 := `{
        "name": "Alice",
        "details": {
            "age": 25,
            "job": "Engineer"
        }
    }`

	var result map[string]interface{}
	jsonUnm(jsonData2, &result)

	// Forced type conversion, ensure type matches before using
	name := result["name"].(string)
	fmt.Println("Name:", name)
	details := result["details"].(map[string]interface{})
	age := details["age"].(float64) // ! Note: Numbers in interface{} are treated as float64
	fmt.Println("Age:", age)
}
