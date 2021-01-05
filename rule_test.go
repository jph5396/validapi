package validapi

import (
	"testing"
)

func TestRegexRule(t *testing.T) {
	success := "^([A-Z])+$"
	fail := "()[A-Z"

	t.Run("Should create", func(t *testing.T) {
		_, err := NewRegexRule(success)
		if err != nil {
			t.Errorf("err should be nil, but got %v", err.Error())
		}
	})

	t.Run("Should fail to create", func(t *testing.T) {
		_, err := NewRegexRule(fail)
		if err == nil {
			t.Errorf("wanted an error, got nil")
		}
	})

	t.Run("Should fail to add to prop", func(t *testing.T) {
		prop := NewProperty("test", Int)
		rule, err := NewRegexRule(success)
		err = prop.AddRules(rule)
		if err == nil {
			t.Error("wanted an error, got nil")
		}
	})

	t.Run("Should pass validation", func(t *testing.T) {
		rule, _ := NewRegexRule(success)
		pass, msg := rule.validate("ABC")
		if !pass {
			t.Errorf("should have passed, but got %v instead", msg)
		}
	})

	t.Run("Should fail validation", func(t *testing.T) {
		rule, _ := NewRegexRule(success)
		pass, _ := rule.validate("A12BC")
		if pass {
			t.Error("should have failed, but passed instead")
		}
	})

}

func TestEnumRule(t *testing.T) {

	t.Run("Build String Enum", func(t *testing.T) {
		_, err := NewEnumRule([]interface{}{"test", "test2"}, String)
		if err != nil {
			t.Errorf("wanted nil got %v", err.Error())
		}
	})

	t.Run("Build Int Enum", func(t *testing.T) {
		_, err := NewEnumRule([]interface{}{1, 123}, Int)
		if err != nil {
			t.Errorf("wanted nil got %v", err.Error())
		}
	})

	t.Run("Should fail to build", func(t *testing.T) {
		_, err := NewEnumRule([]interface{}{1, "123"}, Int)
		if err == nil {
			t.Errorf("wanted err got nil")
		}
	})

	enum, _ := NewEnumRule([]interface{}{1, 123}, Int)

	t.Run("Should Apply to Prop", func(t *testing.T) {
		prop := NewProperty("test", Int)
		err := prop.AddRules(enum)
		if err != nil {
			t.Errorf("should pass but got %v instead", err.Error())
		}
	})

	t.Run("Should validate", func(t *testing.T) {
		pass, msg := enum.validate(1)
		if !pass {
			t.Errorf("Should have passed but got %v", msg)
		}
	})

	t.Run("Should fail to validate", func(t *testing.T) {
		pass, _ := enum.validate(43)
		if pass {
			t.Error("wanted error. got nil")
		}
	})
}

func TestCustomRule(t *testing.T) {
	prop := NewProperty("test", Int)
	stringRule := NewCustomRule("StringRule", String, func(i interface{}) (bool, string) {
		val, ok := i.(string)
		if val == "test" && ok {
			return true, ""
		}
		return false, "rule could not be validated."
	})

	intRule := NewCustomRule("IntRule", Int, func(i interface{}) (bool, string) {
		val, ok := i.(int)
		if val < 10 && ok {
			return true, ""
		}
		return false, "could not validate."
	})
	intRule.SetDescription("should be less than 10")

	t.Run("Apply custom rule to prop", func(t *testing.T) {
		err := prop.AddRules(intRule)
		if err != nil {
			t.Errorf("wanted nil got %v", err.Error())
		}
	})

	t.Run("Fail to apply rule", func(t *testing.T) {
		err := prop.AddRules(stringRule)
		if err == nil {
			t.Errorf("wanted error got nil")
		}
	})

	t.Run("pass validation", func(t *testing.T) {
		pass, msg := intRule.validate(1)
		if !pass {
			t.Errorf("test failed int validation did not pass: %v", msg)
		}
	})
}
