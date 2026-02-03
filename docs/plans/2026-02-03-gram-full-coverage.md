# Grammar Full Coverage Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Close all confirmed grammar gaps, improve the SQL extractor to handle psql variables/backslash commands, and add missing SecLabelStmt alternatives — reaching 99.5%+ regression pass rate.

**Architecture:** Surgical changes to `parser/gram.y` for each grammar gap (adding missing alternatives to existing rules), improvements to `parser/pgregress/extract.go` for psql variable and backslash command detection, and test coverage via existing `tests/grammar/` fixture framework plus `TestPGRegressStats`.

**Tech Stack:** Go, goyacc, PostgreSQL 17.7 grammar reference

---

## Track A: Grammar Fixes

### Task 1: Add TRANSFORM to createfunc_opt_item

**Files:**
- Modify: `parser/gram.y:4647-4663`
- Test: `tests/grammar/sql/phase5_plan.sql`

**Step 1: Write the failing test**

Add to `tests/grammar/sql/phase5_plan.sql`:
```sql
-- Task 1: TRANSFORM in CREATE FUNCTION
CREATE FUNCTION my_func(jsonb) RETURNS int LANGUAGE sql TRANSFORM FOR TYPE jsonb AS $$ SELECT 1; $$;
```

**Step 2: Run test to verify it fails**

Run: `cd pgparser && go test ./tests/grammar/ -run TestGrammarFixtures/phase5_plan.sql -v`
Expected: FAIL with parse error

**Step 3: Add transform_type_list rule and AS func_as rule**

In `parser/gram.y`, first add the `transform_type_list` rule (after `createfunc_opt_item`, around line 4664):

```yacc
transform_type_list:
	FOR TYPE_P Typename
		{ $$ = makeList($3) }
	| transform_type_list ',' FOR TYPE_P Typename
		{ $$ = appendList($1, $5) }
	;
```

Then add the `func_as` rule:

```yacc
func_as:
	Sconst
		{ $$ = makeList(&nodes.String{Str: $1}) }
	| Sconst ',' Sconst
		{ $$ = makeList(&nodes.String{Str: $1}, &nodes.String{Str: $3}) }
	;
```

Then modify `createfunc_opt_item` (lines 4647-4663) by adding two new alternatives before `| common_func_opt_item`:

```yacc
createfunc_opt_item:
	AS func_as
		{
			$$ = &nodes.DefElem{
				Defname: "as",
				Arg:     $2,
			}
		}
	| LANGUAGE NonReservedWord_or_Sconst
		{
			$$ = &nodes.DefElem{
				Defname: "language",
				Arg:     &nodes.String{Str: $2},
			}
		}
	| TRANSFORM transform_type_list
		{
			$$ = &nodes.DefElem{
				Defname: "transform",
				Arg:     $2,
			}
		}
	| WINDOW
		{
			$$ = &nodes.DefElem{
				Defname: "window",
				Arg:     &nodes.Integer{Ival: 1},
			}
		}
	| common_func_opt_item { $$ = $1 }
	;
```

Add `%type` declarations if not present:
- `%type <list> transform_type_list`
- `%type <list> func_as`

**Step 4: Regenerate parser and run test**

Run: `cd pgparser && make generate-parser && go test ./tests/grammar/ -run TestGrammarFixtures/phase5_plan.sql -v`
Expected: PASS

**Step 5: Run full test suite**

Run: `cd pgparser && go test ./... -count=1`
Expected: No regressions

**Step 6: Commit**

```bash
cd pgparser && git add parser/gram.y tests/grammar/sql/phase5_plan.sql && git commit -m "feat(gram): add TRANSFORM and AS func_as to createfunc_opt_item"
```

---

### Task 2: Add LOGGED/UNLOGGED to SeqOptElem

**Files:**
- Modify: `parser/gram.y:6831-6888`
- Test: `tests/grammar/sql/phase5_plan.sql`

