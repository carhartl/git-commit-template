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
	currentDir, _ := os.Getwd()
	dir, err := os.MkdirTemp(os.TempDir(), "test")
	if err != nil {
		panic(err)
	}

	_ = os.Chdir(dir)
	_ = exec.Command("git", "init").Run()

	return func() {
		_ = os.Chdir(currentDir)
		os.Unsetenv("GIT_COMMIT_TEMPLATE_ISSUE_PREFIX")
		defer os.RemoveAll(dir)
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
	expected := "Subject (keep under 50 characters)\n\nContext/description (what and why)"
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
	os.RemoveAll(dir)
}

func ExampleSetTemplateCommand_dryRunBasic() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	_ = set.Parse([]string{"--dry-run"})
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run"}...)

	// Output:
	// Subject (keep under 50 characters)
	//
	// Context/description (what and why)
}

func ExampleSetTemplateCommand_dryRunWithFullIssueRef() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	set.String("issue-ref", "#123", "")
	_ = set.Parse([]string{"--dry-run", "--issue-ref", "#123"})
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue-ref", "#123"}...)

	// Output:
	// Subject (keep under 50 characters)
	//
	// Context/description (what and why)
	//
	// Addresses: #123
}

func ExampleSetTemplateCommand_dryRunWithIssueRefNumber() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	set.String("issue-ref", "#123", "")
	_ = set.Parse([]string{"--dry-run", "--issue-ref", "#123"})
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue-ref", "123"}...)

	// Output:
	// Subject (keep under 50 characters)
	//
	// Context/description (what and why)
	//
	// Addresses: #123
}

func ExampleSetTemplateCommand_dryRunWithIssuePrefixConfig() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	set.String("issue-ref", "#123", "")
	_ = set.Parse([]string{"--dry-run", "--issue-ref", "#123"})
	c := cli.NewContext(app, set, nil)
	os.Setenv("GIT_COMMIT_TEMPLATE_ISSUE_PREFIX", "FOO-")

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue-ref", "123"}...)

	// Output:
	// Subject (keep under 50 characters)
	//
	// Context/description (what and why)
	//
	// Addresses: FOO-123
}

func ExampleSetTemplateCommand_dryRunWithCoAuthor() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	set.String("pair", "John Doe <john@example.com>", "")
	_ = set.Parse([]string{"--dry-run", "--pair", "John Doe <john@example.com>"})
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--pair", "John Doe <john@example.com>"}...)

	// Output:
	// Subject (keep under 50 characters)
	//
	// Context/description (what and why)
	//
	// Co-authored-by: John Doe <john@example.com>
}

func ExampleSetTemplateCommand_dryRunWithAll() {
	teardown := setupSetTemplateTest()
	defer teardown()

	app := &cli.App{Writer: os.Stdout, Commands: []*cli.Command{SetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	set.Bool("dry-run", true, "")
	set.String("issue-ref", "#123", "")
	set.String("pair", "John Doe <john@example.com>", "")
	_ = set.Parse([]string{"--dry-run", "--issue-ref", "#123", "--pair", "John Doe <john@example.com>"})
	c := cli.NewContext(app, set, nil)

	_ = SetTemplateCommand.Run(c, []string{"set", "--dry-run", "--issue-ref", "#123", "--pair", "John Doe <john@example.com>"}...)

	// Output:
	// Subject (keep under 50 characters)
	//
	// Context/description (what and why)
	//
	// Addresses: #123
	//
	// Co-authored-by: John Doe <john@example.com>
}
