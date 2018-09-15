package util

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func isRoot() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		return false
	}

	// 0 = root
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		return false
	}

	return i == 0
}

// EnableAutoStartup enable to auto start
func EnableAutoStartup() error {
	if !isRoot() {
		return errors.New("Need to run with sudo")
	}

	// copy file
	path := os.Args[0]
	pathArr := strings.Split(path, "/")
	fileName := pathArr[len(pathArr)-1]

	input, err := ioutil.ReadFile(os.Args[0])
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/usr/bin/"+fileName, input, 0644)
	if err != nil {
		return err
	}

	// startup
	startupFile := "# bot.conf\n"
	startupFile += "start on filesystem\n"
	startupFile += "exec /usr/bin/" + fileName
	err = ioutil.WriteFile("/etc/init/"+fileName+".conf", []byte(startupFile), 0644)
	if err != nil {
		return err
	}

	err = os.Symlink("/etc/init/"+fileName+".conf", "/etc/init.d/"+fileName)
	if err != nil {
		return err
	}

	return nil
}

// DisableAutoStartup disable to auto start
func DisableAutoStartup() error {
	if !isRoot() {
		return errors.New("Need to run with sudo")
	}
	return nil
}
