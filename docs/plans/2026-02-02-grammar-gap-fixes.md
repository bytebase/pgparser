# Grammar Gap Fixes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Close the remaining grammar gaps in pgparser's gram.y to match PostgreSQL 17's grammar, fixing ~31 real parser failures and ~36 extractor failures.

**Architecture:** Each task adds missing productions to `parser/gram.y` (or fixes `parser/pgregress/extract.go`), regenerates the parser, and validates via regression tests. Tasks are ordered by ease and impact (easiest first).

**Tech Stack:** Go, goyacc, PostgreSQL 17 gram.y reference at `~/Github/postgres/src/backend/parser/gram.y`

---

### Task 1: Add bare `EMPTY` to `json_behavior_type`

**Files:**
- Modify: `parser/gram.y:10449-10457`
- Test: `parser/pgregress/testdata/sql/sqljson_jsontable.sql` (indices 3, 114), `sqljson_queryfuncs.sql` (indices 137, 290)

**Fixes:** 4 failures

**Step 1: Write a failing test**

Create or extend a test in `parser/parsertest/` that exercises bare `EMPTY ON ERROR`:

```go
// In parser/parsertest/sqljson_test.go (or new file)
func TestJsonBareEmpty(t *testing.T) {
    sqls := []string{
        `SELECT * FROM JSON_TABLE('[]', 'strict $.a' COLUMNS (js2 int PATH '$') EMPTY ON ERROR)`,
        `SELECT JSON_QUERY(jsonb '[]', '$[*]' EMPTY ON EMPTY)`,
        `SELECT JSON_VALUE(jsonb '1', '$' EMPTY ON ERROR)`,
    }
    for _, sql := range sqls {
        _, err := parser.Parse(sql)
        if err != nil {
            t.Errorf("Parse(%q) failed: %v", sql[:60], err)
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/parsertest/ -run TestJsonBareEmpty -v`
Expected: FAIL with "syntax error"

**Step 3: Add the missing production**

In `parser/gram.y`, at line 10456 (after `EMPTY_P OBJECT_P` line), add:

```yacc
    | EMPTY_P OBJECT_P { $$ = int64(nodes.JSON_BEHAVIOR_EMPTY_OBJECT) }
    /* non-standard, for Oracle compatibility */
    | EMPTY_P          { $$ = int64(nodes.JSON_BEHAVIOR_EMPTY_ARRAY) }
    ;
```

This matches PG17 `gram.y:16882`.

**Step 4: Regenerate parser and run tests**

Run: `make generate-parser && go test ./parser/parsertest/ -run TestJsonBareEmpty -v`
Expected: PASS

**Step 5: Run full regression suite to verify no regressions**

Run: `go test ./parser/pgregress/ -run TestPGRegress -count=1`
Expected: 4 fewer known failures (indices 3 and 114 in sqljson_jsontable.sql, 137 in sqljson_queryfuncs.sql, 290 in sqljson_queryfuncs.sql)

**Step 6: Commit**

```bash
git add parser/gram.y parser/parser.go parser/parsertest/sqljson_test.go
git commit -m "gram.y: add bare EMPTY to json_behavior_type for Oracle compat"
```

---

### Task 2: Add `oper_argtypes` single-type error production

**Files:**
- Modify: `parser/gram.y:13370-13385`
- Test: regression tests for `errors.sql`

**Fixes:** ~2 failures (better error messages for `DROP OPERATOR +(integer)`)

**Step 1: Write a failing test**

```go
// In parser/parsertest/operator_test.go (or appropriate file)
func TestOperArgTypesSingleError(t *testing.T) {
    // This should parse (the grammar accepts it) but the action produces an error.
    // In PG17, this production exists to give a helpful error message.
    // For pgparser, we just need it to not be a generic syntax error.
    sql := `DROP OPERATOR +(integer)`
    _, err := parser.Parse(sql)
    if err == nil {
        t.Error("Expected error for single-arg operator, got nil")
    }
    // The key test: it should NOT be "syntax error" - it should parse the production
    // and produce a semantic error about missing argument.
    if strings.Contains(err.Error(), "syntax error") {
        t.Errorf("Got generic syntax error instead of specific error: %v", err)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/parsertest/ -run TestOperArgTypesSingleError -v`
Expected: FAIL (currently gives "syntax error")

**Step 3: Add the missing production**

In `parser/gram.y`, at line 13370, add the single-type error production as the FIRST alternative:

