//package driver
//
//import (
//	"database/sql"
//
//	"strings"
//
//	"fmt"
//
//	"github.com/carbin-gun/project/common"
//	"github.com/carbin-gun/project/database"
//	_ "github.com/lib/pq"
//	"github.com/prometheus/log"
//)
//
//const (
//	COLUMNS               = `constraint_catalog,constraint_schema,constraint_name,table_catalog,table_schema,table_name,constraint_type,is_deferrable,initially_deferred`
//	PRIMARY_SQL           = `select $1 from information_shcema.table_constraints where table_schema = $2`
//	PRIMARY_SQL_IN_TABLES = `select $1 from information_shcema.table_constraints where table_schema = $2 and table_name in $3 and constraint_type='PRIMARY KEY' `
//)
//
//type PostgresDriver struct{}
//
//type ColumnKeyValue struct {
//	Name  string
//	Value interface{}
//}
//type QueryRowVisitor func(columns []string, rb []sql.RawBytes) bool
//
//func (postgresDriver PostgresDriver) Load(dsnString string, schema string, tableNames string) (database.Schema, error) {
//	log.Printf("[postgres driver] is loading,schema:%s,tableNames:%s", schema, tableNames)
//	db, err := sql.Open("postgres", dsnString)
//	common.PanicOnError(err, "[postgres driver] connect error")
//	defer db.Close()
//	if err = db.Ping(); err != nil {
//		common.PanicOnError(err, "[postgres driver] Ping error")
//	}
//	ret := postgresDriver.queryColumns(db, schema, tableNames)
//	log.Printf("[Postgres Driver] Loaded schema data of %d tables from database schema[%s]", len(ret), schema)
//	return ret, nil
//}
//
//func (postgresDriver PostgresDriver) queryColumns(db *sql.DB, schema, tableNames string) database.Schema {
//	ret := make(database.Schema)
//	primaryKeys := postgresDriver.QueryPrimaryKeys(db, schema, tableNames)
//	fmt.Println("primaryKeys:%+v", primaryKeys)
//	return ret
//
//}
//
//func (postgresDriver PostgresDriver) QueryPrimaryKeys(db *sql.DB, schema, tableNames string) []database.TableConstraints {
//	var rows *sql.Rows
//	var err error
//	if tableNames != "*" {
//		rows, err = db.Query(PRIMARY_SQL_IN_TABLES, COLUMNS, schema, tableNames)
//	} else {
//		rows, err = db.Query(PRIMARY_SQL, COLUMNS, schema)
//	}
//	common.PanicOnError(err, "[postgres driver]query primary keys error")
//	tableConstraintsSlice := []database.TableConstraints{}
//	EachRow(rows, strings.Split(COLUMNS, ","), func(columns []string, rb []sql.RawBytes) bool {
//		oneRow := ToTableConstraints(columns, rb)
//		tableConstraintsSlice = append(tableConstraintsSlice, oneRow)
//		return true
//	})
//	return tableConstraintsSlice
//}
//
//func ToTableConstraints(columns []string, rb []sql.RawBytes) database.TableConstraints {
//	obj := database.TableConstraints{}
//	if len(columns) == len(rb) {
//		for i := range columns {
//			switch columns[i] {
//			case "constraint_catalog":
//				obj.ConstraintCatalog = database.AsString(rb[i])
//			case "constraint_schema":
//				obj.ConstraintSchema = database.AsString(rb[i])
//			case "constraint_name":
//				obj.ConstraintName = database.AsString(rb[i])
//			case "table_catalog":
//				obj.TableCatalog = database.AsString(rb[i])
//			case "table_schema":
//				obj.TableSchema = database.AsString(rb[i])
//			case "table_name":
//				obj.TableName = database.AsString(rb[i])
//			case "constraint_type":
//				obj.ConstraintType = database.AsString(rb[i])
//			case "is_deferrable":
//				obj.IsDeferrable = database.AsString(rb[i])
//			case "initially_deferred":
//				obj.InitiallyDeferred = database.AsString(rb[i])
//			}
//		}
//	}
//	return obj
//}
//
//func EachRow(rows *sql.Rows, rowNameSameWithQuery []string, visitor QueryRowVisitor) error {
//	defer rows.Close()
//	cols, err := rows.Columns()
//	common.PanicOnError(err, "[postgres driver]query rows.Columns() error")
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
//		if continued := visitor(rowNameSameWithQuery, vals); !continued {
//			break
//		}
//	}
//	return nil
//}
//
////func ProcessEachRow(rows *sql.Rows, rowNameSameWithQuery []string, visitor QueryRowVisitor) []TableConstraints {
////	rowsToReturn := []TableConstraints{}
////	defer rows.Close()
////	cols, err := rows.Columns()
////	common.PanicOnError(err, "[postgres driver]query primary keys rows.Columns() error")
////	vals := make([]sql.RawBytes, len(cols))
////	ints := make([]interface{}, len(cols))
////	for i := range ints {
////		ints[i] = &vals[i]
////	}
////	for rows.Next() {
////		if err := rows.Scan(ints...); err != nil {
////			log.Println(err)
////			return err
////		}
////		oneRow := visitor(rowNameSameWithQuery, vals)
////		rowsToReturn = append(rowsToReturn, oneRow)
////	}
////	return rowsToReturn
////}
