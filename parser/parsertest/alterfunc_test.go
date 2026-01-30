package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// ALTER FUNCTION tests
// =============================================================================

// TestAlterFunctionImmutable tests: ALTER FUNCTION func(integer) IMMUTABLE
func TestAlterFunctionImmutable(t *testing.T) {
	input := "ALTER FUNCTION func(integer) IMMUTABLE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.AlterFunctionStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterFunctionStmt, got %T", result.Items[0])
	}
	if stmt.Objtype != nodes.OBJECT_FUNCTION {
		t.Errorf("expected OBJECT_FUNCTION, got %d", stmt.Objtype)
	}
	if stmt.Func == nil {
		t.Fatal("expected Func to be set")
	}
	if stmt.Func.Objname == nil || len(stmt.Func.Objname.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.Func.Objname.Items[0].(*nodes.String).Str != "func" {
		t.Errorf("expected 'func', got %q", stmt.Func.Objname.Items[0].(*nodes.String).Str)
	}
	if stmt.Actions == nil || len(stmt.Actions.Items) != 1 {
		t.Fatal("expected 1 action")
	}
	de := stmt.Actions.Items[0].(*nodes.DefElem)
	if de.Defname != "volatility" {
		t.Errorf("expected defname 'volatility', got %q", de.Defname)
	}
}

// TestAlterFunctionStableStrict tests: ALTER FUNCTION func(integer) STABLE STRICT
func TestAlterFunctionStableStrict(t *testing.T) {
	input := "ALTER FUNCTION func(integer) STABLE STRICT"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFunctionStmt)
	if stmt.Actions == nil || len(stmt.Actions.Items) != 2 {
		t.Fatalf("expected 2 actions, got %d", stmt.Actions.Len())
	}
}

// TestAlterFunctionSecurityDefiner tests: ALTER FUNCTION func(integer) SECURITY DEFINER
func TestAlterFunctionSecurityDefiner(t *testing.T) {
	input := "ALTER FUNCTION func(integer) SECURITY DEFINER"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFunctionStmt)
	de := stmt.Actions.Items[0].(*nodes.DefElem)
	if de.Defname != "security" {
		t.Errorf("expected defname 'security', got %q", de.Defname)
	}
}

// TestAlterFunctionCost tests: ALTER FUNCTION func(integer) COST 100
func TestAlterFunctionCost(t *testing.T) {
	input := "ALTER FUNCTION func(integer) COST 100"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFunctionStmt)
	de := stmt.Actions.Items[0].(*nodes.DefElem)
	if de.Defname != "cost" {
		t.Errorf("expected defname 'cost', got %q", de.Defname)
	}
}

// TestAlterFunctionRows tests: ALTER FUNCTION func(integer) ROWS 1000
func TestAlterFunctionRows(t *testing.T) {
	input := "ALTER FUNCTION func(integer) ROWS 1000"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFunctionStmt)
	de := stmt.Actions.Items[0].(*nodes.DefElem)
	if de.Defname != "rows" {
		t.Errorf("expected defname 'rows', got %q", de.Defname)
	}
}

// TestAlterProcedureSecurityInvoker tests: ALTER PROCEDURE proc(integer) SECURITY INVOKER
func TestAlterProcedureSecurityInvoker(t *testing.T) {
	input := "ALTER PROCEDURE proc(integer) SECURITY INVOKER"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFunctionStmt)
	if stmt.Objtype != nodes.OBJECT_PROCEDURE {
		t.Errorf("expected OBJECT_PROCEDURE, got %d", stmt.Objtype)
	}
}

// =============================================================================
// DROP FUNCTION/PROCEDURE tests
// =============================================================================

// TestDropFunction tests: DROP FUNCTION func(integer)
func TestDropFunction(t *testing.T) {
	input := "DROP FUNCTION func(integer)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	if stmt.RemoveType != int(nodes.OBJECT_FUNCTION) {
		t.Errorf("expected OBJECT_FUNCTION (%d), got %d", nodes.OBJECT_FUNCTION, stmt.RemoveType)
	}
	if stmt.Missing_ok {
		t.Error("expected Missing_ok to be false")
	}
}