```yacc
oper_argtypes:
    '(' Typename ')'
        {
            pglex := pglex.(*lexer)
            pglex.yyerror("missing argument\nHINT: Use NONE to denote the missing argument of a unary operator.")
            return 1
        }
    | '(' Typename ',' Typename ')'
        {
            $$ = &nodes.List{Items: []nodes.Node{$2, $4}}
        }
    | '(' NONE ',' Typename ')'
        {
            /* left unary */
            $$ = &nodes.List{Items: []nodes.Node{nil, $4}}
        }
    | '(' Typename ',' NONE ')'
        {
            /* right unary */
            $$ = &nodes.List{Items: []nodes.Node{$2, nil}}
        }
    ;
```

Note: the exact error reporting mechanism depends on how pgparser's lexer exposes `yyerror`. Check `parser/gram.y` for how other error productions work (search for `yyerror` or `return 1`).

**Step 4: Regenerate parser and run tests**

Run: `make generate-parser && go test ./parser/parsertest/ -run TestOperArgTypesSingleError -v`
Expected: PASS (or adjust error mechanism based on codebase patterns)

**Step 5: Commit**

```bash
git add parser/gram.y parser/parser.go
git commit -m "gram.y: add oper_argtypes single-type error production"
```

---

### Task 3: Add `func_expr_windowless` to `stats_param`

**Files:**
- Modify: `parser/gram.y:15755-15767`
- Test: `parser/pgregress/testdata/sql/stats_ext.sql`

**Fixes:** ~6 failures

**Step 1: Write a failing test**

```go
func TestCreateStatisticsFuncExpr(t *testing.T) {
    sqls := []string{
        `CREATE STATISTICS s1 ON my_func(a, b) FROM t1`,
        `CREATE STATISTICS s2 (dependencies) ON lower(col1), col2 FROM t1`,
    }
    for _, sql := range sqls {
        _, err := parser.Parse(sql)
        if err != nil {
            t.Errorf("Parse(%q) failed: %v", sql[:50], err)
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/parsertest/ -run TestCreateStatisticsFuncExpr -v`
Expected: FAIL

**Step 3: Add the missing production**

In `parser/gram.y`, after the `ColId` alternative and before the `'(' a_expr ')'` alternative (around line 15762), add:

```yacc
stats_param:
    ColId
        {
            $$ = &nodes.StatsElem{
                Name: $1,
            }
        }
    | func_expr_windowless
        {
            $$ = &nodes.StatsElem{
                Expr: $1,
            }
        }
    | '(' a_expr ')'
        {
            $$ = &nodes.StatsElem{
                Expr: $2,
            }
        }
    ;
```

This matches PG17 `gram.y:4633-4651`.

**Step 4: Regenerate parser and run tests**

Run: `make generate-parser && go test ./parser/parsertest/ -run TestCreateStatisticsFuncExpr -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/gram.y parser/parser.go
git commit -m "gram.y: add func_expr_windowless to stats_param for CREATE STATISTICS"
```

---

### Task 4: Add `JTC_FORMATTED` production to `json_table_column_definition`

**Files:**
- Modify: `parser/gram.y:10660-10712` (add new production)
- Modify: `parser/gram.y` (add mandatory `json_format_clause` rule)
- Test: `parser/pgregress/testdata/sql/sqljson_jsontable.sql` (indices 18, 27, 28, 75, 76, 77, 79, 81, 82, 83, 93)

**Fixes:** 11 failures

**Step 1: Write a failing test**

```go
func TestJsonTableFormattedColumn(t *testing.T) {
    sqls := []string{
        `SELECT * FROM JSON_TABLE(jsonb 'null', 'lax $[*]' COLUMNS (jst text FORMAT JSON PATH '$'))`,
        `SELECT * FROM JSON_TABLE(jsonb '{"a":"1"}', '$' COLUMNS (a text FORMAT JSON PATH '$.a'))`,
    }
    for _, sql := range sqls {
        _, err := parser.Parse(sql)
        if err != nil {
            t.Errorf("Parse(%q) failed: %v", sql[:60], err)
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/parsertest/ -run TestJsonTableFormattedColumn -v`
Expected: FAIL

**Step 3a: Add mandatory `json_format_clause` rule**

Currently pgparser only has `json_format_clause_opt` (line 10410). Add a new mandatory rule. Near line 10410, add:

```yacc
json_format_clause:
    FORMAT JSON
        {
            $$ = &nodes.JsonFormat{
                FormatType: nodes.JS_FORMAT_JSON,
                Location:   -1,
            }
        }
    | FORMAT JSON ENCODING name
        {
            $$ = &nodes.JsonFormat{
                FormatType: nodes.JS_FORMAT_JSON,
                Location:   -1,
            }
        }
    ;

json_format_clause_opt:
    json_format_clause  { $$ = $1 }
    | /* EMPTY */       { $$ = nil }
    ;
```

