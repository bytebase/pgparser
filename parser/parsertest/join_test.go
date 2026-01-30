package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestParseCrossJoin verifies CROSS JOIN produces a JoinExpr with JOIN_INNER and no quals.
func TestParseCrossJoin(t *testing.T) {
	input := "SELECT * FROM t1 CROSS JOIN t2"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	if stmt.FromClause == nil || len(stmt.FromClause.Items) != 1 {
		t.Fatalf("expected 1 item in from clause, got %d", len(stmt.FromClause.Items))
	}

	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_INNER {
		t.Errorf("expected JOIN_INNER (%d), got %d", nodes.JOIN_INNER, join.Jointype)
	}
	if join.IsNatural {
		t.Error("expected IsNatural=false for CROSS JOIN")
	}
	if join.Quals != nil {
		t.Error("expected no quals for CROSS JOIN")
	}
	if join.UsingClause != nil {
		t.Error("expected no UsingClause for CROSS JOIN")
	}

	// Verify left and right args are RangeVars
	larg, ok := join.Larg.(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected Larg to be *nodes.RangeVar, got %T", join.Larg)
	}
	if larg.Relname != "t1" {
		t.Errorf("expected left table 't1', got %q", larg.Relname)
	}

	rarg, ok := join.Rarg.(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected Rarg to be *nodes.RangeVar, got %T", join.Rarg)
	}
	if rarg.Relname != "t2" {
		t.Errorf("expected right table 't2', got %q", rarg.Relname)
	}
}

// TestParseInnerJoinOn verifies INNER JOIN ... ON produces a JoinExpr with JOIN_INNER and Quals set.
func TestParseInnerJoinOn(t *testing.T) {
	input := "SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_INNER {
		t.Errorf("expected JOIN_INNER (%d), got %d", nodes.JOIN_INNER, join.Jointype)
	}
	if join.IsNatural {
		t.Error("expected IsNatural=false")
	}
	if join.Quals == nil {
		t.Fatal("expected Quals to be set for ON clause")
	}

	// Quals should be an A_Expr with "=" operator
	aexpr, ok := join.Quals.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", join.Quals)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
}

// TestParseLeftJoinOn verifies LEFT JOIN produces a JoinExpr with JOIN_LEFT.
func TestParseLeftJoinOn(t *testing.T) {
	input := "SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_LEFT {
		t.Errorf("expected JOIN_LEFT (%d), got %d", nodes.JOIN_LEFT, join.Jointype)
	}
	if join.Quals == nil {
		t.Fatal("expected Quals to be set")
	}
}

// TestParseRightJoinOn verifies RIGHT JOIN produces a JoinExpr with JOIN_RIGHT.
func TestParseRightJoinOn(t *testing.T) {
	input := "SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_RIGHT {
		t.Errorf("expected JOIN_RIGHT (%d), got %d", nodes.JOIN_RIGHT, join.Jointype)
	}
	if join.Quals == nil {
		t.Fatal("expected Quals to be set")
	}
}

// TestParseFullJoinOn verifies FULL JOIN produces a JoinExpr with JOIN_FULL.
func TestParseFullJoinOn(t *testing.T) {
	input := "SELECT * FROM t1 FULL JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_FULL {
		t.Errorf("expected JOIN_FULL (%d), got %d", nodes.JOIN_FULL, join.Jointype)
	}
	if join.Quals == nil {
		t.Fatal("expected Quals to be set")
	}
}

// TestParseLeftOuterJoinOn verifies LEFT OUTER JOIN is the same as LEFT JOIN.
func TestParseLeftOuterJoinOn(t *testing.T) {
	input := "SELECT * FROM t1 LEFT OUTER JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_LEFT {
		t.Errorf("expected JOIN_LEFT (%d) for LEFT OUTER JOIN, got %d", nodes.JOIN_LEFT, join.Jointype)
	}
	if join.Quals == nil {
		t.Fatal("expected Quals to be set")
	}
}

// TestParseJoinUsing verifies JOIN ... USING sets UsingClause.
func TestParseJoinUsing(t *testing.T) {
	input := "SELECT * FROM t1 JOIN t2 USING (id)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_INNER {
		t.Errorf("expected JOIN_INNER (%d), got %d", nodes.JOIN_INNER, join.Jointype)
	}
	if join.UsingClause == nil {
		t.Fatal("expected UsingClause to be set")
	}
	if join.Quals != nil {
		t.Error("expected Quals to be nil when USING is specified")
	}

	// UsingClause should contain one item: "id"
	if len(join.UsingClause.Items) != 1 {
		t.Fatalf("expected 1 item in UsingClause, got %d", len(join.UsingClause.Items))
	}
	str, ok := join.UsingClause.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String in UsingClause, got %T", join.UsingClause.Items[0])
	}
	if str.Str != "id" {
		t.Errorf("expected USING column 'id', got %q", str.Str)
	}
}

