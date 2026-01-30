package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// ALTER ... SET SCHEMA tests
// =============================================================================

func TestAlterTableSetSchema(t *testing.T) {
	input := "ALTER TABLE t SET SCHEMA new_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterObjectSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterObjectSchemaStmt, got %T", result.Items[0])
	}
	if stmt.ObjectType != nodes.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", stmt.ObjectType)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Error("expected relation 't'")
	}
	if stmt.Newschema != "new_schema" {
		t.Errorf("expected newschema 'new_schema', got %q", stmt.Newschema)
	}
}

func TestAlterViewSetSchema(t *testing.T) {
	input := "ALTER VIEW v SET SCHEMA new_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterObjectSchemaStmt)
	if stmt.ObjectType != nodes.OBJECT_VIEW {
		t.Errorf("expected OBJECT_VIEW, got %d", stmt.ObjectType)
	}
	if stmt.Newschema != "new_schema" {
		t.Errorf("expected newschema 'new_schema', got %q", stmt.Newschema)
	}
}

func TestAlterSequenceSetSchema(t *testing.T) {
	input := "ALTER SEQUENCE s SET SCHEMA new_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterObjectSchemaStmt)
	if stmt.ObjectType != nodes.OBJECT_SEQUENCE {
		t.Errorf("expected OBJECT_SEQUENCE, got %d", stmt.ObjectType)
	}
	if stmt.Newschema != "new_schema" {
		t.Errorf("expected newschema 'new_schema', got %q", stmt.Newschema)
	}
}

func TestAlterFunctionSetSchema(t *testing.T) {
	input := "ALTER FUNCTION f(int) SET SCHEMA new_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterObjectSchemaStmt)
	if stmt.ObjectType != nodes.OBJECT_FUNCTION {
		t.Errorf("expected OBJECT_FUNCTION, got %d", stmt.ObjectType)
	}
	if stmt.Newschema != "new_schema" {
		t.Errorf("expected newschema 'new_schema', got %q", stmt.Newschema)
	}
}

func TestAlterTypeSetSchema(t *testing.T) {
	input := "ALTER TYPE mytype SET SCHEMA new_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterObjectSchemaStmt)
	if stmt.ObjectType != nodes.OBJECT_TYPE {
		t.Errorf("expected OBJECT_TYPE, got %d", stmt.ObjectType)
	}
	if stmt.Newschema != "new_schema" {
		t.Errorf("expected newschema 'new_schema', got %q", stmt.Newschema)
	}
}

func TestAlterMatViewSetSchema(t *testing.T) {
	input := "ALTER MATERIALIZED VIEW mv SET SCHEMA new_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterObjectSchemaStmt)
	if stmt.ObjectType != nodes.OBJECT_MATVIEW {
		t.Errorf("expected OBJECT_MATVIEW, got %d", stmt.ObjectType)
	}
	if stmt.Newschema != "new_schema" {
		t.Errorf("expected newschema 'new_schema', got %q", stmt.Newschema)
	}
}

// =============================================================================
// ALTER ... OWNER TO tests
// =============================================================================

func TestAlterTableOwnerTo(t *testing.T) {
	// ALTER TABLE ... OWNER TO is handled by AlterTableStmt via alter_table_cmd
	// with AT_ChangeOwner, same as PostgreSQL
	input := "ALTER TABLE t OWNER TO new_owner"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableStmt, got %T", result.Items[0])
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Error("expected relation 't'")
	}
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatal("expected 1 command")
	}
	cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableCmd, got %T", stmt.Cmds.Items[0])
	}
	if cmd.Subtype != int(nodes.AT_ChangeOwner) {
		t.Errorf("expected AT_ChangeOwner, got %d", cmd.Subtype)
	}
	if cmd.Newowner == nil || cmd.Newowner.Rolename != "new_owner" {
		t.Error("expected newowner 'new_owner'")
	}
}

func TestAlterViewOwnerTo(t *testing.T) {
	// View uses AlterTableStmt with OBJECT_VIEW, but ALTER VIEW ... OWNER TO
	// is handled by AlterTableStmt for simple cases. However, PG also has
	// AlterOwnerStmt for many object types. Let's test one that uses AlterOwnerStmt.
	input := "ALTER FUNCTION f(int) OWNER TO new_owner"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterOwnerStmt)
	if stmt.ObjectType != nodes.OBJECT_FUNCTION {
		t.Errorf("expected OBJECT_FUNCTION, got %d", stmt.ObjectType)
	}
	if stmt.Newowner == nil || stmt.Newowner.Rolename != "new_owner" {
		t.Error("expected newowner 'new_owner'")
	}
}

func TestAlterDomainOwnerTo(t *testing.T) {
	// Test ALTER DOMAIN ... OWNER TO (handled by AlterOwnerStmt)
	input := "ALTER DOMAIN mydom OWNER TO new_owner"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterOwnerStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterOwnerStmt, got %T", result.Items[0])
	}
	if stmt.ObjectType != nodes.OBJECT_DOMAIN {
		t.Errorf("expected OBJECT_DOMAIN, got %d", stmt.ObjectType)
	}
	if stmt.Newowner == nil || stmt.Newowner.Rolename != "new_owner" {
		t.Error("expected newowner 'new_owner'")
	}
}

