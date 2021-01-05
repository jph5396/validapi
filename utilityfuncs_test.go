package validapi

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Name  string
	Score float64
	ID    int
	Bool  bool
}

func TestPropsFromType(t *testing.T) {
	propGroup := PropsFromType(reflect.TypeOf(testStruct{}))

	testData := []struct {
		name string
		typ  Type
	}{
		{"Name", String},
		{"Score", Float},
		{"ID", Int},
		{"Bool", Boolean},
	}

	for _, i := range testData {
		if prop, ok := propGroup.properties[i.name]; ok {
			if prop.getType() != i.typ {
				t.Errorf("types do not match. want %v got %v", i.typ.String(), prop.getType().String())
			}
		} else {
			t.Errorf("prop %v should be in propertygroup but was not", i.name)
		}
	}
}
