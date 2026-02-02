# gram.y Completion Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Complete all 165 missing grammar rules in our gram.y to achieve full PostgreSQL 17.7 SQL syntax coverage, raising regression test pass rate from 68.4% toward 95%+.

**Architecture:** Translate missing rules from PG's Bison gram.y into Go/goyacc syntax. Each task adds a cohesive group of rules, writes tests, builds, and commits. We work in priority order: highest-impact rules first (func_expr_common_subexpr, a_expr, b_expr, window functions, table_ref) then medium-impact (GROUP BY extensions, type system, FOR UPDATE, partitioning, ALTER TABLE) then lower-impact (JSON, XML, CTE extensions).

**Tech Stack:** Go 1.21+, goyacc, PostgreSQL 17.7 gram.y reference

**Current State:**
- Our gram.y: 10,897 lines, 463 rules, 313 productions
- PG gram.y: 19,520 lines, 478+ rules
- Regression: 30,984/45,314 pass (68.4%) across 224 files
- Key gap: 165 missing rule alternatives across ~40 rule names

**Reference files:**
- Our grammar: `parser/gram.y`
- PG grammar: `~/Github/postgres/src/backend/parser/gram.y`
- Node types: `nodes/` directory
- Regression tests: `parser/pgregress/`

**Build & test commands:**
- Build: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/e421-pg-parser-pg/pgparser && go build ./...`
- Test: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/e421-pg-parser-pg/pgparser && go test ./...`
- Regress stats: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/e421-pg-parser-pg/pgparser/parser/pgregress && go test -run TestPGRegressStats -v -count=1`
- Update baseline: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/e421-pg-parser-pg/pgparser/parser/pgregress && go test -run TestPGRegress -update -v -count=1`

---

## Translation Reference

### Go/goyacc equivalents for PG Bison patterns

| PG Bison (C) | Our goyacc (Go) |
|---------------|------------------|
| `makeNode(X)` | `&nodes.X{}` |
| `list_make1(x)` | `makeList($1)` or `&nodes.List{Items: []nodes.Node{x}}` |
| `list_make2(x,y)` | `makeList2(x, y)` or `&nodes.List{Items: []nodes.Node{x, y}}` |
| `lappend(list, x)` | `appendList(list, x)` |
| `NIL` | `nil` |
| `$1` (Node) | `$1` with type assertion if needed |
| `makeFuncCall(name, args, ...)` | `&nodes.FuncCall{Funcname: makeFuncName("name"), Args: args}` |
| `SystemFuncName("x")` | `makeFuncName("pg_catalog", "x")` — check if makeFuncName supports 2-arg form, else build list manually |
| `SystemTypeName("x")` | `makeTypeName("x")` |
| `makeTypeCast(arg, typename, loc)` | `&nodes.TypeCast{Arg: arg, TypeName: typename}` |
| `makeStringConst(s, loc)` | `makeStringConst(s)` |
| `makeIntConst(i, loc)` | `makeIntConst(int64(i))` |
| `makeSQLValueFunction(op, typmod, loc)` | `&nodes.SQLValueFunction{Op: nodes.SVFOp_XXX, Typmod: typmod}` |
| `makeSimpleA_Expr(kind, op, l, r, loc)` | `makeAExpr(nodes.AExprKind_XXX, op, l, r)` |
| `makeAndExpr(l, r, loc)` | `makeBoolExpr(nodes.AND_EXPR, l, r)` |
| `makeOrExpr(l, r, loc)` | `makeBoolExpr(nodes.OR_EXPR, l, r)` |
| `makeNotExpr(e, loc)` | `makeBoolExpr(nodes.NOT_EXPR, e, nil)` |
| `makeXmlExpr(op, name, named, args, loc)` | `&nodes.XmlExpr{Op: nodes.XmlExprOp_XXX, Name: name, NamedArgs: named, Args: args}` |

### Helper functions we may need to ADD to gram.y Go section

Many PG helper functions don't exist yet in our Go code section. For each task, check if the needed helpers exist and add them if not. Key ones to add:

1. `makeSQLValueFunction(op int, typmod int) nodes.Node` — creates `&nodes.SQLValueFunction{Op: op, Typmod: typmod}`
2. `makeTypeCast(arg nodes.Node, typeName *nodes.TypeName) nodes.Node` — creates `&nodes.TypeCast{Arg: arg, TypeName: typeName}`
3. `makeFuncCallStr(name string, args *nodes.List) nodes.Node` — creates FuncCall with system func name
4. `makeCoalesceExpr(args *nodes.List) nodes.Node`
5. `makeMinMaxExpr(op int, args *nodes.List) nodes.Node`
6. `makeNullTest(arg nodes.Node, testtype int) nodes.Node`
7. `makeBooleanTest(arg nodes.Node, testtype int) nodes.Node`
8. `makeCollateClause(arg nodes.Node, collname *nodes.List) nodes.Node`

### Node types we may need to ADD to nodes/

Before each task, check if the required node types exist in `nodes/`. If not, add them. Key ones likely needed:
- `SQLValueFunction` (Op, Typmod)
- `CoalesceExpr` (Args)
- `MinMaxExpr` (Op, Args) — IS_GREATEST, IS_LEAST
- `WindowDef` (Name, Refname, PartitionClause, OrderClause, FrameOptions, StartOffset, EndOffset)
- `LockingClause` (LockedRels, Strength, WaitPolicy)
- `RangeFunction` (Lateral, Ordinality, IsRowsfrom, Functions, Alias, Coldeflist)
- `RangeTableSample` (Relation, Method, Args, Repeatable)
- `RangeTableFunc` (Lateral, Docexpr, Rowexpr, Namespaces, Columns, Alias)
- `RangeTableFuncCol` (Colname, TypeName, ForOrdinality, IsNotNull, Colexpr, Coldefexpr)
- `XmlSerialize` (Xmloption, Expr, TypeName, Indent)
- `JsonObjectConstructor`, `JsonArrayConstructor`, `JsonArrayQueryConstructor`
- `JsonParseExpr`, `JsonScalarExpr`, `JsonSerializeExpr`, `JsonFuncExpr`
- `JsonTable`, `JsonTableColumn`, `JsonTablePathSpec`
- `JsonOutput`, `JsonReturning`, `JsonFormat`, `JsonValueExpr`, `JsonKeyValue`
- `JsonAggConstructor`, `JsonObjectAgg`, `JsonArrayAgg`
- `JsonArgument`, `JsonBehavior`
- `PartitionSpec`, `PartitionElem`, `PartitionBoundSpec`, `PartitionCmd`
- `CTESearchClause`, `CTECycleClause`
- `SetToDefault`

### %union additions likely needed

```
windowdef *nodes.WindowDef
lockclause *nodes.LockingClause
partspec *nodes.PartitionSpec
partelem *nodes.PartitionElem
partbound *nodes.PartitionBoundSpec
```

---

## PRIORITY 1: Highest Impact (covers ~70% of failures)

### Task 1: Add func_expr_common_subexpr — SQL Value Functions

**Impact:** Enables CURRENT_DATE, CURRENT_TIME, CURRENT_TIMESTAMP, LOCALTIME, LOCALTIMESTAMP, CURRENT_USER, SESSION_USER, CURRENT_ROLE, CURRENT_CATALOG, CURRENT_SCHEMA, USER, SYSTEM_USER — used in thousands of test statements.

