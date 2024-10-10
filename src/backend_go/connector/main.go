package main

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

// main initializes the logger, creates an application with the given module, and runs it.
func main() {
	utils.InitLogger(true)
	app := fx.New(Module)

	app.Run()
}