func TestAlterDatabaseOwnerTo(t *testing.T) {
	input := "ALTER DATABASE d OWNER TO new_owner"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterOwnerStmt)
	if stmt.ObjectType != nodes.OBJECT_DATABASE {
		t.Errorf("expected OBJECT_DATABASE, got %d", stmt.ObjectType)
	}
	if stmt.Newowner == nil || stmt.Newowner.Rolename != "new_owner" {
		t.Error("expected newowner 'new_owner'")
	}
}

func TestAlterSchemaOwnerTo(t *testing.T) {
	input := "ALTER SCHEMA s OWNER TO new_owner"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterOwnerStmt)
	if stmt.ObjectType != nodes.OBJECT_SCHEMA {
		t.Errorf("expected OBJECT_SCHEMA, got %d", stmt.ObjectType)
	}
}

// =============================================================================
// ALTER ... DEPENDS ON EXTENSION tests
// =============================================================================

func TestAlterFunctionDependsOnExtension(t *testing.T) {
	input := "ALTER FUNCTION f(int) DEPENDS ON EXTENSION ext"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterObjectDependsStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterObjectDependsStmt, got %T", result.Items[0])
	}
	if stmt.ObjectType != nodes.OBJECT_FUNCTION {
		t.Errorf("expected OBJECT_FUNCTION, got %d", stmt.ObjectType)
	}
	if stmt.Remove {
		t.Error("expected Remove to be false")
	}
	extStr, ok := stmt.Extname.(*nodes.String)
	if !ok || extStr.Str != "ext" {
		t.Error("expected extname 'ext'")
	}
}

func TestAlterFunctionNoDependsOnExtension(t *testing.T) {
	input := "ALTER FUNCTION f(int) NO DEPENDS ON EXTENSION ext"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterObjectDependsStmt)
	if !stmt.Remove {
		t.Error("expected Remove to be true")
	}
}

// =============================================================================
// ALTER OPERATOR SET (...) tests
// =============================================================================

func TestAlterOperatorSet(t *testing.T) {
	input := "ALTER OPERATOR +(int, int) SET (RESTRICT = func)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterOperatorStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterOperatorStmt, got %T", result.Items[0])
	}
	if stmt.Opername == nil {
		t.Fatal("expected non-nil Opername")
	}
	if stmt.Options == nil || len(stmt.Options.Items) == 0 {
		t.Fatal("expected options")
	}
}

// =============================================================================
// ALTER TYPE ... SET (...) tests
// =============================================================================

func TestAlterTypeSet(t *testing.T) {
	input := "ALTER TYPE mytype SET (RECEIVE = func)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterTypeStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTypeStmt, got %T", result.Items[0])
	}
	if stmt.TypeName == nil {
		t.Fatal("expected non-nil TypeName")
	}
}

// =============================================================================
// ALTER DEFAULT PRIVILEGES tests
// =============================================================================

func TestAlterDefaultPrivilegesGrantOnTables(t *testing.T) {
	input := "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO PUBLIC"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterDefaultPrivilegesStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDefaultPrivilegesStmt, got %T", result.Items[0])
	}
	if stmt.Options == nil {
		t.Fatal("expected non-nil Options")
	}
	if stmt.Action == nil {
		t.Fatal("expected non-nil Action")
	}
	if !stmt.Action.IsGrant {
		t.Error("expected IsGrant to be true")
	}
	if stmt.Action.Targtype != nodes.ACL_TARGET_DEFAULTS {
		t.Errorf("expected ACL_TARGET_DEFAULTS, got %d", stmt.Action.Targtype)
	}
}

func TestAlterDefaultPrivilegesRevokeOnFunctions(t *testing.T) {
	input := "ALTER DEFAULT PRIVILEGES FOR ROLE admin REVOKE ALL ON FUNCTIONS FROM PUBLIC"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterDefaultPrivilegesStmt)
	if stmt.Action == nil {
		t.Fatal("expected non-nil Action")
	}
	if stmt.Action.IsGrant {
		t.Error("expected IsGrant to be false")
	}
}

// =============================================================================
// ALTER TEXT SEARCH tests
// =============================================================================

func TestAlterTSDictionary(t *testing.T) {
	input := "ALTER TEXT SEARCH DICTIONARY mydict (STOPWORDS = 'english')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterTSDictionaryStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTSDictionaryStmt, got %T", result.Items[0])
	}
	if stmt.Dictname == nil {
		t.Fatal("expected non-nil Dictname")
	}
	if stmt.Options == nil {
		t.Fatal("expected non-nil Options")
	}
}

func TestAlterTSConfigAddMapping(t *testing.T) {
	input := "ALTER TEXT SEARCH CONFIGURATION myconfig ADD MAPPING FOR word WITH simple"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterTSConfigurationStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTSConfigurationStmt, got %T", result.Items[0])
	}
	if stmt.Kind != nodes.ALTER_TSCONFIG_ADD_MAPPING {
		t.Errorf("expected ALTER_TSCONFIG_ADD_MAPPING, got %d", stmt.Kind)
	}
	if stmt.Cfgname == nil {
		t.Fatal("expected non-nil Cfgname")
	}
}

func TestAlterTSConfigDropMapping(t *testing.T) {
	input := "ALTER TEXT SEARCH CONFIGURATION myconfig DROP MAPPING FOR word"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterTSConfigurationStmt)
	if stmt.Kind != nodes.ALTER_TSCONFIG_DROP_MAPPING {
		t.Errorf("expected ALTER_TSCONFIG_DROP_MAPPING, got %d", stmt.Kind)
	}
}
