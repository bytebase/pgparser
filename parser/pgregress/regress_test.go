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
	var totalStmts, totalPassed, totalKnown, totalNew, totalFixed, totalPsqlVar int

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

			var passed, known, newFail, fixed, psqlVar int
			var failIndices []int

			for i, stmt := range stmts {

				sqlToParse := stmt.SQL
				if stmt.HasPsqlVar {
					sqlToParse, _ = ReplacePsqlVariables(stmt.SQL)
					// Special case for full psql variable statements like ":variable;"
					// which become "psql_var;" after replacement. This is invalid syntax
					// (identifier as statement). Treat it as a passing SELECT 1.
					trimmed := strings.TrimSpace(sqlToParse)
					if trimmed == "psql_var" || trimmed == "psql_var;" {
						sqlToParse = "SELECT 1"
					} else {
						// Handle "EXPLAIN ... :var" -> "EXPLAIN ... psql_var" -> "EXPLAIN ... SELECT 1"
						// We check if EXPLAIN appears before psql_var.
						upper := strings.ToUpper(sqlToParse)
						expIdx := strings.Index(upper, "EXPLAIN")
						varIdx := strings.Index(sqlToParse, "psql_var")
						if expIdx >= 0 && varIdx > expIdx {
							sqlToParse = strings.Replace(sqlToParse, "psql_var", "SELECT 1", 1)
						}
					}
				}
				_, parseErr := parser.Parse(sqlToParse)

				isKnown := intSliceContains(kf, i)

				if parseErr != nil {
					failIndices = append(failIndices, i)
					if stmt.HasPsqlVar {
						psqlVar++
					}
					if isKnown {
						known++
					} else {
						newFail++
						if !*update {
							label := "NEW FAILURE"
							if stmt.HasPsqlVar {
								label = "NEW FAILURE (psql variable)"
							}
							t.Errorf("%s stmt[%d] line %d: %v\n  SQL: %.200s",
								label, i, stmt.StartLine, parseErr, stmt.SQL)
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
			totalPsqlVar += psqlVar

			t.Logf("%s: %d stmts, %d passed, %d failed (%d known, %d new, %d fixed, %d psql-var)",
				base, len(stmts), passed, len(failIndices), known, newFail, fixed, psqlVar)
		})
	}

	// Overall summary
	totalFailed := totalKnown + totalNew
	passRate := float64(0)
	if totalStmts > 0 {
		passRate = 100 * float64(totalPassed) / float64(totalStmts)
	}
	t.Logf("\n=== PG REGRESS SUMMARY ===")
	t.Logf("Files:    %d", len(files))
	t.Logf("Stmts:    %d", totalStmts)
	t.Logf("Passed:   %d (%.1f%%)", totalPassed, passRate)
	t.Logf("Failed:   %d (known: %d, new: %d)", totalFailed, totalKnown, totalNew)
	if totalPsqlVar > 0 {
		t.Logf("Psql-var: %d (psql variable interpolation, not real grammar failures)", totalPsqlVar)
	}
	if totalFixed > 0 {
		t.Logf("Fixed:    %d (remove from known_failures)", totalFixed)
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
		name    string
		total   int
		passed  int
		psqlVar int // failures due to psql variable interpolation
	}
	var stats []fileStat
	var totalStmts, totalPassed, totalPsqlVar int

	for _, file := range files {
		base := filepath.Base(file)
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}

		stmts := ExtractStatements(base, content)
		passed := 0
		psqlVar := 0
		for _, stmt := range stmts {
			sql := stmt.SQL
			if stmt.HasPsqlVar {
				sql, _ = ReplacePsqlVariables(stmt.SQL)
				// Special case for full psql variable statements like ":variable;"
				trimmed := strings.TrimSpace(sql)
				if trimmed == "psql_var" || trimmed == "psql_var;" {
					sql = "SELECT 1"
				} else {
					upper := strings.ToUpper(sql)
					expIdx := strings.Index(upper, "EXPLAIN")
					varIdx := strings.Index(sql, "psql_var")
					if expIdx >= 0 && varIdx > expIdx {
						sql = strings.Replace(sql, "psql_var", "SELECT 1", 1)
					}
				}
			}
			if _, err := parser.Parse(sql); err == nil {
				passed++
			} else if stmt.HasPsqlVar {
				psqlVar++
			} else {
				// FAILURE
			}

		}

		stats = append(stats, fileStat{name: base, total: len(stmts), passed: passed, psqlVar: psqlVar})
		totalStmts += len(stmts)
		totalPassed += passed
		totalPsqlVar += psqlVar
	}

	knownFailures := loadKnownFailures(t)

	// Print per-file results
	for _, s := range stats {
		rate := float64(0)
		if s.total > 0 {
			rate = 100 * float64(s.passed) / float64(s.total)
		}
		status := "OK"
		if s.passed < s.total {
			failures := s.total - s.passed

			// Determine if failures are expected
			kf := knownFailures[s.name]
			// We can't verify EXACTLY which ones failed here without re-running or storing failure indices in stats
			// But heuristically, if failure count matches known failure count, we can mark it differently

			statusLabel := "PARTIAL"
			// Adjust failure count by subtracting known failures
			// This is a rough heuristic for the stats summary
			// For precise checking, TestPGRegress is used.

			// Actually, let's just leave it as is for now, but maybe add "KNOWN" label if count matches?
			if len(kf) > 0 && failures == len(kf)+s.psqlVar { // Assuming psqlVar are separate from parse failures
				// psqlVar are separate because they failed parsing
				// Wait, in the loop: if err!=nil -> passed not incremented.
				// If psqlVar -> psqlVar incremented.
				// failures = total - passed.
				// parseErr != nil -> failures++.
				// psqlVar is a subset of failures?
				// Loop logic:
				// if err == nil: passed++
				// else:
				//    if hasPsqlVar: psqlVar++
				//    else: // regular failure
				//
				// So failures = psqlVar + regular_failures.
				// Regular failures should match len(kf) ideally.

				if failures == len(kf)+s.psqlVar {
					statusLabel = "KNOWN-FAIL"
				}
			}

			if s.psqlVar > 0 {
				status = fmt.Sprintf("%s (%d failures, %d psql-var)", statusLabel, failures, s.psqlVar)
			} else {
				status = fmt.Sprintf("%s (%d failures)", statusLabel, failures)
			}
		}

		t.Logf("%-40s %4d/%4d (%5.1f%%) %s", s.name, s.passed, s.total, rate, status)
	}

	passRate := float64(0)
	if totalStmts > 0 {
		passRate = 100 * float64(totalPassed) / float64(totalStmts)
	}
	t.Logf("\n=== TOTAL: %d/%d statements passed (%.1f%%) across %d files ===",
		totalPassed, totalStmts, passRate, len(files))
	if totalPsqlVar > 0 {
		t.Logf("  Psql-var failures: %d (not real grammar failures)", totalPsqlVar)
	}
}
