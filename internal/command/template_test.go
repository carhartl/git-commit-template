package command

import (
	"flag"
	"os"

	"github.com/urfave/cli/v2"
)

func ExampleTemplateCommand_dryRun() {
	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", false, "")
	_ = set.Parse([]string{"--dry-run"})
	c := cli.NewContext(app, set, nil)

	_ = TemplateCommand.Run(c, []string{"template", "--dry-run"}...)

	// Output:
	// Subject
	//
	// Some context/description
}
