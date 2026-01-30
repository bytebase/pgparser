package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
)

/*****************************************************************************
 *
 *      TRUNCATE tests
 *
 *****************************************************************************/

func TestTruncateSimple(t *testing.T) {
	stmt := parseStmt(t, "TRUNCATE t")
	ts, ok := stmt.(*nodes.TruncateStmt)
	if !ok {
		t.Fatalf("expected *nodes.TruncateStmt, got %T", stmt)
	}
	if ts.Relations == nil || ts.Relations.Len() != 1 {
		t.Fatalf("expected 1 relation, got %v", ts.Relations)
	}
	rv := ts.Relations.Items[0].(*nodes.RangeVar)
	if rv.Relname != "t" {
		t.Errorf("expected relname 't', got %q", rv.Relname)
	}
	if ts.RestartSeqs != false {
		t.Errorf("expected RestartSeqs false, got %v", ts.RestartSeqs)
	}
	if ts.Behavior != nodes.DROP_RESTRICT {
		t.Errorf("expected Behavior DROP_RESTRICT, got %v", ts.Behavior)
	}
}

func TestTruncateTable(t *testing.T) {
	stmt := parseStmt(t, "TRUNCATE TABLE t")
	ts, ok := stmt.(*nodes.TruncateStmt)
	if !ok {
		t.Fatalf("expected *nodes.TruncateStmt, got %T", stmt)
	}
	if ts.Relations == nil || ts.Relations.Len() != 1 {
		t.Fatalf("expected 1 relation, got %v", ts.Relations)
	}
}

func TestTruncateMultiple(t *testing.T) {
	stmt := parseStmt(t, "TRUNCATE t, t2")
	ts, ok := stmt.(*nodes.TruncateStmt)
	if !ok {
		t.Fatalf("expected *nodes.TruncateStmt, got %T", stmt)
	}
	if ts.Relations == nil || ts.Relations.Len() != 2 {
		t.Fatalf("expected 2 relations, got %v", ts.Relations)
	}
}

func TestTruncateCascade(t *testing.T) {
	stmt := parseStmt(t, "TRUNCATE t CASCADE")
	ts, ok := stmt.(*nodes.TruncateStmt)
	if !ok {
		t.Fatalf("expected *nodes.TruncateStmt, got %T", stmt)
	}
	if ts.Behavior != nodes.DROP_CASCADE {
		t.Errorf("expected Behavior DROP_CASCADE, got %v", ts.Behavior)
	}
}

func TestTruncateRestrict(t *testing.T) {
	stmt := parseStmt(t, "TRUNCATE t RESTRICT")
	ts, ok := stmt.(*nodes.TruncateStmt)
	if !ok {
		t.Fatalf("expected *nodes.TruncateStmt, got %T", stmt)
	}
	if ts.Behavior != nodes.DROP_RESTRICT {
		t.Errorf("expected Behavior DROP_RESTRICT, got %v", ts.Behavior)
	}
}

func TestTruncateRestartIdentity(t *testing.T) {
	stmt := parseStmt(t, "TRUNCATE t RESTART IDENTITY")
	ts, ok := stmt.(*nodes.TruncateStmt)
	if !ok {
		t.Fatalf("expected *nodes.TruncateStmt, got %T", stmt)
	}
	if ts.RestartSeqs != true {
		t.Errorf("expected RestartSeqs true, got %v", ts.RestartSeqs)
	}
}

func TestTruncateContinueIdentity(t *testing.T) {
	stmt := parseStmt(t, "TRUNCATE t CONTINUE IDENTITY")
	ts, ok := stmt.(*nodes.TruncateStmt)
	if !ok {
		t.Fatalf("expected *nodes.TruncateStmt, got %T", stmt)
	}
	if ts.RestartSeqs != false {
		t.Errorf("expected RestartSeqs false, got %v", ts.RestartSeqs)
	}
}

/*****************************************************************************
 *
 *      LOCK tests
 *
 *****************************************************************************/