**Files:**
- Modify: `parser/gram.y` — add `func_expr_common_subexpr` rule, update `func_expr` to reference it
- Possibly modify: `nodes/` — ensure `SQLValueFunction` node type exists

**Step 1: Check node types exist**

Check if `nodes.SQLValueFunction` exists. If not, add it with fields: `Op int`, `Typmod int32`. Also check for the SVFOP_* constants.

**Step 2: Add helper functions to gram.y Go section**

Add to the Go code section at the bottom of gram.y:
```go
func makeSQLValueFunction(op int, typmod int) nodes.Node {
    return &nodes.SQLValueFunction{Op: op, Typmod: int32(typmod)}
}

func makeTypeCast(arg nodes.Node, typeName *nodes.TypeName) nodes.Node {
    return &nodes.TypeCast{Arg: arg, TypeName: typeName}
}
```

**Step 3: Add the func_expr_common_subexpr rule — Part 1: SQL value functions**

Add new rule after `func_expr` in gram.y:
```yacc
func_expr_common_subexpr:
    COLLATION FOR '(' a_expr ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "pg_collation_for"),
                Args:     makeList($4),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | CURRENT_DATE
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_DATE, -1) }
    | CURRENT_TIME
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIME, -1) }
    | CURRENT_TIME '(' Iconst ')'
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIME_N, int($3)) }
    | CURRENT_TIMESTAMP
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIMESTAMP, -1) }
    | CURRENT_TIMESTAMP '(' Iconst ')'
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIMESTAMP_N, int($3)) }
    | LOCALTIME
        { $$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIME, -1) }
    | LOCALTIME '(' Iconst ')'
        { $$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIME_N, int($3)) }
    | LOCALTIMESTAMP
        { $$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIMESTAMP, -1) }
    | LOCALTIMESTAMP '(' Iconst ')'
        { $$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIMESTAMP_N, int($3)) }
    | CURRENT_ROLE
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_ROLE, -1) }
    | CURRENT_USER
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_USER, -1) }
    | SESSION_USER
        { $$ = makeSQLValueFunction(nodes.SVFOP_SESSION_USER, -1) }
    | SYSTEM_USER
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "system_user"),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | USER
        { $$ = makeSQLValueFunction(nodes.SVFOP_USER, -1) }
    | CURRENT_CATALOG
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_CATALOG, -1) }
    | CURRENT_SCHEMA
        { $$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_SCHEMA, -1) }
    ;
```

**Step 4: Update func_expr to include func_expr_common_subexpr**

Change:
```yacc
func_expr:
    func_application { $$ = $1 }
    ;
```
To:
```yacc
func_expr:
    func_application { $$ = $1 }
    | func_expr_common_subexpr { $$ = $1 }
    ;
```

**Step 5: Build and test**

Run: `go build ./...` — fix any compilation errors
Run: `go test ./...` — verify existing tests still pass

**Step 6: Commit**

```bash
git add parser/gram.y nodes/
git commit -m "feat(gram): add SQL value functions (CURRENT_DATE/TIME/USER etc)"
```

---

### Task 2: Add func_expr_common_subexpr — CAST, EXTRACT, NORMALIZE, OVERLAY, POSITION, SUBSTRING, TREAT, TRIM

**Impact:** CAST alone appears in ~4000+ failing statements. EXTRACT, SUBSTRING, TRIM are also very common.

**Files:**
- Modify: `parser/gram.y` — extend `func_expr_common_subexpr`, add helper rules

**Step 1: Add helper rules**

Add these helper rules to gram.y:

```yacc
extract_list:
    extract_arg FROM a_expr
        {
            $$ = makeList2(makeStringConst($1), $3)
        }
    ;

extract_arg:
    IDENT       { $$ = $1 }
    | YEAR_P    { $$ = "year" }
    | MONTH_P   { $$ = "month" }
    | DAY_P     { $$ = "day" }
    | HOUR_P    { $$ = "hour" }
    | MINUTE_P  { $$ = "minute" }
    | SECOND_P  { $$ = "second" }
    | Sconst    { $$ = $1 }
    ;

unicode_normal_form:
    NFC  { $$ = "NFC" }
    | NFD  { $$ = "NFD" }
    | NFKC { $$ = "NFKC" }
    | NFKD { $$ = "NFKD" }
    ;

overlay_list:
    a_expr PLACING a_expr FROM a_expr FOR a_expr
        {
            $$ = &nodes.List{Items: []nodes.Node{$1, $3, $5, $7}}
        }
    | a_expr PLACING a_expr FROM a_expr
        {
            $$ = &nodes.List{Items: []nodes.Node{$1, $3, $5}}
        }
    ;

position_list:
    b_expr IN_P b_expr
        {
            $$ = makeList2($3, $1)  /* note: reversed */
        }
    ;

substr_list:
    a_expr FROM a_expr FOR a_expr
        {
            $$ = &nodes.List{Items: []nodes.Node{$1, $3, $5}}
        }
    | a_expr FOR a_expr FROM a_expr
        {
            $$ = &nodes.List{Items: []nodes.Node{$1, $5, $3}}
        }
    | a_expr FROM a_expr
        {
            $$ = makeList2($1, $3)
        }
    | a_expr FOR a_expr
        {
            $$ = &nodes.List{Items: []nodes.Node{
                $1,
                makeIntConst(1),
                &nodes.TypeCast{Arg: $3, TypeName: makeTypeName("int4")},
            }}
        }
    | a_expr SIMILAR a_expr ESCAPE a_expr
        {
            $$ = &nodes.List{Items: []nodes.Node{$1, $3, $5}}
        }
    ;

trim_list:
    a_expr FROM expr_list
        {
            $$ = appendList($3, $1)
        }
    | FROM expr_list
        {
            $$ = $2
        }
    | expr_list
        {
            $$ = $1
        }
    ;
```

**Step 2: Extend func_expr_common_subexpr with CAST, EXTRACT, etc.**

Add these alternatives to `func_expr_common_subexpr`:
```yacc
    | CAST '(' a_expr AS Typename ')'
        { $$ = makeTypeCast($3, $5) }
    | EXTRACT '(' extract_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "extract"),
                Args:     $3,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | NORMALIZE '(' a_expr ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "normalize"),
                Args:     makeList($3),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | NORMALIZE '(' a_expr ',' unicode_normal_form ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "normalize"),
                Args:     makeList2($3, makeStringConst($5)),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | OVERLAY '(' overlay_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "overlay"),
                Args:     $3,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | OVERLAY '(' func_arg_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "overlay"}}},
                Args:     $3,
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
        }
    | POSITION '(' position_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "position"),
                Args:     $3,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | SUBSTRING '(' substr_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "substring"),
                Args:     $3,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | SUBSTRING '(' func_arg_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "substring"}}},
                Args:     $3,
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
        }
    | TREAT '(' a_expr AS Typename ')'
        {
            /* TREAT converts to a function call using the type name */
            last := $5.Names.Items[len($5.Names.Items)-1]
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", last.(*nodes.String).Str),
                Args:     makeList($3),
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
        }
    | TRIM '(' BOTH trim_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "btrim"),
                Args:     $4,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | TRIM '(' LEADING trim_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "ltrim"),
                Args:     $4,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | TRIM '(' TRAILING trim_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "rtrim"),
                Args:     $4,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | TRIM '(' trim_list ')'
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "btrim"),
                Args:     $3,
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
```