// TestDropFunctionIfExists tests: DROP FUNCTION IF EXISTS func(integer)
func TestDropFunctionIfExists(t *testing.T) {
	input := "DROP FUNCTION IF EXISTS func(integer)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
}

// TestDropFunctionCascade tests: DROP FUNCTION func(integer) CASCADE
func TestDropFunctionCascade(t *testing.T) {
	input := "DROP FUNCTION func(integer) CASCADE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if stmt.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected DROP_CASCADE (%d), got %d", nodes.DROP_CASCADE, stmt.Behavior)
	}
}

// TestDropFunctionTwoArgs tests: DROP FUNCTION func(integer, text)
func TestDropFunctionTwoArgs(t *testing.T) {
	input := "DROP FUNCTION func(integer, text)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if stmt.RemoveType != int(nodes.OBJECT_FUNCTION) {
		t.Errorf("expected OBJECT_FUNCTION, got %d", stmt.RemoveType)
	}
	if stmt.Objects == nil || len(stmt.Objects.Items) != 1 {
		t.Fatal("expected 1 object in drop list")
	}
	owa := stmt.Objects.Items[0].(*nodes.ObjectWithArgs)
	if owa.Objargs == nil || len(owa.Objargs.Items) != 2 {
		t.Fatalf("expected 2 args, got %v", owa.Objargs)
	}
}

// TestDropProcedure tests: DROP PROCEDURE proc(integer)
func TestDropProcedure(t *testing.T) {
	input := "DROP PROCEDURE proc(integer)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if stmt.RemoveType != int(nodes.OBJECT_PROCEDURE) {
		t.Errorf("expected OBJECT_PROCEDURE (%d), got %d", nodes.OBJECT_PROCEDURE, stmt.RemoveType)
	}
}

// =============================================================================
// DROP AGGREGATE tests
// =============================================================================

// TestDropAggregate tests: DROP AGGREGATE myagg(integer)
func TestDropAggregate(t *testing.T) {
	input := "DROP AGGREGATE myagg(integer)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	if stmt.RemoveType != int(nodes.OBJECT_AGGREGATE) {
		t.Errorf("expected OBJECT_AGGREGATE (%d), got %d", nodes.OBJECT_AGGREGATE, stmt.RemoveType)
	}
}

// TestDropAggregateIfExistsCascade tests: DROP AGGREGATE IF EXISTS myagg(integer) CASCADE
func TestDropAggregateIfExistsCascade(t *testing.T) {
	input := "DROP AGGREGATE IF EXISTS myagg(integer) CASCADE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if stmt.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected DROP_CASCADE, got %d", stmt.Behavior)
	}
}

// =============================================================================
// DROP OPERATOR tests
// =============================================================================

// TestDropOperator tests: DROP OPERATOR +(integer, integer)
func TestDropOperator(t *testing.T) {
	input := "DROP OPERATOR +(integer, integer)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	if stmt.RemoveType != int(nodes.OBJECT_OPERATOR) {
		t.Errorf("expected OBJECT_OPERATOR (%d), got %d", nodes.OBJECT_OPERATOR, stmt.RemoveType)
	}
}

// TestDropOperatorIfExistsCascade tests: DROP OPERATOR IF EXISTS =(text, text) CASCADE
func TestDropOperatorIfExistsCascade(t *testing.T) {
	input := "DROP OPERATOR IF EXISTS =(text, text) CASCADE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropStmt)
	if !stmt.Missing_ok {
		t.Error("expected Missing_ok to be true")
	}
	if stmt.Behavior != int(nodes.DROP_CASCADE) {
		t.Errorf("expected DROP_CASCADE, got %d", stmt.Behavior)
	}
	if stmt.RemoveType != int(nodes.OBJECT_OPERATOR) {
		t.Errorf("expected OBJECT_OPERATOR, got %d", stmt.RemoveType)
	}
}
