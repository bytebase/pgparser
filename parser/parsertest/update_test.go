package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestUpdateBasic tests: UPDATE t SET name = 'new' WHERE id = 1
func TestUpdateBasic(t *testing.T) {
	input := "UPDATE t SET name = 'new' WHERE id = 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "t" {
		t.Errorf("expected table name 't', got %q", stmt.Relation.Relname)
	}

	// Verify TargetList (SET clause) has 1 item
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 item in TargetList, got %v", stmt.TargetList)
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
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

	// Verify WHERE clause is set
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set")
	}

	aexpr, ok := stmt.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", stmt.WhereClause)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}

	// No FROM, RETURNING, or WITH
	if stmt.FromClause != nil && len(stmt.FromClause.Items) > 0 {
		t.Error("expected no FromClause")
	}
	if stmt.ReturningList != nil && len(stmt.ReturningList.Items) > 0 {
		t.Error("expected no ReturningList")
	}
	if stmt.WithClause != nil {
		t.Error("expected no WithClause")
	}
}

// TestUpdateMultipleColumns tests: UPDATE t SET name = 'new', age = 30 WHERE id = 1
func TestUpdateMultipleColumns(t *testing.T) {
	input := "UPDATE t SET name = 'new', age = 30 WHERE id = 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
	}

	// Verify TargetList has 2 items
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 2 {
		t.Fatalf("expected 2 items in TargetList, got %v", stmt.TargetList)
	}

	// First SET target: name = 'new'
	rt1, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget for first target, got %T", stmt.TargetList.Items[0])
	}
	if rt1.Name != "name" {
		t.Errorf("expected first SET target 'name', got %q", rt1.Name)
	}

	// Second SET target: age = 30
	rt2, ok := stmt.TargetList.Items[1].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget for second target, got %T", stmt.TargetList.Items[1])
	}
	if rt2.Name != "age" {
		t.Errorf("expected second SET target 'age', got %q", rt2.Name)
	}

	intConst, ok := rt2.Val.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for age value, got %T", rt2.Val)
	}
	intVal, ok := intConst.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", intConst.Val)
	}
	if intVal.Ival != 30 {
		t.Errorf("expected 30, got %d", intVal.Ival)
	}

	// Verify WHERE clause is set
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set")
	}
}

// TestUpdateWithFrom tests: UPDATE t SET name = t2.name FROM t2 WHERE t.id = t2.id
func TestUpdateWithFrom(t *testing.T) {
	input := "UPDATE t SET name = t2.name FROM t2 WHERE t.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify TargetList
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 item in TargetList, got %v", stmt.TargetList)
	}

	// Verify FROM clause is set
	if stmt.FromClause == nil || len(stmt.FromClause.Items) != 1 {
		t.Fatalf("expected 1 table in FromClause, got %v", stmt.FromClause)
	}

	rv, ok := stmt.FromClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", stmt.FromClause.Items[0])
	}
	if rv.Relname != "t2" {
		t.Errorf("expected table 't2', got %q", rv.Relname)
	}

	// Verify WHERE clause is set
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set")
	}
}

// TestUpdateReturning tests: UPDATE t SET name = 'new' RETURNING id, name
func TestUpdateReturning(t *testing.T) {
	input := "UPDATE t SET name = 'new' RETURNING id, name"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
	}

	// Verify RETURNING clause has 2 items
	if stmt.ReturningList == nil || len(stmt.ReturningList.Items) != 2 {
		t.Fatalf("expected 2 items in ReturningList, got %v", stmt.ReturningList)
	}

	// First returning item: id
	rt1, ok := stmt.ReturningList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.ReturningList.Items[0])
	}
	colRef1, ok := rt1.Val.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef, got %T", rt1.Val)
	}
	str1, ok := colRef1.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef1.Fields.Items[0])
	}
	if str1.Str != "id" {
		t.Errorf("expected 'id', got %q", str1.Str)
	}

	// Second returning item: name
	rt2, ok := stmt.ReturningList.Items[1].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.ReturningList.Items[1])
	}
	colRef2, ok := rt2.Val.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef, got %T", rt2.Val)
	}
	str2, ok := colRef2.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef2.Fields.Items[0])
	}
	if str2.Str != "name" {
		t.Errorf("expected 'name', got %q", str2.Str)
	}
}

// TestUpdateReturningStar tests: UPDATE t SET name = 'new' RETURNING *
func TestUpdateReturningStar(t *testing.T) {
	input := "UPDATE t SET name = 'new' RETURNING *"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
	}

	// Verify RETURNING clause has 1 item (star)
	if stmt.ReturningList == nil || len(stmt.ReturningList.Items) != 1 {
		t.Fatalf("expected 1 item in ReturningList, got %v", stmt.ReturningList)
	}

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

// TestUpdateWithAlias tests: UPDATE t AS target SET name = 'new'
func TestUpdateWithAlias(t *testing.T) {
	input := "UPDATE t AS target SET name = 'new'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
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

// TestUpdateWithCTE tests: WITH cte AS (SELECT 1) UPDATE t SET name = 'new'
func TestUpdateWithCTE(t *testing.T) {
	input := "WITH cte AS (SELECT 1) UPDATE t SET name = 'new'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
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
}

// TestUpdateNoWhere tests: UPDATE t SET active = false
func TestUpdateNoWhere(t *testing.T) {
	input := "UPDATE t SET active = false"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify TargetList has 1 item
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 item in TargetList, got %v", stmt.TargetList)
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}
	if rt.Name != "active" {
		t.Errorf("expected SET target 'active', got %q", rt.Name)
	}

	// Verify the value is false
	aConst, ok := rt.Val.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", rt.Val)
	}
	boolVal, ok := aConst.Val.(*nodes.Boolean)
	if !ok {
		t.Fatalf("expected *nodes.Boolean, got %T", aConst.Val)
	}
	if boolVal.Boolval != false {
		t.Errorf("expected false, got %v", boolVal.Boolval)
	}

	// Verify WHERE clause is nil
	if stmt.WhereClause != nil {
		t.Error("expected WhereClause to be nil")
	}
}

// TestUpdateTableExpressions is a table-driven test for various UPDATE patterns.
func TestUpdateTableExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic UPDATE", "UPDATE t SET name = 'new' WHERE id = 1"},
		{"multiple columns", "UPDATE t SET name = 'new', age = 30 WHERE id = 1"},
		{"with FROM", "UPDATE t SET name = t2.name FROM t2 WHERE t.id = t2.id"},
		{"RETURNING multiple", "UPDATE t SET name = 'new' RETURNING id, name"},
		{"RETURNING *", "UPDATE t SET name = 'new' RETURNING *"},
		{"with AS alias", "UPDATE t AS target SET name = 'new'"},
		{"WITH UPDATE", "WITH cte AS (SELECT 1) UPDATE t SET name = 'new'"},
		{"no WHERE", "UPDATE t SET active = false"},
		{"schema-qualified", "UPDATE myschema.t SET name = 'new'"},
		{"FROM + WHERE + RETURNING", "UPDATE t SET name = 'new' FROM t2 WHERE t.id = t2.id RETURNING *"},
		{"multiple FROM tables", "UPDATE t SET name = t2.name FROM t2, t3 WHERE t.id = t2.id AND t2.id = t3.id"},
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

			_, ok := result.Items[0].(*nodes.UpdateStmt)
			if !ok {
				t.Fatalf("expected *nodes.UpdateStmt for %q, got %T", tt.input, result.Items[0])
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}
