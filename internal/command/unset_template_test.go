package command

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/urfave/cli/v2"
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
		defer os.RemoveAll(dir)
	}
}

func TestRemoveTemplateFile(t *testing.T) {
	teardown := setupClearTest()
	defer teardown()

	var err error

	_, err = os.Create(".git/.gitmessage.txt")
	if err != nil {
		panic(err)
	}

	app := &cli.App{Writer: io.Discard, Commands: []*cli.Command{UnsetTemplateCommand}}
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(app, set, nil)

	_ = UnsetTemplateCommand.Run(c, []string{"unset"}...)

	// If the respective file is no longer present, os.Stat returns an error, thus we must expect an error here
	// if the file has been correctly removed!
	_, err = os.Stat(".git/.gitmessage.txt")
	if err == nil {
		t.Errorf("Expected .git/.gitmessage.txt to not exist")
	}
}

func TestUnsetCommitTemplateGitConfig(t *testing.T) {
	teardown := setupClearTest()
	defer teardown()

	var err error

	err = exec.Command("git", "config", "--local", "commit.template", ".git/.gitmessage.txt").Run()
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
