package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE AGGREGATE tests
// =============================================================================

// TestCreateAggregateNewStyle tests: CREATE AGGREGATE myagg(integer) (SFUNC = int4_avg_accum, STYPE = _int8)
func TestCreateAggregateNewStyle(t *testing.T) {
	input := "CREATE AGGREGATE myagg(integer) (SFUNC = int4_avg_accum, STYPE = _int8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_AGGREGATE {
		t.Errorf("expected OBJECT_AGGREGATE, got %d", stmt.Kind)
	}
	if stmt.Oldstyle {
		t.Error("expected Oldstyle=false")
	}
	if stmt.Defnames == nil || len(stmt.Defnames.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "myagg" {
		t.Errorf("expected 'myagg', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Args == nil || len(stmt.Args.Items) == 0 {
		t.Fatal("expected aggregate args")
	}
	if stmt.Definition == nil || len(stmt.Definition.Items) < 2 {
		t.Fatal("expected at least 2 definition elements")
	}
}

// TestCreateAggregateMultiArg tests: CREATE AGGREGATE myagg(integer, integer) (SFUNC = f, STYPE = int8)
func TestCreateAggregateMultiArg(t *testing.T) {
	input := "CREATE AGGREGATE myagg(integer, integer) (SFUNC = f, STYPE = int8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_AGGREGATE {
		t.Errorf("expected OBJECT_AGGREGATE, got %d", stmt.Kind)
	}
	if stmt.Args == nil || len(stmt.Args.Items) < 2 {
		t.Fatal("expected at least 2 aggregate args")
	}
}

// TestCreateAggregateStar tests: CREATE AGGREGATE myagg(*) (SFUNC = f, STYPE = int8)
func TestCreateAggregateStar(t *testing.T) {
	input := "CREATE AGGREGATE myagg(*) (SFUNC = f, STYPE = int8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_AGGREGATE {
		t.Errorf("expected OBJECT_AGGREGATE, got %d", stmt.Kind)
	}
	// Args should have marker for star
	if stmt.Args == nil {
		t.Fatal("expected aggregate args (star marker)")
	}
}

// TestCreateAggregateOldStyle tests: CREATE AGGREGATE myagg(basetype = integer, sfunc = f, stype = int8)
func TestCreateAggregateOldStyle(t *testing.T) {
	input := "CREATE AGGREGATE myagg(basetype = integer, sfunc = f, stype = int8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_AGGREGATE {
		t.Errorf("expected OBJECT_AGGREGATE, got %d", stmt.Kind)
	}
	if !stmt.Oldstyle {
		t.Error("expected Oldstyle=true")
	}
	if stmt.Definition == nil || len(stmt.Definition.Items) < 3 {
		t.Fatal("expected at least 3 definition elements")
	}
}

// TestCreateAggregateOrReplace tests: CREATE OR REPLACE AGGREGATE myagg(integer) (SFUNC = f, STYPE = int8)
func TestCreateAggregateOrReplace(t *testing.T) {
	input := "CREATE OR REPLACE AGGREGATE myagg(integer) (SFUNC = f, STYPE = int8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if !stmt.Replace {
		t.Error("expected Replace=true")
	}
}

// =============================================================================
// CREATE OPERATOR tests
// =============================================================================

