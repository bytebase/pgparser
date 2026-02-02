package nodes

// This file contains parse tree node types from PostgreSQL's parsenodes.h.
// These are the "raw" parse tree nodes output by the parser.

// RawStmt wraps a raw (unparsed) statement.
type RawStmt struct {
	Stmt         Node     // raw statement
	StmtLocation ParseLoc // start location, or -1 if unknown
	StmtLen      ParseLoc // length in bytes; 0 means "rest of string"
}

func (n *RawStmt) Tag() NodeTag { return T_RawStmt }

// SelectStmt represents a SELECT statement.
// This is output by the parser for SELECT, VALUES, and table references.
type SelectStmt struct {
	// Fields for "leaf" SelectStmts (simple SELECT or VALUES)
	DistinctClause *List       // NULL, list of DISTINCT ON exprs, or lcons(NIL,NIL) for all (SELECT DISTINCT)
	IntoClause     *IntoClause // target for SELECT INTO
	TargetList     *List       // the target list (of ResTarget)
	FromClause     *List       // the FROM clause
	WhereClause    Node        // WHERE qualification
	GroupClause    *List       // GROUP BY clauses
	GroupDistinct  bool        // Is this GROUP BY DISTINCT?
	HavingClause   Node        // HAVING conditional-expression
	WindowClause   *List       // WINDOW window_name AS (...), ...

	// For VALUES lists
	ValuesLists *List // untransformed list of expression lists

	// Fields for both "leaf" and upper-level SelectStmts
	SortClause    *List       // sort clause (a list of SortBy's)
	LimitOffset   Node        // # of result tuples to skip
	LimitCount    Node        // # of result tuples to return
	LimitOption   LimitOption // limit type
	LockingClause *List       // FOR UPDATE (list of LockingClause's)
	WithClause    *WithClause // WITH clause

	// Fields for upper-level SelectStmts (set operations)
	Op   SetOperation // type of set op
	All  bool         // ALL specified?
	Larg *SelectStmt  // left child
	Rarg *SelectStmt  // right child
}

func (n *SelectStmt) Tag() NodeTag { return T_SelectStmt }

// InsertStmt represents an INSERT statement.
type InsertStmt struct {
	Relation         *RangeVar         // relation to insert into
	Cols             *List             // optional: names of the target columns
	SelectStmt       Node              // the source SELECT/VALUES, or NULL
	OnConflictClause *OnConflictClause // ON CONFLICT clause
	ReturningList    *List             // list of expressions to return
	WithClause       *WithClause       // WITH clause
	Override         OverridingKind    // OVERRIDING clause
}

func (n *InsertStmt) Tag() NodeTag { return T_InsertStmt }

// UpdateStmt represents an UPDATE statement.
type UpdateStmt struct {
	Relation      *RangeVar   // relation to update
	TargetList    *List       // the target list (of ResTarget)
	WhereClause   Node        // qualifications
	FromClause    *List       // optional from clause for more tables
	ReturningList *List       // list of expressions to return
	WithClause    *WithClause // WITH clause
}

func (n *UpdateStmt) Tag() NodeTag { return T_UpdateStmt }

// DeleteStmt represents a DELETE statement.
type DeleteStmt struct {
	Relation      *RangeVar   // relation to delete from
	UsingClause   *List       // optional using clause for more tables
	WhereClause   Node        // qualifications
	ReturningList *List       // list of expressions to return
	WithClause    *WithClause // WITH clause
}

func (n *DeleteStmt) Tag() NodeTag { return T_DeleteStmt }

// CreateStmt represents a CREATE TABLE statement.
type CreateStmt struct {
	Relation       *RangeVar      // relation to create
	TableElts      *List          // column definitions (ColumnDef)
	InhRelations   *List          // relations to inherit from (RangeVar)
	Partbound      Node           // FOR VALUES clause
	Partspec       *PartitionSpec // PARTITION BY clause
	OfTypename     *TypeName      // OF typename
	Constraints    *List          // constraints (Constraint)
	Options        *List          // options from WITH clause
	OnCommit       OnCommitAction // what to do at commit time
	Tablespacename string         // table space to use, or NULL
	AccessMethod   string         // table access method
	IfNotExists    bool           // just do nothing if it already exists?
}

func (n *CreateStmt) Tag() NodeTag { return T_CreateStmt }

// ViewStmt represents a CREATE VIEW statement.
type ViewStmt struct {
	View            *RangeVar // the view to be created
	Aliases         *List     // target column names
	Query           Node      // the SELECT query (as a raw parse tree)
	Replace         bool      // replace an existing view?
	Options         *List     // options from WITH clause
	WithCheckOption int       // WITH CHECK OPTION
}

func (n *ViewStmt) Tag() NodeTag { return T_ViewStmt }

// IndexStmt represents a CREATE INDEX statement.
type IndexStmt struct {
	Idxname        string    // name of new index, or NULL for default
	Relation       *RangeVar // relation to build index on
	AccessMethod   string    // name of access method (eg. btree)
	TableSpace     string    // tablespace, or NULL for default
	IndexParams    *List     // columns to index: a list of IndexElem
	IndexIncludingParams *List     // additional columns to index: a list of IndexElem
	Options        *List     // WITH clause options
	WhereClause    Node      // qualification (partial-index predicate)
	ExcludeOpNames *List     // exclusion operator names, or NIL if none
	Idxcomment     string    // comment to apply to index, or NULL
	IndexOid       Oid       // OID of an existing index, if any
	OldNumber      uint32    // relfilenumber of existing storage, if any
	OldCreateSubid uint32    // rd_createSubid of existing index
	OldFirstRelfilelocatorSubid uint32 // rd_firstRelfilelocatorSubid of existing index
	Unique         bool      // is index unique?
	Nulls_not_distinct bool  // null treatment for UNIQUE constraints
	Primary        bool      // is index a primary key?
	Isconstraint   bool      // is it for a pkey/unique constraint?
	Deferrable     bool      // is the constraint DEFERRABLE?
	Initdeferred   bool      // is the constraint INITIALLY DEFERRED?
	Transformed    bool      // true when transformIndexStmt is finished
	Concurrent     bool      // should this be a concurrent index build?
	IfNotExists    bool      // just do nothing if index already exists?
	ResetDefaultTblspc bool  // reset default_tablespace prior to executing
}

func (n *IndexStmt) Tag() NodeTag { return T_IndexStmt }

// DropStmt represents a DROP statement.
type DropStmt struct {
	Objects    *List  // list of names
	RemoveType int    // object type (ObjectType)
	Behavior   int    // RESTRICT or CASCADE behavior (DropBehavior)
	Missing_ok bool   // skip error if object is missing?
	Concurrent bool   // drop index concurrently?
}

func (n *DropStmt) Tag() NodeTag { return T_DropStmt }

// AlterTableStmt represents an ALTER TABLE statement.
type AlterTableStmt struct {
	Relation   *RangeVar // table to work on
	Cmds       *List     // list of subcommands (AlterTableCmd)
	ObjType    int       // type of object (ObjectType)
	Missing_ok bool      // skip error if table missing
}

func (n *AlterTableStmt) Tag() NodeTag { return T_AlterTableStmt }

// AlterTableCmd represents a subcommand of ALTER TABLE.
type AlterTableCmd struct {
	Subtype   int    // Type of table alteration to apply
	Name      string // column, constraint, or trigger to act on
	Num       int16  // attribute number for columns referenced by number
	Newowner  *RoleSpec
	Def       Node   // definition of new column, index, constraint, etc.
	Behavior  int    // RESTRICT or CASCADE for DROP cases
	Missing_ok bool
}

func (n *AlterTableCmd) Tag() NodeTag { return T_AlterTableCmd }

// AlterTableMoveAllStmt represents ALTER TABLE/INDEX/MATERIALIZED VIEW ALL IN TABLESPACE.
type AlterTableMoveAllStmt struct {
	OrigTablespacename string // source tablespace
	ObjType            int    // Object type to move
	Roles              *List  // List of roles to move objects of
	NewTablespacename  string // target tablespace
	Nowait             bool
}

func (n *AlterTableMoveAllStmt) Tag() NodeTag { return T_AlterTableMoveAllStmt }

// CreateSchemaStmt represents a CREATE SCHEMA statement.
type CreateSchemaStmt struct {
	Schemaname  string    // the name of the schema to create
	Authrole    *RoleSpec // the owner of the created schema
	SchemaElts  *List     // schema components (list of parsetrees)
	IfNotExists bool      // just do nothing if schema already exists?
}

func (n *CreateSchemaStmt) Tag() NodeTag { return T_CreateSchemaStmt }

// RangeVar represents a range variable (table/view reference).
type RangeVar struct {
	Catalogname string   // the catalog (database) name, or NULL
	Schemaname  string   // the schema name, or NULL
	Relname     string   // the relation/sequence name
	Inh         bool     // expand rel by inheritance? recursively act on children?
	Relpersistence byte  // see RELPERSISTENCE_* in pg_class.h
	Alias       *Alias   // table alias & optional column aliases
	Location    ParseLoc // token location, or -1 if unknown
}

func (n *RangeVar) Tag() NodeTag { return T_RangeVar }

// Alias represents an alias (AS clause).
type Alias struct {
	Aliasname string // aliased rel name
	Colnames  *List  // optional list of column aliases
}

func (n *Alias) Tag() NodeTag { return T_Alias }

// IntoClause represents SELECT INTO clause.
type IntoClause struct {
	Rel            *RangeVar      // target relation name
	ColNames       *List          // column names to assign, or NIL
	AccessMethod   string         // table access method
	Options        *List          // options from WITH clause
	OnCommit       OnCommitAction // what to do at commit time
	TableSpaceName string         // table space to use, or NULL
	ViewQuery      Node           // materialized view's SELECT query
	SkipData       bool           // true for WITH NO DATA
}

func (n *IntoClause) Tag() NodeTag { return T_IntoClause }

// ColumnRef represents a column reference (foo.bar.baz).
type ColumnRef struct {
	Fields   *List    // field names (String nodes) or A_Star
	Location ParseLoc // token location, or -1 if unknown
}

func (n *ColumnRef) Tag() NodeTag { return T_ColumnRef }

// ResTarget represents a result target in SELECT's target list, or
// a column name in INSERT/UPDATE.
type ResTarget struct {
	Name        string   // column name or NULL
	Indirection *List    // subscripts, field names, and '*'
	Val         Node     // the value expression to compute or assign
	Location    ParseLoc // token location, or -1 if unknown
}

func (n *ResTarget) Tag() NodeTag { return T_ResTarget }

// MultiAssignRef is an element of a row source expression for UPDATE SET (a,b) = expr.
type MultiAssignRef struct {
	Source   Node // the row-valued expression
	Colno    int  // column number for this target (1..n)
	Ncolumns int  // number of targets in the construct
}

func (n *MultiAssignRef) Tag() NodeTag { return T_MultiAssignRef }

// A_Expr represents an expression with an operator.
type A_Expr struct {
	Kind     A_Expr_Kind // see above
	Name     *List       // possibly-qualified name of operator
	Lexpr    Node        // left argument, or NULL if none
	Rexpr    Node        // right argument, or NULL if none
	Location ParseLoc    // token location, or -1 if unknown
}

func (n *A_Expr) Tag() NodeTag { return T_A_Expr }

// A_Const represents a constant value.
type A_Const struct {
	Isnull   bool     // if true, value is NULL
	Val      Node     // value (Integer, Float, Boolean, String, or BitString node; or NULL if isnull)
	Location ParseLoc // token location, or -1 if unknown
}

func (n *A_Const) Tag() NodeTag { return T_A_Const }

// TypeCast represents a CAST expression.
type TypeCast struct {
	Arg      Node      // the expression being casted
	TypeName *TypeName // the target type
	Location ParseLoc  // token location, or -1 if unknown
}

