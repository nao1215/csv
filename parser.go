package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// parseStructTag parses the struct tag and extracts the header and ruleSet.
// structSlicePointer is a pointer to a slice of structs.
func (c *CSV) parseStructTag(structSlicePointer any) error {
	rv := reflect.ValueOf(structSlicePointer)
	if rv.Kind() != reflect.Ptr {
		return NewError(c.i18nLocalizer, ErrStructSlicePointerID, "")
	}

	elem := rv.Elem()
	switch elem.Kind() {
	case reflect.Slice, reflect.Array:
		elemType := elem.Type().Elem()
		if elemType.Kind() != reflect.Struct {
			return NewError(c.i18nLocalizer, ErrStructSlicePointerID, "")
		}
		ruleSet, err := c.extractRuleSet(elemType)
		if err != nil {
			return err
		}
		c.ruleSet = ruleSet
	default:
		return NewError(c.i18nLocalizer, ErrStructSlicePointerID, fmt.Sprintf("element=%v", elem.Kind()))
	}
	return nil
}

// / extractRuleSet extracts the ruleSet from the struct.
func (c *CSV) extractRuleSet(structType reflect.Type) (ruleSet, error) {
	ruleSet := make(ruleSet, 0, structType.NumField())

	fieldCount := structType.NumField()
	for i := range fieldCount {
		tag := structType.Field(i).Tag
		validators, err := c.parseValidateTag(tag.Get(validateTag.String()))
		if err != nil {
			return nil, err
		}
		ruleSet = append(ruleSet, validators)
	}
	return ruleSet, nil
}

