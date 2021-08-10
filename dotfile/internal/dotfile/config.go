package dotfile

import (
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_sh/common/util"

	"io"
	"os"
	"path/filepath"
)

type Config struct {
	GitDir string `json:"git_dir"`
	WorkTree string `json:"work_tree"`
}

const appConfigDir = "go_dotfile"
const appConfigTemplate = "config_%s"

func BackupFiles(backupDir string, paths []string) {
	if len(paths) == 0 {
		return
	}

	backupErr := "cannot backup. stopping..."

	err := os.MkdirAll(backupDir, 0700)

	util.CheckFatal(err, backupErr)
	for _, path := range paths {
		if len(path) == 0 {
			continue
		}
		source, err := os.Open(path)
		util.CheckFatal(err, backupErr)

		backupFilePath := filepath.Join(backupDir, path)
		backupFileDir := filepath.Dir(backupFilePath)

		err = os.MkdirAll(backupFileDir, 0700)
		util.CheckFatal(err, backupErr+" "+backupFileDir)

		target, err := os.OpenFile(backupFilePath, os.O_CREATE|os.O_RDWR, 0660)
		util.CheckFatal(err, backupErr+" "+backupFileDir)
		_, err = io.Copy(target, source)
		util.CheckFatal(err, backupErr)

		err = target.Close()
		util.CheckFatal(err, backupErr+" "+backupFileDir)
		err = source.Close()
		util.CheckFatal(err, backupErr+" "+backupFileDir)
	}
}
func LoadConfig(aliasName string) *Config {
	cfgPath := resolveConfigPath(ResolveConfigDir(), aliasName)

	f, err := os.Open(cfgPath)
	util.CheckFatal(err, "")
	defer f.Close()

	dec := json.NewDecoder(f)

	conf := Config{}
	err = dec.Decode(&conf)
	util.CheckFatal(err, "")

	return &conf
}

func ResolveConfigDir() string {
	userCfg, err := os.UserConfigDir()
	util.CheckFatal(err, "Cannot resolve UserConfigDir")
	return filepath.Join(userCfg, appConfigDir)
}

func resolveConfigPath(configDir, aliasName string) string {
	return filepath.Join(configDir, fmt.Sprintf(appConfigTemplate, aliasName))
}
func WriteConfig(aliasName string, conf *Config) {
	cfgDir := ResolveConfigDir()
	err := os.MkdirAll(cfgDir, 0700)
	util.CheckFatal(err, "")

	cfgPath := resolveConfigPath(cfgDir, aliasName)

	f, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_RDWR, 0660)
	util.CheckFatal(err, "")
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	err = enc.Encode(&conf)
	util.CheckFatal(err, "")
}
