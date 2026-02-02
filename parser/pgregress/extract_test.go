package pgregress

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtract_SimpleStatements(t *testing.T) {
	stmts := ExtractStatements("test.sql", []byte("SELECT 1;\nSELECT 2;"))
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if stmts[0].SQL != "SELECT 1" {
		t.Errorf("stmt[0] = %q, want %q", stmts[0].SQL, "SELECT 1")
	}
	if stmts[1].SQL != "SELECT 2" {
		t.Errorf("stmt[1] = %q, want %q", stmts[1].SQL, "SELECT 2")
	}
}

func TestExtract_MultiLineStatement(t *testing.T) {
	input := "SELECT *\n  FROM users\n  WHERE id = 1;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	if stmts[0].StartLine != 1 {
		t.Errorf("StartLine = %d, want 1", stmts[0].StartLine)
	}
}

func TestExtract_LineComment(t *testing.T) {
	input := "-- this is a comment\nSELECT 1;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_BlockComment(t *testing.T) {
	input := "SELECT /* comment */ 1;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_NestedBlockComment(t *testing.T) {
	input := "SELECT /* outer /* inner */ still comment */ 1;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_SingleQuotedString(t *testing.T) {
	input := "SELECT 'hello; world';"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d: %v", len(stmts), stmts)
	}
}

func TestExtract_EscapedQuote(t *testing.T) {
	input := "SELECT 'it''s';"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_EString(t *testing.T) {
	input := `SELECT E'hello\nworld';`
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_EStringWithSemicolon(t *testing.T) {
	input := `SELECT E'semi\;colon';`
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_DoubleQuotedIdentifier(t *testing.T) {
	input := `SELECT "col;name" FROM t;`
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_DollarQuote(t *testing.T) {
	input := "CREATE FUNCTION f() RETURNS void AS $$ SELECT 1; $$ LANGUAGE sql;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d: %v", len(stmts), stmts)
	}
}

func TestExtract_DollarQuoteWithTag(t *testing.T) {
	input := "CREATE FUNCTION f() AS $body$ SELECT 1; $body$ LANGUAGE sql;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_NestedDollarQuote(t *testing.T) {
	input := "DO $$ BEGIN EXECUTE $q$ SELECT 1; $q$; END $$;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_MultiLineDollarQuote(t *testing.T) {
	input := `CREATE FUNCTION f() RETURNS trigger AS $$
begin
    if count(*) = 0 then
        raise exception 'error';
    end if;
    return new;
end;
$$ LANGUAGE plpgsql;`
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

func TestExtract_PsqlMetacommand(t *testing.T) {
	input := "\\set foo bar\nSELECT 1;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	if stmts[0].SQL != "SELECT 1" {
		t.Errorf("stmt = %q, want %q", stmts[0].SQL, "SELECT 1")
	}
}

func TestExtract_PsqlGsetMidLine(t *testing.T) {
	input := "SELECT 1 AS x \\gset\nSELECT 2;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d: %v", len(stmts), stmtsSQL(stmts))
	}
}

func TestExtract_PsqlGsetOwnLine(t *testing.T) {
	input := "SELECT 1 AS x\n\\gset\nSELECT 2;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d: %v", len(stmts), stmtsSQL(stmts))
	}
}

func TestExtract_CopyFromStdin(t *testing.T) {
	input := "COPY t FROM STDIN;\n1\t2\n3\t4\n\\.\nSELECT 1;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d: %v", len(stmts), stmtsSQL(stmts))
	}
	if stmts[0].SQL != "COPY t FROM STDIN" {
		t.Errorf("stmt[0] = %q", stmts[0].SQL)
	}
	if stmts[1].SQL != "SELECT 1" {
		t.Errorf("stmt[1] = %q", stmts[1].SQL)
	}
}

func TestExtract_CopyFromStdinWithOptions(t *testing.T) {
	input := "COPY t FROM STDIN WITH (FORMAT csv);\na,b\nc,d\n\\.\nSELECT 2;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d: %v", len(stmts), stmtsSQL(stmts))
	}
}