Also add `json_format_clause` to the `%type <node>` declarations near line 386:
```
%type <node>  json_format_clause json_format_clause_opt
```

**Step 3b: Add `JTC_FORMATTED` production**

In `json_table_column_definition` (line 10660), add a new alternative after the `JTC_REGULAR` production (after line 10685):

```yacc
    | ColId Typename json_format_clause json_table_column_path_clause_opt
        json_wrapper_behavior json_quotes_clause_opt
        json_behavior_clause_opt
        {
            onEmpty, onError := splitJsonBehaviorClause($7)
            $$ = &nodes.JsonTableColumn{
                Coltype:  nodes.JTC_FORMATTED,
                Name:     $1,
                TypeName: $2,
                Format:   $3.(*nodes.JsonFormat),
                Pathspec: asJsonTablePathSpec($4),
                Wrapper:  nodes.JsonWrapper($5),
                Quotes:   nodes.JsonQuotes($6),
                OnEmpty:  onEmpty,
                OnError:  onError,
                Location: -1,
            }
        }
```

**Step 4: Regenerate parser and run tests**

Run: `make generate-parser && go test ./parser/parsertest/ -run TestJsonTableFormattedColumn -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/gram.y parser/parser.go
git commit -m "gram.y: add JTC_FORMATTED production for JSON_TABLE FORMAT JSON columns"
```

---

### Task 5: Add ordered-set aggregate `aggr_args` productions

**Files:**
- Modify: `parser/gram.y:7414-7424`
- Test: `parser/pgregress/testdata/sql/create_aggregate.sql`

**Fixes:** ~5 failures

**Step 1: Write a failing test**

```go
func TestOrderedSetAggrArgs(t *testing.T) {
    sqls := []string{
        // Ordered-set agg with no direct args (ORDER BY only)
        `CREATE AGGREGATE my_percentile(ORDER BY float8) (sfunc = ordered_set_transition, stype = internal, finalfunc = percentile_disc_final, finalfunc_extra)`,
        // Ordered-set agg with both direct and ordered args
        `CREATE AGGREGATE my_percentile2(float8 ORDER BY float8) (sfunc = ordered_set_transition, stype = internal, finalfunc = percentile_disc_final, finalfunc_extra)`,
        // DROP with ordered-set args
        `DROP AGGREGATE my_percentile(ORDER BY float8)`,
    }
    for _, sql := range sqls {
        _, err := parser.Parse(sql)
        if err != nil {
            t.Errorf("Parse(%q) failed: %v", sql[:60], err)
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/parsertest/ -run TestOrderedSetAggrArgs -v`
Expected: FAIL

**Step 3: Modify `aggr_args` to add ORDER BY variants**

In `parser/gram.y`, replace the `aggr_args` rule (lines 7414-7424) with:

```yacc
aggr_args:
    '(' '*' ')'
        {
            /* agg(*) - represented as two-element list: [nil, Integer(-1)] */
            $$ = &nodes.List{Items: []nodes.Node{nil, &nodes.Integer{Ival: -1}}}
        }
    | '(' aggr_args_list ')'
        {
            /* normal agg(args) - represented as [args_list, Integer(-1)] */
            $$ = &nodes.List{Items: []nodes.Node{$2, &nodes.Integer{Ival: -1}}}
        }
    | '(' ORDER BY aggr_args_list ')'
        {
            /* ordered-set agg with no direct args: agg(ORDER BY args) */
            $$ = &nodes.List{Items: []nodes.Node{$4, &nodes.Integer{Ival: 0}}}
        }
    | '(' aggr_args_list ORDER BY aggr_args_list ')'
        {
            /* ordered-set agg: agg(direct_args ORDER BY ordered_args) */
            $$ = makeOrderedSetArgs($2, $5)
        }
    ;
```

**Important note:** The current `aggr_args` returns a flat list. PG17 returns a 2-element list `[args, ndirectargs_marker]`. This change modifies the return type format. The downstream consumer `extractAggrArgTypes` (line 17328) must be updated to handle the new 2-element list format:

```go
func extractAggrArgTypes(args *nodes.List) *nodes.List {
    if args == nil || len(args.Items) == 0 {
        return nil
    }
    // New format: [args_list_or_nil, Integer_marker]
    // First element is either nil (for agg(*)) or a List of FunctionParameter
    argsList, _ := args.Items[0].(*nodes.List)
    if argsList == nil {
        return nil
    }
    result := &nodes.List{}
    for _, item := range argsList.Items {
        if fp, ok := item.(*nodes.FunctionParameter); ok {
            result.Items = append(result.Items, fp.ArgType)
        }
    }
    return result
}
```

