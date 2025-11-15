package csv

import (
	"context"
	"database/sql"
	stdcsv "encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/nao1215/filesql"
)

// DataFrame emulates a subset of pandas.DataFrame backed by lazy SQL execution.
// Operations are accumulated and compiled into a single SQL query at materialization time.
type DataFrame struct {
	source       string  // Primary CSV file path
	ops          []sqlOp // Lazy operations (SQL AST)
	err          error
	warns        *warningBag
	aliasCounter *int
}

type warningBag struct {
	list []string
}

func (wb *warningBag) add(msg string) {
	if wb == nil {
		return
	}
	wb.list = append(wb.list, msg)
}

// sqlOp defines the interface for lazy operations.
type sqlOp interface {
	RequiredFiles() []string
	Apply(q *queryState)
}

// queryState keeps the aggregated query state before compilation.
type queryState struct {
	baseTable string

	selectCols []string // final projection (pandas-like column selection)
	filters    []string
	joins      []joinClause
	sorts      []sortClause

	mutations []mutation
	casts     []castClause
	fillNAs   []fillNAClause

	renames       map[string]string
	renameOrder   []renameEntry
	drops         map[string]struct{}
	suffixes      [2]string
	inlineRenames bool
	warnings      *warningBag
}

func newQueryState(base string, warns *warningBag) queryState {
	return queryState{
		baseTable: base,
		renames:   map[string]string{},
		drops:     map[string]struct{}{},
		suffixes:  defaultSuffixes(),
		warnings:  warns,
	}
}

func (df DataFrame) buildQueryState() queryState {
	state := newQueryState(tableNameFromPath(df.source), df.warns)
	for _, op := range df.ops {
		op.Apply(&state)
	}
	return state
}

// mutation represents a new column defined by an expression.
type mutation struct {
	Column string
	Expr   string
}

// castClause represents a CAST on a column.
type castClause struct {
	Column string
	DType  string
}

// fillNAClause represents a fill operation using IFNULL.
type fillNAClause struct {
	Column string
	Value  any
}

// joinType defines join type.
type joinType string

// Supported join types.
const (
	joinTypeInner joinType = "INNER"
	joinTypeLeft  joinType = "LEFT"
	joinTypeRight joinType = "RIGHT"
	joinTypeFull  joinType = "FULL"
)

// joinClause stores join information.
type joinClause struct {
	OtherTable string
	On         string
	Type       joinType
	SubQuery   string
	Suffixes   [2]string
}

const (
	mergeHowInner = "inner"
	mergeHowLeft  = "left"
	mergeHowRight = "right"
	mergeHowOuter = "outer"
	mergeHowFull  = "full"
)

type renameEntry struct {
	old string
	new string
}

func (q *queryState) resolveOriginalColumn(name string) string {
	current := normalizeColumnName(name)
	for i := len(q.renameOrder) - 1; i >= 0; i-- {
		entry := q.renameOrder[i]
		if entry.new == current {
			current = entry.old
		}
	}
	return current
}

func (q *queryState) resolveCurrentColumn(name string) string {
	current := normalizeColumnName(name)
	for _, entry := range q.renameOrder {
		if entry.old == current {
			current = entry.new
		}
	}
	return current
}

// sortClause stores order by information.
type sortClause struct {
	Column string
	Asc    bool
}

// NewDataFrame behaves similarly to pandas.read_csv, returning a DataFrame backed by the file.
// The DataFrame records operations lazily until materialization (Rows/Head/etc.).
func NewDataFrame(path string) DataFrame {
	return DataFrame{
		source:       path,
		warns:        &warningBag{},
		aliasCounter: new(int),
	}
}

// Filter acts like pandas.DataFrame.query, returning a new lazy DataFrame with an added WHERE clause.
// Expressions are passed directly through to SQLite: never concatenate untrusted user input here without validation.
// Invalid columns are reported through Err/Warns when executed.
func (df DataFrame) Filter(expr string) DataFrame {
	return df.withOp(filterOp{expr: expr})
}

