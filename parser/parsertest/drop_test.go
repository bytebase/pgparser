package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single DropStmt
func parseDropStmt(t *testing.T, input string) *nodes.DropStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to extract a single qualified name (list of string parts) from Objects at a given index
func getObjectName(t *testing.T, stmt *nodes.DropStmt, index int) []string {
	t.Helper()
	if stmt.Objects == nil || index >= len(stmt.Objects.Items) {
		t.Fatalf("expected at least %d object(s), got %v", index+1, stmt.Objects)
	}
	nameList, ok := stmt.Objects.Items[index].(*nodes.List)
	if !ok {
		t.Fatalf("expected *nodes.List for object name at index %d, got %T", index, stmt.Objects.Items[index])
	}
	var parts []string
	for _, item := range nameList.Items {
		s, ok := item.(*nodes.String)
		if !ok {
			t.Fatalf("expected *nodes.String in name list, got %T", item)
		}
		parts = append(parts, s.Str)
	}
	return parts
}

// TestDropTable tests: DROP TABLE foo
func TestDropTable(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE foo")

	if stmt.RemoveType != int(nodes.OBJECT_TABLE) {
		t.Errorf("expected RemoveType OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.RemoveType)
	}
	if stmt.Missing_ok {
		t.Error("expected Missing_ok to be false")
	}
	if stmt.Behavior != int(nodes.DROP_RESTRICT) {
		t.Errorf("expected Behavior DROP_RESTRICT (%d), got %d", nodes.DROP_RESTRICT, stmt.Behavior)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "foo" {
		t.Errorf("expected object name [foo], got %v", parts)
	}
}

// TestDropTableIfExists tests: DROP TABLE IF EXISTS foo
func TestDropTableIfExists(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE IF EXISTS foo")

	if stmt.RemoveType != int(nodes.OBJECT_TABLE) {
		t.Errorf("expected RemoveType OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.RemoveType)
	}
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if stmt.Behavior != int(nodes.DROP_RESTRICT) {
		t.Errorf("expected Behavior DROP_RESTRICT (%d), got %d", nodes.DROP_RESTRICT, stmt.Behavior)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "foo" {
		t.Errorf("expected object name [foo], got %v", parts)
	}
}

// TestDropTableCascade tests: DROP TABLE foo CASCADE
func TestDropTableCascade(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE foo CASCADE")

	if stmt.RemoveType != int(nodes.OBJECT_TABLE) {
		t.Errorf("expected RemoveType OBJECT_TABLE, got %d", stmt.RemoveType)
	}
	if stmt.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected Behavior DROP_CASCADE (%d), got %d", nodes.DROP_CASCADE, stmt.Behavior)
	}
}

// TestDropTableRestrict tests: DROP TABLE foo RESTRICT
func TestDropTableRestrict(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE foo RESTRICT")

	if stmt.RemoveType != int(nodes.OBJECT_TABLE) {
		t.Errorf("expected RemoveType OBJECT_TABLE, got %d", stmt.RemoveType)
	}
	if stmt.Behavior != int(nodes.DROP_RESTRICT) {
		t.Errorf("expected Behavior DROP_RESTRICT (%d), got %d", nodes.DROP_RESTRICT, stmt.Behavior)
	}
}

// TestDropTableQualifiedName tests: DROP TABLE myschema.foo
func TestDropTableQualifiedName(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE myschema.foo")

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 2 {
		t.Fatalf("expected 2 name parts, got %d: %v", len(parts), parts)
	}
	if parts[0] != "myschema" {
		t.Errorf("expected schema 'myschema', got %q", parts[0])
	}
	if parts[1] != "foo" {
		t.Errorf("expected table 'foo', got %q", parts[1])
	}
}

// TestDropTableMultiple tests: DROP TABLE foo, bar
func TestDropTableMultiple(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE foo, bar")

	if stmt.Objects == nil || len(stmt.Objects.Items) != 2 {
		t.Fatalf("expected 2 objects, got %v", stmt.Objects)
	}

	parts0 := getObjectName(t, stmt, 0)
	if len(parts0) != 1 || parts0[0] != "foo" {
		t.Errorf("expected first object [foo], got %v", parts0)
	}

	parts1 := getObjectName(t, stmt, 1)
	if len(parts1) != 1 || parts1[0] != "bar" {
		t.Errorf("expected second object [bar], got %v", parts1)
	}
}

// TestDropView tests: DROP VIEW myview
func TestDropView(t *testing.T) {
	stmt := parseDropStmt(t, "DROP VIEW myview")

	if stmt.RemoveType != int(nodes.OBJECT_VIEW) {
		t.Errorf("expected RemoveType OBJECT_VIEW (%d), got %d", nodes.OBJECT_VIEW, stmt.RemoveType)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "myview" {
		t.Errorf("expected object name [myview], got %v", parts)
	}
}

// TestDropIndex tests: DROP INDEX myindex
func TestDropIndex(t *testing.T) {
	stmt := parseDropStmt(t, "DROP INDEX myindex")

	if stmt.RemoveType != int(nodes.OBJECT_INDEX) {
		t.Errorf("expected RemoveType OBJECT_INDEX (%d), got %d", nodes.OBJECT_INDEX, stmt.RemoveType)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "myindex" {
		t.Errorf("expected object name [myindex], got %v", parts)
	}
}

// TestDropSequence tests: DROP SEQUENCE myseq
func TestDropSequence(t *testing.T) {
	stmt := parseDropStmt(t, "DROP SEQUENCE myseq")

	if stmt.RemoveType != int(nodes.OBJECT_SEQUENCE) {
		t.Errorf("expected RemoveType OBJECT_SEQUENCE (%d), got %d", nodes.OBJECT_SEQUENCE, stmt.RemoveType)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "myseq" {
		t.Errorf("expected object name [myseq], got %v", parts)
	}
}

// TestDropMaterializedView tests: DROP MATERIALIZED VIEW mymatview
func TestDropMaterializedView(t *testing.T) {
	stmt := parseDropStmt(t, "DROP MATERIALIZED VIEW mymatview")

	if stmt.RemoveType != int(nodes.OBJECT_MATVIEW) {
		t.Errorf("expected RemoveType OBJECT_MATVIEW (%d), got %d", nodes.OBJECT_MATVIEW, stmt.RemoveType)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "mymatview" {
		t.Errorf("expected object name [mymatview], got %v", parts)
	}
}

// TestDropTableIfExistsCascade tests: DROP TABLE IF EXISTS foo CASCADE
func TestDropTableIfExistsCascade(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE IF EXISTS foo CASCADE")

	if stmt.RemoveType != int(nodes.OBJECT_TABLE) {
		t.Errorf("expected RemoveType OBJECT_TABLE, got %d", stmt.RemoveType)
	}
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if stmt.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected Behavior DROP_CASCADE (%d), got %d", nodes.DROP_CASCADE, stmt.Behavior)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "foo" {
		t.Errorf("expected object name [foo], got %v", parts)
	}
}

// TestDropForeignTable tests: DROP FOREIGN TABLE ft
func TestDropForeignTable(t *testing.T) {
	stmt := parseDropStmt(t, "DROP FOREIGN TABLE ft")

	if stmt.RemoveType != int(nodes.OBJECT_FOREIGN_TABLE) {
		t.Errorf("expected RemoveType OBJECT_FOREIGN_TABLE (%d), got %d", nodes.OBJECT_FOREIGN_TABLE, stmt.RemoveType)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "ft" {
		t.Errorf("expected object name [ft], got %v", parts)
	}
}

// TestDropMultipleQualified tests: DROP TABLE schema1.t1, schema2.t2
func TestDropMultipleQualified(t *testing.T) {
	stmt := parseDropStmt(t, "DROP TABLE schema1.t1, schema2.t2")

	if stmt.Objects == nil || len(stmt.Objects.Items) != 2 {
		t.Fatalf("expected 2 objects, got %v", stmt.Objects)
	}

	parts0 := getObjectName(t, stmt, 0)
	if len(parts0) != 2 || parts0[0] != "schema1" || parts0[1] != "t1" {
		t.Errorf("expected first object [schema1, t1], got %v", parts0)
	}

	parts1 := getObjectName(t, stmt, 1)
	if len(parts1) != 2 || parts1[0] != "schema2" || parts1[1] != "t2" {
		t.Errorf("expected second object [schema2, t2], got %v", parts1)
	}
}

// TestDropViewIfExists tests: DROP VIEW IF EXISTS v1
func TestDropViewIfExists(t *testing.T) {
	stmt := parseDropStmt(t, "DROP VIEW IF EXISTS v1")

	if stmt.RemoveType != int(nodes.OBJECT_VIEW) {
		t.Errorf("expected RemoveType OBJECT_VIEW (%d), got %d", nodes.OBJECT_VIEW, stmt.RemoveType)
	}
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "v1" {
		t.Errorf("expected object name [v1], got %v", parts)
	}
}

// TestDropIndexIfExistsCascade tests: DROP INDEX IF EXISTS idx1 CASCADE
func TestDropIndexIfExistsCascade(t *testing.T) {
	stmt := parseDropStmt(t, "DROP INDEX IF EXISTS idx1 CASCADE")

	if stmt.RemoveType != int(nodes.OBJECT_INDEX) {
		t.Errorf("expected RemoveType OBJECT_INDEX (%d), got %d", nodes.OBJECT_INDEX, stmt.RemoveType)
	}
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if stmt.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected Behavior DROP_CASCADE (%d), got %d", nodes.DROP_CASCADE, stmt.Behavior)
	}

	parts := getObjectName(t, stmt, 0)
	if len(parts) != 1 || parts[0] != "idx1" {
		t.Errorf("expected object name [idx1], got %v", parts)
	}
}
