package csv

import (
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
	// ErrInvalidAlphaSpaceID is the error ID used when the target is not an alphabetic character or space.
	ErrInvalidAlphaSpaceID = "ErrInvalidAlphaSpace"
	// ErrInvalidAlphaUnicodeID is the error ID used when the target is not a unicode alphabetic character.
	ErrInvalidAlphaUnicodeID = "ErrInvalidAlphaUnicode"
	// ErrInvalidNumericID is the error ID used when the target is not a numeric character.
	ErrInvalidNumericID = "ErrInvalidNumeric"
	// ErrInvalidAlphanumericID is the error ID used when the target is not an alphanumeric character.
	ErrInvalidAlphanumericID = "ErrInvalidAlphanumeric"
	// ErrInvalidAlphanumericUnicodeID is the error ID used when the target is not an alphanumeric unicode character.
	ErrInvalidAlphanumericUnicodeID = "ErrInvalidAlphanumericUnicode"
	// ErrInvalidContainsRuneID is the error ID used when the target does not contain the specified rune.
	ErrInvalidContainsRuneID = "ErrInvalidContainsRune"
	// ErrInvalidContainsRuneFormatID is the error ID used when the containsrune format is invalid.
	ErrInvalidContainsRuneFormatID = "ErrInvalidContainsRuneFormat"
	// ErrRequiredID is the error ID used when the target is required but is empty.
	ErrRequiredID = "ErrRequired"
	// ErrEqualID is the error ID used when the target is not equal to the threshold value.
	ErrEqualID = "ErrEqual"
	// ErrEqualIgnoreCaseID is the error ID used when the target is not equal to the specified value ignoring case.
	ErrEqualIgnoreCaseID = "ErrEqualIgnoreCase"
	// ErrInvalidThresholdID is the error ID used when the threshold value is invalid.
	ErrInvalidThresholdID = "ErrInvalidThreshold"
	// ErrNotEqualID is the error ID used when the target is equal to the threshold value.
	ErrNotEqualID = "ErrNotEqual"
	// ErrInvalidEqualIgnoreCaseFormatID is the error ID used when the eq_ignore_case format is invalid.
	ErrInvalidEqualIgnoreCaseFormatID = "ErrInvalidEqualIgnoreCaseFormat"
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
	// ErrLowercaseID is the error ID used when the target is not a lowercase character.
	ErrLowercaseID = "ErrLowercase"
	// ErrUppercaseID is the error ID used when the target is not an uppercase character.
	ErrUppercaseID = "ErrUppercase"
	// ErrASCIIID is the error ID used when the target is not an ASCII character.
	ErrASCIIID = "ErrASCII"
	// ErrURIID is the error ID used when the target is not a URI.
	ErrURIID = "ErrURI"
	// ErrURLID is the error ID used when the target is not a URL.
	ErrURLID = "ErrURL"
	// ErrHTTPURLID is the error ID used when the target is not an HTTP or HTTPS URL.
	ErrHTTPURLID = "ErrHTTPURL"
	// ErrHTTPSURLID is the error ID used when the target is not an HTTPS URL.
	ErrHTTPSURLID = "ErrHTTPSURL"
	// ErrURLEncodedID is the error ID used when the target is not URL encoded.
	ErrURLEncodedID = "ErrURLEncoded"
	// ErrIPAddrID is the error ID used when the target is not an IP address (ip_addr).
	ErrIPAddrID = "ErrIPAddr"
	// ErrIPv4ID is the error ID used when the target is not an IPv4 address.
	ErrIPv4ID = "ErrIPv4"
	// ErrIPv6ID is the error ID used when the target is not an IPv6 address.
	ErrIPv6ID = "ErrIPv6"
	// ErrUUIDID is the error ID used when the target is not a UUID.
	ErrUUIDID = "ErrUUID"
	// ErrEmailID is the error ID used when the target is not an email.
	ErrEmailID = "ErrEmail"
	// ErrStartsWithID is the error ID used when the target does not start with the specified value.
	ErrStartsWithID = "ErrStartsWith"
	// ErrInvalidStartsWithFormatID is the error ID used when the startswith format is invalid.
	ErrInvalidStartsWithFormatID = "ErrInvalidStartsWithFormat"
	// ErrEndsWithID is the error ID used when the target does not end with the specified value.
	ErrEndsWithID = "ErrEndsWith"
	// ErrInvalidEndsWithFormatID is the error ID used when the endswith format is invalid.
	ErrInvalidEndsWithFormatID = "ErrInvalidEndsWithFormat"
	// ErrContainsID is the error ID used when the target does not contain the specified value.
	ErrContainsID = "ErrContains"
	// ErrInvalidContainsFormatID is the error ID used when the contains format is invalid.
	ErrInvalidContainsFormatID = "ErrInvalidContainsFormat"
	// ErrContainsAnyID is the error ID used when the target does not contain any of the specified values.
	ErrContainsAnyID = "ErrContainsAny"
	// ErrInvalidContainsAnyFormatID is the error ID used when the contains any format is invalid.
	ErrInvalidContainsAnyFormatID = "ErrInvalidContainsAnyFormat"
)
