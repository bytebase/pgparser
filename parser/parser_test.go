package parser

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
)

func TestLexerTokenMapping(t *testing.T) {
	// Debug: Check what tokens the lexer produces and how they map
	input := "SELECT 1"
	lexer := NewLexer(input)

	tok1 := lexer.NextToken()
	t.Logf("Token 1: Type=%d, Str=%q (SELECT keyword should be %d)", tok1.Type, tok1.Str, SELECT)

	tok2 := lexer.NextToken()
	t.Logf("Token 2: Type=%d, Str=%q (ICONST should be %d, lex_ICONST=%d)", tok2.Type, tok2.Str, ICONST, lex_ICONST)

	tok3 := lexer.NextToken()
	t.Logf("Token 3: Type=%d (EOF should be 0, lex_EOF=%d)", tok3.Type, lex_EOF)
}

func TestParserLexerMapping(t *testing.T) {
	// Debug: Check what tokens the parserLexer produces for the parser
	input := "SELECT 1"
	pl := newParserLexer(input)

	var lval pgSymType
	tok1 := pl.Lex(&lval)
	t.Logf("Parser Token 1: Type=%d, lval.str=%q (SELECT should be %d)", tok1, lval.str, SELECT)

	tok2 := pl.Lex(&lval)
	t.Logf("Parser Token 2: Type=%d, lval.ival=%d (ICONST should be %d)", tok2, lval.ival, ICONST)

	tok3 := pl.Lex(&lval)
	t.Logf("Parser Token 3: Type=%d (EOF should be 0)", tok3)
}

func TestParseSimpleSelect(t *testing.T) {
	input := "SELECT 1"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Parse returned nil result (no error reported)")
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(result.Items))
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target in target list")
	}
}

func TestParseSelectWithFrom(t *testing.T) {
	input := "SELECT * FROM users"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Parse returned nil result")
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(result.Items))
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.FromClause == nil || len(stmt.FromClause.Items) != 1 {
		t.Fatalf("expected 1 table in from clause")
	}

	rv, ok := stmt.FromClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", stmt.FromClause.Items[0])
	}

	if rv.Relname != "users" {
		t.Errorf("expected table name 'users', got %q", rv.Relname)
	}
}

func TestParseSelectWithWhere(t *testing.T) {
	input := "SELECT * FROM users WHERE id = 1"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Parse returned nil result")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.WhereClause == nil {
		t.Fatal("expected where clause")
	}

	// Where clause should be an A_Expr with "=" operator
	aexpr, ok := stmt.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", stmt.WhereClause)
	}

	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP kind")
	}
}

func TestParseSelectWithGroupBy(t *testing.T) {
	input := "SELECT a, count(*) FROM t GROUP BY a"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Parse returned nil result")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	// Verify target list has 2 items (a, count(*))
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 2 {
		t.Fatalf("expected 2 targets in target list, got %d", len(stmt.TargetList.Items))
	}

	// Verify FROM clause
	if stmt.FromClause == nil || len(stmt.FromClause.Items) != 1 {
		t.Fatalf("expected 1 table in from clause")
	}

	// Verify GROUP BY clause
	if stmt.GroupClause == nil {
		t.Fatal("expected group clause")
	}

	if len(stmt.GroupClause.Items) != 1 {
		t.Fatalf("expected 1 item in group clause, got %d", len(stmt.GroupClause.Items))
	}

	// The group by item should be a ColumnRef for "a"
	colRef, ok := stmt.GroupClause.Items[0].(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef in group clause, got %T", stmt.GroupClause.Items[0])
	}

	if colRef.Fields == nil || len(colRef.Fields.Items) != 1 {
		t.Fatal("expected 1 field in column ref")
	}

	str, ok := colRef.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef.Fields.Items[0])
	}

	if str.Str != "a" {
		t.Errorf("expected group by column 'a', got %q", str.Str)
	}

	// Verify no HAVING clause
	if stmt.HavingClause != nil {
		t.Error("expected no having clause")
	}
}

