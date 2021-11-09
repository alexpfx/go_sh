package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_sh/common/util"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"passwrapper"
	"path/filepath"
	"time"
)

var letterCharset string
var numberCharset string
var specialCharset string
var configArg string
var length int
var specialCount int
var numberCount int
var upperCaseCount int
var lowerCaseCount int

func main() {
	configDir, err := os.UserConfigDir()
	defaultConfigPath := filepath.Join(configDir, "go_sh/passwrapper/config.json")

	var update bool

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "caminho do arquivo de configuração",
				FilePath:    defaultConfigPath,
				TakesFile:   true,
				Destination: &configArg,
			},
			&cli.BoolFlag{
				Name:        "update",
				Aliases:     []string{"u"},
				Usage:       "gera um novo arquivo de configuração e termina a execução",
				Value:       false,
				Destination: &update,
			},
			&cli.IntFlag{
				Name:        "length",
				Aliases:     []string{"s"},
				Usage:       "tamanho da senha",
				Value:       12,
				Destination: &length,
			},
			&cli.IntFlag{
				Name:        "uppercase",
				Usage:       "quantidade mínima de letras maíusculas",
				Value:       1,
				Destination: &upperCaseCount,
			},
			&cli.IntFlag{
				Name:        "lowercase",
				Usage:       "quantidade mínima de letras minúsculas",
				Value:       1,
				Destination: &lowerCaseCount,
			},
			&cli.IntFlag{
				Name:        "numbers",
				Usage:       "quantidade mínima de números",
				Value:       1,
				Destination: &numberCount,
			},
			&cli.IntFlag{
				Name:        "specials",
				Aliases:     []string{"l"},
				Usage:       "quantidade mínima de caracteres especiais",
				Value:       1,
				Destination: &specialCount,
			},
			&cli.StringFlag{
				Name:        "letterCharset",
				Value:       "abcdefghijklmnopqrstuvxzwy",
				Destination: &letterCharset,
			},
			&cli.StringFlag{
				Name:        "numberCharset",
				Value:       "0123456789",
				Destination: &numberCharset,
			},
			&cli.StringFlag{
				Name:        "specialCharset",
				Value:       "@#$:.!*-",
				Destination: &specialCharset,
			},
		},
		Action: func(c *cli.Context) error {
			if update {
				log.Println("updating")
				updateConfig(defaultConfigPath)
				return nil
			}


			var config passwrapper.Config

			err = json.Unmarshal([]byte(configArg), &config)
			util.CheckFatal(err, "Não pode carregar arquivo de configuração")

			pass := passwrapper.Pass{
				Config: passwrapper.Config{
					LetterCharset:  letterCharset,
					NumberCharset:  numberCharset,
					SpecialCharset: specialCharset,
				},
				Upper:   upperCaseCount,
				Lower:   lowerCaseCount,
				Number:  numberCount,
				Special: specialCount,
				Length:  length,
			}

			fmt.Println(pass.Generate())
			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func updateConfig(configPath string) {
	newConfig := passwrapper.Config{
		LetterCharset:  letterCharset,
		NumberCharset:  numberCharset,
		SpecialCharset: specialCharset,
	}

	prettyJson, err := json.MarshalIndent(newConfig, "", "\t")
	util.CheckFatal(err, "")

	if util.FileExists(configPath) {
		backFile := fmt.Sprintf("%s_%s.bkp", configPath, time.Now().Format("20060102_150405.000000"))
		util.MoveFile(configPath, backFile)
	}
	err = ioutil.WriteFile(configPath, prettyJson, 0644)
	util.CheckFatal(err, "não pode escrever em arquivo de configuração")

}
