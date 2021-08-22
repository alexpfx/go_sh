package pass

import (
	"encoding/json"
	"fmt"
	"github.com/alexpfx/go_sh/common/util"

	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const GpgExt = ".gpg"
const pass = "pass"

type Backup interface {
	Do()
}

type Restore interface {
	Do()
}

func NewBackup(passwordStore, target string) Backup {
	return backup{
		passwordStore: passwordStore,
		target:        target,
	}
}

func NewRestore(targetDir, backupFilePath string, update bool) Restore {
	return restore{
		targetDir:      targetDir,
		backupFilePath: backupFilePath,
		update:         update,
	}
}

type backup struct {
	passwordStore string
	target        string
}

type restore struct {
	targetDir      string
	backupFilePath string
	update              bool
}

func (b backup) Do() {
	gpgFiles := getGpgFiles(b.passwordStore)

	allPassInfos := make([]GpgPassInfo, 0)
	for _, path := range gpgFiles {
		dec := decrypt(path)
		allPassInfos = append(allPassInfos, GpgPassInfo{
			Password: strings.TrimRight(string(dec), "\n"),
			PassName: getPassName(b.passwordStore, path),
			FilePath: path,
		})
	}

	doBackup(allPassInfos, b.target)
}

type GpgPassInfo struct {
	Password string `json:"password,omitempty"`
	PassName string `json:"pass_name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

func List() {

	out, _, err := util.ExecCmd(pass, []string{"list"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}



func (r restore) Do() {
	str := decrypt(r.backupFilePath)
	passInfo := make([]GpgPassInfo, 0)
	err := json.Unmarshal([]byte(str), &passInfo)
	util.CheckFatal(err, "unmarshal error")

	for _, p := range passInfo {
		restoreTarget := filepath.Join(r.targetDir, p.PassName+GpgExt)

		if util.FileExists(restoreTarget) {
			if !r.update {
				continue
			}
		}

		tmpFile := TempFile()
		_, err = fmt.Fprintln(tmpFile, p.Password)
		util.CheckFatal(err, "error when copy to temp file")


		encryptedTempFile := tmpFile.Name() + GpgExt


		util.MoveFile(encryptedTempFile, restoreTarget)

		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}
}

func encryptFile(file string) string {
	//cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "--passphrase", password, "-c", file}
	cmdArgs := []string{"--pinentry-mode=loopback", "-c", file}
	out, serr, err := util.ExecCmd("gpg", cmdArgs)
	util.CheckFatal(err, serr)
	return out
}

func decrypt(file string) string {
	cmdArgs := []string{
		"-d", "--quiet", "--yes",
		"--compress-algo=none", "--pinentry-mode=loopback",
		file,
	}
	out, _, err := util.ExecCmd("gpg2", cmdArgs)

	util.CheckFatal(err, "não pode descriptografar arquivo.")

	return out
}

func getPassName(baseDir, fullPath string) string {
	return strings.TrimSuffix(strings.TrimPrefix(fullPath, baseDir+string(os.PathSeparator)), filepath.Ext(fullPath))
}

func doBackup(passiInfos []GpgPassInfo, target string) {
	tmpFile := TempFile()

	bs, err := json.MarshalIndent(passiInfos, "", "   ")
	util.CheckFatal(err, string(bs))

	err = ioutil.WriteFile(tmpFile.Name(), bs, 0644)
	util.CheckFatal(err, "")

	encryptFile(tmpFile.Name())

	createdFile := tmpFile.Name() + GpgExt

	backupFile := strings.TrimSuffix(target, GpgExt) + GpgExt

	util.MoveFile(createdFile, backupFile)

	err = os.Remove(tmpFile.Name())
	util.CheckFatal(err, "")
}

func getGpgFiles(baseDir string) []string {
	files := make([]string, 0)

	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if d.Type().IsDir() {
			return nil
		}

		if !strings.EqualFold(filepath.Ext(path), GpgExt) {
			return nil
		}
		files = append(files, path)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}


func TempFile() *os.File {
	f, err := ioutil.TempFile("", "gsh*")
	if err != nil {
		log.Fatal(f)
	}
	return f
}

func v(msg ...interface{}) {
	sb := strings.Builder{}

	for _, m := range msg {
		sb.WriteString(fmt.Sprintf("%s ", m))
	}
	log.Println(sb.String())
}