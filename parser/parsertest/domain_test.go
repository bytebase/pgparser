package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestCreateDomainBasic tests: CREATE DOMAIN d AS integer
func TestCreateDomainBasic(t *testing.T) {
	input := "CREATE DOMAIN d AS integer"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.CreateDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateDomainStmt, got %T", result.Items[0])
	}

	if stmt.Domainname == nil || len(stmt.Domainname.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.Domainname.Items[0].(*nodes.String).Str != "d" {
		t.Errorf("expected domain name 'd', got %q", stmt.Domainname.Items[0].(*nodes.String).Str)
	}
	if stmt.Typname == nil {
		t.Fatal("expected type name to be set")
	}
}

// TestCreateDomainSchemaQualified tests: CREATE DOMAIN myschema.d AS text
func TestCreateDomainSchemaQualified(t *testing.T) {
	input := "CREATE DOMAIN myschema.d AS text"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateDomainStmt, got %T", result.Items[0])
	}

	if stmt.Domainname == nil || len(stmt.Domainname.Items) != 2 {
		t.Fatalf("expected 2 name components, got %d", stmt.Domainname.Len())
	}
	if stmt.Domainname.Items[0].(*nodes.String).Str != "myschema" {
		t.Errorf("expected 'myschema', got %q", stmt.Domainname.Items[0].(*nodes.String).Str)
	}
	if stmt.Domainname.Items[1].(*nodes.String).Str != "d" {
		t.Errorf("expected 'd', got %q", stmt.Domainname.Items[1].(*nodes.String).Str)
	}
}

// TestCreateDomainNotNull tests: CREATE DOMAIN d AS integer NOT NULL
func TestCreateDomainNotNull(t *testing.T) {
	input := "CREATE DOMAIN d AS integer NOT NULL"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateDomainStmt, got %T", result.Items[0])
	}

	if stmt.Constraints == nil || len(stmt.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %d", stmt.Constraints.Len())
	}
	c := stmt.Constraints.Items[0].(*nodes.Constraint)
	if c.Contype != nodes.CONSTR_NOTNULL {
		t.Errorf("expected CONSTR_NOTNULL, got %d", c.Contype)
	}
}

// TestCreateDomainDefault tests: CREATE DOMAIN d AS text DEFAULT 'hello'
func TestCreateDomainDefault(t *testing.T) {
	input := "CREATE DOMAIN d AS text DEFAULT 'hello'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateDomainStmt, got %T", result.Items[0])
	}

	if stmt.Constraints == nil || len(stmt.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %d", stmt.Constraints.Len())
	}
	c := stmt.Constraints.Items[0].(*nodes.Constraint)
	if c.Contype != nodes.CONSTR_DEFAULT {
		t.Errorf("expected CONSTR_DEFAULT, got %d", c.Contype)
	}
}

// TestCreateDomainCheck tests: CREATE DOMAIN d AS integer CHECK (VALUE > 0)
func TestCreateDomainCheck(t *testing.T) {
	input := "CREATE DOMAIN d AS integer CHECK (VALUE > 0)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateDomainStmt, got %T", result.Items[0])
	}

	if stmt.Constraints == nil || len(stmt.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %d", stmt.Constraints.Len())
	}
	c := stmt.Constraints.Items[0].(*nodes.Constraint)
	if c.Contype != nodes.CONSTR_CHECK {
		t.Errorf("expected CONSTR_CHECK, got %d", c.Contype)
	}
}

// TestCreateDomainWithoutAs tests: CREATE DOMAIN d integer
func TestCreateDomainWithoutAs(t *testing.T) {
	input := "CREATE DOMAIN d integer"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateDomainStmt, got %T", result.Items[0])
	}

	if stmt.Typname == nil {
		t.Fatal("expected type name to be set")
	}
}

// TestAlterDomainSetDefault tests: ALTER DOMAIN d SET DEFAULT 5
func TestAlterDomainSetDefault(t *testing.T) {
	input := "ALTER DOMAIN d SET DEFAULT 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'T' {
		t.Errorf("expected subtype 'T', got %c", stmt.Subtype)
	}
	if stmt.Def == nil {
		t.Fatal("expected Def to be set")
	}
}

// TestAlterDomainDropDefault tests: ALTER DOMAIN d DROP DEFAULT
func TestAlterDomainDropDefault(t *testing.T) {
	input := "ALTER DOMAIN d DROP DEFAULT"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'T' {
		t.Errorf("expected subtype 'T', got %c", stmt.Subtype)
	}
	if stmt.Def != nil {
		t.Error("expected Def to be nil for DROP DEFAULT")
	}
}

