package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "git-commit-template",
		Usage: "Set up a useful template for commit messages",
		Commands: []*cli.Command{
			SetTemplateCommand,
			UnsetTemplateCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