**Step 3: Build and test**

Run: `go build ./...` — fix compilation errors (especially %type declarations for new rules)
Run: `go test ./...`

**Step 4: Run regression stats**

Run: `go test -run TestPGRegressStats -v -count=1` in `parser/pgregress/`
Note the improvement in pass rate.

**Step 5: Commit**

```bash
git add parser/gram.y
git commit -m "feat(gram): add CAST, EXTRACT, SUBSTRING, TRIM, OVERLAY, POSITION, NORMALIZE"
```

---

### Task 3: Add func_expr_common_subexpr — NULLIF, COALESCE, GREATEST, LEAST

**Impact:** Very commonly used in test SQL. COALESCE alone appears in hundreds of statements.

**Files:**
- Modify: `parser/gram.y`
- Possibly modify: `nodes/` — ensure CoalesceExpr, MinMaxExpr exist

**Step 1: Check/add node types**

Ensure `nodes.CoalesceExpr` (Args, Location) and `nodes.MinMaxExpr` (Op, Args, Location) exist with IS_GREATEST/IS_LEAST constants.

**Step 2: Add alternatives to func_expr_common_subexpr**

```yacc
    | NULLIF '(' a_expr ',' a_expr ')'
        {
            $$ = makeAExpr(nodes.AEXPR_NULLIF, "=", $3, $5)
        }
    | COALESCE '(' expr_list ')'
        {
            $$ = &nodes.CoalesceExpr{Args: $3}
        }
    | GREATEST '(' expr_list ')'
        {
            $$ = &nodes.MinMaxExpr{Op: nodes.IS_GREATEST, Args: $3}
        }
    | LEAST '(' expr_list ')'
        {
            $$ = &nodes.MinMaxExpr{Op: nodes.IS_LEAST, Args: $3}
        }
```

**Step 3: Build, test, commit**

---

### Task 4: Complete a_expr — IS NULL/NOT NULL, IS TRUE/FALSE/UNKNOWN, ISNULL/NOTNULL

**Impact:** IS NULL and IS NOT NULL are used in nearly every test file. IS TRUE/FALSE/UNKNOWN appear in boolean expression tests.

**Files:**
- Modify: `parser/gram.y` — add ~14 alternatives to `a_expr`
- Check: `nodes/` — NullTest, BooleanTest should already exist

**Step 1: Verify NullTest and BooleanTest nodes exist**

Our gram.y already references NullTest and BooleanTest (confirmed in exploration). Check the exact field names and constants (IS_NULL, IS_NOT_NULL, IS_TRUE, IS_NOT_TRUE, IS_FALSE, IS_NOT_FALSE, IS_UNKNOWN, IS_NOT_UNKNOWN).

**Step 2: Add IS NULL / IS NOT NULL / ISNULL / NOTNULL to a_expr**

Add these alternatives to the `a_expr` rule:
```yacc
    | a_expr IS NULL_P                          %prec IS
        {
            $$ = &nodes.NullTest{Arg: $1, Nulltesttype: nodes.IS_NULL}
        }
    | a_expr ISNULL
        {
            $$ = &nodes.NullTest{Arg: $1, Nulltesttype: nodes.IS_NULL}
        }
    | a_expr IS NOT NULL_P                      %prec IS
        {
            $$ = &nodes.NullTest{Arg: $1, Nulltesttype: nodes.IS_NOT_NULL}
        }
    | a_expr NOTNULL
        {
            $$ = &nodes.NullTest{Arg: $1, Nulltesttype: nodes.IS_NOT_NULL}
        }
```

**Step 3: Add IS TRUE / IS NOT TRUE / IS FALSE / IS NOT FALSE / IS UNKNOWN / IS NOT UNKNOWN**

```yacc
    | a_expr IS TRUE_P                          %prec IS
        {
            $$ = &nodes.BooleanTest{Arg: $1, Booltesttype: nodes.IS_TRUE}
        }
    | a_expr IS NOT TRUE_P                      %prec IS
        {
            $$ = &nodes.BooleanTest{Arg: $1, Booltesttype: nodes.IS_NOT_TRUE}
        }
    | a_expr IS FALSE_P                         %prec IS
        {
            $$ = &nodes.BooleanTest{Arg: $1, Booltesttype: nodes.IS_FALSE}
        }
    | a_expr IS NOT FALSE_P                     %prec IS
        {
            $$ = &nodes.BooleanTest{Arg: $1, Booltesttype: nodes.IS_NOT_FALSE}
        }
    | a_expr IS UNKNOWN                         %prec IS
        {
            $$ = &nodes.BooleanTest{Arg: $1, Booltesttype: nodes.IS_UNKNOWN}
        }
    | a_expr IS NOT UNKNOWN                     %prec IS
        {
            $$ = &nodes.BooleanTest{Arg: $1, Booltesttype: nodes.IS_NOT_UNKNOWN}
        }
```

**Step 4: Add row OVERLAPS row**

```yacc
    | row OVERLAPS row
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "overlaps"),
                Args:     concatLists($1, $3),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
```

**Step 5: Build, test, commit**

```bash
git commit -m "feat(gram): add IS NULL/TRUE/FALSE/UNKNOWN, ISNULL/NOTNULL, OVERLAPS to a_expr"
```

---

### Task 5: Complete a_expr — IS DISTINCT FROM, BETWEEN, IN, subquery operators

**Impact:** BETWEEN appears in hundreds of statements, IN() is one of the most common SQL constructs. IS DISTINCT FROM used in join tests.

**Files:**
- Modify: `parser/gram.y`
- Check: `nodes/` — SubLink, SetToDefault

**Step 1: Add in_expr helper rule**

```yacc
in_expr:
    select_with_parens
        {
            $$ = &nodes.SubLink{Subselect: $1}
        }
    | '(' expr_list ')'
        {
            $$ = $2
        }
    ;
```

**Step 2: Add IS DISTINCT FROM to a_expr**

```yacc
    | a_expr IS DISTINCT FROM a_expr            %prec IS
        {
            $$ = makeAExpr(nodes.AEXPR_DISTINCT, "=", $1, $5)
        }
    | a_expr IS NOT DISTINCT FROM a_expr        %prec IS
        {
            $$ = makeAExpr(nodes.AEXPR_NOT_DISTINCT, "=", $1, $6)
        }
```

**Step 3: Add BETWEEN variants**

```yacc
    | a_expr BETWEEN opt_asymmetric b_expr AND a_expr       %prec BETWEEN
        {
            $$ = makeAExpr(nodes.AEXPR_BETWEEN, "BETWEEN", $1,
                &nodes.List{Items: []nodes.Node{$4, $6}})
        }
    | a_expr NOT_LA BETWEEN opt_asymmetric b_expr AND a_expr %prec NOT_LA
        {
            $$ = makeAExpr(nodes.AEXPR_NOT_BETWEEN, "NOT BETWEEN", $1,
                &nodes.List{Items: []nodes.Node{$5, $7}})
        }
    | a_expr BETWEEN SYMMETRIC b_expr AND a_expr            %prec BETWEEN
        {
            $$ = makeAExpr(nodes.AEXPR_BETWEEN_SYM, "BETWEEN SYMMETRIC", $1,
                &nodes.List{Items: []nodes.Node{$4, $6}})
        }
    | a_expr NOT_LA BETWEEN SYMMETRIC b_expr AND a_expr     %prec NOT_LA
        {
            $$ = makeAExpr(nodes.AEXPR_NOT_BETWEEN_SYM, "NOT BETWEEN SYMMETRIC", $1,
                &nodes.List{Items: []nodes.Node{$5, $7}})
        }
```

