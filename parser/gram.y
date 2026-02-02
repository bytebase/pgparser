// PostgreSQL grammar for Go - translated from PostgreSQL's gram.y
// This grammar is compatible with goyacc (golang.org/x/tools/cmd/goyacc)
//
// To generate the parser:
//   go run golang.org/x/tools/cmd/goyacc -o parser.go -p "pg" gram.y

%{
package parser

import (
	"fmt"
	"github.com/pgparser/pgparser/nodes"
)

%}

// Union type for semantic values
%union {
	node      nodes.Node
	list      *nodes.List
	str       string
	ival      int64
	boolean   bool
	ivalp     *int64
	jtype     int        // for JoinType values
	setop     int        // for SetOperation values
	slimit    *SelectLimit  // for LIMIT/OFFSET
	typename  *nodes.TypeName  // for Typename values
	partspec  *nodes.PartitionSpec
	partbound *nodes.PartitionBoundSpec
}

// Token types from the lexer
// Non-keyword tokens must be listed first
%token <str>   IDENT FCONST SCONST BCONST XCONST Op
%token <ival>  ICONST PARAM
%token         TYPECAST DOT_DOT COLON_EQUALS EQUALS_GREATER
%token         LESS_EQUALS GREATER_EQUALS NOT_EQUALS
%token         NOT_LA    /* lookahead token for NOT LIKE etc */
%token         WITH_LA   /* lookahead token for WITH TIME ZONE */
%token         NULLS_LA  /* lookahead token for NULLS FIRST/LAST */

// Keyword tokens (must match keywords.go)
%token <str> ABORT_P ABSENT ABSOLUTE_P ACCESS ACTION ADD_P ADMIN AFTER
%token <str> AGGREGATE ALL ALSO ALTER ALWAYS ANALYSE ANALYZE AND ANY ARRAY AS ASC
%token <str> ASENSITIVE ASSERTION ASSIGNMENT ASYMMETRIC ATOMIC AT ATTACH ATTRIBUTE AUTHORIZATION
%token <str> BACKWARD BEFORE BEGIN_P BETWEEN BIGINT BINARY BIT
%token <str> BOOLEAN_P BOTH BREADTH BY
%token <str> CACHE CALL CALLED CASCADE CASCADED CASE CAST CATALOG_P CHAIN CHAR_P
%token <str> CHARACTER CHARACTERISTICS CHECK CHECKPOINT CLASS CLOSE
%token <str> CLUSTER COALESCE COLLATE COLLATION COLUMN COLUMNS COMMENT COMMENTS COMMIT
%token <str> COMMITTED COMPRESSION CONCURRENTLY CONDITIONAL CONFIGURATION CONFLICT
%token <str> CONNECTION CONSTRAINT CONSTRAINTS CONTENT_P CONTINUE_P CONVERSION_P COPY
%token <str> COST CREATE CROSS CSV CUBE CURRENT_P
%token <str> CURRENT_CATALOG CURRENT_DATE CURRENT_ROLE CURRENT_SCHEMA
%token <str> CURRENT_TIME CURRENT_TIMESTAMP CURRENT_USER CURSOR CYCLE
%token <str> DATA_P DATABASE DAY_P DEALLOCATE DEC DECIMAL_P DECLARE DEFAULT DEFAULTS
%token <str> DEFERRABLE DEFERRED DEFINER DELETE_P DELIMITER DELIMITERS DEPENDS DEPTH DESC
%token <str> DETACH DICTIONARY DISABLE_P DISCARD DISTINCT DO DOCUMENT_P DOMAIN_P
%token <str> DOUBLE_P DROP
%token <str> EACH ELSE EMPTY_P ENABLE_P ENCODING ENCRYPTED END_P ENUM_P ERROR_P ESCAPE
%token <str> EVENT EXCEPT EXCLUDE EXCLUDING EXCLUSIVE EXECUTE EXISTS EXPLAIN EXPRESSION
%token <str> EXTENSION EXTERNAL EXTRACT
%token <str> FALSE_P FAMILY FETCH FILTER FINALIZE FIRST_P FLOAT_P FOLLOWING FOR
%token <str> FORCE FOREIGN FORMAT FORWARD FREEZE FROM FULL FUNCTION FUNCTIONS
%token <str> GENERATED GLOBAL GRANT GRANTED GREATEST GROUP_P GROUPING GROUPS
%token <str> HANDLER HAVING HEADER_P HOLD HOUR_P
%token <str> IDENTITY_P IF_P ILIKE IMMEDIATE IMMUTABLE IMPLICIT_P IMPORT_P IN_P INCLUDE
%token <str> INCLUDING INCREMENT INDENT INDEX INDEXES INHERIT INHERITS INITIALLY INLINE_P
%token <str> INNER_P INOUT INPUT_P INSENSITIVE INSERT INSTEAD INT_P INTEGER
%token <str> INTERSECT INTERVAL INTO INVOKER IS ISNULL ISOLATION
%token <str> JOIN JSON JSON_ARRAY JSON_ARRAYAGG JSON_EXISTS JSON_OBJECT JSON_OBJECTAGG
%token <str> JSON_QUERY JSON_SCALAR JSON_SERIALIZE JSON_TABLE JSON_VALUE
%token <str> KEEP KEY KEYS
%token <str> LABEL LANGUAGE LARGE_P LAST_P LATERAL_P
%token <str> LEADING LEAKPROOF LEAST LEFT LEVEL LIKE LIMIT LISTEN LOAD LOCAL
%token <str> LOCALTIME LOCALTIMESTAMP LOCATION LOCK_P LOCKED LOGGED
%token <str> MAPPING MATCH MATCHED MATERIALIZED MAXVALUE MERGE MERGE_ACTION METHOD
%token <str> MINUTE_P MINVALUE MODE MONTH_P MOVE
%token <str> NAME_P NAMES NATIONAL NATURAL NCHAR NESTED NEW NEXT NFC NFD NFKC NFKD NO
%token <str> NONE NORMALIZE NORMALIZED
%token <str> NOT NOTHING NOTIFY NOTNULL NOWAIT NULL_P NULLIF
%token <str> NULLS_P NUMERIC
%token <str> OBJECT_P OF OFF OFFSET OIDS OLD OMIT ON ONLY OPERATOR OPTION OPTIONS OR
%token <str> ORDER ORDINALITY OTHERS OUT_P OUTER_P
%token <str> OVER OVERLAPS OVERLAY OVERRIDING OWNED OWNER
%token <str> PARALLEL PARAMETER PARSER PARTIAL PARTITION PASSING PASSWORD PATH
%token <str> PLACING PLAN PLANS POLICY
%token <str> POSITION PRECEDING PRECISION PRESERVE PREPARE PREPARED PRIMARY
%token <str> PRIOR PRIVILEGES PROCEDURAL PROCEDURE PROCEDURES PROGRAM PUBLICATION
%token <str> QUOTE QUOTES
%token <str> RANGE READ REAL REASSIGN RECHECK RECURSIVE REF_P REFERENCES REFERENCING
%token <str> REFRESH REINDEX RELATIVE_P RELEASE RENAME REPEATABLE REPLACE REPLICA
%token <str> RESET RESTART RESTRICT RETURN RETURNING RETURNS REVOKE RIGHT ROLE ROLLBACK ROLLUP
%token <str> ROUTINE ROUTINES ROW ROWS RULE
%token <str> SAVEPOINT SCALAR SCHEMA SCHEMAS SCROLL SEARCH SECOND_P SECURITY SELECT
%token <str> SEQUENCE SEQUENCES
%token <str> SERIALIZABLE SERVER SESSION SESSION_USER SET SETS SETOF SHARE SHOW
%token <str> SIMILAR SIMPLE SKIP SMALLINT SNAPSHOT SOME SOURCE SQL_P STABLE STANDALONE_P
%token <str> START STATEMENT STATISTICS STDIN STDOUT STORAGE STORED STRICT_P STRING_P STRIP_P
%token <str> SUBSCRIPTION SUBSTRING SUPPORT SYMMETRIC SYSID SYSTEM_P SYSTEM_USER
%token <str> TABLE TABLES TABLESAMPLE TABLESPACE TARGET TEMP TEMPLATE TEMPORARY TEXT_P THEN
%token <str> TIES TIME TIMESTAMP TO TRAILING TRANSACTION TRANSFORM
%token <str> TREAT TRIGGER TRIM TRUE_P
%token <str> TRUNCATE TRUSTED TYPE_P TYPES_P
%token <str> UESCAPE UNBOUNDED UNCONDITIONAL UNCOMMITTED UNENCRYPTED UNION UNIQUE UNKNOWN
%token <str> UNLISTEN UNLOGGED UNTIL UPDATE USER USING
%token <str> VACUUM VALID VALIDATE VALIDATOR VALUE_P VALUES VARCHAR VARIADIC VARYING
%token <str> VERBOSE VERSION_P VIEW VIEWS VOLATILE
%token <str> WHEN WHERE WHITESPACE_P WINDOW WITH WITHIN WITHOUT WORK WRAPPER WRITE
%token <str> XML_P XMLATTRIBUTES XMLCONCAT XMLELEMENT XMLEXISTS XMLFOREST XMLNAMESPACES
%token <str> XMLPARSE XMLPI XMLROOT XMLSERIALIZE XMLTABLE
%token <str> YEAR_P YES_P
%token <str> ZONE

// Non-terminals with types
%start stmtblock

%type <list>  stmtblock
%type <node>  stmt SelectStmt simple_select select_clause
%type <node>  select_with_parens select_no_parens
%type <node>  a_expr b_expr c_expr columnref AexprConst func_expr func_application func_expr_common_subexpr
%type <node>  target_el where_clause where_or_current_clause
%type <list>  func_name
%type <list>  from_clause
%type <node>  table_ref relation_expr joined_table join_qual
%type <node>  func_table tablesample_clause
%type <boolean>  opt_ordinality
%type <node>  opt_repeatable_clause
%type <ival>  join_type
%type <list>  stmtmulti target_list from_list opt_target_list
%type <list>  group_clause group_by_list
%type <node>  group_by_item having_clause
%type <node>  empty_grouping_set cube_clause rollup_clause grouping_sets_clause
%type <list>  sort_clause opt_sort_clause sortby_list
%type <list>  expr_list opt_indirection indirection func_arg_list
%type <list>  extract_list overlay_list position_list substr_list trim_list
%type <str>   extract_arg unicode_normal_form
%type <node>  sortby indirection_el func_arg_expr
%type <str>   name ColId ColLabel attr_name
%type <node>   opt_alias_clause alias_clause func_alias_clause
%type <str>   unreserved_keyword col_name_keyword type_func_name_keyword reserved_keyword
%type <list>  opt_name_list name_list any_name qualified_name
%type <ival>  opt_asc_desc
%type <ival>  Iconst
%type <str>   Sconst
%type <boolean> opt_all_clause opt_distinct_clause set_quantifier
%type <node>  with_clause opt_with_clause common_table_expr
%type <ival>  opt_materialized
%type <list>  cte_list
%type <list>  window_clause window_definition_list
%type <node>  window_definition window_specification over_clause opt_frame_clause frame_extent frame_bound
%type <str>   opt_existing_window_name
%type <list>  opt_partition_clause
%type <ival>  opt_window_exclusion_clause
%type <node>  within_group_clause filter_clause
%type <node>  func_expr_windowless
%type <slimit> select_limit limit_clause opt_select_limit
%type <list>  for_locking_clause opt_for_locking_clause for_locking_items
%type <node>  for_locking_item
%type <ival>  for_locking_strength opt_nowait_or_skip
%type <list>  locked_rels_list
%type <node>   offset_clause select_limit_value select_offset_value select_fetch_first_value
%type <ival>   sub_type
%type <list>   subquery_Op
%type <node>   case_expr when_clause case_default case_arg
%type <list>   when_clause_list
%type <typename>  Typename SimpleTypename GenericType Numeric opt_float Character
%type <typename>  ConstDatetime ConstInterval Bit BitWithLength BitWithoutLength
%type <typename>  ConstTypename ConstBit ConstCharacter
%type <boolean>  opt_timezone
%type <list>  opt_interval
%type <list>  interval_second
%type <list>      opt_type_modifiers opt_array_bounds
%type <str>       type_function_name
%type <boolean>   opt_varying
%type <node>      array_expr opt_slice_bound
%type <list>      array_expr_list
%type <node>      row explicit_row
%type <list>      implicit_row
%type <node>  InsertStmt insert_rest
%type <node>  insert_target
%type <list>  insert_column_list returning_clause
%type <node>  insert_column_item
%type <node>  opt_on_conflict
%type <node>  values_clause
%type <list>  set_clause_list set_target_list
%type <node>  set_clause set_target
%type <node>  into_clause OptTempTableName
%type <list>  rowsfrom_list opt_col_def_list
%type <node>  rowsfrom_item
%type <node>  opt_search_clause opt_cycle_clause
%type <node>  join_using_alias
%type <node>  UpdateStmt
%type <node>  DeleteStmt
%type <list>  using_clause
%type <node>  relation_expr_opt_alias
%type <node>  CreateStmt
%type <ival>  OptTemp
%type <list>  OptTableElementList TableElementList OptTypedTableElementList TypedTableElementList
%type <node>  TypedTableElement columnOptions
%type <node>  TableElement columnDef TableConstraint TableLikeClause
%type <ival>  TableLikeOptionList TableLikeOption
%type <list>  ColConstraintList opt_column_constraints
%type <node>  ColConstraint ColConstraintElem ConstraintAttr ConstraintElem
%type <node>  DomainConstraint DomainConstraintElem
%type <list>  opt_column_list columnList
%type <partspec>  OptPartitionSpec PartitionSpec_rule
%type <list>  OptWith
%type <ival>  OnCommitOption
%type <str>   OptTableSpace OptAccessMethod
%type <list>  part_params OptInherit opt_collate
%type <node>  part_elem
%type <partbound>  PartitionBoundSpec ForValues
%type <node>  AlterTableStmt RenameStmt DropStmt
%type <list>  alter_table_cmds
%type <node>  alter_table_cmd
%type <ival>  opt_drop_behavior
%type <ival>  object_type_any_name
%type <list>  any_name_list
%type <node>  IndexStmt
%type <boolean>  opt_unique opt_concurrently
%type <str>   opt_single_name access_method_clause
%type <list>  index_params opt_include index_including_params
%type <node>  index_elem
%type <ival>  opt_nulls_order
%type <node>  ViewStmt
%type <ival>  opt_check_option
%type <node>  CreateFunctionStmt
%type <boolean>  opt_or_replace
%type <list>  func_args_with_defaults func_args_with_defaults_list
%type <node>  func_arg_with_default func_arg
%type <list>  table_func_column_list
%type <node>  table_func_column
%type <str>   param_name
%type <ival>  arg_class
%type <typename>  func_return func_type
%type <list>  createfunc_opt_list opt_createfunc_opt_list
%type <node>  createfunc_opt_item common_func_opt_item
%type <node>  opt_routine_body ReturnStmt routine_body_stmt
%type <list>  routine_body_stmt_list
%type <node>  TransactionStmt
%type <boolean>  opt_transaction_chain
%type <node>  ExplainStmt ExplainableStmt
%type <node>  CopyStmt
%type <boolean>  copy_from
%type <str>   copy_file_name
%type <list>  utility_option_list copy_opt_list transaction_mode_list transaction_mode_list_or_empty
%type <node>  copy_opt_item transaction_mode_item
%type <str>   iso_level
%type <node>  utility_option_elem
%type <str>   utility_option_name
%type <node>  utility_option_arg
%type <list>  utility_option_arg_list
%type <str>   opt_boolean_or_string
%type <node>  NumericOnly
%type <node>  GrantStmt RevokeStmt
%type <list>  privileges privilege_list grantee_list
%type <node>  privilege grantee RoleSpec
%type <boolean> opt_grant_grant_option
%type <node>  CheckPointStmt DiscardStmt ListenStmt UnlistenStmt NotifyStmt LoadStmt ClosePortalStmt ConstraintsSetStmt
%type <list>  constraints_set_list
%type <boolean> constraints_set_mode
%type <str>   file_name cursor_name
%type <list>  qualified_name_list
%type <node>  VariableSetStmt VariableShowStmt VariableResetStmt
%type <node>  PrepareStmt ExecuteStmt DeallocateStmt
%type <node>  set_rest set_rest_more generic_set
%type <str>   var_name
%type <list>  var_list
%type <node>  var_value zone_value
%type <str>   NonReservedWord_or_Sconst NonReservedWord
%type <ival>  document_or_content
%type <list>  prep_type_clause execute_param_clause
%type <str>   opt_encoding
%type <node>  reset_rest generic_reset
%type <list>  type_list
%type <ival>  SignedIconst
%type <node>  PreparableStmt
%type <node>  TruncateStmt LockStmt VacuumStmt AnalyzeStmt ClusterStmt ReindexStmt CommentStmt SecLabelStmt
%type <boolean> opt_restart_seqs opt_nowait opt_verbose opt_analyze opt_full opt_freeze
%type <list>  relation_expr_list opt_vacuum_relation_list vacuum_relation_list
%type <node>  vacuum_relation
%type <ival>  opt_lock lock_type reindex_target_type reindex_target_multitable
%type <str>   cluster_index_specification column_compression comment_text security_label opt_provider
%type <ival>  object_type_name object_type_name_on_any_name
%type <list>  opt_reindex_option_list
%type <node>  DeclareCursorStmt FetchStmt MergeStmt CallStmt DoStmt
%type <ival>  cursor_options opt_hold
%type <node>  fetch_args
%type <node>  merge_when_clause
%type <list>  merge_when_list
%type <node>  merge_update merge_delete merge_insert
%type <node>  opt_merge_when_condition
%type <list>  merge_values_clause
%type <ival>  merge_when_tgt_matched merge_when_tgt_not_matched
%type <ival>  override_kind
%type <list>  dostmt_opt_list
%type <node>  dostmt_opt_item
%type <node>  CreateRoleStmt AlterRoleStmt AlterRoleSetStmt DropRoleStmt
%type <node>  CreateUserStmt CreateGroupStmt AlterGroupStmt
%type <node>  GrantRoleStmt RevokeRoleStmt
%type <list>  OptRoleList AlterOptRoleList
%type <node>  CreateOptRoleElem AlterOptRoleElem
%type <str>   RoleId
%type <list>  role_list
%type <node>  SetResetClause
%type <node>  opt_granted_by
%type <list>  grant_role_opt_list
%type <node>  grant_role_opt
%type <node>  grant_role_opt_value
%type <ival>  add_drop
%type <node>  CreatedbStmt AlterDatabaseStmt AlterDatabaseSetStmt DropdbStmt AlterSystemStmt
%type <list>  createdb_opt_list createdb_opt_items
%type <node>  createdb_opt_item
%type <str>   createdb_opt_name
%type <list>  drop_option_list
%type <node>  drop_option
%type <node>  CreateSchemaStmt CreateSeqStmt AlterSeqStmt
%type <node>  CreateDomainStmt AlterDomainStmt AlterEnumStmt AlterCollationStmt AlterCompositeTypeStmt
%type <list>  OptSchemaEltList OptParenthesizedSeqOptList OptSeqOptList SeqOptList
%type <list>  alter_identity_column_option_list
%type <node>  alter_identity_column_option
%type <ival>  generated_when
%type <node>  SeqOptElem
%type <node>  schema_stmt
%type <list>  alter_type_cmds
%type <node>  alter_type_cmd
%type <node>  alter_column_default
%type <boolean>  opt_if_not_exists
%type <node>  TableFuncElement
%type <node>  opt_collate_clause
%type <node>  DefineStmt CompositeTypeStmt CreateEnumStmt CreateRangeStmt
%type <list>  definition def_list old_aggr_definition old_aggr_list
%type <node>  def_elem def_arg old_aggr_elem
%type <list>  OptTableFuncElementList TableFuncElementList
%type <list>  opt_enum_val_list enum_val_list
%type <list>  aggr_args aggr_args_list
%type <node>  aggr_arg
%type <list>  any_operator qual_all_Op qual_Op
%type <str>   all_Op MathOp
%type <node>  AlterFunctionStmt RemoveFuncStmt RemoveAggrStmt RemoveOperStmt
%type <node>  CreateTrigStmt CreateEventTrigStmt AlterEventTrigStmt
%type <node>  RuleStmt CreatePLangStmt
%type <list>  alterfunc_opt_list
%type <list>  function_with_argtypes_list func_args func_args_list
%type <node>  function_with_argtypes
%type <list>  aggregate_with_argtypes_list
%type <node>  aggregate_with_argtypes
%type <list>  operator_with_argtypes_list oper_argtypes
%type <node>  operator_with_argtypes
%type <ival>  TriggerActionTime
%type <list>  TriggerEvents TriggerOneEvent
%type <boolean>  TriggerForSpec TriggerForType
%type <node>  TriggerWhen
%type <list>  TriggerFuncArgs
%type <node>  TriggerFuncArg
%type <list>  TriggerReferencing TriggerTransitions
%type <node>  TriggerTransition
%type <boolean>  TransitionOldOrNew TransitionRowOrTable
%type <str>   TransitionRelName
%type <list>  event_trigger_when_list
%type <node>  event_trigger_when_item
%type <list>  event_trigger_value_list
%type <ival>  enable_trigger
%type <list>  RuleActionList RuleActionMulti
%type <node>  RuleActionStmt RuleActionStmtOrEmpty
%type <ival>  event
%type <boolean>  opt_instead opt_trusted
%type <list>  handler_name opt_inline_handler opt_validator
%type <ival>  ConstraintAttributeSpec ConstraintAttributeElem
%type <ival>  key_match key_actions key_update key_delete key_action
%type <list>  opt_c_include ExclusionConstraintList
%type <list>  ExclusionConstraintElem
%type <str>   OptConsTableSpace ExistingIndex
%type <node>  hash_partbound hash_partbound_elem
%type <node>  OptConstrFromTable
%type <node>  xml_root_version opt_xml_root_standalone
%type <node>  xmlexists_argument
%type <node>  xmltable xmltable_column_el xmltable_column_option_el
%type <list>  xml_attribute_list xml_attributes xmltable_column_list xmltable_column_option_list
%type <list>  xml_namespace_list
%type <node> xml_attribute_el xml_namespace_el
%type <ival>  xml_whitespace_option xml_indent_option

%type <node>  json_value_expr json_output_clause_opt json_format_clause json_format_clause_opt
%type <node>  json_returning_clause_opt
%type <node>  json_behavior json_behavior_clause_opt json_on_error_clause_opt
%type <node>  json_table json_table_column_definition json_table_column_path_clause_opt
%type <list>  json_name_and_value_list json_value_expr_list
%type <list>  json_arguments json_passing_clause_opt
%type <list>  json_table_column_definition_list
%type <list>  json_table_column_option_list
%type <node>  json_name_and_value json_argument json_table_column_option_el
%type <node>  json_aggregate_func
%type <ival>  json_behavior_type json_predicate_type_constraint
%type <ival>  json_wrapper_behavior json_quotes_clause_opt
%type <ival>  json_key_uniqueness_constraint_opt
%type <ival>  json_object_constructor_null_clause_opt json_array_constructor_null_clause_opt
%type <str>   json_table_path_name_opt path_opt
%type <list>  json_array_aggregate_order_by_clause_opt
%type <list>  attrs
%type <node>  CreateFdwStmt AlterFdwStmt
%type <node>  CreateForeignServerStmt AlterForeignServerStmt
%type <node>  CreateForeignTableStmt
%type <node>  CreateUserMappingStmt AlterUserMappingStmt DropUserMappingStmt
%type <node>  ImportForeignSchemaStmt
%type <list>  create_generic_options alter_generic_options
%type <list>  generic_option_list alter_generic_option_list
%type <node>  generic_option_elem alter_generic_option_elem
%type <str>   generic_option_name
%type <node>  generic_option_arg
%type <node>  fdw_option
%type <list>  fdw_options opt_fdw_options
%type <str>   opt_type foreign_server_version opt_foreign_server_version
%type <node>  auth_ident
%type <ival>  import_qualification_type
%type <node>  import_qualification
%type <node>  CreateExtensionStmt AlterExtensionStmt AlterExtensionContentsStmt
%type <list>  create_extension_opt_list alter_extension_opt_list
%type <node>  create_extension_opt_item alter_extension_opt_item
%type <node>  CreateTableSpaceStmt DropTableSpaceStmt
%type <node>  OptTableSpaceOwner
%type <node>  AlterTblSpcStmt
%type <node>  CreateAmStmt
%type <ival>  am_type
%type <node>  CreatePolicyStmt AlterPolicyStmt
%type <boolean>  RowSecurityDefaultPermissive
%type <str>   RowSecurityDefaultForCmd row_security_cmd
%type <list>  RowSecurityDefaultToRole RowSecurityOptionalToRole
%type <node>  RowSecurityOptionalExpr RowSecurityOptionalWithCheck
%type <node>  CreatePublicationStmt AlterPublicationStmt
%type <list>  opt_definition
%type <node>  PublicationObjSpec
%type <list>  pub_obj_list
%type <list>  reloptions opt_reloptions reloption_list
%type <node>  reloption_elem
%type <node>  CreateSubscriptionStmt AlterSubscriptionStmt DropSubscriptionStmt
%type <list>  OptWhereClause
%type <node>  AlterObjectSchemaStmt AlterOwnerStmt AlterObjectDependsStmt
%type <node>  AlterOperatorStmt AlterTypeStmt
%type <node>  AlterDefaultPrivilegesStmt
%type <node>  AlterTSConfigurationStmt AlterTSDictionaryStmt
%type <node>  CreateStatsStmt AlterStatsStmt
%type <node>  CreateOpClassStmt CreateOpFamilyStmt AlterOpFamilyStmt
%type <node>  DropOpClassStmt DropOpFamilyStmt
%type <node>  CreateCastStmt DropCastStmt
%type <node>  CreateTransformStmt DropTransformStmt
%type <node>  CreateConversionStmt
%type <node>  DropOwnedStmt ReassignOwnedStmt
%type <node>  CreateAsStmt CreateMatViewStmt RefreshMatViewStmt
%type <node>  DefACLAction
%type <list>  DefACLOptionList
%type <node>  DefACLOption
%type <ival>  defacl_privilege_target
%type <list>  operator_def_list
%type <node>  operator_def_elem
%type <node>  operator_def_arg
%type <list>  stats_params
%type <node>  stats_param
%type <list>  opclass_item_list opclass_drop_list
%type <node>  opclass_item opclass_drop
%type <boolean>  opt_default
%type <list>  opt_opfamily opclass_purpose
%type <boolean>  opt_recheck
%type <ival>  cast_context
%type <boolean>  opt_if_exists
%type <list>  transform_element_list
%type <boolean>  opt_no
%type <node>  create_as_target create_mv_target
%type <boolean> opt_with_data
%type <boolean> no_inherit
%type <boolean> opt_unique_null_treatment
%type <list>  opt_qualified_name
%type <node>  set_statistics_value
%type <list>  opt_stat_name_list

// Operator precedence - must match PostgreSQL's gram.y
// Lower precedence listed first

%left      UNION EXCEPT
%left      INTERSECT
%left      OR
%left      AND
%right     NOT
%nonassoc  IS ISNULL NOTNULL
%nonassoc  '<' '>' '=' LESS_EQUALS GREATER_EQUALS NOT_EQUALS
%nonassoc  BETWEEN IN_P LIKE ILIKE SIMILAR NOT_LA
%nonassoc  ESCAPE
%nonassoc  UNBOUNDED
%nonassoc  IDENT PARTITION RANGE ROWS GROUPS PRECEDING FOLLOWING CUBE ROLLUP
%nonassoc  SET KEYS OBJECT_P SCALAR VALUE_P WITH WITHOUT PATH
%left      Op OPERATOR
%left      '+' '-'
%left      '*' '/' '%'
%left      '^'
%left      AT
%left      COLLATE
%right     UMINUS
%left      '[' ']'
%left      '(' ')'
%left      TYPECAST
%left      '.'
%left      JOIN CROSS LEFT FULL RIGHT INNER_P NATURAL

%%

// Top-level rule
stmtblock:
	stmtmulti
		{
			setParseResult(pglex, $1)
		}
	;

stmtmulti:
	stmtmulti ';' stmt
		{
			if $3 != nil {
				$$ = appendList($1, $3)
			} else {
				$$ = $1
			}
		}
	| stmt
		{
			if $1 != nil {
				$$ = makeList($1)
			} else {
				$$ = nil
			}
		}
	;

stmt:
	SelectStmt
		{
			$$ = $1
		}
	| InsertStmt
		{
			$$ = $1
		}
	| UpdateStmt
		{
			$$ = $1
		}
	| DeleteStmt
		{
			$$ = $1
		}
	| CreateStmt
		{
			$$ = $1
		}
	| AlterTableStmt
		{
			$$ = $1
		}
	| RenameStmt
		{
			$$ = $1
		}
	| DropStmt
		{
			$$ = $1
		}
	| IndexStmt
		{
			$$ = $1
		}
	| ViewStmt
		{
			$$ = $1
		}
	| CreateFunctionStmt
		{
			$$ = $1
		}
	| TransactionStmt
		{
			$$ = $1
		}
	| ExplainStmt
		{
			$$ = $1
		}
	| CopyStmt
		{
			$$ = $1
		}
	| GrantStmt
		{
			$$ = $1
		}
	| RevokeStmt
		{
			$$ = $1
		}
	| CheckPointStmt
		{
			$$ = $1
		}
	| DiscardStmt
		{
			$$ = $1
		}
	| ListenStmt
		{
			$$ = $1
		}
	| UnlistenStmt
		{
			$$ = $1
		}
	| NotifyStmt
		{
			$$ = $1
		}
	| LoadStmt
		{
			$$ = $1
		}
	| ClosePortalStmt
		{
			$$ = $1
		}
	| ConstraintsSetStmt
		{
			$$ = $1
		}
	| VariableSetStmt
		{
			$$ = $1
		}
	| VariableShowStmt
		{
			$$ = $1
		}
	| VariableResetStmt
		{
			$$ = $1
		}
	| PrepareStmt
		{
			$$ = $1
		}
	| ExecuteStmt
		{
			$$ = $1
		}
	| DeallocateStmt
		{
			$$ = $1
		}
	| TruncateStmt
		{
			$$ = $1
		}
	| LockStmt
		{
			$$ = $1
		}
	| VacuumStmt
		{
			$$ = $1
		}
	| AnalyzeStmt
		{
			$$ = $1
		}
	| ClusterStmt
		{
			$$ = $1
		}
	| ReindexStmt
		{
			$$ = $1
		}
	| CommentStmt
		{
			$$ = $1
		}
	| SecLabelStmt
		{
			$$ = $1
		}
	| DeclareCursorStmt
		{
			$$ = $1
		}
	| FetchStmt
		{
			$$ = $1
		}
	| MergeStmt
		{
			$$ = $1
		}
	| CallStmt
		{
			$$ = $1
		}
	| DoStmt
		{
			$$ = $1
		}
	| CreateRoleStmt
		{
			$$ = $1
		}
	| AlterRoleStmt
		{
			$$ = $1
		}
	| AlterRoleSetStmt
		{
			$$ = $1
		}
	| DropRoleStmt
		{
			$$ = $1
		}
	| CreateUserStmt
		{
			$$ = $1
		}
	| CreateGroupStmt
		{
			$$ = $1
		}
	| AlterGroupStmt
		{
			$$ = $1
		}
	| GrantRoleStmt
		{
			$$ = $1
		}
	| RevokeRoleStmt
		{
			$$ = $1
		}
	| CreatedbStmt
		{
			$$ = $1
		}
	| AlterDatabaseStmt
		{
			$$ = $1
		}
	| AlterDatabaseSetStmt
		{
			$$ = $1
		}
	| DropdbStmt
		{
			$$ = $1
		}
	| AlterSystemStmt
		{
			$$ = $1
		}
	| CreateSchemaStmt
		{
			$$ = $1
		}
	| CreateSeqStmt
		{
			$$ = $1
		}
	| AlterSeqStmt
		{
			$$ = $1
		}
	| CreateDomainStmt
		{
			$$ = $1
		}
	| AlterDomainStmt
		{
			$$ = $1
		}
	| AlterEnumStmt
		{
			$$ = $1
		}
	| AlterCollationStmt
		{
			$$ = $1
		}
	| AlterCompositeTypeStmt
		{
			$$ = $1
		}
	| DefineStmt
		{
			$$ = $1
		}
	| CompositeTypeStmt
		{
			$$ = $1
		}
	| CreateEnumStmt
		{
			$$ = $1
		}
	| CreateRangeStmt
		{
			$$ = $1
		}
	| AlterFunctionStmt
		{
			$$ = $1
		}
	| RemoveFuncStmt
		{
			$$ = $1
		}
	| RemoveAggrStmt
		{
			$$ = $1
		}
	| RemoveOperStmt
		{
			$$ = $1
		}
	| CreateTrigStmt
		{
			$$ = $1
		}
	| CreateEventTrigStmt
		{
			$$ = $1
		}
	| AlterEventTrigStmt
		{
			$$ = $1
		}
	| RuleStmt
		{
			$$ = $1
		}
	| CreatePLangStmt
		{
			$$ = $1
		}
	| CreateFdwStmt
		{
			$$ = $1
		}
	| AlterFdwStmt
		{
			$$ = $1
		}
	| CreateForeignServerStmt
		{
			$$ = $1
		}
	| AlterForeignServerStmt
		{
			$$ = $1
		}
	| CreateForeignTableStmt
		{
			$$ = $1
		}
	| CreateUserMappingStmt
		{
			$$ = $1
		}
	| AlterUserMappingStmt
		{
			$$ = $1
		}
	| DropUserMappingStmt
		{
			$$ = $1
		}
	| ImportForeignSchemaStmt
		{
			$$ = $1
		}
	| CreateExtensionStmt
		{
			$$ = $1
		}
	| AlterExtensionStmt
		{
			$$ = $1
		}
	| AlterExtensionContentsStmt
		{
			$$ = $1
		}
	| CreateTableSpaceStmt
		{
			$$ = $1
		}
	| DropTableSpaceStmt
		{
			$$ = $1
		}
	| AlterTblSpcStmt
		{
			$$ = $1
		}
	| CreateAmStmt
		{
			$$ = $1
		}
	| CreatePolicyStmt
		{
			$$ = $1
		}
	| AlterPolicyStmt
		{
			$$ = $1
		}
	| CreatePublicationStmt
		{
			$$ = $1
		}
	| AlterPublicationStmt
		{
			$$ = $1
		}
	| CreateSubscriptionStmt
		{
			$$ = $1
		}
	| AlterSubscriptionStmt
		{
			$$ = $1
		}
	| DropSubscriptionStmt
		{
			$$ = $1
		}
	| AlterObjectSchemaStmt
		{
			$$ = $1
		}
	| AlterOwnerStmt
		{
			$$ = $1
		}
	| AlterObjectDependsStmt
		{
			$$ = $1
		}
	| AlterOperatorStmt
		{
			$$ = $1
		}
	| AlterTypeStmt
		{
			$$ = $1
		}
	| AlterDefaultPrivilegesStmt
		{
			$$ = $1
		}
	| AlterTSConfigurationStmt
		{
			$$ = $1
		}
	| AlterTSDictionaryStmt
		{
			$$ = $1
		}
	| CreateStatsStmt
		{
			$$ = $1
		}
	| AlterStatsStmt
		{
			$$ = $1
		}
	| CreateOpClassStmt
		{
			$$ = $1
		}
	| CreateOpFamilyStmt
		{
			$$ = $1
		}
	| AlterOpFamilyStmt
		{
			$$ = $1
		}
	| DropOpClassStmt
		{
			$$ = $1
		}
	| DropOpFamilyStmt
		{
			$$ = $1
		}
	| CreateCastStmt
		{
			$$ = $1
		}
	| DropCastStmt
		{
			$$ = $1
		}
	| CreateTransformStmt
		{
			$$ = $1
		}
	| DropTransformStmt
		{
			$$ = $1
		}
	| CreateConversionStmt
		{
			$$ = $1
		}
	| DropOwnedStmt
		{
			$$ = $1
		}
	| ReassignOwnedStmt
		{
			$$ = $1
		}
	| CreateAsStmt
		{
			$$ = $1
		}
	| CreateMatViewStmt
		{
			$$ = $1
		}
	| RefreshMatViewStmt
		{
			$$ = $1
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

/*****************************************************************************
 *
 *      INSERT statement
 *
 *****************************************************************************/

InsertStmt:
	opt_with_clause INSERT INTO insert_target insert_rest opt_on_conflict returning_clause
		{
			n := $5.(*nodes.InsertStmt)
			n.Relation = $4.(*nodes.RangeVar)
			if $6 != nil {
				n.OnConflictClause = $6.(*nodes.OnConflictClause)
			}
			n.ReturningList = $7
			if $1 != nil {
				n.WithClause = $1.(*nodes.WithClause)
			}
			$$ = n
		}
	;

insert_target:
	qualified_name
		{
			$$ = makeRangeVar($1)
		}
	| qualified_name AS ColId
		{
			rv := makeRangeVar($1)
			rv.(*nodes.RangeVar).Alias = &nodes.Alias{Aliasname: $3}
			$$ = rv
		}
	;

insert_rest:
	SelectStmt
		{
			$$ = &nodes.InsertStmt{
				SelectStmt: $1,
			}
		}
	| OVERRIDING override_kind VALUE_P SelectStmt
		{
			$$ = &nodes.InsertStmt{
				Override:   nodes.OverridingKind($2),
				SelectStmt: $4,
			}
		}
	| '(' insert_column_list ')' SelectStmt
		{
			$$ = &nodes.InsertStmt{
				Cols:       $2,
				SelectStmt: $4,
			}
		}
	| '(' insert_column_list ')' OVERRIDING override_kind VALUE_P SelectStmt
		{
			$$ = &nodes.InsertStmt{
				Cols:       $2,
				Override:   nodes.OverridingKind($5),
				SelectStmt: $7,
			}
		}
	| DEFAULT VALUES
		{
			$$ = &nodes.InsertStmt{}
		}
	;

insert_column_list:
	insert_column_item
		{ $$ = makeList($1) }
	| insert_column_list ',' insert_column_item
		{ $$ = appendList($1, $3) }
	;

insert_column_item:
	ColId opt_indirection
		{
			$$ = &nodes.ResTarget{
				Name:        $1,
				Indirection: $2,
			}
		}
	;

opt_on_conflict:
	ON CONFLICT DO NOTHING
		{
			$$ = &nodes.OnConflictClause{
				Action:   ONCONFLICT_NOTHING,
				Location: -1,
			}
		}
	| ON CONFLICT DO UPDATE SET set_clause_list where_clause
		{
			$$ = &nodes.OnConflictClause{
				Action:      ONCONFLICT_UPDATE,
				TargetList:  $6,
				WhereClause: $7,
				Location:    -1,
			}
		}
	| ON CONFLICT '(' index_params ')' DO NOTHING
		{
			$$ = &nodes.OnConflictClause{
				Action:   ONCONFLICT_NOTHING,
				Infer: &nodes.InferClause{
					IndexElems: $4,
					Location:   -1,
				},
				Location: -1,
			}
		}
	| ON CONFLICT '(' index_params ')' DO UPDATE SET set_clause_list where_clause
		{
			$$ = &nodes.OnConflictClause{
				Action:      ONCONFLICT_UPDATE,
				Infer: &nodes.InferClause{
					IndexElems: $4,
					Location:   -1,
				},
				TargetList:  $9,
				WhereClause: $10,
				Location:    -1,
			}
		}
	| ON CONFLICT '(' index_params ')' WHERE a_expr DO NOTHING
		{
			$$ = &nodes.OnConflictClause{
				Action:   ONCONFLICT_NOTHING,
				Infer: &nodes.InferClause{
					IndexElems:  $4,
					WhereClause: $7,
					Location:    -1,
				},
				Location: -1,
			}
		}
	| ON CONFLICT '(' index_params ')' WHERE a_expr DO UPDATE SET set_clause_list where_clause
		{
			$$ = &nodes.OnConflictClause{
				Action:      ONCONFLICT_UPDATE,
				Infer: &nodes.InferClause{
					IndexElems:  $4,
					WhereClause: $7,
					Location:    -1,
				},
				TargetList:  $11,
				WhereClause: $12,
				Location:    -1,
			}
		}
	| ON CONFLICT ON CONSTRAINT name DO NOTHING
		{
			$$ = &nodes.OnConflictClause{
				Action:   ONCONFLICT_NOTHING,
				Infer: &nodes.InferClause{
					Conname:  $5,
					Location: -1,
				},
				Location: -1,
			}
		}
	| ON CONFLICT ON CONSTRAINT name DO UPDATE SET set_clause_list where_clause
		{
			$$ = &nodes.OnConflictClause{
				Action:      ONCONFLICT_UPDATE,
				Infer: &nodes.InferClause{
					Conname:  $5,
					Location: -1,
				},
				TargetList:  $9,
				WhereClause: $10,
				Location:    -1,
			}
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

returning_clause:
	RETURNING target_list { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

set_clause_list:
	set_clause
		{
			if list, ok := $1.(*nodes.List); ok {
				$$ = list
			} else {
				$$ = makeList($1)
			}
		}
	| set_clause_list ',' set_clause
		{
			if list, ok := $3.(*nodes.List); ok {
				$$ = concatLists($1, list)
			} else {
				$$ = appendList($1, $3)
			}
		}
	;

set_clause:
	set_target '=' a_expr
		{
			rt := $1.(*nodes.ResTarget)
			rt.Val = $3
			$$ = rt
		}
	| '(' set_target_list ')' '=' a_expr
		{
			/* multi-column assignment: (a,b) = expr
			 * Create a list of ResTargets, each with a MultiAssignRef val */
			targets := $2.Items
			ncolumns := len(targets)
			result := &nodes.List{}
			for i, t := range targets {
				rt := t.(*nodes.ResTarget)
				rt.Val = &nodes.MultiAssignRef{
					Source:   $5,
					Colno:   i + 1,
					Ncolumns: ncolumns,
				}
				result.Items = append(result.Items, rt)
			}
			$$ = result
		}
	;

set_target:
	ColId opt_indirection
		{
			$$ = &nodes.ResTarget{
				Name:        $1,
				Indirection: $2,
			}
		}
	;

set_target_list:
	set_target
		{ $$ = makeList($1) }
	| set_target_list ',' set_target
		{ $$ = appendList($1, $3) }
	;

/*****************************************************************************
 *
 *      UPDATE statement
 *
 *****************************************************************************/

UpdateStmt:
	opt_with_clause UPDATE relation_expr_opt_alias SET set_clause_list from_clause where_or_current_clause returning_clause
		{
			$$ = &nodes.UpdateStmt{
				Relation:      $3.(*nodes.RangeVar),
				TargetList:    $5,
				FromClause:    $6,
				WhereClause:   $7,
				ReturningList: $8,
			}
			if $1 != nil {
				$$.(*nodes.UpdateStmt).WithClause = $1.(*nodes.WithClause)
			}
		}
	;

/*****************************************************************************
 *
 *      DELETE statement
 *
 *****************************************************************************/

DeleteStmt:
	opt_with_clause DELETE_P FROM relation_expr_opt_alias using_clause where_or_current_clause returning_clause
		{
			$$ = &nodes.DeleteStmt{
				Relation:      $4.(*nodes.RangeVar),
				UsingClause:   $5,
				WhereClause:   $6,
				ReturningList: $7,
			}
			if $1 != nil {
				$$.(*nodes.DeleteStmt).WithClause = $1.(*nodes.WithClause)
			}
		}
	;

using_clause:
	USING from_list { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

relation_expr_opt_alias:
	relation_expr  %prec UMINUS
		{ $$ = $1 }
	| relation_expr ColId
		{
			$1.(*nodes.RangeVar).Alias = &nodes.Alias{Aliasname: $2}
			$$ = $1
		}
	| relation_expr AS ColId
		{
			$1.(*nodes.RangeVar).Alias = &nodes.Alias{Aliasname: $3}
			$$ = $1
		}
	;

/*****************************************************************************
 *
 *      CREATE TABLE statement
 *
 *****************************************************************************/

CreateStmt:
	CREATE OptTemp TABLE qualified_name '(' OptTableElementList ')' OptInherit OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($4)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				TableElts:      $6,
				InhRelations:   $8,
				Partspec:       $9,
				AccessMethod:   $10,
				Options:        $11,
				OnCommit:       nodes.OnCommitAction($12),
				Tablespacename: $13,
			}
		}
	| CREATE OptTemp TABLE IF_P NOT EXISTS qualified_name '(' OptTableElementList ')' OptInherit OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($7)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				TableElts:      $9,
				InhRelations:   $11,
				Partspec:       $12,
				AccessMethod:   $13,
				Options:        $14,
				OnCommit:       nodes.OnCommitAction($15),
				Tablespacename: $16,
				IfNotExists:    true,
			}
		}
	| CREATE OptTemp TABLE qualified_name PARTITION OF qualified_name ForValues OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($4)
			inh := makeRangeVar($7)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				InhRelations:   makeList(inh),
				Partbound:      $8,
				Partspec:       $9,
				AccessMethod:   $10,
				Options:        $11,
				OnCommit:       nodes.OnCommitAction($12),
				Tablespacename: $13,
			}
		}
	| CREATE OptTemp TABLE IF_P NOT EXISTS qualified_name PARTITION OF qualified_name ForValues OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($7)
			inh := makeRangeVar($10)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				InhRelations:   makeList(inh),
				Partbound:      $11,
				Partspec:       $12,
				AccessMethod:   $13,
				Options:        $14,
				OnCommit:       nodes.OnCommitAction($15),
				Tablespacename: $16,
				IfNotExists:    true,
			}
		}
	| CREATE OptTemp TABLE qualified_name PARTITION OF qualified_name '(' TypedTableElementList ')' ForValues OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($4)
			inh := makeRangeVar($7)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				TableElts:      $9,
				InhRelations:   makeList(inh),
				Partbound:      $11,
				Partspec:       $12,
				AccessMethod:   $13,
				Options:        $14,
				OnCommit:       nodes.OnCommitAction($15),
				Tablespacename: $16,
			}
		}
	| CREATE OptTemp TABLE IF_P NOT EXISTS qualified_name PARTITION OF qualified_name '(' TypedTableElementList ')' ForValues OptPartitionSpec OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($7)
			inh := makeRangeVar($10)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				TableElts:      $12,
				InhRelations:   makeList(inh),
				Partbound:      $14,
				Partspec:       $15,
				AccessMethod:   $16,
				Options:        $17,
				OnCommit:       nodes.OnCommitAction($18),
				Tablespacename: $19,
				IfNotExists:    true,
			}
		}
	/* CREATE TABLE OF typename */
	| CREATE OptTemp TABLE qualified_name OF any_name OptTypedTableElementList OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($4)
			tn := makeTypeNameFromNameList($6).(*nodes.TypeName)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				TableElts:      $7,
				OfTypename:     tn,
				AccessMethod:   $8,
				Options:        $9,
				OnCommit:       nodes.OnCommitAction($10),
				Tablespacename: $11,
			}
		}
	| CREATE OptTemp TABLE IF_P NOT EXISTS qualified_name OF any_name OptTypedTableElementList OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVar($7)
			tn := makeTypeNameFromNameList($9).(*nodes.TypeName)
			$$ = &nodes.CreateStmt{
				Relation:       rv.(*nodes.RangeVar),
				TableElts:      $10,
				OfTypename:     tn,
				AccessMethod:   $11,
				Options:        $12,
				OnCommit:       nodes.OnCommitAction($13),
				Tablespacename: $14,
				IfNotExists:    true,
			}
		}
	;

OptTemp:
	TEMPORARY       { $$ = 1 }
	| TEMP          { $$ = 1 }
	| LOCAL TEMPORARY { $$ = 1 }
	| LOCAL TEMP    { $$ = 1 }
	| UNLOGGED      { $$ = 2 }
	| /* EMPTY */   { $$ = 0 }
	;

OptInherit:
	INHERITS '(' qualified_name_list ')'	{ $$ = $3 }
	| /* EMPTY */							{ $$ = nil }
	;

OptPartitionSpec:
	PartitionSpec_rule	{ $$ = $1 }
	| /* EMPTY */		{ $$ = nil }
	;

OptAccessMethod:
	USING name  { $$ = $2 }
	| /* EMPTY */ { $$ = "" }
	;

OptWith:
	WITH reloptions				{ $$ = $2 }
	| WITH OIDS					{ $$ = nil /* deprecated */ }
	| WITHOUT OIDS				{ $$ = nil }
	| /* EMPTY */				{ $$ = nil }
	;

OnCommitOption:
	ON COMMIT DROP				{ $$ = int64(nodes.ONCOMMIT_DROP) }
	| ON COMMIT DELETE_P ROWS	{ $$ = int64(nodes.ONCOMMIT_DELETE_ROWS) }
	| ON COMMIT PRESERVE ROWS	{ $$ = int64(nodes.ONCOMMIT_PRESERVE_ROWS) }
	| /* EMPTY */				{ $$ = int64(nodes.ONCOMMIT_NOOP) }
	;

OptTableSpace:
	TABLESPACE name				{ $$ = $2 }
	| /* EMPTY */				{ $$ = "" }
	;

PartitionSpec_rule:
	PARTITION BY ColId '(' part_params ')'
		{
			$$ = &nodes.PartitionSpec{
				Strategy:   parsePartitionStrategy($3),
				PartParams: $5,
				Location:   -1,
			}
		}
	;

part_params:
	part_elem
		{ $$ = makeList($1) }
	| part_params ',' part_elem
		{ $$ = appendList($1, $3) }
	;

part_elem:
	ColId opt_collate opt_qualified_name
		{
			$$ = &nodes.PartitionElem{
				Name:      $1,
				Collation: $2,
				Opclass:   $3,
				Location:  -1,
			}
		}
	| func_expr_windowless opt_collate opt_qualified_name
		{
			$$ = &nodes.PartitionElem{
				Expr:      $1,
				Collation: $2,
				Opclass:   $3,
				Location:  -1,
			}
		}
	| '(' a_expr ')' opt_collate opt_qualified_name
		{
			$$ = &nodes.PartitionElem{
				Expr:      $2,
				Collation: $4,
				Opclass:   $5,
				Location:  -1,
			}
		}
	;

opt_collate:
	COLLATE any_name	{ $$ = $2 }
	| /* EMPTY */		{ $$ = nil }
	;

ForValues:
	PartitionBoundSpec	{ $$ = $1 }
	| DEFAULT
		{
			$$ = &nodes.PartitionBoundSpec{
				IsDefault: true,
				Location:  -1,
			}
		}
	;

PartitionBoundSpec:
	/* a LIST partition */
	FOR VALUES IN_P '(' expr_list ')'
		{
			$$ = &nodes.PartitionBoundSpec{
				Strategy:   'l',
				Listdatums: $5,
				Location:   -1,
			}
		}
	/* a RANGE partition */
	| FOR VALUES FROM '(' expr_list ')' TO '(' expr_list ')'
		{
			$$ = &nodes.PartitionBoundSpec{
				Strategy:    'r',
				Lowerdatums: $5,
				Upperdatums: $9,
				Location:    -1,
			}
		}
	/* a HASH partition */
	| FOR VALUES WITH '(' hash_partbound ')'
		{
			$$ = $5.(*nodes.PartitionBoundSpec)
		}
	;

hash_partbound:
	hash_partbound_elem ',' hash_partbound_elem
		{
			n := $1.(*nodes.PartitionBoundSpec)
			n2 := $3.(*nodes.PartitionBoundSpec)
			if n.Modulus != 0 {
				n.Remainder = n2.Remainder
			} else {
				n.Modulus = n2.Modulus
			}
			n.Strategy = 'h'
			n.Location = -1
			$$ = n
		}
	;

hash_partbound_elem:
	NonReservedWord Iconst
		{
			n := &nodes.PartitionBoundSpec{}
			switch $1 {
			case "modulus":
				n.Modulus = int($2)
			case "remainder":
				n.Remainder = int($2)
			}
			$$ = n
		}
	;

OptTableElementList:
	TableElementList { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

OptTypedTableElementList:
	'(' TypedTableElementList ')' { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

TypedTableElementList:
	TypedTableElement
		{ $$ = makeList($1) }
	| TypedTableElementList ',' TypedTableElement
		{ $$ = appendList($1, $3) }
	;

TypedTableElement:
	columnOptions { $$ = $1 }
	| TableConstraint { $$ = $1 }
	;

columnOptions:
	ColId opt_column_constraints
		{
			n := &nodes.ColumnDef{
				Colname:     $1,
				TypeName:    nil,
				IsLocal:     true,
				Location:    -1,
			}
			splitColQualList($2, n)
			$$ = n
		}
	| ColId WITH OPTIONS opt_column_constraints
		{
			n := &nodes.ColumnDef{
				Colname:     $1,
				TypeName:    nil,
				IsLocal:     true,
				Location:    -1,
			}
			splitColQualList($4, n)
			$$ = n
		}
	;

TableElementList:
	TableElement
		{ $$ = makeList($1) }
	| TableElementList ',' TableElement
		{ $$ = appendList($1, $3) }
	;

TableElement:
	columnDef { $$ = $1 }
	| TableConstraint { $$ = $1 }
	| TableLikeClause { $$ = $1 }
	;

TableLikeClause:
	LIKE qualified_name
		{
			$$ = &nodes.TableLikeClause{
				Relation: makeRangeVar($2).(*nodes.RangeVar),
			}
		}
	| LIKE qualified_name TableLikeOptionList
		{
			$$ = &nodes.TableLikeClause{
				Relation: makeRangeVar($2).(*nodes.RangeVar),
				Options:  uint32($3),
			}
		}
	;

TableLikeOptionList:
	TableLikeOptionList INCLUDING TableLikeOption
		{ $$ = $1 | $3 }
	| TableLikeOptionList EXCLUDING TableLikeOption
		{ $$ = $1 &^ $3 }
	| INCLUDING TableLikeOption
		{ $$ = $2 }
	| EXCLUDING TableLikeOption
		{ $$ = 0 &^ $2 }
	;

TableLikeOption:
	ALL				{ $$ = 0xFFFFFFFF }
	| COMMENTS		{ $$ = 1 }
	| COMPRESSION	{ $$ = 2 }
	| CONSTRAINTS	{ $$ = 4 }
	| DEFAULTS		{ $$ = 8 }
	| GENERATED		{ $$ = 16 }
	| IDENTITY_P	{ $$ = 32 }
	| INDEXES		{ $$ = 64 }
	| STATISTICS	{ $$ = 128 }
	| STORAGE		{ $$ = 256 }
	;

columnDef:
	ColId Typename opt_column_constraints
		{
			n := &nodes.ColumnDef{
				Colname:     $1,
				TypeName:    $2,
				IsLocal:     true,
				Location:    -1,
			}
			splitColQualList($3, n)
			$$ = n
		}
	;

opt_column_constraints:
	ColConstraintList { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

ColConstraintList:
	ColConstraint
		{ $$ = makeList($1) }
	| ColConstraintList ColConstraint
		{ $$ = appendList($1, $2) }
	;

ColConstraint:
	CONSTRAINT name ColConstraintElem
		{
			n := $3.(*nodes.Constraint)
			n.Conname = $2
			$$ = n
		}
	| ColConstraintElem { $$ = $1 }
	| COLLATE any_name
		{
			$$ = &nodes.CollateClause{
				Collname: $2,
				Location: -1,
			}
		}
	| COMPRESSION ColId
		{
			/* COMPRESSION is stored directly on ColumnDef, use DefElem as carrier */
			$$ = &nodes.DefElem{
				Defname: "compression",
				Arg:     &nodes.String{Str: $2},
			}
		}
	| STORAGE ColId
		{
			/* STORAGE is stored directly on ColumnDef, use DefElem as carrier */
			$$ = &nodes.DefElem{
				Defname: "storage",
				Arg:     &nodes.String{Str: $2},
			}
		}
	| OPTIONS '(' generic_option_list ')'
		{
			/* OPTIONS for foreign table columns */
			$$ = &nodes.DefElem{
				Defname: "fdwoptions",
				Arg:     $3,
			}
		}
	| ConstraintAttr
		{ $$ = $1 }
	;

ConstraintAttr:
	DEFERRABLE
		{
			$$ = &nodes.Constraint{
				Contype:  nodes.CONSTR_ATTR_DEFERRABLE,
				Location: -1,
			}
		}
	| NOT DEFERRABLE
		{
			$$ = &nodes.Constraint{
				Contype:  nodes.CONSTR_ATTR_NOT_DEFERRABLE,
				Location: -1,
			}
		}
	| INITIALLY DEFERRED
		{
			$$ = &nodes.Constraint{
				Contype:  nodes.CONSTR_ATTR_DEFERRED,
				Location: -1,
			}
		}
	| INITIALLY IMMEDIATE
		{
			$$ = &nodes.Constraint{
				Contype:  nodes.CONSTR_ATTR_IMMEDIATE,
				Location: -1,
			}
		}
	;

ColConstraintElem:
	NOT NULL_P
		{
			$$ = &nodes.Constraint{
				Contype:  nodes.CONSTR_NOTNULL,
				Location: -1,
			}
		}
	| NULL_P
		{
			$$ = &nodes.Constraint{
				Contype:  nodes.CONSTR_NULL,
				Location: -1,
			}
		}
	| UNIQUE opt_unique_null_treatment opt_definition OptConsTableSpace
		{
			$$ = &nodes.Constraint{
				Contype:           nodes.CONSTR_UNIQUE,
				NullsNotDistinct:  $2,
				Options:           $3,
				Indexspace:        $4,
				Location:          -1,
			}
		}
	| PRIMARY KEY opt_definition OptConsTableSpace
		{
			$$ = &nodes.Constraint{
				Contype:    nodes.CONSTR_PRIMARY,
				Options:    $3,
				Indexspace:  $4,
				Location:   -1,
			}
		}
	| CHECK '(' a_expr ')' no_inherit
		{
			n := &nodes.Constraint{
				Contype:  nodes.CONSTR_CHECK,
				RawExpr:  $3,
				Location: -1,
			}
			n.IsNoInherit = $5
			$$ = n
		}
	| DEFAULT b_expr
		{
			$$ = &nodes.Constraint{
				Contype:  nodes.CONSTR_DEFAULT,
				RawExpr:  $2,
				Location: -1,
			}
		}
	| REFERENCES qualified_name opt_column_list key_match key_actions ConstraintAttributeSpec
		{
			rv := makeRangeVar($2)
			n := &nodes.Constraint{
				Contype:      nodes.CONSTR_FOREIGN,
				Pktable:      rv.(*nodes.RangeVar),
				PkAttrs:      $3,
				FkMatchtype:  byte($4),
				FkUpdaction:  byte($5 >> 8),
				FkDelaction:  byte($5 & 0xFF),
				Location:     -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $6)
			$$ = n
		}
	| GENERATED ALWAYS AS IDENTITY_P OptParenthesizedSeqOptList
		{
			$$ = &nodes.Constraint{
				Contype:       nodes.CONSTR_IDENTITY,
				GeneratedWhen: 'a',
				Options:       $5,
				Location:      -1,
			}
		}
	| GENERATED BY DEFAULT AS IDENTITY_P OptParenthesizedSeqOptList
		{
			$$ = &nodes.Constraint{
				Contype:       nodes.CONSTR_IDENTITY,
				GeneratedWhen: 'd',
				Options:       $6,
				Location:      -1,
			}
		}
	| GENERATED ALWAYS AS '(' a_expr ')' STORED
		{
			$$ = &nodes.Constraint{
				Contype:       nodes.CONSTR_GENERATED,
				GeneratedWhen: 'a',
				RawExpr:       $5,
				Location:      -1,
			}
		}
	| GENERATED BY DEFAULT AS '(' a_expr ')' STORED
		{
			$$ = &nodes.Constraint{
				Contype:       nodes.CONSTR_GENERATED,
				GeneratedWhen: 'd',
				RawExpr:       $6,
				Location:      -1,
			}
		}
	;

DomainConstraint:
	CONSTRAINT name DomainConstraintElem
		{
			n := $3.(*nodes.Constraint)
			n.Conname = $2
			$$ = n
		}
	| DomainConstraintElem
		{
			$$ = $1
		}
	;

DomainConstraintElem:
	CHECK '(' a_expr ')' ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:        nodes.CONSTR_CHECK,
				RawExpr:        $3,
				Location:       -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $5)
			$$ = n
		}
	| NOT NULL_P ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:        nodes.CONSTR_NOTNULL,
				Location:       -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $3)
			$$ = n
		}
	;

TableConstraint:
	CONSTRAINT name ConstraintElem
		{
			n := $3.(*nodes.Constraint)
			n.Conname = $2
			$$ = n
		}
	| ConstraintElem
		{
			$$ = $1
		}
	;

ConstraintElem:
	UNIQUE opt_unique_null_treatment '(' columnList ')' opt_c_include opt_definition OptConsTableSpace ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:          nodes.CONSTR_UNIQUE,
				NullsNotDistinct: $2,
				Keys:             $4,
				Including:        $6,
				Options:          $7,
				Indexspace:       $8,
				Location:         -1,
				InitiallyValid:   true,
			}
			applyConstraintAttrs(n, $9)
			$$ = n
		}
	| UNIQUE ExistingIndex ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:    nodes.CONSTR_UNIQUE,
				Indexname:  $2,
				Location:   -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $3)
			$$ = n
		}
	| PRIMARY KEY '(' columnList ')' opt_c_include opt_definition OptConsTableSpace ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:    nodes.CONSTR_PRIMARY,
				Keys:       $4,
				Including:  $6,
				Options:    $7,
				Indexspace:  $8,
				Location:   -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $9)
			$$ = n
		}
	| PRIMARY KEY ExistingIndex ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:    nodes.CONSTR_PRIMARY,
				Indexname:  $3,
				Location:   -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $4)
			$$ = n
		}
	| CHECK '(' a_expr ')' ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:  nodes.CONSTR_CHECK,
				RawExpr:  $3,
				Location: -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $5)
			$$ = n
		}
	| FOREIGN KEY '(' columnList ')' REFERENCES qualified_name opt_column_list key_match key_actions ConstraintAttributeSpec
		{
			rv := makeRangeVar($7)
			n := &nodes.Constraint{
				Contype:      nodes.CONSTR_FOREIGN,
				FkAttrs:      $4,
				Pktable:      rv.(*nodes.RangeVar),
				PkAttrs:      $8,
				FkMatchtype:  byte($9),
				FkUpdaction:  byte($10 >> 8),
				FkDelaction:  byte($10 & 0xFF),
				Location:     -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $11)
			$$ = n
		}
	| EXCLUDE access_method_clause '(' ExclusionConstraintList ')' opt_c_include opt_definition OptConsTableSpace where_clause ConstraintAttributeSpec
		{
			n := &nodes.Constraint{
				Contype:        nodes.CONSTR_EXCLUSION,
				AccessMethod:   $2,
				Exclusions:     $4,
				Including:      $6,
				Options:        $7,
				Indexspace:      $8,
				WhereClause:    $9,
				Location:       -1,
				InitiallyValid: true,
			}
			applyConstraintAttrs(n, $10)
			$$ = n
		}
	;

opt_column_list:
	'(' columnList ')' { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

columnList:
	name
		{ $$ = makeList(&nodes.String{Str: $1}) }
	| columnList ',' name
		{ $$ = appendList($1, &nodes.String{Str: $3}) }
	;

key_match:
	MATCH FULL     { $$ = int64('f') }
	| MATCH PARTIAL
		{
			/* PARTIAL is not implemented */
			$$ = int64('p')
		}
	| MATCH SIMPLE { $$ = int64('s') }
	| /* EMPTY */  { $$ = int64('s') }
	;

/*
 * key_actions packs update and delete actions into one int64:
 * high byte = update action, low byte = delete action
 */
key_actions:
	key_update                  { $$ = ($1 << 8) | int64('a') }
	| key_delete                { $$ = (int64('a') << 8) | $1 }
	| key_update key_delete     { $$ = ($1 << 8) | $2 }
	| key_delete key_update     { $$ = ($2 << 8) | $1 }
	| /* EMPTY */               { $$ = (int64('a') << 8) | int64('a') }
	;

key_update:
	ON UPDATE key_action { $$ = $3 }
	;

key_delete:
	ON DELETE_P key_action { $$ = $3 }
	;

key_action:
	NO ACTION                   { $$ = int64('a') }
	| RESTRICT                  { $$ = int64('r') }
	| CASCADE                   { $$ = int64('c') }
	| SET NULL_P opt_column_list    { $$ = int64('n') }
	| SET DEFAULT opt_column_list   { $$ = int64('d') }
	;

opt_c_include:
	INCLUDE '(' columnList ')'  { $$ = $3 }
	| /* EMPTY */               { $$ = nil }
	;

OptConsTableSpace:
	USING INDEX TABLESPACE name { $$ = $4 }
	| /* EMPTY */               { $$ = "" }
	;

ExistingIndex:
	USING INDEX name { $$ = $3 }
	;

ExclusionConstraintList:
	ExclusionConstraintElem
		{ $$ = makeList($1) }
	| ExclusionConstraintList ',' ExclusionConstraintElem
		{ $$ = appendList($1, $3) }
	;

ExclusionConstraintElem:
	index_elem WITH any_operator
		{
			$$ = &nodes.List{Items: []nodes.Node{$1, $3}}
		}
	| index_elem WITH OPERATOR '(' any_operator ')'
		{
			$$ = &nodes.List{Items: []nodes.Node{$1, $5}}
		}
	;

/*****************************************************************************
 *
 *      ALTER TABLE statement
 *
 *****************************************************************************/

AlterTableStmt:
	ALTER TABLE relation_expr alter_table_cmds
		{
			$$ = &nodes.AlterTableStmt{
				Relation:   $3.(*nodes.RangeVar),
				Cmds:       $4,
				ObjType:    int(nodes.OBJECT_TABLE),
			}
		}
	| ALTER TABLE IF_P EXISTS relation_expr alter_table_cmds
		{
			$$ = &nodes.AlterTableStmt{
				Relation:   $5.(*nodes.RangeVar),
				Cmds:       $6,
				ObjType:    int(nodes.OBJECT_TABLE),
				Missing_ok: true,
			}
		}
	| ALTER INDEX qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.AlterTableStmt{
				Relation: rv,
				Cmds:     $4,
				ObjType:  int(nodes.OBJECT_INDEX),
			}
		}
	| ALTER INDEX IF_P EXISTS qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.AlterTableStmt{
				Relation:   rv,
				Cmds:       $6,
				ObjType:    int(nodes.OBJECT_INDEX),
				Missing_ok: true,
			}
		}
	| ALTER INDEX qualified_name ATTACH PARTITION qualified_name
		{
			rvIdx := makeRangeVarFromAnyName($3)
			rvPart := makeRangeVar($6)
			cmd := &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AttachPartition),
				Def:     &nodes.PartitionCmd{
					Name: rvPart.(*nodes.RangeVar),
				},
			}
			$$ = &nodes.AlterTableStmt{
				Relation: rvIdx,
				Cmds:     makeList(cmd),
				ObjType:  int(nodes.OBJECT_INDEX),
			}
		}
	| ALTER SEQUENCE qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.AlterTableStmt{
				Relation: rv,
				Cmds:     $4,
				ObjType:  int(nodes.OBJECT_SEQUENCE),
			}
		}
	| ALTER SEQUENCE IF_P EXISTS qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.AlterTableStmt{
				Relation:   rv,
				Cmds:       $6,
				ObjType:    int(nodes.OBJECT_SEQUENCE),
				Missing_ok: true,
			}
		}
	| ALTER VIEW qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.AlterTableStmt{
				Relation: rv,
				Cmds:     $4,
				ObjType:  int(nodes.OBJECT_VIEW),
			}
		}
	| ALTER VIEW IF_P EXISTS qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.AlterTableStmt{
				Relation:   rv,
				Cmds:       $6,
				ObjType:    int(nodes.OBJECT_VIEW),
				Missing_ok: true,
			}
		}
	| ALTER MATERIALIZED VIEW qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($4)
			$$ = &nodes.AlterTableStmt{
				Relation: rv,
				Cmds:     $5,
				ObjType:  int(nodes.OBJECT_MATVIEW),
			}
		}
	| ALTER MATERIALIZED VIEW IF_P EXISTS qualified_name alter_table_cmds
		{
			rv := makeRangeVarFromAnyName($6)
			$$ = &nodes.AlterTableStmt{
				Relation:   rv,
				Cmds:       $7,
				ObjType:    int(nodes.OBJECT_MATVIEW),
				Missing_ok: true,
			}
		}
	| ALTER FOREIGN TABLE relation_expr alter_table_cmds
		{
			$$ = &nodes.AlterTableStmt{
				Relation: $4.(*nodes.RangeVar),
				Cmds:     $5,
				ObjType:  int(nodes.OBJECT_FOREIGN_TABLE),
			}
		}
	| ALTER FOREIGN TABLE IF_P EXISTS relation_expr alter_table_cmds
		{
			$$ = &nodes.AlterTableStmt{
				Relation:   $6.(*nodes.RangeVar),
				Cmds:       $7,
				ObjType:    int(nodes.OBJECT_FOREIGN_TABLE),
				Missing_ok: true,
			}
		}
	/* ALTER TABLE ALL IN TABLESPACE ... SET TABLESPACE ... */
	| ALTER TABLE ALL IN_P TABLESPACE name SET TABLESPACE name opt_nowait
		{
			$$ = &nodes.AlterTableMoveAllStmt{
				OrigTablespacename: $6,
				ObjType:            int(nodes.OBJECT_TABLE),
				NewTablespacename:  $9,
				Nowait:             $10,
			}
		}
	| ALTER TABLE ALL IN_P TABLESPACE name OWNED BY role_list SET TABLESPACE name opt_nowait
		{
			$$ = &nodes.AlterTableMoveAllStmt{
				OrigTablespacename: $6,
				ObjType:            int(nodes.OBJECT_TABLE),
				Roles:              $9,
				NewTablespacename:  $12,
				Nowait:             $13,
			}
		}
	| ALTER INDEX ALL IN_P TABLESPACE name SET TABLESPACE name opt_nowait
		{
			$$ = &nodes.AlterTableMoveAllStmt{
				OrigTablespacename: $6,
				ObjType:            int(nodes.OBJECT_INDEX),
				NewTablespacename:  $9,
				Nowait:             $10,
			}
		}
	| ALTER INDEX ALL IN_P TABLESPACE name OWNED BY role_list SET TABLESPACE name opt_nowait
		{
			$$ = &nodes.AlterTableMoveAllStmt{
				OrigTablespacename: $6,
				ObjType:            int(nodes.OBJECT_INDEX),
				Roles:              $9,
				NewTablespacename:  $12,
				Nowait:             $13,
			}
		}
	| ALTER MATERIALIZED VIEW ALL IN_P TABLESPACE name SET TABLESPACE name opt_nowait
		{
			$$ = &nodes.AlterTableMoveAllStmt{
				OrigTablespacename: $7,
				ObjType:            int(nodes.OBJECT_MATVIEW),
				NewTablespacename:  $10,
				Nowait:             $11,
			}
		}
	| ALTER MATERIALIZED VIEW ALL IN_P TABLESPACE name OWNED BY role_list SET TABLESPACE name opt_nowait
		{
			$$ = &nodes.AlterTableMoveAllStmt{
				OrigTablespacename: $7,
				ObjType:            int(nodes.OBJECT_MATVIEW),
				Roles:              $10,
				NewTablespacename:  $13,
				Nowait:             $14,
			}
		}
	;

alter_table_cmds:
	alter_table_cmd
		{ $$ = makeList($1) }
	| alter_table_cmds ',' alter_table_cmd
		{ $$ = appendList($1, $3) }
	;

alter_table_cmd:
	ADD_P columnDef
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddColumn),
				Def:     $2,
			}
		}
	| ADD_P COLUMN columnDef
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddColumn),
				Def:     $3,
			}
		}
	| ADD_P COLUMN IF_P NOT EXISTS columnDef
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_AddColumn),
				Def:        $6,
				Missing_ok: true,
			}
		}
	| DROP COLUMN ColId opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_DropColumn),
				Name:     $3,
				Behavior: int($4),
			}
		}
	| DROP COLUMN IF_P EXISTS ColId opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_DropColumn),
				Name:       $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| ALTER COLUMN ColId SET DEFAULT a_expr
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ColumnDefault),
				Name:    $3,
				Def:     $6,
			}
		}
	| ALTER COLUMN ColId DROP DEFAULT
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ColumnDefault),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId SET NOT NULL_P
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetNotNull),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId DROP NOT NULL_P
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropNotNull),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId TYPE_P Typename opt_collate_clause
		{
			coldef := &nodes.ColumnDef{TypeName: $5}
			if $6 != nil {
				coldef.CollClause = $6.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $3,
				Def:     coldef,
			}
		}
	| ADD_P TableConstraint
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddConstraint),
				Def:     $2,
			}
		}
	| DROP CONSTRAINT name opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_DropConstraint),
				Name:     $3,
				Behavior: int($4),
			}
		}
	| DROP CONSTRAINT IF_P EXISTS name opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_DropConstraint),
				Name:       $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| OWNER TO RoleSpec
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_ChangeOwner),
				Newowner: $3.(*nodes.RoleSpec),
			}
		}
	/* ALTER COLUMN ... TYPE ... USING */
	| ALTER COLUMN ColId SET DATA_P TYPE_P Typename opt_collate_clause
		{
			coldef := &nodes.ColumnDef{TypeName: $7}
			if $8 != nil {
				coldef.CollClause = $8.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $3,
				Def:     coldef,
			}
		}
	| ALTER COLUMN ColId TYPE_P Typename opt_collate_clause USING a_expr
		{
			coldef := &nodes.ColumnDef{TypeName: $5}
			if $6 != nil {
				coldef.CollClause = $6.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $3,
				Def:     coldef,
			}
		}
	| ALTER COLUMN ColId SET DATA_P TYPE_P Typename opt_collate_clause USING a_expr
		{
			coldef := &nodes.ColumnDef{TypeName: $7}
			if $8 != nil {
				coldef.CollClause = $8.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $3,
				Def:     coldef,
			}
		}
	| ALTER ColId TYPE_P Typename opt_collate_clause
		{
			coldef := &nodes.ColumnDef{TypeName: $4}
			if $5 != nil {
				coldef.CollClause = $5.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $2,
				Def:     coldef,
			}
		}
	| ALTER ColId SET DATA_P TYPE_P Typename opt_collate_clause
		{
			coldef := &nodes.ColumnDef{TypeName: $6}
			if $7 != nil {
				coldef.CollClause = $7.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $2,
				Def:     coldef,
			}
		}
	| ALTER ColId SET DATA_P TYPE_P Typename opt_collate_clause USING a_expr
		{
			coldef := &nodes.ColumnDef{TypeName: $6}
			if $7 != nil {
				coldef.CollClause = $7.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $2,
				Def:     coldef,
			}
		}
	| ALTER ColId SET DEFAULT a_expr
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ColumnDefault),
				Name:    $2,
				Def:     $5,
			}
		}
	| ALTER ColId DROP DEFAULT
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ColumnDefault),
				Name:    $2,
			}
		}
	| ALTER ColId SET NOT NULL_P
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetNotNull),
				Name:    $2,
			}
		}
	| ALTER ColId DROP NOT NULL_P
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropNotNull),
				Name:    $2,
			}
		}
	| ALTER ColId TYPE_P Typename opt_collate_clause USING a_expr
		{
			coldef := &nodes.ColumnDef{TypeName: $4}
			if $5 != nil {
				coldef.CollClause = $5.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnType),
				Name:    $2,
				Def:     coldef,
			}
		}
	/* ALTER FOREIGN TABLE <name> ALTER [COLUMN] <colname> OPTIONS */
	| ALTER COLUMN ColId alter_generic_options
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnGenericOptions),
				Name:    $3,
				Def:     $4,
			}
		}
	| ALTER ColId alter_generic_options
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterColumnGenericOptions),
				Name:    $2,
				Def:     $3,
			}
		}
	/* VALIDATE CONSTRAINT */
	| VALIDATE CONSTRAINT name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ValidateConstraint),
				Name:    $3,
			}
		}
	/* ALTER CONSTRAINT */
	| ALTER CONSTRAINT name ConstraintAttributeSpec
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AlterConstraint),
				Name:    $3,
			}
		}
	/* INHERIT / NO INHERIT */
	| INHERIT qualified_name
		{
			rv := makeRangeVar($2)
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddInherit),
				Def:     rv,
			}
		}
	| NO INHERIT qualified_name
		{
			rv := makeRangeVar($3)
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropInherit),
				Def:     rv,
			}
		}
	/* ATTACH / DETACH PARTITION */
	| ATTACH PARTITION qualified_name ForValues
		{
			rv := makeRangeVar($3)
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AttachPartition),
				Def:     &nodes.PartitionCmd{
					Name:  rv.(*nodes.RangeVar),
					Bound: $4,
				},
			}
		}
	| DETACH PARTITION qualified_name
		{
			rv := makeRangeVar($3)
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DetachPartition),
				Def:     &nodes.PartitionCmd{
					Name: rv.(*nodes.RangeVar),
				},
			}
		}
	| DETACH PARTITION qualified_name CONCURRENTLY
		{
			rv := makeRangeVar($3)
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DetachPartition),
				Def:     &nodes.PartitionCmd{
					Name:       rv.(*nodes.RangeVar),
					Concurrent: true,
				},
			}
		}
	| DETACH PARTITION qualified_name FINALIZE
		{
			rv := makeRangeVar($3)
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DetachPartitionFinalize),
				Def:     &nodes.PartitionCmd{
					Name: rv.(*nodes.RangeVar),
				},
			}
		}
	/* ENABLE/DISABLE TRIGGER */
	| ENABLE_P TRIGGER name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableTrig),
				Name:    $3,
			}
		}
	| ENABLE_P ALWAYS TRIGGER name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableAlwaysTrig),
				Name:    $4,
			}
		}
	| ENABLE_P REPLICA TRIGGER name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableReplicaTrig),
				Name:    $4,
			}
		}
	| DISABLE_P TRIGGER name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DisableTrig),
				Name:    $3,
			}
		}
	| ENABLE_P TRIGGER ALL
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableTrigAll),
			}
		}
	| DISABLE_P TRIGGER ALL
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DisableTrigAll),
			}
		}
	| ENABLE_P TRIGGER USER
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableTrigUser),
			}
		}
	| DISABLE_P TRIGGER USER
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DisableTrigUser),
			}
		}
	/* ENABLE/DISABLE RULE */
	| ENABLE_P RULE name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableRule),
				Name:    $3,
			}
		}
	| ENABLE_P ALWAYS RULE name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableAlwaysRule),
				Name:    $4,
			}
		}
	| ENABLE_P REPLICA RULE name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableReplicaRule),
				Name:    $4,
			}
		}
	| DISABLE_P RULE name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DisableRule),
				Name:    $3,
			}
		}
	/* ENABLE/DISABLE/FORCE ROW LEVEL SECURITY */
	| ENABLE_P ROW LEVEL SECURITY
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_EnableRowSecurity),
			}
		}
	| DISABLE_P ROW LEVEL SECURITY
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DisableRowSecurity),
			}
		}
	| FORCE ROW LEVEL SECURITY
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ForceRowSecurity),
			}
		}
	| NO FORCE ROW LEVEL SECURITY
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_NoForceRowSecurity),
			}
		}
	/* CLUSTER ON / SET WITHOUT CLUSTER */
	| CLUSTER ON name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ClusterOn),
				Name:    $3,
			}
		}
	| SET WITHOUT CLUSTER
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropCluster),
			}
		}
	/* SET LOGGED / SET UNLOGGED */
	| SET LOGGED
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetLogged),
			}
		}
	| SET UNLOGGED
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetUnLogged),
			}
		}
	/* SET ACCESS METHOD */
	| SET ACCESS METHOD name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetAccessMethod),
				Name:    $4,
			}
		}
	| SET ACCESS METHOD DEFAULT
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetAccessMethod),
				Name:    "",
			}
		}
	/* SET/RESET storage params */
	| SET reloptions
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetRelOptions),
				Def:     $2,
			}
		}
	| RESET reloptions
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ResetRelOptions),
				Def:     $2,
			}
		}
	/* ADD COLUMN IF NOT EXISTS (without COLUMN keyword) */
	| ADD_P IF_P NOT EXISTS columnDef
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_AddColumn),
				Def:        $5,
				Missing_ok: true,
			}
		}
	/* DROP (without COLUMN keyword) */
	| DROP ColId opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_DropColumn),
				Name:     $2,
				Behavior: int($3),
			}
		}
	| DROP IF_P EXISTS ColId opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_DropColumn),
				Name:       $4,
				Behavior:   int($5),
				Missing_ok: true,
			}
		}
	/* REPLICA IDENTITY */
	| REPLICA IDENTITY_P DEFAULT
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ReplicaIdentity),
			}
		}
	| REPLICA IDENTITY_P FULL
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ReplicaIdentity),
			}
		}
	| REPLICA IDENTITY_P NOTHING
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ReplicaIdentity),
			}
		}
	| REPLICA IDENTITY_P USING INDEX name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ReplicaIdentity),
				Name:    $5,
			}
		}
	/* SET STORAGE */
	| ALTER COLUMN ColId SET STORAGE ColId
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStorage),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId SET STORAGE DEFAULT
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStorage),
				Name:    $3,
			}
		}
	| ALTER ColId SET STORAGE ColId
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStorage),
				Name:    $2,
			}
		}
	| ALTER ColId SET STORAGE DEFAULT
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStorage),
				Name:    $2,
			}
		}
	/* SET STATISTICS */
	| ALTER COLUMN ColId SET STATISTICS SignedIconst
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStatistics),
				Name:    $3,
				Def:     makeIntConst($6),
			}
		}
	| ALTER ColId SET STATISTICS SignedIconst
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStatistics),
				Name:    $2,
				Def:     makeIntConst($5),
			}
		}
	| ALTER COLUMN Iconst SET STATISTICS SignedIconst
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStatistics),
				Num:     int16($3),
				Def:     makeIntConst($6),
			}
		}
	| ALTER Iconst SET STATISTICS SignedIconst
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetStatistics),
				Num:     int16($2),
				Def:     makeIntConst($5),
			}
		}
	/* SET COMPRESSION */
	| ALTER COLUMN ColId SET column_compression
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetCompression),
				Name:    $3,
				Def:     &nodes.String{Str: $5},
			}
		}
	| ALTER ColId SET column_compression
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetCompression),
				Name:    $2,
				Def:     &nodes.String{Str: $4},
			}
		}
	/* SET/DROP EXPRESSION */
	| ALTER COLUMN ColId SET EXPRESSION AS '(' a_expr ')'
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetExpression),
				Name:    $3,
				Def:     $8,
			}
		}
	| ALTER ColId SET EXPRESSION AS '(' a_expr ')'
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetExpression),
				Name:    $2,
				Def:     $7,
			}
		}
	| ALTER COLUMN ColId DROP EXPRESSION
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropExpression),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId DROP EXPRESSION IF_P EXISTS
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_DropExpression),
				Name:       $3,
				Missing_ok: true,
			}
		}
	/* ADD GENERATED IDENTITY */
	| ALTER COLUMN ColId ADD_P GENERATED ALWAYS AS IDENTITY_P
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddIdentity),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId ADD_P GENERATED ALWAYS AS IDENTITY_P '(' OptSeqOptList ')'
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddIdentity),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId ADD_P GENERATED BY DEFAULT AS IDENTITY_P
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddIdentity),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId ADD_P GENERATED BY DEFAULT AS IDENTITY_P '(' OptSeqOptList ')'
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddIdentity),
				Name:    $3,
			}
		}
	/* DROP/SET IDENTITY */
	| ALTER COLUMN ColId DROP IDENTITY_P
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropIdentity),
				Name:    $3,
			}
		}
	| ALTER COLUMN ColId DROP IDENTITY_P IF_P EXISTS
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_DropIdentity),
				Name:       $3,
				Missing_ok: true,
			}
		}
	/* ADD GENERATED identity */
	| ALTER COLUMN ColId ADD_P GENERATED generated_when AS IDENTITY_P OptParenthesizedSeqOptList
		{
			c := &nodes.Constraint{
				Contype:       nodes.CONSTR_IDENTITY,
				GeneratedWhen: byte($6),
				Options:       $9,
				Location:      -1,
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddIdentity),
				Name:    $3,
				Def:     c,
			}
		}
	| ALTER ColId ADD_P GENERATED generated_when AS IDENTITY_P OptParenthesizedSeqOptList
		{
			c := &nodes.Constraint{
				Contype:       nodes.CONSTR_IDENTITY,
				GeneratedWhen: byte($5),
				Options:       $8,
				Location:      -1,
			}
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddIdentity),
				Name:    $2,
				Def:     c,
			}
		}
	/* SET GENERATED / identity options */
	| ALTER COLUMN ColId alter_identity_column_option_list
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetIdentity),
				Name:    $3,
				Def:     $4,
			}
		}
	/* ALTER COLUMN ... SET (options) / RESET (options) */
	| ALTER COLUMN ColId SET '(' def_list ')'
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetOptions),
				Name:    $3,
				Def:     $6,
			}
		}
	| ALTER COLUMN ColId RESET '(' def_list ')'
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_ResetOptions),
				Name:    $3,
				Def:     $6,
			}
		}
	/* SET WITH/WITHOUT OIDS (deprecated but still parsed) */
	| SET WITH OIDS
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropOids), /* deprecated, reuse DropOids */
			}
		}
	| SET WITHOUT OIDS
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropOids),
			}
		}
	/* SET TABLESPACE */
	| SET TABLESPACE name
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_SetTableSpace),
				Name:    $3,
			}
		}
	/* OPTIONS for foreign tables */
	| alter_generic_options
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_GenericOptions),
				Def:     $1,
			}
		}
	/* OF typename */
	| OF any_name
		{
			tn := makeTypeNameFromNameList($2)
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_AddOf),
				Def:     tn,
			}
		}
	/* NOT OF */
	| NOT OF
		{
			$$ = &nodes.AlterTableCmd{
				Subtype: int(nodes.AT_DropOf),
			}
		}
	;

opt_drop_behavior:
	CASCADE     { $$ = int64(nodes.DROP_CASCADE) }
	| RESTRICT  { $$ = int64(nodes.DROP_RESTRICT) }
	| /* EMPTY */ { $$ = int64(nodes.DROP_RESTRICT) }
	;

column_compression:
	COMPRESSION ColId      { $$ = $2 }
	| COMPRESSION DEFAULT  { $$ = "default" }
	;

/*****************************************************************************
 *
 *      ALTER TABLE ... RENAME
 *
 *****************************************************************************/

RenameStmt:
	ALTER TABLE relation_expr RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TABLE,
				Relation:   $3.(*nodes.RangeVar),
				Newname:    $6,
			}
		}
	| ALTER TABLE IF_P EXISTS relation_expr RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TABLE,
				Relation:   $5.(*nodes.RangeVar),
				Newname:    $8,
				MissingOk:  true,
			}
		}
	| ALTER TABLE relation_expr RENAME COLUMN ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_TABLE,
				Relation:     $3.(*nodes.RangeVar),
				Subname:      $6,
				Newname:      $8,
			}
		}
	| ALTER TABLE relation_expr RENAME ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_TABLE,
				Relation:     $3.(*nodes.RangeVar),
				Subname:      $5,
				Newname:      $7,
			}
		}
	| ALTER TABLE relation_expr RENAME CONSTRAINT name TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_TABCONSTRAINT,
				RelationType: nodes.OBJECT_TABLE,
				Relation:     $3.(*nodes.RangeVar),
				Subname:      $6,
				Newname:      $8,
			}
		}
	| ALTER TABLE IF_P EXISTS relation_expr RENAME COLUMN ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_TABLE,
				Relation:     $5.(*nodes.RangeVar),
				Subname:      $8,
				Newname:      $10,
				MissingOk:    true,
			}
		}
	| ALTER TABLE IF_P EXISTS relation_expr RENAME ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_TABLE,
				Relation:     $5.(*nodes.RangeVar),
				Subname:      $7,
				Newname:      $9,
				MissingOk:    true,
			}
		}
	| ALTER TABLE IF_P EXISTS relation_expr RENAME CONSTRAINT name TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_TABCONSTRAINT,
				RelationType: nodes.OBJECT_TABLE,
				Relation:     $5.(*nodes.RangeVar),
				Subname:      $8,
				Newname:      $10,
				MissingOk:    true,
			}
		}
	/* RENAME INDEX */
	| ALTER INDEX qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_INDEX,
				Relation:   rv,
				Newname:    $6,
			}
		}
	| ALTER INDEX IF_P EXISTS qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_INDEX,
				Relation:   rv,
				Newname:    $8,
				MissingOk:  true,
			}
		}
	/* RENAME SEQUENCE */
	| ALTER SEQUENCE qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_SEQUENCE,
				Relation:   rv,
				Newname:    $6,
			}
		}
	| ALTER SEQUENCE IF_P EXISTS qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_SEQUENCE,
				Relation:   rv,
				Newname:    $8,
				MissingOk:  true,
			}
		}
	/* RENAME VIEW */
	| ALTER VIEW qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_VIEW,
				Relation:   rv,
				Newname:    $6,
			}
		}
	| ALTER VIEW IF_P EXISTS qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_VIEW,
				Relation:   rv,
				Newname:    $8,
				MissingOk:  true,
			}
		}
	| ALTER VIEW qualified_name RENAME COLUMN ColId TO name
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_VIEW,
				Relation:     rv,
				Subname:      $6,
				Newname:      $8,
			}
		}
	| ALTER VIEW qualified_name RENAME ColId TO name
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_VIEW,
				Relation:     rv,
				Subname:      $5,
				Newname:      $7,
			}
		}
	/* RENAME MATERIALIZED VIEW */
	| ALTER MATERIALIZED VIEW qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($4)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_MATVIEW,
				Relation:   rv,
				Newname:    $7,
			}
		}
	| ALTER MATERIALIZED VIEW IF_P EXISTS qualified_name RENAME TO name
		{
			rv := makeRangeVarFromAnyName($6)
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_MATVIEW,
				Relation:   rv,
				Newname:    $9,
				MissingOk:  true,
			}
		}
	| ALTER MATERIALIZED VIEW qualified_name RENAME COLUMN ColId TO name
		{
			rv := makeRangeVarFromAnyName($4)
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_MATVIEW,
				Relation:     rv,
				Subname:      $7,
				Newname:      $9,
			}
		}
	| ALTER MATERIALIZED VIEW qualified_name RENAME ColId TO name
		{
			rv := makeRangeVarFromAnyName($4)
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_MATVIEW,
				Relation:     rv,
				Subname:      $6,
				Newname:      $8,
			}
		}
	/* RENAME FUNCTION / PROCEDURE / ROUTINE */
	| ALTER FUNCTION function_with_argtypes RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_FUNCTION,
				Object:     $3,
				Newname:    $6,
			}
		}
	| ALTER PROCEDURE function_with_argtypes RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_PROCEDURE,
				Object:     $3,
				Newname:    $6,
			}
		}
	| ALTER ROUTINE function_with_argtypes RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_ROUTINE,
				Object:     $3,
				Newname:    $6,
			}
		}
	/* RENAME AGGREGATE */
	| ALTER AGGREGATE aggregate_with_argtypes RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_AGGREGATE,
				Object:     $3,
				Newname:    $6,
			}
		}
	/* RENAME COLLATION */
	| ALTER COLLATION any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_COLLATION,
				Object:     $3,
				Newname:    $6,
			}
		}
	/* RENAME CONVERSION */
	| ALTER CONVERSION_P any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_CONVERSION,
				Object:     $3,
				Newname:    $6,
			}
		}
	/* RENAME DOMAIN */
	| ALTER DOMAIN_P any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_DOMAIN,
				Object:     $3,
				Newname:    $6,
			}
		}
	| ALTER DOMAIN_P any_name RENAME CONSTRAINT name TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_DOMCONSTRAINT,
				Object:       $3,
				Subname:      $6,
				Newname:      $8,
			}
		}
	/* RENAME SCHEMA */
	| ALTER SCHEMA name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_SCHEMA,
				Subname:    $3,
				Newname:    $6,
			}
		}
	/* RENAME SERVER */
	| ALTER SERVER name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_FOREIGN_SERVER,
				Object:     &nodes.String{Str: $3},
				Newname:    $6,
			}
		}
	/* RENAME FOREIGN DATA WRAPPER */
	| ALTER FOREIGN DATA_P WRAPPER name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_FDW,
				Object:     &nodes.String{Str: $5},
				Newname:    $8,
			}
		}
	/* RENAME TYPE */
	| ALTER TYPE_P any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TYPE,
				Object:     $3,
				Newname:    $6,
			}
		}
	| ALTER TYPE_P any_name RENAME ATTRIBUTE name TO name opt_drop_behavior
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_ATTRIBUTE,
				RelationType: nodes.OBJECT_TYPE,
				Relation:     makeRangeVarFromAnyName($3),
				Subname:      $6,
				Newname:      $8,
				Behavior:     nodes.DropBehavior($9),
			}
		}
	/* RENAME FOREIGN TABLE */
	| ALTER FOREIGN TABLE relation_expr RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:   $4.(*nodes.RangeVar),
				Newname:    $7,
			}
		}
	| ALTER FOREIGN TABLE IF_P EXISTS relation_expr RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:   $6.(*nodes.RangeVar),
				Newname:    $9,
				MissingOk:  true,
			}
		}
	| ALTER FOREIGN TABLE relation_expr RENAME COLUMN ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:     $4.(*nodes.RangeVar),
				Subname:      $7,
				Newname:      $9,
			}
		}
	| ALTER FOREIGN TABLE relation_expr RENAME ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:     $4.(*nodes.RangeVar),
				Subname:      $6,
				Newname:      $8,
			}
		}
	| ALTER FOREIGN TABLE IF_P EXISTS relation_expr RENAME COLUMN ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:     $6.(*nodes.RangeVar),
				Subname:      $9,
				Newname:      $11,
				MissingOk:    true,
			}
		}
	| ALTER FOREIGN TABLE IF_P EXISTS relation_expr RENAME ColId TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType:   nodes.OBJECT_COLUMN,
				RelationType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:     $6.(*nodes.RangeVar),
				Subname:      $8,
				Newname:      $10,
				MissingOk:    true,
			}
		}
	/* RENAME OPERATOR CLASS / FAMILY */
	| ALTER OPERATOR CLASS any_name USING name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_OPCLASS,
				Object:     prependList(&nodes.String{Str: $6}, $4),
				Newname:    $9,
			}
		}
	| ALTER OPERATOR FAMILY any_name USING name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_OPFAMILY,
				Object:     prependList(&nodes.String{Str: $6}, $4),
				Newname:    $9,
			}
		}
	/* RENAME TEXT SEARCH */
	| ALTER TEXT_P SEARCH PARSER any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TSPARSER,
				Object:     $5,
				Newname:    $8,
			}
		}
	| ALTER TEXT_P SEARCH DICTIONARY any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TSDICTIONARY,
				Object:     $5,
				Newname:    $8,
			}
		}
	| ALTER TEXT_P SEARCH TEMPLATE any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TSTEMPLATE,
				Object:     $5,
				Newname:    $8,
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TSCONFIGURATION,
				Object:     $5,
				Newname:    $8,
			}
		}
	/* RENAME PUBLICATION / SUBSCRIPTION */
	| ALTER PUBLICATION name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_PUBLICATION,
				Object:     &nodes.String{Str: $3},
				Newname:    $6,
			}
		}
	| ALTER SUBSCRIPTION name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_SUBSCRIPTION,
				Object:     &nodes.String{Str: $3},
				Newname:    $6,
			}
		}
	/* RENAME RULE */
	| ALTER RULE name ON qualified_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_RULE,
				Relation:   makeRangeVarFromAnyName($5),
				Subname:    $3,
				Newname:    $8,
			}
		}
	/* RENAME TRIGGER */
	| ALTER TRIGGER name ON qualified_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TRIGGER,
				Relation:   makeRangeVarFromAnyName($5),
				Subname:    $3,
				Newname:    $8,
			}
		}
	/* RENAME EVENT TRIGGER */
	| ALTER EVENT TRIGGER name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_EVENT_TRIGGER,
				Object:     &nodes.String{Str: $4},
				Newname:    $7,
			}
		}
	/* RENAME STATISTICS */
	| ALTER STATISTICS any_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_STATISTIC_EXT,
				Object:     $3,
				Newname:    $6,
			}
		}
	/* RENAME POLICY */
	| ALTER POLICY name ON qualified_name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_POLICY,
				Relation:   makeRangeVarFromAnyName($5),
				Subname:    $3,
				Newname:    $8,
			}
		}
	/* RENAME LANGUAGE / DATABASE / TABLESPACE / ROLE */
	| ALTER opt_procedural LANGUAGE name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_LANGUAGE,
				Object:     &nodes.String{Str: $4},
				Newname:    $7,
			}
		}
	| ALTER DATABASE name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_DATABASE,
				Subname:    $3,
				Newname:    $6,
			}
		}
	| ALTER TABLESPACE name RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_TABLESPACE,
				Subname:    $3,
				Newname:    $6,
			}
		}
	| ALTER ROLE RoleSpec RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_ROLE,
				Subname:    $3.(*nodes.RoleSpec).Rolename,
				Newname:    $6,
			}
		}
	| ALTER USER RoleSpec RENAME TO name
		{
			$$ = &nodes.RenameStmt{
				RenameType: nodes.OBJECT_ROLE,
				Subname:    $3.(*nodes.RoleSpec).Rolename,
				Newname:    $6,
			}
		}
	;

/*****************************************************************************
 *
 *      DROP statement
 *
 *****************************************************************************/

DropStmt:
	DROP object_type_any_name IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $5,
				RemoveType: int($2),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP object_type_any_name any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $3,
				RemoveType: int($2),
				Behavior:   int($4),
			}
		}
	/* DROP TYPE / DROP DOMAIN */
	| DROP TYPE_P any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $3,
				RemoveType: int(nodes.OBJECT_TYPE),
				Behavior:   int($4),
			}
		}
	| DROP TYPE_P IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $5,
				RemoveType: int(nodes.OBJECT_TYPE),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP DOMAIN_P any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $3,
				RemoveType: int(nodes.OBJECT_DOMAIN),
				Behavior:   int($4),
			}
		}
	| DROP DOMAIN_P IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $5,
				RemoveType: int(nodes.OBJECT_DOMAIN),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	/* DROP SCHEMA / EXTENSION / LANGUAGE / PUBLICATION / SUBSCRIPTION */
	| DROP SCHEMA name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($3),
				RemoveType: int(nodes.OBJECT_SCHEMA),
				Behavior:   int($4),
			}
		}
	| DROP SCHEMA IF_P EXISTS name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($5),
				RemoveType: int(nodes.OBJECT_SCHEMA),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP EXTENSION name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($3),
				RemoveType: int(nodes.OBJECT_EXTENSION),
				Behavior:   int($4),
			}
		}
	| DROP EXTENSION IF_P EXISTS name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($5),
				RemoveType: int(nodes.OBJECT_EXTENSION),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	/* DROP FUNCTION / PROCEDURE / ROUTINE / AGGREGATE */
	| DROP FUNCTION any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $3,
				RemoveType: int(nodes.OBJECT_FUNCTION),
				Behavior:   int($4),
			}
		}
	| DROP FUNCTION IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $5,
				RemoveType: int(nodes.OBJECT_FUNCTION),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP PROCEDURE any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $3,
				RemoveType: int(nodes.OBJECT_PROCEDURE),
				Behavior:   int($4),
			}
		}
	| DROP PROCEDURE IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $5,
				RemoveType: int(nodes.OBJECT_PROCEDURE),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP ROUTINE any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $3,
				RemoveType: int(nodes.OBJECT_ROUTINE),
				Behavior:   int($4),
			}
		}
	| DROP ROUTINE IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $5,
				RemoveType: int(nodes.OBJECT_ROUTINE),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP AGGREGATE any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $3,
				RemoveType: int(nodes.OBJECT_AGGREGATE),
				Behavior:   int($4),
			}
		}
	| DROP AGGREGATE IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $5,
				RemoveType: int(nodes.OBJECT_AGGREGATE),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	/* DROP TRIGGER / POLICY / RULE ON relation */
	| DROP TRIGGER name ON any_name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    &nodes.List{Items: []nodes.Node{$5}},
				RemoveType: int(nodes.OBJECT_TRIGGER),
				Behavior:   int($6),
			}
		}
	| DROP TRIGGER IF_P EXISTS name ON any_name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    &nodes.List{Items: []nodes.Node{$7}},
				RemoveType: int(nodes.OBJECT_TRIGGER),
				Behavior:   int($8),
				Missing_ok: true,
			}
		}
	| DROP POLICY name ON any_name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    &nodes.List{Items: []nodes.Node{$5}},
				RemoveType: int(nodes.OBJECT_POLICY),
				Behavior:   int($6),
			}
		}
	| DROP POLICY IF_P EXISTS name ON any_name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    &nodes.List{Items: []nodes.Node{$7}},
				RemoveType: int(nodes.OBJECT_POLICY),
				Behavior:   int($8),
				Missing_ok: true,
			}
		}
	| DROP RULE name ON any_name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    &nodes.List{Items: []nodes.Node{$5}},
				RemoveType: int(nodes.OBJECT_RULE),
				Behavior:   int($6),
			}
		}
	| DROP RULE IF_P EXISTS name ON any_name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    &nodes.List{Items: []nodes.Node{$7}},
				RemoveType: int(nodes.OBJECT_RULE),
				Behavior:   int($8),
				Missing_ok: true,
			}
		}
	/* DROP EVENT TRIGGER */
	| DROP EVENT TRIGGER name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList(&nodes.List{Items: []nodes.Node{&nodes.String{Str: $4}}}),
				RemoveType: int(nodes.OBJECT_EVENT_TRIGGER),
				Behavior:   int($5),
			}
		}
	| DROP EVENT TRIGGER IF_P EXISTS name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList(&nodes.List{Items: []nodes.Node{&nodes.String{Str: $6}}}),
				RemoveType: int(nodes.OBJECT_EVENT_TRIGGER),
				Behavior:   int($7),
				Missing_ok: true,
			}
		}
	/* DROP LANGUAGE / PUBLICATION / SUBSCRIPTION */
	| DROP LANGUAGE name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList(&nodes.List{Items: []nodes.Node{&nodes.String{Str: $3}}}),
				RemoveType: int(nodes.OBJECT_LANGUAGE),
				Behavior:   int($4),
			}
		}
	| DROP LANGUAGE IF_P EXISTS name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList(&nodes.List{Items: []nodes.Node{&nodes.String{Str: $5}}}),
				RemoveType: int(nodes.OBJECT_LANGUAGE),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP PUBLICATION name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($3),
				RemoveType: int(nodes.OBJECT_PUBLICATION),
				Behavior:   int($4),
			}
		}
	| DROP PUBLICATION IF_P EXISTS name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($5),
				RemoveType: int(nodes.OBJECT_PUBLICATION),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP SUBSCRIPTION name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($3),
				RemoveType: int(nodes.OBJECT_SUBSCRIPTION),
				Behavior:   int($4),
			}
		}
	| DROP SUBSCRIPTION IF_P EXISTS name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($5),
				RemoveType: int(nodes.OBJECT_SUBSCRIPTION),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	/* DROP INDEX CONCURRENTLY */
	| DROP INDEX CONCURRENTLY any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $4,
				RemoveType: int(nodes.OBJECT_INDEX),
				Behavior:   int($5),
				Concurrent: true,
			}
		}
	| DROP INDEX CONCURRENTLY IF_P EXISTS any_name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    $6,
				RemoveType: int(nodes.OBJECT_INDEX),
				Behavior:   int($7),
				Missing_ok: true,
				Concurrent: true,
			}
		}
	/* DROP FOREIGN DATA WRAPPER / DROP SERVER */
	| DROP FOREIGN DATA_P WRAPPER name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($5),
				RemoveType: int(nodes.OBJECT_FDW),
				Behavior:   int($6),
			}
		}
	| DROP FOREIGN DATA_P WRAPPER IF_P EXISTS name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($7),
				RemoveType: int(nodes.OBJECT_FDW),
				Behavior:   int($8),
				Missing_ok: true,
			}
		}
	| DROP SERVER name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($3),
				RemoveType: int(nodes.OBJECT_FOREIGN_SERVER),
				Behavior:   int($4),
			}
		}
	| DROP SERVER IF_P EXISTS name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($5),
				RemoveType: int(nodes.OBJECT_FOREIGN_SERVER),
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	/* DROP OWNED BY */
	| DROP OWNED BY name_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeNameListAsAnyNameList($4),
				RemoveType: int(nodes.OBJECT_ROLE),
				Behavior:   int($5),
			}
		}
	;

object_type_any_name:
	TABLE               { $$ = int64(nodes.OBJECT_TABLE) }
	| SEQUENCE          { $$ = int64(nodes.OBJECT_SEQUENCE) }
	| VIEW              { $$ = int64(nodes.OBJECT_VIEW) }
	| MATERIALIZED VIEW { $$ = int64(nodes.OBJECT_MATVIEW) }
	| INDEX             { $$ = int64(nodes.OBJECT_INDEX) }
	| FOREIGN TABLE     { $$ = int64(nodes.OBJECT_FOREIGN_TABLE) }
	| COLLATION         { $$ = int64(nodes.OBJECT_COLLATION) }
	| CONVERSION_P      { $$ = int64(nodes.OBJECT_CONVERSION) }
	| STATISTICS        { $$ = int64(nodes.OBJECT_STATISTIC_EXT) }
	| TEXT_P SEARCH PARSER { $$ = int64(nodes.OBJECT_TSPARSER) }
	| TEXT_P SEARCH DICTIONARY { $$ = int64(nodes.OBJECT_TSDICTIONARY) }
	| TEXT_P SEARCH TEMPLATE { $$ = int64(nodes.OBJECT_TSTEMPLATE) }
	| TEXT_P SEARCH CONFIGURATION { $$ = int64(nodes.OBJECT_TSCONFIGURATION) }
	| ACCESS METHOD     { $$ = int64(nodes.OBJECT_ACCESS_METHOD) }
	;

any_name_list:
	any_name
		{
			$$ = &nodes.List{Items: []nodes.Node{$1}}
		}
	| any_name_list ',' any_name
		{
			$1.Items = append($1.Items, $3)
			$$ = $1
		}
	;

/*****************************************************************************
 *
 *      CREATE INDEX statement
 *
 *****************************************************************************/

IndexStmt:
	CREATE opt_unique INDEX opt_concurrently opt_single_name
	ON relation_expr access_method_clause '(' index_params ')' opt_unique_null_treatment opt_include opt_reloptions OptTableSpace where_clause
		{
			rv := $7.(*nodes.RangeVar)
			$$ = &nodes.IndexStmt{
				Idxname:              $5,
				Relation:             rv,
				AccessMethod:         $8,
				IndexParams:          $10,
				Nulls_not_distinct:   $12,
				IndexIncludingParams: $13,
				Options:              $14,
				TableSpace:           $15,
				WhereClause:          $16,
				Unique:               $2,
				Concurrent:           $4,
			}
		}
	| CREATE opt_unique INDEX opt_concurrently IF_P NOT EXISTS name
	ON relation_expr access_method_clause '(' index_params ')' opt_unique_null_treatment opt_include opt_reloptions OptTableSpace where_clause
		{
			rv := $10.(*nodes.RangeVar)
			$$ = &nodes.IndexStmt{
				Idxname:              $8,
				Relation:             rv,
				AccessMethod:         $11,
				IndexParams:          $13,
				Nulls_not_distinct:   $15,
				IndexIncludingParams: $16,
				Options:              $17,
				TableSpace:           $18,
				WhereClause:          $19,
				Unique:               $2,
				Concurrent:           $4,
				IfNotExists:          true,
			}
		}
	;

opt_include:
	INCLUDE '(' index_including_params ')'
		{ $$ = $3 }
	| /* EMPTY */
		{ $$ = nil }
	;

index_including_params:
	index_elem
		{ $$ = makeList($1) }
	| index_including_params ',' index_elem
		{ $$ = appendList($1, $3) }
	;

opt_unique:
	UNIQUE      { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

opt_concurrently:
	CONCURRENTLY { $$ = true }
	| /* EMPTY */  { $$ = false }
	;

opt_single_name:
	name        { $$ = $1 }
	| /* EMPTY */ { $$ = "" }
	;

access_method_clause:
	USING name  { $$ = $2 }
	| /* EMPTY */ { $$ = "" }
	;

index_params:
	index_elem
		{ $$ = makeList($1) }
	| index_params ',' index_elem
		{ $$ = appendList($1, $3) }
	;

index_elem:
	ColId opt_collate opt_qualified_name opt_asc_desc opt_nulls_order
		{
			$$ = &nodes.IndexElem{
				Name:          $1,
				Collation:     $2,
				Opclass:       $3,
				Ordering:      nodes.SortByDir($4),
				NullsOrdering: nodes.SortByNulls($5),
			}
		}
	| ColId opt_collate any_name reloptions opt_asc_desc opt_nulls_order
		{
			$$ = &nodes.IndexElem{
				Name:          $1,
				Collation:     $2,
				Opclass:       $3,
				Opclassopts:   $4,
				Ordering:      nodes.SortByDir($5),
				NullsOrdering: nodes.SortByNulls($6),
			}
		}
	| func_expr_windowless opt_collate opt_qualified_name opt_asc_desc opt_nulls_order
		{
			$$ = &nodes.IndexElem{
				Expr:          $1,
				Collation:     $2,
				Opclass:       $3,
				Ordering:      nodes.SortByDir($4),
				NullsOrdering: nodes.SortByNulls($5),
			}
		}
	| func_expr_windowless opt_collate any_name reloptions opt_asc_desc opt_nulls_order
		{
			$$ = &nodes.IndexElem{
				Expr:          $1,
				Collation:     $2,
				Opclass:       $3,
				Opclassopts:   $4,
				Ordering:      nodes.SortByDir($5),
				NullsOrdering: nodes.SortByNulls($6),
			}
		}
	| '(' a_expr ')' opt_collate opt_qualified_name opt_asc_desc opt_nulls_order
		{
			$$ = &nodes.IndexElem{
				Expr:          $2,
				Collation:     $4,
				Opclass:       $5,
				Ordering:      nodes.SortByDir($6),
				NullsOrdering: nodes.SortByNulls($7),
			}
		}
	| '(' a_expr ')' opt_collate any_name reloptions opt_asc_desc opt_nulls_order
		{
			$$ = &nodes.IndexElem{
				Expr:          $2,
				Collation:     $4,
				Opclass:       $5,
				Opclassopts:   $6,
				Ordering:      nodes.SortByDir($7),
				NullsOrdering: nodes.SortByNulls($8),
			}
		}
	;

opt_nulls_order:
	NULLS_LA FIRST_P  { $$ = int64(nodes.SORTBY_NULLS_FIRST) }
	| NULLS_LA LAST_P { $$ = int64(nodes.SORTBY_NULLS_LAST) }
	| /* EMPTY */      { $$ = int64(nodes.SORTBY_NULLS_DEFAULT) }
	;

/*****************************************************************************
 *
 *      CREATE VIEW statement
 *
 *****************************************************************************/

ViewStmt:
	CREATE OptTemp VIEW qualified_name opt_column_list opt_reloptions
	AS SelectStmt opt_check_option
		{
			rv := makeRangeVar($4).(*nodes.RangeVar)
			$$ = &nodes.ViewStmt{
				View:            rv,
				Aliases:         $5,
				Options:         $6,
				Query:           $8,
				WithCheckOption: int($9),
			}
		}
	| CREATE OR REPLACE OptTemp VIEW qualified_name opt_column_list opt_reloptions
	AS SelectStmt opt_check_option
		{
			rv := makeRangeVar($6).(*nodes.RangeVar)
			$$ = &nodes.ViewStmt{
				View:            rv,
				Aliases:         $7,
				Options:         $8,
				Query:           $10,
				Replace:         true,
				WithCheckOption: int($11),
			}
		}
	| CREATE OptTemp RECURSIVE VIEW qualified_name '(' columnList ')' opt_reloptions
	AS SelectStmt opt_check_option
		{
			rv := makeRangeVar($5).(*nodes.RangeVar)
			/* Create a WITH RECURSIVE CTE wrapper */
			cte := &nodes.CommonTableExpr{
				Ctename:      rv.Relname,
				Ctequery:     $11,
				Aliascolnames:  $7,
			}
			wc := &nodes.WithClause{
				Ctes:      makeList(cte),
				Recursive: true,
			}
			cr := &nodes.ColumnRef{
				Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
				Location: -1,
			}
			rt := &nodes.ResTarget{Val: cr, Location: -1}
			fromRv := &nodes.RangeVar{Relname: rv.Relname, Inh: true, Relpersistence: 'p'}
			sel := &nodes.SelectStmt{
				TargetList: makeList(rt),
				FromClause: makeList(fromRv),
				WithClause: wc,
			}
			$$ = &nodes.ViewStmt{
				View:            rv,
				Aliases:         $7,
				Options:         $9,
				Query:           sel,
				WithCheckOption: int($12),
			}
		}
	| CREATE OR REPLACE OptTemp RECURSIVE VIEW qualified_name '(' columnList ')' opt_reloptions
	AS SelectStmt opt_check_option
		{
			rv := makeRangeVar($7).(*nodes.RangeVar)
			cte := &nodes.CommonTableExpr{
				Ctename:      rv.Relname,
				Ctequery:     $13,
				Aliascolnames:  $9,
			}
			wc := &nodes.WithClause{
				Ctes:      makeList(cte),
				Recursive: true,
			}
			cr := &nodes.ColumnRef{
				Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
				Location: -1,
			}
			rt := &nodes.ResTarget{Val: cr, Location: -1}
			fromRv := &nodes.RangeVar{Relname: rv.Relname, Inh: true, Relpersistence: 'p'}
			sel := &nodes.SelectStmt{
				TargetList: makeList(rt),
				FromClause: makeList(fromRv),
				WithClause: wc,
			}
			$$ = &nodes.ViewStmt{
				View:            rv,
				Aliases:         $9,
				Options:         $11,
				Query:           sel,
				Replace:         true,
				WithCheckOption: int($14),
			}
		}
	;

opt_check_option:
	WITH CHECK OPTION              { $$ = int64(VIEW_CHECK_OPTION_LOCAL) }
	| WITH CASCADED CHECK OPTION   { $$ = int64(VIEW_CHECK_OPTION_CASCADED) }
	| WITH LOCAL CHECK OPTION      { $$ = int64(VIEW_CHECK_OPTION_LOCAL) }
	| /* EMPTY */                  { $$ = int64(VIEW_CHECK_OPTION_NONE) }
	;

/*****************************************************************************
 *
 *      CREATE FUNCTION / CREATE PROCEDURE
 *
 *****************************************************************************/

CreateFunctionStmt:
	CREATE opt_or_replace FUNCTION func_name func_args_with_defaults
	RETURNS func_return opt_createfunc_opt_list opt_routine_body
		{
			n := &nodes.CreateFunctionStmt{
				IsOrReplace: $2,
				Funcname:    $4,
				Parameters:  $5,
				ReturnType:  $7,
				Options:     $8,
				SqlBody:     $9,
			}
			$$ = n
		}
	| CREATE opt_or_replace FUNCTION func_name func_args_with_defaults
	RETURNS TABLE '(' table_func_column_list ')' opt_createfunc_opt_list opt_routine_body
		{
			/* RETURNS TABLE(...) adds table columns to parameter list */
			params := concatLists($5, $9)
			n := &nodes.CreateFunctionStmt{
				IsOrReplace: $2,
				Funcname:    $4,
				Parameters:  params,
				ReturnType:  &nodes.TypeName{Names: makeFuncName("pg_catalog", "record"), Location: -1},
				Options:     $11,
				SqlBody:     $12,
			}
			$$ = n
		}
	| CREATE opt_or_replace FUNCTION func_name func_args_with_defaults
	opt_createfunc_opt_list opt_routine_body
		{
			n := &nodes.CreateFunctionStmt{
				IsOrReplace: $2,
				Funcname:    $4,
				Parameters:  $5,
				Options:     $6,
				SqlBody:     $7,
			}
			$$ = n
		}
	| CREATE opt_or_replace PROCEDURE func_name func_args_with_defaults
	opt_createfunc_opt_list opt_routine_body
		{
			n := &nodes.CreateFunctionStmt{
				IsOrReplace: $2,
				Funcname:    $4,
				Parameters:  $5,
				Options:     $6,
				SqlBody:     $7,
			}
			n.Options = appendList(n.Options, &nodes.DefElem{Defname: "isProcedure", Arg: &nodes.Integer{Ival: 1}})
			$$ = n
		}
	;

opt_or_replace:
	OR REPLACE   { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

func_args_with_defaults:
	'(' func_args_with_defaults_list ')' { $$ = $2 }
	| '(' ')' { $$ = nil }
	;

func_args_with_defaults_list:
	func_arg_with_default
		{ $$ = makeList($1) }
	| func_args_with_defaults_list ',' func_arg_with_default
		{ $$ = appendList($1, $3) }
	;

func_arg_with_default:
	func_arg
		{ $$ = $1 }
	| func_arg DEFAULT a_expr
		{
			fp := $1.(*nodes.FunctionParameter)
			fp.Defexpr = $3
			$$ = fp
		}
	| func_arg '=' a_expr
		{
			fp := $1.(*nodes.FunctionParameter)
			fp.Defexpr = $3
			$$ = fp
		}
	;

func_arg:
	arg_class param_name func_type
		{
			$$ = &nodes.FunctionParameter{
				Name:    $2,
				ArgType: $3,
				Mode:    nodes.FunctionParameterMode($1),
			}
		}
	| param_name arg_class func_type
		{
			$$ = &nodes.FunctionParameter{
				Name:    $1,
				ArgType: $3,
				Mode:    nodes.FunctionParameterMode($2),
			}
		}
	| param_name func_type
		{
			$$ = &nodes.FunctionParameter{
				Name:    $1,
				ArgType: $2,
				Mode:    nodes.FUNC_PARAM_IN,
			}
		}
	| arg_class func_type
		{
			$$ = &nodes.FunctionParameter{
				ArgType: $2,
				Mode:    nodes.FunctionParameterMode($1),
			}
		}
	| func_type
		{
			$$ = &nodes.FunctionParameter{
				ArgType: $1,
				Mode:    nodes.FUNC_PARAM_IN,
			}
		}
	;

arg_class:
	IN_P       { $$ = int64(nodes.FUNC_PARAM_IN) }
	| OUT_P    { $$ = int64(nodes.FUNC_PARAM_OUT) }
	| INOUT    { $$ = int64(nodes.FUNC_PARAM_INOUT) }
	| IN_P OUT_P { $$ = int64(nodes.FUNC_PARAM_INOUT) }
	| VARIADIC { $$ = int64(nodes.FUNC_PARAM_VARIADIC) }
	;

param_name:
	type_function_name { $$ = $1 }
	;

func_return:
	func_type { $$ = $1 }
	;

table_func_column_list:
	table_func_column
		{ $$ = makeList($1) }
	| table_func_column_list ',' table_func_column
		{ $$ = appendList($1, $3) }
	;

table_func_column:
	param_name func_type
		{
			$$ = &nodes.FunctionParameter{
				Name:    $1,
				ArgType: $2,
				Mode:    nodes.FUNC_PARAM_TABLE,
			}
		}
	;

func_type:
	Typename { $$ = $1 }
	| type_function_name attrs '%' TYPE_P
		{
			names := prependList(&nodes.String{Str: $1}, $2)
			tn := makeTypeNameFromNameList(names).(*nodes.TypeName)
			tn.PctType = true
			tn.Location = -1
			$$ = tn
		}
	| SETOF type_function_name attrs '%' TYPE_P
		{
			names := prependList(&nodes.String{Str: $2}, $3)
			tn := makeTypeNameFromNameList(names).(*nodes.TypeName)
			tn.PctType = true
			tn.Setof = true
			tn.Location = -1
			$$ = tn
		}
	;

opt_createfunc_opt_list:
	createfunc_opt_list { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

createfunc_opt_list:
	createfunc_opt_item
		{ $$ = makeList($1) }
	| createfunc_opt_list createfunc_opt_item
		{ $$ = appendList($1, $2) }
	;

opt_routine_body:
	ReturnStmt
		{
			$$ = $1
		}
	| BEGIN_P ATOMIC routine_body_stmt_list END_P
		{
			/* A compound statement stored as a single-item list containing the stmt list */
			$$ = makeList($3)
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

routine_body_stmt_list:
	routine_body_stmt_list routine_body_stmt ';'
		{
			if $2 != nil {
				$$ = appendList($1, $2)
			} else {
				$$ = $1
			}
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

routine_body_stmt:
	stmt
		{
			$$ = $1
		}
	| ReturnStmt
		{
			$$ = $1
		}
	;

ReturnStmt:
	RETURN a_expr
		{
			$$ = &nodes.ReturnStmt{
				Returnval: $2,
			}
		}
	;

createfunc_opt_item:
	LANGUAGE NonReservedWord_or_Sconst
		{
			$$ = &nodes.DefElem{
				Defname: "language",
				Arg:     &nodes.String{Str: $2},
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

common_func_opt_item:
	IMMUTABLE
		{
			$$ = &nodes.DefElem{
				Defname: "volatility",
				Arg:     &nodes.String{Str: "immutable"},
			}
		}
	| STABLE
		{
			$$ = &nodes.DefElem{
				Defname: "volatility",
				Arg:     &nodes.String{Str: "stable"},
			}
		}
	| VOLATILE
		{
			$$ = &nodes.DefElem{
				Defname: "volatility",
				Arg:     &nodes.String{Str: "volatile"},
			}
		}
	| STRICT_P
		{
			$$ = &nodes.DefElem{
				Defname: "strict",
				Arg:     &nodes.Integer{Ival: 1},
			}
		}
	| CALLED ON NULL_P INPUT_P
		{
			$$ = &nodes.DefElem{
				Defname: "strict",
				Arg:     &nodes.Integer{Ival: 0},
			}
		}
	| RETURNS NULL_P ON NULL_P INPUT_P
		{
			$$ = &nodes.DefElem{
				Defname: "strict",
				Arg:     &nodes.Integer{Ival: 1},
			}
		}
	| SECURITY DEFINER
		{
			$$ = &nodes.DefElem{
				Defname: "security",
				Arg:     &nodes.Integer{Ival: 1},
			}
		}
	| SECURITY INVOKER
		{
			$$ = &nodes.DefElem{
				Defname: "security",
				Arg:     &nodes.Integer{Ival: 0},
			}
		}
	| AS Sconst
		{
			$$ = &nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.String{Str: $2},
			}
		}
	| AS Sconst ',' Sconst
		{
			$$ = &nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: $2}, &nodes.String{Str: $4}}},
			}
		}
	| LEAKPROOF
		{
			$$ = &nodes.DefElem{
				Defname: "leakproof",
				Arg:     &nodes.Integer{Ival: 1},
			}
		}
	| NOT LEAKPROOF
		{
			$$ = &nodes.DefElem{
				Defname: "leakproof",
				Arg:     &nodes.Integer{Ival: 0},
			}
		}
	| COST NumericOnly
		{
			$$ = &nodes.DefElem{
				Defname: "cost",
				Arg:     $2,
			}
		}
	| ROWS NumericOnly
		{
			$$ = &nodes.DefElem{
				Defname: "rows",
				Arg:     $2,
			}
		}
	| PARALLEL ColId
		{
			$$ = &nodes.DefElem{
				Defname: "parallel",
				Arg:     &nodes.String{Str: $2},
			}
		}
	| SET set_rest_more
		{
			$$ = &nodes.DefElem{
				Defname: "set",
				Arg:     $2,
			}
		}
	| VariableResetStmt
		{
			$$ = &nodes.DefElem{
				Defname: "set",
				Arg:     $1,
			}
		}
	| SUPPORT any_name
		{
			$$ = &nodes.DefElem{
				Defname: "support",
				Arg:     $2,
			}
		}
	;

/*****************************************************************************
 *
 *      Transaction statements
 *
 *****************************************************************************/

TransactionStmt:
	ABORT_P opt_transaction opt_transaction_chain
		{
			$$ = &nodes.TransactionStmt{
				Kind:  nodes.TRANS_STMT_ROLLBACK,
				Chain: $3,
			}
		}
	| BEGIN_P opt_transaction transaction_mode_list_or_empty
		{
			$$ = &nodes.TransactionStmt{
				Kind:    nodes.TRANS_STMT_BEGIN,
				Options: $3,
			}
		}
	| START TRANSACTION transaction_mode_list_or_empty
		{
			$$ = &nodes.TransactionStmt{
				Kind:    nodes.TRANS_STMT_START,
				Options: $3,
			}
		}
	| PREPARE TRANSACTION Sconst
		{
			$$ = &nodes.TransactionStmt{
				Kind: nodes.TRANS_STMT_PREPARE,
				Gid:  $3,
			}
		}
	| COMMIT PREPARED Sconst
		{
			$$ = &nodes.TransactionStmt{
				Kind: nodes.TRANS_STMT_COMMIT_PREPARED,
				Gid:  $3,
			}
		}
	| ROLLBACK PREPARED Sconst
		{
			$$ = &nodes.TransactionStmt{
				Kind: nodes.TRANS_STMT_ROLLBACK_PREPARED,
				Gid:  $3,
			}
		}
	| COMMIT opt_transaction opt_transaction_chain
		{
			$$ = &nodes.TransactionStmt{
				Kind:  nodes.TRANS_STMT_COMMIT,
				Chain: $3,
			}
		}
	| END_P opt_transaction opt_transaction_chain
		{
			$$ = &nodes.TransactionStmt{
				Kind:  nodes.TRANS_STMT_COMMIT,
				Chain: $3,
			}
		}
	| ROLLBACK opt_transaction opt_transaction_chain
		{
			$$ = &nodes.TransactionStmt{
				Kind:  nodes.TRANS_STMT_ROLLBACK,
				Chain: $3,
			}
		}
	| SAVEPOINT ColId
		{
			$$ = &nodes.TransactionStmt{
				Kind:      nodes.TRANS_STMT_SAVEPOINT,
				Savepoint: $2,
			}
		}
	| RELEASE SAVEPOINT ColId
		{
			$$ = &nodes.TransactionStmt{
				Kind:      nodes.TRANS_STMT_RELEASE,
				Savepoint: $3,
			}
		}
	| RELEASE ColId
		{
			$$ = &nodes.TransactionStmt{
				Kind:      nodes.TRANS_STMT_RELEASE,
				Savepoint: $2,
			}
		}
	| ROLLBACK opt_transaction TO SAVEPOINT ColId
		{
			$$ = &nodes.TransactionStmt{
				Kind:      nodes.TRANS_STMT_ROLLBACK_TO,
				Savepoint: $5,
			}
		}
	| ROLLBACK opt_transaction TO ColId
		{
			$$ = &nodes.TransactionStmt{
				Kind:      nodes.TRANS_STMT_ROLLBACK_TO,
				Savepoint: $4,
			}
		}
	;

opt_transaction:
	WORK        {}
	| TRANSACTION {}
	| /* EMPTY */ {}
	;

opt_transaction_chain:
	AND CHAIN      { $$ = true }
	| AND NO CHAIN { $$ = false }
	| /* EMPTY */  { $$ = false }
	;

/*****************************************************************************
 *
 *      EXPLAIN statement
 *
 *****************************************************************************/

ExplainStmt:
	EXPLAIN ExplainableStmt
		{
			$$ = &nodes.ExplainStmt{
				Query: $2,
			}
		}
	| EXPLAIN ANALYZE ExplainableStmt
		{
			$$ = &nodes.ExplainStmt{
				Query:   $3,
				Options: makeList(&nodes.DefElem{Defname: "analyze"}),
			}
		}
	| EXPLAIN VERBOSE ExplainableStmt
		{
			$$ = &nodes.ExplainStmt{
				Query:   $3,
				Options: makeList(&nodes.DefElem{Defname: "verbose"}),
			}
		}
	| EXPLAIN ANALYZE VERBOSE ExplainableStmt
		{
			$$ = &nodes.ExplainStmt{
				Query: $4,
				Options: &nodes.List{Items: []nodes.Node{
					&nodes.DefElem{Defname: "analyze"},
					&nodes.DefElem{Defname: "verbose"},
				}},
			}
		}
	| EXPLAIN '(' utility_option_list ')' ExplainableStmt
		{
			$$ = &nodes.ExplainStmt{
				Query:   $5,
				Options: $3,
			}
		}
	;

ExplainableStmt:
	SelectStmt          { $$ = $1 }
	| InsertStmt        { $$ = $1 }
	| UpdateStmt        { $$ = $1 }
	| DeleteStmt        { $$ = $1 }
	| MergeStmt         { $$ = $1 }
	| DeclareCursorStmt { $$ = $1 }
	| CreateAsStmt      { $$ = $1 }
	| CreateMatViewStmt { $$ = $1 }
	| RefreshMatViewStmt { $$ = $1 }
	| ExecuteStmt       { $$ = $1 }
	| CreateStmt        { $$ = $1 }
	;

/*****************************************************************************
 *
 *      COPY statement
 *
 *****************************************************************************/

CopyStmt:
	COPY relation_expr opt_column_list copy_from copy_file_name copy_opt_with '(' utility_option_list ')' where_clause
		{
			rv := $2.(*nodes.RangeVar)
			$$ = &nodes.CopyStmt{
				Relation:    rv,
				Attlist:     $3,
				IsFrom:      $4,
				Filename:    $5,
				Options:     $8,
				WhereClause: $10,
			}
		}
	| COPY relation_expr opt_column_list copy_from copy_file_name copy_opt_list where_clause
		{
			rv := $2.(*nodes.RangeVar)
			$$ = &nodes.CopyStmt{
				Relation:    rv,
				Attlist:     $3,
				IsFrom:      $4,
				Filename:    $5,
				Options:     $6,
				WhereClause: $7,
			}
		}
	| COPY relation_expr opt_column_list copy_from PROGRAM Sconst copy_opt_with '(' utility_option_list ')' where_clause
		{
			rv := $2.(*nodes.RangeVar)
			$$ = &nodes.CopyStmt{
				Relation:    rv,
				Attlist:     $3,
				IsFrom:      $4,
				IsProgram:   true,
				Filename:    $6,
				Options:     $9,
				WhereClause: $11,
			}
		}
	| COPY relation_expr opt_column_list copy_from PROGRAM Sconst copy_opt_list where_clause
		{
			rv := $2.(*nodes.RangeVar)
			$$ = &nodes.CopyStmt{
				Relation:    rv,
				Attlist:     $3,
				IsFrom:      $4,
				IsProgram:   true,
				Filename:    $6,
				Options:     $7,
				WhereClause: $8,
			}
		}
	| COPY '(' PreparableStmt ')' TO PROGRAM Sconst copy_opt_with '(' utility_option_list ')'
		{
			$$ = &nodes.CopyStmt{
				Query:     $3,
				IsFrom:    false,
				IsProgram: true,
				Filename:  $7,
				Options:   $10,
			}
		}
	| COPY '(' PreparableStmt ')' TO PROGRAM Sconst copy_opt_list
		{
			$$ = &nodes.CopyStmt{
				Query:     $3,
				IsFrom:    false,
				IsProgram: true,
				Filename:  $7,
				Options:   $8,
			}
		}
	| COPY '(' PreparableStmt ')' TO copy_file_name copy_opt_with '(' utility_option_list ')'
		{
			$$ = &nodes.CopyStmt{
				Query:    $3,
				IsFrom:   false,
				Filename: $6,
				Options:  $9,
			}
		}
	| COPY '(' PreparableStmt ')' TO copy_file_name copy_opt_list
		{
			$$ = &nodes.CopyStmt{
				Query:    $3,
				IsFrom:   false,
				Filename: $6,
				Options:  $7,
			}
		}
	;

copy_from:
	FROM { $$ = true }
	| TO { $$ = false }
	;

copy_file_name:
	Sconst      { $$ = $1 }
	| STDIN     { $$ = "" }
	| STDOUT    { $$ = "" }
	;

copy_opt_with:
	WITH {}
	| /* EMPTY */ {}
	;

/* old-style COPY option syntax */
copy_opt_list:
	copy_opt_list copy_opt_item
		{
			if $2 != nil {
				$$ = appendList($1, $2)
			} else {
				$$ = $1
			}
		}
	| WITH copy_opt_item
		{
			if $2 != nil {
				$$ = makeList($2)
			} else {
				$$ = nil
			}
		}
	| /* EMPTY */ { $$ = nil }
	;

copy_opt_item:
	BINARY
		{
			$$ = &nodes.DefElem{Defname: "format", Arg: &nodes.String{Str: "binary"}}
		}
	| FREEZE
		{
			$$ = &nodes.DefElem{Defname: "freeze", Arg: &nodes.Boolean{Boolval: true}}
		}
	| DELIMITER opt_as Sconst
		{
			$$ = &nodes.DefElem{Defname: "delimiter", Arg: &nodes.String{Str: $3}}
		}
	| NULL_P opt_as Sconst
		{
			$$ = &nodes.DefElem{Defname: "null", Arg: &nodes.String{Str: $3}}
		}
	| CSV
		{
			$$ = &nodes.DefElem{Defname: "format", Arg: &nodes.String{Str: "csv"}}
		}
	| HEADER_P
		{
			$$ = &nodes.DefElem{Defname: "header", Arg: &nodes.Boolean{Boolval: true}}
		}
	| QUOTE opt_as Sconst
		{
			$$ = &nodes.DefElem{Defname: "quote", Arg: &nodes.String{Str: $3}}
		}
	| ESCAPE opt_as Sconst
		{
			$$ = &nodes.DefElem{Defname: "escape", Arg: &nodes.String{Str: $3}}
		}
	| FORCE QUOTE columnList
		{
			$$ = &nodes.DefElem{Defname: "force_quote", Arg: $3}
		}
	| FORCE QUOTE '*'
		{
			$$ = &nodes.DefElem{Defname: "force_quote", Arg: &nodes.A_Star{}}
		}
	| FORCE NOT NULL_P columnList
		{
			$$ = &nodes.DefElem{Defname: "force_not_null", Arg: $4}
		}
	| FORCE NULL_P columnList
		{
			$$ = &nodes.DefElem{Defname: "force_null", Arg: $3}
		}
	| ENCODING Sconst
		{
			$$ = &nodes.DefElem{Defname: "encoding", Arg: &nodes.String{Str: $2}}
		}
	;

utility_option_list:
	utility_option_elem
		{
			$$ = makeList($1)
		}
	| utility_option_list ',' utility_option_elem
		{
			$$ = appendList($1, $3)
		}
	;

utility_option_elem:
	utility_option_name utility_option_arg
		{
			$$ = &nodes.DefElem{
				Defname: $1,
				Arg:     $2,
			}
		}
	;

utility_option_name:
	ColId            { $$ = $1 }
	| ANALYZE        { $$ = "analyze" }
	| VERBOSE        { $$ = "verbose" }
	| FULL           { $$ = "full" }
	| FREEZE         { $$ = "freeze" }
	| FORMAT         { $$ = "format" }
	| NULL_P         { $$ = "null" }
	| DEFAULT        { $$ = "default" }
	| FORCE          { $$ = "force" }
	| CONCURRENTLY   { $$ = "concurrently" }
	;

utility_option_arg:
	opt_boolean_or_string   { $$ = &nodes.String{Str: $1} }
	| NumericOnly           { $$ = $1 }
	| '*'                   { $$ = &nodes.A_Star{} }
	| DEFAULT               { $$ = &nodes.String{Str: "default"} }
	| '(' utility_option_arg_list ')' { $$ = $2 }
	| /* EMPTY */           { $$ = nil }
	;

utility_option_arg_list:
	opt_boolean_or_string
		{ $$ = makeList(&nodes.String{Str: $1}) }
	| utility_option_arg_list ',' opt_boolean_or_string
		{ $$ = appendList($1, &nodes.String{Str: $3}) }
	;

opt_boolean_or_string:
	TRUE_P                      { $$ = "true" }
	| FALSE_P                   { $$ = "false" }
	| ON                        { $$ = "on" }
	| NonReservedWord_or_Sconst { $$ = $1 }
	;

NumericOnly:
	FCONST
		{
			$$ = &nodes.Float{Fval: $1}
		}
	| '+' FCONST
		{
			$$ = &nodes.Float{Fval: $2}
		}
	| '-' FCONST
		{
			f := &nodes.Float{Fval: $2}
			doNegateFloat(f)
			$$ = f
		}
	| SignedIconst
		{
			$$ = &nodes.Integer{Ival: int64($1)}
		}
	;

SignedIconst:
	Iconst      { $$ = $1 }
	| '+' Iconst { $$ = $2 }
	| '-' Iconst { $$ = -$2 }
	;

NonReservedWord_or_Sconst:
	NonReservedWord { $$ = $1 }
	| Sconst        { $$ = $1 }
	;

NonReservedWord:
	IDENT                    { $$ = $1 }
	| unreserved_keyword     { $$ = $1 }
	| col_name_keyword       { $$ = $1 }
	| type_func_name_keyword { $$ = $1 }
	;

/*****************************************************************************
 *
 *      GRANT / REVOKE statements
 *
 *****************************************************************************/

GrantStmt:
	GRANT privileges ON TABLE any_name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TABLE,
				Objects:     makeRangeVarList($5),
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON any_name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TABLE,
				Objects:     makeRangeVarList($4),
				Privileges:  $2,
				Grantees:    $6,
				GrantOption: $7,
			}
		}
	| GRANT privileges ON SEQUENCE any_name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_SEQUENCE,
				Objects:     makeRangeVarList($5),
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON FUNCTION function_with_argtypes_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_FUNCTION,
				Objects:     $5,
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON PROCEDURE function_with_argtypes_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_PROCEDURE,
				Objects:     $5,
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON ROUTINE function_with_argtypes_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_ROUTINE,
				Objects:     $5,
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON DATABASE name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_DATABASE,
				Objects:     makeNameListAsAnyNameList($5),
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON DOMAIN_P any_name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_DOMAIN,
				Objects:     $5,
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON LANGUAGE name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_LANGUAGE,
				Objects:     makeNameListAsAnyNameList($5),
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON SCHEMA name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_SCHEMA,
				Objects:     makeNameListAsAnyNameList($5),
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON TABLESPACE name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TABLESPACE,
				Objects:     makeNameListAsAnyNameList($5),
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON TYPE_P any_name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TYPE,
				Objects:     $5,
				Privileges:  $2,
				Grantees:    $7,
				GrantOption: $8,
			}
		}
	| GRANT privileges ON FOREIGN DATA_P WRAPPER name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_FDW,
				Objects:     makeNameListAsAnyNameList($7),
				Privileges:  $2,
				Grantees:    $9,
				GrantOption: $10,
			}
		}
	| GRANT privileges ON FOREIGN SERVER name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_FOREIGN_SERVER,
				Objects:     makeNameListAsAnyNameList($6),
				Privileges:  $2,
				Grantees:    $8,
				GrantOption: $9,
			}
		}
	| GRANT privileges ON ALL TABLES IN_P SCHEMA name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:     nodes.OBJECT_TABLE,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $2,
				Grantees:    $10,
				GrantOption: $11,
			}
		}
	| GRANT privileges ON ALL SEQUENCES IN_P SCHEMA name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:     nodes.OBJECT_SEQUENCE,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $2,
				Grantees:    $10,
				GrantOption: $11,
			}
		}
	| GRANT privileges ON ALL FUNCTIONS IN_P SCHEMA name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:     nodes.OBJECT_FUNCTION,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $2,
				Grantees:    $10,
				GrantOption: $11,
			}
		}
	| GRANT privileges ON ALL ROUTINES IN_P SCHEMA name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:     nodes.OBJECT_ROUTINE,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $2,
				Grantees:    $10,
				GrantOption: $11,
			}
		}
	| GRANT privileges ON ALL PROCEDURES IN_P SCHEMA name_list TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:     nodes.OBJECT_PROCEDURE,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $2,
				Grantees:    $10,
				GrantOption: $11,
			}
		}
	| GRANT privileges ON LARGE_P OBJECT_P NumericOnly TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_LARGEOBJECT,
				Objects:     makeList($6),
				Privileges:  $2,
				Grantees:    $8,
				GrantOption: $9,
			}
		}
	;

RevokeStmt:
	REVOKE privileges ON TABLE any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_TABLE,
				Objects:    makeRangeVarList($5),
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_TABLE,
				Objects:    makeRangeVarList($4),
				Privileges: $2,
				Grantees:   $6,
				Behavior:   nodes.DropBehavior($7),
			}
		}
	| REVOKE privileges ON SEQUENCE any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_SEQUENCE,
				Objects:    makeRangeVarList($5),
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON FUNCTION function_with_argtypes_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_FUNCTION,
				Objects:    $5,
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON PROCEDURE function_with_argtypes_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_PROCEDURE,
				Objects:    $5,
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON ROUTINE function_with_argtypes_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_ROUTINE,
				Objects:    $5,
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON DATABASE name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_DATABASE,
				Objects:    makeNameListAsAnyNameList($5),
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON DOMAIN_P any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_DOMAIN,
				Objects:    $5,
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON LANGUAGE name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_LANGUAGE,
				Objects:    makeNameListAsAnyNameList($5),
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON SCHEMA name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_SCHEMA,
				Objects:    makeNameListAsAnyNameList($5),
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON TABLESPACE name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_TABLESPACE,
				Objects:    makeNameListAsAnyNameList($5),
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON TYPE_P any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_TYPE,
				Objects:    $5,
				Privileges: $2,
				Grantees:   $7,
				Behavior:   nodes.DropBehavior($8),
			}
		}
	| REVOKE privileges ON FOREIGN DATA_P WRAPPER name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_FDW,
				Objects:    makeNameListAsAnyNameList($7),
				Privileges: $2,
				Grantees:   $9,
				Behavior:   nodes.DropBehavior($10),
			}
		}
	| REVOKE privileges ON FOREIGN SERVER name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_FOREIGN_SERVER,
				Objects:    makeNameListAsAnyNameList($6),
				Privileges: $2,
				Grantees:   $8,
				Behavior:   nodes.DropBehavior($9),
			}
		}
	| REVOKE privileges ON ALL TABLES IN_P SCHEMA name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:    nodes.OBJECT_TABLE,
				Objects:    makeNameListAsAnyNameList($8),
				Privileges: $2,
				Grantees:   $10,
				Behavior:   nodes.DropBehavior($11),
			}
		}
	| REVOKE privileges ON ALL SEQUENCES IN_P SCHEMA name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:    nodes.OBJECT_SEQUENCE,
				Objects:    makeNameListAsAnyNameList($8),
				Privileges: $2,
				Grantees:   $10,
				Behavior:   nodes.DropBehavior($11),
			}
		}
	| REVOKE privileges ON ALL FUNCTIONS IN_P SCHEMA name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:    nodes.OBJECT_FUNCTION,
				Objects:    makeNameListAsAnyNameList($8),
				Privileges: $2,
				Grantees:   $10,
				Behavior:   nodes.DropBehavior($11),
			}
		}
	| REVOKE privileges ON ALL ROUTINES IN_P SCHEMA name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:    nodes.OBJECT_ROUTINE,
				Objects:    makeNameListAsAnyNameList($8),
				Privileges: $2,
				Grantees:   $10,
				Behavior:   nodes.DropBehavior($11),
			}
		}
	| REVOKE privileges ON ALL PROCEDURES IN_P SCHEMA name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_ALL_IN_SCHEMA,
				Objtype:    nodes.OBJECT_PROCEDURE,
				Objects:    makeNameListAsAnyNameList($8),
				Privileges: $2,
				Grantees:   $10,
				Behavior:   nodes.DropBehavior($11),
			}
		}
	| REVOKE privileges ON LARGE_P OBJECT_P NumericOnly FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Targtype:   nodes.ACL_TARGET_OBJECT,
				Objtype:    nodes.OBJECT_LARGEOBJECT,
				Objects:    makeList($6),
				Privileges: $2,
				Grantees:   $8,
				Behavior:   nodes.DropBehavior($9),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON TABLE any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TABLE,
				Objects:     makeRangeVarList($8),
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TABLE,
				Objects:     makeRangeVarList($7),
				Privileges:  $5,
				Grantees:    $9,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($10),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON SEQUENCE any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_SEQUENCE,
				Objects:     makeRangeVarList($8),
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON FUNCTION function_with_argtypes_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_FUNCTION,
				Objects:     $8,
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON PROCEDURE function_with_argtypes_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_PROCEDURE,
				Objects:     $8,
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON ROUTINE function_with_argtypes_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_ROUTINE,
				Objects:     $8,
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON DATABASE name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_DATABASE,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON DOMAIN_P any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_DOMAIN,
				Objects:     $8,
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON LANGUAGE name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_LANGUAGE,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON SCHEMA name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_SCHEMA,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON TABLESPACE name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TABLESPACE,
				Objects:     makeNameListAsAnyNameList($8),
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON TYPE_P any_name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_TYPE,
				Objects:     $8,
				Privileges:  $5,
				Grantees:    $10,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($11),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON FOREIGN DATA_P WRAPPER name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_FDW,
				Objects:     makeNameListAsAnyNameList($10),
				Privileges:  $5,
				Grantees:    $12,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($13),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON FOREIGN SERVER name_list FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				Targtype:    nodes.ACL_TARGET_OBJECT,
				Objtype:     nodes.OBJECT_FOREIGN_SERVER,
				Objects:     makeNameListAsAnyNameList($9),
				Privileges:  $5,
				Grantees:    $11,
				GrantOption: true,
				Behavior:    nodes.DropBehavior($12),
			}
		}
	;

privileges:
	ALL PRIVILEGES '(' columnList ')'
		{
			ap := &nodes.AccessPriv{Cols: $4}
			$$ = makeList(ap)
		}
	| ALL '(' columnList ')'
		{
			ap := &nodes.AccessPriv{Cols: $3}
			$$ = makeList(ap)
		}
	| ALL PRIVILEGES  { $$ = nil }
	| ALL           { $$ = nil }
	| privilege_list { $$ = $1 }
	;

privilege_list:
	privilege
		{ $$ = makeList($1) }
	| privilege_list ',' privilege
		{ $$ = appendList($1, $3) }
	;

privilege:
	SELECT opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "select", Cols: $2} }
	| REFERENCES opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "references", Cols: $2} }
	| CREATE opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "create", Cols: $2} }
	| INSERT opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "insert", Cols: $2} }
	| UPDATE opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "update", Cols: $2} }
	| DELETE_P opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "delete", Cols: $2} }
	| TRIGGER opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "trigger", Cols: $2} }
	| EXECUTE opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "execute", Cols: $2} }
	| TRUNCATE opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "truncate", Cols: $2} }
	| ColId opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: $1, Cols: $2} }
	;

grantee_list:
	grantee
		{ $$ = makeList($1) }
	| grantee_list ',' grantee
		{ $$ = appendList($1, $3) }
	;

grantee:
	RoleSpec { $$ = $1 }
	| GROUP_P RoleSpec { $$ = $2 }
	;

RoleSpec:
	ColId
		{
			$$ = &nodes.RoleSpec{
				Roletype: int(nodes.ROLESPEC_CSTRING),
				Rolename: $1,
			}
		}
	| CURRENT_ROLE
		{
			$$ = &nodes.RoleSpec{
				Roletype: int(nodes.ROLESPEC_CURRENT_ROLE),
			}
		}
	| CURRENT_USER
		{
			$$ = &nodes.RoleSpec{
				Roletype: int(nodes.ROLESPEC_CURRENT_USER),
			}
		}
	| SESSION_USER
		{
			$$ = &nodes.RoleSpec{
				Roletype: int(nodes.ROLESPEC_SESSION_USER),
			}
		}
	;

opt_grant_grant_option:
	WITH GRANT OPTION { $$ = true }
	| /* EMPTY */     { $$ = false }
	;

/*****************************************************************************
 *
 *      GRANT / REVOKE ROLE statements
 *
 *****************************************************************************/

GrantRoleStmt:
	GRANT privilege_list TO role_list opt_granted_by
		{
			$$ = &nodes.GrantRoleStmt{
				IsGrant:      true,
				GrantedRoles: $2,
				GranteeRoles: $4,
				Grantor:      roleSpecOrNil($5),
			}
		}
	| GRANT privilege_list TO role_list WITH grant_role_opt_list opt_granted_by
		{
			$$ = &nodes.GrantRoleStmt{
				IsGrant:      true,
				GrantedRoles: $2,
				GranteeRoles: $4,
				Opt:          $6,
				Grantor:      roleSpecOrNil($7),
			}
		}
	;

RevokeRoleStmt:
	REVOKE privilege_list FROM role_list opt_granted_by opt_drop_behavior
		{
			$$ = &nodes.GrantRoleStmt{
				IsGrant:      false,
				GrantedRoles: $2,
				GranteeRoles: $4,
				Grantor:      roleSpecOrNil($5),
				Behavior:     nodes.DropBehavior($6),
			}
		}
	| REVOKE ColId OPTION FOR privilege_list FROM role_list opt_granted_by opt_drop_behavior
		{
			opt := makeDefElem($2, &nodes.Boolean{Boolval: false})
			$$ = &nodes.GrantRoleStmt{
				IsGrant:      false,
				Opt:          makeList(opt),
				GrantedRoles: $5,
				GranteeRoles: $7,
				Grantor:      roleSpecOrNil($8),
				Behavior:     nodes.DropBehavior($9),
			}
		}
	;

grant_role_opt_list:
	grant_role_opt_list ',' grant_role_opt
		{ $$ = appendList($1, $3) }
	| grant_role_opt
		{ $$ = makeList($1) }
	;

grant_role_opt:
	ColLabel grant_role_opt_value
		{
			$$ = makeDefElem($1, $2)
		}
	;

grant_role_opt_value:
	OPTION   { $$ = &nodes.Boolean{Boolval: true} }
	| TRUE_P { $$ = &nodes.Boolean{Boolval: true} }
	| FALSE_P { $$ = &nodes.Boolean{Boolval: false} }
	;

opt_granted_by:
	GRANTED BY RoleSpec { $$ = $3 }
	| /* EMPTY */       { $$ = nil }
	;

/*****************************************************************************
 *
 *      CREATE ROLE / USER / GROUP
 *
 *****************************************************************************/

CreateRoleStmt:
	CREATE ROLE RoleId opt_with OptRoleList
		{
			$$ = &nodes.CreateRoleStmt{
				StmtType: nodes.ROLESTMT_ROLE,
				Role:     $3,
				Options:  $5,
			}
		}
	;

CreateUserStmt:
	CREATE USER RoleId opt_with OptRoleList
		{
			$$ = &nodes.CreateRoleStmt{
				StmtType: nodes.ROLESTMT_USER,
				Role:     $3,
				Options:  $5,
			}
		}
	;

CreateGroupStmt:
	CREATE GROUP_P RoleId opt_with OptRoleList
		{
			$$ = &nodes.CreateRoleStmt{
				StmtType: nodes.ROLESTMT_GROUP,
				Role:     $3,
				Options:  $5,
			}
		}
	;

opt_with:
	WITH      {}
	| /* EMPTY */ {}
	;

OptRoleList:
	OptRoleList CreateOptRoleElem
		{
			$$ = appendList($1, $2)
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

AlterOptRoleList:
	AlterOptRoleList AlterOptRoleElem
		{
			$$ = appendList($1, $2)
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

AlterOptRoleElem:
	PASSWORD Sconst
		{
			$$ = makeDefElem("password", &nodes.String{Str: $2})
		}
	| PASSWORD NULL_P
		{
			$$ = makeDefElem("password", nil)
		}
	| ENCRYPTED PASSWORD Sconst
		{
			$$ = makeDefElem("password", &nodes.String{Str: $3})
		}
	| UNENCRYPTED PASSWORD Sconst
		{
			pglex.Error("UNENCRYPTED PASSWORD is no longer supported")
			$$ = nil
		}
	| INHERIT
		{
			$$ = makeDefElem("inherit", &nodes.Boolean{Boolval: true})
		}
	| CONNECTION LIMIT SignedIconst
		{
			$$ = makeDefElem("connectionlimit", &nodes.Integer{Ival: int64($3)})
		}
	| VALID UNTIL Sconst
		{
			$$ = makeDefElem("validUntil", &nodes.String{Str: $3})
		}
	| USER role_list
		{
			$$ = makeDefElem("rolemembers", $2)
		}
	| IDENT
		{
			switch $1 {
			case "superuser":
				$$ = makeDefElem("superuser", &nodes.Boolean{Boolval: true})
			case "nosuperuser":
				$$ = makeDefElem("superuser", &nodes.Boolean{Boolval: false})
			case "createrole":
				$$ = makeDefElem("createrole", &nodes.Boolean{Boolval: true})
			case "nocreaterole":
				$$ = makeDefElem("createrole", &nodes.Boolean{Boolval: false})
			case "replication":
				$$ = makeDefElem("isreplication", &nodes.Boolean{Boolval: true})
			case "noreplication":
				$$ = makeDefElem("isreplication", &nodes.Boolean{Boolval: false})
			case "createdb":
				$$ = makeDefElem("createdb", &nodes.Boolean{Boolval: true})
			case "nocreatedb":
				$$ = makeDefElem("createdb", &nodes.Boolean{Boolval: false})
			case "login":
				$$ = makeDefElem("canlogin", &nodes.Boolean{Boolval: true})
			case "nologin":
				$$ = makeDefElem("canlogin", &nodes.Boolean{Boolval: false})
			case "bypassrls":
				$$ = makeDefElem("bypassrls", &nodes.Boolean{Boolval: true})
			case "nobypassrls":
				$$ = makeDefElem("bypassrls", &nodes.Boolean{Boolval: false})
			case "noinherit":
				$$ = makeDefElem("inherit", &nodes.Boolean{Boolval: false})
			default:
				pglex.Error("unrecognized role option \"" + $1 + "\"")
				$$ = nil
			}
		}
	;

CreateOptRoleElem:
	AlterOptRoleElem
		{
			$$ = $1
		}
	| SYSID Iconst
		{
			$$ = makeDefElem("sysid", &nodes.Integer{Ival: int64($2)})
		}
	| ADMIN role_list
		{
			$$ = makeDefElem("adminmembers", $2)
		}
	| ROLE role_list
		{
			$$ = makeDefElem("rolemembers", $2)
		}
	| IN_P ROLE role_list
		{
			$$ = makeDefElem("addroleto", $3)
		}
	| IN_P GROUP_P role_list
		{
			$$ = makeDefElem("addroleto", $3)
		}
	;

/*****************************************************************************
 *
 *      ALTER ROLE / USER (with inline opt_with to avoid reduce-reduce conflicts)
 *
 *****************************************************************************/

AlterRoleStmt:
	ALTER ROLE RoleSpec WITH AlterOptRoleList
		{
			$$ = &nodes.AlterRoleStmt{
				Role:    $3.(*nodes.RoleSpec),
				Action:  1,
				Options: $5,
			}
		}
	| ALTER ROLE RoleSpec AlterOptRoleElem AlterOptRoleList
		{
			$$ = &nodes.AlterRoleStmt{
				Role:    $3.(*nodes.RoleSpec),
				Action:  1,
				Options: prependList($4, $5),
			}
		}
	| ALTER USER RoleSpec WITH AlterOptRoleList
		{
			$$ = &nodes.AlterRoleStmt{
				Role:    $3.(*nodes.RoleSpec),
				Action:  1,
				Options: $5,
			}
		}
	| ALTER USER RoleSpec AlterOptRoleElem AlterOptRoleList
		{
			$$ = &nodes.AlterRoleStmt{
				Role:    $3.(*nodes.RoleSpec),
				Action:  1,
				Options: prependList($4, $5),
			}
		}
	;

/*****************************************************************************
 *
 *      ALTER ROLE SET/RESET (inlined opt_in_database to avoid reduce-reduce conflict)
 *
 *****************************************************************************/

AlterRoleSetStmt:
	ALTER ROLE RoleSpec SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Role:    $3.(*nodes.RoleSpec),
				Setstmt: $4.(*nodes.VariableSetStmt),
			}
		}
	| ALTER ROLE RoleSpec IN_P DATABASE name SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Role:     $3.(*nodes.RoleSpec),
				Database: $6,
				Setstmt:  $7.(*nodes.VariableSetStmt),
			}
		}
	| ALTER ROLE ALL SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Setstmt: $4.(*nodes.VariableSetStmt),
			}
		}
	| ALTER ROLE ALL IN_P DATABASE name SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Database: $6,
				Setstmt:  $7.(*nodes.VariableSetStmt),
			}
		}
	| ALTER USER RoleSpec SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Role:    $3.(*nodes.RoleSpec),
				Setstmt: $4.(*nodes.VariableSetStmt),
			}
		}
	| ALTER USER RoleSpec IN_P DATABASE name SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Role:     $3.(*nodes.RoleSpec),
				Database: $6,
				Setstmt:  $7.(*nodes.VariableSetStmt),
			}
		}
	| ALTER USER ALL SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Setstmt: $4.(*nodes.VariableSetStmt),
			}
		}
	| ALTER USER ALL IN_P DATABASE name SetResetClause
		{
			$$ = &nodes.AlterRoleSetStmt{
				Database: $6,
				Setstmt:  $7.(*nodes.VariableSetStmt),
			}
		}
	;

SetResetClause:
	SET set_rest
		{
			$$ = $2
		}
	| VariableResetStmt
		{
			$$ = $1
		}
	;

/*****************************************************************************
 *
 *      DROP ROLE / USER / GROUP
 *
 *****************************************************************************/

DropRoleStmt:
	DROP ROLE role_list
		{
			$$ = &nodes.DropRoleStmt{
				Roles:     $3,
				MissingOk: false,
			}
		}
	| DROP ROLE IF_P EXISTS role_list
		{
			$$ = &nodes.DropRoleStmt{
				Roles:     $5,
				MissingOk: true,
			}
		}
	| DROP USER role_list
		{
			$$ = &nodes.DropRoleStmt{
				Roles:     $3,
				MissingOk: false,
			}
		}
	| DROP USER IF_P EXISTS role_list
		{
			$$ = &nodes.DropRoleStmt{
				Roles:     $5,
				MissingOk: true,
			}
		}
	| DROP GROUP_P role_list
		{
			$$ = &nodes.DropRoleStmt{
				Roles:     $3,
				MissingOk: false,
			}
		}
	| DROP GROUP_P IF_P EXISTS role_list
		{
			$$ = &nodes.DropRoleStmt{
				Roles:     $5,
				MissingOk: true,
			}
		}
	;

/*****************************************************************************
 *
 *      ALTER GROUP
 *
 *****************************************************************************/

AlterGroupStmt:
	ALTER GROUP_P RoleSpec add_drop USER role_list
		{
			$$ = &nodes.AlterRoleStmt{
				Role:   $3.(*nodes.RoleSpec),
				Action: int($4),
				Options: makeList(makeDefElem("rolemembers", $6)),
			}
		}
	;

add_drop:
	ADD_P  { $$ = 1 }
	| DROP { $$ = -1 }
	;

/*****************************************************************************
 *
 *      CREATE DATABASE
 *
 *****************************************************************************/

CreatedbStmt:
	CREATE DATABASE name opt_with createdb_opt_list
		{
			$$ = &nodes.CreatedbStmt{
				Dbname:  $3,
				Options: $5,
			}
		}
	;

createdb_opt_list:
	createdb_opt_items
		{ $$ = $1 }
	| /* EMPTY */
		{ $$ = nil }
	;

createdb_opt_items:
	createdb_opt_item
		{ $$ = makeList($1) }
	| createdb_opt_items createdb_opt_item
		{ $$ = appendList($1, $2) }
	;

createdb_opt_item:
	createdb_opt_name opt_equal NumericOnly
		{
			$$ = makeDefElem($1, $3)
		}
	| createdb_opt_name opt_equal opt_boolean_or_string
		{
			$$ = makeDefElem($1, &nodes.String{Str: $3})
		}
	| createdb_opt_name opt_equal DEFAULT
		{
			$$ = makeDefElem($1, nil)
		}
	;

/*
 * Ideally we'd use ColId here, but that causes shift/reduce conflicts against
 * the ALTER DATABASE SET/RESET syntaxes. Instead call out specific keywords
 * we need, and allow IDENT so that database option names don't have to be
 * parser keywords unless they are already keywords for other reasons.
 */
createdb_opt_name:
	IDENT                  { $$ = $1 }
	| CONNECTION LIMIT     { $$ = "connection_limit" }
	| ENCODING             { $$ = "encoding" }
	| LOCATION             { $$ = "location" }
	| OWNER                { $$ = "owner" }
	| TABLESPACE           { $$ = "tablespace" }
	| TEMPLATE             { $$ = "template" }
	;

/*
 * Though the equals sign doesn't match other WITH options, pg_dump uses
 * equals for backward compatibility, and it doesn't seem worth removing it.
 */
opt_equal:
	'='
		{}
	| /* EMPTY */
		{}
	;

/*****************************************************************************
 *
 *      ALTER DATABASE
 *
 *****************************************************************************/

AlterDatabaseStmt:
	ALTER DATABASE name WITH createdb_opt_list
		{
			$$ = &nodes.AlterDatabaseStmt{
				Dbname:  $3,
				Options: $5,
			}
		}
	| ALTER DATABASE name createdb_opt_list
		{
			$$ = &nodes.AlterDatabaseStmt{
				Dbname:  $3,
				Options: $4,
			}
		}
	| ALTER DATABASE name SET TABLESPACE name
		{
			$$ = &nodes.AlterDatabaseStmt{
				Dbname:  $3,
				Options: makeList(makeDefElem("tablespace", &nodes.String{Str: $6})),
			}
		}
	;

AlterDatabaseSetStmt:
	ALTER DATABASE name SetResetClause
		{
			$$ = &nodes.AlterDatabaseSetStmt{
				Dbname:  $3,
				Setstmt: $4.(*nodes.VariableSetStmt),
			}
		}
	;

/*****************************************************************************
 *
 *      DROP DATABASE [ IF EXISTS ] dbname [ [ WITH ] ( options ) ]
 *
 *****************************************************************************/

DropdbStmt:
	DROP DATABASE name
		{
			$$ = &nodes.DropdbStmt{
				Dbname:    $3,
				MissingOk: false,
			}
		}
	| DROP DATABASE IF_P EXISTS name
		{
			$$ = &nodes.DropdbStmt{
				Dbname:    $5,
				MissingOk: true,
			}
		}
	| DROP DATABASE name opt_with '(' drop_option_list ')'
		{
			$$ = &nodes.DropdbStmt{
				Dbname:    $3,
				MissingOk: false,
				Options:   $6,
			}
		}
	| DROP DATABASE IF_P EXISTS name opt_with '(' drop_option_list ')'
		{
			$$ = &nodes.DropdbStmt{
				Dbname:    $5,
				MissingOk: true,
				Options:   $8,
			}
		}
	;

drop_option_list:
	drop_option
		{ $$ = makeList($1) }
	| drop_option_list ',' drop_option
		{ $$ = appendList($1, $3) }
	;

drop_option:
	FORCE
		{
			$$ = makeDefElem("force", nil)
		}
	;

/*****************************************************************************
 *
 *      ALTER SYSTEM
 *
 *****************************************************************************/

AlterSystemStmt:
	ALTER SYSTEM_P SET generic_set
		{
			$$ = &nodes.AlterSystemStmt{
				Setstmt: $4.(*nodes.VariableSetStmt),
			}
		}
	| ALTER SYSTEM_P RESET generic_reset
		{
			$$ = &nodes.AlterSystemStmt{
				Setstmt: $4.(*nodes.VariableSetStmt),
			}
		}
	;

/*****************************************************************************
 *
 *      CREATE SCHEMA
 *
 *****************************************************************************/

CreateSchemaStmt:
	CREATE SCHEMA opt_single_name AUTHORIZATION RoleSpec OptSchemaEltList
		{
			$$ = &nodes.CreateSchemaStmt{
				Schemaname: $3,
				Authrole:   $5.(*nodes.RoleSpec),
				SchemaElts: $6,
			}
		}
	| CREATE SCHEMA ColId OptSchemaEltList
		{
			$$ = &nodes.CreateSchemaStmt{
				Schemaname: $3,
				SchemaElts: $4,
			}
		}
	| CREATE SCHEMA IF_P NOT EXISTS opt_single_name AUTHORIZATION RoleSpec OptSchemaEltList
		{
			$$ = &nodes.CreateSchemaStmt{
				Schemaname:  $6,
				Authrole:    $8.(*nodes.RoleSpec),
				SchemaElts:  $9,
				IfNotExists: true,
			}
		}
	| CREATE SCHEMA IF_P NOT EXISTS ColId OptSchemaEltList
		{
			$$ = &nodes.CreateSchemaStmt{
				Schemaname:  $6,
				SchemaElts:  $7,
				IfNotExists: true,
			}
		}
	;

OptSchemaEltList:
	OptSchemaEltList schema_stmt
		{
			$$ = appendList($1, $2)
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

schema_stmt:
	CreateStmt        { $$ = $1 }
	| IndexStmt       { $$ = $1 }
	| CreateSeqStmt   { $$ = $1 }
	| CreateTrigStmt  { $$ = $1 }
	| ViewStmt        { $$ = $1 }
	| GrantStmt       { $$ = $1 }
	;

/*****************************************************************************
 *
 *      CREATE SEQUENCE / ALTER SEQUENCE
 *
 *****************************************************************************/

CreateSeqStmt:
	CREATE OptTemp SEQUENCE qualified_name OptSeqOptList
		{
			rv := makeRangeVar($4)
			rv.(*nodes.RangeVar).Relpersistence = relpersistenceForTemp($2)
			$$ = &nodes.CreateSeqStmt{
				Sequence: rv.(*nodes.RangeVar),
				Options:  $5,
			}
		}
	| CREATE OptTemp SEQUENCE IF_P NOT EXISTS qualified_name OptSeqOptList
		{
			rv := makeRangeVar($7)
			rv.(*nodes.RangeVar).Relpersistence = relpersistenceForTemp($2)
			$$ = &nodes.CreateSeqStmt{
				Sequence:    rv.(*nodes.RangeVar),
				Options:     $8,
				IfNotExists: true,
			}
		}
	;

AlterSeqStmt:
	ALTER SEQUENCE qualified_name SeqOptList
		{
			rv := makeRangeVar($3)
			$$ = &nodes.AlterSeqStmt{
				Sequence: rv.(*nodes.RangeVar),
				Options:  $4,
			}
		}
	| ALTER SEQUENCE IF_P EXISTS qualified_name SeqOptList
		{
			rv := makeRangeVar($5)
			$$ = &nodes.AlterSeqStmt{
				Sequence:  rv.(*nodes.RangeVar),
				Options:   $6,
				MissingOk: true,
			}
		}
	;

OptParenthesizedSeqOptList:
	'(' SeqOptList ')'	{ $$ = $2 }
	| /* EMPTY */		{ $$ = nil }
	;

OptSeqOptList:
	SeqOptList  { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

SeqOptList:
	SeqOptElem
		{ $$ = makeList($1) }
	| SeqOptList SeqOptElem
		{ $$ = appendList($1, $2) }
	;

SeqOptElem:
	AS SimpleTypename
		{
			$$ = makeDefElem("as", $2)
		}
	| CACHE NumericOnly
		{
			$$ = makeDefElem("cache", $2)
		}
	| CYCLE
		{
			$$ = makeDefElem("cycle", &nodes.Boolean{Boolval: true})
		}
	| NO CYCLE
		{
			$$ = makeDefElem("cycle", &nodes.Boolean{Boolval: false})
		}
	| INCREMENT opt_by NumericOnly
		{
			$$ = makeDefElem("increment", $3)
		}
	| MAXVALUE NumericOnly
		{
			$$ = makeDefElem("maxvalue", $2)
		}
	| MINVALUE NumericOnly
		{
			$$ = makeDefElem("minvalue", $2)
		}
	| NO MAXVALUE
		{
			$$ = makeDefElem("maxvalue", nil)
		}
	| NO MINVALUE
		{
			$$ = makeDefElem("minvalue", nil)
		}
	| OWNED BY any_name
		{
			$$ = makeDefElem("owned_by", $3)
		}
	| SEQUENCE NAME_P any_name
		{
			$$ = makeDefElem("sequence_name", $3)
		}
	| START opt_with NumericOnly
		{
			$$ = makeDefElem("start", $3)
		}
	| RESTART
		{
			$$ = makeDefElem("restart", nil)
		}
	| RESTART opt_with NumericOnly
		{
			$$ = makeDefElem("restart", $3)
		}
	;

opt_by:
	BY     { /* nothing */ }
	| /* EMPTY */ { /* nothing */ }
	;

generated_when:
	ALWAYS      { $$ = int64('a') }
	| BY DEFAULT { $$ = int64('d') }
	;

alter_identity_column_option_list:
	alter_identity_column_option
		{ $$ = makeList($1) }
	| alter_identity_column_option_list alter_identity_column_option
		{ $$ = appendList($1, $2) }
	;

alter_identity_column_option:
	RESTART
		{
			$$ = makeDefElem("restart", nil)
		}
	| RESTART opt_with NumericOnly
		{
			$$ = makeDefElem("restart", $3)
		}
	| SET SeqOptElem
		{
			$$ = $2
		}
	| SET GENERATED generated_when
		{
			$$ = makeDefElem("generated", makeIntConst($3))
		}
	;

/*****************************************************************************
 *
 *      CREATE DOMAIN
 *
 *****************************************************************************/

CreateDomainStmt:
	CREATE DOMAIN_P any_name opt_as Typename opt_column_constraints
		{
			$$ = &nodes.CreateDomainStmt{
				Domainname:  $3,
				Typname:     $5,
				Constraints: $6,
			}
		}
	;

opt_as:
	AS     { /* nothing */ }
	| /* EMPTY */ { /* nothing */ }
	;

/*****************************************************************************
 *
 *      ALTER DOMAIN
 *
 *****************************************************************************/

AlterDomainStmt:
	ALTER DOMAIN_P any_name alter_column_default
		{
			n := $4.(*nodes.AlterDomainStmt)
			n.Typname = $3
			$$ = n
		}
	| ALTER DOMAIN_P any_name DROP NOT NULL_P
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype: 'N',
				Typname: $3,
			}
		}
	| ALTER DOMAIN_P any_name SET NOT NULL_P
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype: 'O',
				Typname: $3,
			}
		}
	| ALTER DOMAIN_P any_name ADD_P DomainConstraint
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype: 'C',
				Typname: $3,
				Def:     $5,
			}
		}
	| ALTER DOMAIN_P any_name DROP CONSTRAINT name opt_drop_behavior
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype:  'X',
				Typname:  $3,
				Name:     $6,
				Behavior: nodes.DropBehavior($7),
			}
		}
	| ALTER DOMAIN_P any_name DROP CONSTRAINT IF_P EXISTS name opt_drop_behavior
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype:   'X',
				Typname:   $3,
				Name:      $8,
				Behavior:  nodes.DropBehavior($9),
				MissingOk: true,
			}
		}
	| ALTER DOMAIN_P any_name VALIDATE CONSTRAINT name
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype: 'V',
				Typname: $3,
				Name:    $6,
			}
		}
	;

alter_column_default:
	SET DEFAULT a_expr
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype: 'T',
				Def:     $3,
			}
		}
	| DROP DEFAULT
		{
			$$ = &nodes.AlterDomainStmt{
				Subtype: 'T',
			}
		}
	;

/*****************************************************************************
 *
 *      ALTER TYPE ... ENUM / ALTER TYPE ... COMPOSITE / ALTER COLLATION
 *
 *****************************************************************************/

AlterEnumStmt:
	ALTER TYPE_P any_name ADD_P VALUE_P opt_if_not_exists Sconst
		{
			$$ = &nodes.AlterEnumStmt{
				Typname:            $3,
				Newval:             $7,
				SkipIfNewvalExists: $6,
			}
		}
	| ALTER TYPE_P any_name ADD_P VALUE_P opt_if_not_exists Sconst BEFORE Sconst
		{
			$$ = &nodes.AlterEnumStmt{
				Typname:            $3,
				Newval:             $7,
				NewvalNeighbor:     $9,
				NewvalIsAfter:      false,
				SkipIfNewvalExists: $6,
			}
		}
	| ALTER TYPE_P any_name ADD_P VALUE_P opt_if_not_exists Sconst AFTER Sconst
		{
			$$ = &nodes.AlterEnumStmt{
				Typname:            $3,
				Newval:             $7,
				NewvalNeighbor:     $9,
				NewvalIsAfter:      true,
				SkipIfNewvalExists: $6,
			}
		}
	| ALTER TYPE_P any_name RENAME VALUE_P Sconst TO Sconst
		{
			$$ = &nodes.AlterEnumStmt{
				Typname: $3,
				Oldval:  $6,
				Newval:  $8,
			}
		}
	;

opt_if_not_exists:
	IF_P NOT EXISTS  { $$ = true }
	| /* EMPTY */    { $$ = false }
	;

AlterCollationStmt:
	ALTER COLLATION any_name REFRESH VERSION_P
		{
			$$ = &nodes.AlterCollationStmt{
				Collname: $3,
			}
		}
	;

AlterCompositeTypeStmt:
	ALTER TYPE_P any_name alter_type_cmds
		{
			rv := &nodes.RangeVar{
				Inh:     true,
				Location: -1,
			}
			/* Convert any_name to schema.rel */
			names := $3
			if names != nil && len(names.Items) > 0 {
				switch len(names.Items) {
				case 1:
					rv.Relname = names.Items[0].(*nodes.String).Str
				case 2:
					rv.Schemaname = names.Items[0].(*nodes.String).Str
					rv.Relname = names.Items[1].(*nodes.String).Str
				case 3:
					rv.Catalogname = names.Items[0].(*nodes.String).Str
					rv.Schemaname = names.Items[1].(*nodes.String).Str
					rv.Relname = names.Items[2].(*nodes.String).Str
				}
			}
			$$ = &nodes.AlterTableStmt{
				Relation: rv,
				Cmds:     $4,
				ObjType:  int(nodes.OBJECT_TYPE),
			}
		}
	;

alter_type_cmds:
	alter_type_cmd
		{ $$ = makeList($1) }
	| alter_type_cmds ',' alter_type_cmd
		{ $$ = appendList($1, $3) }
	;

alter_type_cmd:
	ADD_P ATTRIBUTE TableFuncElement opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_AddColumn),
				Def:      $3,
				Behavior: int($4),
			}
		}
	| DROP ATTRIBUTE ColId opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_DropColumn),
				Name:     $3,
				Behavior: int($4),
			}
		}
	| DROP ATTRIBUTE IF_P EXISTS ColId opt_drop_behavior
		{
			$$ = &nodes.AlterTableCmd{
				Subtype:    int(nodes.AT_DropColumn),
				Name:       $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| ALTER ATTRIBUTE ColId SET DATA_P TYPE_P Typename opt_collate_clause opt_drop_behavior
		{
			coldef := &nodes.ColumnDef{
				Colname:    $3,
				TypeName:   $7,
				CollClause: nil,
			}
			if $8 != nil {
				coldef.CollClause = $8.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_AlterColumnType),
				Name:     $3,
				Def:      coldef,
				Behavior: int($9),
			}
		}
	| ALTER ATTRIBUTE ColId TYPE_P Typename opt_collate_clause opt_drop_behavior
		{
			coldef := &nodes.ColumnDef{
				Colname:    $3,
				TypeName:   $5,
				CollClause: nil,
			}
			if $6 != nil {
				coldef.CollClause = $6.(*nodes.CollateClause)
			}
			$$ = &nodes.AlterTableCmd{
				Subtype:  int(nodes.AT_AlterColumnType),
				Name:     $3,
				Def:      coldef,
				Behavior: int($7),
			}
		}
	;

TableFuncElement:
	ColId Typename opt_collate_clause
		{
			coldef := &nodes.ColumnDef{
				Colname:  $1,
				TypeName: $2,
				IsLocal:  true,
				Location: -1,
			}
			if $3 != nil {
				coldef.CollClause = $3.(*nodes.CollateClause)
			}
			$$ = coldef
		}
	;

opt_collate_clause:
	COLLATE any_name
		{
			$$ = &nodes.CollateClause{
				Collname: $2,
				Location: -1,
			}
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

/*****************************************************************************
 *
 *      DefineStmt - CREATE AGGREGATE/OPERATOR/TYPE/TEXT SEARCH/COLLATION
 *
 *****************************************************************************/

DefineStmt:
	CREATE opt_or_replace AGGREGATE func_name aggr_args definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_AGGREGATE,
				Oldstyle:   false,
				Replace:    $2,
				Defnames:   $4,
				Args:       $5,
				Definition: $6,
			}
		}
	| CREATE opt_or_replace AGGREGATE func_name old_aggr_definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_AGGREGATE,
				Oldstyle:   true,
				Replace:    $2,
				Defnames:   $4,
				Definition: $5,
			}
		}
	| CREATE OPERATOR any_operator definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_OPERATOR,
				Defnames:   $3,
				Definition: $4,
			}
		}
	| CREATE TYPE_P any_name definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_TYPE,
				Defnames:   $3,
				Definition: $4,
			}
		}
	| CREATE TYPE_P any_name
		{
			/* Shell type (identified by lack of definition) */
			$$ = &nodes.DefineStmt{
				Kind:     nodes.OBJECT_TYPE,
				Defnames: $3,
			}
		}
	| CREATE TEXT_P SEARCH PARSER any_name definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_TSPARSER,
				Defnames:   $5,
				Definition: $6,
			}
		}
	| CREATE TEXT_P SEARCH DICTIONARY any_name definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_TSDICTIONARY,
				Defnames:   $5,
				Definition: $6,
			}
		}
	| CREATE TEXT_P SEARCH TEMPLATE any_name definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_TSTEMPLATE,
				Defnames:   $5,
				Definition: $6,
			}
		}
	| CREATE TEXT_P SEARCH CONFIGURATION any_name definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_TSCONFIGURATION,
				Defnames:   $5,
				Definition: $6,
			}
		}
	| CREATE COLLATION any_name definition
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_COLLATION,
				Defnames:   $3,
				Definition: $4,
			}
		}
	| CREATE COLLATION IF_P NOT EXISTS any_name definition
		{
			$$ = &nodes.DefineStmt{
				Kind:        nodes.OBJECT_COLLATION,
				Defnames:    $6,
				Definition:  $7,
				IfNotExists: true,
			}
		}
	| CREATE COLLATION any_name FROM any_name
		{
			$$ = &nodes.DefineStmt{
				Kind:       nodes.OBJECT_COLLATION,
				Defnames:   $3,
				Definition: makeList(makeDefElem("from", $5)),
			}
		}
	| CREATE COLLATION IF_P NOT EXISTS any_name FROM any_name
		{
			$$ = &nodes.DefineStmt{
				Kind:        nodes.OBJECT_COLLATION,
				Defnames:    $6,
				Definition:  makeList(makeDefElem("from", $8)),
				IfNotExists: true,
			}
		}
	;

CompositeTypeStmt:
	CREATE TYPE_P any_name AS '(' OptTableFuncElementList ')'
		{
			$$ = &nodes.CompositeTypeStmt{
				Typevar:    makeRangeVarFromAnyName($3),
				Coldeflist: $6,
			}
		}
	;

CreateEnumStmt:
	CREATE TYPE_P any_name AS ENUM_P '(' opt_enum_val_list ')'
		{
			$$ = &nodes.CreateEnumStmt{
				TypeName: $3,
				Vals:     $7,
			}
		}
	;

CreateRangeStmt:
	CREATE TYPE_P any_name AS RANGE definition
		{
			$$ = &nodes.CreateRangeStmt{
				TypeName: $3,
				Params:   $6,
			}
		}
	;

definition:
	'(' def_list ')'
		{
			$$ = $2
		}
	;

def_list:
	def_elem
		{
			$$ = makeList($1)
		}
	| def_list ',' def_elem
		{
			$$ = appendList($1, $3)
		}
	;

def_elem:
	ColLabel '=' def_arg
		{
			$$ = makeDefElem($1, $3)
		}
	| ColLabel
		{
			$$ = makeDefElem($1, nil)
		}
	;

/* Note: any simple identifier will be returned as a type name! */
def_arg:
	func_type
		{
			$$ = $1
		}
	| reserved_keyword
		{
			$$ = &nodes.String{Str: $1}
		}
	| qual_all_Op
		{
			$$ = $1
		}
	| NumericOnly
		{
			$$ = $1
		}
	| Sconst
		{
			$$ = &nodes.String{Str: $1}
		}
	| NONE
		{
			$$ = &nodes.String{Str: "none"}
		}
	;

old_aggr_definition:
	'(' old_aggr_list ')'
		{
			$$ = $2
		}
	;

old_aggr_list:
	old_aggr_elem
		{
			$$ = makeList($1)
		}
	| old_aggr_list ',' old_aggr_elem
		{
			$$ = appendList($1, $3)
		}
	;

/*
 * Must use IDENT here to avoid reduce/reduce conflicts; fortunately none of
 * the item names needed in old aggregate definitions are likely to become
 * SQL keywords.
 */
old_aggr_elem:
	IDENT '=' def_arg
		{
			$$ = makeDefElem($1, $3)
		}
	;

opt_enum_val_list:
	enum_val_list
		{
			$$ = $1
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

enum_val_list:
	Sconst
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| enum_val_list ',' Sconst
		{
			$$ = appendList($1, &nodes.String{Str: $3})
		}
	;

OptTableFuncElementList:
	TableFuncElementList
		{
			$$ = $1
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

TableFuncElementList:
	TableFuncElement
		{
			$$ = makeList($1)
		}
	| TableFuncElementList ',' TableFuncElement
		{
			$$ = appendList($1, $3)
		}
	;

aggr_args:
	'(' '*' ')'
		{
			/* agg(*) - returns 2-element list: [nil, Integer{-1}] */
			$$ = makeList2(nil, &nodes.Integer{Ival: -1})
		}
	| '(' aggr_args_list ')'
		{
			/* normal args - returns 2-element list: [args, Integer{-1}] */
			$$ = makeList2($2, &nodes.Integer{Ival: -1})
		}
	| '(' ORDER BY aggr_args_list ')'
		{
			/* ordered-set agg with no direct args - returns 2-element list: [args, Integer{0}] */
			$$ = makeList2($4, &nodes.Integer{Ival: 0})
		}
	| '(' aggr_args_list ORDER BY aggr_args_list ')'
		{
			/* ordered-set agg with direct args and ordered args */
			$$ = makeOrderedSetArgs($2, $5)
		}
	;

aggr_args_list:
	aggr_arg
		{
			$$ = makeList($1)
		}
	| aggr_args_list ',' aggr_arg
		{
			$$ = appendList($1, $3)
		}
	;

aggr_arg:
	func_arg
		{
			$$ = $1
		}
	;

/*****************************************************************************
 *
 *      any_operator, all_Op, MathOp, qual_all_Op
 *
 *****************************************************************************/

any_operator:
	all_Op
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| ColId '.' any_operator
		{
			$$ = prependList(&nodes.String{Str: $1}, $3)
		}
	;

all_Op:
	Op      { $$ = $1 }
	| MathOp { $$ = $1 }
	;

MathOp:
	'+'             { $$ = "+" }
	| '-'           { $$ = "-" }
	| '*'           { $$ = "*" }
	| '/'           { $$ = "/" }
	| '%'           { $$ = "%" }
	| '^'           { $$ = "^" }
	| '<'           { $$ = "<" }
	| '>'           { $$ = ">" }
	| '='           { $$ = "=" }
	| LESS_EQUALS   { $$ = "<=" }
	| GREATER_EQUALS { $$ = ">=" }
	| NOT_EQUALS    { $$ = "<>" }
	;

qual_all_Op:
	all_Op
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| OPERATOR '(' any_operator ')'
		{
			$$ = $3
		}
	;

qual_Op:
	Op
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| OPERATOR '(' any_operator ')'
		{
			$$ = $3
		}
	;

/*****************************************************************************
 *
 *      RoleId and role_list
 *
 *****************************************************************************/

RoleId:
	RoleSpec
		{
			spc := $1.(*nodes.RoleSpec)
			if spc.Roletype != int(nodes.ROLESPEC_CSTRING) {
				pglex.Error("role name cannot be a reserved keyword here")
			}
			$$ = spc.Rolename
		}
	;

role_list:
	RoleSpec
		{ $$ = makeList($1) }
	| role_list ',' RoleSpec
		{ $$ = appendList($1, $3) }
	;

// SELECT statement
SelectStmt:
	select_no_parens %prec UMINUS
		{
			$$ = $1
		}
	| select_with_parens %prec UMINUS
		{
			$$ = $1
		}
	;

select_with_parens:
	'(' select_no_parens ')'
		{
			$$ = $2
		}
	| '(' select_with_parens ')'
		{
			$$ = $2
		}
	;

select_no_parens:
	simple_select
		{
			$$ = $1
		}
	| select_clause sort_clause
		{
			n := $1.(*nodes.SelectStmt)
			n.SortClause = $2
			$$ = n
		}
	| select_clause opt_sort_clause select_limit
		{
			n := $1.(*nodes.SelectStmt)
			insertSelectOptions(n, $2, nil, $3, nil)
			$$ = n
		}
	| select_clause opt_sort_clause for_locking_clause opt_select_limit
		{
			n := $1.(*nodes.SelectStmt)
			insertSelectOptions(n, $2, $3, $4, nil)
			$$ = n
		}
	| select_clause opt_sort_clause select_limit opt_for_locking_clause
		{
			n := $1.(*nodes.SelectStmt)
			insertSelectOptions(n, $2, $4, $3, nil)
			$$ = n
		}
	| with_clause select_clause
		{
			n := $2.(*nodes.SelectStmt)
			n.WithClause = $1.(*nodes.WithClause)
			$$ = n
		}
	| with_clause select_clause sort_clause
		{
			n := $2.(*nodes.SelectStmt)
			n.WithClause = $1.(*nodes.WithClause)
			n.SortClause = $3
			$$ = n
		}
	| with_clause select_clause opt_sort_clause select_limit
		{
			n := $2.(*nodes.SelectStmt)
			insertSelectOptions(n, $3, nil, $4, $1.(*nodes.WithClause))
			$$ = n
		}
	| with_clause select_clause opt_sort_clause for_locking_clause opt_select_limit
		{
			n := $2.(*nodes.SelectStmt)
			insertSelectOptions(n, $3, $4, $5, $1.(*nodes.WithClause))
			$$ = n
		}
	| with_clause select_clause opt_sort_clause select_limit opt_for_locking_clause
		{
			n := $2.(*nodes.SelectStmt)
			insertSelectOptions(n, $3, $5, $4, $1.(*nodes.WithClause))
			$$ = n
		}
	;

select_clause:
	simple_select
		{
			$$ = $1
		}
	| select_with_parens
		{
			$$ = $1
		}
	;

simple_select:
	SELECT opt_all_clause opt_target_list into_clause from_clause where_clause group_clause having_clause window_clause
		{
			n := &nodes.SelectStmt{
				TargetList: $3,
			}
			if $4 != nil {
				n.IntoClause = $4.(*nodes.IntoClause)
			}
			if $5 != nil {
				n.FromClause = $5
			}
			if $6 != nil {
				n.WhereClause = $6
			}
			if $7 != nil {
				n.GroupClause = $7
			}
			if $8 != nil {
				n.HavingClause = $8
			}
			if $9 != nil {
				n.WindowClause = $9
			}
			$$ = n
		}
	| SELECT distinct_clause target_list into_clause from_clause where_clause group_clause having_clause window_clause
		{
			n := &nodes.SelectStmt{
				TargetList: $3,
			}
			if $4 != nil {
				n.IntoClause = $4.(*nodes.IntoClause)
			}
			if $5 != nil {
				n.FromClause = $5
			}
			if $6 != nil {
				n.WhereClause = $6
			}
			if $7 != nil {
				n.GroupClause = $7
			}
			if $8 != nil {
				n.HavingClause = $8
			}
			if $9 != nil {
				n.WindowClause = $9
			}
			// TODO: handle DISTINCT clause
			$$ = n
		}
	| select_clause UNION set_quantifier select_clause
		{
			$$ = makeSetOp(nodes.SETOP_UNION, $3, $1, $4)
		}
	| select_clause INTERSECT set_quantifier select_clause
		{
			$$ = makeSetOp(nodes.SETOP_INTERSECT, $3, $1, $4)
		}
	| select_clause EXCEPT set_quantifier select_clause
		{
			$$ = makeSetOp(nodes.SETOP_EXCEPT, $3, $1, $4)
		}
	| values_clause
		{
			$$ = $1
		}
	| TABLE relation_expr
		{
			/* same as SELECT * FROM relation_expr */
			cr := &nodes.ColumnRef{
				Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
				Location: -1,
			}
			rt := &nodes.ResTarget{
				Val:      cr,
				Location: -1,
			}
			$$ = &nodes.SelectStmt{
				TargetList: makeList(rt),
				FromClause: makeList($2),
			}
		}
	;

values_clause:
	VALUES '(' expr_list ')'
		{
			n := &nodes.SelectStmt{}
			n.ValuesLists = &nodes.List{Items: []nodes.Node{$3}}
			$$ = n
		}
	| values_clause ',' '(' expr_list ')'
		{
			n := $1.(*nodes.SelectStmt)
			n.ValuesLists.Items = append(n.ValuesLists.Items, $4)
			$$ = n
		}
	;

opt_all_clause:
	ALL { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

opt_distinct_clause:
	distinct_clause { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

distinct_clause:
	DISTINCT
	| DISTINCT ON '(' expr_list ')'
	;

set_quantifier:
	ALL { $$ = true }
	| DISTINCT { $$ = false }
	| /* EMPTY */ { $$ = false }
	;

/*****************************************************************************
 *
 *      WITH clause (Common Table Expressions)
 *
 *****************************************************************************/

opt_with_clause:
	with_clause { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

with_clause:
	WITH cte_list
		{
			$$ = &nodes.WithClause{
				Ctes:      $2,
				Recursive: false,
			}
		}
	| WITH RECURSIVE cte_list
		{
			$$ = &nodes.WithClause{
				Ctes:      $3,
				Recursive: true,
			}
		}
	;

cte_list:
	common_table_expr
		{
			$$ = makeList($1)
		}
	| cte_list ',' common_table_expr
		{
			$$ = appendList($1, $3)
		}
	;

common_table_expr:
	name opt_name_list AS opt_materialized '(' PreparableStmt ')' opt_search_clause opt_cycle_clause
		{
			cte := &nodes.CommonTableExpr{
				Ctename:         $1,
				Aliascolnames:   $2,
				Ctematerialized: int($4),
				Ctequery:        $6,
			}
			if $8 != nil {
				cte.SearchClause = $8
			}
			if $9 != nil {
				cte.CycleClause = $9
			}
			$$ = cte
		}
	| name '(' name_list ')' AS opt_materialized '(' PreparableStmt ')' opt_search_clause opt_cycle_clause
		{
			cte := &nodes.CommonTableExpr{
				Ctename:         $1,
				Aliascolnames:   $3,
				Ctematerialized: int($6),
				Ctequery:        $8,
			}
			if $10 != nil {
				cte.SearchClause = $10
			}
			if $11 != nil {
				cte.CycleClause = $11
			}
			$$ = cte
		}
	;

opt_materialized:
	MATERIALIZED			{ $$ = int64(nodes.CTEMaterializeAlways) }
	| NOT MATERIALIZED		{ $$ = int64(nodes.CTEMaterializeNever) }
	| /* EMPTY */			{ $$ = int64(nodes.CTEMaterializeDefault) }
	;

opt_search_clause:
	SEARCH DEPTH FIRST_P BY columnList SET ColId
		{
			$$ = &nodes.CTESearchClause{
				SearchColList:      $5,
				SearchBreadthFirst: false,
				SearchSeqColumn:    $7,
				Location:           -1,
			}
		}
	| SEARCH BREADTH FIRST_P BY columnList SET ColId
		{
			$$ = &nodes.CTESearchClause{
				SearchColList:      $5,
				SearchBreadthFirst: true,
				SearchSeqColumn:    $7,
				Location:           -1,
			}
		}
	| /* EMPTY */ { $$ = nil }
	;

opt_cycle_clause:
	CYCLE columnList SET ColId TO AexprConst DEFAULT AexprConst USING ColId
		{
			$$ = &nodes.CTECycleClause{
				CycleColList:     $2,
				CycleMarkColumn:  $4,
				CycleMarkValue:   $6,
				CycleMarkDefault: $8,
				CyclePathColumn:  $10,
				Location:         -1,
			}
		}
	| CYCLE columnList SET ColId USING ColId
		{
			$$ = &nodes.CTECycleClause{
				CycleColList:     $2,
				CycleMarkColumn:  $4,
				CycleMarkValue:   makeBoolAConst(1),
				CycleMarkDefault: makeBoolAConst(0),
				CyclePathColumn:  $6,
				Location:         -1,
			}
		}
	| /* EMPTY */ { $$ = nil }
	;

opt_target_list:
	target_list { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

into_clause:
	INTO OptTempTableName
		{
			$$ = &nodes.IntoClause{
				Rel:            $2.(*nodes.RangeVar),
				OnCommit:       nodes.ONCOMMIT_NOOP,
			}
		}
	| /* EMPTY */ { $$ = nil }
	;

OptTempTableName:
	TEMPORARY opt_table qualified_name  { rv := makeRangeVar($3); rv.(*nodes.RangeVar).Relpersistence = 't'; $$ = rv }
	| TEMP opt_table qualified_name     { rv := makeRangeVar($3); rv.(*nodes.RangeVar).Relpersistence = 't'; $$ = rv }
	| LOCAL TEMPORARY opt_table qualified_name  { rv := makeRangeVar($4); rv.(*nodes.RangeVar).Relpersistence = 't'; $$ = rv }
	| LOCAL TEMP opt_table qualified_name  { rv := makeRangeVar($4); rv.(*nodes.RangeVar).Relpersistence = 't'; $$ = rv }
	| GLOBAL TEMPORARY opt_table qualified_name  { rv := makeRangeVar($4); rv.(*nodes.RangeVar).Relpersistence = 't'; $$ = rv }
	| GLOBAL TEMP opt_table qualified_name  { rv := makeRangeVar($4); rv.(*nodes.RangeVar).Relpersistence = 't'; $$ = rv }
	| UNLOGGED opt_table qualified_name  { rv := makeRangeVar($3); rv.(*nodes.RangeVar).Relpersistence = 'u'; $$ = rv }
	| TABLE qualified_name  { $$ = makeRangeVar($2) }
	| qualified_name        { $$ = makeRangeVar($1) }
	;

target_list:
	target_el
		{
			$$ = makeList($1)
		}
	| target_list ',' target_el
		{
			$$ = appendList($1, $3)
		}
	;

target_el:
	a_expr AS ColLabel
		{
			$$ = &nodes.ResTarget{
				Name: $3,
				Val:  $1,
			}
		}
	| a_expr IDENT
		{
			$$ = &nodes.ResTarget{
				Name: $2,
				Val:  $1,
			}
		}
	| a_expr
		{
			$$ = &nodes.ResTarget{
				Val: $1,
			}
		}
	| '*'
		{
			$$ = &nodes.ResTarget{
				Val: &nodes.ColumnRef{
					Fields: &nodes.List{Items: []nodes.Node{&nodes.A_Star{}}},
				},
			}
		}
	;

// FROM clause
from_clause:
	FROM from_list { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

from_list:
	table_ref
		{
			$$ = makeList($1)
		}
	| from_list ',' table_ref
		{
			$$ = appendList($1, $3)
		}
	;

table_ref:
	relation_expr opt_alias_clause
		{
			rv := $1.(*nodes.RangeVar)
			if $2 != nil {
				rv.Alias = $2.(*nodes.Alias)
			}
			$$ = rv
		}
	| select_with_parens opt_alias_clause
		{
			n := &nodes.RangeSubselect{
				Subquery: $1,
			}
			if $2 != nil {
				n.Alias = $2.(*nodes.Alias)
			}
			$$ = n
		}
	| joined_table
		{
			$$ = $1
		}
	| relation_expr opt_alias_clause tablesample_clause
		{
			rv := $1.(*nodes.RangeVar)
			if $2 != nil {
				rv.Alias = $2.(*nodes.Alias)
			}
			n := $3.(*nodes.RangeTableSample)
			n.Relation = rv
			$$ = n
		}
	| func_table func_alias_clause
		{
			n := $1.(*nodes.RangeFunction)
			setFuncAlias(n, $2)
			$$ = n
		}
	| LATERAL_P func_table func_alias_clause
		{
			n := $2.(*nodes.RangeFunction)
			n.Lateral = true
			setFuncAlias(n, $3)
			$$ = n
		}
	| LATERAL_P select_with_parens opt_alias_clause
		{
			n := &nodes.RangeSubselect{
				Lateral:  true,
				Subquery: $2,
			}
			if $3 != nil {
				n.Alias = $3.(*nodes.Alias)
			}
			$$ = n
		}
	| '(' joined_table ')' opt_alias_clause
		{
			j := $2.(*nodes.JoinExpr)
			if $4 != nil {
				j.Alias = $4.(*nodes.Alias)
			}
			$$ = j
		}
	| xmltable opt_alias_clause
		{
			n := $1.(*nodes.RangeTableFunc)
			if $2 != nil {
				n.Alias = $2.(*nodes.Alias)
			}
			$$ = n
		}
	| LATERAL_P xmltable opt_alias_clause
		{
			n := $2.(*nodes.RangeTableFunc)
			n.Lateral = true
			if $3 != nil {
				n.Alias = $3.(*nodes.Alias)
			}
			$$ = n
		}
	| json_table opt_alias_clause
		{
			n := $1.(*nodes.JsonTable)
			if $2 != nil {
				n.Alias = $2.(*nodes.Alias)
			}
			$$ = n
		}
	| LATERAL_P json_table opt_alias_clause
		{
			n := $2.(*nodes.JsonTable)
			n.Lateral = true
			if $3 != nil {
				n.Alias = $3.(*nodes.Alias)
			}
			$$ = n
		}
	;

joined_table:
	table_ref CROSS JOIN table_ref
		{
			$$ = &nodes.JoinExpr{
				Jointype:  nodes.JOIN_INNER,
				IsNatural: false,
				Larg:      $1,
				Rarg:      $4,
			}
		}
	| table_ref join_type JOIN table_ref join_qual
		{
			n := &nodes.JoinExpr{
				Jointype:  nodes.JoinType($2),
				IsNatural: false,
				Larg:      $1,
				Rarg:      $4,
			}
			setJoinQual(n, $5)
			$$ = n
		}
	| table_ref JOIN table_ref join_qual
		{
			n := &nodes.JoinExpr{
				Jointype:  nodes.JOIN_INNER,
				IsNatural: false,
				Larg:      $1,
				Rarg:      $3,
			}
			setJoinQual(n, $4)
			$$ = n
		}
	| table_ref NATURAL join_type JOIN table_ref
		{
			$$ = &nodes.JoinExpr{
				Jointype:  nodes.JoinType($3),
				IsNatural: true,
				Larg:      $1,
				Rarg:      $5,
			}
		}
	| table_ref NATURAL JOIN table_ref
		{
			$$ = &nodes.JoinExpr{
				Jointype:  nodes.JOIN_INNER,
				IsNatural: true,
				Larg:      $1,
				Rarg:      $4,
			}
		}
	;

join_type:
	FULL opt_outer      { $$ = int64(nodes.JOIN_FULL) }
	| LEFT opt_outer    { $$ = int64(nodes.JOIN_LEFT) }
	| RIGHT opt_outer   { $$ = int64(nodes.JOIN_RIGHT) }
	| INNER_P           { $$ = int64(nodes.JOIN_INNER) }
	;

opt_outer:
	OUTER_P
		{ }
	| /* EMPTY */
		{ }
	;

join_qual:
	USING '(' name_list ')' join_using_alias
		{
			/* Wrap USING clause info in a List: [nameList, alias?] */
			if $5 != nil {
				$$ = makeList2($3, $5)
			} else {
				$$ = $3
			}
		}
	| ON a_expr
		{
			$$ = $2
		}
	;

join_using_alias:
	AS ColId
		{
			$$ = &nodes.Alias{Aliasname: $2}
		}
	| /* EMPTY */ { $$ = nil }
	;

relation_expr:
	qualified_name
		{
			$$ = makeRangeVar($1)
		}
	| qualified_name '*'
		{
			rv := makeRangeVar($1)
			rv.(*nodes.RangeVar).Inh = true
			$$ = rv
		}
	| ONLY qualified_name
		{
			rv := makeRangeVar($2)
			rv.(*nodes.RangeVar).Inh = false
			$$ = rv
		}
	;

opt_alias_clause:
	alias_clause { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

func_alias_clause:
	alias_clause { $$ = $1 }
	| AS '(' TableFuncElementList ')'
		{
			$$ = &nodes.List{Items: []nodes.Node{nil, $3}}
		}
	| AS ColId '(' TableFuncElementList ')'
		{
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Alias{Aliasname: $2},
				$4,
			}}
		}
	| ColId '(' TableFuncElementList ')'
		{
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Alias{Aliasname: $1},
				$3,
			}}
		}
	| /* EMPTY */ { $$ = nil }
	;

alias_clause:
	AS ColId '(' name_list ')'
		{
			$$ = &nodes.Alias{Aliasname: $2, Colnames: $4}
		}
	| AS ColId
		{
			$$ = &nodes.Alias{Aliasname: $2}
		}
	| ColId '(' name_list ')'
		{
			$$ = &nodes.Alias{Aliasname: $1, Colnames: $3}
		}
	| ColId
		{
			$$ = &nodes.Alias{Aliasname: $1}
		}
	;

func_table:
	func_expr_windowless opt_ordinality
		{
			$$ = &nodes.RangeFunction{
				Ordinality: $2,
				Functions:  makeList(makeList($1)),
			}
		}
	| ROWS FROM '(' rowsfrom_list ')' opt_ordinality
		{
			$$ = &nodes.RangeFunction{
				IsRowsfrom: true,
				Ordinality: $6,
				Functions:  $4,
			}
		}
	;

rowsfrom_list:
	rowsfrom_item
		{ $$ = makeList($1) }
	| rowsfrom_list ',' rowsfrom_item
		{ $$ = appendList($1, $3) }
	;

rowsfrom_item:
	func_expr_windowless opt_col_def_list
		{
			$$ = makeList2($1, $2)
		}
	;

opt_col_def_list:
	AS '(' TableFuncElementList ')'  { $$ = $3 }
	| /* EMPTY */                    { $$ = nil }
	;

opt_ordinality:
	WITH ORDINALITY		{ $$ = true }
	| /* EMPTY */		{ $$ = false }
	;

tablesample_clause:
	TABLESAMPLE func_name '(' expr_list ')' opt_repeatable_clause
		{
			$$ = &nodes.RangeTableSample{
				Method: $2,
				Args:   $4,
				Repeatable: $6,
				Location: -1,
			}
		}
	;

opt_repeatable_clause:
	REPEATABLE '(' a_expr ')'	{ $$ = $3 }
	| /* EMPTY */				{ $$ = nil }
	;

// WHERE clause
where_clause:
	WHERE a_expr { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

where_or_current_clause:
	WHERE a_expr { $$ = $2 }
	| WHERE CURRENT_P OF cursor_name
		{
			$$ = &nodes.CurrentOfExpr{
				CursorName: $4,
			}
		}
	| /* EMPTY */ { $$ = nil }
	;

// GROUP BY clause
group_clause:
	GROUP_P BY set_quantifier group_by_list { $$ = $4 }
	| /* EMPTY */ { $$ = nil }
	;

group_by_list:
	group_by_item
		{ $$ = makeList($1) }
	| group_by_list ',' group_by_item
		{ $$ = appendList($1, $3) }
	;

group_by_item:
	a_expr { $$ = $1 }
	| empty_grouping_set { $$ = $1 }
	| cube_clause { $$ = $1 }
	| rollup_clause { $$ = $1 }
	| grouping_sets_clause { $$ = $1 }
	;

empty_grouping_set:
	'(' ')'
		{
			$$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_EMPTY, Location: -1}
		}
	;

cube_clause:
	CUBE '(' expr_list ')'
		{
			$$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_CUBE, Content: $3, Location: -1}
		}
	;

rollup_clause:
	ROLLUP '(' expr_list ')'
		{
			$$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_ROLLUP, Content: $3, Location: -1}
		}
	;

grouping_sets_clause:
	GROUPING SETS '(' group_by_list ')'
		{
			$$ = &nodes.GroupingSet{Kind: nodes.GROUPING_SET_SETS, Content: $4, Location: -1}
		}
	;

// HAVING clause
having_clause:
	HAVING a_expr { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

// Sort clause
sort_clause:
	ORDER BY sortby_list { $$ = $3 }
	;

opt_sort_clause:
	sort_clause { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

sortby_list:
	sortby
		{
			$$ = makeList($1)
		}
	| sortby_list ',' sortby
		{
			$$ = appendList($1, $3)
		}
	;

sortby:
	a_expr USING qual_all_Op opt_nulls_order
		{
			$$ = &nodes.SortBy{
				Node:        $1,
				SortbyDir:   nodes.SORTBY_USING,
				SortbyNulls: nodes.SortByNulls($4),
				UseOp:       $3,
			}
		}
	| a_expr opt_asc_desc opt_nulls_order
		{
			$$ = &nodes.SortBy{
				Node:        $1,
				SortbyDir:   nodes.SortByDir($2),
				SortbyNulls: nodes.SortByNulls($3),
			}
		}
	;

opt_asc_desc:
	ASC { $$ = int64(nodes.SORTBY_ASC) }
	| DESC { $$ = int64(nodes.SORTBY_DESC) }
	| /* EMPTY */ { $$ = int64(nodes.SORTBY_DEFAULT) }
	;

select_limit:
	limit_clause offset_clause
		{
			$$ = $1
			$$.LimitOffset = $2
		}
	| offset_clause limit_clause
		{
			$$ = $2
			$$.LimitOffset = $1
		}
	| limit_clause
		{
			$$ = $1
		}
	| offset_clause
		{
			$$ = &SelectLimit{
				LimitOffset: $1,
				LimitOption: nodes.LIMIT_OPTION_COUNT,
			}
		}
	;

limit_clause:
	LIMIT select_limit_value
		{
			$$ = &SelectLimit{
				LimitCount:  $2,
				LimitOption: nodes.LIMIT_OPTION_COUNT,
			}
		}
	| LIMIT select_limit_value ',' select_offset_value
		{
			/* PostgreSQL disallows this syntax with an error, but we parse it.
			 * The LIMIT #,# syntax is deprecated. */
			$$ = &SelectLimit{
				LimitOffset: $2,
				LimitCount:  $4,
				LimitOption: nodes.LIMIT_OPTION_COUNT,
			}
		}
	| FETCH first_or_next select_fetch_first_value row_or_rows ONLY
		{
			$$ = &SelectLimit{
				LimitCount:  $3,
				LimitOption: nodes.LIMIT_OPTION_COUNT,
			}
		}
	| FETCH first_or_next select_fetch_first_value row_or_rows WITH TIES
		{
			$$ = &SelectLimit{
				LimitCount:  $3,
				LimitOption: nodes.LIMIT_OPTION_WITH_TIES,
			}
		}
	| FETCH first_or_next row_or_rows ONLY
		{
			$$ = &SelectLimit{
				LimitCount:  makeIntConst(1),
				LimitOption: nodes.LIMIT_OPTION_COUNT,
			}
		}
	| FETCH first_or_next row_or_rows WITH TIES
		{
			$$ = &SelectLimit{
				LimitCount:  makeIntConst(1),
				LimitOption: nodes.LIMIT_OPTION_WITH_TIES,
			}
		}
	;

offset_clause:
	OFFSET select_offset_value
		{ $$ = $2 }
	| OFFSET select_fetch_first_value row_or_rows
		{ $$ = $2 }
	;

first_or_next:
	FIRST_P
	| NEXT
	;

row_or_rows:
	ROW
	| ROWS
	;

select_limit_value:
	a_expr { $$ = $1 }
	| ALL
		{
			/* LIMIT ALL is represented as a NULL constant */
			$$ = &nodes.A_Const{Isnull: true}
		}
	;

select_offset_value:
	a_expr { $$ = $1 }
	;

select_fetch_first_value:
	c_expr { $$ = $1 }
	| '+' c_expr
		{
			$$ = $2
		}
	| '-' c_expr
		{
			$$ = doNegate($2)
		}
	;

opt_select_limit:
	select_limit		{ $$ = $1 }
	| /* EMPTY */		{ $$ = nil }
	;

/*
 * Locking clause (FOR UPDATE/SHARE etc)
 */
for_locking_clause:
	for_locking_items		{ $$ = $1 }
	| FOR READ ONLY			{ $$ = nil }
	;

opt_for_locking_clause:
	for_locking_clause		{ $$ = $1 }
	| /* EMPTY */			{ $$ = nil }
	;

for_locking_items:
	for_locking_item						{ $$ = makeList($1) }
	| for_locking_items for_locking_item	{ $$ = appendList($1, $2) }
	;

for_locking_item:
	for_locking_strength locked_rels_list opt_nowait_or_skip
		{
			$$ = &nodes.LockingClause{
				LockedRels: $2,
				Strength:   int($1),
				WaitPolicy: int($3),
			}
		}
	;

for_locking_strength:
	FOR UPDATE				{ $$ = int64(nodes.LCS_FORUPDATE) }
	| FOR NO KEY UPDATE		{ $$ = int64(nodes.LCS_FORNOKEYUPDATE) }
	| FOR SHARE				{ $$ = int64(nodes.LCS_FORSHARE) }
	| FOR KEY SHARE			{ $$ = int64(nodes.LCS_FORKEYSHARE) }
	;

locked_rels_list:
	OF qualified_name_list	{ $$ = $2 }
	| /* EMPTY */			{ $$ = nil }
	;

opt_nowait_or_skip:
	NOWAIT					{ $$ = int64(nodes.LockWaitError) }
	| SKIP LOCKED			{ $$ = int64(nodes.LockWaitSkip) }
	| /* EMPTY */			{ $$ = int64(nodes.LockWaitBlock) }
	;

// Expressions
a_expr:
	c_expr { $$ = $1 }
	| a_expr '+' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "+", $1, $3)
		}
	| a_expr '-' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "-", $1, $3)
		}
	| a_expr '*' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "*", $1, $3)
		}
	| a_expr '/' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "/", $1, $3)
		}
	| a_expr '%' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "%", $1, $3)
		}
	| a_expr '^' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "^", $1, $3)
		}
	| a_expr '<' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "<", $1, $3)
		}
	| a_expr '>' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, ">", $1, $3)
		}
	| a_expr '=' a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "=", $1, $3)
		}
	| a_expr LESS_EQUALS a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "<=", $1, $3)
		}
	| a_expr GREATER_EQUALS a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, ">=", $1, $3)
		}
	| a_expr NOT_EQUALS a_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "<>", $1, $3)
		}
	| a_expr qual_Op a_expr				%prec Op
		{
			$$ = makeAExprFromList(nodes.AEXPR_OP, $2, $1, $3)
		}
	| qual_Op a_expr					%prec Op
		{
			$$ = makeAExprFromList(nodes.AEXPR_OP, $1, nil, $2)
		}
	| a_expr AND a_expr
		{
			$$ = makeBoolExpr(nodes.AND_EXPR, $1, $3)
		}
	| a_expr OR a_expr
		{
			$$ = makeBoolExpr(nodes.OR_EXPR, $1, $3)
		}
	| NOT a_expr
		{
			$$ = makeBoolExpr(nodes.NOT_EXPR, $2, nil)
		}
	| NOT_LA a_expr							%prec NOT
		{
			$$ = makeBoolExpr(nodes.NOT_EXPR, $2, nil)
		}
	| a_expr IS NULL_P
		{
			$$ = &nodes.NullTest{
				Arg:         $1,
				Nulltesttype: nodes.IS_NULL,
			}
		}
	| a_expr IS NOT NULL_P
		{
			$$ = &nodes.NullTest{
				Arg:         $1,
				Nulltesttype: nodes.IS_NOT_NULL,
			}
		}
	| a_expr IS TRUE_P
		{
			$$ = &nodes.BooleanTest{
				Arg:          $1,
				Booltesttype: nodes.IS_TRUE,
			}
		}
	| a_expr IS FALSE_P
		{
			$$ = &nodes.BooleanTest{
				Arg:          $1,
				Booltesttype: nodes.IS_FALSE,
			}
		}
	| a_expr IS NOT TRUE_P                     %prec IS
		{
			$$ = &nodes.BooleanTest{
				Arg:          $1,
				Booltesttype: nodes.IS_NOT_TRUE,
			}
		}
	| a_expr IS NOT FALSE_P                    %prec IS
		{
			$$ = &nodes.BooleanTest{
				Arg:          $1,
				Booltesttype: nodes.IS_NOT_FALSE,
			}
		}
	| a_expr IS UNKNOWN                        %prec IS
		{
			$$ = &nodes.BooleanTest{
				Arg:          $1,
				Booltesttype: nodes.IS_UNKNOWN,
			}
		}
	| a_expr IS NOT UNKNOWN                    %prec IS
		{
			$$ = &nodes.BooleanTest{
				Arg:          $1,
				Booltesttype: nodes.IS_NOT_UNKNOWN,
			}
		}
	| a_expr ISNULL
		{
			$$ = &nodes.NullTest{
				Arg:          $1,
				Nulltesttype: nodes.IS_NULL,
			}
		}
	| a_expr NOTNULL
		{
			$$ = &nodes.NullTest{
				Arg:          $1,
				Nulltesttype: nodes.IS_NOT_NULL,
			}
		}
	| a_expr IS DISTINCT FROM a_expr           %prec IS
		{
			$$ = makeAExpr(nodes.AEXPR_DISTINCT, "=", $1, $5)
		}
	| a_expr IS NOT DISTINCT FROM a_expr       %prec IS
		{
			$$ = makeAExpr(nodes.AEXPR_NOT_DISTINCT, "=", $1, $6)
		}
	| row OVERLAPS row
		{
			/* convert to a function call */
			var args *nodes.List
			if r1, ok := $1.(*nodes.RowExpr); ok {
				args = r1.Args
			} else {
				args = makeList($1)
			}
			if r2, ok := $3.(*nodes.RowExpr); ok {
				args = concatLists(args, r2.Args)
			} else {
				args = appendList(args, $3)
			}
			$$ = &nodes.FuncCall{
				Funcname:       makeFuncName("overlaps"),
				Args:           args,
				FuncFormat:     int(nodes.COERCE_SQL_SYNTAX),
				Location:       -1,
			}
		}
	| a_expr IS DOCUMENT_P                             %prec IS
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_DOCUMENT,
				Args:     makeList($1),
				Location: -1,
			}
		}
	| a_expr IS NOT DOCUMENT_P                         %prec IS
		{
			$$ = makeNotExpr(&nodes.XmlExpr{
				Op:       nodes.IS_DOCUMENT,
				Args:     makeList($1),
				Location: -1,
			})
		}
	| a_expr IS json_predicate_type_constraint json_key_uniqueness_constraint_opt   %prec IS
		{
			$$ = &nodes.JsonIsPredicate{
				Expr:       $1,
				ItemType:   nodes.JsonValueType($3),
				UniqueKeys: $4 != 0,
				Location:   -1,
			}
		}
	| a_expr IS NOT json_predicate_type_constraint json_key_uniqueness_constraint_opt   %prec IS
		{
			$$ = makeNotExpr(&nodes.JsonIsPredicate{
				Expr:       $1,
				ItemType:   nodes.JsonValueType($4),
				UniqueKeys: $5 != 0,
				Location:   -1,
			})
		}
	| a_expr IS unicode_normal_form NORMALIZED                    %prec IS
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "is_normalized"),
				Args:       makeList2($1, makeStringConst($3)),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| a_expr IS NOT unicode_normal_form NORMALIZED                %prec IS
		{
			$$ = makeNotExpr(&nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "is_normalized"),
				Args:       makeList2($1, makeStringConst($4)),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			})
		}
	| a_expr IS NORMALIZED                                        %prec IS
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "is_normalized"),
				Args:       makeList2($1, makeStringConst("NFC")),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| a_expr IS NOT NORMALIZED                                    %prec IS
		{
			$$ = makeNotExpr(&nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "is_normalized"),
				Args:       makeList2($1, makeStringConst("NFC")),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			})
		}
	| a_expr LIKE a_expr                              %prec LIKE
		{
			$$ = makeAExpr(nodes.AEXPR_LIKE, "~~", $1, $3)
		}
	| a_expr LIKE a_expr ESCAPE a_expr                 %prec LIKE
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "like_escape"),
				Args:       makeList2($3, $5),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_LIKE, "~~", $1, esc)
		}
	| a_expr NOT_LA LIKE a_expr                        %prec NOT_LA
		{
			$$ = makeAExpr(nodes.AEXPR_LIKE, "!~~", $1, $4)
		}
	| a_expr NOT_LA LIKE a_expr ESCAPE a_expr          %prec NOT_LA
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "like_escape"),
				Args:       makeList2($4, $6),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_LIKE, "!~~", $1, esc)
		}
	| a_expr ILIKE a_expr                              %prec ILIKE
		{
			$$ = makeAExpr(nodes.AEXPR_ILIKE, "~~*", $1, $3)
		}
	| a_expr ILIKE a_expr ESCAPE a_expr                %prec ILIKE
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "like_escape"),
				Args:       makeList2($3, $5),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_ILIKE, "~~*", $1, esc)
		}
	| a_expr NOT_LA ILIKE a_expr                       %prec NOT_LA
		{
			$$ = makeAExpr(nodes.AEXPR_ILIKE, "!~~*", $1, $4)
		}
	| a_expr NOT_LA ILIKE a_expr ESCAPE a_expr         %prec NOT_LA
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "like_escape"),
				Args:       makeList2($4, $6),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_ILIKE, "!~~*", $1, esc)
		}
	| a_expr SIMILAR TO a_expr                         %prec SIMILAR
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
				Args:       makeList($4),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_SIMILAR, "~", $1, esc)
		}
	| a_expr SIMILAR TO a_expr ESCAPE a_expr           %prec SIMILAR
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
				Args:       makeList2($4, $6),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_SIMILAR, "~", $1, esc)
		}
	| a_expr NOT_LA SIMILAR TO a_expr                  %prec NOT_LA
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
				Args:       makeList($5),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_SIMILAR, "!~", $1, esc)
		}
	| a_expr NOT_LA SIMILAR TO a_expr ESCAPE a_expr    %prec NOT_LA
		{
			esc := &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "similar_to_escape"),
				Args:       makeList2($5, $7),
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
			$$ = makeAExpr(nodes.AEXPR_SIMILAR, "!~", $1, esc)
		}
	| a_expr BETWEEN opt_asymmetric b_expr AND a_expr  %prec BETWEEN
		{
			$$ = makeAExpr(nodes.AEXPR_BETWEEN, "BETWEEN", $1,
				&nodes.List{Items: []nodes.Node{$4, $6}})
		}
	| a_expr NOT_LA BETWEEN opt_asymmetric b_expr AND a_expr %prec NOT_LA
		{
			$$ = makeAExpr(nodes.AEXPR_NOT_BETWEEN, "NOT BETWEEN", $1,
				&nodes.List{Items: []nodes.Node{$5, $7}})
		}
	| a_expr BETWEEN SYMMETRIC b_expr AND a_expr       %prec BETWEEN
		{
			$$ = makeAExpr(nodes.AEXPR_BETWEEN_SYM, "BETWEEN SYMMETRIC", $1,
				&nodes.List{Items: []nodes.Node{$4, $6}})
		}
	| a_expr NOT_LA BETWEEN SYMMETRIC b_expr AND a_expr %prec NOT_LA
		{
			$$ = makeAExpr(nodes.AEXPR_NOT_BETWEEN_SYM, "NOT BETWEEN SYMMETRIC", $1,
				&nodes.List{Items: []nodes.Node{$5, $7}})
		}
	| a_expr IN_P '(' expr_list ')'
		{
			$$ = makeAExpr(nodes.AEXPR_IN, "=", $1, makeListNode($4))
		}
	| a_expr NOT_LA IN_P '(' expr_list ')'             %prec NOT_LA
		{
			$$ = makeAExpr(nodes.AEXPR_IN, "<>", $1, makeListNode($5))
		}
	| a_expr IN_P select_with_parens
		{
			$$ = &nodes.SubLink{
				SubLinkType: int(nodes.ANY_SUBLINK),
				Testexpr:    $1,
				Subselect:   $3,
				Location:    -1,
			}
		}
	| a_expr NOT_LA IN_P select_with_parens            %prec NOT_LA
		{
			sublink := &nodes.SubLink{
				SubLinkType: int(nodes.ANY_SUBLINK),
				Testexpr:    $1,
				Subselect:   $4,
				Location:    -1,
			}
			$$ = makeBoolExpr(nodes.NOT_EXPR, sublink, nil)
		}
	| a_expr subquery_Op sub_type select_with_parens %prec Op
		{
			$$ = &nodes.SubLink{
				SubLinkType: int($3),
				Testexpr:    $1,
				OperName:    $2,
				Subselect:   $4,
				Location:    -1,
			}
		}
	| a_expr subquery_Op sub_type '(' a_expr ')' %prec Op
		{
			/* expr op ANY/ALL (expr) — non-subquery form */
			kind := nodes.AEXPR_OP_ANY
			if $3 == int64(nodes.ALL_SUBLINK) {
				kind = nodes.AEXPR_OP_ALL
			}
			$$ = makeAExprFromList(kind, $2, $1, $5)
		}
	| a_expr COLLATE any_name
		{
			$$ = &nodes.CollateClause{
				Arg:      $1,
				Collname: $3,
				Location: -1,
			}
		}
	| a_expr AT TIME ZONE a_expr                       %prec AT
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "timezone"),
				Args:       makeList2($5, $1),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| a_expr AT LOCAL                                  %prec AT
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "timezone"),
				Args:       makeList($1),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| DEFAULT
		{
			$$ = &nodes.SetToDefault{Location: -1}
		}
	| a_expr '[' a_expr ']'
		{
			$$ = &nodes.A_Indirection{
				Arg:         $1,
				Indirection: makeList(&nodes.A_Indices{Uidx: $3}),
			}
		}
	| a_expr '[' opt_slice_bound ':' opt_slice_bound ']'
		{
			$$ = &nodes.A_Indirection{
				Arg: $1,
				Indirection: makeList(&nodes.A_Indices{
					IsSlice: true,
					Lidx:    $3,
					Uidx:    $5,
				}),
			}
		}
	| a_expr TYPECAST Typename
		{
			$$ = &nodes.TypeCast{
				Arg:      $1,
				TypeName: $3,
				Location: -1,
			}
		}
	| '+' a_expr %prec UMINUS
		{
			$$ = $2
		}
	| '-' a_expr %prec UMINUS
		{
			$$ = doNegate($2)
		}
	;

opt_asymmetric:
	ASYMMETRIC  {}
	| /* EMPTY */ {}
	;

sub_type:
	ANY  { $$ = int64(nodes.ANY_SUBLINK) }
	| SOME { $$ = int64(nodes.ANY_SUBLINK) }
	| ALL  { $$ = int64(nodes.ALL_SUBLINK) }
	;

subquery_Op:
	Op
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| LIKE
		{
			$$ = makeList(&nodes.String{Str: "~~"})
		}
	| NOT_LA LIKE
		{
			$$ = makeList(&nodes.String{Str: "!~~"})
		}
	| ILIKE
		{
			$$ = makeList(&nodes.String{Str: "~~*"})
		}
	| NOT_LA ILIKE
		{
			$$ = makeList(&nodes.String{Str: "!~~*"})
		}
	| '='
		{
			$$ = makeList(&nodes.String{Str: "="})
		}
	| '<'
		{
			$$ = makeList(&nodes.String{Str: "<"})
		}
	| '>'
		{
			$$ = makeList(&nodes.String{Str: ">"})
		}
	| LESS_EQUALS
		{
			$$ = makeList(&nodes.String{Str: "<="})
		}
	| GREATER_EQUALS
		{
			$$ = makeList(&nodes.String{Str: ">="})
		}
	| NOT_EQUALS
		{
			$$ = makeList(&nodes.String{Str: "<>"})
		}
	;

case_expr:
	CASE case_arg when_clause_list case_default END_P
		{
			$$ = &nodes.CaseExpr{
				Arg:       $2,
				Args:      $3,
				Defresult: $4,
				Location:  -1,
			}
		}
	;

when_clause_list:
	when_clause
		{
			$$ = makeList($1)
		}
	| when_clause_list when_clause
		{
			$$ = appendList($1, $2)
		}
	;

when_clause:
	WHEN a_expr THEN a_expr
		{
			$$ = &nodes.CaseWhen{
				Expr:     $2,
				Result:   $4,
				Location: -1,
			}
		}
	;

case_default:
	ELSE a_expr { $$ = $2 }
	| /* EMPTY */ { $$ = nil }
	;

case_arg:
	a_expr { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

/*****************************************************************************
 *
 *      ARRAY expressions
 *
 *****************************************************************************/

array_expr:
	'[' expr_list ']'
		{
			$$ = &nodes.A_ArrayExpr{
				Elements: $2,
				Location: -1,
			}
		}
	| '[' array_expr_list ']'
		{
			$$ = &nodes.A_ArrayExpr{
				Elements: $2,
				Location: -1,
			}
		}
	| '[' ']'
		{
			$$ = &nodes.A_ArrayExpr{
				Location: -1,
			}
		}
	;

array_expr_list:
	array_expr
		{
			$$ = makeList($1)
		}
	| array_expr_list ',' array_expr
		{
			$$ = appendList($1, $3)
		}
	;

/*****************************************************************************
 *
 *      ROW expressions
 *
 *****************************************************************************/

row:
	explicit_row
		{
			$$ = $1
		}
	| implicit_row
		{
			$$ = &nodes.RowExpr{
				Args:     $1,
				RowFormat: nodes.COERCE_IMPLICIT_CAST,
				Location: -1,
			}
		}
	;

explicit_row:
	ROW '(' expr_list ')'
		{
			$$ = &nodes.RowExpr{
				Args:     $3,
				RowFormat: nodes.COERCE_EXPLICIT_CALL,
				Location: -1,
			}
		}
	| ROW '(' ')'
		{
			$$ = &nodes.RowExpr{
				RowFormat: nodes.COERCE_EXPLICIT_CALL,
				Location: -1,
			}
		}
	;

implicit_row:
	'(' expr_list ',' a_expr ')'
		{
			$$ = appendList($2, $4)
		}
	;

b_expr:
	c_expr { $$ = $1 }
	| b_expr '+' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "+", $1, $3)
		}
	| b_expr '-' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "-", $1, $3)
		}
	| b_expr '*' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "*", $1, $3)
		}
	| b_expr '/' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "/", $1, $3)
		}
	| b_expr '%' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "%", $1, $3)
		}
	| b_expr '^' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "^", $1, $3)
		}
	| b_expr '<' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "<", $1, $3)
		}
	| b_expr '>' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, ">", $1, $3)
		}
	| b_expr '=' b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "=", $1, $3)
		}
	| b_expr LESS_EQUALS b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "<=", $1, $3)
		}
	| b_expr GREATER_EQUALS b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, ">=", $1, $3)
		}
	| b_expr NOT_EQUALS b_expr
		{
			$$ = makeAExpr(nodes.AEXPR_OP, "<>", $1, $3)
		}
	| b_expr qual_Op b_expr					%prec Op
		{
			$$ = makeAExprFromList(nodes.AEXPR_OP, $2, $1, $3)
		}
	| qual_Op b_expr						%prec Op
		{
			$$ = makeAExprFromList(nodes.AEXPR_OP, $1, nil, $2)
		}
	| b_expr IS DISTINCT FROM b_expr		%prec IS
		{
			$$ = makeAExpr(nodes.AEXPR_DISTINCT, "=", $1, $5)
		}
	| b_expr IS NOT DISTINCT FROM b_expr	%prec IS
		{
			$$ = makeAExpr(nodes.AEXPR_NOT_DISTINCT, "=", $1, $6)
		}
	| b_expr IS DOCUMENT_P					%prec IS
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_DOCUMENT,
				Args:     makeList($1),
				Location: -1,
			}
		}
	| b_expr IS NOT DOCUMENT_P				%prec IS
		{
			$$ = makeNotExpr(&nodes.XmlExpr{
				Op:       nodes.IS_DOCUMENT,
				Args:     makeList($1),
				Location: -1,
			})
		}
	| b_expr TYPECAST Typename
		{
			$$ = &nodes.TypeCast{
				Arg:      $1,
				TypeName: $3,
				Location: -1,
			}
		}
	| '+' b_expr %prec UMINUS
		{
			$$ = $2
		}
	| '-' b_expr %prec UMINUS
		{
			$$ = doNegate($2)
		}
	;

c_expr:
	columnref { $$ = $1 }
	| AexprConst { $$ = $1 }
	| PARAM opt_indirection
		{
			p := &nodes.ParamRef{
				Number:   int($1),
				Location: -1,
			}
			if $2 != nil {
				$$ = &nodes.A_Indirection{
					Arg:         p,
					Indirection: $2,
				}
			} else {
				$$ = p
			}
		}
	| '(' a_expr ')' opt_indirection
		{
			if $4 != nil {
				$$ = &nodes.A_Indirection{
					Arg:         $2,
					Indirection: $4,
				}
			} else {
				$$ = $2
			}
		}
	| func_expr { $$ = $1 }
	| select_with_parens %prec UMINUS
		{
			$$ = &nodes.SubLink{
				SubLinkType: int(nodes.EXPR_SUBLINK),
				Subselect:   $1,
				Location:    -1,
			}
		}
	| EXISTS select_with_parens
		{
			$$ = &nodes.SubLink{
				SubLinkType: int(nodes.EXISTS_SUBLINK),
				Subselect:   $2,
				Location:    -1,
			}
		}
	| case_expr { $$ = $1 }
	| ARRAY select_with_parens
		{
			$$ = &nodes.SubLink{
				SubLinkType: int(nodes.ARRAY_SUBLINK),
				Subselect:   $2,
				Location:    -1,
			}
		}
	| ARRAY array_expr
		{
			$$ = $2
		}
	| CAST '(' a_expr AS Typename ')'
		{
			$$ = &nodes.TypeCast{
				Arg:      $3,
				TypeName: $5,
				Location: -1,
			}
		}
	| explicit_row
		{
			$$ = $1
		}
	| implicit_row
		{
			$$ = &nodes.RowExpr{
				Args:     $1,
				RowFormat: nodes.COERCE_IMPLICIT_CAST,
				Location: -1,
			}
		}
	;

// Function expressions
func_expr:
	func_application within_group_clause filter_clause over_clause
		{
			n := $1.(*nodes.FuncCall)
			if $2 != nil {
				n.AggOrder = $2.(*nodes.List)
				n.AggWithinGroup = true
			}
			if $3 != nil {
				n.AggFilter = $3
			}
			if $4 != nil {
				n.Over = $4
			}
			$$ = n
		}
	| func_expr_common_subexpr { $$ = $1 }
	| json_aggregate_func filter_clause over_clause
		{
			/* json_aggregate_func returns a node; wrap filter/over if present */
			n := $1
			_ = $2
			_ = $3
			$$ = n
		}
	;

func_expr_windowless:
	func_application		{ $$ = $1 }
	| func_expr_common_subexpr	{ $$ = $1 }
	;

within_group_clause:
	WITHIN GROUP_P '(' sort_clause ')'	{ $$ = &nodes.List{Items: $4.Items} }
	| /* EMPTY */						{ $$ = nil }
	;

filter_clause:
	FILTER '(' WHERE a_expr ')'		{ $$ = $4 }
	| /* EMPTY */						{ $$ = nil }
	;

over_clause:
	OVER window_specification	{ $$ = $2 }
	| OVER ColId
		{
			$$ = &nodes.WindowDef{
				Name:         $2,
				FrameOptions: nodes.FRAMEOPTION_DEFAULTS,
				Location:     -1,
			}
		}
	| /* EMPTY */				{ $$ = nil }
	;

window_clause:
	WINDOW window_definition_list	{ $$ = $2 }
	| /* EMPTY */					{ $$ = nil }
	;

window_definition_list:
	window_definition							{ $$ = makeList($1) }
	| window_definition_list ',' window_definition	{ $$ = appendList($1, $3) }
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
			if $3 != nil {
				n.PartitionClause = $3
			}
			if $4 != nil {
				n.OrderClause = $4
			}
			n.Location = -1
			$$ = n
		}
	;

opt_existing_window_name:
	ColId			{ $$ = $1 }
	| /* EMPTY */	%prec Op	{ $$ = "" }
	;

opt_partition_clause:
	PARTITION BY expr_list	{ $$ = $3 }
	| /* EMPTY */			{ $$ = nil }
	;

opt_frame_clause:
	RANGE frame_extent opt_window_exclusion_clause
		{
			n := $2.(*nodes.WindowDef)
			n.FrameOptions |= nodes.FRAMEOPTION_NONDEFAULT | nodes.FRAMEOPTION_RANGE | int($3)
			$$ = n
		}
	| ROWS frame_extent opt_window_exclusion_clause
		{
			n := $2.(*nodes.WindowDef)
			n.FrameOptions |= nodes.FRAMEOPTION_NONDEFAULT | nodes.FRAMEOPTION_ROWS | int($3)
			$$ = n
		}
	| GROUPS frame_extent opt_window_exclusion_clause
		{
			n := $2.(*nodes.WindowDef)
			n.FrameOptions |= nodes.FRAMEOPTION_NONDEFAULT | nodes.FRAMEOPTION_GROUPS | int($3)
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
	EXCLUDE CURRENT_P ROW		{ $$ = int64(nodes.FRAMEOPTION_EXCLUDE_CURRENT_ROW) }
	| EXCLUDE GROUP_P			{ $$ = int64(nodes.FRAMEOPTION_EXCLUDE_GROUP) }
	| EXCLUDE TIES				{ $$ = int64(nodes.FRAMEOPTION_EXCLUDE_TIES) }
	| EXCLUDE NO OTHERS			{ $$ = 0 }
	| /* EMPTY */				{ $$ = 0 }
	;

func_application:
	func_name '(' ')'
		{
			$$ = &nodes.FuncCall{
				Funcname: $1,
			}
		}
	| func_name '(' func_arg_list opt_sort_clause ')'
		{
			n := &nodes.FuncCall{
				Funcname: $1,
				Args:     $3,
			}
			if $4 != nil {
				n.AggOrder = $4
			}
			$$ = n
		}
	| func_name '(' VARIADIC func_arg_expr opt_sort_clause ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:     $1,
				Args:         makeList($4),
				FuncVariadic: true,
				AggOrder:     $5,
			}
		}
	| func_name '(' func_arg_list ',' VARIADIC func_arg_expr opt_sort_clause ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:     $1,
				Args:         appendList($3, $6),
				FuncVariadic: true,
				AggOrder:     $7,
			}
		}
	| func_name '(' '*' ')'
		{
			$$ = &nodes.FuncCall{
				Funcname: $1,
				AggStar:  true,
			}
		}
	| func_name '(' DISTINCT func_arg_list opt_sort_clause ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:    $1,
				Args:        $4,
				AggDistinct: true,
				AggOrder:    $5,
			}
		}
	| func_name '(' ALL func_arg_list opt_sort_clause ')'
		{
			n := &nodes.FuncCall{
				Funcname: $1,
				Args:     $4,
			}
			if $5 != nil {
				n.AggOrder = $5
			}
			$$ = n
		}
	;

/*
 * func_expr_common_subexpr: SQL-standard special function expressions that
 * use reserved keywords as function names. These cannot be ordinary function
 * calls because the keywords would be consumed as part of the syntax.
 */
func_expr_common_subexpr:
	COLLATION FOR '(' a_expr ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "pg_collation_for"),
				Args:       makeList($4),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| CURRENT_DATE
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_DATE, -1)
		}
	| CURRENT_TIME
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIME, -1)
		}
	| CURRENT_TIME '(' Iconst ')'
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIME_N, int($3))
		}
	| CURRENT_TIMESTAMP
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIMESTAMP, -1)
		}
	| CURRENT_TIMESTAMP '(' Iconst ')'
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_TIMESTAMP_N, int($3))
		}
	| LOCALTIME
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIME, -1)
		}
	| LOCALTIME '(' Iconst ')'
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIME_N, int($3))
		}
	| LOCALTIMESTAMP
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIMESTAMP, -1)
		}
	| LOCALTIMESTAMP '(' Iconst ')'
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_LOCALTIMESTAMP_N, int($3))
		}
	| CURRENT_ROLE
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_ROLE, -1)
		}
	| CURRENT_USER
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_USER, -1)
		}
	| SESSION_USER
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_SESSION_USER, -1)
		}
	| SYSTEM_USER
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "system_user"),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| USER
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_USER, -1)
		}
	| CURRENT_CATALOG
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_CATALOG, -1)
		}
	| CURRENT_SCHEMA
		{
			$$ = makeSQLValueFunction(nodes.SVFOP_CURRENT_SCHEMA, -1)
		}
	| CAST '(' a_expr AS Typename ')'
		{
			$$ = makeTypeCast($3, $5)
		}
	| NULLIF '(' a_expr ',' a_expr ')'
		{
			$$ = makeAExpr(nodes.AEXPR_NULLIF, "=", $3, $5)
		}
	| COALESCE '(' expr_list ')'
		{
			$$ = &nodes.CoalesceExpr{
				Args:     $3,
				Location: -1,
			}
		}
	| GREATEST '(' expr_list ')'
		{
			$$ = &nodes.MinMaxExpr{
				Op:       nodes.IS_GREATEST,
				Args:     $3,
				Location: -1,
			}
		}
	| LEAST '(' expr_list ')'
		{
			$$ = &nodes.MinMaxExpr{
				Op:       nodes.IS_LEAST,
				Args:     $3,
				Location: -1,
			}
		}
	| EXTRACT '(' extract_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "extract"),
				Args:       $3,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| NORMALIZE '(' a_expr ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "normalize"),
				Args:       makeList($3),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| NORMALIZE '(' a_expr ',' unicode_normal_form ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "normalize"),
				Args:       makeList2($3, makeStringConst($5)),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| OVERLAY '(' overlay_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "overlay"),
				Args:       $3,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| POSITION '(' position_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "position"),
				Args:       $3,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| SUBSTRING '(' substr_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "substring"),
				Args:       $3,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| TRIM '(' BOTH trim_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "btrim"),
				Args:       $4,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| TRIM '(' LEADING trim_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "ltrim"),
				Args:       $4,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| TRIM '(' TRAILING trim_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "rtrim"),
				Args:       $4,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| TRIM '(' trim_list ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "btrim"),
				Args:       $3,
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| GROUPING '(' expr_list ')'
		{
			$$ = &nodes.GroupingFunc{
				Args:     $3,
				Location: -1,
			}
		}
	| XMLCONCAT '(' expr_list ')'
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_XMLCONCAT,
				Args:     $3,
				Location: -1,
			}
		}
	| XMLELEMENT '(' NAME_P ColLabel ')'
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_XMLELEMENT,
				Name:     $4,
				Location: -1,
			}
		}
	| XMLELEMENT '(' NAME_P ColLabel ',' xml_attributes ')'
		{
			$$ = &nodes.XmlExpr{
				Op:        nodes.IS_XMLELEMENT,
				Name:      $4,
				NamedArgs: $6,
				Location:  -1,
			}
		}
	| XMLELEMENT '(' NAME_P ColLabel ',' expr_list ')'
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_XMLELEMENT,
				Name:     $4,
				Args:     $6,
				Location: -1,
			}
		}
	| XMLELEMENT '(' NAME_P ColLabel ',' xml_attributes ',' expr_list ')'
		{
			$$ = &nodes.XmlExpr{
				Op:        nodes.IS_XMLELEMENT,
				Name:      $4,
				NamedArgs: $6,
				Args:      $8,
				Location:  -1,
			}
		}
	| XMLEXISTS '(' c_expr xmlexists_argument ')'
		{
			/* xmlexists(A PASSING [BY REF] B [BY REF]) is converted to xmlexists(A, B) */
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "xmlexists"),
				Args:       makeList2($3, $4),
				FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
				Location:   -1,
			}
		}
	| XMLFOREST '(' xml_attribute_list ')'
		{
			$$ = &nodes.XmlExpr{
				Op:        nodes.IS_XMLFOREST,
				NamedArgs: $3,
				Location:  -1,
			}
		}
	| XMLPARSE '(' document_or_content a_expr xml_whitespace_option ')'
		{
			x := &nodes.XmlExpr{
				Op:       nodes.IS_XMLPARSE,
				Args:     makeList2($4, makeBoolAConst($5)),
				Xmloption: nodes.XmlOptionType($3),
				Location: -1,
			}
			$$ = x
		}
	| XMLPI '(' NAME_P ColLabel ')'
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_XMLPI,
				Name:     $4,
				Location: -1,
			}
		}
	| XMLPI '(' NAME_P ColLabel ',' a_expr ')'
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_XMLPI,
				Name:     $4,
				Args:     makeList($6),
				Location: -1,
			}
		}
	| XMLROOT '(' a_expr ',' xml_root_version opt_xml_root_standalone ')'
		{
			$$ = &nodes.XmlExpr{
				Op:       nodes.IS_XMLROOT,
				Args:     &nodes.List{Items: []nodes.Node{$3, $5, $6}},
				Location: -1,
			}
		}
	| XMLSERIALIZE '(' document_or_content a_expr AS SimpleTypename xml_indent_option ')'
		{
			$$ = &nodes.XmlSerialize{
				Xmloption: nodes.XmlOptionType($3),
				Expr:      $4,
				TypeName:  $6,
				Indent:    $7 != 0,
				Location:  -1,
			}
		}
	/* SQL/JSON function expressions */
	| JSON_OBJECT '(' func_arg_list ')'
		{
			/* Legacy json_object() function call */
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("pg_catalog", "json_object"),
				Args:       $3,
				FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
				Location:   -1,
			}
		}
	| JSON_OBJECT '(' json_name_and_value_list json_object_constructor_null_clause_opt
		json_key_uniqueness_constraint_opt json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonObjectConstructor{
				Exprs:        $3,
				Output:       asJsonOutput($6),
				AbsentOnNull: $4 != 0,
				UniqueKeys:   $5 != 0,
				Location:     -1,
			}
		}
	| JSON_OBJECT '(' json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonObjectConstructor{
				Output:   asJsonOutput($3),
				Location: -1,
			}
		}
	| JSON_ARRAY '(' json_value_expr_list json_array_constructor_null_clause_opt
		json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonArrayConstructor{
				Exprs:        $3,
				AbsentOnNull: $4 != 0,
				Output:       asJsonOutput($5),
				Location:     -1,
			}
		}
	| JSON_ARRAY '(' select_no_parens json_format_clause_opt json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonArrayQueryConstructor{
				Query:  $3,
				Output: asJsonOutput($5),
				Location: -1,
			}
		}
	| JSON_ARRAY '(' json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonArrayConstructor{
				Output:   asJsonOutput($3),
				Location: -1,
			}
		}
	| JSON '(' json_value_expr json_key_uniqueness_constraint_opt ')'
		{
			$$ = &nodes.JsonParseExpr{
				Expr:       $3.(*nodes.JsonValueExpr),
				UniqueKeys: $4 != 0,
				Location:   -1,
			}
		}
	| JSON_SCALAR '(' a_expr ')'
		{
			$$ = &nodes.JsonScalarExpr{
				Expr:     $3,
				Location: -1,
			}
		}
	| JSON_SERIALIZE '(' json_value_expr json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonSerializeExpr{
				Expr:   $3.(*nodes.JsonValueExpr),
				Output: asJsonOutput($4),
				Location: -1,
			}
		}
	| JSON_QUERY '(' json_value_expr ',' a_expr json_passing_clause_opt
		json_returning_clause_opt json_wrapper_behavior json_quotes_clause_opt
		json_behavior_clause_opt ')'
		{
			onEmpty, onError := splitJsonBehaviorClause($10)
			$$ = &nodes.JsonFuncExpr{
				Op:          nodes.JSON_QUERY_OP,
				ContextItem: $3.(*nodes.JsonValueExpr),
				Pathspec:    $5,
				Passing:     $6,
				Output:      asJsonOutput($7),
				Wrapper:     nodes.JsonWrapper($8),
				Quotes:      nodes.JsonQuotes($9),
				OnEmpty:     onEmpty,
				OnError:     onError,
				Location:    -1,
			}
		}
	| JSON_EXISTS '(' json_value_expr ',' a_expr json_passing_clause_opt
		json_on_error_clause_opt ')'
		{
			$$ = &nodes.JsonFuncExpr{
				Op:          nodes.JSON_EXISTS_OP,
				ContextItem: $3.(*nodes.JsonValueExpr),
				Pathspec:    $5,
				Passing:     $6,
				OnError:     asJsonBehavior($7),
				Location:    -1,
			}
		}
	| JSON_VALUE '(' json_value_expr ',' a_expr json_passing_clause_opt
		json_returning_clause_opt json_behavior_clause_opt ')'
		{
			onEmpty, onError := splitJsonBehaviorClause($8)
			$$ = &nodes.JsonFuncExpr{
				Op:          nodes.JSON_VALUE_OP,
				ContextItem: $3.(*nodes.JsonValueExpr),
				Pathspec:    $5,
				Passing:     $6,
				Output:      asJsonOutput($7),
				OnEmpty:     onEmpty,
				OnError:     onError,
				Location:    -1,
			}
		}
	| MERGE_ACTION '(' ')'
		{
			$$ = &nodes.FuncCall{
				Funcname:   makeFuncName("merge_action"),
				Location:   -1,
			}
		}
	;

/*
 * Helper rules for func_expr_common_subexpr
 */
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
	NFC         { $$ = "NFC" }
	| NFD       { $$ = "NFD" }
	| NFKC      { $$ = "NFKC" }
	| NFKD      { $$ = "NFKD" }
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
			/* note: arguments reversed per PG convention */
			$$ = makeList2($3, $1)
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
				$3,
			}}
		}
	| a_expr SIMILAR a_expr ESCAPE a_expr
		{
			$$ = &nodes.List{Items: []nodes.Node{$1, $3, $5}}
		}
	| expr_list
		{
			/* comma-separated form: substring(x, 1, 3) */
			$$ = $1
		}
	;

trim_list:
	a_expr FROM expr_list
		{
			$$ = prependList($1, $3)
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

/*****************************************************************************
 *
 * XML helper rules
 *
 *****************************************************************************/

xml_root_version:
	VERSION_P a_expr
		{ $$ = $2 }
	| VERSION_P NO VALUE_P
		{ $$ = makeNullAConst() }
	;

opt_xml_root_standalone:
	',' STANDALONE_P YES_P
		{ $$ = makeIntConst(int64(nodes.XML_STANDALONE_YES)) }
	| ',' STANDALONE_P NO
		{ $$ = makeIntConst(int64(nodes.XML_STANDALONE_NO)) }
	| ',' STANDALONE_P NO VALUE_P
		{ $$ = makeIntConst(int64(nodes.XML_STANDALONE_NO_VALUE)) }
	| /* EMPTY */
		{ $$ = makeIntConst(int64(nodes.XML_STANDALONE_OMITTED)) }
	;

xml_attributes:
	XMLATTRIBUTES '(' xml_attribute_list ')'	{ $$ = $3 }
	;

xml_attribute_list:
	xml_attribute_el
		{ $$ = makeList($1) }
	| xml_attribute_list ',' xml_attribute_el
		{ $$ = appendList($1, $3) }
	;

xml_attribute_el:
	a_expr AS ColLabel
		{
			$$ = &nodes.ResTarget{
				Name:     $3,
				Val:      $1,
				Location: -1,
			}
		}
	| a_expr
		{
			$$ = &nodes.ResTarget{
				Val:      $1,
				Location: -1,
			}
		}
	;

document_or_content:
	DOCUMENT_P		{ $$ = int64(nodes.XMLOPTION_DOCUMENT) }
	| CONTENT_P	{ $$ = int64(nodes.XMLOPTION_CONTENT) }
	;

xml_indent_option:
	INDENT			{ $$ = 1 }
	| NO INDENT		{ $$ = 0 }
	| /* EMPTY */	{ $$ = 0 }
	;

xml_whitespace_option:
	PRESERVE WHITESPACE_P		{ $$ = 1 }
	| STRIP_P WHITESPACE_P		{ $$ = 0 }
	| /* EMPTY */				{ $$ = 0 }
	;

xmlexists_argument:
	PASSING c_expr
		{ $$ = $2 }
	| PASSING c_expr xml_passing_mech
		{ $$ = $2 }
	| PASSING xml_passing_mech c_expr
		{ $$ = $3 }
	| PASSING xml_passing_mech c_expr xml_passing_mech
		{ $$ = $3 }
	;

xml_passing_mech:
	BY REF_P
	| BY VALUE_P
	;

/*****************************************************************************
 *
 * XMLTABLE
 *
 *****************************************************************************/

xmltable:
	XMLTABLE '(' c_expr xmlexists_argument COLUMNS xmltable_column_list ')'
		{
			$$ = &nodes.RangeTableFunc{
				Rowexpr:  $3,
				Docexpr:  $4,
				Columns:  $6,
				Location: -1,
			}
		}
	| XMLTABLE '(' XMLNAMESPACES '(' xml_namespace_list ')' ','
		c_expr xmlexists_argument COLUMNS xmltable_column_list ')'
		{
			$$ = &nodes.RangeTableFunc{
				Rowexpr:    $8,
				Docexpr:    $9,
				Columns:    $11,
				Namespaces: $5,
				Location:   -1,
			}
		}
	;

xmltable_column_list:
	xmltable_column_el
		{ $$ = makeList($1) }
	| xmltable_column_list ',' xmltable_column_el
		{ $$ = appendList($1, $3) }
	;

xmltable_column_el:
	ColId Typename
		{
			$$ = &nodes.RangeTableFuncCol{
				Colname:  $1,
				TypeName: $2,
				Location: -1,
			}
		}
	| ColId Typename xmltable_column_option_list
		{
			fc := &nodes.RangeTableFuncCol{
				Colname:  $1,
				TypeName: $2,
				Location: -1,
			}
			for _, item := range $3.Items {
				defel := item.(*nodes.DefElem)
				switch defel.Defname {
				case "default":
					fc.Coldefexpr = defel.Arg
				case "path":
					fc.Colexpr = defel.Arg
				case "__pg__is_not_null":
					if defel.Arg != nil {
						if b, ok := defel.Arg.(*nodes.Boolean); ok && b.Boolval {
							fc.IsNotNull = true
						}
					}
				}
			}
			$$ = fc
		}
	| ColId FOR ORDINALITY
		{
			$$ = &nodes.RangeTableFuncCol{
				Colname:       $1,
				ForOrdinality: true,
				Location:      -1,
			}
		}
	;

xmltable_column_option_list:
	xmltable_column_option_el
		{ $$ = makeList($1) }
	| xmltable_column_option_list xmltable_column_option_el
		{ $$ = appendList($1, $2) }
	;

xmltable_column_option_el:
	IDENT b_expr
		{
			$$ = makeDefElem($1, $2)
		}
	| DEFAULT b_expr
		{
			$$ = makeDefElem("default", $2)
		}
	| NOT NULL_P
		{
			$$ = makeDefElem("__pg__is_not_null", &nodes.Boolean{Boolval: true})
		}
	| NULL_P
		{
			$$ = makeDefElem("__pg__is_not_null", &nodes.Boolean{Boolval: false})
		}
	| PATH b_expr
		{
			$$ = makeDefElem("path", $2)
		}
	;

xml_namespace_list:
	xml_namespace_el
		{ $$ = makeList($1) }
	| xml_namespace_list ',' xml_namespace_el
		{ $$ = appendList($1, $3) }
	;

xml_namespace_el:
	b_expr AS ColLabel
		{
			$$ = &nodes.ResTarget{
				Name:     $3,
				Val:      $1,
				Location: -1,
			}
		}
	| DEFAULT b_expr
		{
			$$ = &nodes.ResTarget{
				Val:      $2,
				Location: -1,
			}
		}
	;

/*****************************************************************************
 *
 * SQL/JSON helper rules
 *
 *****************************************************************************/

json_value_expr:
	a_expr json_format_clause_opt
		{
			$$ = &nodes.JsonValueExpr{
				RawExpr: $1,
			}
		}
	;

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
	json_format_clause      { $$ = $1 }
	| /* EMPTY */           { $$ = nil }
	;

json_returning_clause_opt:
	RETURNING Typename json_format_clause_opt
		{
			$$ = &nodes.JsonOutput{
				TypeName: $2,
			}
		}
	| /* EMPTY */
		{ $$ = nil }
	;

json_output_clause_opt:
	json_returning_clause_opt { $$ = $1 }
	;

json_behavior:
	DEFAULT a_expr
		{
			$$ = &nodes.JsonBehavior{
				Btype: nodes.JSON_BEHAVIOR_DEFAULT,
				Expr:  $2,
				Location: -1,
			}
		}
	| json_behavior_type
		{
			$$ = &nodes.JsonBehavior{
				Btype: nodes.JsonBehaviorType($1),
				Location: -1,
			}
		}
	;

json_behavior_type:
	ERROR_P         { $$ = int64(nodes.JSON_BEHAVIOR_ERROR) }
	| NULL_P        { $$ = int64(nodes.JSON_BEHAVIOR_NULL) }
	| TRUE_P        { $$ = int64(nodes.JSON_BEHAVIOR_TRUE) }
	| FALSE_P       { $$ = int64(nodes.JSON_BEHAVIOR_FALSE) }
	| UNKNOWN       { $$ = int64(nodes.JSON_BEHAVIOR_UNKNOWN) }
	| EMPTY_P ARRAY { $$ = int64(nodes.JSON_BEHAVIOR_EMPTY_ARRAY) }
	| EMPTY_P OBJECT_P { $$ = int64(nodes.JSON_BEHAVIOR_EMPTY_OBJECT) }
	/* non-standard, for Oracle compatibility */
	| EMPTY_P          { $$ = int64(nodes.JSON_BEHAVIOR_EMPTY_ARRAY) }
	;

json_behavior_clause_opt:
	json_behavior ON EMPTY_P json_behavior ON ERROR_P
		{
			/* Returns a list of 2: [on_empty, on_error] */
			$$ = &nodes.List{Items: []nodes.Node{$1, $4}}
		}
	| json_behavior ON EMPTY_P
		{
			$$ = &nodes.List{Items: []nodes.Node{$1, nil}}
		}
	| json_behavior ON ERROR_P
		{
			$$ = &nodes.List{Items: []nodes.Node{nil, $1}}
		}
	| /* EMPTY */
		{ $$ = nil }
	;

json_on_error_clause_opt:
	json_behavior ON ERROR_P
		{ $$ = $1 }
	| /* EMPTY */
		{ $$ = nil }
	;

json_wrapper_behavior:
	WITHOUT WRAPPER                 { $$ = int64(nodes.JSW_NONE) }
	| WITHOUT ARRAY WRAPPER         { $$ = int64(nodes.JSW_NONE) }
	| WITH WRAPPER                  { $$ = int64(nodes.JSW_UNCONDITIONAL) }
	| WITH ARRAY WRAPPER            { $$ = int64(nodes.JSW_UNCONDITIONAL) }
	| WITH CONDITIONAL ARRAY WRAPPER { $$ = int64(nodes.JSW_CONDITIONAL) }
	| WITH CONDITIONAL WRAPPER      { $$ = int64(nodes.JSW_CONDITIONAL) }
	| WITH UNCONDITIONAL ARRAY WRAPPER { $$ = int64(nodes.JSW_UNCONDITIONAL) }
	| WITH UNCONDITIONAL WRAPPER    { $$ = int64(nodes.JSW_UNCONDITIONAL) }
	| /* EMPTY */                   { $$ = int64(nodes.JSW_UNSPEC) }
	;

json_quotes_clause_opt:
	KEEP QUOTES ON SCALAR STRING_P  { $$ = int64(nodes.JS_QUOTES_KEEP) }
	| KEEP QUOTES                   { $$ = int64(nodes.JS_QUOTES_KEEP) }
	| OMIT QUOTES ON SCALAR STRING_P { $$ = int64(nodes.JS_QUOTES_OMIT) }
	| OMIT QUOTES                   { $$ = int64(nodes.JS_QUOTES_OMIT) }
	| /* EMPTY */                   { $$ = int64(nodes.JS_QUOTES_UNSPEC) }
	;

json_predicate_type_constraint:
	JSON       %prec UNBOUNDED      { $$ = int64(nodes.JS_TYPE_ANY) }
	| JSON VALUE_P                  { $$ = int64(nodes.JS_TYPE_ANY) }
	| JSON ARRAY                    { $$ = int64(nodes.JS_TYPE_ARRAY) }
	| JSON OBJECT_P                 { $$ = int64(nodes.JS_TYPE_OBJECT) }
	| JSON SCALAR                   { $$ = int64(nodes.JS_TYPE_SCALAR) }
	;

json_key_uniqueness_constraint_opt:
	WITH UNIQUE KEYS                    { $$ = 1 }
	| WITH UNIQUE      %prec UNBOUNDED  { $$ = 1 }
	| WITHOUT UNIQUE KEYS               { $$ = 0 }
	| WITHOUT UNIQUE   %prec UNBOUNDED  { $$ = 0 }
	| /* EMPTY */      %prec UNBOUNDED  { $$ = 0 }
	;

json_name_and_value_list:
	json_name_and_value
		{ $$ = makeList($1) }
	| json_name_and_value_list ',' json_name_and_value
		{ $$ = appendList($1, $3) }
	;

json_name_and_value:
	c_expr VALUE_P json_value_expr
		{
			$$ = &nodes.JsonKeyValue{
				Key:   $1,
				Value: $3.(*nodes.JsonValueExpr),
			}
		}
	| a_expr ':' json_value_expr
		{
			$$ = &nodes.JsonKeyValue{
				Key:   $1,
				Value: $3.(*nodes.JsonValueExpr),
			}
		}
	;

json_object_constructor_null_clause_opt:
	NULL_P ON NULL_P                { $$ = 0 }
	| ABSENT ON NULL_P              { $$ = 1 }
	| /* EMPTY */                   { $$ = 0 }
	;

json_array_constructor_null_clause_opt:
	NULL_P ON NULL_P                { $$ = 0 }
	| ABSENT ON NULL_P              { $$ = 1 }
	| /* EMPTY */                   { $$ = 1 }
	;

json_value_expr_list:
	json_value_expr
		{ $$ = makeList($1) }
	| json_value_expr_list ',' json_value_expr
		{ $$ = appendList($1, $3) }
	;

json_passing_clause_opt:
	PASSING json_arguments  { $$ = $2 }
	| /* EMPTY */           { $$ = nil }
	;

json_arguments:
	json_argument
		{ $$ = makeList($1) }
	| json_arguments ',' json_argument
		{ $$ = appendList($1, $3) }
	;

json_argument:
	json_value_expr AS ColLabel
		{
			$$ = &nodes.JsonArgument{
				Val:  $1.(*nodes.JsonValueExpr),
				Name: $3,
			}
		}
	;

json_aggregate_func:
	JSON_OBJECTAGG '(' json_name_and_value json_object_constructor_null_clause_opt
		json_key_uniqueness_constraint_opt json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonObjectAgg{
				Constructor: &nodes.JsonAggConstructor{
					Output:   asJsonOutput($6),
					Location: -1,
				},
				Arg:          $3.(*nodes.JsonKeyValue),
				AbsentOnNull: $4 != 0,
				UniqueKeys:   $5 != 0,
			}
		}
	| JSON_ARRAYAGG '(' json_value_expr json_array_aggregate_order_by_clause_opt
		json_array_constructor_null_clause_opt json_returning_clause_opt ')'
		{
			$$ = &nodes.JsonArrayAgg{
				Constructor: &nodes.JsonAggConstructor{
					Output:    asJsonOutput($6),
					Agg_order: $4,
					Location:  -1,
				},
				Arg:          $3.(*nodes.JsonValueExpr),
				AbsentOnNull: $5 != 0,
			}
		}
	;

json_array_aggregate_order_by_clause_opt:
	ORDER BY sortby_list    { $$ = $3 }
	| /* EMPTY */           { $$ = nil }
	;

/*****************************************************************************
 *
 * JSON_TABLE
 *
 *****************************************************************************/

json_table:
	JSON_TABLE '('
		json_value_expr ',' a_expr json_table_path_name_opt
		json_passing_clause_opt
		COLUMNS '(' json_table_column_definition_list ')'
		json_on_error_clause_opt
	')'
		{
			$$ = &nodes.JsonTable{
				ContextItem: $3.(*nodes.JsonValueExpr),
				Pathspec: &nodes.JsonTablePathSpec{
					String:   $5,
					Name:     $6,
					Location: -1,
				},
				Passing:  $7,
				Columns:  $10,
				OnError:  asJsonBehavior($12),
				Location: -1,
			}
		}
	;

json_table_path_name_opt:
	AS name         { $$ = $2 }
	| /* EMPTY */   { $$ = "" }
	;

json_table_column_definition_list:
	json_table_column_definition
		{ $$ = makeList($1) }
	| json_table_column_definition_list ',' json_table_column_definition
		{ $$ = appendList($1, $3) }
	;

json_table_column_definition:
	ColId FOR ORDINALITY
		{
			$$ = &nodes.JsonTableColumn{
				Coltype:  nodes.JTC_FOR_ORDINALITY,
				Name:     $1,
				Location: -1,
			}
		}
	| ColId Typename json_table_column_path_clause_opt
		json_wrapper_behavior json_quotes_clause_opt
		json_behavior_clause_opt
		{
			onEmpty, onError := splitJsonBehaviorClause($6)
			$$ = &nodes.JsonTableColumn{
				Coltype:  nodes.JTC_REGULAR,
				Name:     $1,
				TypeName: $2,
				Pathspec: asJsonTablePathSpec($3),
				Wrapper:  nodes.JsonWrapper($4),
				Quotes:   nodes.JsonQuotes($5),
				OnEmpty:  onEmpty,
				OnError:  onError,
				Location: -1,
			}
		}
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
	| ColId Typename EXISTS json_table_column_path_clause_opt
		json_on_error_clause_opt
		{
			$$ = &nodes.JsonTableColumn{
				Coltype:  nodes.JTC_EXISTS,
				Name:     $1,
				TypeName: $2,
				Pathspec: asJsonTablePathSpec($4),
				OnError:  asJsonBehavior($5),
				Location: -1,
			}
		}
	| NESTED path_opt Sconst json_table_path_name_opt
		COLUMNS '(' json_table_column_definition_list ')'
		{
			$$ = &nodes.JsonTableColumn{
				Coltype: nodes.JTC_NESTED,
				Pathspec: &nodes.JsonTablePathSpec{
					String:   makeStringConst($3),
					Name:     $4,
					Location: -1,
				},
				Columns:  $7,
				Location: -1,
			}
		}
	;

json_table_column_path_clause_opt:
	PATH Sconst
		{
			$$ = &nodes.JsonTablePathSpec{
				String:   makeStringConst($2),
				Location: -1,
			}
		}
	| /* EMPTY */
		{ $$ = nil }
	;

json_table_column_option_list:
	json_table_column_option_el
		{ $$ = makeList($1) }
	| json_table_column_option_list json_table_column_option_el
		{ $$ = appendList($1, $2) }
	;

json_table_column_option_el:
	DEFAULT b_expr
		{ $$ = makeDefElem("default", $2) }
	| PATH b_expr
		{ $$ = makeDefElem("path", $2) }
	| NOT NULL_P
		{ $$ = makeDefElem("__pg__is_not_null", &nodes.Boolean{Boolval: true}) }
	| NULL_P
		{ $$ = makeDefElem("__pg__is_not_null", &nodes.Boolean{Boolval: false}) }
	;

path_opt:
	PATH            { $$ = "" }
	| /* EMPTY */   { $$ = "" }
	;

func_arg_list:
	func_arg_expr
		{
			$$ = makeList($1)
		}
	| func_arg_list ',' func_arg_expr
		{
			$$ = appendList($1, $3)
		}
	;

func_arg_expr:
	a_expr { $$ = $1 }
	| param_name COLON_EQUALS a_expr
		{
			$$ = &nodes.NamedArgExpr{
				Name:      $1,
				Arg:       $3,
				Argnumber: -1,
				Location:  -1,
			}
		}
	| param_name EQUALS_GREATER a_expr
		{
			$$ = &nodes.NamedArgExpr{
				Name:      $1,
				Arg:       $3,
				Argnumber: -1,
				Location:  -1,
			}
		}
	;

func_name:
	type_function_name
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| ColId indirection
		{
			/* schema.funcname or catalog.schema.funcname */
			l := &nodes.List{Items: []nodes.Node{&nodes.String{Str: $1}}}
			if $2 != nil {
				for _, item := range $2.Items {
					l.Items = append(l.Items, item)
				}
			}
			$$ = l
		}
	;

columnref:
	ColId
		{
			$$ = &nodes.ColumnRef{
				Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: $1}}},
			}
		}
	| ColId indirection
		{
			fields := &nodes.List{Items: []nodes.Node{&nodes.String{Str: $1}}}
			if $2 != nil {
				fields.Items = append(fields.Items, $2.Items...)
			}
			$$ = &nodes.ColumnRef{Fields: fields}
		}
	;

indirection:
	indirection_el
		{
			$$ = makeList($1)
		}
	| indirection indirection_el
		{
			$$ = appendList($1, $2)
		}
	;

indirection_el:
	'.' attr_name
		{
			$$ = &nodes.String{Str: $2}
		}
	| '.' '*'
		{
			$$ = &nodes.A_Star{}
		}
	| '[' a_expr ']'
		{
			$$ = &nodes.A_Indices{Uidx: $2}
		}
	| '[' opt_slice_bound ':' opt_slice_bound ']'
		{
			$$ = &nodes.A_Indices{IsSlice: true, Lidx: $2, Uidx: $4}
		}
	;

opt_slice_bound:
	a_expr      { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

opt_indirection:
	indirection { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

attrs:
	'.' attr_name
		{
			$$ = makeList(&nodes.String{Str: $2})
		}
	| attrs '.' attr_name
		{
			$$ = appendList($1, &nodes.String{Str: $3})
		}
	;

AexprConst:
	Iconst
		{
			$$ = &nodes.A_Const{Val: &nodes.Integer{Ival: $1}}
		}
	| FCONST
		{
			$$ = &nodes.A_Const{Val: &nodes.Float{Fval: $1}}
		}
	| Sconst
		{
			$$ = &nodes.A_Const{Val: &nodes.String{Str: $1}}
		}
	| BCONST
		{
			$$ = &nodes.A_Const{Val: &nodes.BitString{Bsval: $1}}
		}
	| XCONST
		{
			$$ = &nodes.A_Const{Val: &nodes.BitString{Bsval: $1}}
		}
	| TRUE_P
		{
			$$ = &nodes.A_Const{Val: &nodes.Boolean{Boolval: true}}
		}
	| FALSE_P
		{
			$$ = &nodes.A_Const{Val: &nodes.Boolean{Boolval: false}}
		}
	| NULL_P
		{
			$$ = &nodes.A_Const{Isnull: true}
		}
	| func_name Sconst
		{
			/* generic type 'literal' syntax */
			t := makeTypeNameFromNameList($1).(*nodes.TypeName)
			$$ = makeStringConstCast($2, t)
		}
	| func_name '(' func_arg_list opt_sort_clause ')' Sconst
		{
			/* generic syntax with a type modifier */
			t := makeTypeNameFromNameList($1).(*nodes.TypeName)
			if $3 != nil {
				t.Typmods = $3
			}
			$$ = makeStringConstCast($6, t)
		}
	| ConstTypename Sconst
		{
			$$ = makeStringConstCast($2, $1)
		}
	| ConstInterval Sconst opt_interval
		{
			t := $1
			if $3 != nil {
				t.Typmods = $3
			}
			$$ = makeStringConstCast($2, t)
		}
	| ConstInterval '(' Iconst ')' Sconst
		{
			t := $1
			t.Typmods = makeList2(makeIntConst(int64(nodes.INTERVAL_FULL_RANGE)), makeIntConst($3))
			$$ = makeStringConstCast($5, t)
		}
	;

Iconst:
	ICONST { $$ = $1 }
	;

Sconst:
	SCONST { $$ = $1 }
	;

// Names and identifiers
ColId:
	IDENT { $$ = $1 }
	| unreserved_keyword { $$ = $1 }
	| col_name_keyword { $$ = $1 }
	;

ColLabel:
	IDENT { $$ = $1 }
	| unreserved_keyword { $$ = $1 }
	| col_name_keyword { $$ = $1 }
	| type_func_name_keyword { $$ = $1 }
	| reserved_keyword { $$ = $1 }
	;

attr_name:
	ColLabel { $$ = $1 }
	;

name:
	ColId { $$ = $1 }
	;

qualified_name:
	ColId
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| ColId '.' attr_name
		{
			l := makeList(&nodes.String{Str: $1})
			$$ = appendList(l, &nodes.String{Str: $3})
		}
	| ColId '.' attr_name '.' attr_name
		{
			l := makeList(&nodes.String{Str: $1})
			l = appendList(l, &nodes.String{Str: $3})
			$$ = appendList(l, &nodes.String{Str: $5})
		}
	;

any_name:
	ColId
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| ColId '.' any_name
		{
			$$ = prependList(&nodes.String{Str: $1}, $3)
		}
	;

opt_name_list:
	name_list { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

name_list:
	name
		{
			$$ = makeList(&nodes.String{Str: $1})
		}
	| name_list ',' name
		{
			$$ = appendList($1, &nodes.String{Str: $3})
		}
	;

expr_list:
	a_expr
		{
			$$ = makeList($1)
		}
	| expr_list ',' a_expr
		{
			$$ = appendList($1, $3)
		}
	;

/*****************************************************************************
 *
 *      Type Name (Typename) rules
 *
 *****************************************************************************/

Typename:
	SimpleTypename opt_array_bounds
		{
			$$ = $1
			if $2 != nil {
				$$.ArrayBounds = $2
			}
		}
	| SETOF SimpleTypename opt_array_bounds
		{
			$$ = $2
			$$.Setof = true
			if $3 != nil {
				$$.ArrayBounds = $3
			}
		}
	| SimpleTypename ARRAY '[' Iconst ']'
		{
			$$ = $1
			$$.ArrayBounds = makeList(&nodes.Integer{Ival: $4})
		}
	| SETOF SimpleTypename ARRAY '[' Iconst ']'
		{
			$$ = $2
			$$.Setof = true
			$$.ArrayBounds = makeList(&nodes.Integer{Ival: $5})
		}
	| SimpleTypename ARRAY
		{
			$$ = $1
			$$.ArrayBounds = makeList(&nodes.Integer{Ival: -1})
		}
	| SETOF SimpleTypename ARRAY
		{
			$$ = $2
			$$.Setof = true
			$$.ArrayBounds = makeList(&nodes.Integer{Ival: -1})
		}
	;

opt_array_bounds:
	opt_array_bounds '[' ']'
		{
			$$ = appendList($1, &nodes.Integer{Ival: -1})
		}
	| opt_array_bounds '[' Iconst ']'
		{
			$$ = appendList($1, &nodes.Integer{Ival: $3})
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

SimpleTypename:
	GenericType       { $$ = $1 }
	| Numeric         { $$ = $1 }
	| Bit             { $$ = $1 }
	| Character       { $$ = $1 }
	| ConstDatetime   { $$ = $1 }
	| ConstInterval opt_interval
		{
			$$ = $1
			if $2 != nil {
				$$.Typmods = $2
			}
		}
	| ConstInterval '(' Iconst ')'
		{
			$$ = $1
			$$.Typmods = makeList2(makeIntConst(int64(nodes.INTERVAL_FULL_RANGE)), makeIntConst($3))
		}
	| BOOLEAN_P       { $$ = makeTypeName("bool") }
	| JSON            { $$ = makeTypeName("json") }
	;

GenericType:
	type_function_name opt_type_modifiers
		{
			$$ = &nodes.TypeName{
				Names:    makeList(&nodes.String{Str: $1}),
				Typmods:  $2,
				Location: -1,
			}
		}
	| type_function_name '.' attr_name opt_type_modifiers
		{
			l := makeList(&nodes.String{Str: $1})
			l = appendList(l, &nodes.String{Str: $3})
			$$ = &nodes.TypeName{
				Names:    l,
				Typmods:  $4,
				Location: -1,
			}
		}
	;

opt_type_modifiers:
	'(' expr_list ')'  { $$ = $2 }
	| /* EMPTY */       { $$ = nil }
	;

Numeric:
	INT_P        { $$ = makeTypeName("int4") }
	| INTEGER    { $$ = makeTypeName("int4") }
	| SMALLINT   { $$ = makeTypeName("int2") }
	| BIGINT     { $$ = makeTypeName("int8") }
	| REAL       { $$ = makeTypeName("float4") }
	| FLOAT_P opt_float
		{
			$$ = $2
		}
	| DOUBLE_P PRECISION  { $$ = makeTypeName("float8") }
	| DECIMAL_P opt_type_modifiers
		{
			$$ = makeTypeName("numeric")
			$$.Typmods = $2
		}
	| DEC opt_type_modifiers
		{
			$$ = makeTypeName("numeric")
			$$.Typmods = $2
		}
	| NUMERIC opt_type_modifiers
		{
			$$ = makeTypeName("numeric")
			$$.Typmods = $2
		}
	;

opt_float:
	'(' Iconst ')'
		{
			if $2 <= 24 {
				$$ = makeTypeName("float4")
			} else {
				$$ = makeTypeName("float8")
			}
		}
	| /* EMPTY */
		{
			$$ = makeTypeName("float8")
		}
	;

Character:
	CHARACTER opt_varying '(' Iconst ')'
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				$$ = makeTypeName("bpchar")
			}
			$$.Typmods = makeList(&nodes.Integer{Ival: $4})
		}
	| CHARACTER opt_varying
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				$$ = makeTypeName("bpchar")
			}
		}
	| CHAR_P opt_varying '(' Iconst ')'
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				$$ = makeTypeName("bpchar")
			}
			$$.Typmods = makeList(&nodes.Integer{Ival: $4})
		}
	| CHAR_P opt_varying
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				$$ = makeTypeName("bpchar")
			}
		}
	| VARCHAR '(' Iconst ')'
		{
			$$ = makeTypeName("varchar")
			$$.Typmods = makeList(&nodes.Integer{Ival: $3})
		}
	| VARCHAR
		{
			$$ = makeTypeName("varchar")
		}
	;

opt_varying:
	VARYING { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

Bit:
	BitWithLength		{ $$ = $1 }
	| BitWithoutLength	{ $$ = $1 }
	;

BitWithLength:
	BIT opt_varying '(' expr_list ')'
		{
			if $2 {
				$$ = makeTypeName("varbit")
			} else {
				$$ = makeTypeName("bit")
			}
			$$.Typmods = $4
		}
	;

BitWithoutLength:
	BIT opt_varying
		{
			if $2 {
				$$ = makeTypeName("varbit")
			} else {
				$$ = makeTypeName("bit")
				$$.Typmods = makeList(makeIntConst(1))
			}
		}
	;

ConstTypename:
	Numeric				{ $$ = $1 }
	| ConstBit			{ $$ = $1 }
	| ConstCharacter	{ $$ = $1 }
	| ConstDatetime		{ $$ = $1 }
	| JSON              { $$ = makeTypeName("json") }
	;

ConstBit:
	BitWithLength		{ $$ = $1 }
	| BitWithoutLength
		{
			$$ = $1
			/* ConstBit defaults to unspecified length for BIT */
		}
	;

ConstCharacter:
	CHARACTER opt_varying '(' Iconst ')'
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				$$ = makeTypeName("bpchar")
			}
			$$.Typmods = makeList(&nodes.Integer{Ival: $4})
		}
	| CHARACTER opt_varying
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				/* ConstCharacter: CHAR without length defaults to unspecified */
				$$ = makeTypeName("bpchar")
			}
		}
	| CHAR_P opt_varying '(' Iconst ')'
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				$$ = makeTypeName("bpchar")
			}
			$$.Typmods = makeList(&nodes.Integer{Ival: $4})
		}
	| CHAR_P opt_varying
		{
			if $2 {
				$$ = makeTypeName("varchar")
			} else {
				$$ = makeTypeName("bpchar")
			}
		}
	| VARCHAR '(' Iconst ')'
		{
			$$ = makeTypeName("varchar")
			$$.Typmods = makeList(&nodes.Integer{Ival: $3})
		}
	| VARCHAR
		{
			$$ = makeTypeName("varchar")
		}
	;

ConstDatetime:
	TIMESTAMP '(' Iconst ')' opt_timezone
		{
			if $5 {
				$$ = makeTypeName("timestamptz")
			} else {
				$$ = makeTypeName("timestamp")
			}
			$$.Typmods = makeList(makeIntConst($3))
		}
	| TIMESTAMP opt_timezone
		{
			if $2 {
				$$ = makeTypeName("timestamptz")
			} else {
				$$ = makeTypeName("timestamp")
			}
		}
	| TIME '(' Iconst ')' opt_timezone
		{
			if $5 {
				$$ = makeTypeName("timetz")
			} else {
				$$ = makeTypeName("time")
			}
			$$.Typmods = makeList(makeIntConst($3))
		}
	| TIME opt_timezone
		{
			if $2 {
				$$ = makeTypeName("timetz")
			} else {
				$$ = makeTypeName("time")
			}
		}
	;

ConstInterval:
	INTERVAL
		{
			$$ = makeTypeName("interval")
		}
	;

opt_timezone:
	WITH_LA TIME ZONE		{ $$ = true }
	| WITHOUT TIME ZONE		{ $$ = false }
	| /* EMPTY */			{ $$ = false }
	;

opt_interval:
	YEAR_P
		{ $$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_YEAR))) }
	| MONTH_P
		{ $$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_MONTH))) }
	| DAY_P
		{ $$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_DAY))) }
	| HOUR_P
		{ $$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_HOUR))) }
	| MINUTE_P
		{ $$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_MINUTE))) }
	| interval_second
		{ $$ = $1 }
	| YEAR_P TO MONTH_P
		{
			$$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_YEAR | nodes.INTERVAL_MASK_MONTH)))
		}
	| DAY_P TO HOUR_P
		{
			$$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_DAY | nodes.INTERVAL_MASK_HOUR)))
		}
	| DAY_P TO MINUTE_P
		{
			$$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_DAY | nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE)))
		}
	| DAY_P TO interval_second
		{
			$$ = $3
			$$.Items[0] = makeIntConst(int64(nodes.INTERVAL_MASK_DAY | nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE | nodes.INTERVAL_MASK_SECOND))
		}
	| HOUR_P TO MINUTE_P
		{
			$$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE)))
		}
	| HOUR_P TO interval_second
		{
			$$ = $3
			$$.Items[0] = makeIntConst(int64(nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE | nodes.INTERVAL_MASK_SECOND))
		}
	| MINUTE_P TO interval_second
		{
			$$ = $3
			$$.Items[0] = makeIntConst(int64(nodes.INTERVAL_MASK_MINUTE | nodes.INTERVAL_MASK_SECOND))
		}
	| /* EMPTY */
		{ $$ = nil }
	;

interval_second:
	SECOND_P
		{
			$$ = makeList(makeIntConst(int64(nodes.INTERVAL_MASK_SECOND)))
		}
	| SECOND_P '(' Iconst ')'
		{
			$$ = makeList2(makeIntConst(int64(nodes.INTERVAL_MASK_SECOND)), makeIntConst($3))
		}
	;

/*****************************************************************************
 *
 *      CHECKPOINT
 *
 *****************************************************************************/

CheckPointStmt:
	CHECKPOINT
		{
			$$ = &nodes.CheckPointStmt{}
		}
	;

/*****************************************************************************
 *
 *      DISCARD
 *
 *****************************************************************************/

DiscardStmt:
	DISCARD ALL
		{
			$$ = &nodes.DiscardStmt{Target: nodes.DISCARD_ALL}
		}
	| DISCARD TEMP
		{
			$$ = &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP}
		}
	| DISCARD TEMPORARY
		{
			$$ = &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP}
		}
	| DISCARD PLANS
		{
			$$ = &nodes.DiscardStmt{Target: nodes.DISCARD_PLANS}
		}
	| DISCARD SEQUENCES
		{
			$$ = &nodes.DiscardStmt{Target: nodes.DISCARD_SEQUENCES}
		}
	;

/*****************************************************************************
 *
 *      LISTEN
 *
 *****************************************************************************/

ListenStmt:
	LISTEN ColId
		{
			$$ = &nodes.ListenStmt{Conditionname: $2}
		}
	;

/*****************************************************************************
 *
 *      UNLISTEN
 *
 *****************************************************************************/

UnlistenStmt:
	UNLISTEN ColId
		{
			$$ = &nodes.UnlistenStmt{Conditionname: $2}
		}
	| UNLISTEN '*'
		{
			$$ = &nodes.UnlistenStmt{Conditionname: ""}
		}
	;

/*****************************************************************************
 *
 *      NOTIFY
 *
 *****************************************************************************/

NotifyStmt:
	NOTIFY ColId
		{
			$$ = &nodes.NotifyStmt{Conditionname: $2}
		}
	| NOTIFY ColId ',' Sconst
		{
			$$ = &nodes.NotifyStmt{Conditionname: $2, Payload: $4}
		}
	;

/*****************************************************************************
 *
 *      LOAD
 *
 *****************************************************************************/

LoadStmt:
	LOAD file_name
		{
			$$ = &nodes.LoadStmt{Filename: $2}
		}
	;

file_name:
	Sconst { $$ = $1 }
	;

/*****************************************************************************
 *
 *      CLOSE cursor
 *
 *****************************************************************************/

ClosePortalStmt:
	CLOSE cursor_name
		{
			$$ = &nodes.ClosePortalStmt{Portalname: $2}
		}
	| CLOSE ALL
		{
			$$ = &nodes.ClosePortalStmt{Portalname: ""}
		}
	;

cursor_name:
	name { $$ = $1 }
	;

/*****************************************************************************
 *
 *      SET CONSTRAINTS
 *
 *****************************************************************************/

ConstraintsSetStmt:
	SET CONSTRAINTS constraints_set_list constraints_set_mode
		{
			$$ = &nodes.ConstraintsSetStmt{
				Constraints: $3,
				Deferred:    $4,
			}
		}
	;

constraints_set_list:
	ALL
		{
			$$ = nil
		}
	| qualified_name_list
		{
			$$ = $1
		}
	;

constraints_set_mode:
	DEFERRED  { $$ = true }
	| IMMEDIATE { $$ = false }
	;

qualified_name_list:
	qualified_name
		{
			$$ = makeList(makeRangeVar($1))
		}
	| qualified_name_list ',' qualified_name
		{
			$$ = appendList($1, makeRangeVar($3))
		}
	;

/*****************************************************************************
 *
 *      SET variable
 *
 *****************************************************************************/

VariableSetStmt:
	SET set_rest
		{
			n := $2.(*nodes.VariableSetStmt)
			n.IsLocal = false
			$$ = n
		}
	| SET LOCAL set_rest
		{
			n := $3.(*nodes.VariableSetStmt)
			n.IsLocal = true
			$$ = n
		}
	| SET SESSION set_rest
		{
			n := $3.(*nodes.VariableSetStmt)
			n.IsLocal = false
			$$ = n
		}
	;

set_rest:
	TRANSACTION transaction_mode_list
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_MULTI,
				Name: "TRANSACTION",
				Args: $2,
			}
		}
	| SESSION CHARACTERISTICS AS TRANSACTION transaction_mode_list
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_MULTI,
				Name: "SESSION CHARACTERISTICS",
				Args: $5,
			}
		}
	| set_rest_more
		{
			$$ = $1
		}
	;

generic_set:
	var_name TO var_list
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: $1,
				Args: $3,
			}
		}
	| var_name '=' var_list
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: $1,
				Args: $3,
			}
		}
	| var_name TO DEFAULT
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_DEFAULT,
				Name: $1,
			}
		}
	| var_name '=' DEFAULT
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_DEFAULT,
				Name: $1,
			}
		}
	;

set_rest_more:
	generic_set
		{
			$$ = $1
		}
	| var_name FROM CURRENT_P
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_CURRENT,
				Name: $1,
			}
		}
	| TIME ZONE zone_value
		{
			n := &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: "timezone",
			}
			if $3 != nil {
				n.Args = makeList($3)
			} else {
				n.Kind = nodes.VAR_SET_DEFAULT
			}
			$$ = n
		}
	| CATALOG_P Sconst
		{
			pglex.Error("current database cannot be changed")
			$$ = nil
		}
	| SCHEMA Sconst
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: "search_path",
				Args: makeList(makeStringConst($2)),
			}
		}
	| NAMES opt_encoding
		{
			n := &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: "client_encoding",
			}
			if $2 != "" {
				n.Args = makeList(makeStringConst($2))
			} else {
				n.Kind = nodes.VAR_SET_DEFAULT
			}
			$$ = n
		}
	| ROLE NonReservedWord_or_Sconst
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: "role",
				Args: makeList(makeStringConst($2)),
			}
		}
	| SESSION AUTHORIZATION NonReservedWord_or_Sconst
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: "session_authorization",
				Args: makeList(makeStringConst($3)),
			}
		}
	| SESSION AUTHORIZATION DEFAULT
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_DEFAULT,
				Name: "session_authorization",
			}
		}
	| XML_P OPTION document_or_content
		{
			var val string
			if $3 == int64(nodes.XMLOPTION_DOCUMENT) {
				val = "DOCUMENT"
			} else {
				val = "CONTENT"
			}
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_VALUE,
				Name: "xmloption",
				Args: makeList(makeStringConst(val)),
			}
		}
	| TRANSACTION SNAPSHOT Sconst
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_SET_MULTI,
				Name: "TRANSACTION SNAPSHOT",
				Args: makeList(makeStringConst($3)),
			}
		}
	;

var_name:
	ColId
		{
			$$ = $1
		}
	| var_name '.' ColId
		{
			$$ = $1 + "." + $3
		}
	;

var_list:
	var_value
		{
			$$ = makeList($1)
		}
	| var_list ',' var_value
		{
			$$ = appendList($1, $3)
		}
	;

var_value:
	opt_boolean_or_string
		{
			$$ = makeStringConst($1)
		}
	| NumericOnly
		{
			$$ = &nodes.A_Const{Val: $1}
		}
	;

zone_value:
	Sconst
		{
			$$ = makeStringConst($1)
		}
	| IDENT
		{
			$$ = makeStringConst($1)
		}
	| NumericOnly
		{
			$$ = &nodes.A_Const{Val: $1}
		}
	| DEFAULT
		{
			$$ = nil
		}
	| LOCAL
		{
			$$ = nil
		}
	;

opt_encoding:
	Sconst      { $$ = $1 }
	| DEFAULT   { $$ = "" }
	| /* EMPTY */ { $$ = "" }
	;

iso_level:
	READ UNCOMMITTED  { $$ = "read uncommitted" }
	| READ COMMITTED  { $$ = "read committed" }
	| REPEATABLE READ { $$ = "repeatable read" }
	| SERIALIZABLE    { $$ = "serializable" }
	;

transaction_mode_item:
	ISOLATION LEVEL iso_level
		{
			$$ = makeDefElem("transaction_isolation", makeStringConst($3))
		}
	| READ ONLY
		{
			$$ = makeDefElem("transaction_read_only", makeIntConst(1))
		}
	| READ WRITE
		{
			$$ = makeDefElem("transaction_read_only", makeIntConst(0))
		}
	| DEFERRABLE
		{
			$$ = makeDefElem("transaction_deferrable", makeIntConst(1))
		}
	| NOT DEFERRABLE
		{
			$$ = makeDefElem("transaction_deferrable", makeIntConst(0))
		}
	;

transaction_mode_list:
	transaction_mode_item
		{
			$$ = makeList($1)
		}
	| transaction_mode_list ',' transaction_mode_item
		{
			$$ = appendList($1, $3)
		}
	| transaction_mode_list transaction_mode_item
		{
			$$ = appendList($1, $2)
		}
	;

transaction_mode_list_or_empty:
	transaction_mode_list { $$ = $1 }
	| /* EMPTY */         { $$ = nil }
	;

/*****************************************************************************
 *
 *      SHOW variable
 *
 *****************************************************************************/

VariableShowStmt:
	SHOW var_name
		{
			$$ = &nodes.VariableShowStmt{
				Name: $2,
			}
		}
	| SHOW TIME ZONE
		{
			$$ = &nodes.VariableShowStmt{
				Name: "timezone",
			}
		}
	| SHOW TRANSACTION ISOLATION LEVEL
		{
			$$ = &nodes.VariableShowStmt{
				Name: "transaction_isolation",
			}
		}
	| SHOW SESSION AUTHORIZATION
		{
			$$ = &nodes.VariableShowStmt{
				Name: "session_authorization",
			}
		}
	| SHOW ALL
		{
			$$ = &nodes.VariableShowStmt{
				Name: "all",
			}
		}
	;

/*****************************************************************************
 *
 *      RESET variable
 *
 *****************************************************************************/

VariableResetStmt:
	RESET reset_rest
		{
			$$ = $2
		}
	;

reset_rest:
	generic_reset
		{
			$$ = $1
		}
	| TIME ZONE
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_RESET,
				Name: "timezone",
			}
		}
	| TRANSACTION ISOLATION LEVEL
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_RESET,
				Name: "transaction_isolation",
			}
		}
	| SESSION AUTHORIZATION
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_RESET,
				Name: "session_authorization",
			}
		}
	;

generic_reset:
	var_name
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_RESET,
				Name: $1,
			}
		}
	| ALL
		{
			$$ = &nodes.VariableSetStmt{
				Kind: nodes.VAR_RESET_ALL,
			}
		}
	;

/*****************************************************************************
 *
 *      PREPARE / EXECUTE / DEALLOCATE
 *
 *****************************************************************************/

PrepareStmt:
	PREPARE name prep_type_clause AS PreparableStmt
		{
			$$ = &nodes.PrepareStmt{
				Name:     $2,
				Argtypes: $3,
				Query:    $5,
			}
		}
	;

prep_type_clause:
	'(' type_list ')'  { $$ = $2 }
	| /* EMPTY */       { $$ = nil }
	;

type_list:
	Typename
		{
			$$ = makeList($1)
		}
	| type_list ',' Typename
		{
			$$ = appendList($1, $3)
		}
	;

PreparableStmt:
	SelectStmt  { $$ = $1 }
	| InsertStmt { $$ = $1 }
	| UpdateStmt { $$ = $1 }
	| DeleteStmt { $$ = $1 }
	| MergeStmt  { $$ = $1 }
	;

ExecuteStmt:
	EXECUTE name execute_param_clause
		{
			$$ = &nodes.ExecuteStmt{
				Name:   $2,
				Params: $3,
			}
		}
	;

execute_param_clause:
	'(' expr_list ')'  { $$ = $2 }
	| /* EMPTY */       { $$ = nil }
	;

DeallocateStmt:
	DEALLOCATE name
		{
			$$ = &nodes.DeallocateStmt{
				Name: $2,
			}
		}
	| DEALLOCATE PREPARE name
		{
			$$ = &nodes.DeallocateStmt{
				Name: $3,
			}
		}
	| DEALLOCATE ALL
		{
			$$ = &nodes.DeallocateStmt{
				IsAll: true,
			}
		}
	| DEALLOCATE PREPARE ALL
		{
			$$ = &nodes.DeallocateStmt{
				IsAll: true,
			}
		}
	;

/*****************************************************************************
 *
 *      TRUNCATE TABLE
 *
 *****************************************************************************/

TruncateStmt:
	TRUNCATE opt_table relation_expr_list opt_restart_seqs opt_drop_behavior
		{
			$$ = &nodes.TruncateStmt{
				Relations:   $3,
				RestartSeqs: $4,
				Behavior:    nodes.DropBehavior($5),
			}
		}
	;

opt_restart_seqs:
	CONTINUE_P IDENTITY_P { $$ = false }
	| RESTART IDENTITY_P  { $$ = true }
	| /* EMPTY */          { $$ = false }
	;

opt_table:
	TABLE
	| /* EMPTY */
	;

relation_expr_list:
	relation_expr
		{
			$$ = makeList($1)
		}
	| relation_expr_list ',' relation_expr
		{
			$$ = appendList($1, $3)
		}
	;

/*****************************************************************************
 *
 *      LOCK TABLE
 *
 *****************************************************************************/

LockStmt:
	LOCK_P opt_table relation_expr_list opt_lock opt_nowait
		{
			$$ = &nodes.LockStmt{
				Relations: $3,
				Mode:      int($4),
				Nowait:    $5,
			}
		}
	;

opt_lock:
	IN_P lock_type MODE { $$ = $2 }
	| /* EMPTY */       { $$ = int64(nodes.AccessExclusiveLock) }
	;

lock_type:
	ACCESS SHARE                 { $$ = int64(nodes.AccessShareLock) }
	| ROW SHARE                  { $$ = int64(nodes.RowShareLock) }
	| ROW EXCLUSIVE              { $$ = int64(nodes.RowExclusiveLock) }
	| SHARE UPDATE EXCLUSIVE     { $$ = int64(nodes.ShareUpdateExclusiveLock) }
	| SHARE                      { $$ = int64(nodes.ShareLock) }
	| SHARE ROW EXCLUSIVE        { $$ = int64(nodes.ShareRowExclusiveLock) }
	| EXCLUSIVE                  { $$ = int64(nodes.ExclusiveLock) }
	| ACCESS EXCLUSIVE           { $$ = int64(nodes.AccessExclusiveLock) }
	;

opt_nowait:
	NOWAIT      { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

/*****************************************************************************
 *
 *      VACUUM / ANALYZE
 *
 *****************************************************************************/

VacuumStmt:
	VACUUM opt_full opt_freeze opt_verbose opt_analyze opt_vacuum_relation_list
		{
			n := &nodes.VacuumStmt{
				IsVacuumCmd: true,
				Rels:        $6,
			}
			var opts []nodes.Node
			if $2 {
				opts = append(opts, makeDefElem("full", nil))
			}
			if $3 {
				opts = append(opts, makeDefElem("freeze", nil))
			}
			if $4 {
				opts = append(opts, makeDefElem("verbose", nil))
			}
			if $5 {
				opts = append(opts, makeDefElem("analyze", nil))
			}
			if len(opts) > 0 {
				n.Options = &nodes.List{Items: opts}
			}
			$$ = n
		}
	| VACUUM '(' utility_option_list ')' opt_vacuum_relation_list
		{
			$$ = &nodes.VacuumStmt{
				Options:     $3,
				Rels:        $5,
				IsVacuumCmd: true,
			}
		}
	;

AnalyzeStmt:
	analyze_keyword opt_verbose opt_vacuum_relation_list
		{
			n := &nodes.VacuumStmt{
				IsVacuumCmd: false,
				Rels:        $3,
			}
			if $2 {
				n.Options = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
			}
			$$ = n
		}
	| analyze_keyword '(' utility_option_list ')' opt_vacuum_relation_list
		{
			$$ = &nodes.VacuumStmt{
				Options:     $3,
				Rels:        $5,
				IsVacuumCmd: false,
			}
		}
	;

analyze_keyword:
	ANALYZE
	| ANALYSE
	;

opt_full:
	FULL        { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

opt_freeze:
	FREEZE      { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

opt_verbose:
	VERBOSE     { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

opt_analyze:
	analyze_keyword { $$ = true }
	| /* EMPTY */   { $$ = false }
	;

vacuum_relation:
	qualified_name opt_column_list
		{
			$$ = &nodes.VacuumRelation{
				Relation: makeRangeVar($1).(*nodes.RangeVar),
				VaCols:   $2,
			}
		}
	;

vacuum_relation_list:
	vacuum_relation
		{
			$$ = makeList($1)
		}
	| vacuum_relation_list ',' vacuum_relation
		{
			$$ = appendList($1, $3)
		}
	;

opt_vacuum_relation_list:
	vacuum_relation_list { $$ = $1 }
	| /* EMPTY */        { $$ = nil }
	;

/*****************************************************************************
 *
 *      CLUSTER
 *
 *****************************************************************************/

ClusterStmt:
	CLUSTER '(' utility_option_list ')' qualified_name cluster_index_specification
		{
			$$ = &nodes.ClusterStmt{
				Relation: makeRangeVar($5).(*nodes.RangeVar),
				Indexname: $6,
				Params:   $3,
			}
		}
	| CLUSTER '(' utility_option_list ')'
		{
			$$ = &nodes.ClusterStmt{
				Params: $3,
			}
		}
	| CLUSTER opt_verbose qualified_name cluster_index_specification
		{
			n := &nodes.ClusterStmt{
				Relation: makeRangeVar($3).(*nodes.RangeVar),
				Indexname: $4,
			}
			if $2 {
				n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
			}
			$$ = n
		}
	| CLUSTER opt_verbose
		{
			n := &nodes.ClusterStmt{}
			if $2 {
				n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
			}
			$$ = n
		}
	/* backwards compatible syntax from before 8.3 */
	| CLUSTER opt_verbose name ON qualified_name
		{
			n := &nodes.ClusterStmt{
				Relation:  makeRangeVar($5).(*nodes.RangeVar),
				Indexname: $3,
			}
			if $2 {
				n.Params = &nodes.List{Items: []nodes.Node{makeDefElem("verbose", nil)}}
			}
			$$ = n
		}
	;

cluster_index_specification:
	USING name   { $$ = $2 }
	| /* EMPTY */ { $$ = "" }
	;

/*****************************************************************************
 *
 *      REINDEX
 *
 *****************************************************************************/

ReindexStmt:
	REINDEX opt_reindex_option_list reindex_target_type opt_concurrently qualified_name
		{
			n := &nodes.ReindexStmt{
				Kind:     nodes.ReindexObjectType($3),
				Relation: makeRangeVar($5).(*nodes.RangeVar),
				Params:   $2,
			}
			if $4 {
				n.Params = appendList(n.Params, makeDefElem("concurrently", nil))
			}
			$$ = n
		}
	| REINDEX opt_reindex_option_list SCHEMA opt_concurrently name
		{
			n := &nodes.ReindexStmt{
				Kind:   nodes.REINDEX_OBJECT_SCHEMA,
				Name:   $5,
				Params: $2,
			}
			if $4 {
				n.Params = appendList(n.Params, makeDefElem("concurrently", nil))
			}
			$$ = n
		}
	| REINDEX opt_reindex_option_list reindex_target_multitable opt_concurrently name
		{
			n := &nodes.ReindexStmt{
				Kind:   nodes.ReindexObjectType($3),
				Name:   $5,
				Params: $2,
			}
			if $4 {
				n.Params = appendList(n.Params, makeDefElem("concurrently", nil))
			}
			$$ = n
		}
	;

reindex_target_type:
	INDEX { $$ = int64(nodes.REINDEX_OBJECT_INDEX) }
	| TABLE { $$ = int64(nodes.REINDEX_OBJECT_TABLE) }
	;

reindex_target_multitable:
	SYSTEM_P { $$ = int64(nodes.REINDEX_OBJECT_SYSTEM) }
	| DATABASE { $$ = int64(nodes.REINDEX_OBJECT_DATABASE) }
	;

opt_reindex_option_list:
	'(' utility_option_list ')' { $$ = $2 }
	| /* EMPTY */               { $$ = nil }
	;

/*****************************************************************************
 *
 *      COMMENT ON
 *
 *****************************************************************************/

CommentStmt:
	COMMENT ON object_type_any_name any_name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.ObjectType($3),
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON COLUMN any_name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_COLUMN,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON object_type_name name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.ObjectType($3),
				Object:  &nodes.String{Str: $4},
				Comment: $6,
			}
		}
	| COMMENT ON TYPE_P any_name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_TYPE,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON DOMAIN_P any_name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_DOMAIN,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON AGGREGATE aggregate_with_argtypes IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_AGGREGATE,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON FUNCTION function_with_argtypes IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_FUNCTION,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON PROCEDURE function_with_argtypes IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_PROCEDURE,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON ROUTINE function_with_argtypes IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_ROUTINE,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON OPERATOR operator_with_argtypes IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_OPERATOR,
				Object:  $4,
				Comment: $6,
			}
		}
	| COMMENT ON CONSTRAINT name ON any_name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_TABCONSTRAINT,
				Object:  appendList($6, &nodes.String{Str: $4}),
				Comment: $8,
			}
		}
	| COMMENT ON CONSTRAINT name ON DOMAIN_P any_name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_DOMCONSTRAINT,
				Object:  &nodes.List{Items: []nodes.Node{makeTypeNameFromNameList($7), &nodes.String{Str: $4}}},
				Comment: $9,
			}
		}
	| COMMENT ON object_type_name_on_any_name name ON any_name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.ObjectType($3),
				Object:  appendList($6, &nodes.String{Str: $4}),
				Comment: $8,
			}
		}
	| COMMENT ON TRANSFORM FOR Typename LANGUAGE name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_TRANSFORM,
				Object:  &nodes.List{Items: []nodes.Node{$5, &nodes.String{Str: $7}}},
				Comment: $9,
			}
		}
	| COMMENT ON OPERATOR CLASS any_name USING name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_OPCLASS,
				Object:  prependList(&nodes.String{Str: $7}, $5),
				Comment: $9,
			}
		}
	| COMMENT ON OPERATOR FAMILY any_name USING name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_OPFAMILY,
				Object:  prependList(&nodes.String{Str: $7}, $5),
				Comment: $9,
			}
		}
	| COMMENT ON LARGE_P OBJECT_P NumericOnly IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_LARGEOBJECT,
				Object:  $5,
				Comment: $7,
			}
		}
	| COMMENT ON FOREIGN DATA_P WRAPPER name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_FDW,
				Object:  &nodes.String{Str: $6},
				Comment: $8,
			}
		}
	| COMMENT ON CAST '(' Typename AS Typename ')' IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_CAST,
				Object:  &nodes.List{Items: []nodes.Node{$5, $7}},
				Comment: $10,
			}
		}
	| COMMENT ON EVENT TRIGGER name IS comment_text
		{
			$$ = &nodes.CommentStmt{
				Objtype: nodes.OBJECT_EVENT_TRIGGER,
				Object:  &nodes.String{Str: $5},
				Comment: $7,
			}
		}
	;

comment_text:
	Sconst      { $$ = $1 }
	| NULL_P    { $$ = "" }
	;

object_type_name:
	SCHEMA       { $$ = int64(nodes.OBJECT_SCHEMA) }
	| DATABASE   { $$ = int64(nodes.OBJECT_DATABASE) }
	| ROLE       { $$ = int64(nodes.OBJECT_ROLE) }
	| TABLESPACE { $$ = int64(nodes.OBJECT_TABLESPACE) }
	| SUBSCRIPTION { $$ = int64(nodes.OBJECT_SUBSCRIPTION) }
	| PUBLICATION  { $$ = int64(nodes.OBJECT_PUBLICATION) }
	| SERVER       { $$ = int64(nodes.OBJECT_FOREIGN_SERVER) }
	;

object_type_name_on_any_name:
	POLICY      { $$ = int64(nodes.OBJECT_POLICY) }
	| RULE      { $$ = int64(nodes.OBJECT_RULE) }
	| TRIGGER   { $$ = int64(nodes.OBJECT_TRIGGER) }
	;

/*****************************************************************************
 *
 *      SECURITY LABEL
 *
 *****************************************************************************/

SecLabelStmt:
	SECURITY LABEL opt_provider ON object_type_any_name any_name IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.ObjectType($5),
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
	| SECURITY LABEL opt_provider ON COLUMN any_name IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_COLUMN,
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
	| SECURITY LABEL opt_provider ON object_type_name name IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.ObjectType($5),
				Object:   &nodes.String{Str: $6},
				Provider: $3,
				Label:    $8,
			}
		}
	| SECURITY LABEL opt_provider ON TYPE_P any_name IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_TYPE,
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
	| SECURITY LABEL opt_provider ON DOMAIN_P any_name IS security_label
		{
			$$ = &nodes.SecLabelStmt{
				Objtype:  nodes.OBJECT_DOMAIN,
				Object:   $6,
				Provider: $3,
				Label:    $8,
			}
		}
	/* TODO: Add SECURITY LABEL ON FUNCTION/AGGREGATE variants in Phase 16 */
	;

/*****************************************************************************
 *
 *      DECLARE CURSOR
 *
 *****************************************************************************/

DeclareCursorStmt:
	DECLARE cursor_name cursor_options CURSOR opt_hold FOR SelectStmt
		{
			$$ = &nodes.DeclareCursorStmt{
				Portalname: $2,
				Options:    int($3 | $5 | nodes.CURSOR_OPT_FAST_PLAN),
				Query:      $7,
			}
		}
	;

cursor_options:
	/* EMPTY */
		{
			$$ = 0
		}
	| cursor_options NO SCROLL
		{
			$$ = $1 | nodes.CURSOR_OPT_NO_SCROLL
		}
	| cursor_options SCROLL
		{
			$$ = $1 | nodes.CURSOR_OPT_SCROLL
		}
	| cursor_options BINARY
		{
			$$ = $1 | nodes.CURSOR_OPT_BINARY
		}
	| cursor_options ASENSITIVE
		{
			$$ = $1 | nodes.CURSOR_OPT_ASENSITIVE
		}
	| cursor_options INSENSITIVE
		{
			$$ = $1 | nodes.CURSOR_OPT_INSENSITIVE
		}
	;

opt_hold:
	/* EMPTY */
		{
			$$ = 0
		}
	| WITH HOLD
		{
			$$ = nodes.CURSOR_OPT_HOLD
		}
	| WITHOUT HOLD
		{
			$$ = 0
		}
	;

/*****************************************************************************
 *
 *      FETCH / MOVE
 *
 *****************************************************************************/

FetchStmt:
	FETCH fetch_args
		{
			n := $2.(*nodes.FetchStmt)
			n.Ismove = false
			$$ = n
		}
	| MOVE fetch_args
		{
			n := $2.(*nodes.FetchStmt)
			n.Ismove = true
			$$ = n
		}
	;

fetch_args:
	cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   1,
				Portalname: $1,
			}
		}
	| from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   1,
				Portalname: $2,
			}
		}
	| NEXT opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   1,
				Portalname: $3,
			}
		}
	| PRIOR opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_BACKWARD,
				HowMany:   1,
				Portalname: $3,
			}
		}
	| FIRST_P opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_ABSOLUTE,
				HowMany:   1,
				Portalname: $3,
			}
		}
	| LAST_P opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_ABSOLUTE,
				HowMany:   -1,
				Portalname: $3,
			}
		}
	| ABSOLUTE_P SignedIconst opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_ABSOLUTE,
				HowMany:   $2,
				Portalname: $4,
			}
		}
	| RELATIVE_P SignedIconst opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_RELATIVE,
				HowMany:   $2,
				Portalname: $4,
			}
		}
	| SignedIconst opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   $1,
				Portalname: $3,
			}
		}
	| ALL opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   nodes.FETCH_ALL,
				Portalname: $3,
			}
		}
	| FORWARD opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   1,
				Portalname: $3,
			}
		}
	| FORWARD SignedIconst opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   $2,
				Portalname: $4,
			}
		}
	| FORWARD ALL opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_FORWARD,
				HowMany:   nodes.FETCH_ALL,
				Portalname: $4,
			}
		}
	| BACKWARD opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_BACKWARD,
				HowMany:   1,
				Portalname: $3,
			}
		}
	| BACKWARD SignedIconst opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_BACKWARD,
				HowMany:   $2,
				Portalname: $4,
			}
		}
	| BACKWARD ALL opt_from_in cursor_name
		{
			$$ = &nodes.FetchStmt{
				Direction:  nodes.FETCH_BACKWARD,
				HowMany:   nodes.FETCH_ALL,
				Portalname: $4,
			}
		}
	;

from_in:
	FROM
	| IN_P
	;

opt_from_in:
	from_in
	| /* EMPTY */
	;

/*****************************************************************************
 *
 *      MERGE INTO
 *
 *****************************************************************************/

MergeStmt:
	opt_with_clause MERGE INTO relation_expr_opt_alias USING table_ref ON a_expr merge_when_list returning_clause
		{
			m := &nodes.MergeStmt{}
			if $1 != nil {
				m.WithClause = $1.(*nodes.WithClause)
			}
			m.Relation = $4.(*nodes.RangeVar)
			m.SourceRelation = $6
			m.JoinCondition = $8
			m.MergeWhenClauses = $9
			m.ReturningList = $10
			$$ = m
		}
	;

merge_when_list:
	merge_when_clause
		{
			$$ = makeList($1)
		}
	| merge_when_list merge_when_clause
		{
			$$ = appendList($1, $2)
		}
	;

merge_when_clause:
	merge_when_tgt_matched opt_merge_when_condition THEN merge_update
		{
			n := $4.(*nodes.MergeWhenClause)
			n.Kind = nodes.MergeMatchKind($1)
			n.Condition = $2
			$$ = n
		}
	| merge_when_tgt_matched opt_merge_when_condition THEN merge_delete
		{
			n := $4.(*nodes.MergeWhenClause)
			n.Kind = nodes.MergeMatchKind($1)
			n.Condition = $2
			$$ = n
		}
	| merge_when_tgt_not_matched opt_merge_when_condition THEN merge_insert
		{
			n := $4.(*nodes.MergeWhenClause)
			n.Kind = nodes.MergeMatchKind($1)
			n.Condition = $2
			$$ = n
		}
	| merge_when_tgt_matched opt_merge_when_condition THEN DO NOTHING
		{
			$$ = &nodes.MergeWhenClause{
				Kind:        nodes.MergeMatchKind($1),
				CommandType: nodes.CMD_NOTHING,
				Condition:   $2,
			}
		}
	| merge_when_tgt_not_matched opt_merge_when_condition THEN DO NOTHING
		{
			$$ = &nodes.MergeWhenClause{
				Kind:        nodes.MergeMatchKind($1),
				CommandType: nodes.CMD_NOTHING,
				Condition:   $2,
			}
		}
	;

merge_when_tgt_matched:
	WHEN MATCHED
		{
			$$ = int64(nodes.MERGE_WHEN_MATCHED)
		}
	| WHEN NOT MATCHED BY SOURCE
		{
			$$ = int64(nodes.MERGE_WHEN_NOT_MATCHED_BY_SOURCE)
		}
	;

merge_when_tgt_not_matched:
	WHEN NOT MATCHED
		{
			$$ = int64(nodes.MERGE_WHEN_NOT_MATCHED_BY_TARGET)
		}
	| WHEN NOT MATCHED BY TARGET
		{
			$$ = int64(nodes.MERGE_WHEN_NOT_MATCHED_BY_TARGET)
		}
	;

opt_merge_when_condition:
	AND a_expr
		{
			$$ = $2
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

merge_update:
	UPDATE SET set_clause_list
		{
			$$ = &nodes.MergeWhenClause{
				CommandType: nodes.CMD_UPDATE,
				Override:    nodes.OVERRIDING_NOT_SET,
				TargetList:  $3,
			}
		}
	;

merge_delete:
	DELETE_P
		{
			$$ = &nodes.MergeWhenClause{
				CommandType: nodes.CMD_DELETE,
				Override:    nodes.OVERRIDING_NOT_SET,
			}
		}
	;

merge_insert:
	INSERT merge_values_clause
		{
			$$ = &nodes.MergeWhenClause{
				CommandType: nodes.CMD_INSERT,
				Override:    nodes.OVERRIDING_NOT_SET,
				Values:      $2,
			}
		}
	| INSERT OVERRIDING override_kind VALUE_P merge_values_clause
		{
			$$ = &nodes.MergeWhenClause{
				CommandType: nodes.CMD_INSERT,
				Override:    nodes.OverridingKind($3),
				Values:      $5,
			}
		}
	| INSERT '(' insert_column_list ')' merge_values_clause
		{
			$$ = &nodes.MergeWhenClause{
				CommandType: nodes.CMD_INSERT,
				Override:    nodes.OVERRIDING_NOT_SET,
				TargetList:  $3,
				Values:      $5,
			}
		}
	| INSERT '(' insert_column_list ')' OVERRIDING override_kind VALUE_P merge_values_clause
		{
			$$ = &nodes.MergeWhenClause{
				CommandType: nodes.CMD_INSERT,
				Override:    nodes.OverridingKind($6),
				TargetList:  $3,
				Values:      $8,
			}
		}
	| INSERT DEFAULT VALUES
		{
			$$ = &nodes.MergeWhenClause{
				CommandType: nodes.CMD_INSERT,
				Override:    nodes.OVERRIDING_NOT_SET,
			}
		}
	;

merge_values_clause:
	VALUES '(' expr_list ')'
		{
			$$ = $3
		}
	;

override_kind:
	USER
		{
			$$ = int64(nodes.OVERRIDING_USER_VALUE)
		}
	| SYSTEM_P
		{
			$$ = int64(nodes.OVERRIDING_SYSTEM_VALUE)
		}
	;

/*****************************************************************************
 *
 *      CALL
 *
 *****************************************************************************/

CallStmt:
	CALL func_application
		{
			$$ = &nodes.CallStmt{
				Funccall: $2.(*nodes.FuncCall),
			}
		}
	;

/*****************************************************************************
 *
 *      DO (anonymous code block)
 *
 *****************************************************************************/

DoStmt:
	DO dostmt_opt_list
		{
			$$ = &nodes.DoStmt{
				Args: $2,
			}
		}
	;

dostmt_opt_list:
	dostmt_opt_item
		{
			$$ = makeList($1)
		}
	| dostmt_opt_list dostmt_opt_item
		{
			$$ = appendList($1, $2)
		}
	;

dostmt_opt_item:
	Sconst
		{
			$$ = &nodes.DefElem{
				Defname: "as",
				Arg:     &nodes.String{Str: $1},
			}
		}
	| LANGUAGE NonReservedWord_or_Sconst
		{
			$$ = &nodes.DefElem{
				Defname: "language",
				Arg:     &nodes.String{Str: $2},
			}
		}
	;

opt_provider:
	FOR NonReservedWord_or_Sconst { $$ = $2 }
	| /* EMPTY */                  { $$ = "" }
	;

security_label:
	Sconst      { $$ = $1 }
	| NULL_P    { $$ = "" }
	;

/*****************************************************************************
 *
 *      ALTER FUNCTION / ALTER PROCEDURE / ALTER ROUTINE
 *
 *****************************************************************************/

AlterFunctionStmt:
	ALTER FUNCTION function_with_argtypes alterfunc_opt_list opt_restrict
		{
			$$ = &nodes.AlterFunctionStmt{
				Objtype: nodes.OBJECT_FUNCTION,
				Func:    $3.(*nodes.ObjectWithArgs),
				Actions: $4,
			}
		}
	| ALTER PROCEDURE function_with_argtypes alterfunc_opt_list opt_restrict
		{
			$$ = &nodes.AlterFunctionStmt{
				Objtype: nodes.OBJECT_PROCEDURE,
				Func:    $3.(*nodes.ObjectWithArgs),
				Actions: $4,
			}
		}
	| ALTER ROUTINE function_with_argtypes alterfunc_opt_list opt_restrict
		{
			$$ = &nodes.AlterFunctionStmt{
				Objtype: nodes.OBJECT_ROUTINE,
				Func:    $3.(*nodes.ObjectWithArgs),
				Actions: $4,
			}
		}
	;

alterfunc_opt_list:
	common_func_opt_item
		{ $$ = makeList($1) }
	| alterfunc_opt_list common_func_opt_item
		{ $$ = appendList($1, $2) }
	;

/* Ignored, merely for SQL compliance */
opt_restrict:
	RESTRICT
	| /* EMPTY */
	;

/*****************************************************************************
 *
 *      DROP FUNCTION / DROP PROCEDURE / DROP ROUTINE
 *
 *****************************************************************************/

RemoveFuncStmt:
	DROP FUNCTION function_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_FUNCTION),
				Objects:    $3,
				Behavior:   int($4),
				Missing_ok: false,
			}
		}
	| DROP FUNCTION IF_P EXISTS function_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_FUNCTION),
				Objects:    $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP PROCEDURE function_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_PROCEDURE),
				Objects:    $3,
				Behavior:   int($4),
				Missing_ok: false,
			}
		}
	| DROP PROCEDURE IF_P EXISTS function_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_PROCEDURE),
				Objects:    $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	| DROP ROUTINE function_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_ROUTINE),
				Objects:    $3,
				Behavior:   int($4),
				Missing_ok: false,
			}
		}
	| DROP ROUTINE IF_P EXISTS function_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_ROUTINE),
				Objects:    $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	;

/*****************************************************************************
 *
 *      DROP AGGREGATE
 *
 *****************************************************************************/

RemoveAggrStmt:
	DROP AGGREGATE aggregate_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_AGGREGATE),
				Objects:    $3,
				Behavior:   int($4),
				Missing_ok: false,
			}
		}
	| DROP AGGREGATE IF_P EXISTS aggregate_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_AGGREGATE),
				Objects:    $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	;

/*****************************************************************************
 *
 *      DROP OPERATOR
 *
 *****************************************************************************/

RemoveOperStmt:
	DROP OPERATOR operator_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_OPERATOR),
				Objects:    $3,
				Behavior:   int($4),
				Missing_ok: false,
			}
		}
	| DROP OPERATOR IF_P EXISTS operator_with_argtypes_list opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_OPERATOR),
				Objects:    $5,
				Behavior:   int($6),
				Missing_ok: true,
			}
		}
	;

/*****************************************************************************
 *
 *      function_with_argtypes and supporting rules
 *
 *****************************************************************************/

function_with_argtypes_list:
	function_with_argtypes
		{ $$ = makeList($1) }
	| function_with_argtypes_list ',' function_with_argtypes
		{ $$ = appendList($1, $3) }
	;

function_with_argtypes:
	func_name func_args
		{
			$$ = &nodes.ObjectWithArgs{
				Objname: $1,
				Objargs: extractArgTypes($2),
			}
		}
	| type_func_name_keyword
		{
			$$ = &nodes.ObjectWithArgs{
				Objname:        &nodes.List{Items: []nodes.Node{&nodes.String{Str: $1}}},
				ArgsUnspecified: true,
			}
		}
	| ColId
		{
			$$ = &nodes.ObjectWithArgs{
				Objname:        &nodes.List{Items: []nodes.Node{&nodes.String{Str: $1}}},
				ArgsUnspecified: true,
			}
		}
	| ColId indirection
		{
			$$ = &nodes.ObjectWithArgs{
				Objname:        checkFuncName(prependList(&nodes.String{Str: $1}, $2)),
				ArgsUnspecified: true,
			}
		}
	;

func_args:
	'(' func_args_list ')'
		{ $$ = $2 }
	| '(' ')'
		{ $$ = nil }
	;

func_args_list:
	func_arg
		{ $$ = makeList($1) }
	| func_args_list ',' func_arg
		{ $$ = appendList($1, $3) }
	;

/*****************************************************************************
 *
 *      aggregate_with_argtypes and supporting rules
 *
 *****************************************************************************/

aggregate_with_argtypes:
	func_name aggr_args
		{
			$$ = &nodes.ObjectWithArgs{
				Objname: $1,
				Objargs: extractAggrArgTypes($2),
			}
		}
	;

aggregate_with_argtypes_list:
	aggregate_with_argtypes
		{ $$ = makeList($1) }
	| aggregate_with_argtypes_list ',' aggregate_with_argtypes
		{ $$ = appendList($1, $3) }
	;

/*****************************************************************************
 *
 *      operator_with_argtypes and supporting rules
 *
 *****************************************************************************/

operator_with_argtypes_list:
	operator_with_argtypes
		{ $$ = makeList($1) }
	| operator_with_argtypes_list ',' operator_with_argtypes
		{ $$ = appendList($1, $3) }
	;

operator_with_argtypes:
	any_operator oper_argtypes
		{
			$$ = &nodes.ObjectWithArgs{
				Objname: $1,
				Objargs: $2,
			}
		}
	;

oper_argtypes:
	'(' Typename ')'
		{
			pglex.Error("missing argument, use NONE to denote the missing argument of a unary operator")
			$$ = nil
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

/*****************************************************************************
 *
 *      CREATE TRIGGER
 *
 *****************************************************************************/

CreateTrigStmt:
	CREATE opt_or_replace TRIGGER name TriggerActionTime TriggerEvents ON
	qualified_name TriggerReferencing TriggerForSpec TriggerWhen
	EXECUTE FUNCTION_or_PROCEDURE func_name '(' TriggerFuncArgs ')'
		{
			eventsInt := $6.Items[0].(*nodes.Integer).Ival
			var columns *nodes.List
			if len($6.Items) > 1 && $6.Items[1] != nil {
				columns = $6.Items[1].(*nodes.List)
			}
			$$ = &nodes.CreateTrigStmt{
				Replace:        $2,
				IsConstraint:   false,
				Trigname:        $4,
				Relation:        makeRangeVarFromAnyName($8),
				Funcname:        $14,
				Args:            $16,
				Row:             $10,
				Timing:          int16($5),
				Events:          int16(eventsInt),
				Columns:         columns,
				WhenClause:      $11,
				TransitionRels:  $9,
				Deferrable:      false,
				Initdeferred:    false,
			}
		}
	| CREATE opt_or_replace CONSTRAINT TRIGGER name AFTER TriggerEvents ON
	qualified_name OptConstrFromTable ConstraintAttributeSpec
	FOR EACH ROW TriggerWhen
	EXECUTE FUNCTION_or_PROCEDURE func_name '(' TriggerFuncArgs ')'
		{
			eventsInt := $7.Items[0].(*nodes.Integer).Ival
			var columns *nodes.List
			if len($7.Items) > 1 && $7.Items[1] != nil {
				columns = $7.Items[1].(*nodes.List)
			}
			casBits := $11
			deferrable := (casBits & int64(nodes.CAS_DEFERRABLE)) != 0
			initdeferred := (casBits & int64(nodes.CAS_INITIALLY_DEFERRED)) != 0
			var constrrel *nodes.RangeVar
			if $10 != nil {
				constrrel = $10.(*nodes.RangeVar)
			}
			$$ = &nodes.CreateTrigStmt{
				Replace:        $2,
				IsConstraint:   true,
				Trigname:        $5,
				Relation:        makeRangeVarFromAnyName($9),
				Funcname:        $18,
				Args:            $20,
				Row:             true,
				Timing:          int16(nodes.TRIGGER_TYPE_AFTER),
				Events:          int16(eventsInt),
				Columns:         columns,
				WhenClause:      $15,
				Deferrable:      deferrable,
				Initdeferred:    initdeferred,
				Constrrel:       constrrel,
			}
		}
	;

TriggerActionTime:
	BEFORE    { $$ = int64(nodes.TRIGGER_TYPE_BEFORE) }
	| AFTER   { $$ = int64(nodes.TRIGGER_TYPE_AFTER) }
	| INSTEAD OF { $$ = int64(nodes.TRIGGER_TYPE_INSTEAD) }
	;

TriggerEvents:
	TriggerOneEvent
		{ $$ = $1 }
	| TriggerEvents OR TriggerOneEvent
		{
			events1 := $1.Items[0].(*nodes.Integer).Ival
			events2 := $3.Items[0].(*nodes.Integer).Ival
			var columns1, columns2 *nodes.List
			if len($1.Items) > 1 && $1.Items[1] != nil {
				columns1 = $1.Items[1].(*nodes.List)
			}
			if len($3.Items) > 1 && $3.Items[1] != nil {
				columns2 = $3.Items[1].(*nodes.List)
			}
			mergedCols := concatLists(columns1, columns2)
			var mergedColsNode nodes.Node
			if mergedCols != nil {
				mergedColsNode = mergedCols
			}
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Integer{Ival: events1 | events2},
				mergedColsNode,
			}}
		}
	;

TriggerOneEvent:
	INSERT
		{
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Integer{Ival: int64(nodes.TRIGGER_TYPE_INSERT)},
				nil,
			}}
		}
	| DELETE_P
		{
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Integer{Ival: int64(nodes.TRIGGER_TYPE_DELETE)},
				nil,
			}}
		}
	| UPDATE
		{
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Integer{Ival: int64(nodes.TRIGGER_TYPE_UPDATE)},
				nil,
			}}
		}
	| UPDATE OF columnList
		{
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Integer{Ival: int64(nodes.TRIGGER_TYPE_UPDATE)},
				$3,
			}}
		}
	| TRUNCATE
		{
			$$ = &nodes.List{Items: []nodes.Node{
				&nodes.Integer{Ival: int64(nodes.TRIGGER_TYPE_TRUNCATE)},
				nil,
			}}
		}
	;

TriggerReferencing:
	REFERENCING TriggerTransitions
		{ $$ = $2 }
	| /* EMPTY */
		{ $$ = nil }
	;

TriggerTransitions:
	TriggerTransition
		{ $$ = makeList($1) }
	| TriggerTransitions TriggerTransition
		{ $$ = appendList($1, $2) }
	;

TriggerTransition:
	TransitionOldOrNew TransitionRowOrTable opt_as TransitionRelName
		{
			$$ = &nodes.TriggerTransition{
				Name:    $4,
				IsNew:   $1,
				IsTable: $2,
			}
		}
	;

TransitionOldOrNew:
	NEW   { $$ = true }
	| OLD { $$ = false }
	;

TransitionRowOrTable:
	TABLE { $$ = true }
	| ROW { $$ = false }
	;

TransitionRelName:
	ColId { $$ = $1 }
	;

TriggerForSpec:
	FOR TriggerForOptEach TriggerForType
		{ $$ = $3 }
	| /* EMPTY */
		{ $$ = false }
	;

TriggerForOptEach:
	EACH
	| /* EMPTY */
	;

TriggerForType:
	ROW       { $$ = true }
	| STATEMENT { $$ = false }
	;

TriggerWhen:
	WHEN '(' a_expr ')'  { $$ = $3 }
	| /* EMPTY */         { $$ = nil }
	;

FUNCTION_or_PROCEDURE:
	FUNCTION
	| PROCEDURE
	;

TriggerFuncArgs:
	TriggerFuncArg
		{ $$ = makeList($1) }
	| TriggerFuncArgs ',' TriggerFuncArg
		{ $$ = appendList($1, $3) }
	| /* EMPTY */
		{ $$ = nil }
	;

TriggerFuncArg:
	Iconst
		{ $$ = &nodes.String{Str: intToString($1)} }
	| FCONST
		{ $$ = &nodes.String{Str: $1} }
	| Sconst
		{ $$ = &nodes.String{Str: $1} }
	| ColLabel
		{ $$ = &nodes.String{Str: $1} }
	;

OptConstrFromTable:
	FROM qualified_name
		{ $$ = makeRangeVar($2) }
	| /* EMPTY */
		{ $$ = nil }
	;

ConstraintAttributeSpec:
	/* EMPTY */
		{ $$ = 0 }
	| ConstraintAttributeSpec ConstraintAttributeElem
		{ $$ = $1 | $2 }
	;

ConstraintAttributeElem:
	NOT DEFERRABLE        { $$ = int64(nodes.CAS_NOT_DEFERRABLE) }
	| DEFERRABLE          { $$ = int64(nodes.CAS_DEFERRABLE) }
	| INITIALLY IMMEDIATE { $$ = int64(nodes.CAS_INITIALLY_IMMEDIATE) }
	| INITIALLY DEFERRED  { $$ = int64(nodes.CAS_INITIALLY_DEFERRED) }
	| NOT VALID           { $$ = int64(nodes.CAS_NOT_VALID) }
	| NO INHERIT          { $$ = int64(nodes.CAS_NO_INHERIT) }
	;

no_inherit:
	NO INHERIT            { $$ = true }
	| /* EMPTY */         { $$ = false }
	;

opt_unique_null_treatment:
	NULLS_LA DISTINCT      { $$ = false }
	| NULLS_LA NOT DISTINCT { $$ = true }
	| /* EMPTY */           { $$ = false }
	;

/*****************************************************************************
 *
 *      CREATE EVENT TRIGGER / ALTER EVENT TRIGGER
 *
 *****************************************************************************/

CreateEventTrigStmt:
	CREATE EVENT TRIGGER name ON ColLabel
	EXECUTE FUNCTION_or_PROCEDURE func_name '(' ')'
		{
			$$ = &nodes.CreateEventTrigStmt{
				Trigname:  $4,
				Eventname: $6,
				Funcname:  $9,
			}
		}
	| CREATE EVENT TRIGGER name ON ColLabel
	WHEN event_trigger_when_list
	EXECUTE FUNCTION_or_PROCEDURE func_name '(' ')'
		{
			$$ = &nodes.CreateEventTrigStmt{
				Trigname:   $4,
				Eventname:  $6,
				Whenclause: $8,
				Funcname:   $11,
			}
		}
	;

event_trigger_when_list:
	event_trigger_when_item
		{ $$ = makeList($1) }
	| event_trigger_when_list AND event_trigger_when_item
		{ $$ = appendList($1, $3) }
	;

event_trigger_when_item:
	ColId IN_P '(' event_trigger_value_list ')'
		{
			$$ = &nodes.DefElem{
				Defname: $1,
				Arg:     $4,
			}
		}
	;

event_trigger_value_list:
	Sconst
		{ $$ = makeList(&nodes.String{Str: $1}) }
	| event_trigger_value_list ',' Sconst
		{ $$ = appendList($1, &nodes.String{Str: $3}) }
	;

AlterEventTrigStmt:
	ALTER EVENT TRIGGER name enable_trigger
		{
			$$ = &nodes.AlterEventTrigStmt{
				Trigname:  $4,
				Tgenabled: byte($5),
			}
		}
	;

enable_trigger:
	ENABLE_P             { $$ = int64(nodes.TRIGGER_FIRES_ON_ORIGIN) }
	| ENABLE_P REPLICA   { $$ = int64(nodes.TRIGGER_FIRES_ON_REPLICA) }
	| ENABLE_P ALWAYS    { $$ = int64(nodes.TRIGGER_FIRES_ALWAYS) }
	| DISABLE_P          { $$ = int64(nodes.TRIGGER_DISABLED) }
	;

/*****************************************************************************
 *
 *      CREATE RULE
 *
 *****************************************************************************/

RuleStmt:
	CREATE opt_or_replace RULE name AS
	ON event TO qualified_name where_clause
	DO opt_instead RuleActionList
		{
			$$ = &nodes.RuleStmt{
				Replace:     $2,
				Relation:    makeRangeVarFromAnyName($9),
				Rulename:    $4,
				WhereClause: $10,
				Event:       nodes.CmdType($7),
				Instead:     $12,
				Actions:     $13,
			}
		}
	;

RuleActionList:
	NOTHING
		{ $$ = nil }
	| RuleActionStmt
		{ $$ = makeList($1) }
	| '(' RuleActionMulti ')'
		{ $$ = $2 }
	;

RuleActionMulti:
	RuleActionMulti ';' RuleActionStmtOrEmpty
		{
			if $3 != nil {
				$$ = appendList($1, $3)
			} else {
				$$ = $1
			}
		}
	| RuleActionStmtOrEmpty
		{
			if $1 != nil {
				$$ = makeList($1)
			} else {
				$$ = nil
			}
		}
	;

RuleActionStmt:
	SelectStmt  { $$ = $1 }
	| InsertStmt { $$ = $1 }
	| UpdateStmt { $$ = $1 }
	| DeleteStmt { $$ = $1 }
	| NotifyStmt { $$ = $1 }
	;

RuleActionStmtOrEmpty:
	RuleActionStmt  { $$ = $1 }
	| /* EMPTY */   { $$ = nil }
	;

event:
	SELECT    { $$ = int64(nodes.CMD_SELECT) }
	| UPDATE  { $$ = int64(nodes.CMD_UPDATE) }
	| DELETE_P { $$ = int64(nodes.CMD_DELETE) }
	| INSERT  { $$ = int64(nodes.CMD_INSERT) }
	;

opt_instead:
	INSTEAD   { $$ = true }
	| ALSO    { $$ = false }
	| /* EMPTY */ { $$ = false }
	;

/*****************************************************************************
 *
 *      CREATE LANGUAGE
 *
 *****************************************************************************/

CreatePLangStmt:
	CREATE opt_or_replace opt_trusted opt_procedural LANGUAGE name
		{
			/* Parameterless form - creates extension */
			$$ = &nodes.CreatePLangStmt{
				Replace:   $2,
				Plname:    $6,
				Pltrusted: $3,
			}
		}
	| CREATE opt_or_replace opt_trusted opt_procedural LANGUAGE name
	  HANDLER handler_name opt_inline_handler opt_validator
		{
			$$ = &nodes.CreatePLangStmt{
				Replace:     $2,
				Plname:      $6,
				Plhandler:   $8,
				Plinline:    $9,
				Plvalidator: $10,
				Pltrusted:   $3,
			}
		}
	;

opt_trusted:
	TRUSTED    { $$ = true }
	| /* EMPTY */ { $$ = false }
	;

handler_name:
	name
		{ $$ = &nodes.List{Items: []nodes.Node{&nodes.String{Str: $1}}} }
	| name attrs
		{ $$ = prependList(&nodes.String{Str: $1}, $2) }
	;

opt_inline_handler:
	INLINE_P handler_name { $$ = $2 }
	| /* EMPTY */         { $$ = nil }
	;

opt_validator:
	VALIDATOR handler_name  { $$ = $2 }
	| NO VALIDATOR          { $$ = nil }
	| /* EMPTY */           { $$ = nil }
	;

opt_procedural:
	PROCEDURAL
	| /* EMPTY */
	;


/*****************************************************************************
 *
 *      FOREIGN DATA WRAPPER statements
 *
 *****************************************************************************/

CreateFdwStmt:
	CREATE FOREIGN DATA_P WRAPPER name opt_fdw_options create_generic_options
		{
			$$ = &nodes.CreateFdwStmt{
				Fdwname:     $5,
				FuncOptions: $6,
				Options:     $7,
			}
		}
	;

fdw_option:
	HANDLER handler_name
		{
			$$ = makeDefElem("handler", $2)
		}
	| NO HANDLER
		{
			$$ = makeDefElem("handler", nil)
		}
	| VALIDATOR handler_name
		{
			$$ = makeDefElem("validator", $2)
		}
	| NO VALIDATOR
		{
			$$ = makeDefElem("validator", nil)
		}
	;

fdw_options:
	fdw_option
		{ $$ = makeList($1) }
	| fdw_options fdw_option
		{ $$ = appendList($1, $2) }
	;

opt_fdw_options:
	fdw_options   { $$ = $1 }
	| /* EMPTY */ { $$ = nil }
	;

/*****************************************************************************
 *
 *      ALTER FOREIGN DATA WRAPPER
 *
 *****************************************************************************/

AlterFdwStmt:
	ALTER FOREIGN DATA_P WRAPPER name opt_fdw_options alter_generic_options
		{
			$$ = &nodes.AlterFdwStmt{
				Fdwname:     $5,
				FuncOptions: $6,
				Options:     $7,
			}
		}
	| ALTER FOREIGN DATA_P WRAPPER name fdw_options
		{
			$$ = &nodes.AlterFdwStmt{
				Fdwname:     $5,
				FuncOptions: $6,
			}
		}
	;

/*****************************************************************************
 *
 *      Generic OPTIONS infrastructure for FDW/SERVER/USER MAPPING
 *
 *****************************************************************************/

create_generic_options:
	OPTIONS '(' generic_option_list ')'
		{ $$ = $3 }
	| /* EMPTY */
		{ $$ = nil }
	;

generic_option_list:
	generic_option_elem
		{ $$ = makeList($1) }
	| generic_option_list ',' generic_option_elem
		{ $$ = appendList($1, $3) }
	;

alter_generic_options:
	OPTIONS '(' alter_generic_option_list ')'
		{ $$ = $3 }
	;

alter_generic_option_list:
	alter_generic_option_elem
		{ $$ = makeList($1) }
	| alter_generic_option_list ',' alter_generic_option_elem
		{ $$ = appendList($1, $3) }
	;

alter_generic_option_elem:
	generic_option_elem
		{
			$$ = $1
		}
	| SET generic_option_elem
		{
			n := $2.(*nodes.DefElem)
			n.Defaction = int(nodes.DEFELEM_SET)
			$$ = n
		}
	| ADD_P generic_option_elem
		{
			n := $2.(*nodes.DefElem)
			n.Defaction = int(nodes.DEFELEM_ADD)
			$$ = n
		}
	| DROP generic_option_name
		{
			$$ = &nodes.DefElem{
				Defname:   $2,
				Defaction: int(nodes.DEFELEM_DROP),
				Location:  -1,
			}
		}
	;

generic_option_elem:
	generic_option_name generic_option_arg
		{
			$$ = &nodes.DefElem{
				Defname:  $1,
				Arg:      $2,
				Location: -1,
			}
		}
	;

generic_option_name:
	ColLabel  { $$ = $1 }
	;

generic_option_arg:
	Sconst  { $$ = &nodes.String{Str: $1} }
	;

/*****************************************************************************
 *
 *      CREATE SERVER / ALTER SERVER
 *
 *****************************************************************************/

CreateForeignServerStmt:
	CREATE SERVER name opt_type opt_foreign_server_version
	  FOREIGN DATA_P WRAPPER name create_generic_options
		{
			$$ = &nodes.CreateForeignServerStmt{
				Servername:  $3,
				Servertype:  $4,
				Version:     $5,
				Fdwname:     $9,
				Options:     $10,
				IfNotExists: false,
			}
		}
	| CREATE SERVER IF_P NOT EXISTS name opt_type opt_foreign_server_version
	  FOREIGN DATA_P WRAPPER name create_generic_options
		{
			$$ = &nodes.CreateForeignServerStmt{
				Servername:  $6,
				Servertype:  $7,
				Version:     $8,
				Fdwname:     $12,
				Options:     $13,
				IfNotExists: true,
			}
		}
	;

opt_type:
	TYPE_P Sconst   { $$ = $2 }
	| /* EMPTY */   { $$ = "" }
	;

foreign_server_version:
	VERSION_P Sconst   { $$ = $2 }
	| VERSION_P NULL_P { $$ = "" }
	;

opt_foreign_server_version:
	foreign_server_version { $$ = $1 }
	| /* EMPTY */          { $$ = "" }
	;

/*****************************************************************************
 *
 *      ALTER SERVER
 *
 *****************************************************************************/

AlterForeignServerStmt:
	ALTER SERVER name foreign_server_version alter_generic_options
		{
			$$ = &nodes.AlterForeignServerStmt{
				Servername: $3,
				Version:    $4,
				Options:    $5,
				HasVersion: true,
			}
		}
	| ALTER SERVER name foreign_server_version
		{
			$$ = &nodes.AlterForeignServerStmt{
				Servername: $3,
				Version:    $4,
				HasVersion: true,
			}
		}
	| ALTER SERVER name alter_generic_options
		{
			$$ = &nodes.AlterForeignServerStmt{
				Servername: $3,
				Options:    $4,
			}
		}
	;

/*****************************************************************************
 *
 *      CREATE FOREIGN TABLE
 *
 *****************************************************************************/

CreateForeignTableStmt:
	CREATE FOREIGN TABLE qualified_name
	  '(' OptTableElementList ')'
	  OptInherit SERVER name create_generic_options
		{
			rv := makeRangeVar($4)
			rv.(*nodes.RangeVar).Relpersistence = 'p'
			$$ = &nodes.CreateForeignTableStmt{
				Base: nodes.CreateStmt{
					Relation:     rv.(*nodes.RangeVar),
					TableElts:    $6,
					InhRelations: $8,
					IfNotExists:  false,
				},
				Servername: $10,
				Options:    $11,
			}
		}
	| CREATE FOREIGN TABLE IF_P NOT EXISTS qualified_name
	  '(' OptTableElementList ')'
	  OptInherit SERVER name create_generic_options
		{
			rv := makeRangeVar($7)
			rv.(*nodes.RangeVar).Relpersistence = 'p'
			$$ = &nodes.CreateForeignTableStmt{
				Base: nodes.CreateStmt{
					Relation:     rv.(*nodes.RangeVar),
					TableElts:    $9,
					InhRelations: $11,
					IfNotExists:  true,
				},
				Servername: $13,
				Options:    $14,
			}
		}
	| CREATE FOREIGN TABLE qualified_name
	  PARTITION OF qualified_name OptTypedTableElementList ForValues
	  SERVER name create_generic_options
		{
			rv := makeRangeVar($4)
			inh := makeRangeVar($7)
			rv.(*nodes.RangeVar).Relpersistence = 'p'
			$$ = &nodes.CreateForeignTableStmt{
				Base: nodes.CreateStmt{
					Relation:     rv.(*nodes.RangeVar),
					InhRelations: makeList(inh),
					TableElts:    $8,
					Partbound:    $9,
					IfNotExists:  false,
				},
				Servername: $11,
				Options:    $12,
			}
		}
	| CREATE FOREIGN TABLE IF_P NOT EXISTS qualified_name
	  PARTITION OF qualified_name OptTypedTableElementList ForValues
	  SERVER name create_generic_options
		{
			rv := makeRangeVar($7)
			inh := makeRangeVar($10)
			rv.(*nodes.RangeVar).Relpersistence = 'p'
			$$ = &nodes.CreateForeignTableStmt{
				Base: nodes.CreateStmt{
					Relation:     rv.(*nodes.RangeVar),
					InhRelations: makeList(inh),
					TableElts:    $11,
					Partbound:    $12,
					IfNotExists:  true,
				},
				Servername: $14,
				Options:    $15,
			}
		}
	;

/*****************************************************************************
 *
 *      IMPORT FOREIGN SCHEMA
 *
 *****************************************************************************/

ImportForeignSchemaStmt:
	IMPORT_P FOREIGN SCHEMA name import_qualification
	  FROM SERVER name INTO name create_generic_options
		{
			var listType nodes.ImportForeignSchemaType
			var tableList *nodes.List
			if $5 != nil {
				qual := $5.(*importQualification)
				listType = qual.listType
				tableList = qual.tableList
			}
			$$ = &nodes.ImportForeignSchemaStmt{
				ServerName:   $8,
				RemoteSchema: $4,
				LocalSchema:  $10,
				ListType:     listType,
				TableList:    tableList,
				Options:      $11,
			}
		}
	;

import_qualification_type:
	LIMIT TO    { $$ = int64(nodes.FDW_IMPORT_SCHEMA_LIMIT_TO) }
	| EXCEPT    { $$ = int64(nodes.FDW_IMPORT_SCHEMA_EXCEPT) }
	;

import_qualification:
	import_qualification_type '(' relation_expr_list ')'
		{
			$$ = &importQualification{
				listType:  nodes.ImportForeignSchemaType($1),
				tableList: $3,
			}
		}
	| /* EMPTY */
		{
			$$ = nil
		}
	;

/*****************************************************************************
 *
 *      CREATE / ALTER / DROP USER MAPPING
 *
 *****************************************************************************/

CreateUserMappingStmt:
	CREATE USER MAPPING FOR auth_ident SERVER name create_generic_options
		{
			$$ = &nodes.CreateUserMappingStmt{
				User:        $5.(*nodes.RoleSpec),
				Servername:  $7,
				Options:     $8,
				IfNotExists: false,
			}
		}
	| CREATE USER MAPPING IF_P NOT EXISTS FOR auth_ident SERVER name create_generic_options
		{
			$$ = &nodes.CreateUserMappingStmt{
				User:        $8.(*nodes.RoleSpec),
				Servername:  $10,
				Options:     $11,
				IfNotExists: true,
			}
		}
	;

auth_ident:
	RoleSpec   { $$ = $1 }
	| USER
		{
			$$ = &nodes.RoleSpec{
				Roletype: int(nodes.ROLESPEC_CURRENT_USER),
			}
		}
	;

DropUserMappingStmt:
	DROP USER MAPPING FOR auth_ident SERVER name
		{
			$$ = &nodes.DropUserMappingStmt{
				User:       $5.(*nodes.RoleSpec),
				Servername: $7,
				MissingOk:  false,
			}
		}
	| DROP USER MAPPING IF_P EXISTS FOR auth_ident SERVER name
		{
			$$ = &nodes.DropUserMappingStmt{
				User:       $7.(*nodes.RoleSpec),
				Servername: $9,
				MissingOk:  true,
			}
		}
	;

AlterUserMappingStmt:
	ALTER USER MAPPING FOR auth_ident SERVER name alter_generic_options
		{
			$$ = &nodes.AlterUserMappingStmt{
				User:       $5.(*nodes.RoleSpec),
				Servername: $7,
				Options:    $8,
			}
		}
	;


/*****************************************************************************
 *
 *      CREATE TABLESPACE / DROP TABLESPACE / ALTER TABLESPACE
 *
 *****************************************************************************/

CreateTableSpaceStmt:
	CREATE TABLESPACE name OptTableSpaceOwner LOCATION Sconst opt_reloptions
		{
			var owner *nodes.RoleSpec
			if $4 != nil {
				owner = $4.(*nodes.RoleSpec)
			}
			$$ = &nodes.CreateTableSpaceStmt{
				Tablespacename: $3,
				Owner:          owner,
				Location:       $6,
				Options:        $7,
			}
		}
	;

OptTableSpaceOwner:
	OWNER RoleSpec { $$ = $2 }
	| /* EMPTY */  { $$ = nil }
	;

DropTableSpaceStmt:
	DROP TABLESPACE name
		{
			$$ = &nodes.DropTableSpaceStmt{
				Tablespacename: $3,
				MissingOk:      false,
			}
		}
	| DROP TABLESPACE IF_P EXISTS name
		{
			$$ = &nodes.DropTableSpaceStmt{
				Tablespacename: $5,
				MissingOk:      true,
			}
		}
	;

AlterTblSpcStmt:
	ALTER TABLESPACE name SET reloptions
		{
			$$ = &nodes.AlterTableSpaceOptionsStmt{
				Tablespacename: $3,
				Options:        $5,
				IsReset:        false,
			}
		}
	| ALTER TABLESPACE name RESET reloptions
		{
			$$ = &nodes.AlterTableSpaceOptionsStmt{
				Tablespacename: $3,
				Options:        $5,
				IsReset:        true,
			}
		}
	;

reloptions:
	'(' reloption_list ')'    { $$ = $2 }
	;

opt_reloptions:
	WITH reloptions   { $$ = $2 }
	| /* EMPTY */     { $$ = nil }
	;

reloption_list:
	reloption_elem
		{ $$ = makeList($1) }
	| reloption_list ',' reloption_elem
		{ $$ = appendList($1, $3) }
	;

reloption_elem:
	ColLabel '=' def_arg
		{
			$$ = makeDefElem($1, $3)
		}
	| ColLabel
		{
			$$ = makeDefElem($1, nil)
		}
	| ColLabel '.' ColLabel '=' def_arg
		{
			$$ = &nodes.DefElem{
				Defnamespace: $1,
				Defname:      $3,
				Arg:          $5,
				Location:     -1,
			}
		}
	| ColLabel '.' ColLabel
		{
			$$ = &nodes.DefElem{
				Defnamespace: $1,
				Defname:      $3,
				Location:     -1,
			}
		}
	;

/*****************************************************************************
 *
 *      CREATE EXTENSION / ALTER EXTENSION
 *
 *****************************************************************************/

CreateExtensionStmt:
	CREATE EXTENSION name opt_with create_extension_opt_list
		{
			$$ = &nodes.CreateExtensionStmt{
				Extname:     $3,
				IfNotExists: false,
				Options:     $5,
			}
		}
	| CREATE EXTENSION IF_P NOT EXISTS name opt_with create_extension_opt_list
		{
			$$ = &nodes.CreateExtensionStmt{
				Extname:     $6,
				IfNotExists: true,
				Options:     $8,
			}
		}
	;

create_extension_opt_list:
	create_extension_opt_list create_extension_opt_item
		{ $$ = appendList($1, $2) }
	| /* EMPTY */
		{ $$ = nil }
	;

create_extension_opt_item:
	SCHEMA name
		{
			$$ = makeDefElem("schema", &nodes.String{Str: $2})
		}
	| VERSION_P NonReservedWord_or_Sconst
		{
			$$ = makeDefElem("new_version", &nodes.String{Str: $2})
		}
	| CASCADE
		{
			$$ = makeDefElem("cascade", &nodes.Boolean{Boolval: true})
		}
	;

/*****************************************************************************
 *
 *      ALTER EXTENSION name UPDATE [ TO version ]
 *
 *****************************************************************************/

AlterExtensionStmt:
	ALTER EXTENSION name UPDATE alter_extension_opt_list
		{
			$$ = &nodes.AlterExtensionStmt{
				Extname: $3,
				Options: $5,
			}
		}
	;

alter_extension_opt_list:
	alter_extension_opt_list alter_extension_opt_item
		{ $$ = appendList($1, $2) }
	| /* EMPTY */
		{ $$ = nil }
	;

alter_extension_opt_item:
	TO NonReservedWord_or_Sconst
		{
			$$ = makeDefElem("new_version", &nodes.String{Str: $2})
		}
	;

/*****************************************************************************
 *
 *      ALTER EXTENSION name ADD/DROP object-identifier
 *
 *****************************************************************************/

AlterExtensionContentsStmt:
	ALTER EXTENSION name add_drop object_type_name name
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.ObjectType($5),
				Object:  &nodes.String{Str: $6},
			}
		}
	| ALTER EXTENSION name add_drop object_type_any_name any_name
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.ObjectType($5),
				Object:  $6,
			}
		}
	| ALTER EXTENSION name add_drop AGGREGATE aggregate_with_argtypes
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_AGGREGATE,
				Object:  $6,
			}
		}
	| ALTER EXTENSION name add_drop FUNCTION function_with_argtypes
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_FUNCTION,
				Object:  $6,
			}
		}
	| ALTER EXTENSION name add_drop PROCEDURE function_with_argtypes
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_PROCEDURE,
				Object:  $6,
			}
		}
	| ALTER EXTENSION name add_drop ROUTINE function_with_argtypes
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_ROUTINE,
				Object:  $6,
			}
		}
	| ALTER EXTENSION name add_drop OPERATOR operator_with_argtypes
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_OPERATOR,
				Object:  $6,
			}
		}
	| ALTER EXTENSION name add_drop OPERATOR CLASS any_name USING name
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_OPCLASS,
				Object:  prependList(&nodes.String{Str: $9}, $7),
			}
		}
	| ALTER EXTENSION name add_drop OPERATOR FAMILY any_name USING name
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_OPFAMILY,
				Object:  prependList(&nodes.String{Str: $9}, $7),
			}
		}
	| ALTER EXTENSION name add_drop DOMAIN_P Typename
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_DOMAIN,
				Object:  $6,
			}
		}
	| ALTER EXTENSION name add_drop TYPE_P Typename
		{
			$$ = &nodes.AlterExtensionContentsStmt{
				Extname: $3,
				Action:  int($4),
				Objtype: nodes.OBJECT_TYPE,
				Object:  $6,
			}
		}
	;

/*****************************************************************************
 *
 *      CREATE ACCESS METHOD
 *
 *****************************************************************************/

CreateAmStmt:
	CREATE ACCESS METHOD name TYPE_P am_type HANDLER handler_name
		{
			$$ = &nodes.CreateAmStmt{
				Amname:      $4,
				Amtype:      byte($6),
				HandlerName: $8,
			}
		}
	;

am_type:
	INDEX     { $$ = nodes.AMTYPE_INDEX }
	| TABLE   { $$ = nodes.AMTYPE_TABLE }
	;

/*****************************************************************************
 *
 *      CREATE POLICY / ALTER POLICY
 *
 *****************************************************************************/

CreatePolicyStmt:
	CREATE POLICY name ON qualified_name RowSecurityDefaultPermissive
		RowSecurityDefaultForCmd RowSecurityDefaultToRole
		RowSecurityOptionalExpr RowSecurityOptionalWithCheck
		{
			$$ = &nodes.CreatePolicyStmt{
				PolicyName: $3,
				Table:      makeRangeVarFromAnyName($5),
				Permissive: $6,
				CmdName:    $7,
				Roles:      $8,
				Qual:       $9,
				WithCheck:  $10,
			}
		}
	;

AlterPolicyStmt:
	ALTER POLICY name ON qualified_name RowSecurityOptionalToRole
		RowSecurityOptionalExpr RowSecurityOptionalWithCheck
		{
			$$ = &nodes.AlterPolicyStmt{
				PolicyName: $3,
				Table:      makeRangeVarFromAnyName($5),
				Roles:      $6,
				Qual:       $7,
				WithCheck:  $8,
			}
		}
	;

RowSecurityOptionalExpr:
	USING '(' a_expr ')'    { $$ = $3 }
	| /* EMPTY */            { $$ = nil }
	;

RowSecurityOptionalWithCheck:
	WITH CHECK '(' a_expr ')'    { $$ = $4 }
	| /* EMPTY */                 { $$ = nil }
	;

RowSecurityDefaultToRole:
	TO role_list    { $$ = $2 }
	| /* EMPTY */
		{
			/* Default is PUBLIC */
			$$ = makeList(&nodes.RoleSpec{
				Roletype: int(nodes.ROLESPEC_PUBLIC),
				Location: -1,
			})
		}
	;

RowSecurityOptionalToRole:
	TO role_list    { $$ = $2 }
	| /* EMPTY */   { $$ = nil }
	;

RowSecurityDefaultPermissive:
	AS IDENT
		{
			if $2 == "permissive" {
				$$ = true
			} else if $2 == "restrictive" {
				$$ = false
			} else {
				pglex.Error("only PERMISSIVE or RESTRICTIVE policies are supported")
				$$ = true
			}
		}
	| /* EMPTY */    { $$ = true }
	;

RowSecurityDefaultForCmd:
	FOR row_security_cmd    { $$ = $2 }
	| /* EMPTY */           { $$ = "all" }
	;

row_security_cmd:
	ALL          { $$ = "all" }
	| SELECT     { $$ = "select" }
	| INSERT     { $$ = "insert" }
	| UPDATE     { $$ = "update" }
	| DELETE_P   { $$ = "delete" }
	;

/*****************************************************************************
 *
 *      CREATE PUBLICATION / ALTER PUBLICATION
 *
 *****************************************************************************/

CreatePublicationStmt:
	CREATE PUBLICATION name opt_definition
		{
			$$ = &nodes.CreatePublicationStmt{
				Pubname: $3,
				Options: $4,
			}
		}
	| CREATE PUBLICATION name FOR ALL TABLES opt_definition
		{
			$$ = &nodes.CreatePublicationStmt{
				Pubname:      $3,
				Options:      $7,
				ForAllTables: true,
			}
		}
	| CREATE PUBLICATION name FOR pub_obj_list opt_definition
		{
			$$ = &nodes.CreatePublicationStmt{
				Pubname:    $3,
				Options:    $6,
				Pubobjects: $5,
			}
		}
	;

AlterPublicationStmt:
	ALTER PUBLICATION name SET definition
		{
			$$ = &nodes.AlterPublicationStmt{
				Pubname: $3,
				Options: $5,
			}
		}
	| ALTER PUBLICATION name ADD_P pub_obj_list
		{
			$$ = &nodes.AlterPublicationStmt{
				Pubname:    $3,
				Pubobjects: $5,
				Action:     nodes.DEFELEM_ADD,
			}
		}
	| ALTER PUBLICATION name SET pub_obj_list
		{
			$$ = &nodes.AlterPublicationStmt{
				Pubname:    $3,
				Pubobjects: $5,
				Action:     nodes.DEFELEM_SET,
			}
		}
	| ALTER PUBLICATION name DROP pub_obj_list
		{
			$$ = &nodes.AlterPublicationStmt{
				Pubname:    $3,
				Pubobjects: $5,
				Action:     nodes.DEFELEM_DROP,
			}
		}
	;

opt_definition:
	WITH definition    { $$ = $2 }
	| /* EMPTY */      { $$ = nil }
	;

PublicationObjSpec:
	TABLE relation_expr opt_column_list OptWhereClause
		{
			var cols *nodes.List
			if $3 != nil {
				cols = $3
			}
			pt := &nodes.PublicationTable{
				Relation:    $2.(*nodes.RangeVar),
				Columns:     cols,
			}
			if $4 != nil {
				pt.WhereClause = $4.Items[0]
			}
			$$ = &nodes.PublicationObjSpec{
				Pubobjtype: nodes.PUBLICATIONOBJ_TABLE,
				Pubtable:   pt,
			}
		}
	| TABLES IN_P SCHEMA ColId
		{
			$$ = &nodes.PublicationObjSpec{
				Pubobjtype: nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA,
				Name:       $4,
			}
		}
	| TABLES IN_P SCHEMA CURRENT_SCHEMA
		{
			$$ = &nodes.PublicationObjSpec{
				Pubobjtype: nodes.PUBLICATIONOBJ_TABLES_IN_CUR_SCHEMA,
			}
		}
	| relation_expr opt_column_list OptWhereClause
		{
			pt := &nodes.PublicationTable{
				Relation: $1.(*nodes.RangeVar),
				Columns:  $2,
			}
			if $3 != nil {
				pt.WhereClause = $3.Items[0]
			}
			$$ = &nodes.PublicationObjSpec{
				Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
				Pubtable:   pt,
			}
		}
	| CURRENT_SCHEMA
		{
			$$ = &nodes.PublicationObjSpec{
				Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
			}
		}
	;

OptWhereClause:
	WHERE '(' a_expr ')'   { $$ = makeList($3) }
	| /* EMPTY */           { $$ = nil }
	;

pub_obj_list:
	PublicationObjSpec
		{ $$ = makeList($1) }
	| pub_obj_list ',' PublicationObjSpec
		{ $$ = appendList($1, $3) }
	;

/*****************************************************************************
 *
 *      CREATE SUBSCRIPTION / ALTER SUBSCRIPTION / DROP SUBSCRIPTION
 *
 *****************************************************************************/

CreateSubscriptionStmt:
	CREATE SUBSCRIPTION name CONNECTION Sconst PUBLICATION name_list opt_definition
		{
			$$ = &nodes.CreateSubscriptionStmt{
				Subname:     $3,
				Conninfo:    $5,
				Publication: $7,
				Options:     $8,
			}
		}
	;

AlterSubscriptionStmt:
	ALTER SUBSCRIPTION name SET definition
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:    nodes.ALTER_SUBSCRIPTION_OPTIONS,
				Subname: $3,
				Options: $5,
			}
		}
	| ALTER SUBSCRIPTION name CONNECTION Sconst
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:     nodes.ALTER_SUBSCRIPTION_CONNECTION,
				Subname:  $3,
				Conninfo: $5,
			}
		}
	| ALTER SUBSCRIPTION name REFRESH PUBLICATION opt_definition
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:    nodes.ALTER_SUBSCRIPTION_REFRESH,
				Subname: $3,
				Options: $6,
			}
		}
	| ALTER SUBSCRIPTION name ADD_P PUBLICATION name_list opt_definition
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:        nodes.ALTER_SUBSCRIPTION_ADD_PUBLICATION,
				Subname:     $3,
				Publication: $6,
				Options:     $7,
			}
		}
	| ALTER SUBSCRIPTION name DROP PUBLICATION name_list opt_definition
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:        nodes.ALTER_SUBSCRIPTION_DROP_PUBLICATION,
				Subname:     $3,
				Publication: $6,
				Options:     $7,
			}
		}
	| ALTER SUBSCRIPTION name SET PUBLICATION name_list opt_definition
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:        nodes.ALTER_SUBSCRIPTION_SET_PUBLICATION,
				Subname:     $3,
				Publication: $6,
				Options:     $7,
			}
		}
	| ALTER SUBSCRIPTION name ENABLE_P
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:    nodes.ALTER_SUBSCRIPTION_ENABLED,
				Subname: $3,
				Options: makeList(makeDefElem("enabled", &nodes.Boolean{Boolval: true})),
			}
		}
	| ALTER SUBSCRIPTION name DISABLE_P
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:    nodes.ALTER_SUBSCRIPTION_ENABLED,
				Subname: $3,
				Options: makeList(makeDefElem("enabled", &nodes.Boolean{Boolval: false})),
			}
		}
	| ALTER SUBSCRIPTION name SKIP definition
		{
			$$ = &nodes.AlterSubscriptionStmt{
				Kind:    nodes.ALTER_SUBSCRIPTION_SKIP,
				Subname: $3,
				Options: $5,
			}
		}
	;

DropSubscriptionStmt:
	DROP SUBSCRIPTION name opt_drop_behavior
		{
			$$ = &nodes.DropSubscriptionStmt{
				Subname:   $3,
				MissingOk: false,
				Behavior:  nodes.DropBehavior($4),
			}
		}
	| DROP SUBSCRIPTION IF_P EXISTS name opt_drop_behavior
		{
			$$ = &nodes.DropSubscriptionStmt{
				Subname:   $5,
				MissingOk: true,
				Behavior:  nodes.DropBehavior($6),
			}
		}
	;


/*****************************************************************************
 *
 * ALTER ... DEPENDS ON EXTENSION
 *
 *****************************************************************************/

AlterObjectDependsStmt:
	ALTER FUNCTION function_with_argtypes opt_no DEPENDS ON EXTENSION name
		{
			$$ = &nodes.AlterObjectDependsStmt{
				ObjectType: nodes.OBJECT_FUNCTION,
				Object:     $3,
				Extname:    &nodes.String{Str: $8},
				Remove:     $4,
			}
		}
	| ALTER PROCEDURE function_with_argtypes opt_no DEPENDS ON EXTENSION name
		{
			$$ = &nodes.AlterObjectDependsStmt{
				ObjectType: nodes.OBJECT_PROCEDURE,
				Object:     $3,
				Extname:    &nodes.String{Str: $8},
				Remove:     $4,
			}
		}
	| ALTER ROUTINE function_with_argtypes opt_no DEPENDS ON EXTENSION name
		{
			$$ = &nodes.AlterObjectDependsStmt{
				ObjectType: nodes.OBJECT_ROUTINE,
				Object:     $3,
				Extname:    &nodes.String{Str: $8},
				Remove:     $4,
			}
		}
	| ALTER TRIGGER name ON qualified_name opt_no DEPENDS ON EXTENSION name
		{
			$$ = &nodes.AlterObjectDependsStmt{
				ObjectType: nodes.OBJECT_TRIGGER,
				Relation:   makeRangeVarFromAnyName($5),
				Object:     makeList(&nodes.String{Str: $3}),
				Extname:    &nodes.String{Str: $10},
				Remove:     $6,
			}
		}
	| ALTER MATERIALIZED VIEW qualified_name opt_no DEPENDS ON EXTENSION name
		{
			$$ = &nodes.AlterObjectDependsStmt{
				ObjectType: nodes.OBJECT_MATVIEW,
				Relation:   makeRangeVarFromAnyName($4),
				Extname:    &nodes.String{Str: $9},
				Remove:     $5,
			}
		}
	| ALTER INDEX qualified_name opt_no DEPENDS ON EXTENSION name
		{
			$$ = &nodes.AlterObjectDependsStmt{
				ObjectType: nodes.OBJECT_INDEX,
				Relation:   makeRangeVarFromAnyName($3),
				Extname:    &nodes.String{Str: $8},
				Remove:     $4,
			}
		}
	;

opt_no:
	NO
		{ $$ = true }
	| /* EMPTY */
		{ $$ = false }
	;

/*****************************************************************************
 *
 * ALTER ... SET SCHEMA
 *
 *****************************************************************************/

AlterObjectSchemaStmt:
	ALTER AGGREGATE aggregate_with_argtypes SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_AGGREGATE,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER COLLATION any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_COLLATION,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER CONVERSION_P any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_CONVERSION,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER DOMAIN_P any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_DOMAIN,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER EXTENSION name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_EXTENSION,
				Object:     &nodes.String{Str: $3},
				Newschema:  $6,
			}
		}
	| ALTER FUNCTION function_with_argtypes SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_FUNCTION,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER OPERATOR operator_with_argtypes SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_OPERATOR,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER OPERATOR CLASS any_name USING name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_OPCLASS,
				Object:     prependList(&nodes.String{Str: $6}, $4),
				Newschema:  $9,
			}
		}
	| ALTER OPERATOR FAMILY any_name USING name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_OPFAMILY,
				Object:     prependList(&nodes.String{Str: $6}, $4),
				Newschema:  $9,
			}
		}
	| ALTER PROCEDURE function_with_argtypes SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_PROCEDURE,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER ROUTINE function_with_argtypes SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_ROUTINE,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER TABLE relation_expr SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TABLE,
				Relation:   $3.(*nodes.RangeVar),
				Newschema:  $6,
			}
		}
	| ALTER TABLE IF_P EXISTS relation_expr SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TABLE,
				Relation:   $5.(*nodes.RangeVar),
				Newschema:  $8,
				MissingOk:  true,
			}
		}
	| ALTER STATISTICS any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_STATISTIC_EXT,
				Object:     $3,
				Newschema:  $6,
			}
		}
	| ALTER TEXT_P SEARCH PARSER any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TSPARSER,
				Object:     $5,
				Newschema:  $8,
			}
		}
	| ALTER TEXT_P SEARCH DICTIONARY any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TSDICTIONARY,
				Object:     $5,
				Newschema:  $8,
			}
		}
	| ALTER TEXT_P SEARCH TEMPLATE any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TSTEMPLATE,
				Object:     $5,
				Newschema:  $8,
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TSCONFIGURATION,
				Object:     $5,
				Newschema:  $8,
			}
		}
	| ALTER SEQUENCE qualified_name SET SCHEMA name
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_SEQUENCE,
				Relation:   rv,
				Newschema:  $6,
			}
		}
	| ALTER SEQUENCE IF_P EXISTS qualified_name SET SCHEMA name
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_SEQUENCE,
				Relation:   rv,
				Newschema:  $8,
				MissingOk:  true,
			}
		}
	| ALTER VIEW qualified_name SET SCHEMA name
		{
			rv := makeRangeVarFromAnyName($3)
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_VIEW,
				Relation:   rv,
				Newschema:  $6,
			}
		}
	| ALTER VIEW IF_P EXISTS qualified_name SET SCHEMA name
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_VIEW,
				Relation:   rv,
				Newschema:  $8,
				MissingOk:  true,
			}
		}
	| ALTER MATERIALIZED VIEW qualified_name SET SCHEMA name
		{
			rv := makeRangeVarFromAnyName($4)
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_MATVIEW,
				Relation:   rv,
				Newschema:  $7,
			}
		}
	| ALTER MATERIALIZED VIEW IF_P EXISTS qualified_name SET SCHEMA name
		{
			rv := makeRangeVarFromAnyName($6)
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_MATVIEW,
				Relation:   rv,
				Newschema:  $9,
				MissingOk:  true,
			}
		}
	| ALTER FOREIGN TABLE relation_expr SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:   $4.(*nodes.RangeVar),
				Newschema:  $7,
			}
		}
	| ALTER FOREIGN TABLE IF_P EXISTS relation_expr SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_FOREIGN_TABLE,
				Relation:   $6.(*nodes.RangeVar),
				Newschema:  $9,
				MissingOk:  true,
			}
		}
	| ALTER TYPE_P any_name SET SCHEMA name
		{
			$$ = &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_TYPE,
				Object:     $3,
				Newschema:  $6,
			}
		}
	;

/*****************************************************************************
 *
 * ALTER OPERATOR ... SET (...)
 *
 *****************************************************************************/

AlterOperatorStmt:
	ALTER OPERATOR operator_with_argtypes SET '(' operator_def_list ')'
		{
			$$ = &nodes.AlterOperatorStmt{
				Opername: $3.(*nodes.ObjectWithArgs),
				Options:  $6,
			}
		}
	;

operator_def_list:
	operator_def_elem
		{ $$ = makeList($1) }
	| operator_def_list ',' operator_def_elem
		{ $$ = appendList($1, $3) }
	;

operator_def_elem:
	ColLabel '=' NONE
		{
			$$ = &nodes.DefElem{
				Defname:  $1,
				Location: -1,
			}
		}
	| ColLabel '=' operator_def_arg
		{
			$$ = &nodes.DefElem{
				Defname:  $1,
				Arg:      $3,
				Location: -1,
			}
		}
	| ColLabel
		{
			$$ = &nodes.DefElem{
				Defname:  $1,
				Location: -1,
			}
		}
	;

operator_def_arg:
	func_type
		{ $$ = $1 }
	| reserved_keyword
		{ $$ = &nodes.String{Str: $1} }
	| qual_all_Op
		{ $$ = $1 }
	| NumericOnly
		{ $$ = $1 }
	| Sconst
		{ $$ = &nodes.String{Str: $1} }
	;

/*****************************************************************************
 *
 * ALTER TYPE ... SET (...)
 *
 *****************************************************************************/

AlterTypeStmt:
	ALTER TYPE_P any_name SET '(' operator_def_list ')'
		{
			$$ = &nodes.AlterTypeStmt{
				TypeName: $3,
				Options:  $6,
			}
		}
	;

/*****************************************************************************
 *
 * ALTER ... OWNER TO
 *
 *****************************************************************************/

AlterOwnerStmt:
	ALTER AGGREGATE aggregate_with_argtypes OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_AGGREGATE,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER COLLATION any_name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_COLLATION,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER CONVERSION_P any_name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_CONVERSION,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER DATABASE name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_DATABASE,
				Object:     &nodes.String{Str: $3},
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER DOMAIN_P any_name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_DOMAIN,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER FUNCTION function_with_argtypes OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_FUNCTION,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER opt_procedural LANGUAGE name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_LANGUAGE,
				Object:     &nodes.String{Str: $4},
				Newowner:   $7.(*nodes.RoleSpec),
			}
		}
	| ALTER LARGE_P OBJECT_P NumericOnly OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_LARGEOBJECT,
				Object:     $4,
				Newowner:   $7.(*nodes.RoleSpec),
			}
		}
	| ALTER OPERATOR operator_with_argtypes OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_OPERATOR,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER OPERATOR CLASS any_name USING name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_OPCLASS,
				Object:     prependList(&nodes.String{Str: $6}, $4),
				Newowner:   $9.(*nodes.RoleSpec),
			}
		}
	| ALTER OPERATOR FAMILY any_name USING name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_OPFAMILY,
				Object:     prependList(&nodes.String{Str: $6}, $4),
				Newowner:   $9.(*nodes.RoleSpec),
			}
		}
	| ALTER PROCEDURE function_with_argtypes OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_PROCEDURE,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER ROUTINE function_with_argtypes OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_ROUTINE,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER SCHEMA name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_SCHEMA,
				Object:     &nodes.String{Str: $3},
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER TYPE_P any_name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_TYPE,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER TABLESPACE name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_TABLESPACE,
				Object:     &nodes.String{Str: $3},
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER STATISTICS any_name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_STATISTIC_EXT,
				Object:     $3,
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER TEXT_P SEARCH DICTIONARY any_name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_TSDICTIONARY,
				Object:     $5,
				Newowner:   $8.(*nodes.RoleSpec),
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_TSCONFIGURATION,
				Object:     $5,
				Newowner:   $8.(*nodes.RoleSpec),
			}
		}
	| ALTER FOREIGN DATA_P WRAPPER name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_FDW,
				Object:     &nodes.String{Str: $5},
				Newowner:   $8.(*nodes.RoleSpec),
			}
		}
	| ALTER SERVER name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_FOREIGN_SERVER,
				Object:     &nodes.String{Str: $3},
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER EVENT TRIGGER name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_EVENT_TRIGGER,
				Object:     &nodes.String{Str: $4},
				Newowner:   $7.(*nodes.RoleSpec),
			}
		}
	| ALTER PUBLICATION name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_PUBLICATION,
				Object:     &nodes.String{Str: $3},
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	| ALTER SUBSCRIPTION name OWNER TO RoleSpec
		{
			$$ = &nodes.AlterOwnerStmt{
				ObjectType: nodes.OBJECT_SUBSCRIPTION,
				Object:     &nodes.String{Str: $3},
				Newowner:   $6.(*nodes.RoleSpec),
			}
		}
	;

/*****************************************************************************
 *
 * ALTER DEFAULT PRIVILEGES
 *
 *****************************************************************************/

AlterDefaultPrivilegesStmt:
	ALTER DEFAULT PRIVILEGES DefACLOptionList DefACLAction
		{
			$$ = &nodes.AlterDefaultPrivilegesStmt{
				Options: $4,
				Action:  $5.(*nodes.GrantStmt),
			}
		}
	;

DefACLOptionList:
	DefACLOptionList DefACLOption
		{ $$ = appendList($1, $2) }
	| /* EMPTY */
		{ $$ = nil }
	;

DefACLOption:
	IN_P SCHEMA name_list
		{
			$$ = &nodes.DefElem{
				Defname: "schemas",
				Arg:     $3,
				Location: -1,
			}
		}
	| FOR ROLE role_list
		{
			$$ = &nodes.DefElem{
				Defname: "roles",
				Arg:     $3,
				Location: -1,
			}
		}
	| FOR USER role_list
		{
			$$ = &nodes.DefElem{
				Defname: "roles",
				Arg:     $3,
				Location: -1,
			}
		}
	;

DefACLAction:
	GRANT privileges ON defacl_privilege_target TO grantee_list opt_grant_grant_option
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     true,
				Privileges:  $2,
				Targtype:    nodes.ACL_TARGET_DEFAULTS,
				Objtype:     nodes.ObjectType($4),
				Grantees:    $6,
				GrantOption: $7,
			}
		}
	| REVOKE privileges ON defacl_privilege_target FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:    false,
				Privileges: $2,
				Targtype:   nodes.ACL_TARGET_DEFAULTS,
				Objtype:    nodes.ObjectType($4),
				Grantees:   $6,
				Behavior:   nodes.DropBehavior($7),
			}
		}
	| REVOKE GRANT OPTION FOR privileges ON defacl_privilege_target FROM grantee_list opt_drop_behavior
		{
			$$ = &nodes.GrantStmt{
				IsGrant:     false,
				GrantOption: true,
				Privileges:  $5,
				Targtype:    nodes.ACL_TARGET_DEFAULTS,
				Objtype:     nodes.ObjectType($7),
				Grantees:    $9,
				Behavior:    nodes.DropBehavior($10),
			}
		}
	;

defacl_privilege_target:
	TABLES
		{ $$ = int64(nodes.OBJECT_TABLE) }
	| FUNCTIONS
		{ $$ = int64(nodes.OBJECT_FUNCTION) }
	| ROUTINES
		{ $$ = int64(nodes.OBJECT_FUNCTION) }
	| SEQUENCES
		{ $$ = int64(nodes.OBJECT_SEQUENCE) }
	| TYPES_P
		{ $$ = int64(nodes.OBJECT_TYPE) }
	| SCHEMAS
		{ $$ = int64(nodes.OBJECT_SCHEMA) }
	;

/*****************************************************************************
 *
 * ALTER TEXT SEARCH DICTIONARY
 *
 *****************************************************************************/

AlterTSDictionaryStmt:
	ALTER TEXT_P SEARCH DICTIONARY any_name definition
		{
			$$ = &nodes.AlterTSDictionaryStmt{
				Dictname: $5,
				Options:  $6,
			}
		}
	;

/*****************************************************************************
 *
 * ALTER TEXT SEARCH CONFIGURATION
 *
 *****************************************************************************/

AlterTSConfigurationStmt:
	ALTER TEXT_P SEARCH CONFIGURATION any_name ADD_P MAPPING FOR name_list WITH any_name_list
		{
			$$ = &nodes.AlterTSConfigurationStmt{
				Kind:      nodes.ALTER_TSCONFIG_ADD_MAPPING,
				Cfgname:   $5,
				Tokentype: $9,
				Dicts:     $11,
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name ALTER MAPPING FOR name_list WITH any_name_list
		{
			$$ = &nodes.AlterTSConfigurationStmt{
				Kind:      nodes.ALTER_TSCONFIG_ALTER_MAPPING_FOR_TOKEN,
				Cfgname:   $5,
				Tokentype: $9,
				Dicts:     $11,
				Override:  true,
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name ALTER MAPPING REPLACE any_name WITH any_name
		{
			$$ = &nodes.AlterTSConfigurationStmt{
				Kind:    nodes.ALTER_TSCONFIG_REPLACE_DICT,
				Cfgname: $5,
				Dicts:   makeList2($9, $11),
				Replace: true,
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name ALTER MAPPING FOR name_list REPLACE any_name WITH any_name
		{
			$$ = &nodes.AlterTSConfigurationStmt{
				Kind:      nodes.ALTER_TSCONFIG_REPLACE_DICT_FOR_TOKEN,
				Cfgname:   $5,
				Tokentype: $9,
				Dicts:     makeList2($11, $13),
				Replace:   true,
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name DROP MAPPING FOR name_list
		{
			$$ = &nodes.AlterTSConfigurationStmt{
				Kind:      nodes.ALTER_TSCONFIG_DROP_MAPPING,
				Cfgname:   $5,
				Tokentype: $9,
			}
		}
	| ALTER TEXT_P SEARCH CONFIGURATION any_name DROP MAPPING IF_P EXISTS FOR name_list
		{
			$$ = &nodes.AlterTSConfigurationStmt{
				Kind:      nodes.ALTER_TSCONFIG_DROP_MAPPING,
				Cfgname:   $5,
				Tokentype: $11,
				MissingOk: true,
			}
		}
	;

/*****************************************************************************
 *
 * CREATE STATISTICS
 *
 *****************************************************************************/

CreateStatsStmt:
	CREATE STATISTICS opt_qualified_name opt_stat_name_list ON stats_params FROM from_list
		{
			$$ = &nodes.CreateStatsStmt{
				Defnames:    $3,
				StatTypes:   $4,
				Exprs:       $6,
				Relations:   $8,
				IfNotExists: false,
			}
		}
	| CREATE STATISTICS IF_P NOT EXISTS any_name opt_stat_name_list ON stats_params FROM from_list
		{
			$$ = &nodes.CreateStatsStmt{
				Defnames:    $6,
				StatTypes:   $7,
				Exprs:       $9,
				Relations:   $11,
				IfNotExists: true,
			}
		}
	;

opt_qualified_name:
	any_name
		{ $$ = $1 }
	| /* EMPTY */
		{ $$ = nil }
	;

opt_stat_name_list:
	'(' name_list ')'
		{ $$ = $2 }
	| /* EMPTY */
		{ $$ = nil }
	;

stats_params:
	stats_param
		{ $$ = makeList($1) }
	| stats_params ',' stats_param
		{ $$ = appendList($1, $3) }
	;

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

/*****************************************************************************
 *
 * ALTER STATISTICS
 *
 *****************************************************************************/

AlterStatsStmt:
	ALTER STATISTICS any_name SET STATISTICS set_statistics_value
		{
			n := &nodes.AlterStatsStmt{
				Defnames:  $3,
				MissingOk: false,
			}
			if iv, ok := $6.(*nodes.Integer); ok {
				n.Stxstattarget = int(iv.Ival)
			}
			$$ = n
		}
	| ALTER STATISTICS IF_P EXISTS any_name SET STATISTICS set_statistics_value
		{
			n := &nodes.AlterStatsStmt{
				Defnames:  $5,
				MissingOk: true,
			}
			if iv, ok := $8.(*nodes.Integer); ok {
				n.Stxstattarget = int(iv.Ival)
			}
			$$ = n
		}
	;

set_statistics_value:
	SignedIconst
		{ $$ = &nodes.Integer{Ival: $1} }
	| DEFAULT
		{ $$ = nil }
	;

/*****************************************************************************
 *
 * CREATE OPERATOR CLASS
 *
 *****************************************************************************/

CreateOpClassStmt:
	CREATE OPERATOR CLASS any_name opt_default FOR TYPE_P Typename USING name opt_opfamily AS opclass_item_list
		{
			$$ = &nodes.CreateOpClassStmt{
				Opclassname:  $4,
				IsDefault:    $5,
				Datatype:     $8,
				Amname:       $10,
				Opfamilyname: $11,
				Items:        $13,
			}
		}
	;

opclass_item_list:
	opclass_item
		{ $$ = makeList($1) }
	| opclass_item_list ',' opclass_item
		{ $$ = appendList($1, $3) }
	;

opclass_item:
	OPERATOR Iconst any_operator opclass_purpose opt_recheck
		{
			owa := &nodes.ObjectWithArgs{
				Objname: $3,
			}
			$$ = &nodes.CreateOpClassItem{
				Itemtype:    nodes.OPCLASS_ITEM_OPERATOR,
				Name:        owa,
				Number:      int($2),
				OrderFamily: $4,
			}
		}
	| OPERATOR Iconst operator_with_argtypes opclass_purpose opt_recheck
		{
			$$ = &nodes.CreateOpClassItem{
				Itemtype:    nodes.OPCLASS_ITEM_OPERATOR,
				Name:        $3.(*nodes.ObjectWithArgs),
				Number:      int($2),
				OrderFamily: $4,
			}
		}
	| FUNCTION Iconst function_with_argtypes
		{
			$$ = &nodes.CreateOpClassItem{
				Itemtype: nodes.OPCLASS_ITEM_FUNCTION,
				Name:     $3.(*nodes.ObjectWithArgs),
				Number:   int($2),
			}
		}
	| FUNCTION Iconst '(' type_list ')' function_with_argtypes
		{
			$$ = &nodes.CreateOpClassItem{
				Itemtype:  nodes.OPCLASS_ITEM_FUNCTION,
				Name:      $6.(*nodes.ObjectWithArgs),
				Number:    int($2),
				ClassArgs: $4,
			}
		}
	| STORAGE Typename
		{
			$$ = &nodes.CreateOpClassItem{
				Itemtype:   nodes.OPCLASS_ITEM_STORAGETYPE,
				Storedtype: $2,
			}
		}
	;

opt_default:
	DEFAULT
		{ $$ = true }
	| /* EMPTY */
		{ $$ = false }
	;

opt_opfamily:
	FAMILY any_name
		{ $$ = $2 }
	| /* EMPTY */
		{ $$ = nil }
	;

opclass_purpose:
	FOR SEARCH
		{ $$ = nil }
	| FOR ORDER BY any_name
		{ $$ = $4 }
	| /* EMPTY */
		{ $$ = nil }
	;

opt_recheck:
	RECHECK
		{ $$ = true }
	| /* EMPTY */
		{ $$ = false }
	;

/*****************************************************************************
 *
 * CREATE OPERATOR FAMILY
 *
 *****************************************************************************/

CreateOpFamilyStmt:
	CREATE OPERATOR FAMILY any_name USING name
		{
			$$ = &nodes.CreateOpFamilyStmt{
				Opfamilyname: $4,
				Amname:       $6,
			}
		}
	;

/*****************************************************************************
 *
 * ALTER OPERATOR FAMILY
 *
 *****************************************************************************/

AlterOpFamilyStmt:
	ALTER OPERATOR FAMILY any_name USING name ADD_P opclass_item_list
		{
			$$ = &nodes.AlterOpFamilyStmt{
				Opfamilyname: $4,
				Amname:       $6,
				IsDrop:       false,
				Items:        $8,
			}
		}
	| ALTER OPERATOR FAMILY any_name USING name DROP opclass_drop_list
		{
			$$ = &nodes.AlterOpFamilyStmt{
				Opfamilyname: $4,
				Amname:       $6,
				IsDrop:       true,
				Items:        $8,
			}
		}
	;

opclass_drop_list:
	opclass_drop
		{ $$ = makeList($1) }
	| opclass_drop_list ',' opclass_drop
		{ $$ = appendList($1, $3) }
	;

opclass_drop:
	OPERATOR Iconst '(' type_list ')'
		{
			$$ = &nodes.CreateOpClassItem{
				Itemtype:  nodes.OPCLASS_ITEM_OPERATOR,
				Number:    int($2),
				ClassArgs: $4,
			}
		}
	| FUNCTION Iconst '(' type_list ')'
		{
			$$ = &nodes.CreateOpClassItem{
				Itemtype:  nodes.OPCLASS_ITEM_FUNCTION,
				Number:    int($2),
				ClassArgs: $4,
			}
		}
	;

/*****************************************************************************
 *
 * DROP OPERATOR CLASS / FAMILY
 *
 *****************************************************************************/

DropOpClassStmt:
	DROP OPERATOR CLASS any_name USING name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeList(prependList(&nodes.String{Str: $6}, $4)),
				RemoveType: int(nodes.OBJECT_OPCLASS),
				Behavior:   int($7),
			}
		}
	| DROP OPERATOR CLASS IF_P EXISTS any_name USING name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeList(prependList(&nodes.String{Str: $8}, $6)),
				RemoveType: int(nodes.OBJECT_OPCLASS),
				Behavior:   int($9),
				Missing_ok: true,
			}
		}
	;

DropOpFamilyStmt:
	DROP OPERATOR FAMILY any_name USING name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeList(prependList(&nodes.String{Str: $6}, $4)),
				RemoveType: int(nodes.OBJECT_OPFAMILY),
				Behavior:   int($7),
			}
		}
	| DROP OPERATOR FAMILY IF_P EXISTS any_name USING name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				Objects:    makeList(prependList(&nodes.String{Str: $8}, $6)),
				RemoveType: int(nodes.OBJECT_OPFAMILY),
				Behavior:   int($9),
				Missing_ok: true,
			}
		}
	;

/*****************************************************************************
 *
 * CREATE CAST / DROP CAST
 *
 *****************************************************************************/

CreateCastStmt:
	CREATE CAST '(' Typename AS Typename ')' WITH FUNCTION function_with_argtypes cast_context
		{
			$$ = &nodes.CreateCastStmt{
				Sourcetype: $4,
				Targettype: $6,
				Func:       $10.(*nodes.ObjectWithArgs),
				Context:    nodes.CoercionContext($11),
			}
		}
	| CREATE CAST '(' Typename AS Typename ')' WITHOUT FUNCTION cast_context
		{
			$$ = &nodes.CreateCastStmt{
				Sourcetype: $4,
				Targettype: $6,
				Context:    nodes.CoercionContext($10),
			}
		}
	| CREATE CAST '(' Typename AS Typename ')' WITH INOUT cast_context
		{
			$$ = &nodes.CreateCastStmt{
				Sourcetype: $4,
				Targettype: $6,
				Context:    nodes.CoercionContext($10),
				Inout:      true,
			}
		}
	;

cast_context:
	AS IMPLICIT_P
		{ $$ = int64(nodes.COERCION_IMPLICIT) }
	| AS ASSIGNMENT
		{ $$ = int64(nodes.COERCION_ASSIGNMENT) }
	| /* EMPTY */
		{ $$ = int64(nodes.COERCION_EXPLICIT) }
	;

DropCastStmt:
	DROP CAST opt_if_exists '(' Typename AS Typename ')' opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_CAST),
				Objects:    makeList(makeList2($5, $7)),
				Behavior:   int($9),
				Missing_ok: $3,
			}
		}
	;

opt_if_exists:
	IF_P EXISTS
		{ $$ = true }
	| /* EMPTY */
		{ $$ = false }
	;

/*****************************************************************************
 *
 * CREATE TRANSFORM / DROP TRANSFORM
 *
 *****************************************************************************/

CreateTransformStmt:
	CREATE opt_or_replace TRANSFORM FOR Typename LANGUAGE name '(' transform_element_list ')'
		{
			items := $9.Items
			var fromsql *nodes.ObjectWithArgs
			var tosql *nodes.ObjectWithArgs
			if items[0] != nil {
				fromsql = items[0].(*nodes.ObjectWithArgs)
			}
			if items[1] != nil {
				tosql = items[1].(*nodes.ObjectWithArgs)
			}
			$$ = &nodes.CreateTransformStmt{
				Replace:  $2,
				TypeName: $5,
				Lang:     $7,
				Fromsql:  fromsql,
				Tosql:    tosql,
			}
		}
	;

transform_element_list:
	FROM SQL_P WITH FUNCTION function_with_argtypes ',' TO SQL_P WITH FUNCTION function_with_argtypes
		{
			$$ = makeList2($5, $11)
		}
	| TO SQL_P WITH FUNCTION function_with_argtypes ',' FROM SQL_P WITH FUNCTION function_with_argtypes
		{
			$$ = makeList2($11, $5)
		}
	| FROM SQL_P WITH FUNCTION function_with_argtypes
		{
			$$ = makeList2($5, nil)
		}
	| TO SQL_P WITH FUNCTION function_with_argtypes
		{
			$$ = makeList2(nil, $5)
		}
	;

DropTransformStmt:
	DROP TRANSFORM opt_if_exists FOR Typename LANGUAGE name opt_drop_behavior
		{
			$$ = &nodes.DropStmt{
				RemoveType: int(nodes.OBJECT_TRANSFORM),
				Objects:    makeList(makeList2($5, &nodes.String{Str: $7})),
				Behavior:   int($8),
				Missing_ok: $3,
			}
		}
	;

/*****************************************************************************
 *
 * CREATE CONVERSION
 *
 *****************************************************************************/

CreateConversionStmt:
	CREATE opt_default CONVERSION_P any_name FOR Sconst TO Sconst FROM any_name
		{
			$$ = &nodes.CreateConversionStmt{
				ConversionName:  $4,
				ForEncodingName: $6,
				ToEncodingName:  $8,
				FuncName:        $10,
				Def:             $2,
			}
		}
	;

/*****************************************************************************
 *
 * DROP OWNED BY / REASSIGN OWNED BY
 *
 *****************************************************************************/

DropOwnedStmt:
	DROP OWNED BY role_list opt_drop_behavior
		{
			$$ = &nodes.DropOwnedStmt{
				Roles:    $4,
				Behavior: nodes.DropBehavior($5),
			}
		}
	;

ReassignOwnedStmt:
	REASSIGN OWNED BY role_list TO RoleSpec
		{
			$$ = &nodes.ReassignOwnedStmt{
				Roles:   $4,
				Newrole: $6.(*nodes.RoleSpec),
			}
		}
	;

/*****************************************************************************
 *
 * CREATE TABLE AS / CREATE MATERIALIZED VIEW
 *
 *****************************************************************************/

CreateAsStmt:
	CREATE OptTemp TABLE create_as_target AS SelectStmt opt_with_data
		{
			into := $4.(*nodes.IntoClause)
			into.Rel.Relpersistence = relpersistenceForTemp($2)
			into.SkipData = !$7
			$$ = &nodes.CreateTableAsStmt{
				Query:       $6,
				Into:        into,
				Objtype:     nodes.OBJECT_TABLE,
				IfNotExists: false,
			}
		}
	| CREATE OptTemp TABLE IF_P NOT EXISTS create_as_target AS SelectStmt opt_with_data
		{
			into := $7.(*nodes.IntoClause)
			into.Rel.Relpersistence = relpersistenceForTemp($2)
			into.SkipData = !$10
			$$ = &nodes.CreateTableAsStmt{
				Query:       $9,
				Into:        into,
				Objtype:     nodes.OBJECT_TABLE,
				IfNotExists: true,
			}
		}
	| CREATE OptTemp TABLE create_as_target AS EXECUTE name execute_param_clause opt_with_data
		{
			into := $4.(*nodes.IntoClause)
			into.Rel.Relpersistence = relpersistenceForTemp($2)
			into.SkipData = !$9
			es := &nodes.ExecuteStmt{
				Name:   $7,
				Params: $8,
			}
			$$ = &nodes.CreateTableAsStmt{
				Query:       es,
				Into:        into,
				Objtype:     nodes.OBJECT_TABLE,
			}
		}
	| CREATE OptTemp TABLE IF_P NOT EXISTS create_as_target AS EXECUTE name execute_param_clause opt_with_data
		{
			into := $7.(*nodes.IntoClause)
			into.Rel.Relpersistence = relpersistenceForTemp($2)
			into.SkipData = !$12
			es := &nodes.ExecuteStmt{
				Name:   $10,
				Params: $11,
			}
			$$ = &nodes.CreateTableAsStmt{
				Query:       es,
				Into:        into,
				Objtype:     nodes.OBJECT_TABLE,
				IfNotExists: true,
			}
		}
	;

create_as_target:
	qualified_name opt_column_list OptAccessMethod OptWith OnCommitOption OptTableSpace
		{
			rv := makeRangeVarFromAnyName($1)
			$$ = &nodes.IntoClause{
				Rel:            rv,
				ColNames:       $2,
				AccessMethod:   $3,
				Options:        $4,
				OnCommit:       nodes.OnCommitAction($5),
				TableSpaceName: $6,
			}
		}
	;

opt_with_data:
	WITH DATA_P
		{ $$ = true }
	| WITH NO DATA_P
		{ $$ = false }
	| /* EMPTY */
		{ $$ = true }
	;

CreateMatViewStmt:
	CREATE MATERIALIZED VIEW create_mv_target AS SelectStmt opt_with_data
		{
			into := $4.(*nodes.IntoClause)
			into.SkipData = !$7
			$$ = &nodes.CreateTableAsStmt{
				Query:       $6,
				Into:        into,
				Objtype:     nodes.OBJECT_MATVIEW,
				IfNotExists: false,
			}
		}
	| CREATE MATERIALIZED VIEW IF_P NOT EXISTS create_mv_target AS SelectStmt opt_with_data
		{
			into := $7.(*nodes.IntoClause)
			into.SkipData = !$10
			$$ = &nodes.CreateTableAsStmt{
				Query:       $9,
				Into:        into,
				Objtype:     nodes.OBJECT_MATVIEW,
				IfNotExists: true,
			}
		}
	;

create_mv_target:
	qualified_name opt_column_list OptAccessMethod opt_reloptions OptTableSpace
		{
			rv := makeRangeVarFromAnyName($1)
			$$ = &nodes.IntoClause{
				Rel:            rv,
				ColNames:       $2,
				AccessMethod:   $3,
				Options:        $4,
				TableSpaceName: $5,
			}
		}
	;

/*****************************************************************************
 *
 * REFRESH MATERIALIZED VIEW
 *
 *****************************************************************************/

RefreshMatViewStmt:
	REFRESH MATERIALIZED VIEW opt_concurrently qualified_name opt_with_data
		{
			rv := makeRangeVarFromAnyName($5)
			$$ = &nodes.RefreshMatViewStmt{
				Concurrent: $4,
				Relation:   rv,
				SkipData:   !$6,
			}
		}
	;


type_function_name:
	IDENT { $$ = $1 }
	| unreserved_keyword { $$ = $1 }
	| type_func_name_keyword { $$ = $1 }
	;

// Keyword categories
// These lists must match the keyword definitions in keywords.go

unreserved_keyword:
	ABORT_P
	| ABSENT
	| ABSOLUTE_P
	| ACCESS
	| ACTION
	| ADD_P
	| ADMIN
	| AFTER
	| AGGREGATE
	| ALSO
	| ALTER
	| ALWAYS
	| ASENSITIVE
	| ASSERTION
	| ASSIGNMENT
	| AT
	| ATOMIC
	| ATTACH
	| ATTRIBUTE
	| BACKWARD
	| BEFORE
	| BEGIN_P
	| BREADTH
	| BY
	| CACHE
	| CALL
	| CALLED
	| CASCADE
	| CASCADED
	| CATALOG_P
	| CHAIN
	| CHARACTERISTICS
	| CHECKPOINT
	| CLASS
	| CLOSE
	| CLUSTER
	| COLUMNS
	| COMMENT
	| COMMENTS
	| COMMIT
	| COMMITTED
	| COMPRESSION
	| CONDITIONAL
	| CONFIGURATION
	| CONFLICT
	| CONNECTION
	| CONSTRAINTS
	| CONTENT_P
	| CONTINUE_P
	| CONVERSION_P
	| COPY
	| COST
	| CSV
	| CUBE
	| CURRENT_P
	| CURSOR
	| CYCLE
	| DATA_P
	| DATABASE
	| DAY_P
	| DEALLOCATE
	| DECLARE
	| DEFAULTS
	| DEFERRED
	| DEFINER
	| DELETE_P
	| DELIMITER
	| DELIMITERS
	| DEPENDS
	| DEPTH
	| DETACH
	| DICTIONARY
	| DISABLE_P
	| DISCARD
	| DOCUMENT_P
	| DOMAIN_P
	| DOUBLE_P
	| DROP
	| EACH
	| EMPTY_P
	| ENABLE_P
	| ENCODING
	| ENCRYPTED
	| ENUM_P
	| ERROR_P
	| ESCAPE
	| EVENT
	| EXCLUDE
	| EXCLUDING
	| EXCLUSIVE
	| EXECUTE
	| EXPLAIN
	| EXPRESSION
	| EXTENSION
	| EXTERNAL
	| FAMILY
	| FILTER
	| FINALIZE
	| FIRST_P
	| FOLLOWING
	| FORCE
	| FORMAT
	| FORWARD
	| FUNCTION
	| FUNCTIONS
	| GENERATED
	| GLOBAL
	| GRANTED
	| GROUPS
	| HANDLER
	| HEADER_P
	| HOLD
	| HOUR_P
	| IDENTITY_P
	| IF_P
	| IMMEDIATE
	| IMMUTABLE
	| IMPLICIT_P
	| IMPORT_P
	| INCLUDE
	| INCLUDING
	| INCREMENT
	| INDENT
	| INDEX
	| INDEXES
	| INHERIT
	| INHERITS
	| INLINE_P
	| INPUT_P
	| INSENSITIVE
	| INSERT
	| INSTEAD
	| INVOKER
	| ISOLATION
	| KEEP
	| KEY
	| KEYS
	| LABEL
	| LANGUAGE
	| LARGE_P
	| LAST_P
	| LEAKPROOF
	| LEVEL
	| LISTEN
	| LOAD
	| LOCAL
	| LOCATION
	| LOCK_P
	| LOCKED
	| LOGGED
	| MAPPING
	| MATCH
	| MATCHED
	| MATERIALIZED
	| MAXVALUE
	| MERGE
	| METHOD
	| MINUTE_P
	| MINVALUE
	| MODE
	| MONTH_P
	| MOVE
	| NAME_P
	| NAMES
	| NESTED
	| NEW
	| NEXT
	| NFC
	| NFD
	| NFKC
	| NFKD
	| NO
	| NORMALIZED
	| NOTHING
	| NOTIFY
	| NOWAIT
	| NULLS_P
	| OBJECT_P
	| OF
	| OFF
	| OIDS
	| OLD
	| OMIT
	| OPERATOR
	| OPTION
	| OPTIONS
	| ORDINALITY
	| OTHERS
	| OVER
	| OVERRIDING
	| OWNED
	| OWNER
	| PARALLEL
	| PARAMETER
	| PARSER
	| PARTIAL
	| PARTITION
	| PASSING
	| PASSWORD
	| PATH
	| PLAN
	| PLANS
	| POLICY
	| PRECEDING
	| PREPARE
	| PREPARED
	| PRESERVE
	| PRIOR
	| PRIVILEGES
	| PROCEDURAL
	| PROCEDURE
	| PROCEDURES
	| PROGRAM
	| PUBLICATION
	| QUOTE
	| QUOTES
	| RANGE
	| READ
	| REASSIGN
	| RECHECK
	| RECURSIVE
	| REF_P
	| REFERENCING
	| REFRESH
	| REINDEX
	| RELATIVE_P
	| RELEASE
	| RENAME
	| REPEATABLE
	| REPLACE
	| REPLICA
	| RESET
	| RESTART
	| RESTRICT
	| RETURN
	| RETURNS
	| REVOKE
	| ROLE
	| ROLLBACK
	| ROLLUP
	| ROUTINE
	| ROUTINES
	| ROWS
	| RULE
	| SAVEPOINT
	| SCALAR
	| SCHEMA
	| SCHEMAS
	| SCROLL
	| SEARCH
	| SECOND_P
	| SECURITY
	| SEQUENCE
	| SEQUENCES
	| SERIALIZABLE
	| SERVER
	| SESSION
	| SET
	| SETS
	| SHARE
	| SHOW
	| SIMPLE
	| SKIP
	| SNAPSHOT
	| SOURCE
	| SQL_P
	| STABLE
	| STANDALONE_P
	| START
	| STATEMENT
	| STATISTICS
	| STDIN
	| STDOUT
	| STORAGE
	| STORED
	| STRICT_P
	| STRING_P
	| STRIP_P
	| SUBSCRIPTION
	| SUPPORT
	| SYSID
	| SYSTEM_P
	| TABLES
	| TABLESPACE
	| TARGET
	| TEMP
	| TEMPLATE
	| TEMPORARY
	| TEXT_P
	| TIES
	| TRANSACTION
	| TRANSFORM
	| TRIGGER
	| TRUNCATE
	| TRUSTED
	| TYPE_P
	| TYPES_P
	| UESCAPE
	| UNBOUNDED
	| UNCOMMITTED
	| UNCONDITIONAL
	| UNENCRYPTED
	| UNKNOWN
	| UNLISTEN
	| UNLOGGED
	| UNTIL
	| UPDATE
	| VACUUM
	| VALID
	| VALIDATE
	| VALIDATOR
	| VALUE_P
	| VARYING
	| VERSION_P
	| VIEW
	| VIEWS
	| VOLATILE
	| WHITESPACE_P
	| WITHIN
	| WITHOUT
	| WORK
	| WRAPPER
	| WRITE
	| XML_P
	| YEAR_P
	| YES_P
	| ZONE
	;

col_name_keyword:
	BETWEEN
	| BIGINT
	| BIT
	| BOOLEAN_P
	| CHAR_P
	| CHARACTER
	| COALESCE
	| DEC
	| DECIMAL_P
	| EXISTS
	| EXTRACT
	| FLOAT_P
	| GREATEST
	| GROUPING
	| INOUT
	| INT_P
	| INTEGER
	| INTERVAL
	| JSON
	| JSON_ARRAY
	| JSON_ARRAYAGG
	| JSON_EXISTS
	| JSON_OBJECT
	| JSON_OBJECTAGG
	| JSON_QUERY
	| JSON_SCALAR
	| JSON_SERIALIZE
	| JSON_TABLE
	| JSON_VALUE
	| LEAST
	| MERGE_ACTION
	| NATIONAL
	| NCHAR
	| NONE
	| NORMALIZE
	| NULLIF
	| NUMERIC
	| OUT_P
	| OVERLAY
	| POSITION
	| PRECISION
	| REAL
	| ROW
	| SETOF
	| SMALLINT
	| SUBSTRING
	| TIME
	| TIMESTAMP
	| TREAT
	| TRIM
	| VALUES
	| VARCHAR
	| XMLATTRIBUTES
	| XMLCONCAT
	| XMLELEMENT
	| XMLEXISTS
	| XMLFOREST
	| XMLNAMESPACES
	| XMLPARSE
	| XMLPI
	| XMLROOT
	| XMLSERIALIZE
	| XMLTABLE
	;

type_func_name_keyword:
	AUTHORIZATION
	| BINARY
	| COLLATION
	| CROSS
	| CURRENT_SCHEMA
	| FREEZE
	| FULL
	| ILIKE
	| INNER_P
	| IS
	| ISNULL
	| JOIN
	| LEFT
	| LIKE
	| NATURAL
	| NOTNULL
	| OUTER_P
	| OVERLAPS
	| RIGHT
	| SIMILAR
	| TABLESAMPLE
	| VERBOSE
	;

reserved_keyword:
	ALL
	| ANALYSE
	| ANALYZE
	| AND
	| ANY
	| ARRAY
	| AS
	| ASC
	| ASYMMETRIC
	| BOTH
	| CASE
	| CAST
	| CHECK
	| COLLATE
	| COLUMN
	| CONSTRAINT
	| CREATE
	| CURRENT_CATALOG
	| CURRENT_DATE
	| CURRENT_ROLE
	| CURRENT_TIME
	| CURRENT_TIMESTAMP
	| CURRENT_USER
	| DEFAULT
	| DEFERRABLE
	| DESC
	| DISTINCT
	| DO
	| ELSE
	| END_P
	| EXCEPT
	| FALSE_P
	| FETCH
	| FOR
	| FOREIGN
	| FROM
	| GRANT
	| GROUP_P
	| HAVING
	| IN_P
	| INITIALLY
	| INTERSECT
	| INTO
	| LATERAL_P
	| LEADING
	| LIMIT
	| LOCALTIME
	| LOCALTIMESTAMP
	| NOT
	| NULL_P
	| OFFSET
	| ON
	| ONLY
	| OR
	| ORDER
	| PLACING
	| PRIMARY
	| REFERENCES
	| RETURNING
	| SELECT
	| SESSION_USER
	| SOME
	| SYMMETRIC
	| SYSTEM_USER
	| TABLE
	| THEN
	| TO
	| TRAILING
	| TRUE_P
	| UNION
	| UNIQUE
	| USER
	| USING
	| VARIADIC
	| WHEN
	| WHERE
	| WINDOW
	| WITH
	;

%%

// OnConflict action constants
const (
	ONCONFLICT_NONE    = 0
	ONCONFLICT_NOTHING = 1
	ONCONFLICT_UPDATE  = 2
)

// ViewCheckOption constants
const (
	VIEW_CHECK_OPTION_NONE     = 0
	VIEW_CHECK_OPTION_LOCAL    = 1
	VIEW_CHECK_OPTION_CASCADED = 2
)

// Helper functions called from grammar actions

// setParseResult stores the parse result in the lexer.
func setParseResult(lex pgLexer, result *nodes.List) {
	if pl, ok := lex.(*parserLexer); ok {
		pl.result = result
	}
}

func makeList(n nodes.Node) *nodes.List {
	if n == nil {
		return &nodes.List{}
	}
	return &nodes.List{Items: []nodes.Node{n}}
}

func appendList(l *nodes.List, n nodes.Node) *nodes.List {
	if l == nil {
		return makeList(n)
	}
	if n != nil {
		l.Items = append(l.Items, n)
	}
	return l
}

func prependList(n nodes.Node, l *nodes.List) *nodes.List {
	if l == nil {
		return makeList(n)
	}
	if n != nil {
		l.Items = append([]nodes.Node{n}, l.Items...)
	}
	return l
}

func makeList2(a nodes.Node, b nodes.Node) *nodes.List {
	return &nodes.List{Items: []nodes.Node{a, b}}
}

func makeListNode(l *nodes.List) nodes.Node {
	return l
}

func makeRangeVar(names *nodes.List) nodes.Node {
	rv := &nodes.RangeVar{Inh: true}
	if names != nil && len(names.Items) > 0 {
		switch len(names.Items) {
		case 1:
			rv.Relname = names.Items[0].(*nodes.String).Str
		case 2:
			rv.Schemaname = names.Items[0].(*nodes.String).Str
			rv.Relname = names.Items[1].(*nodes.String).Str
		case 3:
			rv.Catalogname = names.Items[0].(*nodes.String).Str
			rv.Schemaname = names.Items[1].(*nodes.String).Str
			rv.Relname = names.Items[2].(*nodes.String).Str
		}
	}
	return rv
}

func makeRangeVarList(nameList *nodes.List) *nodes.List {
	if nameList == nil {
		return nil
	}
	result := &nodes.List{}
	for _, item := range nameList.Items {
		rv := makeRangeVar(item.(*nodes.List))
		result.Items = append(result.Items, rv)
	}
	return result
}

func makeAExpr(kind nodes.A_Expr_Kind, op string, lexpr, rexpr nodes.Node) nodes.Node {
	return &nodes.A_Expr{
		Kind:  kind,
		Name:  &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}},
		Lexpr: lexpr,
		Rexpr: rexpr,
	}
}

// makeAExprFromList creates an A_Expr with an operator name list (from qual_Op).
// makeNameListAsAnyNameList wraps each name string into a single-element
// List (matching PG's any_name format for DropStmt.Objects).
func makeNameListAsAnyNameList(names *nodes.List) *nodes.List {
	if names == nil {
		return nil
	}
	result := &nodes.List{}
	for _, item := range names.Items {
		result.Items = append(result.Items, &nodes.List{Items: []nodes.Node{item}})
	}
	return result
}

func makeStringConstCast(s string, typeName *nodes.TypeName) nodes.Node {
	return &nodes.TypeCast{
		Arg:      &nodes.A_Const{Val: &nodes.String{Str: s}},
		TypeName: typeName,
		Location: -1,
	}
}

func parsePartitionStrategy(strategy string) string {
	switch strategy {
	case "list", "LIST", "List":
		return "l"
	case "range", "RANGE", "Range":
		return "r"
	case "hash", "HASH", "Hash":
		return "h"
	default:
		return strategy
	}
}

func makeAExprFromList(kind nodes.A_Expr_Kind, name *nodes.List, lexpr, rexpr nodes.Node) nodes.Node {
	return &nodes.A_Expr{
		Kind:  kind,
		Name:  name,
		Lexpr: lexpr,
		Rexpr: rexpr,
	}
}

func makeBoolExpr(boolop nodes.BoolExprType, arg1, arg2 nodes.Node) nodes.Node {
	be := &nodes.BoolExpr{
		Boolop: boolop,
		Args:   &nodes.List{},
	}
	if arg1 != nil {
		be.Args.Items = append(be.Args.Items, arg1)
	}
	if arg2 != nil {
		be.Args.Items = append(be.Args.Items, arg2)
	}
	return be
}

func makeBetweenArgs(lower, upper nodes.Node) nodes.Node {
	return &nodes.List{Items: []nodes.Node{lower, upper}}
}

func doNegate(n nodes.Node) nodes.Node {
	// For numeric constants, negate in place
	if ac, ok := n.(*nodes.A_Const); ok {
		if i, ok := ac.Val.(*nodes.Integer); ok {
			i.Ival = -i.Ival
			return n
		}
		if f, ok := ac.Val.(*nodes.Float); ok {
			if f.Fval[0] == '-' {
				f.Fval = f.Fval[1:]
			} else {
				f.Fval = "-" + f.Fval
			}
			return n
		}
	}
	// Otherwise, create unary minus expression
	return makeAExpr(nodes.AEXPR_OP, "-", nil, n)
}

func concatLists(a, b *nodes.List) *nodes.List {
	if a == nil { return b }
	if b == nil { return a }
	result := &nodes.List{Items: make([]nodes.Node, 0, len(a.Items)+len(b.Items))}
	result.Items = append(result.Items, a.Items...)
	result.Items = append(result.Items, b.Items...)
	return result
}

func makeSetOp(op nodes.SetOperation, all bool, larg, rarg nodes.Node) nodes.Node {
	n := &nodes.SelectStmt{
		Op:   op,
		All:  all,
		Larg: larg.(*nodes.SelectStmt),
		Rarg: rarg.(*nodes.SelectStmt),
	}
	return n
}

func makeTypeName(typeName string) *nodes.TypeName {
	return &nodes.TypeName{
		Names: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "pg_catalog"},
			&nodes.String{Str: typeName},
		}},
		Location: -1,
	}
}

func makeIntConst(val int64) nodes.Node {
	return &nodes.A_Const{Val: &nodes.Integer{Ival: val}}
}

func makeStringConst(str string) nodes.Node {
	return &nodes.A_Const{Val: &nodes.String{Str: str}}
}

func makeBoolAConst(val int64) nodes.Node {
	if val != 0 {
		return &nodes.A_Const{Val: &nodes.String{Str: "t"}}
	}
	return &nodes.A_Const{Val: &nodes.String{Str: "f"}}
}

func makeNullAConst() nodes.Node {
	return &nodes.A_Const{Isnull: true}
}

func makeNotExpr(expr nodes.Node) nodes.Node {
	return &nodes.BoolExpr{
		Boolop: nodes.NOT_EXPR,
		Args:   &nodes.List{Items: []nodes.Node{expr}},
		Location: -1,
	}
}

func setFuncAlias(n *nodes.RangeFunction, clause nodes.Node) {
	if clause == nil {
		return
	}
	// If it's an Alias, just set the alias
	if alias, ok := clause.(*nodes.Alias); ok {
		n.Alias = alias
		return
	}
	// If it's a List with [alias, coldeflist], it's a func_alias_clause with column defs
	if list, ok := clause.(*nodes.List); ok && len(list.Items) == 2 {
		if alias, ok2 := list.Items[0].(*nodes.Alias); ok2 {
			n.Alias = alias
		}
		if coldeflist, ok2 := list.Items[1].(*nodes.List); ok2 {
			n.Coldeflist = coldeflist
		}
	}
}

func setJoinQual(n *nodes.JoinExpr, qual nodes.Node) {
	if qual == nil {
		return
	}
	if list, ok := qual.(*nodes.List); ok {
		// Check if this is a USING clause with alias: [nameList, alias]
		if len(list.Items) == 2 {
			if innerList, ok2 := list.Items[0].(*nodes.List); ok2 {
				if alias, ok3 := list.Items[1].(*nodes.Alias); ok3 {
					n.UsingClause = innerList
					n.JoinUsing = alias
					return
				}
			}
		}
		n.UsingClause = list
	} else {
		n.Quals = qual
	}
}

func makeTypeNameFromNameList(names *nodes.List) nodes.Node {
	tn := &nodes.TypeName{Location: -1}
	if names != nil {
		tn.Names = names
	}
	return tn
}

func asJsonOutput(n nodes.Node) *nodes.JsonOutput {
	if n == nil {
		return nil
	}
	return n.(*nodes.JsonOutput)
}

func asJsonBehavior(n nodes.Node) *nodes.JsonBehavior {
	if n == nil {
		return nil
	}
	return n.(*nodes.JsonBehavior)
}

func splitJsonBehaviorClause(n nodes.Node) (*nodes.JsonBehavior, *nodes.JsonBehavior) {
	if n == nil {
		return nil, nil
	}
	list := n.(*nodes.List)
	var onEmpty, onError *nodes.JsonBehavior
	if list.Items[0] != nil {
		onEmpty = list.Items[0].(*nodes.JsonBehavior)
	}
	if list.Items[1] != nil {
		onError = list.Items[1].(*nodes.JsonBehavior)
	}
	return onEmpty, onError
}

func asJsonTablePathSpec(n nodes.Node) *nodes.JsonTablePathSpec {
	if n == nil {
		return nil
	}
	return n.(*nodes.JsonTablePathSpec)
}

func doNegateFloat(f *nodes.Float) {
	if len(f.Fval) > 0 && f.Fval[0] == '-' {
		f.Fval = f.Fval[1:]
	} else {
		f.Fval = "-" + f.Fval
	}
}

func makeDefElem(name string, arg nodes.Node) nodes.Node {
	return &nodes.DefElem{
		Defname:  name,
		Arg:      arg,
		Location: -1,
	}
}

// makeFuncName converts a function name string to a *nodes.List of String nodes.
// With one argument, creates an unqualified name.
// With two arguments, creates a schema-qualified name (e.g. "pg_catalog", "func").
func makeFuncName(parts ...string) *nodes.List {
	items := make([]nodes.Node, len(parts))
	for i, p := range parts {
		items[i] = &nodes.String{Str: p}
	}
	return &nodes.List{Items: items}
}

// makeSQLValueFunction creates a SQLValueFunction node.
func makeSQLValueFunction(op nodes.SVFOp, typmod int) nodes.Node {
	return &nodes.SQLValueFunction{Op: op, Typmod: int32(typmod), Location: -1}
}

// makeTypeCast creates a TypeCast node.
func makeTypeCast(arg nodes.Node, typeName *nodes.TypeName) nodes.Node {
	return &nodes.TypeCast{Arg: arg, TypeName: typeName, Location: -1}
}

// roleSpecOrNil safely casts a node to *nodes.RoleSpec, returning nil if the node is nil.
func roleSpecOrNil(n nodes.Node) *nodes.RoleSpec {
	if n == nil {
		return nil
	}
	return n.(*nodes.RoleSpec)
}

// SelectLimit is an internal helper struct for passing LIMIT/OFFSET through grammar rules.
// It is not a Node type - just used during parsing.
type SelectLimit struct {
	LimitOffset nodes.Node
	LimitCount  nodes.Node
	LimitOption nodes.LimitOption
}

// importQualification is an internal helper struct for passing IMPORT FOREIGN SCHEMA
// qualification type and table list through grammar rules.
// It implements the nodes.Node interface so it can be passed through %union.
type importQualification struct {
	listType  nodes.ImportForeignSchemaType
	tableList *nodes.List
}

func (n *importQualification) Tag() nodes.NodeTag { return nodes.T_Invalid }

// splitColQualList separates a list of column qualifiers (constraints and COLLATE)
// into the ColumnDef's constraints and collClause fields.
// This matches PostgreSQL's SplitColQualList function.
func applyConstraintAttrs(n *nodes.Constraint, attrs int64) {
	if attrs&int64(nodes.CAS_DEFERRABLE) != 0 {
		n.Deferrable = true
	}
	if attrs&int64(nodes.CAS_INITIALLY_DEFERRED) != 0 {
		n.Initdeferred = true
		n.Deferrable = true
	}
	if attrs&int64(nodes.CAS_NOT_VALID) != 0 {
		n.SkipValidation = true
		n.InitiallyValid = false
	}
	if attrs&int64(nodes.CAS_NO_INHERIT) != 0 {
		n.IsNoInherit = true
	}
}

func splitColQualList(qualList *nodes.List, coldef *nodes.ColumnDef) {
	if qualList == nil {
		return
	}
	var constraints []nodes.Node
	for _, item := range qualList.Items {
		if cc, ok := item.(*nodes.CollateClause); ok {
			coldef.CollClause = cc
		} else if de, ok := item.(*nodes.DefElem); ok && de.Defname == "compression" {
			if s, ok := de.Arg.(*nodes.String); ok {
				coldef.Compression = s.Str
			}
		} else if de, ok := item.(*nodes.DefElem); ok && de.Defname == "storage" {
			if s, ok := de.Arg.(*nodes.String); ok {
				coldef.StorageName = s.Str
			}
		} else if de, ok := item.(*nodes.DefElem); ok && de.Defname == "fdwoptions" {
			if l, ok := de.Arg.(*nodes.List); ok {
				coldef.Fdwoptions = l
			}
		} else {
			constraints = append(constraints, item)
		}
	}
	if len(constraints) > 0 {
		coldef.Constraints = &nodes.List{Items: constraints}
	}
}

func insertSelectOptions(stmt *nodes.SelectStmt, sortClause *nodes.List, lockingClause *nodes.List,
	limitClause *SelectLimit, withClause *nodes.WithClause) {
	if sortClause != nil {
		stmt.SortClause = sortClause
	}
	if lockingClause != nil {
		stmt.LockingClause = lockingClause
	}
	if limitClause != nil {
		stmt.LimitOffset = limitClause.LimitOffset
		stmt.LimitCount = limitClause.LimitCount
		stmt.LimitOption = limitClause.LimitOption
	}
	if withClause != nil {
		stmt.WithClause = withClause
	}
}

// relpersistenceForTemp returns the relpersistence byte based on temp flag.
// In PostgreSQL: 'p' = permanent (default), 't' = temporary
func relpersistenceForTemp(tempFlag int64) byte {
	if tempFlag == 1 {
		return 't'
	}
	return 'p'
}

// extractArgTypes extracts the type names from a list of FunctionParameter nodes.
func extractArgTypes(args *nodes.List) *nodes.List {
	if args == nil {
		return nil
	}
	result := &nodes.List{}
	for _, item := range args.Items {
		fp, ok := item.(*nodes.FunctionParameter)
		if ok {
			result.Items = append(result.Items, fp.ArgType)
		}
	}
	return result
}

// extractAggrArgTypes extracts the type names from an aggregate argument list.
// aggr_args returns a 2-element list: [args_list, Integer_marker].
// The first element is either nil (for agg(*)) or a *List of FunctionParameter nodes.
// This function extracts the ArgType from each FunctionParameter in the first element.
func extractAggrArgTypes(args *nodes.List) *nodes.List {
	if args == nil || len(args.Items) < 1 {
		return nil
	}
	// First element is the args list (may be nil for agg(*))
	argsList, ok := args.Items[0].(*nodes.List)
	if !ok || argsList == nil {
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

// makeOrderedSetArgs builds the aggr_args result for ordered-set aggregates
// with both direct args and ORDER BY args.
// Returns a 2-element list: [combined_args_list, Integer{ndirectargs}].
// The combined list has the direct args followed by the ORDER BY args.
func makeOrderedSetArgs(directArgs *nodes.List, orderedArgs *nodes.List) *nodes.List {
	combined := &nodes.List{}
	ndirect := 0
	if directArgs != nil {
		combined.Items = append(combined.Items, directArgs.Items...)
		ndirect = len(directArgs.Items)
	}
	if orderedArgs != nil {
		combined.Items = append(combined.Items, orderedArgs.Items...)
	}
	return makeList2(combined, &nodes.Integer{Ival: int64(ndirect)})
}

// checkFuncName validates that a qualified function name has at most 3 parts.
func checkFuncName(names *nodes.List) *nodes.List {
	// Just return as-is; error checking can be done later
	return names
}

// intToString converts an int64 to a string representation.
func intToString(val int64) string {
	return fmt.Sprintf("%d", val)
}

// makeRangeVarFromAnyName creates a RangeVar from a qualified name list (list of String nodes).
// It handles 1-part (name), 2-part (schema.name), and 3-part (catalog.schema.name) names.
func makeRangeVarFromAnyName(names *nodes.List) *nodes.RangeVar {
	rv := &nodes.RangeVar{
		Inh:            true,
		Relpersistence: 'p',
		Location:       -1,
	}
	if names == nil {
		return rv
	}
	switch len(names.Items) {
	case 1:
		rv.Relname = names.Items[0].(*nodes.String).Str
	case 2:
		rv.Schemaname = names.Items[0].(*nodes.String).Str
		rv.Relname = names.Items[1].(*nodes.String).Str
	case 3:
		rv.Catalogname = names.Items[0].(*nodes.String).Str
		rv.Schemaname = names.Items[1].(*nodes.String).Str
		rv.Relname = names.Items[2].(*nodes.String).Str
	}
	return rv
}
