package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE TRIGGER tests
// =============================================================================

// TestCreateTriggerBeforeInsert tests: CREATE TRIGGER trig BEFORE INSERT ON t EXECUTE FUNCTION func()
func TestCreateTriggerBeforeInsert(t *testing.T) {
	input := "CREATE TRIGGER trig BEFORE INSERT ON t EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.CreateTrigStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateTrigStmt, got %T", result.Items[0])
	}
	if stmt.Trigname != "trig" {
		t.Errorf("expected trigname 'trig', got %q", stmt.Trigname)
	}
	if stmt.Timing != int16(nodes.TRIGGER_TYPE_BEFORE) {
		t.Errorf("expected TRIGGER_TYPE_BEFORE (%d), got %d", nodes.TRIGGER_TYPE_BEFORE, stmt.Timing)
	}
	if stmt.Events != int16(nodes.TRIGGER_TYPE_INSERT) {
		t.Errorf("expected TRIGGER_TYPE_INSERT (%d), got %d", nodes.TRIGGER_TYPE_INSERT, stmt.Events)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Error("expected relation 't'")
	}
	if stmt.IsConstraint {
		t.Error("expected IsConstraint to be false")
	}
}

// TestCreateTriggerAfterUpdate tests: CREATE TRIGGER trig AFTER UPDATE ON t EXECUTE FUNCTION func()
func TestCreateTriggerAfterUpdate(t *testing.T) {
	input := "CREATE TRIGGER trig AFTER UPDATE ON t EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if stmt.Timing != int16(nodes.TRIGGER_TYPE_AFTER) {
		t.Errorf("expected TRIGGER_TYPE_AFTER (0), got %d", stmt.Timing)
	}
	if stmt.Events != int16(nodes.TRIGGER_TYPE_UPDATE) {
		t.Errorf("expected TRIGGER_TYPE_UPDATE (%d), got %d", nodes.TRIGGER_TYPE_UPDATE, stmt.Events)
	}
}

// TestCreateTriggerAfterInsertOrUpdate tests: CREATE TRIGGER trig AFTER INSERT OR UPDATE ON t EXECUTE FUNCTION func()
func TestCreateTriggerAfterInsertOrUpdate(t *testing.T) {
	input := "CREATE TRIGGER trig AFTER INSERT OR UPDATE ON t EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	expected := int16(nodes.TRIGGER_TYPE_INSERT | nodes.TRIGGER_TYPE_UPDATE)
	if stmt.Events != expected {
		t.Errorf("expected events %d, got %d", expected, stmt.Events)
	}
}

// TestCreateTriggerBeforeDeleteForEachRow tests: CREATE TRIGGER trig BEFORE DELETE ON t FOR EACH ROW EXECUTE FUNCTION func()
func TestCreateTriggerBeforeDeleteForEachRow(t *testing.T) {
	input := "CREATE TRIGGER trig BEFORE DELETE ON t FOR EACH ROW EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if !stmt.Row {
		t.Error("expected Row to be true")
	}
	if stmt.Events != int16(nodes.TRIGGER_TYPE_DELETE) {
		t.Errorf("expected TRIGGER_TYPE_DELETE, got %d", stmt.Events)
	}
}

// TestCreateTriggerUpdateOfColumns tests: CREATE TRIGGER trig AFTER UPDATE OF col1, col2 ON t EXECUTE FUNCTION func()
func TestCreateTriggerUpdateOfColumns(t *testing.T) {
	input := "CREATE TRIGGER trig AFTER UPDATE OF col1, col2 ON t EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if stmt.Events != int16(nodes.TRIGGER_TYPE_UPDATE) {
		t.Errorf("expected TRIGGER_TYPE_UPDATE, got %d", stmt.Events)
	}
	if stmt.Columns == nil || len(stmt.Columns.Items) != 2 {
		t.Fatalf("expected 2 columns, got %v", stmt.Columns)
	}
}

// TestCreateTriggerInsteadOf tests: CREATE TRIGGER trig INSTEAD OF INSERT ON v EXECUTE FUNCTION func()
func TestCreateTriggerInsteadOf(t *testing.T) {
	input := "CREATE TRIGGER trig INSTEAD OF INSERT ON v EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if stmt.Timing != int16(nodes.TRIGGER_TYPE_INSTEAD) {
		t.Errorf("expected TRIGGER_TYPE_INSTEAD (%d), got %d", nodes.TRIGGER_TYPE_INSTEAD, stmt.Timing)
	}
}

// TestCreateTriggerWithWhen tests: CREATE TRIGGER trig BEFORE INSERT ON t FOR EACH ROW WHEN (NEW.active) EXECUTE FUNCTION func()
func TestCreateTriggerWithWhen(t *testing.T) {
	input := "CREATE TRIGGER trig BEFORE INSERT ON t FOR EACH ROW WHEN (NEW.active) EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if stmt.WhenClause == nil {
		t.Error("expected WhenClause to be set")
	}
	if !stmt.Row {
		t.Error("expected Row to be true")
	}
}

