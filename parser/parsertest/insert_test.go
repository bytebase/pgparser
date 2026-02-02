package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestInsertBasicValues tests: INSERT INTO t VALUES (1, 'hello')
func TestInsertBasicValues(t *testing.T) {
	input := "INSERT INTO t VALUES (1, 'hello')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "t" {
		t.Errorf("expected table name 't', got %q", stmt.Relation.Relname)
	}

	// Verify no explicit columns
	if stmt.Cols != nil && len(stmt.Cols.Items) > 0 {
		t.Error("expected no explicit column list")
	}

	// Verify SelectStmt (VALUES clause) exists
	if stmt.SelectStmt == nil {
		t.Fatal("expected SelectStmt to be set for VALUES clause")
	}

	selectStmt, ok := stmt.SelectStmt.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.SelectStmt)
	}

	// Verify ValuesLists has 1 row
	if selectStmt.ValuesLists == nil || len(selectStmt.ValuesLists.Items) != 1 {
		t.Fatalf("expected 1 VALUES row, got %v", selectStmt.ValuesLists)
	}

	// Verify the row has 2 values
	row, ok := selectStmt.ValuesLists.Items[0].(*nodes.List)
	if !ok {
		t.Fatalf("expected VALUES row to be *nodes.List, got %T", selectStmt.ValuesLists.Items[0])
	}
	if len(row.Items) != 2 {
		t.Fatalf("expected 2 values in row, got %d", len(row.Items))
	}

	// First value should be integer 1
	val1, ok := row.Items[0].(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for first value, got %T", row.Items[0])
	}
	intVal, ok := val1.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", val1.Val)
	}
	if intVal.Ival != 1 {
		t.Errorf("expected integer 1, got %d", intVal.Ival)
	}

	// Second value should be string 'hello'
	val2, ok := row.Items[1].(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for second value, got %T", row.Items[1])
	}
	strVal, ok := val2.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", val2.Val)
	}
	if strVal.Str != "hello" {
		t.Errorf("expected string 'hello', got %q", strVal.Str)
	}

	// No ON CONFLICT, RETURNING, or WITH
	if stmt.OnConflictClause != nil {
		t.Error("expected no OnConflictClause")
	}
	if stmt.ReturningList != nil {
		t.Error("expected no ReturningList")
	}
	if stmt.WithClause != nil {
		t.Error("expected no WithClause")
	}
}

// TestInsertWithColumns tests: INSERT INTO t (id, name) VALUES (1, 'hello')
func TestInsertWithColumns(t *testing.T) {
	input := "INSERT INTO t (id, name) VALUES (1, 'hello')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatalf("expected table name 't'")
	}

	// Verify column list
	if stmt.Cols == nil || len(stmt.Cols.Items) != 2 {
		t.Fatalf("expected 2 columns, got %v", stmt.Cols)
	}

	col1, ok := stmt.Cols.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget for col1, got %T", stmt.Cols.Items[0])
	}
	if col1.Name != "id" {
		t.Errorf("expected column name 'id', got %q", col1.Name)
	}

	col2, ok := stmt.Cols.Items[1].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget for col2, got %T", stmt.Cols.Items[1])
	}
	if col2.Name != "name" {
		t.Errorf("expected column name 'name', got %q", col2.Name)
	}

	// Verify SelectStmt exists (VALUES)
	if stmt.SelectStmt == nil {
		t.Fatal("expected SelectStmt to be set")
	}
}

// TestInsertSelect tests: INSERT INTO t SELECT * FROM t2
func TestInsertSelect(t *testing.T) {
	input := "INSERT INTO t SELECT * FROM t2"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify no columns
	if stmt.Cols != nil && len(stmt.Cols.Items) > 0 {
		t.Error("expected no explicit column list for INSERT...SELECT")
	}

	// Verify the SelectStmt is a SELECT (not VALUES)
	selectStmt, ok := stmt.SelectStmt.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.SelectStmt)
	}

	// The inner SELECT should have a target list and FROM clause
	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) != 1 {
		t.Fatal("expected 1 target in SELECT")
	}

	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) != 1 {
		t.Fatal("expected 1 table in FROM clause")
	}

	rv, ok := selectStmt.FromClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", selectStmt.FromClause.Items[0])
	}
	if rv.Relname != "t2" {
		t.Errorf("expected table 't2', got %q", rv.Relname)
	}
}

