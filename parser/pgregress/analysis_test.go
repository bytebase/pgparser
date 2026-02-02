package pgregress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"

	"github.com/pgparser/pgparser/parser"
)

// TestAnalyzeFailures reads known_failures.json, parses every failing statement,
// categorizes failures by error pattern and extractor issues, then prints a
// comprehensive summary report.
func TestAnalyzeFailures(t *testing.T) {
	// --- Load known_failures.json ---
	data, err := os.ReadFile("known_failures.json")
	if err != nil {
		t.Fatalf("read known_failures.json: %v", err)
	}
	var kf KnownFailures
	if err := json.Unmarshal(data, &kf); err != nil {
		t.Fatalf("parse known_failures.json: %v", err)
	}

	// Count total expected failures
	totalExpected := 0
	for _, indices := range kf {
		totalExpected += len(indices)
	}
	t.Logf("Known failures JSON: %d files, %d total failing statement indices", len(kf), totalExpected)

	// --- Collect all failing statements with parse errors ---
	type failEntry struct {
		sql       string
		file      string
		stmtIdx   int
		startLine int
		errMsg    string
	}

	var failures []failEntry
	var extractionMissing int // stmt index out of range

	sortedFiles := make([]string, 0, len(kf))
	for f := range kf {
		sortedFiles = append(sortedFiles, f)
	}
	sort.Strings(sortedFiles)

	for _, base := range sortedFiles {
		indices := kf[base]
		filePath := filepath.Join("testdata", "sql", base)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Logf("WARNING: cannot read %s: %v", filePath, err)
			continue
		}
		stmts := ExtractStatements(base, content)

		for _, idx := range indices {
			if idx >= len(stmts) {
				extractionMissing++
				continue
			}
			stmt := stmts[idx]
			_, parseErr := parser.Parse(stmt.SQL)
			errMsg := ""
			if parseErr != nil {
				errMsg = parseErr.Error()
			}
			failures = append(failures, failEntry{
				sql:       stmt.SQL,
				file:      base,
				stmtIdx:   idx,
				startLine: stmt.StartLine,
				errMsg:    errMsg,
			})
		}
	}

	t.Logf("Collected %d failing statements (%d indices out of range)", len(failures), extractionMissing)

	// Count how many actually still fail vs now pass
	stillFailing := 0
	nowPassing := 0
	for _, f := range failures {
		if f.errMsg != "" {
			stillFailing++
		} else {
			nowPassing++
		}
	}
	t.Logf("Still failing: %d, Now passing: %d", stillFailing, nowPassing)

	// =====================================================================
	// SECTION 1: EXTRACTOR / ARTIFACT ANALYSIS
	// =====================================================================
	t.Logf("\n\n")
	t.Logf("================================================================")
	t.Logf("=== EXTRACTOR / ARTIFACT ANALYSIS ===")
	t.Logf("================================================================")

	// 1a. psql variable interpolation: :varname, :'varname', :"varname"
	psqlVarRE := regexp.MustCompile(`(?::('[^']*'|"[^"]*"|[A-Za-z_][A-Za-z0-9_]*))`)
	// Be more precise: avoid matching things like ::type (cast), or :=
	psqlVarStrictRE := regexp.MustCompile(`(?:^|[^:A-Za-z0-9_])(:(?:'[^']*'|"[^"]*"|[A-Za-z_][A-Za-z0-9_]*))(?:[^:A-Za-z0-9_]|$)`)
	_ = psqlVarRE

	var psqlVarStmts []failEntry
	for _, f := range failures {
		if f.errMsg == "" {
			continue
		}
		if psqlVarStrictRE.MatchString(f.sql) {
			// Double-check it's not a cast (::) or := operator
			// Remove all :: and := occurrences, then check again
			cleaned := strings.ReplaceAll(f.sql, "::", "XX")
			cleaned = strings.ReplaceAll(cleaned, ":=", "XX")
			if psqlVarStrictRE.MatchString(cleaned) {
				psqlVarStmts = append(psqlVarStmts, f)
			}
		}
	}

	t.Logf("\npsql variable interpolation: %d statements", len(psqlVarStmts))
	// Show up to 10 examples
	for i, f := range psqlVarStmts {
		if i >= 10 {
			t.Logf("  ... and %d more", len(psqlVarStmts)-10)
			break
		}
		snippet := truncateSQL(f.sql, 150)
		// Extract the variable references
		matches := psqlVarStrictRE.FindAllStringSubmatch(
			strings.ReplaceAll(strings.ReplaceAll(f.sql, "::", "XX"), ":=", "XX"), -1)
		vars := make([]string, 0, len(matches))
		for _, m := range matches {
			if len(m) > 1 {
				vars = append(vars, m[1])
			}
		}
		t.Logf("  [%s:%d idx=%d] vars=%v", f.file, f.startLine, f.stmtIdx, vars)
		t.Logf("    SQL: %s", snippet)
	}

	// 1b. Statements that look like data lines (not SQL)
	var dataLineStmts []failEntry
	for _, f := range failures {
		if f.errMsg == "" {
			continue
		}
		trimmed := strings.TrimSpace(f.sql)
		if len(trimmed) == 0 {
			continue
		}
		// Check if it starts with a digit, or looks like tab-separated data
		firstChar := rune(trimmed[0])
		isData := false
		if unicode.IsDigit(firstChar) && !startsWithSQLKeyword(trimmed) {
			isData = true
		}
		// Tab-separated data
		if strings.Contains(trimmed, "\t") && !strings.ContainsAny(strings.ToUpper(trimmed[:min(20, len(trimmed))]), "SELECTINSERTCREATEALTERUPDATEDELETEDROPTRUNCAT") {
			isData = true
		}
		// Backslash-dot (copy terminator leaked)
		if trimmed == "\\." {
			isData = true
		}
		if isData {
			dataLineStmts = append(dataLineStmts, f)
		}
	}

	t.Logf("\nData lines (non-SQL): %d statements", len(dataLineStmts))
	for i, f := range dataLineStmts {
		if i >= 10 {
			t.Logf("  ... and %d more", len(dataLineStmts)-10)
			break
		}
		t.Logf("  [%s:%d idx=%d] %.100s", f.file, f.startLine, f.stmtIdx, f.sql)
	}

	// 1c. Unbalanced quotes or parentheses (possible extractor bugs)
	var unbalancedStmts []failEntry
	for _, f := range failures {
		if f.errMsg == "" {
			continue
		}
		reason := checkUnbalanced(f.sql)
		if reason != "" {
			unbalancedStmts = append(unbalancedStmts, f)
		}
	}

	t.Logf("\nUnbalanced quotes/parens (possible extractor bugs): %d statements", len(unbalancedStmts))
	for i, f := range unbalancedStmts {
		if i >= 10 {
			t.Logf("  ... and %d more", len(unbalancedStmts)-10)
			break
		}
		reason := checkUnbalanced(f.sql)
		snippet := truncateSQL(f.sql, 150)
		t.Logf("  [%s:%d idx=%d] reason=%s", f.file, f.startLine, f.stmtIdx, reason)
		t.Logf("    SQL: %s", snippet)
	}

	// 1d. Psql metacommand artifacts that leaked through
	psqlMetaRE := regexp.MustCompile(`(?m)^\\[a-zA-Z]`)
	var psqlMetaStmts []failEntry
	for _, f := range failures {
		if f.errMsg == "" {
			continue
		}
		if psqlMetaRE.MatchString(f.sql) {
			psqlMetaStmts = append(psqlMetaStmts, f)
		}
	}

	t.Logf("\npsql metacommand artifacts: %d statements", len(psqlMetaStmts))
	for i, f := range psqlMetaStmts {
		if i >= 5 {
			t.Logf("  ... and %d more", len(psqlMetaStmts)-5)
			break
		}
		t.Logf("  [%s:%d idx=%d] %.150s", f.file, f.startLine, f.stmtIdx, f.sql)
	}

	// 1e. Empty or whitespace-only statements
	var emptyStmts []failEntry
	for _, f := range failures {
		if strings.TrimSpace(f.sql) == "" {
			emptyStmts = append(emptyStmts, f)
		}
	}
	t.Logf("\nEmpty/whitespace-only: %d statements", len(emptyStmts))

	// 1f. Comment-only statements (the extractor emitted a "statement" that
	// is entirely SQL comments with no executable SQL)
	var commentOnlyStmts []failEntry
	for _, f := range failures {
		if f.errMsg == "" {
			continue
		}
		if isCommentOnly(f.sql) {
			commentOnlyStmts = append(commentOnlyStmts, f)
		}
	}
	t.Logf("\nComment-only (no executable SQL): %d statements", len(commentOnlyStmts))
	for i, f := range commentOnlyStmts {
		if i >= 10 {
			t.Logf("  ... and %d more", len(commentOnlyStmts)-10)
			break
		}
		snippet := truncateSQL(f.sql, 150)
		t.Logf("  [%s:%d idx=%d] %s", f.file, f.startLine, f.stmtIdx, snippet)
	}

	// =====================================================================
	// SECTION 2: ERROR PATTERN ANALYSIS
	// =====================================================================
	t.Logf("\n\n")
	t.Logf("================================================================")
	t.Logf("=== ERROR PATTERN ANALYSIS ===")
	t.Logf("================================================================")

	// Normalize error messages to group them
	type errorGroup struct {
		pattern  string
		count    int
		examples []failEntry
	}

	patternMap := make(map[string]*errorGroup)

	for _, f := range failures {
		if f.errMsg == "" {
			continue // now passes
		}
		pat := normalizeError(f.errMsg)
		eg, ok := patternMap[pat]
		if !ok {
			eg = &errorGroup{pattern: pat}
			patternMap[pat] = eg
		}
		eg.count++
		if len(eg.examples) < 3 {
			eg.examples = append(eg.examples, f)
		}
	}

	// Sort by count descending
	groups := make([]*errorGroup, 0, len(patternMap))
	for _, eg := range patternMap {
		groups = append(groups, eg)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].count > groups[j].count
	})

	t.Logf("\n%d distinct error patterns (from %d still-failing statements):\n", len(groups), stillFailing)
	for rank, eg := range groups {
		pctVal := 100.0 * float64(eg.count) / float64(stillFailing)
		t.Logf("--- Pattern #%d: %d occurrences (%.1f%%)", rank+1, eg.count, pctVal)
		t.Logf("    Error: %s", eg.pattern)
		for j, ex := range eg.examples {
			snippet := truncateSQL(ex.sql, 150)
			t.Logf("    Example %d [%s:%d idx=%d]: %s", j+1, ex.file, ex.startLine, ex.stmtIdx, snippet)
		}
		t.Logf("")
	}

	// =====================================================================
	// SECTION 2b: SQL CONSTRUCT / KEYWORD HEURISTIC ANALYSIS
	// =====================================================================
	// Since the parser only emits generic "syntax error", identify what SQL
	// constructs/keywords appear in failing statements to guide grammar work.
	t.Logf("\n")
	t.Logf("================================================================")
	t.Logf("=== SQL CONSTRUCT / KEYWORD PRESENCE IN FAILING STATEMENTS ===")
	t.Logf("================================================================")
	t.Logf("(How many failing SQL statements contain these keywords/constructs?)")

	constructKeywords := []string{
		// DDL constructs
		"CREATE TABLE", "ALTER TABLE", "CREATE INDEX", "CREATE VIEW",
		"CREATE FUNCTION", "CREATE PROCEDURE", "CREATE TRIGGER",
		"CREATE TYPE", "CREATE DOMAIN", "CREATE SCHEMA", "CREATE RULE",
		"CREATE POLICY", "CREATE PUBLICATION", "CREATE SUBSCRIPTION",
		"CREATE STATISTICS", "CREATE FOREIGN", "CREATE MATERIALIZED",
		"CREATE AGGREGATE", "CREATE OPERATOR", "CREATE CAST",
		"CREATE EXTENSION", "CREATE CONVERSION", "CREATE SEQUENCE",
		"CREATE TABLESPACE", "CREATE ROLE", "CREATE USER",
		// DML constructs
		"SELECT", "INSERT", "UPDATE", "DELETE", "MERGE",
		"WITH RECURSIVE", "WITH",
		"COPY FROM", "COPY TO", "COPY",
		// Clauses/features
		"RETURNING", "ON CONFLICT", "LATERAL", "TABLESAMPLE",
		"WINDOW", "GROUPING SETS", "CUBE", "ROLLUP",
		"FETCH FIRST", "OFFSET", "LIMIT",
		"PARTITION BY", "GENERATED", "IDENTITY",
		// PG-specific SQL features
		"XMLTABLE", "JSON_TABLE", "JSON_VALUE", "JSON_QUERY",
		"JSON_EXISTS", "JSON_OBJECT", "JSON_ARRAY", "JSON_SCALAR",
		"JSON_SERIALIZE", "XMLNAMESPACES",
		"FOR UPDATE", "FOR SHARE", "FOR NO KEY UPDATE",
		"LISTEN", "NOTIFY", "UNLISTEN",
		"SECURITY LABEL", "COMMENT ON",
		"DECLARE CURSOR", "DECLARE",
		"DO $$", "DO",
		"EXPLAIN", "ANALYZE",
		"VACUUM", "CLUSTER", "REINDEX",
		"GRANT", "REVOKE",
		"SET", "RESET", "SHOW",
		"BEGIN", "COMMIT", "ROLLBACK", "SAVEPOINT",
		"PREPARE", "EXECUTE", "DEALLOCATE",
		// Expressions/operators
		"ARRAY[", "ROW(", "CASE WHEN", "COALESCE", "NULLIF",
		"GREATEST", "LEAST", "EXISTS",
		"IS DISTINCT FROM", "IS NOT DISTINCT FROM",
		"BETWEEN", "SIMILAR TO", "LIKE",
		"OVERLAPS",
		// Type-related
		"::text", "::int", "::json", "::jsonb",
		"CAST(",
	}

	// Only count among genuine SQL failures (exclude comment-only)
	commentOnlySetForConstructs := make(map[string]bool)
	for _, f := range commentOnlyStmts {
		commentOnlySetForConstructs[fmt.Sprintf("%s:%d", f.file, f.stmtIdx)] = true
	}

	type constructCount struct {
		keyword string
		count   int
	}
	var constructCounts []constructCount

	for _, kw := range constructKeywords {
		count := 0
		kwUpper := strings.ToUpper(kw)
		for _, f := range failures {
			if f.errMsg == "" {
				continue
			}
			key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
			if commentOnlySetForConstructs[key] {
				continue
			}
			sqlUpper := strings.ToUpper(f.sql)
			if strings.Contains(sqlUpper, kwUpper) {
				count++
			}
		}
		if count > 0 {
			constructCounts = append(constructCounts, constructCount{kw, count})
		}
	}

	sort.Slice(constructCounts, func(i, j int) bool {
		return constructCounts[i].count > constructCounts[j].count
	})

	t.Logf("")
	for _, cc := range constructCounts {
		t.Logf("  %-35s %5d", cc.keyword, cc.count)
	}

	// =====================================================================
	// SECTION 3: SQL KEYWORD / STATEMENT TYPE BREAKDOWN
	// =====================================================================
	t.Logf("\n")
	t.Logf("================================================================")
	t.Logf("=== FAILING STATEMENT TYPE BREAKDOWN ===")
	t.Logf("================================================================")

	stmtTypeCounts := make(map[string]int)
	for _, f := range failures {
		if f.errMsg == "" {
			continue
		}
		stmtType := classifySQL(f.sql)
		stmtTypeCounts[stmtType]++
	}

	type stmtTypeEntry struct {
		name  string
		count int
	}
	stmtTypes := make([]stmtTypeEntry, 0, len(stmtTypeCounts))
	for name, count := range stmtTypeCounts {
		stmtTypes = append(stmtTypes, stmtTypeEntry{name, count})
	}
	sort.Slice(stmtTypes, func(i, j int) bool {
		return stmtTypes[i].count > stmtTypes[j].count
	})

	t.Logf("\nStatement types among %d still-failing statements:\n", stillFailing)
	for _, st := range stmtTypes {
		pct := 100.0 * float64(st.count) / float64(stillFailing)
		t.Logf("  %-30s %5d (%5.1f%%)", st.name, st.count, pct)
	}

	// =====================================================================
	// SECTION 4: TOP-LEVEL CATEGORY SUMMARY
	// =====================================================================
	t.Logf("\n\n")
	t.Logf("================================================================")
	t.Logf("=== TOP-LEVEL CATEGORY SUMMARY ===")
	t.Logf("================================================================")

	// Categorize each failing statement into one primary bucket
	catPsqlVar := 0
	catExtractor := 0
	catCommentOnly := 0
	catParserGap := 0
	catNowPassing := nowPassing

	// Build sets for dedup
	psqlVarSet := make(map[string]bool)
	for _, f := range psqlVarStmts {
		key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
		psqlVarSet[key] = true
	}
	extractorSet := make(map[string]bool)
	for _, f := range dataLineStmts {
		key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
		extractorSet[key] = true
	}
	for _, f := range unbalancedStmts {
		key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
		extractorSet[key] = true
	}
	for _, f := range psqlMetaStmts {
		key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
		extractorSet[key] = true
	}
	for _, f := range emptyStmts {
		key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
		extractorSet[key] = true
	}
	commentOnlySet2 := make(map[string]bool)
	for _, f := range commentOnlyStmts {
		key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
		commentOnlySet2[key] = true
	}

	for _, f := range failures {
		if f.errMsg == "" {
			continue // counted as nowPassing
		}
		key := fmt.Sprintf("%s:%d", f.file, f.stmtIdx)
		if psqlVarSet[key] {
			catPsqlVar++
		} else if extractorSet[key] {
			catExtractor++
		} else if commentOnlySet2[key] {
			catCommentOnly++
		} else {
			catParserGap++
		}
	}

	t.Logf("\nTotal entries in known_failures.json: %d", totalExpected)
	t.Logf("Statements analyzed:                  %d", len(failures))
	t.Logf("Index out of range (extraction gap):  %d", extractionMissing)
	t.Logf("")
	t.Logf("  Now passing (parser improved):      %5d (%5.1f%%)", catNowPassing, pct(catNowPassing, len(failures)))
	t.Logf("  psql variable interpolation:        %5d (%5.1f%%)", catPsqlVar, pct(catPsqlVar, len(failures)))
	t.Logf("  Extractor bugs / artifacts:         %5d (%5.1f%%)", catExtractor, pct(catExtractor, len(failures)))
	t.Logf("  Comment-only (no SQL content):      %5d (%5.1f%%)", catCommentOnly, pct(catCommentOnly, len(failures)))
	t.Logf("  Parser gaps (grammar not impl):     %5d (%5.1f%%)", catParserGap, pct(catParserGap, len(failures)))

	// =====================================================================
	// SECTION 5: ERROR TOKEN ANALYSIS (what tokens cause most failures)
	// =====================================================================
	t.Logf("\n\n")
	t.Logf("================================================================")
	t.Logf("=== ERROR TOKEN ANALYSIS ===")
	t.Logf("================================================================")

	// Extract the token from "syntax error at or near X" messages
	tokenRE := regexp.MustCompile(`syntax error at or near "([^"]*)"`)
	tokenCounts := make(map[string]int)
	for _, f := range failures {
		if f.errMsg == "" {
			continue
		}
		m := tokenRE.FindStringSubmatch(f.errMsg)
		if len(m) > 1 {
			tokenCounts[m[1]]++
		}
	}

	type tokenEntry struct {
		token string
		count int
	}
	tokens := make([]tokenEntry, 0, len(tokenCounts))
	for tok, cnt := range tokenCounts {
		tokens = append(tokens, tokenEntry{tok, cnt})
	}
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].count > tokens[j].count
	})

	t.Logf("\nTop error tokens (from 'syntax error at or near' messages):\n")
	for i, te := range tokens {
		if i >= 40 {
			t.Logf("  ... and %d more distinct tokens", len(tokens)-40)
			break
		}
		t.Logf("  %-30s %5d", fmt.Sprintf("%q", te.token), te.count)
	}

	// =====================================================================
	// SECTION 6: FILE-LEVEL BREAKDOWN
	// =====================================================================
	t.Logf("\n\n")
	t.Logf("================================================================")
	t.Logf("=== TOP 30 FILES BY FAILURE COUNT ===")
	t.Logf("================================================================")

	fileCounts := make(map[string]int)
	for _, f := range failures {
		if f.errMsg != "" {
			fileCounts[f.file]++
		}
	}

	type fileEntry struct {
		file  string
		count int
	}
	fileEntries := make([]fileEntry, 0, len(fileCounts))
	for f, c := range fileCounts {
		fileEntries = append(fileEntries, fileEntry{f, c})
	}
	sort.Slice(fileEntries, func(i, j int) bool {
		return fileEntries[i].count > fileEntries[j].count
	})

	t.Logf("")
	for i, fe := range fileEntries {
		if i >= 30 {
			break
		}
		t.Logf("  %-40s %5d failures", fe.file, fe.count)
	}
}

