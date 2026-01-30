package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// parseIndexStmt is a helper that parses input and returns the first statement as *nodes.IndexStmt.
func parseIndexStmt(t *testing.T, input string) *nodes.IndexStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.IndexStmt)
	if !ok {
		t.Fatalf("expected *nodes.IndexStmt, got %T", result.Items[0])
	}
	return stmt
}

// TestCreateIndexBasic tests: CREATE INDEX idx ON foo (bar)
func TestCreateIndexBasic(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo (bar)")

	if stmt.Idxname != "idx" {
		t.Errorf("expected Idxname 'idx', got %q", stmt.Idxname)
	}
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "foo" {
		t.Errorf("expected Relation.Relname 'foo', got %q", stmt.Relation.Relname)
	}
	if stmt.Unique {
		t.Error("expected Unique to be false")
	}
	if stmt.Concurrent {
		t.Error("expected Concurrent to be false")
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
	if stmt.AccessMethod != "" {
		t.Errorf("expected empty AccessMethod, got %q", stmt.AccessMethod)
	}
	if stmt.WhereClause != nil {
		t.Error("expected WhereClause to be nil")
	}

	// Check index params
	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 1 {
		t.Fatalf("expected 1 index param, got %v", stmt.IndexParams)
	}
	elem, ok := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem, got %T", stmt.IndexParams.Items[0])
	}
	if elem.Name != "bar" {
		t.Errorf("expected index elem name 'bar', got %q", elem.Name)
	}
	if elem.Ordering != nodes.SORTBY_DEFAULT {
		t.Errorf("expected SORTBY_DEFAULT, got %d", elem.Ordering)
	}
	if elem.NullsOrdering != nodes.SORTBY_NULLS_DEFAULT {
		t.Errorf("expected SORTBY_NULLS_DEFAULT, got %d", elem.NullsOrdering)
	}
}

// TestCreateUniqueIndex tests: CREATE UNIQUE INDEX idx ON foo (bar)
func TestCreateUniqueIndex(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE UNIQUE INDEX idx ON foo (bar)")

	if stmt.Idxname != "idx" {
		t.Errorf("expected Idxname 'idx', got %q", stmt.Idxname)
	}
	if !stmt.Unique {
		t.Error("expected Unique to be true")
	}
	if stmt.Concurrent {
		t.Error("expected Concurrent to be false")
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "foo" {
		t.Errorf("expected relation 'foo', got %v", stmt.Relation)
	}
}

// TestCreateIndexConcurrently tests: CREATE INDEX CONCURRENTLY idx ON foo (bar)
func TestCreateIndexConcurrently(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX CONCURRENTLY idx ON foo (bar)")

	if stmt.Idxname != "idx" {
		t.Errorf("expected Idxname 'idx', got %q", stmt.Idxname)
	}
	if !stmt.Concurrent {
		t.Error("expected Concurrent to be true")
	}
	if stmt.Unique {
		t.Error("expected Unique to be false")
	}
}

// TestCreateIndexIfNotExists tests: CREATE INDEX IF NOT EXISTS idx ON foo (bar)
func TestCreateIndexIfNotExists(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX IF NOT EXISTS idx ON foo (bar)")

	if stmt.Idxname != "idx" {
		t.Errorf("expected Idxname 'idx', got %q", stmt.Idxname)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	if stmt.Unique {
		t.Error("expected Unique to be false")
	}
	if stmt.Concurrent {
		t.Error("expected Concurrent to be false")
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "foo" {
		t.Errorf("expected relation 'foo', got %v", stmt.Relation)
	}
}

// TestCreateIndexMultipleColumns tests: CREATE INDEX idx ON foo (bar, baz)
func TestCreateIndexMultipleColumns(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo (bar, baz)")

	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 2 {
		t.Fatalf("expected 2 index params, got %v", stmt.IndexParams)
	}

	elem1, ok := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem for first param, got %T", stmt.IndexParams.Items[0])
	}
	if elem1.Name != "bar" {
		t.Errorf("expected first elem name 'bar', got %q", elem1.Name)
	}

	elem2, ok := stmt.IndexParams.Items[1].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem for second param, got %T", stmt.IndexParams.Items[1])
	}
	if elem2.Name != "baz" {
		t.Errorf("expected second elem name 'baz', got %q", elem2.Name)
	}
}

// TestCreateIndexAsc tests: CREATE INDEX idx ON foo (bar ASC)
func TestCreateIndexAsc(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo (bar ASC)")

	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 1 {
		t.Fatalf("expected 1 index param, got %v", stmt.IndexParams)
	}
	elem, ok := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem, got %T", stmt.IndexParams.Items[0])
	}
	if elem.Name != "bar" {
		t.Errorf("expected elem name 'bar', got %q", elem.Name)
	}
	if elem.Ordering != nodes.SORTBY_ASC {
		t.Errorf("expected SORTBY_ASC (%d), got %d", nodes.SORTBY_ASC, elem.Ordering)
	}
}

// TestCreateIndexDescNullsFirst tests: CREATE INDEX idx ON foo (bar DESC NULLS FIRST)
func TestCreateIndexDescNullsFirst(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo (bar DESC NULLS FIRST)")

	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 1 {
		t.Fatalf("expected 1 index param, got %v", stmt.IndexParams)
	}
	elem, ok := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem, got %T", stmt.IndexParams.Items[0])
	}
	if elem.Ordering != nodes.SORTBY_DESC {
		t.Errorf("expected SORTBY_DESC (%d), got %d", nodes.SORTBY_DESC, elem.Ordering)
	}
	if elem.NullsOrdering != nodes.SORTBY_NULLS_FIRST {
		t.Errorf("expected SORTBY_NULLS_FIRST (%d), got %d", nodes.SORTBY_NULLS_FIRST, elem.NullsOrdering)
	}
}