Also add `opt_asymmetric` rule if not present:
```yacc
opt_asymmetric:
    ASYMMETRIC  {}
    | /* EMPTY */ {}
    ;
```

**Step 4: Add IN / NOT IN**

```yacc
    | a_expr IN_P in_expr
        {
            if sub, ok := $3.(*nodes.SubLink); ok {
                sub.SubLinkType = nodes.ANY_SUBLINK
                sub.Testexpr = $1
                $$ = sub
            } else {
                $$ = makeAExpr(nodes.AEXPR_IN, "=", $1, $3)
            }
        }
    | a_expr NOT_LA IN_P in_expr                            %prec NOT_LA
        {
            if sub, ok := $4.(*nodes.SubLink); ok {
                sub.SubLinkType = nodes.ANY_SUBLINK
                sub.Testexpr = $1
                $$ = makeBoolExpr(nodes.NOT_EXPR, sub, nil)
            } else {
                $$ = makeAExpr(nodes.AEXPR_IN, "<>", $1, $4)
            }
        }
```

**Step 5: Add subquery operators (ANY/ALL/SOME)**

```yacc
    | a_expr subquery_Op sub_type select_with_parens        %prec Op
        {
            $$ = &nodes.SubLink{
                SubLinkType: $3,
                Testexpr:    $1,
                OperName:    $2,
                Subselect:   $4,
            }
        }
    | a_expr subquery_Op sub_type '(' a_expr ')'            %prec Op
        {
            if $3 == nodes.ANY_SUBLINK {
                $$ = makeAExpr(nodes.AEXPR_OP_ANY, $2, $1, $5)
            } else {
                $$ = makeAExpr(nodes.AEXPR_OP_ALL, $2, $1, $5)
            }
        }
```

Also verify `sub_type` and `subquery_Op` rules exist.

**Step 6: Add DEFAULT**

```yacc
    | DEFAULT
        {
            $$ = &nodes.SetToDefault{}
        }
```

**Step 7: Build, test, commit**

```bash
git commit -m "feat(gram): add BETWEEN, IN, IS DISTINCT FROM, subquery ops, DEFAULT to a_expr"
```

---

### Task 6: Complete a_expr — LIKE, ILIKE, SIMILAR TO with ESCAPE

**Impact:** LIKE and pattern matching appear in hundreds of test statements.

**Files:**
- Modify: `parser/gram.y`

**Step 1: Add LIKE / NOT LIKE (with and without ESCAPE)**

```yacc
    | a_expr LIKE a_expr
        {
            $$ = makeAExpr(nodes.AEXPR_LIKE, "~~", $1, $3)
        }
    | a_expr LIKE a_expr ESCAPE a_expr                      %prec LIKE
        {
            n := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "like_escape"),
                Args:     makeList2($3, $5),
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
            $$ = makeAExpr(nodes.AEXPR_LIKE, "~~", $1, n)
        }
    | a_expr NOT_LA LIKE a_expr                             %prec NOT_LA
        {
            $$ = makeAExpr(nodes.AEXPR_LIKE, "!~~", $1, $4)
        }
    | a_expr NOT_LA LIKE a_expr ESCAPE a_expr               %prec NOT_LA
        {
            n := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "like_escape"),
                Args:     makeList2($4, $6),
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
            $$ = makeAExpr(nodes.AEXPR_LIKE, "!~~", $1, n)
        }
```

**Step 2: Add ILIKE / NOT ILIKE (with and without ESCAPE)**

Same pattern as LIKE but with AEXPR_ILIKE and operators "~~*" / "!~~*".

**Step 3: Add SIMILAR TO / NOT SIMILAR TO (with and without ESCAPE)**

```yacc
    | a_expr SIMILAR TO a_expr                              %prec SIMILAR
        {
            n := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "similar_to_escape"),
                Args:     makeList($4),
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
            $$ = makeAExpr(nodes.AEXPR_SIMILAR, "~", $1, n)
        }
    | a_expr SIMILAR TO a_expr ESCAPE a_expr                %prec SIMILAR
        {
            n := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "similar_to_escape"),
                Args:     makeList2($4, $6),
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
            $$ = makeAExpr(nodes.AEXPR_SIMILAR, "~", $1, n)
        }
    | a_expr NOT_LA SIMILAR TO a_expr                       %prec NOT_LA
        {
            n := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "similar_to_escape"),
                Args:     makeList($5),
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
            $$ = makeAExpr(nodes.AEXPR_SIMILAR, "!~", $1, n)
        }
    | a_expr NOT_LA SIMILAR TO a_expr ESCAPE a_expr         %prec NOT_LA
        {
            n := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "similar_to_escape"),
                Args:     makeList2($5, $7),
                FuncFormat: nodes.COERCE_EXPLICIT_CALL,
            }
            $$ = makeAExpr(nodes.AEXPR_SIMILAR, "!~", $1, n)
        }
```

**Step 4: Build, test, commit**

```bash
git commit -m "feat(gram): add LIKE, ILIKE, SIMILAR TO with ESCAPE to a_expr"
```

---

### Task 7: Complete a_expr — COLLATE, AT TIME ZONE, IS DOCUMENT, IS NORMALIZED, IS JSON, qual_Op

**Impact:** COLLATE and AT TIME ZONE appear in hundreds of statements. IS JSON is PG17-specific.

**Files:**
- Modify: `parser/gram.y`

**Step 1: Add COLLATE**

```yacc
    | a_expr COLLATE any_name
        {
            $$ = &nodes.CollateClause{Arg: $1, Collname: $3}
        }
```

**Step 2: Add AT TIME ZONE / AT LOCAL**

```yacc
    | a_expr AT TIME ZONE a_expr                            %prec AT
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "timezone"),
                Args:     makeList2($5, $1),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | a_expr AT LOCAL                                       %prec AT
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "timezone"),
                Args:     makeList($1),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
```

**Step 3: Add qual_Op alternatives (if missing)**

Check if these exist already:
```yacc
    | a_expr qual_Op a_expr                                 %prec Op
        { $$ = makeAExpr(nodes.AEXPR_OP, $2, $1, $3) }
    | qual_Op a_expr                                        %prec Op
        { $$ = makeAExpr(nodes.AEXPR_OP, $1, nil, $2) }
```

**Step 4: Add IS DOCUMENT / IS NOT DOCUMENT**

```yacc
    | a_expr IS DOCUMENT_P                                  %prec IS
        {
            $$ = &nodes.XmlExpr{Op: nodes.IS_DOCUMENT, Args: makeList($1)}
        }
    | a_expr IS NOT DOCUMENT_P                              %prec IS
        {
            $$ = makeBoolExpr(nodes.NOT_EXPR,
                &nodes.XmlExpr{Op: nodes.IS_DOCUMENT, Args: makeList($1)}, nil)
        }
```

