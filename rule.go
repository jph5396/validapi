package validapi

import (
	"fmt"
	"reflect"
	"regexp"
)

//Rule Interface that defines the common interfaced that should be used when
// implementing a rule type all implemented rules can assume that the provided value
// is already an appropriate type, because the value will have been type checked before
// and invalid ones will not make it to the rule.
// NOTE: rules cannot be applied to the Object Property Type.
type Rule interface {
	//validate the function used to define if an input is valid or not.
	// if an input is not valid, the function should return false, and a string
	// representing the eror message that should be sent as part of the response.
	validate(interface{}) (bool, string)

	//rulevalidation should check to be sure the rule can be applied to given property.
	rulevalidation(Props) error
}

//RegexRule checks to see if the  propety value of a property matches the provided regex string.
type RegexRule struct {
	regexStr string
}

//NewRegexRule accepts a string, confirms it's a valid regex pattern, and returns a rule if the pattern is
//valid. If it is not, it will return a blank rule and an error noting that the regex is invalid.
func NewRegexRule(str string) (RegexRule, error) {

	_, err := regexp.Compile(str)
	if err != nil {
		return RegexRule{}, err
	}

	return RegexRule{
		regexStr: str,
	}, nil
}

func (r RegexRule) validate(i interface{}) (bool, string) {

	value := i.(string)
	//Note: we ignore the error because it should have already been
	// confirmed that the expression string compiles when it was created.
	regex, _ := regexp.Compile(r.regexStr)
	if regex.MatchString(value) {
		return true, ""
	}
	msg := fmt.Sprintf("%v does not match regex pattern %v", value, r.regexStr)
	return false, msg
}

func (r RegexRule) rulevalidation(p Props) error {
	if p.getType() != String {
		err := fmt.Errorf("regex rule cannot be used with property. got type %v, need string", p.getType().String())
		return err
	}
	return nil
}

//EnumRule checks to see if the property value is within a set of valid values.
type EnumRule struct {
	enumType   Type
	enumvalues map[interface{}]struct{}
}

//NewEnumRule checks to see if all members provided are the same type,
// if so, it will return a valid EnumRule. if not, it will return
// a blank rule and a error. type is inferred by checking the first member of the array.
func NewEnumRule(members []interface{}, t Type) (EnumRule, error) {

	enumType := t
	enumvalues := make(map[interface{}]struct{})
	var empty struct{}
	for _, val := range members {
		if reflect.TypeOf(val) != enumType {
			err := fmt.Errorf("enum type mismatch")
			return EnumRule{}, err
		}
		enumvalues[val] = empty
	}
	return EnumRule{
		enumType:   enumType,
		enumvalues: enumvalues,
	}, nil
}

func (r EnumRule) validate(i interface{}) (bool, string) {
	if _, ok := r.enumvalues[i]; ok {
		return true, ""
	}

	msg := fmt.Sprintf("%v not in enum list", i)
	return false, msg
}

func (r EnumRule) rulevalidation(p Props) error {

	if r.enumType != p.getType() {
		err := fmt.Errorf("enum type and prop type do not match. enum %v, prop type %v", r.enumType.String(), p.getType().String())
		return err
	}

	return nil
}

//CustomRule a rule type that allows users to define their own rule.
type CustomRule struct {
	name        string
	description string
	t           Type
	validation  func(interface{}) (bool, string)
}

//NewCustomRule creates a new custom rule that can be applied to properties.
func NewCustomRule(n string, typ Type, v func(interface{}) (bool, string)) CustomRule {
	return CustomRule{
		name:       n,
		t:          typ,
		validation: v,
	}
}
func (cr CustomRule) validate(i interface{}) (bool, string) {
	return cr.validation(i)
}

//SetDescription sets the rules description.
func (cr *CustomRule) SetDescription(desc string) {
	cr.description = desc
}

func (cr CustomRule) rulevalidation(p Props) error {
	if cr.t != p.getType() {
		return fmt.Errorf(" rule %v cannot be applied to prop %v with type %v", cr.name, p.getName(), p.getType().String())
	}
	return nil
}
