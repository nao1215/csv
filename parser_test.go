package csv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/motemen/go-testutil/dataloc"
)

func Test_parseValidateTag(t *testing.T) {
	t.Parallel()
	type args struct {
		tags string
	}
	tests := []struct {
		name string
		args args
		want validators
	}{
		{
			name: "should return a validationRule with all fields set to false",
			args: args{tags: ""},
			want: validators{},
		},
		{
			name: "should return a validationRule with shouldBool set to true",
			args: args{tags: "boolean,alpha"},
			want: validators{
				newBooleanValidator(),
				newAlphaValidator(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseValidateTag(tt.args.tags)
			if err != nil {
				t.Errorf("parseValidateTag() error = %v, test case at %s", err, dataloc.L(tt.name))
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("parseValidateTage() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestCSV_parseStructTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arg     any
		want    ruleSet
		wantErr bool
	}{
		{
			name: "should return an error if the struct is not a pointer",
			arg: &[]struct {
				Name string `validate:"boolean"`
			}{},
			want: ruleSet{
				validators{newBooleanValidator()},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			csv := &CSV{}
			if err := csv.parseStructTag(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("CSV.parseStructTag() error = %v, wantErr %v, test case at %s", err, tt.wantErr, dataloc.L(tt.name))
			}

			if diff := cmp.Diff(csv.ruleSet, tt.want); diff != "" {
				t.Errorf("CSV.parseStructTag() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
