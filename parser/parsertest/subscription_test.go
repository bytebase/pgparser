package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// CREATE SUBSCRIPTION tests
// =============================================================================

// TestCreateSubscriptionBasic tests: CREATE SUBSCRIPTION sub CONNECTION 'host=localhost' PUBLICATION pub1
func TestCreateSubscriptionBasic(t *testing.T) {
	input := "CREATE SUBSCRIPTION sub CONNECTION 'host=localhost' PUBLICATION pub1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}
	stmt, ok := result.Items[0].(*nodes.CreateSubscriptionStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateSubscriptionStmt, got %T", result.Items[0])
	}
	if stmt.Subname != "sub" {
		t.Errorf("expected subname 'sub', got %q", stmt.Subname)
	}
	if stmt.Conninfo != "host=localhost" {
		t.Errorf("expected conninfo 'host=localhost', got %q", stmt.Conninfo)
	}
	if stmt.Publication == nil || stmt.Publication.Len() != 1 {
		t.Fatalf("expected 1 publication, got %d", stmt.Publication.Len())
	}
	pubName := stmt.Publication.Items[0].(*nodes.String).Str
	if pubName != "pub1" {
		t.Errorf("expected publication 'pub1', got %q", pubName)
	}
}

// TestCreateSubscriptionMultiplePubs tests: CREATE SUBSCRIPTION sub CONNECTION 'host=localhost' PUBLICATION pub1, pub2
func TestCreateSubscriptionMultiplePubs(t *testing.T) {
	input := "CREATE SUBSCRIPTION sub CONNECTION 'host=localhost' PUBLICATION pub1, pub2"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateSubscriptionStmt)
	if stmt.Publication == nil || stmt.Publication.Len() != 2 {
		t.Fatalf("expected 2 publications, got %d", stmt.Publication.Len())
	}
}

// TestCreateSubscriptionWithOptions tests: CREATE SUBSCRIPTION sub CONNECTION 'connstr' PUBLICATION pub WITH (slot_name = 'none')
func TestCreateSubscriptionWithOptions(t *testing.T) {
	input := "CREATE SUBSCRIPTION sub CONNECTION 'connstr' PUBLICATION pub WITH (slot_name = 'none')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.CreateSubscriptionStmt)
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatalf("expected 1 option, got %d", stmt.Options.Len())
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "slot_name" {
		t.Errorf("expected defname 'slot_name', got %q", opt.Defname)
	}
}

// =============================================================================
// ALTER SUBSCRIPTION tests
// =============================================================================

// TestAlterSubscriptionConnection tests: ALTER SUBSCRIPTION sub CONNECTION 'new_conn'
func TestAlterSubscriptionConnection(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub CONNECTION 'new_conn'"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if !ok {
		t.Fatalf("expected *nodes.AlterSubscriptionStmt, got %T", result.Items[0])
	}
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_CONNECTION {
		t.Errorf("expected ALTER_SUBSCRIPTION_CONNECTION, got %d", stmt.Kind)
	}
	if stmt.Subname != "sub" {
		t.Errorf("expected subname 'sub', got %q", stmt.Subname)
	}
	if stmt.Conninfo != "new_conn" {
		t.Errorf("expected conninfo 'new_conn', got %q", stmt.Conninfo)
	}
}

// TestAlterSubscriptionSetPublication tests: ALTER SUBSCRIPTION sub SET PUBLICATION pub2
func TestAlterSubscriptionSetPublication(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub SET PUBLICATION pub2"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_SET_PUBLICATION {
		t.Errorf("expected ALTER_SUBSCRIPTION_SET_PUBLICATION, got %d", stmt.Kind)
	}
	if stmt.Publication == nil || stmt.Publication.Len() != 1 {
		t.Fatal("expected 1 publication")
	}
}

