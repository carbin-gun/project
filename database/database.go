package database

import (
	"database/sql"
)

type Column struct {
	Schema       string
	TableName    string
	ColumnName   string
	DefaultValue string
	DataType     string
	ColumnType   string
	ColumnKey    string
	Extra        string
	Comment      string
}

type Table []Column
type Schema map[string]Table

var SupportedDrivers map[string]Driver

type Driver interface {
	Load(dsnString string, schema string, tableNames string) (Schema, error)
	GenerateCode(schema Schema)
}

func init() {
	SupportedDrivers = map[string]Driver{
		"mysql":    MysqlDriver{},
		"postgres": PostgresDriver{},
	}
}

func AsString(rb sql.RawBytes) string {
	if len(rb) > 0 {
		return string(rb)
	}
	return ""
}
