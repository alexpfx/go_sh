package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const devFormat = `.web_url .author.username .commit.username .commit.created_at`

func ParseExistUntracked(workTree string, gitMessage string) []string {
	scanner := bufio.NewScanner(strings.NewReader(gitMessage))
	paths := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "\t") {
			paths = append(paths, filepath.Join(workTree, strings.TrimPrefix(line, "\t")))
		}
	}
	return paths
}

func QuoteArgs(args []string) []string {
	for i, a := range args {
		if strings.ContainsRune(a, ' ') {
			args[i] = strconv.Quote(a)
		}
	}
	return args
}

func ReadStin() string {
	r := bufio.NewReader(os.Stdin)
	text, _ := r.ReadString('\n')
	return text

}

func MoveFile(originalPath string, targetPath string) {
	err := os.MkdirAll(filepath.Dir(targetPath), 0700)
	CheckFatal(err, "")

	originalFile, err := os.Open(originalPath)
	CheckFatal(err, "")

	copyFile, err := os.Create(targetPath)
	CheckFatal(err, "")

	_, err = io.Copy(copyFile, originalFile)
	CheckFatal(err, "")

	err = copyFile.Close()
	CheckFatal(err, "")
	originalFile.Close()
	os.Remove(originalFile.Name())
}

func ExecCmd(cmdStr string, args []string) (stdout string, stderr string, err error) {
	cmd := exec.Command(cmdStr, args...)

	var sout bytes.Buffer
	var serr bytes.Buffer

	cmd.Stdout = &sout
	cmd.Stderr = &serr
	e := cmd.Run()

	return string(sout.Bytes()), string(serr.Bytes()), e
}

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		CheckFatal(err, "")
	}
	return stat.IsDir()
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		CheckFatal(err, "")
	}
	return !stat.IsDir()

}

func ToJsonStr(results interface{}) string {
	bytes, err := json.MarshalIndent(results, "", "   ")
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func CheckFatal(err error, errMsg string) {
	if err == nil {
		return
	}

	if errMsg == "" {
		log.Fatal(err.Error())
		return
	}
	log.Fatal(errMsg)
}
