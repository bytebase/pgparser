package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestAlterTableAddColumn tests: ALTER TABLE t ADD COLUMN name text
func TestAlterTableAddColumn(t *testing.T) {
	input := "ALTER TABLE t ADD COLUMN name text"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableStmt, got %T", result.Items[0])
	}

	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatalf("expected table name 't'")
	}
	if stmt.ObjType != int(nodes.OBJECT_TABLE) {
		t.Errorf("expected ObjType OBJECT_TABLE, got %d", stmt.ObjType)
	}
	if stmt.Missing_ok {
		t.Error("expected Missing_ok to be false")
	}

	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatalf("expected 1 command, got %v", stmt.Cmds)
	}

	cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableCmd, got %T", stmt.Cmds.Items[0])
	}
	if cmd.Subtype != int(nodes.AT_AddColumn) {
		t.Errorf("expected AT_AddColumn, got %d", cmd.Subtype)
	}
	if cmd.Def == nil {
		t.Fatal("expected Def to be set")
	}
	coldef, ok := cmd.Def.(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnDef, got %T", cmd.Def)
	}
	if coldef.Colname != "name" {
		t.Errorf("expected column name 'name', got %q", coldef.Colname)
	}
}

// TestAlterTableAddColumnWithoutKeyword tests: ALTER TABLE t ADD name text (without COLUMN keyword)
func TestAlterTableAddColumnWithoutKeyword(t *testing.T) {
	input := "ALTER TABLE t ADD name text"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_AddColumn) {
		t.Errorf("expected AT_AddColumn, got %d", cmd.Subtype)
	}
	coldef := cmd.Def.(*nodes.ColumnDef)
	if coldef.Colname != "name" {
		t.Errorf("expected column name 'name', got %q", coldef.Colname)
	}
}

// TestAlterTableAddColumnIfNotExists tests: ALTER TABLE t ADD COLUMN IF NOT EXISTS name text
func TestAlterTableAddColumnIfNotExists(t *testing.T) {
	input := "ALTER TABLE t ADD COLUMN IF NOT EXISTS name text"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_AddColumn) {
		t.Errorf("expected AT_AddColumn, got %d", cmd.Subtype)
	}
	if !cmd.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	coldef := cmd.Def.(*nodes.ColumnDef)
	if coldef.Colname != "name" {
		t.Errorf("expected column name 'name', got %q", coldef.Colname)
	}
}

// TestAlterTableDropColumn tests: ALTER TABLE t DROP COLUMN name
func TestAlterTableDropColumn(t *testing.T) {
	input := "ALTER TABLE t DROP COLUMN name"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_DropColumn) {
		t.Errorf("expected AT_DropColumn, got %d", cmd.Subtype)
	}
	if cmd.Name != "name" {
		t.Errorf("expected column name 'name', got %q", cmd.Name)
	}
	if cmd.Behavior != int(nodes.DROP_RESTRICT) {
		t.Errorf("expected DROP_RESTRICT (default), got %d", cmd.Behavior)
	}
}

// TestAlterTableDropColumnIfExistsCascade tests: ALTER TABLE t DROP COLUMN IF EXISTS name CASCADE
func TestAlterTableDropColumnIfExistsCascade(t *testing.T) {
	input := "ALTER TABLE t DROP COLUMN IF EXISTS name CASCADE"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_DropColumn) {
		t.Errorf("expected AT_DropColumn, got %d", cmd.Subtype)
	}
	if cmd.Name != "name" {
		t.Errorf("expected column name 'name', got %q", cmd.Name)
	}
	if !cmd.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if cmd.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected DROP_CASCADE, got %d", cmd.Behavior)
	}
}

