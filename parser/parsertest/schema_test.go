package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestCreateSchemaBasic tests: CREATE SCHEMA myschema
func TestCreateSchemaBasic(t *testing.T) {
	input := "CREATE SCHEMA myschema"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSchemaStmt, got %T", result.Items[0])
	}

	if stmt.Schemaname != "myschema" {
		t.Errorf("expected schema name 'myschema', got %q", stmt.Schemaname)
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
	if stmt.Authrole != nil {
		t.Error("expected Authrole to be nil")
	}
}

// TestCreateSchemaIfNotExists tests: CREATE SCHEMA IF NOT EXISTS myschema
func TestCreateSchemaIfNotExists(t *testing.T) {
	input := "CREATE SCHEMA IF NOT EXISTS myschema"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSchemaStmt, got %T", result.Items[0])
	}

	if stmt.Schemaname != "myschema" {
		t.Errorf("expected schema name 'myschema', got %q", stmt.Schemaname)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
}

// TestCreateSchemaAuthorization tests: CREATE SCHEMA AUTHORIZATION myrole
func TestCreateSchemaAuthorization(t *testing.T) {
	input := "CREATE SCHEMA AUTHORIZATION myrole"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSchemaStmt, got %T", result.Items[0])
	}

	if stmt.Authrole == nil {
		t.Fatal("expected Authrole to be set")
	}
	if stmt.Authrole.Rolename != "myrole" {
		t.Errorf("expected role 'myrole', got %q", stmt.Authrole.Rolename)
	}
}

// TestCreateSchemaWithAuthorization tests: CREATE SCHEMA myschema AUTHORIZATION myrole
func TestCreateSchemaWithAuthorization(t *testing.T) {
	input := "CREATE SCHEMA myschema AUTHORIZATION myrole"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSchemaStmt, got %T", result.Items[0])
	}

	if stmt.Schemaname != "myschema" {
		t.Errorf("expected schema name 'myschema', got %q", stmt.Schemaname)
	}
	if stmt.Authrole == nil {
		t.Fatal("expected Authrole to be set")
	}
	if stmt.Authrole.Rolename != "myrole" {
		t.Errorf("expected role 'myrole', got %q", stmt.Authrole.Rolename)
	}
}

// TestCreateSchemaIfNotExistsAuth tests: CREATE SCHEMA IF NOT EXISTS AUTHORIZATION myrole
func TestCreateSchemaIfNotExistsAuth(t *testing.T) {
	input := "CREATE SCHEMA IF NOT EXISTS AUTHORIZATION myrole"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSchemaStmt, got %T", result.Items[0])
	}

	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	if stmt.Authrole == nil {
		t.Fatal("expected Authrole to be set")
	}
}

// TestCreateSchemaWithTable tests: CREATE SCHEMA myschema CREATE TABLE t(id int)
func TestCreateSchemaWithTable(t *testing.T) {
	input := "CREATE SCHEMA myschema CREATE TABLE t(id int)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSchemaStmt, got %T", result.Items[0])
	}

	if stmt.Schemaname != "myschema" {
		t.Errorf("expected schema name 'myschema', got %q", stmt.Schemaname)
	}
	if stmt.SchemaElts == nil || len(stmt.SchemaElts.Items) != 1 {
		t.Fatalf("expected 1 schema element, got %v", stmt.SchemaElts)
	}
	if _, ok := stmt.SchemaElts.Items[0].(*nodes.CreateStmt); !ok {
		t.Fatalf("expected *nodes.CreateStmt, got %T", stmt.SchemaElts.Items[0])
	}
}

// TestCreateSchemaWithMultipleElements tests: CREATE SCHEMA myschema CREATE TABLE t(id int) CREATE VIEW v AS SELECT 1
func TestCreateSchemaWithMultipleElements(t *testing.T) {
	input := "CREATE SCHEMA myschema CREATE TABLE t(id int) CREATE VIEW v AS SELECT 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSchemaStmt, got %T", result.Items[0])
	}

	if stmt.SchemaElts == nil || len(stmt.SchemaElts.Items) != 2 {
		t.Fatalf("expected 2 schema elements, got %d", stmt.SchemaElts.Len())
	}
	if _, ok := stmt.SchemaElts.Items[0].(*nodes.CreateStmt); !ok {
		t.Fatalf("expected *nodes.CreateStmt for first element, got %T", stmt.SchemaElts.Items[0])
	}
	if _, ok := stmt.SchemaElts.Items[1].(*nodes.ViewStmt); !ok {
		t.Fatalf("expected *nodes.ViewStmt for second element, got %T", stmt.SchemaElts.Items[1])
	}
}