**Step 5: Add IS NORMALIZED variants**

```yacc
    | a_expr IS NORMALIZED                                  %prec IS
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "is_normalized"),
                Args:     makeList($1),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | a_expr IS unicode_normal_form NORMALIZED              %prec IS
        {
            $$ = &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "is_normalized"),
                Args:     makeList2($1, makeStringConst($3)),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
        }
    | a_expr IS NOT NORMALIZED                              %prec IS
        {
            inner := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "is_normalized"),
                Args:     makeList($1),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
            $$ = makeBoolExpr(nodes.NOT_EXPR, inner, nil)
        }
    | a_expr IS NOT unicode_normal_form NORMALIZED          %prec IS
        {
            inner := &nodes.FuncCall{
                Funcname: makeFuncName("pg_catalog", "is_normalized"),
                Args:     makeList2($1, makeStringConst($4)),
                FuncFormat: nodes.COERCE_SQL_SYNTAX,
            }
            $$ = makeBoolExpr(nodes.NOT_EXPR, inner, nil)
        }
```

**Step 6: Add IS JSON predicates (PG17)**

Add `json_predicate_type_constraint` rule and IS JSON alternatives. (Only if JsonIsPredicate node exists; skip if not.)

**Step 7: Build, test, commit**

```bash
git commit -m "feat(gram): add COLLATE, AT TIME ZONE, IS DOCUMENT/NORMALIZED/JSON to a_expr"
```

---

### Task 8: Complete b_expr — IS DISTINCT FROM, IS DOCUMENT, qual_Op

**Impact:** b_expr is used inside BETWEEN and other restricted contexts. Missing alternatives cause cascade failures.

**Files:**
- Modify: `parser/gram.y`

**Step 1: Add missing b_expr alternatives**

Check which of these are already present and add the missing ones:

```yacc
    | b_expr qual_Op b_expr                                 %prec Op
        { $$ = makeAExpr(nodes.AEXPR_OP, $2, $1, $3) }
    | qual_Op b_expr                                        %prec Op
        { $$ = makeAExpr(nodes.AEXPR_OP, $1, nil, $2) }
    | b_expr IS DISTINCT FROM b_expr                        %prec IS
        { $$ = makeAExpr(nodes.AEXPR_DISTINCT, "=", $1, $5) }
    | b_expr IS NOT DISTINCT FROM b_expr                    %prec IS
        { $$ = makeAExpr(nodes.AEXPR_NOT_DISTINCT, "=", $1, $6) }
    | b_expr IS DOCUMENT_P                                  %prec IS
        {
            $$ = &nodes.XmlExpr{Op: nodes.IS_DOCUMENT, Args: makeList($1)}
        }
    | b_expr IS NOT DOCUMENT_P                              %prec IS
        {
            $$ = makeBoolExpr(nodes.NOT_EXPR,
                &nodes.XmlExpr{Op: nodes.IS_DOCUMENT, Args: makeList($1)}, nil)
        }
```

**Step 2: Build, test, commit**

```bash
git commit -m "feat(gram): complete b_expr with IS DISTINCT FROM, IS DOCUMENT, qual_Op"
```

---

### Task 9: Add Window Function Infrastructure

**Impact:** Window functions (OVER, PARTITION BY, frame clauses) appear in ~800+ failing statements across many test files.

**Files:**
- Modify: `parser/gram.y` — add ~11 new rules
- Modify: `nodes/` — ensure WindowDef exists
- May need: %union additions for windowdef type

**Step 1: Ensure WindowDef node type exists**

Check/add `nodes.WindowDef` with fields: Name, Refname, PartitionClause, OrderClause, FrameOptions, StartOffset, EndOffset, Location. Also add FRAMEOPTION_* constants.

**Step 2: Add %union and %type declarations**

If needed, add to %union and declare %type for new rules.

**Step 3: Add window function rules**

```yacc
window_clause:
    WINDOW window_definition_list   { $$ = $2 }
    | /* EMPTY */                   { $$ = nil }
    ;

window_definition_list:
    window_definition                               { $$ = makeList($1) }
    | window_definition_list ',' window_definition  { $$ = appendList($1, $3) }
    ;

window_definition:
    ColId AS window_specification
        {
            n := $3.(*nodes.WindowDef)
            n.Name = $1
            $$ = n
        }
    ;

window_specification:
    '(' opt_existing_window_name opt_partition_clause opt_sort_clause opt_frame_clause ')'
        {
            n := $5.(*nodes.WindowDef)
            n.Refname = $2
            n.PartitionClause = $3
            n.OrderClause = $4
            $$ = n
        }
    ;

opt_existing_window_name:
    ColId           { $$ = $1 }
    | /* EMPTY */   %prec Op { $$ = "" }
    ;

opt_partition_clause:
    PARTITION BY expr_list  { $$ = $3 }
    | /* EMPTY */           { $$ = nil }
    ;

opt_frame_clause:
    RANGE frame_extent opt_window_exclusion_clause
        {
            n := $2.(*nodes.WindowDef)
            n.FrameOptions |= nodes.FRAMEOPTION_NONDEFAULT | nodes.FRAMEOPTION_RANGE | $3
            $$ = n
        }
    | ROWS frame_extent opt_window_exclusion_clause
        {
            n := $2.(*nodes.WindowDef)
            n.FrameOptions |= nodes.FRAMEOPTION_NONDEFAULT | nodes.FRAMEOPTION_ROWS | $3
            $$ = n
        }
    | GROUPS frame_extent opt_window_exclusion_clause
        {
            n := $2.(*nodes.WindowDef)
            n.FrameOptions |= nodes.FRAMEOPTION_NONDEFAULT | nodes.FRAMEOPTION_GROUPS | $3
            $$ = n
        }
    | /* EMPTY */
        {
            $$ = &nodes.WindowDef{FrameOptions: nodes.FRAMEOPTION_DEFAULTS}
        }
    ;

frame_extent:
    frame_bound
        {
            n := $1.(*nodes.WindowDef)
            n.FrameOptions |= nodes.FRAMEOPTION_END_CURRENT_ROW
            $$ = n
        }
    | BETWEEN frame_bound AND frame_bound
        {
            n1 := $2.(*nodes.WindowDef)
            n2 := $4.(*nodes.WindowDef)
            n1.FrameOptions |= n2.FrameOptions << 1
            n1.FrameOptions |= nodes.FRAMEOPTION_BETWEEN
            n1.EndOffset = n2.StartOffset
            $$ = n1
        }
    ;

frame_bound:
    UNBOUNDED PRECEDING
        { $$ = &nodes.WindowDef{FrameOptions: nodes.FRAMEOPTION_START_UNBOUNDED_PRECEDING} }
    | UNBOUNDED FOLLOWING
        { $$ = &nodes.WindowDef{FrameOptions: nodes.FRAMEOPTION_START_UNBOUNDED_FOLLOWING} }
    | CURRENT_P ROW
        { $$ = &nodes.WindowDef{FrameOptions: nodes.FRAMEOPTION_START_CURRENT_ROW} }
    | a_expr PRECEDING
        { $$ = &nodes.WindowDef{FrameOptions: nodes.FRAMEOPTION_START_OFFSET_PRECEDING, StartOffset: $1} }
    | a_expr FOLLOWING
        { $$ = &nodes.WindowDef{FrameOptions: nodes.FRAMEOPTION_START_OFFSET_FOLLOWING, StartOffset: $1} }
    ;

opt_window_exclusion_clause:
    EXCLUDE CURRENT_P ROW   { $$ = nodes.FRAMEOPTION_EXCLUDE_CURRENT_ROW }
    | EXCLUDE GROUP_P       { $$ = nodes.FRAMEOPTION_EXCLUDE_GROUP }
    | EXCLUDE TIES          { $$ = nodes.FRAMEOPTION_EXCLUDE_TIES }
    | EXCLUDE NO OTHERS     { $$ = 0 }
    | /* EMPTY */           { $$ = 0 }
    ;

over_clause:
    OVER window_specification   { $$ = $2 }
    | OVER ColId
        {
            $$ = &nodes.WindowDef{
                Name: $2,
                FrameOptions: nodes.FRAMEOPTION_DEFAULTS,
            }
        }
    | /* EMPTY */               { $$ = nil }
    ;
```

