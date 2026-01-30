package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// DROP OWNED BY tests
// =============================================================================

func TestDropOwned(t *testing.T) {
	input := "DROP OWNED BY myrole"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropOwnedStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropOwnedStmt, got %T", result.Items[0])
	}
	if stmt.Roles == nil || len(stmt.Roles.Items) != 1 {
		t.Fatal("expected 1 role")
	}
	if stmt.Behavior != nodes.DROP_RESTRICT {
		t.Errorf("expected DROP_RESTRICT, got %d", stmt.Behavior)
	}
}

func TestDropOwnedCascade(t *testing.T) {
	input := "DROP OWNED BY myrole CASCADE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropOwnedStmt)
	if stmt.Behavior != nodes.DROP_CASCADE {
		t.Errorf("expected DROP_CASCADE, got %d", stmt.Behavior)
	}
}

func TestDropOwnedMultipleRoles(t *testing.T) {
	input := "DROP OWNED BY role1, role2"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropOwnedStmt)
	if stmt.Roles == nil || len(stmt.Roles.Items) != 2 {
		t.Fatal("expected 2 roles")
	}
}

// =============================================================================
// REASSIGN OWNED BY tests
// =============================================================================

func TestReassignOwned(t *testing.T) {
	input := "REASSIGN OWNED BY myrole TO newrole"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.ReassignOwnedStmt)
	if !ok {
		t.Fatalf("expected *nodes.ReassignOwnedStmt, got %T", result.Items[0])
	}
	if stmt.Roles == nil || len(stmt.Roles.Items) != 1 {
		t.Fatal("expected 1 role")
	}
	if stmt.Newrole == nil || stmt.Newrole.Rolename != "newrole" {
		t.Error("expected newrole 'newrole'")
	}
}

func TestReassignOwnedMultipleRoles(t *testing.T) {
	input := "REASSIGN OWNED BY role1, role2 TO newrole"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.ReassignOwnedStmt)
	if stmt.Roles == nil || len(stmt.Roles.Items) != 2 {
		t.Fatal("expected 2 roles")
	}
}