// TestAlterTableAlterColumnSetDefault tests: ALTER TABLE t ALTER COLUMN name SET DEFAULT 'unknown'
func TestAlterTableAlterColumnSetDefault(t *testing.T) {
	input := "ALTER TABLE t ALTER COLUMN name SET DEFAULT 'unknown'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_ColumnDefault) {
		t.Errorf("expected AT_ColumnDefault, got %d", cmd.Subtype)
	}
	if cmd.Name != "name" {
		t.Errorf("expected column name 'name', got %q", cmd.Name)
	}
	if cmd.Def == nil {
		t.Fatal("expected Def to be set for SET DEFAULT")
	}
	aConst, ok := cmd.Def.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", cmd.Def)
	}
	strVal, ok := aConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", aConst.Val)
	}
	if strVal.Str != "unknown" {
		t.Errorf("expected 'unknown', got %q", strVal.Str)
	}
}

// TestAlterTableAlterColumnDropDefault tests: ALTER TABLE t ALTER COLUMN name DROP DEFAULT
func TestAlterTableAlterColumnDropDefault(t *testing.T) {
	input := "ALTER TABLE t ALTER COLUMN name DROP DEFAULT"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_ColumnDefault) {
		t.Errorf("expected AT_ColumnDefault, got %d", cmd.Subtype)
	}
	if cmd.Name != "name" {
		t.Errorf("expected column name 'name', got %q", cmd.Name)
	}
	if cmd.Def != nil {
		t.Error("expected Def to be nil for DROP DEFAULT")
	}
}

// TestAlterTableAlterColumnSetNotNull tests: ALTER TABLE t ALTER COLUMN name SET NOT NULL
func TestAlterTableAlterColumnSetNotNull(t *testing.T) {
	input := "ALTER TABLE t ALTER COLUMN name SET NOT NULL"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_SetNotNull) {
		t.Errorf("expected AT_SetNotNull, got %d", cmd.Subtype)
	}
	if cmd.Name != "name" {
		t.Errorf("expected column name 'name', got %q", cmd.Name)
	}
}

// TestAlterTableAlterColumnDropNotNull tests: ALTER TABLE t ALTER COLUMN name DROP NOT NULL
func TestAlterTableAlterColumnDropNotNull(t *testing.T) {
	input := "ALTER TABLE t ALTER COLUMN name DROP NOT NULL"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_DropNotNull) {
		t.Errorf("expected AT_DropNotNull, got %d", cmd.Subtype)
	}
	if cmd.Name != "name" {
		t.Errorf("expected column name 'name', got %q", cmd.Name)
	}
}

// TestAlterTableAlterColumnType tests: ALTER TABLE t ALTER COLUMN name TYPE varchar(200)
func TestAlterTableAlterColumnType(t *testing.T) {
	input := "ALTER TABLE t ALTER COLUMN name TYPE varchar(200)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_AlterColumnType) {
		t.Errorf("expected AT_AlterColumnType, got %d", cmd.Subtype)
	}
	if cmd.Name != "name" {
		t.Errorf("expected column name 'name', got %q", cmd.Name)
	}
	if cmd.Def == nil {
		t.Fatal("expected Def to be set for TYPE change")
	}
	coldef, ok := cmd.Def.(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnDef, got %T", cmd.Def)
	}
	if coldef.TypeName == nil {
		t.Fatal("expected TypeName to be set")
	}
}

// TestAlterTableAddConstraint tests: ALTER TABLE t ADD CONSTRAINT pk_id PRIMARY KEY (id)
func TestAlterTableAddConstraint(t *testing.T) {
	input := "ALTER TABLE t ADD CONSTRAINT pk_id PRIMARY KEY (id)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_AddConstraint) {
		t.Errorf("expected AT_AddConstraint, got %d", cmd.Subtype)
	}
	if cmd.Def == nil {
		t.Fatal("expected Def to be set")
	}
	constr, ok := cmd.Def.(*nodes.Constraint)
	if !ok {
		t.Fatalf("expected *nodes.Constraint, got %T", cmd.Def)
	}
	if constr.Contype != nodes.CONSTR_PRIMARY {
		t.Errorf("expected CONSTR_PRIMARY, got %d", constr.Contype)
	}
	if constr.Conname != "pk_id" {
		t.Errorf("expected constraint name 'pk_id', got %q", constr.Conname)
	}
}

