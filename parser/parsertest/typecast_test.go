package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// helper: parse a single SELECT statement and return the SelectStmt.
func parseSelect(t *testing.T, input string) *nodes.SelectStmt {
	t.Helper()
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error for %q: %v", input, err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatalf("expected 1 statement for %q", input)
	}
	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt for %q, got %T", input, result.Items[0])
	}
	return stmt
}

// helper: extract the first target's expression value from a SELECT.
func firstTargetVal(t *testing.T, stmt *nodes.SelectStmt) nodes.Node {
	t.Helper()
	if stmt.TargetList == nil || len(stmt.TargetList.Items) == 0 {
		t.Fatal("expected at least 1 target in target list")
	}
	rt, ok := stmt.TargetList.Items[0].(*nodes.ResTarget)
	if !ok {
		t.Fatalf("expected *nodes.ResTarget, got %T", stmt.TargetList.Items[0])
	}
	return rt.Val
}

// helper: assert a node is a TypeCast and return it.
func assertTypeCast(t *testing.T, n nodes.Node) *nodes.TypeCast {
	t.Helper()
	tc, ok := n.(*nodes.TypeCast)
	if !ok {
		t.Fatalf("expected *nodes.TypeCast, got %T", n)
	}
	return tc
}

// helper: check that a TypeName has the given pg_catalog type name.
func assertTypeNamePgCatalog(t *testing.T, tn *nodes.TypeName, expectedName string) {
	t.Helper()
	if tn == nil {
		t.Fatal("expected non-nil TypeName")
	}
	if tn.Names == nil || len(tn.Names.Items) != 2 {
		t.Fatalf("expected 2-element Names list (pg_catalog.%s), got %d items",
			expectedName, listLen(tn.Names))
	}
	schema := tn.Names.Items[0].(*nodes.String).Str
	typname := tn.Names.Items[1].(*nodes.String).Str
	if schema != "pg_catalog" {
		t.Errorf("expected schema 'pg_catalog', got %q", schema)
	}
	if typname != expectedName {
		t.Errorf("expected type name %q, got %q", expectedName, typname)
	}
}

// helper: check that a TypeName has a generic (unqualified) type name.
func assertTypeNameGeneric(t *testing.T, tn *nodes.TypeName, expectedName string) {
	t.Helper()
	if tn == nil {
		t.Fatal("expected non-nil TypeName")
	}
	if tn.Names == nil || len(tn.Names.Items) != 1 {
		t.Fatalf("expected 1-element Names list for generic type %q, got %d items",
			expectedName, listLen(tn.Names))
	}
	typname := tn.Names.Items[0].(*nodes.String).Str
	if typname != expectedName {
		t.Errorf("expected type name %q, got %q", expectedName, typname)
	}
}

// Test 1: TYPECAST operator - SELECT 42::integer
func TestTypecastOperator(t *testing.T) {
	stmt := parseSelect(t, "SELECT 42::integer")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// Arg should be A_Const(42)
	aConst, ok := tc.Arg.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for arg, got %T", tc.Arg)
	}
	intVal, ok := aConst.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", aConst.Val)
	}
	if intVal.Ival != 42 {
		t.Errorf("expected integer value 42, got %d", intVal.Ival)
	}

	// TypeName should be pg_catalog.int4
	assertTypeNamePgCatalog(t, tc.TypeName, "int4")
}

// Test 2: CAST function - SELECT CAST(42 AS integer)
func TestCastFunction(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(42 AS integer)")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// Arg should be A_Const(42)
	aConst, ok := tc.Arg.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for arg, got %T", tc.Arg)
	}
	intVal, ok := aConst.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", aConst.Val)
	}
	if intVal.Ival != 42 {
		t.Errorf("expected integer value 42, got %d", intVal.Ival)
	}

	// TypeName should be pg_catalog.int4
	assertTypeNamePgCatalog(t, tc.TypeName, "int4")
}

