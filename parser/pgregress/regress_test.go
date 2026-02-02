package pgregress

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/pgparser/pgparser/parser"
)

var update = flag.Bool("update", false, "regenerate known_failures.json from current parse results")

// KnownFailures maps filename -> sorted list of 0-based statement indices that fail.
type KnownFailures map[string][]int

func loadKnownFailures(t *testing.T) KnownFailures {
	t.Helper()
	data, err := os.ReadFile("known_failures.json")
	if err != nil {
		if os.IsNotExist(err) {
			return KnownFailures{}
		}
		t.Fatalf("read known_failures.json: %v", err)
	}
	var kf KnownFailures
	if err := json.Unmarshal(data, &kf); err != nil {
		t.Fatalf("parse known_failures.json: %v", err)
	}
	return kf
}

func intSliceContains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func TestPGRegress(t *testing.T) {
	files, err := filepath.Glob("testdata/sql/*.sql")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(files) == 0 {
		t.Skip("no regression test files found in testdata/sql/")
	}

	knownFailures := loadKnownFailures(t)

	// Collect new failures if updating
	var newFailures KnownFailures
	if *update {
		newFailures = make(KnownFailures)
	}

	// Aggregate statistics
	var totalStmts, totalPassed, totalKnown, totalNew, totalFixed int

	sort.Strings(files)

	for _, file := range files {
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, ".sql")

		t.Run(name, func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read %s: %v", file, err)
			}

			stmts := ExtractStatements(base, content)
			kf := knownFailures[base]

			var passed, known, newFail, fixed int
			var failIndices []int

			for i, stmt := range stmts {
				_, parseErr := parser.Parse(stmt.SQL)
				isKnown := intSliceContains(kf, i)

				if parseErr != nil {
					failIndices = append(failIndices, i)
					if isKnown {
						known++
					} else {
						newFail++
						if !*update {
							t.Errorf("NEW FAILURE stmt[%d] line %d: %v\n  SQL: %.200s",
								i, stmt.StartLine, parseErr, stmt.SQL)
						}
					}
				} else {
					passed++
					if isKnown {
						fixed++
						if !*update {
							t.Errorf("FIXED stmt[%d] line %d now passes, remove from known_failures\n  SQL: %.200s",
								i, stmt.StartLine, stmt.SQL)
						}
					}
				}
			}

			if *update && len(failIndices) > 0 {
				newFailures[base] = failIndices
			}

			totalStmts += len(stmts)
			totalPassed += passed
			totalKnown += known
			totalNew += newFail
			totalFixed += fixed

			t.Logf("%s: %d stmts, %d passed, %d failed (%d known, %d new, %d fixed)",
				base, len(stmts), passed, len(failIndices), known, newFail, fixed)
		})
	}

	// Overall summary
	totalFailed := totalKnown + totalNew
	passRate := float64(0)
	if totalStmts > 0 {
		passRate = 100 * float64(totalPassed) / float64(totalStmts)
	}
	t.Logf("\n=== PG REGRESS SUMMARY ===")
	t.Logf("Files:  %d", len(files))
	t.Logf("Stmts:  %d", totalStmts)
	t.Logf("Passed: %d (%.1f%%)", totalPassed, passRate)
	t.Logf("Failed: %d (known: %d, new: %d)", totalFailed, totalKnown, totalNew)
	if totalFixed > 0 {
		t.Logf("Fixed:  %d (remove from known_failures)", totalFixed)
	}

	// Write updated known_failures.json
	if *update {
		data, err := json.MarshalIndent(newFailures, "", "  ")
		if err != nil {
			t.Fatalf("marshal known_failures: %v", err)
		}
		if err := os.WriteFile("known_failures.json", data, 0644); err != nil {
			t.Fatalf("write known_failures.json: %v", err)
		}
		t.Logf("Updated known_failures.json (%d files with failures)", len(newFailures))
	}
}

// TestPGRegressStats is a read-only variant that just prints statistics
// without failing on any parse errors. Useful for quick compatibility checks.
func TestPGRegressStats(t *testing.T) {
	files, err := filepath.Glob("testdata/sql/*.sql")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(files) == 0 {
		t.Skip("no regression test files found in testdata/sql/")
	}

	sort.Strings(files)

	type fileStat struct {
		name   string
		total  int
		passed int
	}
	var stats []fileStat
	var totalStmts, totalPassed int

	for _, file := range files {
		base := filepath.Base(file)
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}

		stmts := ExtractStatements(base, content)
		passed := 0
		for _, stmt := range stmts {
			if _, err := parser.Parse(stmt.SQL); err == nil {
				passed++
			}
		}

		stats = append(stats, fileStat{name: base, total: len(stmts), passed: passed})
		totalStmts += len(stmts)
		totalPassed += passed
	}

	// Print per-file results
	for _, s := range stats {
		rate := float64(0)
		if s.total > 0 {
			rate = 100 * float64(s.passed) / float64(s.total)
		}
		status := "OK"
		if s.passed < s.total {
			status = fmt.Sprintf("PARTIAL (%d failures)", s.total-s.passed)
		}
		t.Logf("%-40s %4d/%4d (%5.1f%%) %s", s.name, s.passed, s.total, rate, status)
	}

	passRate := float64(0)
	if totalStmts > 0 {
		passRate = 100 * float64(totalPassed) / float64(totalStmts)
	}
	t.Logf("\n=== TOTAL: %d/%d statements passed (%.1f%%) across %d files ===",
		totalPassed, totalStmts, passRate, len(files))
}