// TestParseJoinUsingMultipleColumns verifies JOIN ... USING (a, b) works.
func TestParseJoinUsingMultipleColumns(t *testing.T) {
	input := "SELECT * FROM t1 JOIN t2 USING (a, b)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.UsingClause == nil {
		t.Fatal("expected UsingClause to be set")
	}
	if len(join.UsingClause.Items) != 2 {
		t.Fatalf("expected 2 items in UsingClause, got %d", len(join.UsingClause.Items))
	}

	str0 := join.UsingClause.Items[0].(*nodes.String)
	str1 := join.UsingClause.Items[1].(*nodes.String)
	if str0.Str != "a" {
		t.Errorf("expected first USING column 'a', got %q", str0.Str)
	}
	if str1.Str != "b" {
		t.Errorf("expected second USING column 'b', got %q", str1.Str)
	}
}

// TestParseNaturalJoin verifies NATURAL JOIN sets IsNatural=true.
func TestParseNaturalJoin(t *testing.T) {
	input := "SELECT * FROM t1 NATURAL JOIN t2"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_INNER {
		t.Errorf("expected JOIN_INNER (%d), got %d", nodes.JOIN_INNER, join.Jointype)
	}
	if !join.IsNatural {
		t.Error("expected IsNatural=true for NATURAL JOIN")
	}
	if join.Quals != nil {
		t.Error("expected no Quals for NATURAL JOIN")
	}
	if join.UsingClause != nil {
		t.Error("expected no UsingClause for NATURAL JOIN")
	}
}

// TestParseNaturalLeftJoin verifies NATURAL LEFT JOIN sets both IsNatural and JOIN_LEFT.
func TestParseNaturalLeftJoin(t *testing.T) {
	input := "SELECT * FROM t1 NATURAL LEFT JOIN t2"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_LEFT {
		t.Errorf("expected JOIN_LEFT (%d), got %d", nodes.JOIN_LEFT, join.Jointype)
	}
	if !join.IsNatural {
		t.Error("expected IsNatural=true for NATURAL LEFT JOIN")
	}
}

// TestParseChainedJoins verifies chained joins (left-associative):
// t1 JOIN t2 ON ... JOIN t3 ON ... => (t1 JOIN t2 ON ...) JOIN t3 ON ...
func TestParseChainedJoins(t *testing.T) {
	input := "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id JOIN t3 ON t2.id = t3.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	if stmt.FromClause == nil || len(stmt.FromClause.Items) != 1 {
		t.Fatalf("expected 1 item in from clause")
	}

	// Top-level join: (t1 JOIN t2) JOIN t3
	outerJoin, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if outerJoin.Jointype != nodes.JOIN_INNER {
		t.Errorf("expected outer JOIN_INNER, got %d", outerJoin.Jointype)
	}

	// Right side of outer join should be t3
	rarg, ok := outerJoin.Rarg.(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected outer Rarg to be *nodes.RangeVar, got %T", outerJoin.Rarg)
	}
	if rarg.Relname != "t3" {
		t.Errorf("expected outer right table 't3', got %q", rarg.Relname)
	}

	// Left side of outer join should be another JoinExpr (t1 JOIN t2)
	innerJoin, ok := outerJoin.Larg.(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected outer Larg to be *nodes.JoinExpr (inner join), got %T", outerJoin.Larg)
	}

	if innerJoin.Jointype != nodes.JOIN_INNER {
		t.Errorf("expected inner JOIN_INNER, got %d", innerJoin.Jointype)
	}

	// Left side of inner join should be t1
	innerLarg, ok := innerJoin.Larg.(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected inner Larg to be *nodes.RangeVar, got %T", innerJoin.Larg)
	}
	if innerLarg.Relname != "t1" {
		t.Errorf("expected inner left table 't1', got %q", innerLarg.Relname)
	}

	// Right side of inner join should be t2
	innerRarg, ok := innerJoin.Rarg.(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected inner Rarg to be *nodes.RangeVar, got %T", innerJoin.Rarg)
	}
	if innerRarg.Relname != "t2" {
		t.Errorf("expected inner right table 't2', got %q", innerRarg.Relname)
	}
}

// TestParseImplicitInnerJoinOn verifies that bare JOIN (without INNER keyword) works.
func TestParseImplicitInnerJoinOn(t *testing.T) {
	input := "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_INNER {
		t.Errorf("expected JOIN_INNER (%d) for bare JOIN, got %d", nodes.JOIN_INNER, join.Jointype)
	}
	if join.IsNatural {
		t.Error("expected IsNatural=false for bare JOIN")
	}
	if join.Quals == nil {
		t.Fatal("expected Quals to be set")
	}
}

// TestParseRightOuterJoin verifies RIGHT OUTER JOIN works like RIGHT JOIN.
func TestParseRightOuterJoin(t *testing.T) {
	input := "SELECT * FROM t1 RIGHT OUTER JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_RIGHT {
		t.Errorf("expected JOIN_RIGHT (%d) for RIGHT OUTER JOIN, got %d", nodes.JOIN_RIGHT, join.Jointype)
	}
}

