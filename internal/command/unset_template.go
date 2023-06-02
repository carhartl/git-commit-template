package command

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

var UnsetTemplateCommand = &cli.Command{
	Name:  "unset",
	Usage: "Unset current template",
	Action: func(c *cli.Context) error {
		_, err := os.Stat(".git")
		if err != nil {
			if os.IsNotExist(err) {
				return cli.Exit("Must run in Git repository", 1)
			} else {
				return cli.Exit("Could not detect Git repository", 1)
			}
		}

		// Ignore error here, file may have been removed manually...
		_ = os.Remove(".git/.gitmessage.txt")
		// Git config likewise...
		_ = exec.Command("git", "config", "--local", "--unset", "commit.template").Run()
		return nil
	},
}
