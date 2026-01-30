package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single GrantStmt
func parseGrantStmt(t *testing.T, input string) *nodes.GrantStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.GrantStmt)
	if !ok {
		t.Fatalf("expected *nodes.GrantStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to extract a RoleSpec from grantees list at a given index
func getGrantee(t *testing.T, stmt *nodes.GrantStmt, index int) *nodes.RoleSpec {
	t.Helper()
	if stmt.Grantees == nil || index >= len(stmt.Grantees.Items) {
		t.Fatalf("expected at least %d grantee(s), got %v", index+1, stmt.Grantees)
	}
	role, ok := stmt.Grantees.Items[index].(*nodes.RoleSpec)
	if !ok {
		t.Fatalf("expected *nodes.RoleSpec at index %d, got %T", index, stmt.Grantees.Items[index])
	}
	return role
}

// helper to extract an AccessPriv from privileges list at a given index
func getPrivilege(t *testing.T, stmt *nodes.GrantStmt, index int) *nodes.AccessPriv {
	t.Helper()
	if stmt.Privileges == nil || index >= len(stmt.Privileges.Items) {
		t.Fatalf("expected at least %d privilege(s), got %v", index+1, stmt.Privileges)
	}
	priv, ok := stmt.Privileges.Items[index].(*nodes.AccessPriv)
	if !ok {
		t.Fatalf("expected *nodes.AccessPriv at index %d, got %T", index, stmt.Privileges.Items[index])
	}
	return priv
}

// helper to extract a RangeVar from Objects list at a given index
func getGrantObject(t *testing.T, stmt *nodes.GrantStmt, index int) *nodes.RangeVar {
	t.Helper()
	if stmt.Objects == nil || index >= len(stmt.Objects.Items) {
		t.Fatalf("expected at least %d object(s), got %v", index+1, stmt.Objects)
	}
	rv, ok := stmt.Objects.Items[index].(*nodes.RangeVar)
	if !ok {
		t.Fatalf("expected *nodes.RangeVar at index %d, got %T", index, stmt.Objects.Items[index])
	}
	return rv
}

// TestGrantSelectOnTable tests: GRANT SELECT ON TABLE foo TO user1
func TestGrantSelectOnTable(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT ON TABLE foo TO user1")

	if !stmt.IsGrant {
		t.Error("expected IsGrant to be true")
	}
	if stmt.Targtype != nodes.ACL_TARGET_OBJECT {
		t.Errorf("expected Targtype ACL_TARGET_OBJECT (%d), got %d", nodes.ACL_TARGET_OBJECT, stmt.Targtype)
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected Objtype OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.Objtype)
	}
	if stmt.GrantOption {
		t.Error("expected GrantOption to be false")
	}

	// Check privilege: single SELECT
	if stmt.Privileges == nil || len(stmt.Privileges.Items) != 1 {
		t.Fatalf("expected 1 privilege, got %v", stmt.Privileges)
	}
	priv := getPrivilege(t, stmt, 0)
	if priv.PrivName != "select" {
		t.Errorf("expected privilege 'select', got %q", priv.PrivName)
	}

	// Check object: foo
	if stmt.Objects == nil || len(stmt.Objects.Items) != 1 {
		t.Fatalf("expected 1 object, got %v", stmt.Objects)
	}
	obj := getGrantObject(t, stmt, 0)
	if obj.Relname != "foo" {
		t.Errorf("expected object name 'foo', got %q", obj.Relname)
	}

	// Check grantee: user1
	if stmt.Grantees == nil || len(stmt.Grantees.Items) != 1 {
		t.Fatalf("expected 1 grantee, got %v", stmt.Grantees)
	}
	grantee := getGrantee(t, stmt, 0)
	if grantee.Roletype != int(nodes.ROLESPEC_CSTRING) {
		t.Errorf("expected ROLESPEC_CSTRING (%d), got %d", nodes.ROLESPEC_CSTRING, grantee.Roletype)
	}
	if grantee.Rolename != "user1" {
		t.Errorf("expected grantee 'user1', got %q", grantee.Rolename)
	}
}

// TestGrantAllOnTable tests: GRANT ALL ON TABLE foo TO user1
func TestGrantAllOnTable(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT ALL ON TABLE foo TO user1")

	if !stmt.IsGrant {
		t.Error("expected IsGrant to be true")
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected Objtype OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.Objtype)
	}

	// ALL PRIVILEGES => nil privileges list
	if stmt.Privileges != nil {
		t.Errorf("expected nil Privileges for ALL, got %v", stmt.Privileges)
	}

	grantee := getGrantee(t, stmt, 0)
	if grantee.Rolename != "user1" {
		t.Errorf("expected grantee 'user1', got %q", grantee.Rolename)
	}
}

