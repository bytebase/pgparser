package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE FOREIGN DATA WRAPPER tests
// =============================================================================

// TestCreateFdwBasic tests: CREATE FOREIGN DATA WRAPPER myfdw
func TestCreateFdwBasic(t *testing.T) {
	input := "CREATE FOREIGN DATA WRAPPER myfdw"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.CreateFdwStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateFdwStmt, got %T", result.Items[0])
	}
	if stmt.Fdwname != "myfdw" {
		t.Errorf("expected fdwname 'myfdw', got %q", stmt.Fdwname)
	}
	if stmt.FuncOptions != nil && stmt.FuncOptions.Len() > 0 {
		t.Error("expected no func options")
	}
	if stmt.Options != nil && stmt.Options.Len() > 0 {
		t.Error("expected no options")
	}
}

// TestCreateFdwHandlerValidator tests: CREATE FOREIGN DATA WRAPPER myfdw HANDLER fdw_handler VALIDATOR fdw_validator
func TestCreateFdwHandlerValidator(t *testing.T) {
	input := "CREATE FOREIGN DATA WRAPPER myfdw HANDLER fdw_handler VALIDATOR fdw_validator"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateFdwStmt)
	if stmt.Fdwname != "myfdw" {
		t.Errorf("expected fdwname 'myfdw', got %q", stmt.Fdwname)
	}
	if stmt.FuncOptions == nil || stmt.FuncOptions.Len() != 2 {
		t.Fatalf("expected 2 func options, got %d", stmt.FuncOptions.Len())
	}
	// Check HANDLER option
	handlerOpt := stmt.FuncOptions.Items[0].(*nodes.DefElem)
	if handlerOpt.Defname != "handler" {
		t.Errorf("expected defname 'handler', got %q", handlerOpt.Defname)
	}
	if handlerOpt.Arg == nil {
		t.Error("expected handler arg to be non-nil")
	}
	// Check VALIDATOR option
	validatorOpt := stmt.FuncOptions.Items[1].(*nodes.DefElem)
	if validatorOpt.Defname != "validator" {
		t.Errorf("expected defname 'validator', got %q", validatorOpt.Defname)
	}
}

// TestCreateFdwNoHandlerNoValidator tests: CREATE FOREIGN DATA WRAPPER myfdw NO HANDLER NO VALIDATOR
func TestCreateFdwNoHandlerNoValidator(t *testing.T) {
	input := "CREATE FOREIGN DATA WRAPPER myfdw NO HANDLER NO VALIDATOR"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateFdwStmt)
	if stmt.FuncOptions == nil || stmt.FuncOptions.Len() != 2 {
		t.Fatalf("expected 2 func options, got %d", stmt.FuncOptions.Len())
	}
	handlerOpt := stmt.FuncOptions.Items[0].(*nodes.DefElem)
	if handlerOpt.Defname != "handler" || handlerOpt.Arg != nil {
		t.Error("expected NO HANDLER (handler defelem with nil arg)")
	}
	validatorOpt := stmt.FuncOptions.Items[1].(*nodes.DefElem)
	if validatorOpt.Defname != "validator" || validatorOpt.Arg != nil {
		t.Error("expected NO VALIDATOR (validator defelem with nil arg)")
	}
}

// TestCreateFdwOptions tests: CREATE FOREIGN DATA WRAPPER myfdw OPTIONS (debug 'true')
func TestCreateFdwOptions(t *testing.T) {
	input := "CREATE FOREIGN DATA WRAPPER myfdw OPTIONS (debug 'true')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateFdwStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "debug" {
		t.Errorf("expected option name 'debug', got %q", opt.Defname)
	}
	val := opt.Arg.(*nodes.String)
	if val.Str != "true" {
		t.Errorf("expected option value 'true', got %q", val.Str)
	}
}

// =============================================================================
// ALTER FOREIGN DATA WRAPPER tests
// =============================================================================