func TestParseSelectWithGroupByHaving(t *testing.T) {
	input := "SELECT a, count(*) FROM t GROUP BY a HAVING count(*) > 1"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Parse returned nil result")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	// Verify GROUP BY clause
	if stmt.GroupClause == nil {
		t.Fatal("expected group clause")
	}

	if len(stmt.GroupClause.Items) != 1 {
		t.Fatalf("expected 1 item in group clause, got %d", len(stmt.GroupClause.Items))
	}

	// Verify HAVING clause exists and is an A_Expr (count(*) > 1)
	if stmt.HavingClause == nil {
		t.Fatal("expected having clause")
	}

	aexpr, ok := stmt.HavingClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr in having clause, got %T", stmt.HavingClause)
	}

	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP kind in having clause")
	}

	// Left side should be a FuncCall (count(*))
	funcCall, ok := aexpr.Lexpr.(*nodes.FuncCall)
	if !ok {
		t.Fatalf("expected *nodes.FuncCall on left side of having, got %T", aexpr.Lexpr)
	}

	if !funcCall.AggStar {
		t.Error("expected AggStar to be true for count(*)")
	}

	// Right side should be integer constant 1
	aConst, ok := aexpr.Rexpr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const on right side of having, got %T", aexpr.Rexpr)
	}

	intVal, ok := aConst.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", aConst.Val)
	}

	if intVal.Ival != 1 {
		t.Errorf("expected integer value 1, got %d", intVal.Ival)
	}
}

func TestParseUnion(t *testing.T) {
	input := "SELECT a FROM t1 UNION SELECT b FROM t2"

	result, err := Parse(input)
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

	// Verify set operation fields
	if stmt.Op != nodes.SETOP_UNION {
		t.Errorf("expected SETOP_UNION (%d), got %d", nodes.SETOP_UNION, stmt.Op)
	}
	if stmt.All != false {
		t.Errorf("expected All=false, got true")
	}
	if stmt.Larg == nil {
		t.Fatal("expected non-nil Larg")
	}
	if stmt.Rarg == nil {
		t.Fatal("expected non-nil Rarg")
	}

	// Verify left side is "SELECT a FROM t1"
	if stmt.Larg.TargetList == nil || len(stmt.Larg.TargetList.Items) != 1 {
		t.Fatal("expected 1 target in left arg")
	}
	lTarget := stmt.Larg.TargetList.Items[0].(*nodes.ResTarget)
	lCol := lTarget.Val.(*nodes.ColumnRef)
	lStr := lCol.Fields.Items[0].(*nodes.String)
	if lStr.Str != "a" {
		t.Errorf("expected left target 'a', got %q", lStr.Str)
	}

	// Verify right side is "SELECT b FROM t2"
	if stmt.Rarg.TargetList == nil || len(stmt.Rarg.TargetList.Items) != 1 {
		t.Fatal("expected 1 target in right arg")
	}
	rTarget := stmt.Rarg.TargetList.Items[0].(*nodes.ResTarget)
	rCol := rTarget.Val.(*nodes.ColumnRef)
	rStr := rCol.Fields.Items[0].(*nodes.String)
	if rStr.Str != "b" {
		t.Errorf("expected right target 'b', got %q", rStr.Str)
	}
}

func TestParseUnionAll(t *testing.T) {
	input := "SELECT a FROM t1 UNION ALL SELECT b FROM t2"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.Op != nodes.SETOP_UNION {
		t.Errorf("expected SETOP_UNION, got %d", stmt.Op)
	}
	if stmt.All != true {
		t.Errorf("expected All=true, got false")
	}
	if stmt.Larg == nil || stmt.Rarg == nil {
		t.Fatal("expected non-nil Larg and Rarg")
	}
}

func TestParseIntersect(t *testing.T) {
	input := "SELECT a FROM t1 INTERSECT SELECT b FROM t2"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.Op != nodes.SETOP_INTERSECT {
		t.Errorf("expected SETOP_INTERSECT (%d), got %d", nodes.SETOP_INTERSECT, stmt.Op)
	}
	if stmt.All != false {
		t.Errorf("expected All=false, got true")
	}
	if stmt.Larg == nil || stmt.Rarg == nil {
		t.Fatal("expected non-nil Larg and Rarg")
	}
}