func (n *TypeCast) Tag() NodeTag { return T_TypeCast }

// FuncCall represents a function call.
type FuncCall struct {
	Funcname       *List    // qualified name of function
	Args           *List    // the arguments (list of expressions)
	AggOrder       *List    // ORDER BY (list of SortBy)
	AggFilter      Node     // FILTER clause, if any
	Over           Node     // OVER clause, if any (WindowDef)
	AggWithinGroup bool     // ORDER BY appeared in WITHIN GROUP
	AggStar        bool     // argument was really '*'
	AggDistinct    bool     // arguments were labeled DISTINCT
	FuncVariadic   bool     // last argument was labeled VARIADIC
	FuncFormat     int      // how to display function call (CoercionForm)
	Location       ParseLoc // token location, or -1 if unknown
}

func (n *FuncCall) Tag() NodeTag { return T_FuncCall }

// NamedArgExpr represents a named argument in a function call.
type NamedArgExpr struct {
	Arg       Node     // the argument expression
	Name      string   // the name
	Argnumber int      // argument's number in positional notation
	Location  ParseLoc // argument name location, or -1 if unknown
}

func (n *NamedArgExpr) Tag() NodeTag { return T_NamedArgExpr }

// TypeName represents a data type name.
type TypeName struct {
	Names       *List    // qualified name (list of String nodes)
	TypeOid     Oid      // type identified by OID (InvalidOid if not known)
	Setof       bool     // is a set?
	PctType     bool     // %TYPE specified?
	Typmods     *List    // type modifier expression(s)
	Typemod     int32    // prespecified type modifier
	ArrayBounds *List    // array bounds
	Location    ParseLoc // token location, or -1 if unknown
}

func (n *TypeName) Tag() NodeTag { return T_TypeName }

// ColumnDef represents a column definition in CREATE TABLE.
type ColumnDef struct {
	Colname       string    // name of column
	TypeName      *TypeName // type of column
	Compression   string    // compression method for column
	Inhcount      int       // number of times column is inherited
	IsLocal       bool      // column has local (non-inherited) def'n
	IsNotNull     bool      // NOT NULL constraint
	IsFromType    bool      // column definition came from table type
	Storage       byte      // attstorage setting, or 0 for default
	StorageName   string    // storage directive name or NULL
	RawDefault    Node      // default value (untransformed parse tree)
	CookedDefault Node      // default value (transformed expr tree)
	Identity      byte      // attidentity setting
	IdentitySequence *RangeVar // to store identity sequence name
	Generated     byte      // attgenerated setting
	CollClause    *CollateClause // column collation clause
	CollOid       Oid       // collation OID
	Constraints   *List     // other constraints on column
	Fdwoptions    *List     // per-column FDW options
	Location      ParseLoc  // parse location, or -1 if none/unknown
}

func (n *ColumnDef) Tag() NodeTag { return T_ColumnDef }

// Constraint represents a constraint definition in CREATE TABLE.
type Constraint struct {
	Contype         ConstrType // constraint type (see above)
	Conname         string     // constraint name, or NULL if unnamed
	Deferrable      bool       // DEFERRABLE?
	Initdeferred    bool       // INITIALLY DEFERRED?
	Location        ParseLoc   // token location, or -1 if unknown
	IsNoInherit     bool       // NO INHERIT?
	RawExpr         Node       // CHECK expression (raw parse tree)
	CookedExpr      string     // CHECK expression (cooked)
	GeneratedWhen   byte       // ALWAYS or BY DEFAULT
	NullsNotDistinct bool      // UNIQUE nulls distinct?
	Keys            *List      // PRIMARY KEY/UNIQUE column names
	Including       *List      // PRIMARY KEY/UNIQUE INCLUDE column names
	Exclusions      *List      // exclusion constraint
	Options         *List      // WITH clause options
	Indexname       string     // existing index to use; else NULL
	Indexspace      string     // index tablespace; NULL for default
	ResetDefaultTblspc bool    // reset default_tablespace prior to creating the index
	AccessMethod    string     // index access method; NULL for default
	WhereClause     Node       // WHERE for partial index
	Pktable         *RangeVar  // the table the constraint references
	FkAttrs         *List      // FOREIGN KEY column names
	PkAttrs         *List      // PRIMARY KEY column names
	FkMatchtype     byte       // FULL, PARTIAL, SIMPLE
	FkUpdaction     byte       // ON UPDATE action
	FkDelaction     byte       // ON DELETE action
	FkDelsetcols    *List      // ON DELETE SET column names
	OldConpfeqop    *List      // pg_constraint.conpfeqop of old constraint
	OldPktableOid   Oid        // pg_constraint.confrelid of old constraint
	SkipValidation  bool       // skip validation of existing rows?
	InitiallyValid  bool       // mark the new constraint as valid?
}

func (n *Constraint) Tag() NodeTag { return T_Constraint }

// SortBy represents ORDER BY clause item.
type SortBy struct {
	Node        Node        // expression to sort on
	SortbyDir   SortByDir   // ASC/DESC/USING/default
	SortbyNulls SortByNulls // NULLS FIRST/LAST
	UseOp       *List       // name of op to use, if SORTBY_USING
	Location    ParseLoc    // operator location, or -1 if none/unknown
}

func (n *SortBy) Tag() NodeTag { return T_SortBy }

// WithClause represents WITH clause (common table expressions).
type WithClause struct {
	Ctes      *List // list of CommonTableExprs
	Recursive bool  // true = WITH RECURSIVE
	Location  ParseLoc
}

func (n *WithClause) Tag() NodeTag { return T_WithClause }

// CommonTableExpr represents a single CTE in a WITH clause.
type CommonTableExpr struct {
	Ctename          string   // CTE name
	Aliascolnames    *List    // optional column name list
	Ctematerialized  int      // CTEMaterialize enum
	Ctequery         Node     // the CTE's subquery
	SearchClause     Node     // SEARCH clause
	CycleClause      Node     // CYCLE clause
	Location         ParseLoc // token location, or -1 if unknown
	Cterecursive     bool     // is this CTE actually recursive?
	Cterefcount      int      // number of RTEs referencing this CTE
	Ctecolnames      *List    // list of output column names
	Ctecoltypes      *List    // OID list of output column type OIDs
	Ctecoltypmods    *List    // integer list of output column typmods
	Ctecolcollations *List    // OID list of column collation OIDs
}

func (n *CommonTableExpr) Tag() NodeTag { return T_CommonTableExpr }

// CTESearchClause represents the SEARCH clause in a recursive CTE.
type CTESearchClause struct {
	SearchColList    *List    // list of column names to search by
	SearchBreadthFirst bool  // true = BREADTH FIRST, false = DEPTH FIRST
	SearchSeqColumn  string  // name of the output ordering column
	Location         ParseLoc
}

func (n *CTESearchClause) Tag() NodeTag { return T_CTESearchClause }

// CTECycleClause represents the CYCLE clause in a recursive CTE.
type CTECycleClause struct {
	CycleColList     *List    // list of column names to check for cycles
	CycleMarkColumn  string   // name of the cycle mark column
	CycleMarkValue   Node     // value for cycle mark (default TRUE)
	CycleMarkDefault Node     // default for cycle mark (default FALSE)
	CyclePathColumn  string   // name of the cycle path column
	CycleMarkType    Oid      // type of the cycle mark column
	CycleMarkTypmod  int32
	CycleMarkCollation Oid
	CycleMarkNeop    Oid
	Location         ParseLoc
}

func (n *CTECycleClause) Tag() NodeTag { return T_CTECycleClause }

// RoleSpec represents a role specification.
type RoleSpec struct {
	Roletype int      // type of role (RoleSpecType)
	Rolename string   // filled only for ROLESPEC_CSTRING
	Location ParseLoc // token location, or -1 if unknown
}

func (n *RoleSpec) Tag() NodeTag { return T_RoleSpec }

// CollateClause represents COLLATE clause.
type CollateClause struct {
	Arg      Node     // input expression
	Collname *List    // possibly-qualified collation name
	Location ParseLoc // token location, or -1 if unknown
}

func (n *CollateClause) Tag() NodeTag { return T_CollateClause }

// PartitionSpec represents PARTITION BY clause.
type PartitionSpec struct {
	Strategy   string   // partitioning strategy (PARTITION_STRATEGY_*)
	PartParams *List    // partition key list
	Location   ParseLoc // token location, or -1 if unknown
}

func (n *PartitionSpec) Tag() NodeTag { return T_PartitionSpec }

// PartitionElem represents a single partition key element.
type PartitionElem struct {
	Name      string   // name of column to partition on, or ""
	Expr      Node     // expression to partition on, or nil
	Collation *List    // name of collation; nil = default
	Opclass   *List    // name of desired opclass; nil = default
	Location  ParseLoc // token location, or -1 if unknown
}

func (n *PartitionElem) Tag() NodeTag { return T_PartitionElem }

// PartitionBoundSpec represents a partition bound specification.
type PartitionBoundSpec struct {
	Strategy    byte  // PARTITION_STRATEGY_* code
	IsDefault   bool  // is it a default partition bound?
	Modulus     int   // hash partition modulus
	Remainder   int   // hash partition remainder
	Listdatums  *List // list of Consts (or Exprs) for LIST
	Lowerdatums *List // list of Consts (or Exprs) for RANGE lower
	Upperdatums *List // list of Consts (or Exprs) for RANGE upper
	Location    ParseLoc
}

func (n *PartitionBoundSpec) Tag() NodeTag { return T_PartitionBoundSpec }

// PartitionCmd represents ALTER TABLE ATTACH/DETACH PARTITION.
type PartitionCmd struct {
	Name       *RangeVar           // partition to attach/detach
	Bound      *PartitionBoundSpec // FOR VALUES, if attaching
	Concurrent bool
}

func (n *PartitionCmd) Tag() NodeTag { return T_PartitionCmd }

// OnConflictClause represents ON CONFLICT clause.
type OnConflictClause struct {
	Action      int      // DO NOTHING or DO UPDATE
	Infer       *InferClause
	TargetList  *List    // SET clause for DO UPDATE
	WhereClause Node     // WHERE clause for DO UPDATE
	Location    ParseLoc // token location, or -1 if unknown
}

func (n *OnConflictClause) Tag() NodeTag { return T_OnConflictClause }

// InferClause represents ON CONFLICT index inference clause.
type InferClause struct {
	IndexElems  *List    // IndexElems to infer unique index
	WhereClause Node     // qualification (partial-index predicate)
	Conname     string   // constraint name
	Location    ParseLoc // token location, or -1 if unknown
}

func (n *InferClause) Tag() NodeTag { return T_InferClause }

// DefElem represents a generic definition element.
type DefElem struct {
	Defnamespace string   // namespace (NULL if none)
	Defname      string   // option name
	Arg          Node     // option value (can be integer, string, TypeName, etc)
	Defaction    int      // unspecified action, or SET/ADD/DROP (DefElemAction)
	Location     ParseLoc // token location, or -1 if unknown
}

func (n *DefElem) Tag() NodeTag { return T_DefElem }

// LockingClause represents FOR UPDATE/SHARE clause.
type LockingClause struct {
	LockedRels *List // FOR UPDATE/SHARE relations
	Strength   int   // LockClauseStrength
	WaitPolicy int   // LockWaitPolicy
}

func (n *LockingClause) Tag() NodeTag { return T_LockingClause }

// A_Star represents '*' appearing in expression.
type A_Star struct{}

func (n *A_Star) Tag() NodeTag { return T_A_Star }

// A_Indices represents array subscript or slice.
type A_Indices struct {
	IsSlice bool // true if slice (i.e., [ : ])
	Lidx    Node // slice lower bound, if any
	Uidx    Node // subscript, or slice upper bound if any
}

