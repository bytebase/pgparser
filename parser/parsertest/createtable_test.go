package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestCreateTableBasic tests: CREATE TABLE t (id integer, name varchar(100))
func TestCreateTableBasic(t *testing.T) {
	input := "CREATE TABLE t (id integer, name varchar(100))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateStmt, got %T", result.Items[0])
	}

	// Verify relation
	if stmt.Relation == nil {
		t.Fatal("expected Relation to be set")
	}
	if stmt.Relation.Relname != "t" {
		t.Errorf("expected table name 't', got %q", stmt.Relation.Relname)
	}

	// Verify IfNotExists is false
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}

	// Verify 2 table elements
	if stmt.TableElts == nil || len(stmt.TableElts.Items) != 2 {
		t.Fatalf("expected 2 table elements, got %v", stmt.TableElts)
	}

	// First column: id integer
	col1, ok := stmt.TableElts.Items[0].(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnDef for first element, got %T", stmt.TableElts.Items[0])
	}
	if col1.Colname != "id" {
		t.Errorf("expected column name 'id', got %q", col1.Colname)
	}
	if col1.TypeName == nil {
		t.Fatal("expected TypeName for first column")
	}
	if !col1.IsLocal {
		t.Error("expected IsLocal to be true")
	}

	// Second column: name varchar(100)
	col2, ok := stmt.TableElts.Items[1].(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnDef for second element, got %T", stmt.TableElts.Items[1])
	}
	if col2.Colname != "name" {
		t.Errorf("expected column name 'name', got %q", col2.Colname)
	}
	if col2.TypeName == nil {
		t.Fatal("expected TypeName for second column")
	}
}

// TestCreateTableIfNotExists tests: CREATE TABLE IF NOT EXISTS t (id integer)
func TestCreateTableIfNotExists(t *testing.T) {
	input := "CREATE TABLE IF NOT EXISTS t (id integer)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateStmt, got %T", result.Items[0])
	}

	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}

	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatalf("expected table name 't'")
	}

	if stmt.TableElts == nil || len(stmt.TableElts.Items) != 1 {
		t.Fatalf("expected 1 table element, got %v", stmt.TableElts)
	}
}

// TestCreateTableColumnNotNull tests: CREATE TABLE t (id integer NOT NULL)
func TestCreateTableColumnNotNull(t *testing.T) {
	input := "CREATE TABLE t (id integer NOT NULL)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)

	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", col.Constraints)
	}

	constr, ok := col.Constraints.Items[0].(*nodes.Constraint)
	if !ok {
		t.Fatalf("expected *nodes.Constraint, got %T", col.Constraints.Items[0])
	}
	if constr.Contype != nodes.CONSTR_NOTNULL {
		t.Errorf("expected CONSTR_NOTNULL, got %d", constr.Contype)
	}
}

// TestCreateTableColumnDefault tests: CREATE TABLE t (id integer DEFAULT 0)
func TestCreateTableColumnDefault(t *testing.T) {
	input := "CREATE TABLE t (id integer DEFAULT 0)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)

	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", col.Constraints)
	}

	constr := col.Constraints.Items[0].(*nodes.Constraint)
	if constr.Contype != nodes.CONSTR_DEFAULT {
		t.Errorf("expected CONSTR_DEFAULT, got %d", constr.Contype)
	}
	if constr.RawExpr == nil {
		t.Fatal("expected RawExpr for DEFAULT constraint")
	}

	// Verify default value is 0
	aConst, ok := constr.RawExpr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", constr.RawExpr)
	}
	intVal, ok := aConst.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", aConst.Val)
	}
	if intVal.Ival != 0 {
		t.Errorf("expected default 0, got %d", intVal.Ival)
	}
}

// TestCreateTableColumnPrimaryKey tests: CREATE TABLE t (id integer PRIMARY KEY)
func TestCreateTableColumnPrimaryKey(t *testing.T) {
	input := "CREATE TABLE t (id integer PRIMARY KEY)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)

	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", col.Constraints)
	}

	constr := col.Constraints.Items[0].(*nodes.Constraint)
	if constr.Contype != nodes.CONSTR_PRIMARY {
		t.Errorf("expected CONSTR_PRIMARY, got %d", constr.Contype)
	}
}

// TestCreateTableColumnUnique tests: CREATE TABLE t (name varchar UNIQUE)
func TestCreateTableColumnUnique(t *testing.T) {
	input := "CREATE TABLE t (name varchar UNIQUE)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)

	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", col.Constraints)
	}

	constr := col.Constraints.Items[0].(*nodes.Constraint)
	if constr.Contype != nodes.CONSTR_UNIQUE {
		t.Errorf("expected CONSTR_UNIQUE, got %d", constr.Contype)
	}
}

