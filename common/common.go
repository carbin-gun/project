package common

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type LoggedError struct{ error }

func Errorf(format string, args ...interface{}) {
	// Ensure the user's command prompt starts on the next line.
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	panic(LoggedError{}) // Panic instead of os.Exit so that deferred will run.
}

func Error(format string) {
	fmt.Fprintf(os.Stderr, format)
	panic(LoggedError{})
}

func PanicOnError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Abort: %s: %s\n", msg, err)
		panic(LoggedError{err})
	}

}
func TrueErrorf(charge bool, format string, args ...interface{}) {
	if charge {
		if len(args) == 0 {
			Errorf(format)
		} else {
			Errorf(format, args)
		}

	}
}

func Empty(field string) bool {
	return strings.TrimSpace(field) == ""
}

func CopyDir(destDir, srcDir string, data map[string]interface{}) error {
	var fullSrcDir string
	// Handle symlinked directories.
	f, err := os.Lstat(srcDir)
	if err == nil && f.Mode()&os.ModeSymlink == os.ModeSymlink {
		fullSrcDir, err = os.Readlink(srcDir)
		if err != nil {
			panic(err)
		}
	} else {
		fullSrcDir = srcDir
	}

	return filepath.Walk(fullSrcDir, func(srcPath string, info os.FileInfo, err error) error {
		// Get the relative path from the source base, and the corresponding path in
		// the dest directory.
		relSrcPath := strings.TrimLeft(srcPath[len(fullSrcDir):], string(os.PathSeparator))
		destPath := path.Join(destDir, relSrcPath)

		// Skip dot files and dot directories.
		if strings.HasPrefix(relSrcPath, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create a subdirectory if necessary.
		if info.IsDir() {
			err := os.MkdirAll(path.Join(destDir, relSrcPath), 0777)
			if !os.IsExist(err) {
				PanicOnError(err, "Failed to create directory")
			}
			return nil
		}

		// If this file ends in ".template", render it as a template.
		if strings.HasSuffix(relSrcPath, ".template") {
			renderTemplate(destPath[:len(destPath)-len(".template")], srcPath, data)
			return nil
		}

		// Else, just copy it over.
		CopyFile(destPath, srcPath)
		return nil
	})
}

func CopyFile(destFilename, srcFilename string) {
	destFile, err := os.Create(destFilename)
	PanicOnError(err, "Failed to create file "+destFilename)

	srcFile, err := os.Open(srcFilename)
	PanicOnError(err, "Failed to open file "+srcFilename)

	_, err = io.Copy(destFile, srcFile)
	PanicOnError(err,
		fmt.Sprintf("Failed to copy data from %s to %s", srcFile.Name(), destFile.Name()))

	err = destFile.Close()
	PanicOnError(err, "Failed to close file "+destFile.Name())

	err = srcFile.Close()
	PanicOnError(err, "Failed to close file "+srcFile.Name())
}
func renderTemplate(destPath, srcPath string, data map[string]interface{}) {
	tmpl, err := template.ParseFiles(srcPath)
	PanicOnError(err, "Failed to parse template "+srcPath)

	f, err := os.Create(destPath)
	PanicOnError(err, "Failed to create "+destPath)

	err = tmpl.Execute(f, data)
	PanicOnError(err, "Failed to render template "+srcPath)

	err = f.Close()
	PanicOnError(err, "Failed to close "+f.Name())
}
