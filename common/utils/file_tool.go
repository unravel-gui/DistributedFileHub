package utils

import (
	"os"
	"strconv"
	"strings"
)

// FileExists 判断文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// EnsureFileExists 确保路径存在
func EnsureFileExists(filePath string) error {
	exists := FileExists(filePath)
	if !exists {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func SavePID(file string) error {
	// Save current PID to file
	pid := os.Getpid()
	err := os.WriteFile(file, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return err
	}
	return nil
}

func KillPID(file string) error {
	pidData, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	defer os.Remove(file)
	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		return err
	}
	existingProcess, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = existingProcess.Kill()
	return err
}

func GenNewFileName(fileName string) string {
	ss := strings.Split(fileName, ".")
	fn := "undefined"
	pf := ""
	n := len(ss)
	if n == 1 {
		fn = ss[0]
	} else {
		pf = "." + ss[n-1]
		fn = ss[n-2]
	}
	fn, num := parseFileName(fn)
	newName := fn + "[" + strconv.Itoa(num+1) + "]" + pf
	return newName
}

func parseFileName(fn string) (string, int) {
	n := len(fn) - 1
	if fn[n] != ']' {
		return fn, 0
	}
	i := n - 1
	for ; i >= 0; i-- {
		if fn[i] == '[' {
			break
		}
	}
	if i == -1 {
		return fn, 0
	}
	num, err := strconv.Atoi(fn[i+1 : n])
	if err != nil {
		return fn, 0
	}
	return fn[0:i], num
}