// TestCreateTableColumnCheck tests: CREATE TABLE t (age integer CHECK (age > 0))
func TestCreateTableColumnCheck(t *testing.T) {
	input := "CREATE TABLE t (age integer CHECK (age > 0))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)

	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", col.Constraints)
	}

	constr := col.Constraints.Items[0].(*nodes.Constraint)
	if constr.Contype != nodes.CONSTR_CHECK {
		t.Errorf("expected CONSTR_CHECK, got %d", constr.Contype)
	}
	if constr.RawExpr == nil {
		t.Fatal("expected RawExpr for CHECK constraint")
	}

	// RawExpr should be an A_Expr (age > 0)
	aexpr, ok := constr.RawExpr.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", constr.RawExpr)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
}

// TestCreateTableColumnReferences tests: CREATE TABLE t (user_id integer REFERENCES users(id))
func TestCreateTableColumnReferences(t *testing.T) {
	input := "CREATE TABLE t (user_id integer REFERENCES users(id))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)

	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", col.Constraints)
	}

	constr := col.Constraints.Items[0].(*nodes.Constraint)
	if constr.Contype != nodes.CONSTR_FOREIGN {
		t.Errorf("expected CONSTR_FOREIGN, got %d", constr.Contype)
	}
	if constr.Pktable == nil {
		t.Fatal("expected Pktable for REFERENCES constraint")
	}
	if constr.Pktable.Relname != "users" {
		t.Errorf("expected referenced table 'users', got %q", constr.Pktable.Relname)
	}
	if constr.PkAttrs == nil || len(constr.PkAttrs.Items) != 1 {
		t.Fatalf("expected 1 referenced column, got %v", constr.PkAttrs)
	}
	refCol := constr.PkAttrs.Items[0].(*nodes.String)
	if refCol.Str != "id" {
		t.Errorf("expected referenced column 'id', got %q", refCol.Str)
	}
}

// TestCreateTableTablePrimaryKey tests: CREATE TABLE t (id integer, name text, PRIMARY KEY (id))
func TestCreateTableTablePrimaryKey(t *testing.T) {
	input := "CREATE TABLE t (id integer, name text, PRIMARY KEY (id))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	if stmt.TableElts == nil || len(stmt.TableElts.Items) != 3 {
		t.Fatalf("expected 3 table elements (2 columns + 1 constraint), got %d", len(stmt.TableElts.Items))
	}

	// Third element should be a table constraint
	constr, ok := stmt.TableElts.Items[2].(*nodes.Constraint)
	if !ok {
		t.Fatalf("expected *nodes.Constraint for third element, got %T", stmt.TableElts.Items[2])
	}
	if constr.Contype != nodes.CONSTR_PRIMARY {
		t.Errorf("expected CONSTR_PRIMARY, got %d", constr.Contype)
	}
	if constr.Keys == nil || len(constr.Keys.Items) != 1 {
		t.Fatalf("expected 1 key column, got %v", constr.Keys)
	}
	keyCol := constr.Keys.Items[0].(*nodes.String)
	if keyCol.Str != "id" {
		t.Errorf("expected key column 'id', got %q", keyCol.Str)
	}
}

// TestCreateTableTableUnique tests: CREATE TABLE t (id integer, name text, UNIQUE (name))
func TestCreateTableTableUnique(t *testing.T) {
	input := "CREATE TABLE t (id integer, name text, UNIQUE (name))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	if stmt.TableElts == nil || len(stmt.TableElts.Items) != 3 {
		t.Fatalf("expected 3 table elements, got %d", len(stmt.TableElts.Items))
	}

	constr, ok := stmt.TableElts.Items[2].(*nodes.Constraint)
	if !ok {
		t.Fatalf("expected *nodes.Constraint, got %T", stmt.TableElts.Items[2])
	}
	if constr.Contype != nodes.CONSTR_UNIQUE {
		t.Errorf("expected CONSTR_UNIQUE, got %d", constr.Contype)
	}
	if constr.Keys == nil || len(constr.Keys.Items) != 1 {
		t.Fatalf("expected 1 key, got %v", constr.Keys)
	}
	keyCol := constr.Keys.Items[0].(*nodes.String)
	if keyCol.Str != "name" {
		t.Errorf("expected key column 'name', got %q", keyCol.Str)
	}
}