// TestAlterTableDropConstraint tests: ALTER TABLE t DROP CONSTRAINT pk_id
func TestAlterTableDropConstraint(t *testing.T) {
	input := "ALTER TABLE t DROP CONSTRAINT pk_id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_DropConstraint) {
		t.Errorf("expected AT_DropConstraint, got %d", cmd.Subtype)
	}
	if cmd.Name != "pk_id" {
		t.Errorf("expected constraint name 'pk_id', got %q", cmd.Name)
	}
	if cmd.Behavior != int(nodes.DROP_RESTRICT) {
		t.Errorf("expected DROP_RESTRICT (default), got %d", cmd.Behavior)
	}
}

// TestAlterTableDropConstraintIfExistsCascade tests: ALTER TABLE t DROP CONSTRAINT IF EXISTS pk_id CASCADE
func TestAlterTableDropConstraintIfExistsCascade(t *testing.T) {
	input := "ALTER TABLE t DROP CONSTRAINT IF EXISTS pk_id CASCADE"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_DropConstraint) {
		t.Errorf("expected AT_DropConstraint, got %d", cmd.Subtype)
	}
	if cmd.Name != "pk_id" {
		t.Errorf("expected constraint name 'pk_id', got %q", cmd.Name)
	}
	if !cmd.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if cmd.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected DROP_CASCADE, got %d", cmd.Behavior)
	}
}

// TestAlterTableMultipleCommands tests: ALTER TABLE t ADD COLUMN a integer, DROP COLUMN b
func TestAlterTableMultipleCommands(t *testing.T) {
	input := "ALTER TABLE t ADD COLUMN a integer, DROP COLUMN b"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 2 {
		t.Fatalf("expected 2 commands, got %v", stmt.Cmds)
	}

	// First command: ADD COLUMN a integer
	cmd1 := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd1.Subtype != int(nodes.AT_AddColumn) {
		t.Errorf("cmd1: expected AT_AddColumn, got %d", cmd1.Subtype)
	}
	coldef := cmd1.Def.(*nodes.ColumnDef)
	if coldef.Colname != "a" {
		t.Errorf("cmd1: expected column 'a', got %q", coldef.Colname)
	}

	// Second command: DROP COLUMN b
	cmd2 := stmt.Cmds.Items[1].(*nodes.AlterTableCmd)
	if cmd2.Subtype != int(nodes.AT_DropColumn) {
		t.Errorf("cmd2: expected AT_DropColumn, got %d", cmd2.Subtype)
	}
	if cmd2.Name != "b" {
		t.Errorf("cmd2: expected column 'b', got %q", cmd2.Name)
	}
}

// TestAlterTableIfNotExists tests: ALTER TABLE IF EXISTS t ADD COLUMN name text
func TestAlterTableIfNotExists(t *testing.T) {
	input := "ALTER TABLE IF EXISTS t ADD COLUMN name text"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}
}

// TestAlterTableRenameTo tests: ALTER TABLE t RENAME TO t2
func TestAlterTableRenameTo(t *testing.T) {
	input := "ALTER TABLE t RENAME TO t2"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.RenameStmt)
	if !ok {
		t.Fatalf("expected *nodes.RenameStmt, got %T", result.Items[0])
	}
	if stmt.RenameType != nodes.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", stmt.RenameType)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}
	if stmt.Newname != "t2" {
		t.Errorf("expected new name 't2', got %q", stmt.Newname)
	}
	if stmt.MissingOk {
		t.Error("expected MissingOk to be false")
	}
}

// TestAlterTableRenameIfExists tests: ALTER TABLE IF EXISTS t RENAME TO t2
func TestAlterTableRenameIfExists(t *testing.T) {
	input := "ALTER TABLE IF EXISTS t RENAME TO t2"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.RenameStmt)
	if stmt.RenameType != nodes.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", stmt.RenameType)
	}
	if !stmt.MissingOk {
		t.Error("expected MissingOk to be true")
	}
	if stmt.Newname != "t2" {
		t.Errorf("expected new name 't2', got %q", stmt.Newname)
	}
}

