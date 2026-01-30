package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
)

/*****************************************************************************
 *
 *      MERGE tests
 *
 *****************************************************************************/

func TestMergeUpdateBasic(t *testing.T) {
	stmt := parseStmt(t, "MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN UPDATE SET col = s.col")
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	if ms.Relation == nil {
		t.Fatal("expected Relation to be non-nil")
	}
	if ms.Relation.Relname != "target" {
		t.Errorf("expected Relation.Relname 'target', got %q", ms.Relation.Relname)
	}
	if ms.Relation.Alias == nil || ms.Relation.Alias.Aliasname != "t" {
		t.Errorf("expected alias 't'")
	}
	if ms.SourceRelation == nil {
		t.Fatal("expected SourceRelation to be non-nil")
	}
	if ms.JoinCondition == nil {
		t.Fatal("expected JoinCondition to be non-nil")
	}
	if ms.MergeWhenClauses == nil || len(ms.MergeWhenClauses.Items) != 1 {
		t.Fatalf("expected 1 merge when clause, got %v", ms.MergeWhenClauses)
	}
	wc := ms.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	if wc.Kind != nodes.MERGE_WHEN_MATCHED {
		t.Errorf("expected MERGE_WHEN_MATCHED, got %d", wc.Kind)
	}
	if wc.CommandType != nodes.CMD_UPDATE {
		t.Errorf("expected CMD_UPDATE, got %d", wc.CommandType)
	}
	if wc.TargetList == nil {
		t.Error("expected TargetList to be non-nil for UPDATE")
	}
}

func TestMergeInsertBasic(t *testing.T) {
	stmt := parseStmt(t, "MERGE INTO target USING source ON target.id = source.id WHEN NOT MATCHED THEN INSERT (col) VALUES (source.col)")
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	if ms.MergeWhenClauses == nil || len(ms.MergeWhenClauses.Items) != 1 {
		t.Fatalf("expected 1 merge when clause")
	}
	wc := ms.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	if wc.Kind != nodes.MERGE_WHEN_NOT_MATCHED_BY_TARGET {
		t.Errorf("expected MERGE_WHEN_NOT_MATCHED_BY_TARGET, got %d", wc.Kind)
	}
	if wc.CommandType != nodes.CMD_INSERT {
		t.Errorf("expected CMD_INSERT, got %d", wc.CommandType)
	}
	if wc.TargetList == nil {
		t.Error("expected TargetList to be non-nil for INSERT with columns")
	}
	if wc.Values == nil {
		t.Error("expected Values to be non-nil for INSERT")
	}
}

func TestMergeDeleteBasic(t *testing.T) {
	stmt := parseStmt(t, "MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN DELETE")
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	if ms.MergeWhenClauses == nil || len(ms.MergeWhenClauses.Items) != 1 {
		t.Fatalf("expected 1 merge when clause")
	}
	wc := ms.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	if wc.CommandType != nodes.CMD_DELETE {
		t.Errorf("expected CMD_DELETE, got %d", wc.CommandType)
	}
}

func TestMergeDoNothingMatched(t *testing.T) {
	stmt := parseStmt(t, "MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN DO NOTHING")
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	wc := ms.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	if wc.CommandType != nodes.CMD_NOTHING {
		t.Errorf("expected CMD_NOTHING, got %d", wc.CommandType)
	}
	if wc.Kind != nodes.MERGE_WHEN_MATCHED {
		t.Errorf("expected MERGE_WHEN_MATCHED, got %d", wc.Kind)
	}
}

func TestMergeDoNothingNotMatched(t *testing.T) {
	stmt := parseStmt(t, "MERGE INTO target t USING source s ON t.id = s.id WHEN NOT MATCHED THEN DO NOTHING")
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	wc := ms.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	if wc.CommandType != nodes.CMD_NOTHING {
		t.Errorf("expected CMD_NOTHING, got %d", wc.CommandType)
	}
	if wc.Kind != nodes.MERGE_WHEN_NOT_MATCHED_BY_TARGET {
		t.Errorf("expected MERGE_WHEN_NOT_MATCHED_BY_TARGET, got %d", wc.Kind)
	}
}

func TestMergeMultipleClauses(t *testing.T) {
	sql := `MERGE INTO target t USING source s ON t.id = s.id
		WHEN MATCHED THEN UPDATE SET col = s.col
		WHEN NOT MATCHED THEN INSERT (col) VALUES (s.col)
		WHEN MATCHED THEN DELETE`
	stmt := parseStmt(t, sql)
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	if ms.MergeWhenClauses == nil || len(ms.MergeWhenClauses.Items) != 3 {
		t.Fatalf("expected 3 merge when clauses, got %d", len(ms.MergeWhenClauses.Items))
	}
	// First: WHEN MATCHED THEN UPDATE
	wc0 := ms.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	if wc0.CommandType != nodes.CMD_UPDATE {
		t.Errorf("clause 0: expected CMD_UPDATE, got %d", wc0.CommandType)
	}
	// Second: WHEN NOT MATCHED THEN INSERT
	wc1 := ms.MergeWhenClauses.Items[1].(*nodes.MergeWhenClause)
	if wc1.CommandType != nodes.CMD_INSERT {
		t.Errorf("clause 1: expected CMD_INSERT, got %d", wc1.CommandType)
	}
	// Third: WHEN MATCHED THEN DELETE
	wc2 := ms.MergeWhenClauses.Items[2].(*nodes.MergeWhenClause)
	if wc2.CommandType != nodes.CMD_DELETE {
		t.Errorf("clause 2: expected CMD_DELETE, got %d", wc2.CommandType)
	}
}

