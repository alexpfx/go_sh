package main

import (
	"fmt"
	"github.com/alexpfx/go_sh/common/util"
	"github.com/urfave/cli/v2"
	"os"

	"dotfile"

	"log"
)

const git = "/usr/bin/git"
const defaultRepo = "https://github.com/alexpfx/linux_dotfiles.git"
const defaultGitdir = ".cfg"
const defaultAlias = "cfg"

var version = "development"
var buildTime = "N\\A"

func main() {
	homeDir, err := os.UserHomeDir()
	util.CheckFatal(err, "")

	app := &cli.App{
		Name: "cfg", Usage: "cfg",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "alias", Aliases: []string{"a"}, Usage: "command alias", Value: defaultAlias},
			&cli.BoolFlag{Name: "version", Aliases: []string{"v"}, Usage: "print version and exit", Value: false},
		},
		Action: func(c *cli.Context) error {

			if c.Bool("version") {
				printVersionAndExit()
			}
			conf := dotfile.LoadConfig(c.String("alias"))
			tail := c.Args().Slice()


			aliasArgs := []string{
				"--git-dir=" + conf.GitDir + "/",
				"--work-tree=" + conf.WorkTree,
			}

			if len(tail) == 0 {
				return nil
			}
			out, stderr, err := util.ExecCmd(git, append(aliasArgs, tail...))
			util.CheckFatal(err, stderr)
			fmt.Println(out)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "update",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "gitDir", Aliases: []string{"d"}, Usage: "git dir", Value: defaultGitdir},
					&cli.StringFlag{Name: "workTree", Aliases: []string{"t"}, Usage: "workTree", Value: homeDir},
				},
				Action: func(c *cli.Context) error {
					gitDir := c.String("gitDir")
					workTree := c.String("workTree")
					alias := c.String("alias")

					fmt.Println(gitDir)
					fmt.Println(workTree)
					fmt.Println(alias)
					checkArgs(gitDir, workTree, alias)
					conf := dotfile.Config{
						WorkTree: workTree,
						GitDir:   gitDir,
					}
					dotfile.WriteConfig(alias, &conf)

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

func checkArgs(args ...string) {
	for _, s := range args {
		if s == "" {
			log.Fatal("all parameters must be provided")
		}
	}

}
func printVersionAndExit() {
	fmt.Printf("	Version: %s\n	Build time: %s", version, buildTime)
	os.Exit(0)
}

//72
