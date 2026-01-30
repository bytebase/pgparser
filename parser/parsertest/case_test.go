package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestCaseSearched tests a searched CASE (no argument after CASE keyword).
// SELECT CASE WHEN x = 1 THEN 'one' ELSE 'other' END FROM t
func TestCaseSearched(t *testing.T) {
	input := "SELECT CASE WHEN x = 1 THEN 'one' ELSE 'other' END FROM t"

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

	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	ce, ok := rt.Val.(*nodes.CaseExpr)
	if !ok {
		t.Fatalf("expected *nodes.CaseExpr, got %T", rt.Val)
	}

	// Searched CASE: Arg should be nil
	if ce.Arg != nil {
		t.Errorf("expected CaseExpr.Arg to be nil for searched CASE, got %T", ce.Arg)
	}

	// Should have 1 WHEN clause
	if ce.Args == nil || len(ce.Args.Items) != 1 {
		t.Fatalf("expected 1 CaseWhen, got %d", listLen(ce.Args))
	}

	cw, ok := ce.Args.Items[0].(*nodes.CaseWhen)
	if !ok {
		t.Fatalf("expected *nodes.CaseWhen, got %T", ce.Args.Items[0])
	}

	// The condition should be an A_Expr (x = 1)
	_, ok = cw.Expr.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr for WHEN condition, got %T", cw.Expr)
	}

	// The result should be a string constant 'one'
	aConst, ok := cw.Result.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for THEN result, got %T", cw.Result)
	}
	strVal, ok := aConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", aConst.Val)
	}
	if strVal.Str != "one" {
		t.Errorf("expected THEN result 'one', got %q", strVal.Str)
	}

	// Defresult should be non-nil (ELSE clause)
	if ce.Defresult == nil {
		t.Fatal("expected non-nil Defresult for ELSE clause")
	}
	defConst, ok := ce.Defresult.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for ELSE, got %T", ce.Defresult)
	}
	defStr, ok := defConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String for ELSE value, got %T", defConst.Val)
	}
	if defStr.Str != "other" {
		t.Errorf("expected ELSE result 'other', got %q", defStr.Str)
	}
}

// TestCaseSimple tests a simple CASE with an argument.
// SELECT CASE status WHEN 1 THEN 'active' WHEN 2 THEN 'inactive' END FROM t
func TestCaseSimple(t *testing.T) {
	input := "SELECT CASE status WHEN 1 THEN 'active' WHEN 2 THEN 'inactive' END FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	ce, ok := rt.Val.(*nodes.CaseExpr)
	if !ok {
		t.Fatalf("expected *nodes.CaseExpr, got %T", rt.Val)
	}

	// Simple CASE: Arg should be a ColumnRef for "status"
	if ce.Arg == nil {
		t.Fatal("expected non-nil CaseExpr.Arg for simple CASE")
	}
	colRef, ok := ce.Arg.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for Arg, got %T", ce.Arg)
	}
	if colRef.Fields == nil || len(colRef.Fields.Items) != 1 {
		t.Fatal("expected 1 field in ColumnRef")
	}
	str, ok := colRef.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef.Fields.Items[0])
	}
	if str.Str != "status" {
		t.Errorf("expected Arg column 'status', got %q", str.Str)
	}

	// Should have 2 WHEN clauses
	if ce.Args == nil || len(ce.Args.Items) != 2 {
		t.Fatalf("expected 2 CaseWhen entries, got %d", listLen(ce.Args))
	}

	// Check first WHEN: WHEN 1 THEN 'active'
	cw1 := ce.Args.Items[0].(*nodes.CaseWhen)
	c1, ok := cw1.Expr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for first WHEN expr, got %T", cw1.Expr)
	}
	iv1, ok := c1.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", c1.Val)
	}
	if iv1.Ival != 1 {
		t.Errorf("expected first WHEN value 1, got %d", iv1.Ival)
	}

	// Check second WHEN: WHEN 2 THEN 'inactive'
	cw2 := ce.Args.Items[1].(*nodes.CaseWhen)
	c2, ok := cw2.Expr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for second WHEN expr, got %T", cw2.Expr)
	}
	iv2, ok := c2.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", c2.Val)
	}
	if iv2.Ival != 2 {
		t.Errorf("expected second WHEN value 2, got %d", iv2.Ival)
	}

	// Defresult should be nil (no ELSE clause)
	if ce.Defresult != nil {
		t.Errorf("expected nil Defresult (no ELSE clause), got %T", ce.Defresult)
	}
}

