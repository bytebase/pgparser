-- Phase 4: DDL coverage
CREATE UNLOGGED MATERIALIZED VIEW mv_unlogged AS SELECT 1;

-- Task 3: hash_partbound list
CREATE TABLE hp_test (a int) PARTITION BY HASH (a);
CREATE TABLE hp_child PARTITION OF hp_test FOR VALUES WITH (MODULUS 4, REMAINDER 0);
