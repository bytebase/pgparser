package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// parseViewStmt is a helper that parses input and returns the first statement as *nodes.ViewStmt.
func parseViewStmt(t *testing.T, input string) *nodes.ViewStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.ViewStmt)
	if !ok {
		t.Fatalf("expected *nodes.ViewStmt, got %T", result.Items[0])
	}
	return stmt
}

// TestCreateViewBasic tests: CREATE VIEW v AS SELECT 1
func TestCreateViewBasic(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE VIEW v AS SELECT 1")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected View.Relname 'v', got %q", stmt.View.Relname)
	}
	if stmt.View.Schemaname != "" {
		t.Errorf("expected empty Schemaname, got %q", stmt.View.Schemaname)
	}
	if stmt.Replace {
		t.Error("expected Replace to be false")
	}
	if stmt.Aliases != nil {
		t.Errorf("expected nil Aliases, got %v", stmt.Aliases)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	if stmt.WithCheckOption != parser.VIEW_CHECK_OPTION_NONE {
		t.Errorf("expected WithCheckOption NONE (0), got %d", stmt.WithCheckOption)
	}
}

// TestCreateOrReplaceView tests: CREATE OR REPLACE VIEW v AS SELECT 1
func TestCreateOrReplaceView(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE OR REPLACE VIEW v AS SELECT 1")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected View.Relname 'v', got %q", stmt.View.Relname)
	}
	if !stmt.Replace {
		t.Error("expected Replace to be true")
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
	if stmt.WithCheckOption != parser.VIEW_CHECK_OPTION_NONE {
		t.Errorf("expected WithCheckOption NONE (0), got %d", stmt.WithCheckOption)
	}
}

// TestCreateTempView tests: CREATE TEMP VIEW v AS SELECT 1
func TestCreateTempView(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE TEMP VIEW v AS SELECT 1")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected View.Relname 'v', got %q", stmt.View.Relname)
	}
	if stmt.Replace {
		t.Error("expected Replace to be false")
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
}

// TestCreateViewWithColumns tests: CREATE VIEW v (col1, col2) AS SELECT 1, 2
func TestCreateViewWithColumns(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE VIEW v (col1, col2) AS SELECT 1, 2")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected View.Relname 'v', got %q", stmt.View.Relname)
	}
	if stmt.Aliases == nil {
		t.Fatal("expected Aliases to be set")
	}
	if len(stmt.Aliases.Items) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(stmt.Aliases.Items))
	}
	alias1, ok := stmt.Aliases.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String for alias 1, got %T", stmt.Aliases.Items[0])
	}
	if alias1.Str != "col1" {
		t.Errorf("expected alias 1 'col1', got %q", alias1.Str)
	}
	alias2, ok := stmt.Aliases.Items[1].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String for alias 2, got %T", stmt.Aliases.Items[1])
	}
	if alias2.Str != "col2" {
		t.Errorf("expected alias 2 'col2', got %q", alias2.Str)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
}

// TestCreateViewSchemaQualified tests: CREATE VIEW myschema.v AS SELECT 1
func TestCreateViewSchemaQualified(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE VIEW myschema.v AS SELECT 1")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Schemaname != "myschema" {
		t.Errorf("expected Schemaname 'myschema', got %q", stmt.View.Schemaname)
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected Relname 'v', got %q", stmt.View.Relname)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
}

// TestCreateViewWithCheckOption tests: CREATE VIEW v AS SELECT 1 WITH CHECK OPTION
func TestCreateViewWithCheckOption(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE VIEW v AS SELECT 1 WITH CHECK OPTION")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected View.Relname 'v', got %q", stmt.View.Relname)
	}
	if stmt.WithCheckOption != parser.VIEW_CHECK_OPTION_LOCAL {
		t.Errorf("expected WithCheckOption LOCAL (%d), got %d", parser.VIEW_CHECK_OPTION_LOCAL, stmt.WithCheckOption)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
}

// TestCreateViewWithCascadedCheckOption tests: CREATE VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION
func TestCreateViewWithCascadedCheckOption(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE VIEW v AS SELECT 1 WITH CASCADED CHECK OPTION")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected View.Relname 'v', got %q", stmt.View.Relname)
	}
	if stmt.WithCheckOption != parser.VIEW_CHECK_OPTION_CASCADED {
		t.Errorf("expected WithCheckOption CASCADED (%d), got %d", parser.VIEW_CHECK_OPTION_CASCADED, stmt.WithCheckOption)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
}

// TestCreateViewWithLocalCheckOption tests: CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION
func TestCreateViewWithLocalCheckOption(t *testing.T) {
	stmt := parseViewStmt(t, "CREATE VIEW v AS SELECT 1 WITH LOCAL CHECK OPTION")

	if stmt.View == nil {
		t.Fatal("expected View to be set")
	}
	if stmt.View.Relname != "v" {
		t.Errorf("expected View.Relname 'v', got %q", stmt.View.Relname)
	}
	if stmt.WithCheckOption != parser.VIEW_CHECK_OPTION_LOCAL {
		t.Errorf("expected WithCheckOption LOCAL (%d), got %d", parser.VIEW_CHECK_OPTION_LOCAL, stmt.WithCheckOption)
	}
	if stmt.Query == nil {
		t.Fatal("expected Query to be set")
	}
}
