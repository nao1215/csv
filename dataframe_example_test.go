package csv_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nao1215/csv"
)

func cleanupDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
}

func exampleCSV(data string) (string, func()) {
	dir, err := os.MkdirTemp("", "dfexample")
	if err != nil {
		panic(err)
	}
	path := filepath.Join(dir, "data.csv")
	payload := strings.TrimSpace(data) + "\n"
	if err := os.WriteFile(path, []byte(payload), 0o600); err != nil {
		cleanupDir(dir)
		panic(err)
	}
	cleanup := func() { cleanupDir(dir) }
	return path, cleanup
}

func exampleCSVWithName(name, data string) (string, func()) {
	dir, err := os.MkdirTemp("", "dfexample")
	if err != nil {
		panic(err)
	}
	path := filepath.Join(dir, name)
	payload := strings.TrimSpace(data) + "\n"
	if err := os.WriteFile(path, []byte(payload), 0o600); err != nil {
		cleanupDir(dir)
		panic(err)
	}
	cleanup := func() { cleanupDir(dir) }
	return path, cleanup
}

func ExampleNewDataFrame() {
	path, cleanup := exampleCSV(`
id,name
1,Alice`)
	defer cleanup()

	df := csv.NewDataFrame(path)
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(rows), rows[0]["name"])
	// Output:
	// 1 Alice
}

func ExampleDataFrame_Filter() {
	path, cleanup := exampleCSV(`
id,age
1,10
2,20`)
	defer cleanup()

	df := csv.NewDataFrame(path).Filter("age >= 15")
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["age"])
	// Output:
	// 20
}

func ExampleDataFrame_Select() {
	path, cleanup := exampleCSV(`
id,name,age
1,Alice,23`)
	defer cleanup()

	df := csv.NewDataFrame(path).Select("name")
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["name"])
	// Output:
	// Alice
}

func ExampleDataFrame_Drop() {
	path, cleanup := exampleCSV(`
id,name,age
1,Alice,23`)
	defer cleanup()

	df := csv.NewDataFrame(path).Drop("age")
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["name"])
	// Output:
	// Alice
}

func ExampleDataFrame_Rename() {
	path, cleanup := exampleCSV(`
id,name
1,Alice`)
	defer cleanup()

	df := csv.NewDataFrame(path).Rename(map[string]string{"name": "full_name"})
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["full_name"])
	// Output:
	// Alice
}

func ExampleDataFrame_Mutate() {
	path, cleanup := exampleCSV(`
id,value
1,10`)
	defer cleanup()

	df := csv.NewDataFrame(path).Mutate("double_value", "value * 2")
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["double_value"])
	// Output:
	// 20
}

func ExampleDataFrame_Sort() {
	path, cleanup := exampleCSV(`
id,value
1,10
2,5`)
	defer cleanup()

	df := csv.NewDataFrame(path).Sort("value", true)
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["value"])
	// Output:
	// 5
}

func ExampleDataFrame_DropNA() {
	path, cleanup := exampleCSV(`
id,value
1,
2,5`)
	defer cleanup()

	df := csv.NewDataFrame(path).
		Mutate("clean_value", "NULLIF(value, '')").
		DropNA("clean_value").
		Select("id")
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(rows))
	// Output:
	// 1
}

func ExampleDataFrame_FillNA() {
	path, cleanup := exampleCSV(`
id,value
1,
2,5`)
	defer cleanup()

	df := csv.NewDataFrame(path).FillNA("value", 0)
	fmt.Println(strings.Contains(df.DebugSQL(), "IFNULL(value, 0)"))
	// Output:
	// true
}

func ExampleDataFrame_Join() {
	leftPath, leftCleanup := exampleCSVWithName("users.csv", `
id,name
1,Alice`)
	defer leftCleanup()
	rightPath, rightCleanup := exampleCSVWithName("scores.csv", `
id,score
1,80`)
	defer rightCleanup()

	left := csv.NewDataFrame(leftPath)
	right := csv.NewDataFrame(rightPath)

	df := left.Join(right, "id").Select("name", "score")
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["score"])
	// Output:
	// 80
}

func ExampleDataFrame_Merge() {
	leftPath, leftCleanup := exampleCSVWithName("users.csv", `
id,name
1,Alice`)
	defer leftCleanup()
	rightPath, rightCleanup := exampleCSVWithName("scores.csv", `
id,score
1,80`)
	defer rightCleanup()

	df := csv.NewDataFrame(leftPath).Merge(csv.NewDataFrame(rightPath), csv.MergeOptions{OnKey: "id"})
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["score"])
	// Output:
	// 80
}

