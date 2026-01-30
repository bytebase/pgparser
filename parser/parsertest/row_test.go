package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

func TestParseExplicitRow(t *testing.T) {
	// ROW(1, 2, 3) should produce RowExpr with COERCE_EXPLICIT_CALL, 3 elements
	input := "SELECT ROW(1, 2, 3)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target, got %d", len(stmt.TargetList.Items))
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}

	rowExpr, ok := rt.Val.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr, got %T", rt.Val)
	}

	if rowExpr.RowFormat != nodes.COERCE_EXPLICIT_CALL {
		t.Errorf("expected COERCE_EXPLICIT_CALL (%d), got %d", nodes.COERCE_EXPLICIT_CALL, rowExpr.RowFormat)
	}

	if rowExpr.Args == nil || len(rowExpr.Args.Items) != 3 {
		t.Fatalf("expected 3 elements in ROW, got %v", rowExpr.Args)
	}

	// Verify elements are integer constants 1, 2, 3
	for i, expected := range []int64{1, 2, 3} {
		ac, ok := rowExpr.Args.Items[i].(*nodes.A_Const)
		if !ok {
			t.Fatalf("element %d: expected *nodes.A_Const, got %T", i, rowExpr.Args.Items[i])
		}
		iv, ok := ac.Val.(*nodes.Integer)
		if !ok {
			t.Fatalf("element %d: expected *nodes.Integer, got %T", i, ac.Val)
		}
		if iv.Ival != expected {
			t.Errorf("element %d: expected %d, got %d", i, expected, iv.Ival)
		}
	}
}

func TestParseEmptyRow(t *testing.T) {
	// ROW() should produce RowExpr with COERCE_EXPLICIT_CALL and nil/empty Args
	input := "SELECT ROW()"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target")
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}

	rowExpr, ok := rt.Val.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr, got %T", rt.Val)
	}

	if rowExpr.RowFormat != nodes.COERCE_EXPLICIT_CALL {
		t.Errorf("expected COERCE_EXPLICIT_CALL (%d), got %d", nodes.COERCE_EXPLICIT_CALL, rowExpr.RowFormat)
	}

	if rowExpr.Args != nil && len(rowExpr.Args.Items) != 0 {
		t.Errorf("expected nil or empty Args for ROW(), got %d elements", len(rowExpr.Args.Items))
	}
}

func TestParseImplicitRow(t *testing.T) {
	// (1, 2, 3) should produce RowExpr with COERCE_IMPLICIT_CAST, 3 elements
	input := "SELECT (1, 2, 3)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target, got %d", len(stmt.TargetList.Items))
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}

	rowExpr, ok := rt.Val.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr, got %T", rt.Val)
	}

	if rowExpr.RowFormat != nodes.COERCE_IMPLICIT_CAST {
		t.Errorf("expected COERCE_IMPLICIT_CAST (%d), got %d", nodes.COERCE_IMPLICIT_CAST, rowExpr.RowFormat)
	}

	if rowExpr.Args == nil || len(rowExpr.Args.Items) != 3 {
		t.Fatalf("expected 3 elements in implicit row, got %v", rowExpr.Args)
	}

	// Verify elements are integer constants 1, 2, 3
	for i, expected := range []int64{1, 2, 3} {
		ac, ok := rowExpr.Args.Items[i].(*nodes.A_Const)
		if !ok {
			t.Fatalf("element %d: expected *nodes.A_Const, got %T", i, rowExpr.Args.Items[i])
		}
		iv, ok := ac.Val.(*nodes.Integer)
		if !ok {
			t.Fatalf("element %d: expected *nodes.Integer, got %T", i, ac.Val)
		}
		if iv.Ival != expected {
			t.Errorf("element %d: expected %d, got %d", i, expected, iv.Ival)
		}
	}
}

func TestParseImplicitRowTwoElements(t *testing.T) {
	// (1, 2) should produce RowExpr with COERCE_IMPLICIT_CAST, 2 elements
	// This is the minimum for an implicit row (distinguishes from parenthesized expr)
	input := "SELECT (1, 2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target")
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}

	rowExpr, ok := rt.Val.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr, got %T", rt.Val)
	}

	if rowExpr.RowFormat != nodes.COERCE_IMPLICIT_CAST {
		t.Errorf("expected COERCE_IMPLICIT_CAST (%d), got %d", nodes.COERCE_IMPLICIT_CAST, rowExpr.RowFormat)
	}

	if rowExpr.Args == nil || len(rowExpr.Args.Items) != 2 {
		t.Fatalf("expected 2 elements in implicit row, got %v", rowExpr.Args)
	}

	// Verify elements are integer constants 1, 2
	for i, expected := range []int64{1, 2} {
		ac, ok := rowExpr.Args.Items[i].(*nodes.A_Const)
		if !ok {
			t.Fatalf("element %d: expected *nodes.A_Const, got %T", i, rowExpr.Args.Items[i])
		}
		iv, ok := ac.Val.(*nodes.Integer)
		if !ok {
			t.Fatalf("element %d: expected *nodes.Integer, got %T", i, ac.Val)
		}
		if iv.Ival != expected {
			t.Errorf("element %d: expected %d, got %d", i, expected, iv.Ival)
		}
	}
}

