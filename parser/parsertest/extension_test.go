package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE EXTENSION tests
// =============================================================================

// TestCreateExtensionBasic tests: CREATE EXTENSION pgcrypto
func TestCreateExtensionBasic(t *testing.T) {
	input := "CREATE EXTENSION pgcrypto"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.CreateExtensionStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateExtensionStmt, got %T", result.Items[0])
	}
	if stmt.Extname != "pgcrypto" {
		t.Errorf("expected extname 'pgcrypto', got %q", stmt.Extname)
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
}

// TestCreateExtensionIfNotExists tests: CREATE EXTENSION IF NOT EXISTS pgcrypto
func TestCreateExtensionIfNotExists(t *testing.T) {
	input := "CREATE EXTENSION IF NOT EXISTS pgcrypto"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateExtensionStmt)
	if stmt.Extname != "pgcrypto" {
		t.Errorf("expected extname 'pgcrypto', got %q", stmt.Extname)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
}

// TestCreateExtensionSchema tests: CREATE EXTENSION pgcrypto SCHEMA public
func TestCreateExtensionSchema(t *testing.T) {
	input := "CREATE EXTENSION pgcrypto SCHEMA public"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateExtensionStmt)
	if stmt.Extname != "pgcrypto" {
		t.Errorf("expected extname 'pgcrypto', got %q", stmt.Extname)
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "schema" {
		t.Errorf("expected defname 'schema', got %q", opt.Defname)
	}
	if opt.Arg.(*nodes.String).Str != "public" {
		t.Errorf("expected schema 'public', got %q", opt.Arg.(*nodes.String).Str)
	}
}

// TestCreateExtensionVersion tests: CREATE EXTENSION pgcrypto VERSION '1.3'
func TestCreateExtensionVersion(t *testing.T) {
	input := "CREATE EXTENSION pgcrypto VERSION '1.3'"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateExtensionStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "new_version" {
		t.Errorf("expected defname 'new_version', got %q", opt.Defname)
	}
	if opt.Arg.(*nodes.String).Str != "1.3" {
		t.Errorf("expected version '1.3', got %q", opt.Arg.(*nodes.String).Str)
	}
}

// TestCreateExtensionCascade tests: CREATE EXTENSION pgcrypto CASCADE
func TestCreateExtensionCascade(t *testing.T) {
	input := "CREATE EXTENSION pgcrypto CASCADE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateExtensionStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "cascade" {
		t.Errorf("expected defname 'cascade', got %q", opt.Defname)
	}
}

// TestAlterExtensionUpdate tests: ALTER EXTENSION pgcrypto UPDATE
func TestAlterExtensionUpdate(t *testing.T) {
	input := "ALTER EXTENSION pgcrypto UPDATE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterExtensionStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterExtensionStmt, got %T", result.Items[0])
	}
	if stmt.Extname != "pgcrypto" {
		t.Errorf("expected extname 'pgcrypto', got %q", stmt.Extname)
	}
}

// TestAlterExtensionUpdateTo tests: ALTER EXTENSION pgcrypto UPDATE TO '2.0'
func TestAlterExtensionUpdateTo(t *testing.T) {
	input := "ALTER EXTENSION pgcrypto UPDATE TO '2.0'"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterExtensionStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "new_version" {
		t.Errorf("expected defname 'new_version', got %q", opt.Defname)
	}
	if opt.Arg.(*nodes.String).Str != "2.0" {
		t.Errorf("expected version '2.0', got %q", opt.Arg.(*nodes.String).Str)
	}
}

// TestAlterExtensionAddTable tests: ALTER EXTENSION pgcrypto ADD TABLE t1
func TestAlterExtensionAddTable(t *testing.T) {
	input := "ALTER EXTENSION pgcrypto ADD TABLE t1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterExtensionContentsStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterExtensionContentsStmt, got %T", result.Items[0])
	}
	if stmt.Extname != "pgcrypto" {
		t.Errorf("expected extname 'pgcrypto', got %q", stmt.Extname)
	}
	if stmt.Action != 1 {
		t.Errorf("expected action 1 (ADD), got %d", stmt.Action)
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected objtype OBJECT_TABLE, got %d", stmt.Objtype)
	}
}

// =============================================================================
// CREATE TABLESPACE / DROP TABLESPACE / ALTER TABLESPACE tests
// =============================================================================

