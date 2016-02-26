package database

import "database/sql"

var SupportedDrivers map[string]Driver

type Driver interface {
	Load(dsnString string, schema string, tableNames string) (Schema, error)
	GenerateCode(dbName string, schema Schema, templatePath string, targetDir string)
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