Also add `makeOrderedSetArgs` helper:

```go
func makeOrderedSetArgs(directargs, orderedargs *nodes.List) *nodes.List {
    // Combine direct and ordered args into one list
    combined := &nodes.List{}
    if directargs != nil {
        combined.Items = append(combined.Items, directargs.Items...)
    }
    if orderedargs != nil {
        combined.Items = append(combined.Items, orderedargs.Items...)
    }
    ndirect := 0
    if directargs != nil {
        ndirect = len(directargs.Items)
    }
    return &nodes.List{Items: []nodes.Node{combined, &nodes.Integer{Ival: int64(ndirect)}}}
}
```

**Step 4: Regenerate parser and run all tests**

Run: `make generate-parser && go test ./parser/... -count=1`
Expected: PASS (may need to fix other consumers of `aggr_args` output)

**Step 5: Commit**

```bash
git add parser/gram.y parser/parser.go
git commit -m "gram.y: add ordered-set aggregate aggr_args productions"
```

---

### Task 6: Add `%TYPE` productions to `func_type`

**Files:**
- Modify: `parser/gram.y:4562-4564`
- Test: `parser/pgregress/testdata/sql/misc.sql` (index 29)

**Fixes:** ~3 failures

**Step 1: Write a failing test**

```go
func TestFuncTypePctType(t *testing.T) {
    sqls := []string{
        `CREATE FUNCTION foo(x hobbies_r.name%TYPE) RETURNS hobbies_r.person%TYPE AS 'select 1' LANGUAGE SQL`,
        `CREATE FUNCTION bar(SETOF mytable.col%TYPE) RETURNS void AS 'select 1' LANGUAGE SQL`,
    }
    for _, sql := range sqls {
        _, err := parser.Parse(sql)
        if err != nil {
            t.Errorf("Parse(%q) failed: %v", sql[:60], err)
        }
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/parsertest/ -run TestFuncTypePctType -v`
Expected: FAIL

**Step 3: Add `%TYPE` productions**

In `parser/gram.y`, replace `func_type` (lines 4562-4564) with:

```yacc
func_type:
    Typename { $$ = $1 }
    | type_function_name attrs '%' TYPE_P
        {
            names := &nodes.List{Items: []nodes.Node{&nodes.String{Str: $1}}}
            names.Items = append(names.Items, $2.(*nodes.List).Items...)
            tn := makeTypeNameFromNameList(names).(*nodes.TypeName)
            tn.PctType = true
            tn.Location = -1
            $$ = tn
        }
    | SETOF type_function_name attrs '%' TYPE_P
        {
            names := &nodes.List{Items: []nodes.Node{&nodes.String{Str: $2}}}
            names.Items = append(names.Items, $3.(*nodes.List).Items...)
            tn := makeTypeNameFromNameList(names).(*nodes.TypeName)
            tn.PctType = true
            tn.Setof = true
            tn.Location = -1
            $$ = tn
        }
    ;
```

Note: `attrs` (line 10857) returns a `*nodes.List` of `*nodes.String`. The `type_function_name` (line 4537) returns a `string`. We prepend it using `lcons` equivalent.

**Step 4: Regenerate parser and run tests**

Run: `make generate-parser && go test ./parser/parsertest/ -run TestFuncTypePctType -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/gram.y parser/parser.go
git commit -m "gram.y: add %TYPE productions to func_type for function param/return types"
```

---

### Task 7: Fix extractor `psqlTerminatorRE` to include `\gdesc`, `\gexec`, `\crosstabview`

**Files:**
- Modify: `parser/pgregress/extract.go:35`

**Fixes:** ~9 failures

**Step 1: Verify current regex**

Current regex at line 35: `(?i)\\(g|gx|gset)\b`

**Step 2: Fix the regex**

Change to:
```go
var psqlTerminatorRE = regexp.MustCompile(`(?i)\\(g|gx|gset|gdesc|gexec|crosstabview)\b`)
```

**Step 3: Run regression tests**

Run: `go test ./parser/pgregress/ -run TestPGRegress -count=1`
Expected: ~9 fewer failures in psql.sql

**Step 4: Update known_failures.json**

Run: `go test ./parser/pgregress/ -run TestPGRegress -update`

**Step 5: Commit**

