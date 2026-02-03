-- Phase 5: coverage for remaining grammar-gap plan items

-- JSON behavior: bare EMPTY on empty/error
SELECT JSON_QUERY(jsonb '[]', '$[*]' EMPTY ON EMPTY);
SELECT JSON_VALUE(jsonb '1', '$' EMPTY ON ERROR);
SELECT * FROM JSON_TABLE('[]', 'strict $.a' COLUMNS (js2 int PATH '$' EMPTY ON ERROR));

-- JSON_TABLE formatted column
SELECT * FROM JSON_TABLE(jsonb 'null', 'lax $[*]' COLUMNS (jst text FORMAT JSON PATH '$'));
SELECT * FROM JSON_TABLE(jsonb '{"a":"1"}', '$' COLUMNS (a text FORMAT JSON ENCODING UTF8 PATH '$.a'));

-- CREATE STATISTICS with func_expr_windowless
CREATE STATISTICS s1 ON my_func(a, b) FROM t1;
CREATE STATISTICS s2 (dependencies) ON lower(col1), col2 FROM t1;

-- Ordered-set aggregate args
CREATE AGGREGATE my_percentile(ORDER BY float8) (
  sfunc = ordered_set_transition,
  stype = internal,
  finalfunc = percentile_disc_final,
  finalfunc_extra
);
CREATE AGGREGATE my_percentile2(float8 ORDER BY float8) (
  sfunc = ordered_set_transition,
  stype = internal,
  finalfunc = percentile_disc_final,
  finalfunc_extra
);
DROP AGGREGATE my_percentile(ORDER BY float8);
DROP AGGREGATE my_percentile2(float8 ORDER BY float8);

-- %TYPE for function param/return types
CREATE FUNCTION foo(x hobbies_r.name%TYPE) RETURNS hobbies_r.person%TYPE AS 'select 1' LANGUAGE SQL;
CREATE FUNCTION bar(SETOF mytable.col%TYPE) RETURNS void AS 'select 1' LANGUAGE SQL;

-- Task 1: TRANSFORM in CREATE FUNCTION
CREATE FUNCTION my_func(jsonb) RETURNS int LANGUAGE sql TRANSFORM FOR TYPE jsonb AS $$ SELECT 1; $$;

-- oper_argtypes single-arg error should report missing argument
-- expect-error: missing argument
DROP OPERATOR +(integer);

-- Task 2: LOGGED/UNLOGGED sequences
CREATE SEQUENCE s_logged AS integer LOGGED;
CREATE SEQUENCE s_unlogged AS integer UNLOGGED;