// TestAlterDomainSetNotNull tests: ALTER DOMAIN d SET NOT NULL
func TestAlterDomainSetNotNull(t *testing.T) {
	input := "ALTER DOMAIN d SET NOT NULL"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'O' {
		t.Errorf("expected subtype 'O', got %c", stmt.Subtype)
	}
}

// TestAlterDomainDropNotNull tests: ALTER DOMAIN d DROP NOT NULL
func TestAlterDomainDropNotNull(t *testing.T) {
	input := "ALTER DOMAIN d DROP NOT NULL"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'N' {
		t.Errorf("expected subtype 'N', got %c", stmt.Subtype)
	}
}

// TestAlterDomainAddConstraint tests: ALTER DOMAIN d ADD CONSTRAINT c CHECK (VALUE > 0)
func TestAlterDomainAddConstraint(t *testing.T) {
	input := "ALTER DOMAIN d ADD CONSTRAINT c CHECK (VALUE > 0)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'C' {
		t.Errorf("expected subtype 'C', got %c", stmt.Subtype)
	}
	if stmt.Def == nil {
		t.Fatal("expected Def to be set")
	}
}

// TestAlterDomainDropConstraint tests: ALTER DOMAIN d DROP CONSTRAINT c
func TestAlterDomainDropConstraint(t *testing.T) {
	input := "ALTER DOMAIN d DROP CONSTRAINT c"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'X' {
		t.Errorf("expected subtype 'X', got %c", stmt.Subtype)
	}
	if stmt.Name != "c" {
		t.Errorf("expected name 'c', got %q", stmt.Name)
	}
}

// TestAlterDomainDropConstraintIfExistsCascade tests: ALTER DOMAIN d DROP CONSTRAINT IF EXISTS c CASCADE
func TestAlterDomainDropConstraintIfExistsCascade(t *testing.T) {
	input := "ALTER DOMAIN d DROP CONSTRAINT IF EXISTS c CASCADE"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'X' {
		t.Errorf("expected subtype 'X', got %c", stmt.Subtype)
	}
	if !stmt.MissingOk {
		t.Error("expected MissingOk to be true")
	}
	if stmt.Name != "c" {
		t.Errorf("expected name 'c', got %q", stmt.Name)
	}
	if stmt.Behavior != nodes.DROP_CASCADE {
		t.Errorf("expected DROP_CASCADE, got %d", stmt.Behavior)
	}
}

// TestAlterDomainValidateConstraint tests: ALTER DOMAIN d VALIDATE CONSTRAINT c
func TestAlterDomainValidateConstraint(t *testing.T) {
	input := "ALTER DOMAIN d VALIDATE CONSTRAINT c"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterDomainStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterDomainStmt, got %T", result.Items[0])
	}

	if stmt.Subtype != 'V' {
		t.Errorf("expected subtype 'V', got %c", stmt.Subtype)
	}
	if stmt.Name != "c" {
		t.Errorf("expected name 'c', got %q", stmt.Name)
	}
}

// TestAlterEnumAddValue tests: ALTER TYPE myenum ADD VALUE 'new'
func TestAlterEnumAddValue(t *testing.T) {
	input := "ALTER TYPE myenum ADD VALUE 'new'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterEnumStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterEnumStmt, got %T", result.Items[0])
	}

	if stmt.Newval != "new" {
		t.Errorf("expected new value 'new', got %q", stmt.Newval)
	}
	if stmt.SkipIfNewvalExists {
		t.Error("expected SkipIfNewvalExists to be false")
	}
}

// TestAlterEnumAddValueIfNotExists tests: ALTER TYPE myenum ADD VALUE IF NOT EXISTS 'new'
func TestAlterEnumAddValueIfNotExists(t *testing.T) {
	input := "ALTER TYPE myenum ADD VALUE IF NOT EXISTS 'new'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterEnumStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterEnumStmt, got %T", result.Items[0])
	}

	if stmt.Newval != "new" {
		t.Errorf("expected new value 'new', got %q", stmt.Newval)
	}
	if !stmt.SkipIfNewvalExists {
		t.Error("expected SkipIfNewvalExists to be true")
	}
}

// TestAlterEnumAddValueBefore tests: ALTER TYPE myenum ADD VALUE 'new' BEFORE 'old'
func TestAlterEnumAddValueBefore(t *testing.T) {
	input := "ALTER TYPE myenum ADD VALUE 'new' BEFORE 'old'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterEnumStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterEnumStmt, got %T", result.Items[0])
	}

	if stmt.Newval != "new" {
		t.Errorf("expected new value 'new', got %q", stmt.Newval)
	}
	if stmt.NewvalNeighbor != "old" {
		t.Errorf("expected neighbor 'old', got %q", stmt.NewvalNeighbor)
	}
	if stmt.NewvalIsAfter {
		t.Error("expected NewvalIsAfter to be false")
	}
}

