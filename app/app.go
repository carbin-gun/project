package app

import (
	"github.com/carbin-gun/project/cmd"
	"github.com/codegangsta/cli"
)

const VERSION = "0.0.1"
const AUTHOR = "carbin-gun"
const EMAIL = "cilendeng@gmail.com"

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "project"
	app.Usage = "project,generate the whole project code for you !"
	app.Version = VERSION
	app.Author = AUTHOR
	app.Email = EMAIL

	//app commands
	app.Commands = []cli.Command{
		cmd.New,
	}
	return app
}