```bash
git add parser/pgregress/extract.go parser/pgregress/known_failures.json
git commit -m "extract: add gdesc, gexec, crosstabview to psql terminator regex"
```

---

### Task 8: Fix extractor to handle semicolons inside `CREATE RULE ... DO (...)`

**Files:**
- Modify: `parser/pgregress/extract.go` (in `splitStatements` function, line 306)

**Fixes:** ~12 failures

**Step 1: Add parenthesis depth tracking**

In `splitStatements`, add a `parenDepth` counter to the `stNormal` state. Only treat `;` as a statement terminator when `parenDepth == 0`.

In the `stNormal` case of the switch (around line 304):

```go
case stNormal:
    switch {
    case ch == ';' && parenDepth == 0:
        // Statement terminator (only at top level)
        emit()
    case ch == ';' && parenDepth > 0:
        // Semicolon inside parens (e.g., CREATE RULE body) - keep it
        hasSQL = true
        buf.WriteByte(ch)
    case ch == '(':
        parenDepth++
        hasSQL = true
        buf.WriteByte(ch)
    case ch == ')':
        parenDepth--
        hasSQL = true
        buf.WriteByte(ch)
    // ... rest of existing cases
    }
```

Declare `parenDepth := 0` alongside `blockDepth` at the top of `splitStatements`.

**Step 2: Run regression tests**

Run: `go test ./parser/pgregress/ -run TestPGRegress -count=1`
Expected: ~12 fewer failures in rules.sql and with.sql

**Step 3: Update known_failures.json**

Run: `go test ./parser/pgregress/ -run TestPGRegress -update`

**Step 4: Commit**

```bash
git add parser/pgregress/extract.go parser/pgregress/known_failures.json
git commit -m "extract: track parenthesis depth to handle CREATE RULE multi-statement bodies"
```

---

### Task 9: Fix extractor to handle `BEGIN ATOMIC ... END` function bodies

**Files:**
- Modify: `parser/pgregress/extract.go` (in `splitStatements` function)

**Fixes:** ~15 failures

**Step 1: Add BEGIN ATOMIC tracking**

In `splitStatements`, add state tracking for `BEGIN ATOMIC` blocks. When we see `BEGIN` followed by `ATOMIC` (case-insensitive) at `parenDepth == 0`, increment `atomicDepth`. When we see `END` at the matching depth, decrement. Only treat `;` as a terminator when both `parenDepth == 0` AND `atomicDepth == 0`.

This requires lookahead or buffering. A simpler approach: track if we're inside a `BEGIN ATOMIC` block by scanning for the keyword pair in the accumulated buffer when we see a `;`.

Alternative approach: Add a `beginAtomicDepth` counter. When the scanner encounters `BEGIN` followed by `ATOMIC` (with optional whitespace), set `beginAtomicDepth++`. When it encounters `END` at `beginAtomicDepth > 0`, decrement. Skip semicolons when `beginAtomicDepth > 0`.

Implementation detail: detecting `BEGIN ATOMIC` requires looking at the token stream, not individual characters. A practical approach is to check the buffer contents after writing each word to see if it ends with `BEGIN ATOMIC`.

**Step 2: Run regression tests**

Run: `go test ./parser/pgregress/ -run TestPGRegress -count=1`

**Step 3: Update known_failures.json**

Run: `go test ./parser/pgregress/ -run TestPGRegress -update`

**Step 4: Commit**

```bash
git add parser/pgregress/extract.go parser/pgregress/known_failures.json
git commit -m "extract: handle BEGIN ATOMIC function bodies without splitting on semicolons"
```

---

## Summary of Expected Impact

| Task | Description | Failures Fixed |
|------|------------|---------------|
| 1 | Bare EMPTY in json_behavior_type | 4 |
| 2 | oper_argtypes single-type error | ~2 |
| 3 | stats_param func_expr_windowless | ~6 |
| 4 | JTC_FORMATTED JSON_TABLE column | 11 |
| 5 | Ordered-set aggregate aggr_args | ~5 |
| 6 | func_type %TYPE | ~3 |
| 7 | Extractor: psqlTerminatorRE | ~9 |
| 8 | Extractor: paren depth tracking | ~12 |
| 9 | Extractor: BEGIN ATOMIC | ~15 |
| **Total** | | **~67** |

After all tasks, remaining failures should be:
- psql variable interpolation (~237) - requires preprocessor, out of scope
- Intentionally invalid SQL (~107) - correct behavior
- UESCAPE support (~10) - deferred (requires lexer changes)
- Misc extractor edge cases (~24) - diminishing returns
