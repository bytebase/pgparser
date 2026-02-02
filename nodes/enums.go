package nodes

// CmdType represents the type of a query command.
type CmdType int

const (
	CMD_UNKNOWN CmdType = iota
	CMD_SELECT          // select stmt
	CMD_UPDATE          // update stmt
	CMD_INSERT          // insert stmt
	CMD_DELETE          // delete stmt
	CMD_MERGE           // merge stmt
	CMD_UTILITY         // cmds like create, destroy, copy, vacuum, etc.
	CMD_NOTHING         // dummy command for instead nothing rules with qual
)

// SetOperation represents types of set operations (UNION, INTERSECT, EXCEPT).
type SetOperation int

const (
	SETOP_NONE SetOperation = iota
	SETOP_UNION
	SETOP_INTERSECT
	SETOP_EXCEPT
)

// LimitOption represents LIMIT clause options.
type LimitOption int

const (
	LIMIT_OPTION_COUNT     LimitOption = iota // FETCH FIRST... ONLY
	LIMIT_OPTION_WITH_TIES                    // FETCH FIRST... WITH TIES
)

// SortByDir represents sort ordering direction.
type SortByDir int

const (
	SORTBY_DEFAULT SortByDir = iota
	SORTBY_ASC
	SORTBY_DESC
	SORTBY_USING // not allowed in CREATE INDEX
)

// SortByNulls represents NULLS FIRST/LAST option.
type SortByNulls int

const (
	SORTBY_NULLS_DEFAULT SortByNulls = iota
	SORTBY_NULLS_FIRST
	SORTBY_NULLS_LAST
)

// SetQuantifier represents ALL/DISTINCT option.
type SetQuantifier int

const (
	SET_QUANTIFIER_DEFAULT SetQuantifier = iota
	SET_QUANTIFIER_ALL
	SET_QUANTIFIER_DISTINCT
)

// JoinType represents types of joins.
type JoinType int

const (
	JOIN_INNER JoinType = iota // matching tuple pairs only
	JOIN_LEFT                  // pairs + unmatched LHS tuples
	JOIN_FULL                  // pairs + unmatched LHS + unmatched RHS
	JOIN_RIGHT                 // pairs + unmatched RHS tuples
	JOIN_SEMI                  // LHS tuples that have match(es)
	JOIN_ANTI                  // LHS tuples that don't have a match
	JOIN_RIGHT_SEMI            // RHS tuples that have match(es)
	JOIN_RIGHT_ANTI            // RHS tuples that don't have a match
	JOIN_UNIQUE_OUTER          // LHS path must be made unique
	JOIN_UNIQUE_INNER          // RHS path must be made unique
)

// BoolExprType represents types of boolean expressions.
type BoolExprType int

const (
	AND_EXPR BoolExprType = iota
	OR_EXPR
	NOT_EXPR
)

// A_Expr_Kind represents types of A_Expr.
type A_Expr_Kind int

const (
	AEXPR_OP          A_Expr_Kind = iota // normal operator
	AEXPR_OP_ANY                         // scalar op ANY (array)
	AEXPR_OP_ALL                         // scalar op ALL (array)
	AEXPR_DISTINCT                       // IS DISTINCT FROM - name must be "="
	AEXPR_NOT_DISTINCT                   // IS NOT DISTINCT FROM - name must be "="
	AEXPR_NULLIF                         // NULLIF - name must be "="
	AEXPR_IN                             // [NOT] IN - name must be "=" or "<>"
	AEXPR_LIKE                           // [NOT] LIKE - name must be "~~" or "!~~"
	AEXPR_ILIKE                          // [NOT] ILIKE - name must be "~~*" or "!~~*"
	AEXPR_SIMILAR                        // [NOT] SIMILAR - name must be "~" or "!~"
	AEXPR_BETWEEN                        // name must be "BETWEEN"
	AEXPR_NOT_BETWEEN                    // name must be "NOT BETWEEN"
	AEXPR_BETWEEN_SYM                    // name must be "BETWEEN SYMMETRIC"
	AEXPR_NOT_BETWEEN_SYM                // name must be "NOT BETWEEN SYMMETRIC"
)

