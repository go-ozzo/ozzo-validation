package validation

// modeled after string valdiator

type intValidator func(int64) bool

// int rule checks an int variable using a specified intValidator
type IntRule struct {
	validate intValidator
	err      Error
}

// NewIntRule creates a new validation rule using a function that takes a int value and returns a bool.
// The rule returned will use the function to check if a given int is valid or not.
func NewIntRule(validator intValidator, message string) IntRule {
	return IntRule{
		validate: validator,
		err:      NewError("", message),
	}
}

// NewIntRuleWithError creates a new validation rule using a function that takes a int value and returns a bool.
func NewIntRuleWithError(validator intValidator, err Error) IntRule {
	return IntRule{
		validate: validator,
		err:      err,
	}
}

// Error sets the error message for the rule.
func (i IntRule) Error(message string) IntRule {
	i.err = i.err.SetMessage(message)
	return i
}

// ErrorObject sets the error struct for the rule.
func (i IntRule) ErrorObject(err Error) IntRule {
	i.err = err
	return i
}

func (i IntRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil {
		return nil
	}

	intVal, err := ToInt(value)
	if err != nil {
		return err
	}

	if i.validate(intVal) {
		return nil
	}

	return i.err
}
