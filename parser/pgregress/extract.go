// Package pgregress provides tools for running PostgreSQL regression test
// SQL files against the pgparser parser to verify parse compatibility.
package pgregress

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
)

// ExtractedStmt represents a single SQL statement extracted from a .sql file.
type ExtractedStmt struct {
	SQL        string // The complete SQL text (without trailing semicolon)
	File       string // Source filename (basename)
	StartLine  int    // 1-based line number where this statement starts
	HasPsqlVar bool   // True if the statement contains psql variable interpolation
}

// ExtractStatements splits a PostgreSQL regression test .sql file into
// individual SQL statements. It handles psql metacommands, dollar-quoted
// strings, COPY FROM stdin data blocks, and all PostgreSQL string literal forms.
func ExtractStatements(filename string, content []byte) []ExtractedStmt {
	lines := preprocessLines(content)
	return splitStatements(filename, lines)
}

// processedLine is a line after psql metacommand preprocessing.
type processedLine struct {
	text       string // SQL text for this line (empty if skipped)
	lineNum    int    // original 1-based line number
	terminates bool   // true if a psql terminator (\g, \gset, \gx) ended this line
}

// psqlTerminatorRE matches psql metacommands that act as statement terminators.
var psqlTerminatorRE = regexp.MustCompile(`(?i)\\(g|gx|gset|gdesc|gexec|crosstabview)\b`)

// copyFromStdinRE detects COPY ... FROM STDIN statements.
var copyFromStdinRE = regexp.MustCompile(`(?i)\bCOPY\b[^;]*\bFROM\b\s+STDIN\b`)

