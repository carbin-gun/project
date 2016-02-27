package database

/**
for postgres information_schema.table_constraints and information_schema.key_column_usage
beacause the primary key is in information_schema.table_constraints,but it does not show column name
the information_schema.key_column_usage,shows the column and the two table can do the join op
*/
type TableConstraints struct {
	//the following fields are from information_schema.table_constraints
	ConstraintCatalog string `json:"constraint_catalog"`
	ConstraintSchema  string `json:"constraint_schema"`
	ConstraintName    string `json:"constraint_name"`
	TableCatalog      string `json:"table_catalog"`
	TableSchema       string `json:"table_schema"`
	TableName         string `json:"table_name"`
	ConstraintType    string `json:"constraint_type"`
	IsDeferrable      string `json:"is_deferrable"`
	InitiallyDeferred string `json:"initially_deferred"`
	//the ColumnName field is from information_schema.key_column_usage.
	ColumnName string `json:column_name`
}

/**
for postgres information_schema.columns
*/
type TableColumns struct {
	TableCatalog           string `json:"table_catalog"`
	TableSchema            string `json:"table_schema"`
	TableName              string `json:"table_name"`
	ColumnName             string `json:"column_name"`
	OrdinalPosition        int    `json:"ordinal_position"`
	ColumnDefault          string `json:"column_default"`
	IsNullable             string `json:"is_nullable"`
	DataType               string `json:"data_type"`
	CharacterMaximumLength int    `json:"character_maximum_length"`
	CharacterOctetLength   int    `json:"character_octet_length"`
	NumericPrecision       int    `json:"numeric_precision"`
	NumericPrecisionRadix  int    `json:"numeric_precision_radix"`
	NumericScale           int    `json:"numeric_scale"`
	DatetimePrecision      int    `json:"datetime_precision"`
	IntervalType           string `json:"interval_type"`
	IntervalPrecision      int    `json:"interval_precision"`
	CharacterSetCatalog    string `json:"character_set_catalog"`
	CharacterSetSchema     string `json:"character_set_schema"`
	CharacterSetName       string `json:"character_set_name"`
	CollationCatalog       string `json:"collation_catalog"`
	CollationSchema        string `json:"collation_schema"`
	CollationName          string `json:"collation_name"`
	DomainCatalog          string `json:"domain_catalog"`
	DomainSchema           string `json:"domain_schema"`
	DomainName             string `json:"domain_name"`
	UdtCatalog             string `json:"udt_catalog"`
	UdtSchema              string `json:"udt_schema"`
	UdtName                string `json:"udt_name"`
	ScopeCatalog           string `json:"scope_catalog"`
	ScopeSchema            string `json:"scope_schema"`
	ScopeName              string `json:"scope_name"`
	MaximumCardinality     int    `json:"maximum_cardinality"`
	DtdIdentifier          string `json:"dtd_identifier"`
	IsSelfReferencing      string `json:"is_self_referencing"`
	IsIdentity             string `json:"is_identity"`
	IdentityGeneration     string `json:"identity_generation"`
	IdentityStart          string `json:"identity_start"`
	IdentityIncrement      string `json:"identity_increment"`
	IdentityMaximum        string `json:"identity_maximum"`
	IdentityMinimum        string `json:"identity_minimum"`
	IdentityCycle          string `json:"identity_cycle"`
	IsGenerated            string `json:"is_generated"`
	GenerationExpression   string `json:"generation_expression"`
	IsUpdatable            string `json:"is_updatable"`
}

type IndexInfo interface {
	IsIndex() bool
	IsUniqueIndex() bool
	IsPrimaryKey() bool
}
type PrimaryIndex struct{}
type UniqueIndex struct{}
type OrdinaryIndex struct{}
type NonIndex struct{}

func (p PrimaryIndex) IsIndex() bool {
	return true
}
func (p PrimaryIndex) IsUniqueIndex() bool {
	return true
}
func (p PrimaryIndex) IsPrimaryKey() bool {
	return true
}

func (p UniqueIndex) IsPrimaryKey() bool {
	return false
}
func (p UniqueIndex) IxUniqueIndex() bool {
	return true
}
func (p UniqueIndex) IsIndex() bool {
	return true
}

func (p OrdinaryIndex) IsPrimaryKey() bool {
	return false
}
func (p OrdinaryIndex) IxUniqueIndex() bool {
	return false
}
func (p OrdinaryIndex) IsIndex() bool {
	return true
}
func (p NonIndex) IsPrimaryKey() bool {
	return false
}
func (p NonIndex) IsUniqueIndex() bool {
	return false
}
func (p NonIndex) IsIndex() bool {
	return false
}

type Table []Column
type Schema map[string]Table

/**
universally abstracted to represent a column.
*/
type Column struct {
	Schema       string
	TableName    string
	ColumnName   string
	DefaultValue string
	DataType     string
	ColumnType   string
	ColumnKey    IndexInfo
	Extra        string
	Comment      string
}
