package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE STATISTICS tests
// =============================================================================

func TestCreateStatistics(t *testing.T) {
	input := "CREATE STATISTICS mystat ON col1, col2 FROM t"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateStatsStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateStatsStmt, got %T", result.Items[0])
	}
	if stmt.Defnames == nil {
		t.Fatal("expected non-nil Defnames")
	}
	if stmt.Exprs == nil || len(stmt.Exprs.Items) != 2 {
		t.Fatalf("expected 2 expressions, got %v", stmt.Exprs)
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
}

func TestCreateStatisticsIfNotExists(t *testing.T) {
	input := "CREATE STATISTICS IF NOT EXISTS mystat (dependencies) ON col1, col2 FROM t"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateStatsStmt)
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	if stmt.StatTypes == nil || len(stmt.StatTypes.Items) != 1 {
		t.Fatalf("expected 1 stat type, got %v", stmt.StatTypes)
	}
}

func TestCreateStatisticsMultipleTypes(t *testing.T) {
	input := "CREATE STATISTICS mystat (ndistinct, dependencies) ON col1, col2 FROM t"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateStatsStmt)
	if stmt.StatTypes == nil || len(stmt.StatTypes.Items) != 2 {
		t.Fatalf("expected 2 stat types, got %v", stmt.StatTypes)
	}
}

func TestCreateStatisticsFuncExpr(t *testing.T) {
	sqls := []string{
		`CREATE STATISTICS s1 ON lower(a) FROM t1`,
		`CREATE STATISTICS s2 (dependencies) ON lower(col1), col2 FROM t1`,
		`CREATE STATISTICS s3 ON abs(a), b, c FROM t1`,
	}
	for _, sql := range sqls {
		_, err := parser.Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", sql, err)
		}
	}
}

// =============================================================================
// ALTER STATISTICS tests
// =============================================================================

func TestAlterStatistics(t *testing.T) {
	input := "ALTER STATISTICS mystat SET STATISTICS 100"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterStatsStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterStatsStmt, got %T", result.Items[0])
	}
	if stmt.Defnames == nil {
		t.Fatal("expected non-nil Defnames")
	}
	if stmt.Stxstattarget != 100 {
		t.Errorf("expected stxstattarget 100, got %d", stmt.Stxstattarget)
	}
	if stmt.MissingOk {
		t.Error("expected MissingOk to be false")
	}
}
