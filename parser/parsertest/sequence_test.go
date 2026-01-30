package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// TestCreateSequenceBasic tests: CREATE SEQUENCE seq
func TestCreateSequenceBasic(t *testing.T) {
	input := "CREATE SEQUENCE seq"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Sequence == nil {
		t.Fatal("expected Sequence to be set")
	}
	if stmt.Sequence.Relname != "seq" {
		t.Errorf("expected sequence name 'seq', got %q", stmt.Sequence.Relname)
	}
	if stmt.IfNotExists {
		t.Error("expected IfNotExists to be false")
	}
}

// TestCreateSequenceIfNotExists tests: CREATE SEQUENCE IF NOT EXISTS seq
func TestCreateSequenceIfNotExists(t *testing.T) {
	input := "CREATE SEQUENCE IF NOT EXISTS seq"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
	if stmt.Sequence.Relname != "seq" {
		t.Errorf("expected sequence name 'seq', got %q", stmt.Sequence.Relname)
	}
}

// TestCreateTempSequence tests: CREATE TEMP SEQUENCE seq
func TestCreateTempSequence(t *testing.T) {
	input := "CREATE TEMP SEQUENCE seq"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Sequence == nil {
		t.Fatal("expected Sequence to be set")
	}
	if stmt.Sequence.Relpersistence != 't' {
		t.Errorf("expected relpersistence 't', got %q", stmt.Sequence.Relpersistence)
	}
}

// TestCreateSequenceIncrement tests: CREATE SEQUENCE seq INCREMENT 5
func TestCreateSequenceIncrement(t *testing.T) {
	input := "CREATE SEQUENCE seq INCREMENT 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %v", stmt.Options)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "increment" {
		t.Errorf("expected option name 'increment', got %q", opt.Defname)
	}
}

// TestCreateSequenceIncrementByStart tests: CREATE SEQUENCE seq INCREMENT BY 5 START WITH 100
func TestCreateSequenceIncrementByStart(t *testing.T) {
	input := "CREATE SEQUENCE seq INCREMENT BY 5 START WITH 100"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 2 {
		t.Fatalf("expected 2 options, got %d", stmt.Options.Len())
	}

	opt1 := stmt.Options.Items[0].(*nodes.DefElem)
	if opt1.Defname != "increment" {
		t.Errorf("expected option 'increment', got %q", opt1.Defname)
	}

	opt2 := stmt.Options.Items[1].(*nodes.DefElem)
	if opt2.Defname != "start" {
		t.Errorf("expected option 'start', got %q", opt2.Defname)
	}
}

// TestCreateSequenceMinMaxValue tests: CREATE SEQUENCE seq MINVALUE 1 MAXVALUE 1000
func TestCreateSequenceMinMaxValue(t *testing.T) {
	input := "CREATE SEQUENCE seq MINVALUE 1 MAXVALUE 1000"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 2 {
		t.Fatalf("expected 2 options, got %d", stmt.Options.Len())
	}

	opt1 := stmt.Options.Items[0].(*nodes.DefElem)
	if opt1.Defname != "minvalue" {
		t.Errorf("expected 'minvalue', got %q", opt1.Defname)
	}

	opt2 := stmt.Options.Items[1].(*nodes.DefElem)
	if opt2.Defname != "maxvalue" {
		t.Errorf("expected 'maxvalue', got %q", opt2.Defname)
	}
}

// TestCreateSequenceNoMinMaxValue tests: CREATE SEQUENCE seq NO MINVALUE NO MAXVALUE
func TestCreateSequenceNoMinMaxValue(t *testing.T) {
	input := "CREATE SEQUENCE seq NO MINVALUE NO MAXVALUE"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 2 {
		t.Fatalf("expected 2 options, got %d", stmt.Options.Len())
	}

	opt1 := stmt.Options.Items[0].(*nodes.DefElem)
	if opt1.Defname != "minvalue" {
		t.Errorf("expected 'minvalue', got %q", opt1.Defname)
	}
	if opt1.Arg != nil {
		t.Error("expected nil Arg for NO MINVALUE")
	}

	opt2 := stmt.Options.Items[1].(*nodes.DefElem)
	if opt2.Defname != "maxvalue" {
		t.Errorf("expected 'maxvalue', got %q", opt2.Defname)
	}
	if opt2.Arg != nil {
		t.Error("expected nil Arg for NO MAXVALUE")
	}
}

// TestCreateSequenceCycle tests: CREATE SEQUENCE seq CYCLE
func TestCreateSequenceCycle(t *testing.T) {
	input := "CREATE SEQUENCE seq CYCLE"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}

	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "cycle" {
		t.Errorf("expected 'cycle', got %q", opt.Defname)
	}
}

// TestCreateSequenceCache tests: CREATE SEQUENCE seq CACHE 20
func TestCreateSequenceCache(t *testing.T) {
	input := "CREATE SEQUENCE seq CACHE 20"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}

	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "cache" {
		t.Errorf("expected 'cache', got %q", opt.Defname)
	}
}

// TestCreateSequenceOwnedBy tests: CREATE SEQUENCE seq OWNED BY t.id
func TestCreateSequenceOwnedBy(t *testing.T) {
	input := "CREATE SEQUENCE seq OWNED BY t.id"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}

	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "owned_by" {
		t.Errorf("expected 'owned_by', got %q", opt.Defname)
	}
}

// TestAlterSequenceRestart tests: ALTER SEQUENCE seq RESTART
func TestAlterSequenceRestart(t *testing.T) {
	input := "ALTER SEQUENCE seq RESTART"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterSeqStmt, got %T", result.Items[0])
	}

	if stmt.Sequence == nil || stmt.Sequence.Relname != "seq" {
		t.Fatalf("expected sequence name 'seq'")
	}
	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "restart" {
		t.Errorf("expected 'restart', got %q", opt.Defname)
	}
}

// TestAlterSequenceRestartWith tests: ALTER SEQUENCE seq RESTART WITH 100
func TestAlterSequenceRestartWith(t *testing.T) {
	input := "ALTER SEQUENCE seq RESTART WITH 100"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "restart" {
		t.Errorf("expected 'restart', got %q", opt.Defname)
	}
	if opt.Arg == nil {
		t.Fatal("expected non-nil restart value")
	}
}

// TestAlterSequenceIncrement tests: ALTER SEQUENCE seq INCREMENT BY 10
func TestAlterSequenceIncrement(t *testing.T) {
	input := "ALTER SEQUENCE seq INCREMENT BY 10"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterSeqStmt, got %T", result.Items[0])
	}

	if stmt.Options == nil || len(stmt.Options.Items) != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "increment" {
		t.Errorf("expected 'increment', got %q", opt.Defname)
	}
}

// TestAlterSequenceIfExists tests: ALTER SEQUENCE IF EXISTS seq INCREMENT 5
func TestAlterSequenceIfExists(t *testing.T) {
	input := "ALTER SEQUENCE IF EXISTS seq INCREMENT 5"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.AlterSeqStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterSeqStmt, got %T", result.Items[0])
	}

	if !stmt.MissingOk {
		t.Error("expected MissingOk to be true")
	}
	if stmt.Sequence.Relname != "seq" {
		t.Errorf("expected sequence name 'seq', got %q", stmt.Sequence.Relname)
	}
}