// QuerySource represents possible sources of a Query.
type QuerySource int

const (
	QSRC_ORIGINAL           QuerySource = iota // original parsetree (explicit query)
	QSRC_PARSER                                // added by parse analysis (now unused)
	QSRC_INSTEAD_RULE                          // added by unconditional INSTEAD rule
	QSRC_QUAL_INSTEAD_RULE                     // added by conditional INSTEAD rule
	QSRC_NON_INSTEAD_RULE                      // added by non-INSTEAD rule
)

// OverridingKind represents OVERRIDING clause options.
type OverridingKind int

const (
	OVERRIDING_NOT_SET OverridingKind = iota
	OVERRIDING_USER_VALUE
	OVERRIDING_SYSTEM_VALUE
)

// OnCommitAction represents ON COMMIT actions for temporary tables.
type OnCommitAction int

const (
	ONCOMMIT_NOOP          OnCommitAction = iota // No ON COMMIT clause
	ONCOMMIT_PRESERVE_ROWS                       // ON COMMIT PRESERVE ROWS
	ONCOMMIT_DELETE_ROWS                         // ON COMMIT DELETE ROWS
	ONCOMMIT_DROP                                // ON COMMIT DROP
)

// ConstrType represents constraint types.
type ConstrType int

const (
	CONSTR_NULL ConstrType = iota
	CONSTR_NOTNULL
	CONSTR_DEFAULT
	CONSTR_IDENTITY
	CONSTR_GENERATED
	CONSTR_CHECK
	CONSTR_PRIMARY
	CONSTR_UNIQUE
	CONSTR_EXCLUSION
	CONSTR_FOREIGN
	CONSTR_ATTR_DEFERRABLE
	CONSTR_ATTR_NOT_DEFERRABLE
	CONSTR_ATTR_DEFERRED
	CONSTR_ATTR_IMMEDIATE
)

// CoercionForm represents how to display a node.
type CoercionForm int

const (
	COERCE_EXPLICIT_CALL CoercionForm = iota // display as a function call
	COERCE_EXPLICIT_CAST                     // display as an explicit cast
	COERCE_IMPLICIT_CAST                     // implicit cast, so hide it
	COERCE_SQL_SYNTAX                        // TREAT/XMLROOT/etc SQL syntax
)

// DropBehavior represents RESTRICT vs CASCADE behavior.
type DropBehavior int

const (
	DROP_RESTRICT DropBehavior = iota // drop fails if any dependent objects
	DROP_CASCADE                      // remove dependent objects too
)

// ObjectType represents types of objects.
type ObjectType int

const (
	OBJECT_ACCESS_METHOD ObjectType = iota
	OBJECT_AGGREGATE
	OBJECT_AMOP
	OBJECT_AMPROC
	OBJECT_ATTRIBUTE
	OBJECT_CAST
	OBJECT_COLUMN
	OBJECT_COLLATION
	OBJECT_CONVERSION
	OBJECT_DATABASE
	OBJECT_DEFAULT
	OBJECT_DEFACL
	OBJECT_DOMAIN
	OBJECT_DOMCONSTRAINT
	OBJECT_EVENT_TRIGGER
	OBJECT_EXTENSION
	OBJECT_FDW
	OBJECT_FOREIGN_SERVER
	OBJECT_FOREIGN_TABLE
	OBJECT_FUNCTION
	OBJECT_INDEX
	OBJECT_LANGUAGE
	OBJECT_LARGEOBJECT
	OBJECT_MATVIEW
	OBJECT_OPCLASS
	OBJECT_OPERATOR
	OBJECT_OPFAMILY
	OBJECT_PARAMETER_ACL
	OBJECT_POLICY
	OBJECT_PROCEDURE
	OBJECT_PUBLICATION
	OBJECT_PUBLICATION_NAMESPACE
	OBJECT_PUBLICATION_REL
	OBJECT_ROLE
	OBJECT_ROUTINE
	OBJECT_RULE
	OBJECT_SCHEMA
	OBJECT_SEQUENCE
	OBJECT_STATISTIC_EXT
	OBJECT_SUBSCRIPTION
	OBJECT_TABCONSTRAINT
	OBJECT_TABLE
	OBJECT_TABLESPACE
	OBJECT_TRANSFORM
	OBJECT_TRIGGER
	OBJECT_TSCONFIGURATION
	OBJECT_TSDICTIONARY
	OBJECT_TSPARSER
	OBJECT_TSTEMPLATE
	OBJECT_TYPE
	OBJECT_USER_MAPPING
	OBJECT_VIEW
)