func (n *A_Indices) Tag() NodeTag { return T_A_Indices }

// A_Indirection represents field selection or array subscripting.
type A_Indirection struct {
	Arg         Node  // the thing being selected from
	Indirection *List // subscripts and/or field names and/or *
}

func (n *A_Indirection) Tag() NodeTag { return T_A_Indirection }

// WindowDef represents WINDOW clause definition.
type WindowDef struct {
	Name            string   // window name (NULL in OVER clause)
	Refname         string   // referenced window name, if any
	PartitionClause *List    // PARTITION BY expressions
	OrderClause     *List    // ORDER BY (SortBy)
	FrameOptions    int      // frame_clause options, see WindowDef comments
	StartOffset     Node     // expression for starting bound, if any
	EndOffset       Node     // expression for ending bound, if any
	Location        ParseLoc // parse location, or -1 if none/unknown
}

func (n *WindowDef) Tag() NodeTag { return T_WindowDef }

// JoinExpr represents a JOIN expression.
type JoinExpr struct {
	Jointype    JoinType // type of join
	IsNatural   bool     // Natural join? Will need to shape table
	Larg        Node     // left subtree
	Rarg        Node     // right subtree
	UsingClause *List    // USING clause, if any
	JoinUsing   *Alias   // alias for USING join, if any
	Quals       Node     // qualifications on join, if any
	Alias       *Alias   // user-written alias clause, if any
	Rtindex     int      // RT index assigned for join, or 0
}

func (n *JoinExpr) Tag() NodeTag { return T_JoinExpr }

// FromExpr represents a FROM clause.
type FromExpr struct {
	Fromlist *List // List of join subtrees
	Quals    Node  // qualifiers on join, if any
}

func (n *FromExpr) Tag() NodeTag { return T_FromExpr }

// IndexElem represents a column reference in an index definition.
type IndexElem struct {
	Name          string      // name of attribute to index, or NULL
	Expr          Node        // expression to index, or NULL
	Indexcolname  string      // name for index column; NULL = default
	Collation     *List       // collation for index
	Opclass       *List       // operator class, or NIL
	Opclassopts   *List       // operator class parameters
	Ordering      SortByDir   // ASC/DESC/default
	NullsOrdering SortByNulls // FIRST/LAST/default
}

func (n *IndexElem) Tag() NodeTag { return T_IndexElem }

// ParamRef represents $n parameter reference.
type ParamRef struct {
	Number   int      // number of the parameter
	Location ParseLoc // token location, or -1 if unknown
}

func (n *ParamRef) Tag() NodeTag { return T_ParamRef }

// CurrentOfExpr represents WHERE CURRENT OF cursor_name.
type CurrentOfExpr struct {
	CvarNo     int    // RT index of target relation
	CursorName string // name of referenced cursor
	CursorParam int   // refcursor parameter number
}

func (n *CurrentOfExpr) Tag() NodeTag { return T_CurrentOfExpr }

// SubLink represents a subquery appearing in an expression.
type SubLink struct {
	SubLinkType int      // see SubLinkType above
	SubLinkId   int      // ID (1..n); 0 if not MULTIEXPR
	Testexpr    Node     // outer-query test for ANY/ALL/ROWCOMPARE
	OperName    *List    // originally specified operator name
	Subselect   Node     // subselect as Query* or raw parsetree
	Location    ParseLoc // token location, or -1 if unknown
}

func (n *SubLink) Tag() NodeTag { return T_SubLink }

// BoolExpr represents AND/OR/NOT expression.
type BoolExpr struct {
	Boolop   BoolExprType // AND/OR/NOT
	Args     *List        // arguments to this expression
	Location ParseLoc     // token location, or -1 if unknown
}

func (n *BoolExpr) Tag() NodeTag { return T_BoolExpr }

// NullTestType represents NULL test types.
type NullTestType int

const (
	IS_NULL NullTestType = iota
	IS_NOT_NULL
)

// NullTest represents IS [NOT] NULL test.
type NullTest struct {
	Arg          Node         // input expression
	Nulltesttype NullTestType // IS NULL, IS NOT NULL
	Argisrow     bool         // T to perform field-by-field null checks
	Location     ParseLoc     // token location, or -1 if unknown
}

func (n *NullTest) Tag() NodeTag { return T_NullTest }

// BoolTestType represents boolean test types.
type BoolTestType int

const (
	IS_TRUE BoolTestType = iota
	IS_NOT_TRUE
	IS_FALSE
	IS_NOT_FALSE
	IS_UNKNOWN
	IS_NOT_UNKNOWN
)

// BooleanTest represents IS [NOT] TRUE/FALSE/UNKNOWN test.
type BooleanTest struct {
	Arg          Node         // input expression
	Booltesttype BoolTestType // test type
	Location     ParseLoc     // token location, or -1 if unknown
}

func (n *BooleanTest) Tag() NodeTag { return T_BooleanTest }

// RangeSubselect represents a subquery appearing in a FROM clause.
type RangeSubselect struct {
	Lateral  bool   // does it have LATERAL prefix?
	Subquery Node   // the untransformed sub-select clause
	Alias    *Alias // table alias & optional column aliases
}

func (n *RangeSubselect) Tag() NodeTag { return T_RangeSubselect }

// RangeFunction represents a function appearing in FROM clause.
type RangeFunction struct {
	Lateral    bool   // does it have LATERAL prefix?
	Ordinality bool   // does it have WITH ORDINALITY suffix?
	IsRowsfrom bool   // is this a ROWS FROM() clause?
	Functions  *List  // list of RangeFunction items
	Alias      *Alias // table alias & optional column aliases
	Coldeflist *List  // list of ColumnDef nodes for ROWS FROM()
}

func (n *RangeFunction) Tag() NodeTag { return T_RangeFunction }

// RangeTableSample represents TABLESAMPLE appearing in FROM clause.
type RangeTableSample struct {
	Relation   Node     // relation to sample
	Method     *List    // sampling method name (possibly schema qualified)
	Args       *List    // argument(s) for sampling method
	Repeatable Node     // REPEATABLE expression, or NULL if none
	Location   ParseLoc // method name location, or -1 if unknown
}

func (n *RangeTableSample) Tag() NodeTag { return T_RangeTableSample }

// TableLikeClause represents LIKE clause in CREATE TABLE.
type TableLikeClause struct {
	Relation      *RangeVar // relation to clone
	Options       uint32    // OR of TableLikeOption flags
	RelationOid   Oid       // set during parse analysis to the OID of the relation
	Columns       *List     // list of ColumnDef nodes
	AncillaryData *List     // list of DefElem nodes for INDEX etc
}

func (n *TableLikeClause) Tag() NodeTag { return T_TableLikeClause }

// CaseExpr represents a CASE expression.
type CaseExpr struct {
	Casetype   Oid      // type of expression result
	Casecollid Oid      // OID of collation, or InvalidOid if none
	Arg        Node     // implicit equality comparison argument
	Args       *List    // the arguments (list of CaseWhen)
	Defresult  Node     // the default result (ELSE clause)
	Location   ParseLoc // token location, or -1 if unknown
}

func (n *CaseExpr) Tag() NodeTag { return T_CaseExpr }

// CaseWhen represents a WHEN clause in a CASE expression.
type CaseWhen struct {
	Expr     Node     // condition expression
	Result   Node     // substitution result
	Location ParseLoc // token location, or -1 if unknown
}

func (n *CaseWhen) Tag() NodeTag { return T_CaseWhen }

// CoalesceExpr represents a COALESCE expression.
type CoalesceExpr struct {
	Coalescetype   Oid      // type of expression result
	Coalescecollid Oid      // OID of collation, or InvalidOid if none
	Args           *List    // the arguments
	Location       ParseLoc // token location, or -1 if unknown
}

func (n *CoalesceExpr) Tag() NodeTag { return T_CoalesceExpr }

// MinMaxExpr represents a GREATEST or LEAST expression.
type MinMaxExpr struct {
	Minmaxtype   Oid          // common type of arguments and result
	Minmaxcollid Oid          // OID of collation of result
	Op           MinMaxOp     // GREATEST or LEAST
	Args         *List        // the arguments
	Location     ParseLoc     // token location, or -1 if unknown
}

// MinMaxOp represents GREATEST vs LEAST.
type MinMaxOp int

const (
	IS_GREATEST MinMaxOp = iota
	IS_LEAST
)

func (n *MinMaxExpr) Tag() NodeTag { return T_MinMaxExpr }

// NullIfExpr represents a NULLIF expression.
// This is represented as an OpExpr in the parse tree.
type NullIfExpr struct {
	Opno         Oid      // PG_OPERATOR OID of the operator
	Opfuncid     Oid      // PG_PROC OID of underlying function
	Opresulttype Oid      // PG_TYPE OID of result value
	Opretset     bool     // true if operator returns set
	Opcollid     Oid      // OID of collation of result
	Inputcollid  Oid      // OID of collation that operator should use
	Args         *List    // arguments to the operator (min 2)
	Location     ParseLoc // token location, or -1 if unknown
}

func (n *NullIfExpr) Tag() NodeTag { return T_NullIfExpr }

// RowExpr represents a ROW() or (a, b, c) expression.
type RowExpr struct {
	Args       *List       // the fields
	RowTypeid  Oid         // RECORDOID or a composite type's ID
	RowFormat  CoercionForm // how to display this node
	Colnames   *List       // list of String, or NIL
	Location   ParseLoc    // token location, or -1 if unknown
}

func (n *RowExpr) Tag() NodeTag { return T_RowExpr }

// ArrayExpr represents an ARRAY[] construct.
type ArrayExpr struct {
	ArrayTypeid  Oid      // type of expression result
	ArrayCollid  Oid      // OID of collation, or InvalidOid if none
	ElementTypeid Oid     // common type of array elements
	Elements     *List    // list of Array elements
	Multidims    bool     // true if elements are sub-arrays
	Location     ParseLoc // token location, or -1 if unknown
}

func (n *ArrayExpr) Tag() NodeTag { return T_ArrayExpr }

// A_ArrayExpr represents an ARRAY[] construct in raw parse tree.
type A_ArrayExpr struct {
	Elements *List    // array element expressions
	Location ParseLoc // token location, or -1 if unknown
}

func (n *A_ArrayExpr) Tag() NodeTag { return T_A_ArrayExpr }

// GroupingFunc represents a GROUPING(...) expression.
type GroupingFunc struct {
	Args       *List    // arguments, not evaluated but kept for benefit of EXPLAIN etc.
	Refs       *List    // ressortgrouprefs of arguments
	Agglevelsup uint32  // same as Aggref.agglevelsup
	Location   ParseLoc // token location, or -1 if unknown
}

func (n *GroupingFunc) Tag() NodeTag { return T_GroupingFunc }

// GroupingSet represents a CUBE, ROLLUP, or GROUPING SETS clause.
type GroupingSet struct {
	Kind     GroupingSetKind // GROUPING SETS, CUBE, ROLLUP
	Content  *List           // content of the set
	Location ParseLoc        // token location, or -1 if unknown
}

// GroupingSetKind represents the kind of grouping set.
type GroupingSetKind int

const (
	GROUPING_SET_EMPTY GroupingSetKind = iota
	GROUPING_SET_SIMPLE
	GROUPING_SET_ROLLUP
	GROUPING_SET_CUBE
	GROUPING_SET_SETS
)

func (n *GroupingSet) Tag() NodeTag { return T_GroupingSet }

