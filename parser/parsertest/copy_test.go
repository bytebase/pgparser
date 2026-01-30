package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single CopyStmt
func parseCopyStmt(t *testing.T, input string) *nodes.CopyStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.CopyStmt)
	if !ok {
		t.Fatalf("expected *nodes.CopyStmt, got %T", result.Items[0])
	}
	return stmt
}

// TestCopyFrom tests: COPY foo FROM '/tmp/data.csv'
func TestCopyFrom(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo FROM '/tmp/data.csv'")

	if !stmt.IsFrom {
		t.Error("expected IsFrom=true")
	}
	if stmt.Filename != "/tmp/data.csv" {
		t.Errorf("expected filename /tmp/data.csv, got %s", stmt.Filename)
	}
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "foo" {
		t.Errorf("expected table name 'foo', got %q", stmt.Relation.Relname)
	}
	if stmt.Options != nil {
		t.Errorf("expected no options, got %v", stmt.Options)
	}
	if stmt.Attlist != nil {
		t.Errorf("expected no column list, got %v", stmt.Attlist)
	}
	if stmt.WhereClause != nil {
		t.Errorf("expected no where clause, got %v", stmt.WhereClause)
	}
}

// TestCopyTo tests: COPY foo TO '/tmp/data.csv'
func TestCopyTo(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo TO '/tmp/data.csv'")

	if stmt.IsFrom {
		t.Error("expected IsFrom=false")
	}
	if stmt.Filename != "/tmp/data.csv" {
		t.Errorf("expected filename /tmp/data.csv, got %s", stmt.Filename)
	}
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "foo" {
		t.Errorf("expected table name 'foo', got %q", stmt.Relation.Relname)
	}
}

// TestCopyWithColumns tests: COPY foo (col1, col2) FROM '/tmp/data.csv'
func TestCopyWithColumns(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo (col1, col2) FROM '/tmp/data.csv'")

	if !stmt.IsFrom {
		t.Error("expected IsFrom=true")
	}
	if stmt.Filename != "/tmp/data.csv" {
		t.Errorf("expected filename /tmp/data.csv, got %s", stmt.Filename)
	}
	if stmt.Attlist == nil {
		t.Fatal("expected Attlist to be set")
	}
	if len(stmt.Attlist.Items) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.Attlist.Items))
	}
	col1, ok := stmt.Attlist.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", stmt.Attlist.Items[0])
	}
	if col1.Str != "col1" {
		t.Errorf("expected column 'col1', got %q", col1.Str)
	}
	col2, ok := stmt.Attlist.Items[1].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", stmt.Attlist.Items[1])
	}
	if col2.Str != "col2" {
		t.Errorf("expected column 'col2', got %q", col2.Str)
	}
}

// TestCopyFromStdin tests: COPY foo FROM STDIN
func TestCopyFromStdin(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo FROM STDIN")

	if !stmt.IsFrom {
		t.Error("expected IsFrom=true")
	}
	if stmt.Filename != "" {
		t.Errorf("expected empty filename for STDIN, got %q", stmt.Filename)
	}
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "foo" {
		t.Errorf("expected table name 'foo', got %q", stmt.Relation.Relname)
	}
}

// TestCopyToStdout tests: COPY foo TO STDOUT
func TestCopyToStdout(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo TO STDOUT")

	if stmt.IsFrom {
		t.Error("expected IsFrom=false")
	}
	if stmt.Filename != "" {
		t.Errorf("expected empty filename for STDOUT, got %q", stmt.Filename)
	}
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "foo" {
		t.Errorf("expected table name 'foo', got %q", stmt.Relation.Relname)
	}
}

// TestCopyFromWithOptions tests: COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv)
func TestCopyFromWithOptions(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv)")

	if !stmt.IsFrom {
		t.Error("expected IsFrom=true")
	}
	if stmt.Filename != "/tmp/data.csv" {
		t.Errorf("expected filename /tmp/data.csv, got %s", stmt.Filename)
	}
	if stmt.Options == nil {
		t.Fatal("expected Options to be set")
	}
	if len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %d", len(stmt.Options.Items))
	}
	de := findOption(t, stmt.Options, "format")
	if de == nil {
		t.Fatal("expected to find 'format' option")
	}
	if de.Arg == nil {
		t.Fatal("expected 'format' option to have an Arg")
	}
	strArg, ok := de.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String arg, got %T", de.Arg)
	}
	if strArg.Str != "csv" {
		t.Errorf("expected format arg 'csv', got %q", strArg.Str)
	}
}