// TestAlterTableRenameColumnExplicit tests: ALTER TABLE t RENAME COLUMN old_name TO new_name
func TestAlterTableRenameColumnExplicit(t *testing.T) {
	input := "ALTER TABLE t RENAME COLUMN old_name TO new_name"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.RenameStmt)
	if stmt.RenameType != nodes.OBJECT_COLUMN {
		t.Errorf("expected OBJECT_COLUMN, got %d", stmt.RenameType)
	}
	if stmt.RelationType != nodes.OBJECT_TABLE {
		t.Errorf("expected RelationType OBJECT_TABLE, got %d", stmt.RelationType)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatal("expected table name 't'")
	}
	if stmt.Subname != "old_name" {
		t.Errorf("expected old name 'old_name', got %q", stmt.Subname)
	}
	if stmt.Newname != "new_name" {
		t.Errorf("expected new name 'new_name', got %q", stmt.Newname)
	}
}

// TestAlterTableRenameColumnImplicit tests: ALTER TABLE t RENAME old_name TO new_name (without COLUMN keyword)
func TestAlterTableRenameColumnImplicit(t *testing.T) {
	input := "ALTER TABLE t RENAME old_name TO new_name"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.RenameStmt)
	if stmt.RenameType != nodes.OBJECT_COLUMN {
		t.Errorf("expected OBJECT_COLUMN, got %d", stmt.RenameType)
	}
	if stmt.Subname != "old_name" {
		t.Errorf("expected old name 'old_name', got %q", stmt.Subname)
	}
	if stmt.Newname != "new_name" {
		t.Errorf("expected new name 'new_name', got %q", stmt.Newname)
	}
}

// TestAlterTableRenameConstraint tests: ALTER TABLE t RENAME CONSTRAINT old_con TO new_con
func TestAlterTableRenameConstraint(t *testing.T) {
	input := "ALTER TABLE t RENAME CONSTRAINT old_con TO new_con"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.RenameStmt)
	if stmt.RenameType != nodes.OBJECT_TABCONSTRAINT {
		t.Errorf("expected OBJECT_TABCONSTRAINT, got %d", stmt.RenameType)
	}
	if stmt.RelationType != nodes.OBJECT_TABLE {
		t.Errorf("expected RelationType OBJECT_TABLE, got %d", stmt.RelationType)
	}
	if stmt.Subname != "old_con" {
		t.Errorf("expected old name 'old_con', got %q", stmt.Subname)
	}
	if stmt.Newname != "new_con" {
		t.Errorf("expected new name 'new_con', got %q", stmt.Newname)
	}
}

// TestAlterTableSchemaQualified tests: ALTER TABLE myschema.t ADD COLUMN name text
func TestAlterTableSchemaQualified(t *testing.T) {
	input := "ALTER TABLE myschema.t ADD COLUMN name text"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Schemaname != "myschema" {
		t.Errorf("expected schema 'myschema', got %q", stmt.Relation.Schemaname)
	}
	if stmt.Relation.Relname != "t" {
		t.Errorf("expected table 't', got %q", stmt.Relation.Relname)
	}
}

// TestAlterTableDropColumnRestrict tests: ALTER TABLE t DROP COLUMN name RESTRICT
func TestAlterTableDropColumnRestrict(t *testing.T) {
	input := "ALTER TABLE t DROP COLUMN name RESTRICT"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Behavior != int(nodes.DROP_RESTRICT) {
		t.Errorf("expected DROP_RESTRICT, got %d", cmd.Behavior)
	}
}

// TestAlterTableDropColumnCascade tests: ALTER TABLE t DROP COLUMN name CASCADE
func TestAlterTableDropColumnCascade(t *testing.T) {
	input := "ALTER TABLE t DROP COLUMN name CASCADE"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected DROP_CASCADE, got %d", cmd.Behavior)
	}
}

