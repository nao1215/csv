package csv

import (
	"testing"

	"github.com/motemen/go-testutil/dataloc"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// helperLocalizer is a helper function that returns a new localizer.
func helperLocalizer(t *testing.T) *i18n.Localizer {
	t.Helper()
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "i18n/en.yaml"); err != nil {
		t.Fatalf("load en locale: %v", err)
	}
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "i18n/ja.yaml"); err != nil {
		t.Fatalf("load ja locale: %v", err)
	}
	return i18n.NewLocalizer(bundle, "en")
}

func Test_booleanValidator_Do(t *testing.T) {
	t.Parallel()

	type args struct {
		target any
	}
	tests := []struct {
		name    string
		b       *booleanValidator
		args    args
		wantErr bool
	}{
		{
			name:    "should return nil if target is a boolean: true",
			b:       newBooleanValidator(),
			args:    args{target: "true"},
			wantErr: false,
		},
		{
			name:    "should return nil if target is a boolean: false",
			b:       newBooleanValidator(),
			args:    args{target: "false"},
			wantErr: false,
		},
		{
			name:    "should return nil if target is an int and is 0",
			b:       newBooleanValidator(),
			args:    args{target: "0"},
			wantErr: false,
		},
		{
			name:    "should return nil if target is an int and is 1",
			b:       newBooleanValidator(),
			args:    args{target: "1"},
			wantErr: false,
		},
		{
			name:    "should return an error if target is an int and is not 0 or 1",
			b:       newBooleanValidator(),
			args:    args{target: "2"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := &booleanValidator{}
			if err := b.Do(helperLocalizer(t), tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("booleanValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_alphaValidator_Do(t *testing.T) {
	t.Parallel()

	type args struct {
		target any
	}
	tests := []struct {
		name    string
		a       *alphabetValidator
		args    args
		wantErr bool
	}{
		{
			name:    "should return nil if target is a string and is a multiple alphabetic characters",
			a:       newAlphaValidator(),
			args:    args{target: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"},
			wantErr: false,
		},
		{
			name: "should return nil if target is empty string",
			a:    newAlphaValidator(),
			args: args{target: ""},
		},
		{
			name:    "should return an error if target is not a string",
			a:       newAlphaValidator(),
			args:    args{target: 1},
			wantErr: true,
		},
		{
			name:    "should return an error if target contains number",
			a:       newAlphaValidator(),
			args:    args{target: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1"},
			wantErr: true,
		},
		{
			name:    "should return an error if target contains special character",
			a:       newAlphaValidator(),
			args:    args{target: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &alphabetValidator{}
			if err := a.Do(helperLocalizer(t), tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("alphaValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_alphaSpaceValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		a       *alphaSpaceValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target has alphabets and spaces",
			a:       newAlphaSpaceValidator(),
			arg:     "hello world",
			wantErr: false,
		},
		{
			name:    "should return nil if target is empty string",
			a:       newAlphaSpaceValidator(),
			arg:     "",
			wantErr: false,
		},
		{
			name:    "should return error if contains number",
			a:       newAlphaSpaceValidator(),
			arg:     "hello world1",
			wantErr: true,
		},
		{
			name:    "should return error if not a string",
			a:       newAlphaSpaceValidator(),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("alphaSpaceValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_alphaUnicodeValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		a       *alphaUnicodeValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target is unicode letters",
			a:       newAlphaUnicodeValidator(),
			arg:     "東京ΑΒΓабв",
			wantErr: false,
		},
		{
			name:    "should return error if contains number",
			a:       newAlphaUnicodeValidator(),
			arg:     "東京1",
			wantErr: true,
		},
		{
			name:    "should return error if not a string",
			a:       newAlphaUnicodeValidator(),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.a.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("alphaUnicodeValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_numericValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n       *numericValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target is a string and is a numeric character",
			n:       newNumericValidator(),
			arg:     "1234567890",
			wantErr: false,
		},
		{
			name:    "should return an error if target is not a string",
			n:       newNumericValidator(),
			arg:     1,
			wantErr: true,
		},
		{
			name:    "should return an error if target is not a numeric character",
			n:       newNumericValidator(),
			arg:     "1234567890a",
			wantErr: true,
		},
		{
			name:    "should return an error if target is an empty string",
			n:       newNumericValidator(),
			arg:     "",
			wantErr: false,
		},
		{
			name:    "should return error if target is a string and is a float",
			n:       newNumericValidator(),
			arg:     "0.0",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			n := &numericValidator{}
			if err := n.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("numericValidator.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_numberValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n       *numberValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil for integer",
			n:       newNumberValidator(),
			arg:     "123",
			wantErr: false,
		},
		{
			name:    "should return nil for signed integer",
			n:       newNumberValidator(),
			arg:     "-10",
			wantErr: false,
		},
		{
			name:    "should return nil for decimal",
			n:       newNumberValidator(),
			arg:     "+3.14",
			wantErr: false,
		},
		{
			name:    "should return error for trailing dot",
			n:       newNumberValidator(),
			arg:     "1.",
			wantErr: true,
		},
		{
			name:    "should return error for leading dot",
			n:       newNumberValidator(),
			arg:     ".5",
			wantErr: true,
		},
		{
			name:    "should return error if not string",
			n:       newNumberValidator(),
			arg:     1,
			wantErr: true,
		},
		{
			name:    "should return error if empty string",
			n:       newNumberValidator(),
			arg:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.n.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("numberValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_containsRuneValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		c       *containsRuneValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target contains rune",
			c:       newContainsRuneValidator('界'),
			arg:     "こんにちは世界",
			wantErr: false,
		},
		{
			name:    "should return error if target does not contain rune",
			c:       newContainsRuneValidator('界'),
			arg:     "こんにちは",
			wantErr: true,
		},
		{
			name:    "should return error if target is not string",
			c:       newContainsRuneValidator('界'),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.c.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("containsRuneValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_alphanumericUnicodeValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		a       *alphanumericUnicodeValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target is alphanumeric unicode",
			a:       newAlphanumericUnicodeValidator(),
			arg:     "東京123abc",
			wantErr: false,
		},
		{
			name:    "should return error if target contains symbol",
			a:       newAlphanumericUnicodeValidator(),
			arg:     "東京123abc!",
			wantErr: true,
		},
		{
			name:    "should return error if target is not string",
			a:       newAlphanumericUnicodeValidator(),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.a.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("alphanumericUnicodeValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_startsWithValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       *startsWithValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target starts with prefix",
			s:       newStartsWithValidator("pre"),
			arg:     "prefix-value",
			wantErr: false,
		},
		{
			name:    "should return error if target does not start with prefix",
			s:       newStartsWithValidator("pre"),
			arg:     "value",
			wantErr: true,
		},
		{
			name:    "should return error if target is not a string",
			s:       newStartsWithValidator("pre"),
			arg:     10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.s.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("startsWithValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_endsWithValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		e       *endsWithValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target ends with suffix",
			e:       newEndsWithValidator("fix"),
			arg:     "suffix",
			wantErr: false,
		},
		{
			name:    "should return error if target does not end with suffix",
			e:       newEndsWithValidator("fix"),
			arg:     "value",
			wantErr: true,
		},
		{
			name:    "should return error if target is not a string",
			e:       newEndsWithValidator("fix"),
			arg:     10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.e.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("endsWithValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_endsNotWithValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		e       *endsNotWithValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target does not end with suffix",
			e:       newEndsNotWithValidator("fix"),
			arg:     "value",
			wantErr: false,
		},
		{
			name:    "should return error if target ends with suffix",
			e:       newEndsNotWithValidator("fix"),
			arg:     "suffix",
			wantErr: true,
		},
		{
			name:    "should return error if target is not a string",
			e:       newEndsNotWithValidator("fix"),
			arg:     10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.e.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("endsNotWithValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_startsNotWithValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       *startsNotWithValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target does not start with prefix",
			s:       newStartsNotWithValidator("pre"),
			arg:     "value",
			wantErr: false,
		},
		{
			name:    "should return error if target starts with prefix",
			s:       newStartsNotWithValidator("pre"),
			arg:     "prefix",
			wantErr: true,
		},
		{
			name:    "should return error if target is not string",
			s:       newStartsNotWithValidator("pre"),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.s.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("startsNotWithValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_excludesValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		e       *excludesValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target does not contain excluded substring",
			e:       newExcludesValidator("bad"),
			arg:     "good value",
			wantErr: false,
		},
		{
			name:    "should return error if target contains excluded substring",
			e:       newExcludesValidator("bad"),
			arg:     "this is bad value",
			wantErr: true,
		},
		{
			name:    "should return error if target is not a string",
			e:       newExcludesValidator("bad"),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.e.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("excludesValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_excludesAllValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		e       *excludesAllValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target does not contain any excluded runes",
			e:       newExcludesAllValidator("!@"),
			arg:     "hello world",
			wantErr: false,
		},
		{
			name:    "should return error if target contains excluded rune",
			e:       newExcludesAllValidator("!@"),
			arg:     "hello@world",
			wantErr: true,
		},
		{
			name:    "should return error if target is not a string",
			e:       newExcludesAllValidator("!@"),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.e.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("excludesAllValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_excludesRuneValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		e       *excludesRuneValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target does not contain rune",
			e:       newExcludesRuneValidator('禁'),
			arg:     "許可",
			wantErr: false,
		},
		{
			name:    "should return error if target contains rune",
			e:       newExcludesRuneValidator('禁'),
			arg:     "禁止",
			wantErr: true,
		},
		{
			name:    "should return error if target is not string",
			e:       newExcludesRuneValidator('禁'),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.e.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("excludesRuneValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_multibyteValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		m       *multibyteValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target contains multibyte",
			m:       newMultibyteValidator(),
			arg:     "こんにちは",
			wantErr: false,
		},
		{
			name:    "should return error if target is ASCII only",
			m:       newMultibyteValidator(),
			arg:     "hello",
			wantErr: true,
		},
		{
			name:    "should return error if target is not string",
			m:       newMultibyteValidator(),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.m.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("multibyteValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_cidrValidators_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		v       validator
		arg     any
		wantErr bool
	}{
		{
			name:    "cidr ok ipv4",
			v:       newCIDRValidator(),
			arg:     "192.168.0.0/24",
			wantErr: false,
		},
		{
			name:    "cidr ok ipv6",
			v:       newCIDRValidator(),
			arg:     "2001:db8::/32",
			wantErr: false,
		},
		{
			name:    "cidr invalid string",
			v:       newCIDRValidator(),
			arg:     "invalid",
			wantErr: true,
		},
		{
			name:    "cidrv4 ok",
			v:       newCIDRv4Validator(),
			arg:     "10.0.0.0/8",
			wantErr: false,
		},
		{
			name:    "cidrv4 reject ipv6",
			v:       newCIDRv4Validator(),
			arg:     "2001:db8::/32",
			wantErr: true,
		},
		{
			name:    "cidrv6 ok",
			v:       newCIDRv6Validator(),
			arg:     "2001:db8::/48",
			wantErr: false,
		},
		{
			name:    "cidrv6 reject ipv4",
			v:       newCIDRv6Validator(),
			arg:     "192.168.0.0/24",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.v.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("cidr validation error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}
func Test_printASCIIValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		p       *printASCIIValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target is printable ASCII",
			p:       newPrintASCIIValidator(),
			arg:     "Hello, World! 123",
			wantErr: false,
		},
		{
			name:    "should return error if target contains multibyte",
			p:       newPrintASCIIValidator(),
			arg:     "こんにちは",
			wantErr: true,
		},
		{
			name:    "should return error if target contains control char",
			p:       newPrintASCIIValidator(),
			arg:     "Hello\tWorld",
			wantErr: true,
		},
		{
			name:    "should return error if target is not string",
			p:       newPrintASCIIValidator(),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.p.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("printASCIIValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_urlValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		u       *urlValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target is a valid http url",
			u:       newURLValidator(),
			arg:     "https://example.com/index.html",
			wantErr: false,
		},
		{
			name:    "should return nil if target is a valid file url",
			u:       newURLValidator(),
			arg:     "file:///tmp/data.csv",
			wantErr: false,
		},
		{
			name:    "should return error if target is missing scheme",
			u:       newURLValidator(),
			arg:     "example.com",
			wantErr: true,
		},
		{
			name:    "should return error if target is empty",
			u:       newURLValidator(),
			arg:     "",
			wantErr: true,
		},
		{
			name:    "should return error if target is not a string",
			u:       newURLValidator(),
			arg:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.u.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("urlValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_uriValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		u       *uriValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for valid URI without host",
			u:       newURIValidator(),
			arg:     "custom-scheme:/resource",
			wantErr: false,
		},
		{
			name:    "returns nil for valid URI with host",
			u:       newURIValidator(),
			arg:     "ftp://example.com/files",
			wantErr: false,
		},
		{
			name:    "returns error for empty string",
			u:       newURIValidator(),
			arg:     "",
			wantErr: true,
		},
		{
			name:    "returns error for malformed URI",
			u:       newURIValidator(),
			arg:     "://missing-scheme",
			wantErr: true,
		},
		{
			name:    "returns error for non-string",
			u:       newURIValidator(),
			arg:     10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.u.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("uriValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_httpURLValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		u       *httpURLValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for http url",
			u:       newHTTPURLValidator(),
			arg:     "http://example.com/path",
			wantErr: false,
		},
		{
			name:    "returns nil for https url",
			u:       newHTTPURLValidator(),
			arg:     "https://example.com/path",
			wantErr: false,
		},
		{
			name:    "returns error for missing host",
			u:       newHTTPURLValidator(),
			arg:     "https:///missing-host",
			wantErr: true,
		},
		{
			name:    "returns error for other scheme",
			u:       newHTTPURLValidator(),
			arg:     "file:///tmp/data",
			wantErr: true,
		},
		{
			name:    "returns error for non-string",
			u:       newHTTPURLValidator(),
			arg:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.u.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("httpURLValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_httpsURLValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		u       *httpsURLValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for https url",
			u:       newHTTPSURLValidator(),
			arg:     "https://example.com/path",
			wantErr: false,
		},
		{
			name:    "returns error for http url",
			u:       newHTTPSURLValidator(),
			arg:     "http://example.com/path",
			wantErr: true,
		},
		{
			name:    "returns error for missing host",
			u:       newHTTPSURLValidator(),
			arg:     "https:///missing-host",
			wantErr: true,
		},
		{
			name:    "returns error for non-string",
			u:       newHTTPSURLValidator(),
			arg:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.u.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("httpsURLValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_urlEncodedValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		u       *urlEncodedValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for encoded string",
			u:       newURLEncodedValidator(),
			arg:     "foo%20bar%2Fbaz",
			wantErr: false,
		},
		{
			name:    "returns nil for plain string with no percent",
			u:       newURLEncodedValidator(),
			arg:     "simple-string",
			wantErr: false,
		},
		{
			name:    "returns error for broken escape",
			u:       newURLEncodedValidator(),
			arg:     "bad%2Gvalue",
			wantErr: true,
		},
		{
			name:    "returns error for single percent",
			u:       newURLEncodedValidator(),
			arg:     "bad%value",
			wantErr: true,
		},
		{
			name:    "returns error for non-string",
			u:       newURLEncodedValidator(),
			arg:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.u.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("urlEncodedValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_dataURIValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		u       *dataURIValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for valid data uri",
			u:       newDataURIValidator(),
			arg:     "data:text/plain;base64,SGVsbG8=",
			wantErr: false,
		},
		{
			name:    "returns error for malformed scheme",
			u:       newDataURIValidator(),
			arg:     "text/plain;base64,SGVsbG8=",
			wantErr: true,
		},
		{
			name:    "returns error for invalid base64",
			u:       newDataURIValidator(),
			arg:     "data:text/plain;base64,%%%%",
			wantErr: true,
		},
		{
			name:    "returns error for non-string",
			u:       newDataURIValidator(),
			arg:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.u.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("dataURIValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_hostnameValidators_Do(t *testing.T) {
	t.Parallel()

	hostnameTests := []struct {
		name    string
		v       *hostnameValidator
		arg     any
		wantErr bool
	}{
		{"hostname ok", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), "example.com", false},
		{"hostname starts with digit invalid", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), "1example.com", true},
		{"hostname underscore invalid", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), "exa_mple.com", true},
		{"hostname RFC1123 digit ok", newHostnameValidator(hostnameRFC1123LabelRegexp, ErrHostnameRFC1123ID), "1example.com", false},
		{"hostname RFC1123 trailing dot", newHostnameValidator(hostnameRFC1123LabelRegexp, ErrHostnameRFC1123ID), "example.com.", true},
		{"hostname non-string", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), 123, true},
	}

	for _, tt := range hostnameTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.v.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("hostnameValidator.Do() error = %v, wantErr %v, test %s", err, tt.wantErr, tt.name)
			}
		})
	}

	hostPortTests := []struct {
		name    string
		arg     any
		wantErr bool
	}{
		{"host with port ok", "example.com:80", false},
		{"ipv4 with port ok", "127.0.0.1:8080", false},
		{"ipv6 with port ok", "[2001:db8::1]:443", false},
		{"missing port", "example.com", true},
		{"bad port", "example.com:99999", true},
		{"bad host", "exa_mple.com:80", true},
		{"non-string", 123, true},
	}

	for _, tt := range hostPortTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := newHostnamePortValidator().Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("hostnamePortValidator.Do() error = %v, wantErr %v, test %s", err, tt.wantErr, tt.name)
			}
		})
	}
}

func Test_fqdnValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		f       *fqdnValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for valid fqdn",
			f:       newFQDNValidator(),
			arg:     "example.com",
			wantErr: false,
		},
		{
			name:    "returns nil for subdomain fqdn",
			f:       newFQDNValidator(),
			arg:     "api.dev.example.com",
			wantErr: false,
		},
		{
			name:    "returns error for trailing dot",
			f:       newFQDNValidator(),
			arg:     "example.com.",
			wantErr: true,
		},
		{
			name:    "returns error for single label",
			f:       newFQDNValidator(),
			arg:     "localhost",
			wantErr: true,
		},
		{
			name:    "returns error for invalid chars",
			f:       newFQDNValidator(),
			arg:     "exa_mple.com",
			wantErr: true,
		},
		{
			name:    "returns error for non-string",
			f:       newFQDNValidator(),
			arg:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.f.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("fqdnValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_hostnameValidators_Do(t *testing.T) {
	t.Parallel()

	hostnameTests := []struct {
		name    string
		v       *hostnameValidator
		arg     any
		wantErr bool
	}{
		{"hostname ok", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), "example.com", false},
		{"hostname starts with digit invalid", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), "1example.com", true},
		{"hostname underscore invalid", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), "exa_mple.com", true},
		{"hostname RFC1123 digit ok", newHostnameValidator(hostnameRFC1123LabelRegexp, ErrHostnameRFC1123ID), "1example.com", false},
		{"hostname RFC1123 trailing dot", newHostnameValidator(hostnameRFC1123LabelRegexp, ErrHostnameRFC1123ID), "example.com.", true},
		{"hostname non-string", newHostnameValidator(hostnameRFC952LabelRegexp, ErrHostnameID), 123, true},
	}

	for _, tt := range hostnameTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.v.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("hostnameValidator.Do() error = %v, wantErr %v, test %s", err, tt.wantErr, tt.name)
			}
		})
	}

	hostPortTests := []struct {
		name    string
		arg     any
		wantErr bool
	}{
		{"host with port ok", "example.com:80", false},
		{"ipv4 with port ok", "127.0.0.1:8080", false},
		{"ipv6 with port ok", "[2001:db8::1]:443", false},
		{"missing port", "example.com", true},
		{"bad port", "example.com:99999", true},
		{"bad host", "exa_mple.com:80", true},
		{"non-string", 123, true},
	}

	for _, tt := range hostPortTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := newHostnamePortValidator().Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("hostnamePortValidator.Do() error = %v, wantErr %v, test %s", err, tt.wantErr, tt.name)
			}
		})
	}
}

func Test_ipAddrValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		i       *ipAddrValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for ipv4",
			i:       newIPAddrValidator(),
			arg:     "10.0.0.1",
			wantErr: false,
		},
		{
			name:    "returns nil for ipv6",
			i:       newIPAddrValidator(),
			arg:     "2001:db8::1",
			wantErr: false,
		},
		{
			name:    "returns error for invalid ip",
			i:       newIPAddrValidator(),
			arg:     "999.0.0.1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.i.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("ipAddrValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_ip4AddrValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		i       *ip4AddrValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for ipv4",
			i:       newIP4AddrValidator(),
			arg:     "192.168.1.1",
			wantErr: false,
		},
		{
			name:    "returns error for ipv6",
			i:       newIP4AddrValidator(),
			arg:     "2001:db8::1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.i.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("ip4AddrValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_ip6AddrValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		i       *ip6AddrValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for ipv6",
			i:       newIP6AddrValidator(),
			arg:     "2001:db8::1",
			wantErr: false,
		},
		{
			name:    "returns error for ipv4",
			i:       newIP6AddrValidator(),
			arg:     "192.168.1.1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.i.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("ip6AddrValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_uuidValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		u       *uuidValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "returns nil for valid uuid",
			u:       newUUIDValidator(),
			arg:     "123e4567-e89b-12d3-a456-426614174000",
			wantErr: false,
		},
		{
			name:    "returns error for invalid uuid",
			u:       newUUIDValidator(),
			arg:     "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "returns error for non-string",
			u:       newUUIDValidator(),
			arg:     123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.u.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("uuidValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_equalIgnoreCaseValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		e       *equalIgnoreCaseValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target equals expected ignoring case",
			e:       newEqualIgnoreCaseValidator("Value"),
			arg:     "value",
			wantErr: false,
		},
		{
			name:    "should return error if target does not equal expected ignoring case",
			e:       newEqualIgnoreCaseValidator("Value"),
			arg:     "different",
			wantErr: true,
		},
		{
			name:    "should return error if target is not a string",
			e:       newEqualIgnoreCaseValidator("Value"),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.e.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("equalIgnoreCaseValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}

func Test_notEqualIgnoreCaseValidator_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n       *notEqualIgnoreCaseValidator
		arg     any
		wantErr bool
	}{
		{
			name:    "should return nil if target not equal ignoring case",
			n:       newNotEqualIgnoreCaseValidator("Value"),
			arg:     "other",
			wantErr: false,
		},
		{
			name:    "should return error if target equals ignoring case",
			n:       newNotEqualIgnoreCaseValidator("Value"),
			arg:     "value",
			wantErr: true,
		},
		{
			name:    "should return error if target is not string",
			n:       newNotEqualIgnoreCaseValidator("Value"),
			arg:     1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.n.Do(helperLocalizer(t), tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("notEqualIgnoreCaseValidator.Do() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}
		})
	}
}
