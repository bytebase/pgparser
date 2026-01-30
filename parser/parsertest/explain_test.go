package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single ExplainStmt
func parseExplainStmt(t *testing.T, input string) *nodes.ExplainStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.ExplainStmt)
	if !ok {
		t.Fatalf("expected *nodes.ExplainStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to find a DefElem by name in an options list
func findOption(t *testing.T, opts *nodes.List, name string) *nodes.DefElem {
	t.Helper()
	if opts == nil {
		return nil
	}
	for _, item := range opts.Items {
		de, ok := item.(*nodes.DefElem)
		if !ok {
			t.Fatalf("expected *nodes.DefElem in options, got %T", item)
		}
		if de.Defname == name {
			return de
		}
	}
	return nil
}

// TestExplainSelect tests: EXPLAIN SELECT 1
func TestExplainSelect(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN SELECT 1")

	// Query should be a SelectStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}

	// No options for plain EXPLAIN
	if stmt.Options != nil && len(stmt.Options.Items) > 0 {
		t.Errorf("expected no options for plain EXPLAIN, got %d", len(stmt.Options.Items))
	}
}

// TestExplainAnalyze tests: EXPLAIN ANALYZE SELECT 1
func TestExplainAnalyze(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN ANALYZE SELECT 1")

	// Query should be a SelectStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}

	// Options should contain "analyze"
	if stmt.Options == nil || len(stmt.Options.Items) == 0 {
		t.Fatal("expected options to contain 'analyze'")
	}

	de := findOption(t, stmt.Options, "analyze")
	if de == nil {
		t.Fatal("expected to find 'analyze' option")
	}
}

// TestExplainVerbose tests: EXPLAIN VERBOSE SELECT 1
func TestExplainVerbose(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN VERBOSE SELECT 1")

	// Query should be a SelectStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}

	// Options should contain "verbose"
	if stmt.Options == nil || len(stmt.Options.Items) == 0 {
		t.Fatal("expected options to contain 'verbose'")
	}

	de := findOption(t, stmt.Options, "verbose")
	if de == nil {
		t.Fatal("expected to find 'verbose' option")
	}
}

// TestExplainAnalyzeVerbose tests: EXPLAIN ANALYZE VERBOSE SELECT 1
func TestExplainAnalyzeVerbose(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN ANALYZE VERBOSE SELECT 1")

	// Query should be a SelectStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}

	// Options should contain both "analyze" and "verbose"
	if stmt.Options == nil || len(stmt.Options.Items) < 2 {
		t.Fatalf("expected at least 2 options, got %v", stmt.Options)
	}

	analyzeOpt := findOption(t, stmt.Options, "analyze")
	if analyzeOpt == nil {
		t.Fatal("expected to find 'analyze' option")
	}

	verboseOpt := findOption(t, stmt.Options, "verbose")
	if verboseOpt == nil {
		t.Fatal("expected to find 'verbose' option")
	}
}

// TestExplainParenthesizedOptions tests: EXPLAIN (ANALYZE, VERBOSE) SELECT 1
func TestExplainParenthesizedOptions(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN (ANALYZE, VERBOSE) SELECT 1")

	// Query should be a SelectStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}

	// Options should contain both "analyze" and "verbose"
	if stmt.Options == nil || len(stmt.Options.Items) < 2 {
		t.Fatalf("expected at least 2 options, got %v", stmt.Options)
	}

	analyzeOpt := findOption(t, stmt.Options, "analyze")
	if analyzeOpt == nil {
		t.Fatal("expected to find 'analyze' option")
	}

	verboseOpt := findOption(t, stmt.Options, "verbose")
	if verboseOpt == nil {
		t.Fatal("expected to find 'verbose' option")
	}
}

// TestExplainFormatJSON tests: EXPLAIN (FORMAT JSON) SELECT 1
func TestExplainFormatJSON(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN (FORMAT JSON) SELECT 1")

	// Query should be a SelectStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}

	// Options should contain "format" with arg "json"
	if stmt.Options == nil || len(stmt.Options.Items) == 0 {
		t.Fatal("expected options to be set")
	}

	formatOpt := findOption(t, stmt.Options, "format")
	if formatOpt == nil {
		t.Fatal("expected to find 'format' option")
	}

	// The arg should be a String with value "json"
	if formatOpt.Arg == nil {
		t.Fatal("expected 'format' option to have an Arg")
	}
	strArg, ok := formatOpt.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String arg for format, got %T", formatOpt.Arg)
	}
	if strArg.Str != "json" {
		t.Errorf("expected format arg 'json', got %q", strArg.Str)
	}
}

// TestExplainInsert tests: EXPLAIN INSERT INTO foo VALUES (1)
func TestExplainInsert(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN INSERT INTO foo VALUES (1)")

	// Query should be an InsertStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	insertStmt, ok := stmt.Query.(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", stmt.Query)
	}

	// Verify the relation
	if insertStmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if insertStmt.Relation.Relname != "foo" {
		t.Errorf("expected table name 'foo', got %q", insertStmt.Relation.Relname)
	}

	// No options for plain EXPLAIN
	if stmt.Options != nil && len(stmt.Options.Items) > 0 {
		t.Errorf("expected no options for plain EXPLAIN, got %d", len(stmt.Options.Items))
	}
}