// TestInsertDefaultValues tests: INSERT INTO t DEFAULT VALUES
func TestInsertDefaultValues(t *testing.T) {
	input := "INSERT INTO t DEFAULT VALUES"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// DEFAULT VALUES means nil SelectStmt
	if stmt.SelectStmt != nil {
		t.Error("expected nil SelectStmt for DEFAULT VALUES")
	}

	// No columns
	if stmt.Cols != nil && len(stmt.Cols.Items) > 0 {
		t.Error("expected no column list")
	}
}

// TestInsertReturning tests: INSERT INTO t VALUES (1) RETURNING id
func TestInsertReturning(t *testing.T) {
	input := "INSERT INTO t VALUES (1) RETURNING id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify RETURNING clause
	if stmt.ReturningList == nil || len(stmt.ReturningList.Items) != 1 {
		t.Fatalf("expected 1 item in RETURNING list, got %v", stmt.ReturningList)
	}

	rt, ok := stmt.ReturningList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.ReturningList.Items[0])
	}

	colRef, ok := rt.Val.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef, got %T", rt.Val)
	}

	str, ok := colRef.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef.Fields.Items[0])
	}
	if str.Str != "id" {
		t.Errorf("expected 'id', got %q", str.Str)
	}
}

// TestInsertOnConflictDoNothing tests: INSERT INTO t VALUES (1) ON CONFLICT DO NOTHING
func TestInsertOnConflictDoNothing(t *testing.T) {
	input := "INSERT INTO t VALUES (1) ON CONFLICT DO NOTHING"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify ON CONFLICT clause
	if stmt.OnConflictClause == nil {
		t.Fatal("expected OnConflictClause to be set")
	}

	if stmt.OnConflictClause.Action != parser.ONCONFLICT_NOTHING {
		t.Errorf("expected parser.ONCONFLICT_NOTHING (%d), got %d", parser.ONCONFLICT_NOTHING, stmt.OnConflictClause.Action)
	}

	// No infer clause for bare DO NOTHING
	if stmt.OnConflictClause.Infer != nil {
		t.Error("expected no Infer clause for bare ON CONFLICT DO NOTHING")
	}
}

// TestInsertOnConflictDoUpdate tests: INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'new'
func TestInsertOnConflictDoUpdate(t *testing.T) {
	input := "INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'new'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify ON CONFLICT clause
	if stmt.OnConflictClause == nil {
		t.Fatal("expected OnConflictClause to be set")
	}

	if stmt.OnConflictClause.Action != parser.ONCONFLICT_UPDATE {
		t.Errorf("expected parser.ONCONFLICT_UPDATE (%d), got %d", parser.ONCONFLICT_UPDATE, stmt.OnConflictClause.Action)
	}

	// Verify SET clause
	if stmt.OnConflictClause.TargetList == nil || len(stmt.OnConflictClause.TargetList.Items) != 1 {
		t.Fatalf("expected 1 SET clause item, got %v", stmt.OnConflictClause.TargetList)
	}

	rt, ok := stmt.OnConflictClause.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.OnConflictClause.TargetList.Items[0])
	}
	if rt.Name != "name" {
		t.Errorf("expected SET target 'name', got %q", rt.Name)
	}

	// Verify SET value is string 'new'
	aConst, ok := rt.Val.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", rt.Val)
	}
	strVal, ok := aConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", aConst.Val)
	}
	if strVal.Str != "new" {
		t.Errorf("expected 'new', got %q", strVal.Str)
	}
}

// TestInsertMultipleValuesRows tests: INSERT INTO t VALUES (1, 'a'), (2, 'b'), (3, 'c')
func TestInsertMultipleValuesRows(t *testing.T) {
	input := "INSERT INTO t VALUES (1, 'a'), (2, 'b'), (3, 'c')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify SelectStmt is a VALUES clause with 3 rows
	selectStmt, ok := stmt.SelectStmt.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.SelectStmt)
	}

	if selectStmt.ValuesLists == nil || len(selectStmt.ValuesLists.Items) != 3 {
		t.Fatalf("expected 3 VALUES rows, got %v", selectStmt.ValuesLists)
	}

	// Verify each row has 2 values
	for i, item := range selectStmt.ValuesLists.Items {
		row, ok := item.(*nodes.List)
		if !ok {
			t.Fatalf("expected VALUES row %d to be *nodes.List, got %T", i, item)
		}
		if len(row.Items) != 2 {
			t.Errorf("expected 2 values in row %d, got %d", i, len(row.Items))
		}
	}

	// Verify first row values are (1, 'a')
	row0 := selectStmt.ValuesLists.Items[0].(*nodes.List)
	val0, ok := row0.Items[0].(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", row0.Items[0])
	}
	intVal, ok := val0.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", val0.Val)
	}
	if intVal.Ival != 1 {
		t.Errorf("expected 1, got %d", intVal.Ival)
	}
}

