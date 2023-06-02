package command

import (
	"os"
	"os/exec"
	"text/template"

	"github.com/urfave/cli/v2"
)

var message = template.Must(
	template.New("message").Parse("Subject (keep under 50 characters)\n\nContext/description (what and why){{if .Issue}}\n\nAddresses: {{.Issue}}{{end}}{{if .Pair}}\n\nCo-authored-by: {{.Pair}}{{end}}"),
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
		var f *os.File
		if c.Bool("dry-run") {
			f = os.Stdout
		} else {
			var err error
			f, err = os.Create(".git/.gitmessage.txt")
			if err != nil {
				panic(err)
			}
			defer f.Close()
		}

		_ = message.Execute(f, struct {
			Issue string
			Pair  string
		}{c.String("issue-ref"), c.String("pair")})

		if !c.Bool("dry-run") {
			err := exec.Command("git", "config", "--local", "commit.template", ".git/.gitmessage.txt").Run()
			if err != nil {
				panic(err)
			}
		}

		return nil
	},
}
