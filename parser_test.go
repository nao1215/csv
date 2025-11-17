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
		{
			name: "should return url validator",
			args: args{tags: "url"},
			want: validators{
				newURLValidator(),
			},
		},
		{
			name: "should return uri validator",
			args: args{tags: "uri"},
			want: validators{
				newURIValidator(),
			},
		},
		{
			name: "should return http url validator",
			args: args{tags: "http_url"},
			want: validators{
				newHTTPURLValidator(),
			},
		},
		{
			name: "should return https url validator",
			args: args{tags: "https_url"},
			want: validators{
				newHTTPSURLValidator(),
			},
		},
		{
			name: "should return url encoded validator",
			args: args{tags: "url_encoded"},
			want: validators{
				newURLEncodedValidator(),
			},
		},
		{
			name: "should return ip_addr validator",
			args: args{tags: "ip_addr"},
			want: validators{
				newIPAddrValidator(),
			},
		},
		{
			name: "should return ip4_addr validator",
			args: args{tags: "ip4_addr"},
			want: validators{
				newIP4AddrValidator(),
			},
		},
		{
			name: "should return ip6_addr validator",
			args: args{tags: "ip6_addr"},
			want: validators{
				newIP6AddrValidator(),
			},
		},
		{
			name: "should return uuid validator",
			args: args{tags: "uuid"},
			want: validators{
				newUUIDValidator(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &CSV{}

			got, err := c.parseValidateTag(tt.args.tags)
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
