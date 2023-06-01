package command

import (
	"flag"
	"os"

	"github.com/urfave/cli/v2"
)

func ExampleTemplateCommand_dryRunBasic() {
	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
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
	set.Bool("dry-run", true, "")
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

func ExampleTemplateCommand_dryRunWithCoAuthor() {
	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	set.String("pair", "John Doe <john@example.com>", "")
	_ = set.Parse([]string{"--dry-run", "--pair", "John Doe <john@example.com>"})
	c := cli.NewContext(app, set, nil)

	_ = TemplateCommand.Run(c, []string{"template", "--dry-run", "--pair", "John Doe <john@example.com>"}...)

	// Output:
	// Subject
	//
	// Some context/description
	//
	// Co-authored-by: John Doe <john@example.com>
}

func ExampleTemplateCommand_dryRunWithAll() {
	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	set.String("issue-ref", "#123", "")
	set.String("pair", "John Doe <john@example.com>", "")
	_ = set.Parse([]string{"--dry-run", "--issue-ref", "#123", "--pair", "John Doe <john@example.com>"})
	c := cli.NewContext(app, set, nil)

	_ = TemplateCommand.Run(c, []string{"template", "--dry-run", "--issue-ref", "#123", "--pair", "John Doe <john@example.com>"}...)

	// Output:
	// Subject
	//
	// Some context/description
	//
	// Addresses: #123
	//
	// Co-authored-by: John Doe <john@example.com>
}