// TestInsertWithCTE tests: WITH cte AS (SELECT 1) INSERT INTO t SELECT * FROM cte
func TestInsertWithCTE(t *testing.T) {
	input := "WITH cte AS (SELECT 1) INSERT INTO t SELECT * FROM cte"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify WITH clause
	if stmt.WithClause == nil {
		t.Fatal("expected WithClause to be set")
	}

	if stmt.WithClause.Ctes == nil || len(stmt.WithClause.Ctes.Items) != 1 {
		t.Fatalf("expected 1 CTE, got %v", stmt.WithClause.Ctes)
	}

	cte, ok := stmt.WithClause.Ctes.Items[0].(*nodes.CommonTableExpr)
	if !ok {
		t.Fatalf("expected *nodes.CommonTableExpr, got %T", stmt.WithClause.Ctes.Items[0])
	}
	if cte.Ctename != "cte" {
		t.Errorf("expected CTE name 'cte', got %q", cte.Ctename)
	}

	// Verify inner SELECT
	selectStmt, ok := stmt.SelectStmt.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.SelectStmt)
	}
	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) != 1 {
		t.Fatal("expected 1 table in FROM clause")
	}
	rv, ok := selectStmt.FromClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", selectStmt.FromClause.Items[0])
	}
	if rv.Relname != "cte" {
		t.Errorf("expected table 'cte', got %q", rv.Relname)
	}
}

// TestInsertWithAlias tests: INSERT INTO t AS target VALUES (1)
func TestInsertWithAlias(t *testing.T) {
	input := "INSERT INTO t AS target VALUES (1)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify alias
	if stmt.Relation.Alias == nil {
		t.Fatal("expected Alias to be set")
	}
	if stmt.Relation.Alias.Aliasname != "target" {
		t.Errorf("expected alias 'target', got %q", stmt.Relation.Alias.Aliasname)
	}
}

// TestInsertOnConflictWithInfer tests: INSERT INTO t VALUES (1) ON CONFLICT (id) DO NOTHING
func TestInsertOnConflictWithInfer(t *testing.T) {
	input := "INSERT INTO t VALUES (1) ON CONFLICT (id) DO NOTHING"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	if stmt.OnConflictClause == nil {
		t.Fatal("expected OnConflictClause to be set")
	}

	if stmt.OnConflictClause.Action != parser.ONCONFLICT_NOTHING {
		t.Errorf("expected parser.ONCONFLICT_NOTHING, got %d", stmt.OnConflictClause.Action)
	}

	// Verify infer clause
	if stmt.OnConflictClause.Infer == nil {
		t.Fatal("expected Infer clause to be set")
	}

	if stmt.OnConflictClause.Infer.IndexElems == nil || len(stmt.OnConflictClause.Infer.IndexElems.Items) != 1 {
		t.Fatalf("expected 1 index element, got %v", stmt.OnConflictClause.Infer.IndexElems)
	}

	elem, ok := stmt.OnConflictClause.Infer.IndexElems.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem, got %T", stmt.OnConflictClause.Infer.IndexElems.Items[0])
	}
	if elem.Name != "id" {
		t.Errorf("expected 'id', got %q", elem.Name)
	}
}

// TestInsertOnConflictDoUpdateWithWhere tests:
// INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'new' WHERE id > 0
func TestInsertOnConflictDoUpdateWithWhere(t *testing.T) {
	input := "INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'new' WHERE id > 0"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	if stmt.OnConflictClause == nil {
		t.Fatal("expected OnConflictClause to be set")
	}

	if stmt.OnConflictClause.Action != parser.ONCONFLICT_UPDATE {
		t.Errorf("expected parser.ONCONFLICT_UPDATE, got %d", stmt.OnConflictClause.Action)
	}

	// Verify WHERE clause on ON CONFLICT
	if stmt.OnConflictClause.WhereClause == nil {
		t.Fatal("expected WhereClause in OnConflictClause")
	}

	aexpr, ok := stmt.OnConflictClause.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", stmt.OnConflictClause.WhereClause)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
}

