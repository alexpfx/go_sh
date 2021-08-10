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

type Backup interface {
	Do()
}

type Restore interface {
	Do()
}

func NewBackup(passwordStore, target string) backup {
	return backup{
		passwordStore: passwordStore,
		target:        target,
	}
}

func NewRestore(prefix, passwordStoreDir, backupFilePath string, update bool) restore {
	return restore{
		targetPasswordStore: passwordStoreDir,
		backupFilePath:      backupFilePath,
		prefix:              prefix,
		update:              update,
	}
}

type backup struct {
	passwordStore string
	target        string
}

type restore struct {
	targetPasswordStore string
	backupFilePath      string
	prefix              string
	update              bool
}

func (b backup) Do() {
	v("iniciando backup...")
	v("store: ", b.passwordStore)
	v("target: ", b.target)

	if util.FileExists(b.target) {
		f("arquivo de destino existe")
	}

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
	out, _, err := util.ExecCmd("pass", []string{"list"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func (r restore) Do() {
	v("iniciando restore...")

	str := decrypt(r.backupFilePath)
	v("descriptografando arquivo... ", r.backupFilePath)
	passInfo := make([]GpgPassInfo, 0)

	err := json.Unmarshal([]byte(str), &passInfo)
	util.CheckFatal(err, "")
	for _, p := range passInfo {
		v("extraindo... ", p.PassName)
		targetGpgFile := filepath.Join(r.targetPasswordStore, r.prefix+p.PassName+GpgExt)

		if util.FileExists(targetGpgFile) {
			if !r.update {
				v("Arquivo existe e não será sobreescrito: ", targetGpgFile)
				continue
			}
		}

		tmpFile := TempFile()
		v("criou arquivo temporario... ", tmpFile.Name())


		_, err := fmt.Fprintln(tmpFile, p.Password)
		util.CheckFatal(err, "")

		encryptFile(tmpFile.Name())

		v("movendo arquivo...", targetGpgFile, " ... ", tmpFile.Name()+GpgExt)
		util.MoveFile(targetGpgFile, tmpFile.Name()+GpgExt)

		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}
}

func encryptFile(file string) string {
	//cmdArgs := []string{"gpg", "--pinentry-mode=loopback", "--passphrase", password, "-c", file}
	v("codificando arquivo...", file)
	cmdArgs := []string{"--pinentry-mode=loopback", "-c", file}
	out, _, err := util.ExecCmd("gpg", cmdArgs)
	if err != nil {
		e("erro ao codificar arquivo: ")
		e(err)
	}
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
	v("arquivo temporario criado...", tmpFile.Name())

	bs, err := json.MarshalIndent(passiInfos, "", "   ")
	if err != nil {
		e("erro ao serializar: ")
		f(err.Error())
	}

	err = ioutil.WriteFile(tmpFile.Name(), bs, 0644)
	if err != nil {
		e("erro ao gravar arquivo temporario: ")
		f(err.Error())
	}

	encryptFile(tmpFile.Name())

	createdFile := tmpFile.Name() + GpgExt

	v("movendo arquivo criptografado")
	backupFile := strings.TrimSuffix(target, GpgExt) + GpgExt

	v("movendo arquivo...", createdFile, " ... ", backupFile)
	util.MoveFile(backupFile, createdFile)
	i("backup criado em ", backupFile)

	v("removendo arquivo temporário")
	err = os.Remove(tmpFile.Name())
	if err != nil {
		e("erro ao gravar remover arquivo temporário: ")
		f(err.Error())
	}
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
		v("obteve arquivo gpg: \n", path)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}

func i(msg ...interface{}) {
	for _, m := range msg {
		fmt.Println(m)
	}

}

func v(msg ...interface{}) {
	sb := strings.Builder{}

	for _, m := range msg {
		sb.WriteString(fmt.Sprintf("%s ", m))
	}
	log.Println(sb.String())
}
func f(err string) {
	log.Fatalln(err)

}
func e(msg ...interface{}) {
	sb := strings.Builder{}

	for _, m := range msg {
		sb.WriteString(fmt.Sprintf("%s ", m))
	}

	log.Println(sb.String())
}

func TempFile() *os.File {
	f, err := ioutil.TempFile("", "gsh*")
	if err != nil {
		log.Fatal(f)
	}
	return f
}
