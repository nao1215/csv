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
input := `id,name,age
1,Gina,23
a,Yulia,25
3,Den1s,30
`

type person struct {
    ID   int    `validate:"numeric"`
    Name string `validate:"alpha"`
    Age  int    `validate:"gt=24"`
}

var people []person
buf := bytes.NewBufferString(input)
csvReader, _ := csv.NewCSV(buf)
errs := csvReader.Decode(&people)

for _, err := range errs {
    fmt.Println(err.Error())
}
````

---

### Validation Tags

#### String rules

| Tag Name     | Description                              |
| ------------ | ---------------------------------------- |
| alpha        | Alphabetic characters only               |
| alphanumeric | Alphanumeric characters                  |
| ascii        | ASCII characters only                    |
| boolean      | Boolean values                           |
| contains     | Contains substring                       |
| containsany  | Contains any of the specified characters |
| startswith   | Starts with the specified substring      |
| endswith     | Ends with the specified substring        |
| lowercase    | Lowercase only                           |
| numeric      | Numeric only                             |
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
| uuid       | UUID string (with hyphens)                                                                  |

#### Network rules

| Tag Name | Description                      |
| -------- | -------------------------------- |
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