// SubLinkType represents types of SubLink.
type SubLinkType int

const (
	EXISTS_SUBLINK SubLinkType = iota
	ALL_SUBLINK
	ANY_SUBLINK
	ROWCOMPARE_SUBLINK
	EXPR_SUBLINK
	MULTIEXPR_SUBLINK
	ARRAY_SUBLINK
	CTE_SUBLINK // for SubPlans only
)

// RoleSpecType represents types of role specification.
type RoleSpecType int

const (
	ROLESPEC_CSTRING      RoleSpecType = iota // role name as string
	ROLESPEC_CURRENT_ROLE                     // CURRENT_ROLE
	ROLESPEC_CURRENT_USER                     // CURRENT_USER (synonym)
	ROLESPEC_SESSION_USER                     // SESSION_USER
	ROLESPEC_PUBLIC                           // PUBLIC
)

// AlterTableType represents ALTER TABLE command types.
type AlterTableType int

const (
	AT_AddColumn AlterTableType = iota
	AT_AddColumnToView
	AT_ColumnDefault
	AT_CookedColumnDefault
	AT_DropNotNull
	AT_SetNotNull
	AT_SetExpression
	AT_DropExpression
	AT_CheckNotNull
	AT_SetStatistics
	AT_SetOptions
	AT_ResetOptions
	AT_SetStorage
	AT_SetCompression
	AT_DropColumn
	AT_AddIndex
	AT_ReAddIndex
	AT_AddConstraint
	AT_ReAddConstraint
	AT_ReAddDomainConstraint
	AT_AlterConstraint
	AT_ValidateConstraint
	AT_AddIndexConstraint
	AT_DropConstraint
	AT_ReAddComment
	AT_AlterColumnType
	AT_AlterColumnGenericOptions
	AT_ChangeOwner
	AT_ClusterOn
	AT_DropCluster
	AT_SetLogged
	AT_SetUnLogged
	AT_DropOids
	AT_SetAccessMethod
	AT_SetTableSpace
	AT_SetRelOptions
	AT_ResetRelOptions
	AT_ReplaceRelOptions
	AT_EnableTrig
	AT_EnableAlwaysTrig
	AT_EnableReplicaTrig
	AT_DisableTrig
	AT_EnableTrigAll
	AT_DisableTrigAll
	AT_EnableTrigUser
	AT_DisableTrigUser
	AT_EnableRule
	AT_EnableAlwaysRule
	AT_EnableReplicaRule
	AT_DisableRule
	AT_AddInherit
	AT_DropInherit
	AT_AddOf
	AT_DropOf
	AT_ReplicaIdentity
	AT_EnableRowSecurity
	AT_DisableRowSecurity
	AT_ForceRowSecurity
	AT_NoForceRowSecurity
	AT_GenericOptions
	AT_AttachPartition
	AT_DetachPartition
	AT_DetachPartitionFinalize
	AT_AddIdentity
	AT_SetIdentity
	AT_DropIdentity
	AT_ReAddStatistics
)

// LockClauseStrength represents FOR UPDATE/SHARE strength.
type LockClauseStrength int

