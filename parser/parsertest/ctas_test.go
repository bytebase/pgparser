package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE TABLE AS tests
// =============================================================================

func TestCreateTableAs(t *testing.T) {
	input := "CREATE TABLE t AS SELECT 1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateTableAsStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateTableAsStmt, got %T", result.Items[0])
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", stmt.Objtype)
	}
	if stmt.Into == nil {
		t.Fatal("expected non-nil Into")
	}
	if stmt.Into.Rel == nil || stmt.Into.Rel.Relname != "t" {
		t.Error("expected relation name 't'")
	}
	if stmt.Query == nil {
		t.Fatal("expected non-nil Query")
	}
}

func TestCreateTableAsFrom(t *testing.T) {
	input := "CREATE TABLE t AS SELECT * FROM other"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTableAsStmt)
	if stmt.Into.Rel.Relname != "t" {
		t.Errorf("expected 't', got %q", stmt.Into.Rel.Relname)
	}
}

func TestCreateTableIfNotExistsAs(t *testing.T) {
	input := "CREATE TABLE IF NOT EXISTS t AS SELECT 1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTableAsStmt)
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
}

func TestCreateTableAsWithData(t *testing.T) {
	input := "CREATE TABLE t AS SELECT 1 WITH DATA"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTableAsStmt)
	if stmt.Into.SkipData {
		t.Error("expected SkipData to be false")
	}
}

func TestCreateTableAsWithNoData(t *testing.T) {
	input := "CREATE TABLE t AS SELECT 1 WITH NO DATA"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTableAsStmt)
	if !stmt.Into.SkipData {
		t.Error("expected SkipData to be true")
	}
}

// =============================================================================
// CREATE MATERIALIZED VIEW tests
// =============================================================================

func TestCreateMatView(t *testing.T) {
	input := "CREATE MATERIALIZED VIEW mv AS SELECT 1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateTableAsStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateTableAsStmt, got %T", result.Items[0])
	}
	if stmt.Objtype != nodes.OBJECT_MATVIEW {
		t.Errorf("expected OBJECT_MATVIEW, got %d", stmt.Objtype)
	}
	if stmt.Into == nil || stmt.Into.Rel == nil || stmt.Into.Rel.Relname != "mv" {
		t.Error("expected relation name 'mv'")
	}
}

func TestCreateMatViewIfNotExists(t *testing.T) {
	input := "CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT 1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTableAsStmt)
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
}

func TestCreateMatViewWithNoData(t *testing.T) {
	input := "CREATE MATERIALIZED VIEW mv AS SELECT 1 WITH NO DATA"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTableAsStmt)
	if !stmt.Into.SkipData {
		t.Error("expected SkipData to be true")
	}
}

// =============================================================================
// REFRESH MATERIALIZED VIEW tests
// =============================================================================

func TestRefreshMatView(t *testing.T) {
	input := "REFRESH MATERIALIZED VIEW mv"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.RefreshMatViewStmt)
	if !ok {
		t.Fatalf("expected *nodes.RefreshMatViewStmt, got %T", result.Items[0])
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "mv" {
		t.Error("expected relation 'mv'")
	}
	if stmt.Concurrent {
		t.Error("expected Concurrent to be false")
	}
}

func TestRefreshMatViewConcurrently(t *testing.T) {
	input := "REFRESH MATERIALIZED VIEW CONCURRENTLY mv"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RefreshMatViewStmt)
	if !stmt.Concurrent {
		t.Error("expected Concurrent to be true")
	}
}

func TestRefreshMatViewWithData(t *testing.T) {
	input := "REFRESH MATERIALIZED VIEW mv WITH DATA"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RefreshMatViewStmt)
	if stmt.SkipData {
		t.Error("expected SkipData to be false")
	}
}

func TestRefreshMatViewWithNoData(t *testing.T) {
	input := "REFRESH MATERIALIZED VIEW mv WITH NO DATA"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RefreshMatViewStmt)
	if !stmt.SkipData {
		t.Error("expected SkipData to be true")
	}
}