// Select mirrors pandas column projection on a lazy SQL query.
// Nonexistent columns are skipped while emitting warnings retrievable via Warnings().
func (df DataFrame) Select(cols ...string) DataFrame {
	return df.withOp(selectOp{cols: cols})
}

// Drop removes columns similar to pandas.DataFrame.drop, but missing columns emit warnings instead of failing.
func (df DataFrame) Drop(cols ...string) DataFrame {
	return df.withOp(dropOp{cols: cols})
}

// Rename performs column renaming like pandas.DataFrame.rename in a lazy fashion.
// Missing columns are ignored, and a warning is recorded instead of raising an error.
func (df DataFrame) Rename(mapping map[string]string) DataFrame {
	return df.withOp(renameOp{mapping: mapping})
}

// Mutate behaves like pandas.DataFrame.assign, creating derived columns backed by SQL expressions.
// Expressions reference the current column names and are evaluated lazily.
func (df DataFrame) Mutate(col string, expr string) DataFrame {
	return df.withOp(mutateOp{col: col, expr: expr})
}

// Sort matches pandas.DataFrame.sort_values for a single column, ordering results at execution time.
func (df DataFrame) Sort(col string, asc bool) DataFrame {
	return df.withOp(sortOp{col: col, asc: asc})
}

// Cast casts a column akin to pandas.Series.astype without eagerly validating column existence.
// Missing columns are skipped with warnings.
func (df DataFrame) Cast(col string, dtype string) DataFrame {
	return df.withOp(castOp{col: col, dtype: dtype})
}

// DropNA behaves like pandas.DataFrame.dropna(subset=cols) with AND semantics across provided columns.
func (df DataFrame) DropNA(cols ...string) DataFrame {
	return df.withOp(dropNaOp{cols: cols})
}

// FillNA mirrors pandas.Series.fillna for a column but records warnings if the column is absent.
func (df DataFrame) FillNA(col string, value any) DataFrame {
	return df.withOp(fillNaOp{col: col, value: value})
}

// Join performs an INNER JOIN similar to pandas.merge(how="inner") using a single key.
func (df DataFrame) Join(other DataFrame, on string) DataFrame {
	return df.joinWithType(other, []string{on}, joinTypeInner, defaultSuffixes())
}

// LeftJoin performs a LEFT JOIN similar to pandas.merge(..., how="left").
func (df DataFrame) LeftJoin(other DataFrame, on string) DataFrame {
	return df.joinWithType(other, []string{on}, joinTypeLeft, defaultSuffixes())
}

// RightJoin performs a RIGHT JOIN similar to pandas.merge(..., how="right").
// Support depends on the underlying SQL engine configured via filesql; SQLite backends may not support RIGHT JOIN.
// TODO: emulate RIGHT JOIN for SQLite (e.g., via UNION) so behavior matches pandas even on limited engines.
func (df DataFrame) RightJoin(other DataFrame, on string) DataFrame {
	return df.joinWithType(other, []string{on}, joinTypeRight, defaultSuffixes())
}

// FullJoin performs a FULL OUTER JOIN similar to pandas.merge(..., how="outer").
// Support depends on the SQL engine; SQLite backends typically do not support FULL OUTER JOIN.
// TODO: emulate FULL OUTER JOIN (LEFT/RIGHT UNION) for SQLite so callers see pandas-like results.
func (df DataFrame) FullJoin(other DataFrame, on string) DataFrame {
	return df.joinWithType(other, []string{on}, joinTypeFull, defaultSuffixes())
}

// MergeOptions configures Merge; set either On or OnKey and optionally How/Suffixes.
type MergeOptions struct {
	On    []string
	OnKey string
	// How accepts "", "inner", "left", "right", "outer", or "full".
	How      string
	Suffixes [2]string
}