// TestGrantAllPrivilegesOnTable tests: GRANT ALL PRIVILEGES ON TABLE foo TO user1
func TestGrantAllPrivilegesOnTable(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT ALL PRIVILEGES ON TABLE foo TO user1")

	if !stmt.IsGrant {
		t.Error("expected IsGrant to be true")
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected Objtype OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.Objtype)
	}

	// ALL PRIVILEGES => nil privileges list
	if stmt.Privileges != nil {
		t.Errorf("expected nil Privileges for ALL PRIVILEGES, got %v", stmt.Privileges)
	}
}

// TestGrantMultiplePrivileges tests: GRANT SELECT, INSERT ON TABLE foo TO user1, user2
func TestGrantMultiplePrivileges(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT, INSERT ON TABLE foo TO user1, user2")

	if !stmt.IsGrant {
		t.Error("expected IsGrant to be true")
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected Objtype OBJECT_TABLE, got %d", stmt.Objtype)
	}

	// Check privileges: SELECT, INSERT
	if stmt.Privileges == nil || len(stmt.Privileges.Items) != 2 {
		t.Fatalf("expected 2 privileges, got %v", stmt.Privileges)
	}
	priv0 := getPrivilege(t, stmt, 0)
	if priv0.PrivName != "select" {
		t.Errorf("expected first privilege 'select', got %q", priv0.PrivName)
	}
	priv1 := getPrivilege(t, stmt, 1)
	if priv1.PrivName != "insert" {
		t.Errorf("expected second privilege 'insert', got %q", priv1.PrivName)
	}

	// Check grantees: user1, user2
	if stmt.Grantees == nil || len(stmt.Grantees.Items) != 2 {
		t.Fatalf("expected 2 grantees, got %v", stmt.Grantees)
	}
	g0 := getGrantee(t, stmt, 0)
	if g0.Rolename != "user1" {
		t.Errorf("expected first grantee 'user1', got %q", g0.Rolename)
	}
	g1 := getGrantee(t, stmt, 1)
	if g1.Rolename != "user2" {
		t.Errorf("expected second grantee 'user2', got %q", g1.Rolename)
	}
}

// TestGrantSelectDefaultTable tests: GRANT SELECT ON foo TO user1 (default TABLE target)
func TestGrantSelectDefaultTable(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT ON foo TO user1")

	if !stmt.IsGrant {
		t.Error("expected IsGrant to be true")
	}
	// When TABLE keyword is omitted, Objtype should still be OBJECT_TABLE
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected Objtype OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.Objtype)
	}

	priv := getPrivilege(t, stmt, 0)
	if priv.PrivName != "select" {
		t.Errorf("expected privilege 'select', got %q", priv.PrivName)
	}

	obj := getGrantObject(t, stmt, 0)
	if obj.Relname != "foo" {
		t.Errorf("expected object name 'foo', got %q", obj.Relname)
	}

	grantee := getGrantee(t, stmt, 0)
	if grantee.Rolename != "user1" {
		t.Errorf("expected grantee 'user1', got %q", grantee.Rolename)
	}
}

// TestGrantToCurrentUser tests: GRANT SELECT ON TABLE foo TO CURRENT_USER
func TestGrantToCurrentUser(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT ON TABLE foo TO CURRENT_USER")

	if !stmt.IsGrant {
		t.Error("expected IsGrant to be true")
	}

	grantee := getGrantee(t, stmt, 0)
	if grantee.Roletype != int(nodes.ROLESPEC_CURRENT_USER) {
		t.Errorf("expected ROLESPEC_CURRENT_USER (%d), got %d", nodes.ROLESPEC_CURRENT_USER, grantee.Roletype)
	}
	if grantee.Rolename != "" {
		t.Errorf("expected empty Rolename for CURRENT_USER, got %q", grantee.Rolename)
	}
}

