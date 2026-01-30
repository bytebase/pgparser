package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestParseAndSerialize tests parsing SQL and serializing the result.
func TestParseAndSerialize(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple select", "SELECT 1"},
		{"select with column", "SELECT a"},
		{"select star", "SELECT *"},
		{"select from", "SELECT * FROM users"},
		{"select with alias", "SELECT * FROM users AS u"},
		{"select with where", "SELECT * FROM users WHERE id = 1"},
		{"select with and", "SELECT * FROM users WHERE id = 1 AND name = 'test'"},
		{"select with or", "SELECT * FROM users WHERE id = 1 OR id = 2"},
		{"select with order by", "SELECT * FROM users ORDER BY id"},
		{"select with order by desc", "SELECT * FROM users ORDER BY id DESC"},
		{"select with multiple columns", "SELECT a, b, c FROM t"},
		{"select with expression", "SELECT a + b FROM t"},
		{"select with function", "SELECT count(*) FROM t"},
		{"select with null check", "SELECT * FROM t WHERE a IS NULL"},
		{"select with not null check", "SELECT * FROM t WHERE a IS NOT NULL"},
		{"select with like", "SELECT * FROM t WHERE name LIKE '%test%'"},
		{"select with between", "SELECT * FROM t WHERE a BETWEEN 1 AND 10"},
		{"select with in", "SELECT * FROM t WHERE a IN (1, 2, 3)"},
		{"select with subquery", "SELECT * FROM (SELECT 1) AS sub"},
		{"select with group by", "SELECT a, count(*) FROM t GROUP BY a"},
		{"select with group by having", "SELECT a, count(*) FROM t GROUP BY a HAVING count(*) > 1"},
		{"select with multiple group by", "SELECT a, b, sum(c) FROM t GROUP BY a, b"},
		{"union", "SELECT a FROM t1 UNION SELECT b FROM t2"},
		{"union all", "SELECT a FROM t1 UNION ALL SELECT b FROM t2"},
		{"union distinct", "SELECT a FROM t1 UNION DISTINCT SELECT b FROM t2"},
		{"intersect", "SELECT a FROM t1 INTERSECT SELECT b FROM t2"},
		{"except", "SELECT a FROM t1 EXCEPT SELECT b FROM t2"},
		{"chained union", "SELECT a FROM t1 UNION SELECT b FROM t2 UNION SELECT c FROM t3"},
		{"union with order by", "SELECT a FROM t1 UNION SELECT b FROM t2 ORDER BY a"},
		{"intersect all", "SELECT a FROM t1 INTERSECT ALL SELECT b FROM t2"},
		{"except all", "SELECT a FROM t1 EXCEPT ALL SELECT b FROM t2"},
		{"basic CTE", "WITH cte AS (SELECT 1) SELECT * FROM cte"},
		{"recursive CTE", "WITH RECURSIVE cte AS (SELECT 1 UNION ALL SELECT n FROM cte) SELECT * FROM cte"},
		{"multiple CTEs", "WITH cte1 AS (SELECT 1), cte2 AS (SELECT 2) SELECT * FROM cte1, cte2"},
		{"CTE with column names", "WITH cte(a, b) AS (SELECT 1, 2) SELECT * FROM cte"},
		{"CTE with order by", "WITH cte AS (SELECT 1 AS n) SELECT * FROM cte ORDER BY n"},
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

			// Verify we got a SelectStmt
			stmt, ok := result.Items[0].(*nodes.SelectStmt)
			if !ok {
				t.Fatalf("Expected SelectStmt for %q, got %T", tt.input, result.Items[0])
			}

			// Verify target list is not empty for simple selects,
			// or Larg/Rarg are set for set operations
			if stmt.Op == nodes.SETOP_NONE {
				if stmt.TargetList == nil || len(stmt.TargetList.Items) == 0 {
					t.Errorf("Expected non-empty target list for %q", tt.input)
				}
			} else {
				if stmt.Larg == nil || stmt.Rarg == nil {
					t.Errorf("Expected Larg and Rarg for set operation in %q", tt.input)
				}
			}

			// Serialize the result
			serialized := nodes.NodeToString(stmt)
			if serialized == "" {
				t.Errorf("Serialization returned empty string for %q", tt.input)
			}

			t.Logf("Input: %s", tt.input)
			t.Logf("Serialized: %s", serialized)
		})
	}
}

// TestParseMultipleStatements tests parsing multiple statements separated by semicolons.
func TestParseMultipleStatements(t *testing.T) {
	input := "SELECT 1; SELECT 2; SELECT 3"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil {
		t.Fatal("Parse returned nil")
	}

	if len(result.Items) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(result.Items))
	}

	for i, item := range result.Items {
		_, ok := item.(*nodes.SelectStmt)
		if !ok {
			t.Errorf("Statement %d: expected SelectStmt, got %T", i, item)
		}
	}
}

// TestParseExpressions tests parsing various expressions.
func TestParseExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"arithmetic add", "SELECT 1 + 2"},
		{"arithmetic subtract", "SELECT 3 - 1"},
		{"arithmetic multiply", "SELECT 2 * 3"},
		{"arithmetic divide", "SELECT 6 / 2"},
		{"arithmetic modulo", "SELECT 7 % 3"},
		{"comparison equal", "SELECT 1 = 1"},
		{"comparison not equal", "SELECT 1 <> 2"},
		{"comparison less", "SELECT 1 < 2"},
		{"comparison greater", "SELECT 2 > 1"},
		{"comparison less equal", "SELECT 1 <= 2"},
		{"comparison greater equal", "SELECT 2 >= 1"},
		{"unary minus", "SELECT -1"},
		{"unary plus", "SELECT +1"},
		{"boolean true", "SELECT TRUE"},
		{"boolean false", "SELECT FALSE"},
		{"null", "SELECT NULL"},
		{"string literal", "SELECT 'hello'"},
		{"float literal", "SELECT 3.14"},
		{"parenthesized", "SELECT (1 + 2) * 3"},
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

			t.Logf("Input: %s", tt.input)
			t.Logf("Result: %+v", result.Items[0])
		})
	}
}

// TestParseQualifiedNames tests parsing qualified names (schema.table).
func TestParseQualifiedNames(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		schema string
		table  string
	}{
		{"simple", "SELECT * FROM users", "", "users"},
		{"qualified", "SELECT * FROM public.users", "public", "users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			stmt := result.Items[0].(*nodes.SelectStmt)
			if stmt.FromClause == nil || len(stmt.FromClause.Items) == 0 {
				t.Fatalf("Expected from clause for %q", tt.input)
			}

			rv := stmt.FromClause.Items[0].(*nodes.RangeVar)
			if rv.Schemaname != tt.schema {
				t.Errorf("Expected schema %q, got %q", tt.schema, rv.Schemaname)
			}
			if rv.Relname != tt.table {
				t.Errorf("Expected table %q, got %q", tt.table, rv.Relname)
			}
		})
	}
}

// TestParseColumnAlias tests parsing column aliases.
func TestParseColumnAlias(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"as alias", "SELECT 1 AS num", "num"},
		{"implicit alias", "SELECT 1 num", "num"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			stmt := result.Items[0].(*nodes.SelectStmt)
			if stmt.TargetList == nil || len(stmt.TargetList.Items) == 0 {
				t.Fatalf("Expected target list for %q", tt.input)
			}

			target := stmt.TargetList.Items[0].(*nodes.ResTarget)
			if target.Name != tt.expected {
				t.Errorf("Expected alias %q, got %q", tt.expected, target.Name)
			}
		})
	}
}
