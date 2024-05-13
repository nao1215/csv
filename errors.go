package csv

import "errors"

var (
	// ErrStructSlicePointer is returned when the value is not a pointer to a struct slice.
	ErrStructSlicePointer = errors.New("value is not a pointer to a struct slice")
	// ErrInvalidBoolean is returned when the target is not a boolean.
	ErrInvalidBoolean = errors.New("target is not a boolean")
	// ErrInvalidAlphabet is returned when the target is not an alphabetic character.
	ErrInvalidAlphabet = errors.New("target is not an alphabetic character")
	// ErrInvalidNumeric is returned when the target is not a numeric character.
	ErrInvalidNumeric = errors.New("target is not a numeric character")
	// ErrInvalidAlphanumeric is returned when the target is not an alphanumeric character.
	ErrInvalidAlphanumeric = errors.New("target is not an alphanumeric character")
	// ErrRequired is returned when the target is required but is empty.
	ErrRequired = errors.New("target is required but is empty")
	// ErrEqual is returned when the target is not equal to the value.
	ErrEqual = errors.New("target is not equal to the threshold value")
	// ErrInvalidThreshold is returned when the target is not greater than the value.
	ErrInvalidThreshold = errors.New("threshold value is invalid")
	// ErrInvalidThresholdFormat is returned when the threshold value is not an integer.
	ErrInvalidThresholdFormat = errors.New("threshold format is invalid")
	// ErrNotEqual is returned when the target is equal to the value.
	ErrNotEqual = errors.New("target is equal to threshold the value")
	// ErrGreaterThan is returned when the target is not greater than the value.
	ErrGreaterThan = errors.New("target is not greater than the threshold value")
	// ErrGreaterThanEqual is returned when the target is not greater than or equal to the value.
	ErrGreaterThanEqual = errors.New("target is not greater than or equal to the threshold value")
	// ErrLessThan is returned when the target is not less than the value.
	ErrLessThan = errors.New("target is not less than the threshold value")
	// ErrLessThanEqual is returned when the target is not less than or equal to the value.
	ErrLessThanEqual = errors.New("target is not less than or equal to the threshold value")
	// ErrMin is returned when the target is less than the minimum value.
	ErrMin = errors.New("target is less than the minimum value")
	// ErrMax is returned when the target is greater than the maximum value.
	ErrMax = errors.New("target is greater than the maximum value")
	// ErrLength is returned when the target length is not equal to the value.
	ErrLength = errors.New("target length is not equal to the threshold value")
)
