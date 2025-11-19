package csv

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rivo/uniseg"
)

const fileScheme = "file"
const uuidRegexPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
const dataURIRegexPattern = `^data:[^;]+;base64,[A-Za-z0-9+/]+={0,2}$`

// validator is a struct that contains the validation rules for a column.
type validators []validator

// validator is the interface that wraps the Do method.
type validator interface {
	Do(localizer *i18n.Localizer, target any) error
}

// booleanValidator is a struct that contains the validation rules for a boolean column.
type booleanValidator struct{}

// newBooleanValidator returns a new booleanValidator.
func newBooleanValidator() *booleanValidator {
	return &booleanValidator{}
}

// Do validates the target as a boolean.
// If the target is an int, it will be validated as a boolean if it's 0 or 1.
func (b *booleanValidator) Do(localizer *i18n.Localizer, target any) error {
	if v, ok := target.(string); ok {
		if v == "true" || v == "false" || v == "0" || v == "1" {
			return nil
		}
	}
	return NewError(localizer, ErrInvalidBooleanID, fmt.Sprintf("value=%v", target))
}

// alphabetValidator is a struct that contains the validation rules for an alpha column.
type alphabetValidator struct{}

// newAlphaValidator returns a new alphaValidator.
func newAlphaValidator() *alphabetValidator {
	return &alphabetValidator{}
}

// Do validates the target string only contains alphabetic character.
func (a *alphabetValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidAlphabetID, fmt.Sprintf("value=%v", target))
	}

	for _, r := range v {
		if !isAlpha(r) {
			return NewError(localizer, ErrInvalidAlphabetID, fmt.Sprintf("value=%v", target))
		}
	}
	return nil
}

// isAlpha returns true if the rune is an alphabetic character.
func isAlpha(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

// alphaUnicodeValidator validates unicode alphabetic characters.
type alphaUnicodeValidator struct{}

// newAlphaUnicodeValidator returns a new alphaUnicodeValidator.
func newAlphaUnicodeValidator() *alphaUnicodeValidator {
	return &alphaUnicodeValidator{}
}

// Do validates the target string only contains unicode letters.
func (a *alphaUnicodeValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidAlphaUnicodeID, fmt.Sprintf("value=%v", target))
	}

	for _, r := range v {
		if !unicode.IsLetter(r) {
			return NewError(localizer, ErrInvalidAlphaUnicodeID, fmt.Sprintf("value=%v", target))
		}
	}
	return nil
}

// alphaSpaceValidator validates strings that contain only alphabetic characters or spaces.
type alphaSpaceValidator struct{}

// newAlphaSpaceValidator returns a new alphaSpaceValidator.
func newAlphaSpaceValidator() *alphaSpaceValidator {
	return &alphaSpaceValidator{}
}

// Do validates the target string only contains alphabetic characters or spaces.
func (a *alphaSpaceValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidAlphaSpaceID, fmt.Sprintf("value=%v", target))
	}

	for _, r := range v {
		if !isAlpha(r) && r != ' ' {
			return NewError(localizer, ErrInvalidAlphaSpaceID, fmt.Sprintf("value=%v", target))
		}
	}
	return nil
}

// numericValidator is a struct that contains the validation rules for a numeric column.
type numericValidator struct{}

// newNumericValidator returns a new numericValidator.
func newNumericValidator() *numericValidator {
	return &numericValidator{}
}