// WindowClause represents a WINDOW clause entry.
type WindowClause struct {
	Name            string   // window name (NULL if none)
	Refname         string   // referenced window name (NULL if none)
	PartitionClause *List    // PARTITION BY list
	OrderClause     *List    // ORDER BY list
	FrameOptions    int      // frame_clause options, see WindowDef
	StartOffset     Node     // expression for starting bound, if any
	EndOffset       Node     // expression for ending bound, if any
	RunCondition    *List    // qual to help short-circuit execution
	StartInRangeFunc Oid     // in_range function for start bound
	EndInRangeFunc   Oid     // in_range function for end bound
	InRangeColl      Oid     // collation for in_range comparisons
	InRangeAsc       bool    // use ASC sort order for in_range?
	InRangeNullsFirst bool   // nulls sort first for in_range?
	Winref           uint32  // ID referenced by window functions
	Copiedorder      bool    // did we copy orderClause from refname?
}

func (n *WindowClause) Tag() NodeTag { return T_WindowClause }

// MergeStmt represents a MERGE statement.
type MergeStmt struct {
	Relation         *RangeVar   // target relation to merge into
	SourceRelation   Node        // source relation
	JoinCondition    Node        // join condition between source and target
	MergeWhenClauses *List       // list of MergeWhenClause
	ReturningList    *List       // list of expressions to return
	WithClause       *WithClause // WITH clause
}

func (n *MergeStmt) Tag() NodeTag { return T_MergeStmt }

// MergeWhenClause represents a WHEN clause in MERGE statement.
type MergeWhenClause struct {
	Kind        MergeMatchKind // MATCHED, NOT MATCHED BY SOURCE, NOT MATCHED BY TARGET
	Condition   Node           // condition expression
	TargetList  *List          // target list for INSERT/UPDATE
	Values      *List          // VALUES for INSERT (NULL for UPDATE/DELETE)
	Override    OverridingKind // OVERRIDING clause for INSERT
	CommandType CmdType        // INSERT/UPDATE/DELETE
}

// MergeMatchKind represents the match kind for MERGE.
type MergeMatchKind int

const (
	MERGE_WHEN_MATCHED MergeMatchKind = iota
	MERGE_WHEN_NOT_MATCHED_BY_SOURCE
	MERGE_WHEN_NOT_MATCHED_BY_TARGET
)

func (n *MergeWhenClause) Tag() NodeTag { return T_MergeWhenClause }

// TruncateStmt represents a TRUNCATE statement.
type TruncateStmt struct {
	Relations   *List          // list of relation names to truncate
	RestartSeqs bool           // restart owned sequences?
	Behavior    DropBehavior   // RESTRICT or CASCADE behavior
}

func (n *TruncateStmt) Tag() NodeTag { return T_TruncateStmt }

// CommentStmt represents a COMMENT statement.
type CommentStmt struct {
	Objtype ObjectType // object kind
	Object  Node       // qualified name of object
	Comment string     // comment to insert, or NULL to drop
}

func (n *CommentStmt) Tag() NodeTag { return T_CommentStmt }

// CreateSeqStmt represents a CREATE SEQUENCE statement.
type CreateSeqStmt struct {
	Sequence    *RangeVar // sequence to create
	Options     *List     // list of DefElem
	OwnerId     Oid       // owner's OID (if specified)
	ForIdentity bool      // for GENERATED ... AS IDENTITY
	IfNotExists bool      // just do nothing if it already exists?
}

func (n *CreateSeqStmt) Tag() NodeTag { return T_CreateSeqStmt }

// AlterSeqStmt represents an ALTER SEQUENCE statement.
type AlterSeqStmt struct {
	Sequence    *RangeVar // sequence to alter
	Options     *List     // list of DefElem
	ForIdentity bool      // for GENERATED ... AS IDENTITY
	MissingOk   bool      // skip if sequence doesn't exist?
}

func (n *AlterSeqStmt) Tag() NodeTag { return T_AlterSeqStmt }

// CreateFunctionStmt represents a CREATE FUNCTION statement.
type CreateFunctionStmt struct {
	IsOrReplace bool       // T = replace if already exists
	Funcname    *List      // qualified name of function to create
	Parameters  *List      // list of FunctionParameter
	ReturnType  *TypeName  // return type (NULL if void)
	Options     *List      // list of DefElem
	SqlBody     Node       // SQL body, or NULL
}

func (n *CreateFunctionStmt) Tag() NodeTag { return T_CreateFunctionStmt }

// ReturnStmt represents a RETURN statement in SQL-standard function bodies.
type ReturnStmt struct {
	Returnval Node // return value expression
}

func (n *ReturnStmt) Tag() NodeTag { return T_ReturnStmt }

// FunctionParameter represents a parameter in CREATE FUNCTION.
type FunctionParameter struct {
	Name    string           // parameter name, or NULL if not given
	ArgType *TypeName        // type name
	Mode    FunctionParameterMode // IN/OUT/etc
	Defexpr Node             // default value, or NULL
}

// FunctionParameterMode represents the mode of a function parameter.
type FunctionParameterMode byte

const (
	FUNC_PARAM_IN FunctionParameterMode = 'i'
	FUNC_PARAM_OUT FunctionParameterMode = 'o'
	FUNC_PARAM_INOUT FunctionParameterMode = 'b'
	FUNC_PARAM_VARIADIC FunctionParameterMode = 'v'
	FUNC_PARAM_TABLE FunctionParameterMode = 't'
	FUNC_PARAM_DEFAULT FunctionParameterMode = 'd'
)

func (n *FunctionParameter) Tag() NodeTag { return T_FunctionParameter }

// DoStmt represents a DO statement (anonymous code block).
type DoStmt struct {
	Args *List // list of DefElem
}

func (n *DoStmt) Tag() NodeTag { return T_DoStmt }

// CreateEnumStmt represents a CREATE TYPE ... AS ENUM statement.
type CreateEnumStmt struct {
	TypeName *List // qualified name (list of String)
	Vals     *List // enum values (list of String)
}

func (n *CreateEnumStmt) Tag() NodeTag { return T_CreateEnumStmt }

// AlterEnumStmt represents an ALTER TYPE ... ENUM statement.
type AlterEnumStmt struct {
	Typname           *List  // qualified name (list of String)
	Oldval            string // old enum value name (for RENAME)
	Newval            string // new enum value name
	NewvalNeighbor    string // neighboring enum value for ADD
	NewvalIsAfter     bool   // place new value after neighbor?
	SkipIfNewvalExists bool  // no error if new val exists?
}

func (n *AlterEnumStmt) Tag() NodeTag { return T_AlterEnumStmt }

// CreateDomainStmt represents a CREATE DOMAIN statement.
type CreateDomainStmt struct {
	Domainname  *List        // qualified name
	Typname     *TypeName    // base type
	CollClause  *CollateClause // collation clause
	Constraints *List        // list of Constraint nodes
}

func (n *CreateDomainStmt) Tag() NodeTag { return T_CreateDomainStmt }

// AlterDomainStmt represents an ALTER DOMAIN statement.
type AlterDomainStmt struct {
	Subtype     byte         // 'T' = default, 'N' = NOT NULL, 'O' = drop NOT NULL, 'C' = add constraint, 'X' = drop constraint
	Typname     *List        // qualified name
	Name        string       // constraint name, or NULL
	Def         Node         // definition of default or constraint
	Behavior    DropBehavior // cascade behavior
	MissingOk   bool         // skip if domain doesn't exist?
}

func (n *AlterDomainStmt) Tag() NodeTag { return T_AlterDomainStmt }

// CreateTrigStmt represents a CREATE TRIGGER statement.
type CreateTrigStmt struct {
	Replace       bool      // replace trigger if already exists?
	IsConstraint  bool      // is this a constraint trigger?
	Trigname      string    // trigger name
	Relation      *RangeVar // relation trigger is on
	Funcname      *List     // function to call
	Args          *List     // arguments to the trigger function
	Row           bool      // ROW or STATEMENT trigger
	Timing        int16     // BEFORE, AFTER, or INSTEAD
	Events        int16     // INSERT, UPDATE, DELETE, TRUNCATE
	Columns       *List     // column names, or NIL for all columns
	WhenClause    Node      // WHEN clause
	TransitionRels *List    // list of TransitionTableSpec
	Deferrable    bool      // constraint trigger is deferrable?
	Initdeferred  bool      // constraint trigger is initially deferred?
	Constrrel     *RangeVar // constraint's referenced rel, for FK
}

func (n *CreateTrigStmt) Tag() NodeTag { return T_CreateTrigStmt }

// GrantStmt represents GRANT and REVOKE statements.
type GrantStmt struct {
	IsGrant     bool       // true = GRANT, false = REVOKE
	Targtype    GrantTargetType // type of the grant target
	Objtype     ObjectType // kind of object being operated on
	Objects     *List      // list of object names
	Privileges  *List      // list of AccessPriv nodes
	Grantees    *List      // list of RoleSpec nodes
	GrantOption bool       // grant or revoke grant option
	Grantor     *RoleSpec  // set grantor to other than current role
	Behavior    DropBehavior // drop behavior (RESTRICT/CASCADE)
}

// GrantTargetType represents grant target type.
type GrantTargetType int

const (
	ACL_TARGET_OBJECT GrantTargetType = iota // grant on specific objects
	ACL_TARGET_ALL_IN_SCHEMA // grant on all objects in given schemas
	ACL_TARGET_DEFAULTS      // ALTER DEFAULT PRIVILEGES
)

func (n *GrantStmt) Tag() NodeTag { return T_GrantStmt }

// AccessPriv represents a single privilege in GRANT/REVOKE.
type AccessPriv struct {
	PrivName string // privilege name, NULL for ALL PRIVILEGES
	Cols     *List  // list of String
}

func (n *AccessPriv) Tag() NodeTag { return T_AccessPriv }

// CopyStmt represents a COPY statement.
type CopyStmt struct {
	Relation  *RangeVar // relation to copy to/from
	Query     Node      // the query (SELECT or DML statement)
	Attlist   *List     // list of column names, or NIL for all
	IsFrom    bool      // TO or FROM
	IsProgram bool      // is 'filename' a program?
	Filename  string    // filename, or NULL for stdin/stdout
	Options   *List     // list of DefElem
	WhereClause Node    // WHERE condition (COPY FROM only)
}

func (n *CopyStmt) Tag() NodeTag { return T_CopyStmt }

// ExplainStmt represents an EXPLAIN statement.
type ExplainStmt struct {
	Query   Node  // the query
	Options *List // list of DefElem
}

func (n *ExplainStmt) Tag() NodeTag { return T_ExplainStmt }

// CreateTableAsStmt represents a CREATE TABLE AS statement.
type CreateTableAsStmt struct {
	Query        Node        // the query
	Into         *IntoClause // destination table
	Objtype      ObjectType  // OBJECT_TABLE or OBJECT_MATVIEW
	IsSelectInto bool        // was SELECT INTO
	IfNotExists  bool        // just do nothing if table exists?
}

func (n *CreateTableAsStmt) Tag() NodeTag { return T_CreateTableAsStmt }

// RefreshMatViewStmt represents a REFRESH MATERIALIZED VIEW statement.
type RefreshMatViewStmt struct {
	Concurrent bool      // allow concurrent access?
	SkipData   bool      // don't run SELECT query
	Relation   *RangeVar // relation to refresh
}

func (n *RefreshMatViewStmt) Tag() NodeTag { return T_RefreshMatViewStmt }

// VacuumStmt represents a VACUUM or ANALYZE statement.
type VacuumStmt struct {
	Options  *List // list of DefElem
	Rels     *List // list of VacuumRelation, or NIL for all
	IsVacuumCmd bool // true for VACUUM, false for ANALYZE
}

func (n *VacuumStmt) Tag() NodeTag { return T_VacuumStmt }

// VacuumRelation represents a single relation to vacuum/analyze.
type VacuumRelation struct {
	Relation  *RangeVar // relation to process, or NULL for current database
	Oid       Oid       // OID of relation to process (filled in later)
	VaCols    *List     // list of column names, or NIL for all
}

