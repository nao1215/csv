package csv

// crossFieldRule represents a validation rule that requires access to multiple fields
// in the same struct (same CSV row). This is intentionally limited to flat structs.
type crossFieldRule struct {
	// fieldIndex is the index of the current field (column) this rule applies to.
	fieldIndex int
	// targetField is the name of the field to compare with (struct field name).
	targetField string
	// op specifies the comparison operator (eq, etc.).
	op crossFieldOp
}

// crossFieldOp enumerates supported cross-field operations (same row only; flat structs).
type crossFieldOp string

const (
	crossFieldOpEqual    crossFieldOp = "eqfield"
	crossFieldOpNotEqual crossFieldOp = "nefield"
	crossFieldOpContains crossFieldOp = "fieldcontains"
	crossFieldOpExcludes crossFieldOp = "fieldexcludes"
	crossFieldOpGte      crossFieldOp = "gtefield"
	crossFieldOpGt       crossFieldOp = "gtfield"
	crossFieldOpLte      crossFieldOp = "ltefield"
	crossFieldOpLt       crossFieldOp = "ltfield"
)

// crossFieldRuleSet is a slice of crossFieldRule.
// crossFieldRuleSet is per-field slice of crossFieldRule.
type crossFieldRuleSet [][]crossFieldRule
