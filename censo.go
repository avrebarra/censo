package censo

import (
	"reflect"
)

// This package uses Go's reflections https://blog.golang.org/laws-of-reflection
// It's not bcs anybody are tired, reflections is actually hard to comprehend
// it is normal. It's better not to use reflection anywhere, not transparent.

// TODO: implement deep matching + check for possible panics sources

// C represent censorship schema
type C struct {
	Field       string
	Placeholder interface{}
}

// CBas will create new basic censorship schema -> will set to null value of field
func CBas(f string) (c C) {
	return C{Field: f}
}

// CSim will create new simple censorship schema
func CSim(f string, p interface{}) (c C) {
	return C{f, p}
}

// Censor will censor matching schema in target and replace
func Censor(target interface{}, schemas []C) (err error) {
	sv := reflect.ValueOf(target).Elem()
	st := sv.Type()

	if sv.Kind() == reflect.Struct {
		for i := 0; i < st.NumField(); i++ {
			fieldVal := sv.Field(i)

			// Search field in schema list
			for _, sch := range schemas {
				if fieldVal.CanSet() && st.Field(i).Name == sch.Field {
					rep := reflect.ValueOf(sch.Placeholder)

					// check if replacement value matches as field's value's type
					if sch.Placeholder != nil && rep.Type() == fieldVal.Type() {
						fieldVal.Set(rep)
					} else {
						// set to zero value
						fieldVal.Set(reflect.Zero(fieldVal.Type()))
					}
				}
			}
		}
	}
	return
}