// TestInsertReturningMultiple tests: INSERT INTO t VALUES (1, 'a') RETURNING id, name
func TestInsertReturningMultiple(t *testing.T) {
	input := "INSERT INTO t VALUES (1, 'a') RETURNING id, name"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	if stmt.ReturningList == nil || len(stmt.ReturningList.Items) != 2 {
		t.Fatalf("expected 2 items in RETURNING list, got %v", stmt.ReturningList)
	}
}

// TestInsertReturningAll tests: INSERT INTO t VALUES (1) RETURNING *
func TestInsertReturningAll(t *testing.T) {
	input := "INSERT INTO t VALUES (1) RETURNING *"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	if stmt.ReturningList == nil || len(stmt.ReturningList.Items) != 1 {
		t.Fatalf("expected 1 item in RETURNING list, got %v", stmt.ReturningList)
	}

	// Should be a ResTarget with * (ColumnRef with A_Star)
	rt, ok := stmt.ReturningList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.ReturningList.Items[0])
	}

	colRef, ok := rt.Val.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef, got %T", rt.Val)
	}

	if len(colRef.Fields.Items) != 1 {
		t.Fatalf("expected 1 field, got %d", len(colRef.Fields.Items))
	}

	_, ok = colRef.Fields.Items[0].(*nodes.A_Star)
	if !ok {
		t.Fatalf("expected *nodes.A_Star, got %T", colRef.Fields.Items[0])
	}
}

// TestInsertSchemaQualifiedTable tests: INSERT INTO myschema.t VALUES (1)
func TestInsertSchemaQualifiedTable(t *testing.T) {
	input := "INSERT INTO myschema.t VALUES (1)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", result.Items[0])
	}

	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Schemaname != "myschema" {
		t.Errorf("expected schema 'myschema', got %q", stmt.Relation.Schemaname)
	}
	if stmt.Relation.Relname != "t" {
		t.Errorf("expected table 't', got %q", stmt.Relation.Relname)
	}
}

// TestInsertTableExpressions is a table-driven test for various INSERT patterns.
func TestInsertTableExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic VALUES", "INSERT INTO t VALUES (1, 'hello')"},
		{"with columns", "INSERT INTO t (id, name) VALUES (1, 'hello')"},
		{"INSERT SELECT", "INSERT INTO t SELECT * FROM t2"},
		{"DEFAULT VALUES", "INSERT INTO t DEFAULT VALUES"},
		{"RETURNING", "INSERT INTO t VALUES (1) RETURNING id"},
		{"RETURNING *", "INSERT INTO t VALUES (1) RETURNING *"},
		{"ON CONFLICT DO NOTHING", "INSERT INTO t VALUES (1) ON CONFLICT DO NOTHING"},
		{"ON CONFLICT DO UPDATE", "INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'new'"},
		{"ON CONFLICT with infer", "INSERT INTO t VALUES (1) ON CONFLICT (id) DO NOTHING"},
		{"ON CONFLICT infer + UPDATE", "INSERT INTO t VALUES (1) ON CONFLICT (id) DO UPDATE SET name = 'new'"},
		{"ON CONFLICT + WHERE", "INSERT INTO t VALUES (1) ON CONFLICT DO UPDATE SET name = 'x' WHERE id > 0"},
		{"multiple VALUES rows", "INSERT INTO t VALUES (1, 'a'), (2, 'b'), (3, 'c')"},
		{"WITH INSERT", "WITH cte AS (SELECT 1) INSERT INTO t SELECT * FROM cte"},
		{"INSERT with alias", "INSERT INTO t AS target VALUES (1)"},
		{"schema-qualified table", "INSERT INTO myschema.t VALUES (1)"},
		{"RETURNING multiple cols", "INSERT INTO t VALUES (1, 'a') RETURNING id, name"},
		{"columns + RETURNING", "INSERT INTO t (id) VALUES (1) RETURNING *"},
		{"ON CONFLICT ON CONSTRAINT", "INSERT INTO t VALUES (1) ON CONFLICT ON CONSTRAINT t_pkey DO NOTHING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			if result == nil || len(result.Items) == 0 {
				t.Fatalf("Parse returned empty result for %q", tt.input)
			}

			_, ok := result.Items[0].(*nodes.InsertStmt)
			if !ok {
				t.Fatalf("expected *nodes.InsertStmt for %q, got %T", tt.input, result.Items[0])
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}
