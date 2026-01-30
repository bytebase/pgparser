package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
)

/*****************************************************************************
 *
 *      DECLARE CURSOR tests
 *
 *****************************************************************************/

func TestDeclareCursorBasic(t *testing.T) {
	stmt := parseStmt(t, "DECLARE c CURSOR FOR SELECT 1")
	dc, ok := stmt.(*nodes.DeclareCursorStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeclareCursorStmt, got %T", stmt)
	}
	if dc.Portalname != "c" {
		t.Errorf("expected Portalname 'c', got %q", dc.Portalname)
	}
	// Should have FAST_PLAN set
	if dc.Options&nodes.CURSOR_OPT_FAST_PLAN == 0 {
		t.Errorf("expected CURSOR_OPT_FAST_PLAN to be set, Options=%d", dc.Options)
	}
	if dc.Query == nil {
		t.Error("expected Query to be non-nil")
	}
}

func TestDeclareCursorBinary(t *testing.T) {
	stmt := parseStmt(t, "DECLARE c BINARY CURSOR FOR SELECT 1")
	dc, ok := stmt.(*nodes.DeclareCursorStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeclareCursorStmt, got %T", stmt)
	}
	if dc.Options&nodes.CURSOR_OPT_BINARY == 0 {
		t.Errorf("expected CURSOR_OPT_BINARY to be set, Options=%d", dc.Options)
	}
}

func TestDeclareCursorScroll(t *testing.T) {
	stmt := parseStmt(t, "DECLARE c SCROLL CURSOR FOR SELECT 1")
	dc, ok := stmt.(*nodes.DeclareCursorStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeclareCursorStmt, got %T", stmt)
	}
	if dc.Options&nodes.CURSOR_OPT_SCROLL == 0 {
		t.Errorf("expected CURSOR_OPT_SCROLL to be set, Options=%d", dc.Options)
	}
}

func TestDeclareCursorNoScroll(t *testing.T) {
	stmt := parseStmt(t, "DECLARE c NO SCROLL CURSOR FOR SELECT 1")
	dc, ok := stmt.(*nodes.DeclareCursorStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeclareCursorStmt, got %T", stmt)
	}
	if dc.Options&nodes.CURSOR_OPT_NO_SCROLL == 0 {
		t.Errorf("expected CURSOR_OPT_NO_SCROLL to be set, Options=%d", dc.Options)
	}
}

func TestDeclareCursorInsensitive(t *testing.T) {
	stmt := parseStmt(t, "DECLARE c INSENSITIVE CURSOR FOR SELECT 1")
	dc, ok := stmt.(*nodes.DeclareCursorStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeclareCursorStmt, got %T", stmt)
	}
	if dc.Options&nodes.CURSOR_OPT_INSENSITIVE == 0 {
		t.Errorf("expected CURSOR_OPT_INSENSITIVE to be set, Options=%d", dc.Options)
	}
}

func TestDeclareCursorWithHold(t *testing.T) {
	stmt := parseStmt(t, "DECLARE c CURSOR WITH HOLD FOR SELECT 1")
	dc, ok := stmt.(*nodes.DeclareCursorStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeclareCursorStmt, got %T", stmt)
	}
	if dc.Options&nodes.CURSOR_OPT_HOLD == 0 {
		t.Errorf("expected CURSOR_OPT_HOLD to be set, Options=%d", dc.Options)
	}
}

func TestDeclareCursorWithoutHold(t *testing.T) {
	stmt := parseStmt(t, "DECLARE c CURSOR WITHOUT HOLD FOR SELECT 1")
	dc, ok := stmt.(*nodes.DeclareCursorStmt)
	if !ok {
		t.Fatalf("expected *nodes.DeclareCursorStmt, got %T", stmt)
	}
	if dc.Options&nodes.CURSOR_OPT_HOLD != 0 {
		t.Errorf("expected CURSOR_OPT_HOLD NOT to be set, Options=%d", dc.Options)
	}
}

/*****************************************************************************
 *
 *      FETCH / MOVE tests
 *
 *****************************************************************************/