// Merge merges two DataFrames similarly to pandas.merge while validating join keys eagerly.
// Missing join keys record an error retrievable via Err().
func (df DataFrame) Merge(other DataFrame, opts MergeOptions) DataFrame {
	keys := opts.On
	if len(keys) == 0 && opts.OnKey != "" {
		keys = []string{opts.OnKey}
	}
	if len(keys) == 0 {
		next := df
		next.err = errors.New("merge requires at least one join key")
		return next
	}

	suffixes := opts.Suffixes
	if suffixes == ([2]string{}) {
		suffixes = defaultSuffixes()
	}

	var jt joinType
	switch strings.ToLower(opts.How) {
	case "", mergeHowInner:
		jt = joinTypeInner
	case mergeHowLeft:
		jt = joinTypeLeft
	case mergeHowRight:
		jt = joinTypeRight
	case mergeHowOuter, mergeHowFull:
		jt = joinTypeFull
	default:
		jt = joinTypeInner
	}
	return df.joinWithType(other, keys, jt, suffixes)
}

func defaultSuffixes() [2]string { return [2]string{"_x", "_y"} }

func joinConditionFromKeys(base string, other string, keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s.%s = %s.%s", base, k, other, k))
	}
	return strings.Join(parts, " AND ")
}

func (df DataFrame) joinWithType(other DataFrame, keys []string, jt joinType, suffixes [2]string) DataFrame {
	aliasBase := tableNameFromPath(other.source)
	if len(keys) == 0 {
		return df
	}
	alias := df.nextAlias(aliasBase)

	subQueryState := newQueryState(aliasBase, df.warns)
	subQueryState.inlineRenames = true

	for _, op := range other.ops {
		op.Apply(&subQueryState)
	}

	subQuery := compileSQL(subQueryState, nil)

	return df.withOp(joinOp{
		alias:    alias,
		subquery: subQuery,
		files:    other.collectRequiredFiles(),
		on:       joinConditionFromKeys(tableNameFromPath(df.source), alias, keys),
		joinType: jt,
		suffixes: suffixes,
	})
}

// Rows materializes the DataFrame similar to pandas.DataFrame.to_dict("records").
func (df DataFrame) Rows() ([]map[string]any, error) {
	return df.execute(context.Background(), nil)
}

// Head returns the first n rows, matching pandas.DataFrame.head.
func (df DataFrame) Head(n int) ([]map[string]any, error) {
	if n < 0 {
		n = 0
	}
	limit := n
	return df.execute(context.Background(), &limit)
}

// Tail mirrors pandas.DataFrame.tail by fetching all rows before slicing.
func (df DataFrame) Tail(n int) ([]map[string]any, error) {
	rows, err := df.Rows()
	if err != nil {
		return nil, err
	}
	if n <= 0 || n >= len(rows) {
		return rows, nil
	}
	return rows[len(rows)-n:], nil
}

// Print renders the DataFrame similar to pandas.DataFrame.to_string.
func (df DataFrame) Print(w io.Writer) error {
	rows, err := df.Rows()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	keys := orderedKeys(rows[0])
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, strings.Join(keys, "\t")); err != nil {
		return err
	}
	for _, row := range rows {
		values := make([]string, len(keys))
		for i, key := range keys {
			values[i] = fmt.Sprint(row[key])
		}
		if _, err := fmt.Fprintln(tw, strings.Join(values, "\t")); err != nil {
			return err
		}
	}
	return tw.Flush()
}

