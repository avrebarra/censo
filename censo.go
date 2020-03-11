package censo

import (
	"reflect"
)

// This package uses Go's reflections https://blog.golang.org/laws-of-reflection
// It's not bcs anybody are tired, reflections is actually hard to comprehend
// it is normal.
//
// Lets be fair it's better not to use reflection anywhere, it's not transparent

const (
	universalFieldMatcherSign = "*"
)

// C represent Censo's censorship schema
type C struct {
	Field       string      // Field defines what field name to apply censor
	Placeholder interface{} // Placeholder defines value to replace orig value
}

// CSet is set/array of C
type CSet []C

// CFPMap is a mapping of field name (string) to placeholder/replacer object
type CFPMap map[string]interface{}

// Censor will censor matching schema in target and replace with defined placeh
// older value.
// This func was designed as ignorant func. It doesn't actually care/doesnt halt
// processing even if any error happened. Worst case happened, the original data
// would not be touched at all (even if error happened).
func Censor(target interface{}, set CSet) (err error) {
	return censor(
		reflect.ValueOf(target).Elem(),
		CSetToCPMap(set),
		"",
	)
}

// PowerCensor let you censor data and decide by yourself what censor strategy
// to apply depending on field name + field value.
func PowerCensor(target interface{}, cf func(fieldname string, fieldval interface{}) (placeholder interface{})) (err error) {
	powerC := C{
		Field:       universalFieldMatcherSign,
		Placeholder: cf,
	}

	return censor(
		reflect.ValueOf(target).Elem(),
		CSetToCPMap([]C{powerC}),
		"",
	)
}

func censor(sv reflect.Value, cfpmap CFPMap, keyprefix string) (err error) {
	// try to cast interfaces to map[string]interface
	v, ok := sv.Interface().(map[string]interface{})
	if ok {
		sv = reflect.ValueOf(v)
	}

	// perform censoring
	switch sv.Kind() {
	case reflect.Struct:
		err = censorStruct(sv, cfpmap, keyprefix)
		return

	case reflect.Map:
		_, ok := sv.Interface().(map[string]interface{})
		if !ok {
			return ErrNotCensorable
		}
		err = censorMap(sv, cfpmap, keyprefix)
		return

	default:
		return ErrNotCensorable

	}
}

func censorStruct(sv reflect.Value, cfpmap CFPMap, keyprefix string) (err error) {
	st := sv.Type()

	for i := 0; i < st.NumField(); i++ {
		// decide key to search in cfp
		cfpkey := keyprefix + st.Field(i).Name

		// fetch original field value
		fieldVal := sv.Field(i)

		// if field value is map/struct do recursive process
		if fieldVal.Kind() == reflect.Struct || fieldVal.Kind() == reflect.Map {
			censor(fieldVal, cfpmap, cfpkey+"/")
			continue
		}

		// fetch placeholder value
		placeholder, cfpkeyfound := cfpmap[cfpkey]
		if !cfpkeyfound {
			placeholder, cfpkeyfound = cfpmap[universalFieldMatcherSign]
		}
		if !cfpkeyfound {
			continue // if cfpkey or universal not listed, skip field
		}
		placeholderv := reflect.ValueOf(placeholder)

		// if placeholder is a replacerfunc, create placeholder first
		if placeholderv.Kind() == reflect.Func {
			if powerfunc, ok := placeholder.(func(fieldname string, fieldval interface{}) (placeholder interface{})); ok {
				placeholder = powerfunc(cfpkey, fieldVal.Interface())
				placeholderv = reflect.ValueOf(placeholder)
			} else if placeholderfunc, ok := placeholder.(func(i interface{}) (o interface{})); ok {
				placeholder = placeholderfunc(fieldVal.Interface())
				placeholderv = reflect.ValueOf(placeholder)
			} else {
				continue
			}
		}

		// replace orig data with placeholder
		if fieldVal.CanSet() {
			rep := reflect.ValueOf(placeholder)

			if placeholder != nil && rep.Type() == fieldVal.Type() {
				// direct set if replacement matches original value's type
				fieldVal.Set(rep)
			} else {
				// set to zero value
				fieldVal.Set(reflect.Zero(fieldVal.Type()))
			}
		}
	}

	return
}

func censorMap(sv reflect.Value, cfpmap CFPMap, keyprefix string) (err error) {
	for _, fname := range sv.MapKeys() {
		// decide key to search in cfp
		cfpkey := keyprefix + fname.String()

		// fetch original field value
		fieldVal := sv.MapIndex(fname)

		// if field value is map/struct do recursive process
		if fieldVal.Elem().Kind() == reflect.Struct || fieldVal.Elem().Kind() == reflect.Map {
			censor(fieldVal.Elem(), cfpmap, cfpkey+"/")
			continue
		}

		// fetch placeholder value
		placeholder, cfpkeyfound := cfpmap[cfpkey]
		if !cfpkeyfound {
			placeholder, cfpkeyfound = cfpmap[universalFieldMatcherSign]
		}
		if !cfpkeyfound {
			continue // if cfpkey or universal not listed, skip field
		}
		placeholderv := reflect.ValueOf(placeholder)

		// if placeholder is a replacerfunc, create placeholder first
		if placeholderv.Kind() == reflect.Func {
			if powerfunc, ok := placeholder.(func(fieldname string, fieldval interface{}) (placeholder interface{})); ok {
				placeholder = powerfunc(cfpkey, fieldVal.Interface())
				placeholderv = reflect.ValueOf(placeholder)
			} else if placeholderfunc, ok := placeholder.(func(i interface{}) (o interface{})); ok {
				placeholder = placeholderfunc(fieldVal.Interface())
				placeholderv = reflect.ValueOf(placeholder)
			} else {
				continue
			}
		}

		// replace orig data with placeholder
		if placeholder != nil && (placeholderv.Type() == fieldVal.Elem().Type() || placeholderv.Type().ConvertibleTo(fieldVal.Type())) {
			// replace with value if applicable
			fv := placeholderv.Convert(fieldVal.Elem().Type())
			sv.SetMapIndex(fname, fv)
		} else {
			// set to zero value if not applicable or not defined
			sv.SetMapIndex(fname, reflect.Zero(fieldVal.Elem().Type()))
		}
	}

	return
}