// Do validates the target as a numeric.
func (n *numericValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidNumericID, fmt.Sprintf("value=%v", target))
	}

	if v == "" {
		return nil
	}

	if _, err := strconv.Atoi(v); err != nil {
		return NewError(localizer, ErrInvalidNumericID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// isNumeric returns true if the rune is a numeric character.
func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

// numberValidator is a struct that contains the validation rules for a number column.
// It accepts signed integers and decimals (ASCII), rejecting malformed numbers (e.g., ".5", "5.", "1.2.3").
type numberValidator struct {
	regexp *regexp.Regexp
}

// newNumberValidator returns a new numberValidator.
func newNumberValidator() *numberValidator {
	const numberPattern = `^[-+]?[0-9]+(\.[0-9]+)?$`
	return &numberValidator{
		regexp: regexp.MustCompile(numberPattern),
	}
}

// Do validates the target as a number (integer or decimal with optional sign).
func (n *numberValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidNumberID, fmt.Sprintf("value=%v", target))
	}

	if !n.regexp.MatchString(v) {
		return NewError(localizer, ErrInvalidNumberID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// alphanumericValidator is a struct that contains the validation rules for an alphanumeric column.
type alphanumericValidator struct{}

// newAlphanumericValidator returns a new alphanumericValidator.
func newAlphanumericValidator() *alphanumericValidator {
	return &alphanumericValidator{}
}

// Do validates the target string only contains alphanumeric character.
func (a *alphanumericValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidAlphanumericID, fmt.Sprintf("value=%v", target))
	}

	for _, r := range v {
		if !isAlpha(r) && !isNumeric(r) {
			return NewError(localizer, ErrInvalidAlphanumericID, fmt.Sprintf("value=%v", target))
		}
	}
	return nil
}

// alphanumericUnicodeValidator validates alphanumeric unicode strings.
type alphanumericUnicodeValidator struct{}

// newAlphanumericUnicodeValidator returns a new alphanumericUnicodeValidator.
func newAlphanumericUnicodeValidator() *alphanumericUnicodeValidator {
	return &alphanumericUnicodeValidator{}
}

// Do validates the target string only contains unicode letters or digits.
func (a *alphanumericUnicodeValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidAlphanumericUnicodeID, fmt.Sprintf("value=%v", target))
	}

	for _, r := range v {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return NewError(localizer, ErrInvalidAlphanumericUnicodeID, fmt.Sprintf("value=%v", target))
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
func (r *requiredValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrRequiredID, fmt.Sprintf("value=%v", target))
	}

	if v == "" {
		return NewError(localizer, ErrRequiredID, fmt.Sprintf("value=%v", target))
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
func (e *equalValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrEqualID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrEqualID, fmt.Sprintf("value=%v", target))
	}
	if value != e.threshold {
		return NewError(localizer, ErrEqualID, fmt.Sprintf("threshold=%v, value=%v", e.threshold, value))
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
func (n *notEqualValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrNotEqualID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrNotEqualID, fmt.Sprintf("value=%v", target))
	}

	if value == n.threshold {
		return NewError(localizer, ErrNotEqualID, fmt.Sprintf("threshold=%v, value=%v", n.threshold, value))
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
func (g *greaterThanValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrGreaterThanID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrGreaterThanID, fmt.Sprintf("value=%v", target))
	}

	if value <= g.threshold {
		return NewError(localizer, ErrGreaterThanID, fmt.Sprintf("threshold=%v, value=%v", g.threshold, value))
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
func (g *greaterThanEqualValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrGreaterThanEqualID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrGreaterThanEqualID, fmt.Sprintf("value=%v", target))
	}

	if value < g.threshold {
		return NewError(localizer, ErrGreaterThanEqualID, fmt.Sprintf("threshold=%v, value=%v", g.threshold, value))
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
func (l *lessThanValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrLessThanID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrLessThanID, fmt.Sprintf("value=%v", target))
	}
	if value >= l.threshold {
		return NewError(localizer, ErrLessThanID, fmt.Sprintf("threshold=%v, value=%v", l.threshold, value))
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
func (l *lessThanEqualValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrLessThanEqualID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrLessThanEqualID, fmt.Sprintf("value=%v", target))
	}

	if value > l.threshold {
		return NewError(localizer, ErrLessThanEqualID, fmt.Sprintf("threshold=%v, value=%v", l.threshold, value))
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
func (m *minValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrMinID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrMinID, fmt.Sprintf("value=%v", target))
	}

	if value < m.threshold {
		return NewError(localizer, ErrMinID, fmt.Sprintf("threshold=%v, value=%v", m.threshold, value))
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
func (m *maxValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrMaxID, fmt.Sprintf("value=%v", target))
	}

	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return NewError(localizer, ErrMaxID, fmt.Sprintf("value=%v", target))
	}

	if value > m.threshold {
		return NewError(localizer, ErrMaxID, fmt.Sprintf("threshold=%v, value=%v", m.threshold, value))
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
func (l *lengthValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrLengthID, fmt.Sprintf("value=%v", target))
	}

	count := uniseg.GraphemeClusterCount(v)
	if count != int(l.threshold) {
		return NewError(localizer, ErrLengthID, fmt.Sprintf("length threshold=%v, value=%v", l.threshold, target))
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
func (o *oneOfValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrOneOfID, fmt.Sprintf("value=%v", target))
	}

	for _, s := range o.oneOf {
		if v == s {
			return nil
		}
	}
	return NewError(localizer, ErrOneOfID, fmt.Sprintf("oneof=%s, value=%v", strings.Join(o.oneOf, " "), target))
}

// lowercaseValidator is a struct that contains the validation rules for a lowercase column.
type lowercaseValidator struct{}

// newLowercaseValidator returns a new lowercaseValidator.
func newLowercaseValidator() *lowercaseValidator {
	return &lowercaseValidator{}
}