// TestAlterFdwOptionsAdd tests: ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD debug 'true')
func TestAlterFdwOptionsAdd(t *testing.T) {
	input := "ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD debug 'true')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterFdwStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterFdwStmt, got %T", result.Items[0])
	}
	if stmt.Fdwname != "myfdw" {
		t.Errorf("expected fdwname 'myfdw', got %q", stmt.Fdwname)
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "debug" {
		t.Errorf("expected option name 'debug', got %q", opt.Defname)
	}
	if opt.Defaction != int(nodes.DEFELEM_ADD) {
		t.Errorf("expected DEFELEM_ADD (%d), got %d", nodes.DEFELEM_ADD, opt.Defaction)
	}
}

// TestAlterFdwHandler tests: ALTER FOREIGN DATA WRAPPER myfdw HANDLER new_handler
func TestAlterFdwHandler(t *testing.T) {
	input := "ALTER FOREIGN DATA WRAPPER myfdw HANDLER new_handler"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFdwStmt)
	if stmt.FuncOptions == nil || stmt.FuncOptions.Len() != 1 {
		t.Fatalf("expected 1 func option, got %d", stmt.FuncOptions.Len())
	}
	handlerOpt := stmt.FuncOptions.Items[0].(*nodes.DefElem)
	if handlerOpt.Defname != "handler" {
		t.Errorf("expected defname 'handler', got %q", handlerOpt.Defname)
	}
}

// TestAlterFdwOptionsSet tests: ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (SET debug 'false')
func TestAlterFdwOptionsSet(t *testing.T) {
	input := "ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (SET debug 'false')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFdwStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defaction != int(nodes.DEFELEM_SET) {
		t.Errorf("expected DEFELEM_SET (%d), got %d", nodes.DEFELEM_SET, opt.Defaction)
	}
}

// TestAlterFdwOptionsDrop tests: ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (DROP debug)
func TestAlterFdwOptionsDrop(t *testing.T) {
	input := "ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (DROP debug)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFdwStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "debug" {
		t.Errorf("expected option name 'debug', got %q", opt.Defname)
	}
	if opt.Defaction != int(nodes.DEFELEM_DROP) {
		t.Errorf("expected DEFELEM_DROP (%d), got %d", nodes.DEFELEM_DROP, opt.Defaction)
	}
	if opt.Arg != nil {
		t.Error("expected nil Arg for DROP")
	}
}

// =============================================================================
// CREATE SERVER tests
// =============================================================================

// TestCreateServerBasic tests: CREATE SERVER myserver FOREIGN DATA WRAPPER myfdw
func TestCreateServerBasic(t *testing.T) {
	input := "CREATE SERVER myserver FOREIGN DATA WRAPPER myfdw"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateForeignServerStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateForeignServerStmt, got %T", result.Items[0])
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.Fdwname != "myfdw" {
		t.Errorf("expected fdwname 'myfdw', got %q", stmt.Fdwname)
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
}

// TestCreateServerIfNotExists tests: CREATE SERVER IF NOT EXISTS myserver FOREIGN DATA WRAPPER myfdw
func TestCreateServerIfNotExists(t *testing.T) {
	input := "CREATE SERVER IF NOT EXISTS myserver FOREIGN DATA WRAPPER myfdw"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateForeignServerStmt)
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
}

// TestCreateServerFull tests: CREATE SERVER myserver TYPE 'postgres' VERSION '15' FOREIGN DATA WRAPPER myfdw OPTIONS (host 'localhost', port '5432')
func TestCreateServerFull(t *testing.T) {
	input := "CREATE SERVER myserver TYPE 'postgres' VERSION '15' FOREIGN DATA WRAPPER myfdw OPTIONS (host 'localhost', port '5432')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateForeignServerStmt)
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.Servertype != "postgres" {
		t.Errorf("expected servertype 'postgres', got %q", stmt.Servertype)
	}
	if stmt.Version != "15" {
		t.Errorf("expected version '15', got %q", stmt.Version)
	}
	if stmt.Fdwname != "myfdw" {
		t.Errorf("expected fdwname 'myfdw', got %q", stmt.Fdwname)
	}
	if stmt.Options == nil || stmt.Options.Len() != 2 {
		t.Fatalf("expected 2 options, got %d", stmt.Options.Len())
	}
	opt1 := stmt.Options.Items[0].(*nodes.DefElem)
	if opt1.Defname != "host" {
		t.Errorf("expected option name 'host', got %q", opt1.Defname)
	}
	opt2 := stmt.Options.Items[1].(*nodes.DefElem)
	if opt2.Defname != "port" {
		t.Errorf("expected option name 'port', got %q", opt2.Defname)
	}
}

