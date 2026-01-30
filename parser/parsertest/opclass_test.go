package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE OPERATOR CLASS tests
// =============================================================================

func TestCreateOpClass(t *testing.T) {
	input := "CREATE OPERATOR CLASS myclass DEFAULT FOR TYPE int4 USING btree AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int4, int4)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateOpClassStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateOpClassStmt, got %T", result.Items[0])
	}
	if !stmt.IsDefault {
		t.Error("expected IsDefault to be true")
	}
	if stmt.Amname != "btree" {
		t.Errorf("expected amname 'btree', got %q", stmt.Amname)
	}
	if stmt.Items == nil || len(stmt.Items.Items) < 2 {
		t.Fatal("expected at least 2 items")
	}
}

func TestCreateOpClassNonDefault(t *testing.T) {
	input := "CREATE OPERATOR CLASS myclass FOR TYPE int4 USING btree AS OPERATOR 1 <"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateOpClassStmt)
	if stmt.IsDefault {
		t.Error("expected IsDefault to be false")
	}
}

// =============================================================================
// CREATE OPERATOR FAMILY tests
// =============================================================================

func TestCreateOpFamily(t *testing.T) {
	input := "CREATE OPERATOR FAMILY myfam USING btree"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateOpFamilyStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateOpFamilyStmt, got %T", result.Items[0])
	}
	if stmt.Amname != "btree" {
		t.Errorf("expected amname 'btree', got %q", stmt.Amname)
	}
}

// =============================================================================
// ALTER OPERATOR FAMILY tests
// =============================================================================

func TestAlterOpFamilyAdd(t *testing.T) {
	input := "ALTER OPERATOR FAMILY myfam USING btree ADD OPERATOR 1 <(int4, int4)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterOpFamilyStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterOpFamilyStmt, got %T", result.Items[0])
	}
	if stmt.IsDrop {
		t.Error("expected IsDrop to be false")
	}
	if stmt.Amname != "btree" {
		t.Errorf("expected amname 'btree', got %q", stmt.Amname)
	}
}

func TestAlterOpFamilyDrop(t *testing.T) {
	input := "ALTER OPERATOR FAMILY myfam USING btree DROP OPERATOR 1 (int4, int4)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterOpFamilyStmt)
	if !stmt.IsDrop {
		t.Error("expected IsDrop to be true")
	}
}

// =============================================================================
// DROP OPERATOR CLASS/FAMILY tests
// =============================================================================

func TestDropOpClass(t *testing.T) {
	input := "DROP OPERATOR CLASS myclass USING btree"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	if stmt.RemoveType != int(nodes.OBJECT_OPCLASS) {
		t.Errorf("expected OBJECT_OPCLASS, got %d", stmt.RemoveType)
	}
}

func TestDropOpClassIfExists(t *testing.T) {
	input := "DROP OPERATOR CLASS IF EXISTS myclass USING btree CASCADE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if stmt.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected DROP_CASCADE, got %d", stmt.Behavior)
	}
}

func TestDropOpFamily(t *testing.T) {
	input := "DROP OPERATOR FAMILY myfam USING btree"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if stmt.RemoveType != int(nodes.OBJECT_OPFAMILY) {
		t.Errorf("expected OBJECT_OPFAMILY, got %d", stmt.RemoveType)
	}
}
