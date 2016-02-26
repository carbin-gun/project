package database

import (
	"os"
	"text/template"
)

type CodeGenResult struct {
	Name  string
	Error error
}

func GenByDefault(dbName string, schema Schema, templatePath string, targetDir string) {

	CreateIfNotExist(targetDir)
	genResults := make(chan CodeGenResult)
	template := parseTemplate(templatePath)
	for tableName, columns := range schema {
		//Concurrently generate codes for tables.
		go func(tableName string, tableColumns Table) {
			err := generateModel(dbName, tableName, tableColumns, template)
			genResult := CodeGenResult{Name: tableName, Error: err}
			genResults <- genResult
		}(tableName, columns)

	}

}

func parseTemplate(path string) *template.Template {
	if path == "" {
		return nil
	}
	return template.Must(template.ParseFiles(path))
}

//generateModel will return a go filed named by table name,the content is go model struct and database access code.
func generateModel(dbName, tableName string, tableColumns Table, template *template.Template) error {

	return nil
}

func CreateIfNotExist(targetDir string) {
	if fs, err := os.Stat(targetDir); err != nil || !fs.IsDir() {
		os.Mkdir(targetDir, os.ModeDir|os.ModePerm)
	}
}
