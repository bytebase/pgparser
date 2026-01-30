package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

func TestParseLimitOnly(t *testing.T) {
	input := "SELECT * FROM t LIMIT 10"

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

	// Verify LimitCount is set to integer 10
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	ac, ok := stmt.LimitCount.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for LimitCount, got %T", stmt.LimitCount)
	}
	iv, ok := ac.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", ac.Val)
	}
	if iv.Ival != 10 {
		t.Errorf("expected LimitCount=10, got %d", iv.Ival)
	}

	// Verify LimitOffset is nil
	if stmt.LimitOffset != nil {
		t.Error("expected LimitOffset to be nil")
	}

	// Verify LimitOption
	if stmt.LimitOption != nodes.LIMIT_OPTION_COUNT {
		t.Errorf("expected LIMIT_OPTION_COUNT, got %d", stmt.LimitOption)
	}
}

func TestParseOffsetOnly(t *testing.T) {
	input := "SELECT * FROM t OFFSET 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitOffset is set to integer 5
	if stmt.LimitOffset == nil {
		t.Fatal("expected LimitOffset to be set")
	}
	ac, ok := stmt.LimitOffset.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for LimitOffset, got %T", stmt.LimitOffset)
	}
	iv, ok := ac.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", ac.Val)
	}
	if iv.Ival != 5 {
		t.Errorf("expected LimitOffset=5, got %d", iv.Ival)
	}

	// Verify LimitCount is nil (offset only)
	if stmt.LimitCount != nil {
		t.Error("expected LimitCount to be nil")
	}
}

func TestParseLimitOffset(t *testing.T) {
	input := "SELECT * FROM t LIMIT 10 OFFSET 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is 10
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	lcAc := stmt.LimitCount.(*nodes.A_Const)
	lcIv := lcAc.Val.(*nodes.Integer)
	if lcIv.Ival != 10 {
		t.Errorf("expected LimitCount=10, got %d", lcIv.Ival)
	}

	// Verify LimitOffset is 5
	if stmt.LimitOffset == nil {
		t.Fatal("expected LimitOffset to be set")
	}
	loAc := stmt.LimitOffset.(*nodes.A_Const)
	loIv := loAc.Val.(*nodes.Integer)
	if loIv.Ival != 5 {
		t.Errorf("expected LimitOffset=5, got %d", loIv.Ival)
	}

	// Verify LimitOption
	if stmt.LimitOption != nodes.LIMIT_OPTION_COUNT {
		t.Errorf("expected LIMIT_OPTION_COUNT, got %d", stmt.LimitOption)
	}
}

func TestParseOffsetLimit(t *testing.T) {
	// Reversed order: OFFSET before LIMIT
	input := "SELECT * FROM t OFFSET 5 LIMIT 10"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is 10
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	lcAc := stmt.LimitCount.(*nodes.A_Const)
	lcIv := lcAc.Val.(*nodes.Integer)
	if lcIv.Ival != 10 {
		t.Errorf("expected LimitCount=10, got %d", lcIv.Ival)
	}

	// Verify LimitOffset is 5
	if stmt.LimitOffset == nil {
		t.Fatal("expected LimitOffset to be set")
	}
	loAc := stmt.LimitOffset.(*nodes.A_Const)
	loIv := loAc.Val.(*nodes.Integer)
	if loIv.Ival != 5 {
		t.Errorf("expected LimitOffset=5, got %d", loIv.Ival)
	}

	// Verify LimitOption
	if stmt.LimitOption != nodes.LIMIT_OPTION_COUNT {
		t.Errorf("expected LIMIT_OPTION_COUNT, got %d", stmt.LimitOption)
	}
}

func TestParseLimitAll(t *testing.T) {
	input := "SELECT * FROM t LIMIT ALL"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is a NULL A_Const (LIMIT ALL)
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	ac, ok := stmt.LimitCount.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for LimitCount, got %T", stmt.LimitCount)
	}
	if !ac.Isnull {
		t.Error("expected LimitCount to be a NULL A_Const for LIMIT ALL")
	}

	// Verify LimitOffset is nil
	if stmt.LimitOffset != nil {
		t.Error("expected LimitOffset to be nil")
	}
}

func TestParseFetchFirstRowsOnly(t *testing.T) {
	input := "SELECT * FROM t FETCH FIRST 5 ROWS ONLY"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is 5
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	ac := stmt.LimitCount.(*nodes.A_Const)
	iv := ac.Val.(*nodes.Integer)
	if iv.Ival != 5 {
		t.Errorf("expected LimitCount=5, got %d", iv.Ival)
	}

	// Verify LimitOption is COUNT (ONLY)
	if stmt.LimitOption != nodes.LIMIT_OPTION_COUNT {
		t.Errorf("expected LIMIT_OPTION_COUNT, got %d", stmt.LimitOption)
	}
}

