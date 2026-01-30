package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestDeleteBasic tests: DELETE FROM t WHERE id = 1
func TestDeleteBasic(t *testing.T) {
	input := "DELETE FROM t WHERE id = 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "t" {
		t.Errorf("expected table name 't', got %q", stmt.Relation.Relname)
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

	// No USING, RETURNING, or WITH
	if stmt.UsingClause != nil && len(stmt.UsingClause.Items) > 0 {
		t.Error("expected no UsingClause")
	}
	if stmt.ReturningList != nil && len(stmt.ReturningList.Items) > 0 {
		t.Error("expected no ReturningList")
	}
	if stmt.WithClause != nil {
		t.Error("expected no WithClause")
	}
}

// TestDeleteAllRows tests: DELETE FROM t (no WHERE clause)
func TestDeleteAllRows(t *testing.T) {
	input := "DELETE FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "t" {
		t.Errorf("expected table name 't', got %q", stmt.Relation.Relname)
	}

	// Verify WHERE clause is nil
	if stmt.WhereClause != nil {
		t.Error("expected WhereClause to be nil")
	}
}

// TestDeleteWithUsing tests: DELETE FROM t USING t2 WHERE t.id = t2.id
func TestDeleteWithUsing(t *testing.T) {
	input := "DELETE FROM t USING t2 WHERE t.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}

	// Verify USING clause is set
	if stmt.UsingClause == nil || len(stmt.UsingClause.Items) != 1 {
		t.Fatalf("expected 1 table in UsingClause, got %v", stmt.UsingClause)
	}

	rv, ok := stmt.UsingClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", stmt.UsingClause.Items[0])
	}
	if rv.Relname != "t2" {
		t.Errorf("expected table 't2', got %q", rv.Relname)
	}

	// Verify WHERE clause is set
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set")
	}
}

// TestDeleteReturning tests: DELETE FROM t WHERE id = 1 RETURNING id
func TestDeleteReturning(t *testing.T) {
	input := "DELETE FROM t WHERE id = 1 RETURNING id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
	}

	// Verify RETURNING clause has 1 item
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
	str, ok := colRef.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef.Fields.Items[0])
	}
	if str.Str != "id" {
		t.Errorf("expected 'id', got %q", str.Str)
	}
}

// TestDeleteReturningStar tests: DELETE FROM t RETURNING *
func TestDeleteReturningStar(t *testing.T) {
	input := "DELETE FROM t RETURNING *"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
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

// TestDeleteWithAlias tests: DELETE FROM t AS target WHERE target.id = 1
func TestDeleteWithAlias(t *testing.T) {
	input := "DELETE FROM t AS target WHERE target.id = 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
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

// TestDeleteWithCTE tests: WITH cte AS (SELECT 1) DELETE FROM t WHERE id IN (SELECT * FROM cte)
func TestDeleteWithCTE(t *testing.T) {
	input := "WITH cte AS (SELECT 1) DELETE FROM t WHERE id IN (SELECT * FROM cte)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
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

	// Verify WHERE clause is set (IN subquery)
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set")
	}
}

// TestDeleteSchemaQualified tests: DELETE FROM myschema.t WHERE id = 1
func TestDeleteSchemaQualified(t *testing.T) {
	input := "DELETE FROM myschema.t WHERE id = 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", result.Items[0])
	}

	// Verify relation with schema
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

// TestDeleteTableExpressions is a table-driven test for various DELETE patterns.
func TestDeleteTableExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic DELETE", "DELETE FROM t WHERE id = 1"},
		{"DELETE all rows", "DELETE FROM t"},
		{"DELETE with USING", "DELETE FROM t USING t2 WHERE t.id = t2.id"},
		{"DELETE RETURNING id", "DELETE FROM t WHERE id = 1 RETURNING id"},
		{"DELETE RETURNING *", "DELETE FROM t RETURNING *"},
		{"DELETE with alias", "DELETE FROM t AS target WHERE target.id = 1"},
		{"WITH DELETE", "WITH cte AS (SELECT 1) DELETE FROM t WHERE id IN (SELECT * FROM cte)"},
		{"DELETE schema-qualified", "DELETE FROM myschema.t WHERE id = 1"},
		{"DELETE USING multiple tables", "DELETE FROM t USING t2, t3 WHERE t.id = t2.id AND t2.id = t3.id"},
		{"DELETE USING + RETURNING", "DELETE FROM t USING t2 WHERE t.id = t2.id RETURNING *"},
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

			_, ok := result.Items[0].(*nodes.DeleteStmt)
			if !ok {
				t.Fatalf("expected *nodes.DeleteStmt for %q, got %T", tt.input, result.Items[0])
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}
