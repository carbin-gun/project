package database

import (
	"database/sql"

	"strings"

	"fmt"

	"github.com/carbin-gun/project/common"
	_ "github.com/lib/pq"
	"github.com/prometheus/log"
)

const (
	COLUMNS               = `constraint_catalog,constraint_schema,constraint_name,table_catalog,table_schema,table_name,constraint_type,is_deferrable,initially_deferred`
	PRIMARY_SQL           = `select ` + COLUMNS + ` from information_schema.table_constraints where table_schema = $1 and constraint_type='PRIMARY KEY'`
	PRIMARY_SQL_IN_TABLES = `select ` + COLUMNS + ` from information_schema.table_constraints where table_schema = $1 and table_name in $2 and constraint_type='PRIMARY KEY' `
)

type PostgresDriver struct{}

type ColumnKeyValue struct {
	Name  string
	Value interface{}
}
type QueryRowVisitor func(columns []string, rb []sql.RawBytes) bool

func (postgresDriver PostgresDriver) Load(dsnString string, schema string, tableNames string) (Schema, error) {
	log.Printf("[postgres driver] is loading,schema:%s,tableNames:%s", schema, tableNames)
	db, err := sql.Open("postgres", dsnString)
	common.PanicOnError(err, "[postgres driver] connect error")
	defer db.Close()
	if err = db.Ping(); err != nil {
		common.PanicOnError(err, "[postgres driver] Ping error")
	}
	ret := postgresDriver.queryColumns(db, schema, tableNames)
	log.Printf("[Postgres Driver] Loaded schema data of %d tables from database schema[%s]", len(ret), schema)
	return ret, nil
}

func (postgresDriver PostgresDriver) queryColumns(db *sql.DB, schema, tableNames string) Schema {
	ret := make(Schema)
	tableConstraints := postgresDriver.QueryPrimaryKeys(db, schema, tableNames)
	for _, constraint := range tableConstraints {
		fmt.Println(constraint.TableName)
	}
	return ret

}

func (postgresDriver PostgresDriver) QueryPrimaryKeys(db *sql.DB, schema, tableNames string) []TableConstraints {
	var rows *sql.Rows
	var err error
	if tableNames != "*" {
		rows, err = db.Query(PRIMARY_SQL_IN_TABLES, schema, tableNames)
	} else {
		rows, err = db.Query(PRIMARY_SQL, schema)
	}
	common.PanicOnError(err, "[postgres driver]query primary keys error")
	tableConstraintsSlice := []TableConstraints{}
	EachRow(rows, strings.Split(COLUMNS, ","), func(columns []string, rb []sql.RawBytes) bool {
		oneRow := ToTableConstraints(columns, rb)
		tableConstraintsSlice = append(tableConstraintsSlice, oneRow)
		return true
	})
	return tableConstraintsSlice
}

func ToTableConstraints(columns []string, rb []sql.RawBytes) TableConstraints {
	obj := TableConstraints{}
	if len(columns) == len(rb) {
		for i := range columns {
			switch columns[i] {
			case "constraint_catalog":
				obj.ConstraintCatalog = AsString(rb[i])
			case "constraint_schema":
				obj.ConstraintSchema = AsString(rb[i])
			case "constraint_name":
				obj.ConstraintName = AsString(rb[i])
			case "table_catalog":
				obj.TableCatalog = AsString(rb[i])
			case "table_schema":
				obj.TableSchema = AsString(rb[i])
			case "table_name":
				obj.TableName = AsString(rb[i])
			case "constraint_type":
				obj.ConstraintType = AsString(rb[i])
			case "is_deferrable":
				obj.IsDeferrable = AsString(rb[i])
			case "initially_deferred":
				obj.InitiallyDeferred = AsString(rb[i])
			}
		}
	}
	return obj
}

func EachRow(rows *sql.Rows, rowNameSameWithQuery []string, visitor QueryRowVisitor) error {
	defer rows.Close()
	cols, err := rows.Columns()
	common.PanicOnError(err, "[postgres driver]query rows.Columns() error")
	vals := make([]sql.RawBytes, len(cols))
	ints := make([]interface{}, len(cols))
	for i := range ints {
		ints[i] = &vals[i]
	}
	for rows.Next() {
		if err := rows.Scan(ints...); err != nil {
			log.Println(err)
			return err
		}
		if continued := visitor(rowNameSameWithQuery, vals); !continued {
			break
		}
	}
	return nil
}

//func ProcessEachRow(rows *sql.Rows, rowNameSameWithQuery []string, visitor QueryRowVisitor) []TableConstraints {
//	rowsToReturn := []TableConstraints{}
//	defer rows.Close()
//	cols, err := rows.Columns()
//	common.PanicOnError(err, "[postgres driver]query primary keys rows.Columns() error")
//	vals := make([]sql.RawBytes, len(cols))
//	ints := make([]interface{}, len(cols))
//	for i := range ints {
//		ints[i] = &vals[i]
//	}
//	for rows.Next() {
//		if err := rows.Scan(ints...); err != nil {
//			log.Println(err)
//			return err
//		}
//		oneRow := visitor(rowNameSameWithQuery, vals)
//		rowsToReturn = append(rowsToReturn, oneRow)
//	}
//	return rowsToReturn
//}
