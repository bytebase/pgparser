package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single CreatedbStmt
func parseCreatedbStmt(t *testing.T, input string) *nodes.CreatedbStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.CreatedbStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreatedbStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single AlterDatabaseStmt
func parseAlterDatabaseStmt(t *testing.T, input string) *nodes.AlterDatabaseStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.AlterDatabaseStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDatabaseStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single AlterDatabaseSetStmt
func parseAlterDatabaseSetStmt(t *testing.T, input string) *nodes.AlterDatabaseSetStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.AlterDatabaseSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDatabaseSetStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single DropdbStmt
func parseDropdbStmt(t *testing.T, input string) *nodes.DropdbStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.DropdbStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropdbStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single AlterSystemStmt
func parseAlterSystemStmt(t *testing.T, input string) *nodes.AlterSystemStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.AlterSystemStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterSystemStmt, got %T", result.Items[0])
	}
	return stmt
}

// === CREATE DATABASE tests ===

func TestCreateDatabaseSimple(t *testing.T) {
	stmt := parseCreatedbStmt(t, "CREATE DATABASE mydb")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Options != nil && len(stmt.Options.Items) != 0 {
		t.Errorf("expected no options, got %v", stmt.Options)
	}
}

func TestCreateDatabaseOwner(t *testing.T) {
	stmt := parseCreatedbStmt(t, "CREATE DATABASE mydb OWNER myuser")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "owner" {
		t.Errorf("expected option name 'owner', got %q", opt.Defname)
	}
	arg, ok := opt.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected string arg, got %T", opt.Arg)
	}
	if arg.Str != "myuser" {
		t.Errorf("expected 'myuser', got %q", arg.Str)
	}
}

func TestCreateDatabaseEncoding(t *testing.T) {
	stmt := parseCreatedbStmt(t, "CREATE DATABASE mydb ENCODING 'UTF8'")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "encoding" {
		t.Errorf("expected option name 'encoding', got %q", opt.Defname)
	}
}

func TestCreateDatabaseTemplate(t *testing.T) {
	stmt := parseCreatedbStmt(t, "CREATE DATABASE mydb TEMPLATE template0")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "template" {
		t.Errorf("expected option name 'template', got %q", opt.Defname)
	}
}

func TestCreateDatabaseMultipleOptions(t *testing.T) {
	stmt := parseCreatedbStmt(t, "CREATE DATABASE mydb OWNER myuser ENCODING 'UTF8' TEMPLATE template0")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 3 {
		t.Fatalf("expected 3 options, got %d", len(stmt.Options.Items))
	}
	// Check option names
	names := []string{"owner", "encoding", "template"}
	for i, name := range names {
		opt := stmt.Options.Items[i].(*nodes.DefElem)
		if opt.Defname != name {
			t.Errorf("option[%d]: expected %q, got %q", i, name, opt.Defname)
		}
	}
}

func TestCreateDatabaseConnectionLimit(t *testing.T) {
	stmt := parseCreatedbStmt(t, "CREATE DATABASE mydb CONNECTION LIMIT 100")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "connection_limit" {
		t.Errorf("expected option name 'connection_limit', got %q", opt.Defname)
	}
	arg, ok := opt.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected Integer arg, got %T", opt.Arg)
	}
	if arg.Ival != 100 {
		t.Errorf("expected 100, got %d", arg.Ival)
	}
}

// === ALTER DATABASE tests ===

func TestAlterDatabaseSetVariable(t *testing.T) {
	stmt := parseAlterDatabaseSetStmt(t, "ALTER DATABASE mydb SET search_path TO public")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected VAR_SET_VALUE, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "search_path" {
		t.Errorf("expected 'search_path', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterDatabaseReset(t *testing.T) {
	stmt := parseAlterDatabaseSetStmt(t, "ALTER DATABASE mydb RESET search_path")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_RESET {
		t.Errorf("expected VAR_RESET, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "search_path" {
		t.Errorf("expected 'search_path', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterDatabaseResetAll(t *testing.T) {
	stmt := parseAlterDatabaseSetStmt(t, "ALTER DATABASE mydb RESET ALL")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_RESET_ALL {
		t.Errorf("expected VAR_RESET_ALL, got %d", stmt.Setstmt.Kind)
	}
}

func TestAlterDatabaseSetTablespace(t *testing.T) {
	stmt := parseAlterDatabaseStmt(t, "ALTER DATABASE mydb SET TABLESPACE newtbs")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "tablespace" {
		t.Errorf("expected option name 'tablespace', got %q", opt.Defname)
	}
	arg, ok := opt.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected string arg, got %T", opt.Arg)
	}
	if arg.Str != "newtbs" {
		t.Errorf("expected 'newtbs', got %q", arg.Str)
	}
}

// === DROP DATABASE tests ===

func TestDropDatabase(t *testing.T) {
	stmt := parseDropdbStmt(t, "DROP DATABASE mydb")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.MissingOk {
		t.Error("expected MissingOk false")
	}
	if stmt.Options != nil && len(stmt.Options.Items) != 0 {
		t.Errorf("expected no options, got %v", stmt.Options)
	}
}

func TestDropDatabaseIfExists(t *testing.T) {
	stmt := parseDropdbStmt(t, "DROP DATABASE IF EXISTS mydb")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if !stmt.MissingOk {
		t.Error("expected MissingOk true")
	}
}

func TestDropDatabaseForce(t *testing.T) {
	stmt := parseDropdbStmt(t, "DROP DATABASE mydb (FORCE)")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if stmt.MissingOk {
		t.Error("expected MissingOk false")
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "force" {
		t.Errorf("expected option name 'force', got %q", opt.Defname)
	}
}

func TestDropDatabaseIfExistsForce(t *testing.T) {
	stmt := parseDropdbStmt(t, "DROP DATABASE IF EXISTS mydb (FORCE)")
	if stmt.Dbname != "mydb" {
		t.Errorf("expected dbname 'mydb', got %q", stmt.Dbname)
	}
	if !stmt.MissingOk {
		t.Error("expected MissingOk true")
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "force" {
		t.Errorf("expected option name 'force', got %q", opt.Defname)
	}
}

// === ALTER SYSTEM tests ===

func TestAlterSystemSet(t *testing.T) {
	stmt := parseAlterSystemStmt(t, "ALTER SYSTEM SET wal_level = replica")
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected VAR_SET_VALUE, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "wal_level" {
		t.Errorf("expected 'wal_level', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterSystemSetMaxConnections(t *testing.T) {
	stmt := parseAlterSystemStmt(t, "ALTER SYSTEM SET max_connections = 200")
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected VAR_SET_VALUE, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "max_connections" {
		t.Errorf("expected 'max_connections', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterSystemReset(t *testing.T) {
	stmt := parseAlterSystemStmt(t, "ALTER SYSTEM RESET wal_level")
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_RESET {
		t.Errorf("expected VAR_RESET, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "wal_level" {
		t.Errorf("expected 'wal_level', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterSystemResetAll(t *testing.T) {
	stmt := parseAlterSystemStmt(t, "ALTER SYSTEM RESET ALL")
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_RESET_ALL {
		t.Errorf("expected VAR_RESET_ALL, got %d", stmt.Setstmt.Kind)
	}
}
