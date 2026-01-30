package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper to parse and extract a single CreateRoleStmt
func parseCreateRoleStmt(t *testing.T, input string) *nodes.CreateRoleStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.CreateRoleStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateRoleStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single AlterRoleStmt
func parseAlterRoleStmt(t *testing.T, input string) *nodes.AlterRoleStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.AlterRoleStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterRoleStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single AlterRoleSetStmt
func parseAlterRoleSetStmt(t *testing.T, input string) *nodes.AlterRoleSetStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.AlterRoleSetStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterRoleSetStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single DropRoleStmt
func parseDropRoleStmt(t *testing.T, input string) *nodes.DropRoleStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.DropRoleStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropRoleStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to parse and extract a single GrantRoleStmt
func parseGrantRoleStmt(t *testing.T, input string) *nodes.GrantRoleStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.GrantRoleStmt)
	if !ok {
		t.Fatalf("expected *nodes.GrantRoleStmt, got %T", result.Items[0])
	}
	return stmt
}

// helper to extract a DefElem from a list at a given index
func getDefElem(t *testing.T, list *nodes.List, index int) *nodes.DefElem {
	t.Helper()
	if list == nil || index >= len(list.Items) {
		t.Fatalf("expected at least %d DefElem(s), got %v", index+1, list)
	}
	de, ok := list.Items[index].(*nodes.DefElem)
	if !ok {
		t.Fatalf("expected *nodes.DefElem at index %d, got %T", index, list.Items[index])
	}
	return de
}

// helper to extract a RoleSpec from a list at a given index
func getRoleSpec(t *testing.T, list *nodes.List, index int) *nodes.RoleSpec {
	t.Helper()
	if list == nil || index >= len(list.Items) {
		t.Fatalf("expected at least %d RoleSpec(s), got %v", index+1, list)
	}
	rs, ok := list.Items[index].(*nodes.RoleSpec)
	if !ok {
		t.Fatalf("expected *nodes.RoleSpec at index %d, got %T", index, list.Items[index])
	}
	return rs
}

/*****************************************************************************
 * CREATE ROLE tests
 *****************************************************************************/

func TestCreateRoleBasic(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole")
	if stmt.StmtType != nodes.ROLESTMT_ROLE {
		t.Errorf("expected ROLESTMT_ROLE, got %d", stmt.StmtType)
	}
	if stmt.Role != "myrole" {
		t.Errorf("expected role 'myrole', got %q", stmt.Role)
	}
	if stmt.Options != nil && len(stmt.Options.Items) > 0 {
		t.Errorf("expected no options, got %v", stmt.Options)
	}
}

func TestCreateRoleWithLogin(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole WITH LOGIN")
	if stmt.Role != "myrole" {
		t.Errorf("expected role 'myrole', got %q", stmt.Role)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "canlogin" {
		t.Errorf("expected option 'canlogin', got %q", de.Defname)
	}
	bval, ok := de.Arg.(*nodes.Boolean)
	if !ok || !bval.Boolval {
		t.Errorf("expected canlogin=true, got %v", de.Arg)
	}
}

func TestCreateRoleSuperuserCreatedb(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole SUPERUSER CREATEDB")
	if stmt.Options == nil || len(stmt.Options.Items) != 2 {
		t.Fatalf("expected 2 options, got %v", stmt.Options)
	}
	de0 := getDefElem(t, stmt.Options, 0)
	if de0.Defname != "superuser" {
		t.Errorf("expected option 'superuser', got %q", de0.Defname)
	}
	de1 := getDefElem(t, stmt.Options, 1)
	if de1.Defname != "createdb" {
		t.Errorf("expected option 'createdb', got %q", de1.Defname)
	}
}

func TestCreateRolePassword(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole PASSWORD 'secret'")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "password" {
		t.Errorf("expected option 'password', got %q", de.Defname)
	}
	sval, ok := de.Arg.(*nodes.String)
	if !ok || sval.Str != "secret" {
		t.Errorf("expected password='secret', got %v", de.Arg)
	}
}

func TestCreateRoleConnectionLimit(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole CONNECTION LIMIT 5")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "connectionlimit" {
		t.Errorf("expected option 'connectionlimit', got %q", de.Defname)
	}
	ival, ok := de.Arg.(*nodes.Integer)
	if !ok || ival.Ival != 5 {
		t.Errorf("expected connectionlimit=5, got %v", de.Arg)
	}
}