// TestCreateIndexPartial tests: CREATE INDEX idx ON foo (bar) WHERE bar > 10
func TestCreateIndexPartial(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo (bar) WHERE bar > 10")

	if stmt.Idxname != "idx" {
		t.Errorf("expected Idxname 'idx', got %q", stmt.Idxname)
	}
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set for partial index")
	}
	// The WHERE clause should be an A_Expr for "bar > 10"
	aexpr, ok := stmt.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr for WhereClause, got %T", stmt.WhereClause)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
}

// TestCreateIndexUsing tests: CREATE INDEX idx ON foo USING btree (bar)
func TestCreateIndexUsing(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo USING btree (bar)")

	if stmt.AccessMethod != "btree" {
		t.Errorf("expected AccessMethod 'btree', got %q", stmt.AccessMethod)
	}
	if stmt.Idxname != "idx" {
		t.Errorf("expected Idxname 'idx', got %q", stmt.Idxname)
	}
}

// TestCreateIndexUnnamed tests: CREATE INDEX ON foo (bar) - no name
func TestCreateIndexUnnamed(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX ON foo (bar)")

	if stmt.Idxname != "" {
		t.Errorf("expected empty Idxname for unnamed index, got %q", stmt.Idxname)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "foo" {
		t.Errorf("expected relation 'foo', got %v", stmt.Relation)
	}
	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 1 {
		t.Fatalf("expected 1 index param, got %v", stmt.IndexParams)
	}
}

// TestCreateIndexSchemaQualified tests: CREATE INDEX idx ON myschema.foo (bar)
func TestCreateIndexSchemaQualified(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON myschema.foo (bar)")

	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Schemaname != "myschema" {
		t.Errorf("expected Schemaname 'myschema', got %q", stmt.Relation.Schemaname)
	}
	if stmt.Relation.Relname != "foo" {
		t.Errorf("expected Relname 'foo', got %q", stmt.Relation.Relname)
	}
}

// TestCreateIndexExpression tests: CREATE INDEX idx ON foo ((col1 + col2))
func TestCreateIndexExpression(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo ((col1 + col2))")

	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 1 {
		t.Fatalf("expected 1 index param, got %v", stmt.IndexParams)
	}
	elem, ok := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem, got %T", stmt.IndexParams.Items[0])
	}
	// Expression index: Name should be empty, Expr should be set
	if elem.Name != "" {
		t.Errorf("expected empty Name for expression index, got %q", elem.Name)
	}
	if elem.Expr == nil {
		t.Fatal("expected Expr to be set for expression index")
	}
	// The expression should be an A_Expr for "col1 + col2"
	aexpr, ok := elem.Expr.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr for expression, got %T", elem.Expr)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
}

// TestCreateIndexNullsLast tests: CREATE INDEX idx ON foo (bar NULLS LAST)
func TestCreateIndexNullsLast(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE INDEX idx ON foo (bar NULLS LAST)")

	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 1 {
		t.Fatalf("expected 1 index param, got %v", stmt.IndexParams)
	}
	elem, ok := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem, got %T", stmt.IndexParams.Items[0])
	}
	if elem.NullsOrdering != nodes.SORTBY_NULLS_LAST {
		t.Errorf("expected SORTBY_NULLS_LAST (%d), got %d", nodes.SORTBY_NULLS_LAST, elem.NullsOrdering)
	}
}

// TestCreateUniqueConcurrentIfNotExists tests combined options:
// CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx ON foo USING hash (bar DESC NULLS FIRST) WHERE bar IS NOT NULL
func TestCreateUniqueConcurrentIfNotExists(t *testing.T) {
	stmt := parseIndexStmt(t, "CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx ON foo USING hash (bar DESC NULLS FIRST) WHERE bar IS NOT NULL")

	if !stmt.Unique {
		t.Error("expected Unique to be true")
	}
	if !stmt.Concurrent {
		t.Error("expected Concurrent to be true")
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	if stmt.Idxname != "idx" {
		t.Errorf("expected Idxname 'idx', got %q", stmt.Idxname)
	}
	if stmt.AccessMethod != "hash" {
		t.Errorf("expected AccessMethod 'hash', got %q", stmt.AccessMethod)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "foo" {
		t.Errorf("expected relation 'foo', got %v", stmt.Relation)
	}
	if stmt.WhereClause == nil {
		t.Fatal("expected WhereClause to be set")
	}

	// Check index param
	if stmt.IndexParams == nil || len(stmt.IndexParams.Items) != 1 {
		t.Fatalf("expected 1 index param, got %v", stmt.IndexParams)
	}
	elem, ok := stmt.IndexParams.Items[0].(*nodes.IndexElem)
	if !ok {
		t.Fatalf("expected *nodes.IndexElem, got %T", stmt.IndexParams.Items[0])
	}
	if elem.Ordering != nodes.SORTBY_DESC {
		t.Errorf("expected SORTBY_DESC, got %d", elem.Ordering)
	}
	if elem.NullsOrdering != nodes.SORTBY_NULLS_FIRST {
		t.Errorf("expected SORTBY_NULLS_FIRST, got %d", elem.NullsOrdering)
	}
}
