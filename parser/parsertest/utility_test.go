package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single statement of a given type
func parseStmt(t *testing.T, input string) nodes.Node {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	return result.Items[0]
}

/*****************************************************************************
 *
 *      CHECKPOINT tests
 *
 *****************************************************************************/

// TestCheckpoint tests: CHECKPOINT
func TestCheckpoint(t *testing.T) {
	stmt := parseStmt(t, "CHECKPOINT")
	_, ok := stmt.(*nodes.CheckPointStmt)
	if !ok {
		t.Fatalf("expected *nodes.CheckPointStmt, got %T", stmt)
	}
}

/*****************************************************************************
 *
 *      DISCARD tests
 *
 *****************************************************************************/

// TestDiscardAll tests: DISCARD ALL
func TestDiscardAll(t *testing.T) {
	stmt := parseStmt(t, "DISCARD ALL")
	ds, ok := stmt.(*nodes.DiscardStmt)
	if !ok {
		t.Fatalf("expected *nodes.DiscardStmt, got %T", stmt)
	}
	if ds.Target != nodes.DISCARD_ALL {
		t.Errorf("expected Target DISCARD_ALL (%d), got %d", nodes.DISCARD_ALL, ds.Target)
	}
}

// TestDiscardPlans tests: DISCARD PLANS
func TestDiscardPlans(t *testing.T) {
	stmt := parseStmt(t, "DISCARD PLANS")
	ds, ok := stmt.(*nodes.DiscardStmt)
	if !ok {
		t.Fatalf("expected *nodes.DiscardStmt, got %T", stmt)
	}
	if ds.Target != nodes.DISCARD_PLANS {
		t.Errorf("expected Target DISCARD_PLANS (%d), got %d", nodes.DISCARD_PLANS, ds.Target)
	}
}

// TestDiscardSequences tests: DISCARD SEQUENCES
func TestDiscardSequences(t *testing.T) {
	stmt := parseStmt(t, "DISCARD SEQUENCES")
	ds, ok := stmt.(*nodes.DiscardStmt)
	if !ok {
		t.Fatalf("expected *nodes.DiscardStmt, got %T", stmt)
	}
	if ds.Target != nodes.DISCARD_SEQUENCES {
		t.Errorf("expected Target DISCARD_SEQUENCES (%d), got %d", nodes.DISCARD_SEQUENCES, ds.Target)
	}
}

// TestDiscardTemp tests: DISCARD TEMP
func TestDiscardTemp(t *testing.T) {
	stmt := parseStmt(t, "DISCARD TEMP")
	ds, ok := stmt.(*nodes.DiscardStmt)
	if !ok {
		t.Fatalf("expected *nodes.DiscardStmt, got %T", stmt)
	}
	if ds.Target != nodes.DISCARD_TEMP {
		t.Errorf("expected Target DISCARD_TEMP (%d), got %d", nodes.DISCARD_TEMP, ds.Target)
	}
}

// TestDiscardTemporary tests: DISCARD TEMPORARY
func TestDiscardTemporary(t *testing.T) {
	stmt := parseStmt(t, "DISCARD TEMPORARY")
	ds, ok := stmt.(*nodes.DiscardStmt)
	if !ok {
		t.Fatalf("expected *nodes.DiscardStmt, got %T", stmt)
	}
	if ds.Target != nodes.DISCARD_TEMP {
		t.Errorf("expected Target DISCARD_TEMP (%d), got %d", nodes.DISCARD_TEMP, ds.Target)
	}
}

/*****************************************************************************
 *
 *      LISTEN tests
 *
 *****************************************************************************/

// TestListen tests: LISTEN mychannel
func TestListen(t *testing.T) {
	stmt := parseStmt(t, "LISTEN mychannel")
	ls, ok := stmt.(*nodes.ListenStmt)
	if !ok {
		t.Fatalf("expected *nodes.ListenStmt, got %T", stmt)
	}
	if ls.Conditionname != "mychannel" {
		t.Errorf("expected Conditionname 'mychannel', got %q", ls.Conditionname)
	}
}

/*****************************************************************************
 *
 *      UNLISTEN tests
 *
 *****************************************************************************/

// TestUnlisten tests: UNLISTEN mychannel
func TestUnlisten(t *testing.T) {
	stmt := parseStmt(t, "UNLISTEN mychannel")
	us, ok := stmt.(*nodes.UnlistenStmt)
	if !ok {
		t.Fatalf("expected *nodes.UnlistenStmt, got %T", stmt)
	}
	if us.Conditionname != "mychannel" {
		t.Errorf("expected Conditionname 'mychannel', got %q", us.Conditionname)
	}
}

// TestUnlistenAll tests: UNLISTEN *
func TestUnlistenAll(t *testing.T) {
	stmt := parseStmt(t, "UNLISTEN *")
	us, ok := stmt.(*nodes.UnlistenStmt)
	if !ok {
		t.Fatalf("expected *nodes.UnlistenStmt, got %T", stmt)
	}
	if us.Conditionname != "" {
		t.Errorf("expected Conditionname '' (empty for UNLISTEN *), got %q", us.Conditionname)
	}
}

/*****************************************************************************
 *
 *      NOTIFY tests
 *
 *****************************************************************************/

// TestNotify tests: NOTIFY mychannel
func TestNotify(t *testing.T) {
	stmt := parseStmt(t, "NOTIFY mychannel")
	ns, ok := stmt.(*nodes.NotifyStmt)
	if !ok {
		t.Fatalf("expected *nodes.NotifyStmt, got %T", stmt)
	}
	if ns.Conditionname != "mychannel" {
		t.Errorf("expected Conditionname 'mychannel', got %q", ns.Conditionname)
	}
	if ns.Payload != "" {
		t.Errorf("expected Payload '' (empty), got %q", ns.Payload)
	}
}