func (n *VacuumRelation) Tag() NodeTag { return T_VacuumRelation }

// TransactionStmt represents a transaction control statement.
type TransactionStmt struct {
	Kind      TransactionStmtKind // type of transaction statement
	Options   *List               // for BEGIN/START TRANSACTION
	Savepoint string              // for SAVEPOINT, ROLLBACK TO, RELEASE
	Gid       string              // for two-phase commit
	Chain     bool                // AND CHAIN option
	Location  ParseLoc            // token location, or -1 if unknown
}

// TransactionStmtKind represents the kind of transaction statement.
type TransactionStmtKind int

const (
	TRANS_STMT_BEGIN TransactionStmtKind = iota
	TRANS_STMT_START
	TRANS_STMT_COMMIT
	TRANS_STMT_ROLLBACK
	TRANS_STMT_SAVEPOINT
	TRANS_STMT_RELEASE
	TRANS_STMT_ROLLBACK_TO
	TRANS_STMT_PREPARE
	TRANS_STMT_COMMIT_PREPARED
	TRANS_STMT_ROLLBACK_PREPARED
)

func (n *TransactionStmt) Tag() NodeTag { return T_TransactionStmt }

// PrepareStmt represents a PREPARE statement.
type PrepareStmt struct {
	Name     string    // name of plan
	Argtypes *List     // list of TypeName
	Query    Node      // the query itself
}

func (n *PrepareStmt) Tag() NodeTag { return T_PrepareStmt }

// ExecuteStmt represents an EXECUTE statement.
type ExecuteStmt struct {
	Name   string // name of plan
	Params *List  // list of expressions for parameter values
}

func (n *ExecuteStmt) Tag() NodeTag { return T_ExecuteStmt }

// DeallocateStmt represents a DEALLOCATE statement.
type DeallocateStmt struct {
	Name    string // name of plan to deallocate, or NULL for all
	IsAll   bool   // true if DEALLOCATE ALL
	Location ParseLoc // token location
}

func (n *DeallocateStmt) Tag() NodeTag { return T_DeallocateStmt }

// LockStmt represents a LOCK TABLE statement.
type LockStmt struct {
	Relations *List    // list of RangeVar
	Mode      int      // lock mode
	Nowait    bool     // no wait option
}

func (n *LockStmt) Tag() NodeTag { return T_LockStmt }

// SetOperationStmt represents a set-operation (UNION/INTERSECT/EXCEPT) tree.
type SetOperationStmt struct {
	Op            SetOperation // type of set op
	All           bool         // ALL specified?
	Larg          Node         // left child
	Rarg          Node         // right child
	ColTypes      *List        // OID list of output column types
	ColTypmods    *List        // integer list of output column typmods
	ColCollations *List        // OID list of output column collations
	GroupClauses  *List        // list of SortGroupClause
}

func (n *SetOperationStmt) Tag() NodeTag { return T_SetOperationStmt }

// SortGroupClause represents a single element of ORDER BY, GROUP BY, etc.
type SortGroupClause struct {
	TleSortGroupRef uint32 // reference into targetlist
	Eqop            Oid    // operator for equality comparison
	Sortop          Oid    // operator for sorting
	Nulls_first     bool   // sort nulls first
	Hashable        bool   // can hash for grouping?
}

func (n *SortGroupClause) Tag() NodeTag { return T_SortGroupClause }

// RenameStmt represents ALTER ... RENAME statement.
type RenameStmt struct {
	RenameType   ObjectType // OBJECT_TABLE, OBJECT_COLUMN, etc
	RelationType ObjectType // if column, what's the relation type?
	Relation     *RangeVar  // in case it's a table
	Object       Node       // qualified name of object
	Subname      string     // name of contained object (column, rule, etc)
	Newname      string     // new name
	Behavior     DropBehavior // RESTRICT or CASCADE
	MissingOk    bool       // skip error if missing?
}

func (n *RenameStmt) Tag() NodeTag { return T_RenameStmt }

// AlterObjectSchemaStmt represents ALTER ... SET SCHEMA statement.
type AlterObjectSchemaStmt struct {
	ObjectType ObjectType // OBJECT_TABLE, etc
	Relation   *RangeVar  // table, sequence, view, matview, index
	Object     Node       // qualified name of object
	Newschema  string     // new schema name
	MissingOk  bool       // skip error if missing?
}

func (n *AlterObjectSchemaStmt) Tag() NodeTag { return T_AlterObjectSchemaStmt }

// AlterOwnerStmt represents ALTER ... OWNER TO statement.
type AlterOwnerStmt struct {
	ObjectType ObjectType // OBJECT_TABLE, etc
	Relation   *RangeVar  // for relation types
	Object     Node       // qualified name of object
	Newowner   *RoleSpec  // new owner
}

func (n *AlterOwnerStmt) Tag() NodeTag { return T_AlterOwnerStmt }

// ClusterStmt represents a CLUSTER statement.
type ClusterStmt struct {
	Relation  *RangeVar // table to cluster
	Indexname string    // name of index to use, or NULL
	Params    *List     // list of DefElem nodes (PG17+)
}

func (n *ClusterStmt) Tag() NodeTag { return T_ClusterStmt }

// ReindexStmt represents a REINDEX statement.
type ReindexStmt struct {
	Kind      ReindexObjectType // REINDEX_OBJECT_INDEX, etc
	Relation  *RangeVar         // table or index to reindex
	Name      string            // name of database/schema to reindex
	Params    *List             // list of DefElem
}

// ReindexObjectType represents the type of object to reindex.
type ReindexObjectType int

const (
	REINDEX_OBJECT_INDEX ReindexObjectType = iota
	REINDEX_OBJECT_TABLE
	REINDEX_OBJECT_SCHEMA
	REINDEX_OBJECT_SYSTEM
	REINDEX_OBJECT_DATABASE
)

func (n *ReindexStmt) Tag() NodeTag { return T_ReindexStmt }

// CheckPointStmt - CHECKPOINT
type CheckPointStmt struct {
}

func (n *CheckPointStmt) Tag() NodeTag { return T_CheckPointStmt }

// DiscardStmt - DISCARD
type DiscardStmt struct {
	Target DiscardMode
}

func (n *DiscardStmt) Tag() NodeTag { return T_DiscardStmt }

// ListenStmt - LISTEN
type ListenStmt struct {
	Conditionname string
}

func (n *ListenStmt) Tag() NodeTag { return T_ListenStmt }

// UnlistenStmt - UNLISTEN
type UnlistenStmt struct {
	Conditionname string // empty string means UNLISTEN *
}

func (n *UnlistenStmt) Tag() NodeTag { return T_UnlistenStmt }

// NotifyStmt - NOTIFY
type NotifyStmt struct {
	Conditionname string
	Payload       string
}

func (n *NotifyStmt) Tag() NodeTag { return T_NotifyStmt }

// LoadStmt - LOAD
type LoadStmt struct {
	Filename string
}

func (n *LoadStmt) Tag() NodeTag { return T_LoadStmt }

// ClosePortalStmt - CLOSE cursor
type ClosePortalStmt struct {
	Portalname string // empty string means CLOSE ALL
}

func (n *ClosePortalStmt) Tag() NodeTag { return T_ClosePortalStmt }

// ConstraintsSetStmt - SET CONSTRAINTS
type ConstraintsSetStmt struct {
	Constraints *List
	Deferred    bool
}

func (n *ConstraintsSetStmt) Tag() NodeTag { return T_ConstraintsSetStmt }

// VariableSetStmt - SET variable
type VariableSetStmt struct {
	Kind    VariableSetKind
	Name    string
	Args    *List
	IsLocal bool
}

func (n *VariableSetStmt) Tag() NodeTag { return T_VariableSetStmt }

// VariableShowStmt - SHOW variable
type VariableShowStmt struct {
	Name string
}

func (n *VariableShowStmt) Tag() NodeTag { return T_VariableShowStmt }

// DeclareCursorStmt represents a DECLARE cursor statement.
type DeclareCursorStmt struct {
	Portalname string // name of the portal (cursor)
	Options    int    // bitmask of CURSOR_OPT_*
	Query      Node   // the query (SelectStmt)
}

func (n *DeclareCursorStmt) Tag() NodeTag { return T_DeclareCursorStmt }

// FetchStmt represents a FETCH or MOVE statement.
type FetchStmt struct {
	Direction  FetchDirection // see above
	HowMany   int64          // number of rows, or FETCH_ALL
	Portalname string         // name of portal (cursor)
	Ismove     bool           // true if MOVE
}

func (n *FetchStmt) Tag() NodeTag { return T_FetchStmt }

// CallStmt represents a CALL statement.
type CallStmt struct {
	Funccall *FuncCall // the function call
}

func (n *CallStmt) Tag() NodeTag { return T_CallStmt }

// SecLabelStmt represents a SECURITY LABEL statement.
type SecLabelStmt struct {
	Objtype  ObjectType // object kind
	Object   Node       // qualified name of object
	Provider string     // security label provider
	Label    string     // new security label, or empty to drop
}

func (n *SecLabelStmt) Tag() NodeTag { return T_SecLabelStmt }

// CreateRoleStmt represents CREATE ROLE/USER/GROUP statements.
type CreateRoleStmt struct {
	StmtType RoleStmtType // ROLESTMT_ROLE, ROLESTMT_USER, ROLESTMT_GROUP
	Role     string       // role name
	Options  *List        // list of DefElem
}

func (n *CreateRoleStmt) Tag() NodeTag { return T_CreateRoleStmt }

// AlterRoleStmt represents ALTER ROLE/USER/GROUP statements.
type AlterRoleStmt struct {
	Role    *RoleSpec // role specification
	Options *List     // list of DefElem
	Action  int       // +1 = add, -1 = drop
}

func (n *AlterRoleStmt) Tag() NodeTag { return T_AlterRoleStmt }

// AlterRoleSetStmt represents ALTER ROLE SET/RESET statements.
type AlterRoleSetStmt struct {
	Role     *RoleSpec        // role, or NULL for ALL
	Database string           // database name, or empty
	Setstmt  *VariableSetStmt // SET/RESET statement
}

func (n *AlterRoleSetStmt) Tag() NodeTag { return T_AlterRoleSetStmt }

// DropRoleStmt represents DROP ROLE/USER/GROUP statements.
type DropRoleStmt struct {
	Roles     *List // list of RoleSpec nodes
	MissingOk bool  // skip error if role is missing?
}

func (n *DropRoleStmt) Tag() NodeTag { return T_DropRoleStmt }

// GrantRoleStmt represents GRANT role TO role / REVOKE role FROM role.
type GrantRoleStmt struct {
	GrantedRoles *List        // list of roles to grant/revoke (AccessPriv nodes)
	GranteeRoles *List        // list of member roles (RoleSpec nodes)
	IsGrant      bool         // true = GRANT, false = REVOKE
	Opt          *List        // grant options (list of DefElem), PG17+
	Grantor      *RoleSpec    // set grantor to other than current role
	Behavior     DropBehavior // drop behavior for REVOKE
}

func (n *GrantRoleStmt) Tag() NodeTag { return T_GrantRoleStmt }

// CreatedbStmt - CREATE DATABASE
type CreatedbStmt struct {
	Dbname  string
	Options *List // list of DefElem
}

func (n *CreatedbStmt) Tag() NodeTag { return T_CreatedbStmt }

// AlterDatabaseStmt - ALTER DATABASE
type AlterDatabaseStmt struct {
	Dbname  string
	Options *List
}

func (n *AlterDatabaseStmt) Tag() NodeTag { return T_AlterDatabaseStmt }