// TestCreateTableTableForeignKey tests: CREATE TABLE t (id integer, FOREIGN KEY (id) REFERENCES t2(id))
func TestCreateTableTableForeignKey(t *testing.T) {
	input := "CREATE TABLE t (id integer, FOREIGN KEY (id) REFERENCES t2(id))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	if stmt.TableElts == nil || len(stmt.TableElts.Items) != 2 {
		t.Fatalf("expected 2 table elements, got %d", len(stmt.TableElts.Items))
	}

	constr, ok := stmt.TableElts.Items[1].(*nodes.Constraint)
	if !ok {
		t.Fatalf("expected *nodes.Constraint, got %T", stmt.TableElts.Items[1])
	}
	if constr.Contype != nodes.CONSTR_FOREIGN {
		t.Errorf("expected CONSTR_FOREIGN, got %d", constr.Contype)
	}

	// FK columns
	if constr.FkAttrs == nil || len(constr.FkAttrs.Items) != 1 {
		t.Fatalf("expected 1 FK column, got %v", constr.FkAttrs)
	}
	fkCol := constr.FkAttrs.Items[0].(*nodes.String)
	if fkCol.Str != "id" {
		t.Errorf("expected FK column 'id', got %q", fkCol.Str)
	}

	// Referenced table
	if constr.Pktable == nil || constr.Pktable.Relname != "t2" {
		t.Errorf("expected referenced table 't2'")
	}

	// Referenced columns
	if constr.PkAttrs == nil || len(constr.PkAttrs.Items) != 1 {
		t.Fatalf("expected 1 PK column, got %v", constr.PkAttrs)
	}
	pkCol := constr.PkAttrs.Items[0].(*nodes.String)
	if pkCol.Str != "id" {
		t.Errorf("expected PK column 'id', got %q", pkCol.Str)
	}
}

// TestCreateTableNamedConstraint tests: CREATE TABLE t (id integer CONSTRAINT pk_id PRIMARY KEY)
func TestCreateTableNamedConstraint(t *testing.T) {
	input := "CREATE TABLE t (id integer CONSTRAINT pk_id PRIMARY KEY)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	col := stmt.TableElts.Items[0].(*nodes.ColumnDef)

	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected 1 constraint, got %v", col.Constraints)
	}

	constr := col.Constraints.Items[0].(*nodes.Constraint)
	if constr.Contype != nodes.CONSTR_PRIMARY {
		t.Errorf("expected CONSTR_PRIMARY, got %d", constr.Contype)
	}
	if constr.Conname != "pk_id" {
		t.Errorf("expected constraint name 'pk_id', got %q", constr.Conname)
	}
}

// TestCreateTableMultipleConstraints tests: CREATE TABLE t (id integer NOT NULL PRIMARY KEY, name varchar(100) NOT NULL DEFAULT 'unknown')
func TestCreateTableMultipleConstraints(t *testing.T) {
	input := "CREATE TABLE t (id integer NOT NULL PRIMARY KEY, name varchar(100) NOT NULL DEFAULT 'unknown')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	if stmt.TableElts == nil || len(stmt.TableElts.Items) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(stmt.TableElts.Items))
	}

	// First column: id integer NOT NULL PRIMARY KEY
	col1 := stmt.TableElts.Items[0].(*nodes.ColumnDef)
	if col1.Colname != "id" {
		t.Errorf("expected 'id', got %q", col1.Colname)
	}
	if col1.Constraints == nil || len(col1.Constraints.Items) != 2 {
		t.Fatalf("expected 2 constraints on id, got %v", col1.Constraints)
	}
	c1 := col1.Constraints.Items[0].(*nodes.Constraint)
	if c1.Contype != nodes.CONSTR_NOTNULL {
		t.Errorf("expected CONSTR_NOTNULL, got %d", c1.Contype)
	}
	c2 := col1.Constraints.Items[1].(*nodes.Constraint)
	if c2.Contype != nodes.CONSTR_PRIMARY {
		t.Errorf("expected CONSTR_PRIMARY, got %d", c2.Contype)
	}

	// Second column: name varchar(100) NOT NULL DEFAULT 'unknown'
	col2 := stmt.TableElts.Items[1].(*nodes.ColumnDef)
	if col2.Colname != "name" {
		t.Errorf("expected 'name', got %q", col2.Colname)
	}
	if col2.Constraints == nil || len(col2.Constraints.Items) != 2 {
		t.Fatalf("expected 2 constraints on name, got %v", col2.Constraints)
	}
	c3 := col2.Constraints.Items[0].(*nodes.Constraint)
	if c3.Contype != nodes.CONSTR_NOTNULL {
		t.Errorf("expected CONSTR_NOTNULL, got %d", c3.Contype)
	}
	c4 := col2.Constraints.Items[1].(*nodes.Constraint)
	if c4.Contype != nodes.CONSTR_DEFAULT {
		t.Errorf("expected CONSTR_DEFAULT, got %d", c4.Contype)
	}
	if c4.RawExpr == nil {
		t.Fatal("expected RawExpr for DEFAULT constraint")
	}
	aConst, ok := c4.RawExpr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", c4.RawExpr)
	}
	strVal, ok := aConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", aConst.Val)
	}
	if strVal.Str != "unknown" {
		t.Errorf("expected 'unknown', got %q", strVal.Str)
	}
}