// =============================================================================
// ALTER SERVER tests
// =============================================================================

// TestAlterServerVersion tests: ALTER SERVER myserver VERSION '16'
func TestAlterServerVersion(t *testing.T) {
	input := "ALTER SERVER myserver VERSION '16'"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterForeignServerStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterForeignServerStmt, got %T", result.Items[0])
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.Version != "16" {
		t.Errorf("expected version '16', got %q", stmt.Version)
	}
	if !stmt.HasVersion {
		t.Error("expected HasVersion to be true")
	}
}

// TestAlterServerOptions tests: ALTER SERVER myserver OPTIONS (SET host 'newhost')
func TestAlterServerOptions(t *testing.T) {
	input := "ALTER SERVER myserver OPTIONS (SET host 'newhost')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterForeignServerStmt)
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "host" {
		t.Errorf("expected option name 'host', got %q", opt.Defname)
	}
	if opt.Defaction != int(nodes.DEFELEM_SET) {
		t.Errorf("expected DEFELEM_SET, got %d", opt.Defaction)
	}
}

// TestAlterServerVersionAndOptions tests: ALTER SERVER myserver VERSION '16' OPTIONS (SET host 'newhost')
func TestAlterServerVersionAndOptions(t *testing.T) {
	input := "ALTER SERVER myserver VERSION '16' OPTIONS (SET host 'newhost')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterForeignServerStmt)
	if stmt.Version != "16" {
		t.Errorf("expected version '16', got %q", stmt.Version)
	}
	if !stmt.HasVersion {
		t.Error("expected HasVersion to be true")
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
}

// =============================================================================
// CREATE FOREIGN TABLE tests
// =============================================================================

// TestCreateForeignTableBasic tests: CREATE FOREIGN TABLE ft (id integer, name text) SERVER myserver
func TestCreateForeignTableBasic(t *testing.T) {
	input := "CREATE FOREIGN TABLE ft (id integer, name text) SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateForeignTableStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateForeignTableStmt, got %T", result.Items[0])
	}
	if stmt.Base.Relation == nil || stmt.Base.Relation.Relname != "ft" {
		t.Error("expected relation name 'ft'")
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.Base.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
	// Check columns
	if stmt.Base.TableElts == nil || stmt.Base.TableElts.Len() != 2 {
		t.Fatalf("expected 2 columns, got %d", stmt.Base.TableElts.Len())
	}
	col1 := stmt.Base.TableElts.Items[0].(*nodes.ColumnDef)
	if col1.Colname != "id" {
		t.Errorf("expected column name 'id', got %q", col1.Colname)
	}
	col2 := stmt.Base.TableElts.Items[1].(*nodes.ColumnDef)
	if col2.Colname != "name" {
		t.Errorf("expected column name 'name', got %q", col2.Colname)
	}
}

// TestCreateForeignTableIfNotExistsWithOptions tests: CREATE FOREIGN TABLE IF NOT EXISTS ft (id integer) SERVER myserver OPTIONS (table_name 'remote_table')
func TestCreateForeignTableIfNotExistsWithOptions(t *testing.T) {
	input := "CREATE FOREIGN TABLE IF NOT EXISTS ft (id integer) SERVER myserver OPTIONS (table_name 'remote_table')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateForeignTableStmt)
	if !stmt.Base.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "table_name" {
		t.Errorf("expected option name 'table_name', got %q", opt.Defname)
	}
}