// ToCSV persists the DataFrame to disk, analogous to pandas.DataFrame.to_csv.
func (df DataFrame) ToCSV(path string) (err error) {
	rows, err := df.Rows()
	if err != nil {
		return err
	}

	cleanPath := filepath.Clean(path)
	file, err := os.Create(cleanPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	writer := stdcsv.NewWriter(file)
	if len(rows) == 0 {
		writer.Flush()
		return writer.Error()
	}

	keys := orderedKeys(rows[0])
	if err := writer.Write(keys); err != nil {
		return err
	}
	for _, row := range rows {
		record := make([]string, len(keys))
		for i, key := range keys {
			record[i] = fmt.Sprint(row[key])
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

// DebugSQL returns the lazily constructed SQL statement for inspection.
// It does not execute the statement and returns an empty string when planning previously failed.
func (df DataFrame) DebugSQL() string {
	if df.err != nil {
		return ""
	}
	state := df.buildQueryState()
	return compileSQL(state, nil)
}

// Columns returns the ordered column labels, mimicking pandas.DataFrame.columns.
func (df DataFrame) Columns() []string {
	rows, err := df.Head(1)
	if err != nil || len(rows) == 0 {
		return nil
	}
	return orderedKeys(rows[0])
}

// Shape returns (rows, cols) like pandas.DataFrame.shape, materializing the DataFrame if needed.
func (df DataFrame) Shape() (int, int) {
	rows, err := df.Rows()
	if err != nil {
		return 0, 0
	}
	count := len(rows)
	cols := 0
	if count > 0 {
		cols = len(rows[0])
	}
	return count, cols
}

// execute compiles operations into SQL, executes it, and returns result rows.
func (df DataFrame) execute(ctx context.Context, limit *int) ([]map[string]any, error) {
	if df.err != nil {
		return nil, df.err
	}
	state := df.buildQueryState()
	files := df.collectRequiredFiles()

	db, err := filesql.OpenContext(ctx, files...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	sqlStmt := compileSQL(state, limit)

	rows, err := db.QueryContext(ctx, sqlStmt)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	result, err := scanAll(rows)
	if err != nil {
		return nil, err
	}
	if err := transformRows(result, state); err != nil {
		return nil, err
	}
	return result, rows.Err()
}

// collectRequiredFiles gathers all files needed for execution.
func (df DataFrame) collectRequiredFiles() []string {
	files := map[string]struct{}{
		df.source: {},
	}

	for _, op := range df.ops {
		for _, f := range op.RequiredFiles() {
			files[f] = struct{}{}
		}
	}

	result := make([]string, 0, len(files))
	for f := range files {
		result = append(result, f)
	}
	return result
}

func (df DataFrame) withOp(op sqlOp) DataFrame {
	newOps := make([]sqlOp, len(df.ops), len(df.ops)+1)
	copy(newOps, df.ops)
	newOps = append(newOps, op)
	return DataFrame{
		source:       df.source,
		ops:          newOps,
		err:          df.err,
		warns:        df.warns,
		aliasCounter: df.aliasCounter,
	}
}

// Err surfaces deferred planning errors (for example, invalid Merge/Join options).
func (df DataFrame) Err() error {
	return df.err
}

// Warnings returns accumulated non-fatal warnings (missing rename/drop/cast/fill targets).
func (df DataFrame) Warnings() []string {
	if df.warns == nil || len(df.warns.list) == 0 {
		return nil
	}
	out := make([]string, len(df.warns.list))
	copy(out, df.warns.list)
	return out
}

// filterOp appends WHERE clause.
type filterOp struct {
	expr string
}

func (o filterOp) RequiredFiles() []string { return nil }
func (o filterOp) Apply(q *queryState) {
	q.filters = append(q.filters, o.expr)
}

// selectOp sets projection columns.
type selectOp struct {
	cols []string
}

func (o selectOp) RequiredFiles() []string { return nil }
func (o selectOp) Apply(q *queryState) {
	seen := map[string]struct{}{}
	cols := make([]string, 0, len(o.cols))
	for _, col := range o.cols {
		name := normalizeColumnName(col)
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		cols = append(cols, name)
	}
	q.selectCols = cols
}

// dropOp records columns to drop from the result.
type dropOp struct {
	cols []string
}

func (o dropOp) RequiredFiles() []string { return nil }
func (o dropOp) Apply(q *queryState) {
	if q.drops == nil {
		q.drops = map[string]struct{}{}
	}
	for _, c := range o.cols {
		name := normalizeColumnName(c)
		if dest, ok := q.renames[name]; ok {
			name = dest
		}
		if len(q.selectCols) > 0 && !containsString(q.selectCols, name) {
			// Column not present in current selection; ignore drop request, mirroring pandas warnings.
			q.warn(fmt.Sprintf("drop skipped for column %s: not selected", name))
			continue
		}
		q.drops[name] = struct{}{}
	}
	if len(q.selectCols) > 0 {
		next := make([]string, 0, len(q.selectCols))
		for _, col := range q.selectCols {
			if _, ok := q.drops[col]; ok {
				continue
			}
			next = append(next, col)
		}
		q.selectCols = next
	}
}

// renameOp records column rename mapping.
type renameOp struct {
	mapping map[string]string
}

func (o renameOp) RequiredFiles() []string { return nil }
func (o renameOp) Apply(q *queryState) {
	if q.renames == nil {
		q.renames = map[string]string{}
	}
	for k, v := range o.mapping {
		requestedOld := normalizeColumnName(k)
		newCol := normalizeColumnName(v)
		if q.drops != nil {
			if _, ok := q.drops[newCol]; ok {
				continue
			}
		}
		origOld := q.resolveOriginalColumn(requestedOld)
		currentName := q.resolveCurrentColumn(origOld)
		if len(q.selectCols) > 0 && !containsString(q.selectCols, currentName) {
			// pandas would raise KeyError; we surface the issue via warnings instead.
			q.warn(fmt.Sprintf("rename skipped for column %s: not selected", requestedOld))
			continue
		}
		if currentName == newCol {
			continue
		}
		q.renames[origOld] = newCol
		updated := false
		for i := range q.renameOrder {
			if q.renameOrder[i].old == origOld {
				q.renameOrder[i].new = newCol
				updated = true
				break
			}
		}
		if !updated {
			q.renameOrder = append(q.renameOrder, renameEntry{old: origOld, new: newCol})
		}
		for i, col := range q.selectCols {
			if col == currentName {
				q.selectCols[i] = newCol
			}
		}
	}
}

// mutateOp appends mutation definition.
type mutateOp struct {
	col  string
	expr string
}

func (o mutateOp) RequiredFiles() []string { return nil }
func (o mutateOp) Apply(q *queryState) {
	target := normalizeColumnName(q.resolveCurrentColumn(o.col))
	expr := rewriteExpressionWithOriginalNames(o.expr, q.renameOrder)
	q.mutations = append(q.mutations, mutation{
		Column: target,
		Expr:   expr,
	})
	if len(q.selectCols) > 0 {
		found := false
		for _, col := range q.selectCols {
			if col == target {
				found = true
				break
			}
		}
		if !found {
			q.selectCols = append(q.selectCols, target)
		}
	}
}

// sortOp appends ORDER BY clause.
type sortOp struct {
	col string
	asc bool
}

func (o sortOp) RequiredFiles() []string { return nil }
func (o sortOp) Apply(q *queryState) {
	col := normalizeColumnName(o.col)
	q.sorts = append(q.sorts, sortClause{
		Column: col,
		Asc:    o.asc,
	})
}

// castOp appends CAST clause.
type castOp struct {
	col   string
	dtype string
}

func (o castOp) RequiredFiles() []string { return nil }
func (o castOp) Apply(q *queryState) {
	col := q.resolveOriginalColumn(o.col)
	if !q.columnInSelection(col) && !q.columnInSelection(o.col) {
		// pandas would raise; we defer to warnings to keep pipelines flowing.
		q.warn(fmt.Sprintf("cast ignored for column %s: not selected", o.col))
		return
	}
	q.casts = append(q.casts, castClause{
		Column: col,
		DType:  o.dtype,
	})
}

// dropNaOp appends IS NOT NULL filters.
type dropNaOp struct {
	cols []string
}

func (o dropNaOp) RequiredFiles() []string { return nil }
func (o dropNaOp) Apply(q *queryState) {
	// NOTE: Multiple columns result in an AND condition (all columns must be non-NULL).
	for _, c := range o.cols {
		col := q.resolveOriginalColumn(c)
		q.filters = append(q.filters, col+" IS NOT NULL")
	}
}

// fillNaOp appends IFNULL expressions.
type fillNaOp struct {
	col   string
	value any
}

func (o fillNaOp) RequiredFiles() []string { return nil }
func (o fillNaOp) Apply(q *queryState) {
	col := q.resolveOriginalColumn(o.col)
	if !q.columnInSelection(col) && !q.columnInSelection(o.col) {
		// pandas would raise; we log the issue instead.
		q.warn(fmt.Sprintf("fillna ignored for column %s: not selected", o.col))
		return
	}
	q.fillNAs = append(q.fillNAs, fillNAClause{
		Column: col,
		Value:  o.value,
	})
}

// joinOp registers JOIN.
type joinOp struct {
	alias    string
	subquery string
	files    []string
	on       string
	joinType joinType
	suffixes [2]string
}

func (o joinOp) RequiredFiles() []string { return o.files }
func (o joinOp) Apply(q *queryState) {
	// NOTE: Later joins override any previously configured suffix preference; last applied join wins.
	if o.suffixes != ([2]string{}) {
		q.suffixes = o.suffixes
	}
	q.joins = append(q.joins, joinClause{
		OtherTable: o.alias,
		On:         o.on,
		Type:       o.joinType,
		SubQuery:   o.subquery,
		Suffixes:   o.suffixes,
	})
}

// compileSQL builds SQL from the aggregated queryState.
func compileSQL(q queryState, limit *int) string {
	selectExprs := buildSelectExpressions(q)

	var sb strings.Builder
	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(selectExprs, ", "))
	sb.WriteString(" FROM ")
	sb.WriteString(q.baseTable)

	for _, j := range q.joins {
		sb.WriteString(" ")
		sb.WriteString(joinKeyword(j.Type))
		if j.SubQuery != "" {
			sb.WriteString(" (")
			sb.WriteString(j.SubQuery)
			sb.WriteString(") AS ")
			sb.WriteString(j.OtherTable)
		} else {
			sb.WriteString(" ")
			sb.WriteString(j.OtherTable)
		}
		sb.WriteString(" ON ")
		sb.WriteString(joinCondition(q.baseTable, j))
	}

	if len(q.filters) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(q.filters, " AND "))
	}

	if len(q.sorts) > 0 {
		sb.WriteString(" ORDER BY ")
		for i, s := range q.sorts {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(resolveSortColumn(&q, s.Column))
			if !s.Asc {
				sb.WriteString(" DESC")
			}
		}
	}

	if limit != nil {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(*limit))
	}
	return sb.String()
}

func buildSelectExpressions(q queryState) []string {
	exprs := make([]string, 0)
	exprs = append(exprs, q.baseTable+".*")
	for _, j := range q.joins {
		exprs = append(exprs, j.OtherTable+".*")
	}

	for _, c := range q.casts {
		exprs = append(exprs, fmt.Sprintf("CAST(%s AS %s) AS %s", c.Column, c.DType, c.Column))
	}
	for _, f := range q.fillNAs {
		exprs = append(exprs, fmt.Sprintf("IFNULL(%s, %s) AS %s", f.Column, literal(f.Value), f.Column))
	}
	for _, m := range q.mutations {
		exprs = append(exprs, fmt.Sprintf("%s AS %s", m.Expr, m.Column))
	}
	var renameExprs []string
	if q.inlineRenames {
		for old, newName := range q.renames {
			renameExprs = append(renameExprs, fmt.Sprintf("%s AS %s", old, newName))
		}
	}
	exprs = append(exprs, renameExprs...)

	if len(exprs) == 0 {
		return []string{"*"}
	}
	return exprs
}

func joinKeyword(t joinType) string {
	switch t {
	case joinTypeInner:
		return "INNER JOIN"
	case joinTypeLeft:
		return "LEFT JOIN"
	case joinTypeRight:
		return "RIGHT JOIN"
	case joinTypeFull:
		return "FULL OUTER JOIN"
	default:
		return "INNER JOIN"
	}
}

func resolveSortColumn(q *queryState, name string) string {
	col := q.resolveOriginalColumn(name)
	if col == "" {
		return name
	}
	return col
}

func joinCondition(base string, j joinClause) string {
	cond := strings.TrimSpace(j.On)
	if cond == "" {
		return cond
	}
	if isExplicitJoinCondition(cond) {
		return cond
	}
	return fmt.Sprintf("%s.%s = %s.%s", base, cond, j.OtherTable, cond)
}

var (
	simpleJoinRegex   = regexp.MustCompile(`^\w+(\.\w+)?\s*=\s*\w+(\.\w+)?$`)
	duplicateKeyRegex = regexp.MustCompile(`^(.+)_([0-9]+)$`)
)
var identifierRegex = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)

