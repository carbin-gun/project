package database

type ModelMeta struct {
	Name          string
	DbName        string
	TableName     string
	PrimaryFields PrimaryFields
	Fields        []ModelField
	Uniques       []ModelField
}
type PrimaryFields []*ModelField
type ModelField struct {
	Name            string
	ColumnName      string
	Type            string
	JsonMeta        string
	IsPrimaryKey    bool
	IsUniqueKey     bool
	IsAutoIncrement bool
	DefaultValue    string
	Extra           string
	Comment         string
}

type Table []Column
type Schema map[string]Table
