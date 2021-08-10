package util

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

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

func MoveFile(targetPath string, originalPath string) {

	originalFile, err := os.Open(originalPath)
	if err != nil {
		log.Fatal(err)
	}

	defer originalFile.Close()
	defer os.Remove(originalFile.Name())

	copyFile, err := os.Create(targetPath)
	if err != nil {
		log.Fatal(err)
	}
	defer copyFile.Close()

	_, err = io.Copy(copyFile, originalFile)

	if err != nil {
		log.Fatal(err)
	}
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

func CheckFatal(err error, errMsg string) {
	if err == nil {
		return
	}

	if errMsg == "" {
		log.Fatal(err.Error())
		return
	}
	log.Fatal(errMsg, " ", err.Error())

}