**Step 4: Update func_expr to use over_clause**

Change `func_expr` from:
```yacc
func_expr:
    func_application { $$ = $1 }
    | func_expr_common_subexpr { $$ = $1 }
    ;
```
To:
```yacc
func_expr:
    func_application within_group_clause filter_clause over_clause
        {
            n := $1.(*nodes.FuncCall)
            if $2 != nil {
                n.AggOrder = $2
                n.AggWithinGroup = true
            }
            n.AggFilter = $3
            n.Over = $4
            $$ = n
        }
    | func_expr_common_subexpr { $$ = $1 }
    ;
```

Also add `within_group_clause` and `filter_clause` if they don't exist:
```yacc
within_group_clause:
    WITHIN GROUP_P '(' sort_clause ')'  { $$ = $4 }
    | /* EMPTY */                       { $$ = nil }
    ;

filter_clause:
    FILTER '(' WHERE a_expr ')'     { $$ = $4 }
    | /* EMPTY */                   { $$ = nil }
    ;
```

**Step 5: Update simple_select to include window_clause**

Add `window_clause` to the SELECT production in `simple_select`.

**Step 6: Build, test, commit**

```bash
git commit -m "feat(gram): add window function infrastructure (OVER, PARTITION BY, frame clauses)"
```

---

### Task 10: Expand table_ref — func_table, LATERAL, tablesample, xmltable

