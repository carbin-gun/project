package database

import (
	"database/sql"
	"log"

	"strings"

	"fmt"

	"github.com/carbin-gun/project/common"
	_ "github.com/lib/pq"
)

const (
	PRIMARY_KEYS_COLUMNS  = `a.constraint_catalog,a.constraint_schema,a.constraint_name,a.table_catalog,a.table_schema,a.table_name,a.constraint_type,a.is_deferrable,a.initially_deferred,b.column_name`
	TABLE_COLUMNS_COLUMNS = `table_catalog, table_schema, table_name, column_name, ordinal_position, column_default, is_nullable, data_type, character_maximum_length, character_octet_length, numeric_precision, numeric_precision_radix, numeric_scale, datetime_precision, interval_type, interval_precision, character_set_catalog, character_set_schema, character_set_name, collation_catalog, collation_schema, collation_name, domain_catalog, domain_schema, domain_name, udt_catalog, udt_schema, udt_name, scope_catalog, scope_schema, scope_name, maximum_cardinality, dtd_identifier, is_self_referencing, is_identity, identity_generation, identity_start, identity_increment, identity_maximum, identity_minimum, identity_cycle, is_generated, generation_expression, is_updatable`
	//TABLE_CONSTRAINTS_COLUMNS = `constraint_catalog,constraint_schema,constraint_name,table_catalog,table_schema,table_name,constraint_type,is_deferrable,initially_deferred`
	//PRIMARY_SQL               = `select ` + TABLE_CONSTRAINTS_COLUMNS + ` from information_schema.table_constraints where table_schema = $1 and constraint_type='PRIMARY KEY'`
	//PRIMARY_SQL_IN_TABLES     = `select ` + TABLE_CONSTRAINTS_COLUMNS + ` from information_schema.table_constraints where table_schema = $1 and table_name in $2 and constraint_type='PRIMARY KEY' `
	COLUMNS_SQL                = `select ` + TABLE_COLUMNS_COLUMNS + ` from information_schema.columns where table_schema=$1 order by ordinal_position`
	COLUMNS_SQL_IN_TABLES      = `select ` + TABLE_COLUMNS_COLUMNS + ` from information_schema.columns where table_schema=$1 and table_name in $2 order by ordinal_position`
	PRIMARY_KEYS_SQL           = `select ` + PRIMARY_KEYS_COLUMNS + ` from information_schema.table_constraints a inner join information_schema.key_column_usage b on a.constraint_name = b.constraint_name and a.table_schema=b.table_schema and a.table_name=b.table_name and a.table_schema=$1 and a.constraint_type='PRIMARY KEY' `
	PRIMARY_KEYS_SQL_IN_TABLES = `select ` + PRIMARY_KEYS_COLUMNS + ` from information_schema.table_constraints a inner join information_schema.key_column_usage b on a.constraint_name = b.constraint_name and a.table_schema=b.table_schema and a.table_name=b.table_name and a.table_schema=$1  and table_name in $2 and a.constraint_type='PRIMARY KEY' `
)

type StringSet map[string]struct{}

type PostgresDriver struct{}

type ColumnKeyValue struct {
	Name  string
	Value interface{}
}

//The primary key is built as the pattern 'tableName.columnName'

const PRIMARY_KEY_PATTERN string = "%s.%s"

func buildPrimaryKey(tableName, columnName string) string {
	return fmt.Sprintf(PRIMARY_KEY_PATTERN, tableName, columnName)
}

type QueryRowVisitor func(columns []string, rb []sql.RawBytes) bool

//Load query the table info for the code generation
func (postgresDriver PostgresDriver) Load(dsnString string, schema string, tableNames string) (Schema, error) {
	log.Printf("[postgres driver] is loading,schema:%s,tableNames:%s", schema, tableNames)
	db, err := sql.Open("postgres", dsnString)
	common.PanicOnError(err, "[postgres driver] connect error")
	defer db.Close()
	if err = db.Ping(); err != nil {
		common.PanicOnError(err, "[postgres driver] Ping error")
	}
	ret := postgresDriver.Query(db, schema, tableNames)
	log.Printf("[Postgres Driver] Loaded schema data of %d tables from database schema[%s]", len(ret), schema)
	return ret, nil
}

//GenerateCode generate the codes according to the instance of schema
func (postgresDriver PostgresDriver) GenerateCode(dbName string, schema Schema, templatePath string, targetDir string) {
	GenByDefault(dbName, schema, templatePath, targetDir)
}