// TestCreateForeignTableSchemaQualified tests: CREATE FOREIGN TABLE myschema.ft (id integer) SERVER myserver
func TestCreateForeignTableSchemaQualified(t *testing.T) {
	input := "CREATE FOREIGN TABLE myschema.ft (id integer) SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateForeignTableStmt)
	if stmt.Base.Relation.Schemaname != "myschema" {
		t.Errorf("expected schema 'myschema', got %q", stmt.Base.Relation.Schemaname)
	}
	if stmt.Base.Relation.Relname != "ft" {
		t.Errorf("expected relation 'ft', got %q", stmt.Base.Relation.Relname)
	}
}

// =============================================================================
// CREATE USER MAPPING tests
// =============================================================================

// TestCreateUserMappingCurrentUser tests: CREATE USER MAPPING FOR CURRENT_USER SERVER myserver OPTIONS (user 'me', password 'secret')
func TestCreateUserMappingCurrentUser(t *testing.T) {
	input := "CREATE USER MAPPING FOR CURRENT_USER SERVER myserver OPTIONS (user 'me', password 'secret')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateUserMappingStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateUserMappingStmt, got %T", result.Items[0])
	}
	if stmt.User == nil || stmt.User.Roletype != int(nodes.ROLESPEC_CURRENT_USER) {
		t.Error("expected ROLESPEC_CURRENT_USER")
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.Options == nil || stmt.Options.Len() != 2 {
		t.Fatalf("expected 2 options, got %d", stmt.Options.Len())
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
}

// TestCreateUserMappingIfNotExistsPublic tests: CREATE USER MAPPING IF NOT EXISTS FOR PUBLIC SERVER myserver
func TestCreateUserMappingIfNotExistsPublic(t *testing.T) {
	input := "CREATE USER MAPPING IF NOT EXISTS FOR PUBLIC SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateUserMappingStmt)
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	// PUBLIC is handled by RoleSpec
}

// TestCreateUserMappingNamedUser tests: CREATE USER MAPPING FOR myuser SERVER myserver
func TestCreateUserMappingNamedUser(t *testing.T) {
	input := "CREATE USER MAPPING FOR myuser SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateUserMappingStmt)
	if stmt.User == nil || stmt.User.Roletype != int(nodes.ROLESPEC_CSTRING) {
		t.Errorf("expected ROLESPEC_CSTRING, got %d", stmt.User.Roletype)
	}
	if stmt.User.Rolename != "myuser" {
		t.Errorf("expected rolename 'myuser', got %q", stmt.User.Rolename)
	}
}

// TestCreateUserMappingForUser tests: CREATE USER MAPPING FOR USER SERVER myserver
// USER is a reserved keyword requiring special handling via auth_ident
func TestCreateUserMappingForUser(t *testing.T) {
	input := "CREATE USER MAPPING FOR USER SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateUserMappingStmt)
	if stmt.User == nil || stmt.User.Roletype != int(nodes.ROLESPEC_CURRENT_USER) {
		t.Errorf("expected ROLESPEC_CURRENT_USER, got %d", stmt.User.Roletype)
	}
}

// =============================================================================
// ALTER USER MAPPING tests
// =============================================================================

// TestAlterUserMapping tests: ALTER USER MAPPING FOR CURRENT_USER SERVER myserver OPTIONS (SET password 'new')
func TestAlterUserMapping(t *testing.T) {
	input := "ALTER USER MAPPING FOR CURRENT_USER SERVER myserver OPTIONS (SET password 'new')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterUserMappingStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterUserMappingStmt, got %T", result.Items[0])
	}
	if stmt.User == nil || stmt.User.Roletype != int(nodes.ROLESPEC_CURRENT_USER) {
		t.Error("expected ROLESPEC_CURRENT_USER")
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defaction != int(nodes.DEFELEM_SET) {
		t.Errorf("expected DEFELEM_SET, got %d", opt.Defaction)
	}
}

// =============================================================================
// DROP USER MAPPING tests
// =============================================================================