// preprocessLines handles psql metacommands and COPY FROM stdin data blocks.
// It returns processed lines with metacommands removed and COPY data skipped.
func preprocessLines(content []byte) []processedLine {
	rawLines := bytes.Split(content, []byte("\n"))
	result := make([]processedLine, 0, len(rawLines))
	inCopyData := false

	for i, raw := range rawLines {
		lineNum := i + 1
		line := string(raw)

		// In COPY data mode, skip until \. terminator
		if inCopyData {
			trimmed := strings.TrimRight(line, "\r")
			if trimmed == "\\." {
				inCopyData = false
			}
			// Skip this line either way (data or terminator)
			result = append(result, processedLine{lineNum: lineNum})
			continue
		}

		trimmedLeft := strings.TrimLeftFunc(line, unicode.IsSpace)

		// Check for psql metacommand (line starts with \)
		if strings.HasPrefix(trimmedLeft, `\`) {
			pl := processedLine{lineNum: lineNum}
			// If this is a terminator like \g, \gset, \gx, mark it
			if psqlTerminatorRE.MatchString(trimmedLeft) {
				pl.terminates = true
			}
			result = append(result, pl)
			continue
		}

		// Check for mid-line psql metacommand (e.g., "SELECT 1 \gset")
		// We need to find \ that starts a metacommand but is NOT inside a string.
		sqlPart, terminated := splitMidLineMeta(line)
		if terminated {
			result = append(result, processedLine{
				text:       sqlPart,
				lineNum:    lineNum,
				terminates: true,
			})
			continue
		}

		// Normal SQL line
		result = append(result, processedLine{
			text:    line,
			lineNum: lineNum,
		})

		// Check if this line completes a COPY FROM STDIN statement.
		// We check if the line ends with a semicolon (possibly after whitespace)
		// and the accumulated context looks like COPY FROM STDIN.
		// This is a heuristic: we'll do the real check in splitStatements
		// after we see the full statement. But for the line-level preprocessing,
		// we set inCopyData after the statement-ending semicolon on a COPY FROM STDIN.
		// Actually, the real detection must happen at the statement level.
		// We'll handle this by scanning for the pattern across the recent lines.
		// For simplicity, check if this line contains "stdin" + ";" pattern.
		if looksLikeCopyStdinEnd(line) {
			inCopyData = true
		}
	}

	return result
}

// looksLikeCopyStdinEnd is a quick heuristic to detect if a line might end
// a COPY FROM STDIN statement. We look for STDIN followed by optional clauses and ;.
// This is imperfect but sufficient for regression test files.
func looksLikeCopyStdinEnd(line string) bool {
	upper := strings.ToUpper(line)
	// Check this single line for the pattern
	if strings.Contains(upper, "STDIN") && strings.Contains(line, ";") {
		return true
	}
	return false
}

// splitMidLineMeta checks if a line contains a mid-line psql metacommand
// like "SELECT 1 \gset". Returns the SQL portion and whether it terminated.
// It does a simple scan to avoid splitting inside strings.
func splitMidLineMeta(line string) (sqlPart string, terminated bool) {
	n := len(line)
	inSingle := false
	inDouble := false
	inDollar := false
	dollarTag := ""
	i := 0

	for i < n {
		ch := line[i]

		if inSingle {
			if ch == '\'' {
				if i+1 < n && line[i+1] == '\'' {
					i += 2
					continue
				}
				inSingle = false
			}
			i++
			continue
		}

		if inDouble {
			if ch == '"' {
				if i+1 < n && line[i+1] == '"' {
					i += 2
					continue
				}
				inDouble = false
			}
			i++
			continue
		}

		if inDollar {
			// Look for closing dollar tag
			if ch == '$' && i+len(dollarTag) <= n && line[i:i+len(dollarTag)] == dollarTag {
				i += len(dollarTag)
				inDollar = false
				continue
			}
			i++
			continue
		}

		switch ch {
		case '\'':
			// Check for E'...'
			if i > 0 && (line[i-1] == 'E' || line[i-1] == 'e') {
				// E-string, still enters single-quote mode
			}
			inSingle = true
			i++
		case '"':
			inDouble = true
			i++
		case '$':
			tag := tryParseDollarTagStr(line, i)
			if tag != "" {
				dollarTag = tag
				inDollar = true
				i += len(tag)
			} else {
				i++
			}
		case '-':
			if i+1 < n && line[i+1] == '-' {
				// Line comment - rest of line is comment, no metacommand possible
				return "", false
			}
			i++
		case '\\':
			// Potential psql metacommand
			rest := line[i:]
			if psqlTerminatorRE.MatchString(rest) {
				return strings.TrimRight(line[:i], " \t"), true
			}
			// Other backslash mid-line (e.g., \echo) - also treat as terminator
			// for the SQL portion if there's SQL before it
			if i > 0 && strings.TrimSpace(line[:i]) != "" {
				// Check if it looks like a psql command
				if len(rest) > 1 && (rest[1] >= 'a' && rest[1] <= 'z' || rest[1] >= 'A' && rest[1] <= 'Z') {
					return strings.TrimRight(line[:i], " \t"), false
				}
			}
			i++
		default:
			i++
		}
	}

	return "", false
}

// extractState represents the scanner state for statement splitting.
type extractState int

const (
	stNormal       extractState = iota
	stSingleQuote               // inside '...'
	stEscapeString              // inside E'...' with backslash escapes
	stDoubleQuote               // inside "..."
	stDollarQuote               // inside $$...$$ or $tag$...$tag$
	stBlockComment              // inside /* ... */
)

// splitStatements takes preprocessed lines and splits them into SQL statements
// by scanning for semicolons outside of strings, comments, and dollar-quoted blocks.
func splitStatements(filename string, lines []processedLine) []ExtractedStmt {
	var result []ExtractedStmt
	var buf strings.Builder
	state := stNormal
	blockDepth := 0
	dollarTag := ""
	stmtStartLine := 0
	hasContent := false // true if buf has non-whitespace content
	hasSQL := false     // true if buf has non-comment SQL content
	inCopyData := false
	parenDepth := 0         // tracks nesting depth of (...) in normal SQL
	atomicDepth := 0        // tracks nesting depth of BEGIN ATOMIC ... END blocks
	lastWord := ""          // previous completed word (uppercase)
	currentWord := []byte{} // word currently being accumulated

	// finishWord is called when a non-identifier character is seen in stNormal.
	// It checks whether the last two words form "BEGIN ATOMIC" or the current
	// word is "END" to manage atomicDepth.
	finishWord := func() {
		if len(currentWord) == 0 {
			return
		}
		word := strings.ToUpper(string(currentWord))
		currentWord = currentWord[:0]

		if word == "ATOMIC" && lastWord == "BEGIN" {
			atomicDepth++
		} else if word == "END" && atomicDepth > 0 && parenDepth == 0 {
			atomicDepth--
		}
		lastWord = word
	}

	emit := func() {
		sql := strings.TrimSpace(buf.String())
		if sql != "" && hasSQL {
			result = append(result, ExtractedStmt{
				SQL:        sql,
				File:       filename,
				StartLine:  stmtStartLine,
				HasPsqlVar: ContainsPsqlVariable(sql),
			})
			// Check if this was COPY FROM STDIN
			if copyFromStdinRE.MatchString(sql) {
				inCopyData = true
			}
		}
		buf.Reset()
		hasContent = false
		hasSQL = false
		stmtStartLine = 0
		parenDepth = 0
		atomicDepth = 0
		lastWord = ""
		currentWord = currentWord[:0]
	}

	for _, pl := range lines {
		// Skip empty preprocessed lines (metacommands, copy data)
		if pl.text == "" {
			if inCopyData {
				// The preprocessor already handled \. detection,
				// but we need to know when copy data ends.
				// Check the raw state: if preprocessor set text="" for
				// a non-copy-data reason, we need to keep inCopyData.
				// Actually preprocessor handles inCopyData separately.
				// Let's just check: if we are in copy data and see
				// the terminator signal (next non-empty line), we exit.
				// Simpler: preprocessor already skipped copy data.
				// We just need to reset inCopyData when we get a non-empty line.
				continue
			}
			if pl.terminates {
				// Metacommand that terminates a statement (e.g., \gset on its own line)
				emit()
			}
			continue
		}

		// We got a non-empty line, so COPY data mode is over
		if inCopyData {
			inCopyData = false
		}

		text := pl.text
		n := len(text)

		for i := 0; i < n; i++ {
			ch := text[i]

			// Track statement start line
			if !hasContent && !isSpaceByte(ch) {
				stmtStartLine = pl.lineNum
				hasContent = true
			}

			switch state {
			case stNormal:
				// Track words for BEGIN ATOMIC / END detection
				if isIdentCont(ch) {
					currentWord = append(currentWord, ch)
				} else {
					finishWord()
				}

				switch {
				case ch == ';':
					// Only split on semicolons outside parens and outside BEGIN ATOMIC blocks
					if parenDepth == 0 && atomicDepth == 0 {
						// Statement terminator
						emit()
					} else {
						buf.WriteByte(ch)
					}

				case ch == '(':
					parenDepth++
					hasSQL = true
					buf.WriteByte(ch)

				case ch == ')':
					if parenDepth > 0 {
						parenDepth--
					}
					hasSQL = true
					buf.WriteByte(ch)

				case ch == '-' && i+1 < n && text[i+1] == '-':
					// Line comment: include rest of line in buffer
					buf.WriteString(text[i:])
					i = n // skip rest of line
					// Note: do NOT set hasSQL — comments alone don't make a statement

				case ch == '/' && i+1 < n && text[i+1] == '*':
					// Block comment start
					state = stBlockComment
					blockDepth = 1
					buf.WriteString("/*")
					i++ // skip *
					// Note: do NOT set hasSQL — block comments alone don't make a statement

				case ch == '\'' && i > 0 && (text[i-1] == 'E' || text[i-1] == 'e'):
					// E-string (the E was already written to buf)
					state = stEscapeString
					hasSQL = true
					buf.WriteByte(ch)

				case ch == '\'':
					state = stSingleQuote
					hasSQL = true
					buf.WriteByte(ch)

				case ch == '"':
					state = stDoubleQuote
					hasSQL = true
					buf.WriteByte(ch)

				case ch == '$':
					tag := tryParseDollarTagStr(text, i)
					if tag != "" {
						dollarTag = tag
						state = stDollarQuote
						hasSQL = true
						buf.WriteString(tag)
						i += len(tag) - 1 // -1 because loop increments
					} else {
						hasSQL = true
						buf.WriteByte(ch)
					}

				case ch == '\\' && i+1 < n && text[i+1] == ';':
					// psql \; separator: emit current statement and skip the \;
					emit()
					i++ // skip the ; as well

				default:
					if !isSpaceByte(ch) {
						hasSQL = true
					}
					buf.WriteByte(ch)
				}

			case stSingleQuote:
				buf.WriteByte(ch)
				if ch == '\'' {
					if i+1 < n && text[i+1] == '\'' {
						// Escaped quote ''
						buf.WriteByte('\'')
						i++
					} else {
						state = stNormal
					}
				}

			case stEscapeString:
				buf.WriteByte(ch)
				if ch == '\\' {
					// Backslash escape: consume next char
					if i+1 < n {
						i++
						buf.WriteByte(text[i])
					}
				} else if ch == '\'' {
					if i+1 < n && text[i+1] == '\'' {
						buf.WriteByte('\'')
						i++
					} else {
						state = stNormal
					}
				}

			case stDoubleQuote:
				buf.WriteByte(ch)
				if ch == '"' {
					if i+1 < n && text[i+1] == '"' {
						buf.WriteByte('"')
						i++
					} else {
						state = stNormal
					}
				}

			case stDollarQuote:
				if ch == '$' && i+len(dollarTag) <= n && text[i:i+len(dollarTag)] == dollarTag {
					buf.WriteString(dollarTag)
					i += len(dollarTag) - 1
					state = stNormal
				} else {
					buf.WriteByte(ch)
				}

			case stBlockComment:
				if ch == '/' && i+1 < n && text[i+1] == '*' {
					blockDepth++
					buf.WriteString("/*")
					i++
				} else if ch == '*' && i+1 < n && text[i+1] == '/' {
					blockDepth--
					buf.WriteString("*/")
					i++
					if blockDepth == 0 {
						state = stNormal
					}
				} else {
					buf.WriteByte(ch)
				}
			}
		}

		// Flush trailing word at end of line for BEGIN ATOMIC detection.
		if state == stNormal {
			finishWord()
		}

		// End of line: add newline to buffer if we have content
		// (preserves multi-line statement formatting)
		if hasContent || buf.Len() > 0 {
			buf.WriteByte('\n')
		}

		// Check for psql terminator on this line
		if pl.terminates && state == stNormal {
			emit()
		}
	}

	// Emit any remaining statement
	emit()

	return result
}

// tryParseDollarTagStr tries to parse a dollar-quote tag at position i in s.
// Returns the full tag (e.g., "$$" or "$body$") or "" if not a valid dollar tag.
func tryParseDollarTagStr(s string, i int) string {
	if i >= len(s) || s[i] != '$' {
		return ""
	}

	// Try to match $[tag]$ where tag is [a-zA-Z_][a-zA-Z0-9_]*
	j := i + 1

	// Empty tag: $$
	if j < len(s) && s[j] == '$' {
		return "$$"
	}

	// Tag must start with letter or underscore
	if j >= len(s) {
		return ""
	}
	if !isIdentStart(s[j]) {
		return ""
	}

	// Consume identifier characters
	j++
	for j < len(s) && isIdentCont(s[j]) {
		j++
	}

	// Must end with $
	if j >= len(s) || s[j] != '$' {
		return ""
	}

	return s[i : j+1]
}

func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentCont(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9')
}

func isSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// ContainsPsqlVariable detects psql variable interpolation patterns in SQL text.
// These patterns (:varname, :'varname', :"varname") are expanded by the psql client
// before being sent to the server, so they are not valid SQL and cannot be parsed
// by a standard SQL parser.
//
// The function uses a best-effort approach: it scans through the SQL text tracking
// string literal context (single-quoted, double-quoted, dollar-quoted) and comments
// to avoid false positives from colons inside strings or comments. It correctly
// skips :: (PostgreSQL type cast) and := (PL/pgSQL assignment).
func ContainsPsqlVariable(sql string) bool {
	_, replaced := ReplacePsqlVariables(sql)
	return replaced
}

// ReplacePsqlVariables replaces psql variable interpolation patterns with
// valid SQL placeholders, allowing the statement to be parsed for grammar checking.
// Returns the sanitized SQL and true if any replacements were made.
//
// Replacements:
//
//	:varname   -> psql_var
//	:'varname' -> 'psql_var'
//	:"varname" -> "psql_var"
func ReplacePsqlVariables(sql string) (string, bool) {
	// Quick exit
	if !strings.Contains(sql, ":") {
		return sql, false
	}

	var buf strings.Builder
	buf.Grow(len(sql))

	n := len(sql)
	i := 0
	replaced := false

	for i < n {
		ch := sql[i]

		switch {
		// Single-quoted string
		case ch == '\'':
			start := i
			i++
			for i < n {
				if sql[i] == '\'' {
					i++
					if i < n && sql[i] == '\'' {
						i++
						continue
					}
					break
				}
				if sql[i] == '\\' {
					i++
					if i < n {
						i++
					}
					continue
				}
				i++
			}
			buf.WriteString(sql[start:i])

		// Double-quoted identifier
		case ch == '"':
			start := i
			i++
			for i < n {
				if sql[i] == '"' {
					i++
					if i < n && sql[i] == '"' {
						i++
						continue
					}
					break
				}
				i++
			}
			buf.WriteString(sql[start:i])

		// Dollar-quoted string
		case ch == '$':
			tag := tryParseDollarTagStr(sql, i)
			if tag != "" {
				start := i
				i += len(tag)
				for i < n {
					if sql[i] == '$' && i+len(tag) <= n && sql[i:i+len(tag)] == tag {
						i += len(tag)
						break
					}
					i++
				}
				buf.WriteString(sql[start:i])
			} else {
				buf.WriteByte(ch)
				i++
			}

		// Line comment
		case ch == '-' && i+1 < n && sql[i+1] == '-':
			start := i
			i += 2
			for i < n && sql[i] != '\n' {
				i++
			}
			buf.WriteString(sql[start:i])

		// Block comment
		case ch == '/' && i+1 < n && sql[i+1] == '*':
			start := i
			i += 2
			depth := 1
			for i < n && depth > 0 {
				if sql[i] == '/' && i+1 < n && sql[i+1] == '*' {
					depth++
					i += 2
				} else if sql[i] == '*' && i+1 < n && sql[i+1] == '/' {
					depth--
					i += 2
				} else {
					i++
				}
			}
			buf.WriteString(sql[start:i])

		// Colon: check for psql variable
		case ch == ':':
			isVar := false
			replType := 0 // 0: :var, 1: :'var', 2: :"var"
			varEnd := i + 1

			if i+1 < n {
				next := sql[i+1]
				// :: type cast
				if next == ':' {
					buf.WriteString("::")
					i += 2
					continue
				}
				// := assignment
				if next == '=' {
					buf.WriteString(":=")
					i += 2
					continue
				}
				// :varname
				if isIdentStart(next) {
					isVar = true
					replType = 0
					varEnd = i + 2
					for varEnd < n && isIdentCont(sql[varEnd]) {
						varEnd++
					}
				} else if next == '\'' {
					// :'varname'
					isVar = true
					replType = 1
					// Scan until closing quote
					varEnd = i + 2 // skip :'
					// Note: psql variable names inside quotes can contain anything but quote?
					// Standard psql says "The variable name must be surrounded by single quotes"
					// We'll just scan until the next single quote.
					for varEnd < n {
						if sql[varEnd] == '\'' {
							// Check for escaped quote? psql variables usually simple.
							// Assuming simple variable names for regression tests.
							varEnd++
							break
						}
						varEnd++
					}
				} else if next == '"' {
					// :"varname"
					isVar = true
					replType = 2
					varEnd = i + 2 // skip :"
					for varEnd < n {
						if sql[varEnd] == '"' {
							varEnd++
							break
						}
						varEnd++
					}
				}
			}

			if isVar {
				replaced = true
				switch replType {
				case 0:
					buf.WriteString("psql_var")
				case 1:
					buf.WriteString("'psql_var'")
				case 2:
					buf.WriteString("\"psql_var\"")
				}
				i = varEnd
			} else {
				buf.WriteByte(ch)
				i++
			}

		default:
			buf.WriteByte(ch)
			i++
		}
	}

	return buf.String(), replaced
}
