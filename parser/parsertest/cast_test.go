package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE CAST tests
// =============================================================================

func TestCreateCastWithFunction(t *testing.T) {
	input := "CREATE CAST (bigint AS int4) WITH FUNCTION int4(bigint)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateCastStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateCastStmt, got %T", result.Items[0])
	}
	if stmt.Sourcetype == nil {
		t.Fatal("expected non-nil Sourcetype")
	}
	if stmt.Targettype == nil {
		t.Fatal("expected non-nil Targettype")
	}
	if stmt.Func == nil {
		t.Fatal("expected non-nil Func")
	}
	if stmt.Inout {
		t.Error("expected Inout to be false")
	}
	if stmt.Context != nodes.COERCION_EXPLICIT {
		t.Errorf("expected COERCION_EXPLICIT, got %d", stmt.Context)
	}
}

func TestCreateCastWithoutFunction(t *testing.T) {
	input := "CREATE CAST (bigint AS int4) WITHOUT FUNCTION"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateCastStmt)
	if stmt.Func != nil {
		t.Error("expected nil Func")
	}
	if stmt.Inout {
		t.Error("expected Inout to be false")
	}
}

func TestCreateCastWithInout(t *testing.T) {
	input := "CREATE CAST (bigint AS int4) WITH INOUT"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateCastStmt)
	if !stmt.Inout {
		t.Error("expected Inout to be true")
	}
}

func TestCreateCastAsAssignment(t *testing.T) {
	input := "CREATE CAST (bigint AS int4) WITH FUNCTION int4(bigint) AS ASSIGNMENT"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateCastStmt)
	if stmt.Context != nodes.COERCION_ASSIGNMENT {
		t.Errorf("expected COERCION_ASSIGNMENT, got %d", stmt.Context)
	}
}

func TestCreateCastAsImplicit(t *testing.T) {
	input := "CREATE CAST (bigint AS int4) WITH FUNCTION int4(bigint) AS IMPLICIT"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateCastStmt)
	if stmt.Context != nodes.COERCION_IMPLICIT {
		t.Errorf("expected COERCION_IMPLICIT, got %d", stmt.Context)
	}
}

// =============================================================================
// DROP CAST tests
// =============================================================================

func TestDropCast(t *testing.T) {
	input := "DROP CAST (bigint AS int4)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	if stmt.RemoveType != int(nodes.OBJECT_CAST) {
		t.Errorf("expected OBJECT_CAST, got %d", stmt.RemoveType)
	}
}

func TestDropCastIfExists(t *testing.T) {
	input := "DROP CAST IF EXISTS (bigint AS int4) CASCADE"
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
// CREATE CONVERSION tests
// =============================================================================

func TestCreateConversion(t *testing.T) {
	input := "CREATE CONVERSION myconv FOR 'UTF8' TO 'LATIN1' FROM utf8_to_iso8859_1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateConversionStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateConversionStmt, got %T", result.Items[0])
	}
	if stmt.ForEncodingName != "UTF8" {
		t.Errorf("expected 'UTF8', got %q", stmt.ForEncodingName)
	}
	if stmt.ToEncodingName != "LATIN1" {
		t.Errorf("expected 'LATIN1', got %q", stmt.ToEncodingName)
	}
	if stmt.Def {
		t.Error("expected Def to be false")
	}
}

func TestCreateDefaultConversion(t *testing.T) {
	input := "CREATE DEFAULT CONVERSION myconv FOR 'UTF8' TO 'LATIN1' FROM utf8_to_iso8859_1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateConversionStmt)
	if !stmt.Def {
		t.Error("expected Def to be true")
	}
}

// =============================================================================
// CREATE TRANSFORM tests
// =============================================================================

func TestCreateTransform(t *testing.T) {
	input := "CREATE TRANSFORM FOR int LANGUAGE plpgsql (FROM SQL WITH FUNCTION from_func(int), TO SQL WITH FUNCTION to_func(int))"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateTransformStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateTransformStmt, got %T", result.Items[0])
	}
	if stmt.Lang != "plpgsql" {
		t.Errorf("expected lang 'plpgsql', got %q", stmt.Lang)
	}
	if stmt.Fromsql == nil {
		t.Error("expected non-nil Fromsql")
	}
	if stmt.Tosql == nil {
		t.Error("expected non-nil Tosql")
	}
}

func TestDropTransform(t *testing.T) {
	input := "DROP TRANSFORM FOR int LANGUAGE plpgsql"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	if stmt.RemoveType != int(nodes.OBJECT_TRANSFORM) {
		t.Errorf("expected OBJECT_TRANSFORM, got %d", stmt.RemoveType)
	}
}
