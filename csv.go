package csv

import (
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// LocaleFS embeds translation files located in the i18n directory.
//
//go:embed i18n/*
var LocaleFS embed.FS

// CSV is a struct that implements CSV Reader and Writer.
type CSV struct {
	// headerless is a flag that indicates the csv file has no header.
	headerless bool
	// reader is the csv reader.
	reader *csv.Reader
	// header is a type that represents the header of a csv.
	header header
	// ruleSets is slice of ruleSet.
	// The order of the ruleSet is the same as the order of the columns in the csv.
	ruleSet ruleSet
	// crossFieldRules holds validation rules that require multiple fields within the same row.
	// The index matches the column index (field index) for which the rule applies.
	crossFieldRules crossFieldRuleSet
	// i18nBundle is the i18n bundle. It is used to translate error messages.
	// The default language is English.
	i18nBundle *i18n.Bundle
	// i18nLocalizer is the i18n localizer. It is used to localize error messages.
	// The default language is English.
	i18nLocalizer *i18n.Localizer
}

type (
	// header is a type that represents the header of a CSV file.
	header []column
	// column is a type that represents a column in a CSV file.
	column string
	// ruleSet is a map that contains the validation rules for each column.
	ruleSet []validators
)

// NewCSV returns a new CSV struct.
func NewCSV(r io.Reader, opts ...Option) (*CSV, error) {
	csv := &CSV{
		reader: csv.NewReader(r),
	}

	if err := csv.newI18n(); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(csv); err != nil {
			return nil, err
		}
	}
	return csv, nil
}

// newI18n initializes the i18n bundle and localizer.
func (c *CSV) newI18n() error {
	c.i18nBundle = i18n.NewBundle(language.English)
	c.i18nBundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	for _, lang := range []string{"en", "ja", "ru"} {
		if _, err := c.i18nBundle.LoadMessageFileFS(LocaleFS, fmt.Sprintf("i18n/%s.yaml", lang)); err != nil {
			return NewError(c.i18nLocalizer, "ErrLoadMessageFile", err.Error())
		}
	}
	c.i18nLocalizer = i18n.NewLocalizer(c.i18nBundle, "en")
	return nil
}

// Decode reads the CSV and returns the columns that have syntax errors on a per-line basis.
// The strutSlicePointer is a pointer to structure slice where validation rules are set in struct tags.
func (c *CSV) Decode(structSlicePointer any) []error {
	errors := make([]error, 0)
	if err := c.parseStructTag(structSlicePointer); err != nil {
		errors = append(errors, err)
		return errors
	}

	firstLine := 1
	if !c.headerless {
		firstLine = 2 // first line is 2 because the header is on line 1.
		if err := c.readHeader(); err != nil {
			errors = append(errors, err)
			return errors
		}
	}

	structSlicePtrValue := reflect.ValueOf(structSlicePointer)
	structSliceValue := structSlicePtrValue.Elem()

	for line := firstLine; ; line++ {
		record, err := c.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, err)
			break
		}

		structValue := reflect.New(structSliceValue.Type().Elem()).Elem()
		for i, v := range record {
			validators := c.ruleSet[i]
			for _, validator := range validators {
				if err := validator.Do(c.i18nLocalizer, v); err != nil {
					errors = append(errors, fmt.Errorf("line:%d column %s: %w", line, c.header[i], err))
				}
			}
			_ = setStructFieldValue(structValue, i, v) //nolint:errcheck // user will not see this error.
		}
		c.validateCrossFieldRules(structValue, line, &errors)
		structSliceValue.Set(reflect.Append(structSliceValue, structValue))
	}
	return errors
}

// validateCrossFieldRules runs cross-field validations for a single struct value (one CSV row).
// It assumes flat structs and compares to other fields in the same struct.
func (c *CSV) validateCrossFieldRules(structValue reflect.Value, line int, errors *[]error) {
	for fieldIdx, rules := range c.crossFieldRules {
		if len(rules) == 0 {
			continue
		}
		srcField := structValue.Field(fieldIdx)
		srcName := structValue.Type().Field(fieldIdx).Name
		colName := srcName
		if fieldIdx < len(c.header) {
			colName = string(c.header[fieldIdx])
		}

		for _, rule := range rules {
			targetField := structValue.FieldByName(rule.targetField)
			if !targetField.IsValid() {
				err := NewError(c.i18nLocalizer, ErrInvalidEqualFieldFormatID, rule.targetField)
				*errors = append(*errors, fmt.Errorf("line:%d column %s: %w", line, colName, err))
				continue
			}

			switch rule.op {
			case crossFieldOpEqual:
				if !compareValuesEqual(srcField.Interface(), targetField.Interface()) {
					err := NewError(
						c.i18nLocalizer,
						ErrEqualFieldID,
						fmt.Sprintf("field=%s, other=%s", srcName, rule.targetField),
					)
					*errors = append(*errors, fmt.Errorf("line:%d column %s: %w", line, colName, err))
				}
			case crossFieldOpContains:
				srcStr, okSrc := srcField.Interface().(string)
				targetStr, okTgt := targetField.Interface().(string)
				if !okSrc || !okTgt {
					err := NewError(c.i18nLocalizer, ErrFieldContainsID, fmt.Sprintf("value=%v", srcField.Interface()))
					*errors = append(*errors, fmt.Errorf("line:%d column %s: %w", line, colName, err))
					continue
				}
				if !strings.Contains(srcStr, targetStr) {
					err := NewError(
						c.i18nLocalizer,
						ErrFieldContainsID,
						fmt.Sprintf("field=%s, other=%s", srcName, rule.targetField),
					)
					*errors = append(*errors, fmt.Errorf("line:%d column %s: %w", line, colName, err))
				}
			case crossFieldOpExcludes:
				srcStr, okSrc := srcField.Interface().(string)
				targetStr, okTgt := targetField.Interface().(string)
				if !okSrc || !okTgt {
					err := NewError(c.i18nLocalizer, ErrFieldExcludesID, fmt.Sprintf("value=%v", srcField.Interface()))
					*errors = append(*errors, fmt.Errorf("line:%d column %s: %w", line, colName, err))
					continue
				}
				if strings.Contains(srcStr, targetStr) {
					err := NewError(
						c.i18nLocalizer,
						ErrFieldExcludesID,
						fmt.Sprintf("field=%s, other=%s", srcName, rule.targetField),
					)
					*errors = append(*errors, fmt.Errorf("line:%d column %s: %w", line, colName, err))
				}
			}
		}
	}
}

// readHeader reads the header of the CSV file.
func (c *CSV) readHeader() error {
	record, err := c.reader.Read()
	if err != nil {
		return err
	}

	columns := make([]column, 0, len(record))
	for _, v := range record {
		columns = append(columns, column(v))
	}
	c.header = columns
	return nil
}

// setStructFieldValue sets the value of a field in a struct.
func setStructFieldValue(structValue reflect.Value, index int, value string) error {
	if index >= structValue.NumField() {
		return errors.New("index out of range for struct")
	}

	fieldValue := structValue.Field(index)
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, fieldValue.Type().Bits())
		if err != nil {
			return err
		}
		fieldValue.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, fieldValue.Type().Bits())
		if err != nil {
			return err
		}
		fieldValue.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, fieldValue.Type().Bits())
		if err != nil {
			return err
		}
		fieldValue.SetFloat(floatValue)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind().String())
	}
	return nil
}
