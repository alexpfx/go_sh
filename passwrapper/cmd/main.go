package main

import (
	"fmt"
	"log"
	"os"
	"passwrapper"

	"github.com/urfave/cli/v2"
)

var letterCharset string
var numberCharset string
var specialCharset string
var length int
var specialCount int
var numberCount int
var upperCaseCount int
var lowerCaseCount int

func main() {
	
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "length",
				Aliases:     []string{"s"},
				Usage:       "tamanho da senha",
				Value:       14,
				Destination: &length,
			},
			&cli.IntFlag{
				Name:        "uppercase",
				Usage:       "quantidade mínima de letras maíusculas",
				Value:       2,
				Destination: &upperCaseCount,
			},
			&cli.IntFlag{
				Name:        "lowercase",
				Usage:       "quantidade mínima de letras minúsculas",
				Value:       2,
				Destination: &lowerCaseCount,
			},
			&cli.IntFlag{
				Name:        "numbers",
				Usage:       "quantidade mínima de números",
				Value:       2,
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}