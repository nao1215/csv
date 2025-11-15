package csv

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCSV_Decode(t *testing.T) {
	t.Parallel()
	t.Run("read `id,name,age` header with value", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "sample.csv"))
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Fatalf("failed to close file: %v", err)
			}
		}()

		c, err := NewCSV(f)
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID   int    `validate:"numeric"`
			Name string `validate:"alpha"`
			Age  int    `validate:"numeric"`
		}
		people := make([]person, 0)

		errs := c.Decode(&people)
		if len(errs) != 0 {
			for _, err := range errs {
				t.Errorf("CSV.Decode() got errors: %v", err)
			}
		}

		want := []person{
			{ID: 1, Name: "Gina", Age: 23},
			{ID: 2, Name: "Yulia", Age: 25},
			{ID: 3, Name: "Denis", Age: 30},
		}
		if diff := cmp.Diff(people, want); diff != "" {
			t.Errorf("CSV.Decode() mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("read `id,name,age` header with value and headerless", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "sample_headerless.csv"))
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Fatalf("failed to close file: %v", err)
			}
		}()

		c, err := NewCSV(f, WithHeaderless())
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID   int    `validate:"numeric"`
			Name string `validate:"alpha"`
			Age  int    `validate:"numeric"`
		}
		people := make([]person, 0)

		errs := c.Decode(&people)
		if len(errs) != 0 {
			t.Errorf("CSV.Decode() got errors: %v", errs)
		}

		want := []person{
			{ID: 1, Name: "Gina", Age: 23},
			{ID: 2, Name: "Yulia", Age: 25},
			{ID: 3, Name: "Denis", Age: 30},
		}
		if diff := cmp.Diff(people, want); diff != "" {
			t.Errorf("CSV.Decode() mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("read `id,name,age` header with tab separator", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "sample.tsv"))
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Fatalf("failed to close file: %v", err)
			}
		}()

		c, err := NewCSV(f, WithTabDelimiter())
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID   int    `validate:"numeric"`
			Name string `validate:"alpha"`
			Age  int    `validate:"numeric"`
		}
		people := make([]person, 0)

		errs := c.Decode(&people)
		if len(errs) != 0 {
			t.Errorf("CSV.Decode() got errors: %v", errs)
		}

		want := []person{
			{ID: 1, Name: "Gina", Age: 23},
			{ID: 2, Name: "Yulia", Age: 25},
			{ID: 3, Name: "Denis", Age: 30},
		}
		if diff := cmp.Diff(people, want); diff != "" {
			t.Errorf("CSV.Decode() mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("validate min, max: success case", func(t *testing.T) {
		t.Parallel()

		input := `id,age
1,0
2,1
3,120
4,119
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID  int // no validate
			Age int `validate:"min=0,max=120.0"`
		}

		people := make([]person, 0)
		errs := c.Decode(&people)
		if len(errs) != 0 {
			t.Errorf("CSV.Decode() got errors: %v", errs)
		}

		want := []person{
			{ID: 1, Age: 0},
			{ID: 2, Age: 1},
			{ID: 3, Age: 120},
			{ID: 4, Age: 119},
		}

		if diff := cmp.Diff(people, want); diff != "" {
			t.Errorf("CSV.Decode() mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("validate len: success case", func(t *testing.T) {
		t.Parallel()

		input := `id,name
1,abc
2,あいう
3,👩‍❤‍💋‍👩🇷🇺😂
`
		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID   int    // no validate
			Name string `validate:"len=3"`
		}
		persons := make([]person, 0)

		errs := c.Decode(&persons)
		if len(errs) != 0 {
			t.Errorf("CSV.Decode() got errors: %v", errs)
		}

		want := []person{
			{ID: 1, Name: "abc"},
			{ID: 2, Name: "あいう"},
			{ID: 3, Name: "👩‍❤‍💋‍👩🇷🇺😂"},
		}

		if diff := cmp.Diff(persons, want); diff != "" {
			t.Errorf("CSV.Decode() mismatch (-got +want):\n%s", diff)
		}
	})

	t.Run("validate oneof: success case", func(t *testing.T) {
		t.Parallel()

		input := `id,gender
1,male
2,female
3,prefer_not_to
`
		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID     int    // no validate
			Gender string `validate:"oneof=male female prefer_not_to"`
		}

		people := make([]person, 0)
		errs := c.Decode(&people)
		if len(errs) != 0 {
			t.Errorf("CSV.Decode() got errors: %v", errs)
		}

		want := []person{
			{ID: 1, Gender: "male"},
			{ID: 2, Gender: "female"},
			{ID: 3, Gender: "prefer_not_to"},
		}

		if diff := cmp.Diff(people, want); diff != "" {
			t.Errorf("CSV.Decode() mismatch (-got +want):\n%s", diff)
		}
	})
}

func Test_ErrCheck(t *testing.T) {
	t.Parallel()

	t.Run("error: `id,name,age,password` header", func(t *testing.T) {
		t.Parallel()

		f, err := os.Open(filepath.Join("testdata", "all_error.csv"))
		if err != nil {
			t.Fatal(err)
		}

		c, err := NewCSV(f)
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID       int    `validate:"numeric,gte=1"`
			Name     string `validate:"alpha"`
			Age      int    `validate:"numeric,gt=-1,lt=120,gte=0"`
			Password string `validate:"required,alphanumeric"`
			IsAdmin  bool   `validate:"boolean"`
			Zero     int    `validate:"numeric,eq=0,lte=1,ne=1"`
		}
		people := make([]person, 0)

		got := c.Decode(&people)
		for i, err := range got {
			switch i {
			case 0:
				if err.Error() != "line:2 column id: target is not greater than or equal to the threshold value: threshold=1, value=0" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 1:
				if err.Error() != "line:3 column password: target is not an alphanumeric character: value=password-bad" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 2:
				if err.Error() != "line:4 column password: target is required but is empty: value=" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 3:
				if err.Error() != "line:5 column name: target is not an alphabetic character: value=1Joyless" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 4:
				if err.Error() != "line:5 column zero: target is not equal to the threshold value: threshold=0, value=1" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 5:
				if err.Error() != "line:5 column zero: target is equal to the threshold value: threshold=1, value=1" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 6:
				if err.Error() != "line:6 column age: target is not less than the threshold value: threshold=120, value=120" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 7:
				if err.Error() != "line:7 column is_admin: target is not a boolean: value=2" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 8:
				if err.Error() != "line:8 column age: target is not greater than the threshold value: threshold=-1, value=-1" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 9:
				if err.Error() != "line:8 column age: target is not greater than or equal to the threshold value: threshold=0, value=-1" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 10:
				if err.Error() != "line:9 column id: target is not a numeric character: value=a" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate min, max: error case", func(t *testing.T) {
		t.Parallel()

		input := `id,age
1,0
2,-1
3,120
4,120.1
`
		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID  int // no validate
			Age int `validate:"min=0,max=120.0"`
		}

		people := make([]person, 0)
		errs := c.Decode(&people)

		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:3 column age: target is less than the minimum value: threshold=0, value=-1" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 1:
				if err.Error() != "line:5 column age: target is greater than the maximum value: threshold=120, value=120.1" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate len: error case", func(t *testing.T) {
		t.Parallel()

		input := `id,name
1,abcd
2,あいうえ
3,👩‍❤‍💋‍👩🇷🇺😂🏯
`
		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID   int    // no validate
			Name string `validate:"len=3"`
		}
		persons := make([]person, 0)

		errs := c.Decode(&persons)

		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:2 column name: target length is not equal to the threshold value: length threshold=3, value=abcd" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 1:
				if err.Error() != "line:3 column name: target length is not equal to the threshold value: length threshold=3, value=あいうえ" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 2:
				if err.Error() != "line:4 column name: target length is not equal to the threshold value: length threshold=3, value=👩‍❤‍💋‍👩🇷🇺😂🏯" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate oneof: error case", func(t *testing.T) {
		t.Parallel()

		input := `id,gender
1,smale
2,child
3,prefer_not_tooa
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			ID     int    // no validate
			Gender string `validate:"oneof=male female prefer_not_to"`
		}

		people := make([]person, 0)
		errs := c.Decode(&people)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:2 column gender: target is not one of the values: oneof=male female prefer_not_to, value=smale" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 1:
				if err.Error() != "line:3 column gender: target is not one of the values: oneof=male female prefer_not_to, value=child" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 2:
				if err.Error() != "line:4 column gender: target is not one of the values: oneof=male female prefer_not_to, value=prefer_not_tooa" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate lowercase", func(t *testing.T) {
		t.Parallel()

		input := `name
Abc
abc
ABC
あいう
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			Name string `validate:"lowercase"`
		}

		persons := make([]person, 0)
		errs := c.Decode(&persons)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:2 column name: target is not a lowercase character: value=Abc" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 1:
				if err.Error() != "line:4 column name: target is not a lowercase character: value=ABC" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 2:
				if err.Error() != "line:5 column name: target is not a lowercase character: value=あいう" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate uppercase", func(t *testing.T) {
		t.Parallel()

		input := `name
Abc
abc
ABC
あいう
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type person struct {
			Name string `validate:"uppercase"`
		}

		persons := make([]person, 0)
		errs := c.Decode(&persons)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:2 column name: target is not an uppercase character: value=Abc" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 1:
				if err.Error() != "line:3 column name: target is not an uppercase character: value=abc" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 2:
				if err.Error() != "line:5 column name: target is not an uppercase character: value=あいう" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate ascii", func(t *testing.T) {
		t.Parallel()

		input := fmt.Sprintf(
			"name\n%s\n%s\n",
			"\"!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\`]^_abcdefghijklmnopqrstuvwxyz{|}~\"",
			"あいう",
		)

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}
		type ascii struct {
			Name string `validate:"ascii"`
		}

		asciis := make([]ascii, 0)
		errs := c.Decode(&asciis)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:3 column name: target is not an ASCII character: value=あいう" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate email", func(t *testing.T) {
		t.Parallel()

		input := `email
simple@example.com
very.common@example.com
disposable.style.email.with+symbol@example.com
user.name+tag+sorting@example.com
admin@mailserver1
badあ@example.com
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type email struct {
			Email string `validate:"email"`
		}

		emails := make([]email, 0)
		errs := c.Decode(&emails)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:6 column email: target is not a valid email address: value=admin@mailserver1" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			case 1:
				if err.Error() != "line:7 column email: target is not a valid email address: value=badあ@example.com" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate contains", func(t *testing.T) {
		t.Parallel()

		input := `name
search for a needle in a haystack
example sentence
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type contains struct {
			Name string `validate:"contains=needle"`
		}

		containsList := make([]contains, 0)
		errs := c.Decode(&containsList)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:3 column name: target does not contain the specified value: contains=needle, value=example sentence" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("invalid contains tag format", func(t *testing.T) {
		t.Parallel()

		input := `name
search for a needle in a haystack
example sentence
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type contains struct {
			Name string `validate:"contains=needle bad_value"`
		}

		containsList := make([]contains, 0)
		errs := c.Decode(&containsList)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "'contains' tag format is invalid: contains=needle bad_value" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})

	t.Run("validate containsany", func(t *testing.T) {
		t.Parallel()

		input := `name
you can't find a needle in a haystack
example sentence
I sleep in a bed
`

		c, err := NewCSV(bytes.NewBufferString(input))
		if err != nil {
			t.Fatal(err)
		}

		type containsAny struct {
			Name string `validate:"containsany=needle bed"`
		}

		containsAnyList := make([]containsAny, 0)
		errs := c.Decode(&containsAnyList)
		for i, err := range errs {
			switch i {
			case 0:
				if err.Error() != "line:3 column name: target does not contain any of the specified values: containsany=needle bed, value=example sentence" {
					t.Errorf("CSV.Decode() got errors: %v", err)
				}
			}
		}
	})
}