// --- Helper functions ---

func pct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return 100.0 * float64(n) / float64(total)
}

func truncateSQL(sql string, maxLen int) string {
	// Collapse multiple whitespace/newlines
	s := strings.Join(strings.Fields(sql), " ")
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// normalizeError groups error messages by replacing specific tokens/positions
// with placeholders to find patterns.
func normalizeError(msg string) string {
	// "syntax error at or near "SOMETHING"" -> "syntax error at or near <TOKEN>"
	re1 := regexp.MustCompile(`syntax error at or near "[^"]*"`)
	result := re1.ReplaceAllString(msg, `syntax error at or near <TOKEN>`)

	// "parse error (ret=N)" -> "parse error (ret=N)"
	re2 := regexp.MustCompile(`\(ret=\d+\)`)
	result = re2.ReplaceAllString(result, "(ret=N)")

	return result
}

// classifySQL returns a broad classification of the SQL statement type.
func classifySQL(sql string) string {
	// Strip all leading comments (line and block) to find the first real SQL token
	stripped := stripLeadingComments(sql)
	upper := strings.ToUpper(strings.TrimSpace(stripped))
	fields := strings.Fields(upper)
	if len(fields) == 0 {
		return "(comment-only)"
	}

	first := fields[0]

	switch first {
	case "SELECT", "WITH":
		return "SELECT/WITH"
	case "INSERT":
		return "INSERT"
	case "UPDATE":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	case "CREATE":
		if len(fields) > 1 {
			second := fields[1]
			if second == "OR" && len(fields) > 3 {
				return "CREATE " + fields[3]
			}
			if second == "UNIQUE" || second == "INDEX" {
				return "CREATE INDEX"
			}
			if second == "TEMP" || second == "TEMPORARY" || second == "UNLOGGED" || second == "LOCAL" || second == "GLOBAL" {
				if len(fields) > 2 {
					return "CREATE " + second + " " + fields[2]
				}
			}
			return "CREATE " + second
		}
		return "CREATE"
	case "ALTER":
		if len(fields) > 1 {
			return "ALTER " + fields[1]
		}
		return "ALTER"
	case "DROP":
		if len(fields) > 1 {
			return "DROP " + fields[1]
		}
		return "DROP"
	case "SET":
		return "SET"
	case "RESET":
		return "RESET"
	case "SHOW":
		return "SHOW"
	case "GRANT":
		return "GRANT"
	case "REVOKE":
		return "REVOKE"
	case "BEGIN", "START":
		return "BEGIN/START"
	case "COMMIT", "END":
		return "COMMIT/END"
	case "ROLLBACK":
		return "ROLLBACK"
	case "SAVEPOINT":
		return "SAVEPOINT"
	case "PREPARE":
		return "PREPARE"
	case "EXECUTE":
		return "EXECUTE"
	case "DEALLOCATE":
		return "DEALLOCATE"
	case "EXPLAIN":
		return "EXPLAIN"
	case "ANALYZE", "ANALYSE":
		return "ANALYZE"
	case "VACUUM":
		return "VACUUM"
	case "COPY":
		return "COPY"
	case "DO":
		return "DO"
	case "CALL":
		return "CALL"
	case "LOCK":
		return "LOCK"
	case "COMMENT":
		return "COMMENT"
	case "NOTIFY":
		return "NOTIFY"
	case "LISTEN":
		return "LISTEN"
	case "UNLISTEN":
		return "UNLISTEN"
	case "FETCH":
		return "FETCH"
	case "MOVE":
		return "MOVE"
	case "CLOSE":
		return "CLOSE"
	case "DECLARE":
		return "DECLARE"
	case "CLUSTER":
		return "CLUSTER"
	case "REINDEX":
		return "REINDEX"
	case "LOAD":
		return "LOAD"
	case "REASSIGN":
		return "REASSIGN"
	case "SECURITY":
		return "SECURITY LABEL"
	case "MERGE":
		return "MERGE"
	case "TABLE":
		return "TABLE"
	case "VALUES":
		return "VALUES"
	case "TRUNCATE":
		return "TRUNCATE"
	case "DISCARD":
		return "DISCARD"
	case "IMPORT":
		return "IMPORT"
	case "REFRESH":
		return "REFRESH"
	case "CHECKPOINT":
		return "CHECKPOINT"
	case "RELEASE":
		return "RELEASE"
	default:
		// Check for some patterns
		if strings.HasPrefix(first, "(") {
			return "(subquery/expr)"
		}
		return "OTHER: " + first
	}
}

// stripLeadingComments removes all leading SQL comments (line -- and block /* */)
// from the beginning of a SQL string, returning the remainder.
func stripLeadingComments(sql string) string {
	s := sql
	for {
		s = strings.TrimSpace(s)
		if s == "" {
			return ""
		}
		// Line comment
		if strings.HasPrefix(s, "--") {
			idx := strings.IndexByte(s, '\n')
			if idx < 0 {
				return "" // entire string is a comment
			}
			s = s[idx+1:]
			continue
		}
		// Block comment
		if strings.HasPrefix(s, "/*") {
			depth := 1
			i := 2
			for i < len(s) && depth > 0 {
				if i+1 < len(s) && s[i] == '/' && s[i+1] == '*' {
					depth++
					i += 2
				} else if i+1 < len(s) && s[i] == '*' && s[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			s = s[i:]
			continue
		}
		// Not a comment
		return s
	}
}

// startsWithSQLKeyword checks if the trimmed string starts with a common SQL keyword.
func startsWithSQLKeyword(s string) bool {
	upper := strings.ToUpper(s)
	keywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP",
		"SET", "RESET", "SHOW", "GRANT", "REVOKE", "BEGIN", "COMMIT",
		"ROLLBACK", "SAVEPOINT", "PREPARE", "EXECUTE", "DEALLOCATE",
		"EXPLAIN", "ANALYZE", "VACUUM", "COPY", "DO", "CALL", "LOCK",
		"COMMENT", "NOTIFY", "LISTEN", "UNLISTEN", "FETCH", "MOVE",
		"CLOSE", "DECLARE", "WITH", "TABLE", "VALUES", "TRUNCATE",
		"CLUSTER", "REINDEX", "LOAD", "MERGE", "DISCARD", "IMPORT",
		"REFRESH", "CHECKPOINT", "RELEASE", "REASSIGN", "SECURITY",
		"START", "END", "ABORT",
	}
	for _, kw := range keywords {
		if strings.HasPrefix(upper, kw) && (len(upper) == len(kw) || !isIdentByte(upper[len(kw)])) {
			return true
		}
	}
	return false
}

func isIdentByte(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '_'
}

// checkUnbalanced does a rough check for unbalanced quotes or parentheses.
// Returns a reason string if unbalanced, empty string if OK.
func checkUnbalanced(sql string) string {
	// Track state
	parenDepth := 0
	bracketDepth := 0
	inSingle := false
	inDouble := false
	inDollar := false
	dollarTag := ""
	inEscape := false
	inLineComment := false
	blockDepth := 0

	runes := []byte(sql)
	n := len(runes)

	for i := 0; i < n; i++ {
		ch := runes[i]

		// Line comment
		if inLineComment {
			if ch == '\n' {
				inLineComment = false
			}
			continue
		}

		// Block comment
		if blockDepth > 0 {
			if ch == '/' && i+1 < n && runes[i+1] == '*' {
				blockDepth++
				i++
			} else if ch == '*' && i+1 < n && runes[i+1] == '/' {
				blockDepth--
				i++
			}
			continue
		}

		// E-string escape mode
		if inEscape {
			if ch == '\\' {
				i++ // skip next
				continue
			}
			if ch == '\'' {
				if i+1 < n && runes[i+1] == '\'' {
					i++
					continue
				}
				inEscape = false
			}
			continue
		}

		// Single-quoted string
		if inSingle {
			if ch == '\'' {
				if i+1 < n && runes[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
			}
			continue
		}

		// Double-quoted identifier
		if inDouble {
			if ch == '"' {
				if i+1 < n && runes[i+1] == '"' {
					i++
					continue
				}
				inDouble = false
			}
			continue
		}

		// Dollar-quoted string
		if inDollar {
			if ch == '$' && i+len(dollarTag) <= n && string(runes[i:i+len(dollarTag)]) == dollarTag {
				i += len(dollarTag) - 1
				inDollar = false
			}
			continue
		}

		// Normal mode
		switch {
		case ch == '-' && i+1 < n && runes[i+1] == '-':
			inLineComment = true
			i++
		case ch == '/' && i+1 < n && runes[i+1] == '*':
			blockDepth++
			i++
		case ch == '\'' && i > 0 && (runes[i-1] == 'E' || runes[i-1] == 'e'):
			inEscape = true
		case ch == '\'':
			inSingle = true
		case ch == '"':
			inDouble = true
		case ch == '$':
			tag := tryParseDollarTagStr(sql, i)
			if tag != "" {
				dollarTag = tag
				inDollar = true
				i += len(tag) - 1
			}
		case ch == '(':
			parenDepth++
		case ch == ')':
			parenDepth--
		case ch == '[':
			bracketDepth++
		case ch == ']':
			bracketDepth--
		}
	}

	var reasons []string
	if inSingle || inEscape {
		reasons = append(reasons, "unclosed single quote")
	}
	if inDouble {
		reasons = append(reasons, "unclosed double quote")
	}
	if inDollar {
		reasons = append(reasons, "unclosed dollar quote ("+dollarTag+")")
	}
	if blockDepth > 0 {
		reasons = append(reasons, fmt.Sprintf("unclosed block comment (depth %d)", blockDepth))
	}
	if parenDepth > 0 {
		reasons = append(reasons, fmt.Sprintf("unclosed parens (depth %d)", parenDepth))
	}
	if parenDepth < 0 {
		reasons = append(reasons, fmt.Sprintf("extra close parens (%d)", -parenDepth))
	}
	if bracketDepth > 0 {
		reasons = append(reasons, fmt.Sprintf("unclosed brackets (depth %d)", bracketDepth))
	}
	if bracketDepth < 0 {
		reasons = append(reasons, fmt.Sprintf("extra close brackets (%d)", -bracketDepth))
	}

	return strings.Join(reasons, "; ")
}

// isCommentOnly returns true if the SQL string contains only comments
// (line comments and/or block comments) and no executable SQL.
func isCommentOnly(sql string) bool {
	s := strings.TrimSpace(sql)
	if s == "" {
		return true
	}

	i := 0
	n := len(s)
	for i < n {
		ch := s[i]

		// Skip whitespace
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			i++
			continue
		}

		// Line comment
		if ch == '-' && i+1 < n && s[i+1] == '-' {
			// Skip to end of line
			for i < n && s[i] != '\n' {
				i++
			}
			continue
		}

		// Block comment
		if ch == '/' && i+1 < n && s[i+1] == '*' {
			depth := 1
			i += 2
			for i < n && depth > 0 {
				if s[i] == '/' && i+1 < n && s[i+1] == '*' {
					depth++
					i += 2
				} else if s[i] == '*' && i+1 < n && s[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			continue
		}

		// Found a non-comment, non-whitespace character
		return false
	}

	return true
}
