package cmd

import (
	"bytes"
	"fmt"
	"go/build"
	"os/exec"
	"path/filepath"

	"os"

	"github.com/carbin-gun/project/common"
	"github.com/codegangsta/cli"
	"github.com/revel/revel"
)

var New = cli.Command{
	Name:      "new",
	Usage:     "create a new app framework according to the template.",
	ShortName: "N",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "template,t",
			Value: "default",
			Usage: "template the new app will copy from",
		},
	},
	Action: NewApp,
}

var (

	// go related paths
	gopath  string
	gocmd   string
	srcRoot string

	// revel related paths
	revelPkg     *build.Package
	appPath      string
	appName      string
	basePath     string
	importPath   string
	skeletonPath string
)

func NewApp(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) == 0 {
		fmt.Println("Usage: project new <project>")
		return
	}

	fmt.Println("args:", args)
	// checking and setting go paths
	initGoPaths()
	//
	//// checking and setting application
	setApplicationPath(args)
	//
	//// checking and setting skeleton
	setSkeletonPath(args)
	//
	//// copy files to new app directory
	copyNewAppFiles()
	//
	//// goodbye world
	fmt.Fprintln(os.Stdout, "Your application is ready:\n  ", appPath)
	fmt.Fprintln(os.Stdout, "\nYou can run it with:\n   project run", importPath)
}

// lookup and set Go related variables
func initGoPaths() {
	// lookup go path
	gopath = build.Default.GOPATH
	if gopath == "" {
		panic("Abort: GOPATH environment variable is not set. " +
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
	}

	// set go src path
	srcRoot = filepath.Join(filepath.SplitList(gopath)[0], "src")

	// check for go executable
	var err error
	gocmd, err = exec.LookPath("go")
	if err != nil {
		panic("Go executable not found in PATH.")
	}

}
func setApplicationPath(args []string) {
	var err error
	importPath = args[0]
	if filepath.IsAbs(importPath) {
		common.Errorf("Abort: '%s' looks like a directory.  Please provide a Go import path instead.",
			importPath)
	}

	_, err = build.Import(importPath, "", build.FindOnly)
	if err == nil {
		common.Errorf("Abort: Import path \"%s\" already exists.\n", importPath)
	}

	revelPkg, err = build.Import(revel.REVEL_IMPORT_PATH, "", build.FindOnly)
	if err != nil {
		common.Errorf("Abort: Could not find Revel source code: %s\n", err)
	}

	appPath = filepath.Join(srcRoot, filepath.FromSlash(importPath))
	appName = filepath.Base(appPath)
	basePath = filepath.ToSlash(filepath.Dir(importPath))
	if basePath == "." {
		// we need to remove the a single '.' when
		// the app is in the $GOROOT/src directory
		basePath = ""
	} else {
		// we need to append a '/' when the app is
		// is a subdirectory such as $GOROOT/src/path/to/revelapp
		basePath += "/"
	}
}
func setSkeletonPath(args []string) {
	var err error
	if len(args) == 2 { // user specified
		skeletonName := args[1]
		_, err = build.Import(skeletonName, "", build.FindOnly)
		if err != nil {
			// Execute "go get <pkg>"
			getCmd := exec.Command(gocmd, "get", "-d", skeletonName)
			fmt.Println("Exec:", getCmd.Args)
			getOutput, err := getCmd.CombinedOutput()
			fmt.Println("Exec Result:", string(getOutput), ",error:", err)
			// check getOutput for no buildible string
			bpos := bytes.Index(getOutput, []byte("no buildable Go source files in"))
			if err != nil && bpos == -1 {
				common.Errorf("Abort: Could not find or 'go get' Skeleton  source code: %s\n%s\n", getOutput, skeletonName)
			}
		}
		// use the
		skeletonPath = filepath.Join(srcRoot, skeletonName)

	} else {
		// use the revel default
		skeletonPath = filepath.Join(revelPkg.Dir, "skeleton")
	}
}

func copyNewAppFiles() {
	var err error
	err = os.MkdirAll(appPath, 0777)
	common.PanicOnError(err, "Failed to create directory "+appPath)

	common.CopyDir(appPath, skeletonPath, map[string]interface{}{
		// app.conf
		"AppName":  appName,
		"BasePath": basePath,
	})

	// Dotfiles are skipped by mustCopyDir, so we have to explicitly copy the .gitignore.
	gitignore := ".gitignore"
	common.CopyFile(filepath.Join(appPath, gitignore), filepath.Join(skeletonPath, gitignore))

}
