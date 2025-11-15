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
