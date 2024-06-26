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
	// line:2 column age: target is not greater than the threshold value: threshold=24.000000, value=23.000000
	// line:3 column id: target is not a numeric character: value=a
	// line:4 column name: target is not an alphabetic character: value=Den1s
}
