package command

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/urfave/cli/v2"
)

const (
	templateFile = ".git/.gitmessage.txt"
)

func setupClearTest() func() {
	currentDir, _ := os.Getwd()
	dir, err := os.MkdirTemp(os.TempDir(), "test")
	if err != nil {
		panic(err)
	}

	_ = os.Chdir(dir)
	_ = exec.Command("git", "init").Run()

	return func() {
		_ = os.Chdir(currentDir)
		_ = os.RemoveAll(dir)
	}
}

func TestRemoveTemplateFile(t *testing.T) {
	teardown := setupClearTest()
	defer teardown()

	var err error

	_, err = os.Create(templateFile)
	if err != nil {
		panic(err)
	}

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{UnsetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = UnsetTemplateCommand.Run(c, []string{"unset"}...)

	// If the respective file is no longer present, os.Stat returns an error, thus we must expect an error here
	// if the file has been correctly removed!
	_, err = os.Stat(templateFile)
	if err == nil {
		t.Errorf("Expected %s to not exist", templateFile)
	}
}

func TestUnsetCommitTemplateGitConfig(t *testing.T) {
	teardown := setupClearTest()
	defer teardown()

	var err error

	err = exec.Command("git", "config", "--local", "commit.template", templateFile).Run()
	if err != nil {
		panic(err)
	}

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{UnsetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = UnsetTemplateCommand.Run(c, []string{"unset"}...)

	// If the respective config is not set, git exits with 1, thus we must expect an error here
	// if the config has been correctly removed!
	_, err = exec.Command("git", "config", "--local", "commit.template").Output()
	if err == nil {
		t.Errorf("Git config not unset")
	}
}

func TestUnsetOutsideOfGitRepo(t *testing.T) {
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

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{UnsetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	err = UnsetTemplateCommand.Run(c, []string{"unset"}...)
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
