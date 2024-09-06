package csv

// tag is struct tag name.
type tag string

const (
	// validateTag is the struct tag name for validation rules.
	validateTag tag = "validate"
)

// tagValue is the struct tag value.
type tagValue string

const (
	// booleanTagValue is the struct tag value for boolean rule.
	booleanTagValue tagValue = "boolean"
	// alphaTagValue is the struct tag name for alpha only fields.
	alphaTagValue tagValue = "alpha"
	// numericTagValue is the struct tag name for numeric fields.
	numericTagValue tagValue = "numeric"
	// alphanumericTagValue is the struct tag name for alphanumeric fields.
	alphanumericTagValue tagValue = "alphanumeric"
	// requiredTagValue is the struct tag name for required fields.
	requiredTagValue tagValue = "required"
	// equalTagValue is the struct tag name for equal fields.
	equalTagValue tagValue = "eq"
	// notEqualTagValue is the struct tag name for not equal fields.
	notEqualTagValue tagValue = "ne"
	// greaterThanTagValue is the struct tag name for greater than fields.
	greaterThanTagValue tagValue = "gt"
	// greaterThanEqualTagValue is the struct tag name for greater than or equal fields.
	greaterThanEqualTagValue tagValue = "gte"
	// lessThanTagValue is the struct tag name for less than fields.
	lessThanTagValue tagValue = "lt"
	// lessThanEqualTagValue is the struct tag name for less than or equal fields.
	lessThanEqualTagValue tagValue = "lte"
	// minTagValue is the struct tag name for minimum fields.
	minTagValue tagValue = "min"
	// maxTagValue is the struct tag name for maximum fields.
	maxTagValue tagValue = "max"
	// lengthTagValue is the struct tag name for length fields.
	lengthTagValue tagValue = "len"
	// oneOfTagValue is the struct tag name for one of fields.
	oneOfTagValue tagValue = "oneof"
	// lowercaseTagValue is the struct tag name for lowercase fields.
	lowercaseTagValue tagValue = "lowercase"
	// uppercaseTagValue is the struct tag name for uppercase fields.
	uppercaseTagValue tagValue = "uppercase"
	// asciiTagValue is the struct tag name for ascii fields.
	asciiTagValue tagValue = "ascii"
	// emailTagValue is the struct tag name for email fields.
	emailTagValue tagValue = "email"
	// containsTagValue is the struct tag name for contains fields.
	containsTagValue tagValue = "contains"
	// containsAnyTagValue is the struct tag name for contains any fields.
	containsAnyTagValue tagValue = "containsany"
)

// String returns the string representation of the tag.
func (t tag) String() string {
	return string(t)
}

// String returns the string representation of the tag value.
func (t tagValue) String() string {
	return string(t)
}
