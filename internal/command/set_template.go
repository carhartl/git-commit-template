// TODO: read pairs from config with fuzzy matching
package command

import (
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/urfave/cli/v2"
)

type EnvConfig struct {
	IssuePrefix string `env:"GIT_COMMIT_TEMPLATE_ISSUE_PREFIX" env-default:"#"`
}

var message = template.Must(
	template.New("message").Parse("Subject (keep under 50 characters)\n\nContext/description (what and why){{if .Issue}}\n\nAddresses: {{.Issue}}{{end}}{{if .Pair}}\n\nCo-authored-by: {{.Pair}}{{end}}"),
)

func issueRef(s string, prefix string) string {
	if strings.HasPrefix(s, prefix) || len(s) == 0 {
		return s
	}
	return prefix + s
}

var SetTemplateCommand = &cli.Command{
	Name:  "set",
	Usage: "Set up template for commit message",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "dry-run", Aliases: []string{"d"}, Usage: "Print template to stdout", Value: false},
		&cli.StringFlag{Name: "issue-ref", Aliases: []string{"i"}, Usage: "Issue reference to add to template", Value: ""},
		&cli.StringFlag{Name: "pair", Aliases: []string{"p"}, Usage: "Co-author to add to template", Value: ""},
	},
	Before: GitCheck,
	Action: func(c *cli.Context) error {
		var cfg EnvConfig
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return cli.Exit("Reading environment variables failed", 1)
		}

		var f *os.File
		if c.Bool("dry-run") {
			f = os.Stdout
		} else {
			var err error
			f, err = os.Create(".git/.gitmessage.txt")
			if err != nil {
				return cli.Exit("Saving template file failed", 1)
			}
			defer f.Close()
		}

		_ = message.Execute(f, struct {
			Issue string
			Pair  string
		}{issueRef(c.String("issue-ref"), cfg.IssuePrefix), c.String("pair")})

		if !c.Bool("dry-run") {
			err := exec.Command("git", "config", "--local", "commit.template", ".git/.gitmessage.txt").Run()
			if err != nil {
				return cli.Exit("Updating Git config failed", 1)
			}
		}

		return nil
	},
}
