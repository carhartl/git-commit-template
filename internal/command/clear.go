package command

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

var ClearCommand = &cli.Command{
	Name:  "clear",
	Usage: "Unset current template",
	Action: func(c *cli.Context) error {
		// Ignore error here, file may have been removed manually...
		_ = os.Remove(".git/.gitmessage.txt")
		// Git config likewise...
		_ = exec.Command("git", "config", "--local", "--unset", "commit.template").Run()
		return nil
	},
}