// TestGrantWithGrantOption tests: GRANT SELECT ON TABLE foo TO user1 WITH GRANT OPTION
func TestGrantWithGrantOption(t *testing.T) {
	stmt := parseGrantStmt(t, "GRANT SELECT ON TABLE foo TO user1 WITH GRANT OPTION")

	if !stmt.IsGrant {
		t.Error("expected IsGrant to be true")
	}
	if !stmt.GrantOption {
		t.Error("expected GrantOption to be true")
	}

	priv := getPrivilege(t, stmt, 0)
	if priv.PrivName != "select" {
		t.Errorf("expected privilege 'select', got %q", priv.PrivName)
	}

	obj := getGrantObject(t, stmt, 0)
	if obj.Relname != "foo" {
		t.Errorf("expected object name 'foo', got %q", obj.Relname)
	}

	grantee := getGrantee(t, stmt, 0)
	if grantee.Rolename != "user1" {
		t.Errorf("expected grantee 'user1', got %q", grantee.Rolename)
	}
}

// TestRevokeSelect tests: REVOKE SELECT ON TABLE foo FROM user1
func TestRevokeSelect(t *testing.T) {
	stmt := parseGrantStmt(t, "REVOKE SELECT ON TABLE foo FROM user1")

	if stmt.IsGrant {
		t.Error("expected IsGrant to be false for REVOKE")
	}
	if stmt.Targtype != nodes.ACL_TARGET_OBJECT {
		t.Errorf("expected Targtype ACL_TARGET_OBJECT (%d), got %d", nodes.ACL_TARGET_OBJECT, stmt.Targtype)
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected Objtype OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.Objtype)
	}
	// Default behavior: RESTRICT
	if stmt.Behavior != nodes.DROP_RESTRICT {
		t.Errorf("expected Behavior DROP_RESTRICT (%d), got %d", nodes.DROP_RESTRICT, stmt.Behavior)
	}
	if stmt.GrantOption {
		t.Error("expected GrantOption to be false")
	}

	priv := getPrivilege(t, stmt, 0)
	if priv.PrivName != "select" {
		t.Errorf("expected privilege 'select', got %q", priv.PrivName)
	}

	obj := getGrantObject(t, stmt, 0)
	if obj.Relname != "foo" {
		t.Errorf("expected object name 'foo', got %q", obj.Relname)
	}

	grantee := getGrantee(t, stmt, 0)
	if grantee.Rolename != "user1" {
		t.Errorf("expected grantee 'user1', got %q", grantee.Rolename)
	}
}

// TestRevokeAllCascade tests: REVOKE ALL ON TABLE foo FROM user1 CASCADE
func TestRevokeAllCascade(t *testing.T) {
	stmt := parseGrantStmt(t, "REVOKE ALL ON TABLE foo FROM user1 CASCADE")

	if stmt.IsGrant {
		t.Error("expected IsGrant to be false for REVOKE")
	}
	if stmt.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected Objtype OBJECT_TABLE (%d), got %d", nodes.OBJECT_TABLE, stmt.Objtype)
	}

	// ALL => nil privileges
	if stmt.Privileges != nil {
		t.Errorf("expected nil Privileges for ALL, got %v", stmt.Privileges)
	}

	// CASCADE
	if stmt.Behavior != nodes.DROP_CASCADE {
		t.Errorf("expected Behavior DROP_CASCADE (%d), got %d", nodes.DROP_CASCADE, stmt.Behavior)
	}

	grantee := getGrantee(t, stmt, 0)
	if grantee.Rolename != "user1" {
		t.Errorf("expected grantee 'user1', got %q", grantee.Rolename)
	}
}

