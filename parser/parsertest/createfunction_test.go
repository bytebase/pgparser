package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// parseCreateFunctionStmt is a helper that parses input and returns the first statement as *nodes.CreateFunctionStmt.
func parseCreateFunctionStmt(t *testing.T, input string) *nodes.CreateFunctionStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement, got %v", result)
	}
	stmt, ok := result.Items[0].(*nodes.CreateFunctionStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateFunctionStmt, got %T", result.Items[0])
	}
	return stmt
}

// findDefElem finds a DefElem by name in the Options list.
func findDefElem(opts *nodes.List, name string) *nodes.DefElem {
	if opts == nil {
		return nil
	}
	for _, item := range opts.Items {
		de, ok := item.(*nodes.DefElem)
		if ok && de.Defname == name {
			return de
		}
	}
	return nil
}

// TestCreateFunctionBasic tests: CREATE FUNCTION foo() RETURNS void LANGUAGE sql AS 'SELECT 1'
func TestCreateFunctionBasic(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS void LANGUAGE sql AS 'SELECT 1'")

	if stmt.IsOrReplace {
		t.Error("expected IsOrReplace to be false")
	}
	if stmt.Funcname == nil || len(stmt.Funcname.Items) != 1 {
		t.Fatalf("expected 1 name component, got %v", stmt.Funcname)
	}
	fname := stmt.Funcname.Items[0].(*nodes.String).Str
	if fname != "foo" {
		t.Errorf("expected function name 'foo', got %q", fname)
	}
	if stmt.Parameters != nil && len(stmt.Parameters.Items) > 0 {
		t.Errorf("expected no parameters, got %d", len(stmt.Parameters.Items))
	}
	if stmt.ReturnType == nil {
		t.Fatal("expected ReturnType to be set")
	}
	langDef := findDefElem(stmt.Options, "language")
	if langDef == nil {
		t.Fatal("expected language option")
	}
	langStr, ok := langDef.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected language arg to be *nodes.String, got %T", langDef.Arg)
	}
	if langStr.Str != "sql" {
		t.Errorf("expected language 'sql', got %q", langStr.Str)
	}
	asDef := findDefElem(stmt.Options, "as")
	if asDef == nil {
		t.Fatal("expected as option")
	}
	asList, ok := asDef.Arg.(*nodes.List)
	if !ok {
		t.Fatalf("expected as arg to be *nodes.List, got %T", asDef.Arg)
	}
	if len(asList.Items) != 1 {
		t.Fatalf("expected 1 item in AS list, got %d", len(asList.Items))
	}
	asStr, ok := asList.Items[0].(*nodes.String)
	if !ok {
		t.Fatalf("expected AS list item to be *nodes.String, got %T", asList.Items[0])
	}
	if asStr.Str != "SELECT 1" {
		t.Errorf("expected as 'SELECT 1', got %q", asStr.Str)
	}
}

// TestCreateFunctionWithParams tests: CREATE FUNCTION add(a integer, b integer) RETURNS integer LANGUAGE sql AS 'SELECT a + b'
func TestCreateFunctionWithParams(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION add(a integer, b integer) RETURNS integer LANGUAGE sql AS 'SELECT a + b'")

	if stmt.Funcname == nil || len(stmt.Funcname.Items) != 1 {
		t.Fatalf("expected 1 name component, got %v", stmt.Funcname)
	}
	fname := stmt.Funcname.Items[0].(*nodes.String).Str
	if fname != "add" {
		t.Errorf("expected function name 'add', got %q", fname)
	}
	if stmt.Parameters == nil || len(stmt.Parameters.Items) != 2 {
		t.Fatalf("expected 2 parameters, got %v", stmt.Parameters)
	}

	// Check first parameter
	p1 := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	if p1.Name != "a" {
		t.Errorf("expected param 1 name 'a', got %q", p1.Name)
	}
	if p1.ArgType == nil {
		t.Fatal("expected param 1 type to be set")
	}
	if p1.Mode != nodes.FUNC_PARAM_IN {
		t.Errorf("expected param 1 mode IN (%d), got %d", nodes.FUNC_PARAM_IN, p1.Mode)
	}

	// Check second parameter
	p2 := stmt.Parameters.Items[1].(*nodes.FunctionParameter)
	if p2.Name != "b" {
		t.Errorf("expected param 2 name 'b', got %q", p2.Name)
	}
	if p2.ArgType == nil {
		t.Fatal("expected param 2 type to be set")
	}
	if p2.Mode != nodes.FUNC_PARAM_IN {
		t.Errorf("expected param 2 mode IN (%d), got %d", nodes.FUNC_PARAM_IN, p2.Mode)
	}

	if stmt.ReturnType == nil {
		t.Fatal("expected ReturnType to be set")
	}
}