func TestParseRowComparison(t *testing.T) {
	// ROW(a, b) = ROW(1, 2) should produce A_Expr with RowExpr operands
	input := "SELECT * FROM t WHERE ROW(a, b) = ROW(1, 2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	// WHERE clause should be A_Expr with "=" operator
	aexpr, ok := stmt.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr in where clause, got %T", stmt.WhereClause)
	}

	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}

	// Left side should be RowExpr with COERCE_EXPLICIT_CALL
	leftRow, ok := aexpr.Lexpr.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr on left side, got %T", aexpr.Lexpr)
	}
	if leftRow.RowFormat != nodes.COERCE_EXPLICIT_CALL {
		t.Errorf("left: expected COERCE_EXPLICIT_CALL, got %d", leftRow.RowFormat)
	}
	if leftRow.Args == nil || len(leftRow.Args.Items) != 2 {
		t.Fatalf("left: expected 2 elements, got %v", leftRow.Args)
	}

	// Right side should be RowExpr with COERCE_EXPLICIT_CALL
	rightRow, ok := aexpr.Rexpr.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr on right side, got %T", aexpr.Rexpr)
	}
	if rightRow.RowFormat != nodes.COERCE_EXPLICIT_CALL {
		t.Errorf("right: expected COERCE_EXPLICIT_CALL, got %d", rightRow.RowFormat)
	}
	if rightRow.Args == nil || len(rightRow.Args.Items) != 2 {
		t.Fatalf("right: expected 2 elements, got %v", rightRow.Args)
	}
}

func TestParseImplicitRowComparison(t *testing.T) {
	// (x, y) = (1, 2) should use implicit row on both sides
	input := "SELECT * FROM t WHERE (x, y) = (1, 2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	// WHERE clause should be A_Expr with "=" operator
	aexpr, ok := stmt.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr in where clause, got %T", stmt.WhereClause)
	}

	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}

	// Left side should be RowExpr with COERCE_IMPLICIT_CAST
	leftRow, ok := aexpr.Lexpr.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr on left side, got %T", aexpr.Lexpr)
	}
	if leftRow.RowFormat != nodes.COERCE_IMPLICIT_CAST {
		t.Errorf("left: expected COERCE_IMPLICIT_CAST, got %d", leftRow.RowFormat)
	}
	if leftRow.Args == nil || len(leftRow.Args.Items) != 2 {
		t.Fatalf("left: expected 2 elements, got %v", leftRow.Args)
	}

	// Right side should be RowExpr with COERCE_IMPLICIT_CAST
	rightRow, ok := aexpr.Rexpr.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr on right side, got %T", aexpr.Rexpr)
	}
	if rightRow.RowFormat != nodes.COERCE_IMPLICIT_CAST {
		t.Errorf("right: expected COERCE_IMPLICIT_CAST, got %d", rightRow.RowFormat)
	}
	if rightRow.Args == nil || len(rightRow.Args.Items) != 2 {
		t.Fatalf("right: expected 2 elements, got %v", rightRow.Args)
	}
}

func TestParseExplicitRowSingleElement(t *testing.T) {
	// ROW(42) should produce RowExpr with COERCE_EXPLICIT_CALL, 1 element
	input := "SELECT ROW(42)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target")
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}

	rowExpr, ok := rt.Val.(*nodes.RowExpr)
	if !ok {
		t.Fatalf("expected *nodes.RowExpr, got %T", rt.Val)
	}

	if rowExpr.RowFormat != nodes.COERCE_EXPLICIT_CALL {
		t.Errorf("expected COERCE_EXPLICIT_CALL (%d), got %d", nodes.COERCE_EXPLICIT_CALL, rowExpr.RowFormat)
	}

	if rowExpr.Args == nil || len(rowExpr.Args.Items) != 1 {
		t.Fatalf("expected 1 element in ROW, got %v", rowExpr.Args)
	}

	ac, ok := rowExpr.Args.Items[0].(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", rowExpr.Args.Items[0])
	}
	iv, ok := ac.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", ac.Val)
	}
	if iv.Ival != 42 {
		t.Errorf("expected 42, got %d", iv.Ival)
	}
}

func TestParseSingleParenExprNotRow(t *testing.T) {
	// (42) should NOT be a RowExpr - it should be a simple parenthesized expression
	input := "SELECT (42)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target")
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}

	// Single parenthesized expression should NOT be a RowExpr
	_, isRow := rt.Val.(*nodes.RowExpr)
	if isRow {
		t.Error("expected (42) to NOT be a RowExpr - single parenthesized expressions are not rows")
	}

	// It should be a plain A_Const
	ac, ok := rt.Val.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", rt.Val)
	}
	iv, ok := ac.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", ac.Val)
	}
	if iv.Ival != 42 {
		t.Errorf("expected 42, got %d", iv.Ival)
	}
}
