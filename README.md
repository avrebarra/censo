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

type Dummy struct {
	FieldA string
	FieldB int
	FieldC string
}

var censorship []censo.C = []censo.C{
	censo.CBas("FieldB"),
	censo.CSim("FieldA", "****"),
}

func main() {
	data := Dummy{
		FieldA: "real_value",
		FieldB: 1234,
		FieldC: "real_value",
	}

	err := censo.Censor(&data, censorship)
	if err != nil {
		panic(err)
	}

	fmt.Println("FieldA:", data.FieldA)
	fmt.Println("FieldB:", data.FieldB)
	fmt.Println("FieldC:", data.FieldC)

	// Outputs:
	// FieldA: ****
	// FieldB: 0
	// FieldC: real_value
}
~~~