// TestNotifyWithPayload tests: NOTIFY mychannel, 'hello world'
func TestNotifyWithPayload(t *testing.T) {
	stmt := parseStmt(t, "NOTIFY mychannel, 'hello world'")
	ns, ok := stmt.(*nodes.NotifyStmt)
	if !ok {
		t.Fatalf("expected *nodes.NotifyStmt, got %T", stmt)
	}
	if ns.Conditionname != "mychannel" {
		t.Errorf("expected Conditionname 'mychannel', got %q", ns.Conditionname)
	}
	if ns.Payload != "hello world" {
		t.Errorf("expected Payload 'hello world', got %q", ns.Payload)
	}
}

/*****************************************************************************
 *
 *      LOAD tests
 *
 *****************************************************************************/

// TestLoad tests: LOAD '/usr/lib/libpq.so'
func TestLoad(t *testing.T) {
	stmt := parseStmt(t, "LOAD '/usr/lib/libpq.so'")
	ls, ok := stmt.(*nodes.LoadStmt)
	if !ok {
		t.Fatalf("expected *nodes.LoadStmt, got %T", stmt)
	}
	if ls.Filename != "/usr/lib/libpq.so" {
		t.Errorf("expected Filename '/usr/lib/libpq.so', got %q", ls.Filename)
	}
}

/*****************************************************************************
 *
 *      CLOSE tests
 *
 *****************************************************************************/

// TestCloseCursor tests: CLOSE mycursor
func TestCloseCursor(t *testing.T) {
	stmt := parseStmt(t, "CLOSE mycursor")
	cs, ok := stmt.(*nodes.ClosePortalStmt)
	if !ok {
		t.Fatalf("expected *nodes.ClosePortalStmt, got %T", stmt)
	}
	if cs.Portalname != "mycursor" {
		t.Errorf("expected Portalname 'mycursor', got %q", cs.Portalname)
	}
}

// TestCloseAll tests: CLOSE ALL
func TestCloseAll(t *testing.T) {
	stmt := parseStmt(t, "CLOSE ALL")
	cs, ok := stmt.(*nodes.ClosePortalStmt)
	if !ok {
		t.Fatalf("expected *nodes.ClosePortalStmt, got %T", stmt)
	}
	if cs.Portalname != "" {
		t.Errorf("expected Portalname '' (empty for CLOSE ALL), got %q", cs.Portalname)
	}
}

/*****************************************************************************
 *
 *      SET CONSTRAINTS tests
 *
 *****************************************************************************/

// TestSetConstraintsAllDeferred tests: SET CONSTRAINTS ALL DEFERRED
func TestSetConstraintsAllDeferred(t *testing.T) {
	stmt := parseStmt(t, "SET CONSTRAINTS ALL DEFERRED")
	cs, ok := stmt.(*nodes.ConstraintsSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.ConstraintsSetStmt, got %T", stmt)
	}
	if cs.Constraints != nil {
		t.Errorf("expected Constraints nil (ALL), got %v", cs.Constraints)
	}
	if !cs.Deferred {
		t.Error("expected Deferred to be true")
	}
}

// TestSetConstraintsAllImmediate tests: SET CONSTRAINTS ALL IMMEDIATE
func TestSetConstraintsAllImmediate(t *testing.T) {
	stmt := parseStmt(t, "SET CONSTRAINTS ALL IMMEDIATE")
	cs, ok := stmt.(*nodes.ConstraintsSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.ConstraintsSetStmt, got %T", stmt)
	}
	if cs.Constraints != nil {
		t.Errorf("expected Constraints nil (ALL), got %v", cs.Constraints)
	}
	if cs.Deferred {
		t.Error("expected Deferred to be false")
	}
}

// TestSetConstraintsNameDeferred tests: SET CONSTRAINTS myconstraint DEFERRED
func TestSetConstraintsNameDeferred(t *testing.T) {
	stmt := parseStmt(t, "SET CONSTRAINTS myconstraint DEFERRED")
	cs, ok := stmt.(*nodes.ConstraintsSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.ConstraintsSetStmt, got %T", stmt)
	}
	if cs.Constraints == nil || len(cs.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", cs.Constraints)
	}
	rv, ok := cs.Constraints.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar, got %T", cs.Constraints.Items[0])
	}
	if rv.Relname != "myconstraint" {
		t.Errorf("expected constraint name 'myconstraint', got %q", rv.Relname)
	}
	if !cs.Deferred {
		t.Error("expected Deferred to be true")
	}
}

// TestSetConstraintsMultipleImmediate tests: SET CONSTRAINTS c1, c2 IMMEDIATE
func TestSetConstraintsMultipleImmediate(t *testing.T) {
	stmt := parseStmt(t, "SET CONSTRAINTS c1, c2 IMMEDIATE")
	cs, ok := stmt.(*nodes.ConstraintsSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.ConstraintsSetStmt, got %T", stmt)
	}
	if cs.Constraints == nil || len(cs.Constraints.Items) != 2 {
		t.Fatalf("expected 2 constraints, got %v", cs.Constraints)
	}
	rv0, ok := cs.Constraints.Items[0].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar at index 0, got %T", cs.Constraints.Items[0])
	}
	if rv0.Relname != "c1" {
		t.Errorf("expected constraint name 'c1', got %q", rv0.Relname)
	}
	rv1, ok := cs.Constraints.Items[1].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar at index 1, got %T", cs.Constraints.Items[1])
	}
	if rv1.Relname != "c2" {
		t.Errorf("expected constraint name 'c2', got %q", rv1.Relname)
	}
	if cs.Deferred {
		t.Error("expected Deferred to be false")
	}
}