**Impact:** Function calls in FROM clause, LATERAL subqueries, and TABLESAMPLE appear in many test files (~400+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — ensure RangeFunction, RangeTableSample, RangeTableFunc exist

**Step 1: Add func_table and func_expr_windowless rules**

```yacc
func_expr_windowless:
    func_application    { $$ = $1 }
    | func_expr_common_subexpr { $$ = $1 }
    ;

func_table:
    func_expr_windowless opt_ordinality
        {
            $$ = &nodes.RangeFunction{
                Ordinality: $2,
                Functions:  makeList(makeList2($1, nil)),
            }
        }
    | ROWS FROM '(' rowsfrom_list ')' opt_ordinality
        {
            $$ = &nodes.RangeFunction{
                Ordinality: $6,
                IsRowsfrom: true,
                Functions:  $4,
            }
        }
    ;

opt_ordinality:
    WITH_LA ORDINALITY  { $$ = true }
    | /* EMPTY */       { $$ = false }
    ;

func_alias_clause:
    alias_clause
        { $$ = makeList2($1, nil) }
    | AS '(' TableFuncElementList ')'
        { $$ = makeList2(nil, $3) }
    | AS ColId '(' TableFuncElementList ')'
        {
            a := &nodes.Alias{Aliasname: $2}
            $$ = makeList2(a, $4)
        }
    | ColId '(' TableFuncElementList ')'
        {
            a := &nodes.Alias{Aliasname: $1}
            $$ = makeList2(a, $3)
        }
    ;
```

**Step 2: Add new table_ref alternatives**

Add to `table_ref`:
```yacc
    | relation_expr opt_alias_clause tablesample_clause
        {
            n := $3.(*nodes.RangeTableSample)
            $1.(*nodes.RangeVar).Alias = $2
            n.Relation = $1
            $$ = n
        }
    | func_table func_alias_clause
        {
            n := $1.(*nodes.RangeFunction)
            n.Alias = $2[0]  /* first element is alias */
            n.Coldeflist = $2[1]  /* second is column def list */
            $$ = n
        }
    | LATERAL_P func_table func_alias_clause
        {
            n := $2.(*nodes.RangeFunction)
            n.Lateral = true
            n.Alias = $3[0]
            n.Coldeflist = $3[1]
            $$ = n
        }
    | LATERAL_P select_with_parens opt_alias_clause
        {
            $$ = &nodes.RangeSubselect{
                Lateral:  true,
                Subquery: $2,
                Alias:    $3,
            }
        }
    | xmltable opt_alias_clause
        {
            n := $1.(*nodes.RangeTableFunc)
            n.Alias = $2
            $$ = n
        }
    | LATERAL_P xmltable opt_alias_clause
        {
            n := $2.(*nodes.RangeTableFunc)
            n.Lateral = true
            n.Alias = $3
            $$ = n
        }
    | '(' joined_table ')' alias_clause
        {
            $2.(*nodes.JoinExpr).Alias = $4
            $$ = $2
        }
```

**Step 3: Add tablesample_clause rule**

```yacc
tablesample_clause:
    TABLESAMPLE func_name '(' expr_list ')' opt_repeatable_clause
        {
            $$ = &nodes.RangeTableSample{
                Method: $2,
                Args:   $4,
                Repeatable: $6,
            }
        }
    ;

opt_repeatable_clause:
    REPEATABLE '(' a_expr ')'   { $$ = $3 }
    | /* EMPTY */               { $$ = nil }
    ;
```

**Step 4: Build, test, commit**

```bash
git commit -m "feat(gram): expand table_ref with func_table, LATERAL, tablesample, xmltable"
```

---

## PRIORITY 2: Medium Impact

### Task 11: Add FOR UPDATE / FOR SHARE locking clauses

**Impact:** FOR UPDATE, FOR NO KEY UPDATE, FOR SHARE, FOR KEY SHARE appear in ~300+ statements.

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — ensure LockingClause exists

**Step 1: Add locking rules**

```yacc
for_locking_clause:
    for_locking_items               { $$ = $1 }
    | FOR READ ONLY                 { $$ = nil }
    ;

for_locking_items:
    for_locking_item                        { $$ = makeList($1) }
    | for_locking_items for_locking_item    { $$ = appendList($1, $2) }
    ;

for_locking_item:
    for_locking_strength locked_rels_list opt_nowait_or_skip
        {
            $$ = &nodes.LockingClause{
                LockedRels: $2,
                Strength:   $1,
                WaitPolicy: $3,
            }
        }
    ;

for_locking_strength:
    FOR UPDATE              { $$ = nodes.LCS_FORUPDATE }
    | FOR NO KEY UPDATE     { $$ = nodes.LCS_FORNOKEYUPDATE }
    | FOR SHARE             { $$ = nodes.LCS_FORSHARE }
    | FOR KEY SHARE         { $$ = nodes.LCS_FORKEYSHARE }
    ;

locked_rels_list:
    OF qualified_name_list      { $$ = $2 }
    | /* EMPTY */               { $$ = nil }
    ;

opt_nowait_or_skip:
    NOWAIT                      { $$ = nodes.LockWaitError }
    | SKIP LOCKED               { $$ = nodes.LockWaitSkip }
    | /* EMPTY */               { $$ = nodes.LockWaitBlock }
    ;
```

**Step 2: Wire into select_limit or simple_select**

Add `for_locking_clause` to the SELECT parsing path (check where `select_no_parens` handles it).

**Step 3: Build, test, commit**

```bash
git commit -m "feat(gram): add FOR UPDATE/SHARE/KEY SHARE locking clauses"
```

---

### Task 12: Extend GROUP BY — CUBE, ROLLUP, GROUPING SETS, empty grouping set

**Impact:** CUBE, ROLLUP, GROUPING SETS appear in analytics tests (~200+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — ensure GroupingSet node exists

**Step 1: Add grouping rules**

```yacc
group_by_item:
    a_expr                  { $$ = $1 }
    | empty_grouping_set    { $$ = $1 }
    | cube_clause           { $$ = $1 }
    | rollup_clause         { $$ = $1 }
    | grouping_sets_clause  { $$ = $1 }
    ;

empty_grouping_set:
    '(' ')'
        {
            $$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_EMPTY}
        }
    ;

cube_clause:
    CUBE '(' expr_list ')'
        {
            $$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_CUBE, Content: $3}
        }
    ;

rollup_clause:
    ROLLUP '(' expr_list ')'
        {
            $$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_ROLLUP, Content: $3}
        }
    ;

grouping_sets_clause:
    GROUPING SETS '(' group_by_list ')'
        {
            $$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_SETS, Content: $4}
        }
    ;
```

**Step 2: Build, test, commit**

```bash
git commit -m "feat(gram): add CUBE, ROLLUP, GROUPING SETS, empty grouping set"
```

---

### Task 13: Complete SimpleTypename — ConstDatetime, ConstInterval, Bit, Character, JsonType

**Impact:** TIMESTAMP, TIME, INTERVAL, BIT, JSON type names appear in many CREATE TABLE and CAST statements (~300+ failures).

**Files:**
- Modify: `parser/gram.y`

**Step 1: Check which SimpleTypename alternatives already exist**

Our grammar has 4 alternatives. PG has 8. Add the missing ones: Bit, ConstDatetime, ConstInterval with opt_interval, JsonType.

**Step 2: Add missing type rules**

```yacc
SimpleTypename:
    GenericType             { $$ = $1 }
    | Numeric               { $$ = $1 }
    | Bit                   { $$ = $1 }
    | Character             { $$ = $1 }
    | ConstDatetime         { $$ = $1 }
    | ConstInterval opt_interval
        {
            $$ = $1
            $$.(*nodes.TypeName).Typmods = $2
        }
    | ConstInterval '(' Iconst ')'
        {
            $$ = $1
            $$.(*nodes.TypeName).Typmods = &nodes.List{Items: []nodes.Node{
                makeIntConst(nodes.INTERVAL_FULL_RANGE),
                makeIntConst($3),
            }}
        }
    | JsonType              { $$ = $1 }
    ;
```

Add each sub-rule (ConstDatetime, ConstInterval, Bit, BitWithLength, BitWithoutLength, JsonType, opt_interval, interval_second) following the PG grammar verbatim, translated to Go.

**Step 3: Build, test, commit**

```bash
git commit -m "feat(gram): complete SimpleTypename with datetime, interval, bit, json types"
```

---

### Task 14: Add Partition rules — PartitionSpec, PartitionBoundSpec, partition_cmd

**Impact:** CREATE TABLE ... PARTITION BY and ALTER TABLE ATTACH/DETACH PARTITION appear in ~200+ statements.

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — ensure partition node types exist

**Step 1: Add PartitionSpec, part_params, part_elem**

Translate from PG grammar, adding PartitionSpec, PartitionElem node types if needed.

**Step 2: Add PartitionBoundSpec**

Add FOR VALUES IN/FROM..TO/WITH, DEFAULT alternatives.

**Step 3: Add partition_cmd for ALTER TABLE**

Add ATTACH PARTITION, DETACH PARTITION alternatives.

**Step 4: Wire OptPartitionSpec into CreateStmt**

**Step 5: Build, test, commit**

```bash
git commit -m "feat(gram): add partition rules (PARTITION BY, FOR VALUES, ATTACH/DETACH)"
```

---

### Task 15: Expand CreateStmt — OF type, PARTITION OF, IF NOT EXISTS variants

**Impact:** CREATE TABLE OF, CREATE TABLE PARTITION OF, and IF NOT EXISTS variants (~100+ failures).

**Files:**
- Modify: `parser/gram.y`

**Step 1: Add missing CreateStmt alternatives**

Our grammar has 2 alternatives (basic and IF NOT EXISTS). PG has 6. Add: OF type (2 variants), PARTITION OF (2 variants).

**Step 2: Build, test, commit**

```bash
git commit -m "feat(gram): expand CreateStmt with OF type, PARTITION OF variants"
```

---

### Task 16: Expand alter_table_cmd — trigger/rule enable/disable, identity, RLS, etc.

**Impact:** 47 missing alternatives. Covers ALTER TABLE trigger/rule management, identity columns, row-level security, replica identity, etc. (~400+ failures).

**Files:**
- Modify: `parser/gram.y`

**Step 1: Add column-level ALTER alternatives**

Add the missing column-level operations: SET EXPRESSION, DROP EXPRESSION, SET STATISTICS (by column number), SET/RESET options, SET STORAGE, SET COMPRESSION, ADD GENERATED AS IDENTITY, SET IDENTITY, DROP IDENTITY.

**Step 2: Add trigger/rule management alternatives**

Add ENABLE/DISABLE TRIGGER (name/ALL/USER), ENABLE ALWAYS/REPLICA TRIGGER, ENABLE/DISABLE RULE, ENABLE ALWAYS/REPLICA RULE.

**Step 3: Add table-level ALTER alternatives**

Add: SET WITHOUT OIDS, CLUSTER ON, SET WITHOUT CLUSTER, SET LOGGED/UNLOGGED, INHERIT/NO INHERIT, OF/NOT OF, SET ACCESS METHOD, REPLICA IDENTITY, ENABLE/DISABLE/FORCE/NO FORCE ROW LEVEL SECURITY, VALIDATE CONSTRAINT, ALTER CONSTRAINT.

**Step 4: Build, test, commit**

```bash
git commit -m "feat(gram): expand alter_table_cmd with triggers, identity, RLS, etc."
```

---

### Task 17: Add CTE enhancements — MATERIALIZED, SEARCH, CYCLE

**Impact:** WITH MATERIALIZED / NOT MATERIALIZED and SEARCH/CYCLE clauses (~100+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — ensure CTESearchClause, CTECycleClause exist

**Step 1: Add opt_materialized rule**

```yacc
opt_materialized:
    MATERIALIZED            { $$ = nodes.CTEMaterializeAlways }
    | NOT MATERIALIZED      { $$ = nodes.CTEMaterializeNever }
    | /* EMPTY */           { $$ = nodes.CTEMaterializeDefault }
    ;
```

**Step 2: Wire into common_table_expr**

**Step 3: Add opt_search_clause and opt_cycle_clause**

Translate from PG grammar.

**Step 4: Build, test, commit**

```bash
git commit -m "feat(gram): add CTE MATERIALIZED, SEARCH, CYCLE clauses"
```

---

### Task 18: Add Foreign Key action rules — key_match, key_actions, key_update, key_delete, key_action

**Impact:** ON UPDATE/DELETE CASCADE/SET NULL/RESTRICT/NO ACTION and MATCH FULL/SIMPLE appear in many CREATE TABLE tests (~150+ failures).

**Files:**
- Modify: `parser/gram.y`

**Step 1: Add key rules**

Translate key_match, key_actions, key_update, key_delete, key_action from PG grammar. Wire into ConstraintElem for FOREIGN KEY constraints.

**Step 2: Build, test, commit**

```bash
git commit -m "feat(gram): add foreign key action rules (ON UPDATE/DELETE, MATCH)"
```

---

## PRIORITY 3: Lower Impact (JSON, XML, advanced features)

### Task 19: Add XML function rules to func_expr_common_subexpr

**Impact:** XMLCONCAT, XMLELEMENT, XMLEXISTS, XMLFOREST, XMLPARSE, XMLPI, XMLROOT, XMLSERIALIZE (~150+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — ensure XmlExpr, XmlSerialize exist

**Step 1: Add XML alternatives to func_expr_common_subexpr**

Add all XML alternatives: XMLCONCAT, XMLELEMENT (4 variants), XMLEXISTS, XMLFOREST, XMLPARSE, XMLPI (2 variants), XMLROOT, XMLSERIALIZE.

**Step 2: Add helper rules**

Add xml_attributes, xml_attribute_list, xml_attribute_el, xml_whitespace_option, xml_root_version, opt_xml_root_standalone, xmlexists_argument, xml_passing_mech, document_or_content.

**Step 3: Build, test, commit**

```bash
git commit -m "feat(gram): add XML function rules (XMLELEMENT, XMLFOREST, XMLPARSE, etc.)"
```

---

### Task 20: Add XMLTABLE rule

**Impact:** XMLTABLE in FROM clause (~50+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — ensure RangeTableFunc, RangeTableFuncCol exist

**Step 1: Add xmltable rule and helpers**

Add xmltable, xmltable_column_list, xmltable_column_el, xmltable_column_option_list, xmltable_column_option_el, xml_namespace_list, xml_namespace_el.

**Step 2: Build, test, commit**

```bash
git commit -m "feat(gram): add XMLTABLE rule for FROM clause"
```

---

### Task 21: Add JSON constructor rules

**Impact:** JSON_OBJECT, JSON_ARRAY, JSON, JSON_SCALAR, JSON_SERIALIZE (~100+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — add JSON node types

**Step 1: Add JSON helper rules**

Add: json_value_expr, json_value_expr_list, json_format_clause, json_format_clause_opt, json_returning_clause_opt, json_key_uniqueness_constraint_opt, json_name_and_value, json_name_and_value_list, json_object_constructor_null_clause_opt, json_array_constructor_null_clause_opt.

**Step 2: Add JSON constructor alternatives to func_expr_common_subexpr**

Add JSON_OBJECT (3 variants), JSON_ARRAY (3 variants), JSON(), JSON_SCALAR(), JSON_SERIALIZE().

**Step 3: Build, test, commit**

```bash
git commit -m "feat(gram): add JSON constructor rules (JSON_OBJECT, JSON_ARRAY, etc.)"
```

---

### Task 22: Add JSON query/function rules

**Impact:** JSON_QUERY, JSON_EXISTS, JSON_VALUE, JSON_TABLE (~80+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — add JsonFuncExpr, JsonTable node types

**Step 1: Add JSON behavior/passing rules**

Add: json_behavior, json_behavior_type, json_behavior_clause_opt, json_on_error_clause_opt, json_passing_clause_opt, json_arguments, json_argument, json_quotes_clause_opt, json_wrapper_behavior, json_predicate_type_constraint.

**Step 2: Add JSON_QUERY, JSON_EXISTS, JSON_VALUE to func_expr_common_subexpr**

**Step 3: Add json_table rule for FROM clause**

Add: json_table, json_table_column_definition, json_table_column_definition_list, json_table_column_path_clause_opt, json_table_path_name_opt.

**Step 4: Add json_table to table_ref (alongside xmltable)**

**Step 5: Build, test, commit**

```bash
git commit -m "feat(gram): add JSON query functions and JSON_TABLE"
```

---

### Task 23: Add JSON aggregate rules

**Impact:** JSON_OBJECTAGG, JSON_ARRAYAGG (~30+ failures).

**Files:**
- Modify: `parser/gram.y`
- Modify: `nodes/` — add JsonObjectAgg, JsonArrayAgg, JsonAggConstructor

**Step 1: Add json_aggregate_func rule**

```yacc
json_aggregate_func:
    JSON_OBJECTAGG '(' json_name_and_value
        json_object_constructor_null_clause_opt
        json_key_uniqueness_constraint_opt
        json_returning_clause_opt ')'
        { /* build JsonObjectAgg node */ }
    | JSON_ARRAYAGG '(' json_value_expr
        json_array_aggregate_order_by_clause_opt
        json_array_constructor_null_clause_opt
        json_returning_clause_opt ')'
        { /* build JsonArrayAgg node */ }
    ;
```

**Step 2: Wire into func_expr**

Add `json_aggregate_func filter_clause over_clause` as a func_expr alternative.

**Step 3: Build, test, commit**

```bash
git commit -m "feat(gram): add JSON aggregate functions (JSON_OBJECTAGG, JSON_ARRAYAGG)"
```

---

### Task 24: Add MERGE_ACTION and remaining special functions

**Impact:** MERGE_ACTION() for MERGE statement support (~20+ failures).

**Step 1: Add MERGE_ACTION to func_expr_common_subexpr**

**Step 2: Build, test, commit**

---

### Task 25: Final — Update regression baseline

**Step 1: Run full regression stats**

```bash
cd parser/pgregress && go test -run TestPGRegressStats -v -count=1
```

Note the new pass rate (target: 90%+).

**Step 2: Update known_failures.json**

```bash
cd parser/pgregress && go test -run TestPGRegress -update -v -count=1
```

**Step 3: Run full test suite**

```bash
go test ./...
```

**Step 4: Final commit**

```bash
git add parser/pgregress/known_failures.json
git commit -m "chore: update known_failures.json baseline after gram.y completion"
```

---

## Summary

| Priority | Tasks | Est. New Statements Passing |
|----------|-------|-----------------------------|
| P1 (Tasks 1-10) | func_expr_common_subexpr, a_expr, b_expr, window functions, table_ref | ~9,000-10,000 |
| P2 (Tasks 11-18) | FOR UPDATE, GROUP BY ext, types, partitions, CreateStmt, ALTER TABLE, CTE, FK | ~2,500-3,500 |
| P3 (Tasks 19-24) | XML functions, XMLTABLE, JSON constructors, JSON query, JSON aggs, MERGE_ACTION | ~1,000-1,500 |

**Expected final pass rate: 90-95%** (up from 68.4%)

The remaining 5-10% failures will likely be:
- psql variable interpolation (2.8% — inherent, cannot fix in parser)
- Very obscure syntax forms
- Extractor edge cases
