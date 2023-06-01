package command

import (
	"os"
	"text/template"

	"github.com/urfave/cli/v2"
)

var message = template.Must(
	template.New("message").Parse("Subject\n\nSome context/description\n{{if . }}\nAddresses: {{.}}\n{{end}}"),
)

var TemplateCommand = &cli.Command{
	Name:  "template",
	Usage: "Set up template for commit message",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "dry-run", Aliases: []string{"d"}, Usage: "Print template to stdout", Value: false},
		&cli.StringFlag{Name: "issue-ref", Aliases: []string{"i"}, Usage: "Issue reference to add to template", Value: ""},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("dry-run") {
			_ = message.Execute(os.Stdout, c.String("issue-ref"))
		}

		return nil
	},
}
