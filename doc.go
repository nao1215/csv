// Package csv provides two main features:
//
// 1. CSV/TSV validation using struct tags
//   - Validates row/column values based on `validate:` rules.
//   - Returns detailed errors (row, column, rule violation).
//   - Supports multiple languages for error messages.
//
// 2. A pandas-like DataFrame API backed by SQL (filesql)
//   - Enables filtering, selecting, joining, mutating, sorting, casting,
//     and cleaning CSV/TSV data.
//   - Operations are lazy and compiled into a single SQL query on execution.
//   - Useful for lightweight data manipulation without Python.
//
// # Validation
//
// Define rules using struct tags:
//
//	type User struct {
//	    ID    int    `validate:"numeric"`
//	    Name  string `validate:"alpha"`
//	    Score int    `validate:"gte=0,lte=100"`
//	}
//
// Decode() reads the CSV/TSV and applies validation before populating the struct slice.
//
// # DataFrame
//
// DataFrame offers a chainable API:
//
//	df := csv.NewDataFrame("data.csv").
//	    Select("name", "age").
//	    Filter("age >= 20").
//	    Mutate("decade", "age / 10")
//
//	rows, _ := df.Rows()
//
// All transformations are evaluated lazily and executed as SQL via filesql.
//
// # Scope
//
// The package currently supports CSV and TSV files.
// DataFrame is intended as a lightweight, pandas-inspired data manipulation layer
// for Go, motivated by combining csv processing with the author's filesql engine.
//
// For full examples and details, see the README.
package csv
