package csv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/uniseg"
)

// validator is a struct that contains the validation rules for a column.
type validators []validator

// validator is the interface that wraps the Do method.
type validator interface {
	Do(target any) error
}

// booleanValidator is a struct that contains the validation rules for a boolean column.
type booleanValidator struct{}

// newBooleanValidator returns a new booleanValidator.
func newBooleanValidator() *booleanValidator {
	return &booleanValidator{}
}

// Do validates the target as a boolean.
// If the target is an int, it will be validated as a boolean if it's 0 or 1.
func (b *booleanValidator) Do(target any) error {
	if v, ok := target.(string); ok {
		if v == "true" || v == "false" || v == "0" || v == "1" {
			return nil
		}
	}
	return fmt.Errorf("%w: value=%v", ErrInvalidBoolean, target) //nolint
}

// alphabetValidator is a struct that contains the validation rules for an alpha column.
type alphabetValidator struct{}

// newAlphaValidator returns a new alphaValidator.
func newAlphaValidator() *alphabetValidator {
	return &alphabetValidator{}
}

// Do validates the target string only contains alphabetic character.
func (a *alphabetValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrInvalidAlphabet, target) //nolint
	}

	for _, r := range v {
		if !isAlpha(r) {
			return fmt.Errorf("%w: value=%v", ErrInvalidAlphabet, target) //nolint
		}
	}
	return nil
}

// isAlpha returns true if the rune is an alphabetic character.
func isAlpha(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

// numericValidator is a struct that contains the validation rules for a numeric column.
type numericValidator struct{}

// newNumericValidator returns a new numericValidator.
func newNumericValidator() *numericValidator {
	return &numericValidator{}
}

// Do validates the target as a numeric.
func (n *numericValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrInvalidNumeric, target) //nolint
	}

	if v == "" {
		return nil
	}

	if _, err := strconv.Atoi(v); err != nil {
		return fmt.Errorf("%w: value=%s", ErrInvalidNumeric, v) //nolint
	}
	return nil
}

// isNumeric returns true if the rune is a numeric character.
func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

// alphanumericValidator is a struct that contains the validation rules for an alphanumeric column.
type alphanumericValidator struct{}

// newAlphanumericValidator returns a new alphanumericValidator.
func newAlphanumericValidator() *alphanumericValidator {
	return &alphanumericValidator{}
}

// Do validates the target string only contains alphanumeric character.
func (a *alphanumericValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrInvalidAlphanumeric, target) //nolint
	}

	for _, r := range v {
		if !isAlpha(r) && !isNumeric(r) {
			return fmt.Errorf("%w: value=%v", ErrInvalidAlphanumeric, target) //nolint
		}
	}
	return nil
}

// requiredValidator is a struct that contains the validation rules for a required column.
type requiredValidator struct{}

// newRequiredValidator returns a new requiredValidator.
func newRequiredValidator() *requiredValidator {
	return &requiredValidator{}
}

// Do validates the target is not empty.
func (r *requiredValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrRequired, target) //nolint
	}

	if v == "" {
		return fmt.Errorf("%w: value=%v", ErrRequired, target) //nolint
	}
	return nil
}

// equalValidator is a struct that contains the validation rules for an equal column.
type equalValidator struct {
	threshold float64
}

// newEqualValidator returns a new equalValidator.
func newEqualValidator(threshold float64) *equalValidator {
	return &equalValidator{threshold: threshold}
}

// Do validates the target is equal to the threshold.
func (e *equalValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrEqual, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrEqual, target) //nolint
	}

	if value != e.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrEqual, e.threshold, value) //nolint
	}
	return nil
}

// notEqualValidator is a struct that contains the validation rules for a not equal column.
type notEqualValidator struct {
	threshold float64
}

// newNotEqualValidator returns a new notEqualValidator.
func newNotEqualValidator(threshold float64) *notEqualValidator {
	return &notEqualValidator{threshold: threshold}
}

// Do validates the target is not equal to the threshold.
func (n *notEqualValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrNotEqual, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrNotEqual, target) //nolint
	}

	if value == n.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrNotEqual, n.threshold, value) //nolint
	}
	return nil
}

// greaterThanValidator is a struct that contains the validation rules for a greater than column.
type greaterThanValidator struct {
	threshold float64
}

// newGreaterThanValidator returns a new greaterThanValidator.
func newGreaterThanValidator(threshold float64) *greaterThanValidator {
	return &greaterThanValidator{threshold: threshold}
}