// TestAlterSubscriptionRefresh tests: ALTER SUBSCRIPTION sub REFRESH PUBLICATION
func TestAlterSubscriptionRefresh(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub REFRESH PUBLICATION"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_REFRESH {
		t.Errorf("expected ALTER_SUBSCRIPTION_REFRESH, got %d", stmt.Kind)
	}
}

// TestAlterSubscriptionEnable tests: ALTER SUBSCRIPTION sub ENABLE
func TestAlterSubscriptionEnable(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub ENABLE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_ENABLED {
		t.Errorf("expected ALTER_SUBSCRIPTION_ENABLED, got %d", stmt.Kind)
	}
	// Check that the option is "enabled" = true
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatal("expected 1 option")
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	if opt.Defname != "enabled" {
		t.Errorf("expected defname 'enabled', got %q", opt.Defname)
	}
	bval, ok := opt.Arg.(*nodes.Boolean)
	if !ok || !bval.Boolval {
		t.Error("expected enabled = true")
	}
}

// TestAlterSubscriptionDisable tests: ALTER SUBSCRIPTION sub DISABLE
func TestAlterSubscriptionDisable(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub DISABLE"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_ENABLED {
		t.Errorf("expected ALTER_SUBSCRIPTION_ENABLED, got %d", stmt.Kind)
	}
	opt := stmt.Options.Items[0].(*nodes.DefElem)
	bval := opt.Arg.(*nodes.Boolean)
	if bval.Boolval {
		t.Error("expected enabled = false")
	}
}

// TestAlterSubscriptionSetOptions tests: ALTER SUBSCRIPTION sub SET (slot_name = 'myslot')
func TestAlterSubscriptionSetOptions(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub SET (slot_name = 'myslot')"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_OPTIONS {
		t.Errorf("expected ALTER_SUBSCRIPTION_OPTIONS, got %d", stmt.Kind)
	}
	if stmt.Options == nil || stmt.Options.Len() != 1 {
		t.Fatal("expected 1 option")
	}
}

// TestAlterSubscriptionAddPublication tests: ALTER SUBSCRIPTION sub ADD PUBLICATION pub2
func TestAlterSubscriptionAddPublication(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub ADD PUBLICATION pub2"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_ADD_PUBLICATION {
		t.Errorf("expected ALTER_SUBSCRIPTION_ADD_PUBLICATION, got %d", stmt.Kind)
	}
}

// TestAlterSubscriptionDropPublication tests: ALTER SUBSCRIPTION sub DROP PUBLICATION pub1
func TestAlterSubscriptionDropPublication(t *testing.T) {
	input := "ALTER SUBSCRIPTION sub DROP PUBLICATION pub1"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.AlterSubscriptionStmt)
	if stmt.Kind != nodes.ALTER_SUBSCRIPTION_DROP_PUBLICATION {
		t.Errorf("expected ALTER_SUBSCRIPTION_DROP_PUBLICATION, got %d", stmt.Kind)
	}
}

// =============================================================================
// DROP SUBSCRIPTION tests
// =============================================================================

// TestDropSubscription tests: DROP SUBSCRIPTION sub
func TestDropSubscription(t *testing.T) {
	input := "DROP SUBSCRIPTION sub"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt, ok := result.Items[0].(*nodes.DropSubscriptionStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropSubscriptionStmt, got %T", result.Items[0])
	}
	if stmt.Subname != "sub" {
		t.Errorf("expected subname 'sub', got %q", stmt.Subname)
	}
	if stmt.MissingOk {
		t.Error("expected MissingOk to be false")
	}
}

// TestDropSubscriptionIfExists tests: DROP SUBSCRIPTION IF EXISTS sub
func TestDropSubscriptionIfExists(t *testing.T) {
	input := "DROP SUBSCRIPTION IF EXISTS sub"
	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	stmt := result.Items[0].(*nodes.DropSubscriptionStmt)
	if !stmt.MissingOk {
		t.Error("expected MissingOk to be true")
	}
}