func (postgresDriver PostgresDriver) Query(db *sql.DB, schema, tableNames string) Schema {
	ret := make(Schema)
	tableConstraints, primaryKeys := postgresDriver.QueryPrimaryKeys(db, schema, tableNames)
	for _, constraint := range tableConstraints {
		log.Printf("queried table [%s.%s] in schema [%s]\n", constraint.TableCatalog, constraint.TableName, constraint.ConstraintSchema)
	}
	tableColumns := postgresDriver.QueryColumns(db, schema, tableNames)
	for _, item := range tableColumns {
		extra := ""
		if strings.HasPrefix(item.ColumnDefault, "nextval(") {
			extra = "AUTO_INCREMENT"
		}
		schemaItem := Column{
			Schema:       item.TableSchema,
			TableName:    item.TableName,
			ColumnName:   item.ColumnName,
			DefaultValue: item.ColumnDefault,
			DataType:     postgresDriver.ConvertToGoType(item.DataType),
			ColumnKey:    setUpColumnKey(primaryKeys, item.TableName, item.ColumnName),
			Extra:        extra,
		}
		ret[item.TableName] = append(ret[item.TableName], schemaItem)

	}
	return ret

}

func setUpColumnKey(primaries StringSet, tableName, columnName string) IndexInfo {
	pk := buildPrimaryKey(tableName, columnName)
	if _, ok := primaries[pk]; ok {
		index := PrimaryIndex{}
		return index
	} else {
		index := NonIndex{}
		return index
	}
}

/**
convert from database type to golang type
*/
func (p PostgresDriver) ConvertToGoType(dt string) string {
	kFieldTypes := map[string]string{
		"bigint":    "int64",
		"int":       "int",
		"integer":   "int",
		"smallint":  "int",
		"character": "string",
		"text":      "string",
		"timestamp": "time.Time",
		"numeric":   "float64",
		"boolean":   "bool",
	}
	dt = strings.Split(dt, " ")[0]
	if fieldType, ok := kFieldTypes[strings.ToLower(dt)]; !ok {
		return "string"
	} else {
		return fieldType
	}
}

func (postgresDriver PostgresDriver) QueryColumns(db *sql.DB, schema, tableNames string) []TableColumns {
	var rows *sql.Rows
	var err error
	if tableNames != "*" {
		rows, err = db.Query(COLUMNS_SQL_IN_TABLES, schema, tableNames)
	} else {
		rows, err = db.Query(COLUMNS_SQL, schema)
	}
	defer rows.Close()
	common.PanicOnError(err, "[postgres driver]QueryColumns error")
	tableColumnsSlice := []TableColumns{}
	EachRow(rows, strings.Split(TABLE_COLUMNS_COLUMNS, ","), func(columns []string, rb []sql.RawBytes) bool {
		oneRow := ToTableColumns(columns, rb)
		tableColumnsSlice = append(tableColumnsSlice, oneRow)
		return true
	})
	return tableColumnsSlice

}

func (postgresDriver PostgresDriver) QueryPrimaryKeys(db *sql.DB, schema, tableNames string) ([]TableConstraints, StringSet) {
	var rows *sql.Rows
	var err error
	if tableNames != `*` {
		rows, err = db.Query(PRIMARY_KEYS_SQL_IN_TABLES, schema, tableNames)
	} else {
		rows, err = db.Query(PRIMARY_KEYS_SQL, schema)
	}
	defer rows.Close()
	common.PanicOnError(err, "[postgres driver]QueryPrimaryKeys error")
	tableConstraintsSlice := []TableConstraints{}
	primaryKeys := make(StringSet)

	EachRow(rows, strings.Split(PRIMARY_KEYS_COLUMNS, ","), func(columns []string, rb []sql.RawBytes) bool {
		oneRow := ToTableConstraints(columns, rb)
		tableConstraintsSlice = append(tableConstraintsSlice, oneRow)
		primaryKeys[buildPrimaryKey(oneRow.TableName, oneRow.ColumnName)] = struct{}{}
		return true
	})
	return tableConstraintsSlice, primaryKeys
}

