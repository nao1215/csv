package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDataFrame(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "rows returns all records",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample.csv"))
				rows, err := df.Rows()
				if err != nil {
					t.Fatalf("Rows returned error: %v", err)
				}

				if len(rows) != 3 {
					t.Fatalf("expect 3 rows, got %d", len(rows))
				}

				if got := fmt.Sprint(rows[0]["name"]); got != "Gina" {
					t.Fatalf("unexpected first row name: %s", got)
				}
			},
		},
		{
			name: "head and tail",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample.csv"))

				head, err := df.Head(2)
				if err != nil {
					t.Fatalf("Head returned error: %v", err)
				}
				if len(head) != 2 {
					t.Fatalf("expect head 2 rows, got %d", len(head))
				}
				if got := fmt.Sprint(head[1]["name"]); got != "Yulia" {
					t.Fatalf("unexpected second row name in head: %s", got)
				}

				tail, err := df.Tail(2)
				if err != nil {
					t.Fatalf("Tail returned error: %v", err)
				}
				if len(tail) != 2 {
					t.Fatalf("expect tail 2 rows, got %d", len(tail))
				}
				if got := fmt.Sprint(tail[0]["name"]); got != "Yulia" || fmt.Sprint(tail[1]["name"]) != "Denis" {
					t.Fatalf("unexpected tail rows: %+v", tail)
				}
			},
		},
		{
			name: "filter and select",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample.csv")).Filter("age >= 25").Select("name", "age")
				rows, err := df.Rows()
				if err != nil {
					t.Fatalf("Rows returned error: %v", err)
				}

				if len(rows) != 2 {
					t.Fatalf("expect 2 rows after filter, got %d", len(rows))
				}
				if len(rows[0]) != 2 {
					t.Fatalf("expect only selected columns, got %d columns", len(rows[0]))
				}
				if got := fmt.Sprint(rows[0]["name"]); got != "Yulia" {
					t.Fatalf("unexpected first row after filter: %s", got)
				}
			},
		},
		{
			name: "mutate rename and drop",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample.csv")).
					Mutate("age_next", "age + 1").
					Rename(map[string]string{"name": "username"}).
					Drop("age")

				rows, err := df.Head(1)
				if err != nil {
					t.Fatalf("Head returned error: %v", err)
				}
				if len(rows) != 1 {
					t.Fatalf("expect 1 row, got %d", len(rows))
				}

				row := rows[0]
				if _, ok := row["age"]; ok {
					t.Fatalf("age should be dropped: %+v", row)
				}
				if got := fmt.Sprint(row["username"]); got != "Gina" {
					t.Fatalf("unexpected renamed column value: %s", got)
				}
				if got := fmt.Sprint(row["age_next"]); got != "24" {
					t.Fatalf("unexpected mutated value: %s", got)
				}
			},
		},
		{
			name: "join",
			fn: func(t *testing.T) {
				t.Helper()
				users := NewDataFrame(filepath.Join("testdata", "sample.csv"))
				orders := NewDataFrame(filepath.Join("testdata", "orders.csv"))

				rows, err := users.Join(orders, "id").
					Select("sample.name", "orders.total").
					Rows()
				if err != nil {
					t.Fatalf("Join Rows returned error: %v", err)
				}

				if len(rows) != 2 {
					t.Fatalf("expect 2 joined rows, got %d", len(rows))
				}
				if got := fmt.Sprint(rows[0]["name"]); got != "Gina" || fmt.Sprint(rows[1]["name"]) != "Denis" {
					t.Fatalf("unexpected join result: %+v", rows)
				}
				if got := fmt.Sprint(rows[0]["total"]); got != "100" {
					t.Fatalf("unexpected join column value: %s", got)
				}
			},
		},
		{
			name: "to csv",
			fn: func(t *testing.T) {
				t.Helper()
				output := filepath.Join(t.TempDir(), "out.csv")
				df := NewDataFrame(filepath.Join("testdata", "sample.csv")).Filter("age >= 25")

				if err := df.ToCSV(output); err != nil {
					t.Fatalf("ToCSV returned error: %v", err)
				}

				file, err := os.Open(filepath.Clean(output))
				if err != nil {
					t.Fatalf("failed to open output csv: %v", err)
				}
				defer file.Close()

				r := csv.NewReader(file)
				records, err := r.ReadAll()
				if err != nil {
					t.Fatalf("failed to read output csv: %v", err)
				}

				if len(records) != 3 { // header + 2 rows
					t.Fatalf("unexpected record count: %d", len(records))
				}
				if records[0][0] != "age" || records[0][1] != "id" || records[0][2] != "name" {
					t.Fatalf("unexpected header: %v", records[0])
				}
			},
		},
		{
			name: "print to writer",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample.csv")).Select("id", "name")
				var buf bytes.Buffer
				if err := df.Print(&buf); err != nil {
					t.Fatalf("Print returned error: %v", err)
				}
				if buf.Len() == 0 {
					t.Fatalf("expected output from Print, got none")
				}
			},
		},
		{
			name: "sort descending",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample.csv")).Sort("age", false)
				rows, err := df.Rows()
				if err != nil {
					t.Fatalf("Rows returned error: %v", err)
				}
				if got := fmt.Sprint(rows[0]["name"]); got != "Denis" {
					t.Fatalf("expected Denis first after sort desc, got %s", got)
				}
			},
		},
		{
			name: "cast to integer",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample.csv")).Cast("age", "INTEGER")
				rows, err := df.Rows()
				if err != nil {
					t.Fatalf("Rows returned error: %v", err)
				}
				if _, ok := rows[0]["age"].(int64); !ok {
					t.Fatalf("expected age to be int64 after cast, got %T", rows[0]["age"])
				}
			},
		},
		{
			name: "dropna filters nulls",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample_with_null.csv")).DropNA("age")
				qs := queryState{
					baseTable: tableNameFromPath(df.source),
					renames:   map[string]string{},
					drops:     map[string]struct{}{},
				}
				for _, op := range df.ops {
					op.Apply(&qs)
				}
				sql := compileSQL(qs, nil)
				if !strings.Contains(sql, "age IS NOT NULL") {
					t.Fatalf("expected DropNA to add filter, got %s", sql)
				}
			},
		},
		{
			name: "fillna replaces nulls",
			fn: func(t *testing.T) {
				t.Helper()
				df := NewDataFrame(filepath.Join("testdata", "sample_with_null.csv")).FillNA("age", 0)
				qs := queryState{
					baseTable: tableNameFromPath(df.source),
					renames:   map[string]string{},
					drops:     map[string]struct{}{},
				}
				for _, op := range df.ops {
					op.Apply(&qs)
				}
				sql := compileSQL(qs, nil)
				if !strings.Contains(sql, "IFNULL(age, 0)") {
					t.Fatalf("expected FillNA to add IFNULL, got %s", sql)
				}
			},
		},
		{
			name: "join preserves duplicate columns with suffix",
			fn: func(t *testing.T) {
				t.Helper()
				users := NewDataFrame(filepath.Join("testdata", "sample.csv"))
				orders := NewDataFrame(filepath.Join("testdata", "orders.csv"))

				rows, err := users.Join(orders, "id").Head(1)
				if err != nil {
					t.Fatalf("Join returned error: %v", err)
				}
				row := rows[0]
				if _, ok := row["id_x"]; !ok {
					t.Fatalf("base id column with suffix missing: %+v", row)
				}
				if _, ok := row["id_y"]; !ok {
					t.Fatalf("joined id column with suffix missing: %+v", row)
				}
			},
		},
		{
			name: "merge multiple keys keeps suffixes",
			fn: func(t *testing.T) {
				t.Helper()
				users := NewDataFrame(filepath.Join("testdata", "sample.csv"))
				orders := NewDataFrame(filepath.Join("testdata", "orders.csv"))

				rows, err := users.Merge(orders, MergeOptions{
					On: []string{"id"},
				}).Head(1)
				if err != nil {
					t.Fatalf("Merge returned error: %v", err)
				}
				row := rows[0]
				if _, ok := row["id_x"]; !ok {
					t.Fatalf("expected left suffix id_x: %+v", row)
				}
				if _, ok := row["id_y"]; !ok {
					t.Fatalf("expected right suffix id_y: %+v", row)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.fn(t)
		})
	}
}