// TestAlterEnumAddValueAfter tests: ALTER TYPE myenum ADD VALUE 'new' AFTER 'old'
func TestAlterEnumAddValueAfter(t *testing.T) {
	input := "ALTER TYPE myenum ADD VALUE 'new' AFTER 'old'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterEnumStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterEnumStmt, got %T", result.Items[0])
	}

	if stmt.Newval != "new" {
		t.Errorf("expected new value 'new', got %q", stmt.Newval)
	}
	if stmt.NewvalNeighbor != "old" {
		t.Errorf("expected neighbor 'old', got %q", stmt.NewvalNeighbor)
	}
	if !stmt.NewvalIsAfter {
		t.Error("expected NewvalIsAfter to be true")
	}
}

// TestAlterEnumRenameValue tests: ALTER TYPE myenum RENAME VALUE 'old' TO 'new'
func TestAlterEnumRenameValue(t *testing.T) {
	input := "ALTER TYPE myenum RENAME VALUE 'old' TO 'new'"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterEnumStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterEnumStmt, got %T", result.Items[0])
	}

	if stmt.Oldval != "old" {
		t.Errorf("expected old value 'old', got %q", stmt.Oldval)
	}
	if stmt.Newval != "new" {
		t.Errorf("expected new value 'new', got %q", stmt.Newval)
	}
}

// TestAlterCollationRefreshVersion tests: ALTER COLLATION myschema.mycoll REFRESH VERSION
func TestAlterCollationRefreshVersion(t *testing.T) {
	input := "ALTER COLLATION myschema.mycoll REFRESH VERSION"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterCollationStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterCollationStmt, got %T", result.Items[0])
	}

	if stmt.Collname == nil || len(stmt.Collname.Items) != 2 {
		t.Fatalf("expected 2 name components, got %d", stmt.Collname.Len())
	}
	if stmt.Collname.Items[0].(*nodes.String).Str != "myschema" {
		t.Errorf("expected 'myschema', got %q", stmt.Collname.Items[0].(*nodes.String).Str)
	}
	if stmt.Collname.Items[1].(*nodes.String).Str != "mycoll" {
		t.Errorf("expected 'mycoll', got %q", stmt.Collname.Items[1].(*nodes.String).Str)
	}
}

// TestAlterCompositeTypeAddAttribute tests: ALTER TYPE mytype ADD ATTRIBUTE a integer
func TestAlterCompositeTypeAddAttribute(t *testing.T) {
	input := "ALTER TYPE mytype ADD ATTRIBUTE a integer"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableStmt, got %T", result.Items[0])
	}

	if stmt.ObjType != int(nodes.OBJECT_TYPE) {
		t.Errorf("expected OBJECT_TYPE, got %d", stmt.ObjType)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "mytype" {
		t.Fatal("expected relation name 'mytype'")
	}
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatalf("expected 1 command, got %d", stmt.Cmds.Len())
	}
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_AddColumn) {
		t.Errorf("expected AT_AddColumn, got %d", cmd.Subtype)
	}
}

// TestAlterCompositeTypeDropAttribute tests: ALTER TYPE mytype DROP ATTRIBUTE a
func TestAlterCompositeTypeDropAttribute(t *testing.T) {
	input := "ALTER TYPE mytype DROP ATTRIBUTE a"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableStmt, got %T", result.Items[0])
	}

	if stmt.ObjType != int(nodes.OBJECT_TYPE) {
		t.Errorf("expected OBJECT_TYPE, got %d", stmt.ObjType)
	}
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatalf("expected 1 command, got %d", stmt.Cmds.Len())
	}
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_DropColumn) {
		t.Errorf("expected AT_DropColumn, got %d", cmd.Subtype)
	}
	if cmd.Name != "a" {
		t.Errorf("expected name 'a', got %q", cmd.Name)
	}
}

// TestAlterCompositeTypeAlterAttribute tests: ALTER TYPE mytype ALTER ATTRIBUTE a TYPE text
func TestAlterCompositeTypeAlterAttribute(t *testing.T) {
	input := "ALTER TYPE mytype ALTER ATTRIBUTE a TYPE text"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterTableStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterTableStmt, got %T", result.Items[0])
	}

	if stmt.ObjType != int(nodes.OBJECT_TYPE) {
		t.Errorf("expected OBJECT_TYPE, got %d", stmt.ObjType)
	}
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatalf("expected 1 command, got %d", stmt.Cmds.Len())
	}
	cmd := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if cmd.Subtype != int(nodes.AT_AlterColumnType) {
		t.Errorf("expected AT_AlterColumnType, got %d", cmd.Subtype)
	}
	if cmd.Name != "a" {
		t.Errorf("expected name 'a', got %q", cmd.Name)
	}
}