func TestExtract_EmptyStatements(t *testing.T) {
	input := "SELECT 1;;\nSELECT 2;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
}

func TestExtract_StartLineTracking(t *testing.T) {
	input := "-- comment\n\nSELECT 1;\n\nSELECT\n  2;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if stmts[0].StartLine != 1 {
		t.Errorf("stmt[0].StartLine = %d, want 1", stmts[0].StartLine)
	}
	if stmts[1].StartLine != 5 {
		t.Errorf("stmt[1].StartLine = %d, want 5", stmts[1].StartLine)
	}
}

func TestExtract_OnlyComments(t *testing.T) {
	input := "-- just a comment\n-- another comment"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 0 {
		t.Fatalf("expected 0 statements, got %d: %v", len(stmts), stmtsSQL(stmts))
	}
}

func TestExtract_DollarSignParam(t *testing.T) {
	// $1 is a parameter, not a dollar-quote
	input := "PREPARE p AS SELECT $1;"
	stmts := ExtractStatements("test.sql", []byte(input))
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
}

// Smoke tests on real PG regression files
func TestExtract_RealSelectFile(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "sql", "select.sql"))
	if err != nil {
		t.Skip("testdata not available:", err)
	}
	stmts := ExtractStatements("select.sql", content)
	if len(stmts) == 0 {
		t.Fatal("extracted 0 statements from select.sql")
	}
	t.Logf("select.sql: %d statements extracted", len(stmts))
}

func TestExtract_RealCopyFile(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "sql", "copy.sql"))
	if err != nil {
		t.Skip("testdata not available:", err)
	}
	stmts := ExtractStatements("copy.sql", content)
	if len(stmts) == 0 {
		t.Fatal("extracted 0 statements from copy.sql")
	}
	// Verify no COPY data lines leaked through as statements
	for i, s := range stmts {
		upper := stmtUpper(s.SQL)
		// Data lines would be like "1\t2" which is not valid SQL keyword start
		if len(upper) > 0 && upper[0] >= '0' && upper[0] <= '9' {
			// Could be a legitimate expression, but flag for review
			t.Logf("  stmt[%d] line %d starts with digit: %.80s", i, s.StartLine, s.SQL)
		}
	}
	t.Logf("copy.sql: %d statements extracted", len(stmts))
}

func TestExtract_RealPlpgsqlFile(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "sql", "plpgsql.sql"))
	if err != nil {
		t.Skip("testdata not available:", err)
	}
	stmts := ExtractStatements("plpgsql.sql", content)
	if len(stmts) == 0 {
		t.Fatal("extracted 0 statements from plpgsql.sql")
	}
	// Verify dollar-quoted function bodies are not split
	for i, s := range stmts {
		if containsDollarQuoteStart(s.SQL) && !containsDollarQuoteEnd(s.SQL) {
			t.Errorf("stmt[%d] line %d has unmatched dollar quote: %.200s", i, s.StartLine, s.SQL)
		}
	}
	t.Logf("plpgsql.sql: %d statements extracted", len(stmts))
}

// Helper functions

func stmtsSQL(stmts []ExtractedStmt) []string {
	result := make([]string, len(stmts))
	for i, s := range stmts {
		if len(s.SQL) > 80 {
			result[i] = s.SQL[:80] + "..."
		} else {
			result[i] = s.SQL
		}
	}
	return result
}

func stmtUpper(s string) string {
	if len(s) > 100 {
		s = s[:100]
	}
	r := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			r[i] = c - 32
		} else {
			r[i] = c
		}
	}
	return string(r)
}

func containsDollarQuoteStart(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '$' {
			tag := tryParseDollarTagStr(s, i)
			if tag != "" {
				return true
			}
		}
	}
	return false
}

func containsDollarQuoteEnd(s string) bool {
	// Simple check: count dollar tags - they should be even
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '$' {
			tag := tryParseDollarTagStr(s, i)
			if tag != "" {
				count++
				i += len(tag) - 1
			}
		}
	}
	return count%2 == 0
}
