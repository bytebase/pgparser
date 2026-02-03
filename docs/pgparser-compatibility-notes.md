# pgparser Compatibility Notes (PG parser + nodeToString)

This document captures our current analysis of gaps between pgparser and
PostgreSQL's native parser, with emphasis on **100% syntax compatibility** and
**nodeToString output parity**.

## Final Goal

1) **Syntax coverage**: pgparser accepts the same SQL grammar as PG.
2) **AST shape parity**: nodes and fields match PG parse trees.
3) **nodeToString parity**: output strings exactly match PG.

All three are required; parsing success alone is insufficient.

## Current State (High-Level)

- Grammar coverage is incomplete (many missing rule alternatives).
- Some rules are refactored/flattened due to goyacc vs bison differences.
- AST shapes sometimes diverge (e.g., qualified_name, row constructors).
- Location tracking is largely missing (goyacc lacks %locations).

## Why PG vs goyacc Matters

- **%locations**: bison emits token positions; goyacc doesn't. Many nodes in
  pgparser use `Location: -1`, so nodeToString cannot match PG.
- **Conflict resolution**: PG grammar relies on bison behavior and `%prec`.
  goyacc sometimes requires rule inlining/rewrites; syntax may be equivalent
  but structure differs.
- **Error productions**: PG has error-recovery branches that are absent or
  simplified in pgparser.

## Expression System (Focused Findings)

### Missing or unsupported syntax (parsing fails)

- `UNIQUE (SELECT ...)` in `a_expr` (PG-only).
- `TREAT (expr AS type)` (PG-only in func_expr_common_subexpr).
- `ROW()` empty constructor (PG-only in row).
- `OPERATOR(schema.op) ANY/ALL` (PG-only in subquery_Op).
- `(SELECT ...)[...]` / `(SELECT ...).field` (PG-only in c_expr).
- `NATIONAL CHAR` / `NCHAR` type keywords (PG-only in character rules).

### Structural/AST differences (parse succeeds, output diverges)

- `qualified_name`: PG produces `RangeVar`; pgparser produces `[]String`.
- `ConstCharacter` / `ConstTypename`: typmod defaults handled differently.
- `JsonType`: pgparser inlines JSON token; PG uses SystemTypeName path.
- `AexprConst`: pgparser lacks some PG-only validation paths for type modifiers.

### Syntax rewrites that are *likely* equivalent

- `IN` handling: PG uses `in_expr` helper; pgparser inlines alternatives in
  `a_expr`. Generally equivalent but verify `OperName` population.
- `subquery_Op`: pgparser expands basic operators instead of `all_Op/MathOp`.

## nodeToString Parity Challenges

### Root causes

- **Location fields**: PG sets location via `@n` (bison); pgparser lacks this.
- **AST shape**: different node types or list shapes produce different output.
- **Field order/omissions**: outfuncs order and default values must match PG.

### Requirements to reach parity

- Implement **location propagation** in lexer/parser (simulate `%locations`).
- Align **node constructors** to match PG raw parse trees.
- Ensure `nodes/outfuncs.go` matches PG `outfuncs.c` ordering and defaults.

## Recommended Roadmap

### Phase 1: Syntax coverage

- Continue executing `docs/plans/2026-01-30-gram-completion.md`.
- Focus on missing/high-impact rules: expressions, window, table_ref, COPY,
  partitioning, JSON/XML.

### Phase 2: AST shape parity

- Systematically compare PG vs pgparser raw parse trees.
- Fix structural differences (RangeVar, row constructors, operator names,
  SetToDefault, TypeName, etc.).

### Phase 3: nodeToString parity

- Add token location tracking (simulated `%locations`).
- Ensure all nodes emit fields in PG order with matching defaults.

## Tooling We Added

We added a baseline PG parse helper and a diff tool to compare outputs.

### PG helper (raw_parser + nodeToString)

- `tools/pg_parse_helper/pg_parse_helper.c`
- `tools/pg_parse_helper/build.sh`

Build (requires PG source configured):

```
PG_SRC=~/Github/postgres tools/pg_parse_helper/build.sh
```

### Diff harness (PG output vs pgparser output)

- `tools/pg_parse_diff/main.go`
- `tools/pg_parse_diff/smoke.sql`

Usage:

```
go run tools/pg_parse_diff/main.go --file tools/pg_parse_diff/smoke.sql
```

Or for regression SQL:

```
go run tools/pg_parse_diff/main.go --dir parser/pgregress/testdata/sql
```

This compares **raw parse trees** (same level as pgparser's `parser.Parse`).

## Next Focus Areas

1) Close the high-impact missing syntax in expressions and table refs.
2) Fix AST shape divergences that create nodeToString diffs (e.g., RangeVar).
3) Implement location propagation to match PG output.

