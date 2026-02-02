package pgregress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/pgparser/pgparser/parser"
)

// TestShowFailures2 reads known_failures.json and prints ALL failing SQL
// statements for a curated set of files, with SQL truncated to 300 chars.
func TestShowFailures2(t *testing.T) {
	// Target files to inspect
	targetFiles := []string{
		"update.sql",
		"merge.sql",
		"foreign_data.sql",
		"xml.sql",
		"sqljson.sql",
		"create_am.sql",
		"sqljson_jsontable.sql",
		"strings.sql",
		"stats_ext.sql",
		"copy.sql",
		"tsearch.sql",
		"rules.sql",
		"identity.sql",
	}

	// Load known_failures.json
	data, err := os.ReadFile("known_failures.json")
	if err != nil {
		t.Fatalf("read known_failures.json: %v", err)
	}
	var kf KnownFailures
	if err := json.Unmarshal(data, &kf); err != nil {
		t.Fatalf("parse known_failures.json: %v", err)
	}

	sort.Strings(targetFiles)

	totalFailures := 0
	totalStillFailing := 0

	for _, base := range targetFiles {
		indices, ok := kf[base]
		if !ok {
			t.Logf("\n=== %s: NOT IN known_failures.json ===\n", base)
			continue
		}

		filePath := filepath.Join("testdata", "sql", base)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Logf("\n=== %s: CANNOT READ FILE: %v ===\n", base, err)
			continue
		}

		stmts := ExtractStatements(base, content)

		t.Logf("\n========================================================================")
		t.Logf("=== %s: %d known failures, %d total statements ===", base, len(indices), len(stmts))
		t.Logf("========================================================================")

		fileStillFailing := 0

		for _, idx := range indices {
			totalFailures++

			if idx >= len(stmts) {
				t.Logf("  [idx=%d] OUT OF RANGE (only %d statements extracted)", idx, len(stmts))
				continue
			}

			stmt := stmts[idx]
			_, parseErr := parser.Parse(stmt.SQL)

			status := "NOW PASSES"
			errStr := ""
			if parseErr != nil {
				status = "FAILS"
				errStr = parseErr.Error()
				fileStillFailing++
				totalStillFailing++
			}

			// Truncate SQL to 300 chars (collapse whitespace)
			sqlSnippet := strings.Join(strings.Fields(stmt.SQL), " ")
			if len(sqlSnippet) > 300 {
				sqlSnippet = sqlSnippet[:300] + "..."
			}

			t.Logf("")
			t.Logf("  [idx=%d line=%d] %s", idx, stmt.StartLine, status)
			if errStr != "" {
				t.Logf("    Error: %s", errStr)
			}
			t.Logf("    SQL: %s", sqlSnippet)
		}

		t.Logf("")
		t.Logf("  --- %s summary: %d/%d still failing ---", base, fileStillFailing, len(indices))
	}

	t.Logf("\n========================================================================")
	t.Logf("=== GRAND TOTAL: %d/%d entries still failing across %d files ===",
		totalStillFailing, totalFailures, len(targetFiles))
	t.Logf("========================================================================")

	// Print a compact per-file summary table at the end
	t.Logf("\n--- Per-file summary ---")
	t.Logf("%-30s %8s %8s", "File", "Known", "Status")
	for _, base := range targetFiles {
		indices, ok := kf[base]
		if !ok {
			t.Logf("%-30s %8s %8s", base, "-", "not in json")
			continue
		}

		filePath := filepath.Join("testdata", "sql", base)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Logf("%-30s %8d %8s", base, len(indices), "read error")
			continue
		}

		stmts := ExtractStatements(base, content)
		stillFailing := 0
		for _, idx := range indices {
			if idx < len(stmts) {
				if _, parseErr := parser.Parse(stmts[idx].SQL); parseErr != nil {
					stillFailing++
				}
			}
		}

		t.Logf("%-30s %8d %s",
			base, len(indices),
			fmt.Sprintf("%d still fail, %d now pass", stillFailing, len(indices)-stillFailing))
	}
}