// TestCreateOperatorBasic tests: CREATE OPERATOR === (leftarg = text, rightarg = text, procedure = texteq)
func TestCreateOperatorBasic(t *testing.T) {
	input := "CREATE OPERATOR === (leftarg = text, rightarg = text, procedure = texteq)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_OPERATOR {
		t.Errorf("expected OBJECT_OPERATOR, got %d", stmt.Kind)
	}
	if stmt.Defnames == nil || len(stmt.Defnames.Items) != 1 {
		t.Fatal("expected 1 name component for operator")
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "===" {
		t.Errorf("expected '===', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Definition == nil || len(stmt.Definition.Items) < 3 {
		t.Fatal("expected at least 3 definition elements")
	}
}

// TestCreateOperatorPlus tests: CREATE OPERATOR + (leftarg = integer, rightarg = integer, procedure = int4pl)
func TestCreateOperatorPlus(t *testing.T) {
	input := "CREATE OPERATOR + (leftarg = integer, rightarg = integer, procedure = int4pl)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_OPERATOR {
		t.Errorf("expected OBJECT_OPERATOR, got %d", stmt.Kind)
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "+" {
		t.Errorf("expected '+', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
}

// TestCreateOperatorSchemaQualified tests: CREATE OPERATOR myschema.+ (leftarg = integer, rightarg = integer, procedure = int4pl)
func TestCreateOperatorSchemaQualified(t *testing.T) {
	input := "CREATE OPERATOR myschema.+ (leftarg = integer, rightarg = integer, procedure = int4pl)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_OPERATOR {
		t.Errorf("expected OBJECT_OPERATOR, got %d", stmt.Kind)
	}
	if stmt.Defnames == nil || len(stmt.Defnames.Items) != 2 {
		t.Fatalf("expected 2 name components, got %d", stmt.Defnames.Len())
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "myschema" {
		t.Errorf("expected 'myschema', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Defnames.Items[1].(*nodes.String).Str != "+" {
		t.Errorf("expected '+', got %q", stmt.Defnames.Items[1].(*nodes.String).Str)
	}
}

// =============================================================================
// CREATE TYPE tests
// =============================================================================

// TestCreateTypeShell tests: CREATE TYPE mytype (shell type, no definition)
func TestCreateTypeShell(t *testing.T) {
	input := "CREATE TYPE mytype"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_TYPE {
		t.Errorf("expected OBJECT_TYPE, got %d", stmt.Kind)
	}
	if stmt.Defnames == nil || len(stmt.Defnames.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "mytype" {
		t.Errorf("expected 'mytype', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Definition != nil {
		t.Error("shell type should have nil definition")
	}
}

// TestCreateTypeBase tests: CREATE TYPE mytype (INPUT = mytype_in, OUTPUT = mytype_out)
func TestCreateTypeBase(t *testing.T) {
	input := "CREATE TYPE mytype (INPUT = mytype_in, OUTPUT = mytype_out)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_TYPE {
		t.Errorf("expected OBJECT_TYPE, got %d", stmt.Kind)
	}
	if stmt.Definition == nil || len(stmt.Definition.Items) != 2 {
		t.Fatal("expected 2 definition elements")
	}
}

// TestCreateTypeComposite tests: CREATE TYPE complex AS (r float8, i float8)
func TestCreateTypeComposite(t *testing.T) {
	input := "CREATE TYPE complex AS (r float8, i float8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CompositeTypeStmt)
	if !ok {
		t.Fatalf("expected *nodes.CompositeTypeStmt, got %T", result.Items[0])
	}

	if stmt.Typevar == nil {
		t.Fatal("expected Typevar to be set")
	}
	if stmt.Typevar.Relname != "complex" {
		t.Errorf("expected 'complex', got %q", stmt.Typevar.Relname)
	}
	if stmt.Coldeflist == nil || len(stmt.Coldeflist.Items) != 2 {
		t.Fatal("expected 2 columns")
	}

	col1 := stmt.Coldeflist.Items[0].(*nodes.ColumnDef)
	if col1.Colname != "r" {
		t.Errorf("expected column 'r', got %q", col1.Colname)
	}

	col2 := stmt.Coldeflist.Items[1].(*nodes.ColumnDef)
	if col2.Colname != "i" {
		t.Errorf("expected column 'i', got %q", col2.Colname)
	}
}

// TestCreateTypeEnum tests: CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy')
func TestCreateTypeEnum(t *testing.T) {
	input := "CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateEnumStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateEnumStmt, got %T", result.Items[0])
	}

	if stmt.TypeName == nil || len(stmt.TypeName.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.TypeName.Items[0].(*nodes.String).Str != "mood" {
		t.Errorf("expected 'mood', got %q", stmt.TypeName.Items[0].(*nodes.String).Str)
	}
	if stmt.Vals == nil || len(stmt.Vals.Items) != 3 {
		t.Fatalf("expected 3 enum values, got %d", stmt.Vals.Len())
	}
	vals := []string{"sad", "ok", "happy"}
	for i, v := range vals {
		got := stmt.Vals.Items[i].(*nodes.String).Str
		if got != v {
			t.Errorf("expected val[%d]=%q, got %q", i, v, got)
		}
	}
}

// TestCreateTypeEnumEmpty tests: CREATE TYPE mood AS ENUM ()
func TestCreateTypeEnumEmpty(t *testing.T) {
	input := "CREATE TYPE mood AS ENUM ()"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateEnumStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateEnumStmt, got %T", result.Items[0])
	}

	if stmt.Vals != nil && len(stmt.Vals.Items) != 0 {
		t.Errorf("expected empty vals for empty enum, got %d", len(stmt.Vals.Items))
	}
}

// TestCreateTypeRange tests: CREATE TYPE float8_range AS RANGE (SUBTYPE = float8)
func TestCreateTypeRange(t *testing.T) {
	input := "CREATE TYPE float8_range AS RANGE (SUBTYPE = float8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateRangeStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateRangeStmt, got %T", result.Items[0])
	}

	if stmt.TypeName == nil || len(stmt.TypeName.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.TypeName.Items[0].(*nodes.String).Str != "float8_range" {
		t.Errorf("expected 'float8_range', got %q", stmt.TypeName.Items[0].(*nodes.String).Str)
	}
	if stmt.Params == nil || len(stmt.Params.Items) < 1 {
		t.Fatal("expected at least 1 range parameter")
	}
}

// =============================================================================
// CREATE TEXT SEARCH tests
// =============================================================================

// TestCreateTextSearchParser tests: CREATE TEXT SEARCH PARSER myparser (START = prsd_start, GETTOKEN = prsd_nexttoken, END = prsd_end, LEXTYPES = prsd_lextype)
func TestCreateTextSearchParser(t *testing.T) {
	input := "CREATE TEXT SEARCH PARSER myparser (START = prsd_start, GETTOKEN = prsd_nexttoken, END = prsd_end, LEXTYPES = prsd_lextype)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_TSPARSER {
		t.Errorf("expected OBJECT_TSPARSER, got %d", stmt.Kind)
	}
	if stmt.Defnames == nil || len(stmt.Defnames.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "myparser" {
		t.Errorf("expected 'myparser', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Definition == nil || len(stmt.Definition.Items) != 4 {
		t.Fatalf("expected 4 definition elements, got %d", stmt.Definition.Len())
	}
}

// TestCreateTextSearchDictionary tests: CREATE TEXT SEARCH DICTIONARY mydict (TEMPLATE = simple, STOPWORDS = english)
func TestCreateTextSearchDictionary(t *testing.T) {
	input := "CREATE TEXT SEARCH DICTIONARY mydict (TEMPLATE = simple, STOPWORDS = english)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_TSDICTIONARY {
		t.Errorf("expected OBJECT_TSDICTIONARY, got %d", stmt.Kind)
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "mydict" {
		t.Errorf("expected 'mydict', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Definition == nil || len(stmt.Definition.Items) != 2 {
		t.Fatalf("expected 2 definition elements, got %d", stmt.Definition.Len())
	}
}

// TestCreateTextSearchTemplate tests: CREATE TEXT SEARCH TEMPLATE mytempl (INIT = dsimple_init, LEXIZE = dsimple_lexize)
func TestCreateTextSearchTemplate(t *testing.T) {
	input := "CREATE TEXT SEARCH TEMPLATE mytempl (INIT = dsimple_init, LEXIZE = dsimple_lexize)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_TSTEMPLATE {
		t.Errorf("expected OBJECT_TSTEMPLATE, got %d", stmt.Kind)
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "mytempl" {
		t.Errorf("expected 'mytempl', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
}

// TestCreateTextSearchConfiguration tests: CREATE TEXT SEARCH CONFIGURATION myconfig (PARSER = default)
func TestCreateTextSearchConfiguration(t *testing.T) {
	input := "CREATE TEXT SEARCH CONFIGURATION myconfig (PARSER = default)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_TSCONFIGURATION {
		t.Errorf("expected OBJECT_TSCONFIGURATION, got %d", stmt.Kind)
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "myconfig" {
		t.Errorf("expected 'myconfig', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
}

// =============================================================================
// CREATE COLLATION tests
// =============================================================================

// TestCreateCollationDefinition tests: CREATE COLLATION mycoll (LOCALE = 'en_US.utf8')
func TestCreateCollationDefinition(t *testing.T) {
	input := "CREATE COLLATION mycoll (LOCALE = 'en_US.utf8')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_COLLATION {
		t.Errorf("expected OBJECT_COLLATION, got %d", stmt.Kind)
	}
	if stmt.Defnames == nil || len(stmt.Defnames.Items) != 1 {
		t.Fatal("expected 1 name component")
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "mycoll" {
		t.Errorf("expected 'mycoll', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Definition == nil || len(stmt.Definition.Items) != 1 {
		t.Fatal("expected 1 definition element")
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists=false")
	}
}

// TestCreateCollationIfNotExists tests: CREATE COLLATION IF NOT EXISTS mycoll (LOCALE = 'en_US.utf8')
func TestCreateCollationIfNotExists(t *testing.T) {
	input := "CREATE COLLATION IF NOT EXISTS mycoll (LOCALE = 'en_US.utf8')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_COLLATION {
		t.Errorf("expected OBJECT_COLLATION, got %d", stmt.Kind)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
}

// TestCreateCollationFrom tests: CREATE COLLATION mycoll FROM "C"
func TestCreateCollationFrom(t *testing.T) {
	input := `CREATE COLLATION mycoll FROM "C"`

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_COLLATION {
		t.Errorf("expected OBJECT_COLLATION, got %d", stmt.Kind)
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "mycoll" {
		t.Errorf("expected 'mycoll', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	// Definition should have "from" element
	if stmt.Definition == nil || len(stmt.Definition.Items) != 1 {
		t.Fatal("expected 1 definition element for FROM clause")
	}
	defElem := stmt.Definition.Items[0].(*nodes.DefElem)
	if defElem.Defname != "from" {
		t.Errorf("expected defname 'from', got %q", defElem.Defname)
	}
}

// TestCreateCollationIfNotExistsFrom tests: CREATE COLLATION IF NOT EXISTS mycoll FROM "C"
func TestCreateCollationIfNotExistsFrom(t *testing.T) {
	input := `CREATE COLLATION IF NOT EXISTS mycoll FROM "C"`

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_COLLATION {
		t.Errorf("expected OBJECT_COLLATION, got %d", stmt.Kind)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists=true")
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "mycoll" {
		t.Errorf("expected 'mycoll', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	// Definition should have "from" element
	if stmt.Definition == nil || len(stmt.Definition.Items) != 1 {
		t.Fatal("expected 1 definition element for FROM clause")
	}
	defElem := stmt.Definition.Items[0].(*nodes.DefElem)
	if defElem.Defname != "from" {
		t.Errorf("expected defname 'from', got %q", defElem.Defname)
	}
}

// =============================================================================
// Additional edge case tests
// =============================================================================

// TestCreateTypeCompositeEmpty tests: CREATE TYPE empty_type AS ()
func TestCreateTypeCompositeEmpty(t *testing.T) {
	input := "CREATE TYPE empty_type AS ()"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CompositeTypeStmt)
	if !ok {
		t.Fatalf("expected *nodes.CompositeTypeStmt, got %T", result.Items[0])
	}

	if stmt.Typevar.Relname != "empty_type" {
		t.Errorf("expected 'empty_type', got %q", stmt.Typevar.Relname)
	}
	// Empty column list is valid
}

// TestCreateTypeSchemaQualified tests: CREATE TYPE myschema.mytype (INPUT = mytype_in, OUTPUT = mytype_out)
func TestCreateTypeSchemaQualified(t *testing.T) {
	input := "CREATE TYPE myschema.mytype (INPUT = mytype_in, OUTPUT = mytype_out)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Kind != nodes.OBJECT_TYPE {
		t.Errorf("expected OBJECT_TYPE, got %d", stmt.Kind)
	}
	if stmt.Defnames == nil || len(stmt.Defnames.Items) != 2 {
		t.Fatalf("expected 2 name components, got %d", stmt.Defnames.Len())
	}
	if stmt.Defnames.Items[0].(*nodes.String).Str != "myschema" {
		t.Errorf("expected 'myschema', got %q", stmt.Defnames.Items[0].(*nodes.String).Str)
	}
	if stmt.Defnames.Items[1].(*nodes.String).Str != "mytype" {
		t.Errorf("expected 'mytype', got %q", stmt.Defnames.Items[1].(*nodes.String).Str)
	}
}

// TestDefElemNoValue tests that def_elem without '=' works: CREATE TYPE t (PASSEDBYVALUE)
func TestDefElemNoValue(t *testing.T) {
	input := "CREATE TYPE mytype (PASSEDBYVALUE)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Definition == nil || len(stmt.Definition.Items) != 1 {
		t.Fatal("expected 1 definition element")
	}
	defElem := stmt.Definition.Items[0].(*nodes.DefElem)
	if defElem.Defname != "passedbyvalue" {
		t.Errorf("expected 'passedbyvalue', got %q", defElem.Defname)
	}
	if defElem.Arg != nil {
		t.Error("expected nil Arg for key-only def_elem")
	}
}

// TestDefElemNumericValue tests: CREATE TYPE mytype (INTERNALLENGTH = 16)
func TestDefElemNumericValue(t *testing.T) {
	input := "CREATE TYPE mytype (INTERNALLENGTH = 16)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Definition == nil || len(stmt.Definition.Items) != 1 {
		t.Fatal("expected 1 definition element")
	}
	defElem := stmt.Definition.Items[0].(*nodes.DefElem)
	if defElem.Defname != "internallength" {
		t.Errorf("expected 'internallength', got %q", defElem.Defname)
	}
}

// TestDefElemStringValue tests: CREATE TEXT SEARCH DICTIONARY d (TEMPLATE = simple, STOPWORDS = 'english')
func TestDefElemStringValue(t *testing.T) {
	input := "CREATE TEXT SEARCH DICTIONARY d (TEMPLATE = simple, STOPWORDS = 'english')"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.DefineStmt)
	if !ok {
		t.Fatalf("expected *nodes.DefineStmt, got %T", result.Items[0])
	}

	if stmt.Definition == nil || len(stmt.Definition.Items) != 2 {
		t.Fatal("expected 2 definition elements")
	}
	// Check second elem has string value
	defElem := stmt.Definition.Items[1].(*nodes.DefElem)
	if defElem.Defname != "stopwords" {
		t.Errorf("expected 'stopwords', got %q", defElem.Defname)
	}
	str, ok := defElem.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String arg, got %T", defElem.Arg)
	}
	if str.Str != "english" {
		t.Errorf("expected 'english', got %q", str.Str)
	}
}

// TestCreateTypeRangeMultipleParams tests: CREATE TYPE r AS RANGE (SUBTYPE = float8, SUBTYPE_OPCLASS = float8_ops)
func TestCreateTypeRangeMultipleParams(t *testing.T) {
	input := "CREATE TYPE r AS RANGE (SUBTYPE = float8, SUBTYPE_OPCLASS = float8_ops)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateRangeStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateRangeStmt, got %T", result.Items[0])
	}

	if stmt.Params == nil || len(stmt.Params.Items) != 2 {
		t.Fatalf("expected 2 range params, got %d", stmt.Params.Len())
	}
}

// TestCreateCompositeTypeSchemaQualified tests: CREATE TYPE myschema.complex AS (r float8, i float8)
func TestCreateCompositeTypeSchemaQualified(t *testing.T) {
	input := "CREATE TYPE myschema.complex AS (r float8, i float8)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CompositeTypeStmt)
	if !ok {
		t.Fatalf("expected *nodes.CompositeTypeStmt, got %T", result.Items[0])
	}

	if stmt.Typevar.Schemaname != "myschema" {
		t.Errorf("expected schema 'myschema', got %q", stmt.Typevar.Schemaname)
	}
	if stmt.Typevar.Relname != "complex" {
		t.Errorf("expected 'complex', got %q", stmt.Typevar.Relname)
	}
}
