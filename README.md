<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/csv.svg)](https://pkg.go.dev/github.com/nao1215/csv)
[![MultiPlatformUnitTest](https://github.com/nao1215/csv/actions/workflows/multi_platform_ut.yml/badge.svg)](https://github.com/nao1215/csv/actions/workflows/multi_platform_ut.yml)
[![reviewdog](https://github.com/nao1215/csv/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/csv/actions/workflows/reviewdog.yml)
![Coverage](https://raw.githubusercontent.com/nao1215/octocovs-central-repo/main/badges/nao1215/csv/coverage.svg)

## What is csv package?

The csv package is a library for performing validation when reading CSV or TSV files. Validation rules are specified using struct tags. The csv package read returns which columns of which rows do not adhere to the specified rules.
  
We are implementing internationalization (i18n) for error messages. 

### Supported OS and Go versions

- Linux
- macOS
- Windows
- go version 1.24 or later

### Supported languages

- English
- Japanese
- Russian

If you want to add a new language, please create a pull request.
Ref. https://github.com/nao1215/csv/pull/8

## Why need csv package?

I was frustrated with error-filled CSV files written by non-engineers.

I encountered a use case of "importing one CSV file into multiple DB tables". Unfortunately, I couldn't directly import the CSV file into the DB tables. So, I attempted to import the CSV file through a Go-based application.

What frustrated me was not knowing where the errors in the CSV file were. Existing libraries didn't provide output like "The value in the Mth column of the Nth row is incorrect". I attempted to import multiple times and processed error messages one by one. Eventually, I started writing code to parse each column, which wasted a considerable amount of time.

Based on the above experience, I decided to create a generic CSV validation tool.

## How to use

Please attach the "validate:" tag to your structure and write the validation rules after it. It's crucial that the "order of columns" matches the "order of field definitions" in the structure. The csv package does not automatically adjust the order.

When using csv.Decode, please pass a pointer to a slice of structures tagged with struct tags. The csv package will perform validation based on the struct tags and save the read results to the slice of structures if there are no errors. If there are errors, it will return them as []error.

### Example: english error message
```go
package main

import (
	"bytes"
	"fmt"

	"github.com/nao1215/csv"
)

func main() {
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
```

### Example: japanese error message

```go
func main() {
	input := `id,name,age
1,Gina,23
a,Yulia,25
3,Den1s,30
`
	buf := bytes.NewBufferString(input)
	c, err := csv.NewCSV(buf, csv.WithJapaneseLanguage()) // Set Japanese language option
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
```

### Struct tags

You set the validation rules following the "validate:" tag according to the rules in the table below. If you need to set multiple rules, please enumerate them separated by commas.

#### Strings

| Tag Name          | Description                                       |
|-------------------|---------------------------------------------------|
| alpha             | Check whether value is alphabetic or not           |
| alphanumeric     | Check whether value is alphanumeric or not        |
| ascii             | Check whether value is ASCII or not                |
| boolean           | Check whether value is boolean or not.           |
| contains          | Check whether value contains the specified substring <br> e.g. `validate:"contains=abc"` |
| containsany       | Check whether value contains any of the specified characters <br> e.g. `validate:"containsany=abc def"` |
| lowercase         | Check whether value is lowercase or not           |
| numeric           | Check whether value is numeric or not              |
| uppercase         | Check whether value is uppercase or not           |

#### Format

| Tag Name          | Description                                       |
|-------------------|---------------------------------------------------|
| email             | Check whether value is an email address or not     |

#### Comparisons

| Tag Name          | Description                                       |
|-------------------|---------------------------------------------------|
| eq                | Check whether value is equal to the specified value.<br> e.g. `validate:"eq=1"` |
| gt                | Check whether value is greater than the specified value <br> e.g. `validate:"gt=1"` |
| gte               | Check whether value is greater than or equal to the specified value <br> e.g. `validate:"gte=1"` |
| lt                | Check whether value is less than the specified value <br> e.g. `validate:"lt=1"` |
| lte               | Check whether value is less than or equal to the specified value <br> e.g. `validate:"lte=1"` |
| ne                | Check whether value is not equal to the specified value <br> e.g. `validate:"ne=1"` |

#### Other

| Tag Name          | Description                                       |
|-------------------|---------------------------------------------------|
| len 			    | Check whether the length of the value is equal to the specified value <br> e.g. `validate:"len=10"` |
| max               | Check whether value is less than or equal to the specified value <br> e.g. `validate:"max=100"` |
| min               | Check whether value is greater than or equal to the specified value <br> e.g. `validate:"min=1"` |
| oneof             | Check whether value is included in the specified values <br> e.g. `validate:"oneof=male female prefer_not_to"` |
| required          | Check whether value is empty or not                |

## License
[MIT License](./LICENSE)

## Contribution

First off, thanks for taking the time to contribute! See [CONTRIBUTING.md](./CONTRIBUTING.md) for more information. Contributions are not only related to development. For example, GitHub Star motivates me to develop! Please feel free to contribute to this project.

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/csv&type=Date)](https://star-history.com/#nao1215/csv&Date)

### Special Thanks

I was inspired by the following OSS. Thank you for your great work!
- [go-playground/validator](https://github.com/go-playground/validator)
- [shogo82148/go-header-csv](https://github.com/shogo82148/go-header-csv)

### Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://debimate.jp/"><img src="https://avatars.githubusercontent.com/u/22737008?v=4?s=75" width="75px;" alt="CHIKAMATSU Naohiro"/><br /><sub><b>CHIKAMATSU Naohiro</b></sub></a><br /><a href="https://github.com/nao1215/csv/commits?author=nao1215" title="Documentation">📖</a></td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td align="center" size="13px" colspan="7">
        <img src="https://raw.githubusercontent.com/all-contributors/all-contributors-cli/1b8533af435da9854653492b1327a23a4dbd0a10/assets/logo-small.svg">
          <a href="https://all-contributors.js.org/docs/en/bot/usage">Add your contributions</a>
        </img>
      </td>
    </tr>
  </tfoot>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
