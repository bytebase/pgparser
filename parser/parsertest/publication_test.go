package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE PUBLICATION tests
// =============================================================================

// TestCreatePublicationBasic tests: CREATE PUBLICATION pub
func TestCreatePublicationBasic(t *testing.T) {
	input := "CREATE PUBLICATION pub"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.CreatePublicationStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreatePublicationStmt, got %T", result.Items[0])
	}
	if stmt.Pubname != "pub" {
		t.Errorf("expected pubname 'pub', got %q", stmt.Pubname)
	}
	if stmt.ForAllTables {
		t.Error("expected ForAllTables to be false")
	}
}

// TestCreatePublicationForAllTables tests: CREATE PUBLICATION pub FOR ALL TABLES
func TestCreatePublicationForAllTables(t *testing.T) {
	input := "CREATE PUBLICATION pub FOR ALL TABLES"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePublicationStmt)
	if !stmt.ForAllTables {
		t.Error("expected ForAllTables to be true")
	}
}

// TestCreatePublicationForTable tests: CREATE PUBLICATION pub FOR TABLE t1
func TestCreatePublicationForTable(t *testing.T) {
	input := "CREATE PUBLICATION pub FOR TABLE t1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePublicationStmt)
	if stmt.Pubobjects == nil || stmt.Pubobjects.Len() != 1 {
		t.Fatalf("expected 1 pub object, got %d", stmt.Pubobjects.Len())
	}
	obj := stmt.Pubobjects.Items[0].(*nodes.PublicationObjSpec)
	if obj.Pubobjtype != nodes.PUBLICATIONOBJ_TABLE {
		t.Errorf("expected PUBLICATIONOBJ_TABLE, got %d", obj.Pubobjtype)
	}
	if obj.Pubtable == nil || obj.Pubtable.Relation == nil {
		t.Fatal("expected pubtable with relation")
	}
	if obj.Pubtable.Relation.Relname != "t1" {
		t.Errorf("expected relation 't1', got %q", obj.Pubtable.Relation.Relname)
	}
}

// TestCreatePublicationForMultipleTables tests: CREATE PUBLICATION pub FOR TABLE t1, TABLE t2
func TestCreatePublicationForMultipleTables(t *testing.T) {
	input := "CREATE PUBLICATION pub FOR TABLE t1, TABLE t2"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePublicationStmt)
	if stmt.Pubobjects == nil || stmt.Pubobjects.Len() != 2 {
		t.Fatalf("expected 2 pub objects, got %d", stmt.Pubobjects.Len())
	}
}

// TestCreatePublicationForTableWhere tests: CREATE PUBLICATION pub FOR TABLE t1 WHERE (col > 0)
func TestCreatePublicationForTableWhere(t *testing.T) {
	input := "CREATE PUBLICATION pub FOR TABLE t1 WHERE (col > 0)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePublicationStmt)
	if stmt.Pubobjects == nil || stmt.Pubobjects.Len() != 1 {
		t.Fatalf("expected 1 pub object, got %d", stmt.Pubobjects.Len())
	}
	obj := stmt.Pubobjects.Items[0].(*nodes.PublicationObjSpec)
	if obj.Pubtable == nil || obj.Pubtable.WhereClause == nil {
		t.Error("expected WHERE clause to be non-nil")
	}
}

// TestCreatePublicationForTablesInSchema tests: CREATE PUBLICATION pub FOR TABLES IN SCHEMA public
func TestCreatePublicationForTablesInSchema(t *testing.T) {
	input := "CREATE PUBLICATION pub FOR TABLES IN SCHEMA public"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePublicationStmt)
	if stmt.Pubobjects == nil || stmt.Pubobjects.Len() != 1 {
		t.Fatalf("expected 1 pub object, got %d", stmt.Pubobjects.Len())
	}
	obj := stmt.Pubobjects.Items[0].(*nodes.PublicationObjSpec)
	if obj.Pubobjtype != nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA {
		t.Errorf("expected PUBLICATIONOBJ_TABLES_IN_SCHEMA, got %d", obj.Pubobjtype)
	}
	if obj.Name != "public" {
		t.Errorf("expected schema name 'public', got %q", obj.Name)
	}
}

// TestCreatePublicationWithOptions tests: CREATE PUBLICATION pub FOR ALL TABLES WITH (publish = 'insert')
func TestCreatePublicationWithOptions(t *testing.T) {
	input := "CREATE PUBLICATION pub FOR ALL TABLES WITH (publish = 'insert')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePublicationStmt)
	if !stmt.ForAllTables {
		t.Error("expected ForAllTables to be true")
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "publish" {
		t.Errorf("expected defname 'publish', got %q", opt.Defname)
	}
}

// =============================================================================
// ALTER PUBLICATION tests
// =============================================================================

// TestAlterPublicationSetOptions tests: ALTER PUBLICATION pub SET (publish = 'insert,update')
func TestAlterPublicationSetOptions(t *testing.T) {
	input := "ALTER PUBLICATION pub SET (publish = 'insert,update')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterPublicationStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterPublicationStmt, got %T", result.Items[0])
	}
	if stmt.Pubname != "pub" {
		t.Errorf("expected pubname 'pub', got %q", stmt.Pubname)
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
}

// TestAlterPublicationAddTable tests: ALTER PUBLICATION pub ADD TABLE t3
func TestAlterPublicationAddTable(t *testing.T) {
	input := "ALTER PUBLICATION pub ADD TABLE t3"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterPublicationStmt)
	if stmt.Pubobjects == nil || stmt.Pubobjects.Len() != 1 {
		t.Fatalf("expected 1 pub object, got %d", stmt.Pubobjects.Len())
	}
	if stmt.Action != nodes.DEFELEM_ADD {
		t.Errorf("expected action DEFELEM_ADD, got %d", stmt.Action)
	}
}

// TestAlterPublicationDropTable tests: ALTER PUBLICATION pub DROP TABLE t1
func TestAlterPublicationDropTable(t *testing.T) {
	input := "ALTER PUBLICATION pub DROP TABLE t1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterPublicationStmt)
	if stmt.Action != nodes.DEFELEM_DROP {
		t.Errorf("expected action DEFELEM_DROP, got %d", stmt.Action)
	}
}