const (
	LCS_NONE LockClauseStrength = iota
	LCS_FORKEYSHARE
	LCS_FORSHARE
	LCS_FORNOKEYUPDATE
	LCS_FORUPDATE
)

// LockWaitPolicy represents NOWAIT/SKIP LOCKED option.
type LockWaitPolicy int

const (
	LockWaitBlock LockWaitPolicy = iota // default behavior: wait for lock
	LockWaitSkip                        // SKIP LOCKED
	LockWaitError                       // NOWAIT
)

// CTEMaterialize represents CTE materialization options.
type CTEMaterialize int

const (
	CTEMaterializeDefault CTEMaterialize = iota
	CTEMaterializeAlways
	CTEMaterializeNever
)

// DiscardMode represents DISCARD target types.
type DiscardMode int

const (
	DISCARD_ALL       DiscardMode = iota
	DISCARD_PLANS
	DISCARD_SEQUENCES
	DISCARD_TEMP
)

// VariableSetKind represents SET variable kinds.
type VariableSetKind int

const (
	VAR_SET_VALUE   VariableSetKind = iota // SET var = value
	VAR_SET_DEFAULT                        // SET var TO DEFAULT
	VAR_SET_CURRENT                        // SET var FROM CURRENT
	VAR_SET_MULTI                          // special case for SET TRANSACTION
	VAR_RESET                              // RESET var
	VAR_RESET_ALL                          // RESET ALL
)

// RoleStmtType for CREATE ROLE/USER/GROUP
type RoleStmtType int

const (
	ROLESTMT_ROLE  RoleStmtType = iota
	ROLESTMT_USER
	ROLESTMT_GROUP
)

// Lock mode constants (from PostgreSQL's lockdefs.h)
const (
	NoLock                  = 0
	AccessShareLock         = 1
	RowShareLock            = 2
	RowExclusiveLock        = 3
	ShareUpdateExclusiveLock = 4
	ShareLock               = 5
	ShareRowExclusiveLock   = 6
	ExclusiveLock           = 7
	AccessExclusiveLock     = 8
)

// ClusterOption bitmask values
const (
	CLUOPT_VERBOSE = 1 << 0
)

// FetchDirection for FETCH/MOVE statements.
type FetchDirection int

const (
	FETCH_FORWARD  FetchDirection = iota
	FETCH_BACKWARD
	FETCH_ABSOLUTE
	FETCH_RELATIVE
)

// FETCH_ALL is the special value meaning fetch all rows.
const FETCH_ALL int64 = 0x7FFFFFFFFFFFFFFF // LONG_MAX

// Cursor option bitmask constants.
const (
	CURSOR_OPT_BINARY       = 0x0001
	CURSOR_OPT_SCROLL       = 0x0002
	CURSOR_OPT_NO_SCROLL    = 0x0004
	CURSOR_OPT_INSENSITIVE  = 0x0008
	CURSOR_OPT_ASENSITIVE   = 0x0010
	CURSOR_OPT_HOLD         = 0x0020
	CURSOR_OPT_FAST_PLAN    = 0x0100
	CURSOR_OPT_GENERIC_PLAN = 0x0200
	CURSOR_OPT_CUSTOM_PLAN  = 0x0400
	CURSOR_OPT_PARALLEL_OK  = 0x0800
)

// Trigger type bitmask constants (from trigger.h).
const (
	TRIGGER_TYPE_ROW      = 1 << 0
	TRIGGER_TYPE_BEFORE   = 1 << 1
	TRIGGER_TYPE_INSERT   = 1 << 2
	TRIGGER_TYPE_DELETE   = 1 << 3
	TRIGGER_TYPE_UPDATE   = 1 << 4
	TRIGGER_TYPE_TRUNCATE = 1 << 5
	TRIGGER_TYPE_INSTEAD  = 1 << 6
	TRIGGER_TYPE_AFTER    = 0 // default (not BEFORE, not INSTEAD)
)

