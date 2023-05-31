package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var ClearCommand = &cli.Command{
	Name:  "clear",
	Usage: "Unset current template",
	Action: func(c *cli.Context) error {
		fmt.Println("TODO")
		return nil
	},
}
