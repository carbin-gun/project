package database

type MysqlDriver struct{}

func (mysql MysqlDriver) Load(dsnString string, schema string, tableNames string) (Schema, error) {
	ret := make(Schema)
	return ret, nil
}

//GenerateCode generate the codes according to the instance of schema
func (mysql MysqlDriver) GenerateCode(dbName string, schema Schema, templatePath string, targetDir string) {
	GenByDefault(dbName, schema, templatePath, targetDir)
}
