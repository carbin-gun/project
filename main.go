package main

import (
	"os"
	"runtime"

	"github.com/carbin-gun/project/app"
	"github.com/carbin-gun/project/common"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			//if it's a LoggedError,it means the error messages have been logged already.
			if _, ok := err.(common.LoggedError); !ok {
				// This panic was not expected / logged.
				panic(err)
			}
			os.Exit(1)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())
	app.New().Run(os.Args)
}
