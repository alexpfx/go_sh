package main

import (
	"fmt"
	"github.com/alexpfx/go_sh/common/util"
	"github.com/urfave/cli/v2"
	"os"

	"github.com/alexpfx/go_sh/dotfile/internal/dotfile"

	"log"
)

const git = "/usr/bin/git"
const defaultAlias = "cfg"

var version = "development"
var buildTime = "N\\A"

func main() {
	app := &cli.App{
		Name:  "cfg_repo",
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
				Name: "cfg", Usage: "cfg",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "alias", Aliases: []string{"a"}, Usage: "command alias", Value: defaultAlias},
					&cli.StringFlag{Name: "gitDir", Aliases: []string{"d"}, Usage: "git dir"},
					&cli.StringFlag{Name: "workTree", Aliases: []string{"t"}, Usage: "workTree"},
					&cli.BoolFlag{Name: "update_config", Aliases: []string{"u"}, Usage: "write new config file and exit", Value: false},
				},
				Action: func(c *cli.Context) error {
					updateConfig := c.Bool("update_config")
					gitDir := c.String("gitDir")
					workTree := c.String("workTree")
					alias := c.String("alias")

					if updateConfig {
						checkArgs(gitDir, workTree, alias)
						conf := dotfile.Config{
							WorkTree: workTree,
							GitDir:   gitDir,
						}
						dotfile.WriteConfig(alias, &conf)
						return nil
					}

					tail := c.Args().Tail()
					conf := dotfile.LoadConfig(alias)

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
			},
		},
	}

	err := app.Run(os.Args)
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
