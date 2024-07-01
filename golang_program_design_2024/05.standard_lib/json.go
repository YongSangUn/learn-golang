package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	marshalError()
	custom()
	encodeJson()
	decodeJson()
}

func jsonM(data interface{}) {
	if jsonData, err := json.Marshal(data); err != nil {
		log.Printf("marshaling failed: %s", err)
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
		log.Printf("unmarshaling failed: %s", err)
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

type UserError struct {
	Name string
	Age  int
	// Assuming there is a field here that cannot be serialized.
	// Data chan struct{} // The channel cannot be represented in JSON.
}

func marshalError() {
	// marshaling error
	u1 := UserError{
		Name: "Alice",
		Age:  30,
		// Data: make(chan struct{}),
	}
	bytes, err := json.Marshal(u1)
	if err != nil {
		log.Printf("marshaling failed: %v", err)
	}

	// unmarshaling error
	fmt.Println(string(bytes))
	// "age" should be an integer, but a string is given here.
	var data = []byte(`{"name":"Alice","age":"unknown"}`)
	var u2 UserError
	err = json.Unmarshal(data, &u2)
	if err != nil {
		// log.Printf("unmarshaling failed: %v", err)
	}
	fmt.Printf("%+v\n", u2)
}

// These processes can be customized by implementing the json.Marshaler and json.Unmarshaler interfaces.

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (c Color) MarshalJSON() ([]byte, error) {
	hex := fmt.Sprintf("\"#%02x%02x%02x\"", c.Red, c.Green, c.Blue)
	return []byte(hex), nil
}

func (c *Color) UnmarshalJSON(data []byte) error {
	_, err := fmt.Sscanf(string(data), "\"#%02x%02x%02x\"", &c.Red, &c.Green, &c.Blue)
	return err
}

func custom() {
	c := Color{Red: 255, Green: 99, Blue: 71}

	jsonColor, _ := json.Marshal(c)
	fmt.Println(string(jsonColor))

	var newColor Color
	json.Unmarshal(jsonColor, &newColor)
	fmt.Printf("%#v\n", newColor)
}

// JSON data can be directly written to any object that implements the io.Writer interface,
// meaning that JSON data can be directly encoded to files, network connections, and more.
func encodeJson() {
	users := []UserError{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	file, _ := os.Create("users.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		log.Fatalf("encode error: %v", err)
	}

}

// json.Decoder can read JSON data directly from any object that implements the io.Reader interface, seeking and parsing JSON objects and arrays.
func decodeJson() {
	file, _ := os.Open("users.json")

	var u []UserError
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&u); err != nil {
		log.Fatalf("decode error: %v", err)
	}
	fmt.Println(u)
}
