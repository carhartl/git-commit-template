package main

import (
	"log"
	"os"

	"github.com/carhartl/git-commit-message/internal/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "git-commit-message",
		Commands: []*cli.Command{
			command.TemplateCommand,
			command.ClearCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
