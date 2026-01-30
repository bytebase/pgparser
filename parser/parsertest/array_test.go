package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

func TestParseArrayConstructor(t *testing.T) {
	// ARRAY[1, 2, 3] should produce A_ArrayExpr with 3 elements
	input := "SELECT ARRAY[1, 2, 3]"

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

	arrExpr, ok := rt.Val.(*nodes.A_ArrayExpr)
	if !ok {
		t.Fatalf("expected *nodes.A_ArrayExpr, got %T", rt.Val)
	}

	if arrExpr.Elements == nil || len(arrExpr.Elements.Items) != 3 {
		t.Fatalf("expected 3 elements in ARRAY, got %v", arrExpr.Elements)
	}

	// Verify elements are integer constants 1, 2, 3
	for i, expected := range []int64{1, 2, 3} {
		ac, ok := arrExpr.Elements.Items[i].(*nodes.A_Const)
		if !ok {
			t.Fatalf("element %d: expected *nodes.A_Const, got %T", i, arrExpr.Elements.Items[i])
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

func TestParseArraySubquery(t *testing.T) {
	// ARRAY(SELECT id FROM t) should produce SubLink with ARRAY_SUBLINK
	input := "SELECT ARRAY(SELECT id FROM t)"

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

	sublink, ok := rt.Val.(*nodes.SubLink)
	if !ok {
		t.Fatalf("expected *nodes.SubLink, got %T", rt.Val)
	}

	if sublink.SubLinkType != int(nodes.ARRAY_SUBLINK) {
		t.Errorf("expected ARRAY_SUBLINK (%d), got %d", int(nodes.ARRAY_SUBLINK), sublink.SubLinkType)
	}

	// Subselect should be a SelectStmt
	innerStmt, ok := sublink.Subselect.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt in Subselect, got %T", sublink.Subselect)
	}

	// Inner select should have FROM t
	if innerStmt.FromClause == nil || len(innerStmt.FromClause.Items) != 1 {
		t.Fatal("expected 1 table in inner select FROM clause")
	}
	rv, ok := innerStmt.FromClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", innerStmt.FromClause.Items[0])
	}
	if rv.Relname != "t" {
		t.Errorf("expected inner FROM table 't', got %q", rv.Relname)
	}
}

func TestParseEmptyArray(t *testing.T) {
	// ARRAY[] should produce A_ArrayExpr with nil/empty Elements
	input := "SELECT ARRAY[]"

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

	arrExpr, ok := rt.Val.(*nodes.A_ArrayExpr)
	if !ok {
		t.Fatalf("expected *nodes.A_ArrayExpr, got %T", rt.Val)
	}

	if arrExpr.Elements != nil && len(arrExpr.Elements.Items) != 0 {
		t.Errorf("expected nil or empty Elements for ARRAY[], got %d elements", len(arrExpr.Elements.Items))
	}
}

func TestParseArraySubscript(t *testing.T) {
	// arr[1] should produce A_Indirection with A_Indices(Uidx=1)
	input := "SELECT arr[1] FROM t"

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

	indir, ok := rt.Val.(*nodes.A_Indirection)
	if !ok {
		t.Fatalf("expected *nodes.A_Indirection, got %T", rt.Val)
	}

	// Arg should be a ColumnRef for "arr"
	colRef, ok := indir.Arg.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for Arg, got %T", indir.Arg)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "arr" {
		t.Errorf("expected column 'arr', got %q", str.Str)
	}

	// Indirection should have 1 A_Indices element
	if indir.Indirection == nil || len(indir.Indirection.Items) != 1 {
		t.Fatalf("expected 1 indirection element, got %v", indir.Indirection)
	}

	indices, ok := indir.Indirection.Items[0].(*nodes.A_Indices)
	if !ok {
		t.Fatalf("expected *nodes.A_Indices, got %T", indir.Indirection.Items[0])
	}

	if indices.IsSlice {
		t.Error("expected IsSlice=false for subscript")
	}

	// Uidx should be integer constant 1
	ac, ok := indices.Uidx.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for Uidx, got %T", indices.Uidx)
	}
	iv, ok := ac.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", ac.Val)
	}
	if iv.Ival != 1 {
		t.Errorf("expected Uidx=1, got %d", iv.Ival)
	}
}