// TestCaseMultipleWhen tests CASE with multiple WHEN clauses and ELSE.
// SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END FROM t
func TestCaseMultipleWhen(t *testing.T) {
	input := "SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	ce, ok := rt.Val.(*nodes.CaseExpr)
	if !ok {
		t.Fatalf("expected *nodes.CaseExpr, got %T", rt.Val)
	}

	// Should have 2 WHEN clauses
	if ce.Args == nil || len(ce.Args.Items) != 2 {
		t.Fatalf("expected 2 CaseWhen entries, got %d", listLen(ce.Args))
	}

	// Check first WHEN condition: a > 0
	cw1 := ce.Args.Items[0].(*nodes.CaseWhen)
	aexpr1, ok := cw1.Expr.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr for first WHEN, got %T", cw1.Expr)
	}
	if aexpr1.Name == nil || len(aexpr1.Name.Items) == 0 {
		t.Fatal("expected operator name in A_Expr")
	}
	opStr1 := aexpr1.Name.Items[0].(*nodes.String).Str
	if opStr1 != ">" {
		t.Errorf("expected '>' operator, got %q", opStr1)
	}

	// Check second WHEN condition: a < 0
	cw2 := ce.Args.Items[1].(*nodes.CaseWhen)
	aexpr2, ok := cw2.Expr.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr for second WHEN, got %T", cw2.Expr)
	}
	opStr2 := aexpr2.Name.Items[0].(*nodes.String).Str
	if opStr2 != "<" {
		t.Errorf("expected '<' operator, got %q", opStr2)
	}

	// Defresult should be non-nil
	if ce.Defresult == nil {
		t.Fatal("expected non-nil Defresult for ELSE clause")
	}
	defConst, ok := ce.Defresult.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for ELSE, got %T", ce.Defresult)
	}
	defStr, ok := defConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String for ELSE value, got %T", defConst.Val)
	}
	if defStr.Str != "zero" {
		t.Errorf("expected ELSE result 'zero', got %q", defStr.Str)
	}
}

// TestCaseWithoutElse tests CASE without an ELSE clause.
// SELECT CASE WHEN x = 1 THEN 'one' END FROM t
func TestCaseWithoutElse(t *testing.T) {
	input := "SELECT CASE WHEN x = 1 THEN 'one' END FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	ce, ok := rt.Val.(*nodes.CaseExpr)
	if !ok {
		t.Fatalf("expected *nodes.CaseExpr, got %T", rt.Val)
	}

	// Should have 1 WHEN clause
	if ce.Args == nil || len(ce.Args.Items) != 1 {
		t.Fatalf("expected 1 CaseWhen, got %d", listLen(ce.Args))
	}

	// Defresult should be nil
	if ce.Defresult != nil {
		t.Errorf("expected nil Defresult (no ELSE), got %T", ce.Defresult)
	}
}