// AlterDatabaseSetStmt - ALTER DATABASE SET/RESET
type AlterDatabaseSetStmt struct {
	Dbname  string
	Setstmt *VariableSetStmt
}

func (n *AlterDatabaseSetStmt) Tag() NodeTag { return T_AlterDatabaseSetStmt }

// DropdbStmt - DROP DATABASE
type DropdbStmt struct {
	Dbname    string
	MissingOk bool
	Options   *List
}

func (n *DropdbStmt) Tag() NodeTag { return T_DropdbStmt }

// AlterSystemStmt - ALTER SYSTEM SET/RESET
type AlterSystemStmt struct {
	Setstmt *VariableSetStmt
}

func (n *AlterSystemStmt) Tag() NodeTag { return T_AlterSystemStmt }

// AlterCollationStmt represents ALTER COLLATION ... REFRESH VERSION.
type AlterCollationStmt struct {
	Collname *List // qualified name (list of String)
}

func (n *AlterCollationStmt) Tag() NodeTag { return T_AlterCollationStmt }

// DefineStmt represents CREATE AGGREGATE/OPERATOR/TYPE/TEXT SEARCH/COLLATION.
type DefineStmt struct {
	Kind        ObjectType // OBJECT_AGGREGATE, OBJECT_OPERATOR, etc.
	Oldstyle    bool       // old-style (pre-8.2) aggregate syntax
	Defnames    *List      // qualified name (list of String)
	Args        *List      // arguments (for aggregates)
	Definition  *List      // definition (list of DefElem)
	IfNotExists bool       // IF NOT EXISTS
	Replace     bool       // OR REPLACE
}

func (n *DefineStmt) Tag() NodeTag { return T_DefineStmt }

// CompositeTypeStmt represents CREATE TYPE name AS (column_list).
type CompositeTypeStmt struct {
	Typevar    *RangeVar // type name as RangeVar
	Coldeflist *List     // list of ColumnDef
}

func (n *CompositeTypeStmt) Tag() NodeTag { return T_CompositeTypeStmt }

// CreateRangeStmt represents CREATE TYPE name AS RANGE (params).
type CreateRangeStmt struct {
	TypeName *List // qualified name (list of String)
	Params   *List // list of DefElem
}

func (n *CreateRangeStmt) Tag() NodeTag { return T_CreateRangeStmt }

// ObjectWithArgs represents a function/operator name with argument types.
// Used in ALTER FUNCTION, DROP FUNCTION/PROCEDURE/AGGREGATE/OPERATOR etc.
type ObjectWithArgs struct {
	Objname        *List // qualified name (list of String)
	Objargs        *List // argument types (list of TypeName)
	ArgsUnspecified bool  // true if no argument list was given
}

func (n *ObjectWithArgs) Tag() NodeTag { return T_ObjectWithArgs }

// AlterFunctionStmt represents an ALTER FUNCTION/PROCEDURE/ROUTINE statement.
type AlterFunctionStmt struct {
	Objtype ObjectType      // OBJECT_FUNCTION, OBJECT_PROCEDURE, OBJECT_ROUTINE
	Func    *ObjectWithArgs // function with args
	Actions *List           // list of DefElem
}

func (n *AlterFunctionStmt) Tag() NodeTag { return T_AlterFunctionStmt }

// CreateEventTrigStmt represents a CREATE EVENT TRIGGER statement.
type CreateEventTrigStmt struct {
	Trigname   string // trigger name
	Eventname  string // event name (e.g., ddl_command_start)
	Whenclause *List  // list of DefElem (tag filters)
	Funcname   *List  // qualified function name
}

func (n *CreateEventTrigStmt) Tag() NodeTag { return T_CreateEventTrigStmt }

// AlterEventTrigStmt represents an ALTER EVENT TRIGGER statement.
type AlterEventTrigStmt struct {
	Trigname  string // trigger name
	Tgenabled byte   // 'O'=enable, 'D'=disable, 'R'=replica, 'A'=always
}

func (n *AlterEventTrigStmt) Tag() NodeTag { return T_AlterEventTrigStmt }

// RuleStmt represents a CREATE RULE statement.
type RuleStmt struct {
	Relation    *RangeVar // relation the rule is for
	Rulename    string    // name of the rule
	WhereClause Node      // qualifications
	Event       CmdType   // SELECT, INSERT, etc
	Instead     bool      // is a DO INSTEAD rule?
	Actions     *List     // the action statements
	Replace     bool      // OR REPLACE
}

func (n *RuleStmt) Tag() NodeTag { return T_RuleStmt }

// CreatePLangStmt represents a CREATE LANGUAGE statement.
type CreatePLangStmt struct {
	Replace     bool   // OR REPLACE
	Plname      string // PL name
	Plhandler   *List  // PL call handler function (qual. name)
	Plinline    *List  // optional inline handler function (qual. name)
	Plvalidator *List  // optional validator function (qual. name)
	Pltrusted   bool   // PL is trusted
}

func (n *CreatePLangStmt) Tag() NodeTag { return T_CreatePLangStmt }

// TriggerTransition represents a transition table specification in CREATE TRIGGER.
type TriggerTransition struct {
	Name    string // transition relation name
	IsNew   bool   // is it NEW TABLE or OLD TABLE?
	IsTable bool   // is it TABLE or ROW?
}

func (n *TriggerTransition) Tag() NodeTag { return T_TriggerTransition }

// CreateFdwStmt represents a CREATE FOREIGN DATA WRAPPER statement.
type CreateFdwStmt struct {
	Fdwname     string // foreign-data wrapper name
	FuncOptions *List  // HANDLER, VALIDATOR, etc.
	Options     *List  // generic options to FDW
}

func (n *CreateFdwStmt) Tag() NodeTag { return T_CreateFdwStmt }

// AlterFdwStmt represents an ALTER FOREIGN DATA WRAPPER statement.
type AlterFdwStmt struct {
	Fdwname     string // foreign-data wrapper name
	FuncOptions *List  // HANDLER, VALIDATOR, etc.
	Options     *List  // generic options to FDW
}

func (n *AlterFdwStmt) Tag() NodeTag { return T_AlterFdwStmt }

// CreateForeignServerStmt represents a CREATE SERVER statement.
type CreateForeignServerStmt struct {
	Servername  string // server name
	Servertype  string // optional server type
	Version     string // optional server version
	Fdwname     string // FDW name
	IfNotExists bool   // just do nothing if it already exists?
	Options     *List  // generic options to server
}

func (n *CreateForeignServerStmt) Tag() NodeTag { return T_CreateForeignServerStmt }

// AlterForeignServerStmt represents an ALTER SERVER statement.
type AlterForeignServerStmt struct {
	Servername string // server name
	Version    string // optional server version
	Options    *List  // generic options to server
	HasVersion bool   // version was specified
}

func (n *AlterForeignServerStmt) Tag() NodeTag { return T_AlterForeignServerStmt }

// CreateForeignTableStmt represents a CREATE FOREIGN TABLE statement.
type CreateForeignTableStmt struct {
	Base       CreateStmt // base CREATE TABLE fields
	Servername string     // server name
	Options    *List      // generic options to foreign table
}

func (n *CreateForeignTableStmt) Tag() NodeTag { return T_CreateForeignTableStmt }

// CreateUserMappingStmt represents a CREATE USER MAPPING statement.
type CreateUserMappingStmt struct {
	User        *RoleSpec // user role
	Servername  string    // server name
	IfNotExists bool      // just do nothing if it already exists?
	Options     *List     // generic options to user mapping
}

func (n *CreateUserMappingStmt) Tag() NodeTag { return T_CreateUserMappingStmt }

// AlterUserMappingStmt represents an ALTER USER MAPPING statement.
type AlterUserMappingStmt struct {
	User       *RoleSpec // user role
	Servername string    // server name
	Options    *List     // generic options to user mapping
}

func (n *AlterUserMappingStmt) Tag() NodeTag { return T_AlterUserMappingStmt }

// DropUserMappingStmt represents a DROP USER MAPPING statement.
type DropUserMappingStmt struct {
	User       *RoleSpec // user role
	Servername string    // server name
	MissingOk  bool      // skip error if missing?
}

func (n *DropUserMappingStmt) Tag() NodeTag { return T_DropUserMappingStmt }

// ImportForeignSchemaStmt represents an IMPORT FOREIGN SCHEMA statement.
type ImportForeignSchemaStmt struct {
	ServerName   string                 // FDW server name
	RemoteSchema string                 // remote schema name to import
	LocalSchema  string                 // local schema to import into
	ListType     ImportForeignSchemaType // type of table list filter
	TableList    *List                  // list of tables to import or exclude
	Options      *List                  // generic options
}

func (n *ImportForeignSchemaStmt) Tag() NodeTag { return T_ImportForeignSchemaStmt }

// CreateExtensionStmt represents a CREATE EXTENSION statement.
type CreateExtensionStmt struct {
	Extname     string // name of the extension to create
	IfNotExists bool   // just do nothing if it already exists?
	Options     *List  // list of DefElem nodes
}

func (n *CreateExtensionStmt) Tag() NodeTag { return T_CreateExtensionStmt }

// AlterExtensionStmt represents an ALTER EXTENSION UPDATE statement.
type AlterExtensionStmt struct {
	Extname string // name of the extension to alter
	Options *List  // list of DefElem nodes
}

func (n *AlterExtensionStmt) Tag() NodeTag { return T_AlterExtensionStmt }

// AlterExtensionContentsStmt represents ALTER EXTENSION ADD/DROP object.
type AlterExtensionContentsStmt struct {
	Extname string     // name of the extension
	Action  int        // +1 = ADD, -1 = DROP
	Objtype ObjectType // object type
	Object  Node       // qualified name of the object
}

func (n *AlterExtensionContentsStmt) Tag() NodeTag { return T_AlterExtensionContentsStmt }

// CreateTableSpaceStmt represents a CREATE TABLESPACE statement.
type CreateTableSpaceStmt struct {
	Tablespacename string    // name of the tablespace to create
	Owner          *RoleSpec // owner of the tablespace
	Location       string    // directory path for tablespace
	Options        *List     // WITH clause options
}

func (n *CreateTableSpaceStmt) Tag() NodeTag { return T_CreateTableSpaceStmt }

// DropTableSpaceStmt represents a DROP TABLESPACE statement.
type DropTableSpaceStmt struct {
	Tablespacename string // name of the tablespace to drop
	MissingOk      bool   // skip error if missing?
}

func (n *DropTableSpaceStmt) Tag() NodeTag { return T_DropTableSpaceStmt }

// AlterTableSpaceOptionsStmt represents ALTER TABLESPACE SET/RESET statement.
type AlterTableSpaceOptionsStmt struct {
	Tablespacename string // name of the tablespace
	Options        *List  // list of DefElem
	IsReset        bool   // true for RESET, false for SET
}

func (n *AlterTableSpaceOptionsStmt) Tag() NodeTag { return T_AlterTableSpaceOptionsStmt }

// CreateAmStmt represents a CREATE ACCESS METHOD statement.
type CreateAmStmt struct {
	Amname      string // access method name
	HandlerName *List  // handler function name
	Amtype      byte   // 'i' for INDEX, 't' for TABLE
}

func (n *CreateAmStmt) Tag() NodeTag { return T_CreateAmStmt }

// CreatePolicyStmt represents a CREATE POLICY statement.
type CreatePolicyStmt struct {
	PolicyName string    // policy name
	Table      *RangeVar // the table the policy applies to
	CmdName    string    // the command name (all, select, insert, update, delete)
	Permissive bool      // restrictive or permissive policy
	Roles      *List     // list of roles (RoleSpec)
	Qual       Node      // USING qualification
	WithCheck  Node      // WITH CHECK qualification
}

func (n *CreatePolicyStmt) Tag() NodeTag { return T_CreatePolicyStmt }

