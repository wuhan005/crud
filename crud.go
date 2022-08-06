package main

import (
	"os"

	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/crud/internal/cmd"
)

func main() {
	defer log.Stop()
	if err := log.NewConsole(); err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "crud"
	app.Commands = []*cli.Command{
		cmd.Generate,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to run application: %v", err)
	}
}