//ToTableColumns convert low-level,rawBytes slice to one struct instance for TableColumns
func ToTableColumns(columns []string, rb []sql.RawBytes) TableColumns {
	obj := TableColumns{}
	if len(columns) == len(rb) {
		for i := range columns {
			columnName := strings.TrimSpace(columns[i])
			switch columnName {
			case "table_catalog":
				obj.TableCatalog = AsString(rb[i])
			case "table_schema":
				obj.TableSchema = AsString(rb[i])
			case "table_name":
				obj.TableName = AsString(rb[i])
			case "column_name":
				obj.ColumnName = AsString(rb[i])
			case "ordinal_position":
				obj.OrdinalPosition = AsInt(rb[i])
			case "column_default":
				obj.ColumnDefault = AsString(rb[i])
			case "is_nullable":
				obj.IsNullable = AsString(rb[i])
			case "data_type":
				obj.DataType = AsString(rb[i])
			case "character_maximum_length":
				obj.CharacterMaximumLength = AsInt(rb[i])
			case "character_octet_length":
				obj.CharacterOctetLength = AsInt(rb[i])
			case "numeric_precision":
				obj.NumericPrecision = AsInt(rb[i])
			case "numeric_precision_radix":
				obj.NumericPrecisionRadix = AsInt(rb[i])
			case "numeric_scale":
				obj.NumericScale = AsInt(rb[i])
			case "datetime_precision":
				obj.DatetimePrecision = AsInt(rb[i])
			case "interval_type":
				obj.IntervalType = AsString(rb[i])
			case "interval_precision":
				obj.IntervalPrecision = AsInt(rb[i])
			case "character_set_catalog":
				obj.CharacterSetCatalog = AsString(rb[i])
			case "character_set_schema":
				obj.CharacterSetSchema = AsString(rb[i])
			case "character_set_name":
				obj.CharacterSetName = AsString(rb[i])
			case "collation_catalog":
				obj.CollationCatalog = AsString(rb[i])
			case "collation_schema":
				obj.CollationSchema = AsString(rb[i])
			case "collation_name":
				obj.CollationName = AsString(rb[i])
			case "domain_catalog":
				obj.DomainCatalog = AsString(rb[i])
			case "domain_schema":
				obj.DomainSchema = AsString(rb[i])
			case "domain_name":
				obj.DomainName = AsString(rb[i])
			case "udt_catalog":
				obj.UdtCatalog = AsString(rb[i])
			case "udt_schema":
				obj.UdtSchema = AsString(rb[i])
			case "udt_name":
				obj.UdtName = AsString(rb[i])
			case "scope_catalog":
				obj.ScopeCatalog = AsString(rb[i])
			case "scope_schema":
				obj.ScopeSchema = AsString(rb[i])
			case "scope_name":
				obj.ScopeName = AsString(rb[i])
			case "maximum_cardinality":
				obj.MaximumCardinality = AsInt(rb[i])
			case "dtd_identifier":
				obj.DtdIdentifier = AsString(rb[i])
			case "is_self_referencing":
				obj.IsSelfReferencing = AsString(rb[i])
			case "is_identity":
				obj.IsIdentity = AsString(rb[i])
			case "identity_generation":
				obj.IdentityGeneration = AsString(rb[i])
			case "identity_start":
				obj.IdentityStart = AsString(rb[i])
			case "identity_increment":
				obj.IdentityIncrement = AsString(rb[i])
			case "identity_maximum":
				obj.IdentityMaximum = AsString(rb[i])
			case "identity_minimum":
				obj.IdentityMinimum = AsString(rb[i])
			case "identity_cycle":
				obj.IdentityCycle = AsString(rb[i])
			case "is_generated":
				obj.IsGenerated = AsString(rb[i])
			case "generation_expression":
				obj.GenerationExpression = AsString(rb[i])
			case "is_updatable":
				obj.IsUpdatable = AsString(rb[i])
			}
		}
	}
	return obj
}

//ToTableConstraints convert low-level,rawBytes slice to one struct instance for TableCOnstraints
func ToTableConstraints(columns []string, rb []sql.RawBytes) TableConstraints {
	obj := TableConstraints{}
	if len(columns) == len(rb) {
		for i := range columns {
			prefixLen := 2 //there will be "a." or "b." as the prefix,so we should remove it.
			columnName := strings.TrimSpace(columns[i])[prefixLen:]
			loopVal := AsString(rb[i])
			switch columnName {
			case "constraint_catalog":
				obj.ConstraintCatalog = loopVal
			case "constraint_schema":
				obj.ConstraintSchema = loopVal
			case "constraint_name":
				obj.ConstraintName = loopVal
			case "table_catalog":
				obj.TableCatalog = loopVal
			case "table_schema":
				obj.TableSchema = loopVal
			case "table_name":
				obj.TableName = loopVal
			case "constraint_type":
				obj.ConstraintType = loopVal
			case "is_deferrable":
				obj.IsDeferrable = loopVal
			case "initially_deferred":
				obj.InitiallyDeferred = loopVal
			case "column_name":
				obj.ColumnName = loopVal

			}
		}
	}
	return obj
}

func EachRow(rows *sql.Rows, rowNameSameWithQuery []string, visitor QueryRowVisitor) error {
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