func TestParseExcept(t *testing.T) {
	input := "SELECT a FROM t1 EXCEPT SELECT b FROM t2"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.Op != nodes.SETOP_EXCEPT {
		t.Errorf("expected SETOP_EXCEPT (%d), got %d", nodes.SETOP_EXCEPT, stmt.Op)
	}
	if stmt.All != false {
		t.Errorf("expected All=false, got true")
	}
	if stmt.Larg == nil || stmt.Rarg == nil {
		t.Fatal("expected non-nil Larg and Rarg")
	}
}

func TestParseUnionChained(t *testing.T) {
	// UNION is left-associative, so:
	// SELECT a FROM t1 UNION SELECT b FROM t2 UNION SELECT c FROM t3
	// should parse as:
	// (SELECT a FROM t1 UNION SELECT b FROM t2) UNION SELECT c FROM t3
	input := "SELECT a FROM t1 UNION SELECT b FROM t2 UNION SELECT c FROM t3"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Top level should be UNION
	if stmt.Op != nodes.SETOP_UNION {
		t.Errorf("expected top-level SETOP_UNION, got %d", stmt.Op)
	}

	// Right side should be a simple SELECT (SELECT c FROM t3)
	if stmt.Rarg == nil {
		t.Fatal("expected non-nil Rarg")
	}
	if stmt.Rarg.Op != nodes.SETOP_NONE {
		t.Errorf("expected right side to be simple SELECT (SETOP_NONE), got %d", stmt.Rarg.Op)
	}
	rTarget := stmt.Rarg.TargetList.Items[0].(*nodes.ResTarget)
	rCol := rTarget.Val.(*nodes.ColumnRef)
	rStr := rCol.Fields.Items[0].(*nodes.String)
	if rStr.Str != "c" {
		t.Errorf("expected right target 'c', got %q", rStr.Str)
	}

	// Left side should be another UNION (SELECT a FROM t1 UNION SELECT b FROM t2)
	if stmt.Larg == nil {
		t.Fatal("expected non-nil Larg")
	}
	if stmt.Larg.Op != nodes.SETOP_UNION {
		t.Errorf("expected left side to be SETOP_UNION, got %d", stmt.Larg.Op)
	}

	// Left-left should be "SELECT a FROM t1"
	if stmt.Larg.Larg == nil {
		t.Fatal("expected non-nil Larg.Larg")
	}
	llTarget := stmt.Larg.Larg.TargetList.Items[0].(*nodes.ResTarget)
	llCol := llTarget.Val.(*nodes.ColumnRef)
	llStr := llCol.Fields.Items[0].(*nodes.String)
	if llStr.Str != "a" {
		t.Errorf("expected left-left target 'a', got %q", llStr.Str)
	}

	// Left-right should be "SELECT b FROM t2"
	if stmt.Larg.Rarg == nil {
		t.Fatal("expected non-nil Larg.Rarg")
	}
	lrTarget := stmt.Larg.Rarg.TargetList.Items[0].(*nodes.ResTarget)
	lrCol := lrTarget.Val.(*nodes.ColumnRef)
	lrStr := lrCol.Fields.Items[0].(*nodes.String)
	if lrStr.Str != "b" {
		t.Errorf("expected left-right target 'b', got %q", lrStr.Str)
	}
}

func TestParseIntersectPrecedence(t *testing.T) {
	// INTERSECT has higher precedence than UNION, so:
	// SELECT a FROM t1 UNION SELECT b FROM t2 INTERSECT SELECT c FROM t3
	// should parse as:
	// SELECT a FROM t1 UNION (SELECT b FROM t2 INTERSECT SELECT c FROM t3)
	input := "SELECT a FROM t1 UNION SELECT b FROM t2 INTERSECT SELECT c FROM t3"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Top level should be UNION (lower precedence)
	if stmt.Op != nodes.SETOP_UNION {
		t.Errorf("expected top-level SETOP_UNION, got %d", stmt.Op)
	}

	// Left side should be simple SELECT (SELECT a FROM t1)
	if stmt.Larg == nil {
		t.Fatal("expected non-nil Larg")
	}
	if stmt.Larg.Op != nodes.SETOP_NONE {
		t.Errorf("expected left side to be simple SELECT (SETOP_NONE), got %d", stmt.Larg.Op)
	}

	// Right side should be INTERSECT (higher precedence binds tighter)
	if stmt.Rarg == nil {
		t.Fatal("expected non-nil Rarg")
	}
	if stmt.Rarg.Op != nodes.SETOP_INTERSECT {
		t.Errorf("expected right side to be SETOP_INTERSECT, got %d", stmt.Rarg.Op)
	}
}

