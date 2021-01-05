package validapi

import (
	"fmt"
	"reflect"
)

//exposed funcs to make working with this package easier.

//PropsFromType receives a reflect.Type of a struct and
// returns a propertygroup based of the field name and types of the struct
// useful for creating propertygroups that dont need any specific rules applied to them.
// will panic if the provided type's kind is not a struct.
func PropsFromType(t reflect.Type) PropertyGroup {
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("myapi.PropsFromType received kind %v. wanted struct", t.Kind()))
	}

	propGroup := NewPropertyGroup()
	for i := 0; i < t.NumField(); i++ {
		var prop Property
		field := t.Field(i)
		name := field.Name
		switch field.Type {
		case Int:
			prop = NewProperty(name, Int)
		case String:
			prop = NewProperty(name, String)
		case Float:
			prop = NewProperty(name, Float)
		case Boolean:
			prop = NewProperty(name, Boolean)
		default:
			panic(fmt.Errorf("type of %v not supported", field.Type.String()))
		}

		err := propGroup.AddProperties(prop)
		if err != nil {
			panic(err)
		}
	}

	return propGroup
}