func TestParseArraySlice(t *testing.T) {
	// arr[1:3] should produce A_Indirection with A_Indices(IsSlice=true, Lidx=1, Uidx=3)
	input := "SELECT arr[1:3] FROM t"

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

	indir, ok := rt.Val.(*nodes.A_Indirection)
	if !ok {
		t.Fatalf("expected *nodes.A_Indirection, got %T", rt.Val)
	}

	// Arg should be a ColumnRef for "arr"
	colRef, ok := indir.Arg.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for Arg, got %T", indir.Arg)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "arr" {
		t.Errorf("expected column 'arr', got %q", str.Str)
	}

	// Indirection should have 1 A_Indices element
	if indir.Indirection == nil || len(indir.Indirection.Items) != 1 {
		t.Fatalf("expected 1 indirection element, got %v", indir.Indirection)
	}

	indices, ok := indir.Indirection.Items[0].(*nodes.A_Indices)
	if !ok {
		t.Fatalf("expected *nodes.A_Indices, got %T", indir.Indirection.Items[0])
	}

	if !indices.IsSlice {
		t.Error("expected IsSlice=true for slice")
	}

	// Lidx should be integer constant 1
	acLidx, ok := indices.Lidx.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for Lidx, got %T", indices.Lidx)
	}
	ivLidx, ok := acLidx.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer for Lidx, got %T", acLidx.Val)
	}
	if ivLidx.Ival != 1 {
		t.Errorf("expected Lidx=1, got %d", ivLidx.Ival)
	}

	// Uidx should be integer constant 3
	acUidx, ok := indices.Uidx.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for Uidx, got %T", indices.Uidx)
	}
	ivUidx, ok := acUidx.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer for Uidx, got %T", acUidx.Val)
	}
	if ivUidx.Ival != 3 {
		t.Errorf("expected Uidx=3, got %d", ivUidx.Ival)
	}
}

func TestParseNestedArray(t *testing.T) {
	// ARRAY[[1,2],[3,4]] should produce A_ArrayExpr with nested A_ArrayExpr elements
	input := "SELECT ARRAY[[1,2],[3,4]]"

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

	outerArr, ok := rt.Val.(*nodes.A_ArrayExpr)
	if !ok {
		t.Fatalf("expected *nodes.A_ArrayExpr, got %T", rt.Val)
	}

	if outerArr.Elements == nil || len(outerArr.Elements.Items) != 2 {
		t.Fatalf("expected 2 elements in outer ARRAY, got %v", outerArr.Elements)
	}

	// Each element should be an A_ArrayExpr
	for i, item := range outerArr.Elements.Items {
		innerArr, ok := item.(*nodes.A_ArrayExpr)
		if !ok {
			t.Fatalf("element %d: expected *nodes.A_ArrayExpr, got %T", i, item)
		}
		if innerArr.Elements == nil || len(innerArr.Elements.Items) != 2 {
			t.Fatalf("element %d: expected 2 inner elements, got %v", i, innerArr.Elements)
		}
	}

	// Verify first inner array has elements 1, 2
	inner0 := outerArr.Elements.Items[0].(*nodes.A_ArrayExpr)
	ac0 := inner0.Elements.Items[0].(*nodes.A_Const)
	iv0 := ac0.Val.(*nodes.Integer)
	if iv0.Ival != 1 {
		t.Errorf("expected inner[0][0]=1, got %d", iv0.Ival)
	}
	ac1 := inner0.Elements.Items[1].(*nodes.A_Const)
	iv1 := ac1.Val.(*nodes.Integer)
	if iv1.Ival != 2 {
		t.Errorf("expected inner[0][1]=2, got %d", iv1.Ival)
	}

	// Verify second inner array has elements 3, 4
	inner1 := outerArr.Elements.Items[1].(*nodes.A_ArrayExpr)
	ac2 := inner1.Elements.Items[0].(*nodes.A_Const)
	iv2 := ac2.Val.(*nodes.Integer)
	if iv2.Ival != 3 {
		t.Errorf("expected inner[1][0]=3, got %d", iv2.Ival)
	}
	ac3 := inner1.Elements.Items[1].(*nodes.A_Const)
	iv3 := ac3.Val.(*nodes.Integer)
	if iv3.Ival != 4 {
		t.Errorf("expected inner[1][1]=4, got %d", iv3.Ival)
	}
}

func TestParseArrayInWhere(t *testing.T) {
	// arr[1] = 5 in WHERE clause
	input := "SELECT * FROM t WHERE arr[1] = 5"

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

	// Left side should be A_Indirection (arr[1])
	indir, ok := aexpr.Lexpr.(*nodes.A_Indirection)
	if !ok {
		t.Fatalf("expected *nodes.A_Indirection on left side, got %T", aexpr.Lexpr)
	}

	// Arg should be a ColumnRef for "arr"
	colRef, ok := indir.Arg.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for Arg, got %T", indir.Arg)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "arr" {
		t.Errorf("expected column 'arr', got %q", str.Str)
	}

	// Indirection should have A_Indices with Uidx=1
	if indir.Indirection == nil || len(indir.Indirection.Items) != 1 {
		t.Fatalf("expected 1 indirection element")
	}
	indices, ok := indir.Indirection.Items[0].(*nodes.A_Indices)
	if !ok {
		t.Fatalf("expected *nodes.A_Indices, got %T", indir.Indirection.Items[0])
	}
	ac, ok := indices.Uidx.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for Uidx, got %T", indices.Uidx)
	}
	iv, ok := ac.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", ac.Val)
	}
	if iv.Ival != 1 {
		t.Errorf("expected Uidx=1, got %d", iv.Ival)
	}

	// Right side should be integer constant 5
	acRight, ok := aexpr.Rexpr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const on right side, got %T", aexpr.Rexpr)
	}
	ivRight, ok := acRight.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", acRight.Val)
	}
	if ivRight.Ival != 5 {
		t.Errorf("expected right side value 5, got %d", ivRight.Ival)
	}
}
