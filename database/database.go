package database

import (
	"database/sql"
	"strconv"
)

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

func AsInt64(rb sql.RawBytes) int64 {
	if len(rb) > 0 {
		if n, err := strconv.ParseInt(string(rb), 10, 64); err == nil {
			return n
		}
	}
	return 0
}

func AsInt(rb sql.RawBytes) int {
	return int(AsInt64(rb))
}
