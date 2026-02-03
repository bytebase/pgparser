-- Phase 4: admin/privilege grammar coverage
GRANT INSERT ON atest2 TO regress_priv_user4 GRANTED BY CURRENT_USER;
REVOKE ALL ON LARGE OBJECT 2001, 2002 FROM PUBLIC GRANTED BY CURRENT_ROLE;
REVOKE GRANT OPTION FOR TRUNCATE ON atest2 FROM regress_priv_user4 GRANTED BY regress_priv_user5;

-- Task 4: SecLabelStmt missing alternatives
SECURITY LABEL ON FUNCTION my_func(int) IS 'secret';
SECURITY LABEL ON PROCEDURE my_proc(text) IS 'secret';
SECURITY LABEL ON ROUTINE my_routine(int, text) IS 'secret';
SECURITY LABEL ON AGGREGATE my_agg(int) IS 'secret';
SECURITY LABEL ON LARGE OBJECT 12345 IS 'secret';
