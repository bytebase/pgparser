package nodes

// NodeTag identifies the type of a node.
// Values must match PostgreSQL's NodeTag enum.
type NodeTag int

// NodeTag values for all node types.
// These must match PostgreSQL's generated nodetags.h.
// We define them manually for the nodes we need.
const (
	T_Invalid NodeTag = 0

	// Value nodes (from value.h)
	T_Integer NodeTag = iota + 1
	T_Float
	T_Boolean
	T_String
	T_BitString

	// List nodes
	T_List
	T_IntList
	T_OidList

	// Primitive nodes (from primnodes.h)
	T_Alias
	T_RangeVar
	T_TableFunc
	T_IntoClause
	T_Var
	T_Const
	T_Param
	T_Aggref
	T_GroupingFunc
	T_WindowFunc
	T_SubscriptingRef
	T_FuncExpr
	T_NamedArgExpr
	T_OpExpr
	T_DistinctExpr
	T_NullIfExpr
	T_ScalarArrayOpExpr
	T_BoolExpr
	T_SubLink
	T_SubPlan
	T_AlternativeSubPlan
	T_FieldSelect
	T_FieldStore
	T_RelabelType
	T_CoerceViaIO
	T_ArrayCoerceExpr
	T_ConvertRowtypeExpr
	T_CollateExpr
	T_CaseExpr
	T_CaseWhen
	T_CaseTestExpr
	T_ArrayExpr
	T_RowExpr
	T_RowCompareExpr
	T_CoalesceExpr
	T_MinMaxExpr
	T_SQLValueFunction
	T_XmlExpr
	T_JsonFormat
	T_JsonReturning
	T_JsonValueExpr
	T_JsonConstructorExpr
	T_JsonIsPredicate
	T_JsonBehavior
	T_JsonExpr
	T_JsonTablePath
	T_JsonTablePathScan
	T_JsonTableSiblingJoin
	T_NullTest
	T_BooleanTest
	T_CoerceToDomain
	T_CoerceToDomainValue
	T_SetToDefault
	T_CurrentOfExpr
	T_NextValueExpr
	T_InferenceElem
	T_TargetEntry
	T_RangeTblRef
	T_JoinExpr
	T_FromExpr
	T_OnConflictExpr
	T_MergeAction

	// Parse nodes (from parsenodes.h) - Raw parse tree nodes
	T_Query
	T_TypeName
	T_ColumnRef
	T_ParamRef
	T_A_Expr
	T_A_Const
	T_TypeCast
	T_CollateClause
	T_RoleSpec
	T_FuncCall
	T_A_Star
	T_A_Indices
	T_A_Indirection
	T_A_ArrayExpr
	T_ResTarget
	T_MultiAssignRef
	T_SortBy
	T_WindowDef
	T_RangeSubselect
	T_RangeFunction
	T_RangeTableFunc
	T_RangeTableFuncCol
	T_RangeTableSample
	T_ColumnDef
	T_TableLikeClause
	T_IndexElem
	T_DefElem
	T_LockingClause
	T_XmlSerialize
	T_PartitionElem
	T_PartitionSpec
	T_PartitionBoundSpec
	T_PartitionRangeDatum
	T_SinglePartitionSpec
	T_PartitionCmd
	T_RangeTblEntry
	T_RTEPermissionInfo
	T_RangeTblFunction
	T_TableSampleClause
	T_WithCheckOption
	T_SortGroupClause
	T_GroupingSet
	T_WindowClause
	T_RowMarkClause
	T_WithClause
	T_InferClause
	T_OnConflictClause
	T_CTESearchClause
	T_CTECycleClause
	T_CommonTableExpr
	T_MergeWhenClause
	T_TriggerTransition
	T_JsonOutput
	T_JsonArgument
	T_JsonFuncExpr
	T_JsonTablePathSpec
	T_JsonTable
	T_JsonTableColumn
	T_JsonKeyValue
	T_JsonParseExpr
	T_JsonScalarExpr
	T_JsonSerializeExpr
	T_JsonObjectConstructor
	T_JsonArrayConstructor
	T_JsonArrayQueryConstructor
	T_JsonAggConstructor
	T_JsonObjectAgg
	T_JsonArrayAgg
	T_RawStmt
	T_InsertStmt
	T_DeleteStmt
	T_UpdateStmt
	T_MergeStmt
	T_SelectStmt
	T_SetOperationStmt
	T_ReturnStmt
	T_PLAssignStmt
	T_CreateSchemaStmt
	T_AlterTableStmt
	T_ReplicaIdentityStmt
	T_AlterTableCmd
	T_AlterCollationStmt
	T_AlterDomainStmt
	T_GrantStmt
	T_ObjectWithArgs
	T_AccessPriv
	T_GrantRoleStmt
	T_AlterDefaultPrivilegesStmt
	T_CopyStmt
	T_VariableSetStmt
	T_VariableShowStmt
	T_CreateStmt
	T_Constraint
	T_CreateTableSpaceStmt
	T_DropTableSpaceStmt
	T_AlterTableSpaceOptionsStmt
	T_AlterTableMoveAllStmt
	T_CreateExtensionStmt
	T_AlterExtensionStmt
	T_AlterExtensionContentsStmt
	T_CreateFdwStmt
	T_AlterFdwStmt
	T_CreateForeignServerStmt
	T_AlterForeignServerStmt
	T_CreateForeignTableStmt
	T_CreateUserMappingStmt
	T_AlterUserMappingStmt
	T_DropUserMappingStmt
	T_ImportForeignSchemaStmt
	T_CreatePolicyStmt
	T_AlterPolicyStmt
	T_CreateAmStmt
	T_CreateTrigStmt
	T_CreateEventTrigStmt
	T_AlterEventTrigStmt
	T_CreatePLangStmt
	T_CreateRoleStmt
	T_AlterRoleStmt
	T_AlterRoleSetStmt
	T_DropRoleStmt
	T_CreateSeqStmt
	T_AlterSeqStmt
	T_DefineStmt
	T_CreateDomainStmt
	T_CreateOpClassStmt
	T_CreateOpClassItem
	T_CreateOpFamilyStmt
	T_AlterOpFamilyStmt
	T_DropStmt
	T_TruncateStmt
	T_CommentStmt
	T_SecLabelStmt
	T_DeclareCursorStmt
	T_ClosePortalStmt
	T_FetchStmt
	T_IndexStmt
	T_CreateStatsStmt
	T_StatsElem
	T_AlterStatsStmt
	T_CreateFunctionStmt
	T_FunctionParameter
	T_AlterFunctionStmt
	T_DoStmt
	T_InlineCodeBlock
	T_CallStmt
	T_CallContext
	T_RenameStmt
	T_AlterObjectDependsStmt
	T_AlterObjectSchemaStmt
	T_AlterOwnerStmt
	T_AlterOperatorStmt
	T_AlterTypeStmt
	T_RuleStmt
	T_NotifyStmt
	T_ListenStmt
	T_UnlistenStmt
	T_TransactionStmt
	T_CompositeTypeStmt
	T_CreateEnumStmt
	T_CreateRangeStmt
	T_AlterEnumStmt
	T_ViewStmt
	T_LoadStmt
	T_CreatedbStmt
	T_AlterDatabaseStmt
	T_AlterDatabaseRefreshCollStmt
	T_AlterDatabaseSetStmt
	T_DropdbStmt
	T_AlterSystemStmt
	T_ClusterStmt
	T_VacuumStmt
	T_VacuumRelation
	T_ExplainStmt
	T_CreateTableAsStmt
	T_RefreshMatViewStmt
	T_CheckPointStmt
	T_DiscardStmt
	T_LockStmt
	T_ConstraintsSetStmt
	T_ReindexStmt
	T_CreateConversionStmt
	T_CreateCastStmt
	T_CreateTransformStmt
	T_PrepareStmt
	T_ExecuteStmt
	T_DeallocateStmt
	T_DropOwnedStmt
	T_ReassignOwnedStmt
	T_AlterTSDictionaryStmt
	T_AlterTSConfigurationStmt
	T_PublicationTable
	T_PublicationObjSpec
	T_CreatePublicationStmt
	T_AlterPublicationStmt
	T_CreateSubscriptionStmt
	T_AlterSubscriptionStmt
	T_DropSubscriptionStmt
)