func TestInternalHelpers(t *testing.T) {
	t.Parallel()

	t.Run("literal", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			in   any
			want string
		}{
			{nil, "NULL"},
			{"O'Hara", "'O''Hara'"},
			{[]byte("bytes"), "'bytes'"},
			{true, "1"},
			{false, "0"},
			{42, "42"},
		}
		for _, tt := range tests {
			if got := literal(tt.in); got != tt.want {
				t.Fatalf("literal(%v)=%s, want %s", tt.in, got, tt.want)
			}
		}
	})

	t.Run("apply operations builds query", func(t *testing.T) {
		t.Parallel()

		qs := queryState{
			baseTable: "users",
			renames:   map[string]string{},
			drops:     map[string]struct{}{},
		}

		ops := []sqlOp{
			filterOp{expr: "age > 20"},
			selectOp{cols: []string{"id", "name", "age"}},
			dropOp{cols: []string{"unused"}},
			renameOp{mapping: map[string]string{"name": "username"}},
			mutateOp{col: "age_next", expr: "age + 1"},
			sortOp{col: "age", asc: false},
			castOp{col: "age", dtype: "INTEGER"},
			dropNaOp{cols: []string{"age"}},
			fillNaOp{col: "name", value: "unknown"},
		}

		for _, op := range ops {
			op.Apply(&qs)
		}

		sql := compileSQL(qs, nil)

		if !strings.Contains(sql, "age IS NOT NULL") || !strings.Contains(sql, "age > 20") {
			t.Fatalf("expected filters in SQL, got %s", sql)
		}
		if !strings.Contains(sql, "ORDER BY age DESC") {
			t.Fatalf("expected order by in SQL, got %s", sql)
		}
		if !strings.Contains(sql, "CAST(age AS INTEGER)") {
			t.Fatalf("expected cast in SQL, got %s", sql)
		}
		if !strings.Contains(sql, "IFNULL(name, 'unknown')") {
			t.Fatalf("expected fillna in SQL, got %s", sql)
		}
		if !strings.Contains(sql, "age + 1 AS age_next") {
			t.Fatalf("expected mutation in SQL, got %s", sql)
		}
	})
}

