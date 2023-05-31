package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var TemplateCommand = &cli.Command{
	Name:  "template",
	Usage: "Set up template for commit message",
	Action: func(c *cli.Context) error {
		fmt.Println("TODO")
		return nil
	},
}
