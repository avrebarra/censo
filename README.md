![GitHub Logo](./censo.jpg)

# Censo [![GoDoc](https://godoc.org/github.com/shrotavre/censo?status.svg)](http://godoc.org/github.com/shrotavre/censo) [![Go Report Card](https://goreportcard.com/badge/shrotavre/censo)](https://goreportcard.com/report/github.com/shrotavre/censo)

Censo is a simple Go object/struct field omitter library made using reflections under the hood.

## Usage

~~~ go
// main.go
package main

import (
	"fmt"

	"github.com/shrotavre/censo"
)

type Parent struct {
	FieldA int
	FieldB string
	FieldC string

	First Child
}

type Child struct {
	FieldA string
}

var censorship []censo.C = []censo.C{
	censo.CBas("FieldA"),
	censo.CSim("FieldB", "****"),
	censo.CSim("First/FieldA", "****"),
}

func main() {
	data := Parent{
		FieldA: 1234,
		FieldB: "real_value",
		FieldC: "kinda_real_value",

		First: Child{
			FieldA: "very_real_value",
		},
	}

	err := censo.Censor(&data, censorship)
	if err != nil {
		panic(err)
	}

	fmt.Println("FieldA:", data.FieldA)
	fmt.Println("FieldB:", data.FieldB)
	fmt.Println("First/FieldA:", data.First.FieldA)

	// Outputs:
	// FieldA: 0
	// FieldB: ****
	// First/FieldA: ****
}

~~~