func ExampleDataFrame_Rows() {
	path, cleanup := exampleCSV(`
id,name
1,Alice`)
	defer cleanup()

	df := csv.NewDataFrame(path)
	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(rows))
	// Output:
	// 1
}

func ExampleDataFrame_Print() {
	path, cleanup := exampleCSV(`
id,name
1,Alice`)
	defer cleanup()

	df := csv.NewDataFrame(path)
	var buf bytes.Buffer
	if err := df.Print(&buf); err != nil {
		panic(err)
	}
	fmt.Print(buf.String())
	// Output:
	// id  name
	// 1   Alice
}

func ExampleDataFrame_DebugSQL() {
	path, cleanup := exampleCSV(`
id,name
1,Alice`)
	defer cleanup()

	df := csv.NewDataFrame(path).Select("name")
	fmt.Println(df.DebugSQL())
	// Output:
	// SELECT data.* FROM data
}

func ExampleDataFrame_Columns() {
	path, cleanup := exampleCSV(`
id,name
1,Alice`)
	defer cleanup()

	df := csv.NewDataFrame(path)
	fmt.Println(df.Columns())
	// Output:
	// [id name]
}

func ExampleDataFrame_Shape() {
	path, cleanup := exampleCSV(`
id,name
1,Alice`)
	defer cleanup()

	df := csv.NewDataFrame(path)
	r, c := df.Shape()
	fmt.Println(r, c)
	// Output:
	// 1 2
}

func Example_dataFrame_basic() {
	path, cleanup := exampleCSV(`
id,name,age
1,Alice,23
2,Bob,30`)
	defer cleanup()

	df := csv.NewDataFrame(path).
		Select("name", "age").
		Filter("age >= 25").
		Mutate("decade", "age / 10").
		Sort("age", true)

	rows, err := df.Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(rows[0]["name"], rows[0]["decade"])
	// Output:
	// Bob 3
}

func Example_dataFrame_join() {
	usersPath, cleanupUsers := exampleCSVWithName("users.csv", `
id,name
1,Alice
2,Bob`)
	defer cleanupUsers()
	purchasesPath, cleanupPurchases := exampleCSVWithName("purchases.csv", `
id,total
1,100`)
	defer cleanupPurchases()

	result, err := csv.NewDataFrame(usersPath).
		LeftJoin(csv.NewDataFrame(purchasesPath), "id").
		Rows()
	if err != nil {
		panic(err)
	}
	fmt.Println(result[0]["total"])
	// Output:
	// 100
}

func Example_dataFrame_cleaning() {
	path, cleanup := exampleCSV(`
id,name,score
1,Alice,
2,Bob,80`)
	defer cleanup()

	df := csv.NewDataFrame(path).
		FillNA("score", 0).
		Cast("score", "INTEGER").
		DropNA("name").
		Rename(map[string]string{"score": "final_score"}).
		Select("id", "name", "final_score")

	shapeRows, shapeCols := df.Shape()
	fmt.Printf("%d rows, %d cols\n", shapeRows, shapeCols)
	fmt.Println(df.Columns())
	// Output:
	// 2 rows, 3 cols
	// [final_score id name]
}

func ExampleDataFrame_joinFilterSort() {
	users := csv.NewDataFrame(filepath.Join("testdata", "sample.csv")).
		Select("id", "name", "age").
		Mutate("age_bucket", "CASE WHEN age >= 30 THEN '30s' ELSE '20s' END")

	orders := csv.NewDataFrame(filepath.Join("testdata", "orders.csv")).
		Filter("total >= 100").
		Mutate("gross_total", "total + 5")

	depts := csv.NewDataFrame(filepath.Join("testdata", "departments.csv")).
		Select("id", "dept").
		Rename(map[string]string{"dept": "dept_name"})

	df := users.
		Join(orders, "id").
		Join(depts, "id").
		Filter("age >= 23").
		Sort("gross_total", false).
		Select("name", "dept_name", "gross_total", "age_bucket")

	var buf bytes.Buffer
	if err := df.Print(&buf); err != nil {
		panic(err)
	}
	fmt.Print(buf.String())

	// Output:
	// age_bucket  dept_name    gross_total  name
	// 30s         Engineering  155          Denis
	// 20s         Sales        105          Gina
}
