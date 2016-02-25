package database

type MysqlDriver struct{}

func (mysql MysqlDriver) Load(dsnString string, schema string, tableNames string) (Schema, error) {
	ret := make(Schema)
	return ret, nil
}
