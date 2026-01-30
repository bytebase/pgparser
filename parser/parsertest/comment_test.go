package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/nodes"
)

/*****************************************************************************
 *
 *      COMMENT ON tests
 *
 *****************************************************************************/

func TestCommentOnTable(t *testing.T) {
	stmt := parseStmt(t, "COMMENT ON TABLE t IS 'my comment'")
	cs, ok := stmt.(*nodes.CommentStmt)
	if !ok {
		t.Fatalf("expected *nodes.CommentStmt, got %T", stmt)
	}
	if cs.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", cs.Objtype)
	}
	if cs.Comment != "my comment" {
		t.Errorf("expected comment 'my comment', got %q", cs.Comment)
	}
}

func TestCommentOnTableNull(t *testing.T) {
	stmt := parseStmt(t, "COMMENT ON TABLE t IS NULL")
	cs, ok := stmt.(*nodes.CommentStmt)
	if !ok {
		t.Fatalf("expected *nodes.CommentStmt, got %T", stmt)
	}
	if cs.Comment != "" {
		t.Errorf("expected empty comment (NULL), got %q", cs.Comment)
	}
}

func TestCommentOnColumn(t *testing.T) {
	stmt := parseStmt(t, "COMMENT ON COLUMN t.c IS 'column comment'")
	cs, ok := stmt.(*nodes.CommentStmt)
	if !ok {
		t.Fatalf("expected *nodes.CommentStmt, got %T", stmt)
	}
	if cs.Objtype != nodes.OBJECT_COLUMN {
		t.Errorf("expected OBJECT_COLUMN, got %d", cs.Objtype)
	}
	if cs.Comment != "column comment" {
		t.Errorf("expected comment 'column comment', got %q", cs.Comment)
	}
}

func TestCommentOnView(t *testing.T) {
	stmt := parseStmt(t, "COMMENT ON VIEW v IS 'view comment'")
	cs, ok := stmt.(*nodes.CommentStmt)
	if !ok {
		t.Fatalf("expected *nodes.CommentStmt, got %T", stmt)
	}
	if cs.Objtype != nodes.OBJECT_VIEW {
		t.Errorf("expected OBJECT_VIEW, got %d", cs.Objtype)
	}
}

func TestCommentOnIndex(t *testing.T) {
	stmt := parseStmt(t, "COMMENT ON INDEX idx IS 'index comment'")
	cs, ok := stmt.(*nodes.CommentStmt)
	if !ok {
		t.Fatalf("expected *nodes.CommentStmt, got %T", stmt)
	}
	if cs.Objtype != nodes.OBJECT_INDEX {
		t.Errorf("expected OBJECT_INDEX, got %d", cs.Objtype)
	}
}

func TestCommentOnSchema(t *testing.T) {
	stmt := parseStmt(t, "COMMENT ON SCHEMA myschema IS 'schema comment'")
	cs, ok := stmt.(*nodes.CommentStmt)
	if !ok {
		t.Fatalf("expected *nodes.CommentStmt, got %T", stmt)
	}
	if cs.Objtype != nodes.OBJECT_SCHEMA {
		t.Errorf("expected OBJECT_SCHEMA, got %d", cs.Objtype)
	}
	if cs.Comment != "schema comment" {
		t.Errorf("expected comment 'schema comment', got %q", cs.Comment)
	}
}

func TestCommentOnDatabase(t *testing.T) {
	stmt := parseStmt(t, "COMMENT ON DATABASE mydb IS 'db comment'")
	cs, ok := stmt.(*nodes.CommentStmt)
	if !ok {
		t.Fatalf("expected *nodes.CommentStmt, got %T", stmt)
	}
	if cs.Objtype != nodes.OBJECT_DATABASE {
		t.Errorf("expected OBJECT_DATABASE, got %d", cs.Objtype)
	}
}

/*****************************************************************************
 *
 *      SECURITY LABEL tests
 *
 *****************************************************************************/

func TestSecLabelOnTable(t *testing.T) {
	stmt := parseStmt(t, "SECURITY LABEL ON TABLE t IS 'label'")
	sl, ok := stmt.(*nodes.SecLabelStmt)
	if !ok {
		t.Fatalf("expected *nodes.SecLabelStmt, got %T", stmt)
	}
	if sl.Objtype != nodes.OBJECT_TABLE {
		t.Errorf("expected OBJECT_TABLE, got %d", sl.Objtype)
	}
	if sl.Label != "label" {
		t.Errorf("expected label 'label', got %q", sl.Label)
	}
	if sl.Provider != "" {
		t.Errorf("expected empty provider, got %q", sl.Provider)
	}
}

func TestSecLabelWithProvider(t *testing.T) {
	stmt := parseStmt(t, "SECURITY LABEL FOR selinux ON TABLE t IS 'label'")
	sl, ok := stmt.(*nodes.SecLabelStmt)
	if !ok {
		t.Fatalf("expected *nodes.SecLabelStmt, got %T", stmt)
	}
	if sl.Provider != "selinux" {
		t.Errorf("expected provider 'selinux', got %q", sl.Provider)
	}
}

func TestSecLabelOnColumn(t *testing.T) {
	stmt := parseStmt(t, "SECURITY LABEL ON COLUMN t.c IS 'label'")
	sl, ok := stmt.(*nodes.SecLabelStmt)
	if !ok {
		t.Fatalf("expected *nodes.SecLabelStmt, got %T", stmt)
	}
	if sl.Objtype != nodes.OBJECT_COLUMN {
		t.Errorf("expected OBJECT_COLUMN, got %d", sl.Objtype)
	}
}

func TestSecLabelOnSchema(t *testing.T) {
	stmt := parseStmt(t, "SECURITY LABEL ON SCHEMA myschema IS 'label'")
	sl, ok := stmt.(*nodes.SecLabelStmt)
	if !ok {
		t.Fatalf("expected *nodes.SecLabelStmt, got %T", stmt)
	}
	if sl.Objtype != nodes.OBJECT_SCHEMA {
		t.Errorf("expected OBJECT_SCHEMA, got %d", sl.Objtype)
	}
}
