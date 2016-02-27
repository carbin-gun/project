package cmd

import (
	"fmt"
	"strings"

	"os"
	"path"

	"github.com/carbin-gun/project/common"
	"github.com/carbin-gun/project/database"
	"github.com/codegangsta/cli"
)

var DATABASE_SUPPORTED = []string{"mysql", "postgres"}

var (
	db       string //database dsn,e.g. root@tcp(127.0.0.1:3306)/test?charset=utf-8
	driver   string //database driver.now supported:mysql,postgres
	schema   string //represent schema for postgres,database name for mysql
	dir      string //directory for where to put the codes into
	tables   string //tables that you can specify to generate code from
	template string //template to generate code from
)

var Database = cli.Command{
	Name:      "database",
	ShortName: "db",
	Aliases:   []string{"d"},
	Usage:     "generate model & database access go files automatically",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "db",
			Value: "",
			Usage: fmt.Sprintf("query from database type,supported [%s] for now", strings.Join(DATABASE_SUPPORTED, ",")),
		},
		cli.StringFlag{
			Name:  "driver",
			Value: "",
			Usage: fmt.Sprintf("query from database type,supported [%s] for now", strings.Join(DATABASE_SUPPORTED, ",")),
		},
		cli.StringFlag{
			Name:  "schema",
			Value: "",
			Usage: " represent Schema for postgresql, but represent database name for mysql",
		},
		cli.StringFlag{
			Name:  "dir",
			Value: "",
			Usage: "Go source code package for generated models code",
		},
		cli.StringFlag{
			Name:  "tables",
			Value: "*",
			Usage: "generate code according to  specific tables by user ",
		},
		cli.StringFlag{
			Name:  "template",
			Value: "",
			Usage: "Passing the template to generate code, or use the default one",
		},
	},
	Action: generateByDatabase,
}

func generateByDatabase(ctx *cli.Context) {
	paramsCheck(ctx)
	specifyDriver, ok := database.SupportedDrivers[strings.ToLower(driver)]
	if !ok {
		common.Errorf("Not support driver,currently support [%s]", strings.Join(DATABASE_SUPPORTED, ","))
	}
	schemaInfo, err := specifyDriver.Load(db, schema, tables)
	common.PanicOnError(err, "[database driver]Load error")
	specifyDriver.GenerateCode(db, schemaInfo, template, dir)

}

func paramsCheck(ctx *cli.Context) {
	db = ctx.String("db")
	common.TrueErrorf(common.Empty(db), "Please provide the db dsn-like thing.e.g. -db=\"root@/blog\" for mysql,-db=\"dbname=blog sslmode=disable\"for postgres")
	driver = ctx.String("driver")
	common.TrueErrorf(common.Empty(driver), "Please provide the driver,you can specify mysql or postgres for now")
	schema = ctx.String("schema")
	common.TrueErrorf(common.Empty(schema), "Please provide schema name,which represents schema in PostgreSQL,but represents database name in MySQL ")
	dir = ctx.String("dir")
	common.TrueErrorf(common.Empty(dir), "Please provide the dir you want to generate codes into")
	tables = ctx.String("tables")
	common.TrueErrorf(common.Empty(tables), "Please provide the tables you want to generate codes from.if all ,provide *,more than one tables,use comma as delimiter")
	template = ctx.String("template")
	if template != "" {
		if file, err := os.Stat(path.Join(dir, template)); err != nil || !file.IsDir() {
			common.PanicOnError(err, "Please provide a existing template file.")
		}
	}
	//fail fast,if it's a not supported driver
	if _, ok := database.SupportedDrivers[strings.ToLower(driver)]; !ok {
		common.Errorf("Not support driver,currently support [%s]", strings.Join(DATABASE_SUPPORTED, ","))
	}
}