func TestLockTable(t *testing.T) {
	stmt := parseStmt(t, "LOCK TABLE t")
	ls, ok := stmt.(*nodes.LockStmt)
	if !ok {
		t.Fatalf("expected *nodes.LockStmt, got %T", stmt)
	}
	if ls.Relations == nil || ls.Relations.Len() != 1 {
		t.Fatalf("expected 1 relation, got %v", ls.Relations)
	}
	rv := ls.Relations.Items[0].(*nodes.RangeVar)
	if rv.Relname != "t" {
		t.Errorf("expected relname 't', got %q", rv.Relname)
	}
	if ls.Mode != nodes.AccessExclusiveLock {
		t.Errorf("expected mode AccessExclusiveLock (%d), got %d", nodes.AccessExclusiveLock, ls.Mode)
	}
	if ls.Nowait != false {
		t.Errorf("expected Nowait false, got %v", ls.Nowait)
	}
}

func TestLockNoTableKeyword(t *testing.T) {
	stmt := parseStmt(t, "LOCK t")
	ls, ok := stmt.(*nodes.LockStmt)
	if !ok {
		t.Fatalf("expected *nodes.LockStmt, got %T", stmt)
	}
	if ls.Relations == nil || ls.Relations.Len() != 1 {
		t.Fatalf("expected 1 relation, got %v", ls.Relations)
	}
}

func TestLockAccessShare(t *testing.T) {
	stmt := parseStmt(t, "LOCK TABLE t IN ACCESS SHARE MODE")
	ls, ok := stmt.(*nodes.LockStmt)
	if !ok {
		t.Fatalf("expected *nodes.LockStmt, got %T", stmt)
	}
	if ls.Mode != nodes.AccessShareLock {
		t.Errorf("expected mode AccessShareLock (%d), got %d", nodes.AccessShareLock, ls.Mode)
	}
}

func TestLockRowExclusive(t *testing.T) {
	stmt := parseStmt(t, "LOCK TABLE t IN ROW EXCLUSIVE MODE")
	ls, ok := stmt.(*nodes.LockStmt)
	if !ok {
		t.Fatalf("expected *nodes.LockStmt, got %T", stmt)
	}
	if ls.Mode != nodes.RowExclusiveLock {
		t.Errorf("expected mode RowExclusiveLock (%d), got %d", nodes.RowExclusiveLock, ls.Mode)
	}
}

func TestLockExclusiveNowait(t *testing.T) {
	stmt := parseStmt(t, "LOCK TABLE t IN EXCLUSIVE MODE NOWAIT")
	ls, ok := stmt.(*nodes.LockStmt)
	if !ok {
		t.Fatalf("expected *nodes.LockStmt, got %T", stmt)
	}
	if ls.Mode != nodes.ExclusiveLock {
		t.Errorf("expected mode ExclusiveLock (%d), got %d", nodes.ExclusiveLock, ls.Mode)
	}
	if ls.Nowait != true {
		t.Errorf("expected Nowait true, got %v", ls.Nowait)
	}
}

func TestLockMultipleShareMode(t *testing.T) {
	stmt := parseStmt(t, "LOCK TABLE t, t2 IN SHARE MODE")
	ls, ok := stmt.(*nodes.LockStmt)
	if !ok {
		t.Fatalf("expected *nodes.LockStmt, got %T", stmt)
	}
	if ls.Relations == nil || ls.Relations.Len() != 2 {
		t.Fatalf("expected 2 relations, got %v", ls.Relations)
	}
	if ls.Mode != nodes.ShareLock {
		t.Errorf("expected mode ShareLock (%d), got %d", nodes.ShareLock, ls.Mode)
	}
}

/*****************************************************************************
 *
 *      VACUUM tests
 *
 *****************************************************************************/

func TestVacuumSimple(t *testing.T) {
	stmt := parseStmt(t, "VACUUM")
	vs, ok := stmt.(*nodes.VacuumStmt)
	if !ok {
		t.Fatalf("expected *nodes.VacuumStmt, got %T", stmt)
	}
	if vs.IsVacuumCmd != true {
		t.Errorf("expected IsVacuumCmd true, got %v", vs.IsVacuumCmd)
	}
	if vs.Rels != nil && vs.Rels.Len() > 0 {
		t.Errorf("expected nil Rels, got %v", vs.Rels)
	}
}

func TestVacuumTable(t *testing.T) {
	stmt := parseStmt(t, "VACUUM t")
	vs, ok := stmt.(*nodes.VacuumStmt)
	if !ok {
		t.Fatalf("expected *nodes.VacuumStmt, got %T", stmt)
	}
	if vs.IsVacuumCmd != true {
		t.Errorf("expected IsVacuumCmd true, got %v", vs.IsVacuumCmd)
	}
	if vs.Rels == nil || vs.Rels.Len() != 1 {
		t.Fatalf("expected 1 relation, got %v", vs.Rels)
	}
}

