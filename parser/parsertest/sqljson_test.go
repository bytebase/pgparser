package parsertest

import (
	"testing"

	"github.com/pgparser/pgparser/parser"
)

// =============================================================================
// SQL/JSON bare EMPTY behavior tests
// =============================================================================

func TestJsonTableFormattedColumn(t *testing.T) {
	sqls := []string{
		`SELECT * FROM JSON_TABLE(jsonb 'null', 'lax $[*]' COLUMNS (jst text FORMAT JSON PATH '$'))`,
		`SELECT * FROM JSON_TABLE(jsonb '{"a":"1"}', '$' COLUMNS (a text FORMAT JSON PATH '$.a'))`,
		`SELECT * FROM JSON_TABLE(jsonb '{"a":"1"}', '$' COLUMNS (a text FORMAT JSON PATH '$.a' WITH WRAPPER))`,
	}
	for _, sql := range sqls {
		_, err := parser.Parse(sql)
		if err != nil {
			label := sql
			if len(label) > 60 {
				label = label[:60]
			}
			t.Errorf("Parse(%q) failed: %v", label, err)
		}
	}
}

func TestJsonBareEmpty(t *testing.T) {
	sqls := []string{
		`SELECT * FROM JSON_TABLE('[]', 'strict $.a' COLUMNS (js2 int PATH '$') EMPTY ON ERROR)`,
		`SELECT JSON_QUERY(jsonb '[]', '$[*]' EMPTY ON EMPTY)`,
		`SELECT JSON_VALUE(jsonb '1', '$' EMPTY ON ERROR)`,
	}
	for _, sql := range sqls {
		_, err := parser.Parse(sql)
		if err != nil {
			label := sql
			if len(label) > 60 {
				label = label[:60]
			}
			t.Errorf("Parse(%q) failed: %v", label, err)
		}
	}
}