func TestSelectThenRenameKeepsColumn(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "select_rename.csv")).
		Select("A").
		Rename(map[string]string{"A": "B"})

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if _, ok := rows[0]["B"]; !ok {
		t.Fatalf("expected column B to exist after rename, row: %+v", rows[0])
	}
}

func TestSelectThenMutateKeepsDerivedColumn(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "select_mutate.csv")).
		Select("A").
		Mutate("C", "A + 1")

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	if _, ok := rows[0]["C"]; !ok {
		t.Fatalf("expected derived column C, row: %+v", rows[0])
	}
}

func TestMergeSelectUndefinedColumnShouldError(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "merge_left.csv")).Merge(
		NewDataFrame(filepath.Join("testdata", "merge_right.csv")),
		MergeOptions{
			On:       []string{"id"},
			Suffixes: [2]string{"_L", "_R"},
		},
	).Select("val")

	if _, err := df.Rows(); err == nil {
		t.Fatal("expected error selecting undefined column without suffix, got none")
	}
}

func TestDropDoesNotRemoveSuffixColumns(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "merge_left.csv")).Merge(
		NewDataFrame(filepath.Join("testdata", "merge_right.csv")),
		MergeOptions{On: []string{"id"}},
	).Drop("val")

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	if _, ok := rows[0]["val_x"]; !ok {
		t.Fatalf("expected val_x to remain, row: %+v", rows[0])
	}
	if _, ok := rows[0]["val_y"]; !ok {
		t.Fatalf("expected val_y to remain, row: %+v", rows[0])
	}
}

func TestRenameThenCastShouldUseNewName(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "select_mutate.csv")).
		Rename(map[string]string{"A": "B"}).
		Cast("B", "INTEGER")

	if _, err := df.Rows(); err != nil {
		t.Fatalf("expected rename then cast to succeed, got error: %v", err)
	}
}