// parseValidateTag parses the validate tag.
// This function return a set of Validate functions based on
// the rules specified in the validation tag.
func (c *CSV) parseValidateTag(tags string) (validators, error) {
	tagList := strings.Split(tags, ",")
	validatorList := make(validators, 0, len(tagList))

	for _, t := range tagList {
		switch {
		case strings.HasPrefix(t, booleanTagValue.String()):
			validatorList = append(validatorList, newBooleanValidator())
		case strings.HasPrefix(t, alphaSpaceTagValue.String()):
			validatorList = append(validatorList, newAlphaSpaceValidator())
		case strings.HasPrefix(t, alphanumericUnicodeTagValue.String()):
			validatorList = append(validatorList, newAlphanumericUnicodeValidator())
		case strings.HasPrefix(t, alphaTagValue.String()) && !strings.HasPrefix(t, alphanumericTagValue.String()):
			validatorList = append(validatorList, newAlphaValidator())
		case strings.HasPrefix(t, numericTagValue.String()):
			validatorList = append(validatorList, newNumericValidator())
		case strings.HasPrefix(t, alphanumericTagValue.String()):
			validatorList = append(validatorList, newAlphanumericValidator())
		case strings.HasPrefix(t, requiredTagValue.String()):
			validatorList = append(validatorList, newRequiredValidator())
		case strings.HasPrefix(t, equalIgnoreCaseTagValue.String()):
			values, err := c.parseSpecifiedValues(t)
			if err != nil {
				return nil, err
			}
			if len(values) != 1 || values[0] == "" {
				return nil, NewError(c.i18nLocalizer, ErrInvalidEqualIgnoreCaseFormatID, t)
			}
			validatorList = append(validatorList, newEqualIgnoreCaseValidator(values[0]))
		case strings.HasPrefix(t, equalTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newEqualValidator(threshold))
		case strings.HasPrefix(t, notEqualTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newNotEqualValidator(threshold))
		case strings.HasPrefix(t, greaterThanTagValue.String()) && !strings.HasPrefix(t, greaterThanEqualTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newGreaterThanValidator(threshold))
		case strings.HasPrefix(t, greaterThanEqualTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newGreaterThanEqualValidator(threshold))
		case strings.HasPrefix(t, lessThanTagValue.String()) && !strings.HasPrefix(t, lessThanEqualTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newLessThanValidator(threshold))
		case strings.HasPrefix(t, lessThanEqualTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newLessThanEqualValidator(threshold))
		case strings.HasPrefix(t, minTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newMinValidator(threshold))
		case strings.HasPrefix(t, maxTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newMaxValidator(threshold))
		case strings.HasPrefix(t, lengthTagValue.String()):
			threshold, err := c.parseThreshold(t)
			if err != nil {
				return nil, err
			}
			validatorList = append(validatorList, newLengthValidator(threshold))
		case strings.HasPrefix(t, oneOfTagValue.String()):
			oneOf, err := c.parseSpecifiedValues(t)
			if err != nil {
				return nil, NewError(c.i18nLocalizer, ErrInvalidOneOfFormatID, t)
			}
			validatorList = append(validatorList, newOneOfValidator(oneOf))
		case strings.HasPrefix(t, lowercaseTagValue.String()):
			validatorList = append(validatorList, newLowercaseValidator())
		case strings.HasPrefix(t, uppercaseTagValue.String()):
			validatorList = append(validatorList, newUppercaseValidator())
		case strings.HasPrefix(t, asciiTagValue.String()):
			validatorList = append(validatorList, newASCIIValidator())
		case strings.HasPrefix(t, urlEncodedTagValue.String()):
			validatorList = append(validatorList, newURLEncodedValidator())
		case strings.HasPrefix(t, uriTagValue.String()):
			validatorList = append(validatorList, newURIValidator())
		case strings.HasPrefix(t, urlTagValue.String()):
			validatorList = append(validatorList, newURLValidator())
		case strings.HasPrefix(t, httpURLTagValue.String()):
			validatorList = append(validatorList, newHTTPURLValidator())
		case strings.HasPrefix(t, httpsURLTagValue.String()):
			validatorList = append(validatorList, newHTTPSURLValidator())
		case strings.HasPrefix(t, ip4AddrTagValue.String()):
			validatorList = append(validatorList, newIP4AddrValidator())
		case strings.HasPrefix(t, ip6AddrTagValue.String()):
			validatorList = append(validatorList, newIP6AddrValidator())
		case strings.HasPrefix(t, ipAddrTagValue.String()):
			validatorList = append(validatorList, newIPAddrValidator())
		case strings.HasPrefix(t, uuidTagValue.String()):
			validatorList = append(validatorList, newUUIDValidator())
		case strings.HasPrefix(t, emailTagValue.String()):
			validatorList = append(validatorList, newEmailValidator())
		case strings.HasPrefix(t, startsWithTagValue.String()):
			values, err := c.parseSpecifiedValues(t)
			if err != nil {
				return nil, err
			}
			if len(values) != 1 || values[0] == "" {
				return nil, NewError(c.i18nLocalizer, ErrInvalidStartsWithFormatID, t)
			}
			validatorList = append(validatorList, newStartsWithValidator(values[0]))
		case strings.HasPrefix(t, endsWithTagValue.String()):
			values, err := c.parseSpecifiedValues(t)
			if err != nil {
				return nil, err
			}
			if len(values) != 1 || values[0] == "" {
				return nil, NewError(c.i18nLocalizer, ErrInvalidEndsWithFormatID, t)
			}
			validatorList = append(validatorList, newEndsWithValidator(values[0]))
		case strings.HasPrefix(t, containsTagValue.String()) && !strings.HasPrefix(t, containsAnyTagValue.String()):
			values, err := c.parseSpecifiedValues(t)
			if err != nil {
				return nil, err
			}
			if len(values) != 1 {
				return nil, NewError(c.i18nLocalizer, ErrInvalidContainsFormatID, t)
			}
			validatorList = append(validatorList, newContainsValidator(values[0]))
		case strings.HasPrefix(t, containsAnyTagValue.String()):
			values, err := c.parseSpecifiedValues(t)
			if err != nil {
				return nil, err
			}
			if len(values) == 0 {
				return nil, NewError(c.i18nLocalizer, ErrInvalidContainsAnyFormatID, t)
			}
			validatorList = append(validatorList, newContainsAnyValidator(values))
		}
	}
	return validatorList, nil
}

// parseThreshold parses the threshold value.
// tagValue is the value of the struct tag. e.g. eq=10, gt=5.2
func (c *CSV) parseThreshold(tagValue string) (float64, error) {
	parts := strings.Split(tagValue, "=")

	if len(parts) == 2 {
		num, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, NewError(c.i18nLocalizer, ErrInvalidThresholdFormatID, tagValue)
		}
		return num, nil
	}
	return 0, NewError(c.i18nLocalizer, ErrInvalidThresholdFormatID, tagValue)
}

// parseSpecifiedValues parses the tag values.
// tagValue is the value of the struct tag. e.g. oneof=male female prefer_not_to
func (c *CSV) parseSpecifiedValues(tagValue string) ([]string, error) {
	parts := strings.Split(tagValue, "=")

	if len(parts) == 2 {
		return strings.Split(parts[1], " "), nil
	}
	return nil, errors.New("invalid tag values format")
}
