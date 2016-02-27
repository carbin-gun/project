package database

import (
	"bufio"
	"fmt"
	"strings"
	"text/template"
)

type ModelMeta struct {
	Pkg           string
	Name          string
	DbName        string
	TableName     string
	PrimaryFields PrimaryFields
	Fields        []ModelField
	Uniques       []ModelField
}
type PrimaryFields []*ModelField
type ModelField struct {
	Name            string
	ColumnName      string
	Type            string
	JsonMeta        string
	IsPrimaryKey    bool
	IsUniqueKey     bool
	IsAutoIncrement bool
	DefaultValue    string
	Extra           string
	Comment         string
}

func (pf PrimaryFields) FormatObject() func(string) string {
	return func(name string) string {
		// "<Article ArticleId=%v UserId=%v>", obj.ArticleId, obj.UserId
		formats := make([]string, len(pf))
		for i, field := range pf {
			formats[i] = fmt.Sprintf("%s=%%v", field.Name)
		}
		outputs := make([]string, 1+len(pf))
		outputs[0] = fmt.Sprintf("\"<%s %s>\"", name, strings.Join(formats, " "))
		for i, field := range pf {
			outputs[i+1] = fmt.Sprintf("obj.%s", field.Name)
		}
		return strings.Join(outputs, ", ")
	}
}

func (pf PrimaryFields) FormatIncrementId() func() string {
	// obj.Id = {{if eq .PrimaryField.Type "int64"}}id{{else}}{{.PrimaryField.Type}}(id){{end}}
	return func() string {
		for _, field := range pf {
			if field.IsAutoIncrement {
				if field.Type == "int64" {
					return fmt.Sprintf("obj.%s = id", field.Name)
				} else {
					return fmt.Sprintf("obj.%s = %s(id)", field.Name, field.Type)
				}
			}
		}
		return ""
	}
}

func (pf PrimaryFields) FormatFilters() func(string) string {
	// filter := {{.Name}}Objs.Filter{{.PrimaryField.Name}}("=", obj.{{.PrimaryField.Name}})
	return func(name string) string {
		filters := make([]string, len(pf))
		for i, field := range pf {
			if i == 0 {
				filters[i] = fmt.Sprintf("filter := %sObjs.Filter%s(\"=\", obj.%s)", name, field.Name, field.Name)
			} else {
				filters[i] = fmt.Sprintf("filter = filter.And(%sObjs.Filter%s(\"=\", obj.%s))", name, field.Name, field.Name)
			}
		}
		return strings.Join(filters, "\n")
	}
}

func (m ModelMeta) HasAutoIncrementPrimaryKey() bool {
	for _, pField := range m.PrimaryFields {
		if pField.IsAutoIncrement {
			return true
		}
	}
	return false
}

func (m ModelMeta) AllFields() string {
	fields := make([]string, len(m.Fields))
	for i, f := range m.Fields {
		fields[i] = fmt.Sprintf("\"%s\"", f.Name)
	}
	return strings.Join(fields, ", ")
}

func (m ModelMeta) InsertableFields() string {
	fields := make([]string, 0, len(m.Fields))
	for _, f := range m.Fields {
		if f.IsPrimaryKey && f.IsAutoIncrement {
			continue
		}
		autoTimestamp := strings.ToUpper(f.DefaultValue) == "CURRENT_TIMESTAMP" ||
			strings.ToUpper(f.DefaultValue) == "NOW()"
		if f.Type == "time.Time" && autoTimestamp {
			continue
		}
		fields = append(fields, fmt.Sprintf("\"%s\"", f.Name))
	}
	return strings.Join(fields, ", ")
}

func (m ModelMeta) GetInsertableFields() []ModelField {
	fields := make([]ModelField, 0, len(m.Fields))
	for _, f := range m.Fields {
		if f.IsPrimaryKey && f.IsAutoIncrement {
			continue
		}
		autoTimestamp := strings.ToUpper(f.DefaultValue) == "CURRENT_TIMESTAMP" ||
			strings.ToUpper(f.DefaultValue) == "NOW()"
		if f.Type == "time.Time" && autoTimestamp {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

func (m ModelMeta) UpdatableFields() string {
	fields := make([]string, 0, len(m.Fields))
	for _, f := range m.Fields {
		if f.IsPrimaryKey {
			continue
		}
		autoUpdateTime := strings.ToUpper(f.Extra) == "ON UPDATE CURRENT_TIMESTAMP"
		if autoUpdateTime {
			continue
		}
		fields = append(fields, fmt.Sprintf("\"%s\"", f.Name))
	}
	return strings.Join(fields, ", ")
}

func (m ModelMeta) GetUpdatableFields() []ModelField {
	fields := make([]ModelField, 0, len(m.Fields))
	for _, f := range m.Fields {
		if f.IsPrimaryKey {
			continue
		}
		autoUpdateTime := strings.ToUpper(f.Extra) == "ON UPDATE CURRENT_TIMESTAMP"
		if autoUpdateTime {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

func (m ModelMeta) GenHeader(w *bufio.Writer, tmpl *template.Template, importTime bool) error {
	return m.getTemplate(tmpl, "header", HeaderTemplate).Execute(w, map[string]interface{}{
		"DbName":     m.DbName,
		"TableName":  m.TableName,
		"PkgName":    m.Pkg,
		"ImportTime": importTime,
	})
}

func (m ModelMeta) getTemplate(tmpl *template.Template, name string, defaultTmpl *template.Template) *template.Template {
	if tmpl != nil {
		if definedTmpl := tmpl.Lookup(name); definedTmpl != nil {
			return definedTmpl
		}
	}
	return defaultTmpl
}

func (m ModelMeta) GenStruct(w *bufio.Writer, tmpl *template.Template) error {
	return m.getTemplate(tmpl, "struct", StructTemplate).Execute(w, m)
}

func (m ModelMeta) GenObjectApi(w *bufio.Writer, tmpl *template.Template) error {
	return m.getTemplate(tmpl, "obj_api", ObjectAPITemplate).Execute(w, m)
}

func (m ModelMeta) GenQueryApi(w *bufio.Writer, tmpl *template.Template) error {
	return m.getTemplate(tmpl, "query_api", QueryAPITemplate).Execute(w, m)
}

func (m ModelMeta) GenManagedObjApi(w *bufio.Writer, tmpl *template.Template) error {
	return m.getTemplate(tmpl, "managed_api", ManageAPITemplate).Execute(w, m)
}
func (f ModelField) ConverterFuncName() string {
	convertors := map[string]string{
		"int64":     "AsInt64",
		"int":       "AsInt",
		"string":    "AsString",
		"time.Time": "AsTime",
		"float64":   "AsFloat64",
		"bool":      "AsBool",
	}
	if c, ok := convertors[f.Type]; ok {
		return c
	}
	return "AsString"
}
