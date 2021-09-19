package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_sh/cb_fix"
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
				Name:     "rules",
				Aliases: []string{"r"},
				Usage:    "regras",
				FilePath: filepath.Join(configDir, "go_sh/cb_fix/rules.json"),
				TakesFile: true,
			},
			&cli.BoolFlag{
				Name:     "debug",
				Aliases: []string{"d"},
				Usage:    "debug",
			},
			&cli.BoolFlag{
				Name:     "exit_on_first",
				Aliases: []string{"x"},
				Usage:    "exits on first match",
			},

		},
		Action: func(c *cli.Context) error {
			input := getInput()
			if input == ""{
				return nil
			}

			rStr := c.String("rules")

			rules := make([]cb_fix.Rule, 0)

			err = json.Unmarshal([]byte(rStr), &rules)
			if err != nil{
				return err
			}

			for _, rule := range rules {
				rx := regexp.MustCompile(rule.Copy)
				found := rx.FindString(input)
				if found == ""{
					continue
				}
				rx = regexp.MustCompile(rule.Match)
				replacedStr := rx.ReplaceAllString(found, rule.Replace)
				fmt.Println(replacedStr)
				break
			}

			return nil
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getInput() string {
	r := bufio.NewReader(os.Stdin)
	text, _ := r.ReadString('\n')
	return text

}