// TestExplainDelete tests: EXPLAIN DELETE FROM foo
func TestExplainDelete(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN DELETE FROM foo")

	// Query should be a DeleteStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	deleteStmt, ok := stmt.Query.(*nodes.DeleteStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeleteStmt, got %T", stmt.Query)
	}

	// Verify the relation
	if deleteStmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if deleteStmt.Relation.Relname != "foo" {
		t.Errorf("expected table name 'foo', got %q", deleteStmt.Relation.Relname)
	}

	// No options for plain EXPLAIN
	if stmt.Options != nil && len(stmt.Options.Items) > 0 {
		t.Errorf("expected no options for plain EXPLAIN, got %d", len(stmt.Options.Items))
	}
}

// TestExplainUpdate tests: EXPLAIN UPDATE foo SET x = 1
func TestExplainUpdate(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN UPDATE foo SET x = 1")

	// Query should be an UpdateStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	updateStmt, ok := stmt.Query.(*nodes.UpdateStmt)
	if !ok {
		t.Fatalf("expected *nodes.UpdateStmt, got %T", stmt.Query)
	}

	// Verify the relation
	if updateStmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if updateStmt.Relation.Relname != "foo" {
		t.Errorf("expected table name 'foo', got %q", updateStmt.Relation.Relname)
	}

	// Verify TargetList (SET clause) has 1 item
	if updateStmt.TargetList == nil || len(updateStmt.TargetList.Items) != 1 {
		t.Fatalf("expected 1 item in TargetList, got %v", updateStmt.TargetList)
	}

	rt, ok := updateStmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", updateStmt.TargetList.Items[0])
	}
	if rt.Name != "x" {
		t.Errorf("expected SET target 'x', got %q", rt.Name)
	}

	// No options for plain EXPLAIN
	if stmt.Options != nil && len(stmt.Options.Items) > 0 {
		t.Errorf("expected no options for plain EXPLAIN, got %d", len(stmt.Options.Items))
	}
}

// TestExplainAnalyzeInsert tests: EXPLAIN ANALYZE INSERT INTO foo VALUES (1)
func TestExplainAnalyzeInsert(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN ANALYZE INSERT INTO foo VALUES (1)")

	// Query should be an InsertStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected *nodes.InsertStmt, got %T", stmt.Query)
	}

	// Options should contain "analyze"
	if stmt.Options == nil || len(stmt.Options.Items) == 0 {
		t.Fatal("expected options to contain 'analyze'")
	}

	de := findOption(t, stmt.Options, "analyze")
	if de == nil {
		t.Fatal("expected to find 'analyze' option")
	}
}