// TestCreateOrReplaceTrigger tests: CREATE OR REPLACE TRIGGER trig BEFORE INSERT ON t EXECUTE FUNCTION func()
func TestCreateOrReplaceTrigger(t *testing.T) {
	input := "CREATE OR REPLACE TRIGGER trig BEFORE INSERT ON t EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if !stmt.Replace {
		t.Error("expected Replace to be true")
	}
}

// TestCreateConstraintTrigger tests: CREATE CONSTRAINT TRIGGER trig AFTER INSERT ON t DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION func()
func TestCreateConstraintTrigger(t *testing.T) {
	input := "CREATE CONSTRAINT TRIGGER trig AFTER INSERT ON t DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if !stmt.IsConstraint {
		t.Error("expected IsConstraint to be true")
	}
	if !stmt.Deferrable {
		t.Error("expected Deferrable to be true")
	}
	if !stmt.Initdeferred {
		t.Error("expected Initdeferred to be true")
	}
	if !stmt.Row {
		t.Error("expected Row to be true")
	}
}

// TestCreateTriggerWithProcedure tests: CREATE TRIGGER trig BEFORE INSERT ON t EXECUTE PROCEDURE func()
func TestCreateTriggerWithProcedure(t *testing.T) {
	input := "CREATE TRIGGER trig BEFORE INSERT ON t EXECUTE PROCEDURE func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateTrigStmt)
	if stmt.Trigname != "trig" {
		t.Errorf("expected trigname 'trig', got %q", stmt.Trigname)
	}
}

// =============================================================================
// CREATE EVENT TRIGGER tests
// =============================================================================

// TestCreateEventTrigger tests: CREATE EVENT TRIGGER etrig ON ddl_command_start EXECUTE FUNCTION func()
func TestCreateEventTrigger(t *testing.T) {
	input := "CREATE EVENT TRIGGER etrig ON ddl_command_start EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreateEventTrigStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateEventTrigStmt, got %T", result.Items[0])
	}
	if stmt.Trigname != "etrig" {
		t.Errorf("expected trigname 'etrig', got %q", stmt.Trigname)
	}
	if stmt.Eventname != "ddl_command_start" {
		t.Errorf("expected eventname 'ddl_command_start', got %q", stmt.Eventname)
	}
	if stmt.Whenclause != nil {
		t.Error("expected Whenclause to be nil")
	}
}

// TestCreateEventTriggerWithWhen tests: CREATE EVENT TRIGGER etrig ON ddl_command_start WHEN TAG IN ('CREATE TABLE', 'DROP TABLE') EXECUTE FUNCTION func()
func TestCreateEventTriggerWithWhen(t *testing.T) {
	input := "CREATE EVENT TRIGGER etrig ON ddl_command_start WHEN TAG IN ('CREATE TABLE', 'DROP TABLE') EXECUTE FUNCTION func()"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateEventTrigStmt)
	if stmt.Whenclause == nil || len(stmt.Whenclause.Items) != 1 {
		t.Fatal("expected 1 when clause item")
	}
	de := stmt.Whenclause.Items[0].(*nodes.DefElem)
	// ColId downcases identifiers - TAG becomes "tag"
	if de.Defname != "tag" {
		t.Errorf("expected defname 'tag', got %q", de.Defname)
	}
	vals := de.Arg.(*nodes.List)
	if len(vals.Items) != 2 {
		t.Fatalf("expected 2 tag values, got %d", len(vals.Items))
	}
}

// =============================================================================
// ALTER EVENT TRIGGER tests
// =============================================================================

// TestAlterEventTriggerEnable tests: ALTER EVENT TRIGGER etrig ENABLE
func TestAlterEventTriggerEnable(t *testing.T) {
	input := "ALTER EVENT TRIGGER etrig ENABLE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterEventTrigStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterEventTrigStmt, got %T", result.Items[0])
	}
	if stmt.Trigname != "etrig" {
		t.Errorf("expected trigname 'etrig', got %q", stmt.Trigname)
	}
	if stmt.Tgenabled != byte(nodes.TRIGGER_FIRES_ON_ORIGIN) {
		t.Errorf("expected TRIGGER_FIRES_ON_ORIGIN (%d), got %d", nodes.TRIGGER_FIRES_ON_ORIGIN, stmt.Tgenabled)
	}
}

// TestAlterEventTriggerDisable tests: ALTER EVENT TRIGGER etrig DISABLE
func TestAlterEventTriggerDisable(t *testing.T) {
	input := "ALTER EVENT TRIGGER etrig DISABLE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterEventTrigStmt)
	if stmt.Tgenabled != byte(nodes.TRIGGER_DISABLED) {
		t.Errorf("expected TRIGGER_DISABLED (%d), got %d", nodes.TRIGGER_DISABLED, stmt.Tgenabled)
	}
}
