# Grammar Full Coverage Design (99%+ Target)

Date: 2026-02-03

## Current State

- **Regression pass rate: 99.2%** (44,945 / 45,321)
- **Remaining failures: 376**
- **gram.y: 17,706 lines, 687 rules** (vs PG's 19,520 lines, 727 rules)
- **59 PG rules missing** from Go grammar (but most are handled via inlining)
- **67 rules** with different alternative counts between PG and Go

## Failure Root Cause Breakdown

| Category | Count | Track |
|----------|------:|-------|
| psql variable interpolation (`:var`, `\set`, `\gset`) | ~165 | Extractor |
| Intentionally invalid SQL (errors.sql) | ~42 | No fix needed |
| Extractor statement-splitting bugs | ~25 | Extractor |
| psql backslash commands in SQL files | ~53 | Extractor |
| UESCAPE syntax (`U&'...' UESCAPE '!'`) | 9 | Lexer |
| Real grammar gaps | ~15 | gram.y |
| Ambiguous / investigation needed | ~20 | Mixed |

## Three Tracks

### Track A: Fix Real Grammar Gaps (~15 failures → gram.y)

These are confirmed parse failures from actual SQL that the parser rejects:

#### A1: `SECURITY LABEL ON PROCEDURE/ROUTINE` (2 failures)

Missing 5 alternatives in `SecLabelStmt`: PROCEDURE, ROUTINE, LARGE OBJECT, TABLESPACE, PUBLICATION.

**PG gram.y reference** (line ~6900):
```
SecLabelStmt:
    SECURITY LABEL opt_provider ON PROCEDURE function_with_argtypes IS security_label_value
    SECURITY LABEL opt_provider ON ROUTINE function_with_argtypes IS security_label_value
    SECURITY LABEL opt_provider ON LARGE_P OBJECT_P NumericOnly IS security_label_value
    SECURITY LABEL opt_provider ON TABLESPACE name IS security_label_value
    SECURITY LABEL opt_provider ON PUBLICATION name IS security_label_value
```

**Files:** `parser/gram.y` SecLabelStmt rule

#### A2: `CREATE TRIGGER ... OF column` (1 failure)

Column-specific triggers: `CREATE TRIGGER trig BEFORE INSERT OF a ON t1`.

PG's TriggerEvents allows `INSERT OF column_list`. Go grammar likely missing the `OF columnList` variant in the trigger event production.

**Files:** `parser/gram.y` TriggerEvents rule

#### A3: `CREATE FUNCTION ... TRANSFORM FOR TYPE` (1 failure)

Missing `TRANSFORM transform_type_list` alternative in `createfunc_opt_item`.

**PG gram.y reference** (line ~9766):
```
createfunc_opt_item:
    ...
    | TRANSFORM transform_type_list { ... }
```

Also need `transform_type_list` rule if missing.

**Files:** `parser/gram.y` createfunc_opt_item rule

#### A4: `CREATE SEQUENCE ... LOGGED/UNLOGGED` (2 failures)

Missing `LOGGED` and `UNLOGGED` alternatives in `SeqOptElem`.

**PG gram.y reference** (line ~5860):
```
SeqOptElem:
    ...
    | LOGGED { ... }
    | UNLOGGED { ... }
```

**Files:** `parser/gram.y` SeqOptElem rule

#### A5: `TABLESAMPLE` on derived table / subquery (1 failure)

`SELECT * FROM (SELECT ...) q TABLESAMPLE BERNOULLI(5)` fails. The grammar only allows TABLESAMPLE on `relation_expr`, not on parenthesized subqueries.

PG's `table_ref` has:
```
    relation_expr opt_alias_clause tablesample_clause
```
but the Go version may not apply tablesample correctly to aliased derived tables.

**Files:** `parser/gram.y` table_ref rule

#### A6: `hash_partbound` multiple modulus/remainder (1 failure)

`FOR VALUES WITH (MODULUS 4, REMAINDER 0, MODULUS 8, REMAINDER 1)` — the `hash_partbound` rule needs a list-building alternative.

**PG gram.y reference** (line ~4193):
```
hash_partbound:
    hash_partbound_elem
    | hash_partbound ',' hash_partbound_elem
```

**Files:** `parser/gram.y` hash_partbound rule

#### A7: `CREATE LANGUAGE ... VALIDATOR` (1 failure)

Missing `validator_clause` rule. PG has:
```
validator_clause:
    VALIDATOR handler_name { ... }
    | NO VALIDATOR { ... }
```

**Files:** `parser/gram.y` CreatePLangStmt + validator_clause

#### A8: Numeric literals with underscores: `_1_000.5` (1 failure)

This is a **lexer issue**, not a grammar issue. PG17 supports `_` in numeric literals. The Go lexer doesn't recognize this.

**Files:** `parser/lexer.go` (numeric literal scanning)

#### A9: `SELECT * FROM rank() OVER (...)` — function in FROM (1 failure)

This should be a semantic error, not a syntax error. The grammar should parse it (it's valid syntax) but the function call with OVER clause can't be used in FROM.

**Files:** `parser/gram.y` table_ref / func_table rule

#### A10: `WINDOW_partition_err` — comma after PARTITION BY (1 failure)

`PARTITION BY four, ORDER BY ten` — trailing comma before ORDER BY. This is intentionally invalid SQL but the parser should give a better error.

**Investigation needed** to decide if this is a grammar fix or acceptable.

#### A11: `CREATE STATISTICS` error productions (3 failures)

Incomplete `CREATE STATISTICS` syntax (missing FROM, missing ON, etc.) intended to produce error messages. PG gram.y has these as error-reporting productions.

**Files:** `parser/gram.y` CreateStatsStmt rule

#### A12: `CREATE FOREIGN TABLE ft1 ()` — empty column list without SERVER (1 failure)

`CREATE FOREIGN TABLE ft1 ()` (no SERVER clause) should produce an error but currently fails as syntax error. PG allows this syntax at parse level and rejects at semantic level.

**Files:** `parser/gram.y` CreateForeignTableStmt rule

#### A13: `ALTER TABLE ... RENAME TRIGGER a ON ONLY parent` (1 failure from triggers.sql)

`alter trigger a on only grandparent rename to b` — missing the `ONLY` variant in trigger rename.

**Files:** `parser/gram.y` RenameStmt or relevant rule

### Track B: Fix Extractor (~190 failures → extract.go)

The SQL statement extractor in `parser/pgregress/extract.go` mishandles several psql features:

#### B1: psql Variable Interpolation (~165 failures)

Statements like `:show_data`, `:init_range_parted`, `SELECT :varname`, `:'filename'` are psql variable expansions. The extractor should either:
- **Option 1**: Skip statements containing unresolved psql variables (mark them as "requires psql preprocessing")
- **Option 2**: Replace `:varname` with placeholder values for parse testing

**Recommendation**: Option 1 — skip and classify as "psql-dependent", remove from failure count.

#### B2: psql Backslash Commands (~25 failures)

`\set`, `\getenv`, `\lo_import`, `\lo_unlink`, `\copy`, etc. Some of these are followed by SQL on the next line, causing the extractor to merge psql commands with SQL.

Fix: Improve `psqlTerminatorRE` and line-by-line psql command detection.

**Files:** `parser/pgregress/extract.go`

### Track C: Structural Alignment (59 missing rules → gram.y)

Even though most missing rules don't cause test failures (because they're inlined), having structural alignment with PG's gram.y makes future maintenance easier. These fall into categories:

#### C1: Re-extract inlined rules (refactor, no behavioral change)

Rules that Go inlined but PG keeps separate. Re-extracting them would make the grammar easier to maintain:

- `privilege_target` (21 alternatives) — extract from `GrantStmt`/`RevokeStmt`
- `drop_type_name` (8 alternatives) — extract from `DropStmt`
- `in_expr` (2 alternatives) — extract from `a_expr`
- `character` / `CharacterWithLength` / `CharacterWithoutLength` — extract from `Character`/`ConstCharacter`

**Risk**: Moderate — refactoring grammar rules can introduce shift/reduce conflicts.
**Recommendation**: Defer to a separate cleanup pass.

#### C2: Add truly missing rules (no inlined equivalent)

- `bare_label_keyword` (452 alternatives) — PG17 keyword classification system
- `BareColLabel` (2 alternatives) — uses `bare_label_keyword`
- `PLAssignStmt` / `PLpgSQL_Expr` / `plassign_target` / `plassign_equals` — PL/pgSQL-specific, not needed for SQL parser
- `parse_toplevel` / `toplevel_stmt` — Go uses `stmtblock` instead
- `CreateAssertionStmt` — PG doesn't implement assertions; grammar exists only for error message
- `copy_options` / `copy_generic_opt_*` — Go handles COPY options differently

**Recommendation**: Only add `bare_label_keyword` / `BareColLabel` if testing shows keyword classification issues. Skip PL/pgSQL rules. Add `CreateAssertionStmt` for error message parity.

#### C3: Add missing alternatives to existing rules

Low-impact alternatives missing from existing rules:

| Rule | Missing Alternatives | Impact |
|------|---------------------|--------|
| `SecLabelStmt` | +5 object types | Covered in Track A |
| `DefineStmt` | +3 (but Go extracted as separate stmts) | None |
| `createfunc_opt_item` | +2 (TRANSFORM, FunctionSetResetClause) | Covered in Track A |
| `SeqOptElem` | +2 (LOGGED/UNLOGGED) | Covered in Track A |
| `PublicationObjSpec` | +2 (ColId, ColId *) | Minor |
| `json_table_column_definition` | +1 (JTC_FORMATTED) | Minor |
| `zone_value` | +2 (ConstInterval, I_or_F_const) | Minor |
| `json_name_and_value` | +1 | Minor |
| `type_func_name_keyword` | +1 keyword | Minor |

## Phased Execution

### Phase 1: Grammar Fixes (A1-A7, A11-A13) — ~15 changes

Fix all confirmed grammar parse failures. Each is a small, surgical change to gram.y.

**Expected result**: 99.2% → 99.5% (fix ~15 real failures)

### Phase 2: Extractor Improvements (B1-B2) — extract.go changes

1. Detect and skip psql variable-dependent statements
2. Improve psql backslash command handling

**Expected result**: 99.5% → 99.9% (reclassify ~190 as psql-dependent)

### Phase 3: Lexer Fix (A8) — UESCAPE + numeric underscores

1. Support `_` in numeric literals (PG17 feature)
2. Support `U&'...' UESCAPE` syntax

**Expected result**: Fix 10 additional failures

### Phase 4: Structural Alignment (C2-C3) — optional maintenance

Add remaining missing alternatives for structural parity. No behavioral impact.

## Testing Strategy

- Each gram.y change: write a targeted test in `tests/grammar/` exercising the specific syntax
- Each extractor fix: verify via `TestPGRegressStats` improvement
- Final: run `go test ./... -count=1` and verify no regressions

## Acceptance Criteria

- `go build ./...` and `go test ./...` pass
- Regression pass rate ≥ 99.5% (excluding psql-dependent and intentionally-invalid)
- All 59 structurally-missing rules either added or documented as intentionally inlined
