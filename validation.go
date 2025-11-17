package csv

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rivo/uniseg"
)

const fileScheme = "file"

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

// ipValidator validates IPv4 or IPv6 addresses.
type ipValidator struct{}

// newIPValidator returns a new ipValidator.
func newIPValidator() *ipValidator {
	return &ipValidator{}
}

// Do validates the target is a valid IP address.
func (i *ipValidator) Do(localizer *i18n.Localizer, target any) error {
	v, ok := target.(string)
	if !ok {
		return NewError(localizer, ErrIPID, fmt.Sprintf("value=%v", target))
	}

	if v == "" || net.ParseIP(v) == nil {
		return NewError(localizer, ErrIPID, fmt.Sprintf("value=%v", target))
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

var urlEncodedRegexp = regexp.MustCompile(`^(?:[^%]|%[0-9A-Fa-f]{2})*$`)

// urlEncodedValidator validates URL-encoded strings (no invalid % escapes).
type urlEncodedValidator struct{}

// newURLEncodedValidator returns a new urlEncodedValidator.
func newURLEncodedValidator() *urlEncodedValidator {
	return &urlEncodedValidator{}
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