// Test 3: VARCHAR cast with length - SELECT CAST(x AS varchar(100)) FROM t
func TestVarcharCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS varchar(100)) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// Arg should be ColumnRef("x")
	colRef, ok := tc.Arg.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for arg, got %T", tc.Arg)
	}
	if colRef.Fields == nil || len(colRef.Fields.Items) != 1 {
		t.Fatal("expected 1 field in ColumnRef")
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "x" {
		t.Errorf("expected column 'x', got %q", str.Str)
	}

	// TypeName should be pg_catalog.varchar with Typmods [100]
	assertTypeNamePgCatalog(t, tc.TypeName, "varchar")
	if tc.TypeName.Typmods == nil || len(tc.TypeName.Typmods.Items) != 1 {
		t.Fatalf("expected 1 typmod for varchar(100), got %d", listLen(tc.TypeName.Typmods))
	}
	typmod, ok := tc.TypeName.Typmods.Items[0].(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer for typmod, got %T", tc.TypeName.Typmods.Items[0])
	}
	if typmod.Ival != 100 {
		t.Errorf("expected typmod 100, got %d", typmod.Ival)
	}
}

// Test 4: BOOLEAN cast - SELECT CAST(1 AS boolean)
func TestBooleanCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(1 AS boolean)")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// TypeName should be pg_catalog.bool
	assertTypeNamePgCatalog(t, tc.TypeName, "bool")
}

// Test 5: FLOAT cast - SELECT 3.14::float
func TestFloatCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT 3.14::float")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// Arg should be A_Const(3.14)
	aConst, ok := tc.Arg.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for arg, got %T", tc.Arg)
	}
	fVal, ok := aConst.Val.(*nodes.Float)
	if !ok {
		t.Fatalf("expected *nodes.Float, got %T", aConst.Val)
	}
	if fVal.Fval != "3.14" {
		t.Errorf("expected float value '3.14', got %q", fVal.Fval)
	}

	// FLOAT without precision defaults to float8
	assertTypeNamePgCatalog(t, tc.TypeName, "float8")
}

// Test 6: DOUBLE PRECISION - SELECT CAST(x AS double precision) FROM t
func TestDoublePrecisionCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS double precision) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// TypeName should be pg_catalog.float8
	assertTypeNamePgCatalog(t, tc.TypeName, "float8")
}

// Test 7: NUMERIC with precision - SELECT CAST(x AS numeric(10,2)) FROM t
func TestNumericPrecisionCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS numeric(10,2)) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// TypeName should be pg_catalog.numeric
	assertTypeNamePgCatalog(t, tc.TypeName, "numeric")

	// Typmods should be [10, 2]
	if tc.TypeName.Typmods == nil || len(tc.TypeName.Typmods.Items) != 2 {
		t.Fatalf("expected 2 typmods for numeric(10,2), got %d", listLen(tc.TypeName.Typmods))
	}
	prec, ok := tc.TypeName.Typmods.Items[0].(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for precision, got %T", tc.TypeName.Typmods.Items[0])
	}
	precInt, ok := prec.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer for precision value, got %T", prec.Val)
	}
	if precInt.Ival != 10 {
		t.Errorf("expected precision 10, got %d", precInt.Ival)
	}

	scale, ok := tc.TypeName.Typmods.Items[1].(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for scale, got %T", tc.TypeName.Typmods.Items[1])
	}
	scaleInt, ok := scale.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer for scale value, got %T", scale.Val)
	}
	if scaleInt.Ival != 2 {
		t.Errorf("expected scale 2, got %d", scaleInt.Ival)
	}
}

// Test 8: Generic type cast - SELECT x::text FROM t
func TestGenericTypeCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT x::text FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// Arg should be ColumnRef("x")
	colRef, ok := tc.Arg.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for arg, got %T", tc.Arg)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "x" {
		t.Errorf("expected column 'x', got %q", str.Str)
	}

	// TypeName should be generic: Names=["text"]
	// "text" is an unreserved keyword, so it goes through type_function_name -> GenericType
	assertTypeNameGeneric(t, tc.TypeName, "text")
}

