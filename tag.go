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
	// alphaSpaceTagValue is the struct tag name for alpha with spaces fields.
	alphaSpaceTagValue tagValue = "alphaspace"
	// alphaUnicodeTagValue is the struct tag name for unicode alpha only fields.
	alphaUnicodeTagValue tagValue = "alphaunicode"
	// numericTagValue is the struct tag name for numeric fields.
	numericTagValue tagValue = "numeric"
	// alphanumericTagValue is the struct tag name for alphanumeric fields.
	alphanumericTagValue tagValue = "alphanumeric"
	// alphanumericUnicodeTagValue is the struct tag name for alphanumeric unicode fields.
	alphanumericUnicodeTagValue tagValue = "alphanumunicode"
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
	// equalIgnoreCaseTagValue is the struct tag name for case-insensitive equal fields.
	equalIgnoreCaseTagValue tagValue = "eq_ignore_case"
	// oneOfTagValue is the struct tag name for one of fields.
	oneOfTagValue tagValue = "oneof"
	// lowercaseTagValue is the struct tag name for lowercase fields.
	lowercaseTagValue tagValue = "lowercase"
	// uppercaseTagValue is the struct tag name for uppercase fields.
	uppercaseTagValue tagValue = "uppercase"
	// asciiTagValue is the struct tag name for ascii fields.
	asciiTagValue tagValue = "ascii"
	// notEqualIgnoreCaseTagValue is the struct tag name for case-insensitive not equal fields.
	notEqualIgnoreCaseTagValue tagValue = "ne_ignore_case"
	// numberTagValue is the struct tag name for number fields.
	numberTagValue tagValue = "number"
	// containsRuneTagValue is the struct tag name for contains rune fields.
	containsRuneTagValue tagValue = "containsrune"
	// uriTagValue is the struct tag name for uri fields.
	uriTagValue tagValue = "uri"
	// urlTagValue is the struct tag name for url fields.
	urlTagValue tagValue = "url"
	// httpURLTagValue is the struct tag name for http or https url fields.
	httpURLTagValue tagValue = "http_url"
	// httpsURLTagValue is the struct tag name for https-only url fields.
	httpsURLTagValue tagValue = "https_url"
	// urlEncodedTagValue is the struct tag name for url encoded fields.
	urlEncodedTagValue tagValue = "url_encoded"
	// ipAddrTagValue is the struct tag name for ip_addr fields (IPv4 or IPv6).
	ipAddrTagValue tagValue = "ip_addr"
	// ip4AddrTagValue is the struct tag name for ip4_addr fields (IPv4 only).
	ip4AddrTagValue tagValue = "ip4_addr"
	// ip6AddrTagValue is the struct tag name for ip6_addr fields (IPv6 only).
	ip6AddrTagValue tagValue = "ip6_addr"
	// uuidTagValue is the struct tag name for uuid fields.
	uuidTagValue tagValue = "uuid"
	// emailTagValue is the struct tag name for email fields.
	emailTagValue tagValue = "email"
	// startsWithTagValue is the struct tag name for startswith fields.
	startsWithTagValue tagValue = "startswith"
	// startsNotWithTagValue is the struct tag name for startsnotwith fields.
	startsNotWithTagValue tagValue = "startsnotwith"
	// endsWithTagValue is the struct tag name for endswith fields.
	endsWithTagValue tagValue = "endswith"
	// endsNotWithTagValue is the struct tag name for endsnotwith fields.
	endsNotWithTagValue tagValue = "endsnotwith"
	// excludesTagValue is the struct tag name for excludes fields.
	excludesTagValue tagValue = "excludes"
	// excludesAllTagValue is the struct tag name for excludesall fields.
	excludesAllTagValue tagValue = "excludesall"
	// excludesRuneTagValue is the struct tag name for excludesrune fields.
	excludesRuneTagValue tagValue = "excludesrune"
	// multibyteTagValue is the struct tag name for multibyte fields.
	multibyteTagValue tagValue = "multibyte"
	// printASCIITagValue is the struct tag name for printable ascii fields.
	printASCIITagValue tagValue = "printascii"
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
