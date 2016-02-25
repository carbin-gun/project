package database

type TableConstraints struct {
	ConstraintCatalog string `json:"constraint_catalog"`
	ConstraintSchema  string `json:"constraint_schema"`
	ConstraintName    string `json:"constraint_name"`
	TableCatalog      string `json:"table_catalog"`
	TableSchema       string `json:"table_schema"`
	TableName         string `json:"table_name"`
	ConstraintType    string `json:"constraint_type"`
	IsDeferrable      string `json:"is_deferrable"`
	InitiallyDeferred string `json:"initially_deferred"`
}