// TestParseFullOuterJoin verifies FULL OUTER JOIN works like FULL JOIN.
func TestParseFullOuterJoin(t *testing.T) {
	input := "SELECT * FROM t1 FULL OUTER JOIN t2 ON t1.id = t2.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr, got %T", stmt.FromClause.Items[0])
	}

	if join.Jointype != nodes.JOIN_FULL {
		t.Errorf("expected JOIN_FULL (%d) for FULL OUTER JOIN, got %d", nodes.JOIN_FULL, join.Jointype)
	}
}

// TestParseJoinWithWhereClause verifies JOIN combined with WHERE clause.
func TestParseJoinWithWhereClause(t *testing.T) {
	input := "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id WHERE t1.active = true"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// FROM clause should have a JoinExpr
	join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
	if !ok {
		t.Fatalf("expected *nodes.JoinExpr in FROM, got %T", stmt.FromClause.Items[0])
	}
	if join.Quals == nil {
		t.Fatal("expected join Quals to be set")
	}

	// WHERE clause should also be set
	if stmt.WhereClause == nil {
		t.Fatal("expected WHERE clause to be set")
	}
}

// TestParseJoinAllTypes runs a table-driven test for all join types to ensure parsing succeeds.
func TestParseJoinAllTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		joinType nodes.JoinType
		natural  bool
		hasQuals bool
		hasUsing bool
	}{
		{
			name:     "CROSS JOIN",
			input:    "SELECT * FROM t1 CROSS JOIN t2",
			joinType: nodes.JOIN_INNER,
			natural:  false,
			hasQuals: false,
			hasUsing: false,
		},
		{
			name:     "INNER JOIN ON",
			input:    "SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_INNER,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
		{
			name:     "LEFT JOIN ON",
			input:    "SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_LEFT,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
		{
			name:     "RIGHT JOIN ON",
			input:    "SELECT * FROM t1 RIGHT JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_RIGHT,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
		{
			name:     "FULL JOIN ON",
			input:    "SELECT * FROM t1 FULL JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_FULL,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
		{
			name:     "LEFT OUTER JOIN ON",
			input:    "SELECT * FROM t1 LEFT OUTER JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_LEFT,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
		{
			name:     "RIGHT OUTER JOIN ON",
			input:    "SELECT * FROM t1 RIGHT OUTER JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_RIGHT,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
		{
			name:     "FULL OUTER JOIN ON",
			input:    "SELECT * FROM t1 FULL OUTER JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_FULL,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
		{
			name:     "JOIN USING",
			input:    "SELECT * FROM t1 JOIN t2 USING (id)",
			joinType: nodes.JOIN_INNER,
			natural:  false,
			hasQuals: false,
			hasUsing: true,
		},
		{
			name:     "NATURAL JOIN",
			input:    "SELECT * FROM t1 NATURAL JOIN t2",
			joinType: nodes.JOIN_INNER,
			natural:  true,
			hasQuals: false,
			hasUsing: false,
		},
		{
			name:     "NATURAL LEFT JOIN",
			input:    "SELECT * FROM t1 NATURAL LEFT JOIN t2",
			joinType: nodes.JOIN_LEFT,
			natural:  true,
			hasQuals: false,
			hasUsing: false,
		},
		{
			name:     "NATURAL RIGHT JOIN",
			input:    "SELECT * FROM t1 NATURAL RIGHT JOIN t2",
			joinType: nodes.JOIN_RIGHT,
			natural:  true,
			hasQuals: false,
			hasUsing: false,
		},
		{
			name:     "NATURAL FULL JOIN",
			input:    "SELECT * FROM t1 NATURAL FULL JOIN t2",
			joinType: nodes.JOIN_FULL,
			natural:  true,
			hasQuals: false,
			hasUsing: false,
		},
		{
			name:     "bare JOIN ON",
			input:    "SELECT * FROM t1 JOIN t2 ON t1.id = t2.id",
			joinType: nodes.JOIN_INNER,
			natural:  false,
			hasQuals: true,
			hasUsing: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			stmt := result.Items[0].(*nodes.SelectStmt)
			if stmt.FromClause == nil || len(stmt.FromClause.Items) < 1 {
				t.Fatalf("expected at least 1 item in from clause for %q", tt.input)
			}

			join, ok := stmt.FromClause.Items[0].(*nodes.JoinExpr)
			if !ok {
				t.Fatalf("expected *nodes.JoinExpr for %q, got %T", tt.input, stmt.FromClause.Items[0])
			}

			if join.Jointype != tt.joinType {
				t.Errorf("expected JoinType %d, got %d", tt.joinType, join.Jointype)
			}
			if join.IsNatural != tt.natural {
				t.Errorf("expected IsNatural=%v, got %v", tt.natural, join.IsNatural)
			}
			if tt.hasQuals && join.Quals == nil {
				t.Error("expected Quals to be set")
			}
			if !tt.hasQuals && join.Quals != nil {
				t.Error("expected Quals to be nil")
			}
			if tt.hasUsing && join.UsingClause == nil {
				t.Error("expected UsingClause to be set")
			}
			if !tt.hasUsing && join.UsingClause != nil {
				t.Error("expected UsingClause to be nil")
			}
		})
	}
}
