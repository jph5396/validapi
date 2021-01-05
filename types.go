package validapi

import (
	"reflect"
)

//Type used as a wrapper for reflect.Types in case I wish to change
// how types are defined in the package. variables to be used with
// property are exported with this Type.
type Type reflect.Type

//String string type variable.
var String Type = reflect.TypeOf("")

//Int int type variable.
var Int Type = reflect.TypeOf(1)

//Float float type variable.
var Float Type = reflect.TypeOf(1.12)

//Boolean boolean type variable.
var Boolean Type = reflect.TypeOf(true)

//Group PropertyGroup type variable.
var Group = reflect.TypeOf(PropertyGroup{})