func TestVacuumVerbose(t *testing.T) {
	stmt := parseStmt(t, "VACUUM (VERBOSE) t")
	vs, ok := stmt.(*nodes.VacuumStmt)
	if !ok {
		t.Fatalf("expected *nodes.VacuumStmt, got %T", stmt)
	}
	if vs.Options == nil || vs.Options.Len() < 1 {
		t.Fatalf("expected at least 1 option, got %v", vs.Options)
	}
}

func TestVacuumAnalyze(t *testing.T) {
	stmt := parseStmt(t, "VACUUM (ANALYZE) t")
	vs, ok := stmt.(*nodes.VacuumStmt)
	if !ok {
		t.Fatalf("expected *nodes.VacuumStmt, got %T", stmt)
	}
	if vs.IsVacuumCmd != true {
		t.Errorf("expected IsVacuumCmd true, got %v", vs.IsVacuumCmd)
	}
}

/*****************************************************************************
 *
 *      ANALYZE tests
 *
 *****************************************************************************/

func TestAnalyzeSimple(t *testing.T) {
	stmt := parseStmt(t, "ANALYZE")
	vs, ok := stmt.(*nodes.VacuumStmt)
	if !ok {
		t.Fatalf("expected *nodes.VacuumStmt, got %T", stmt)
	}
	if vs.IsVacuumCmd != false {
		t.Errorf("expected IsVacuumCmd false, got %v", vs.IsVacuumCmd)
	}
}

func TestAnalyzeTable(t *testing.T) {
	stmt := parseStmt(t, "ANALYZE t")
	vs, ok := stmt.(*nodes.VacuumStmt)
	if !ok {
		t.Fatalf("expected *nodes.VacuumStmt, got %T", stmt)
	}
	if vs.IsVacuumCmd != false {
		t.Errorf("expected IsVacuumCmd false, got %v", vs.IsVacuumCmd)
	}
	if vs.Rels == nil || vs.Rels.Len() != 1 {
		t.Fatalf("expected 1 relation, got %v", vs.Rels)
	}
}

func TestAnalyzeColumns(t *testing.T) {
	stmt := parseStmt(t, "ANALYZE t(col1, col2)")
	vs, ok := stmt.(*nodes.VacuumStmt)
	if !ok {
		t.Fatalf("expected *nodes.VacuumStmt, got %T", stmt)
	}
	if vs.Rels == nil || vs.Rels.Len() != 1 {
		t.Fatalf("expected 1 relation, got %v", vs.Rels)
	}
	vr := vs.Rels.Items[0].(*nodes.VacuumRelation)
	if vr.VaCols == nil || vr.VaCols.Len() != 2 {
		t.Fatalf("expected 2 columns, got %v", vr.VaCols)
	}
}

/*****************************************************************************
 *
 *      CLUSTER tests
 *
 *****************************************************************************/

func TestClusterTable(t *testing.T) {
	stmt := parseStmt(t, "CLUSTER t")
	cs, ok := stmt.(*nodes.ClusterStmt)
	if !ok {
		t.Fatalf("expected *nodes.ClusterStmt, got %T", stmt)
	}
	if cs.Relation == nil || cs.Relation.Relname != "t" {
		t.Errorf("expected relation 't', got %v", cs.Relation)
	}
	if cs.Indexname != "" {
		t.Errorf("expected empty indexname, got %q", cs.Indexname)
	}
}

func TestClusterTableUsing(t *testing.T) {
	stmt := parseStmt(t, "CLUSTER t USING idx")
	cs, ok := stmt.(*nodes.ClusterStmt)
	if !ok {
		t.Fatalf("expected *nodes.ClusterStmt, got %T", stmt)
	}
	if cs.Relation == nil || cs.Relation.Relname != "t" {
		t.Errorf("expected relation 't', got %v", cs.Relation)
	}
	if cs.Indexname != "idx" {
		t.Errorf("expected indexname 'idx', got %q", cs.Indexname)
	}
}