// TestCopyFromWithMultipleOptions tests: COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv, HEADER true)
func TestCopyFromWithMultipleOptions(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv, HEADER true)")

	if !stmt.IsFrom {
		t.Error("expected IsFrom=true")
	}
	if stmt.Options == nil {
		t.Fatal("expected Options to be set")
	}
	if len(stmt.Options.Items) != 2 {
		t.Fatalf("expected 2 options, got %d", len(stmt.Options.Items))
	}

	formatOpt := findOption(t, stmt.Options, "format")
	if formatOpt == nil {
		t.Fatal("expected to find 'format' option")
	}

	headerOpt := findOption(t, stmt.Options, "header")
	if headerOpt == nil {
		t.Fatal("expected to find 'header' option")
	}
}

// TestCopySelectTo tests: COPY (SELECT * FROM foo) TO '/tmp/data.csv'
func TestCopySelectTo(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY (SELECT * FROM foo) TO '/tmp/data.csv'")

	if stmt.IsFrom {
		t.Error("expected IsFrom=false")
	}
	if stmt.Filename != "/tmp/data.csv" {
		t.Errorf("expected filename /tmp/data.csv, got %s", stmt.Filename)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}
	if stmt.Relation != nil {
		t.Errorf("expected no Relation for query-based COPY, got %v", stmt.Relation)
	}
}

// TestCopySelectToWithOptions tests: COPY (SELECT * FROM foo) TO '/tmp/data.csv' WITH (FORMAT csv)
func TestCopySelectToWithOptions(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY (SELECT * FROM foo) TO '/tmp/data.csv' WITH (FORMAT csv)")

	if stmt.IsFrom {
		t.Error("expected IsFrom=false")
	}
	if stmt.Filename != "/tmp/data.csv" {
		t.Errorf("expected filename /tmp/data.csv, got %s", stmt.Filename)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	_, ok := stmt.Query.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", stmt.Query)
	}
	if stmt.Options == nil {
		t.Fatal("expected Options to be set")
	}
	formatOpt := findOption(t, stmt.Options, "format")
	if formatOpt == nil {
		t.Fatal("expected to find 'format' option")
	}
}

// TestCopyFromWithWhere tests: COPY foo FROM '/tmp/data.csv' WHERE x > 0
func TestCopyFromWithWhere(t *testing.T) {
	stmt := parseCopyStmt(t, "COPY foo FROM '/tmp/data.csv' WHERE x > 0")

	if !stmt.IsFrom {
		t.Error("expected IsFrom=true")
	}
	if stmt.Filename != "/tmp/data.csv" {
		t.Errorf("expected filename /tmp/data.csv, got %s", stmt.Filename)
	}
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set")
	}
	// WhereClause should be an A_Expr for "x > 0"
	aexpr, ok := stmt.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", stmt.WhereClause)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
}

// TestCopyTableExpressions is a table-driven test for various COPY patterns.
func TestCopyTableExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"COPY FROM file", "COPY foo FROM '/tmp/data.csv'"},
		{"COPY TO file", "COPY foo TO '/tmp/data.csv'"},
		{"COPY with columns", "COPY foo (col1, col2) FROM '/tmp/data.csv'"},
		{"COPY FROM STDIN", "COPY foo FROM STDIN"},
		{"COPY TO STDOUT", "COPY foo TO STDOUT"},
		{"COPY FROM with FORMAT", "COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv)"},
		{"COPY FROM with multiple options", "COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv, HEADER true)"},
		{"COPY SELECT TO file", "COPY (SELECT * FROM foo) TO '/tmp/data.csv'"},
		{"COPY SELECT TO with options", "COPY (SELECT * FROM foo) TO '/tmp/data.csv' WITH (FORMAT csv)"},
		{"COPY FROM with WHERE", "COPY foo FROM '/tmp/data.csv' WHERE x > 0"},
		{"COPY FROM with options and WHERE", "COPY foo FROM '/tmp/data.csv' WITH (FORMAT csv) WHERE x > 0"},
		{"COPY without WITH keyword", "COPY foo FROM '/tmp/data.csv' (FORMAT csv)"},
		{"COPY SELECT without WITH", "COPY (SELECT 1) TO '/tmp/out.csv' (FORMAT csv)"},
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

			_, ok := result.Items[0].(*nodes.CopyStmt)
			if !ok {
				t.Fatalf("expected *nodes.CopyStmt for %q, got %T", tt.input, result.Items[0])
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}