// Do validates the target is a lowercase string.
func (l *lowercaseValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrLowercaseID, fmt.Sprintf("value=%v", target))
	}

	if v != strings.ToLower(v) {
		return NewError(localizer, ErrLowercaseID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// uppercaseValidator is a struct that contains the validation rules for an uppercase column.
type uppercaseValidator struct{}

// newUppercaseValidator returns a new uppercaseValidator.
func newUppercaseValidator() *uppercaseValidator {
	return &uppercaseValidator{}
}

// Do validates the target is an uppercase string.
func (u *uppercaseValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrUppercaseID, fmt.Sprintf("value=%v", target))
	}

	if v != strings.ToUpper(v) {
		return NewError(localizer, ErrUppercaseID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// asciiValidator is a struct that contains the validation rules for an ASCII column.
type asciiValidator struct{}

// newASCIIValidator returns a new asciiValidator.
func newASCIIValidator() *asciiValidator {
	return &asciiValidator{}
}

// Do validates the target is an ASCII string.
func (a *asciiValidator) Do(localizer *i18n.Localizer, target any) error {
	const maxASCII = 127

	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrASCIIID, fmt.Sprintf("value=%v", target))
	}

	for _, r := range v {
		if r > maxASCII {
			return NewError(localizer, ErrASCIIID, fmt.Sprintf("value=%v", target))
		}
	}
	return nil
}

// uriValidator validates generic URIs (scheme required, host optional).
// It mirrors go-playground/validator's `uri` rule and accepts values like "foo://bar",
// while still rejecting empty strings or malformed request URIs.
type uriValidator struct{}

// newURIValidator returns a new uriValidator.
func newURIValidator() *uriValidator {
	return &uriValidator{}
}

// Do validates the target is a URI.
func (u *uriValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrURIID, fmt.Sprintf("value=%v", target))
	}

	if v == "" {
		return NewError(localizer, ErrURIID, fmt.Sprintf("value=%v", target))
	}

	if i := strings.Index(v, "#"); i > -1 {
		v = v[:i]
	}

	if _, err := url.ParseRequestURI(v); err != nil {
		return NewError(localizer, ErrURIID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// urlValidator is a struct that contains the validation rules for a URL column.
// It requires a scheme and either a host (for http/https/etc.) or a non-empty path for file:,
// matching the simplified URL check from go-playground/validator's `url` rule.
type urlValidator struct{}

// newURLValidator returns a new urlValidator.
func newURLValidator() *urlValidator {
	return &urlValidator{}
}

// Do validates the target is a URL.
func (u *urlValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrURLID, fmt.Sprintf("value=%v", target))
	}

	if v == "" {
		return NewError(localizer, ErrURLID, fmt.Sprintf("value=%v", target))
	}

	parsed, err := url.Parse(strings.ToLower(v))
	if err != nil || parsed.Scheme == "" {
		return NewError(localizer, ErrURLID, fmt.Sprintf("value=%v", target))
	}

	isFileScheme := parsed.Scheme == fileScheme
	if (isFileScheme && (parsed.Path == "" || parsed.Path == "/")) || (!isFileScheme && parsed.Host == "" && parsed.Fragment == "" && parsed.Opaque == "") {
		return NewError(localizer, ErrURLID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// ipAddrValidator accepts both IPv4 and IPv6 (alias of go-playground's `ip_addr`) but uses a dedicated error ID.
type ipAddrValidator struct{}

// newIPAddrValidator returns a new ipAddrValidator.
func newIPAddrValidator() *ipAddrValidator {
	return &ipAddrValidator{}
}

// Do validates the target is a valid IP address (IPv4 or IPv6).
func (i *ipAddrValidator) Do(localizer *i18n.Localizer, target any) error {
	_, err := validateIP(localizer, target, ErrIPAddrID)
	return err
}

// ip4AddrValidator accepts IPv4 addresses only (alias of go-playground's `ip4_addr`).
type ip4AddrValidator struct{}

// newIP4AddrValidator returns a new ip4AddrValidator.
func newIP4AddrValidator() *ip4AddrValidator {
	return &ip4AddrValidator{}
}

// Do validates the target is a valid IPv4 address.
func (i *ip4AddrValidator) Do(localizer *i18n.Localizer, target any) error {
	return validateIPv4(localizer, target, ErrIPv4ID)
}

// ip6AddrValidator accepts IPv6 addresses only (alias of go-playground's `ip6_addr`).
type ip6AddrValidator struct{}

// newIP6AddrValidator returns a new ip6AddrValidator.
func newIP6AddrValidator() *ip6AddrValidator {
	return &ip6AddrValidator{}
}

// Do validates the target is a valid IPv6 address.
func (i *ip6AddrValidator) Do(localizer *i18n.Localizer, target any) error {
	return validateIPv6(localizer, target, ErrIPv6ID)
}

var uuidRegexp = regexp.MustCompile(uuidRegexPattern)

// validateIP validates IPv4 or IPv6 strings and returns an error with the provided error ID.
func validateIP(localizer *i18n.Localizer, target any, errorID string) (net.IP, error) {
	v, ok := target.(string)
	if !ok {
		return nil, NewError(localizer, errorID, fmt.Sprintf("value=%v", target))
	}

	if v == "" {
		return nil, NewError(localizer, errorID, fmt.Sprintf("value=%v", target))
	}

	parsed := net.ParseIP(v)
	if parsed == nil {
		return nil, NewError(localizer, errorID, fmt.Sprintf("value=%v", target))
	}
	return parsed, nil
}

// validateIPv4 validates IPv4 strings.
func validateIPv4(localizer *i18n.Localizer, target any, errorID string) error {
	parsed, err := validateIP(localizer, target, errorID)
	if err != nil {
		return err
	}

	if parsed.To4() == nil {
		return NewError(localizer, errorID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// validateIPv6 validates IPv6 strings.
func validateIPv6(localizer *i18n.Localizer, target any, errorID string) error {
	parsed, err := validateIP(localizer, target, errorID)
	if err != nil {
		return err
	}

	if parsed.To4() != nil {
		return NewError(localizer, errorID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// uuidValidator validates UUID strings (accepts any version).
type uuidValidator struct{}

// newUUIDValidator returns a new uuidValidator.
func newUUIDValidator() *uuidValidator {
	return &uuidValidator{}
}

// Do validates the target is a UUID.
func (u *uuidValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrUUIDID, fmt.Sprintf("value=%v", target))
	}

	if !uuidRegexp.MatchString(v) {
		return NewError(localizer, ErrUUIDID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// httpURLValidator validates http or https URLs with a required host.
type httpURLValidator struct{}

// newHTTPURLValidator returns a new httpURLValidator.
func newHTTPURLValidator() *httpURLValidator {
	return &httpURLValidator{}
}

// Do validates the target is an HTTP or HTTPS URL with host present.
func (u *httpURLValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrHTTPURLID, fmt.Sprintf("value=%v", target))
	}

	parsed, err := url.Parse(strings.ToLower(v))
	if err != nil || parsed.Host == "" {
		return NewError(localizer, ErrHTTPURLID, fmt.Sprintf("value=%v", target))
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return NewError(localizer, ErrHTTPURLID, fmt.Sprintf("value=%v", target))
	}

	// Reuse the general URL shape check for consistency.
	isFileScheme := parsed.Scheme == fileScheme
	if (isFileScheme && (parsed.Path == "" || parsed.Path == "/")) || (!isFileScheme && parsed.Host == "" && parsed.Fragment == "" && parsed.Opaque == "") {
		return NewError(localizer, ErrHTTPURLID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// httpsURLValidator validates https URLs with a required host.
type httpsURLValidator struct{}

// newHTTPSURLValidator returns a new httpsURLValidator.
func newHTTPSURLValidator() *httpsURLValidator {
	return &httpsURLValidator{}
}

// Do validates the target is an HTTPS URL with host present.
func (u *httpsURLValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrHTTPSURLID, fmt.Sprintf("value=%v", target))
	}

	parsed, err := url.Parse(strings.ToLower(v))
	if err != nil || parsed.Host == "" || parsed.Scheme != "https" {
		return NewError(localizer, ErrHTTPSURLID, fmt.Sprintf("value=%v", target))
	}

	isFileScheme := parsed.Scheme == fileScheme
	if (isFileScheme && (parsed.Path == "" || parsed.Path == "/")) || (!isFileScheme && parsed.Host == "" && parsed.Fragment == "" && parsed.Opaque == "") {
		return NewError(localizer, ErrHTTPSURLID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

var (
	urlEncodedRegexp           = regexp.MustCompile(`^(?:[^%]|%[0-9A-Fa-f]{2})*$`)
	dataURIRegex               = regexp.MustCompile(dataURIRegexPattern)
	fqdnLabelRegexp            = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
	hostnameRFC952LabelRegexp  = regexp.MustCompile(`^[A-Za-z](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?$`)
	hostnameRFC1123LabelRegexp = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?$`)
)

// urlEncodedValidator validates URL-encoded strings (no invalid % escapes).
type urlEncodedValidator struct{}

// newURLEncodedValidator returns a new urlEncodedValidator.
func newURLEncodedValidator() *urlEncodedValidator {
	return &urlEncodedValidator{}
}

type dataURIValidator struct{}

// newDataURIValidator returns a new dataURIValidator.
func newDataURIValidator() *dataURIValidator {
	return &dataURIValidator{}
}

type hostnameValidator struct {
	labelRegexp       *regexp.Regexp
	errID             string
	allowLeadingDigit bool
}

// newHostnameValidator returns a new hostnameValidator with the given label regex and error id.
func newHostnameValidator(labelRegexp *regexp.Regexp, errID string) *hostnameValidator {
	return &hostnameValidator{
		labelRegexp:       labelRegexp,
		errID:             errID,
		allowLeadingDigit: true,
	}
}

type hostnamePortValidator struct{}

// newHostnamePortValidator returns a new hostnamePortValidator.
func newHostnamePortValidator() *hostnamePortValidator {
	return &hostnamePortValidator{}
}

// Do validates the target is URL encoded.
func (u *urlEncodedValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrURLEncodedID, fmt.Sprintf("value=%v", target))
	}

	if !urlEncodedRegexp.MatchString(v) {
		return NewError(localizer, ErrURLEncodedID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// Do validates the target is a Data URI with base64 payload.
func (d *dataURIValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrDataURIID, fmt.Sprintf("value=%v", target))
	}

	if !dataURIRegex.MatchString(v) {
		return NewError(localizer, ErrDataURIID, fmt.Sprintf("value=%v", target))
	}

	parts := strings.SplitN(v, ",", 2)
	if len(parts) != 2 {
		return NewError(localizer, ErrDataURIID, fmt.Sprintf("value=%v", target))
	}

	if _, err := base64.StdEncoding.DecodeString(parts[1]); err != nil {
		return NewError(localizer, ErrDataURIID, fmt.Sprintf("value=%v", target))
	}

	return nil
}

// Do validates the target is a hostname according to the provided label regexp.
func (h *hostnameValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, h.errID, fmt.Sprintf("value=%v", target))
	}

	// Reject leading/trailing dot and require at least one dot for FQDN-like hostnames.
	if strings.HasPrefix(v, ".") || strings.HasSuffix(v, ".") {
		return NewError(localizer, h.errID, fmt.Sprintf("value=%v", target))
	}

	labels := strings.Split(v, ".")
	if len(labels) < 1 {
		return NewError(localizer, h.errID, fmt.Sprintf("value=%v", target))
	}

	totalLen := 0
	for _, label := range labels {
		totalLen += len(label) + 1
		if !h.labelRegexp.MatchString(label) {
			return NewError(localizer, h.errID, fmt.Sprintf("value=%v", target))
		}
	}
	if totalLen-1 > 253 {
		return NewError(localizer, h.errID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// Do validates the target is a hostname:port where host is IP or RFC1123 hostname.
func (h *hostnamePortValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrHostnamePortID, fmt.Sprintf("value=%v", target))
	}

	host, portStr, err := net.SplitHostPort(v)
	if err != nil {
		return NewError(localizer, ErrHostnamePortID, fmt.Sprintf("value=%v", target))
	}

	if p, err := strconv.Atoi(portStr); err != nil || p < 1 || p > 65535 {
		return NewError(localizer, ErrHostnamePortID, fmt.Sprintf("value=%v", target))
	}

	// IPv6 with brackets.
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		if ip := net.ParseIP(strings.Trim(host, "[]")); ip != nil {
			return nil
		}
		return NewError(localizer, ErrHostnamePortID, fmt.Sprintf("value=%v", target))
	}

	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	if err := newHostnameValidator(hostnameRFC1123LabelRegexp, ErrHostnamePortID).Do(localizer, host); err != nil {
		return NewError(localizer, ErrHostnamePortID, fmt.Sprintf("value=%v", target))
	}

	return nil
}

// emailValidator is a struct that contains the validation rules for an email column.
type emailValidator struct {
	regexp *regexp.Regexp
}

// newEmailValidator returns a new emailValidator.
func newEmailValidator() *emailValidator {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return &emailValidator{
		regexp: regexp.MustCompile(emailRegexPattern),
	}
}

// Do validates the target is an email.
func (e *emailValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrEmailID, fmt.Sprintf("value=%v", target))
	}

	if !e.regexp.MatchString(v) {
		return NewError(localizer, ErrEmailID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// startsWithValidator is a struct that contains the validation rules for a startswith column.
type startsWithValidator struct {
	prefix string
}

// newStartsWithValidator returns a new startsWithValidator.
func newStartsWithValidator(prefix string) *startsWithValidator {
	return &startsWithValidator{prefix: prefix}
}

// Do validates the target starts with the prefix.
func (s *startsWithValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrStartsWithID, fmt.Sprintf("value=%v", target))
	}

	if !strings.HasPrefix(v, s.prefix) {
		return NewError(localizer, ErrStartsWithID, fmt.Sprintf("startswith=%s, value=%v", s.prefix, target))
	}
	return nil
}

// startsNotWithValidator validates that the target does not start with the prefix.
type startsNotWithValidator struct {
	prefix string
}

// newStartsNotWithValidator returns a new startsNotWithValidator.
func newStartsNotWithValidator(prefix string) *startsNotWithValidator {
	return &startsNotWithValidator{prefix: prefix}
}

// Do validates the target does not start with the prefix.
func (s *startsNotWithValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrStartsNotWithID, fmt.Sprintf("value=%v", target))
	}

	if strings.HasPrefix(v, s.prefix) {
		return NewError(localizer, ErrStartsNotWithID, fmt.Sprintf("startsnotwith=%s, value=%v", s.prefix, target))
	}
	return nil
}

// equalIgnoreCaseValidator validates that two strings are equal, ignoring case.
type equalIgnoreCaseValidator struct {
	expected string
}

// newEqualIgnoreCaseValidator returns a new equalIgnoreCaseValidator.
func newEqualIgnoreCaseValidator(expected string) *equalIgnoreCaseValidator {
	return &equalIgnoreCaseValidator{expected: expected}
}

// Do validates the target matches the expected value ignoring case.
func (e *equalIgnoreCaseValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrEqualIgnoreCaseID, fmt.Sprintf("value=%v", target))
	}

	if !strings.EqualFold(v, e.expected) {
		return NewError(localizer, ErrEqualIgnoreCaseID, fmt.Sprintf("eq_ignore_case=%s, value=%v", e.expected, target))
	}
	return nil
}

// compareValuesEqual reports whether two values are equal for eqfield comparison.
func compareValuesEqual(a, b any) bool {
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		return ok && va == vb
	case int, int8, int16, int32, int64:
		return toInt64(a) == toInt64(b)
	case uint, uint8, uint16, uint32, uint64:
		return toUint64(a) == toUint64(b)
	case float32, float64:
		return toFloat64(a) == toFloat64(b)
	case bool:
		vb, ok := b.(bool)
		return ok && va == vb
	default:
		return false
	}
}

// compareValuesGTE reports whether a >= b for supported types.
// Strings are compared by length to mirror go-playground behavior for gt/gte on strings.
func compareValuesGTE(a, b any) (bool, bool) {
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		if !ok {
			return false, false
		}
		return len(va) >= len(vb), true
	case int, int8, int16, int32, int64:
		return toInt64(a) >= toInt64(b), true
	case uint, uint8, uint16, uint32, uint64:
		return toUint64(a) >= toUint64(b), true
	case float32, float64:
		return toFloat64(a) >= toFloat64(b), true
	default:
		return false, false
	}
}

// compareValuesGT reports whether a > b for supported types.
// Strings are compared by length to mirror go-playground behavior for gt/gte on strings.
func compareValuesGT(a, b any) (bool, bool) {
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		if !ok {
			return false, false
		}
		return len(va) > len(vb), true
	case int, int8, int16, int32, int64:
		return toInt64(a) > toInt64(b), true
	case uint, uint8, uint16, uint32, uint64:
		return toUint64(a) > toUint64(b), true
	case float32, float64:
		return toFloat64(a) > toFloat64(b), true
	default:
		return false, false
	}
}

// compareValuesLTE reports whether a <= b for supported types.
// Strings are compared by length.
func compareValuesLTE(a, b any) (bool, bool) {
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		if !ok {
			return false, false
		}
		return len(va) <= len(vb), true
	case int, int8, int16, int32, int64:
		return toInt64(a) <= toInt64(b), true
	case uint, uint8, uint16, uint32, uint64:
		return toUint64(a) <= toUint64(b), true
	case float32, float64:
		return toFloat64(a) <= toFloat64(b), true
	default:
		return false, false
	}
}

// compareValuesLT reports whether a < b for supported types.
// Strings are compared by length.
func compareValuesLT(a, b any) (bool, bool) {
	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		if !ok {
			return false, false
		}
		return len(va) < len(vb), true
	case int, int8, int16, int32, int64:
		return toInt64(a) < toInt64(b), true
	case uint, uint8, uint16, uint32, uint64:
		return toUint64(a) < toUint64(b), true
	case float32, float64:
		return toFloat64(a) < toFloat64(b), true
	default:
		return false, false
	}
}

// Helpers to normalize numeric types.
func toInt64(v any) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return n
	}
	return 0
}

func toUint64(v any) uint64 {
	switch n := v.(type) {
	case uint:
		return uint64(n)
	case uint8:
		return uint64(n)
	case uint16:
		return uint64(n)
	case uint32:
		return uint64(n)
	case uint64:
		return n
	}
	return 0
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float32:
		return float64(n)
	case float64:
		return n
	}
	return 0
}

// notEqualIgnoreCaseValidator validates that two strings are not equal, ignoring case.
type notEqualIgnoreCaseValidator struct {
	expected string
}

// newNotEqualIgnoreCaseValidator returns a new notEqualIgnoreCaseValidator.
func newNotEqualIgnoreCaseValidator(expected string) *notEqualIgnoreCaseValidator {
	return &notEqualIgnoreCaseValidator{expected: expected}
}

// Do validates the target does not match the expected value ignoring case.
func (n *notEqualIgnoreCaseValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrNotEqualIgnoreCaseID, fmt.Sprintf("value=%v", target))
	}

	if strings.EqualFold(v, n.expected) {
		return NewError(localizer, ErrNotEqualIgnoreCaseID, fmt.Sprintf("ne_ignore_case=%s, value=%v", n.expected, target))
	}
	return nil
}

// endsWithValidator is a struct that contains the validation rules for an endswith column.
type endsWithValidator struct {
	suffix string
}

// newEndsWithValidator returns a new endsWithValidator.
func newEndsWithValidator(suffix string) *endsWithValidator {
	return &endsWithValidator{suffix: suffix}
}

// Do validates the target ends with the suffix.
func (e *endsWithValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrEndsWithID, fmt.Sprintf("value=%v", target))
	}

	if !strings.HasSuffix(v, e.suffix) {
		return NewError(localizer, ErrEndsWithID, fmt.Sprintf("endswith=%s, value=%v", e.suffix, target))
	}
	return nil
}

// endsNotWithValidator is a struct that contains the validation rules for an endsnotwith column.
type endsNotWithValidator struct {
	suffix string
}

// newEndsNotWithValidator returns a new endsNotWithValidator.
func newEndsNotWithValidator(suffix string) *endsNotWithValidator {
	return &endsNotWithValidator{suffix: suffix}
}

// Do validates the target does not end with the suffix.
func (e *endsNotWithValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrEndsNotWithID, fmt.Sprintf("value=%v", target))
	}

	if strings.HasSuffix(v, e.suffix) {
		return NewError(localizer, ErrEndsNotWithID, fmt.Sprintf("endsnotwith=%s, value=%v", e.suffix, target))
	}
	return nil
}

// excludesValidator is a struct that contains the validation rules for an excludes column.
type excludesValidator struct {
	excludes string
}

// newExcludesValidator returns a new excludesValidator.
func newExcludesValidator(excludes string) *excludesValidator {
	return &excludesValidator{excludes: excludes}
}

// Do validates the target does not contain the excluded value.
func (e *excludesValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrExcludesID, fmt.Sprintf("value=%v", target))
	}

	if strings.Contains(v, e.excludes) {
		return NewError(localizer, ErrExcludesID, fmt.Sprintf("excludes=%s, value=%v", e.excludes, target))
	}
	return nil
}

// excludesAllValidator is a struct that contains the validation rules for excludesall column.
// It fails if the target contains any rune from the excludes set (go-playground/validator parity).
type excludesAllValidator struct {
	excludes string
}

// newExcludesAllValidator returns a new excludesAllValidator.
func newExcludesAllValidator(excludes string) *excludesAllValidator {
	return &excludesAllValidator{excludes: excludes}
}

// Do validates the target does not contain any rune from excludes.
func (e *excludesAllValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrExcludesAllID, fmt.Sprintf("value=%v", target))
	}

	if v == "" || e.excludes == "" {
		return nil
	}

	if strings.ContainsAny(v, e.excludes) {
		return NewError(localizer, ErrExcludesAllID, fmt.Sprintf("excludesall=%s, value=%v", e.excludes, target))
	}
	return nil
}

// excludesRuneValidator validates that target does NOT contain the specified single rune.
type excludesRuneValidator struct {
	r rune
}

// newExcludesRuneValidator returns a new excludesRuneValidator.
func newExcludesRuneValidator(r rune) *excludesRuneValidator {
	return &excludesRuneValidator{r: r}
}

// Do validates the target string does not contain the rune.
func (e *excludesRuneValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrExcludesRuneID, fmt.Sprintf("value=%v", target))
	}

	if strings.ContainsRune(v, e.r) {
		return NewError(localizer, ErrExcludesRuneID, fmt.Sprintf("excludesrune=%c, value=%v", e.r, target))
	}
	return nil
}

// multibyteValidator validates that target contains multibyte characters.
type multibyteValidator struct{}

// newMultibyteValidator returns a new multibyteValidator.
func newMultibyteValidator() *multibyteValidator {
	return &multibyteValidator{}
}

// Do validates the target contains at least one multibyte character.
func (m *multibyteValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrMultibyteID, fmt.Sprintf("value=%v", target))
	}

	if utf8.RuneCountInString(v) == len(v) || v == "" {
		return NewError(localizer, ErrMultibyteID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// cidrValidator validates IPv4 or IPv6 CIDR.
type cidrValidator struct{}

// newCIDRValidator returns a new cidrValidator.
func newCIDRValidator() *cidrValidator {
	return &cidrValidator{}
}

// Do validates the target is a valid CIDR.
func (c *cidrValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ip, err := parseCIDR(localizer, target, ErrCIDRID)
	if err != nil {
		return err
	}
	if v == "" || ip == nil {
		return NewError(localizer, ErrCIDRID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// printASCIIValidator validates printable ASCII strings.
type printASCIIValidator struct{}

// cidrv4Validator validates IPv4 CIDR.
type cidrv4Validator struct{}

// newCIDRv4Validator returns a new cidrv4Validator.
func newCIDRv4Validator() *cidrv4Validator {
	return &cidrv4Validator{}
}

// Do validates the target is a valid IPv4 CIDR.
func (c *cidrv4Validator) Do(localizer *i18n.Localizer, target any) error {
	v, ip, err := parseCIDR(localizer, target, ErrCIDRv4ID)
	if err != nil {
		return err
	}
	if ip.To4() == nil {
		return NewError(localizer, ErrCIDRv4ID, fmt.Sprintf("value=%v", target))
	}
	if v == "" {
		return NewError(localizer, ErrCIDRv4ID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// cidrv6Validator validates IPv6 CIDR.
type cidrv6Validator struct{}

// newCIDRv6Validator returns a new cidrv6Validator.
func newCIDRv6Validator() *cidrv6Validator {
	return &cidrv6Validator{}
}

// Do validates the target is a valid IPv6 CIDR.
func (c *cidrv6Validator) Do(localizer *i18n.Localizer, target any) error {
	v, ip, err := parseCIDR(localizer, target, ErrCIDRv6ID)
	if err != nil {
		return err
	}
	if ip.To4() != nil {
		return NewError(localizer, ErrCIDRv6ID, fmt.Sprintf("value=%v", target))
	}
	if v == "" {
		return NewError(localizer, ErrCIDRv6ID, fmt.Sprintf("value=%v", target))
	}
	return nil
}

// newPrintASCIIValidator returns a new printASCIIValidator.
func newPrintASCIIValidator() *printASCIIValidator {
	return &printASCIIValidator{}
}

// Do validates the target contains only printable ASCII characters (0x20-0x7E).
func (p *printASCIIValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrPrintASCIIID, fmt.Sprintf("value=%v", target))
	}

	for _, r := range v {
		if r < 0x20 || r > 0x7e {
			return NewError(localizer, ErrPrintASCIIID, fmt.Sprintf("value=%v", target))
		}
	}
	return nil
}

// parseCIDR parses CIDR string and returns original string and parsed IP.
func parseCIDR(localizer *i18n.Localizer, target any, errorID string) (string, net.IP, error) {
	v, ok := target.(string)
	if !ok {
		return "", nil, NewError(localizer, errorID, fmt.Sprintf("value=%v", target))
	}
	ip, _, err := net.ParseCIDR(v)
	if err != nil {
		return v, nil, NewError(localizer, errorID, fmt.Sprintf("value=%v", target))
	}
	return v, ip, nil
}

// containsValidator is a struct that contains the validation rules for a contains column.
type containsValidator struct {
	contains string
}

// newContainsValidator returns a new containsValidator.
func newContainsValidator(contains string) *containsValidator {
	return &containsValidator{contains: contains}
}

// Do validates the target contains the contains value.
func (c *containsValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrContainsID, fmt.Sprintf("value=%v", target))
	}

	if !strings.Contains(v, c.contains) {
		return NewError(localizer, ErrContainsID, fmt.Sprintf("contains=%s, value=%v", c.contains, target))
	}
	return nil
}

// containsAnyValidator is a struct that contains the validation rules for a contains any column.
type containsAnyValidator struct {
	contains []string
}

// newContainsAnyValidator returns a new containsAnyValidator.
func newContainsAnyValidator(contains []string) *containsAnyValidator {
	return &containsAnyValidator{contains: contains}
}

// Do validates the target contains any of the contains values.
func (c *containsAnyValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrContainsAnyID, fmt.Sprintf("value=%v", target))
	}

	for _, s := range c.contains {
		if strings.Contains(v, s) {
			return nil
		}
	}
	return NewError(localizer, ErrContainsAnyID, fmt.Sprintf("containsany=%s, value=%v", strings.Join(c.contains, " "), target))
}

// containsRuneValidator validates that target contains the specified rune.
type containsRuneValidator struct {
	r rune
}

// newContainsRuneValidator returns a new containsRuneValidator.
func newContainsRuneValidator(r rune) *containsRuneValidator {
	return &containsRuneValidator{r: r}
}

// Do validates the target string contains the rune.
func (c *containsRuneValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrInvalidContainsRuneID, fmt.Sprintf("value=%v", target))
	}

	if !strings.ContainsRune(v, c.r) {
		return NewError(localizer, ErrInvalidContainsRuneID, fmt.Sprintf("containsrune=%c, value=%v", c.r, target))
	}
	return nil
}
