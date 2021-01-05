package validapi

import (
	"testing"
)

func TestProperty(t *testing.T) {
	prop1 := NewProperty("Name", String)
	prop2 := NewProperty("ID", Int)

	//basic test.
	err := prop1.validate("test", "test")
	if err != nil {
		t.Errorf("got %v wanted nil", err.Error())
	}

	// make sure int validation will still return an err
	// when a string is passed.
	err = prop2.validate("int test", "test")
	if err == nil {
		t.Error("wanted error got nil")
	}

	//confirm a float64 can be passed and verified as int.
	err = prop2.validate("intasfloat", float64(202001))
	if err != nil {
		t.Error(err.Error())
	}

	enumvals := []interface{}{1, 2, 3}
	//test rule on prop.
	rule, err := NewEnumRule(enumvals, Int)
	if err != nil {
		t.Errorf("could not build rule for testing. err: %v", err.Error())
	}
	err = prop2.AddRules(rule)
	if err != nil {
		t.Errorf("could not place rule on prop. %v", err.Error())
	}

	err = prop2.validate("int", 1)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = prop2.validate("int fail ", 12)
	if err == nil {
		t.Errorf("got nil wanted error")
	}
}

func TestComplextProperty(t *testing.T) {
	// build properties.
	prop1 := NewProperty("Name", String)
	prop2 := NewProperty("ID", Int)
	prop3 := NewProperty("score", Float)

	propgroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	propgroup.AddProperties(prop1, prop2, prop3)
	objProp := NewObjectProperty("User", false)
	objProp.UsePropertyGroup(propgroup)

	var test = map[string]interface{}{
		"User": map[string]interface{}{
			"Name":  "Jimbo",
			"ID":    43,
			"score": 23.45,
		},
		"UserFail": map[string]interface{}{
			"Name":  "Tomas",
			"ID":    43.4,
			"score": "test",
		},
	}

	err := objProp.validate("User", test["User"])
	if err != nil {
		t.Errorf("wanted nil, got %v", err.Error())
	}

	err = objProp.validate("UserFail", test["UserFail"])
	if err == nil {
		t.Errorf("wanted error got nil")
	}

}

func TestNestedProperties(t *testing.T) {
	prop1 := NewProperty("Name", String)
	prop2 := NewProperty("ID", Int)
	prop3 := NewProperty("score", Float)

	propgroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	propgroup.AddProperties(prop1, prop2, prop3)
	objProp := NewObjectProperty("User", false)
	objProp.UsePropertyGroup(propgroup)

	computer := NewProperty("computer", String)
	status := NewProperty("status", String)
	supergroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	supergroup.AddProperties(computer, status, objProp)
	superObjProp := NewObjectProperty("PCID", false)
	superObjProp.UsePropertyGroup(supergroup)
	var test = map[string]interface{}{
		"computer": "testPC",
		"status":   "working",
		"User": map[string]interface{}{
			"Name":  "Jimbo",
			"ID":    43,
			"score": 23.45,
		},
	}

	err := superObjProp.validate("nesteduser", test)
	if err != nil {
		t.Errorf("wanted nil got %v", err.Error())
	}

}

func TestArrayOfProperties(t *testing.T) {
	prop1 := NewProperty("Name", String)
	prop2 := NewProperty("ID", Int)
	prop3 := NewProperty("score", Float)

	propgroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	propgroup.AddProperties(prop1, prop2, prop3)
	objProp := NewObjectProperty("Users", true)
	objProp.UsePropertyGroup(propgroup)

	var test = map[string]interface{}{
		"Users": []map[string]interface{}{
			{
				"Name":  "Jimbo",
				"ID":    43,
				"score": 23.45,
			},
			{
				"Name":  "Steven",
				"ID":    22,
				"score": 223.45,
			},
			{
				"Name":  "Paul",
				"ID":    11,
				"score": 2.45,
			},
		},
	}

	err := objProp.validate("list", test["Users"])
	if err != nil {
		t.Errorf("Wanted nil got %v", err.Error())
	}
}