func TestCreateRoleValidUntil(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole VALID UNTIL '2025-12-31'")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "validUntil" {
		t.Errorf("expected option 'validUntil', got %q", de.Defname)
	}
	sval, ok := de.Arg.(*nodes.String)
	if !ok || sval.Str != "2025-12-31" {
		t.Errorf("expected validUntil='2025-12-31', got %v", de.Arg)
	}
}

func TestCreateRoleNosuperuserNologin(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole NOSUPERUSER NOLOGIN")
	if stmt.Options == nil || len(stmt.Options.Items) != 2 {
		t.Fatalf("expected 2 options, got %v", stmt.Options)
	}
	de0 := getDefElem(t, stmt.Options, 0)
	if de0.Defname != "superuser" {
		t.Errorf("expected option 'superuser', got %q", de0.Defname)
	}
	bval0, ok := de0.Arg.(*nodes.Boolean)
	if !ok || bval0.Boolval {
		t.Errorf("expected superuser=false, got %v", de0.Arg)
	}
	de1 := getDefElem(t, stmt.Options, 1)
	if de1.Defname != "canlogin" {
		t.Errorf("expected option 'canlogin', got %q", de1.Defname)
	}
	bval1, ok := de1.Arg.(*nodes.Boolean)
	if !ok || bval1.Boolval {
		t.Errorf("expected canlogin=false, got %v", de1.Arg)
	}
}

func TestCreateRoleInRole(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole IN ROLE parent_role")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "addroleto" {
		t.Errorf("expected option 'addroleto', got %q", de.Defname)
	}
	roleList, ok := de.Arg.(*nodes.List)
	if !ok || len(roleList.Items) != 1 {
		t.Fatalf("expected 1 role in addroleto list, got %v", de.Arg)
	}
	rs := getRoleSpec(t, roleList, 0)
	if rs.Rolename != "parent_role" {
		t.Errorf("expected 'parent_role', got %q", rs.Rolename)
	}
}

func TestCreateRoleAdmin(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE ROLE myrole ADMIN admin_role")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "adminmembers" {
		t.Errorf("expected option 'adminmembers', got %q", de.Defname)
	}
}

/*****************************************************************************
 * CREATE USER tests
 *****************************************************************************/

func TestCreateUserBasic(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE USER myuser")
	if stmt.StmtType != nodes.ROLESTMT_USER {
		t.Errorf("expected ROLESTMT_USER, got %d", stmt.StmtType)
	}
	if stmt.Role != "myuser" {
		t.Errorf("expected role 'myuser', got %q", stmt.Role)
	}
}