// TestDropUserMapping tests: DROP USER MAPPING FOR CURRENT_USER SERVER myserver
func TestDropUserMapping(t *testing.T) {
	input := "DROP USER MAPPING FOR CURRENT_USER SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropUserMappingStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropUserMappingStmt, got %T", result.Items[0])
	}
	if stmt.User == nil || stmt.User.Roletype != int(nodes.ROLESPEC_CURRENT_USER) {
		t.Error("expected ROLESPEC_CURRENT_USER")
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
	if stmt.MissingOk {
		t.Error("expected MissingOk to be false")
	}
}

// TestDropUserMappingIfExists tests: DROP USER MAPPING IF EXISTS FOR PUBLIC SERVER myserver
func TestDropUserMappingIfExists(t *testing.T) {
	input := "DROP USER MAPPING IF EXISTS FOR PUBLIC SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropUserMappingStmt)
	if !stmt.MissingOk {
		t.Error("expected MissingOk to be true")
	}
}

// =============================================================================
// IMPORT FOREIGN SCHEMA tests
// =============================================================================

// TestImportForeignSchemaBasic tests: IMPORT FOREIGN SCHEMA remote_schema FROM SERVER myserver INTO local_schema
func TestImportForeignSchemaBasic(t *testing.T) {
	input := "IMPORT FOREIGN SCHEMA remote_schema FROM SERVER myserver INTO local_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.ImportForeignSchemaStmt)
	if !ok {
		t.Fatalf("expected *nodes.ImportForeignSchemaStmt, got %T", result.Items[0])
	}
	if stmt.RemoteSchema != "remote_schema" {
		t.Errorf("expected remote_schema 'remote_schema', got %q", stmt.RemoteSchema)
	}
	if stmt.ServerName != "myserver" {
		t.Errorf("expected server_name 'myserver', got %q", stmt.ServerName)
	}
	if stmt.LocalSchema != "local_schema" {
		t.Errorf("expected local_schema 'local_schema', got %q", stmt.LocalSchema)
	}
	if stmt.ListType != nodes.FDW_IMPORT_SCHEMA_ALL {
		t.Errorf("expected FDW_IMPORT_SCHEMA_ALL, got %d", stmt.ListType)
	}
}

// TestImportForeignSchemaLimitTo tests: IMPORT FOREIGN SCHEMA remote_schema LIMIT TO (t1, t2) FROM SERVER myserver INTO local_schema
func TestImportForeignSchemaLimitTo(t *testing.T) {
	input := "IMPORT FOREIGN SCHEMA remote_schema LIMIT TO (t1, t2) FROM SERVER myserver INTO local_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.ImportForeignSchemaStmt)
	if stmt.ListType != nodes.FDW_IMPORT_SCHEMA_LIMIT_TO {
		t.Errorf("expected FDW_IMPORT_SCHEMA_LIMIT_TO, got %d", stmt.ListType)
	}
	if stmt.TableList == nil || stmt.TableList.Len() != 2 {
		t.Fatalf("expected 2 tables in limit to, got %d", stmt.TableList.Len())
	}
}

// TestImportForeignSchemaExcept tests: IMPORT FOREIGN SCHEMA remote_schema EXCEPT (t3) FROM SERVER myserver INTO local_schema
func TestImportForeignSchemaExcept(t *testing.T) {
	input := "IMPORT FOREIGN SCHEMA remote_schema EXCEPT (t3) FROM SERVER myserver INTO local_schema"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.ImportForeignSchemaStmt)
	if stmt.ListType != nodes.FDW_IMPORT_SCHEMA_EXCEPT {
		t.Errorf("expected FDW_IMPORT_SCHEMA_EXCEPT, got %d", stmt.ListType)
	}
	if stmt.TableList == nil || stmt.TableList.Len() != 1 {
		t.Fatalf("expected 1 table in except, got %d", stmt.TableList.Len())
	}
}

// TestImportForeignSchemaWithOptions tests: IMPORT FOREIGN SCHEMA remote_schema FROM SERVER myserver INTO local_schema OPTIONS (import_default 'true')
func TestImportForeignSchemaWithOptions(t *testing.T) {
	input := "IMPORT FOREIGN SCHEMA remote_schema FROM SERVER myserver INTO local_schema OPTIONS (import_default 'true')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.ImportForeignSchemaStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "import_default" {
		t.Errorf("expected option name 'import_default', got %q", opt.Defname)
	}
}