// AlterPolicyStmt represents an ALTER POLICY statement.
type AlterPolicyStmt struct {
	PolicyName string    // policy name
	Table      *RangeVar // the table the policy applies to
	Roles      *List     // list of roles (RoleSpec)
	Qual       Node      // USING qualification
	WithCheck  Node      // WITH CHECK qualification
}

func (n *AlterPolicyStmt) Tag() NodeTag { return T_AlterPolicyStmt }

// CreatePublicationStmt represents a CREATE PUBLICATION statement.
type CreatePublicationStmt struct {
	Pubname      string // publication name
	Options      *List  // list of DefElem nodes
	Pubobjects   *List  // list of PublicationObjSpec
	ForAllTables bool   // FOR ALL TABLES
}

func (n *CreatePublicationStmt) Tag() NodeTag { return T_CreatePublicationStmt }

// AlterPublicationStmt represents an ALTER PUBLICATION statement.
type AlterPublicationStmt struct {
	Pubname      string        // publication name
	Options      *List         // list of DefElem nodes
	Pubobjects   *List         // list of PublicationObjSpec
	ForAllTables bool          // FOR ALL TABLES
	Action       DefElemAction // SET, ADD, DROP
}

func (n *AlterPublicationStmt) Tag() NodeTag { return T_AlterPublicationStmt }

// PublicationObjSpec represents a publication object specification.
type PublicationObjSpec struct {
	Pubobjtype PublicationObjSpecType // object type
	Name       string                 // schema name for TABLES IN SCHEMA
	Pubtable   *PublicationTable      // table specification
	Location   ParseLoc               // token location
}

func (n *PublicationObjSpec) Tag() NodeTag { return T_PublicationObjSpec }

// PublicationTable represents a table in a publication.
type PublicationTable struct {
	Relation    *RangeVar // relation to publish
	WhereClause Node      // WHERE clause for row filter
	Columns     *List     // column list filter
}

func (n *PublicationTable) Tag() NodeTag { return T_PublicationTable }

// CreateSubscriptionStmt represents a CREATE SUBSCRIPTION statement.
type CreateSubscriptionStmt struct {
	Subname     string // subscription name
	Conninfo    string // connection string
	Publication *List  // list of publication names (String nodes)
	Options     *List  // list of DefElem
}

func (n *CreateSubscriptionStmt) Tag() NodeTag { return T_CreateSubscriptionStmt }

// AlterSubscriptionStmt represents an ALTER SUBSCRIPTION statement.
type AlterSubscriptionStmt struct {
	Kind        AlterSubscriptionType // kind of alter
	Subname     string                // subscription name
	Conninfo    string                // connection string
	Publication *List                 // list of publication names
	Options     *List                 // list of DefElem
}

func (n *AlterSubscriptionStmt) Tag() NodeTag { return T_AlterSubscriptionStmt }

// DropSubscriptionStmt represents a DROP SUBSCRIPTION statement.
type DropSubscriptionStmt struct {
	Subname   string       // subscription name
	MissingOk bool         // skip error if missing?
	Behavior  DropBehavior // RESTRICT or CASCADE
}

func (n *DropSubscriptionStmt) Tag() NodeTag { return T_DropSubscriptionStmt }

// AlterObjectDependsStmt represents ALTER ... DEPENDS ON EXTENSION statement.
type AlterObjectDependsStmt struct {
	ObjectType ObjectType // OBJECT_FUNCTION, etc
	Relation   *RangeVar  // for relation types
	Object     Node       // qualified name of object
	Extname    Node       // extension name (String node)
	Remove     bool       // true if NO DEPENDS
}

func (n *AlterObjectDependsStmt) Tag() NodeTag { return T_AlterObjectDependsStmt }

// AlterOperatorStmt represents ALTER OPERATOR ... SET (...) statement.
type AlterOperatorStmt struct {
	Opername *ObjectWithArgs // operator name and argument types
	Options  *List           // list of DefElem
}

func (n *AlterOperatorStmt) Tag() NodeTag { return T_AlterOperatorStmt }

// AlterTypeStmt represents ALTER TYPE ... SET (...) statement.
type AlterTypeStmt struct {
	TypeName *List // qualified name (list of String)
	Options  *List // list of DefElem
}

func (n *AlterTypeStmt) Tag() NodeTag { return T_AlterTypeStmt }

// AlterDefaultPrivilegesStmt represents ALTER DEFAULT PRIVILEGES statement.
type AlterDefaultPrivilegesStmt struct {
	Options *List      // list of DefElem
	Action  *GrantStmt // the GRANT/REVOKE action
}

func (n *AlterDefaultPrivilegesStmt) Tag() NodeTag { return T_AlterDefaultPrivilegesStmt }

// AlterTSDictionaryStmt represents ALTER TEXT SEARCH DICTIONARY statement.
type AlterTSDictionaryStmt struct {
	Dictname *List // qualified name
	Options  *List // definition list
}

func (n *AlterTSDictionaryStmt) Tag() NodeTag { return T_AlterTSDictionaryStmt }

// AlterTSConfigurationStmt represents ALTER TEXT SEARCH CONFIGURATION statement.
type AlterTSConfigurationStmt struct {
	Kind      AlterTSConfigType // ADD/ALTER/DROP MAPPING etc
	Cfgname   *List             // qualified config name
	Tokentype *List             // list of token type names
	Dicts     *List             // list of dictionary names
	Override  bool              // if ALTER MAPPING
	Replace   bool              // if replacing dicts
	MissingOk bool              // for IF EXISTS
}

func (n *AlterTSConfigurationStmt) Tag() NodeTag { return T_AlterTSConfigurationStmt }

// CreateStatsStmt represents CREATE STATISTICS statement.
type CreateStatsStmt struct {
	Defnames    *List  // qualified name
	StatTypes   *List  // stat types (list of String)
	Exprs       *List  // expressions (list of StatsElem)
	Relations   *List  // FROM clause (list of RangeVar)
	Stxcomment  string // comment
	IfNotExists bool   // IF NOT EXISTS
}

func (n *CreateStatsStmt) Tag() NodeTag { return T_CreateStatsStmt }

// StatsElem represents a statistics element.
type StatsElem struct {
	Name string // column name, or NULL
	Expr Node   // expression, or NULL
}

func (n *StatsElem) Tag() NodeTag { return T_StatsElem }

// AlterStatsStmt represents ALTER STATISTICS statement.
type AlterStatsStmt struct {
	Defnames      *List // qualified name
	MissingOk     bool  // IF EXISTS
	Stxstattarget int   // new statistics target
}

func (n *AlterStatsStmt) Tag() NodeTag { return T_AlterStatsStmt }

// CreateOpClassStmt represents CREATE OPERATOR CLASS statement.
type CreateOpClassStmt struct {
	Opclassname  *List     // qualified name
	Opfamilyname *List     // qualified opfamily name, or NIL
	Amname       string    // access method name
	Datatype     *TypeName // datatype of indexed column
	Items        *List     // list of CreateOpClassItem
	IsDefault    bool      // DEFAULT
}

func (n *CreateOpClassStmt) Tag() NodeTag { return T_CreateOpClassStmt }

// CreateOpClassItem represents an item in CREATE OPERATOR CLASS.
type CreateOpClassItem struct {
	Itemtype   int             // see OPCLASS_ITEM_* constants
	Name       *ObjectWithArgs // operator or function
	Number     int             // strategy number or support proc number
	OrderFamily *List          // opfamily for ordering
	ClassArgs  *List           // type arguments
	Storedtype *TypeName       // storage type
}

func (n *CreateOpClassItem) Tag() NodeTag { return T_CreateOpClassItem }

// CreateOpFamilyStmt represents CREATE OPERATOR FAMILY statement.
type CreateOpFamilyStmt struct {
	Opfamilyname *List  // qualified name
	Amname       string // access method name
}

func (n *CreateOpFamilyStmt) Tag() NodeTag { return T_CreateOpFamilyStmt }

// AlterOpFamilyStmt represents ALTER OPERATOR FAMILY statement.
type AlterOpFamilyStmt struct {
	Opfamilyname *List  // qualified name
	Amname       string // access method name
	IsDrop       bool   // true for DROP, false for ADD
	Items        *List  // list of CreateOpClassItem
}

func (n *AlterOpFamilyStmt) Tag() NodeTag { return T_AlterOpFamilyStmt }

// CreateCastStmt represents CREATE CAST statement.
type CreateCastStmt struct {
	Sourcetype *TypeName       // source type
	Targettype *TypeName       // target type
	Func       *ObjectWithArgs // cast function, or NULL
	Context    CoercionContext // cast context
	Inout      bool            // WITH INOUT
}

func (n *CreateCastStmt) Tag() NodeTag { return T_CreateCastStmt }

// CreateTransformStmt represents CREATE TRANSFORM statement.
type CreateTransformStmt struct {
	Replace  bool            // OR REPLACE
	TypeName *TypeName       // target type
	Lang     string          // language name
	Fromsql  *ObjectWithArgs // FROM SQL function
	Tosql    *ObjectWithArgs // TO SQL function
}

func (n *CreateTransformStmt) Tag() NodeTag { return T_CreateTransformStmt }

// CreateConversionStmt represents CREATE CONVERSION statement.
type CreateConversionStmt struct {
	ConversionName  *List  // conversion name
	ForEncodingName string // source encoding
	ToEncodingName  string // target encoding
	FuncName        *List  // function name
	Def             bool   // DEFAULT?
}

func (n *CreateConversionStmt) Tag() NodeTag { return T_CreateConversionStmt }

// DropOwnedStmt represents DROP OWNED BY statement.
type DropOwnedStmt struct {
	Roles    *List        // list of RoleSpec
	Behavior DropBehavior // RESTRICT or CASCADE
}

func (n *DropOwnedStmt) Tag() NodeTag { return T_DropOwnedStmt }

// ReassignOwnedStmt represents REASSIGN OWNED BY statement.
type ReassignOwnedStmt struct {
	Roles   *List     // list of RoleSpec
	Newrole *RoleSpec // new owner
}

func (n *ReassignOwnedStmt) Tag() NodeTag { return T_ReassignOwnedStmt }

// SQLValueFunction represents SQL-standard functions that don't require
// a function call syntax, e.g. CURRENT_DATE, CURRENT_USER, etc.
type SQLValueFunction struct {
	Op       SVFOp    // which function this is
	Typmod   int32    // typmod to apply, or -1
	Location ParseLoc // token location, or -1
}

func (n *SQLValueFunction) Tag() NodeTag { return T_SQLValueFunction }

// SetToDefault represents a DEFAULT marker in expressions.
type SetToDefault struct {
	TypeId   Oid      // type for substituted value
	Typmod   int32    // typemod for substituted value
	Collation Oid     // collation for the datatype
	Location ParseLoc // token location, or -1
}

func (n *SetToDefault) Tag() NodeTag { return T_SetToDefault }

// XmlExprOp represents the type of XML expression.
type XmlExprOp int

const (
	IS_XMLCONCAT    XmlExprOp = iota // XMLCONCAT(args)
	IS_XMLELEMENT                    // XMLELEMENT(name, xml_attributes, args)
	IS_XMLFOREST                     // XMLFOREST(xml_attributes)
	IS_XMLPARSE                      // XMLPARSE(text, is_doc, preserve_ws)
	IS_XMLPI                         // XMLPI(name [, args])
	IS_XMLROOT                       // XMLROOT(xml, version, standalone)
	IS_XMLSERIALIZE                  // XMLSERIALIZE(is_document, xmlval, indent)
	IS_DOCUMENT                      // xmlval IS DOCUMENT
)

// XmlOptionType for DOCUMENT or CONTENT
type XmlOptionType int