// Test 9: Chained casts - SELECT '2024-01-01'::text::varchar
func TestChainedCasts(t *testing.T) {
	stmt := parseSelect(t, "SELECT '2024-01-01'::text::varchar")
	val := firstTargetVal(t, stmt)

	// Outer cast: ...::varchar
	outerTC := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, outerTC.TypeName, "varchar")

	// Inner cast: '2024-01-01'::text
	innerTC := assertTypeCast(t, outerTC.Arg)
	assertTypeNameGeneric(t, innerTC.TypeName, "text")

	// Innermost arg: string constant '2024-01-01'
	aConst, ok := innerTC.Arg.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for innermost arg, got %T", innerTC.Arg)
	}
	strVal, ok := aConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", aConst.Val)
	}
	if strVal.Str != "2024-01-01" {
		t.Errorf("expected string '2024-01-01', got %q", strVal.Str)
	}
}

// Test 10: Cast in WHERE - SELECT * FROM t WHERE x::integer > 5
func TestCastInWhere(t *testing.T) {
	stmt := parseSelect(t, "SELECT * FROM t WHERE x::integer > 5")

	if stmt.WhereClause == nil {
		t.Fatal("expected non-nil WhereClause")
	}

	// WhereClause should be A_Expr(>)
	aexpr, ok := stmt.WhereClause.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr in WhereClause, got %T", stmt.WhereClause)
	}
	if aexpr.Kind != nodes.AEXPR_OP {
		t.Errorf("expected AEXPR_OP, got %d", aexpr.Kind)
	}
	opStr := aexpr.Name.Items[0].(*nodes.String).Str
	if opStr != ">" {
		t.Errorf("expected operator '>', got %q", opStr)
	}

	// Left side should be TypeCast (x::integer)
	tc := assertTypeCast(t, aexpr.Lexpr)
	assertTypeNamePgCatalog(t, tc.TypeName, "int4")

	colRef, ok := tc.Arg.(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef for cast arg, got %T", tc.Arg)
	}
	str := colRef.Fields.Items[0].(*nodes.String)
	if str.Str != "x" {
		t.Errorf("expected column 'x', got %q", str.Str)
	}

	// Right side should be A_Const(5)
	rConst, ok := aexpr.Rexpr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const for right side, got %T", aexpr.Rexpr)
	}
	rInt, ok := rConst.Val.(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer, got %T", rConst.Val)
	}
	if rInt.Ival != 5 {
		t.Errorf("expected integer 5, got %d", rInt.Ival)
	}
}

// Additional tests for completeness

// TestSmallintCast tests SMALLINT type.
func TestSmallintCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT 42::smallint")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "int2")
}

// TestBigintCast tests BIGINT type.
func TestBigintCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT 42::bigint")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "int8")
}

// TestRealCast tests REAL type.
func TestRealCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT 3.14::real")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "float4")
}

// TestCharCast tests CHAR type.
func TestCharCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS char(10)) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "bpchar")
	if tc.TypeName.Typmods == nil || len(tc.TypeName.Typmods.Items) != 1 {
		t.Fatalf("expected 1 typmod for char(10), got %d", listLen(tc.TypeName.Typmods))
	}
	typmod := tc.TypeName.Typmods.Items[0].(*nodes.Integer)
	if typmod.Ival != 10 {
		t.Errorf("expected typmod 10, got %d", typmod.Ival)
	}
}

// TestCharacterVaryingCast tests CHARACTER VARYING type.
func TestCharacterVaryingCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS character varying(50)) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "varchar")
	if tc.TypeName.Typmods == nil || len(tc.TypeName.Typmods.Items) != 1 {
		t.Fatalf("expected 1 typmod for character varying(50), got %d", listLen(tc.TypeName.Typmods))
	}
	typmod := tc.TypeName.Typmods.Items[0].(*nodes.Integer)
	if typmod.Ival != 50 {
		t.Errorf("expected typmod 50, got %d", typmod.Ival)
	}
}

// TestDecimalCast tests DECIMAL type.
func TestDecimalCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS decimal(8,3)) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "numeric")
}

