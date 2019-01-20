package main

import (
	"github.com/sirupsen/logrus"
	rr "github.com/agui2200/roadrunner/cmd/rr/cmd"

	// services (plugins)
	"github.com/agui2200/roadrunner/service/env"
	"github.com/agui2200/roadrunner/service/http"
	"github.com/agui2200/roadrunner/service/rpc"
	"github.com/agui2200/roadrunner/service/static"

	// additional commands and debug handlers
	_ "github.com/agui2200/roadrunner/cmd/rr/http"
)

func main() {
	rr.Container.Register(env.ID, &env.Service{})
	rr.Container.Register(rpc.ID, &rpc.Service{})
	rr.Container.Register(http.ID, &http.Service{})
	rr.Container.Register(static.ID, &static.Service{})

	rr.Logger.Formatter = &logrus.TextFormatter{ForceColors: true}

	// you can register additional commands using cmd.CLI
	rr.Execute()
}