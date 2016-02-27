package database

import (
	"fmt"
	"os"
	"path"
)

func CreateModelFile(dir, tableName string) (*os.File, error) {
	modelFilePath := path.Join(dir, tableName+".go")
	file, err := os.Create(modelFilePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func CreateDirIfNotExist(targetDir string) {
	if _, err := os.Stat(targetDir); err != nil {
		fmt.Println("stat dir error:", err)
		os.MkdirAll(targetDir, os.ModeDir|os.ModePerm)
	}
}

func ToCapitalCase(name string) string {
	// cp___hello_12jiu -> CpHello12Jiu
	data := []byte(name)
	segStart := true
	endPos := 0
	for i := 0; i < len(data); i++ {
		ch := data[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			if segStart {
				if ch >= 'a' && ch <= 'z' {
					ch = ch - 'a' + 'A'
				}
				segStart = false
			} else {
				if ch >= 'A' && ch <= 'Z' {
					ch = ch - 'A' + 'a'
				}
			}
			data[endPos] = ch
			endPos++
		} else if ch >= '0' && ch <= '9' {
			data[endPos] = ch
			endPos++
			segStart = true
		} else {
			segStart = true
		}
	}
	return string(data[:endPos])
}
