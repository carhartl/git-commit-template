package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var TemplateCommand = &cli.Command{
	Name:  "template",
	Usage: "Set up template for commit message",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "dry-run", Aliases: []string{"d"}, Usage: "Print template to stdout", Value: false},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("dry-run") {
			fmt.Fprintf(c.App.Writer, "Subject\n\nSome context/description\n")
		}

		return nil
	},
}
