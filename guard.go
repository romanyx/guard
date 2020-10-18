package guard

import (
	"fmt"
	"reflect"
)

// MustNotNil panics if the interface is nil
// or its value is nil.
func MustNotNil(
	pN int,
	name string,
	i interface{},
) bool {
	if i == nil {
		panic(fmt.Sprintf(
			"parameter %s number %d is nil",
			name,
			pN,
		))
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr,
		reflect.Map,
		reflect.Array,
		reflect.Chan,
		reflect.Slice:
		if reflect.ValueOf(i).IsNil() {
			panic(fmt.Sprintf(
				"parameter %s number %d value is nil",
				name,
				pN,
			))
		}
	}

	return false
}
