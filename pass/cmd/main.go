package main

import (
	"github.com/alexpfx/go_sh/pass/internal/pass"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

const defaultBackupFile = "pass-backup"

func main() {
	log.SetFlags(log.Llongfile)

	homeDir, _ := os.UserHomeDir()
	var defaultPasswordStore = filepath.Join(homeDir, ".password-store/")
	//currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//util.CheckFatal(err, "error when get current directory")

	var defaultRestoreDir = filepath.Join(homeDir, "pass_restore/")
	var debugMode bool
	app := &cli.App{
		Name:  "go_pass",
		Usage: "scripts de linux",
		Commands: []*cli.Command{
			{
				Name:  "pass",
				Usage: "comandos relativos ao comando pass",
				Subcommands: []*cli.Command{
					{
						Name: "list",
						Action: func(c *cli.Context) error {
							pass.List()
							return nil
						},
					},
					{
						Name: "restore",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "backup-file",
								Aliases: []string{"t"},
								Value:   defaultBackupFile + ".gpg",
							},
							&cli.BoolFlag{
								Name:    "update",
								Aliases: []string{"U"},
								Value:   true,
							},
							&cli.StringFlag{
								Name:  "prefix",
								Value: "restored/",
							},
							&cli.StringFlag{
								Name:    "restore-dir",
								Aliases: []string{"d"},
								Value:   defaultRestoreDir,
							},
							&cli.StringFlag{
								Name:    "public-key",
								Aliases: []string{"k"},
								Value:   "",
							},
						},
						Action: func(c *cli.Context) error {
							backupFile := c.String("backup-file")
							targetDir := c.String("restore-dir")
							forceUpdate := c.Bool("update")

							restore := pass.NewRestore(targetDir, backupFile, forceUpdate)
							restore.Do()

							return nil
						},
					},
					{
						Name: "backup",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "password-store",
								Aliases: []string{"d"},
								EnvVars: []string{"PASSWORD_STORE_DIR"},
								Value:   defaultPasswordStore,
							},
							&cli.StringFlag{
								Name:    "backup-file",
								Aliases: []string{"t"},
								Value:   defaultBackupFile,
							},
						},
						Action: func(c *cli.Context) error {
							passwordStore := c.String("password-store")
							target := c.String("backup-file")

							if debugMode {
								log.SetFlags(log.Llongfile)
								//log.SetFlags(0)
								//log.SetOutput(io.Discard)
							}

							backup := pass.NewBackup(passwordStore, target)
							backup.Do()

							return nil
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"D"},
				Value: false,
				Destination: &debugMode,
			},
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
