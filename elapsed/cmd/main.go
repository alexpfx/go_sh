package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type Record struct {
	Time int64
}

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tag",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:    "file",
				Usage:   "sufixo do arquivo",
				Aliases: []string{"f"},
			},
		},

		Action: func(c *cli.Context) error {
			fname := "./.last_execution" + c.String("file")
			f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0755)

			if err != nil {
				return err
			}

			var inputData Record
			err = binary.Read(f, binary.BigEndian, &inputData)
			if err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("erro de leitura: %s", err)
			}
			f.Close()

			elapsed := time.Since(time.Unix(inputData.Time, 0))
			prefix := c.String("prefix")
			if prefix != "" {
				fmt.Print(prefix + " ")
			}
			fmt.Println(elapsed)

			f, err = os.OpenFile(fname, os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			outputData := Record{
				Time: time.Now().Unix(),
			}

			err = binary.Write(f, binary.BigEndian, outputData)
			if err != nil {
				return fmt.Errorf("erro de escrita %s", err)
			}

			return nil

		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