func TestMultipleMergesKeepIndependentSuffixes(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "merge_left.csv")).
		Merge(NewDataFrame(filepath.Join("testdata", "merge_b.csv")), MergeOptions{
			On:       []string{"id"},
			Suffixes: [2]string{"_A", "_B"},
		}).
		Merge(NewDataFrame(filepath.Join("testdata", "merge_c.csv")), MergeOptions{
			On:       []string{"id"},
			Suffixes: [2]string{"_X", "_Y"},
		})

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	row := rows[0]
	for _, col := range []string{"val_A", "val_B", "val_X", "val_Y"} {
		if _, ok := row[col]; !ok {
			t.Fatalf("expected column %s to exist, row: %+v", col, row)
		}
	}
}

func TestFillNAKeepsNumericType(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "fillna_nums.csv")).Cast("A", "INTEGER").FillNA("A", 0)

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	if _, ok := rows[1]["A"].(int64); !ok {
		t.Fatalf("expected filled value to remain integer, got %T", rows[1]["A"])
	}
}

func TestRenameThenSelect(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "rename_select.csv")).
		Rename(map[string]string{"A": "B"}).
		Select("B")

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if _, ok := rows[0]["B"]; !ok {
		t.Fatalf("expected renamed column B, row: %+v", rows[0])
	}
}

func TestSelectDottedColumnName(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "dotted.csv")).Select("A.B")

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	if _, ok := rows[0]["A.B"]; !ok {
		t.Fatalf("expected dotted column name to be preserved, row: %+v", rows[0])
	}
}

func TestRenameThenDrop(t *testing.T) {
	t.Parallel()

	df := NewDataFrame(filepath.Join("testdata", "drop_rename.csv")).
		Rename(map[string]string{"A": "Z"}).
		Drop("Z")

	rows, err := df.Rows()
	if err != nil {
		t.Fatalf("Rows returned error: %v", err)
	}
	if len(rows[0]) != 0 {
		t.Fatalf("expected all columns dropped, row: %+v", rows[0])
	}
}

func TestRenameChainThenSelectKeepsFinalName(t *testing.T) {
	t.Run("rename chain leaves only final column", func(t *testing.T) {
		t.Parallel()

		df := NewDataFrame(filepath.Join("testdata", "rename_chain.csv")).
			Rename(map[string]string{"a": "x"}).
			Rename(map[string]string{"x": "y"}).
			Select("y")

		rows, err := df.Rows()
		if err != nil {
			t.Fatalf("Rows returned error: %v", err)
		}
		if len(rows) == 0 {
			t.Fatalf("expected rows, got none")
		}
		row := rows[0]
		if len(row) != 1 {
			t.Fatalf("expected only final column y, row: %+v", row)
		}
		if _, ok := row["y"]; !ok {
			t.Fatalf("expected column y after chained renames, row: %+v", row)
		}
	})
}

func TestSelectDeduplicatesMixedCaseColumns(t *testing.T) {
	t.Run("select removes case-insensitive duplicates", func(t *testing.T) {
		t.Parallel()

		df := NewDataFrame(filepath.Join("testdata", "select_case_dup.csv")).
			Select("A", "a", "b", "A")

		rows, err := df.Rows()
		if err != nil {
			t.Fatalf("Rows returned error: %v", err)
		}
		if len(rows) == 0 {
			t.Fatalf("expected rows, got none")
		}
		row := rows[0]
		if len(row) != 2 {
			t.Fatalf("expected only columns A and b, row: %+v", row)
		}
		if _, ok := row["A"]; !ok {
			t.Fatalf("expected column A to remain, row: %+v", row)
		}
		if _, ok := row["b"]; !ok {
			t.Fatalf("expected column b to remain, row: %+v", row)
		}
	})
}

func TestRenameThenMutateUsesRenamedColumn(t *testing.T) {
	t.Run("mutate operates on renamed column", func(t *testing.T) {
		t.Parallel()

		df := NewDataFrame(filepath.Join("testdata", "ab_simple.csv")).
			Rename(map[string]string{"a": "x"}).
			Mutate("x", "x + 5")

		rows, err := df.Rows()
		if err != nil {
			t.Fatalf("Rows returned error: %v", err)
		}
		if _, ok := rows[0]["x"]; !ok {
			t.Fatalf("expected mutate to use renamed column x, row: %+v", rows[0])
		}
		if _, ok := rows[0]["a"]; ok {
			t.Fatalf("did not expect original column a to remain, row: %+v", rows[0])
		}
	})
}

