package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE POLICY tests
// =============================================================================

// TestCreatePolicyUsing tests: CREATE POLICY p ON t USING (true)
func TestCreatePolicyUsing(t *testing.T) {
	input := "CREATE POLICY p ON t USING (true)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.CreatePolicyStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreatePolicyStmt, got %T", result.Items[0])
	}
	if stmt.PolicyName != "p" {
		t.Errorf("expected policy name 'p', got %q", stmt.PolicyName)
	}
	if stmt.Table == nil || stmt.Table.Relname != "t" {
		t.Error("expected table 't'")
	}
	if stmt.CmdName != "all" {
		t.Errorf("expected cmd_name 'all', got %q", stmt.CmdName)
	}
	if !stmt.Permissive {
		t.Error("expected permissive to be true by default")
	}
	if stmt.Qual == nil {
		t.Error("expected USING qual to be non-nil")
	}
}

// TestCreatePolicyForSelectToPublic tests: CREATE POLICY p ON t FOR SELECT TO PUBLIC USING (true)
func TestCreatePolicyForSelectToPublic(t *testing.T) {
	input := "CREATE POLICY p ON t FOR SELECT TO PUBLIC USING (true)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePolicyStmt)
	if stmt.CmdName != "select" {
		t.Errorf("expected cmd_name 'select', got %q", stmt.CmdName)
	}
	if stmt.Roles == nil || stmt.Roles.Len() != 1 {
		t.Fatal("expected 1 role")
	}
}

// TestCreatePolicyPermissiveAll tests: CREATE POLICY p ON t AS PERMISSIVE FOR ALL TO role1 USING (col > 0)
func TestCreatePolicyPermissiveAll(t *testing.T) {
	input := "CREATE POLICY p ON t AS permissive FOR ALL TO role1 USING (col > 0)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePolicyStmt)
	if !stmt.Permissive {
		t.Error("expected permissive to be true")
	}
	if stmt.CmdName != "all" {
		t.Errorf("expected cmd_name 'all', got %q", stmt.CmdName)
	}
	if stmt.Roles == nil || stmt.Roles.Len() != 1 {
		t.Fatal("expected 1 role")
	}
}

// TestCreatePolicyInsertWithCheck tests: CREATE POLICY p ON t FOR INSERT WITH CHECK (col > 0)
func TestCreatePolicyInsertWithCheck(t *testing.T) {
	input := "CREATE POLICY p ON t FOR INSERT WITH CHECK (col > 0)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePolicyStmt)
	if stmt.CmdName != "insert" {
		t.Errorf("expected cmd_name 'insert', got %q", stmt.CmdName)
	}
	if stmt.WithCheck == nil {
		t.Error("expected WITH CHECK to be non-nil")
	}
}

// TestCreatePolicyRestrictive tests: CREATE POLICY p ON t AS RESTRICTIVE FOR UPDATE USING (true) WITH CHECK (true)
func TestCreatePolicyRestrictive(t *testing.T) {
	input := "CREATE POLICY p ON t AS restrictive FOR UPDATE USING (true) WITH CHECK (true)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePolicyStmt)
	if stmt.Permissive {
		t.Error("expected permissive to be false (restrictive)")
	}
	if stmt.CmdName != "update" {
		t.Errorf("expected cmd_name 'update', got %q", stmt.CmdName)
	}
	if stmt.Qual == nil {
		t.Error("expected USING qual to be non-nil")
	}
	if stmt.WithCheck == nil {
		t.Error("expected WITH CHECK to be non-nil")
	}
}

// TestCreatePolicyForDelete tests: CREATE POLICY p ON t FOR DELETE USING (true)
func TestCreatePolicyForDelete(t *testing.T) {
	input := "CREATE POLICY p ON t FOR DELETE USING (true)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePolicyStmt)
	if stmt.CmdName != "delete" {
		t.Errorf("expected cmd_name 'delete', got %q", stmt.CmdName)
	}
}

// =============================================================================
// ALTER POLICY tests
// =============================================================================

// TestAlterPolicyUsing tests: ALTER POLICY p ON t USING (col > 5)
func TestAlterPolicyUsing(t *testing.T) {
	input := "ALTER POLICY p ON t USING (col > 5)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterPolicyStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterPolicyStmt, got %T", result.Items[0])
	}
	if stmt.PolicyName != "p" {
		t.Errorf("expected policy name 'p', got %q", stmt.PolicyName)
	}
	if stmt.Table == nil || stmt.Table.Relname != "t" {
		t.Error("expected table 't'")
	}
	if stmt.Qual == nil {
		t.Error("expected USING qual to be non-nil")
	}
}

// TestAlterPolicyToRoles tests: ALTER POLICY p ON t TO role1, role2
func TestAlterPolicyToRoles(t *testing.T) {
	input := "ALTER POLICY p ON t TO role1, role2"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterPolicyStmt)
	if stmt.Roles == nil || stmt.Roles.Len() != 2 {
		t.Fatalf("expected 2 roles, got %d", stmt.Roles.Len())
	}
}

// TestAlterPolicyWithCheck tests: ALTER POLICY p ON t USING (true) WITH CHECK (col > 0)
func TestAlterPolicyWithCheck(t *testing.T) {
	input := "ALTER POLICY p ON t USING (true) WITH CHECK (col > 0)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterPolicyStmt)
	if stmt.Qual == nil {
		t.Error("expected USING qual to be non-nil")
	}
	if stmt.WithCheck == nil {
		t.Error("expected WITH CHECK to be non-nil")
	}
}
