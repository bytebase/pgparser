package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE RULE tests
// =============================================================================

// TestCreateRuleDoNothing tests: CREATE RULE r AS ON INSERT TO t DO NOTHING
func TestCreateRuleDoNothing(t *testing.T) {
	input := "CREATE RULE r AS ON INSERT TO t DO NOTHING"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.RuleStmt)
	if !ok {
		t.Fatalf("expected *nodes.RuleStmt, got %T", result.Items[0])
	}
	if stmt.Rulename != "r" {
		t.Errorf("expected rulename 'r', got %q", stmt.Rulename)
	}
	if stmt.Event != nodes.CMD_INSERT {
		t.Errorf("expected CMD_INSERT, got %d", stmt.Event)
	}
	if stmt.Relation == nil || stmt.Relation.Relname != "t" {
		t.Error("expected relation 't'")
	}
	if stmt.Instead {
		t.Error("expected Instead to be false")
	}
	if stmt.Actions != nil && len(stmt.Actions.Items) > 0 {
		t.Error("expected empty actions for DO NOTHING")
	}
}

// TestCreateRuleDoInsteadNothing tests: CREATE RULE r AS ON INSERT TO t DO INSTEAD NOTHING
func TestCreateRuleDoInsteadNothing(t *testing.T) {
	input := "CREATE RULE r AS ON INSERT TO t DO INSTEAD NOTHING"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RuleStmt)
	if !stmt.Instead {
		t.Error("expected Instead to be true")
	}
}

// TestCreateRuleDoAlsoInsert tests: CREATE RULE r AS ON INSERT TO t DO ALSO INSERT INTO log DEFAULT VALUES
func TestCreateRuleDoAlsoInsert(t *testing.T) {
	input := "CREATE RULE r AS ON INSERT TO t DO ALSO INSERT INTO log DEFAULT VALUES"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RuleStmt)
	if stmt.Instead {
		t.Error("expected Instead to be false for ALSO")
	}
	if stmt.Actions == nil || len(stmt.Actions.Items) != 1 {
		t.Fatal("expected 1 action")
	}
	_, ok := stmt.Actions.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Fatalf("expected InsertStmt action, got %T", stmt.Actions.Items[0])
	}
}

// TestCreateRuleMultipleActions tests: CREATE RULE r AS ON DELETE TO t DO (INSERT INTO log DEFAULT VALUES)
func TestCreateRuleMultipleActions(t *testing.T) {
	input := "CREATE RULE r AS ON DELETE TO t DO (INSERT INTO log DEFAULT VALUES)"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RuleStmt)
	if stmt.Event != nodes.CMD_DELETE {
		t.Errorf("expected CMD_DELETE, got %d", stmt.Event)
	}
	if stmt.Actions == nil || len(stmt.Actions.Items) != 1 {
		t.Fatal("expected 1 action")
	}
}

// TestCreateOrReplaceRule tests: CREATE OR REPLACE RULE r AS ON UPDATE TO t WHERE OLD.id > 0 DO NOTHING
func TestCreateOrReplaceRule(t *testing.T) {
	input := "CREATE OR REPLACE RULE r AS ON UPDATE TO t WHERE OLD.id > 0 DO NOTHING"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RuleStmt)
	if !stmt.Replace {
		t.Error("expected Replace to be true")
	}
	if stmt.WhereClause == nil {
		t.Error("expected WhereClause to be set")
	}
	if stmt.Event != nodes.CMD_UPDATE {
		t.Errorf("expected CMD_UPDATE, got %d", stmt.Event)
	}
}

// TestCreateRuleOnSelect tests: CREATE RULE r AS ON SELECT TO v DO INSTEAD SELECT * FROM t
func TestCreateRuleOnSelect(t *testing.T) {
	input := "CREATE RULE r AS ON SELECT TO v DO INSTEAD SELECT * FROM t"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.RuleStmt)
	if stmt.Event != nodes.CMD_SELECT {
		t.Errorf("expected CMD_SELECT, got %d", stmt.Event)
	}
	if !stmt.Instead {
		t.Error("expected Instead to be true")
	}
	if stmt.Actions == nil || len(stmt.Actions.Items) != 1 {
		t.Fatal("expected 1 action")
	}
}

// =============================================================================
// CREATE LANGUAGE tests
// =============================================================================

// TestCreateLanguageSimple tests: CREATE LANGUAGE plpgsql
func TestCreateLanguageSimple(t *testing.T) {
	input := "CREATE LANGUAGE plpgsql"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.CreatePLangStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreatePLangStmt, got %T", result.Items[0])
	}
	if stmt.Plname != "plpgsql" {
		t.Errorf("expected plname 'plpgsql', got %q", stmt.Plname)
	}
}

// TestCreateTrustedLanguageWithHandler tests: CREATE TRUSTED LANGUAGE plpgsql HANDLER plpgsql_call_handler INLINE plpgsql_inline_handler VALIDATOR plpgsql_validator
func TestCreateTrustedLanguageWithHandler(t *testing.T) {
	input := "CREATE TRUSTED LANGUAGE plpgsql HANDLER plpgsql_call_handler INLINE plpgsql_inline_handler VALIDATOR plpgsql_validator"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePLangStmt)
	if !stmt.Pltrusted {
		t.Error("expected Pltrusted to be true")
	}
	if stmt.Plhandler == nil || len(stmt.Plhandler.Items) != 1 {
		t.Fatal("expected handler name")
	}
	if stmt.Plhandler.Items[0].(*nodes.String).Str != "plpgsql_call_handler" {
		t.Errorf("expected handler 'plpgsql_call_handler', got %q", stmt.Plhandler.Items[0].(*nodes.String).Str)
	}
	if stmt.Plinline == nil || len(stmt.Plinline.Items) != 1 {
		t.Fatal("expected inline handler name")
	}
	if stmt.Plvalidator == nil || len(stmt.Plvalidator.Items) != 1 {
		t.Fatal("expected validator name")
	}
}

// TestCreateOrReplaceLanguage tests: CREATE OR REPLACE LANGUAGE plpgsql HANDLER plpgsql_call_handler
func TestCreateOrReplaceLanguage(t *testing.T) {
	input := "CREATE OR REPLACE LANGUAGE plpgsql HANDLER plpgsql_call_handler"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreatePLangStmt)
	if !stmt.Replace {
		t.Error("expected Replace to be true")
	}
	if stmt.Plhandler == nil {
		t.Fatal("expected handler to be set")
	}
}
