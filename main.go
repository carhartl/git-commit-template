package main

import (
	"log"
	"os"

	"github.com/carhartl/git-commit-template/internal/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "git-commit-template",
		Usage:   "Set up a useful template for commit messages",
		Version: "v0.2.0",
		Commands: []*cli.Command{
			command.SetTemplateCommand,
			command.UnsetTemplateCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
