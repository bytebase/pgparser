package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

func TestParseScalarSubquery(t *testing.T) {
	// SELECT (SELECT 1) -- scalar subquery as expression
	input := "SELECT (SELECT 1)"

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

	// Target list should have 1 item
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target, got %d", len(stmt.TargetList.Items))
	}

	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}

	// The expression should be a SubLink with EXPR_SUBLINK
	sublink, ok := rt.Val.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink, got %T", rt.Val)
	}

	if sublink.SubLinkType != int(nodes.EXPR_SUBLINK) {
		t.Errorf("expected EXPR_SUBLINK (%d), got %d", int(nodes.EXPR_SUBLINK), sublink.SubLinkType)
	}

	// Subselect should be a SelectStmt
	innerStmt, ok := sublink.Subselect.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt in Subselect, got %T", sublink.Subselect)
	}

	// Inner select should be "SELECT 1"
	if innerStmt.TargetList == nil || len(innerStmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target in inner select")
	}

	innerRT := innerStmt.TargetList.Items[0].(*nodes.ResTarget)
	aConst, ok := innerRT.Val.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const in inner select, got %T", innerRT.Val)
	}
	intVal, ok := aConst.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", aConst.Val)
	}
	if intVal.Ival != 1 {
		t.Errorf("expected inner select value 1, got %d", intVal.Ival)
	}

	// Testexpr should be nil for scalar subquery
	if sublink.Testexpr != nil {
		t.Error("expected Testexpr to be nil for scalar subquery")
	}
}

func TestParseExistsSubquery(t *testing.T) {
	input := "SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// WHERE clause should be a SubLink with EXISTS_SUBLINK
	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	if sublink.SubLinkType != int(nodes.EXISTS_SUBLINK) {
		t.Errorf("expected EXISTS_SUBLINK (%d), got %d", int(nodes.EXISTS_SUBLINK), sublink.SubLinkType)
	}

	// Subselect should be a SelectStmt
	innerStmt, ok := sublink.Subselect.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt in Subselect, got %T", sublink.Subselect)
	}

	// Inner select should have FROM t2
	if innerStmt.FromClause == nil || len(innerStmt.FromClause.Items) != 1 {
		t.Fatal("expected 1 table in inner select FROM clause")
	}
	rv, ok := innerStmt.FromClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", innerStmt.FromClause.Items[0])
	}
	if rv.Relname != "t2" {
		t.Errorf("expected inner FROM table 't2', got %q", rv.Relname)
	}

	// Testexpr should be nil for EXISTS
	if sublink.Testexpr != nil {
		t.Error("expected Testexpr to be nil for EXISTS subquery")
	}
}

func TestParseInSubquery(t *testing.T) {
	input := "SELECT * FROM t WHERE id IN (SELECT id FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	// IN subquery becomes ANY_SUBLINK
	if sublink.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Errorf("expected ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sublink.SubLinkType)
	}

	// Testexpr should be the column ref "id"
	if sublink.Testexpr == nil {
		t.Fatal("expected Testexpr to be set for IN subquery")
	}
	colRef, ok := sublink.Testexpr.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for Testexpr, got %T", sublink.Testexpr)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "id" {
		t.Errorf("expected Testexpr column 'id', got %q", str.Str)
	}

	// Subselect should be present
	if sublink.Subselect == nil {
		t.Fatal("expected Subselect to be set")
	}
	innerStmt, ok := sublink.Subselect.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt in Subselect, got %T", sublink.Subselect)
	}

	// Inner select should have FROM t2
	if innerStmt.FromClause == nil || len(innerStmt.FromClause.Items) != 1 {
		t.Fatal("expected 1 table in inner select FROM clause")
	}
	rv := innerStmt.FromClause.Items[0].(*nodes.RangeVar)
	if rv.Relname != "t2" {
		t.Errorf("expected inner FROM table 't2', got %q", rv.Relname)
	}
}

