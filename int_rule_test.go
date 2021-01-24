package validation

import (
	"github.com/brianvoe/gofakeit/v5"
	"testing"
)

// mock true return in rule
func testIntValidatorTrue(int int64) bool {
	return true
}

// mock false return in rule
func testIntValidatorFalse(int int64) bool {
	return false
}

var errShouldBeTrue = NewError("test_err_should_be_true", "error should return true")

func TestNewIntRule(t *testing.T) {
	// test raw int rule false
	rule := NewIntRule(testIntValidatorFalse, "this is a test, this should be false")
	err := rule.Validate(gofakeit.Int64())
	if err == nil {
		t.Error("expected error when using testIntValidatorFalse")
	}

	// test raw int rule true
	rule = NewIntRule(testIntValidatorTrue, "this is a test, this should be false")
	rule.ErrorObject(errShouldBeTrue)
	rule.Error("this is a test, this should be false")
	err = rule.Validate(gofakeit.Int64())
	if err != nil {
		t.Error("expected error when using testIntValidatorTrue")
	}

	// test intrule with error
	rule = NewIntRuleWithError(testIntValidatorFalse, errShouldBeTrue)
	err = rule.Validate(gofakeit.Int64())
	if err == nil {
		t.Error("expected error")
	}

	// test wrong type
	rule = NewIntRuleWithError(testIntValidatorTrue, errShouldBeTrue)
	err = rule.Validate(gofakeit.Name())
	if err == nil {
		t.Error("wrong type should result in err")
	}

}