func TestDropAfterRenameRemovesBothNames(t *testing.T) {
	t.Run("drop removes renamed aliases", func(t *testing.T) {
		t.Parallel()

		df := NewDataFrame(filepath.Join("testdata", "ab_simple.csv")).
			Rename(map[string]string{"a": "x"}).
			Drop("x")

		rows, err := df.Rows()
		if err != nil {
			t.Fatalf("Rows returned error: %v", err)
		}
		if len(rows) == 0 {
			t.Fatalf("expected rows, got none")
		}
		row := rows[0]
		if _, ok := row["x"]; ok {
			t.Fatalf("expected renamed column x to be dropped, row: %+v", row)
		}
		if _, ok := row["a"]; ok {
			t.Fatalf("expected original column a to be dropped as well, row: %+v", row)
		}
	})
}

func TestApplyJoinSuffixesHandlesNumericSuffixBounds(t *testing.T) {
	t.Run("apply join suffixes respects numeric bounds", func(t *testing.T) {
		t.Parallel()

		rows := []map[string]any{
			{
				"col":      "base",
				"col_1":    "join1",
				"col_2":    "join2",
				"col_2025": "actual",
			},
		}
		joins := []joinClause{
			{Suffixes: [2]string{"_base", "_j1"}},
			{Suffixes: [2]string{"_base2", "_j2"}},
		}

		applyJoinSuffixes(rows, joins)

		row := rows[0]
		for _, col := range []string{"col_base", "col_j1", "col_base2", "col_j2"} {
			if _, ok := row[col]; !ok {
				t.Fatalf("expected column %s to exist after suffixing, row: %+v", col, row)
			}
		}
		if _, ok := row["col_2025"]; !ok {
			t.Fatalf("expected literal column col_2025 to remain, row: %+v", row)
		}
		if _, ok := row["col"]; ok {
			t.Fatalf("expected original duplicate key to be removed, row: %+v", row)
		}
		if _, ok := row["col_2"]; ok {
			t.Fatalf("expected duplicate placeholder col_2 to be removed, row: %+v", row)
		}
	})
}

func TestJoinConditionDoesNotTreatCompositeExpressionsAsSimple(t *testing.T) {
	t.Run("composite join condition stays custom", func(t *testing.T) {
		t.Parallel()

		if isExplicitJoinCondition("a = b AND c = d") {
			t.Fatalf("expected composite join condition to be treated as custom expression")
		}
	})
}

func TestJoinConditionRejectsIncompleteExpressions(t *testing.T) {
	t.Run("incomplete join condition stays custom", func(t *testing.T) {
		t.Parallel()

		if isExplicitJoinCondition(" id = ") {
			t.Fatalf("expected incomplete join condition to be treated as custom expression")
		}
	})
}

func TestFillNADoesNotResurrectUnselectedColumns(t *testing.T) {
	t.Run("FillNA skips unselected columns", func(t *testing.T) {
		t.Parallel()

		df := NewDataFrame(filepath.Join("testdata", "ab_simple.csv")).
			Select("a").
			FillNA("b", 0)

		rows, err := df.Rows()
		if err != nil {
			t.Fatalf("Rows returned error: %v", err)
		}
		row := rows[0]
		if len(row) != 1 {
			t.Fatalf("expected only column a after FillNA, row: %+v", row)
		}
		if _, ok := row["a"]; !ok {
			t.Fatalf("expected column a to remain, row: %+v", row)
		}
		if _, ok := row["b"]; ok {
			t.Fatalf("did not expect column b to be resurrected, row: %+v", row)
		}
	})
}

func TestCastDoesNotResurrectUnselectedColumns(t *testing.T) {
	t.Run("Cast skips unselected columns", func(t *testing.T) {
		t.Parallel()

		df := NewDataFrame(filepath.Join("testdata", "ab_simple.csv")).
			Select("a").
			Cast("b", "TEXT")

		rows, err := df.Rows()
		if err != nil {
			t.Fatalf("Rows returned error: %v", err)
		}
		row := rows[0]
		if len(row) != 1 {
			t.Fatalf("expected only column a after Cast, row: %+v", row)
		}
		if _, ok := row["a"]; !ok {
			t.Fatalf("expected column a to remain, row: %+v", row)
		}
		if _, ok := row["b"]; ok {
			t.Fatalf("did not expect column b to be resurrected, row: %+v", row)
		}
	})
}