func isExplicitJoinCondition(cond string) bool {
	return simpleJoinRegex.MatchString(cond)
}

func literal(v any) string {
	if v == nil {
		return "NULL"
	}
	switch val := v.(type) {
	case string:
		return "'" + strings.ReplaceAll(val, "'", "''") + "'"
	case []byte:
		return "'" + strings.ReplaceAll(string(val), "'", "''") + "'"
	case bool:
		if val {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprint(val)
	}
}

// scanAll scans all rows into []map[string]any using column names as keys.
func scanAll(rows *sql.Rows) ([]map[string]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	resolved := make([]string, len(cols))
	counts := make(map[string]int, len(cols))
	for i, col := range cols {
		if c := counts[col]; c > 0 {
			resolved[i] = fmt.Sprintf("%s_%d", col, c)
		} else {
			resolved[i] = col
		}
		counts[col]++
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(cols))
		for i, name := range resolved {
			switch v := values[i].(type) {
			case []byte:
				row[name] = string(v)
			default:
				row[name] = v
			}
		}
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func transformRows(rows []map[string]any, state queryState) error {
	applyJoinSuffixes(rows, state.joins)

	for _, row := range rows {
		if err := applyRename(row, state.renameOrder); err != nil {
			return err
		}
		if err := applyDrops(row, state.drops); err != nil {
			return err
		}
		if len(state.selectCols) > 0 {
			next := make(map[string]any, len(state.selectCols))
			for _, col := range state.selectCols {
				key, ok := resolveRowKey(row, col)
				if !ok {
					return fmt.Errorf("column %s not found", col)
				}
				next[key] = row[key]
			}
			for k := range row {
				delete(row, k)
			}
			for k, v := range next {
				row[k] = v
			}
		}
	}
	return nil
}

func applyJoinSuffixes(rows []map[string]any, joins []joinClause) {
	if len(joins) == 0 {
		return
	}
	for _, row := range rows {
		consumed := map[string]bool{}
		keys := make([]string, 0, len(row))
		for k := range row {
			keys = append(keys, k)
		}
		for _, key := range keys {
			base, idx, ok := parseDuplicateKey(key)
			if !ok {
				continue
			}
			joinIdx := idx - 1
			if joinIdx < 0 || joinIdx >= len(joins) {
				continue
			}
			joinVal, ok := row[key]
			if !ok {
				continue
			}
			suffixes := effectiveSuffixes(joins[joinIdx].Suffixes)
			leftName := base + suffixes[0]
			rightName := base + suffixes[1]
			if baseVal, ok := row[base]; ok {
				row[leftName] = baseVal
			}
			row[rightName] = joinVal
			delete(row, key)
			consumed[base] = true
		}
		for base := range consumed {
			delete(row, base)
		}
	}
}

func effectiveSuffixes(sfx [2]string) [2]string {
	if sfx == ([2]string{}) {
		return defaultSuffixes()
	}
	return sfx
}

func applyRename(row map[string]any, entries []renameEntry) error {
	for _, entry := range entries {
		if val, ok := row[entry.old]; ok {
			row[entry.new] = val
			delete(row, entry.old)
		}
	}
	return nil
}

func applyDrops(row map[string]any, drops map[string]struct{}) error {
	for col := range drops {
		if key, ok := resolveRowKey(row, col); ok {
			delete(row, key)
		}
	}
	return nil
}

func parseDuplicateKey(key string) (string, int, bool) {
	matches := duplicateKeyRegex.FindStringSubmatch(key)
	if len(matches) != 3 {
		return "", 0, false
	}
	num, err := strconv.Atoi(matches[2])
	if err != nil || num <= 0 || num > 10 {
		return "", 0, false
	}
	return matches[1], num, true
}

func tableNameFromPath(path string) string {
	base := filepath.Base(path)
	for {
		ext := filepath.Ext(base)
		if ext == "" {
			break
		}
		trimmed := strings.TrimSuffix(base, ext)
		base = trimmed
	}
	return base
}

func orderedKeys(row map[string]any) []string {
	keys := make([]string, 0, len(row))
	for k := range row {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func normalizeColumnName(name string) string {
	if idx := strings.LastIndex(name, "."); idx >= 0 && idx < len(name)-1 {
		prefix := name[:idx]
		if prefix == strings.ToLower(prefix) {
			return name[idx+1:]
		}
	}
	return name
}

// resolveRowKey searches for the provided name, favoring fully-qualified names,
// then bare column names, and finally suffix variants (e.g. *_x, *_y).
func resolveRowKey(row map[string]any, name string) (string, bool) {
	candidates := []string{name}
	if idx := strings.LastIndex(name, "."); idx >= 0 && idx < len(name)-1 {
		candidates = append(candidates, name[idx+1:])
	}
	seen := map[string]struct{}{}
	check := func(candidate string) (string, bool) {
		if candidate == "" {
			return "", false
		}
		if _, dup := seen[candidate]; dup {
			return "", false
		}
		seen[candidate] = struct{}{}
		if _, ok := row[candidate]; ok {
			return candidate, true
		}
		return "", false
	}
	for _, candidate := range candidates {
		if found, ok := check(candidate); ok {
			return found, true
		}
	}
	if strings.Contains(name, ".") || strings.HasSuffix(name, "_x") || strings.HasSuffix(name, "_y") {
		suffixes := defaultSuffixes()
		for _, base := range candidates {
			for _, suffix := range suffixes {
				if found, ok := check(base + suffix); ok {
					return found, true
				}
			}
		}
	}
	return "", false
}

func containsString(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func rewriteExpressionWithOriginalNames(expr string, entries []renameEntry) string {
	result := expr
	for _, entry := range entries {
		if entry.new == entry.old || entry.new == "" {
			continue
		}
		result = identifierRegex.ReplaceAllStringFunc(result, func(token string) string {
			if token == entry.new {
				return entry.old
			}
			return token
		})
	}
	return result
}
func (q *queryState) columnInSelection(name string) bool {
	if len(q.selectCols) == 0 {
		return true
	}
	if containsString(q.selectCols, name) {
		return true
	}
	current := q.resolveCurrentColumn(name)
	if containsString(q.selectCols, current) {
		return true
	}
	original := q.resolveOriginalColumn(name)
	return containsString(q.selectCols, original)
}

func (q *queryState) warn(msg string) {
	if q.warnings != nil {
		q.warnings.add(msg)
	}
}

func (df DataFrame) nextAlias(base string) string {
	if df.aliasCounter == nil {
		counter := 0
		df.aliasCounter = &counter
	}
	*df.aliasCounter++
	return fmt.Sprintf("%s_alias_%d", base, *df.aliasCounter)
}