// =============================================================================
// Additional edge case tests
// =============================================================================

// TestCreateFdwHandlerOnly tests: CREATE FOREIGN DATA WRAPPER myfdw HANDLER fdw_handler
func TestCreateFdwHandlerOnly(t *testing.T) {
	input := "CREATE FOREIGN DATA WRAPPER myfdw HANDLER fdw_handler"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateFdwStmt)
	if stmt.FuncOptions == nil || stmt.FuncOptions.Len() != 1 {
		t.Fatalf("expected 1 func option, got %d", stmt.FuncOptions.Len())
	}
	handlerOpt := stmt.FuncOptions.Items[0].(*nodes.DefElem)
	if handlerOpt.Defname != "handler" {
		t.Errorf("expected defname 'handler', got %q", handlerOpt.Defname)
	}
}

// TestCreateFdwSchemaQualifiedHandler tests: CREATE FOREIGN DATA WRAPPER myfdw HANDLER myschema.fdw_handler
func TestCreateFdwSchemaQualifiedHandler(t *testing.T) {
	input := "CREATE FOREIGN DATA WRAPPER myfdw HANDLER myschema.fdw_handler"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateFdwStmt)
	if stmt.FuncOptions == nil || stmt.FuncOptions.Len() != 1 {
		t.Fatalf("expected 1 func option, got %d", stmt.FuncOptions.Len())
	}
	handlerOpt := stmt.FuncOptions.Items[0].(*nodes.DefElem)
	if handlerOpt.Defname != "handler" {
		t.Errorf("expected defname 'handler', got %q", handlerOpt.Defname)
	}
	// The arg should be a List with two string nodes (schema-qualified name)
	handlerName, ok := handlerOpt.Arg.(*nodes.List)
	if !ok {
		t.Fatalf("expected *nodes.List for handler name, got %T", handlerOpt.Arg)
	}
	if handlerName.Len() != 2 {
		t.Errorf("expected 2-part handler name, got %d parts", handlerName.Len())
	}
}

// TestAlterFdwMultipleAlterOptions tests: ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD key1 'val1', SET key2 'val2', DROP key3)
func TestAlterFdwMultipleAlterOptions(t *testing.T) {
	input := "ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD key1 'val1', SET key2 'val2', DROP key3)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterFdwStmt)
	if stmt.Options == nil || stmt.Options.Len() != 3 {
		t.Fatalf("expected 3 options, got %d", stmt.Options.Len())
	}
	opt1 := stmt.Options.Items[0].(*nodes.DefElem)
	if opt1.Defname != "key1" || opt1.Defaction != int(nodes.DEFELEM_ADD) {
		t.Errorf("expected ADD key1, got action=%d name=%q", opt1.Defaction, opt1.Defname)
	}
	opt2 := stmt.Options.Items[1].(*nodes.DefElem)
	if opt2.Defname != "key2" || opt2.Defaction != int(nodes.DEFELEM_SET) {
		t.Errorf("expected SET key2, got action=%d name=%q", opt2.Defaction, opt2.Defname)
	}
	opt3 := stmt.Options.Items[2].(*nodes.DefElem)
	if opt3.Defname != "key3" || opt3.Defaction != int(nodes.DEFELEM_DROP) {
		t.Errorf("expected DROP key3, got action=%d name=%q", opt3.Defaction, opt3.Defname)
	}
}

// TestCreateForeignTableEmpty tests: CREATE FOREIGN TABLE ft () SERVER myserver
func TestCreateForeignTableEmpty(t *testing.T) {
	input := "CREATE FOREIGN TABLE ft () SERVER myserver"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateForeignTableStmt)
	if stmt.Base.Relation.Relname != "ft" {
		t.Errorf("expected relation 'ft', got %q", stmt.Base.Relation.Relname)
	}
	if stmt.Servername != "myserver" {
		t.Errorf("expected servername 'myserver', got %q", stmt.Servername)
	}
}