func TestSortUsesRenamedColumn(t *testing.T) {
	t.Run("sort orders by renamed column", func(t *testing.T) {
		t.Parallel()

		df := NewDataFrame(filepath.Join("testdata", "rename_sort.csv")).
			Rename(map[string]string{"a": "x"}).
			Rename(map[string]string{"x": "y"}).
			Sort("y", true)

		rows, err := df.Rows()
		if err != nil {
			t.Fatalf("Rows returned error: %v", err)
		}
		if len(rows) != 2 {
			t.Fatalf("expected two rows, got %d", len(rows))
		}
		if fmt.Sprint(rows[0]["y"]) != "1" || fmt.Sprint(rows[1]["y"]) != "2" {
			t.Fatalf("expected rows to be sorted by renamed column y, rows: %+v", rows)
		}
	})
}

func TestApplyJoinSuffixesStableOrdering(t *testing.T) {
	t.Run("applies suffix order consistently per join", func(t *testing.T) {
		rows := []map[string]any{
			{
				"col":   "left",
				"col_1": "right1",
				"col_2": "right2",
			},
		}
		joins := []joinClause{
			{Suffixes: [2]string{"_L1", "_R1"}},
			{Suffixes: [2]string{"_L2", "_R2"}},
		}

		applyJoinSuffixes(rows, joins)

		row := rows[0]
		if got := row["col_L1"]; got != "left" {
			t.Fatalf("expected col_L1 to hold left value, got %v", got)
		}
		if got := row["col_R1"]; got != "right1" {
			t.Fatalf("expected col_R1 to hold first right value, got %v", got)
		}
		if got := row["col_L2"]; got != "left" {
			t.Fatalf("expected col_L2 to still mirror left value, got %v", got)
		}
		if got := row["col_R2"]; got != "right2" {
			t.Fatalf("expected col_R2 to hold second right value, got %v", got)
		}
	})
}

func TestBuildSelectExpressionsKeepsMutationBeforeRename(t *testing.T) {
	t.Run("mutations precede inline renames", func(t *testing.T) {
		q := queryState{
			baseTable:     "users",
			inlineRenames: true,
			renames:       map[string]string{"A": "B"},
			mutations: []mutation{
				{Column: "C", Expr: "A + 1"},
			},
		}

		exprs := buildSelectExpressions(q)
		mutIdx := indexOfExpr(exprs, "AS C")
		renameIdx := indexOfExpr(exprs, "AS B")
		if mutIdx == -1 || renameIdx == -1 {
			t.Fatalf("expected both mutation and rename expressions, got %v", exprs)
		}
		if mutIdx >= renameIdx {
			t.Fatalf("expected mutation before rename, got %v", exprs)
		}
	})
}

func TestRewriteExpressionWithOriginalNamesRespectsBoundaries(t *testing.T) {
	t.Run("does not replace substrings within identifiers", func(t *testing.T) {
		entries := []renameEntry{{old: "rate", new: "rateLimit"}}
		expr := "rateLimit + rateLimiter"
		got := rewriteExpressionWithOriginalNames(expr, entries)
		if got != "rate + rateLimiter" {
			t.Fatalf("unexpected rewrite result: %s", got)
		}
	})
}

func TestResolveRowKeySuffixPreference(t *testing.T) {
	t.Run("falls back to default suffixes but not auto indexes", func(t *testing.T) {
		row := map[string]any{
			"foo_x": 10,
			"foo_y": 20,
			"foo_2": 30,
		}
		key, ok := resolveRowKey(row, "tbl.foo")
		if !ok {
			t.Fatalf("expected foo suffix to resolve")
		}
		if key != "foo_x" {
			t.Fatalf("expected foo_x to be chosen first, got %s", key)
		}
	})
}

func indexOfExpr(exprs []string, needle string) int {
	for i, expr := range exprs {
		if strings.Contains(expr, needle) {
			return i
		}
	}
	return -1
}