// TestRevokeSelectRestrict tests: REVOKE SELECT ON TABLE foo FROM user1 RESTRICT
func TestRevokeSelectRestrict(t *testing.T) {
	stmt := parseGrantStmt(t, "REVOKE SELECT ON TABLE foo FROM user1 RESTRICT")

	if stmt.IsGrant {
		t.Error("expected IsGrant to be false for REVOKE")
	}
	if stmt.Behavior != nodes.DROP_RESTRICT {
		t.Errorf("expected Behavior DROP_RESTRICT (%d), got %d", nodes.DROP_RESTRICT, stmt.Behavior)
	}

	priv := getPrivilege(t, stmt, 0)
	if priv.PrivName != "select" {
		t.Errorf("expected privilege 'select', got %q", priv.PrivName)
	}
}

// TestGrantRevokeTableDriven is a table-driven test for various GRANT/REVOKE patterns.
func TestGrantRevokeTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isGrant bool
	}{
		// GRANT variants
		{"grant select on table", "GRANT SELECT ON TABLE foo TO user1", true},
		{"grant all on table", "GRANT ALL ON TABLE foo TO user1", true},
		{"grant all privileges on table", "GRANT ALL PRIVILEGES ON TABLE foo TO user1", true},
		{"grant select insert on table", "GRANT SELECT, INSERT ON TABLE foo TO user1, user2", true},
		{"grant select default table", "GRANT SELECT ON foo TO user1", true},
		{"grant to current_user", "GRANT SELECT ON TABLE foo TO CURRENT_USER", true},
		{"grant to session_user", "GRANT SELECT ON TABLE foo TO SESSION_USER", true},
		{"grant with grant option", "GRANT SELECT ON TABLE foo TO user1 WITH GRANT OPTION", true},
		{"grant update on table", "GRANT UPDATE ON TABLE foo TO user1", true},
		{"grant delete on table", "GRANT DELETE ON TABLE foo TO user1", true},
		{"grant references on table", "GRANT REFERENCES ON TABLE foo TO user1", true},
		{"grant trigger on table", "GRANT TRIGGER ON TABLE foo TO user1", true},
		{"grant insert on table", "GRANT INSERT ON TABLE foo TO user1", true},

		// REVOKE variants
		{"revoke select on table", "REVOKE SELECT ON TABLE foo FROM user1", false},
		{"revoke all on table", "REVOKE ALL ON TABLE foo FROM user1", false},
		{"revoke all on table cascade", "REVOKE ALL ON TABLE foo FROM user1 CASCADE", false},
		{"revoke select on table restrict", "REVOKE SELECT ON TABLE foo FROM user1 RESTRICT", false},
		{"revoke all privileges on table", "REVOKE ALL PRIVILEGES ON TABLE foo FROM user1", false},
		{"revoke select default table", "REVOKE SELECT ON foo FROM user1", false},
		{"revoke from current_user", "REVOKE SELECT ON TABLE foo FROM CURRENT_USER", false},
		{"revoke from session_user", "REVOKE SELECT ON TABLE foo FROM SESSION_USER", false},
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

			stmt, ok := result.Items[0].(*nodes.GrantStmt)
			if !ok {
				t.Fatalf("expected *nodes.GrantStmt for %q, got %T", tt.input, result.Items[0])
			}

			if stmt.IsGrant != tt.isGrant {
				t.Errorf("expected IsGrant=%v for %q, got %v", tt.isGrant, tt.input, stmt.IsGrant)
			}

			t.Logf("OK: %s", tt.input)
		})
	}
}

// TestRevokeGrantOptionFor tests: REVOKE GRANT OPTION FOR SELECT ON TABLE foo FROM user1
// PostgreSQL supports revoking just the grant option without revoking the privilege itself.
func TestRevokeGrantOptionFor(t *testing.T) {
	stmt := parseGrantStmt(t, "REVOKE GRANT OPTION FOR SELECT ON TABLE foo FROM user1")

	if stmt.IsGrant {
		t.Error("expected IsGrant to be false for REVOKE")
	}
	if !stmt.GrantOption {
		t.Error("expected GrantOption to be true for REVOKE GRANT OPTION FOR")
	}

	priv := getPrivilege(t, stmt, 0)
	if priv.PrivName != "select" {
		t.Errorf("expected privilege 'select', got %q", priv.PrivName)
	}
}
