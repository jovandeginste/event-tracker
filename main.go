package main

import (
	appassets "github.com/jovandeginste/event-tracker/assets"
	"github.com/jovandeginste/event-tracker/pkg/app"
)

var (
	gitRef     = "0.0.0-dev"
	gitRefName = "local"
	gitRefType = "local"
	gitCommit  = "local"
	buildTime  = "now"
)

func main() {
	a := app.NewApp(app.Version{
		BuildTime: buildTime,
		Ref:       gitRef,
		RefName:   gitRefName,
		RefType:   gitRefType,
		Sha:       gitCommit,
	})
	a.Config.Bind = ":8080"
	a.Config.Logging = true
	a.Assets = appassets.FS()

	if err := a.Initialize(); err != nil {
		panic(err)
	}

	if err := a.Serve(); err != nil {
		panic(err)
	}
}