func TestCreateUserWithPassword(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE USER myuser WITH PASSWORD 'secret'")
	if stmt.StmtType != nodes.ROLESTMT_USER {
		t.Errorf("expected ROLESTMT_USER, got %d", stmt.StmtType)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "password" {
		t.Errorf("expected option 'password', got %q", de.Defname)
	}
}

/*****************************************************************************
 * CREATE GROUP tests
 *****************************************************************************/

func TestCreateGroupBasic(t *testing.T) {
	stmt := parseCreateRoleStmt(t, "CREATE GROUP mygroup")
	if stmt.StmtType != nodes.ROLESTMT_GROUP {
		t.Errorf("expected ROLESTMT_GROUP, got %d", stmt.StmtType)
	}
	if stmt.Role != "mygroup" {
		t.Errorf("expected role 'mygroup', got %q", stmt.Role)
	}
}

/*****************************************************************************
 * ALTER ROLE tests
 *****************************************************************************/

func TestAlterRoleSuperuser(t *testing.T) {
	stmt := parseAlterRoleStmt(t, "ALTER ROLE myrole SUPERUSER")
	if stmt.Role == nil || stmt.Role.Rolename != "myrole" {
		t.Errorf("expected role 'myrole', got %v", stmt.Role)
	}
	if stmt.Action != 1 {
		t.Errorf("expected Action=1, got %d", stmt.Action)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "superuser" {
		t.Errorf("expected option 'superuser', got %q", de.Defname)
	}
}

func TestAlterRoleNosuperuser(t *testing.T) {
	stmt := parseAlterRoleStmt(t, "ALTER ROLE myrole NOSUPERUSER")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "superuser" {
		t.Errorf("expected option 'superuser', got %q", de.Defname)
	}
	bval, ok := de.Arg.(*nodes.Boolean)
	if !ok || bval.Boolval {
		t.Errorf("expected superuser=false, got %v", de.Arg)
	}
}

func TestAlterRolePassword(t *testing.T) {
	stmt := parseAlterRoleStmt(t, "ALTER ROLE myrole PASSWORD 'newpass'")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "password" {
		t.Errorf("expected option 'password', got %q", de.Defname)
	}
	sval, ok := de.Arg.(*nodes.String)
	if !ok || sval.Str != "newpass" {
		t.Errorf("expected password='newpass', got %v", de.Arg)
	}
}

func TestAlterRoleWithPassword(t *testing.T) {
	stmt := parseAlterRoleStmt(t, "ALTER ROLE myrole WITH PASSWORD 'newpass'")
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "password" {
		t.Errorf("expected option 'password', got %q", de.Defname)
	}
}

/*****************************************************************************
 * ALTER ROLE SET/RESET tests
 *****************************************************************************/

func TestAlterRoleSetSearchPath(t *testing.T) {
	stmt := parseAlterRoleSetStmt(t, "ALTER ROLE myrole SET search_path TO public")
	if stmt.Role == nil || stmt.Role.Rolename != "myrole" {
		t.Errorf("expected role 'myrole', got %v", stmt.Role)
	}
	if stmt.Database != "" {
		t.Errorf("expected empty database, got %q", stmt.Database)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected VAR_SET_VALUE, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "search_path" {
		t.Errorf("expected name 'search_path', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterRoleResetSearchPath(t *testing.T) {
	stmt := parseAlterRoleSetStmt(t, "ALTER ROLE myrole RESET search_path")
	if stmt.Role == nil || stmt.Role.Rolename != "myrole" {
		t.Errorf("expected role 'myrole', got %v", stmt.Role)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_RESET {
		t.Errorf("expected VAR_RESET, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "search_path" {
		t.Errorf("expected name 'search_path', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterRoleResetAll(t *testing.T) {
	stmt := parseAlterRoleSetStmt(t, "ALTER ROLE myrole RESET ALL")
	if stmt.Role == nil || stmt.Role.Rolename != "myrole" {
		t.Errorf("expected role 'myrole', got %v", stmt.Role)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_RESET_ALL {
		t.Errorf("expected VAR_RESET_ALL, got %d", stmt.Setstmt.Kind)
	}
}

func TestAlterRoleInDatabaseSet(t *testing.T) {
	stmt := parseAlterRoleSetStmt(t, "ALTER ROLE myrole IN DATABASE mydb SET search_path TO public")
	if stmt.Role == nil || stmt.Role.Rolename != "myrole" {
		t.Errorf("expected role 'myrole', got %v", stmt.Role)
	}
	if stmt.Database != "mydb" {
		t.Errorf("expected database 'mydb', got %q", stmt.Database)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Kind != nodes.VAR_SET_VALUE {
		t.Errorf("expected VAR_SET_VALUE, got %d", stmt.Setstmt.Kind)
	}
	if stmt.Setstmt.Name != "search_path" {
		t.Errorf("expected name 'search_path', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterRoleAllSetTimezone(t *testing.T) {
	stmt := parseAlterRoleSetStmt(t, "ALTER ROLE ALL SET timezone TO 'UTC'")
	if stmt.Role != nil {
		t.Errorf("expected nil role for ALL, got %v", stmt.Role)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Name != "timezone" {
		t.Errorf("expected name 'timezone', got %q", stmt.Setstmt.Name)
	}
}

func TestAlterRoleAllInDatabaseSetTimezone(t *testing.T) {
	stmt := parseAlterRoleSetStmt(t, "ALTER ROLE ALL IN DATABASE mydb SET timezone TO 'UTC'")
	if stmt.Role != nil {
		t.Errorf("expected nil role for ALL, got %v", stmt.Role)
	}
	if stmt.Database != "mydb" {
		t.Errorf("expected database 'mydb', got %q", stmt.Database)
	}
	if stmt.Setstmt == nil {
		t.Fatal("expected non-nil Setstmt")
	}
	if stmt.Setstmt.Name != "timezone" {
		t.Errorf("expected name 'timezone', got %q", stmt.Setstmt.Name)
	}
}

/*****************************************************************************
 * DROP ROLE tests
 *****************************************************************************/

func TestDropRole(t *testing.T) {
	stmt := parseDropRoleStmt(t, "DROP ROLE myrole")
	if stmt.MissingOk {
		t.Error("expected MissingOk=false")
	}
	if stmt.Roles == nil || len(stmt.Roles.Items) != 1 {
		t.Fatalf("expected 1 role, got %v", stmt.Roles)
	}
	rs := getRoleSpec(t, stmt.Roles, 0)
	if rs.Rolename != "myrole" {
		t.Errorf("expected 'myrole', got %q", rs.Rolename)
	}
}

func TestDropRoleIfExists(t *testing.T) {
	stmt := parseDropRoleStmt(t, "DROP ROLE IF EXISTS myrole")
	if !stmt.MissingOk {
		t.Error("expected MissingOk=true")
	}
	if stmt.Roles == nil || len(stmt.Roles.Items) != 1 {
		t.Fatalf("expected 1 role, got %v", stmt.Roles)
	}
}

func TestDropRoleMultiple(t *testing.T) {
	stmt := parseDropRoleStmt(t, "DROP ROLE myrole, myrole2")
	if stmt.MissingOk {
		t.Error("expected MissingOk=false")
	}
	if stmt.Roles == nil || len(stmt.Roles.Items) != 2 {
		t.Fatalf("expected 2 roles, got %v", stmt.Roles)
	}
	rs0 := getRoleSpec(t, stmt.Roles, 0)
	if rs0.Rolename != "myrole" {
		t.Errorf("expected 'myrole', got %q", rs0.Rolename)
	}
	rs1 := getRoleSpec(t, stmt.Roles, 1)
	if rs1.Rolename != "myrole2" {
		t.Errorf("expected 'myrole2', got %q", rs1.Rolename)
	}
}

func TestDropUser(t *testing.T) {
	stmt := parseDropRoleStmt(t, "DROP USER myuser")
	if stmt.MissingOk {
		t.Error("expected MissingOk=false")
	}
	if stmt.Roles == nil || len(stmt.Roles.Items) != 1 {
		t.Fatalf("expected 1 role, got %v", stmt.Roles)
	}
	rs := getRoleSpec(t, stmt.Roles, 0)
	if rs.Rolename != "myuser" {
		t.Errorf("expected 'myuser', got %q", rs.Rolename)
	}
}

func TestDropGroup(t *testing.T) {
	stmt := parseDropRoleStmt(t, "DROP GROUP mygroup")
	if stmt.MissingOk {
		t.Error("expected MissingOk=false")
	}
	if stmt.Roles == nil || len(stmt.Roles.Items) != 1 {
		t.Fatalf("expected 1 role, got %v", stmt.Roles)
	}
	rs := getRoleSpec(t, stmt.Roles, 0)
	if rs.Rolename != "mygroup" {
		t.Errorf("expected 'mygroup', got %q", rs.Rolename)
	}
}

/*****************************************************************************
 * GRANT role TO tests
 *****************************************************************************/

func TestGrantRoleToUser(t *testing.T) {
	stmt := parseGrantRoleStmt(t, "GRANT myrole TO myuser")
	if !stmt.IsGrant {
		t.Error("expected IsGrant=true")
	}
	if stmt.GrantedRoles == nil || len(stmt.GrantedRoles.Items) != 1 {
		t.Fatalf("expected 1 granted role, got %v", stmt.GrantedRoles)
	}
	if stmt.GranteeRoles == nil || len(stmt.GranteeRoles.Items) != 1 {
		t.Fatalf("expected 1 grantee role, got %v", stmt.GranteeRoles)
	}
	// Granted role is an AccessPriv (from privilege_list)
	grantedPriv, ok := stmt.GrantedRoles.Items[0].(*nodes.AccessPriv)
	if !ok {
		t.Fatalf("expected *nodes.AccessPriv for granted role, got %T", stmt.GrantedRoles.Items[0])
	}
	if grantedPriv.PrivName != "myrole" {
		t.Errorf("expected granted role 'myrole', got %q", grantedPriv.PrivName)
	}
	// Grantee is a RoleSpec
	grantee := getRoleSpec(t, stmt.GranteeRoles, 0)
	if grantee.Rolename != "myuser" {
		t.Errorf("expected grantee 'myuser', got %q", grantee.Rolename)
	}
}

func TestGrantRoleWithAdminOption(t *testing.T) {
	stmt := parseGrantRoleStmt(t, "GRANT myrole TO myuser WITH ADMIN OPTION")
	if !stmt.IsGrant {
		t.Error("expected IsGrant=true")
	}
	if stmt.Opt == nil || len(stmt.Opt.Items) != 1 {
		t.Fatalf("expected 1 grant option, got %v", stmt.Opt)
	}
	de := getDefElem(t, stmt.Opt, 0)
	if de.Defname != "admin" {
		t.Errorf("expected option name 'admin', got %q", de.Defname)
	}
	bval, ok := de.Arg.(*nodes.Boolean)
	if !ok || !bval.Boolval {
		t.Errorf("expected admin=true, got %v", de.Arg)
	}
}

func TestGrantMultipleRoles(t *testing.T) {
	stmt := parseGrantRoleStmt(t, "GRANT role1, role2 TO user1, user2")
	if !stmt.IsGrant {
		t.Error("expected IsGrant=true")
	}
	if stmt.GrantedRoles == nil || len(stmt.GrantedRoles.Items) != 2 {
		t.Fatalf("expected 2 granted roles, got %v", stmt.GrantedRoles)
	}
	if stmt.GranteeRoles == nil || len(stmt.GranteeRoles.Items) != 2 {
		t.Fatalf("expected 2 grantee roles, got %v", stmt.GranteeRoles)
	}
}

/*****************************************************************************
 * REVOKE role FROM tests
 *****************************************************************************/

func TestRevokeRoleFromUser(t *testing.T) {
	stmt := parseGrantRoleStmt(t, "REVOKE myrole FROM myuser")
	if stmt.IsGrant {
		t.Error("expected IsGrant=false")
	}
	if stmt.GrantedRoles == nil || len(stmt.GrantedRoles.Items) != 1 {
		t.Fatalf("expected 1 revoked role, got %v", stmt.GrantedRoles)
	}
	grantedPriv, ok := stmt.GrantedRoles.Items[0].(*nodes.AccessPriv)
	if !ok {
		t.Fatalf("expected *nodes.AccessPriv for revoked role, got %T", stmt.GrantedRoles.Items[0])
	}
	if grantedPriv.PrivName != "myrole" {
		t.Errorf("expected revoked role 'myrole', got %q", grantedPriv.PrivName)
	}
	grantee := getRoleSpec(t, stmt.GranteeRoles, 0)
	if grantee.Rolename != "myuser" {
		t.Errorf("expected grantee 'myuser', got %q", grantee.Rolename)
	}
}

func TestRevokeAdminOptionFor(t *testing.T) {
	stmt := parseGrantRoleStmt(t, "REVOKE ADMIN OPTION FOR myrole FROM myuser")
	if stmt.IsGrant {
		t.Error("expected IsGrant=false")
	}
	if stmt.Opt == nil || len(stmt.Opt.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Opt)
	}
	de := getDefElem(t, stmt.Opt, 0)
	if de.Defname != "admin" {
		t.Errorf("expected option name 'admin', got %q", de.Defname)
	}
	bval, ok := de.Arg.(*nodes.Boolean)
	if !ok || bval.Boolval {
		t.Errorf("expected admin=false (revoke), got %v", de.Arg)
	}
}

/*****************************************************************************
 * ALTER GROUP tests
 *****************************************************************************/

func TestAlterGroupAddUser(t *testing.T) {
	stmt := parseAlterRoleStmt(t, "ALTER GROUP mygroup ADD USER myuser")
	if stmt.Role == nil || stmt.Role.Rolename != "mygroup" {
		t.Errorf("expected role 'mygroup', got %v", stmt.Role)
	}
	if stmt.Action != 1 {
		t.Errorf("expected Action=1 (add), got %d", stmt.Action)
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	de := getDefElem(t, stmt.Options, 0)
	if de.Defname != "rolemembers" {
		t.Errorf("expected option 'rolemembers', got %q", de.Defname)
	}
}

func TestAlterGroupDropUser(t *testing.T) {
	stmt := parseAlterRoleStmt(t, "ALTER GROUP mygroup DROP USER myuser")
	if stmt.Role == nil || stmt.Role.Rolename != "mygroup" {
		t.Errorf("expected role 'mygroup', got %v", stmt.Role)
	}
	if stmt.Action != -1 {
		t.Errorf("expected Action=-1 (drop), got %d", stmt.Action)
	}
}

/*****************************************************************************
 * Table-driven parse-only tests
 *****************************************************************************/

func TestRoleManagementTableDriven(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// CREATE ROLE
		{"create role basic", "CREATE ROLE myrole"},
		{"create role with login", "CREATE ROLE myrole WITH LOGIN"},
		{"create role superuser createdb", "CREATE ROLE myrole SUPERUSER CREATEDB"},
		{"create role password", "CREATE ROLE myrole PASSWORD 'secret'"},
		{"create role connection limit", "CREATE ROLE myrole CONNECTION LIMIT 5"},
		{"create role valid until", "CREATE ROLE myrole VALID UNTIL '2025-12-31'"},
		{"create role nosuperuser nologin", "CREATE ROLE myrole NOSUPERUSER NOLOGIN"},
		{"create role in role", "CREATE ROLE myrole IN ROLE parent_role"},
		{"create role admin", "CREATE ROLE myrole ADMIN admin_role"},
		{"create role encrypted password", "CREATE ROLE myrole ENCRYPTED PASSWORD 'secret'"},
		{"create role inherit", "CREATE ROLE myrole INHERIT"},
		{"create role noinherit", "CREATE ROLE myrole NOINHERIT"},
		{"create role createrole", "CREATE ROLE myrole CREATEROLE"},
		{"create role nocreaterole", "CREATE ROLE myrole NOCREATEROLE"},
		{"create role replication", "CREATE ROLE myrole REPLICATION"},
		{"create role noreplication", "CREATE ROLE myrole NOREPLICATION"},
		{"create role bypassrls", "CREATE ROLE myrole BYPASSRLS"},
		{"create role nobypassrls", "CREATE ROLE myrole NOBYPASSRLS"},
		{"create role password null", "CREATE ROLE myrole PASSWORD NULL"},
		{"create role sysid", "CREATE ROLE myrole SYSID 100"},
		{"create role in group", "CREATE ROLE myrole IN GROUP parent_group"},
		{"create role role list", "CREATE ROLE myrole ROLE member1, member2"},
		{"create role multiple options", "CREATE ROLE myrole WITH SUPERUSER CREATEDB LOGIN PASSWORD 'secret'"},

		// CREATE USER
		{"create user basic", "CREATE USER myuser"},
		{"create user with password", "CREATE USER myuser WITH PASSWORD 'secret'"},
		{"create user login", "CREATE USER myuser LOGIN"},

		// CREATE GROUP
		{"create group basic", "CREATE GROUP mygroup"},
		{"create group with user", "CREATE GROUP mygroup USER member1"},

		// ALTER ROLE
		{"alter role superuser", "ALTER ROLE myrole SUPERUSER"},
		{"alter role nosuperuser", "ALTER ROLE myrole NOSUPERUSER"},
		{"alter role password", "ALTER ROLE myrole PASSWORD 'newpass'"},
		{"alter role with password", "ALTER ROLE myrole WITH PASSWORD 'newpass'"},
		{"alter user superuser", "ALTER USER myuser SUPERUSER"},

		// ALTER ROLE SET/RESET
		{"alter role set", "ALTER ROLE myrole SET search_path TO public"},
		{"alter role reset", "ALTER ROLE myrole RESET search_path"},
		{"alter role reset all", "ALTER ROLE myrole RESET ALL"},
		{"alter role in database set", "ALTER ROLE myrole IN DATABASE mydb SET search_path TO public"},
		{"alter role all set", "ALTER ROLE ALL SET timezone TO 'UTC'"},
		{"alter role all in database set", "ALTER ROLE ALL IN DATABASE mydb SET timezone TO 'UTC'"},
		{"alter user set", "ALTER USER myuser SET search_path TO public"},
		{"alter user all set", "ALTER USER ALL SET timezone TO 'UTC'"},

		// DROP ROLE/USER/GROUP
		{"drop role", "DROP ROLE myrole"},
		{"drop role if exists", "DROP ROLE IF EXISTS myrole"},
		{"drop role multiple", "DROP ROLE myrole, myrole2"},
		{"drop user", "DROP USER myuser"},
		{"drop user if exists", "DROP USER IF EXISTS myuser"},
		{"drop group", "DROP GROUP mygroup"},
		{"drop group if exists", "DROP GROUP IF EXISTS mygroup"},

		// GRANT role TO
		{"grant role to user", "GRANT myrole TO myuser"},
		{"grant role with admin option", "GRANT myrole TO myuser WITH ADMIN OPTION"},
		{"grant multiple roles", "GRANT role1, role2 TO user1, user2"},

		// REVOKE role FROM
		{"revoke role from user", "REVOKE myrole FROM myuser"},
		{"revoke admin option for", "REVOKE ADMIN OPTION FOR myrole FROM myuser"},

		// ALTER GROUP
		{"alter group add user", "ALTER GROUP mygroup ADD USER myuser"},
		{"alter group drop user", "ALTER GROUP mygroup DROP USER myuser"},
		{"alter group add multiple users", "ALTER GROUP mygroup ADD USER user1, user2"},
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
			t.Logf("OK: %s -> %T", tt.input, result.Items[0])
		})
	}
}