// TestCaseNested tests nested CASE expressions.
// SELECT CASE WHEN x = 1 THEN CASE WHEN y = 1 THEN 'a' ELSE 'b' END ELSE 'c' END FROM t
func TestCaseNested(t *testing.T) {
	input := "SELECT CASE WHEN x = 1 THEN CASE WHEN y = 1 THEN 'a' ELSE 'b' END ELSE 'c' END FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	outerCE, ok := rt.Val.(*nodes.CaseExpr)
	if !ok {
		t.Fatalf("expected *nodes.CaseExpr, got %T", rt.Val)
	}

	// Outer CASE should have 1 WHEN clause
	if outerCE.Args == nil || len(outerCE.Args.Items) != 1 {
		t.Fatalf("expected 1 outer CaseWhen, got %d", listLen(outerCE.Args))
	}

	outerCW := outerCE.Args.Items[0].(*nodes.CaseWhen)

	// The Result of the outer WHEN should itself be a CaseExpr (the nested CASE)
	innerCE, ok := outerCW.Result.(*nodes.CaseExpr)
	if !ok {
		t.Fatalf("expected nested *nodes.CaseExpr as THEN result, got %T", outerCW.Result)
	}

	// Inner CASE should have 1 WHEN clause
	if innerCE.Args == nil || len(innerCE.Args.Items) != 1 {
		t.Fatalf("expected 1 inner CaseWhen, got %d", listLen(innerCE.Args))
	}

	// Inner CASE should have ELSE 'b'
	if innerCE.Defresult == nil {
		t.Fatal("expected non-nil Defresult for inner CASE")
	}
	innerDefConst, ok := innerCE.Defresult.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for inner ELSE, got %T", innerCE.Defresult)
	}
	innerDefStr, ok := innerDefConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", innerDefConst.Val)
	}
	if innerDefStr.Str != "b" {
		t.Errorf("expected inner ELSE 'b', got %q", innerDefStr.Str)
	}

	// Outer CASE should have ELSE 'c'
	if outerCE.Defresult == nil {
		t.Fatal("expected non-nil Defresult for outer CASE")
	}
	outerDefConst, ok := outerCE.Defresult.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for outer ELSE, got %T", outerCE.Defresult)
	}
	outerDefStr, ok := outerDefConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", outerDefConst.Val)
	}
	if outerDefStr.Str != "c" {
		t.Errorf("expected outer ELSE 'c', got %q", outerDefStr.Str)
	}
}

// TestCaseInWhere tests CASE expression used in a WHERE clause.
// SELECT * FROM t WHERE CASE WHEN x > 0 THEN true ELSE false END
func TestCaseInWhere(t *testing.T) {
	input := "SELECT * FROM t WHERE CASE WHEN x > 0 THEN true ELSE false END"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// WhereClause should be a CaseExpr
	if stmt.WhereClause == nil {
		t.Fatal("expected non-nil WhereClause")
	}

	ce, ok := stmt.WhereClause.(*nodes.CaseExpr)
	if !ok {
		t.Fatalf("expected *nodes.CaseExpr in WHERE clause, got %T", stmt.WhereClause)
	}

	// Searched CASE: Arg should be nil
	if ce.Arg != nil {
		t.Errorf("expected nil Arg for searched CASE in WHERE, got %T", ce.Arg)
	}

	// Should have 1 WHEN clause
	if ce.Args == nil || len(ce.Args.Items) != 1 {
		t.Fatalf("expected 1 CaseWhen, got %d", listLen(ce.Args))
	}

	cw := ce.Args.Items[0].(*nodes.CaseWhen)

	// THEN result should be boolean true
	thenConst, ok := cw.Result.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for THEN true, got %T", cw.Result)
	}
	boolVal, ok := thenConst.Val.(*nodes.Boolean)
	if !ok {
		t.Fatalf("expected *nodes.Boolean, got %T", thenConst.Val)
	}
	if !boolVal.Boolval {
		t.Error("expected THEN result to be true")
	}

	// ELSE should be boolean false
	if ce.Defresult == nil {
		t.Fatal("expected non-nil Defresult")
	}
	elseConst, ok := ce.Defresult.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for ELSE false, got %T", ce.Defresult)
	}
	elseBool, ok := elseConst.Val.(*nodes.Boolean)
	if !ok {
		t.Fatalf("expected *nodes.Boolean, got %T", elseConst.Val)
	}
	if elseBool.Boolval {
		t.Error("expected ELSE result to be false")
	}
}

// listLen returns the number of items in a list, or 0 if nil.
func listLen(l *nodes.List) int {
	if l == nil {
		return 0
	}
	return len(l.Items)
}