**Step 1: Write the failing test**

Add to `tests/grammar/sql/phase5_plan.sql`:
```sql
-- Task 2: LOGGED/UNLOGGED sequences
CREATE SEQUENCE s_logged AS integer LOGGED;
CREATE SEQUENCE s_unlogged AS integer UNLOGGED;
```

**Step 2: Run test to verify it fails**

Run: `cd pgparser && go test ./tests/grammar/ -run TestGrammarFixtures/phase5_plan.sql -v`
Expected: FAIL with parse error

**Step 3: Add LOGGED and UNLOGGED alternatives**

In `parser/gram.y`, add two alternatives to `SeqOptElem` after the `RESTART opt_with NumericOnly` alternative (before line 6888's `;`):

```yacc
	| LOGGED
		{
			$$ = makeDefElem("logged", &nodes.Boolean{Boolval: true})
		}
	| UNLOGGED
		{
			$$ = makeDefElem("logged", &nodes.Boolean{Boolval: false})
		}
```

**Step 4: Regenerate parser and run test**

Run: `cd pgparser && make generate-parser && go test ./tests/grammar/ -run TestGrammarFixtures/phase5_plan.sql -v`
Expected: PASS

**Step 5: Run full test suite**

Run: `cd pgparser && go test ./... -count=1`
Expected: No regressions

**Step 6: Commit**

```bash
cd pgparser && git add parser/gram.y tests/grammar/sql/phase5_plan.sql && git commit -m "feat(gram): add LOGGED/UNLOGGED to SeqOptElem"
```

---

### Task 3: Fix hash_partbound to support list

**Files:**
- Modify: `parser/gram.y:1628-1656`
- Test: `tests/grammar/sql/phase4_ddl.sql`

**Step 1: Write the failing test**

Add to `tests/grammar/sql/phase4_ddl.sql`:
```sql
-- Task 3: hash_partbound list
CREATE TABLE hp_parent (a int) PARTITION BY HASH (a);
CREATE TABLE hp_child PARTITION OF hp_parent FOR VALUES WITH (MODULUS 4, REMAINDER 0);
```

**Step 2: Run test to verify it passes (baseline)**

Run: `cd pgparser && go test ./tests/grammar/ -run TestGrammarFixtures/phase4_ddl.sql -v`
Expected: PASS (2-element form already works)

**Step 3: Write the test for multi-element form**

This is the actual failing case. However, PostgreSQL's grammar only uses `hash_partbound` as a list of `hash_partbound_elem` — the semantics of multiple modulus/remainder pairs is not standard. The PG grammar just builds a list and validates later. Let's match PG's grammar structure.

Replace `hash_partbound` rule (lines 1628-1642) with list-building form:

```yacc
hash_partbound:
	hash_partbound_elem
		{
			$$ = makeList($1)
		}
	| hash_partbound ',' hash_partbound_elem
		{
			$$ = appendList($1, $3)
		}
	;
```

And update `hash_partbound_elem` (lines 1644-1656) to return a `*nodes.DefElem`:

```yacc
hash_partbound_elem:
	NonReservedWord Iconst
		{
			$$ = makeDefElem($1, &nodes.Integer{Ival: $2})
		}
	;
```

Then update the `ForValues` rule reference (line 1622-1625) that uses `hash_partbound`. Currently:
```yacc
	| FOR VALUES WITH '(' hash_partbound ')'
		{
			$$ = $5.(*nodes.PartitionBoundSpec)
		}
```

This needs to change to construct a `PartitionBoundSpec` from the list of DefElems:

```yacc
	| FOR VALUES WITH '(' hash_partbound ')'
		{
			n := &nodes.PartitionBoundSpec{
				Strategy: 'h',
				Location: -1,
			}
			for _, item := range $5.(*nodes.List).Items {
				de := item.(*nodes.DefElem)
				switch de.Defname {
				case "modulus":
					n.Modulus = int(de.Arg.(*nodes.Integer).Ival)
				case "remainder":
					n.Remainder = int(de.Arg.(*nodes.Integer).Ival)
				}
			}
			$$ = n
		}
```

Update `%type` if needed: `hash_partbound` should be `<list>`, `hash_partbound_elem` should be `<node>`.

**Step 4: Regenerate parser and run test**

Run: `cd pgparser && make generate-parser && go test ./tests/grammar/ -run TestGrammarFixtures/phase4_ddl.sql -v`
Expected: PASS

**Step 5: Run full test suite**

Run: `cd pgparser && go test ./... -count=1`
Expected: No regressions

**Step 6: Commit**

```bash
cd pgparser && git add parser/gram.y tests/grammar/sql/phase4_ddl.sql && git commit -m "feat(gram): refactor hash_partbound to support arbitrary-length list"
```

---

### Task 4: Add missing SecLabelStmt alternatives

**Files:**
- Modify: `parser/gram.y:12858-12905`
- Test: `tests/grammar/sql/phase4_admin.sql`

**Step 1: Write the failing test**

Add to `tests/grammar/sql/phase4_admin.sql`:
```sql
-- Task 4: SecLabelStmt missing alternatives
SECURITY LABEL ON FUNCTION my_func(int) IS 'secret';
SECURITY LABEL ON PROCEDURE my_proc(text) IS 'secret';
SECURITY LABEL ON ROUTINE my_routine(int, text) IS 'secret';
SECURITY LABEL ON AGGREGATE my_agg(int) IS 'secret';
SECURITY LABEL ON LARGE OBJECT 12345 IS 'secret';
```

**Step 2: Run test to verify it fails**

Run: `cd pgparser && go test ./tests/grammar/ -run TestGrammarFixtures/phase4_admin.sql -v`
Expected: FAIL with parse error on FUNCTION/PROCEDURE/ROUTINE/AGGREGATE/LARGE OBJECT

**Step 3: Add 5 missing alternatives to SecLabelStmt**

In `parser/gram.y`, add before the closing `;` of `SecLabelStmt` (line 12905), replacing the TODO comment at line 12904:

```yacc
	| SECURITY LABEL opt_provider ON AGGREGATE aggregate_with_argtypes IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_AGGREGATE,
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
	| SECURITY LABEL opt_provider ON FUNCTION function_with_argtypes IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_FUNCTION,
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
	| SECURITY LABEL opt_provider ON LARGE_P OBJECT_P NumericOnly IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_LARGEOBJECT,
				Object:   $7,
				Provider: $3,
				Label:    $9,
			}
		}
	| SECURITY LABEL opt_provider ON PROCEDURE function_with_argtypes IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_PROCEDURE,
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
	| SECURITY LABEL opt_provider ON ROUTINE function_with_argtypes IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_ROUTINE,
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
```

Verify that `nodes.OBJECT_AGGREGATE`, `nodes.OBJECT_FUNCTION`, `nodes.OBJECT_LARGEOBJECT`, `nodes.OBJECT_PROCEDURE`, `nodes.OBJECT_ROUTINE` constants exist. If not, add them to the ObjectType enum in `nodes/`.

**Step 4: Regenerate parser and run test**

Run: `cd pgparser && make generate-parser && go test ./tests/grammar/ -run TestGrammarFixtures/phase4_admin.sql -v`
Expected: PASS

**Step 5: Run full test suite**

Run: `cd pgparser && go test ./... -count=1`
Expected: No regressions

**Step 6: Commit**

```bash
cd pgparser && git add parser/gram.y tests/grammar/sql/phase4_admin.sql && git commit -m "feat(gram): add FUNCTION/PROCEDURE/ROUTINE/AGGREGATE/LARGE OBJECT to SecLabelStmt"
```

---

## Track B: Extractor Improvements

### Task 5: Add psql variable detection to extractor

**Files:**
- Modify: `parser/pgregress/extract.go:1-553`
- Test: `parser/pgregress/extract_test.go`

**Step 1: Write the failing test**

Add test cases to `parser/pgregress/extract_test.go` that verify psql variable statements are detected:

```go
func TestPsqlVariableDetection(t *testing.T) {
	content := []byte(`SELECT 1;
SELECT :varname;
SELECT :'filename';
INSERT INTO t VALUES (:x, :y);
SELECT 1;
`)
	stmts := ExtractStatements("test.sql", content)
	// Statement indices 1, 2, 3 contain psql variables
	for i, stmt := range stmts {
		hasPsql := containsPsqlVariable(stmt.SQL)
		switch i {
		case 0, 4:
			if hasPsql {
				t.Errorf("stmt %d should NOT have psql variable: %s", i, stmt.SQL)
			}
		case 1, 2, 3:
			if !hasPsql {
				t.Errorf("stmt %d SHOULD have psql variable: %s", i, stmt.SQL)
			}
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestPsqlVariableDetection -v`
Expected: FAIL (containsPsqlVariable doesn't exist yet)

**Step 3: Implement psql variable detection**

Add to `parser/pgregress/extract.go`:

```go
// psqlVariableRE detects psql variable interpolation patterns.
// Matches :varname, :'varname', :"varname" but not ::type_cast or :=.
var psqlVariableRE = regexp.MustCompile(`(?:^|[^:]):[A-Za-z_][A-Za-z0-9_]*(?:'|"|(?:[^:]|$))`)

// containsPsqlVariable returns true if the SQL text contains psql variable
// interpolation patterns that cannot be parsed as standard SQL.
func containsPsqlVariable(sql string) bool {
	// Quick check for colon (most SQL won't have this)
	if !strings.Contains(sql, ":") {
		return false
	}
	// Check for :varname patterns (not :: type cast, not := assignment)
	// Patterns: :name, :'name', :"name"
	for i := 0; i < len(sql); i++ {
		if sql[i] != ':' {
			continue
		}
		// Skip :: (type cast)
		if i+1 < len(sql) && sql[i+1] == ':' {
			i++
			continue
		}
		// Skip := (assignment)
		if i+1 < len(sql) && sql[i+1] == '=' {
			continue
		}
		// Check for :identifier or :'...' or :"..."
		if i+1 < len(sql) {
			next := sql[i+1]
			if next == '\'' || next == '"' {
				return true
			}
			if (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || next == '_' {
				// Verify it's not inside a string literal (simple heuristic)
				return true
			}
		}
	}
	return false
}
```

Export `ContainsPsqlVariable` (capitalized) for use in test framework.

**Step 4: Run test to verify it passes**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestPsqlVariableDetection -v`
Expected: PASS

**Step 5: Commit**

```bash
cd pgparser && git add parser/pgregress/extract.go parser/pgregress/extract_test.go && git commit -m "feat(extract): add psql variable interpolation detection"
```

---

### Task 6: Integrate psql variable detection into regression test framework

**Files:**
- Modify: `parser/pgregress/regress_test.go`
- Modify: `parser/pgregress/extract.go` (export function)

**Step 1: Understand current framework**

The `TestPGRegress` in `regress_test.go` iterates over extracted statements, parses each, and compares against `known_failures.json`. Statements with psql variables fail at parse time and are counted as failures.

**Step 2: Add psql variable classification to ExtractedStmt**

Add a `HasPsqlVar` field to `ExtractedStmt`:

```go
type ExtractedStmt struct {
	SQL        string
	File       string
	StartLine  int
	HasPsqlVar bool // true if SQL contains psql variable interpolation
}
```

Set `HasPsqlVar` in `splitStatements()` when building each statement, using `containsPsqlVariable()`.

**Step 3: Update regress_test.go to skip psql variable statements**

In `TestPGRegress` and `TestPGRegressStats`, when a statement has `HasPsqlVar == true`, skip it instead of counting it as a failure. Track separately:

```go
var psqlVarSkipped int
// ... in loop:
if stmt.HasPsqlVar {
    psqlVarSkipped++
    continue
}
```

**Step 4: Run full regression test**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestPGRegressStats -v`
Expected: ~165 fewer failures (reclassified as psql-dependent)

**Step 5: Update known_failures.json**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestPGRegress -update`
This should regenerate `known_failures.json` with the psql variable statements removed.

**Step 6: Commit**

```bash
cd pgparser && git add parser/pgregress/ && git commit -m "feat(extract): skip psql variable statements in regression tests"
```

---

### Task 7: Improve backslash command handling in extractor

**Files:**
- Modify: `parser/pgregress/extract.go:64-73`
- Test: `parser/pgregress/extract_test.go`

**Step 1: Write the failing test**

```go
func TestBackslashCommandHandling(t *testing.T) {
	content := []byte(`\\set varname value
SELECT 1;
\\getenv envvar PG_LIBDIR
SELECT 2;
\\lo_import 'file.txt'
SELECT 3;
`)
	stmts := ExtractStatements("test.sql", content)
	// Should get exactly 3 clean SQL statements
	var sqlStmts []string
	for _, s := range stmts {
		if s.SQL != "" {
			sqlStmts = append(sqlStmts, s.SQL)
		}
	}
	if len(sqlStmts) != 3 {
		t.Fatalf("expected 3 statements, got %d: %v", len(sqlStmts), sqlStmts)
	}
}
```

**Step 2: Run test to verify current behavior**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestBackslashCommandHandling -v`
Check if it already passes.

**Step 3: If needed, improve backslash detection**

The current code (line 65) checks `strings.HasPrefix(trimmedLeft, "\\")`. This should already handle `\set`, `\getenv`, `\lo_import`, etc. If it does, no change needed.

If some backslash commands are being merged with subsequent SQL, add more specific patterns:

```go
// psqlMetacmdRE matches psql metacommands at the start of a line.
var psqlMetacmdRE = regexp.MustCompile(`^\\[a-zA-Z]`)
```

And replace the `strings.HasPrefix(trimmedLeft, "\\")` check with the regex.

**Step 4: Run test**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestBackslashCommandHandling -v`
Expected: PASS

**Step 5: Run full suite and update known_failures**

Run: `cd pgparser && go test ./... -count=1`

**Step 6: Commit**

```bash
cd pgparser && git add parser/pgregress/ && git commit -m "fix(extract): improve psql backslash command detection"
```

---

## Track C: Final Verification

### Task 8: Full regression verification and known_failures update

**Files:**
- Modify: `parser/pgregress/known_failures.json`

**Step 1: Run full test suite**

Run: `cd pgparser && go test ./... -count=1 -v 2>&1 | tail -50`

**Step 2: Run regression stats**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestPGRegressStats -v`
Expected: Pass rate should be above 99.5%

**Step 3: Update known_failures.json**

Run: `cd pgparser && go test ./parser/pgregress/ -run TestPGRegress -update`

**Step 4: Review and commit**

```bash
cd pgparser && git add parser/pgregress/known_failures.json && git commit -m "chore: update known_failures.json after grammar and extractor fixes"
```

---

## Summary of Expected Impact

| Task | Change | Failures Fixed |
|------|--------|----------------|
| 1 | TRANSFORM + AS func_as in createfunc_opt_item | ~3 |
| 2 | LOGGED/UNLOGGED in SeqOptElem | ~2 |
| 3 | hash_partbound list support | ~1 |
| 4 | SecLabelStmt 5 alternatives | ~0 (proactive) |
| 5-6 | psql variable detection | ~165 reclassified |
| 7 | Backslash command handling | ~25 reclassified |
| 8 | Final verification | — |

**Expected final pass rate**: 99.5%+ (excluding psql-dependent and intentionally-invalid SQL)