func TestParseUnionDistinct(t *testing.T) {
	// UNION DISTINCT is the same as plain UNION (All=false)
	input := "SELECT a FROM t1 UNION DISTINCT SELECT b FROM t2"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	if stmt.Op != nodes.SETOP_UNION {
		t.Errorf("expected SETOP_UNION, got %d", stmt.Op)
	}
	if stmt.All != false {
		t.Errorf("expected All=false for UNION DISTINCT, got true")
	}
}

func TestParseSetOpWithOrderBy(t *testing.T) {
	// ORDER BY applies to the entire result of the UNION
	input := "SELECT a FROM t1 UNION SELECT b FROM t2 ORDER BY a"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// The top-level statement should have a sort clause
	// In PostgreSQL, ORDER BY on a UNION applies to the whole result.
	// The way our grammar works, select_no_parens handles "select_clause sort_clause"
	// so the top-level stmt is a SelectStmt with Op=SETOP_UNION and SortClause set.
	if stmt.Op != nodes.SETOP_UNION {
		t.Errorf("expected SETOP_UNION, got %d", stmt.Op)
	}
	if stmt.SortClause == nil {
		t.Error("expected SortClause to be set for UNION ... ORDER BY")
	}
}

func TestParseSelectWithMultipleGroupBy(t *testing.T) {
	input := "SELECT a, b, sum(c) FROM t GROUP BY a, b"

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Parse returned nil result")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	// Verify target list has 3 items (a, b, sum(c))
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 3 {
		t.Fatalf("expected 3 targets in target list, got %d", len(stmt.TargetList.Items))
	}

	// Verify GROUP BY clause has 2 items
	if stmt.GroupClause == nil {
		t.Fatal("expected group clause")
	}

	if len(stmt.GroupClause.Items) != 2 {
		t.Fatalf("expected 2 items in group clause, got %d", len(stmt.GroupClause.Items))
	}

	// Verify first group by item is column "a"
	colRef1, ok := stmt.GroupClause.Items[0].(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for first group item, got %T", stmt.GroupClause.Items[0])
	}

	str1, ok := colRef1.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef1.Fields.Items[0])
	}

	if str1.Str != "a" {
		t.Errorf("expected first group by column 'a', got %q", str1.Str)
	}

	// Verify second group by item is column "b"
	colRef2, ok := stmt.GroupClause.Items[1].(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for second group item, got %T", stmt.GroupClause.Items[1])
	}

	str2, ok := colRef2.Fields.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", colRef2.Fields.Items[0])
	}

	if str2.Str != "b" {
		t.Errorf("expected second group by column 'b', got %q", str2.Str)
	}
}

func TestOrderedSetAggrArgs(t *testing.T) {
	sqls := []string{
		// CREATE AGGREGATE with ordered-set args (ORDER BY only)
		`CREATE AGGREGATE my_percentile(ORDER BY float8) (sfunc = ordered_set_transition, stype = internal, finalfunc = percentile_disc_final, finalfunc_extra)`,
		// CREATE AGGREGATE with both direct and ordered args
		`CREATE AGGREGATE my_percentile2(float8 ORDER BY float8) (sfunc = ordered_set_transition, stype = internal, finalfunc = percentile_disc_final, finalfunc_extra)`,
		// DROP with ordered-set args
		`DROP AGGREGATE my_percentile(ORDER BY float8)`,
		// ALTER with ordered-set args
		`ALTER AGGREGATE my_percentile(ORDER BY float8) RENAME TO new_name`,
	}
	for _, sql := range sqls {
		label := sql
		if len(label) > 40 {
			label = label[:40]
		}
		t.Run(label, func(t *testing.T) {
			_, err := Parse(sql)
			if err != nil {
				preview := sql
				if len(preview) > 60 {
					preview = preview[:60]
				}
				t.Errorf("Parse(%q) failed: %v", preview, err)
			}
		})
	}
}
