package reflection

import "reflect"

type caller struct {
	value interface{}
}

type reflector struct {
	value     interface{}
	entryName map[string]reflect.Value
	entryType map[reflect.Type]reflect.Value
}
