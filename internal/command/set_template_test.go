package command

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func setupSetTemplateTest() func() {
	cwd, _ := os.Getwd()
	dir, err := os.MkdirTemp(os.TempDir(), "test")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	_ = cmd.Run()

	_ = os.Chdir(dir) // Required for the git check hook of command to pass

	return func() {
		_ = os.Chdir(cwd)
		_ = os.Unsetenv("GIT_COMMIT_TEMPLATE_AUTHOR_FILE")
		_ = os.Unsetenv("GIT_COMMIT_TEMPLATE_ISSUE_PREFIX")
		_ = os.Unsetenv("GIT_COMMIT_TEMPLATE_TEMPLATE")
		_ = os.RemoveAll(dir)
	}
}

func TestCreateTemplateFile(t *testing.T) {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set"}...)

	_, err := os.Stat(".git/.gitmessage.txt")
	if err != nil {
		t.Errorf("Expected .git/.gitmessage.txt to exist, but got %v", err)
	}
}

func TestWriteMessageToTemplateFile(t *testing.T) {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set"}...)

	content, err := os.ReadFile(".git/.gitmessage.txt")
	if err != nil {
		panic(err)
	}
	message := string(content)
	expected := "Subject"
	if message != expected {
		t.Errorf("Expected message content does not match")
	}
}

func TestSetCommitTemplateGitConfig(t *testing.T) {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set"}...)

	config, err := exec.Command("git", "config", "--local", "commit.template").Output()
	if err != nil {
		t.Errorf("Expected Git config not applied")
		return
	}
	expected := ".git/.gitmessage.txt"
	actual := strings.TrimSpace(string(config))
	if actual != expected {
		t.Errorf("Unexpected Git config applied: wanted %q, got %q", expected, actual)
	}
}

func TestCoAuthorWithFuzzyMatchingFileConfig(t *testing.T) {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)
	_ = os.Setenv("GIT_COMMIT_TEMPLATE_AUTHOR_FILE", ".git-commit-template-authors")
	if err := os.WriteFile(".git-commit-template-authors", []byte("John Doe <john@example.com>"), 0666); err != nil {
		panic(err)
	}

	_ = SetTemplateCommand.Run(c, []string{"set", "--pair", "john"}...)

	content, err := os.ReadFile(".git/.gitmessage.txt")
	if err != nil {
		panic(err)
	}
	message := string(content)
	expected := "Co-authored-by: John Doe <john@example.com>"
	if !strings.Contains(message, expected) {
		t.Errorf("Expected Co-authored-by line not part of message")
	}
}

func TestSetTemplateCommandOutsideOfGitRepo(t *testing.T) {
	currentDir, _ := os.Getwd()
	dir, err := os.MkdirTemp(os.TempDir(), "test")
	if err != nil {
		panic(err)
	}

	_ = os.Chdir(dir)

	origExiter := cli.OsExiter
	defer func() {
		cli.OsExiter = origExiter
	}()
	var exitCodeFromOsExiter int
	cli.OsExiter = func(exitCode int) {
		exitCodeFromOsExiter = exitCode
	}

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	err = SetTemplateCommand.Run(c, []string{"set"}...)
	if exitCodeFromOsExiter != 1 {
		t.Errorf("Expected app to exit with 1: %d", err.(cli.ExitCoder).ExitCode())
	}
	expected := "Must run in Git repository"
	if err.Error() != expected {
		t.Errorf("Expected error not observed, wanted %q, got %q", expected, err.Error())
	}

	_ = os.Chdir(currentDir)
	_ = os.RemoveAll(dir)
}

func ExampleSetTemplateCommand_dryRunBasic() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run"}...)

	// Output:
	// Subject
}

func ExampleSetTemplateCommand_dryRunWithFullIssueRef() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue", "#123"}...)

	// Output:
	// Subject
	//
	// Addresses: #123
}

func ExampleSetTemplateCommand_dryRunWithIssueRefNumber() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue", "123"}...)

	// Output:
	// Subject
	//
	// Addresses: #123
}

func ExampleSetTemplateCommand_dryRunWithIssuePrefixConfig() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)
	_ = os.Setenv("GIT_COMMIT_TEMPLATE_ISSUE_PREFIX", "FOO-")

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue", "123"}...)

	// Output:
	// Subject
	//
	// Addresses: FOO-123
}

func ExampleSetTemplateCommand_dryRunWithTemplateConfig() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)
	_ = os.Setenv("GIT_COMMIT_TEMPLATE_TEMPLATE", "{{if .Issue}}{{.Issue}}{{end}} Subject{{if .CoAuthors}}\n\n{{end}}{{range .CoAuthors}}\nCo-authored-by: {{.}}{{end}}")

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue", "123"}...)

	// Output:
	// #123 Subject
}

func ExampleSetTemplateCommand_dryRunWithCoAuthor() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--pair", "John Doe <john@example.com>"}...)

	// Output:
	// Subject
	//
	// Co-authored-by: John Doe <john@example.com>
}

func ExampleSetTemplateCommand_dryRunWithMultipleCoAuthors() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--pair", "John Doe <john@example.com>", "--pair", "Jane Foo <jane@example.com>"}...)

	// Output:
	// Subject
	//
	// Co-authored-by: John Doe <john@example.com>
	// Co-authored-by: Jane Foo <jane@example.com>
}

func ExampleSetTemplateCommand_dryRunWithAll() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue", "#123", "--pair", "John Doe <john@example.com>"}...)

	// Output:
	// Subject
	//
	// Addresses: #123
	//
	// Co-authored-by: John Doe <john@example.com>
}