// TestCreateTablespaceBasic tests: CREATE TABLESPACE ts LOCATION '/data'
func TestCreateTablespaceBasic(t *testing.T) {
	input := "CREATE TABLESPACE ts LOCATION '/data'"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateTableSpaceStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateTableSpaceStmt, got %T", result.Items[0])
	}
	if stmt.Tablespacename != "ts" {
		t.Errorf("expected tablespacename 'ts', got %q", stmt.Tablespacename)
	}
	if stmt.Location != "/data" {
		t.Errorf("expected location '/data', got %q", stmt.Location)
	}
	if stmt.Owner != nil {
		t.Error("expected owner to be nil")
	}
}

// TestCreateTablespaceOwner tests: CREATE TABLESPACE ts OWNER myuser LOCATION '/data'
func TestCreateTablespaceOwner(t *testing.T) {
	input := "CREATE TABLESPACE ts OWNER myuser LOCATION '/data'"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTableSpaceStmt)
	if stmt.Owner == nil {
		t.Fatal("expected owner to be non-nil")
	}
	if stmt.Owner.Rolename != "myuser" {
		t.Errorf("expected owner 'myuser', got %q", stmt.Owner.Rolename)
	}
}

// TestDropTablespace tests: DROP TABLESPACE ts
func TestDropTablespace(t *testing.T) {
	input := "DROP TABLESPACE ts"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropTableSpaceStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropTableSpaceStmt, got %T", result.Items[0])
	}
	if stmt.Tablespacename != "ts" {
		t.Errorf("expected tablespacename 'ts', got %q", stmt.Tablespacename)
	}
	if stmt.MissingOk {
		t.Error("expected MissingOk to be false")
	}
}

// TestDropTablespaceIfExists tests: DROP TABLESPACE IF EXISTS ts
func TestDropTablespaceIfExists(t *testing.T) {
	input := "DROP TABLESPACE IF EXISTS ts"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropTableSpaceStmt)
	if !stmt.MissingOk {
		t.Error("expected MissingOk to be true")
	}
}

// TestAlterTablespaceSet tests: ALTER TABLESPACE ts SET (seq_page_cost = 2)
func TestAlterTablespaceSet(t *testing.T) {
	input := "ALTER TABLESPACE ts SET (seq_page_cost = 2)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterTableSpaceOptionsStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableSpaceOptionsStmt, got %T", result.Items[0])
	}
	if stmt.Tablespacename != "ts" {
		t.Errorf("expected tablespacename 'ts', got %q", stmt.Tablespacename)
	}
	if stmt.IsReset {
		t.Error("expected IsReset to be false")
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
}

// TestAlterTablespaceReset tests: ALTER TABLESPACE ts RESET (seq_page_cost)
func TestAlterTablespaceReset(t *testing.T) {
	input := "ALTER TABLESPACE ts RESET (seq_page_cost)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterTableSpaceOptionsStmt)
	if !stmt.IsReset {
		t.Error("expected IsReset to be true")
	}
}

// =============================================================================
// CREATE ACCESS METHOD tests
// =============================================================================

// TestCreateAmIndex tests: CREATE ACCESS METHOD myam TYPE INDEX HANDLER handler_func
func TestCreateAmIndex(t *testing.T) {
	input := "CREATE ACCESS METHOD myam TYPE INDEX HANDLER handler_func"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateAmStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateAmStmt, got %T", result.Items[0])
	}
	if stmt.Amname != "myam" {
		t.Errorf("expected amname 'myam', got %q", stmt.Amname)
	}
	if stmt.Amtype != 'i' {
		t.Errorf("expected amtype 'i', got %c", stmt.Amtype)
	}
	if stmt.HandlerName == nil || stmt.HandlerName.Len() != 1 {
		t.Fatal("expected handler name with 1 element")
	}
	handlerStr := stmt.HandlerName.Items[0].(*nodes.String).Str
	if handlerStr != "handler_func" {
		t.Errorf("expected handler 'handler_func', got %q", handlerStr)
	}
}

// TestCreateAmTable tests: CREATE ACCESS METHOD myam TYPE TABLE HANDLER handler_func
func TestCreateAmTable(t *testing.T) {
	input := "CREATE ACCESS METHOD myam TYPE TABLE HANDLER handler_func"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateAmStmt)
	if stmt.Amtype != 't' {
		t.Errorf("expected amtype 't', got %c", stmt.Amtype)
	}
}