const (
	XMLOPTION_DOCUMENT XmlOptionType = iota
	XMLOPTION_CONTENT
)

// XML standalone constants
const (
	XML_STANDALONE_YES      = 0
	XML_STANDALONE_NO       = 1
	XML_STANDALONE_NO_VALUE = 2
	XML_STANDALONE_OMITTED  = 3
)

// XmlExpr represents various SQL/XML functions requiring special grammar.
type XmlExpr struct {
	Op        XmlExprOp     // xml function ID
	Name      string        // name in xml(NAME foo ...) syntaxes
	NamedArgs *List         // non-XML expressions for xml_attributes
	ArgNames  *List         // parallel list of String values
	Args      *List         // list of expressions
	Xmloption XmlOptionType // DOCUMENT or CONTENT
	Indent    bool          // INDENT option for XMLSERIALIZE
	Type      Oid           // target type for XMLSERIALIZE
	Typmod    int32         // target typmod for XMLSERIALIZE
	Location  ParseLoc      // token location, or -1
}

func (n *XmlExpr) Tag() NodeTag { return T_XmlExpr }

// XmlSerialize represents XMLSERIALIZE(DOCUMENT|CONTENT expr AS typename).
type XmlSerialize struct {
	Xmloption XmlOptionType // DOCUMENT or CONTENT
	Expr      Node          // the XML expression
	TypeName  *TypeName     // target type name
	Indent    bool          // INDENT option
	Location  ParseLoc      // token location, or -1
}

func (n *XmlSerialize) Tag() NodeTag { return T_XmlSerialize }

// RangeTableFunc represents raw form of table functions such as XMLTABLE.
type RangeTableFunc struct {
	Lateral    bool     // does it have LATERAL prefix?
	Docexpr    Node     // document expression
	Rowexpr    Node     // row generator expression
	Namespaces *List    // list of namespaces as ResTarget
	Columns    *List    // list of RangeTableFuncCol
	Alias      *Alias   // table alias & optional column aliases
	Location   ParseLoc // token location, or -1
}

func (n *RangeTableFunc) Tag() NodeTag { return T_RangeTableFunc }

// RangeTableFuncCol represents one column in a RangeTableFunc.
type RangeTableFuncCol struct {
	Colname       string    // name of generated column
	TypeName      *TypeName // type of generated column
	ForOrdinality bool      // does it have FOR ORDINALITY?
	IsNotNull     bool      // does it have NOT NULL?
	Colexpr       Node      // column filter expression (PATH)
	Coldefexpr    Node      // column default value expression
	Location      ParseLoc  // token location, or -1
}

func (n *RangeTableFuncCol) Tag() NodeTag { return T_RangeTableFuncCol }

// ===== SQL/JSON node types =====

// JsonEncoding represents JSON encoding type.
type JsonEncoding int

const (
	JS_ENC_DEFAULT JsonEncoding = iota
	JS_ENC_UTF8
	JS_ENC_UTF16
	JS_ENC_UTF32
)

// JsonFormatType represents JSON format type.
type JsonFormatType int

const (
	JS_FORMAT_DEFAULT JsonFormatType = iota
	JS_FORMAT_JSON
	JS_FORMAT_JSONB
)

// JsonFormat represents a JSON format clause.
type JsonFormat struct {
	FormatType JsonFormatType
	Encoding   JsonEncoding
	Location   ParseLoc
}

func (n *JsonFormat) Tag() NodeTag { return T_JsonFormat }

// JsonReturning represents the RETURNING clause of a JSON function.
type JsonReturning struct {
	Format *JsonFormat
	Typid  Oid
	Typmod int32
}

func (n *JsonReturning) Tag() NodeTag { return T_JsonReturning }

// JsonValueExpr represents a JSON value expression.
type JsonValueExpr struct {
	RawExpr       Node
	FormattedExpr Node
	Format        *JsonFormat
}

func (n *JsonValueExpr) Tag() NodeTag { return T_JsonValueExpr }

// JsonOutput represents the output clause of JSON constructors.
type JsonOutput struct {
	TypeName  *TypeName
	Returning *JsonReturning
}

func (n *JsonOutput) Tag() NodeTag { return T_JsonOutput }

// JsonArgument represents an argument in PASSING clause.
type JsonArgument struct {
	Val  *JsonValueExpr
	Name string
}

func (n *JsonArgument) Tag() NodeTag { return T_JsonArgument }

// JsonQuotes represents KEEP/OMIT QUOTES option.
type JsonQuotes int

const (
	JS_QUOTES_UNSPEC JsonQuotes = iota
	JS_QUOTES_KEEP
	JS_QUOTES_OMIT
)

// JsonWrapper represents wrapper behavior.
type JsonWrapper int

const (
	JSW_UNSPEC JsonWrapper = iota
	JSW_NONE
	JSW_CONDITIONAL
	JSW_UNCONDITIONAL
)

// JsonBehaviorType represents JSON behavior types.
type JsonBehaviorType int

const (
	JSON_BEHAVIOR_NULL JsonBehaviorType = iota
	JSON_BEHAVIOR_ERROR
	JSON_BEHAVIOR_EMPTY
	JSON_BEHAVIOR_TRUE
	JSON_BEHAVIOR_FALSE
	JSON_BEHAVIOR_UNKNOWN
	JSON_BEHAVIOR_EMPTY_ARRAY
	JSON_BEHAVIOR_EMPTY_OBJECT
	JSON_BEHAVIOR_DEFAULT
)

// JsonBehavior represents ON ERROR / ON EMPTY behavior.
type JsonBehavior struct {
	Btype    JsonBehaviorType
	Expr     Node
	Coerce   Node
	Location ParseLoc
}

func (n *JsonBehavior) Tag() NodeTag { return T_JsonBehavior }

// JsonExprOp represents JSON function operation type.
type JsonExprOp int

const (
	JSON_EXISTS_OP JsonExprOp = iota
	JSON_QUERY_OP
	JSON_VALUE_OP
	JSON_TABLE_OP
)

// JsonFuncExpr represents JSON_VALUE, JSON_QUERY, JSON_EXISTS function calls.
type JsonFuncExpr struct {
	Op          JsonExprOp
	ColumnName  string
	ContextItem *JsonValueExpr
	Pathspec    Node
	Passing     *List
	Output      *JsonOutput
	OnEmpty     *JsonBehavior
	OnError     *JsonBehavior
	Wrapper     JsonWrapper
	Quotes      JsonQuotes
	Location    ParseLoc
}

func (n *JsonFuncExpr) Tag() NodeTag { return T_JsonFuncExpr }

// JsonTablePathSpec represents a path specification in JSON_TABLE.
type JsonTablePathSpec struct {
	String       Node
	Name         string
	NameLocation ParseLoc
	Location     ParseLoc
}

func (n *JsonTablePathSpec) Tag() NodeTag { return T_JsonTablePathSpec }

// JsonTableColumnType represents the type of a JSON_TABLE column.
type JsonTableColumnType int

const (
	JTC_FOR_ORDINALITY JsonTableColumnType = iota
	JTC_REGULAR
	JTC_EXISTS
	JTC_FORMATTED
	JTC_NESTED
)

// JsonTableColumn represents a column definition in JSON_TABLE.
type JsonTableColumn struct {
	Coltype  JsonTableColumnType
	Name     string
	TypeName *TypeName
	Pathspec *JsonTablePathSpec
	Format   *JsonFormat
	Wrapper  JsonWrapper
	Quotes   JsonQuotes
	Columns  *List
	OnEmpty  *JsonBehavior
	OnError  *JsonBehavior
	Location ParseLoc
}

func (n *JsonTableColumn) Tag() NodeTag { return T_JsonTableColumn }

// JsonTable represents JSON_TABLE function.
type JsonTable struct {
	ContextItem *JsonValueExpr
	Pathspec    *JsonTablePathSpec
	Passing     *List
	Columns     *List
	OnError     *JsonBehavior
	Alias       *Alias
	Lateral     bool
	Location    ParseLoc
}

func (n *JsonTable) Tag() NodeTag { return T_JsonTable }

// JsonKeyValue represents a key-value pair in JSON_OBJECT.
type JsonKeyValue struct {
	Key   Node
	Value *JsonValueExpr
}

func (n *JsonKeyValue) Tag() NodeTag { return T_JsonKeyValue }

// JsonParseExpr represents JSON() parse expression.
type JsonParseExpr struct {
	Expr       *JsonValueExpr
	Output     *JsonOutput
	UniqueKeys bool
	Location   ParseLoc
}

func (n *JsonParseExpr) Tag() NodeTag { return T_JsonParseExpr }

// JsonScalarExpr represents JSON_SCALAR() expression.
type JsonScalarExpr struct {
	Expr     Node
	Output   *JsonOutput
	Location ParseLoc
}

func (n *JsonScalarExpr) Tag() NodeTag { return T_JsonScalarExpr }

// JsonSerializeExpr represents JSON_SERIALIZE() expression.
type JsonSerializeExpr struct {
	Expr     *JsonValueExpr
	Output   *JsonOutput
	Location ParseLoc
}

func (n *JsonSerializeExpr) Tag() NodeTag { return T_JsonSerializeExpr }

// JsonObjectConstructor represents JSON_OBJECT() constructor.
type JsonObjectConstructor struct {
	Exprs        *List // list of JsonKeyValue
	Output       *JsonOutput
	AbsentOnNull bool
	UniqueKeys   bool
	Location     ParseLoc
}

func (n *JsonObjectConstructor) Tag() NodeTag { return T_JsonObjectConstructor }

// JsonArrayConstructor represents JSON_ARRAY() with value list.
type JsonArrayConstructor struct {
	Exprs        *List // list of JsonValueExpr
	Output       *JsonOutput
	AbsentOnNull bool
	Location     ParseLoc
}

func (n *JsonArrayConstructor) Tag() NodeTag { return T_JsonArrayConstructor }

// JsonArrayQueryConstructor represents JSON_ARRAY() with subquery.
type JsonArrayQueryConstructor struct {
	Query        Node
	Output       *JsonOutput
	Format       *JsonFormat
	AbsentOnNull bool
	Location     ParseLoc
}

func (n *JsonArrayQueryConstructor) Tag() NodeTag { return T_JsonArrayQueryConstructor }

// JsonAggConstructor represents common aggregate constructor fields.
type JsonAggConstructor struct {
	Output   *JsonOutput
	Agg_filter Node
	Agg_order *List
	Over     *WindowDef
	Location ParseLoc
}

func (n *JsonAggConstructor) Tag() NodeTag { return T_JsonAggConstructor }

// JsonObjectAgg represents JSON_OBJECTAGG() aggregate.
type JsonObjectAgg struct {
	Constructor  *JsonAggConstructor
	Arg          *JsonKeyValue
	AbsentOnNull bool
	UniqueKeys   bool
}

func (n *JsonObjectAgg) Tag() NodeTag { return T_JsonObjectAgg }

// JsonArrayAgg represents JSON_ARRAYAGG() aggregate.
type JsonArrayAgg struct {
	Constructor  *JsonAggConstructor
	Arg          *JsonValueExpr
	AbsentOnNull bool
}

func (n *JsonArrayAgg) Tag() NodeTag { return T_JsonArrayAgg }

// JsonValueType for IS JSON predicate.
type JsonValueType int

const (
	JS_TYPE_ANY    JsonValueType = iota
	JS_TYPE_OBJECT
	JS_TYPE_ARRAY
	JS_TYPE_SCALAR
)

// JsonIsPredicate represents expr IS JSON predicate.
type JsonIsPredicate struct {
	Expr       Node
	Format     *JsonFormat
	ItemType   JsonValueType
	UniqueKeys bool
	Location   ParseLoc
}

func (n *JsonIsPredicate) Tag() NodeTag { return T_JsonIsPredicate }