func TestParseFetchNextRowOnly(t *testing.T) {
	// FETCH NEXT ROW ONLY implies LIMIT 1
	input := "SELECT * FROM t FETCH NEXT ROW ONLY"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is 1
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	ac := stmt.LimitCount.(*nodes.A_Const)
	iv := ac.Val.(*nodes.Integer)
	if iv.Ival != 1 {
		t.Errorf("expected LimitCount=1, got %d", iv.Ival)
	}

	// Verify LimitOption is COUNT (ONLY)
	if stmt.LimitOption != nodes.LIMIT_OPTION_COUNT {
		t.Errorf("expected LIMIT_OPTION_COUNT, got %d", stmt.LimitOption)
	}
}

func TestParseOrderByLimit(t *testing.T) {
	input := "SELECT * FROM t ORDER BY id LIMIT 10 OFFSET 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify ORDER BY clause
	if stmt.SortClause == nil {
		t.Fatal("expected SortClause to be set")
	}
	if len(stmt.SortClause.Items) != 1 {
		t.Fatalf("expected 1 sort item, got %d", len(stmt.SortClause.Items))
	}
	sb := stmt.SortClause.Items[0].(*nodes.SortBy)
	colRef := sb.Node.(*nodes.ColumnRef)
	colStr := colRef.Fields.Items[0].(*nodes.String)
	if colStr.Str != "id" {
		t.Errorf("expected sort by 'id', got %q", colStr.Str)
	}

	// Verify LimitCount is 10
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	lcAc := stmt.LimitCount.(*nodes.A_Const)
	lcIv := lcAc.Val.(*nodes.Integer)
	if lcIv.Ival != 10 {
		t.Errorf("expected LimitCount=10, got %d", lcIv.Ival)
	}

	// Verify LimitOffset is 5
	if stmt.LimitOffset == nil {
		t.Fatal("expected LimitOffset to be set")
	}
	loAc := stmt.LimitOffset.(*nodes.A_Const)
	loIv := loAc.Val.(*nodes.Integer)
	if loIv.Ival != 5 {
		t.Errorf("expected LimitOffset=5, got %d", loIv.Ival)
	}
}

func TestParseCTEWithLimit(t *testing.T) {
	input := "WITH cte AS (SELECT 1) SELECT * FROM cte LIMIT 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify WITH clause
	if stmt.WithClause == nil {
		t.Fatal("expected WithClause to be set")
	}
	if stmt.WithClause.Ctes == nil || len(stmt.WithClause.Ctes.Items) != 1 {
		t.Fatal("expected 1 CTE")
	}
	cte := stmt.WithClause.Ctes.Items[0].(*nodes.CommonTableExpr)
	if cte.Ctename != "cte" {
		t.Errorf("expected CTE name 'cte', got %q", cte.Ctename)
	}

	// Verify LimitCount is 5
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	ac := stmt.LimitCount.(*nodes.A_Const)
	iv := ac.Val.(*nodes.Integer)
	if iv.Ival != 5 {
		t.Errorf("expected LimitCount=5, got %d", iv.Ival)
	}
}

func TestParseFetchFirstWithTies(t *testing.T) {
	input := "SELECT * FROM t FETCH FIRST 3 ROWS WITH TIES"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is 3
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	ac := stmt.LimitCount.(*nodes.A_Const)
	iv := ac.Val.(*nodes.Integer)
	if iv.Ival != 3 {
		t.Errorf("expected LimitCount=3, got %d", iv.Ival)
	}

	// Verify LimitOption is WITH_TIES
	if stmt.LimitOption != nodes.LIMIT_OPTION_WITH_TIES {
		t.Errorf("expected LIMIT_OPTION_WITH_TIES, got %d", stmt.LimitOption)
	}
}

func TestParseFetchNextRowWithTies(t *testing.T) {
	// FETCH NEXT ROW WITH TIES implies LIMIT 1 WITH TIES
	input := "SELECT * FROM t FETCH NEXT ROW WITH TIES"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is 1
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	ac := stmt.LimitCount.(*nodes.A_Const)
	iv := ac.Val.(*nodes.Integer)
	if iv.Ival != 1 {
		t.Errorf("expected LimitCount=1, got %d", iv.Ival)
	}

	// Verify LimitOption is WITH_TIES
	if stmt.LimitOption != nodes.LIMIT_OPTION_WITH_TIES {
		t.Errorf("expected LIMIT_OPTION_WITH_TIES, got %d", stmt.LimitOption)
	}
}

func TestParseLimitExpression(t *testing.T) {
	// LIMIT with an expression instead of a simple constant
	input := "SELECT * FROM t LIMIT 5 + 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Verify LimitCount is an expression (A_Expr with +)
	if stmt.LimitCount == nil {
		t.Fatal("expected LimitCount to be set")
	}
	aexpr, ok := stmt.LimitCount.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr for LimitCount, got %T", stmt.LimitCount)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
}