func TestClusterVerbose(t *testing.T) {
	stmt := parseStmt(t, "CLUSTER (VERBOSE) t USING idx")
	cs, ok := stmt.(*nodes.ClusterStmt)
	if !ok {
		t.Fatalf("expected *nodes.ClusterStmt, got %T", stmt)
	}
	if cs.Relation == nil || cs.Relation.Relname != "t" {
		t.Errorf("expected relation 't', got %v", cs.Relation)
	}
	if cs.Indexname != "idx" {
		t.Errorf("expected indexname 'idx', got %q", cs.Indexname)
	}
	if cs.Params == nil || cs.Params.Len() < 1 {
		t.Fatalf("expected params with verbose, got %v", cs.Params)
	}
}

/*****************************************************************************
 *
 *      REINDEX tests
 *
 *****************************************************************************/

func TestReindexTable(t *testing.T) {
	stmt := parseStmt(t, "REINDEX TABLE t")
	rs, ok := stmt.(*nodes.ReindexStmt)
	if !ok {
		t.Fatalf("expected *nodes.ReindexStmt, got %T", stmt)
	}
	if rs.Kind != nodes.REINDEX_OBJECT_TABLE {
		t.Errorf("expected kind REINDEX_OBJECT_TABLE, got %d", rs.Kind)
	}
	if rs.Relation == nil || rs.Relation.Relname != "t" {
		t.Errorf("expected relation 't', got %v", rs.Relation)
	}
}

func TestReindexIndex(t *testing.T) {
	stmt := parseStmt(t, "REINDEX INDEX idx")
	rs, ok := stmt.(*nodes.ReindexStmt)
	if !ok {
		t.Fatalf("expected *nodes.ReindexStmt, got %T", stmt)
	}
	if rs.Kind != nodes.REINDEX_OBJECT_INDEX {
		t.Errorf("expected kind REINDEX_OBJECT_INDEX, got %d", rs.Kind)
	}
}

func TestReindexSchema(t *testing.T) {
	stmt := parseStmt(t, "REINDEX SCHEMA myschema")
	rs, ok := stmt.(*nodes.ReindexStmt)
	if !ok {
		t.Fatalf("expected *nodes.ReindexStmt, got %T", stmt)
	}
	if rs.Kind != nodes.REINDEX_OBJECT_SCHEMA {
		t.Errorf("expected kind REINDEX_OBJECT_SCHEMA, got %d", rs.Kind)
	}
	if rs.Name != "myschema" {
		t.Errorf("expected name 'myschema', got %q", rs.Name)
	}
}

func TestReindexDatabase(t *testing.T) {
	stmt := parseStmt(t, "REINDEX DATABASE mydb")
	rs, ok := stmt.(*nodes.ReindexStmt)
	if !ok {
		t.Fatalf("expected *nodes.ReindexStmt, got %T", stmt)
	}
	if rs.Kind != nodes.REINDEX_OBJECT_DATABASE {
		t.Errorf("expected kind REINDEX_OBJECT_DATABASE, got %d", rs.Kind)
	}
	if rs.Name != "mydb" {
		t.Errorf("expected name 'mydb', got %q", rs.Name)
	}
}

func TestReindexSystem(t *testing.T) {
	stmt := parseStmt(t, "REINDEX SYSTEM mydb")
	rs, ok := stmt.(*nodes.ReindexStmt)
	if !ok {
		t.Fatalf("expected *nodes.ReindexStmt, got %T", stmt)
	}
	if rs.Kind != nodes.REINDEX_OBJECT_SYSTEM {
		t.Errorf("expected kind REINDEX_OBJECT_SYSTEM, got %d", rs.Kind)
	}
	if rs.Name != "mydb" {
		t.Errorf("expected name 'mydb', got %q", rs.Name)
	}
}

func TestReindexVerboseTable(t *testing.T) {
	stmt := parseStmt(t, "REINDEX (VERBOSE) TABLE t")
	rs, ok := stmt.(*nodes.ReindexStmt)
	if !ok {
		t.Fatalf("expected *nodes.ReindexStmt, got %T", stmt)
	}
	if rs.Kind != nodes.REINDEX_OBJECT_TABLE {
		t.Errorf("expected kind REINDEX_OBJECT_TABLE, got %d", rs.Kind)
	}
	if rs.Params == nil || rs.Params.Len() < 1 {
		t.Fatalf("expected params with verbose, got %v", rs.Params)
	}
}
