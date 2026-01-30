package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestParseBasicCTE tests parsing a basic CTE: WITH cte AS (SELECT 1) SELECT * FROM cte
func TestParseBasicCTE(t *testing.T) {
	input := "WITH cte AS (SELECT 1) SELECT * FROM cte"

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

	// Verify WITH clause exists
	if stmt.WithClause == nil {
		t.Fatal("expected WithClause to be set")
	}

	wc := stmt.WithClause
	if wc.Recursive {
		t.Error("expected Recursive=false for non-recursive CTE")
	}

	// Verify CTE list has 1 entry
	if wc.Ctes == nil || len(wc.Ctes.Items) != 1 {
		t.Fatalf("expected 1 CTE, got %d", len(wc.Ctes.Items))
	}

	cte, ok := wc.Ctes.Items[0].(*nodes.CommonTableExpr)
	if !ok {
		t.Fatalf("expected *nodes.CommonTableExpr, got %T", wc.Ctes.Items[0])
	}

	if cte.Ctename != "cte" {
		t.Errorf("expected CTE name 'cte', got %q", cte.Ctename)
	}

	// CTE query should be a SelectStmt
	cteQuery, ok := cte.Ctequery.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected CTE query to be *nodes.SelectStmt, got %T", cte.Ctequery)
	}

	// CTE query should have a target list with 1 item (the constant 1)
	if cteQuery.TargetList == nil || len(cteQuery.TargetList.Items) != 1 {
		t.Fatalf("expected 1 target in CTE query, got %d", len(cteQuery.TargetList.Items))
	}

	// Verify main query references cte in FROM clause
	if stmt.FromClause == nil || len(stmt.FromClause.Items) != 1 {
		t.Fatal("expected 1 table in FROM clause")
	}

	rv, ok := stmt.FromClause.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", stmt.FromClause.Items[0])
	}

	if rv.Relname != "cte" {
		t.Errorf("expected table name 'cte', got %q", rv.Relname)
	}

	// No alias column names
	if cte.Aliascolnames != nil && len(cte.Aliascolnames.Items) > 0 {
		t.Error("expected no alias column names for basic CTE")
	}
}

// TestParseRecursiveCTE tests parsing a recursive CTE.
func TestParseRecursiveCTE(t *testing.T) {
	input := "WITH RECURSIVE cte AS (SELECT 1 UNION ALL SELECT n FROM cte) SELECT * FROM cte"

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

	// Verify WITH clause exists and is recursive
	if stmt.WithClause == nil {
		t.Fatal("expected WithClause to be set")
	}

	wc := stmt.WithClause
	if !wc.Recursive {
		t.Error("expected Recursive=true for WITH RECURSIVE")
	}

	// Verify CTE list has 1 entry
	if wc.Ctes == nil || len(wc.Ctes.Items) != 1 {
		t.Fatalf("expected 1 CTE, got %d", len(wc.Ctes.Items))
	}

	cte, ok := wc.Ctes.Items[0].(*nodes.CommonTableExpr)
	if !ok {
		t.Fatalf("expected *nodes.CommonTableExpr, got %T", wc.Ctes.Items[0])
	}

	if cte.Ctename != "cte" {
		t.Errorf("expected CTE name 'cte', got %q", cte.Ctename)
	}

	// CTE query should be a UNION ALL SelectStmt
	cteQuery, ok := cte.Ctequery.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected CTE query to be *nodes.SelectStmt, got %T", cte.Ctequery)
	}

	if cteQuery.Op != nodes.SETOP_UNION {
		t.Errorf("expected SETOP_UNION, got %d", cteQuery.Op)
	}

	if !cteQuery.All {
		t.Error("expected All=true for UNION ALL")
	}
}

// TestParseMultipleCTEs tests parsing multiple CTEs.
func TestParseMultipleCTEs(t *testing.T) {
	input := "WITH cte1 AS (SELECT 1), cte2 AS (SELECT 2) SELECT * FROM cte1, cte2"

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

	if stmt.WithClause == nil {
		t.Fatal("expected WithClause to be set")
	}

	wc := stmt.WithClause
	if wc.Recursive {
		t.Error("expected Recursive=false")
	}

	// Verify CTE list has 2 entries
	if wc.Ctes == nil || len(wc.Ctes.Items) != 2 {
		t.Fatalf("expected 2 CTEs, got %d", len(wc.Ctes.Items))
	}

	cte1, ok := wc.Ctes.Items[0].(*nodes.CommonTableExpr)
	if !ok {
		t.Fatalf("expected *nodes.CommonTableExpr for cte1, got %T", wc.Ctes.Items[0])
	}
	if cte1.Ctename != "cte1" {
		t.Errorf("expected CTE name 'cte1', got %q", cte1.Ctename)
	}

	cte2, ok := wc.Ctes.Items[1].(*nodes.CommonTableExpr)
	if !ok {
		t.Fatalf("expected *nodes.CommonTableExpr for cte2, got %T", wc.Ctes.Items[1])
	}
	if cte2.Ctename != "cte2" {
		t.Errorf("expected CTE name 'cte2', got %q", cte2.Ctename)
	}

	// Verify main query has 2 tables in FROM clause
	if stmt.FromClause == nil || len(stmt.FromClause.Items) != 2 {
		t.Fatalf("expected 2 tables in FROM clause, got %d", len(stmt.FromClause.Items))
	}
}

