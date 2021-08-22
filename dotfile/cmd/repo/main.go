package main

import (
	"fmt"
	"github.com/alexpfx/go_sh/common/util"
	"github.com/alexpfx/go_sh/dotfile/internal/dotfile"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const git = "/usr/bin/git"
const defaultAlias = "cfg"
const defaultRepo = "https://github.com/alexpfx/linux_dotfiles.git"
const defaultGitdir = ".cfg"

var version = "development"
var buildTime = "N\\A"

func main() {
	homeDir, err := os.UserHomeDir()
	util.CheckFatal(err, "")

	app := &cli.App{
		Name: "cfg_repo",
		Usage: "init a repository",
		Commands: []*cli.Command{
			{
				Name: "version", Usage: "print build version and exit",
				Action: func(context *cli.Context) error {
					printVersionAndExit()
					return nil
				},
			},
			{
				Name: "init", Usage: "init a repo",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "alias", Aliases: []string{"a"}, Usage: "command alias", Value: defaultAlias},
					&cli.BoolFlag{Name: "force", Aliases: []string{"f"}, Usage: "remove ditDir if it exinitCmd.BoolVar already exists", Value: false},
					&cli.StringFlag{Name: "repository", Aliases: []string{"r"}, Usage: "repository", Value: defaultRepo},
					&cli.StringFlag{Name: "gitDir", Aliases: []string{"d"}, Usage: "git dir", Value: defaultGitdir},
					&cli.StringFlag{Name: "workTree", Aliases: []string{"t"}, Usage: "workTree", Value: homeDir},
				},
				Action: func(c *cli.Context) error {
					initRepoCmd(
						c.String("repo"),
						c.String("gitDir"),
						c.String("workTree"),
						c.String("alias"),
						c.Bool("force"),
					)

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func initRepoCmd(repo, gitDir, workTree, alias string, force bool) {
	conf := dotfile.Config{
		WorkTree: workTree,
		GitDir:   gitDir,
	}

	if force && util.DirExists(gitDir) {
		err := os.RemoveAll(gitDir)
		util.CheckFatal(err, "cannot remove gitDir")
	}

	_, serr, err := util.ExecCmd(git, []string{"clone", "--bare", repo, gitDir})
	util.CheckFatal(err, serr)

	aliasArgs := []string{
		"--git-dir=" + conf.GitDir + "/",
		"--work-tree=" + conf.WorkTree,
	}

	_, serr, err = util.ExecCmd(git, append(aliasArgs, "config", "--local", "status.showUntrackedFiles", "no"))

	dotfile.WriteConfig(alias, &conf)

	checkout(alias, aliasArgs, workTree)
}

func printVersionAndExit() {
	fmt.Printf("	Version: %s\n	Build time: %s", version, buildTime)
	os.Exit(0)
}

func checkout(alias string, aliasArgs []string, workTree string) {
	var existUntracked []string
	_, serr, err := util.ExecCmd(git, append(aliasArgs, "checkout"))

	if err != nil {
		existUntracked = util.ParseExistUntracked(workTree, serr)
		if len(existUntracked) == 0 {
			util.CheckFatal(err, err.Error())
		}

		dotfile.BackupFiles(fmt.Sprintf(".%s%s_bkp/", workTree, alias), existUntracked)

		for _, untracked := range existUntracked {
			os.RemoveAll(untracked)
		}

		_, serr, err = util.ExecCmd(git, append(aliasArgs, "checkout"))
		util.CheckFatal(err, serr)
	}
}

// 113
