![GitHub Logo](./censo.jpg)

# Censo [![GoDoc](https://godoc.org/github.com/shrotavre/censo?status.svg)](http://godoc.org/github.com/shrotavre/censo) [![Go Report Card](https://goreportcard.com/badge/shrotavre/censo)](https://goreportcard.com/report/github.com/shrotavre/censo)

Censo is a simple Go object/struct field omitter library.

## Usage
### Censorship Schema Definition

Example on how to construct censorship schema.

~~~ go
// main.go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shrotavre/censo"
)

func main() {
	// * Censo's censor model
	// CBasic -> reset to field's zero value
	// CSimple -> replace with value, or fallback to CBas
	// CFunc -> replace with value generated from function passed or fallback to CBas

	cschema := []censo.C{
		censo.CBas("FieldA"),               
		censo.CSim("FieldB", "****"),       
		censo.CFunc("FieldC", func(i interface{}) (o interface{}) {
			o = i

			if v, ok := i.(string); ok && strings.Contains(v, "real") {
				o = strings.ReplaceAll(v, "real", "fake")
			}

			return
		}),
	}
}

~~~

### Struct Censoring

Simple example on censoring a struct data.

~~~ go
// main.go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

func main() {
	targetstruct := Parent{
		FieldA: 1234,
		FieldB: "real_value",
		FieldC: "kinda_real_value",
		First: Child{
			FieldA: "very_real_value",
		},
	}

	err := censo.Censor(&targetstruct, cschema)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result:", targetstruct)
	// TODO: Add outputs with previously made schema
}

~~~

### JSON Map Censoring

Simple example on censoring a map data.

~~~ go
func main() {
	jsontarget := `{"FieldA":"real_value","FieldC":"real_value","First":{"FieldA":"very_real_value"}}`
	var targetmap map[string]interface{}
	err = json.Unmarshal([]byte(jsontarget), &targetmap)
	if err != nil {
		panic(err)
	}

	err = censo.Censor(&targetmap, cschema)
	if err != nil {
		panic(err)
	}

	fmt.Println("Result:", targetmap)

	// Outputs with previously made schema:
	// Target Map: map[FieldA: FieldC:fake_value First:map[FieldA:****]]
}
~~~

### Power Censoring*
PowerCensor will let you censor data and decide by yourself what censor strategy
to apply depending on field name + field value.

~~~ go
func main() {
	targetstruct = Parent{
		FieldA: 1234,
		FieldB: "real_value",
		FieldC: "kinda_real_value",
		First: Child{
			FieldA: "very_real_value",
		},
	}

	err = censo.PowerCensor(&targetstruct, func(fieldname string, fieldval interface{}) (placeholder interface{}) {
		placeholder = fieldval

		if _, ok := placeholder.(string); ok {
			placeholder = fieldname
		} else if _, ok := placeholder.(int); ok {
			placeholder = 9999
		}

		return
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("FieldA:", targetstruct.FieldA)
	fmt.Println("FieldB:", targetstruct.FieldB)
	fmt.Println("FieldC:", targetstruct.FieldC)
	fmt.Println("First/FieldA:", targetstruct.First.FieldA)

	// Outputs:
	// FieldA: 9999
	// FieldB: FieldB
	// FieldC: FieldC
	// First/FieldA: First/FieldA
}
~~~
