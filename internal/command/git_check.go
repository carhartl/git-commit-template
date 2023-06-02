package command

import (
	"os"

	"github.com/urfave/cli/v2"
)

var GitCheck = func(cCtx *cli.Context) error {
	_, err := os.Stat(".git")
	if err != nil {
		if os.IsNotExist(err) {
			return cli.Exit("Must run in Git repository", 1)
		} else {
			return cli.Exit("Could not detect Git repository", 1)
		}
	}
	return nil
}