// TestCreateFunctionUnnamedParam tests: CREATE FUNCTION foo(integer) RETURNS integer LANGUAGE plpgsql AS 'body'
func TestCreateFunctionUnnamedParam(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo(integer) RETURNS integer LANGUAGE plpgsql AS 'body'")

	if stmt.Parameters == nil || len(stmt.Parameters.Items) != 1 {
		t.Fatalf("expected 1 parameter, got %v", stmt.Parameters)
	}
	p1 := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	if p1.Name != "" {
		t.Errorf("expected empty param name, got %q", p1.Name)
	}
	if p1.ArgType == nil {
		t.Fatal("expected param type to be set")
	}
	if p1.Mode != nodes.FUNC_PARAM_IN {
		t.Errorf("expected param mode IN (%d), got %d", nodes.FUNC_PARAM_IN, p1.Mode)
	}
}

// TestCreateOrReplaceFunction tests: CREATE OR REPLACE FUNCTION foo() RETURNS void LANGUAGE sql AS 'SELECT 1'
func TestCreateOrReplaceFunction(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE OR REPLACE FUNCTION foo() RETURNS void LANGUAGE sql AS 'SELECT 1'")

	if !stmt.IsOrReplace {
		t.Error("expected IsOrReplace to be true")
	}
	if stmt.Funcname == nil || len(stmt.Funcname.Items) != 1 {
		t.Fatalf("expected 1 name component, got %v", stmt.Funcname)
	}
	fname := stmt.Funcname.Items[0].(*nodes.String).Str
	if fname != "foo" {
		t.Errorf("expected function name 'foo', got %q", fname)
	}
}

// TestCreateFunctionInOutParams tests: CREATE FUNCTION foo(IN a integer, OUT b integer) LANGUAGE sql AS 'body'
func TestCreateFunctionInOutParams(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo(IN a integer, OUT b integer) LANGUAGE sql AS 'body'")

	if stmt.Parameters == nil || len(stmt.Parameters.Items) != 2 {
		t.Fatalf("expected 2 parameters, got %v", stmt.Parameters)
	}

	p1 := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	if p1.Name != "a" {
		t.Errorf("expected param 1 name 'a', got %q", p1.Name)
	}
	if p1.Mode != nodes.FUNC_PARAM_IN {
		t.Errorf("expected param 1 mode IN (%d), got %d", nodes.FUNC_PARAM_IN, p1.Mode)
	}

	p2 := stmt.Parameters.Items[1].(*nodes.FunctionParameter)
	if p2.Name != "b" {
		t.Errorf("expected param 2 name 'b', got %q", p2.Name)
	}
	if p2.Mode != nodes.FUNC_PARAM_OUT {
		t.Errorf("expected param 2 mode OUT (%d), got %d", nodes.FUNC_PARAM_OUT, p2.Mode)
	}
}

// TestCreateFunctionDefaultParam tests: CREATE FUNCTION foo(x integer DEFAULT 0) RETURNS integer LANGUAGE sql AS 'SELECT x'
func TestCreateFunctionDefaultParam(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo(x integer DEFAULT 0) RETURNS integer LANGUAGE sql AS 'SELECT x'")

	if stmt.Parameters == nil || len(stmt.Parameters.Items) != 1 {
		t.Fatalf("expected 1 parameter, got %v", stmt.Parameters)
	}

	p1 := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	if p1.Name != "x" {
		t.Errorf("expected param name 'x', got %q", p1.Name)
	}
	if p1.Defexpr == nil {
		t.Error("expected default expression to be set")
	}
}

// TestCreateFunctionImmutable tests: CREATE FUNCTION foo() RETURNS integer LANGUAGE sql IMMUTABLE AS 'SELECT 1'
func TestCreateFunctionImmutable(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS integer LANGUAGE sql IMMUTABLE AS 'SELECT 1'")

	volDef := findDefElem(stmt.Options, "volatility")
	if volDef == nil {
		t.Fatal("expected volatility option")
	}
	volStr, ok := volDef.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected volatility arg to be *nodes.String, got %T", volDef.Arg)
	}
	if volStr.Str != "immutable" {
		t.Errorf("expected volatility 'immutable', got %q", volStr.Str)
	}
}

// TestCreateFunctionStrict tests: CREATE FUNCTION foo() RETURNS integer LANGUAGE sql STRICT AS 'SELECT 1'
func TestCreateFunctionStrict(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS integer LANGUAGE sql STRICT AS 'SELECT 1'")

	strictDef := findDefElem(stmt.Options, "strict")
	if strictDef == nil {
		t.Fatal("expected strict option")
	}
	strictVal, ok := strictDef.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected strict arg to be *nodes.Integer, got %T", strictDef.Arg)
	}
	if strictVal.Ival != 1 {
		t.Errorf("expected strict value 1, got %d", strictVal.Ival)
	}
}

