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

func setupTemplateTest() func() {
	currentDir, _ := os.Getwd()
	dir, err := os.MkdirTemp(os.TempDir(), "test")
	if err != nil {
		panic(err)
	}

	_ = os.Chdir(dir)
	_ = exec.Command("git", "init").Run()

	return func() {
		_ = os.Chdir(currentDir)
		defer os.RemoveAll(dir)
	}
}

func TestCreateTemplateFile(t *testing.T) {
	teardown := setupTemplateTest()
	defer teardown()

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = TemplateCommand.Run(c, []string{"template"}...)

	_, err := os.Stat(".git/.gitmessage.txt")
	if err != nil {
		t.Errorf("Expected .git/.gitmessage.txt to exist, but got %v", err)
	}
}

func TestWriteMessageToTemplateFile(t *testing.T) {
	teardown := setupTemplateTest()
	defer teardown()

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = TemplateCommand.Run(c, []string{"template"}...)

	content, err := os.ReadFile(".git/.gitmessage.txt")
	if err != nil {
		panic(err)
	}
	message := string(content)
	expected := "Subject\n\nSome context/description\n"
	if message != expected {
		t.Errorf("Expected message content does not match")
	}
}

func TestSetCommitTemplateGitConfig(t *testing.T) {
	teardown := setupTemplateTest()
	defer teardown()

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{TemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = TemplateCommand.Run(c, []string{"template"}...)

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
