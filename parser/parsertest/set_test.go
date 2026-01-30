package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
)

/*****************************************************************************
 *
 *      SET variable tests
 *
 *****************************************************************************/

func TestSetSearchPathTO(t *testing.T) {
	stmt := parseStmt(t, "SET search_path TO public")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "search_path" {
		t.Errorf("expected Name 'search_path', got %q", vs.Name)
	}
	if vs.IsLocal {
		t.Errorf("expected IsLocal false")
	}
	if vs.Args == nil || len(vs.Args.Items) != 1 {
		t.Fatalf("expected 1 arg, got %v", vs.Args)
	}
}

func TestSetSearchPathEquals(t *testing.T) {
	stmt := parseStmt(t, "SET search_path = 'public'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "search_path" {
		t.Errorf("expected Name 'search_path', got %q", vs.Name)
	}
	if vs.Args == nil || len(vs.Args.Items) != 1 {
		t.Fatalf("expected 1 arg, got %v", vs.Args)
	}
}

func TestSetLocalTimezone(t *testing.T) {
	stmt := parseStmt(t, "SET LOCAL timezone = 'UTC'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
	if !vs.IsLocal {
		t.Errorf("expected IsLocal true")
	}
}

func TestSetSessionTimezone(t *testing.T) {
	stmt := parseStmt(t, "SET SESSION timezone = 'UTC'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
	if vs.IsLocal {
		t.Errorf("expected IsLocal false for SESSION")
	}
}

func TestSetSearchPathMultipleValues(t *testing.T) {
	stmt := parseStmt(t, "SET search_path TO public, pg_catalog")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Args == nil || len(vs.Args.Items) != 2 {
		t.Fatalf("expected 2 args, got %v", vs.Args)
	}
}

func TestSetSearchPathDefault(t *testing.T) {
	stmt := parseStmt(t, "SET search_path TO DEFAULT")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_DEFAULT {
		t.Errorf("expected Kind VAR_SET_DEFAULT, got %d", vs.Kind)
	}
	if vs.Name != "search_path" {
		t.Errorf("expected Name 'search_path', got %q", vs.Name)
	}
}

func TestSetTimeZoneString(t *testing.T) {
	stmt := parseStmt(t, "SET TIME ZONE 'UTC'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
	if vs.Args == nil || len(vs.Args.Items) != 1 {
		t.Fatalf("expected 1 arg, got %v", vs.Args)
	}
}

func TestSetTimeZoneLOCAL(t *testing.T) {
	stmt := parseStmt(t, "SET TIME ZONE LOCAL")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_DEFAULT {
		t.Errorf("expected Kind VAR_SET_DEFAULT, got %d", vs.Kind)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
}

func TestSetTimeZoneDefault(t *testing.T) {
	stmt := parseStmt(t, "SET TIME ZONE DEFAULT")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_DEFAULT {
		t.Errorf("expected Kind VAR_SET_DEFAULT, got %d", vs.Kind)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
}

func TestSetTimeZoneNumeric(t *testing.T) {
	stmt := parseStmt(t, "SET TIME ZONE 0")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
}

func TestSetNames(t *testing.T) {
	stmt := parseStmt(t, "SET NAMES 'UTF8'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "client_encoding" {
		t.Errorf("expected Name 'client_encoding', got %q", vs.Name)
	}
}

func TestSetSessionAuthorization(t *testing.T) {
	stmt := parseStmt(t, "SET SESSION AUTHORIZATION 'user1'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "session_authorization" {
		t.Errorf("expected Name 'session_authorization', got %q", vs.Name)
	}
}

func TestSetSessionAuthorizationDefault(t *testing.T) {
	stmt := parseStmt(t, "SET SESSION AUTHORIZATION DEFAULT")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_DEFAULT {
		t.Errorf("expected Kind VAR_SET_DEFAULT, got %d", vs.Kind)
	}
	if vs.Name != "session_authorization" {
		t.Errorf("expected Name 'session_authorization', got %q", vs.Name)
	}
}

func TestSetRole(t *testing.T) {
	stmt := parseStmt(t, "SET ROLE 'admin'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "role" {
		t.Errorf("expected Name 'role', got %q", vs.Name)
	}
}

func TestSetXmlOptionDocument(t *testing.T) {
	stmt := parseStmt(t, "SET XML OPTION DOCUMENT")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "xmloption" {
		t.Errorf("expected Name 'xmloption', got %q", vs.Name)
	}
}

func TestSetXmlOptionContent(t *testing.T) {
	stmt := parseStmt(t, "SET XML OPTION CONTENT")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "xmloption" {
		t.Errorf("expected Name 'xmloption', got %q", vs.Name)
	}
	// Check arg value
	if vs.Args != nil && len(vs.Args.Items) == 1 {
		ac, ok := vs.Args.Items[0].(*nodes.A_Const)
		if ok {
			s, ok := ac.Val.(*nodes.String)
			if ok && s.Str != "CONTENT" {
				t.Errorf("expected xmloption arg 'CONTENT', got %q", s.Str)
			}
		}
	}
}

func TestSetTransactionIsolationLevel(t *testing.T) {
	stmt := parseStmt(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_MULTI {
		t.Errorf("expected Kind VAR_SET_MULTI, got %d", vs.Kind)
	}
	if vs.Name != "TRANSACTION" {
		t.Errorf("expected Name 'TRANSACTION', got %q", vs.Name)
	}
}

func TestSetTransactionReadOnly(t *testing.T) {
	stmt := parseStmt(t, "SET TRANSACTION READ ONLY")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_MULTI {
		t.Errorf("expected Kind VAR_SET_MULTI, got %d", vs.Kind)
	}
	if vs.Name != "TRANSACTION" {
		t.Errorf("expected Name 'TRANSACTION', got %q", vs.Name)
	}
}

/*****************************************************************************
 *
 *      SHOW variable tests
 *
 *****************************************************************************/

func TestShowSearchPath(t *testing.T) {
	stmt := parseStmt(t, "SHOW search_path")
	vs, ok := stmt.(*nodes.VariableShowStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableShowStmt, got %T", stmt)
	}
	if vs.Name != "search_path" {
		t.Errorf("expected Name 'search_path', got %q", vs.Name)
	}
}

func TestShowTimeZone(t *testing.T) {
	stmt := parseStmt(t, "SHOW TIME ZONE")
	vs, ok := stmt.(*nodes.VariableShowStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableShowStmt, got %T", stmt)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
}

func TestShowTransactionIsolationLevel(t *testing.T) {
	stmt := parseStmt(t, "SHOW TRANSACTION ISOLATION LEVEL")
	vs, ok := stmt.(*nodes.VariableShowStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableShowStmt, got %T", stmt)
	}
	if vs.Name != "transaction_isolation" {
		t.Errorf("expected Name 'transaction_isolation', got %q", vs.Name)
	}
}

func TestShowSessionAuthorization(t *testing.T) {
	stmt := parseStmt(t, "SHOW SESSION AUTHORIZATION")
	vs, ok := stmt.(*nodes.VariableShowStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableShowStmt, got %T", stmt)
	}
	if vs.Name != "session_authorization" {
		t.Errorf("expected Name 'session_authorization', got %q", vs.Name)
	}
}

func TestShowAll(t *testing.T) {
	stmt := parseStmt(t, "SHOW ALL")
	vs, ok := stmt.(*nodes.VariableShowStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableShowStmt, got %T", stmt)
	}
	if vs.Name != "all" {
		t.Errorf("expected Name 'all', got %q", vs.Name)
	}
}

/*****************************************************************************
 *
 *      RESET variable tests
 *
 *****************************************************************************/

func TestResetSearchPath(t *testing.T) {
	stmt := parseStmt(t, "RESET search_path")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_RESET {
		t.Errorf("expected Kind VAR_RESET, got %d", vs.Kind)
	}
	if vs.Name != "search_path" {
		t.Errorf("expected Name 'search_path', got %q", vs.Name)
	}
}

func TestResetTimeZone(t *testing.T) {
	stmt := parseStmt(t, "RESET TIME ZONE")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_RESET {
		t.Errorf("expected Kind VAR_RESET, got %d", vs.Kind)
	}
	if vs.Name != "timezone" {
		t.Errorf("expected Name 'timezone', got %q", vs.Name)
	}
}

func TestResetTransactionIsolationLevel(t *testing.T) {
	stmt := parseStmt(t, "RESET TRANSACTION ISOLATION LEVEL")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_RESET {
		t.Errorf("expected Kind VAR_RESET, got %d", vs.Kind)
	}
	if vs.Name != "transaction_isolation" {
		t.Errorf("expected Name 'transaction_isolation', got %q", vs.Name)
	}
}

func TestResetSessionAuthorization(t *testing.T) {
	stmt := parseStmt(t, "RESET SESSION AUTHORIZATION")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_RESET {
		t.Errorf("expected Kind VAR_RESET, got %d", vs.Kind)
	}
	if vs.Name != "session_authorization" {
		t.Errorf("expected Name 'session_authorization', got %q", vs.Name)
	}
}

func TestResetAll(t *testing.T) {
	stmt := parseStmt(t, "RESET ALL")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_RESET_ALL {
		t.Errorf("expected Kind VAR_RESET_ALL, got %d", vs.Kind)
	}
}

/*****************************************************************************
 *
 *      PREPARE / EXECUTE / DEALLOCATE tests
 *
 *****************************************************************************/

func TestPrepareSimple(t *testing.T) {
	stmt := parseStmt(t, "PREPARE p AS SELECT 1")
	ps, ok := stmt.(*nodes.PrepareStmt)
	if !ok {
		t.Fatalf("expected *nodes.PrepareStmt, got %T", stmt)
	}
	if ps.Name != "p" {
		t.Errorf("expected Name 'p', got %q", ps.Name)
	}
	if ps.Argtypes != nil {
		t.Errorf("expected nil Argtypes, got %v", ps.Argtypes)
	}
	if ps.Query == nil {
		t.Fatal("expected non-nil Query")
	}
}

func TestPrepareWithTypes(t *testing.T) {
	stmt := parseStmt(t, "PREPARE p(integer) AS SELECT 1")
	ps, ok := stmt.(*nodes.PrepareStmt)
	if !ok {
		t.Fatalf("expected *nodes.PrepareStmt, got %T", stmt)
	}
	if ps.Name != "p" {
		t.Errorf("expected Name 'p', got %q", ps.Name)
	}
	if ps.Argtypes == nil || len(ps.Argtypes.Items) != 1 {
		t.Fatalf("expected 1 argtype, got %v", ps.Argtypes)
	}
	if ps.Query == nil {
		t.Fatal("expected non-nil Query")
	}
}

func TestExecuteSimple(t *testing.T) {
	stmt := parseStmt(t, "EXECUTE p")
	es, ok := stmt.(*nodes.ExecuteStmt)
	if !ok {
		t.Fatalf("expected *nodes.ExecuteStmt, got %T", stmt)
	}
	if es.Name != "p" {
		t.Errorf("expected Name 'p', got %q", es.Name)
	}
	if es.Params != nil {
		t.Errorf("expected nil Params, got %v", es.Params)
	}
}

func TestExecuteWithParams(t *testing.T) {
	stmt := parseStmt(t, "EXECUTE p(1, 'hello')")
	es, ok := stmt.(*nodes.ExecuteStmt)
	if !ok {
		t.Fatalf("expected *nodes.ExecuteStmt, got %T", stmt)
	}
	if es.Name != "p" {
		t.Errorf("expected Name 'p', got %q", es.Name)
	}
	if es.Params == nil || len(es.Params.Items) != 2 {
		t.Fatalf("expected 2 params, got %v", es.Params)
	}
}

func TestDeallocateName(t *testing.T) {
	stmt := parseStmt(t, "DEALLOCATE p")
	ds, ok := stmt.(*nodes.DeallocateStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeallocateStmt, got %T", stmt)
	}
	if ds.Name != "p" {
		t.Errorf("expected Name 'p', got %q", ds.Name)
	}
	if ds.IsAll {
		t.Errorf("expected IsAll false")
	}
}

func TestDeallocatePrepareName(t *testing.T) {
	stmt := parseStmt(t, "DEALLOCATE PREPARE p")
	ds, ok := stmt.(*nodes.DeallocateStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeallocateStmt, got %T", stmt)
	}
	if ds.Name != "p" {
		t.Errorf("expected Name 'p', got %q", ds.Name)
	}
	if ds.IsAll {
		t.Errorf("expected IsAll false")
	}
}

func TestDeallocateAll(t *testing.T) {
	stmt := parseStmt(t, "DEALLOCATE ALL")
	ds, ok := stmt.(*nodes.DeallocateStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeallocateStmt, got %T", stmt)
	}
	if !ds.IsAll {
		t.Errorf("expected IsAll true")
	}
}

func TestDeallocatePrepareAll(t *testing.T) {
	stmt := parseStmt(t, "DEALLOCATE PREPARE ALL")
	ds, ok := stmt.(*nodes.DeallocateStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeallocateStmt, got %T", stmt)
	}
	if !ds.IsAll {
		t.Errorf("expected IsAll true")
	}
}

/*****************************************************************************
 *
 *      SET CONSTRAINTS still works (regression test)
 *
 *****************************************************************************/

func TestSetConstraintsStillWorks(t *testing.T) {
	stmt := parseStmt(t, "SET CONSTRAINTS ALL DEFERRED")
	cs, ok := stmt.(*nodes.ConstraintsSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.ConstraintsSetStmt, got %T", stmt)
	}
	if !cs.Deferred {
		t.Errorf("expected Deferred true")
	}
}

func TestSetConstraintsImmediate(t *testing.T) {
	stmt := parseStmt(t, "SET CONSTRAINTS my_fk IMMEDIATE")
	cs, ok := stmt.(*nodes.ConstraintsSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.ConstraintsSetStmt, got %T", stmt)
	}
	if cs.Deferred {
		t.Errorf("expected Deferred false")
	}
}

/*****************************************************************************
 *
 *      Dotted var_name test
 *
 *****************************************************************************/

func TestSetDottedVarName(t *testing.T) {
	stmt := parseStmt(t, "SET myapp.debug TO true")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Name != "myapp.debug" {
		t.Errorf("expected Name 'myapp.debug', got %q", vs.Name)
	}
}

/*****************************************************************************
 *
 *      SET SCHEMA test
 *
 *****************************************************************************/

func TestSetSchema(t *testing.T) {
	stmt := parseStmt(t, "SET SCHEMA 'myschema'")
	vs, ok := stmt.(*nodes.VariableSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.VariableSetStmt, got %T", stmt)
	}
	if vs.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected Kind VAR_SET_VALUE, got %d", vs.Kind)
	}
	if vs.Name != "search_path" {
		t.Errorf("expected Name 'search_path', got %q", vs.Name)
	}
}