// TestCreateFunctionSecurityDefiner tests: CREATE FUNCTION foo() RETURNS integer SECURITY DEFINER LANGUAGE sql AS 'SELECT 1'
func TestCreateFunctionSecurityDefiner(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS integer SECURITY DEFINER LANGUAGE sql AS 'SELECT 1'")

	secDef := findDefElem(stmt.Options, "security")
	if secDef == nil {
		t.Fatal("expected security option")
	}
	secVal, ok := secDef.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected security arg to be *nodes.Integer, got %T", secDef.Arg)
	}
	if secVal.Ival != 1 {
		t.Errorf("expected security value 1 (DEFINER), got %d", secVal.Ival)
	}
}

// TestCreateProcedure tests: CREATE PROCEDURE do_something() LANGUAGE sql AS 'SELECT 1'
func TestCreateProcedure(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE PROCEDURE do_something() LANGUAGE sql AS 'SELECT 1'")

	if stmt.IsOrReplace {
		t.Error("expected IsOrReplace to be false")
	}
	if stmt.Funcname == nil || len(stmt.Funcname.Items) != 1 {
		t.Fatalf("expected 1 name component, got %v", stmt.Funcname)
	}
	fname := stmt.Funcname.Items[0].(*nodes.String).Str
	if fname != "do_something" {
		t.Errorf("expected procedure name 'do_something', got %q", fname)
	}
	if stmt.ReturnType != nil {
		t.Errorf("expected nil ReturnType for procedure, got %v", stmt.ReturnType)
	}
	// Check isProcedure marker
	procDef := findDefElem(stmt.Options, "isProcedure")
	if procDef == nil {
		t.Fatal("expected isProcedure option for procedure")
	}
	procVal, ok := procDef.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected isProcedure arg to be *nodes.Integer, got %T", procDef.Arg)
	}
	if procVal.Ival != 1 {
		t.Errorf("expected isProcedure value 1, got %d", procVal.Ival)
	}
}

// TestCreateFunctionNoReturns tests: CREATE FUNCTION foo() LANGUAGE sql AS 'SELECT 1'
// (function without RETURNS clause)
func TestCreateFunctionNoReturns(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() LANGUAGE sql AS 'SELECT 1'")

	if stmt.Funcname == nil || len(stmt.Funcname.Items) != 1 {
		t.Fatalf("expected 1 name component, got %v", stmt.Funcname)
	}
	fname := stmt.Funcname.Items[0].(*nodes.String).Str
	if fname != "foo" {
		t.Errorf("expected function name 'foo', got %q", fname)
	}
	if stmt.ReturnType != nil {
		t.Errorf("expected nil ReturnType, got %v", stmt.ReturnType)
	}
}

// TestCreateFunctionStableVolatile tests multiple volatility options
func TestCreateFunctionStableVolatile(t *testing.T) {
	tests := []struct {
		sql      string
		expected string
	}{
		{"CREATE FUNCTION f() RETURNS int LANGUAGE sql STABLE AS '1'", "stable"},
		{"CREATE FUNCTION f() RETURNS int LANGUAGE sql VOLATILE AS '1'", "volatile"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			stmt := parseCreateFunctionStmt(t, tt.sql)
			volDef := findDefElem(stmt.Options, "volatility")
			if volDef == nil {
				t.Fatal("expected volatility option")
			}
			volStr, ok := volDef.Arg.(*nodes.String)
			if !ok {
				t.Fatalf("expected volatility arg to be *nodes.String, got %T", volDef.Arg)
			}
			if volStr.Str != tt.expected {
				t.Errorf("expected volatility %q, got %q", tt.expected, volStr.Str)
			}
		})
	}
}

// TestCreateFunctionCalledOnNullInput tests: CREATE FUNCTION foo() RETURNS integer CALLED ON NULL INPUT LANGUAGE sql AS 'SELECT 1'
func TestCreateFunctionCalledOnNullInput(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS integer CALLED ON NULL INPUT LANGUAGE sql AS 'SELECT 1'")

	strictDef := findDefElem(stmt.Options, "strict")
	if strictDef == nil {
		t.Fatal("expected strict option")
	}
	strictVal, ok := strictDef.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected strict arg to be *nodes.Integer, got %T", strictDef.Arg)
	}
	if strictVal.Ival != 0 {
		t.Errorf("expected strict value 0 (CALLED ON NULL INPUT), got %d", strictVal.Ival)
	}
}

// TestCreateFunctionReturnsNullOnNullInput tests: CREATE FUNCTION foo() RETURNS integer RETURNS NULL ON NULL INPUT LANGUAGE sql AS 'SELECT 1'
func TestCreateFunctionReturnsNullOnNullInput(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS integer RETURNS NULL ON NULL INPUT LANGUAGE sql AS 'SELECT 1'")

	strictDef := findDefElem(stmt.Options, "strict")
	if strictDef == nil {
		t.Fatal("expected strict option")
	}
	strictVal, ok := strictDef.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected strict arg to be *nodes.Integer, got %T", strictDef.Arg)
	}
	if strictVal.Ival != 1 {
		t.Errorf("expected strict value 1 (RETURNS NULL ON NULL INPUT), got %d", strictVal.Ival)
	}
}

