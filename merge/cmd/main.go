package main

import (
	"encoding/json"
	"github.com/alexpfx/go_sh/common/util"
	"github.com/alexpfx/go_sh/merge"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	var lastNumberRegex = regexp.MustCompile("[0-9]+")

	configDir, err := os.UserConfigDir()
	var configArg string
	var input string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "configArg",
				Aliases:     []string{"c"},
				Usage:       "caminho do arquivo de configuração",
				FilePath:    filepath.Join(configDir, "go_sh/merge/configArg.json"),
				TakesFile:   true,
				Destination: &configArg,
			},
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "input",
				Value:       "",
				Destination: &input,
			},
			&cli.StringFlag{
				Name:    "privateToken",
				Usage:   "personal gitlab private token",
				EnvVars: []string{"PRIVATE_TOKEN"},
			},
		},
		Action: func(c *cli.Context) error {
			if configArg == "" {
				return nil
			}

			var config merge.Config

			json.Unmarshal([]byte(configArg), &config)

			if &input == nil || input == "" {
				input = util.ReadStin()
			}
			if input == "" {
				return nil
			}
			mergeId := lastNumberRegex.FindString(input)

			mr := merge.MrFetch{
				Config: config,
			}

			mr.Fetch(mergeId)

			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