// Do validates the target is greater than the threshold.
func (g *greaterThanValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrGreaterThan, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrGreaterThan, target) //nolint
	}

	if value <= g.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrGreaterThan, g.threshold, value) //nolint
	}
	return nil
}

// greaterThanEqualValidator is a struct that contains the validation rules for a greater than or equal column.
type greaterThanEqualValidator struct {
	threshold float64
}

// newGreaterThanEqualValidator returns a new greaterThanEqualValidator.
func newGreaterThanEqualValidator(threshold float64) *greaterThanEqualValidator {
	return &greaterThanEqualValidator{threshold: threshold}
}

// Do validates the target is greater than or equal to the threshold.
func (g *greaterThanEqualValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrGreaterThanEqual, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrGreaterThanEqual, target) //nolint
	}

	if value < g.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrGreaterThanEqual, g.threshold, value) //nolint
	}
	return nil
}

// lessThanValidator is a struct that contains the validation rules for a less than column.
type lessThanValidator struct {
	threshold float64
}

// newLessThanValidator returns a new lessThanValidator.
func newLessThanValidator(threshold float64) *lessThanValidator {
	return &lessThanValidator{threshold: threshold}
}

// Do validates the target is less than the threshold.
func (l *lessThanValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrLessThan, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrLessThan, target) //nolint
	}
	if value >= l.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrLessThan, l.threshold, value) //nolint
	}
	return nil
}

// lessThanEqualValidator is a struct that contains the validation rules for a less than or equal column.
type lessThanEqualValidator struct {
	threshold float64
}

// newLessThanEqualValidator returns a new lessThanEqualValidator.
func newLessThanEqualValidator(threshold float64) *lessThanEqualValidator {
	return &lessThanEqualValidator{threshold: threshold}
}

// Do validates the target is less than or equal to the threshold.
func (l *lessThanEqualValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrLessThanEqual, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrLessThanEqual, target) //nolint
	}

	if value > l.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrLessThanEqual, l.threshold, value) //nolint
	}
	return nil
}

// minValidator is a struct that contains the validation rules for a minimum column.
type minValidator struct {
	threshold float64
}

// newMinValidator returns a new minValidator.
func newMinValidator(threshold float64) *minValidator {
	return &minValidator{threshold: threshold}
}

// Do validates the target is greater than or equal to the threshold.
func (m *minValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrMin, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrMin, target) //nolint
	}

	if value < m.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrMin, m.threshold, value) //nolint
	}
	return nil
}

// maxValidator is a struct that contains the validation rules for a maximum column.
type maxValidator struct {
	threshold float64
}

// newMaxValidator returns a new maxValidator.
func newMaxValidator(threshold float64) *maxValidator {
	return &maxValidator{threshold: threshold}
}

// Do validates the target is less than or equal to the threshold.
func (m *maxValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrMax, target) //nolint
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("%w: value=%v", ErrMax, target) //nolint
	}

	if value > m.threshold {
		return fmt.Errorf("%w: threshold=%f, value=%f", ErrMax, m.threshold, value) //nolint
	}
	return nil
}

// lengthValidator is a struct that contains the validation rules for a length column.
type lengthValidator struct {
	threshold float64
}

// newLengthValidator returns a new lengthValidator.
func newLengthValidator(threshold float64) *lengthValidator {
	return &lengthValidator{threshold: threshold}
}

// Do validates the target length is equal to the threshold.
func (l *lengthValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrLength, target) //nolint
	}

	count := uniseg.GraphemeClusterCount(v)
	if count != int(l.threshold) {
		return fmt.Errorf("%w: length threshold=%d, value=%s", ErrLength, int(l.threshold), v) //nolint
	}
	return nil
}

// oneOfValidator is a struct that contains the validation rules for a one of column.
type oneOfValidator struct {
	oneOf []string
}

// newOneOfValidator returns a new oneOfValidator.
func newOneOfValidator(oneOf []string) *oneOfValidator {
	return &oneOfValidator{oneOf: oneOf}
}

// Do validates the target is one of the oneOf values.
func (o *oneOfValidator) Do(target any) error {
	v, ok := target.(string)
	if !ok {
		return fmt.Errorf("%w: value=%v", ErrOneOf, target) //nolint
	}

	for _, s := range o.oneOf {
		if v == s {
			return nil
		}
	}
	return fmt.Errorf("%w: oneof=%s, value=%v", ErrOneOf, strings.Join(o.oneOf, " "), target) //nolint
}