// TestCreateTableSchemaQualified tests: CREATE TABLE myschema.t (id integer)
func TestCreateTableSchemaQualified(t *testing.T) {
	input := "CREATE TABLE myschema.t (id integer)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
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

// TestCreateTempTable tests: CREATE TEMP TABLE t (id integer)
func TestCreateTempTable(t *testing.T) {
	input := "CREATE TEMP TABLE t (id integer)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateStmt, got %T", result.Items[0])
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Fatalf("expected table name 't'")
	}
}

// TestCreateTableTableDriven is a table-driven test for various CREATE TABLE patterns.
func TestCreateTableTableDriven(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic two columns", "CREATE TABLE t (id integer, name varchar(100))"},
		{"IF NOT EXISTS", "CREATE TABLE IF NOT EXISTS t (id integer)"},
		{"column NOT NULL", "CREATE TABLE t (id integer NOT NULL)"},
		{"column NULL", "CREATE TABLE t (id integer NULL)"},
		{"column DEFAULT", "CREATE TABLE t (id integer DEFAULT 0)"},
		{"column DEFAULT string", "CREATE TABLE t (name text DEFAULT 'hello')"},
		{"column PRIMARY KEY", "CREATE TABLE t (id integer PRIMARY KEY)"},
		{"column UNIQUE", "CREATE TABLE t (name varchar UNIQUE)"},
		{"column CHECK", "CREATE TABLE t (age integer CHECK (age > 0))"},
		{"column REFERENCES", "CREATE TABLE t (user_id integer REFERENCES users(id))"},
		{"column REFERENCES no columns", "CREATE TABLE t (user_id integer REFERENCES users)"},
		{"table PRIMARY KEY", "CREATE TABLE t (id integer, PRIMARY KEY (id))"},
		{"table UNIQUE", "CREATE TABLE t (id integer, UNIQUE (id))"},
		{"table CHECK", "CREATE TABLE t (id integer, CHECK (id > 0))"},
		{"table FOREIGN KEY", "CREATE TABLE t (id integer, FOREIGN KEY (id) REFERENCES t2(id))"},
		{"table FOREIGN KEY no ref cols", "CREATE TABLE t (id integer, FOREIGN KEY (id) REFERENCES t2)"},
		{"named column constraint", "CREATE TABLE t (id integer CONSTRAINT pk_id PRIMARY KEY)"},
		{"named table constraint", "CREATE TABLE t (id integer, CONSTRAINT pk_t PRIMARY KEY (id))"},
		{"multiple column constraints", "CREATE TABLE t (id integer NOT NULL PRIMARY KEY)"},
		{"multiple columns with constraints", "CREATE TABLE t (id integer NOT NULL PRIMARY KEY, name varchar(100) NOT NULL DEFAULT 'unknown')"},
		{"schema-qualified", "CREATE TABLE myschema.t (id integer)"},
		{"TEMP TABLE", "CREATE TEMP TABLE t (id integer)"},
		{"TEMPORARY TABLE", "CREATE TEMPORARY TABLE t (id integer)"},
		{"UNLOGGED TABLE", "CREATE UNLOGGED TABLE t (id integer)"},
		{"LOCAL TEMP TABLE", "CREATE LOCAL TEMP TABLE t (id integer)"},
		{"LOCAL TEMPORARY TABLE", "CREATE LOCAL TEMPORARY TABLE t (id integer)"},
		{"empty table", "CREATE TABLE t ()"},
		{"composite PK", "CREATE TABLE t (a integer, b integer, PRIMARY KEY (a, b))"},
		{"composite UNIQUE", "CREATE TABLE t (a integer, b integer, UNIQUE (a, b))"},
		{"composite FK", "CREATE TABLE t (a integer, b integer, FOREIGN KEY (a, b) REFERENCES t2(x, y))"},
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

			_, ok := result.Items[0].(*nodes.CreateStmt)
			if !ok {
				t.Fatalf("expected *nodes.CreateStmt for %q, got %T", tt.input, result.Items[0])
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}