func TestParseNotInSubquery(t *testing.T) {
	input := "SELECT * FROM t WHERE id NOT IN (SELECT id FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	// NOT IN subquery becomes BoolExpr(NOT) wrapping SubLink(ANY)
	boolExpr, ok := stmt.WhereClause.(*nodes.BoolExpr)
	if !ok {
		t.Fatalf("expected *nodes.BoolExpr in where clause, got %T", stmt.WhereClause)
	}

	if boolExpr.Boolop != nodes.NOT_EXPR {
		t.Errorf("expected NOT_EXPR, got %d", boolExpr.Boolop)
	}

	if boolExpr.Args == nil || len(boolExpr.Args.Items) != 1 {
		t.Fatalf("expected 1 arg in BoolExpr, got %d", len(boolExpr.Args.Items))
	}

	sublink, ok := boolExpr.Args.Items[0].(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink wrapped by NOT, got %T", boolExpr.Args.Items[0])
	}

	if sublink.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Errorf("expected ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sublink.SubLinkType)
	}

	// Testexpr should be the column ref "id"
	if sublink.Testexpr == nil {
		t.Fatal("expected Testexpr to be set")
	}
	colRef, ok := sublink.Testexpr.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for Testexpr, got %T", sublink.Testexpr)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "id" {
		t.Errorf("expected Testexpr column 'id', got %q", str.Str)
	}

	// Subselect should be present
	if sublink.Subselect == nil {
		t.Fatal("expected Subselect to be set")
	}
}

func TestParseEqualAnySubquery(t *testing.T) {
	input := "SELECT * FROM t WHERE id = ANY (SELECT id FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	if sublink.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Errorf("expected ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sublink.SubLinkType)
	}

	// Testexpr should be the column ref "id"
	if sublink.Testexpr == nil {
		t.Fatal("expected Testexpr to be set")
	}
	colRef, ok := sublink.Testexpr.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for Testexpr, got %T", sublink.Testexpr)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "id" {
		t.Errorf("expected Testexpr column 'id', got %q", str.Str)
	}

	// OperName should be "="
	if sublink.OperName == nil || len(sublink.OperName.Items) != 1 {
		t.Fatalf("expected 1 item in OperName, got %v", sublink.OperName)
	}
	opStr, ok := sublink.OperName.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String in OperName, got %T", sublink.OperName.Items[0])
	}
	if opStr.Str != "=" {
		t.Errorf("expected OperName '=', got %q", opStr.Str)
	}

	// Subselect should be present
	if sublink.Subselect == nil {
		t.Fatal("expected Subselect to be set")
	}
}

func TestParseEqualAllSubquery(t *testing.T) {
	input := "SELECT * FROM t WHERE id = ALL (SELECT id FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	if sublink.SubLinkType != int(nodes.ALL_SUBLINK) {
		t.Errorf("expected ALL_SUBLINK (%d), got %d", int(nodes.ALL_SUBLINK), sublink.SubLinkType)
	}

	// Testexpr should be the column ref "id"
	if sublink.Testexpr == nil {
		t.Fatal("expected Testexpr to be set")
	}

	// OperName should be "="
	if sublink.OperName == nil || len(sublink.OperName.Items) != 1 {
		t.Fatalf("expected 1 item in OperName")
	}
	opStr := sublink.OperName.Items[0].(*nodes.String)
	if opStr.Str != "=" {
		t.Errorf("expected OperName '=', got %q", opStr.Str)
	}

	// Subselect should be present
	if sublink.Subselect == nil {
		t.Fatal("expected Subselect to be set")
	}
}

func TestParseInSubqueryWithLimit(t *testing.T) {
	input := "SELECT * FROM t WHERE id IN (SELECT id FROM t2 LIMIT 5)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	if sublink.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Errorf("expected ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sublink.SubLinkType)
	}

	// Subselect should be a SelectStmt with LIMIT 5
	innerStmt, ok := sublink.Subselect.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt in Subselect, got %T", sublink.Subselect)
	}

	if innerStmt.LimitCount == nil {
		t.Fatal("expected inner select to have LimitCount")
	}
	ac, ok := innerStmt.LimitCount.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for LimitCount, got %T", innerStmt.LimitCount)
	}
	iv, ok := ac.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", ac.Val)
	}
	if iv.Ival != 5 {
		t.Errorf("expected inner LIMIT 5, got %d", iv.Ival)
	}
}

