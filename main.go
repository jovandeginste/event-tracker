package main

import (
	"github.com/jovandeginste/event-tracker/pkgs/app"
)

func main() {
	a := app.New()
	a.Config.Bind = ":8080"
	a.Config.Logging = true

	if err := a.Initialize(); err != nil {
		panic(err)
	}

	if err := a.ImportDir("./testfiles"); err != nil {
		panic(err)
	}

	if err := a.Serve(); err != nil {
		panic(err)
	}
}
