package validapi

import (
	"fmt"
	"reflect"
	"strconv"
)

//Props the interface that should represent a single property in a json object.
// currently Property and Object Property implement it.
type Props interface {
	getName() string
	getType() Type
	validate(string, interface{}) error
}

//Property represents a single property in a request body.
type Property struct {
	Name     string
	propType Type
	rules    []Rule
}

//NewProperty creates a property with a blank rule set.
func NewProperty(name string, typ Type) *Property {
	return &Property{
		Name:     name,
		propType: typ,
		rules:    []Rule{},
	}
}

func (p Property) getName() string {
	return p.Name
}

func (p Property) getType() Type {
	return p.propType
}

//AddRules will take the rules provided and add them to the Property,
// checking if they are valid first. If not, it will print a msg stating
// it has been ignored.
func (p *Property) AddRules(rules ...Rule) *Property {
	for _, r := range rules {
		err := r.rulevalidation(p)
		if err == nil {
			p.rules = append(p.rules, r)
		} else {
			panic(fmt.Errorf("could not add rules to Property %v. error: %v", p.Name, err.Error()))
		}
	}
	return p
}

func (p Property) validate(key string, value interface{}) error {
	valueType := reflect.TypeOf(value)

	//When Json is decoded in go, all JSON numbers are converted to float64 types.
	// this means we need to handle integer checking differently.
	if p.propType == Int {
		if !valueType.ConvertibleTo(Int) {
			return fmt.Errorf("%v: invalid type. got %v, want %v", key, valueType.String(), p.propType.String())
		}

	} else if valueType != p.propType {
		return fmt.Errorf("%v: invalid type. got %v, want %v", key, valueType.String(), p.propType.String())
	}
	for _, rule := range p.rules {
		err := rule.validate(value)
		if err != nil {
			return fmt.Errorf("%v: %v", key, err.Error())
		}
	}
	return nil
}

//PropertyGroup wrapper used to be sure property names are unique when applied to a route.
type PropertyGroup struct {
	properties map[string]Props
}

//NewPropertyGroup creates a PropertyGroup with no properties.
func NewPropertyGroup() *PropertyGroup {
	return &PropertyGroup{properties: make(map[string]Props)}
}

//AddProperties attempts to add properties to PropertyGroup. It will throw an error if any Properties have
// conflicting names or aliases.
func (pg *PropertyGroup) AddProperties(props ...Props) *PropertyGroup {
	for _, prop := range props {
		if _, present := pg.properties[prop.getName()]; !present {
			pg.properties[prop.getName()] = prop
		} else {
			panic(fmt.Errorf("duplicated Prop name: %v", prop.getName()))
		}
	}
	return pg
}

func (pg *PropertyGroup) validateGroup(body map[string]interface{}) error {
	for key, val := range body {
		if property, ok := pg.properties[key]; ok {
			err := property.validate(key, val)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf(" %v is not a valid Property", key)
		}
	}

	return nil
}

//ObjectProperty represents a property that would be an object type in json
// instead of a basic type. contains a group of properties that will validate the
// items contained in the json object.
type ObjectProperty struct {
	Name     string
	propType Type
	slice    bool
	group    *PropertyGroup
}

func (o ObjectProperty) getName() string {
	return o.Name
}

func (o ObjectProperty) getType() Type {
	return o.propType
}

//NewObjectProperty creates a new Object Property with the name provided and sets the slice var
func NewObjectProperty(name string, slice bool) *ObjectProperty {
	return &ObjectProperty{
		Name:     name,
		slice:    slice,
		propType: Group,
		group:    NewPropertyGroup(),
	}
}

//UsePropertyGroup sets the property group on the object. this replaces the entire propertygroup
// that is already created when calling NewObjectProperty. the primary use case for it is when
// reusing a propertygroup from another route.
func (o *ObjectProperty) UsePropertyGroup(pg *PropertyGroup) *ObjectProperty {
	o.group = pg
	return o
}

//AddProperties add Base Properties to the property group of the object property.
func (o *ObjectProperty) AddProperties(p ...Props) *ObjectProperty {
	o.group.AddProperties(p...)
	return o
}

func (o ObjectProperty) validate(key string, val interface{}) error {
	if o.slice && reflect.TypeOf(val).Kind() == reflect.Slice {
		reflectVal := reflect.ValueOf(val)
		for i := 0; i < reflectVal.Len(); i++ {
			err := objectvalidator(strconv.Itoa(i), reflectVal.Index(i).Interface(), o.group.properties)
			if err != nil {
				return fmt.Errorf("%v.%v", key, err.Error())
			}
		}
	} else {
		err := objectvalidator(key, val, o.group.properties)
		if err != nil {
			return err
		}
	}
	return nil
}

// function used by objectProperty.validate to validate a value. It has been
// written here so it can be used in both standard and array instances.
func objectvalidator(key string, val interface{}, props map[string]Props) error {
	if reflect.TypeOf(val).Kind() == reflect.Map {
		mapIter := reflect.ValueOf(val).MapRange()
		for mapIter.Next() {
			k := mapIter.Key().String()
			v := mapIter.Value().Interface()
			if prop, ok := props[k]; ok {
				err := prop.validate(k, v)
				if err != nil {
					return fmt.Errorf("%v.%v", key, err.Error())
				}
			} else {
				return fmt.Errorf("%v is not a valid prop", k)
			}
		}
	} else {
		return fmt.Errorf("%v not a valid type. got %v want Object", key, reflect.TypeOf(val).Kind().String())
	}
	return nil
}