func TestMergeSubquery(t *testing.T) {
	sql := "MERGE INTO target t USING (SELECT * FROM source) s ON t.id = s.id WHEN MATCHED THEN UPDATE SET col = s.col"
	stmt := parseStmt(t, sql)
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	if ms.SourceRelation == nil {
		t.Fatal("expected SourceRelation to be non-nil")
	}
	// SourceRelation should be a RangeSubselect
	_, ok = ms.SourceRelation.(*nodes.RangeSubselect)
	if !ok {
		t.Errorf("expected source to be *nodes.RangeSubselect, got %T", ms.SourceRelation)
	}
}

func TestMergeWithCondition(t *testing.T) {
	sql := "MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED AND t.col > 0 THEN UPDATE SET col = s.col"
	stmt := parseStmt(t, sql)
	ms, ok := stmt.(*nodes.MergeStmt)
	if !ok {
		t.Fatalf("expected *nodes.MergeStmt, got %T", stmt)
	}
	wc := ms.MergeWhenClauses.Items[0].(*nodes.MergeWhenClause)
	if wc.Condition == nil {
		t.Error("expected Condition to be non-nil")
	}
}

/*****************************************************************************
 *
 *      CALL tests
 *
 *****************************************************************************/

func TestCallNoArgs(t *testing.T) {
	stmt := parseStmt(t, "CALL myfunc()")
	cs, ok := stmt.(*nodes.CallStmt)
	if !ok {
		t.Fatalf("expected *nodes.CallStmt, got %T", stmt)
	}
	if cs.Funccall == nil {
		t.Fatal("expected Funccall to be non-nil")
	}
	if cs.Funccall.Funcname == nil || len(cs.Funccall.Funcname.Items) != 1 {
		t.Fatalf("expected 1 name component")
	}
	fname := cs.Funccall.Funcname.Items[0].(*nodes.String).Str
	if fname != "myfunc" {
		t.Errorf("expected function name 'myfunc', got %q", fname)
	}
}

func TestCallWithArgs(t *testing.T) {
	stmt := parseStmt(t, "CALL myfunc(1, 'hello')")
	cs, ok := stmt.(*nodes.CallStmt)
	if !ok {
		t.Fatalf("expected *nodes.CallStmt, got %T", stmt)
	}
	if cs.Funccall == nil {
		t.Fatal("expected Funccall to be non-nil")
	}
	if cs.Funccall.Args == nil || len(cs.Funccall.Args.Items) != 2 {
		t.Fatalf("expected 2 arguments, got %v", cs.Funccall.Args)
	}
}

/*****************************************************************************
 *
 *      DO tests
 *
 *****************************************************************************/

func TestDoBasic(t *testing.T) {
	stmt := parseStmt(t, "DO $$ BEGIN NULL; END $$")
	ds, ok := stmt.(*nodes.DoStmt)
	if !ok {
		t.Fatalf("expected *nodes.DoStmt, got %T", stmt)
	}
	if ds.Args == nil || len(ds.Args.Items) != 1 {
		t.Fatalf("expected 1 arg (code block), got %v", ds.Args)
	}
	de := ds.Args.Items[0].(*nodes.DefElem)
	if de.Defname != "as" {
		t.Errorf("expected Defname 'as', got %q", de.Defname)
	}
}

func TestDoWithLanguage(t *testing.T) {
	stmt := parseStmt(t, "DO LANGUAGE plpgsql $$ BEGIN NULL; END $$")
	ds, ok := stmt.(*nodes.DoStmt)
	if !ok {
		t.Fatalf("expected *nodes.DoStmt, got %T", stmt)
	}
	if ds.Args == nil || len(ds.Args.Items) != 2 {
		t.Fatalf("expected 2 args (language + code), got %v", ds.Args)
	}
	// First should be language
	de0 := ds.Args.Items[0].(*nodes.DefElem)
	if de0.Defname != "language" {
		t.Errorf("expected first arg Defname 'language', got %q", de0.Defname)
	}
	langStr := de0.Arg.(*nodes.String).Str
	if langStr != "plpgsql" {
		t.Errorf("expected language 'plpgsql', got %q", langStr)
	}
	// Second should be code block
	de1 := ds.Args.Items[1].(*nodes.DefElem)
	if de1.Defname != "as" {
		t.Errorf("expected second arg Defname 'as', got %q", de1.Defname)
	}
}

func TestDoLanguageAfterCode(t *testing.T) {
	stmt := parseStmt(t, "DO $$ BEGIN NULL; END $$ LANGUAGE plpgsql")
	ds, ok := stmt.(*nodes.DoStmt)
	if !ok {
		t.Fatalf("expected *nodes.DoStmt, got %T", stmt)
	}
	if ds.Args == nil || len(ds.Args.Items) != 2 {
		t.Fatalf("expected 2 args, got %v", ds.Args)
	}
	// First should be code block
	de0 := ds.Args.Items[0].(*nodes.DefElem)
	if de0.Defname != "as" {
		t.Errorf("expected first arg Defname 'as', got %q", de0.Defname)
	}
	// Second should be language
	de1 := ds.Args.Items[1].(*nodes.DefElem)
	if de1.Defname != "language" {
		t.Errorf("expected second arg Defname 'language', got %q", de1.Defname)
	}
}
