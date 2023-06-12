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

const (
	dryRunFlag   = "dry-run"
	issueRefFlag = "issue"
	coAuthorFlag = "pair"
)

type Config struct {
	IssuePrefix string `split_words:"true" default:"#"`
	AuthorFile  string `split_words:"true" default:"$HOME/.git-commit-template-authors"`
}

var message = template.Must(
	template.New("message").Parse("Subject (keep under 50 characters)\n\nContext/description (what and why){{if .Issue}}\n\nAddresses: {{.Issue}}{{end}}{{if .CoAuthors}}\n{{end}}{{range .CoAuthors}}\nCo-authored-by: {{.}}{{end}}"),
)

func issueRef(s string, prefix string) string {
	if strings.HasPrefix(s, prefix) || len(s) == 0 {
		return s
	}
	return prefix + s
}

func findCoAuthors(s []string, path string) []string {
	if fileName, found := strings.CutPrefix(path, "$HOME/"); found {
		home, _ := os.UserHomeDir()
		path = home + "/" + fileName
	}

	f, _ := os.Open(path)
	defer f.Close()

	var coAuthors []string
	for _, author := range s {
		coAuthor := findCoAuthor(author, f)
		if len(coAuthor) > 0 {
			coAuthors = append(coAuthors, coAuthor)
		}
	}

	return coAuthors
}

func findCoAuthor(s string, f *os.File) string {
	if len(s) == 0 {
		return s
	}

	if match, _ := regexp.MatchString("^.+ <.+@.+>$", s); match {
		return s
	}

	if f == nil {
		return ""
	}

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
		&cli.BoolFlag{Name: dryRunFlag, Aliases: []string{"d"}, Usage: "Print template to stdout", Value: false},
		&cli.StringFlag{Name: issueRefFlag, Aliases: []string{"i"}, Usage: "Issue reference to add to template", Value: ""},
		&cli.StringSliceFlag{Name: coAuthorFlag, Aliases: []string{"p"}, Usage: "Co-author(s) to add to template", Value: cli.NewStringSlice()},
	},
	Before: GitCheck,
	Action: func(c *cli.Context) error {
		var cfg Config
		err := envconfig.Process("GIT_COMMIT_TEMPLATE", &cfg)
		if err != nil {
			return cli.Exit("Reading environment variables failed", 1)
		}

		var f *os.File
		if c.Bool(dryRunFlag) {
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
			Issue     string
			CoAuthors []string
		}{issueRef(c.String(issueRefFlag), cfg.IssuePrefix), findCoAuthors(c.StringSlice(coAuthorFlag), cfg.AuthorFile)})

		if !c.Bool(dryRunFlag) {
			err := exec.Command("git", "config", "--local", "commit.template", ".git/.gitmessage.txt").Run()
			if err != nil {
				return cli.Exit("Updating Git config failed", 1)
			}
		}

		return nil
	},
}
