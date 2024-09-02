package csv

import (
	"errors"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Error is an error that is used to localize error messages.
type Error struct {
	id         string
	subMessage string
	localizer  *i18n.Localizer
}

// Error returns the localized error message.
func (e *Error) Error() string {
	if e.subMessage != "" {
		return fmt.Sprintf(
			"%s: %s",
			e.localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: e.id,
			}),
			e.subMessage,
		)
	}
	return e.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: e.id,
	})
}

// Is reports whether the target error is the same as the error.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.id == t.id
}

// NewError returns a new Error.
func NewError(localizer *i18n.Localizer, id, subMessage string) *Error {
	return &Error{
		id:         id,
		subMessage: subMessage,
		localizer:  localizer,
	}
}

var (
	// ErrStructSlicePointerID is the error ID used when the value is not a pointer to a struct slice.
	ErrStructSlicePointerID = "ErrStructSlicePointer"
	// ErrInvalidOneOfFormatID is the error ID used when the target is not one of the specified values.
	ErrInvalidOneOfFormatID = "ErrInvalidOneOfFormat"
	// ErrInvalidThresholdFormatID is the error ID used when the threshold format is invalid.
	ErrInvalidThresholdFormatID = "ErrInvalidThresholdFormat"
	// ErrInvalidBooleanID is the error ID used when the target is not a boolean.
	ErrInvalidBooleanID = "ErrInvalidBoolean"
	// ErrInvalidAlphabetID is the error ID used when the target is not an alphabetic character.
	ErrInvalidAlphabetID = "ErrInvalidAlphabet"
	// ErrInvalidNumericID is the error ID used when the target is not a numeric character.
	ErrInvalidNumericID = "ErrInvalidNumeric"
	// ErrInvalidAlphanumericID is the error ID used when the target is not an alphanumeric character.
	ErrInvalidAlphanumericID = "ErrInvalidAlphanumeric"
	// ErrRequiredID is the error ID used when the target is required but is empty.
	ErrRequiredID = "ErrRequired"
	// ErrEqualID is the error ID used when the target is not equal to the threshold value.
	ErrEqualID = "ErrEqual"
	// ErrInvalidThresholdID is the error ID used when the threshold value is invalid.
	ErrInvalidThresholdID = "ErrInvalidThreshold"
	// ErrNotEqualID is the error ID used when the target is equal to the threshold value.
	ErrNotEqualID = "ErrNotEqual"
	// ErrGreaterThanID is the error ID used when the target is not greater than the threshold value.
	ErrGreaterThanID = "ErrGreaterThan"
	// ErrGreaterThanEqualID is the error ID used when the target is not greater than or equal to the threshold value.
	ErrGreaterThanEqualID = "ErrGreaterThanEqual"
	// ErrLessThanID is the error ID used when the target is not less than the threshold value.
	ErrLessThanID = "ErrLessThan"
	// ErrLessThanEqualID is the error ID used when the target is not less than or equal to the threshold value.
	ErrLessThanEqualID = "ErrLessThanEqual"
	// ErrMinID is the error ID used when the target is less than the minimum value.
	ErrMinID = "ErrMin"
	// ErrMaxID is the error ID used when the target is greater than the maximum value.
	ErrMaxID = "ErrMax"
	// ErrLengthID is the error ID used when the target length is not equal to the threshold value.
	ErrLengthID = "ErrLength"
	// ErrOneOfID is the error ID used when the target is not one of the specified values.
	ErrOneOfID = "ErrOneOf"
	// ErrInvalidStructID is the error ID used when the target is not a struct.
	ErrInvalidStructID = "ErrInvalidStruct"
	// ErrUnsupportedTypeID is the error ID used when the target is an unsupported type.
	ErrUnsupportedTypeID = "ErrUnsupportedType"
)

var (
	// ErrStructSlicePointer is returned when the value is not a pointer to a struct slice.
	ErrStructSlicePointer = errors.New("value is not a pointer to a struct slice")
	// ErrInvalidOneOfFormat is returned when the target is not one of the values.
	ErrInvalidOneOfFormat = errors.New("target is not one of the values")
	// ErrInvalidThresholdFormat is returned when the threshold value is not an integer.
	ErrInvalidThresholdFormat = errors.New("threshold format is invalid")
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
	// ErrOneOf is returned when the target is not one of the values.
	ErrOneOf = errors.New("target is not one of the values")
)
