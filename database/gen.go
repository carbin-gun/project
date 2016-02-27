package database

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"text/template"
)

type CodeGenResult struct {
	TableName string
	Error     error
}

func GenByDefault(dbName string, schema Schema, templatePath string, targetDir string) {
	CreateDirIfNotExist(targetDir)
	genResults := make(chan CodeGenResult)
	template := parseTemplate(templatePath)
	for tableName, columns := range schema {
		//Concurrently generate codes for tables.
		go func(dir, tableName string, tableColumns Table) {
			err := generateModel(dir, dbName, tableName, tableColumns, template)
			genResult := CodeGenResult{TableName: tableName, Error: err}
			genResults <- genResult
		}(targetDir, tableName, columns)

	}
	for i := 0; i < len(schema); i++ {
		result := <-genResults
		if result.Error != nil {
			log.Printf("Error on CodeGen for table [%s], %s", result.TableName, result.Error)
		} else {
			log.Printf("Done for table [%s] ,model file [%s/%s.go]", result.TableName, targetDir, result.TableName)
		}
	}
	close(genResults)
	fmt.Println("\n==============================================\n\n    All Done Successfully,Congratulations !!! \n\n==============================================")
}

func parseTemplate(path string) *template.Template {
	if path == "" {
		return nil
	}
	return template.Must(template.ParseFiles(path))
}

//generateModel will return a go filed named by table name,the content is go model struct and database access code.
func generateModel(targetDir, dbName, tableName string, tableColumns Table, template *template.Template) error {
	modelFile, err := CreateModelFile(targetDir, tableName)
	if err != nil {
		return err

	}
	w := bufio.NewWriter(modelFile)

	defer func() {
		w.Flush()
		modelFile.Close()
	}()
	modelMeta := ModelMeta{
		Pkg:       targetDir,
		Name:      ToCapitalCase(tableName),
		DbName:    dbName,
		TableName: tableName,
		Fields:    make([]ModelField, len(tableColumns)),
		Uniques:   make([]ModelField, 0, len(tableColumns)),
	}
	needTime := false
	for i, col := range tableColumns {
		field := ModelField{
			Name:            ToCapitalCase(col.ColumnName),
			ColumnName:      col.ColumnName,
			Type:            col.DataType,
			JsonMeta:        fmt.Sprintf("`json:\"%s\"`", col.ColumnName),
			IsPrimaryKey:    col.ColumnKey.IsPrimaryKey(),
			IsUniqueKey:     col.ColumnKey.IsUniqueIndex(),
			IsAutoIncrement: strings.ToUpper(col.Extra) == "AUTO_INCREMENT",
			DefaultValue:    col.DefaultValue,
			Extra:           col.Extra,
			Comment:         col.Comment,
		}
		if field.Type == "time.Time" {
			needTime = true
		}
		if field.IsPrimaryKey {
			modelMeta.PrimaryFields = append(modelMeta.PrimaryFields, &field)
		}

		if field.IsUniqueKey {
			modelMeta.Uniques = append(modelMeta.Uniques, field)
		}

		modelMeta.Fields[i] = field
	}
	if err := modelMeta.GenHeader(w, template, needTime); err != nil {
		return fmt.Errorf("[%s] Fail to gen model header, %s", tableName, err)
	}
	if err := modelMeta.GenStruct(w, template); err != nil {
		return fmt.Errorf("[%s] Fail to gen model struct, %s", tableName, err)
	}
	if err := modelMeta.GenObjectApi(w, template); err != nil {
		return fmt.Errorf("[%s] Fail to gen model object api, %s", tableName, err)
	}
	if err := modelMeta.GenQueryApi(w, template); err != nil {
		return fmt.Errorf("[%s] Fail to gen model query api, %s", tableName, err)
	}
	if err := modelMeta.GenManagedObjApi(w, template); err != nil {
		return fmt.Errorf("[%s] Fail to gen model managed objects api, %s", tableName, err)
	}

	return nil
}
