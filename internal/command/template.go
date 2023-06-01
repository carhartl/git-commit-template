package command

import (
	"os"
	"text/template"

	"github.com/urfave/cli/v2"
)

var message = template.Must(
	template.New("message").Parse("Subject\n\nSome context/description\n{{if .Issue}}\nAddresses: {{.Issue}}\n{{end}}{{if .Pair}}\nCo-authored-by: {{.Pair}}\n{{end}}"),
)

var TemplateCommand = &cli.Command{
	Name:  "template",
	Usage: "Set up template for commit message",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "dry-run", Aliases: []string{"d"}, Usage: "Print template to stdout", Value: false},
		&cli.StringFlag{Name: "issue-ref", Aliases: []string{"i"}, Usage: "Issue reference to add to template", Value: ""},
		&cli.StringFlag{Name: "pair", Aliases: []string{"p"}, Usage: "Co-author to add to template", Value: ""},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("dry-run") {
			_ = message.Execute(os.Stdout, struct {
				Issue string
				Pair  string
			}{c.String("issue-ref"), c.String("pair")})
		} else {
			f, err := os.Create(".git/.gitmessage.txt")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_ = message.Execute(f, struct {
				Issue string
				Pair  string
			}{c.String("issue-ref"), c.String("pair")})
		}
		return nil
	},
}
