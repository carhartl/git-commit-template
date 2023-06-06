package command

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli/v2"
)

type Config struct {
	IssuePrefix string `split_words:"true" default:"#"`
	AuthorFile  string `split_words:"true" default:"$HOME/.git-commit-template-authors"`
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

func findCoAuthor(s string, path string) string {
	if len(s) == 0 {
		return s
	}

	if match, _ := regexp.MatchString("^.+ <.+@.+>$", s); match {
		return s
	}

	if fileName, found := strings.CutPrefix(path, "$HOME/"); found {
		home, _ := os.UserHomeDir()
		path = home + "/" + fileName
	}

	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(s)) {
			return scanner.Text()
		}
	}
	return ""
}

var SetTemplateCommand = &cli.Command{
	Name:  "set",
	Usage: "Set up template for commit message",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "dry-run", Aliases: []string{"d"}, Usage: "Print template to stdout", Value: false},
		&cli.StringFlag{Name: "issue", Aliases: []string{"i"}, Usage: "Issue reference to add to template", Value: ""},
		&cli.StringFlag{Name: "pair", Aliases: []string{"p"}, Usage: "Co-author to add to template", Value: ""},
	},
	Before: GitCheck,
	Action: func(c *cli.Context) error {
		var cfg Config
		err := envconfig.Process("GIT_COMMIT_TEMPLATE", &cfg)
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
		}{issueRef(c.String("issue"), cfg.IssuePrefix), findCoAuthor(c.String("pair"), cfg.AuthorFile)})

		if !c.Bool("dry-run") {
			err := exec.Command("git", "config", "--local", "commit.template", ".git/.gitmessage.txt").Run()
			if err != nil {
				return cli.Exit("Updating Git config failed", 1)
			}
		}

		return nil
	},
}