// TestParseCTEWithColumnNames tests parsing a CTE with column name list.
func TestParseCTEWithColumnNames(t *testing.T) {
	input := "WITH cte(a, b) AS (SELECT 1, 2) SELECT * FROM cte"

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

	if stmt.WithClause == nil {
		t.Fatal("expected WithClause to be set")
	}

	wc := stmt.WithClause
	if wc.Ctes == nil || len(wc.Ctes.Items) != 1 {
		t.Fatalf("expected 1 CTE, got %d", len(wc.Ctes.Items))
	}

	cte, ok := wc.Ctes.Items[0].(*nodes.CommonTableExpr)
	if !ok {
		t.Fatalf("expected *nodes.CommonTableExpr, got %T", wc.Ctes.Items[0])
	}

	if cte.Ctename != "cte" {
		t.Errorf("expected CTE name 'cte', got %q", cte.Ctename)
	}

	// Verify column names
	if cte.Aliascolnames == nil || len(cte.Aliascolnames.Items) != 2 {
		t.Fatalf("expected 2 alias column names, got %v", cte.Aliascolnames)
	}

	col1, ok := cte.Aliascolnames.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String for first column name, got %T", cte.Aliascolnames.Items[0])
	}
	if col1.Str != "a" {
		t.Errorf("expected first column name 'a', got %q", col1.Str)
	}

	col2, ok := cte.Aliascolnames.Items[1].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String for second column name, got %T", cte.Aliascolnames.Items[1])
	}
	if col2.Str != "b" {
		t.Errorf("expected second column name 'b', got %q", col2.Str)
	}
}

// TestParseCTEWithOrderBy tests parsing CTE with ORDER BY on the main query.
func TestParseCTEWithOrderBy(t *testing.T) {
	input := "WITH cte AS (SELECT 1 AS n) SELECT * FROM cte ORDER BY n"

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

	// Verify WITH clause exists
	if stmt.WithClause == nil {
		t.Fatal("expected WithClause to be set")
	}

	// Verify ORDER BY clause exists
	if stmt.SortClause == nil {
		t.Fatal("expected SortClause to be set")
	}

	if len(stmt.SortClause.Items) != 1 {
		t.Fatalf("expected 1 sort item, got %d", len(stmt.SortClause.Items))
	}
}

// TestParseCTETableExpressions is a table-driven test for various CTE patterns.
func TestParseCTETableExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic CTE", "WITH cte AS (SELECT 1) SELECT * FROM cte"},
		{"recursive CTE", "WITH RECURSIVE cte AS (SELECT 1 UNION ALL SELECT n FROM cte) SELECT * FROM cte"},
		{"multiple CTEs", "WITH cte1 AS (SELECT 1), cte2 AS (SELECT 2) SELECT * FROM cte1, cte2"},
		{"CTE with column names", "WITH cte(a, b) AS (SELECT 1, 2) SELECT * FROM cte"},
		{"CTE with ORDER BY", "WITH cte AS (SELECT 1 AS n) SELECT * FROM cte ORDER BY n"},
		{"CTE with WHERE", "WITH cte AS (SELECT 1 AS n) SELECT * FROM cte WHERE n = 1"},
		{"CTE with complex subquery", "WITH cte AS (SELECT a, b FROM t WHERE a > 1) SELECT * FROM cte"},
		{"CTE used multiple times", "WITH cte AS (SELECT 1 AS n) SELECT * FROM cte AS c1, cte AS c2"},
		{"recursive CTE with UNION", "WITH RECURSIVE cte AS (SELECT 1 UNION SELECT n FROM cte) SELECT * FROM cte"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			if result == nil {
				t.Fatalf("Parse returned nil for %q", tt.input)
			}

			if len(result.Items) == 0 {
				t.Fatalf("Parse returned empty result for %q", tt.input)
			}

			stmt, ok := result.Items[0].(*nodes.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt for %q, got %T", tt.input, result.Items[0])
			}

			if stmt.WithClause == nil {
				t.Fatalf("Expected WithClause for %q", tt.input)
			}

			if stmt.WithClause.Ctes == nil || len(stmt.WithClause.Ctes.Items) == 0 {
				t.Fatalf("Expected at least one CTE for %q", tt.input)
			}

			t.Logf("Input: %s", tt.input)
			t.Logf("WithClause: Recursive=%v, CTEs=%d", stmt.WithClause.Recursive, len(stmt.WithClause.Ctes.Items))
			for i, item := range stmt.WithClause.Ctes.Items {
				cte := item.(*nodes.CommonTableExpr)
				t.Logf("  CTE[%d]: name=%q", i, cte.Ctename)
			}
		})
	}
}