func TestParseGreaterThanSomeSubquery(t *testing.T) {
	// SOME is a synonym for ANY
	input := "SELECT * FROM t WHERE id > SOME (SELECT id FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	// SOME is the same as ANY
	if sublink.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Errorf("expected ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sublink.SubLinkType)
	}

	// OperName should be ">"
	if sublink.OperName == nil || len(sublink.OperName.Items) != 1 {
		t.Fatal("expected OperName to be set")
	}
	opStr := sublink.OperName.Items[0].(*nodes.String)
	if opStr.Str != ">" {
		t.Errorf("expected OperName '>', got %q", opStr.Str)
	}
}

func TestParseLessThanAllSubquery(t *testing.T) {
	input := "SELECT * FROM t WHERE id < ALL (SELECT id FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	if sublink.SubLinkType != int(nodes.ALL_SUBLINK) {
		t.Errorf("expected ALL_SUBLINK (%d), got %d", int(nodes.ALL_SUBLINK), sublink.SubLinkType)
	}

	// OperName should be "<"
	if sublink.OperName == nil || len(sublink.OperName.Items) != 1 {
		t.Fatal("expected OperName to be set")
	}
	opStr := sublink.OperName.Items[0].(*nodes.String)
	if opStr.Str != "<" {
		t.Errorf("expected OperName '<', got %q", opStr.Str)
	}
}

func TestParseNotEqualsAnySubquery(t *testing.T) {
	input := "SELECT * FROM t WHERE id <> ANY (SELECT id FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	sublink, ok := stmt.WhereClause.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink in where clause, got %T", stmt.WhereClause)
	}

	if sublink.SubLinkType != int(nodes.ANY_SUBLINK) {
		t.Errorf("expected ANY_SUBLINK (%d), got %d", int(nodes.ANY_SUBLINK), sublink.SubLinkType)
	}

	// OperName should be "<>"
	if sublink.OperName == nil || len(sublink.OperName.Items) != 1 {
		t.Fatal("expected OperName to be set")
	}
	opStr := sublink.OperName.Items[0].(*nodes.String)
	if opStr.Str != "<>" {
		t.Errorf("expected OperName '<>', got %q", opStr.Str)
	}
}

func TestParseExistsWithAnd(t *testing.T) {
	// Test EXISTS combined with another condition via AND
	input := "SELECT * FROM t WHERE x = 1 AND EXISTS (SELECT 1 FROM t2)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	// Should be a BoolExpr(AND) with two args
	boolExpr, ok := stmt.WhereClause.(*nodes.BoolExpr)
	if !ok {
		t.Fatalf("expected *nodes.BoolExpr, got %T", stmt.WhereClause)
	}

	if boolExpr.Boolop != nodes.AND_EXPR {
		t.Errorf("expected AND_EXPR, got %d", boolExpr.Boolop)
	}

	if boolExpr.Args == nil || len(boolExpr.Args.Items) != 2 {
		t.Fatalf("expected 2 args in AND, got %d", len(boolExpr.Args.Items))
	}

	// Second arg should be EXISTS SubLink
	sublink, ok := boolExpr.Args.Items[1].(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink as second AND arg, got %T", boolExpr.Args.Items[1])
	}

	if sublink.SubLinkType != int(nodes.EXISTS_SUBLINK) {
		t.Errorf("expected EXISTS_SUBLINK (%d), got %d", int(nodes.EXISTS_SUBLINK), sublink.SubLinkType)
	}
}

func TestParseScalarSubqueryInExpression(t *testing.T) {
	// Test scalar subquery used in an expression context
	input := "SELECT (SELECT 1) + 2"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target")
	}

	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)

	// Should be an A_Expr (SubLink + 2)
	aexpr, ok := rt.Val.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", rt.Val)
	}

	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}

	// Left side should be a SubLink
	sublink, ok := aexpr.Lexpr.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink on left side, got %T", aexpr.Lexpr)
	}

	if sublink.SubLinkType != int(nodes.EXPR_SUBLINK) {
		t.Errorf("expected EXPR_SUBLINK (%d), got %d", int(nodes.EXPR_SUBLINK), sublink.SubLinkType)
	}

	// Right side should be integer 2
	aConst, ok := aexpr.Rexpr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const on right side, got %T", aexpr.Rexpr)
	}
	iv := aConst.Val.(*nodes.Integer)
	if iv.Ival != 2 {
		t.Errorf("expected right-side value 2, got %d", iv.Ival)
	}
}
