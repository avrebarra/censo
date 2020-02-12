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

// CSet is set/array of C
type CSet []C

// CBas will create new basic censorship schema -> will set to null value of field
func CBas(f string) (c C) {
	return C{Field: f}
}

// CSim will create new simple censorship schema
func CSim(f string, p interface{}) (c C) {
	return C{Field: f, Placeholder: p}
}

// Censor will censor matching schema in target and replace
func Censor(target interface{}, set CSet) (err error) {
	sv := reflect.ValueOf(target).Elem()
	cpmap := convertCSetToCPMap(set)

	censor(sv, cpmap, "")

	return
}

func censor(sv reflect.Value, cpmap map[string]interface{}, keyprefix string) {
	// id := time.Now().UnixNano()
	// fmt.Println("CENSOR", id, sv, sv.Kind())

	if sv.Kind() == reflect.Struct {
		censorStruct(sv, cpmap, keyprefix)
	} else if sv.Kind() == reflect.Map {
		censorMap(sv, cpmap, keyprefix)
	}

	// fmt.Println("RESULT", id, sv)
}

func censorStruct(sv reflect.Value, cpmap map[string]interface{}, keyprefix string) {
	st := sv.Type()

	for i := 0; i < st.NumField(); i++ {
		cpkey := keyprefix + st.Field(i).Name

		fieldVal := sv.Field(i)
		placeholder, found := cpmap[cpkey]

		if fieldVal.Kind() == reflect.Struct || fieldVal.Kind() == reflect.Map {
			censor(fieldVal, cpmap, cpkey+"/")
		} else if fieldVal.CanSet() && found {
			rep := reflect.ValueOf(placeholder)

			// check if replacement value matches as field's value's type
			if placeholder != nil && rep.Type() == fieldVal.Type() {
				fieldVal.Set(rep)
			} else { // set to zero value
				fieldVal.Set(reflect.Zero(fieldVal.Type()))
			}
		}
	}
}

func censorMap(sv reflect.Value, cpmap map[string]interface{}, keyprefix string) {
	// fmt.Println("  MAP", sv)
	// fmt.Println("  MAP", cpmap)
	for _, fname := range sv.MapKeys() {
		cpkey := keyprefix + fname.String()

		fieldVal := sv.MapIndex(fname)

		placeholder, found := cpmap[cpkey]
		placeholderv := reflect.ValueOf(placeholder)

		// fmt.Println("   K", fname, cpkey)
		// fmt.Println(fieldVal.Elem().Kind())
		// fmt.Println()

		if (fieldVal.Elem().Kind() == reflect.Struct || fieldVal.Elem().Kind() == reflect.Map) && !found {
			censor(fieldVal.Elem(), cpmap, cpkey+"/")
			continue
		}

		if !found {
			continue
		}

		// replace orig data with C
		if placeholder != nil && (placeholderv.Type() == fieldVal.Elem().Type() || placeholderv.Type().ConvertibleTo(fieldVal.Type())) {
			fv := placeholderv.Convert(fieldVal.Elem().Type())
			sv.SetMapIndex(fname, fv)
		} else { // set to zero value
			sv.SetMapIndex(fname, reflect.Zero(fieldVal.Elem().Type()))
		}
	}
}
