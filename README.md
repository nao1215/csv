<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/csv.svg)](https://pkg.go.dev/github.com/nao1215/csv)
[![MultiPlatformUnitTest](https://github.com/nao1215/csv/actions/workflows/multi_platform_ut.yml/badge.svg)](https://github.com/nao1215/csv/actions/workflows/multi_platform_ut.yml)
[![reviewdog](https://github.com/nao1215/csv/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/csv/actions/workflows/reviewdog.yml)
![Coverage](https://raw.githubusercontent.com/nao1215/octocovs-central-repo/main/badges/nao1215/csv/coverage.svg)

![csv-logo](./doc/images/csv-logo-small.png)

## What is the csv package?

The csv package provides two complementary features:

1. **CSV Validation using struct tags**  
   - Detect invalid values with precise error positions  
   - Multi-language error messages (English, Japanese, Russian)  
   - Validation rules expressed through `validate:` tags  

2. **A pandas-like DataFrame engine for CSV processing**  
   - Inspired by the author’s experience learning machine learning with pandas  
   - Built on top of [`nao1215/filesql`](https://github.com/nao1215/filesql)  
   - Brings SQL-backed, pandas-like transformations into pure Go  
   - Suitable for data exploration, ETL preprocessing, and joining multiple CSV files without Python

Both modules can be used independently.

---

## Supported OS and Go versions

- Linux  
- macOS  
- Windows  
- Go 1.24 or later

## Supported languages (validation only)

- English  
- Japanese  
- Russian  

---

## CSV Validation

Validation rules are written using struct tags.  
The order of struct fields **must match** the column order of the CSV file.

### Example (English messages)

```go
input := `id,name,age,password,password_confirm,role,note,nickname,ip,cidr,url
1,Alice,17,Secret123,Secret12,superuser,"TODO: fix",alice!,999.0.0.1,10.0.0.0/33,http://example.com
-5,Bob,30,short,short,admin,"Note: ready",Bob123,192.168.0.1,192.168.0.0/24,https://example.com
`

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

var accounts []account
buf := bytes.NewBufferString(input)
csvReader, _ := csv.NewCSV(buf)
errs := csvReader.Decode(&accounts)

for _, err := range errs {
    fmt.Println(err.Error())
}
````

Output:

```
line:2 column age: target is not greater than or equal to the threshold value: threshold=18, value=17
line:2 column password: target is not greater than or equal to the threshold value: value=Secret123
line:2 column role: target is not one of the values: oneof=admin user, value=superuser
line:2 column note: target contains a prohibited substring: excludes=TODO, value=TODO: fix
line:2 column note: target does not start with the specified value: startswith=Note, value=TODO: fix
line:2 column nickname: target is not an alphanumeric unicode character: value=alice!
line:2 column ip: target is not a valid IPv4 address: value=999.0.0.1
line:2 column cidr: target is not a valid IPv4 CIDR: value=10.0.0.0/33
line:2 column url: target is not a valid HTTPS URL: value=http://example.com
line:2 column password_confirm: field is not equal to the specified field: field=PasswordConfirm, other=Password
line:3 column id: target is not greater than or equal to the threshold value: threshold=1, value=-5
line:3 column password: target is not greater than or equal to the threshold value: value=short
```

---

### Validation Tags

#### Fields rules

| Tag Name     | Description                              |
| ------------ | ---------------------------------------- |
| eqfield    | Equal to another field in the same row |
| fieldcontains | Field contains the value of another field (same row) |
| fieldexcludes | Field does not contain the value of another field (same row) |
| gtefield      | Field Greater Than or Equal To Another Field (same row)|
| gtfield  | Greater than another field in the same row |
| ltefield | Less or equal to another field in the same row |
| ltfield  | Less than another field in the same row |
| nefield    | Not equal to another field in the same row |

#### String rules

| Tag Name     | Description                              |
| ------------ | ---------------------------------------- |
| alpha        | Alphabetic characters only               |
| alphaunicode | Unicode alphabetic characters            |
| alphanumeric | Alphanumeric characters                  |
| alphanumunicode | Unicode letters and digits            |
| alphaspace   | Alphabetic characters and spaces         |
| ascii        | ASCII characters only                    |
| boolean      | Boolean values                           |
| contains     | Contains substring                       |
| containsany  | Contains any of the specified characters |
| containsrune | Contains the specified rune              |
| endsnotwith  | Must not end with the specified substring|
| endswith     | Ends with the specified substring        |
| excludes     | Must not contain the specified substring |
| excludesall  | Must not contain any of the specified runes |
| excludesrune | Must not contain the specified rune      |
| lowercase    | Lowercase only                           |
| multibyte    | Contains at least one multibyte character|
| number       | Signed integer or decimal number         |
| numeric      | Numeric only                             |
| datauri      | Valid Data URI (data:*;base64,…)         |
| printascii   | Printable ASCII characters only          |
| startsnotwith| Must not start with the specified substring |
| startswith   | Starts with the specified substring      |
| uppercase    | Uppercase only                           |

#### Format rules

| Tag Name   | Description                                                                                 |
| ---------- | ------------------------------------------------------------------------------------------- |
| email      | Valid email address                                                                         |
| uri        | Valid URI (scheme required, host optional)                                                  |
| url        | Valid URL (scheme required; `file:` allows path-only, other schemes require URL shape)      |
| http_url   | Valid HTTP(S) URL with host                                                                 |
| https_url  | Valid HTTPS URL with host                                                                   |
| url_encoded| URL-encoded string (percent escapes, no malformed `%` sequences)                            |
| datauri    | Valid Data URI (data:*;base64,…)                                                            |
| hostname   | Hostname (RFC 952)                                                                          |
| hostname_rfc1123 | Hostname (RFC 1123)                                                                   |
| hostname_port | Host and port combination                                                                |
| uuid       | UUID string (with hyphens)                                                                  |

#### Network rules

| Tag Name | Description                      |
| -------- | -------------------------------- |
| cidr         | Valid CIDR (IPv4 or IPv6)                |
| cidrv4       | Valid IPv4 CIDR                          |
| cidrv6       | Valid IPv6 CIDR                          |
| ip_addr  | IPv4 or IPv6 address             |
| ip4_addr | IPv4 address only                |
| ip6_addr | IPv6 address only                |

#### Comparison rules

| Tag Name | Description                  |
| -------- | ---------------------------- |
| eq       | Equal to the specified value |
| eq_ignore_case | Equal to the specified value (case-insensitive) |
| gt       | Greater than                 |
| gte      | Greater or equal             |
| lt       | Less than                    |
| lte      | Less or equal                |
| ne_ignore_case | Not equal to the specified value (case-insensitive) |
| ne       | Not equal                    |

#### Other rules

| Tag Name | Description                    |
| -------- | ------------------------------ |
| len      | Exact length                   |
| max      | Maximum value                  |
| min      | Minimum value                  |
| oneof    | Must match one of given values |
| required | Must not be empty              |

---

## DataFrame Engine

The DataFrame API provides pandas-like transformations for CSV files using lazy SQL execution.

### DataFrame Methods

| Method                                             | Description                        |
| -------------------------------------------------- | ---------------------------------- |
| `func (df DataFrame) Select(cols ...string) DataFrame`   | Select columns                     |
| `func (df DataFrame) Filter(expr string) DataFrame`      | Filter rows using SQL expressions  |
| `func (df DataFrame) Drop(cols ...string) DataFrame`     | Drop columns                       |
| `func (df DataFrame) Rename(map[string]string) DataFrame`| Rename columns                     |
| `func (df DataFrame) Mutate(col, expr string) DataFrame` | Create derived columns             |
| `func (df DataFrame) Sort(col string, asc bool) DataFrame` | Sort rows                       |
| `func (df DataFrame) DropNA(cols ...string) DataFrame`   | Remove rows with NULL values       |
| `func (df DataFrame) FillNA(col string, value any) DataFrame` | Fill NULL values            |
| `func (df DataFrame) Cast(col, dtype string) DataFrame`  | Cast column type                   |
| `func (df DataFrame) Join(other DataFrame, key string) DataFrame` | Inner join          |
| `func (df DataFrame) LeftJoin(other DataFrame, key string) DataFrame` | Left join      |
| `func (df DataFrame) RightJoin(other DataFrame, key string) DataFrame` | Right join    |
| `func (df DataFrame) FullJoin(other DataFrame, key string) DataFrame` | Full join       |
| `func (df DataFrame) Merge(other DataFrame, opts MergeOptions) DataFrame` | pandas-like merge |
| `func (df DataFrame) Rows() ([]map[string]any, error)`   | Evaluate query and return results  |
| `func (df DataFrame) Head(n int) ([]map[string]any, error)` | Return first n rows           |
| `func (df DataFrame) Tail(n int) ([]map[string]any, error)` | Return last n rows            |
| `func (df DataFrame) Print(w io.Writer) error`           | Pretty-print table                 |
| `func (df DataFrame) ToCSV(path string) error`           | Write results as CSV               |
| `func (df DataFrame) DebugSQL() string`                  | Show generated SQL                 |
| `func (df DataFrame) Columns() []string`                 | List of columns                    |
| `func (df DataFrame) Shape() (int, int)`                 | (rows, columns)                    |

> **Notes**
>
> - `DataFrame.Filter` inlines the provided SQL expression verbatim. Always sanitize or parameterize user input before passing it to `Filter`.
> - `RightJoin`/`FullJoin` require backend support. SQLite (used by filesql) does not implement RIGHT/FULL OUTER JOIN, so those methods will return SQL errors on SQLite-based setups.

---

### DataFrame Example

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/nao1215/csv"
)

func main() {
	if err := os.WriteFile("people.csv", []byte("id,name,age\n1,Alice,23\n2,Bob,30\n"), 0o600); err != nil {
		log.Fatal(err)
	}

    df := csv.NewDataFrame("people.csv").
        Select("name", "age").
        Filter("age >= 25").
        Mutate("decade", "age / 10").
        Sort("age", true)

    rows, err := df.Rows()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(rows[0]["name"], rows[0]["decade"])
}
```

---

## License

MIT License. See [LICENSE](./LICENSE).

---

## Contribution

Contributions are welcome!
See [CONTRIBUTING.md](./CONTRIBUTING.md).

Star support is greatly appreciated.

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/csv\&type=Date)](https://star-history.com/#nao1215/csv&Date)

---

## Special Thanks

I was inspired by the following OSS. Thank you for your great work!!

* [go-playground/validator](https://github.com/go-playground/validator)
* [shogo82148/go-header-csv](https://github.com/shogo82148/go-header-csv)
* [pandas](https://pandas.pydata.org/)

---

## Contributors ✨

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->

<!-- prettier-ignore-start -->

<!-- markdownlint-disable -->

<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://debimate.jp/"><img src="https://avatars.githubusercontent.com/u/22737008?v=4?s=75" width="75px;" alt="CHIKAMATSU Naohiro"/><br /><sub><b>CHIKAMATSU Naohiro</b></sub></a><br /><a href="https://github.com/nao1215/csv/commits?author=nao1215" title="Documentation">📖</a></td>
    </tr>
  </tbody>
</table>
<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the all-contributors specification.
