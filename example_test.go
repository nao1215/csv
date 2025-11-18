//go:build linux || darwin

package csv_test

import (
	"bytes"
	"fmt"

	"github.com/nao1215/csv"
)

func ExampleCSV() {
	input := `id,name,age,password,password_confirm,role,note,nickname,ip,cidr,url
1,Alice,17,Secret123,Secret12,superuser,"TODO: fix",alice!,999.0.0.1,10.0.0.0/33,http://example.com
-5,Bob,30,short,short,admin,"Note: ready",Bob123,192.168.0.1,192.168.0.0/24,https://example.com
`
	buf := bytes.NewBufferString(input)
	c, err := csv.NewCSV(buf)
	if err != nil {
		panic(err)
	}

	type account struct {
		ID              int    `validate:"number,gte=1"`
		Name            string `validate:"alpha"`
		Age             int    `validate:"number,gte=18,lte=65"`
		Password        string `validate:"required,gte=8"`
		PasswordConfirm string `validate:"eqfield=Password"`
		Role            string `validate:"oneof=admin user"`
		Note            string `validate:"excludes=TODO,startswith=Note"`
		Nickname        string `validate:"alphanumunicode"`
		IP              string `validate:"ip4_addr"`
		CIDR            string `validate:"cidrv4"`
		URL             string `validate:"https_url"`
	}
	accounts := make([]account, 0)

	errs := c.Decode(&accounts)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	// Output:
	// line:2 column age: target is not greater than or equal to the threshold value: threshold=18, value=17
	// line:2 column password: target is not greater than or equal to the threshold value: value=Secret123
	// line:2 column role: target is not one of the values: oneof=admin user, value=superuser
	// line:2 column note: target contains a prohibited substring: excludes=TODO, value=TODO: fix
	// line:2 column note: target does not start with the specified value: startswith=Note, value=TODO: fix
	// line:2 column nickname: target is not an alphanumeric unicode character: value=alice!
	// line:2 column ip: target is not a valid IPv4 address: value=999.0.0.1
	// line:2 column cidr: target is not a valid IPv4 CIDR: value=10.0.0.0/33
	// line:2 column url: target is not a valid HTTPS URL: value=http://example.com
	// line:2 column password_confirm: field is not equal to the specified field: field=PasswordConfirm, other=Password
	// line:3 column id: target is not greater than or equal to the threshold value: threshold=1, value=-5
	// line:3 column password: target is not greater than or equal to the threshold value: value=short
}

func ExampleWithJapaneseLanguage() {
	input := `id,name,age
1,Gina,23
a,Yulia,25
3,Den1s,30
`
	buf := bytes.NewBufferString(input)
	c, err := csv.NewCSV(buf, csv.WithJapaneseLanguage())
	if err != nil {
		panic(err)
	}

	type person struct {
		ID   int    `validate:"numeric"`
		Name string `validate:"alpha"`
		Age  int    `validate:"gt=24"`
	}
	people := make([]person, 0)

	errs := c.Decode(&people)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	// Output:
	// line:2 column age: 値がしきい値より大きくありません: threshold=24, value=23
	// line:3 column id: 値が数字ではありません: value=a
	// line:4 column name: 値がアルファベット文字ではありません: value=Den1s
}

func ExampleWithRussianLanguage() {
	input := `id,name,age
1,Gina,23
a,Yulia,25
3,Den1s,30
`
	buf := bytes.NewBufferString(input)
	c, err := csv.NewCSV(buf, csv.WithRussianLanguage())
	if err != nil {
		panic(err)
	}

	type person struct {
		ID   int    `validate:"numeric"`
		Name string `validate:"alpha"`
		Age  int    `validate:"gt=24"`
	}
	people := make([]person, 0)

	errs := c.Decode(&people)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	// Output:
	// line:2 column age: целевое значение не больше порогового значения: threshold=24, value=23
	// line:3 column id: целевое значение не является числовым символом: value=a
	// line:4 column name: целевое значение не является алфавитным символом: value=Den1s
}