// Trigger fire condition constants (from trigger.h).
const (
	TRIGGER_FIRES_ON_ORIGIN  = 'O'
	TRIGGER_FIRES_ALWAYS     = 'A'
	TRIGGER_FIRES_ON_REPLICA = 'R'
	TRIGGER_DISABLED         = 'D'
)

// ConstraintAttributeSpec bit constants (from gram.y).
const (
	CAS_NOT_DEFERRABLE    = 1 << 0
	CAS_DEFERRABLE        = 1 << 1
	CAS_INITIALLY_IMMEDIATE = 1 << 2
	CAS_INITIALLY_DEFERRED  = 1 << 3
	CAS_NOT_VALID         = 1 << 4
	CAS_NO_INHERIT        = 1 << 5
)

// ImportForeignSchemaType represents the type of foreign schema import.
type ImportForeignSchemaType int

const (
	FDW_IMPORT_SCHEMA_ALL      ImportForeignSchemaType = iota
	FDW_IMPORT_SCHEMA_LIMIT_TO
	FDW_IMPORT_SCHEMA_EXCEPT
)

// DefElemAction represents the action for ALTER ... OPTIONS.
type DefElemAction int

const (
	DEFELEM_UNSPEC DefElemAction = iota
	DEFELEM_SET
	DEFELEM_ADD
	DEFELEM_DROP
)

// PublicationObjSpecType represents the type of publication object.
type PublicationObjSpecType int

const (
	PUBLICATIONOBJ_TABLE               PublicationObjSpecType = iota
	PUBLICATIONOBJ_TABLES_IN_SCHEMA
	PUBLICATIONOBJ_TABLES_IN_CUR_SCHEMA
	PUBLICATIONOBJ_CONTINUATION
)

// AlterSubscriptionType represents the kind of ALTER SUBSCRIPTION.
type AlterSubscriptionType int

const (
	ALTER_SUBSCRIPTION_OPTIONS        AlterSubscriptionType = iota
	ALTER_SUBSCRIPTION_CONNECTION
	ALTER_SUBSCRIPTION_SET_PUBLICATION
	ALTER_SUBSCRIPTION_ADD_PUBLICATION
	ALTER_SUBSCRIPTION_DROP_PUBLICATION
	ALTER_SUBSCRIPTION_REFRESH
	ALTER_SUBSCRIPTION_ENABLED
	ALTER_SUBSCRIPTION_SKIP
)

// AlterPublicationAction represents the action for ALTER PUBLICATION.
type AlterPublicationAction int

const (
	AP_AddObjects  AlterPublicationAction = iota
	AP_DropObjects
	AP_SetObjects
)

// AMTYPE constants for access method types.
const (
	AMTYPE_INDEX = 'i'
	AMTYPE_TABLE = 't'
)

// CoercionContext represents cast context.
type CoercionContext int

const (
	COERCION_IMPLICIT   CoercionContext = iota // coercion in context of expression
	COERCION_ASSIGNMENT                        // coercion in context of assignment
	COERCION_PLPGSQL                           // explicit coercion in PL/pgSQL
	COERCION_EXPLICIT                          // explicit cast operation
)

// AlterTSConfigType represents ALTER TEXT SEARCH CONFIGURATION kinds.
type AlterTSConfigType int

const (
	ALTER_TSCONFIG_ADD_MAPPING            AlterTSConfigType = iota
	ALTER_TSCONFIG_ALTER_MAPPING_FOR_TOKEN
	ALTER_TSCONFIG_REPLACE_DICT
	ALTER_TSCONFIG_REPLACE_DICT_FOR_TOKEN
	ALTER_TSCONFIG_DROP_MAPPING
)

// OPCLASS_ITEM_* constants for CreateOpClassItem.
const (
	OPCLASS_ITEM_OPERATOR    = 1
	OPCLASS_ITEM_FUNCTION    = 2
	OPCLASS_ITEM_STORAGETYPE = 3
)

// SVFOp represents SQL-standard function types for SQLValueFunction.
type SVFOp int

