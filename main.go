package main

import (
	"fmt"
)

// Struct - Defined by user
type Address struct {
	Street string
	City   string
	State  string
}

type ReturnData struct {
	Name    string
	Age     int
	Address Address
} 

func (r ReturnData) Error() string {
	return "Error parsing new data"
}

func main() {
	var userType ReturnData

	refilledStructure, err := decodeData(`- Name: John Smith - Age: 30 - Address: 123 Main St, New York, NY`, &userType)
	if err != nil {
		panic(err)
	}

	hydrated, ok := refilledStructure.(*ReturnData)
	if !ok {
		panic("failed to hydrate")
	}

	// Range over the fields of the struct directly
	fmt.Printf("--- Name: %s\n", hydrated.Name)
	fmt.Printf("--- Age: %d\n", hydrated.Age)
	fmt.Printf("--- Address: %s, %s, %s\n", hydrated.Address.Street, hydrated.Address.City, hydrated.Address.State)
}
