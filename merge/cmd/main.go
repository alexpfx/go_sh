package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_sh/common/util"
	"github.com/alexpfx/go_sh/merge"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)


//echo 10697 | merge | jq -r '[.merge.web_url, .merge.author.username, .merge.commit.username, .merge.commit.created_at] | @tsv' | xsel -b
func main() {
	var lastNumberRegex = regexp.MustCompile("[0-9]+")

	configDir, err := os.UserConfigDir()
	var configArg string
	var input string
	var token string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "caminho do arquivo de configuração",
				FilePath:    filepath.Join(configDir, "go_sh/merge/config.json"),
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
				Destination: &token,
				Usage:   "personal gitlab private token",
				EnvVars: []string{"PRIVATE_TOKEN"},
			},
		},
		Action: func(c *cli.Context) error {
			if configArg == "" {
				return nil
			}

			var config merge.Config

			err = json.Unmarshal([]byte(configArg), &config)

			util.CheckFatal(err, "")

			if &input == nil || input == "" {
				input = util.ReadStin()
			}
			if input == "" {
				return nil
			}
			mergeId := lastNumberRegex.FindString(input)
			_, err = strconv.Atoi(mergeId)
			if err != nil {
				fmt.Printf("Id inválido: %s ", mergeId)
			}

			mr := merge.Fetcher{
				Config: config,
				Token: token,
			}

			fetchResult := mr.Fetch(mergeId)
			for _, r := range fetchResult {
				fmt.Println(util.ToJsonStr(r))
			}

			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
