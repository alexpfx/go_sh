package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_sh/cbfix"
	"github.com/alexpfx/go_sh/common/util"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	configDir, err := os.UserConfigDir()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "rules",
				Aliases:   []string{"r"},
				Usage:     "caminho do arquivo json com as regras",
				FilePath:  filepath.Join(configDir, "go_sh/cbfix/rules.json"),
				TakesFile: true,
			},
			&cli.StringFlag{
				Name: "input", Aliases: []string{"i"}, Usage: "especifica os valores de entrada",
			},

			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "debug",
			},
			/*&cli.BoolFlag{
				Name:    "exit_on_first",
				Aliases: []string{"x"},
				Usage:   "exits on first match",
				Value:   true,
			},*/
		},
		Action: func(c *cli.Context) error {
			var input string
			input = c.String("input")
			if input == "" {
				input = util.ReadStin()
			}
			if input == "" {
				return nil
			}

			rStr := c.String("rules")

			rules := make([]cbfix.Rule, 0)

			err = json.Unmarshal([]byte(rStr), &rules)
			if err != nil {
				log.Fatal(err)
			}

			found := false
			for _, rule := range rules {
				rx := regexp.MustCompile(rule.Copy)
				str := rx.FindString(input)
				if str == "" {
					continue
				}
				rx = regexp.MustCompile(rule.Match)
				replacedStr := rx.ReplaceAllString(str, rule.Replace)
				fmt.Println(replacedStr)
				found = true
				break
			}
			if !found{
				fmt.Println(input)
			}

			return nil
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
