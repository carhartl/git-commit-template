package command

import (
	"flag"
	"os"

	"github.com/urfave/cli/v2"
)

func ExampleTemplateCommand_dryRunBasic() {
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

func ExampleTemplateCommand_dryRunWithIssueRef() {
	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", false, "")
	set.String("issue-ref", "#123", "")
	_ = set.Parse([]string{"--dry-run", "--issue-ref", "#123"})
	c := cli.NewContext(app, set, nil)

	_ = TemplateCommand.Run(c, []string{"template", "--dry-run", "--issue-ref", "#123"}...)

	// Output:
	// Subject
	//
	// Some context/description
	//
	// Addresses: #123
}
