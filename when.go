package validation

// When returns a validation rule that checks if the specified condition
//is true, validate the value by the specified rules.
func When(condition bool, rules ...Rule) WhenRule {
	return WhenRule{
		condition: condition,
		rules:     rules,
	}
}

// WhenRule is a validation rule that validate element if condition is true.
type WhenRule struct {
	condition bool
	rules     []Rule
}

// Validate checks if the condition is true, validate value by specified rules.
func (r WhenRule) Validate(value interface{}) error {
	if r.condition {
		return Validate(value, r.rules...)
	}

	return nil
}
