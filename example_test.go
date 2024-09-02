//go:build linux || darwin

package csv_test

import (
	"bytes"
	"fmt"

	"github.com/nao1215/csv"
)

func ExampleCSV() {
	input := `id,name,age
1,Gina,23
a,Yulia,25
3,Den1s,30
`
	buf := bytes.NewBufferString(input)
	c, err := csv.NewCSV(buf)
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
	// line:2 column age: target is not greater than the threshold value: threshold=24, value=23
	// line:3 column id: target is not a numeric character: value=a
	// line:4 column name: target is not an alphabetic character: value=Den1s
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
	// line:2 column age: ターゲットがしきい値より大きくありません: threshold=24, value=23
	// line:3 column id: ターゲットが数字ではありません: value=a
	// line:4 column name: ターゲットがアルファベット文字ではありません: value=Den1s
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