// TestCreateFunctionSecurityInvoker tests: CREATE FUNCTION foo() RETURNS integer SECURITY INVOKER LANGUAGE sql AS 'SELECT 1'
func TestCreateFunctionSecurityInvoker(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS integer SECURITY INVOKER LANGUAGE sql AS 'SELECT 1'")

	secDef := findDefElem(stmt.Options, "security")
	if secDef == nil {
		t.Fatal("expected security option")
	}
	secVal, ok := secDef.Arg.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected security arg to be *nodes.Integer, got %T", secDef.Arg)
	}
	if secVal.Ival != 0 {
		t.Errorf("expected security value 0 (INVOKER), got %d", secVal.Ival)
	}
}

// TestCreateFunctionAsTwoStrings tests: CREATE FUNCTION foo() RETURNS integer AS 'obj_file', 'link_symbol' LANGUAGE C
func TestCreateFunctionAsTwoStrings(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo() RETURNS integer AS 'obj_file', 'link_symbol' LANGUAGE C")

	asDef := findDefElem(stmt.Options, "as")
	if asDef == nil {
		t.Fatal("expected as option")
	}
	asList, ok := asDef.Arg.(*nodes.List)
	if !ok {
		t.Fatalf("expected as arg to be *nodes.List, got %T", asDef.Arg)
	}
	if len(asList.Items) != 2 {
		t.Fatalf("expected 2 items in AS list, got %d", len(asList.Items))
	}
	s1 := asList.Items[0].(*nodes.String).Str
	s2 := asList.Items[1].(*nodes.String).Str
	if s1 != "obj_file" {
		t.Errorf("expected first AS string 'obj_file', got %q", s1)
	}
	if s2 != "link_symbol" {
		t.Errorf("expected second AS string 'link_symbol', got %q", s2)
	}
}

// TestFuncTypePctType tests %TYPE syntax in function parameter and return types.
// e.g., CREATE FUNCTION foo(hobbies_r.name%TYPE) RETURNS hobbies_r.person%TYPE ...
func TestFuncTypePctType(t *testing.T) {
	sqls := []string{
		`CREATE FUNCTION foo(hobbies_r.name%TYPE) RETURNS hobbies_r.person%TYPE AS 'select 1' LANGUAGE SQL`,
		`CREATE FUNCTION bar(x mytable.mycolumn%TYPE) RETURNS void AS 'select 1' LANGUAGE SQL`,
	}
	for _, sql := range sqls {
		t.Run(sql[:40], func(t *testing.T) {
			_, err := parser.Parse(sql)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", sql[:60], err)
			}
		})
	}
}

// TestFuncTypePctTypeReturnPctType tests that the return type has PctType=true when using %TYPE.
func TestFuncTypePctTypeReturnPctType(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, `CREATE FUNCTION foo(x mytable.mycolumn%TYPE) RETURNS mytable.mycolumn%TYPE AS 'select 1' LANGUAGE SQL`)

	// Check return type has PctType set
	if stmt.ReturnType == nil {
		t.Fatal("expected ReturnType to be set")
	}
	if !stmt.ReturnType.PctType {
		t.Error("expected ReturnType.PctType to be true")
	}

	// Check parameter type has PctType set
	if stmt.Parameters == nil || len(stmt.Parameters.Items) != 1 {
		t.Fatalf("expected 1 parameter, got %v", stmt.Parameters)
	}
	p1 := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	if p1.ArgType == nil {
		t.Fatal("expected param type to be set")
	}
	if !p1.ArgType.PctType {
		t.Error("expected param ArgType.PctType to be true")
	}
}

// TestCreateFunctionDefaultEqualsParam tests: CREATE FUNCTION foo(x integer = 0) RETURNS integer LANGUAGE sql AS 'SELECT x'
func TestCreateFunctionDefaultEqualsParam(t *testing.T) {
	stmt := parseCreateFunctionStmt(t, "CREATE FUNCTION foo(x integer = 0) RETURNS integer LANGUAGE sql AS 'SELECT x'")

	if stmt.Parameters == nil || len(stmt.Parameters.Items) != 1 {
		t.Fatalf("expected 1 parameter, got %v", stmt.Parameters)
	}
	p1 := stmt.Parameters.Items[0].(*nodes.FunctionParameter)
	if p1.Name != "x" {
		t.Errorf("expected param name 'x', got %q", p1.Name)
	}
	if p1.Defexpr == nil {
		t.Error("expected default expression to be set")
	}
}