// TestAlterTableTableDriven is a table-driven test for various ALTER TABLE patterns.
func TestAlterTableTableDriven(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"add column with COLUMN keyword", "ALTER TABLE t ADD COLUMN name text"},
		{"add column without COLUMN keyword", "ALTER TABLE t ADD name text"},
		{"add column if not exists", "ALTER TABLE t ADD COLUMN IF NOT EXISTS name text"},
		{"drop column", "ALTER TABLE t DROP COLUMN name"},
		{"drop column cascade", "ALTER TABLE t DROP COLUMN name CASCADE"},
		{"drop column restrict", "ALTER TABLE t DROP COLUMN name RESTRICT"},
		{"drop column if exists", "ALTER TABLE t DROP COLUMN IF EXISTS name"},
		{"drop column if exists cascade", "ALTER TABLE t DROP COLUMN IF EXISTS name CASCADE"},
		{"alter column set default int", "ALTER TABLE t ALTER COLUMN id SET DEFAULT 0"},
		{"alter column set default string", "ALTER TABLE t ALTER COLUMN name SET DEFAULT 'unknown'"},
		{"alter column drop default", "ALTER TABLE t ALTER COLUMN name DROP DEFAULT"},
		{"alter column set not null", "ALTER TABLE t ALTER COLUMN name SET NOT NULL"},
		{"alter column drop not null", "ALTER TABLE t ALTER COLUMN name DROP NOT NULL"},
		{"alter column type", "ALTER TABLE t ALTER COLUMN name TYPE varchar(200)"},
		{"alter column type integer", "ALTER TABLE t ALTER COLUMN id TYPE bigint"},
		{"add constraint primary key", "ALTER TABLE t ADD CONSTRAINT pk_id PRIMARY KEY (id)"},
		{"add constraint unique", "ALTER TABLE t ADD UNIQUE (name)"},
		{"add constraint check", "ALTER TABLE t ADD CHECK (age > 0)"},
		{"add constraint foreign key", "ALTER TABLE t ADD FOREIGN KEY (user_id) REFERENCES users(id)"},
		{"drop constraint", "ALTER TABLE t DROP CONSTRAINT pk_id"},
		{"drop constraint cascade", "ALTER TABLE t DROP CONSTRAINT pk_id CASCADE"},
		{"drop constraint if exists", "ALTER TABLE t DROP CONSTRAINT IF EXISTS pk_id"},
		{"drop constraint if exists cascade", "ALTER TABLE t DROP CONSTRAINT IF EXISTS pk_id CASCADE"},
		{"multiple commands", "ALTER TABLE t ADD COLUMN a integer, DROP COLUMN b, ALTER COLUMN c SET NOT NULL"},
		{"if exists", "ALTER TABLE IF EXISTS t ADD COLUMN name text"},
		{"schema-qualified", "ALTER TABLE myschema.t ADD COLUMN name text"},
		{"rename table", "ALTER TABLE t RENAME TO t2"},
		{"rename table if exists", "ALTER TABLE IF EXISTS t RENAME TO t2"},
		{"rename column explicit", "ALTER TABLE t RENAME COLUMN old_col TO new_col"},
		{"rename column implicit", "ALTER TABLE t RENAME old_col TO new_col"},
		{"rename constraint", "ALTER TABLE t RENAME CONSTRAINT old_con TO new_con"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse error for %q: %v", tt.input, err)
			}

			if result == nil || len(result.Items) == 0 {
				t.Fatalf("Parse returned empty result for %q", tt.input)
			}

			// Check it's either AlterTableStmt or RenameStmt
			switch result.Items[0].(type) {
			case *nodes.AlterTableStmt:
				// OK
			case *nodes.RenameStmt:
				// OK
			default:
				t.Fatalf("expected *nodes.AlterTableStmt or *nodes.RenameStmt for %q, got %T", tt.input, result.Items[0])
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}