// NodeTagName returns the string name of a NodeTag.
func NodeTagName(tag NodeTag) string {
	switch tag {
	case T_Invalid:
		return "Invalid"
	case T_Integer:
		return "Integer"
	case T_Float:
		return "Float"
	case T_Boolean:
		return "Boolean"
	case T_String:
		return "String"
	case T_BitString:
		return "BitString"
	case T_List:
		return "List"
	case T_Query:
		return "Query"
	case T_SelectStmt:
		return "SelectStmt"
	case T_InsertStmt:
		return "InsertStmt"
	case T_UpdateStmt:
		return "UpdateStmt"
	case T_DeleteStmt:
		return "DeleteStmt"
	case T_CreateStmt:
		return "CreateStmt"
	case T_ViewStmt:
		return "ViewStmt"
	case T_IndexStmt:
		return "IndexStmt"
	case T_CheckPointStmt:
		return "CheckPointStmt"
	case T_DiscardStmt:
		return "DiscardStmt"
	case T_ListenStmt:
		return "ListenStmt"
	case T_UnlistenStmt:
		return "UnlistenStmt"
	case T_NotifyStmt:
		return "NotifyStmt"
	case T_LoadStmt:
		return "LoadStmt"
	case T_ClosePortalStmt:
		return "ClosePortalStmt"
	case T_ConstraintsSetStmt:
		return "ConstraintsSetStmt"
	case T_VariableSetStmt:
		return "VariableSetStmt"
	case T_VariableShowStmt:
		return "VariableShowStmt"
	case T_PrepareStmt:
		return "PrepareStmt"
	case T_ExecuteStmt:
		return "ExecuteStmt"
	case T_DeallocateStmt:
		return "DeallocateStmt"
	case T_TruncateStmt:
		return "TruncateStmt"
	case T_CommentStmt:
		return "CommentStmt"
	case T_SecLabelStmt:
		return "SecLabelStmt"
	case T_LockStmt:
		return "LockStmt"
	case T_VacuumStmt:
		return "VacuumStmt"
	case T_ClusterStmt:
		return "ClusterStmt"
	case T_ReindexStmt:
		return "ReindexStmt"
	case T_DeclareCursorStmt:
		return "DeclareCursorStmt"
	case T_FetchStmt:
		return "FetchStmt"
	case T_MergeStmt:
		return "MergeStmt"
	case T_MergeWhenClause:
		return "MergeWhenClause"
	case T_CallStmt:
		return "CallStmt"
	case T_DoStmt:
		return "DoStmt"
	case T_CreateRoleStmt:
		return "CreateRoleStmt"
	case T_AlterRoleStmt:
		return "AlterRoleStmt"
	case T_AlterRoleSetStmt:
		return "AlterRoleSetStmt"
	case T_DropRoleStmt:
		return "DropRoleStmt"
	case T_GrantRoleStmt:
		return "GrantRoleStmt"
	case T_CreatedbStmt:
		return "CreatedbStmt"
	case T_AlterDatabaseStmt:
		return "AlterDatabaseStmt"
	case T_AlterDatabaseSetStmt:
		return "AlterDatabaseSetStmt"
	case T_DropdbStmt:
		return "DropdbStmt"
	case T_AlterSystemStmt:
		return "AlterSystemStmt"
	case T_CreateSchemaStmt:
		return "CreateSchemaStmt"
	case T_CreateSeqStmt:
		return "CreateSeqStmt"
	case T_AlterSeqStmt:
		return "AlterSeqStmt"
	case T_CreateDomainStmt:
		return "CreateDomainStmt"
	case T_AlterDomainStmt:
		return "AlterDomainStmt"
	case T_AlterEnumStmt:
		return "AlterEnumStmt"
	case T_AlterCollationStmt:
		return "AlterCollationStmt"
	case T_DefineStmt:
		return "DefineStmt"
	case T_CompositeTypeStmt:
		return "CompositeTypeStmt"
	case T_CreateEnumStmt:
		return "CreateEnumStmt"
	case T_CreateRangeStmt:
		return "CreateRangeStmt"
	case T_AlterFunctionStmt:
		return "AlterFunctionStmt"
	case T_ObjectWithArgs:
		return "ObjectWithArgs"
	case T_CreateTrigStmt:
		return "CreateTrigStmt"
	case T_CreateEventTrigStmt:
		return "CreateEventTrigStmt"
	case T_AlterEventTrigStmt:
		return "AlterEventTrigStmt"
	case T_RuleStmt:
		return "RuleStmt"
	case T_CreatePLangStmt:
		return "CreatePLangStmt"
	case T_DropStmt:
		return "DropStmt"
	case T_CreateFunctionStmt:
		return "CreateFunctionStmt"
	case T_TriggerTransition:
		return "TriggerTransition"
	case T_CreateFdwStmt:
		return "CreateFdwStmt"
	case T_AlterFdwStmt:
		return "AlterFdwStmt"
	case T_CreateForeignServerStmt:
		return "CreateForeignServerStmt"
	case T_AlterForeignServerStmt:
		return "AlterForeignServerStmt"
	case T_CreateForeignTableStmt:
		return "CreateForeignTableStmt"
	case T_CreateUserMappingStmt:
		return "CreateUserMappingStmt"
	case T_AlterUserMappingStmt:
		return "AlterUserMappingStmt"
	case T_DropUserMappingStmt:
		return "DropUserMappingStmt"
	case T_ImportForeignSchemaStmt:
		return "ImportForeignSchemaStmt"
	case T_CreateExtensionStmt:
		return "CreateExtensionStmt"
	case T_AlterExtensionStmt:
		return "AlterExtensionStmt"
	case T_AlterExtensionContentsStmt:
		return "AlterExtensionContentsStmt"
	case T_CreateTableSpaceStmt:
		return "CreateTableSpaceStmt"
	case T_DropTableSpaceStmt:
		return "DropTableSpaceStmt"
	case T_AlterTableSpaceOptionsStmt:
		return "AlterTableSpaceOptionsStmt"
	case T_CreateAmStmt:
		return "CreateAmStmt"
	case T_CreatePolicyStmt:
		return "CreatePolicyStmt"
	case T_AlterPolicyStmt:
		return "AlterPolicyStmt"
	case T_CreatePublicationStmt:
		return "CreatePublicationStmt"
	case T_AlterPublicationStmt:
		return "AlterPublicationStmt"
	case T_PublicationObjSpec:
		return "PublicationObjSpec"
	case T_PublicationTable:
		return "PublicationTable"
	case T_CreateSubscriptionStmt:
		return "CreateSubscriptionStmt"
	case T_AlterSubscriptionStmt:
		return "AlterSubscriptionStmt"
	case T_DropSubscriptionStmt:
		return "DropSubscriptionStmt"
	case T_AlterObjectSchemaStmt:
		return "AlterObjectSchemaStmt"
	case T_AlterOwnerStmt:
		return "AlterOwnerStmt"
	case T_AlterObjectDependsStmt:
		return "AlterObjectDependsStmt"
	case T_AlterOperatorStmt:
		return "AlterOperatorStmt"
	case T_AlterTypeStmt:
		return "AlterTypeStmt"
	case T_AlterDefaultPrivilegesStmt:
		return "AlterDefaultPrivilegesStmt"
	case T_AlterTSConfigurationStmt:
		return "AlterTSConfigurationStmt"
	case T_AlterTSDictionaryStmt:
		return "AlterTSDictionaryStmt"
	case T_CreateStatsStmt:
		return "CreateStatsStmt"
	case T_StatsElem:
		return "StatsElem"
	case T_AlterStatsStmt:
		return "AlterStatsStmt"
	case T_CreateOpClassStmt:
		return "CreateOpClassStmt"
	case T_CreateOpClassItem:
		return "CreateOpClassItem"
	case T_CreateOpFamilyStmt:
		return "CreateOpFamilyStmt"
	case T_AlterOpFamilyStmt:
		return "AlterOpFamilyStmt"
	case T_CreateCastStmt:
		return "CreateCastStmt"
	case T_CreateTransformStmt:
		return "CreateTransformStmt"
	case T_CreateConversionStmt:
		return "CreateConversionStmt"
	case T_DropOwnedStmt:
		return "DropOwnedStmt"
	case T_ReassignOwnedStmt:
		return "ReassignOwnedStmt"
	case T_CreateTableAsStmt:
		return "CreateTableAsStmt"
	case T_RefreshMatViewStmt:
		return "RefreshMatViewStmt"
	// Add more as needed
	default:
		return "Unknown"
	}
}