func TestFetchCursorName(t *testing.T) {
	stmt := parseStmt(t, "FETCH c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Portalname != "c" {
		t.Errorf("expected Portalname 'c', got %q", fs.Portalname)
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 1 {
		t.Errorf("expected HowMany 1, got %d", fs.HowMany)
	}
	if fs.Ismove {
		t.Error("expected Ismove false")
	}
}

func TestFetchNextFrom(t *testing.T) {
	stmt := parseStmt(t, "FETCH NEXT FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Portalname != "c" {
		t.Errorf("expected Portalname 'c', got %q", fs.Portalname)
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 1 {
		t.Errorf("expected HowMany 1, got %d", fs.HowMany)
	}
}

func TestFetchPriorFrom(t *testing.T) {
	stmt := parseStmt(t, "FETCH PRIOR FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_BACKWARD {
		t.Errorf("expected Direction FETCH_BACKWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 1 {
		t.Errorf("expected HowMany 1, got %d", fs.HowMany)
	}
}

func TestFetchFirstFrom(t *testing.T) {
	stmt := parseStmt(t, "FETCH FIRST FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_ABSOLUTE {
		t.Errorf("expected Direction FETCH_ABSOLUTE, got %d", fs.Direction)
	}
	if fs.HowMany != 1 {
		t.Errorf("expected HowMany 1, got %d", fs.HowMany)
	}
}

func TestFetchLastFrom(t *testing.T) {
	stmt := parseStmt(t, "FETCH LAST FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_ABSOLUTE {
		t.Errorf("expected Direction FETCH_ABSOLUTE, got %d", fs.Direction)
	}
	if fs.HowMany != -1 {
		t.Errorf("expected HowMany -1, got %d", fs.HowMany)
	}
}

func TestFetchAbsolute(t *testing.T) {
	stmt := parseStmt(t, "FETCH ABSOLUTE 5 FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_ABSOLUTE {
		t.Errorf("expected Direction FETCH_ABSOLUTE, got %d", fs.Direction)
	}
	if fs.HowMany != 5 {
		t.Errorf("expected HowMany 5, got %d", fs.HowMany)
	}
}

func TestFetchRelative(t *testing.T) {
	stmt := parseStmt(t, "FETCH RELATIVE -2 FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_RELATIVE {
		t.Errorf("expected Direction FETCH_RELATIVE, got %d", fs.Direction)
	}
	if fs.HowMany != -2 {
		t.Errorf("expected HowMany -2, got %d", fs.HowMany)
	}
}

func TestFetchCount(t *testing.T) {
	stmt := parseStmt(t, "FETCH 10 FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 10 {
		t.Errorf("expected HowMany 10, got %d", fs.HowMany)
	}
}

func TestFetchAll(t *testing.T) {
	stmt := parseStmt(t, "FETCH ALL FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != nodes.FETCH_ALL {
		t.Errorf("expected HowMany FETCH_ALL, got %d", fs.HowMany)
	}
}

func TestFetchForward(t *testing.T) {
	stmt := parseStmt(t, "FETCH FORWARD FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 1 {
		t.Errorf("expected HowMany 1, got %d", fs.HowMany)
	}
}

func TestFetchForwardCount(t *testing.T) {
	stmt := parseStmt(t, "FETCH FORWARD 5 FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 5 {
		t.Errorf("expected HowMany 5, got %d", fs.HowMany)
	}
}

func TestFetchForwardAll(t *testing.T) {
	stmt := parseStmt(t, "FETCH FORWARD ALL FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != nodes.FETCH_ALL {
		t.Errorf("expected HowMany FETCH_ALL, got %d", fs.HowMany)
	}
}

func TestFetchBackward(t *testing.T) {
	stmt := parseStmt(t, "FETCH BACKWARD FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_BACKWARD {
		t.Errorf("expected Direction FETCH_BACKWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 1 {
		t.Errorf("expected HowMany 1, got %d", fs.HowMany)
	}
}

func TestFetchBackwardCount(t *testing.T) {
	stmt := parseStmt(t, "FETCH BACKWARD 5 FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_BACKWARD {
		t.Errorf("expected Direction FETCH_BACKWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 5 {
		t.Errorf("expected HowMany 5, got %d", fs.HowMany)
	}
}

func TestFetchBackwardAll(t *testing.T) {
	stmt := parseStmt(t, "FETCH BACKWARD ALL FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if fs.Direction != nodes.FETCH_BACKWARD {
		t.Errorf("expected Direction FETCH_BACKWARD, got %d", fs.Direction)
	}
	if fs.HowMany != nodes.FETCH_ALL {
		t.Errorf("expected HowMany FETCH_ALL, got %d", fs.HowMany)
	}
}

func TestMoveNextFrom(t *testing.T) {
	stmt := parseStmt(t, "MOVE NEXT FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if !fs.Ismove {
		t.Error("expected Ismove true")
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 1 {
		t.Errorf("expected HowMany 1, got %d", fs.HowMany)
	}
}

func TestMoveForwardCount(t *testing.T) {
	stmt := parseStmt(t, "MOVE FORWARD 10 FROM c")
	fs, ok := stmt.(*nodes.FetchStmt)
	if !ok {
		t.Fatalf("expected *nodes.FetchStmt, got %T", stmt)
	}
	if !fs.Ismove {
		t.Error("expected Ismove true")
	}
	if fs.Direction != nodes.FETCH_FORWARD {
		t.Errorf("expected Direction FETCH_FORWARD, got %d", fs.Direction)
	}
	if fs.HowMany != 10 {
		t.Errorf("expected HowMany 10, got %d", fs.HowMany)
	}
}