const (
	SVFOP_CURRENT_DATE      SVFOp = iota
	SVFOP_CURRENT_TIME
	SVFOP_CURRENT_TIME_N
	SVFOP_CURRENT_TIMESTAMP
	SVFOP_CURRENT_TIMESTAMP_N
	SVFOP_LOCALTIME
	SVFOP_LOCALTIME_N
	SVFOP_LOCALTIMESTAMP
	SVFOP_LOCALTIMESTAMP_N
	SVFOP_CURRENT_ROLE
	SVFOP_CURRENT_USER
	SVFOP_USER
	SVFOP_SESSION_USER
	SVFOP_CURRENT_CATALOG
	SVFOP_CURRENT_SCHEMA
)

// FRAMEOPTION_* constants for WindowDef.FrameOptions (bitmask).
// These match PostgreSQL's FRAMEOPTION_* defines in parsenodes.h.
const (
	FRAMEOPTION_NONDEFAULT                 = 0x00001 // any specified?
	FRAMEOPTION_RANGE                      = 0x00002 // RANGE behavior
	FRAMEOPTION_ROWS                       = 0x00004 // ROWS behavior
	FRAMEOPTION_GROUPS                     = 0x00008 // GROUPS behavior
	FRAMEOPTION_BETWEEN                    = 0x00010 // BETWEEN given?
	FRAMEOPTION_START_UNBOUNDED_PRECEDING  = 0x00020 // start is U. P.
	FRAMEOPTION_END_UNBOUNDED_PRECEDING    = 0x00040 // (disallowed)
	FRAMEOPTION_START_UNBOUNDED_FOLLOWING  = 0x00080 // (disallowed)
	FRAMEOPTION_END_UNBOUNDED_FOLLOWING    = 0x00100 // end is U. F.
	FRAMEOPTION_START_CURRENT_ROW          = 0x00200 // start is C. R.
	FRAMEOPTION_END_CURRENT_ROW            = 0x00400 // end is C. R.
	FRAMEOPTION_START_OFFSET_PRECEDING     = 0x00800 // start is O. P.
	FRAMEOPTION_END_OFFSET_PRECEDING       = 0x01000 // end is O. P.
	FRAMEOPTION_START_OFFSET_FOLLOWING     = 0x02000 // start is O. F.
	FRAMEOPTION_END_OFFSET_FOLLOWING       = 0x04000 // end is O. F.
	FRAMEOPTION_EXCLUDE_CURRENT_ROW        = 0x08000 // omit C.R.
	FRAMEOPTION_EXCLUDE_GROUP              = 0x10000 // omit C.R. & peers
	FRAMEOPTION_EXCLUDE_TIES               = 0x20000 // omit peers only

	FRAMEOPTION_START_OFFSET = FRAMEOPTION_START_OFFSET_PRECEDING | FRAMEOPTION_START_OFFSET_FOLLOWING
	FRAMEOPTION_END_OFFSET   = FRAMEOPTION_END_OFFSET_PRECEDING | FRAMEOPTION_END_OFFSET_FOLLOWING

	FRAMEOPTION_DEFAULTS = FRAMEOPTION_RANGE | FRAMEOPTION_START_UNBOUNDED_PRECEDING | FRAMEOPTION_END_CURRENT_ROW
)

// Interval field codes (from postgres datetime.h).
// Used in INTERVAL type modifiers.
const (
	INTERVAL_MASK_YEAR   = 1 << 2
	INTERVAL_MASK_MONTH  = 1 << 1
	INTERVAL_MASK_DAY    = 1 << 3
	INTERVAL_MASK_HOUR   = 1 << 10
	INTERVAL_MASK_MINUTE = 1 << 11
	INTERVAL_MASK_SECOND = 1 << 12
	INTERVAL_FULL_RANGE  = 0x7FFF
)

// RELPERSISTENCE_* constants.
const (
	RELPERSISTENCE_PERMANENT = 'p'
	RELPERSISTENCE_UNLOGGED  = 'u'
	RELPERSISTENCE_TEMP      = 't'
)