// TestExplainParenthesizedAnalyzeVerbose tests: EXPLAIN (ANALYZE true, VERBOSE true) SELECT 1
func TestExplainParenthesizedAnalyzeVerbose(t *testing.T) {
	stmt := parseExplainStmt(t, "EXPLAIN (ANALYZE true, VERBOSE true) SELECT 1")

	// Query should be a SelectStmt
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}

	// Options should contain "analyze" and "verbose" with string args
	if stmt.Options == nil || len(stmt.Options.Items) < 2 {
		t.Fatalf("expected at least 2 options, got %v", stmt.Options)
	}

	analyzeOpt := findOption(t, stmt.Options, "analyze")
	if analyzeOpt == nil {
		t.Fatal("expected to find 'analyze' option")
	}
	if analyzeOpt.Arg == nil {
		t.Fatal("expected 'analyze' option to have an Arg")
	}

	verboseOpt := findOption(t, stmt.Options, "verbose")
	if verboseOpt == nil {
		t.Fatal("expected to find 'verbose' option")
	}
	if verboseOpt.Arg == nil {
		t.Fatal("expected 'verbose' option to have an Arg")
	}
}

// TestExplainTableExpressions is a table-driven test for various EXPLAIN patterns.
func TestExplainTableExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"plain EXPLAIN SELECT", "EXPLAIN SELECT 1"},
		{"EXPLAIN ANALYZE SELECT", "EXPLAIN ANALYZE SELECT 1"},
		{"EXPLAIN VERBOSE SELECT", "EXPLAIN VERBOSE SELECT 1"},
		{"EXPLAIN ANALYZE VERBOSE SELECT", "EXPLAIN ANALYZE VERBOSE SELECT 1"},
		{"EXPLAIN parenthesized ANALYZE VERBOSE", "EXPLAIN (ANALYZE, VERBOSE) SELECT 1"},
		{"EXPLAIN FORMAT JSON", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"EXPLAIN FORMAT YAML", "EXPLAIN (FORMAT YAML) SELECT 1"},
		{"EXPLAIN FORMAT TEXT", "EXPLAIN (FORMAT TEXT) SELECT 1"},
		{"EXPLAIN FORMAT XML", "EXPLAIN (FORMAT XML) SELECT 1"},
		{"EXPLAIN COSTS off", "EXPLAIN (COSTS off) SELECT 1"},
		{"EXPLAIN BUFFERS", "EXPLAIN (BUFFERS) SELECT 1"},
		{"EXPLAIN TIMING", "EXPLAIN (TIMING) SELECT 1"},
		{"EXPLAIN INSERT", "EXPLAIN INSERT INTO foo VALUES (1)"},
		{"EXPLAIN DELETE", "EXPLAIN DELETE FROM foo"},
		{"EXPLAIN UPDATE", "EXPLAIN UPDATE foo SET x = 1"},
		{"EXPLAIN ANALYZE INSERT", "EXPLAIN ANALYZE INSERT INTO foo VALUES (1)"},
		{"EXPLAIN ANALYZE DELETE", "EXPLAIN ANALYZE DELETE FROM foo"},
		{"EXPLAIN ANALYZE UPDATE", "EXPLAIN ANALYZE UPDATE foo SET x = 1"},
		{"EXPLAIN with multiple options", "EXPLAIN (ANALYZE true, VERBOSE true, FORMAT JSON) SELECT 1"},
		{"EXPLAIN SELECT with WHERE", "EXPLAIN SELECT * FROM foo WHERE id = 1"},
		{"EXPLAIN SELECT with JOIN", "EXPLAIN SELECT * FROM foo JOIN bar ON foo.id = bar.id"},
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

			_, ok := result.Items[0].(*nodes.ExplainStmt)
			if !ok {
				t.Fatalf("expected *nodes.ExplainStmt for %q, got %T", tt.input, result.Items[0])
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}