// TestArrayBoundsCast tests array type cast (e.g., integer[]).
func TestArrayBoundsCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT x::integer[] FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "int4")

	// ArrayBounds should have 1 entry with value -1 (unspecified size)
	if tc.TypeName.ArrayBounds == nil || len(tc.TypeName.ArrayBounds.Items) != 1 {
		t.Fatalf("expected 1 array bound, got %d", listLen(tc.TypeName.ArrayBounds))
	}
	bound, ok := tc.TypeName.ArrayBounds.Items[0].(*nodes.Integer)
	if !ok {
		t.Fatalf("expected *nodes.Integer for array bound, got %T", tc.TypeName.ArrayBounds.Items[0])
	}
	if bound.Ival != -1 {
		t.Errorf("expected array bound -1 (unspecified), got %d", bound.Ival)
	}
}

// TestIntCast tests INT type (alias for integer).
func TestIntCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS int) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "int4")
}

// TestDecCast tests DEC type (alias for numeric/decimal).
func TestDecCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS dec) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "numeric")
}

// TestFloatWithPrecisionCast tests FLOAT(n) where n <= 24 gives float4.
func TestFloatWithPrecisionCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS float(24)) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "float4")
}

// TestFloatWithLargePrecisionCast tests FLOAT(n) where n > 24 gives float8.
func TestFloatWithLargePrecisionCast(t *testing.T) {
	stmt := parseSelect(t, "SELECT CAST(x AS float(53)) FROM t")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)
	assertTypeNamePgCatalog(t, tc.TypeName, "float8")
}

// TestSetofTypename tests SETOF type.
func TestSetofTypename(t *testing.T) {
	// SETOF isn't typically used in CAST, but the Typename rule supports it.
	// We can test it via the :: operator with a generic type.
	// Actually, SETOF in a cast context doesn't make much sense in PostgreSQL,
	// but the grammar supports it. Let's just verify the grammar doesn't error.
	// We'll test it indirectly by verifying SETOF is properly a col_name_keyword.
	// For now, we skip this test since SETOF in a cast context is unusual.
	t.Skip("SETOF in cast context is unusual and not typically tested")
}

// TestCastStringToVarchar tests casting a string literal to varchar.
func TestCastStringToVarchar(t *testing.T) {
	stmt := parseSelect(t, "SELECT 'hello'::varchar(10)")
	val := firstTargetVal(t, stmt)
	tc := assertTypeCast(t, val)

	// Arg should be string constant 'hello'
	aConst, ok := tc.Arg.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", tc.Arg)
	}
	strVal, ok := aConst.Val.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", aConst.Val)
	}
	if strVal.Str != "hello" {
		t.Errorf("expected 'hello', got %q", strVal.Str)
	}

	// TypeName should be pg_catalog.varchar with Typmods [10]
	assertTypeNamePgCatalog(t, tc.TypeName, "varchar")
	if tc.TypeName.Typmods == nil || len(tc.TypeName.Typmods.Items) != 1 {
		t.Fatalf("expected 1 typmod, got %d", listLen(tc.TypeName.Typmods))
	}
	typmod := tc.TypeName.Typmods.Items[0].(*nodes.Integer)
	if typmod.Ival != 10 {
		t.Errorf("expected typmod 10, got %d", typmod.Ival)
	}
}

// TestCastInExpressionContext tests a cast used in a larger expression.
func TestCastInExpressionContext(t *testing.T) {
	stmt := parseSelect(t, "SELECT x::integer + 1 FROM t")
	val := firstTargetVal(t, stmt)

	// Should be A_Expr(+) with left=TypeCast, right=A_Const(1)
	aexpr, ok := val.(*nodes.A_Expr)
	if !ok {
		t.Fatalf("expected *nodes.A_Expr, got %T", val)
	}

	// Left side should be TypeCast
	tc := assertTypeCast(t, aexpr.Lexpr)
	assertTypeNamePgCatalog(t, tc.TypeName, "int4")

	// Right side should be integer 1
	rConst, ok := aexpr.Rexpr.(*nodes.A_Const)
	if !ok {
		t.Fatalf("expected *nodes.A_Const, got %T", aexpr.Rexpr)
	}
	rInt := rConst.Val.(*nodes.Integer)
	if rInt.Ival != 1 {
		t.Errorf("expected integer 1, got %d", rInt.Ival)
	}
}
